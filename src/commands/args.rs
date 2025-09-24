use crate::export::ExportFormat;
use clap::{Args, ValueEnum};
use image_manager_lib::SimilarityThreshold;

#[derive(Args)]
pub struct OrganizeArgs {
    #[arg(help = "Directory to scan for images (default: current directory)")]
    pub directory: std::path::PathBuf,

    #[arg(
        short = 'r',
        long,
        help = "Scan directories recursively (default: false)"
    )]
    pub recursive: bool,

    #[arg(long, value_enum, help = "Filter by specific image format")]
    pub format: Option<ImageFormatFilter>,

    #[arg(long, help = "Export results to file")]
    pub export: Option<std::path::PathBuf>,

    #[arg(
        long,
        value_enum,
        default_value = "csv",
        help = "Export format (csv or json)"
    )]
    pub export_format: ExportFormat,

    #[arg(
        long,
        help = "Target directory for organized files (required with --copy)"
    )]
    pub target_path: Option<std::path::PathBuf>,

    #[arg(long, help = "Copy files to target directory (default: preview only)")]
    pub copy: bool,
}

impl Default for OrganizeArgs {
    fn default() -> Self {
        Self {
            directory: std::path::PathBuf::from("."),
            recursive: false,
            format: None,
            export: None,
            export_format: ExportFormat::Csv,
            target_path: None,
            copy: false,
        }
    }
}

#[derive(Args)]
pub struct DuplicatesArgs {
    #[arg(help = "Directory to scan for duplicate images (default: current directory)")]
    pub directory: std::path::PathBuf,

    #[arg(
        short = 'r',
        long,
        help = "Scan directories recursively (default: false)"
    )]
    pub recursive: bool,

    #[arg(
        long,
        help = "Similarity threshold for duplicate detection (0.0-1.0, e.g., 0.85)"
    )]
    pub threshold: Option<f32>,

    #[arg(
        long,
        value_enum,
        help = "Preset similarity threshold level (overrides --threshold, default: medium)"
    )]
    pub sensitivity: Option<ThresholdLevel>,

    #[arg(long, help = "Export results to file")]
    pub export: Option<std::path::PathBuf>,

    #[arg(
        long,
        value_enum,
        default_value = "json",
        help = "Export format (csv or json)"
    )]
    pub export_format: ExportFormat,

    #[arg(
        long,
        value_enum,
        default_value = "size_filtered",
        help = "Duplicate detection mode (default: size_filtered)"
    )]
    pub mode: DuplicateScanMode,
}

impl Default for DuplicatesArgs {
    fn default() -> Self {
        Self {
            directory: std::path::PathBuf::from("."),
            recursive: false,
            threshold: None,
            sensitivity: None,
            export: None,
            export_format: ExportFormat::Json,
            mode: DuplicateScanMode::SizeFiltered,
        }
    }
}

impl DuplicatesArgs {
    pub fn get_similarity_threshold(&self) -> Result<SimilarityThreshold, String> {
        if let Some(preset_level) = self.sensitivity {
            Ok(preset_level.into())
        } else if let Some(custom_threshold) = self.threshold {
            SimilarityThreshold::new(custom_threshold)
                .map_err(|e| format!("Invalid threshold: {}", e))
        } else {
            Ok(SimilarityThreshold::medium())
        }
    }
}

#[derive(ValueEnum, Clone)]
pub enum ImageFormatFilter {
    Jpeg,
    Png,
    Gif,
    Tiff,
    WebP,
    Bmp,
    Ico,
}

impl From<ImageFormatFilter> for image_manager_lib::config::ImageFormat {
    fn from(filter: ImageFormatFilter) -> Self {
        match filter {
            ImageFormatFilter::Jpeg => image_manager_lib::config::ImageFormat::Jpeg,
            ImageFormatFilter::Png => image_manager_lib::config::ImageFormat::Png,
            ImageFormatFilter::Gif => image_manager_lib::config::ImageFormat::Gif,
            ImageFormatFilter::Tiff => image_manager_lib::config::ImageFormat::Tiff,
            ImageFormatFilter::WebP => image_manager_lib::config::ImageFormat::WebP,
            ImageFormatFilter::Bmp => image_manager_lib::config::ImageFormat::Bmp,
            ImageFormatFilter::Ico => image_manager_lib::config::ImageFormat::Ico,
        }
    }
}

#[derive(ValueEnum, Clone, Copy, Debug)]
pub enum ThresholdLevel {
    #[value(name = "low")]
    Low,
    #[value(name = "medium")]
    Medium,
    #[value(name = "high")]
    High,
}

impl From<ThresholdLevel> for SimilarityThreshold {
    fn from(level: ThresholdLevel) -> Self {
        match level {
            ThresholdLevel::Low => SimilarityThreshold::low(),
            ThresholdLevel::Medium => SimilarityThreshold::medium(),
            ThresholdLevel::High => SimilarityThreshold::high(),
        }
    }
}

#[derive(ValueEnum, Clone, Copy, Debug)]
pub enum DuplicateScanMode {
    #[value(name = "size_filtered")]
    SizeFiltered,
    #[value(name = "complete")]
    Complete,
}

impl From<DuplicateScanMode> for image_manager_lib::config::DuplicateMode {
    fn from(mode: DuplicateScanMode) -> Self {
        match mode {
            DuplicateScanMode::SizeFiltered => {
                image_manager_lib::config::DuplicateMode::SizeFiltered
            }
            DuplicateScanMode::Complete => image_manager_lib::config::DuplicateMode::Complete,
        }
    }
}
