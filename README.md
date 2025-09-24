# 🖼️ ImageManager

**ImageManager** is a cross-platform command-line tool for organizing and cleaning up image libraries. It automatically sorts images into a structured folder hierarchy and detects (and optionally removes) duplicate files safely and efficiently.

Built in Go — lightweight, fast, and ready for automation.

## 📋 Table of Contents

- [🔑 Features](#-features)
- [⚙️ Installation](#️-installation)
- [🚀 Usage](#-usage)
  - [📁 Sort Images](#-sort-images)
  - [🧹 Find Duplicates](#-find-duplicates)
  - [🌐 Global Flags](#-global-flags)
  - [📜 Subcommands](#-subcommands)
- [💡 Common Usage Patterns](#-common-usage-patterns)
- [🛠 Configuration](#-configuration-configyaml)
- [📝 Logging](#-logging)
- [🌐 Internationalization](#-internationalization)
- [🤝 Contributing](#-contributing)

---

## 🔑 Features

- **Image Sorting**
  Automatically organizes photos into a `YYYY/MM/DD` folder structure based on capture date.

- **Duplicate Detection**
  Identifies duplicate images using content-based file hashing.

- **Flexible Date Extraction**
  Determines an image's date using a prioritized fallback strategy:
  1. EXIF metadata (default)
  2. File modification time
  2. File creation time

- **Customizable Duplicate Handling**
  - **Action Mode**: *Dry Run* (simulation) or *Move to Trash*
  - **Retention Strategy**: Keep the *oldest* file or the one with the *shortest file path*

- **Safe by Design**
  Built-in safety mechanisms like *dryRun* mode and *trash-first* deletion help prevent accidental data loss.

- **Multilingual Support**
  Currently supports English (`en`) and German (`de`).

- **Comprehensive Logging**
  All operations are logged for transparency and troubleshooting.

---

## ⚙️ Installation

### Build from Source

**Prerequisites**: Go 1.24+

```bash
git clone git@github.com:schneiel/ImageManagerGo.git
cd ImageManagerGo

# Build application
go build -o imagemanager main.go
```

The `imagemanager` binary will be created in the root directory.

---

## 🚀 Usage

ImageManager uses subcommands to perform operations. The two main ones are `sort` and `dedup`.

### 📁 Sort Images

Organize images from a source directory into `YYYY/MM/DD` folders under a destination path:

```bash
# Basic usage (dry run mode)
./imagemanager sort --source "/path/to/images" --destination "/path/to/archive"

# Actually copy files
./imagemanager sort --source "/path/to/images" --destination "/path/to/archive" --actionStrategy "copy"
```

### 🧹 Find Duplicates

Scan a directory for duplicate images and process them based on the configured strategy:

```bash
# Basic usage (dry run mode)
./imagemanager dedup --source "/path/to/images"

# Actually move duplicates to trash
./imagemanager dedup --source "/path/to/images" --actionStrategy "moveToTrash"

# Advanced usage with custom settings
./imagemanager dedup --source "/path/to/images" --actionStrategy "moveToTrash" --keepStrategy "keepShortestPath" --workers 4 --threshold 5
```

---

### 🌐 Global Flags

These flags can be used with any subcommand:

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--config` | `-c` | Path to a configuration file | `config.yaml` |
| `--language` | `-l` | Language for localization (available: `en`, `de`) | `en` |

---

### 📜 Subcommands

#### `sort`

Organizes images into a date-based folder structure.

**Flags:**

| Flag | Short | Description | Options/Default |
|------|-------|-------------|-----------------|
| `--source` | `-s` | Source directory containing images to sort | (overrides config) |
| `--destination` | `-d` | Destination directory for sorted images | (overrides config) |
| `--actionStrategy` | | Action Strategy | `dryRun`, `copy` (default: `dryRun`) |

**Example:**
```bash
./imagemanager sort --source "/path/to/images" --destination "/path/to/archive" --actionStrategy "copy"
```

#### `dedup`

Detects and handles duplicate images.

**Flags:**

| Flag | Short | Description | Options/Default |
|------|-------|-------------|-----------------|
| `--source` | `-s` | Source directory to scan for duplicates | (overrides config) |
| `--actionStrategy` | | Action strategy | `dryRun`, `moveToTrash` (default: `dryRun`) |
| `--keepStrategy` | `-k` | Strategy for which file to keep | `keepOldest`, `keepShortestPath` (default: `keepOldest`) |
| `--trashPath` | `-t` | Path to move duplicates to | (default: `.trash`) |
| `--workers` | `-w` | Number of worker threads for hashing | (default: 8) |
| `--threshold` | | Similarity threshold for images | (default: 1) |

**Example:**
```bash
./imagemanager dedup --source "/path/to/images" --actionStrategy "moveToTrash" --keepStrategy "keepOldest" --workers 4 --threshold 5
```

---

## 💡 Common Usage Patterns

### Quick Start Commands

```bash
# Test sorting (dry run)
./imagemanager sort -s "/path/to/images" -d "/path/to/archive"

# Actually sort images
./imagemanager sort -s "/path/to/images" -d "/path/to/archive" --actionStrategy "copy"

# Test duplicate detection (dry run)
./imagemanager dedup -s "/path/to/images"

# Remove duplicates (move to trash)
./imagemanager dedup -s "/path/to/images" --actionStrategy "moveToTrash"

# High-performance duplicate detection
./imagemanager dedup -s "/path/to/images" -w 16 --threshold 0

# Use German language
./imagemanager -l de sort -s "/path/to/images" -d "/path/to/archive"
```

### Command Combinations

```bash
# Sort first, then deduplicate
./imagemanager sort -s "/unsorted" -d "/sorted" --actionStrategy "copy"
./imagemanager dedup -s "/sorted" --actionStrategy "moveToTrash"
```

---

## 🛠 Configuration (`config.yaml`)

You can define advanced options in a `config.yaml` file.

### Example:

```yaml
allowedImageExtensions:
  - ".jpg"
  - ".jpeg"
  - ".png"

# Configuration for file paths and names
files:
  applicationLog: "application.log"
  dedupDryRunLog: "dedup_dry_run_log.csv"
  sortDryRunLog: "sort_dry_run_log.csv"

# Configuration for the deduplicator command
deduplicator:
  source: "/path/to/images"
  # Defines the action to take when duplicates are found.
  # Valid options: "moveToTrash", "dryRun"
  actionStrategy: "dryRun"
  # Defines which file to keep when duplicates are found.
  # Valid options: "keepOldest", "keepShortestPath"
  keepStrategy: "keepOldest"
  # Path for trash directory (default: ".trash")
  trashPath: ".trash"
  # Number of worker threads for processing (default: 8)
  workers: 8
  # Similarity threshold for duplicate detection (default: 1)
  threshold: 1

# Configuration for date extraction during sorting
sorter:
  source: "/path/to/source"
  destination: "/path/to/destination"
  # Valid options: "copy", "dryRun"
  actionStrategy: "dryRun"
  # Defines the order of strategies to use for extracting an image's date.
  # The first strategy that succeeds is used.
  # Valid options: "exif", "creationTime", "modTime"
  date:
    strategyOrder:
      - "exif"
      - "modTime"
      - "creationTime"
    exifStrategies:
      - fieldName: "DateTimeOriginal"
        layout: "2006:01:02 15:04:05"
      - fieldName: "DateTimeDigitized"
        layout: "2006:01:02 15:04:05"
      - fieldName: "SubSecDateTimeOriginal"
        layout: "2006:01:02 15:04:05.00"
      - fieldName: "DateTime"
        layout: "2006:01:02 15:04:05-07:00"
```

---

## 📝 Logging

Log and output files are created during execution. These filenames are configurable in the `files` section of your config:

- **Application Log** (default: `application.log`): General logs and errors
- **Dedup Dry Run Log** (default: `dedup_dry_run_log.csv`): List of duplicates found (in dryRun mode), including paths and hashes
- **Sort Dry Run Log** (default: `sort_dry_run_log.csv`): List of old paths and expected new paths after sortation

### Configurable File Paths

You can customize all output file names in your `config.yaml`:

```yaml
files:
  applicationLog: "my_custom_app.log"
  dedupDryRunLog: "duplicates_report.csv"
  sortDryRunLog: "sort_operations.csv"
```

---

## 🌐 Internationalization

Language selection via `--language` flag.

**Supported Languages:**
- English (`en`)
- German (`de`)

**Example:**
```bash
./imagemanager --language="de" dedup --source "/images"
```

---

## 📌 Planned Features

- Configuration validation
- GUI interface
- Additional file formats and hash strategies
- Expanded language support
- Progress bars and better user feedback

---

## 🤝 Contributing

Contributions, bug reports, and feature requests are welcome! Please open an issue or submit a pull request.

### Before submitting:

1. Ensure your code follows Go best practices
2. Add appropriate documentation for new features
3. Test your changes thoroughly

---