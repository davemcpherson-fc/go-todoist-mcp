package main

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds all configuration for the application.
type Config struct {
	Todoist struct {
		APIToken string `yaml:"api_token"`
		BaseURL  string `yaml:"base_url"`
	} `yaml:"todoist"`
	Server struct {
		Port string `yaml:"port"`
	} `yaml:"server"`
}

// LoadConfig reads configuration from 'config.yaml' and environment variables.
func LoadConfig() *Config {
	f, err := os.Open("config.yaml")
	if err != nil {
		log.Fatalf("FATAL: Could not open config.yaml: %v", err)
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(&cfg); err != nil {
		log.Fatalf("FATAL: Could not decode config.yaml: %v", err)
	}

	// Override API token with environment variable if it exists
	if apiToken := os.Getenv("TODOIST_API_TOKEN"); apiToken != "" {
		cfg.Todoist.APIToken = apiToken
	}

	if cfg.Todoist.APIToken == "REPLACE_WITH_YOUR_TODOIST_API_TOKEN" || cfg.Todoist.APIToken == "" {
		log.Fatal("FATAL: Todoist API token is not set in config.yaml or as TODOIST_API_TOKEN env var")
	}

	return &cfg
}
