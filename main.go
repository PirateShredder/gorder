package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// categoryMap defines grouping of extensions into categories
var categoryMap = map[string][]string{
	"Images":        {"jpg", "jpeg", "png", "gif", "webp", "bmp", "tiff", "tif", "svg", "heic", "heif", "ico", "raw", "cr2", "nef", "orf", "arw", "psb", "dds", "hdr", "jp2"},
	"Videos":        {"mp4", "mov", "avi", "mkv", "wmv", "flv", "mpeg", "mpg", "m4v", "3gp", "webm", "vob", "ts", "m2ts", "rm", "rmvb", "asf"},
	"Audio":         {"mp3", "wav", "aac", "flac", "ogg", "m4a", "wma", "alac", "aiff", "amr", "mid", "midi", "opus", "pcm"},
	"Documents":     {"doc", "docx", "pdf", "txt", "rtf", "odt", "md", "epub", "tex", "ps", "pages", "djvu", "fodt", "rtfd"},
	"Spreadsheets":  {"xls", "xlsx", "csv", "ods", "tsv", "xlsm", "xlsb", "numbers"},
	"Presentations": {"ppt", "pptx", "odp", "key", "pps", "ppsx"},
	"Archives":      {"zip", "tar", "tar.gz", "tgz", "rar", "7z", "xz", "iso", "bz2", "gz", "lz", "lzma", "cab", "zst", "arj"},
	"Executables":   {"exe", "msi", "bat", "cmd", "apk", "aab", "ipa", "dmg", "pkg", "app", "deb", "rpm", "flatpak", "snap", "jar", "war", "bin", "sh"},
	"Web":           {"html", "htm", "css", "js", "ts", "jsx", "tsx"},
	"Data":          {"json", "xml", "yaml", "yml", "ini", "toml", "ndjson"},
	"Code":          {"c", "h", "cpp", "hpp", "cs", "java", "kt", "py", "rb", "php", "go", "rs", "swift", "scala", "lua", "pl", "ps1", "sql", "r", "m", "asm", "dart"},
	"Design":        {"psd", "psb", "ai", "eps", "indd", "xd", "fig", "sketch", "cdr", "afdesign", "afphoto", "afpub"},
	"Fonts":         {"ttf", "otf", "woff", "woff2", "eot", "fon", "pfb", "pfa"},
	"3D":            {"blend", "fbx", "obj", "stl", "3ds", "dae", "ply", "glb", "gltf", "max", "usd", "usdz"},
	"CAD":           {"dwg", "dxf", "dwt", "stp", "step", "iges", "igs", "sldprt", "sldasm", "ipt", "iam"},
	"Config":        {"log", "cfg", "conf", "env", "editorconfig", "properties", "jsonc", "reg"},
	"Database":      {"db", "sqlite", "sqlite3", "mdb", "accdb", "dbf", "parquet", "feather", "hdf5", "h5"},
	"Backup":        {"bak", "tmp", "old", "backup", "swp", "swo"},
}

type moveAction struct {
	from string
	to   string
}

var logFile *os.File
var moveLog []moveAction

func main() {
	// Custom usage/help message
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Gorder - The Advanced File Organizer

USAGE:
    gorder [options]

DESCRIPTION:
    Intelligently organizes files using multiple strategies: extension-based,
    category-based, date-based grouping, or fetch/flatten mode. Includes
    advanced features like duplicate detection, report generation, and undo.

ORGANIZATION MODES:
    -c, -categories              Group files by categories (Images, Documents, etc.)
    --date-mode <mode>           Group by date: 'year', 'month', 'day', or 'week'

FILE HANDLING:
    -d, -dry, -dryrun           Preview changes without moving files
    -f, -full                   Use full extensions (e.g., tar.gz instead of gz)
    -q, -quiet                  Use simple folder names without gorder_ prefix
    --noext-folder <name>       Folder name for files without extensions
    --case-sensitive            Treat extensions as case-sensitive

FILTERING:
    -i, -include <list>         Comma-separated list of extensions to include
    -e, -exclude <list>         Comma-separated list of extensions to exclude

ADVANCED OPTIONS:
    -r, -recursive              Process subdirectories recursively
    -t, -target <dir>           Target directory for organized folders
    -p, --fetch, --flatten      Flatten directory by pulling files from subdirs
    --cleanup                   Remove empty subdirectories (use with -p)
    -u, -undo                   Undo the last organization operation

ANALYSIS & REPORTS:
    -R, --report                Generate detailed directory analysis (gorder_report.md)
    -D, --duplicates            Detect and report duplicate files (gorder_dups.md)
    --delete-dups               Delete duplicates (use with -D, requires confirmation)

EXAMPLES:
    gorder                      # Organize files by extension (default)
    gorder -c                   # Organize by categories
    gorder -d -c                # Preview category organization
    gorder --date-mode month    # Organize by month
    gorder -r -c                # Recursively organize by categories
    gorder -p --cleanup         # Flatten directory structure
    gorder -R                   # Generate directory report
    gorder -D                   # Find duplicate files
    gorder -u                   # Undo last operation

For more information, see README.md

`)
	}

	// Flags
	dryRun := flag.Bool("dry", false, "Show what would happen, but don't move files")
	flag.BoolVar(dryRun, "d", false, "Show what would happen, but don't move files (shorthand)")
	flag.BoolVar(dryRun, "dryrun", false, "Show what would happen, but don't move files")

	useFullExt := flag.Bool("full", false, "Use full extension (e.g. .tar.gz) instead of only last extension")
	flag.BoolVar(useFullExt, "f", false, "Use full extension (e.g. .tar.gz) instead of only last extension (shorthand)")

	useCategories := flag.Bool("categories", false, "Group files by categories instead of individual extensions")
	flag.BoolVar(useCategories, "c", false, "Group files by categories instead of individual extensions (shorthand)")

	includeList := flag.String("include", "", "Comma-separated list of extensions or patterns to include (e.g., '.,.gitignore')")
	flag.StringVar(includeList, "i", "", "Comma-separated list of extensions or patterns to include (shorthand)")

	excludeList := flag.String("exclude", "", "Comma-separated list of extensions or patterns to exclude")
	flag.StringVar(excludeList, "e", "", "Comma-separated list of extensions or patterns to exclude (shorthand)")

	noExtFolder := flag.String("noext-folder", "", "Folder name for files without extensions (e.g., 'NoExtension')")

	caseSensitive := flag.Bool("case-sensitive", false, "Treat extensions as case-sensitive (e.g., .JPG vs .jpg)")

	dateMode := flag.String("date-mode", "", "Group by date: 'year', 'month', 'day', or 'week'")

	recursive := flag.Bool("recursive", false, "Process subdirectories recursively")
	flag.BoolVar(recursive, "r", false, "Process subdirectories recursively (shorthand)")

	targetDir := flag.String("target", ".", "Target directory for organized folders")
	flag.StringVar(targetDir, "t", ".", "Target directory for organized folders (shorthand)")

	undo := flag.Bool("undo", false, "Undo the last organization operation")
	flag.BoolVar(undo, "u", false, "Undo the last organization operation (shorthand)")

	fetch := flag.Bool("fetch", false, "Flatten directory structure by moving all files from subdirectories to current directory")
	flag.BoolVar(fetch, "flatten", false, "Flatten directory structure by moving all files from subdirectories to current directory")
	flag.BoolVar(fetch, "p", false, "Flatten directory structure by moving all files from subdirectories to current directory (shorthand)")

	cleanup := flag.Bool("cleanup", false, "Remove empty subdirectories after fetch operation (use with --fetch)")

	quiet := flag.Bool("quiet", false, "Use simple folder names without gorder_ prefix (e.g., 'jpg' instead of 'gorder_jpg')")
	flag.BoolVar(quiet, "q", false, "Use simple folder names without gorder_ prefix (shorthand)")

	report := flag.Bool("report", false, "Generate a detailed report (gorder_report.md) of directory contents")
	flag.BoolVar(report, "R", false, "Generate a detailed report (gorder_report.md) of directory contents (shorthand)")

	duplicates := flag.Bool("duplicates", false, "Detect and report duplicate files (gorder_dups.md)")
	flag.BoolVar(duplicates, "D", false, "Detect and report duplicate files (shorthand)")

	deleteDups := flag.Bool("delete-dups", false, "Delete duplicate files (use with --duplicates, keeps first instance)")

	flag.Parse()

	// Handle undo mode
	if *undo {
		performUndo()
		return
	}

	// Handle fetch mode
	if *fetch {
		performFetch(*dryRun, *cleanup)
		return
	}

	// Handle report generation
	if *report {
		generateReport()
		return
	}

	// Handle duplicate detection
	if *duplicates {
		findDuplicates(*deleteDups)
		return
	}

	// Initialize log file for undo functionality
	if !*dryRun {
		var err error
		logFile, err = os.Create(".gorder_log.txt")
		if err != nil {
			log.Printf("Warning: Could not create log file: %v", err)
		} else {
			defer logFile.Close()
		}
	}

	// Parse include/exclude lists
	includeSet := parseList(*includeList)
	excludeSet := parseList(*excludeList)

	// Build extension to category map if using categories
	extToCat := make(map[string]string)
	if *useCategories {
		for category, exts := range categoryMap {
			for _, ext := range exts {
				extToCat[ext] = category
			}
		}
	}

	// Process directory
	if *recursive {
		processDirectoryRecursive(".", *dryRun, *useFullExt, *useCategories, *caseSensitive, *dateMode, *noExtFolder, includeSet, excludeSet, extToCat, *targetDir, *quiet)
	} else {
		processDirectory(".", *dryRun, *useFullExt, *useCategories, *caseSensitive, *dateMode, *noExtFolder, includeSet, excludeSet, extToCat, *targetDir, *quiet)
	}
}

func processDirectory(dir string, dryRun, useFullExt, useCategories, caseSensitive bool, dateMode, noExtFolder string, includeSet, excludeSet map[string]bool, extToCat map[string]string, targetDir string, quiet bool) {
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		name := file.Name()

		// Skip log files
		if name == ".gorder_log.txt" {
			continue
		}

		// Check if file should be processed based on include/exclude
		if !shouldProcess(name, includeSet, excludeSet) {
			continue
		}

		// Get file info for date-based grouping
		var folderName string
		if dateMode != "" {
			info, err := file.Info()
			if err != nil {
				log.Printf("Cannot get file info for %s: %v\n", name, err)
				continue
			}
			folderName = getDateFolder(info.ModTime(), dateMode)
		} else {
			ext := getExtension(name, useFullExt)

			// Handle files without extension
			if ext == "" {
				if noExtFolder != "" {
					folderName = noExtFolder
				} else {
					continue
				}
			} else {
				if !caseSensitive {
					ext = strings.ToLower(ext)
				}

				// Determine folder name
				if useCategories {
					if cat, ok := extToCat[strings.ToLower(ext)]; ok {
						folderName = cat
					} else {
						folderName = ext
					}
				} else {
					if quiet {
						folderName = ext
					} else {
						folderName = fmt.Sprintf("gorder_%s", ext)
					}
				}
			}
		}

		// Create folder if needed
		fullFolderPath := filepath.Join(targetDir, folderName)
		if _, err := os.Stat(fullFolderPath); errors.Is(err, os.ErrNotExist) {
			if !dryRun {
				if err := os.MkdirAll(fullFolderPath, 0755); err != nil {
					log.Printf("Cannot create folder %s: %v\n", fullFolderPath, err)
					continue
				}
			}
			fmt.Printf("[+] Created folder: %s\n", fullFolderPath)
		}

		oldPath := filepath.Join(dir, name)
		newPath := filepath.Join(fullFolderPath, name)

		// Handle name collision
		newPath = avoidCollision(newPath)

		if dryRun {
			fmt.Printf("[DRY] Would move %s → %s\n", oldPath, newPath)
			continue
		}

		if err := os.Rename(oldPath, newPath); err != nil {
			log.Printf("Error moving %s: %v\n", oldPath, err)
		} else {
			fmt.Printf("Moved %s → %s\n", oldPath, newPath)
			// Log for undo
			if logFile != nil {
				fmt.Fprintf(logFile, "%s|%s\n", newPath, oldPath)
			}
		}
	}
}

func processDirectoryRecursive(dir string, dryRun, useFullExt, useCategories, caseSensitive bool, dateMode, noExtFolder string, includeSet, excludeSet map[string]bool, extToCat map[string]string, targetDir string, quiet bool) {
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		name := info.Name()

		// Skip log files and files in target directory
		if name == ".gorder_log.txt" || strings.HasPrefix(path, targetDir+string(os.PathSeparator)) {
			return nil
		}

		// Check if file should be processed
		if !shouldProcess(name, includeSet, excludeSet) {
			return nil
		}

		// Get folder name based on mode
		var folderName string
		if dateMode != "" {
			folderName = getDateFolder(info.ModTime(), dateMode)
		} else {
			ext := getExtension(name, useFullExt)

			if ext == "" {
				if noExtFolder != "" {
					folderName = noExtFolder
				} else {
					return nil
				}
			} else {
				if !caseSensitive {
					ext = strings.ToLower(ext)
				}

				if useCategories {
					if cat, ok := extToCat[strings.ToLower(ext)]; ok {
						folderName = cat
					} else {
						folderName = ext
					}
				} else {
					if quiet {
						folderName = ext
					} else {
						folderName = fmt.Sprintf("gorder_%s", ext)
					}
				}
			}
		}

		// Create folder if needed
		fullFolderPath := filepath.Join(targetDir, folderName)
		if _, err := os.Stat(fullFolderPath); errors.Is(err, os.ErrNotExist) {
			if !dryRun {
				if err := os.MkdirAll(fullFolderPath, 0755); err != nil {
					log.Printf("Cannot create folder %s: %v\n", fullFolderPath, err)
					return nil
				}
			}
			fmt.Printf("[+] Created folder: %s\n", fullFolderPath)
		}

		newPath := filepath.Join(fullFolderPath, name)

		// Handle name collision
		newPath = avoidCollision(newPath)

		if dryRun {
			fmt.Printf("[DRY] Would move %s → %s\n", path, newPath)
			return nil
		}

		if err := os.Rename(path, newPath); err != nil {
			log.Printf("Error moving %s: %v\n", path, err)
		} else {
			fmt.Printf("Moved %s → %s\n", path, newPath)
			// Log for undo
			if logFile != nil {
				fmt.Fprintf(logFile, "%s|%s\n", newPath, path)
			}
		}

		return nil
	})

	if err != nil {
		log.Printf("Error walking directory: %v", err)
	}
}

func getExtension(name string, full bool) string {
	if full {
		parts := strings.Split(name, ".")
		if len(parts) > 2 {
			return strings.Join(parts[len(parts)-2:], ".")
		}
	}
	return strings.TrimPrefix(filepath.Ext(name), ".")
}

func avoidCollision(path string) string {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return path
	}

	ext := filepath.Ext(path)
	base := strings.TrimSuffix(filepath.Base(path), ext)
	dir := filepath.Dir(path)

	// Try sequential numbering
	for i := 1; ; i++ {
		newName := fmt.Sprintf("%s (%d)%s", base, i, ext)
		newPath := filepath.Join(dir, newName)
		if _, err := os.Stat(newPath); errors.Is(err, os.ErrNotExist) {
			return newPath
		}
	}
}

func shouldProcess(name string, includeSet, excludeSet map[string]bool) bool {
	// Skip hidden files by default unless explicitly included
	if strings.HasPrefix(name, ".") {
		if len(includeSet) == 0 || !includeSet[name] {
			return false
		}
	}

	// Check exclude list
	ext := filepath.Ext(name)
	if excludeSet[ext] || excludeSet[name] {
		return false
	}

	// If include list is specified, only process included files
	if len(includeSet) > 0 {
		return includeSet[ext] || includeSet[name] || includeSet["."]
	}

	return true
}

func parseList(list string) map[string]bool {
	result := make(map[string]bool)
	if list == "" {
		return result
	}

	items := strings.Split(list, ",")
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item != "" {
			result[item] = true
		}
	}
	return result
}

func getDateFolder(modTime time.Time, mode string) string {
	switch mode {
	case "year":
		return fmt.Sprintf("%d", modTime.Year())
	case "month":
		return fmt.Sprintf("%04d-%02d", modTime.Year(), modTime.Month())
	case "day":
		return fmt.Sprintf("%04d-%02d-%02d", modTime.Year(), modTime.Month(), modTime.Day())
	case "week":
		_, week := modTime.ISOWeek()
		return fmt.Sprintf("Week_%02d", week)
	default:
		return "Unknown"
	}
}

func performUndo() {
	// Read the log file
	file, err := os.Open(".gorder_log.txt")
	if err != nil {
		log.Fatal("Cannot open log file. No previous operation to undo.")
	}
	defer file.Close()

	var actions []moveAction
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "|")
		if len(parts) == 2 {
			actions = append(actions, moveAction{from: parts[0], to: parts[1]})
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("Error reading log file:", err)
	}

	if len(actions) == 0 {
		fmt.Println("No actions to undo.")
		return
	}

	fmt.Printf("Undoing %d file moves...\n", len(actions))

	successCount := 0
	for _, action := range actions {
		if err := os.Rename(action.from, action.to); err != nil {
			log.Printf("Error moving %s back to %s: %v\n", action.from, action.to, err)
		} else {
			fmt.Printf("Restored %s → %s\n", action.from, action.to)
			successCount++
		}
	}

	fmt.Printf("\nUndo complete: %d/%d files restored.\n", successCount, len(actions))

	// Remove the log file after successful undo
	if successCount > 0 {
		os.Remove(".gorder_log.txt")
	}
}

func performFetch(dryRun, cleanup bool) {
	fmt.Println("Fetching files from subdirectories...")

	var filesToMove []moveAction

	// Walk through all subdirectories
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the current directory
		if path == "." {
			return nil
		}

		// Skip the log file
		if filepath.Base(path) == ".gorder_log.txt" {
			return nil
		}

		// If it's a file and not in current directory, add to move list
		if !info.IsDir() {
			// Check if file is in a subdirectory
			dir := filepath.Dir(path)
			if dir != "." {
				// Target is current directory with collision avoidance
				targetPath := filepath.Join(".", filepath.Base(path))
				filesToMove = append(filesToMove, moveAction{
					from: path,
					to:   targetPath,
				})
			}
		}

		return nil
	})

	if err != nil {
		log.Fatal("Error walking directory:", err)
	}

	if len(filesToMove) == 0 {
		fmt.Println("No files found in subdirectories.")
		return
	}

	fmt.Printf("Found %d files in subdirectories\n", len(filesToMove))

	// Move files
	successCount := 0
	for _, action := range filesToMove {
		// Apply collision avoidance
		finalPath := avoidCollision(action.to)

		if dryRun {
			fmt.Printf("[DRY] Would move %s → %s\n", action.from, finalPath)
			successCount++
		} else {
			if err := os.Rename(action.from, finalPath); err != nil {
				log.Printf("Error moving %s: %v\n", action.from, err)
			} else {
				fmt.Printf("Moved %s → %s\n", action.from, finalPath)
				successCount++
			}
		}
	}

	fmt.Printf("\nFetch complete: %d/%d files moved to current directory\n", successCount, len(filesToMove))

	// Cleanup empty directories if requested
	if cleanup && !dryRun {
		fmt.Println("\nCleaning up empty directories...")
		removeEmptyDirs(".")
	} else if cleanup && dryRun {
		fmt.Println("\n[DRY] Would clean up empty directories")
	}
}

func removeEmptyDirs(root string) {
	// Walk bottom-up to remove nested empty directories
	var dirs []string

	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && path != "." {
			dirs = append(dirs, path)
		}

		return nil
	})

	// Reverse order to process deepest directories first
	for i := len(dirs) - 1; i >= 0; i-- {
		dir := dirs[i]

		// Check if directory is empty
		entries, err := os.ReadDir(dir)
		if err != nil {
			log.Printf("Error reading directory %s: %v\n", dir, err)
			continue
		}

		if len(entries) == 0 {
			if err := os.Remove(dir); err != nil {
				log.Printf("Error removing empty directory %s: %v\n", dir, err)
			} else {
				fmt.Printf("Removed empty directory: %s\n", dir)
			}
		}
	}
}

type fileInfo struct {
	path string
	size int64
	ext  string
}

type extStats struct {
	count int
	size  int64
}

func generateReport() {
	fmt.Println("Generating directory report...")

	var allFiles []fileInfo
	extMap := make(map[string]*extStats)
	var totalSize int64
	var totalDirs int
	var hiddenFiles int
	var noExtFiles int

	// Walk through all files and directories
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if path != "." {
				totalDirs++
			}
			return nil
		}

		// Track hidden files
		if strings.HasPrefix(filepath.Base(path), ".") {
			hiddenFiles++
			return nil
		}

		ext := filepath.Ext(path)
		if ext == "" {
			noExtFiles++
			ext = "(no extension)"
		} else {
			ext = strings.ToLower(ext)
		}

		size := info.Size()
		totalSize += size

		// Track file
		allFiles = append(allFiles, fileInfo{
			path: path,
			size: size,
			ext:  ext,
		})

		// Track extension stats
		if _, ok := extMap[ext]; !ok {
			extMap[ext] = &extStats{}
		}
		extMap[ext].count++
		extMap[ext].size += size

		return nil
	})

	if err != nil {
		log.Fatal("Error scanning directory:", err)
	}

	// Sort files by size (largest first)
	sort.Slice(allFiles, func(i, j int) bool {
		return allFiles[i].size > allFiles[j].size
	})

	// Sort extensions by total size (for bar chart)
	type extEntry struct {
		ext   string
		stats *extStats
	}
	var extList []extEntry
	for ext, stats := range extMap {
		extList = append(extList, extEntry{ext, stats})
	}
	sort.Slice(extList, func(i, j int) bool {
		return extList[i].stats.size > extList[j].stats.size
	})

	// Generate report
	reportFile, err := os.Create("gorder_report.md")
	if err != nil {
		log.Fatal("Error creating report file:", err)
	}
	defer reportFile.Close()

	w := bufio.NewWriter(reportFile)
	defer w.Flush()

	// Header
	fmt.Fprintf(w, "# Gorder Directory Report\n\n")
	fmt.Fprintf(w, "Generated: %s\n\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Fprintf(w, "---\n\n")

	// Directory Summary
	fmt.Fprintf(w, "## 📊 Directory Summary\n\n")
	fmt.Fprintf(w, "- **Total Files**: %d\n", len(allFiles))
	fmt.Fprintf(w, "- **Total Size**: %s\n", formatSize(totalSize))
	fmt.Fprintf(w, "- **Total Directories**: %d\n", totalDirs)
	fmt.Fprintf(w, "- **Hidden Files**: %d (skipped from analysis)\n", hiddenFiles)
	fmt.Fprintf(w, "- **Files Without Extension**: %d\n\n", noExtFiles)

	// File Type Distribution
	fmt.Fprintf(w, "## 📁 File Type Distribution\n\n")
	fmt.Fprintf(w, "| Extension | Count | Total Size |\n")
	fmt.Fprintf(w, "|-----------|-------|------------|\n")
	for _, entry := range extList {
		fmt.Fprintf(w, "| %s | %d | %s |\n", entry.ext, entry.stats.count, formatSize(entry.stats.size))
	}
	fmt.Fprintf(w, "\n")

	// Top 5 File Types Bar Chart
	if len(extList) > 0 {
		fmt.Fprintf(w, "## 📈 Top 5 File Types (by size)\n\n")
		maxWidth := 50
		topN := 5
		if len(extList) < topN {
			topN = len(extList)
		}
		maxSize := extList[0].stats.size
		for i := 0; i < topN; i++ {
			entry := extList[i]
			barWidth := int(float64(entry.stats.size) / float64(maxSize) * float64(maxWidth))
			if barWidth == 0 && entry.stats.size > 0 {
				barWidth = 1
			}
			bar := strings.Repeat("█", barWidth)
			fmt.Fprintf(w, "%15s | %-50s %s\n", entry.ext, bar, formatSize(entry.stats.size))
		}
		fmt.Fprintf(w, "\n")
	}

	// Top 10 Largest Files
	fmt.Fprintf(w, "## 📦 Top 10 Largest Files\n\n")
	fmt.Fprintf(w, "| # | File Path | Size |\n")
	fmt.Fprintf(w, "|---|-----------|------|\n")
	topFiles := 10
	if len(allFiles) < topFiles {
		topFiles = len(allFiles)
	}
	for i := 0; i < topFiles; i++ {
		file := allFiles[i]
		fmt.Fprintf(w, "| %d | %s | %s |\n", i+1, file.path, formatSize(file.size))
	}

	fmt.Printf("\n✅ Report generated: gorder_report.md\n")
	fmt.Printf("   Total files analyzed: %d\n", len(allFiles))
	fmt.Printf("   Total size: %s\n", formatSize(totalSize))
}

func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func findDuplicates(deleteDups bool) {
	fmt.Println("Scanning for duplicate files...")

	type fileHash struct {
		path string
		size int64
	}

	hashMap := make(map[string][]fileHash)
	var totalFiles int
	var totalSize int64

	// Calculate hash for all files
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// Skip hidden files
		if strings.HasPrefix(filepath.Base(path), ".") {
			return nil
		}

		// Skip report files
		if filepath.Base(path) == "gorder_report.md" || filepath.Base(path) == "gorder_dups.md" {
			return nil
		}

		totalFiles++
		totalSize += info.Size()

		// Calculate MD5 hash
		hash, err := hashFile(path)
		if err != nil {
			log.Printf("Error hashing %s: %v\n", path, err)
			return nil
		}

		hashMap[hash] = append(hashMap[hash], fileHash{
			path: path,
			size: info.Size(),
		})

		return nil
	})

	if err != nil {
		log.Fatal("Error scanning for duplicates:", err)
	}

	// Find duplicates (hashes with more than one file)
	var duplicateGroups [][]fileHash
	var duplicateCount int
	var duplicateSize int64

	for _, files := range hashMap {
		if len(files) > 1 {
			duplicateGroups = append(duplicateGroups, files)
			duplicateCount += len(files) - 1 // Don't count the first instance
			for i := 1; i < len(files); i++ {
				duplicateSize += files[i].size
			}
		}
	}

	if len(duplicateGroups) == 0 {
		fmt.Println("\n✅ No duplicate files found!")
		return
	}

	// Sort groups by size (largest first)
	sort.Slice(duplicateGroups, func(i, j int) bool {
		return duplicateGroups[i][0].size > duplicateGroups[j][0].size
	})

	// Generate report
	reportFile, err := os.Create("gorder_dups.md")
	if err != nil {
		log.Fatal("Error creating duplicates report:", err)
	}
	defer reportFile.Close()

	w := bufio.NewWriter(reportFile)
	defer w.Flush()

	// Header
	fmt.Fprintf(w, "# Gorder Duplicate Files Report\n\n")
	fmt.Fprintf(w, "Generated: %s\n\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Fprintf(w, "---\n\n")

	// Summary
	fmt.Fprintf(w, "## 📊 Summary\n\n")
	fmt.Fprintf(w, "- **Total Files Scanned**: %d\n", totalFiles)
	fmt.Fprintf(w, "- **Duplicate Groups**: %d\n", len(duplicateGroups))
	fmt.Fprintf(w, "- **Duplicate Files**: %d\n", duplicateCount)
	fmt.Fprintf(w, "- **Wasted Space**: %s\n\n", formatSize(duplicateSize))

	// Duplicate Groups
	fmt.Fprintf(w, "## 🔍 Duplicate Groups\n\n")

	for i, group := range duplicateGroups {
		fmt.Fprintf(w, "### Group %d (Size: %s, %d copies)\n\n", i+1, formatSize(group[0].size), len(group))
		for j, file := range group {
			if j == 0 {
				fmt.Fprintf(w, "- **[KEEP]** %s\n", file.path)
			} else {
				fmt.Fprintf(w, "- %s\n", file.path)
			}
		}
		fmt.Fprintf(w, "\n")
	}

	fmt.Printf("\n✅ Duplicates report generated: gorder_dups.md\n")
	fmt.Printf("   Duplicate groups: %d\n", len(duplicateGroups))
	fmt.Printf("   Duplicate files: %d\n", duplicateCount)
	fmt.Printf("   Wasted space: %s\n", formatSize(duplicateSize))

	// Handle deletion if requested
	if deleteDups {
		fmt.Printf("\n⚠️  Delete duplicates mode enabled!\n")
		fmt.Printf("   This will delete %d duplicate files (keeping first instance of each).\n", duplicateCount)
		fmt.Printf("   Type 'yes' to confirm deletion: ")

		var response string
		fmt.Scanln(&response)

		if response == "yes" {
			deletedCount := 0
			deletedSize := int64(0)

			for _, group := range duplicateGroups {
				// Keep first file, delete the rest
				for i := 1; i < len(group); i++ {
					if err := os.Remove(group[i].path); err != nil {
						log.Printf("Error deleting %s: %v\n", group[i].path, err)
					} else {
						fmt.Printf("Deleted: %s\n", group[i].path)
						deletedCount++
						deletedSize += group[i].size
					}
				}
			}

			fmt.Printf("\n✅ Deletion complete!\n")
			fmt.Printf("   Files deleted: %d\n", deletedCount)
			fmt.Printf("   Space freed: %s\n", formatSize(deletedSize))
		} else {
			fmt.Println("Deletion cancelled.")
		}
	}
}

func hashFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
