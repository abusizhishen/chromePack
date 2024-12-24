package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run script.go <new_folder_name>")
		return
	}

	newFolderName := os.Args[1]
	filesToRemove := []string{".git", ".idea"}

	sourceDir := path.Join(getHomeDir(), "WebstormProjects", "chrome_ext_base")
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		return
	}

	newFolderPath := filepath.Join(filepath.Dir(currentDir), newFolderName)
	fmt.Printf("Copying folder to: %s\n", newFolderPath)

	err = copyDir(sourceDir, newFolderPath)
	if err != nil {
		fmt.Printf("Error copying folder: %v\n", err)
		return
	}

	fmt.Println("Removing specified files...")
	for _, file := range filesToRemove {
		filePath := filepath.Join(newFolderPath, file)
		if err := os.RemoveAll(filePath); err != nil {
			fmt.Printf("Error removing file %s: %v\n", file, err)
		}
	}

	fmt.Println("Running initialization command...")
	cmd := exec.Command("git", "init")
	cmd.Dir = newFolderPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error running initialization command: %v\n", err)
		return
	}

	fmt.Println("Process completed successfully.")
}
