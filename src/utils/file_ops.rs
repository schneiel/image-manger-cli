use anyhow::Result;
use std::path::{Path, PathBuf};

pub mod config {
    pub const MAX_FILENAME_ATTEMPTS: usize = 1000;
}

pub fn get_unique_filename(target_path: &Path) -> Result<PathBuf> {
    let stem = target_path
        .file_stem()
        .and_then(|s| s.to_str())
        .unwrap_or("file");
    let extension = target_path.extension().and_then(|ext| ext.to_str());

    let mut counter = 1;
    let mut new_path = target_path.to_path_buf();

    while new_path.exists() {
        let new_stem = format!("{}_{}", stem, counter);
        new_path = target_path.with_file_name(new_stem);
        if let Some(ext) = extension {
            new_path = new_path.with_extension(ext);
        }
        counter += 1;

        if counter > config::MAX_FILENAME_ATTEMPTS {
            return Err(anyhow::anyhow!(
                "Too many files with similar names exist (limit: {})",
                config::MAX_FILENAME_ATTEMPTS
            ));
        }
    }

    Ok(new_path)
}

pub fn get_target_directory(base_path: &Path) -> Result<PathBuf> {
    Ok(base_path.to_path_buf())
}

#[cfg(test)]
mod tests {
    use super::*;
    use std::fs;
    use tempfile::TempDir;

    #[test]
    fn test_get_unique_filename() {
        let temp_dir = TempDir::new().unwrap();
        let base_path = temp_dir.path().join("test.txt");

        let result = get_unique_filename(&base_path).unwrap();
        assert_eq!(result, base_path);

        fs::write(&base_path, "test").unwrap();
        let result = get_unique_filename(&base_path).unwrap();
        assert_eq!(result, temp_dir.path().join("test_1.txt"));

        fs::write(&result, "test").unwrap();
        let result2 = get_unique_filename(&base_path).unwrap();
        assert_eq!(result2, temp_dir.path().join("test_2.txt"));
    }

    #[test]
    fn test_get_target_directory() {
        let temp_dir = TempDir::new().unwrap();
        let base_path = temp_dir.path();

        let result = get_target_directory(base_path).unwrap();
        assert_eq!(result, base_path);
    }
}
