package visual

import (
	"dtap-go/pkg/config"
)

// ConfigAdapter converts between different config types and SkillMarkerConfig
type ConfigAdapter struct{}

// NewConfigAdapter creates a new config adapter
func NewConfigAdapter() *ConfigAdapter {
	return &ConfigAdapter{}
}

// ConvertFromNormalConfig converts a normal config to SkillMarkerConfig
func (ca *ConfigAdapter) ConvertFromNormalConfig(normalConfig *config.Config) *SkillMarkerConfig {
	skillConfig := NewSkillMarkerConfig()

	skillConfig.Enabled = normalConfig.SkillMarkers.Enabled
	skillConfig.Opacity = normalConfig.SkillMarkers.Opacity
	skillConfig.Size = normalConfig.SkillMarkers.Size

	// Convert skills map
	skillConfig.Skills = make(map[string]SkillConfig)
	for skillName, skillData := range normalConfig.SkillMarkers.Skills {
		skillConfig.Skills[skillName] = SkillConfig{
			Color:   skillData.Color,
			Enabled: skillData.Enabled,
		}
	}

	return skillConfig
}

// ConvertRangeDisplayFromNormalConfig converts normal config range display to RangeDisplayConfig
func (ca *ConfigAdapter) ConvertRangeDisplayFromNormalConfig(normalConfig *config.Config) *RangeDisplayConfig {
	return &RangeDisplayConfig{
		Enabled:   normalConfig.RangeDisplay.Enabled,
		Color:     normalConfig.RangeDisplay.Color,
		Opacity:   normalConfig.RangeDisplay.Opacity,
		Thickness: normalConfig.RangeDisplay.Thickness,
		Segments:  normalConfig.RangeDisplay.Segments,
	}
}

// ConvertFromWOEConfig converts a WoE config to SkillMarkerConfig
func (ca *ConfigAdapter) ConvertFromWOEConfig(woeConfig *config.WOEConfig) *SkillMarkerConfig {
	skillConfig := NewSkillMarkerConfig()

	skillConfig.Enabled = woeConfig.SkillMarkers.Enabled
	skillConfig.Opacity = woeConfig.SkillMarkers.Opacity
	skillConfig.Size = woeConfig.SkillMarkers.Size

	// Convert skills map
	skillConfig.Skills = make(map[string]SkillConfig)
	for skillName, skillData := range woeConfig.SkillMarkers.Skills {
		skillConfig.Skills[skillName] = SkillConfig{
			Color:   skillData.Color,
			Enabled: skillData.Enabled,
		}
	}

	return skillConfig
}

// ConvertRangeDisplayFromWOEConfig converts WoE config range display to RangeDisplayConfig
func (ca *ConfigAdapter) ConvertRangeDisplayFromWOEConfig(woeConfig *config.WOEConfig) *RangeDisplayConfig {
	return &RangeDisplayConfig{
		Enabled:   woeConfig.RangeDisplay.Enabled,
		Color:     woeConfig.RangeDisplay.Color,
		Opacity:   woeConfig.RangeDisplay.Opacity,
		Thickness: woeConfig.RangeDisplay.Thickness,
		Segments:  woeConfig.RangeDisplay.Segments,
	}
}

// UpdateNormalConfig updates a normal config with skill marker settings
func (ca *ConfigAdapter) UpdateNormalConfig(normalConfig *config.Config, skillConfig *SkillMarkerConfig) {
	normalConfig.SkillMarkers.Enabled = skillConfig.Enabled
	normalConfig.SkillMarkers.Opacity = skillConfig.Opacity
	normalConfig.SkillMarkers.Size = skillConfig.Size

	// Update skills map
	if normalConfig.SkillMarkers.Skills == nil {
		normalConfig.SkillMarkers.Skills = make(map[string]struct {
			Color   string `json:"color"`
			Enabled bool   `json:"enabled"`
		})
	}

	for skillName, skillData := range skillConfig.Skills {
		normalConfig.SkillMarkers.Skills[skillName] = struct {
			Color   string `json:"color"`
			Enabled bool   `json:"enabled"`
		}{
			Color:   skillData.Color,
			Enabled: skillData.Enabled,
		}
	}
}

// UpdateWOEConfig updates a WoE config with skill marker settings
func (ca *ConfigAdapter) UpdateWOEConfig(woeConfig *config.WOEConfig, skillConfig *SkillMarkerConfig) {
	woeConfig.SkillMarkers.Enabled = skillConfig.Enabled
	woeConfig.SkillMarkers.Opacity = skillConfig.Opacity
	woeConfig.SkillMarkers.Size = skillConfig.Size

	// Update skills map
	if woeConfig.SkillMarkers.Skills == nil {
		woeConfig.SkillMarkers.Skills = make(map[string]struct {
			Color   string `json:"color"`
			Enabled bool   `json:"enabled"`
		})
	}

	for skillName, skillData := range skillConfig.Skills {
		woeConfig.SkillMarkers.Skills[skillName] = struct {
			Color   string `json:"color"`
			Enabled bool   `json:"enabled"`
		}{
			Color:   skillData.Color,
			Enabled: skillData.Enabled,
		}
	}
}
