use image_manager_lib::ProgressHandle;
use indicatif::{ProgressBar, ProgressStyle};

pub mod config {
    use std::time::Duration;

    pub const DEFAULT_PROGRESS_INTERVAL: Duration = Duration::from_millis(100);
    pub const SPINNER_UPDATE_INTERVAL: Duration = Duration::from_millis(120);
    pub const MAX_DISPLAY_ITEMS: usize = 10;
}

pub fn create_scanner_progress() -> ProgressBar {
    let progress = ProgressBar::new_spinner();
    progress.set_style(
        ProgressStyle::default_spinner()
            .template("{spinner:.green} {msg}")
            .unwrap(),
    );
    progress.enable_steady_tick(config::DEFAULT_PROGRESS_INTERVAL);
    progress
}

pub fn create_processor_progress() -> ProgressBar {
    let progress = ProgressBar::new_spinner();
    progress.set_style(
        ProgressStyle::default_spinner()
            .template("{spinner:.green} {msg:.cyan}")
            .unwrap(),
    );
    progress.enable_steady_tick(config::SPINNER_UPDATE_INTERVAL);
    progress
}

pub fn create_copy_progress(total: u64) -> ProgressBar {
    let progress = ProgressBar::new(total);
    progress.set_style(
        ProgressStyle::default_bar()
            .template("{spinner:.green} [{bar:40.cyan/blue}] {pos}/{len} {msg}")
            .unwrap()
            .progress_chars("#>-"),
    );
    progress
}

pub fn start_progress_monitoring(
    progress_handle: ProgressHandle,
    initial_message: &str,
) -> std::thread::JoinHandle<()> {
    let spinner = create_processor_progress();
    spinner.set_message(initial_message.to_string());
    let spinner_clone = spinner.clone();

    std::thread::spawn(move || {
        while !progress_handle.is_complete() {
            let info = progress_handle.get_progress();
            let current_file = info.current_file.as_deref().unwrap_or("processing...");

            spinner_clone.set_message(format!(
                "{}: {:.1}% - {}",
                info.phase.name(),
                info.percentage.unwrap_or(0.0),
                current_file
            ));

            std::thread::sleep(config::DEFAULT_PROGRESS_INTERVAL);
        }
        spinner_clone.finish_with_message("Operation completed");
    })
}
