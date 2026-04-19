use anyhow::{Context, Result};
use console::style;
use filesort::{FileSortConfig, FileSorter, ProgressHandle};
use std::collections::HashMap;
use std::fs;

use super::OrganizeArgs;
use crate::export::{data::TargetConfig, export_data, ExportData};
use crate::output::print_organize_preview;
use crate::progress::{config, create_scanner_progress, start_filesort_progress_monitoring};
use crate::utils::{date_utils, file_ops, validation};
use crate::FILES;

pub fn handle_organize(args: &OrganizeArgs) -> Result<()> {
    validation::validate_organize_args(args)?;

    let progress = create_scanner_progress();
    progress.set_message("Initializing file sorter...");

    let mut config = FileSortConfig {
        recursive_scan: args.recursive,
        parallel_processing: true,
        max_file_size: None,
        extensions_filter: None,
    };

    if !args.extension.is_empty() {
        config.extensions_filter = Some(args.extension.clone());
    }

    let sorter = FileSorter::with_config(config);
    progress.finish_with_message("File sorter initialized");

    let progress_handle = ProgressHandle::new();
    let progress_for_monitoring = progress_handle.clone();

    let monitor_handle =
        start_filesort_progress_monitoring(progress_for_monitoring, "Organizing files...");

    let operation_start = std::time::Instant::now();
    let (organized_files, errors) = sorter
        .organize_by_date_with_progress(&args.directory, &progress_handle)
        .with_context(|| {
            format!(
                "Failed to organize files in directory: {}",
                args.directory.display()
            )
        })?;

    if let Err(e) = monitor_handle.join() {
        tracing::warn!("progress monitoring thread failed: {:?}", e);
    }

    let elapsed = operation_start.elapsed();
    println!(
        "\n{} Organization completed in {:.1}s",
        style("✓").green(),
        elapsed.as_secs_f64()
    );

    if organized_files.is_empty() && errors.is_empty() {
        println!(
            "\n{} {}",
            style("📭").yellow(),
            style("No files found in directory").bold()
        );
        return Ok(());
    }

    if let Some(export_path) = &args.export {
        let total_processed: usize = organized_files.values().map(std::vec::Vec::len).sum();
        let target_config = TargetConfig {
            base_path: args.target_path.clone(),
        };

        let export_data_obj = ExportData::organize(
            &organized_files,
            target_config,
            args.directory.clone(),
            total_processed,
        );

        export_data(&export_data_obj, export_path, args.export_format)?;

        println!(
            "\n{} {}",
            style("📄").green(),
            style("Export completed").green()
        );
        println!("   Format: {}", style(args.export_format.name()).cyan());
        println!("   Location: {}", style(export_path.display()).cyan());
    }

    let final_organized_files = if args.copy {
        if let Some(target_path) = &args.target_path {
            copy_files_to_target(&organized_files, target_path)?
        } else {
            return Err(anyhow::anyhow!(
                "--copy flag requires --target-path to be specified"
            ));
        }
    } else {
        organized_files
    };

    display_organize_results(&final_organized_files, &errors, args)?;

    Ok(())
}

fn display_organize_results(
    organized_files: &HashMap<String, Vec<std::path::PathBuf>>,
    errors: &[filesort::ProcessingError],
    args: &OrganizeArgs,
) -> Result<()> {
    println!(
        "\n{} {}",
        FILES,
        style("Organization Preview").bold().cyan()
    );
    println!("{}", style("━".repeat(50)).dim());
    print_organize_preview(organized_files, errors, args.target_path.as_ref());

    let error_strings: Vec<String> = errors
        .iter()
        .map(std::string::ToString::to_string)
        .collect();

    display_errors(&error_strings, "Processing Errors");

    if args.copy {
        if let Some(target_path) = &args.target_path {
            let target_dir = file_ops::get_target_directory(target_path);
            println!(
                "\n{} {}",
                style("📁").blue(),
                style("Files Copied Successfully").bold().blue()
            );
            println!(
                "   Target directory: {}",
                style(target_dir.display()).cyan()
            );
            println!(
                "   Total files copied: {}",
                style(
                    organized_files
                        .values()
                        .map(std::vec::Vec::len)
                        .sum::<usize>()
                        .to_string()
                )
                .green()
            );
        } else {
            return Err(anyhow::anyhow!(
                "Copy flag is set but no target path provided"
            ));
        }
    }

    Ok(())
}

fn copy_files_to_target(
    organized_files: &HashMap<String, Vec<std::path::PathBuf>>,
    target_base: &std::path::Path,
) -> Result<HashMap<String, Vec<std::path::PathBuf>>> {
    let target_dir = file_ops::get_target_directory(target_base);

    let total_files: usize = organized_files.values().map(std::vec::Vec::len).sum();
    if total_files == 0 {
        return Ok(HashMap::new());
    }

    fs::create_dir_all(&target_dir).with_context(|| {
        format!(
            "Failed to create target directory: {}",
            target_dir.display()
        )
    })?;

    let progress = crate::progress::create_copy_progress(total_files as u64);
    progress.set_message("Copying files...");

    let mut copied_files = HashMap::new();
    let mut copy_errors = Vec::new();

    for (date, files) in organized_files {
        let mut files_for_date = Vec::new();

        if let Some((year, month, day)) = date_utils::parse_date_string(date) {
            let date_dir = target_dir.join(year).join(month).join(day);
            if let Err(e) = fs::create_dir_all(&date_dir) {
                copy_errors.push(format!(
                    "Failed to create directory {}: {}",
                    date_dir.display(),
                    e
                ));
                continue;
            }

            for file in files {
                let file_name = file.file_name().map_or_else(
                    || String::from("unnamed"),
                    |n| n.to_string_lossy().to_string(),
                );
                progress.set_message(format!("Copying {file_name}"));

                let target_file = date_dir.join(
                    file.file_name()
                        .unwrap_or_else(|| std::ffi::OsStr::new("unnamed")),
                );

                let final_target_file = if target_file.exists() {
                    match file_ops::get_unique_filename(&target_file) {
                        Ok(path) => path,
                        Err(e) => {
                            copy_errors.push(format!(
                                "Failed to generate unique filename for {}: {}",
                                target_file.display(),
                                e
                            ));
                            continue;
                        }
                    }
                } else {
                    target_file
                };

                match fs::copy(file, &final_target_file) {
                    Ok(_) => {
                        files_for_date.push(final_target_file);
                    }
                    Err(e) => {
                        copy_errors.push(format!(
                            "Failed to copy {} to {}: {}",
                            file.display(),
                            final_target_file.display(),
                            e
                        ));
                    }
                }
                progress.inc(1);
            }
        }

        copied_files.insert(date.clone(), files_for_date);
    }

    progress.finish();

    if !copy_errors.is_empty() {
        display_errors(&copy_errors, "Copy Errors");
    }

    Ok(copied_files)
}

fn display_errors(errors: &[String], error_type: &str) {
    if !errors.is_empty() {
        println!("\n{} {}", style("⚠️").yellow(), style(error_type).yellow());
        println!("{}", style("━".repeat(30)).dim());

        for error in errors.iter().take(config::MAX_DISPLAY_ITEMS) {
            println!("  {}", style(format!("• {error}")).red());
        }

        if errors.len() > config::MAX_DISPLAY_ITEMS {
            println!(
                "  {} ... and {} more errors",
                style("•").red(),
                errors.len() - config::MAX_DISPLAY_ITEMS
            );
        }
    }
}
