package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type fileInfo struct {
	name string
	size int64
	path string
}

func main() {
	// Define command-line flags
	dirFlag := flag.String("dir", ".", "directory to search for duplicates")
	deleteFlag := flag.Bool("delete", false, "whether to delete duplicates")
	interactiveFlag := flag.Bool("interactive", false, "whether to prompt before deleting duplicates")
	helpFlag := flag.Bool("help", false, "display help")

	// Parse command-line flags
	flag.Parse()

	// Display help message
	if *helpFlag {
		fmt.Println("Usage: go run main.go [OPTIONS]\n\nOptions:")
		flag.PrintDefaults()
		return
	}

	// Get the root directory
	rootDir := *dirFlag

	// Walk the directory tree and find duplicates
	duplicateGroups := findDuplicates(rootDir)

	// If no duplicates were found, exit
	if len(duplicateGroups) == 0 {
		fmt.Println("No duplicate files found.")
		return
	}

	// Print the duplicates
	printDuplicates(duplicateGroups)

	// Delete duplicates if the delete flag is set
	if *deleteFlag {
		deleteFiles(duplicateGroups, *interactiveFlag)
		deleteEmptyDirs(rootDir)
	}
}

// Function to find duplicate files
func findDuplicates(rootDir string) map[string][]fileInfo {
	// Create a map to store file groups
	fileGroups := make(map[string][]fileInfo)

	// Create a mutex to synchronize access to the map
	var mu sync.Mutex

	// Create a WaitGroup to keep track of the goroutines
	var wg sync.WaitGroup

	// Traverse the directory tree concurrently
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		// If the path is a regular file, add it to the appropriate group
		if info.Mode().IsRegular() {
			wg.Add(1)
			go func() {
				defer wg.Done()

				// Get the file size and name without extension
				size := info.Size()
				name := info.Name()
				cutName := strings.TrimSuffix(name, filepath.Ext(name))

				// Create a new fileInfo struct
				file := fileInfo{name: name, size: size, path: path}

				// Add the file to the appropriate group
				key := fmt.Sprintf("%s%d", cutName, file.size)

				mu.Lock()
				fileGroups[key] = append(fileGroups[key], file)
				mu.Unlock()
			}()
		}
		return nil
	})

	// Wait for all goroutines to finish
	wg.Wait()

	// Check for errors during traversal
	if err != nil {
		fmt.Printf("Error traversing directory: %v\n", err)
	}

	// Filter out groups with only one file
	duplicateGroups := make(map[string][]fileInfo)
	for key, group := range fileGroups {
		if len(group) > 1 {
			duplicateGroups[key] = group
		}
	}

	return duplicateGroups
}

// Function to print duplicate files
func printDuplicates(duplicateGroups map[string][]fileInfo) {
	fmt.Println("Duplicate files found:")
	for _, group := range duplicateGroups {
		fmt.Printf("- %d files:\n", len(group))
		for _, file := range group {
			fmt.Printf("  - %s (%s)\n", file.path, formatSize(file.size))
		}
	}
	fmt.Println("")
}

// Function to delete duplicate files
func deleteFiles(duplicateGroups map[string][]fileInfo, interactive bool) {
	// Loop through each group of duplicates
	for _, paths := range duplicateGroups {
		// Loop through each file in the group
		if len(paths) > 1 {
			// Keep the first path and delete the others
			keepFile := paths[0]

			for _, file := range paths[1:] {
				// If we are not deleting all files, prompt the user for confirmation
				if interactive {
					// Print the file path and prompt the user for confirmation
					fmt.Printf("Do you want to remove duplicate %s for %s? [y/N]: ", file.path, keepFile.path)

					var response string
					fmt.Scanln(&response)
					if response != "y" && response != "Y" {
						// If the user does not confirm, skip to the next file
						continue
					}
				}
				// Attempt to delete the file
				err := os.Remove(file.path)
				if err != nil {
					fmt.Printf("Error deleting file %s: %s\n", file.path, err.Error())
					continue
				}
				// Print a message indicating the file was deleted
				fmt.Printf("Deleted file: %s\n", file.path)
			}
		}
	}
}

func deleteEmptyDirs(rootDir string) error {
	// Create a stack of directories to traverse
	stack := []string{rootDir}

	// Traverse the directory tree
	for len(stack) > 0 {
		// Get the last directory from the stack
		dir := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		// Read the contents of the directory
		entries, err := os.ReadDir(dir)
		if err != nil {
			return err
		}

		// Check if the directory is empty
		if len(entries) == 0 {
			// Remove the empty directory
			if err := os.Remove(dir); err != nil {
				return err
			}
			fmt.Printf("Deleted directory: %s\n", dir)

			// Add the parent directory to the stack
			parentDir := filepath.Dir(dir)
			if parentDir != rootDir {
				stack = append(stack, parentDir)
			}
		} else {
			// Add any subdirectories to the stack
			for _, entry := range entries {
				if entry.IsDir() {
					stack = append(stack, filepath.Join(dir, entry.Name()))
				}
			}
		}
	}

	return nil
}

func formatSize(size int64) string {
	const unit = 1000
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "kMGTPE"[exp])
}
