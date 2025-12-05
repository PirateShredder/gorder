# Gorder - Installation and Usage Guide

## Table of Contents

- [System Requirements](#system-requirements)
- [Installation](#installation)
  - [Method 1: Install via Go](#method-1-install-via-go)
  - [Method 2: Build from Source](#method-2-build-from-source)
  - [Method 3: Download Binary](#method-3-download-binary)
- [Quick Start](#quick-start)
- [Organization Modes](#organization-modes)
  - [Extension Mode (Default)](#extension-mode-default)
  - [Category Mode](#category-mode)
  - [Date Mode](#date-mode)
- [Command Reference](#command-reference)
  - [Core Options](#core-options)
  - [File Selection](#file-selection)
  - [Advanced Options](#advanced-options)
- [Usage Examples](#usage-examples)
  - [Basic Usage](#basic-usage)
  - [Photo Management](#photo-management)
  - [Downloads Cleanup](#downloads-cleanup)
  - [Project Organization](#project-organization)
  - [Backup and Archive](#backup-and-archive)
- [Advanced Workflows](#advanced-workflows)
- [Understanding Categories](#understanding-categories)
- [Undo Operations](#undo-operations)
- [Best Practices](#best-practices)
- [Troubleshooting](#troubleshooting)
- [FAQ](#faq)

---

## System Requirements

- **Operating System**: Windows, macOS, Linux, or any OS supported by Go
- **Go Version**: 1.16 or higher (only needed for building from source)
- **Disk Space**: ~3 MB for binary
- **Permissions**: Read/write access to directories you want to organize

---

## Installation

### Method 1: Install via Go

If you have Go installed on your system:

```bash
go install github.com/PirateShredder/gorder@latest
```

This installs `gorder` to `$GOPATH/bin`. Ensure this directory is in your `PATH`.

**Verify installation:**
```bash
gorder -h
```

### Method 2: Build from Source

1. **Clone the repository:**
   ```bash
   git clone https://github.com/PirateShredder/gorder.git
   cd gorder
   ```

2. **Build the binary:**
   ```bash
   go build -o gorder main.go
   ```

3. **Install to system (optional):**
   ```bash
   # macOS/Linux
   sudo mv gorder /usr/local/bin/

   # Or add to your PATH
   export PATH=$PATH:$(pwd)
   ```

4. **Windows users:**
   ```cmd
   go build -o gorder.exe main.go
   # Move gorder.exe to a directory in your PATH
   ```

### Method 3: Download Binary

Download pre-compiled binaries from the [Releases](https://github.com/PirateShredder/gorder/releases) page.

1. Download the appropriate binary for your OS
2. Extract the archive
3. Move the binary to a directory in your PATH
4. Make it executable (macOS/Linux): `chmod +x gorder`

---

## Quick Start

### First Run

Navigate to a directory you want to organize:

```bash
cd ~/Downloads
```

**Preview what will happen (dry-run):**
```bash
gorder -d
```

**Organize files by extension:**
```bash
gorder
```

**Organize files by category:**
```bash
gorder -c
```

**Undo if you change your mind:**
```bash
gorder -u
```

---

## Organization Modes

### Extension Mode (Default)

Creates a folder for each file extension with the `gorder_` prefix.

**Example:**
```
Before:
├── photo.jpg
├── document.pdf
├── song.mp3

After:
├── gorder_jpg/
│   └── photo.jpg
├── gorder_pdf/
│   └── document.pdf
└── gorder_mp3/
    └── song.mp3
```

**Usage:**
```bash
gorder                # Basic extension mode
gorder -f             # Use full extensions (.tar.gz instead of .gz)
gorder --case-sensitive  # Treat .JPG and .jpg separately
gorder -q             # Use simple folder names (jpg/ instead of gorder_jpg/)
```

### Category Mode

Groups related file types into logical categories (Images, Documents, Videos, etc.).

**Example:**
```
Before:
├── photo.jpg
├── image.png
├── document.pdf
├── notes.txt
├── video.mp4

After:
├── Images/
│   ├── photo.jpg
│   └── image.png
├── Documents/
│   ├── document.pdf
│   └── notes.txt
└── Videos/
    └── video.mp4
```

**Usage:**
```bash
gorder -c             # Category mode (shorthand)
gorder --categories   # Category mode (full flag)
```

**Supported Categories:**
- **Images**: jpg, png, gif, svg, webp, bmp, tiff, heic, raw, and 15+ more
- **Videos**: mp4, mov, avi, mkv, wmv, flv, webm, and 10+ more
- **Audio**: mp3, wav, aac, flac, ogg, m4a, and 7+ more
- **Documents**: pdf, doc, docx, txt, rtf, odt, md, epub, and 6+ more
- **Spreadsheets**: xls, xlsx, csv, ods, numbers, and 3+ more
- **Presentations**: ppt, pptx, odp, key, pps, ppsx
- **Archives**: zip, tar, rar, 7z, gz, iso, and 11+ more
- **Executables**: exe, dmg, apk, deb, jar, and 12+ more
- **Web**: html, css, js, ts, jsx, tsx
- **Data**: json, xml, yaml, yml, toml
- **Code**: c, cpp, py, java, go, rs, swift, and 14+ more
- **Design**: psd, ai, eps, fig, sketch, and 7+ more
- **Fonts**: ttf, otf, woff, woff2
- **3D**: blend, fbx, obj, stl, gltf, and 7+ more
- **CAD**: dwg, dxf, step, and 8+ more
- **Config**: log, ini, cfg, conf, env
- **Database**: db, sqlite, mdb, sql, and 6+ more
- **Backup**: bak, tmp, old, backup, swp

Files with extensions not in these categories will be placed in individual folders.

### Date Mode

Organizes files based on their modification timestamp.

**Modes:**
- `year`: Group by year (2024/, 2025/)
- `month`: Group by year-month (2024-12/, 2025-01/)
- `day`: Group by full date (2024-12-04/)
- `week`: Group by ISO week number (Week_49/)

**Example (month mode):**
```
Before:
├── vacation_2024.jpg (modified: Jan 2024)
├── report.pdf (modified: Dec 2024)
├── photo.jpg (modified: Jan 2024)

After:
├── 2024-01/
│   ├── vacation_2024.jpg
│   └── photo.jpg
└── 2024-12/
    └── report.pdf
```

**Usage:**
```bash
gorder --date-mode year    # Organize by year
gorder --date-mode month   # Organize by month
gorder --date-mode day     # Organize by day
gorder --date-mode week    # Organize by ISO week
```

---

## Command Reference

### Core Options

| Flag | Shorthand | Description |
|------|-----------|-------------|
| `--dry` | `-d`, `--dryrun` | Preview changes without moving files |
| `--categories` | `-c` | Use category-based grouping |
| `--date-mode <mode>` | - | Group by date (year/month/day/week) |
| `--full` | `-f` | Use full extensions (e.g., .tar.gz) |
| `--undo` | `-u` | Undo the last organization operation |

### File Selection

| Flag | Shorthand | Description |
|------|-----------|-------------|
| `--include <list>` | `-i` | Only process specified extensions (comma-separated) |
| `--exclude <list>` | `-e` | Skip specified extensions (comma-separated) |
| `--noext-folder <name>` | - | Folder name for files without extensions |
| `--case-sensitive` | - | Treat .JPG and .jpg as different |
| `--quiet` | `-q` | Use simple folder names without `gorder_` prefix |

### Advanced Options

| Flag | Shorthand | Description |
|------|-----------|-------------|
| `--recursive` | `-r` | Process subdirectories recursively |
| `--target <dir>` | `-t` | Target directory for organized folders |

---

## Usage Examples

### Basic Usage

**Preview before organizing:**
```bash
cd ~/Desktop
gorder -d
```

**Organize by extension:**
```bash
gorder
```

**Organize by category:**
```bash
gorder -c
```

### Photo Management

**Organize photos by month:**
```bash
cd ~/Pictures/Unsorted
gorder --date-mode month
```

**Only organize image files:**
```bash
gorder -i .jpg,.png,.gif,.raw -c
```

**Organize camera files by year (including RAW formats):**
```bash
gorder --date-mode year -i .jpg,.cr2,.nef,.arw
```

**Separate JPEG and JPG files (case-sensitive cameras):**
```bash
gorder --case-sensitive -i .jpg,.JPG
```

**Use simple folder names without prefix:**
```bash
gorder -q  # Creates jpg/, pdf/, txt/ instead of gorder_jpg/, etc.
```

### Downloads Cleanup

**Organize downloads by category, excluding temporary files:**
```bash
cd ~/Downloads
gorder -c -e .tmp,.part,.crdownload
```

**Only organize archives and executables:**
```bash
gorder -c -i .zip,.rar,.7z,.exe,.dmg
```

**Preview organization with all options:**
```bash
gorder -d -c -e .tmp
```

### Project Organization

**Organize entire project recursively:**
```bash
cd ~/Projects/myproject
gorder -r -c
```

**Organize into separate directory:**
```bash
gorder -r -c -t ~/Projects/myproject_organized
```

**Organize code files only:**
```bash
gorder -c -i .py,.js,.go,.rs,.java
```

**Handle files without extensions (README, LICENSE, Makefile):**
```bash
gorder -c --noext-folder ProjectFiles
```

### Backup and Archive

**Organize by year for archival:**
```bash
cd ~/OldFiles
gorder --date-mode year
```

**Organize with full extension support (important for .tar.gz files):**
```bash
gorder -f -c
```

**Exclude backup files from organization:**
```bash
gorder -c -e .bak,.old,.backup
```

---

## Advanced Workflows

### Workflow 1: Careful Photo Organization

```bash
# Step 1: Preview what will happen
cd ~/Photos/Import
gorder -d --date-mode month -i .jpg,.png,.heic

# Step 2: If satisfied, organize
gorder --date-mode month -i .jpg,.png,.heic

# Step 3: If you made a mistake, undo
gorder -u
```

### Workflow 2: Deep Project Cleanup

```bash
# Organize entire project tree by category into separate folder
cd ~/Projects/legacy-project
gorder -r -c -t ~/Projects/legacy-organized -e .tmp,.swp,.log

# Review results
ls -la ~/Projects/legacy-organized

# Undo if needed
cd ~/Projects/legacy-project
gorder -u
```

### Workflow 3: Downloads Folder Maintenance

```bash
# Weekly downloads cleanup
cd ~/Downloads

# First, remove obvious junk
gorder -c -i .tmp,.part,.crdownload -t ~/.Trash

# Then organize everything else by category
gorder -c -e .tmp,.part

# Keep archives separate with full extensions
cd ~/Downloads/Archives
gorder -f --noext-folder Misc
```

### Workflow 4: Camera Card Import

```bash
# Import from camera card
cd /Volumes/CAMERA_SD

# Organize by day (useful for multi-day shoots)
gorder --date-mode day -i .jpg,.cr2 -t ~/Photos/Import/2024-12-shoot

# Verify import
ls -la ~/Photos/Import/2024-12-shoot

# Separate RAW and JPEG for editing workflow
cd ~/Photos/Import/2024-12-shoot
gorder -c -i .cr2 -t RAW
gorder -c -i .jpg -t JPEG
```

---

## Understanding Categories

The category system groups 100+ file extensions into 15 logical categories. Here's how it works:

**Category Mapping (main.go:16-35):**
- Each category contains a list of related file extensions
- During categorization, `gorder` looks up the file's extension in this map
- If found, the file goes to that category folder
- If not found, it falls back to creating an individual `gorder_<ext>` folder

**Example:**
```
File: vacation.jpg
Extension: jpg
Lookup: categoryMap["Images"] contains "jpg"
Result: Moved to Images/vacation.jpg

File: custom.xyz
Extension: xyz
Lookup: Not in any category
Result: Moved to gorder_xyz/custom.xyz
```

**To see all categories:**
```bash
gorder -h | grep -A1 categories
```

---

## Undo Operations

Every organization operation creates a log file (`.gorder_log.txt`) that tracks all file moves.

**How it works:**
1. Each file move is logged: `new_path|original_path`
2. Dry-run mode (`-d`) does NOT create a log
3. Running `gorder -u` reads the log and reverses all moves
4. After successful undo, the log file is deleted

**Usage:**
```bash
# Organize files
gorder -c

# Changed your mind? Undo
gorder -u
# Output: "Undoing 47 file moves..."
# Output: "Undo complete: 47/47 files restored."
```

**Important notes:**
- Undo only works for the most recent operation
- If you run `gorder` again, the old log is overwritten
- Dry-run doesn't affect the log file
- Manual file moves after organization may cause undo issues

**Log file location:**
```
.gorder_log.txt (in the directory where you ran gorder)
```

**Log file format:**
```
Images/photo.jpg|photo.jpg
Documents/report.pdf|report.pdf
Videos/movie.mp4|movie.mp4
```

---

## Best Practices

### 1. Always Dry-Run First

Before organizing important directories:
```bash
gorder -d -c  # Preview the changes
```

### 2. Start Small

Test on a small directory before running on your entire file system:
```bash
mkdir test_gorder
cd test_gorder
# Add some test files
gorder -d -c
```

### 3. Use Include/Exclude Wisely

Be specific about what you want to organize:
```bash
# Good: Specific file types
gorder -i .jpg,.png

# Better: Exclude temporary files
gorder -c -e .tmp,.swp,.log
```

### 4. Backup Important Data

For critical files, create a backup before organizing:
```bash
cp -r ~/Important ~/Important_backup
cd ~/Important
gorder -c
```

### 5. Check Results Before Continuing

After organizing:
```bash
ls -la       # Check folder structure
gorder -u    # Undo if something looks wrong
```

### 6. Use Target Directory for Safety

Organize into a separate directory to preserve original structure:
```bash
gorder -c -t ~/Organized
# Original files stay in place until you verify
```

### 7. Understand Your File Types

Know what extensions you have:
```bash
# Find all extensions in current directory
find . -type f | perl -ne 'print $1."\n" if m/\.([^.\/]+)$/' | sort -u
```

### 8. Handle Archives Carefully

Use full extension mode for archives:
```bash
gorder -f  # Preserves .tar.gz instead of just .gz
```

### 9. Regular Maintenance

Set up a regular cleanup routine:
```bash
# Weekly downloads cleanup script
#!/bin/bash
cd ~/Downloads
gorder -d -c -e .tmp  # Preview
read -p "Proceed? (y/n) " -n 1 -r
if [[ $REPLY =~ ^[Yy]$ ]]; then
    gorder -c -e .tmp
fi
```

### 10. Document Your Workflow

Keep notes on your organization strategy:
```bash
# In ~/organization_notes.md
Downloads: Use category mode, exclude .tmp files
Photos: Use date-mode month, include .jpg,.heic
Projects: Use recursive with target directory
```

---

## Troubleshooting

### Problem: "Permission denied" error

**Cause:** No write permission in the directory

**Solution:**
```bash
# Check permissions
ls -la

# Fix permissions (if you own the directory)
chmod u+w .

# Or run from a directory where you have permissions
```

### Problem: Files not being organized

**Cause:** Files may be hidden or excluded

**Solution:**
```bash
# Check for hidden files
ls -la

# Include hidden files explicitly
gorder -i .gitignore,. --noext-folder Hidden

# Check if files are being excluded
gorder -d  # Look for skipped files
```

### Problem: Wrong category assignments

**Cause:** Extension not in category map

**Solution:**
```bash
# Extensions not in categories fall back to individual folders
# This is expected behavior

# If you want custom categories, you'd need to modify main.go
# Or use extension mode: gorder (without -c)
```

### Problem: "Cannot create folder" error

**Cause:** Target directory doesn't exist or invalid path

**Solution:**
```bash
# Create target directory first
mkdir -p ~/Organized

# Then run gorder
gorder -c -t ~/Organized
```

### Problem: Undo doesn't work

**Cause:** Log file missing or corrupted

**Solution:**
```bash
# Check if log exists
ls -la .gorder_log.txt

# If missing, you can't undo (sorry!)
# If corrupted, manual restoration needed

# Prevention: Don't delete .gorder_log.txt manually
```

### Problem: Duplicate files created

**Cause:** Collision handling is working as designed

**Solution:**
```bash
# Collisions are renamed: file.txt -> file (1).txt
# This prevents data loss

# To find duplicates:
find . -name "* (*.*"

# To merge manually if they're actually duplicates
```

### Problem: Recursive mode too aggressive

**Cause:** Organizing subdirectories you didn't want to touch

**Solution:**
```bash
# Use non-recursive mode for current directory only
gorder -c  # No -r flag

# Or use include/exclude to be selective
gorder -r -c -e .git,.node_modules
```

### Problem: Running out of disk space

**Cause:** Large file operations

**Solution:**
```bash
# Check disk space first
df -h .

# Files are moved (not copied), so space shouldn't increase
# If it does, check for the log file size
du -h .gorder_log.txt
```

---

## FAQ

### Q: Does gorder move or copy files?

**A:** Gorder **moves** files (not copies). The original file is relocated to the organized folder. No duplicates are created unless there's a collision.

### Q: Can I organize files across different drives?

**A:** Yes, but it depends on your OS. On the same filesystem, it's a fast move. Across filesystems, it may be slower as it becomes a copy+delete operation.

### Q: What happens if I run gorder multiple times?

**A:** Each run organizes the current state. Previously organized files are skipped if they're already in folders. The log file is overwritten, so you can only undo the most recent operation.

### Q: Can I customize the category mappings?

**A:** Not from the command line. You'd need to edit `main.go` (lines 16-35) and rebuild. Look for the `categoryMap` variable.

### Q: Does gorder work with symbolic links?

**A:** Gorder treats symlinks as regular files and will move them. The link is moved, not the target file.

### Q: Will gorder organize already organized folders?

**A:** If you run gorder in a directory with `gorder_*` folders or category folders, it will skip those directories and only process files in the current directory.

### Q: Can I organize files by size?

**A:** Not currently. Only extension, category, or date modes are supported.

### Q: What about files with multiple extensions like .tar.gz?

**A:** Use the `-f` flag: `gorder -f`. This treats `.tar.gz` as one extension instead of just `.gz`.

### Q: Is there a GUI version?

**A:** No, gorder is command-line only. However, it's easy to wrap in a shell script or create a simple GUI wrapper.

### Q: Can I exclude entire folders from recursive processing?

**A:** Not directly by folder name. You could use `-e` with common extensions in those folders, or manually move them before running gorder.

### Q: What happens to file metadata (dates, permissions)?

**A:** File metadata is preserved during the move operation. Modification times, creation times, and permissions remain unchanged.

### Q: Can I run gorder automatically on a schedule?

**A:** Yes! Use cron (Linux/macOS) or Task Scheduler (Windows):
```bash
# Cron example (every Sunday at 2 AM)
0 2 * * 0 cd ~/Downloads && /usr/local/bin/gorder -c -e .tmp
```

### Q: Does gorder support regex patterns for include/exclude?

**A:** No, it only supports exact extension matches and comma-separated lists. For regex, you'd need to pipe file lists or modify the source code.

### Q: Can I undo after running gorder multiple times?

**A:** No, only the most recent operation can be undone. The log file is overwritten each time.

### Q: What's the maximum number of files gorder can handle?

**A:** There's no hard limit. It's constrained only by available memory and filesystem limitations. Tested with 10,000+ files successfully.

### Q: Does gorder work on network drives/NAS?

**A:** Yes, if the network drive is mounted and you have read/write permissions. Performance depends on network speed.

### Q: Can I contribute new categories or features?

**A:** Yes! The project is open source. Fork the repository, make changes, and submit a pull request.

---

## Getting Help

If you encounter issues not covered in this guide:

1. Check the error message carefully
2. Run with `-d` (dry-run) to diagnose
3. Review the troubleshooting section above
4. Check existing issues on GitHub
5. Open a new issue with:
   - Your OS and gorder version
   - The exact command you ran
   - The error message
   - Expected vs actual behavior

**Happy organizing! 🎉**
