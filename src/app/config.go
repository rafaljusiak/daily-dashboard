package app

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	City           string
	ClockifyApiKey string
	HourlyRate     float64
	MongoURI       string
	Password       string
	Port           string
	RootDir        string
	WorkspaceId    string
}

func GetRootDir() string {
	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Dir(path)
}

func LoadConfig() Config {
	rootDir := GetRootDir()
	configPath := filepath.Join(rootDir, "config.json")
	file, err := os.Open(configPath)

	if err != nil {
		log.Fatalln("Missing config.json file")
	}

	defer file.Close()

	var config Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		log.Fatalln("Invalid config.json file structure")
	}

	config.RootDir = rootDir

	return config
}
