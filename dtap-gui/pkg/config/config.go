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
		Enabled  bool   `json:"enabled"`
		Delay    int    `json:"delay_ms"`
		GloomKey string `json:"gloom_key"`
		AspdKey  string `json:"aspd_key"`
	} `json:"autobuff"`

	Spammer struct {
		Enabled     bool     `json:"enabled"`
		Keys        []string `json:"keys"`
		Delay       float64  `json:"delay_ms"` // Cambiado a float64 para soportar decimales
		SkillCancel bool     `json:"skill_cancel"`
		DoubleClick bool     `json:"double_click"`
		MouseClick  bool     `json:"mouse_click"`
		Mode        string   `json:"mode"` // "simple", "adaptive", "ultra", "ultra_fast"
	} `json:"spammer"`

	UI struct {
		ShowDebug bool   `json:"show_debug"`
		LogLevel  string `json:"log_level"`
	} `json:"ui"`
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
			Enabled  bool   `json:"enabled"`
			Delay    int    `json:"delay_ms"`
			GloomKey string `json:"gloom_key"`
			AspdKey  string `json:"aspd_key"`
		}{
			Enabled:  false,
			Delay:    5000,
			GloomKey: "F3",
			AspdKey:  "F4",
		},
		Spammer: struct {
			Enabled     bool     `json:"enabled"`
			Keys        []string `json:"keys"`
			Delay       float64  `json:"delay_ms"`
			SkillCancel bool     `json:"skill_cancel"`
			DoubleClick bool     `json:"double_click"`
			MouseClick  bool     `json:"mouse_click"`
			Mode        string   `json:"mode"`
		}{
			Enabled:     false,
			Keys:        []string{},
			Delay:       1.0, // Default to 1ms
			SkillCancel: false,
			DoubleClick: false,
			MouseClick:  true,     // Enable mouse click by default
			Mode:        "simple", // Default to simple mode
		},
		UI: struct {
			ShowDebug bool   `json:"show_debug"`
			LogLevel  string `json:"log_level"`
		}{
			ShowDebug: false,
			LogLevel:  "info",
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
