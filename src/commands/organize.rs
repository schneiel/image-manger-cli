use anyhow::{Context, Result};
use console::style;
use std::collections::HashMap;
use std::fs;
use std::path::PathBuf;

use image_manager_lib::{ImageManager, ImageManagerConfig};

use super::OrganizeArgs;
use crate::export::{data::TargetConfig, export_data, ExportData};
use crate::output::print_organize_preview;
use crate::progress::{config, create_scanner_progress, start_progress_monitoring};
use crate::utils::{date_utils, file_ops, validation};
use crate::FILES;

pub fn handle_organize(args: OrganizeArgs) -> Result<()> {
    validation::validate_organize_args(&args)?;

    let progress = create_scanner_progress();
    progress.set_message("Initializing image manager...");

    let mut config = ImageManagerConfig {
        recursive_scan: args.recursive,
        parallel_processing: true,
        ..Default::default()
    };

    if let Some(ref format_filter) = args.format {
        config.supported_formats = vec![format_filter.clone().into()];
    }

    let manager = ImageManager::with_config(config.clone());
    progress.finish_with_message("Image manager initialized");

    let progress_handle = image_manager_lib::ProgressHandle::new();
    let progress_for_monitoring = progress_handle.clone();

    let monitor_handle = start_progress_monitoring(progress_for_monitoring, "Organizing images...");

    let operation_start = std::time::Instant::now();
    let (organized_images, errors) = manager
        .organize_by_date_with_progress(&args.directory, &progress_handle)
        .with_context(|| {
            format!(
                "Failed to organize images in directory: {}",
                args.directory.display()
            )
        })?;

    let _ = monitor_handle.join();

    let elapsed = operation_start.elapsed();
    println!(
        "\n{} Organization completed in {:.1}s",
        style("✓").green(),
        elapsed.as_secs_f64()
    );

    if organized_images.is_empty() && errors.is_empty() {
        println!(
            "\n{} {}",
            style("📭").yellow(),
            style("No supported images found in directory").bold()
        );
        return Ok(());
    }

    if let Some(export_path) = &args.export {
        let total_processed: usize = organized_images.values().map(|v| v.len()).sum();
        let target_config = TargetConfig {
            base_path: args.target_path.clone(),
        };

        let export_data_obj = ExportData::organize(
            organized_images.clone(),
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

    let final_organized_images = if args.copy {
        if let Some(target_path) = &args.target_path {
            copy_files_to_target(&organized_images, target_path)?
        } else {
            return Err(anyhow::anyhow!(
                "--copy flag requires --target-path to be specified"
            ));
        }
    } else {
        organized_images
    };

    display_organize_results(&final_organized_images, &errors, &args)?;

    Ok(())
}

fn display_organize_results(
    organized_images: &HashMap<String, Vec<PathBuf>>,
    errors: &[image_manager_lib::ProcessingError],
    args: &OrganizeArgs,
) -> Result<()> {
    println!(
        "\n{} {}",
        FILES,
        style("Organization Preview").bold().cyan()
    );
    println!("{}", style("━".repeat(50)).dim());
    print_organize_preview(organized_images, errors, args.target_path.as_ref());

    let error_strings: Vec<String> = errors.iter().map(|e| e.to_string()).collect();

    display_errors(&error_strings, "Processing Errors");

    if args.copy {
        if let Some(target_path) = &args.target_path {
            let target_dir = file_ops::get_target_directory(target_path)?;
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
                    organized_images
                        .values()
                        .map(|v| v.len())
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
    organized_images: &HashMap<String, Vec<PathBuf>>,
    target_base: &std::path::Path,
) -> Result<HashMap<String, Vec<PathBuf>>> {
    let target_dir = file_ops::get_target_directory(target_base)?;

    let total_files: usize = organized_images.values().map(|v| v.len()).sum();
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

    for (date, files) in organized_images {
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
                progress.set_message(format!(
                    "Copying {}",
                    file.file_name().unwrap_or_default().to_string_lossy()
                ));

                let target_file = date_dir.join(file.file_name().unwrap_or_default());

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
            println!("  {}", style(format!("• {}", error)).red());
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
