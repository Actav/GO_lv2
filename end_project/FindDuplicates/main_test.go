package main

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestFindDuplicates(t *testing.T) {
	// Create a temporary directory with files and duplicates
	tempDir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatalf("Error creating temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create files with duplicates
	file1 := createFile(tempDir, "file1.txt", 100)
	file2 := createFile(tempDir, "file2.txt", 200)
	file3 := createFile(tempDir, "file1.old", 100)
	file4 := createFile(tempDir, "file2.old", 200)

	// Call findDuplicates function
	result := findDuplicates(tempDir)

	// Create expected result
	expected := map[string][]fileInfo{
		"file1100": []fileInfo{file1, file3},
		"file2200": []fileInfo{file2, file4},
	}

	// Check that the number of keys in the result is the same as in the expected map
	if len(result) != len(expected) {
		t.Errorf("Expected %v keys in result but got %v", len(expected), len(result))
	}

	// Check each key in the expected result
	for key, expectedFiles := range expected {
		// Check if the key exists in the result
		if resultFiles, ok := result[key]; ok {
			// Check if the number of files for the key is the same in the result and expected
			if len(resultFiles) != len(expectedFiles) {
				t.Errorf("Expected %v files for key %v but got %v", len(expectedFiles), key, len(resultFiles))
			}
			// Check if the files in the result are the same as in the expected result
			for _, expectedFile := range expectedFiles {
				if !containsFileInfo(resultFiles, expectedFile) {
					t.Errorf("Expected file %v for key %v but not found in result", expectedFile, key)
				}
			}
		} else {
			t.Errorf("Expected key %v but not found in result", key)
		}
	}
}

func TestFindDuplicatesWithError(t *testing.T) {
	// Create a temporary directory
	tmpDir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Call findDuplicates function with non-existent directory
	result := findDuplicates(tmpDir)

	// Check that function returned empty map
	if len(result) != 0 {
		t.Errorf("Expected empty map but got %v", result)
	}
}

func TestDeleteEmptyDirs(t *testing.T) {
	// Create temporary directory and files
	rootDir, err := ioutil.TempDir("", "test_delete_empty_dirs")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(rootDir)

	dir1 := filepath.Join(rootDir, "dir1")
	err = os.MkdirAll(dir1, 0777)
	if err != nil {
		t.Fatal(err)
	}

	dir2 := filepath.Join(rootDir, "dir2")
	err = os.MkdirAll(dir2, 0777)
	if err != nil {
		t.Fatal(err)
	}

	file2 := filepath.Join(dir2, "file2.txt")
	err = ioutil.WriteFile(file2, []byte("test"), 0666)
	if err != nil {
		t.Fatal(err)
	}

	// Ensure directories are deleted
	err = deleteEmptyDirs(rootDir)
	if err != nil {
		t.Fatal(err)
	}

	// Check if directories were deleted
	if _, err := os.Stat(dir1); err == nil {
		t.Errorf("Expected dir1 %s to be deleted", dir1)
	}

	if _, err := os.Stat(dir2); err != nil {
		t.Errorf("Expected dir2 %s to not be deleted", dir2)
	}
}

func TestFormatSize(t *testing.T) {
	testCases := []struct {
		size int64
		want string
	}{
		{0, "0 B"},
		{500, "500 B"},
		{999, "999 B"},
		{1000, "1.0 kB"},
		{1500, "1.5 kB"},
		{999999, "1000.0 kB"},
		{1000000, "1.0 MB"},
		{1500000, "1.5 MB"},
		{1000000000, "1.0 GB"},
		{1500000000, "1.5 GB"},
	}

	for _, tc := range testCases {
		got := formatSize(tc.size)
		if got != tc.want {
			t.Errorf("formatSize(%d) = %s; want %s", tc.size, got, tc.want)
		}
	}
}

// containsFileInfo returns true if the given slice of fileInfo contains the given fileInfo
func containsFileInfo(files []fileInfo, fileInfo fileInfo) bool {
	for _, f := range files {
		if f.name == fileInfo.name && f.size == fileInfo.size && f.path == fileInfo.path {
			return true
		}
	}
	return false
}

func createFile(dir string, name string, size int64) fileInfo {
	// Create file path
	path := filepath.Join(dir, name)

	// Create file with random data
	file, err := os.Create(path)
	if err != nil {
		panic(fmt.Sprintf("Error creating file %s: %v", path, err))
	}
	defer file.Close()

	// Write random data to file
	data := make([]byte, size)
	_, err = rand.Read(data)
	if err != nil {
		panic(fmt.Sprintf("Error writing data to file %s: %v", path, err))
	}
	_, err = file.Write(data)
	if err != nil {
		panic(fmt.Sprintf("Error writing data to file %s: %v", path, err))
	}

	// Create and return fileInfo struct
	return fileInfo{name: name, size: size, path: path}
}
