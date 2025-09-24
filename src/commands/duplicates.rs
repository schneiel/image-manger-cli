use anyhow::{Context, Result};
use console::style;
use image_manager_lib::{ImageManager, ImageManagerConfig};

use super::DuplicatesArgs;
use crate::export::{data::DuplicateGroup, export_data, ExportData};
use crate::output::print_duplicates_preview;
use crate::progress::{config, create_scanner_progress, start_progress_monitoring};
use crate::utils::validation;
use crate::DUPLICATE;

pub fn handle_duplicates(args: DuplicatesArgs) -> Result<()> {
    validation::validate_duplicates_args(&args)?;

    let progress = create_scanner_progress();
    progress.set_message("Initializing image manager...");

    let mut config = ImageManagerConfig {
        recursive_scan: args.recursive,
        similarity_threshold: args
            .get_similarity_threshold()
            .map_err(|e| anyhow::anyhow!("Invalid similarity threshold: {}", e))?,
        parallel_processing: true,
        ..Default::default()
    };

    config.duplicate_mode = args.mode.into();

    let manager = ImageManager::with_config(config.clone());
    progress.finish_with_message("Image manager initialized");

    let progress_handle = image_manager_lib::ProgressHandle::new();
    let progress_for_monitoring = progress_handle.clone();

    let monitor_handle =
        start_progress_monitoring(progress_for_monitoring, "Scanning for duplicate images...");

    let operation_start = std::time::Instant::now();
    let (duplicate_groups, errors) = manager
        .find_duplicates_with_progress(&args.directory, &progress_handle)
        .with_context(|| "Failed to find duplicates")?;

    let _ = monitor_handle.join();

    let elapsed = operation_start.elapsed();
    println!(
        "\n{} Duplicate detection completed in {:.1}s",
        style("✓").green(),
        elapsed.as_secs_f64()
    );

    display_duplicates_results(
        &duplicate_groups,
        &errors,
        &args,
        &config.similarity_threshold,
    )?;

    Ok(())
}

fn display_duplicates_results(
    duplicate_groups: &image_manager_lib::duplicates::DuplicateGroups,
    errors: &[image_manager_lib::ProcessingError],
    args: &DuplicatesArgs,
    similarity_threshold: &image_manager_lib::SimilarityThreshold,
) -> Result<()> {
    println!(
        "\n{} {}",
        DUPLICATE,
        style("Duplicate Detection Preview").bold().cyan()
    );
    println!("{}", style("━".repeat(50)).dim());

    print_duplicates_preview(duplicate_groups, errors, similarity_threshold);

    if let Some(export_path) = &args.export {
        let total_processed: usize = duplicate_groups.iter().map(|group| group.len()).sum();

        let export_duplicate_groups: Vec<DuplicateGroup> = duplicate_groups
            .iter()
            .enumerate()
            .map(|(index, group)| DuplicateGroup {
                group_id: format!("group_{}", index + 1),
                files: group.clone(),
                similarity: similarity_threshold.value(),
            })
            .collect();

        let export_data_obj = ExportData::duplicates(
            export_duplicate_groups,
            similarity_threshold.value(),
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

    let error_strings: Vec<String> = errors.iter().map(|e| e.to_string()).collect();

    display_errors(&error_strings, "Processing Errors");

    Ok(())
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
