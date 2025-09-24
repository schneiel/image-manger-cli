use anyhow::{Context, Result};
use std::fs::File;
use std::io::Write;
use std::path::Path;

use super::data::{ExportData, ExportDataType};
use super::trait_impl::Exporter;

pub struct CsvExporter;

impl Exporter for CsvExporter {
    fn export(&self, data: &ExportData, path: &Path) -> Result<()> {
        let mut file = File::create(path)
            .with_context(|| format!("Failed to create CSV file: {}", path.display()))?;

        match &data.data {
            ExportDataType::Organize {
                file_records,
                target_config,
            } => {
                self.export_organize_csv(&mut file, file_records, target_config)?;
            }
            ExportDataType::Duplicates {
                file_records,
                similarity_threshold,
            } => {
                self.export_duplicates_csv(&mut file, file_records, *similarity_threshold)?;
            }
        }

        Ok(())
    }
}

impl CsvExporter {
    fn export_organize_csv(
        &self,
        file: &mut File,
        file_records: &[crate::export::data::OrganizeFileRecord],
        _target_config: &crate::export::data::TargetConfig,
    ) -> Result<()> {
        writeln!(
            file,
            "Original Path,Target Path,Date Directory,File Name,File Size (bytes),File Extension"
        )?;

        for record in file_records {
            writeln!(
                file,
                "\"{}\",\"{}\",\"{}\",\"{}\",{},\"{}\"",
                record.original_path.display(),
                record.target_path.display(),
                record.date_directory,
                record.file_name,
                record.file_size_bytes,
                record.file_extension
            )?;
        }

        Ok(())
    }

    fn export_duplicates_csv(
        &self,
        file: &mut File,
        file_records: &[crate::export::data::DuplicateFileRecord],
        _similarity_threshold: f32,
    ) -> Result<()> {
        writeln!(file, "Group ID,File Path,Position in Group,Group Size,Similarity,File Size (bytes),File Extension")?;

        for record in file_records {
            writeln!(
                file,
                "\"{}\",\"{}\",{},{},{:.4},{},\"{}\"",
                record.group_id,
                record.file_path.display(),
                record.position_in_group,
                record.group_size,
                record.similarity,
                record.file_size_bytes,
                record.file_extension
            )?;
        }

        Ok(())
    }
}

pub struct JsonExporter;

impl Exporter for JsonExporter {
    fn export(&self, data: &ExportData, path: &Path) -> Result<()> {
        let json_string = serde_json::to_string_pretty(data)
            .with_context(|| "Failed to serialize data to JSON")?;

        std::fs::write(path, json_string)
            .with_context(|| format!("Failed to write JSON file: {}", path.display()))?;

        Ok(())
    }
}
