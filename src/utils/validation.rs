use anyhow::Result;
use std::path::Path;

pub fn validate_directory(path: &Path, description: &str) -> Result<()> {
    if !path.exists() {
        return Err(anyhow::anyhow!(
            "{} does not exist: {}",
            description,
            path.display()
        ));
    }

    if !path.is_dir() {
        return Err(anyhow::anyhow!(
            "{} is not a directory: {}",
            description,
            path.display()
        ));
    }

    Ok(())
}

pub fn validate_different_directories(source: &Path, target: &Path) -> Result<()> {
    let source_canonical = source
        .canonicalize()
        .unwrap_or_else(|_| source.to_path_buf());
    let target_canonical = target
        .canonicalize()
        .unwrap_or_else(|_| target.to_path_buf());

    if source_canonical == target_canonical {
        return Err(anyhow::anyhow!(
            "Source and target directories cannot be same"
        ));
    }

    Ok(())
}

pub fn validate_organize_args(args: &crate::commands::OrganizeArgs) -> Result<()> {
    validate_directory(&args.directory, "Source directory")?;

    if let Some(target_path) = &args.target_path {
        if target_path.exists() && target_path.is_dir() {
            validate_different_directories(&args.directory, target_path)?;
        }
    }

    if args.copy && args.target_path.is_none() {
        return Err(anyhow::anyhow!(
            "--copy flag requires --target-path to be specified"
        ));
    }

    Ok(())
}

pub fn validate_similarity_threshold(threshold: f32) -> Result<()> {
    if !(0.0..=1.0).contains(&threshold) {
        return Err(anyhow::anyhow!(
            "Similarity threshold must be between 0.0 and 1.0, got: {}",
            threshold
        ));
    }
    Ok(())
}

pub fn validate_duplicates_args(args: &crate::commands::DuplicatesArgs) -> Result<()> {
    validate_directory(&args.directory, "Source directory")?;

    if let Some(threshold) = args.threshold {
        validate_similarity_threshold(threshold)?;
    }

    Ok(())
}

#[cfg(test)]
mod tests {
    use super::*;
    use std::fs;
    use tempfile::TempDir;

    #[test]
    fn test_validate_directory_valid() {
        let temp_dir = TempDir::new().unwrap();
        assert!(validate_directory(temp_dir.path(), "Test directory").is_ok());
    }

    #[test]
    fn test_validate_directory_nonexistent() {
        let path = Path::new("/nonexistent/directory");
        assert!(validate_directory(path, "Test directory").is_err());
    }

    #[test]
    fn test_validate_directory_is_file() {
        let temp_dir = TempDir::new().unwrap();
        let file_path = temp_dir.path().join("test.txt");
        fs::write(&file_path, "test").unwrap();

        assert!(validate_directory(&file_path, "Test directory").is_err());
    }

    #[test]
    fn test_validate_similarity_threshold() {
        assert!(validate_similarity_threshold(0.5).is_ok());
        assert!(validate_similarity_threshold(0.0).is_ok());
        assert!(validate_similarity_threshold(1.0).is_ok());
        assert!(validate_similarity_threshold(-0.1).is_err());
        assert!(validate_similarity_threshold(1.1).is_err());
    }

    #[test]
    fn test_validate_different_directories() {
        let temp_dir = TempDir::new().unwrap();
        let dir1 = temp_dir.path().join("dir1");
        let dir2 = temp_dir.path().join("dir2");
        fs::create_dir(&dir1).unwrap();
        fs::create_dir(&dir2).unwrap();

        assert!(validate_different_directories(&dir1, &dir2).is_ok());
        assert!(validate_different_directories(&dir1, &dir1).is_err());
    }
}
