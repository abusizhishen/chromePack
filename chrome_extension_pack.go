package main

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Output        string   `yaml:"output"`
	FilesToRemove []string `yaml:"files_to_remove"`
}

func main() {
	log("Starting script execution")

	projectPath, err := os.Getwd()
	if err != nil {
		log(fmt.Sprintf("Error getting current directory: %v", err))
		return
	}

	projectName := filepath.Base(projectPath)
	configFilePath := filepath.Join(getHomeDir(), ".chromePack.yaml")
	log("Loading configuration")
	config := loadConfig(configFilePath, projectName)

	outputProjectDir := filepath.Join(config.Output, projectName)

	log("Removing existing output directory")
	os.RemoveAll(outputProjectDir)

	log("Copying project to output directory")
	err = copyDir(projectPath, outputProjectDir)
	if err != nil {
		log(fmt.Sprintf("Error copying project: %v", err))
		return
	}

	log("Removing specified files from output directory")
	for _, file := range config.FilesToRemove {
		filePath := filepath.Join(outputProjectDir, file)
		os.RemoveAll(filePath)
	}

	log("Creating zip file")
	zipFilePath := filepath.Join(config.Output, fmt.Sprintf("%s.zip", projectName))
	err = createZip(zipFilePath, outputProjectDir)
	if err != nil {
		log(fmt.Sprintf("Error creating zip file: %v", err))
		return
	}

	log("Cleaning up temporary files")
	os.RemoveAll(outputProjectDir)
	log(fmt.Sprintf("Project archived successfully at %s", zipFilePath))
}

func loadConfig(configFilePath string, projectName string) Config {
	defaultConfig := Config{
		Output:        filepath.Join(getHomeDir(), "Downloads"),
		FilesToRemove: []string{".git", ".gitignore", ".idea", ".DS_Store", "build.sh"},
	}

	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		log("Configuration file not found. Creating default configuration file.")
		data, err := yaml.Marshal(defaultConfig)
		if err != nil {
			log(fmt.Sprintf("Error marshaling default config: %v", err))
			return defaultConfig
		}
		err = ioutil.WriteFile(configFilePath, data, 0644)
		if err != nil {
			log(fmt.Sprintf("Error writing default config file: %v", err))
		}
		log(fmt.Sprintf("Default config file created at %s", configFilePath))
		return defaultConfig
	}

	log("Loading existing configuration file")
	data, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		log(fmt.Sprintf("Error reading config file: %v", err))
		return defaultConfig
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log(fmt.Sprintf("Error parsing config file: %v", err))
		return defaultConfig
	}

	return config
}

func createZip(zipPath, srcDir string) error {
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	return filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath := strings.TrimPrefix(path, srcDir+string(filepath.Separator))
		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		zipEntry, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}

		_, err = io.Copy(zipEntry, file)
		return err
	})
}

func log(message string) {
	fmt.Println(message)
}
