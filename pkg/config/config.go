package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config represents the application configuration
type Config struct {
	AutoPot struct {
		HPPotKey   string `json:"hp_pot_key"`
		SPPotKey   string `json:"sp_pot_key"`
		HPPercent  int    `json:"hp_percent"`
		SPPercent  int    `json:"sp_percent"`
		CheckDelay int    `json:"check_delay_ms"`
		Enabled    bool   `json:"enabled"`
	} `json:"autopot"`

	AutoBuff struct {
		Enabled bool `json:"enabled"`
		Delay   int  `json:"delay_ms"`
	} `json:"autobuff"`

	Spammer struct {
		Enabled     bool     `json:"enabled"`
		Keys        []string `json:"keys"`
		Delay       int      `json:"delay_ms"`
		SkillCancel bool     `json:"skill_cancel"`
		DoubleClick bool     `json:"double_click"`
	} `json:"spammer"`

	UI struct {
		ShowDebug bool   `json:"show_debug"`
		LogLevel  string `json:"log_level"`
	} `json:"ui"`

	SkillMarkers struct {
		Enabled bool    `json:"enabled"`
		Opacity float32 `json:"opacity"`
		Size    float32 `json:"size"`
		Skills  map[string]struct {
			Color   string `json:"color"`
			Enabled bool   `json:"enabled"`
		} `json:"skills"`
	} `json:"skill_markers"`

	RangeDisplay struct {
		Enabled   bool    `json:"enabled"`
		Color     string  `json:"color"`
		Opacity   float32 `json:"opacity"`
		Thickness float32 `json:"thickness"`
		Segments  int     `json:"segments"`
	} `json:"range_display"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		AutoPot: struct {
			HPPotKey   string `json:"hp_pot_key"`
			SPPotKey   string `json:"sp_pot_key"`
			HPPercent  int    `json:"hp_percent"`
			SPPercent  int    `json:"sp_percent"`
			CheckDelay int    `json:"check_delay_ms"`
			Enabled    bool   `json:"enabled"`
		}{
			HPPotKey:   "F1",
			SPPotKey:   "F2",
			HPPercent:  50,
			SPPercent:  30,
			CheckDelay: 100,
			Enabled:    true,
		},
		AutoBuff: struct {
			Enabled bool `json:"enabled"`
			Delay   int  `json:"delay_ms"`
		}{
			Enabled: false,
			Delay:   5000,
		},
		Spammer: struct {
			Enabled     bool     `json:"enabled"`
			Keys        []string `json:"keys"`
			Delay       int      `json:"delay_ms"`
			SkillCancel bool     `json:"skill_cancel"`
			DoubleClick bool     `json:"double_click"`
		}{
			Enabled:     false,
			Keys:        []string{},
			Delay:       50,
			SkillCancel: false,
			DoubleClick: false,
		},
		UI: struct {
			ShowDebug bool   `json:"show_debug"`
			LogLevel  string `json:"log_level"`
		}{
			ShowDebug: false,
			LogLevel:  "info",
		},
		SkillMarkers: struct {
			Enabled bool    `json:"enabled"`
			Opacity float32 `json:"opacity"`
			Size    float32 `json:"size"`
			Skills  map[string]struct {
				Color   string `json:"color"`
				Enabled bool   `json:"enabled"`
			} `json:"skills"`
		}{
			Enabled: true,
			Opacity: 0.5,
			Size:    4.8,
			Skills: map[string]struct {
				Color   string `json:"color"`
				Enabled bool   `json:"enabled"`
			}{
				"stone_wall":  {Color: "0x7F008888", Enabled: true},
				"fire_wall":   {Color: "0x7F880000", Enabled: true},
				"sanctuary":   {Color: "0x7F00FFFF", Enabled: true},
				"safety_wall": {Color: "0x7F880088", Enabled: true},
				"fire_pillar": {Color: "0x7F888800", Enabled: true},
				"magnus_ex":   {Color: "0x7F00FFFF", Enabled: true},
				"pneuma":      {Color: "0x7F00FFFF", Enabled: true},
				"ice_wall":    {Color: "0x7F880088", Enabled: true},
				"quagmire":    {Color: "0x7F448844", Enabled: true},
				"warp_portal": {Color: "0x7FFFFFFF", Enabled: true},
			},
		},
		RangeDisplay: struct {
			Enabled   bool    `json:"enabled"`
			Color     string  `json:"color"`
			Opacity   float32 `json:"opacity"`
			Thickness float32 `json:"thickness"`
			Segments  int     `json:"segments"`
		}{
			Enabled:   true,
			Color:     "0x7F00FF00", // Green with 50% opacity
			Opacity:   0.5,
			Thickness: 2.0,
			Segments:  32,
		},
	}
}

// LoadConfig loads configuration from file
func LoadConfig(filename string) (*Config, error) {
	// Check if config file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// Create default config file
		config := DefaultConfig()
		if err := config.SaveConfig(filename); err != nil {
			return nil, fmt.Errorf("error creating default config: %v", err)
		}
		return config, nil
	}

	// Read config file
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	// Parse JSON
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error parsing config file: %v", err)
	}

	return &config, nil
}

// SaveConfig saves configuration to file
func (c *Config) SaveConfig(filename string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("error creating config directory: %v", err)
	}

	// Marshal to JSON with indentation
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling config: %v", err)
	}

	// Write to file
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("error writing config file: %v", err)
	}

	return nil
}

// GetConfigPath returns the default config file path
func GetConfigPath() string {
	// Obtener el directorio donde se ejecuta el programa
	execDir, err := os.Getwd()
	if err != nil {
		return "dtap-config.json"
	}
	return filepath.Join(execDir, "dtap-config.json")
}
