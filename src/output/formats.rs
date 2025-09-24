use console::style;
use std::collections::HashMap;
use std::path::PathBuf;

use image_manager_lib::ProcessingError;

pub fn print_organize_preview(
    organized_images: &HashMap<String, Vec<PathBuf>>,
    errors: &[ProcessingError],
    target_path: Option<&PathBuf>,
) {
    if organized_images.is_empty() && errors.is_empty() {
        println!(
            "\n{} {}",
            style("📭").yellow(),
            style("No supported images found in directory").bold()
        );
        return;
    }

    println!(
        "\n{} {}",
        style("📁").cyan(),
        style("Organization Preview").bold().cyan()
    );
    println!("{}", style("━".repeat(50)).dim());

    for (date, files) in organized_images {
        println!("\n{} {}", style("📅").blue(), style(date).bold());

        if let Some(target_path) = target_path {
            let target_dir_name = target_path
                .file_name()
                .and_then(|name| name.to_str())
                .unwrap_or("untitled");
            println!(
                "   Target: {}/{}/{}",
                style(target_dir_name).green(),
                style(date).cyan(),
                style(files.len()).yellow()
            );
        } else {
            println!("   Files: {}", style(files.len()).yellow());
        }

        for (i, file) in files.iter().enumerate() {
            println!(
                "   {}. {}",
                style(i + 1).dim(),
                style(file.file_name().unwrap_or_default().to_string_lossy()).cyan()
            );
        }
    }

    print_errors(errors);
}

pub fn print_duplicates_preview(
    duplicate_groups: &image_manager_lib::duplicates::DuplicateGroups,
    errors: &[ProcessingError],
    similarity_threshold: &image_manager_lib::SimilarityThreshold,
) {
    if duplicate_groups.is_empty() && errors.is_empty() {
        println!(
            "\n{} {}",
            style("📭").yellow(),
            style("No duplicate images found").bold()
        );
        return;
    }

    println!(
        "\n{} {}",
        style("🔄").cyan(),
        style("Duplicate Detection Preview").bold().cyan()
    );
    println!("{}", style("━".repeat(50)).dim());
    println!(
        "Similarity threshold: {}",
        style(format!("{:.2}%", similarity_threshold.value() * 100.0)).green()
    );

    for (group_index, group) in duplicate_groups.iter().enumerate() {
        if group.len() > 1 {
            println!(
                "\n{} {}",
                style("Group").blue(),
                style(group_index + 1).bold()
            );
            println!("   Files: {}", style(group.len()).yellow());

            for (file_index, file) in group.iter().enumerate() {
                let size_str = if let Ok(metadata) = std::fs::metadata(file) {
                    format!(" ({})", style(format_bytes(metadata.len())).dim())
                } else {
                    String::new()
                };

                println!(
                    "   {}. {}{}",
                    style(file_index + 1).dim(),
                    style(file.display()).cyan(),
                    size_str
                );
            }
        }
    }

    print_errors(errors);
}

pub fn print_errors(errors: &[ProcessingError]) {
    if !errors.is_empty() {
        println!(
            "\n{} {}",
            style("⚠️").yellow(),
            style("Processing Errors").yellow()
        );
        println!("{}", style("━".repeat(30)).dim());

        for error in errors.iter().take(10) {
            println!("  {}", style(format!("• {}", error)).red());
        }

        if errors.len() > 10 {
            println!(
                "  {} ... and {} more errors",
                style("•").red(),
                errors.len() - 10
            );
        }
    }
}

fn format_bytes(bytes: u64) -> String {
    const UNITS: &[&str] = &["B", "KB", "MB", "GB"];
    let mut size = bytes as f64;
    let mut unit_index = 0;

    while size >= 1024.0 && unit_index < UNITS.len() - 1 {
        size /= 1024.0;
        unit_index += 1;
    }

    if unit_index == 0 {
        format!("{} {}", bytes, UNITS[unit_index])
    } else {
        format!("{:.1} {}", size, UNITS[unit_index])
    }
}
