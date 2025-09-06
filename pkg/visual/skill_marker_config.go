package visual

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// SkillMarkerConfig represents the configuration for skill markers
type SkillMarkerConfig struct {
	Enabled bool                   `json:"enabled"`
	Opacity float32                `json:"opacity"` // 0.0 to 1.0
	Size    float32                `json:"size"`    // Size of the marker in game units
	Skills  map[string]SkillConfig `json:"skills"`
}

// SkillConfig represents configuration for a specific skill
type SkillConfig struct {
	Color   string `json:"color"` // Hex color in format "0xAARRGGBB"
	Enabled bool   `json:"enabled"`
}

// Default skill configurations based on SimpleROHook
var DefaultSkillConfigs = map[string]SkillConfig{
	"stone_wall":  {Color: "0x7F008888", Enabled: true}, // SW - Grey
	"fire_wall":   {Color: "0x7F880000", Enabled: true}, // FW - Red
	"sanctuary":   {Color: "0x7F00FFFF", Enabled: true}, // Sanctuary - Cyan
	"safety_wall": {Color: "0x7F880088", Enabled: true}, // SG - Purple
	"fire_pillar": {Color: "0x7F888800", Enabled: true}, // FP - Yellow
	"magnus_ex":   {Color: "0x7F00FFFF", Enabled: true}, // ME - Cyan
	"pneuma":      {Color: "0x7F00FFFF", Enabled: true}, // Pneuma - Cyan
	"ice_wall":    {Color: "0x7F880088", Enabled: true}, // IW - Purple
	"quagmire":    {Color: "0x7F448844", Enabled: true}, // QM - Green
	"warp_portal": {Color: "0x7FFFFFFF", Enabled: true}, // Warp Portal - White
	"blessing":    {Color: "0x7F888888", Enabled: true}, // Blessing - Grey
}

// SkillIDMapping maps skill IDs to skill names
var SkillIDMapping = map[int]string{
	0x7E: "stone_wall",
	0x7F: "fire_wall",
	0x80: "warp_portal",
	0x81: "warp_portal",
	0x82: "blessing",
	0x83: "sanctuary",
	0x84: "magnus_ex",
	0x85: "pneuma",
	0x86: "safety_wall",
	0x87: "fire_pillar",
	0x88: "fire_pillar",
	0x89: "blessing",
	0x8A: "blessing",
	0x8B: "blessing",
	0x8C: "blessing",
	0x8D: "ice_wall",
	0x8E: "quagmire",
	0x8F: "blessing",
}

// NewSkillMarkerConfig creates a new skill marker configuration with defaults
func NewSkillMarkerConfig() *SkillMarkerConfig {
	return &SkillMarkerConfig{
		Enabled: true,
		Opacity: 0.5,
		Size:    4.8,
		Skills:  make(map[string]SkillConfig),
	}
}

// LoadSkillMarkerConfig loads skill marker configuration from file
func LoadSkillMarkerConfig(configPath string) (*SkillMarkerConfig, error) {
	config := NewSkillMarkerConfig()

	// Set default skills
	for name, skillConfig := range DefaultSkillConfigs {
		config.Skills[name] = skillConfig
	}

	// Try to load from file if it exists
	if _, err := os.Stat(configPath); err == nil {
		data, err := os.ReadFile(configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}

		if err := json.Unmarshal(data, config); err != nil {
			return nil, fmt.Errorf("failed to parse config file: %w", err)
		}
	}

	return config, nil
}

// SaveSkillMarkerConfig saves skill marker configuration to file
func SaveSkillMarkerConfig(config *SkillMarkerConfig, configPath string) error {
	// Ensure directory exists
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GetSkillConfig returns the configuration for a specific skill ID
func (c *SkillMarkerConfig) GetSkillConfig(skillID int) (SkillConfig, bool) {
	skillName, exists := SkillIDMapping[skillID]
	if !exists {
		return SkillConfig{}, false
	}

	skillConfig, exists := c.Skills[skillName]
	if !exists {
		return SkillConfig{}, false
	}

	return skillConfig, true
}

// IsSkillEnabled checks if a skill marker is enabled
func (c *SkillMarkerConfig) IsSkillEnabled(skillID int) bool {
	if !c.Enabled {
		return false
	}

	skillConfig, exists := c.GetSkillConfig(skillID)
	return exists && skillConfig.Enabled
}

// ParseColor parses a hex color string to RGBA values
func ParseColor(colorStr string) (r, g, b, a uint8, err error) {
	if len(colorStr) != 10 || colorStr[:2] != "0x" {
		return 0, 0, 0, 0, fmt.Errorf("invalid color format: %s", colorStr)
	}

	var color uint32
	if _, err := fmt.Sscanf(colorStr, "0x%08X", &color); err != nil {
		return 0, 0, 0, 0, fmt.Errorf("failed to parse color: %w", err)
	}

	a = uint8((color >> 24) & 0xFF)
	r = uint8((color >> 16) & 0xFF)
	g = uint8((color >> 8) & 0xFF)
	b = uint8(color & 0xFF)

	return r, g, b, a, nil
}
