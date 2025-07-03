package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// type AreaOption string

const (
	Badminton string = "Badminton"
	Billiard  string = "Billiards"
)
const (
	AreaHervanta   string = "hervanta"
	AreaKauppi     string = "kauppi"
	AreaCityCentre string = "citycentre"
)

const (
	BookingCourtType string = "BookingCourt"
	ReserveCourtType string = "ReserveCourt"
)

var sportAreaMap = map[string]string{
	Badminton: "Ball games",
	Billiard:  "Other",
}

var sportDialogMap = map[string]string{
	Badminton: "Sulkapallo",
	Billiard:  "Biljardi",
}

var areaTypeDisplay = map[string]string{
	AreaHervanta:   "Hervanta",
	AreaKauppi:     "Kauppi",
	AreaCityCentre: "City centre",
}

var sportBookingTypeMap = map[string]string{
	Badminton: BookingCourtType,
	Billiard:  ReserveCourtType,
}

type Config struct {
	BaseURL       string         `yaml:"base_url"`
	Email         string         `yaml:"email"`
	Password      string         `yaml:"password"`
	StateFileName string         `yaml:"state_file_name"`
	ActivitySlots []ActivitySlot `yaml:"activity_slots"`
}

type ActivitySlot struct {
	Day        string `yaml:"day"`
	Date       string `yaml:"date"`
	Hour       string `yaml:"hour"`
	Activity   string `yaml:"activity"`
	CourseArea string `yaml:"course_area"`
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
		ActivitySlots: config.ActivitySlots,
	}, nil
}

func (c *Config) Validate() error {

	for _, slot := range c.ActivitySlots {

		if _, ok := areaTypeDisplay[slot.CourseArea]; !ok {
			options := make([]string, 0, len(areaTypeDisplay))
			for k := range areaTypeDisplay {
				options = append(options, string(k))
			}
			return fmt.Errorf("invalid course_area: %s\n available options are: %s", slot.CourseArea, strings.Join(options, " | "))
		}

	}
	return nil
}

func (a *ActivitySlot) DisplaySportDialogMap(sport string) string {
	return sportDialogMap[sport]
}

func (a *ActivitySlot) DisplayCourseOption(sport string) string {
	return sportAreaMap[sport]
}

func (a *ActivitySlot) DisplayCourseArea(area string) string {
	return areaTypeDisplay[area]
}

func (a *ActivitySlot) DisplayBookingType(sport string) string {
	if t, ok := sportBookingTypeMap[a.Activity]; ok {
		return t
	}
	return "Unknown"
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
