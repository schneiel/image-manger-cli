use anyhow::{Context, Result};
use std::path::Path;

use super::data::ExportData;
use super::formats::{CsvExporter, JsonExporter};

pub trait Exporter {
    fn export(&self, data: &ExportData, path: &Path) -> Result<()>;
}

#[derive(Debug, Clone, Copy, PartialEq, Eq, clap::ValueEnum)]
pub enum ExportFormat {
    Csv,
    Json,
}

impl ExportFormat {
    pub fn create_exporter(self) -> Box<dyn Exporter> {
        match self {
            ExportFormat::Csv => Box::new(CsvExporter),
            ExportFormat::Json => Box::new(JsonExporter),
        }
    }

    pub fn name(self) -> &'static str {
        match self {
            ExportFormat::Csv => "CSV",
            ExportFormat::Json => "JSON",
        }
    }
}

pub fn export_data(data: &ExportData, path: &Path, format: ExportFormat) -> Result<()> {
    let exporter = format.create_exporter();

    exporter.export(data, path).with_context(|| {
        format!(
            "Failed to export data to {} format: {}",
            format.name(),
            path.display()
        )
    })
}
