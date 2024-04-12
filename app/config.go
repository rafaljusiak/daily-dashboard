package app

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	ClockifyApiKey string
	HourlyRate     float64
	WorkspaceId    string
}

func LoadConfig() Config {
	path := "./config.json"
	file, err := os.Open(path)

	if err != nil {
		log.Fatalln("Missing ./config.json file")
	}

	defer file.Close()

	var config Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		log.Fatalln("Invalid config.json file structure")
	}

	return config
}
