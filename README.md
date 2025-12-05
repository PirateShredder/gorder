# 📂 Gorder - The Advanced File Organizer

**Tired of cluttered directories? `gorder` is a powerful, feature-rich command-line tool that automatically organizes files using intelligent categorization, date-based grouping, and extensive customization options.**

[![Go Report Card](https://goreportcard.com/badge/github.com/PirateShredder/gorder)](https://goreportcard.com/report/github.com/PirateShredder/gorder)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

---

## 🚀 What it does

`gorder` intelligently organizes files in your directory using multiple organization strategies:

- **Extension-based**: Group files by file type (default)
- **Category-based**: Group common file types into logical categories (Images, Documents, Videos, etc.)
- **Date-based**: Organize files by modification date (year, month, day, or week)
- **Fetch/Flatten**: Pull all files from subdirectories to current directory
- **Report Generation**: Analyze directory contents with detailed statistics and visualizations
- **Duplicate Detection**: Find and optionally remove duplicate files
- **Custom filtering**: Include/exclude specific file types
- **Recursive processing**: Organize entire directory trees
- **Undo support**: Reverse any organization operation

**Before `gorder`:**
```
.
├── document.pdf
├── image.jpg
├── notes.txt
├── presentation.pptx
├── report.docx
├── video.mp4
└── music.mp3
```

**After `gorder` (default mode):**
```
.
├── gorder_docx/
│   └── report.docx
├── gorder_jpg/
│   └── image.jpg
├── gorder_mp3/
│   └── music.mp3
├── gorder_mp4/
│   └── video.mp4
├── gorder_pdf/
│   └── document.pdf
├── gorder_pptx/
│   └── presentation.pptx
└── gorder_txt/
    └── notes.txt
```

**After `gorder -c` (category mode):**
```
.
├── Audio/
│   └── music.mp3
├── Documents/
│   ├── document.pdf
│   ├── notes.txt
│   └── report.docx
├── Images/
│   └── image.jpg
├── Presentations/
│   └── presentation.pptx
└── Videos/
    └── video.mp4
```

## 📦 Installation

Ensure you have Go installed and your `GOPATH` is set up correctly.

```sh
go install github.com/PirateShredder/gorder@latest
```

## 🛠️ Usage

### Basic Usage

Simply navigate to the directory you want to organize and run:

```sh
gorder
```

For help and a list of all options:

```sh
gorder -h
# or
gorder --help
```

### Command-Line Options

#### Organization Modes

- **`-c`, `-categories`**: Group files by categories instead of individual extensions
  ```sh
  gorder -c
  # Images: jpg, png, gif, etc. → Images/
  # Documents: pdf, txt, doc, etc. → Documents/
  # Videos: mp4, avi, mkv, etc. → Videos/
  ```

- **`--date-mode <mode>`**: Group files by modification date
  - `year`: Group by year (e.g., `2024/`, `2025/`)
  - `month`: Group by year-month (e.g., `2024-12/`, `2025-01/`)
  - `day`: Group by full date (e.g., `2024-12-04/`)
  - `week`: Group by calendar week (e.g., `Week_49/`)
  ```sh
  gorder --date-mode month
  ```

#### File Handling

- **`-d`, `-dry`, `-dryrun`**: Preview changes without moving files
  ```sh
  gorder -d  # See what would happen
  ```

- **`-f`, `-full`**: Use full extensions (e.g., `tar.gz` instead of just `gz`)
  ```sh
  gorder -f
  ```

- **`--noext-folder <name>`**: Specify folder for files without extensions
  ```sh
  gorder --noext-folder NoExtension
  # README, LICENSE, Makefile → NoExtension/
  ```

- **`--case-sensitive`**: Treat extensions as case-sensitive
  ```sh
  gorder --case-sensitive
  # .jpg → gorder_jpg/
  # .JPG → gorder_JPG/
  ```

- **`-q`, `--quiet`**: Use simple folder names without `gorder_` prefix
  ```sh
  gorder -q
  # .jpg → jpg/  (instead of gorder_jpg/)
  # .pdf → pdf/  (instead of gorder_pdf/)
  ```

#### Filtering

- **`-i`, `-include <list>`**: Only process specified extensions (comma-separated)
  ```sh
  gorder -i .jpg,.png,.gif  # Only organize images
  gorder -i .,.gitignore    # Include files without extension
  ```

- **`-e`, `-exclude <list>`**: Exclude specified extensions (comma-separated)
  ```sh
  gorder -e .tmp,.bak  # Skip temporary and backup files
  ```

#### Advanced Options

- **`-r`, `-recursive`**: Process subdirectories recursively
  ```sh
  gorder -r  # Organize all files in all subdirectories
  ```

- **`-t`, `-target <dir>`**: Specify target directory for organized folders
  ```sh
  gorder -t ~/organized  # Create organized folders in ~/organized/
  ```

- **`-p`, `--fetch`, `--flatten`**: Flatten directory structure by moving all files from subdirectories to current directory
  ```sh
  gorder -p                # Pull all files from subdirs to current directory
  gorder -p --cleanup      # Also remove empty directories after fetch
  ```

- **`--cleanup`**: Remove empty subdirectories after fetch operation (use with `--fetch`)

- **`-u`, `-undo`**: Undo the last organization operation
  ```sh
  gorder -u  # Restore files to original locations
  ```

- **`-R`, `--report`**: Generate detailed directory analysis report
  ```sh
  gorder -R  # Creates gorder_report.md with statistics and visualizations
  ```

- **`-D`, `--duplicates`**: Detect and report duplicate files
  ```sh
  gorder -D                 # Creates gorder_dups.md with duplicate groups
  gorder -D --delete-dups   # Find and delete duplicates (with confirmation)
  ```

- **`--delete-dups`**: Delete duplicate files (use with `--duplicates`, keeps first instance)

### Example Workflows

**Organize photos by month:**
```sh
gorder --date-mode month -i .jpg,.png,.raw
```

**Clean up downloads with categories, excluding temporary files:**
```sh
gorder -c -e .tmp,.part
```

**Organize entire project recursively into separate folder:**
```sh
gorder -r -c -t ~/organized_project
```

**Preview organization before applying:**
```sh
gorder -d -c  # Check what will happen
gorder -c     # Apply if satisfied
```

**Organize with full extensions and handle files without extensions:**
```sh
gorder -f --noext-folder Scripts
```

**Undo if you made a mistake:**
```sh
gorder -u
```

**Flatten deeply nested directory structure:**
```sh
gorder -p --cleanup  # Pull all files to current dir and remove empty folders
```

**Use simple folder names without prefix:**
```sh
gorder -q  # Creates jpg/, pdf/, txt/ instead of gorder_jpg/, etc.
```

**Generate comprehensive directory report:**
```sh
gorder -R  # Analyze files, create report with statistics and visualizations
```

**Find and review duplicate files:**
```sh
gorder -D  # Scan for duplicates, save report to gorder_dups.md
```

**Clean up duplicate files:**
```sh
gorder -D --delete-dups  # Find duplicates and delete all but first instance
```

## 📂 Supported Categories

When using `-c` or `-categories`, files are grouped into these categories:

- **Images**: jpg, png, gif, svg, webp, bmp, tiff, heic, raw, etc.
- **Videos**: mp4, mov, avi, mkv, wmv, flv, webm, etc.
- **Audio**: mp3, wav, aac, flac, ogg, m4a, etc.
- **Documents**: pdf, doc, docx, txt, rtf, odt, md, epub, etc.
- **Spreadsheets**: xls, xlsx, csv, ods, numbers, etc.
- **Presentations**: ppt, pptx, odp, key, etc.
- **Archives**: zip, tar, rar, 7z, gz, iso, etc.
- **Executables**: exe, dmg, apk, deb, jar, etc.
- **Web**: html, css, js, ts, jsx, tsx
- **Data**: json, xml, yaml, yml, toml
- **Code**: c, cpp, py, java, go, rs, swift, etc.
- **Design**: psd, ai, eps, fig, sketch, etc.
- **Fonts**: ttf, otf, woff, woff2, etc.
- **3D**: blend, fbx, obj, stl, gltf, etc.
- **CAD**: dwg, dxf, step, etc.
- **Config**: log, ini, cfg, conf, env
- **Database**: db, sqlite, mdb, sql, etc.
- **Backup**: bak, tmp, old, backup, swp

## 🔧 Building from Source

1. **Clone the repository:**
   ```sh
   git clone https://github.com/PirateShredder/gorder.git
   cd gorder
   ```
2. **Build the binary:**
   ```sh
   go build -o gorder main.go
   ```
3. **Run it:**
   ```sh
   ./gorder
   ```

## ⚙️ Features

- ✅ Multiple organization strategies (extension, category, date)
- ✅ Fetch/flatten mode to pull files from subdirectories
- ✅ Quiet mode for simple folder names (without `gorder_` prefix)
- ✅ Dry-run mode to preview changes
- ✅ Recursive directory processing
- ✅ Include/exclude filters
- ✅ Custom target directory
- ✅ Smart collision handling with sequential numbering
- ✅ Undo functionality with action logging
- ✅ Automatic cleanup of empty directories
- ✅ Full extension support (e.g., `.tar.gz`)
- ✅ Case-sensitive mode
- ✅ No-extension file handling
- ✅ Comprehensive category mapping
- ✅ Report generation with file statistics and visualizations
- ✅ Duplicate file detection with optional deletion

## 📜 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
