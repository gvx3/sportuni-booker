package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	BaseURL       string         `yaml:"base_url"`
	Email         string         `yaml:"email"`
	Password      string         `yaml:"password"`
	StateFileName string         `yaml:"state_file_name"`
	ActivitySlots []ActivitySlot `yaml:"activity_slots"`
}

type ActivitySlot struct {
	Day      string `yaml:"day"`
	Date     string `yaml:"date"`
	Hour     string `yaml:"hour"`
	Activity string `yaml:"activity"`
}

func NewConfig() (*Config, error) {

	configPath, err := getConfigPath()
	if err != nil {
		return nil, fmt.Errorf("failed to determine config path: %w", err)
	}

	config, err := loadConfigFromFile(configPath)
	if err == nil {
		return config, nil
	}

	return &Config{
		BaseURL:       getEnv("SPORTUNI_BASE_URL", "https://www.tuni.fi/sportuni/omasivu/?newPage=selection&lang=en"),
		Email:         getEnv("SPORTUNI_EMAIL", ""),
		Password:      getEnv("SPORTUNI_PASSWORD", ""),
		StateFileName: getEnv("SPORTUNI_STATE_FILE", "ms_user.json"),
	}, nil
}

func getConfigPath() (string, error) {

	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %w", err)
	}

	configPath := filepath.Join(currentDir, "config.yaml")
	if _, err := os.Stat(configPath); err == nil {
		return configPath, nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	homeConfigPath := filepath.Join(homeDir, ".sportuni", "config.yaml")
	if _, err := os.Stat(homeConfigPath); err == nil {
		return homeConfigPath, nil
	}

	return configPath, nil
}

func loadConfigFromFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error parsing config file: %w", err)
	}

	return &config, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
