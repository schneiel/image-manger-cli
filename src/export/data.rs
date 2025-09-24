use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use std::path::PathBuf;

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ExportData {
    pub metadata: ExportMetadata,
    pub data: ExportDataType,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ExportMetadata {
    pub timestamp: DateTime<Utc>,
    pub command: String,
    pub version: String,
    pub source_directory: PathBuf,
    pub total_processed: usize,
    pub command_metadata: HashMap<String, serde_json::Value>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(tag = "type")]
pub enum ExportDataType {
    Organize {
        file_records: Vec<OrganizeFileRecord>,
        target_config: TargetConfig,
    },
    Duplicates {
        file_records: Vec<DuplicateFileRecord>,
        similarity_threshold: f32,
    },
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct TargetConfig {
    pub base_path: Option<PathBuf>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct OrganizeFileRecord {
    pub original_path: PathBuf,
    pub target_path: PathBuf,
    pub date_directory: String,
    pub file_name: String,
    pub file_size_bytes: u64,
    pub file_extension: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct DuplicateFileRecord {
    pub file_path: PathBuf,
    pub group_id: String,
    pub position_in_group: usize,
    pub group_size: usize,
    pub similarity: f32,
    pub file_size_bytes: u64,
    pub file_extension: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct DuplicateGroup {
    pub group_id: String,
    pub files: Vec<PathBuf>,
    pub similarity: f32,
}

impl ExportData {
    pub fn organize(
        organized_files: HashMap<String, Vec<PathBuf>>,
        target_config: TargetConfig,
        source_directory: PathBuf,
        total_processed: usize,
    ) -> Self {
        let mut file_records = Vec::new();

        for (date, files) in &organized_files {
            for file_path in files {
                let file_name = file_path
                    .file_name()
                    .and_then(|n| n.to_str())
                    .unwrap_or("unknown")
                    .to_string();

                let file_extension = file_path
                    .extension()
                    .and_then(|e| e.to_str())
                    .unwrap_or("")
                    .to_string();

                let file_size = std::fs::metadata(file_path).map(|m| m.len()).unwrap_or(0);

                let target_path = if let Some(base_path) = &target_config.base_path {
                    let dir_name = base_path
                        .file_name()
                        .and_then(|name| name.to_str())
                        .unwrap_or("untitled");
                    PathBuf::from(format!("{}/{}/{}", dir_name, date, file_name))
                } else {
                    let source_dir_name = source_directory
                        .file_name()
                        .and_then(|name| name.to_str())
                        .unwrap_or("untitled");
                    PathBuf::from(format!("{}/{}/{}", source_dir_name, date, file_name))
                };

                file_records.push(OrganizeFileRecord {
                    original_path: file_path.clone(),
                    target_path,
                    date_directory: date.clone(),
                    file_name,
                    file_extension,
                    file_size_bytes: file_size,
                });
            }
        }

        let mut command_metadata = HashMap::new();
        if let Some(ref base_path) = target_config.base_path {
            command_metadata.insert(
                "target_path".to_string(),
                serde_json::json!(base_path.to_string_lossy()),
            );
        }

        Self {
            metadata: ExportMetadata {
                timestamp: Utc::now(),
                command: "organize".to_string(),
                version: env!("CARGO_PKG_VERSION").to_string(),
                source_directory,
                total_processed,
                command_metadata,
            },
            data: ExportDataType::Organize {
                file_records,
                target_config,
            },
        }
    }

    pub fn duplicates(
        duplicate_groups: Vec<DuplicateGroup>,
        similarity_threshold: f32,
        source_directory: PathBuf,
        total_processed: usize,
    ) -> Self {
        let mut file_records = Vec::new();

        for group in &duplicate_groups {
            for (position, file_path) in group.files.iter().enumerate() {
                let file_extension = file_path
                    .extension()
                    .and_then(|e| e.to_str())
                    .unwrap_or("")
                    .to_string();

                let file_size = std::fs::metadata(file_path).map(|m| m.len()).unwrap_or(0);

                file_records.push(DuplicateFileRecord {
                    file_path: file_path.clone(),
                    group_id: group.group_id.clone(),
                    position_in_group: position + 1,
                    group_size: group.files.len(),
                    similarity: group.similarity,
                    file_size_bytes: file_size,
                    file_extension,
                });
            }
        }

        let mut command_metadata = HashMap::new();
        command_metadata.insert(
            "similarity_threshold".to_string(),
            serde_json::json!(similarity_threshold),
        );
        command_metadata.insert(
            "duplicate_groups_count".to_string(),
            serde_json::json!(duplicate_groups.len()),
        );

        Self {
            metadata: ExportMetadata {
                timestamp: Utc::now(),
                command: "duplicates".to_string(),
                version: env!("CARGO_PKG_VERSION").to_string(),
                source_directory,
                total_processed,
                command_metadata,
            },
            data: ExportDataType::Duplicates {
                file_records,
                similarity_threshold,
            },
        }
    }
}
