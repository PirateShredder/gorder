package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	// 1. Read all files and folders in the current directory
	files, err := os.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	// 2. Loop through all the files
	for _, file := range files {
		if !file.IsDir() {
			// 3. Get the file extension
			ext := filepath.Ext(file.Name())
			if ext == "" {
				continue // Skip files with no extension
			}
		fileType := ext[1:] // Remove the dot

			// 4. Create a folder for the file type if it doesn't exist
			if _, err := os.Stat(fileType); os.IsNotExist(err) {
				if err := os.Mkdir(fileType, 0755); err != nil {
					log.Fatal(err)
				}
			}

			// 5. Move the file to the new folder
			oldPath := file.Name()
			newPath := filepath.Join(fileType, file.Name())
			if err := os.Rename(oldPath, newPath); err != nil {
				log.Printf("Error moving file %s: %v", oldPath, err)
			} else {
				fmt.Printf("Moved %s to %s\n", oldPath, newPath)
			}
		}
	}
}
