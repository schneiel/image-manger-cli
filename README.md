# 🖼️ ImageManager

**ImageManager** is a cross-platform command-line tool for organizing and cleaning up image libraries. It automatically sorts images into a structured folder hierarchy and detects (and optionally removes) duplicate files safely and efficiently.

Built in Go — lightweight, fast, and ready for automation.

---

## 🔑 Features

- **Image Sorting**  
  Automatically organizes photos into a `YYYY/MM` folder structure based on capture date.

- **Duplicate Detection**  
  Identifies duplicate images using content-based file hashing.

- **Flexible Date Extraction**  
  Determines an image’s date using a prioritized fallback strategy:
  1. EXIF metadata (default)
  2. File modification time
  2. File creation time

- **Customizable Duplicate Handling**  
  - **Action Mode**: *Dry Run* (simulation) or *Move to Trash*
  - **Retention Strategy**: Keep the *oldest* file or the one with the *shortest file path*

- **Safe by Design**  
  Built-in safety mechanisms like *dry-run* mode and *trash-first* deletion help prevent accidental data loss.

- **Multilingual Support**  
  Currently supports English (`en`) and German (`de`).

- **Comprehensive Logging**  
  All operations are logged for transparency and troubleshooting.

---

## ⚙️ Installation

Clone the repository and build the binary using Go:

```bash
git clone git@github.com:schneiel/ImageManagerGo.git
cd ImageManagerGo
go build -o ImageManager ./main.go
```

The `ImageManager` binary will be created in the current directory.

---

## 🚀 Usage

ImageManager uses subcommands to perform operations. The two main ones are `sort` and `dedup`.

### 📁 Sort Images

Organize images from a source directory into `YYYY/MM` folders under a destination path:

```bash
./ImageManager sort --source "/path/to/images" --destination "/path/to/archive"
```

### 🧹 Find Duplicates

Scan a directory for duplicate images and process them based on the configured strategy:

```bash
./ImageManager dedup --source "/path/to/images" --lang "en"
```

---

### 🌐 Global Flags

These flags can be used with any subcommand:

- `--config`, `-c`: Path to a configuration file (default: `./config.yaml`)
- `--dry-run`: Simulates actions without making changes (no files are moved or deleted)

---

### 📜 Subcommands

#### `sort`

Organizes images into a date-based folder structure.

**Flags:**

- `--source`: Source directory (overrides `source` in config)
- `--destination`: Destination directory (overrides `destination` in config)

#### `dedup`

Detects and handles duplicate images.

**Flags:**

- `--source`: (Required) Directory to scan (overrides `source` in config)

---

## 🛠 Configuration (`config.yaml`)

You can define advanced options in a `config.yaml` file.

### Example:

```yaml
source: "/path/to/images"
destination: "/path/to/images_sorted"

dry_run: true

deduplicator:
  actionStrategy: "moveToTrash"         # Options: "dryRun", "moveToTrash"
  keepStrategy: "keepOldest"            # Options: "keepOldest", "keepShortestPath"

date:
  strategies:
    - "exif"
    - "creationTime"
    - "modTime"
```

---

## 📝 Logging

Two log files are created during execution:

- `application.log`: General logs and errors
- `dedup_dry_run_log.csv`: List of duplicates found (in dry-run mode), including paths and hashes
- `sort_dry_run_log.csv`: List of old paths and expected new paths after sortation 

---

## 🌐 Internationalization

Language selection via `--lang` flag.

**Supported Languages:**
- English (`en`)
- German (`de`)

**Example:**
```bash
./ImageManager --lang="de" dedup --source "/images"
```

---

## 📌 Planned Improvements (WIP)

- Configuration validation
- GUI
- Additional file formats and hash strategies
- Expanded language support

---

## 🤝 Contributing

Contributions, bug reports, and feature requests are welcome! Please open an issue or submit a pull request.

---
