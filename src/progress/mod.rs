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
            .expect("spinner template is valid by construction"),
    );
    progress.enable_steady_tick(config::DEFAULT_PROGRESS_INTERVAL);
    progress
}

pub fn create_processor_progress() -> ProgressBar {
    let progress = ProgressBar::new_spinner();
    progress.set_style(
        ProgressStyle::default_spinner()
            .template("{spinner:.green} {msg}")
            .expect("spinner template is valid by construction"),
    );
    progress.enable_steady_tick(config::SPINNER_UPDATE_INTERVAL);
    progress
}

pub fn create_copy_progress(total: u64) -> ProgressBar {
    let progress = ProgressBar::new(total);
    progress.set_style(
        ProgressStyle::default_bar()
            .template("{spinner:.green} [{bar:40.cyan/blue}] {pos}/{len} {msg}")
            .expect("bar template is valid by construction")
            .progress_chars("#>-"),
    );
    progress
}

pub fn start_progress_monitoring(
    progress_handle: imgdedup::ProgressHandle,
    initial_message: &str,
) -> std::thread::JoinHandle<()> {
    start_progress_monitoring_generic(progress_handle, initial_message)
}

pub fn start_filesort_progress_monitoring(
    progress_handle: filesort::ProgressHandle,
    initial_message: &str,
) -> std::thread::JoinHandle<()> {
    start_progress_monitoring_generic(progress_handle, initial_message)
}

fn start_progress_monitoring_generic<T: ProgressTrackerTrait>(
    progress_handle: T,
    initial_message: &str,
) -> std::thread::JoinHandle<()> {
    let spinner = create_processor_progress();
    spinner.set_message(initial_message.to_string());

    std::thread::spawn(move || {
        while !progress_handle.is_complete() {
            let info = progress_handle.get_progress();
            let current_file = info.current_file.as_deref().unwrap_or("processing...");

            spinner.set_message(format!(
                "{}: {:.1}% - {}",
                info.phase,
                info.percentage.unwrap_or(0.0),
                current_file
            ));

            std::thread::sleep(config::DEFAULT_PROGRESS_INTERVAL);
        }
        spinner.finish_with_message("Operation completed");
    })
}

pub trait ProgressTrackerTrait: Send + Sync + 'static {
    fn is_complete(&self) -> bool;
    fn get_progress(&self) -> ProgressInfoTrait;
}

#[derive(Clone)]
pub struct ProgressInfoTrait {
    pub phase: &'static str,
    pub percentage: Option<f64>,
    pub current_file: Option<String>,
}

impl ProgressTrackerTrait for imgdedup::ProgressHandle {
    fn is_complete(&self) -> bool {
        self.is_complete()
    }

    fn get_progress(&self) -> ProgressInfoTrait {
        let info = self.get_progress();
        ProgressInfoTrait {
            phase: info.phase.name(),
            percentage: info.percentage,
            current_file: info.current_file,
        }
    }
}

impl ProgressTrackerTrait for filesort::ProgressHandle {
    fn is_complete(&self) -> bool {
        self.is_complete()
    }

    fn get_progress(&self) -> ProgressInfoTrait {
        let info = self.get_progress();
        ProgressInfoTrait {
            phase: info.phase.name(),
            percentage: info.percentage,
            current_file: info.current_file,
        }
    }
}
