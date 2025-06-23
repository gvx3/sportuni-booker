package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type CourseOption string
type AreaOption string

const (
	CourseBallGame CourseOption = "ball_games"
	// CourseOther    CourseOption = "other"
)
const (
	AreaHervanta   AreaOption = "hervanta"
	AreaKauppi     AreaOption = "kauppi"
	AreaCityCentre AreaOption = "citycentre"
)

var courseTypeDisplay = map[CourseOption]string{
	CourseBallGame: "Ball games",
	// CourseOther:    "Other",
}

var areaTypeDisplay = map[AreaOption]string{
	AreaHervanta:   "Hervanta",
	AreaKauppi:     "Kauppi",
	AreaCityCentre: "City centre",
}

type Config struct {
	BaseURL       string         `yaml:"base_url"`
	Email         string         `yaml:"email"`
	Password      string         `yaml:"password"`
	StateFileName string         `yaml:"state_file_name"`
	ActivitySlots []ActivitySlot `yaml:"activity_slots"`
	CourseType    CourseOption   `yaml:"course_type"`
	CourseArea    AreaOption     `yaml:"course_area"`
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
		if err := config.Validate(); err != nil {
			return nil, fmt.Errorf("config validation failed: %w", err)
		}
		return config, nil
	}

	return &Config{
		BaseURL:       getEnv("SPORTUNI_BASE_URL", "https://www.tuni.fi/sportuni/omasivu/?newPage=selection&lang=en"),
		Email:         getEnv("SPORTUNI_EMAIL", ""),
		Password:      getEnv("SPORTUNI_PASSWORD", ""),
		StateFileName: getEnv("SPORTUNI_STATE_FILE", "ms_user.json"),
		CourseType:    CourseBallGame, // default
		CourseArea:    AreaHervanta,   // default
	}, nil
}

func (c *Config) Validate() error {
	if _, ok := courseTypeDisplay[c.CourseType]; !ok {
		options := make([]string, 0, len(courseTypeDisplay))
		for k := range courseTypeDisplay {
			options = append(options, string(k))
		}
		return fmt.Errorf("invalid course_type: %s\n available options are: %s", c.CourseType, strings.Join(options, " | "))
	}

	if _, ok := areaTypeDisplay[c.CourseArea]; !ok {
		options := make([]string, 0, len(areaTypeDisplay))
		for k := range areaTypeDisplay {
			options = append(options, string(k))
		}
		return fmt.Errorf("invalid course_area: %s\n available options are: %s", c.CourseArea, strings.Join(options, " | "))
	}
	return nil
}

func (c *Config) DisplayCourseOption() string {
	return courseTypeDisplay[c.CourseType]
}

func (c *Config) DisplayCourseArea() string {
	return areaTypeDisplay[c.CourseArea]
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
