package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// WOEConfig represents optimized configuration for War of Emperium
type WOEConfig struct {
	// Ultra-fast autopot settings
	AutoPot struct {
		HPPotKey       string `json:"hp_pot_key"`
		SPPotKey       string `json:"sp_pot_key"`
		HPPercent      int    `json:"hp_percent"`
		SPPercent      int    `json:"sp_percent"`
		EmergencyHP    int    `json:"emergency_hp_percent"`
		EmergencySP    int    `json:"emergency_sp_percent"`
		CheckDelay     int    `json:"check_delay_ms"`
		EmergencyDelay int    `json:"emergency_delay_ms"`
		Enabled        bool   `json:"enabled"`
		SmartMode      bool   `json:"smart_mode"`
		PredictiveMode bool   `json:"predictive_mode"`
	} `json:"autopot"`

	// Ultra-fast spammer settings
	Spammer struct {
		Enabled       bool     `json:"enabled"`
		Keys          []string `json:"keys"`
		Delay         int      `json:"delay_ms"`
		BurstDelay    int      `json:"burst_delay_ms"`
		SkillCancel   bool     `json:"skill_cancel"`
		DoubleClick   bool     `json:"double_click"`
		BurstMode     bool     `json:"burst_mode"`
		ComboMode     bool     `json:"combo_mode"`
		UltraFastMode bool     `json:"ultra_fast_mode"`
	} `json:"spammer"`

	// Performance optimization settings
	Performance struct {
		TargetFPS         int  `json:"target_fps"`
		AdaptiveMode      bool `json:"adaptive_mode"`
		MemoryCache       bool `json:"memory_cache"`
		InputOptimization bool `json:"input_optimization"`
		BackgroundRefresh bool `json:"background_refresh"`
	} `json:"performance"`

	// WoE specific settings
	WOE struct {
		CombatMode     bool `json:"combat_mode"`
		DefenseMode    bool `json:"defense_mode"`
		OffenseMode    bool `json:"offense_mode"`
		AutoBuff       bool `json:"auto_buff"`
		BuffDelay      int  `json:"buff_delay_ms"`
		EnemyDetection bool `json:"enemy_detection"`
		SkillRotation  bool `json:"skill_rotation"`
	} `json:"woe"`

	// Monitoring settings
	Monitoring struct {
		Enabled        bool `json:"enabled"`
		ReportInterval int  `json:"report_interval_seconds"`
		LogPerformance bool `json:"log_performance"`
		ShowStats      bool `json:"show_stats"`
	} `json:"monitoring"`

	// UI settings
	UI struct {
		ShowDebug       bool   `json:"show_debug"`
		LogLevel        string `json:"log_level"`
		ShowPerformance bool   `json:"show_performance"`
		CompactMode     bool   `json:"compact_mode"`
	} `json:"ui"`

	// Skill markers for WoE tactical advantage
	SkillMarkers struct {
		Enabled        bool    `json:"enabled"`
		Opacity        float32 `json:"opacity"`
		Size           float32 `json:"size"`
		WoEOptimized   bool    `json:"woe_optimized"`
		ShowEnemyOnly  bool    `json:"show_enemy_only"`
		ShowAlliedOnly bool    `json:"show_allied_only"`
		Skills         map[string]struct {
			Color   string `json:"color"`
			Enabled bool   `json:"enabled"`
		} `json:"skills"`
	} `json:"skill_markers"`

	// Range display for WoE tactical advantage
	RangeDisplay struct {
		Enabled          bool    `json:"enabled"`
		Color            string  `json:"color"`
		Opacity          float32 `json:"opacity"`
		Thickness        float32 `json:"thickness"`
		Segments         int     `json:"segments"`
		WoEOptimized     bool    `json:"woe_optimized"`
		ShowOnlyInCombat bool    `json:"show_only_in_combat"`
	} `json:"range_display"`
}

// DefaultWOEConfig returns the default WoE-optimized configuration
func DefaultWOEConfig() *WOEConfig {
	return &WOEConfig{
		AutoPot: struct {
			HPPotKey       string `json:"hp_pot_key"`
			SPPotKey       string `json:"sp_pot_key"`
			HPPercent      int    `json:"hp_percent"`
			SPPercent      int    `json:"sp_percent"`
			EmergencyHP    int    `json:"emergency_hp_percent"`
			EmergencySP    int    `json:"emergency_sp_percent"`
			CheckDelay     int    `json:"check_delay_ms"`
			EmergencyDelay int    `json:"emergency_delay_ms"`
			Enabled        bool   `json:"enabled"`
			SmartMode      bool   `json:"smart_mode"`
			PredictiveMode bool   `json:"predictive_mode"`
		}{
			HPPotKey:       "F1",
			SPPotKey:       "F2",
			HPPercent:      60, // Higher threshold for WoE
			SPPercent:      40, // Higher threshold for WoE
			EmergencyHP:    20, // Emergency threshold
			EmergencySP:    15, // Emergency threshold
			CheckDelay:     50, // Faster checking
			EmergencyDelay: 10, // Ultra-fast emergency
			Enabled:        true,
			SmartMode:      true,
			PredictiveMode: true,
		},
		Spammer: struct {
			Enabled       bool     `json:"enabled"`
			Keys          []string `json:"keys"`
			Delay         int      `json:"delay_ms"`
			BurstDelay    int      `json:"burst_delay_ms"`
			SkillCancel   bool     `json:"skill_cancel"`
			DoubleClick   bool     `json:"double_click"`
			BurstMode     bool     `json:"burst_mode"`
			ComboMode     bool     `json:"combo_mode"`
			UltraFastMode bool     `json:"ultra_fast_mode"`
		}{
			Enabled:       false,
			Keys:          []string{},
			Delay:         25, // Ultra-fast for WoE
			BurstDelay:    10, // Even faster in burst mode
			SkillCancel:   true,
			DoubleClick:   true,
			BurstMode:     true,
			ComboMode:     false,
			UltraFastMode: true,
		},
		Performance: struct {
			TargetFPS         int  `json:"target_fps"`
			AdaptiveMode      bool `json:"adaptive_mode"`
			MemoryCache       bool `json:"memory_cache"`
			InputOptimization bool `json:"input_optimization"`
			BackgroundRefresh bool `json:"background_refresh"`
		}{
			TargetFPS:         60,
			AdaptiveMode:      true,
			MemoryCache:       true,
			InputOptimization: true,
			BackgroundRefresh: true,
		},
		WOE: struct {
			CombatMode     bool `json:"combat_mode"`
			DefenseMode    bool `json:"defense_mode"`
			OffenseMode    bool `json:"offense_mode"`
			AutoBuff       bool `json:"auto_buff"`
			BuffDelay      int  `json:"buff_delay_ms"`
			EnemyDetection bool `json:"enemy_detection"`
			SkillRotation  bool `json:"skill_rotation"`
		}{
			CombatMode:     true,
			DefenseMode:    false,
			OffenseMode:    true,
			AutoBuff:       false,
			BuffDelay:      5000,
			EnemyDetection: false,
			SkillRotation:  false,
		},
		Monitoring: struct {
			Enabled        bool `json:"enabled"`
			ReportInterval int  `json:"report_interval_seconds"`
			LogPerformance bool `json:"log_performance"`
			ShowStats      bool `json:"show_stats"`
		}{
			Enabled:        true,
			ReportInterval: 30,
			LogPerformance: true,
			ShowStats:      true,
		},
		UI: struct {
			ShowDebug       bool   `json:"show_debug"`
			LogLevel        string `json:"log_level"`
			ShowPerformance bool   `json:"show_performance"`
			CompactMode     bool   `json:"compact_mode"`
		}{
			ShowDebug:       false,
			LogLevel:        "info",
			ShowPerformance: true,
			CompactMode:     true,
		},
		SkillMarkers: struct {
			Enabled        bool    `json:"enabled"`
			Opacity        float32 `json:"opacity"`
			Size           float32 `json:"size"`
			WoEOptimized   bool    `json:"woe_optimized"`
			ShowEnemyOnly  bool    `json:"show_enemy_only"`
			ShowAlliedOnly bool    `json:"show_allied_only"`
			Skills         map[string]struct {
				Color   string `json:"color"`
				Enabled bool   `json:"enabled"`
			} `json:"skills"`
		}{
			Enabled:        true,
			Opacity:        0.7, // Más opaco para WoE
			Size:           4.8,
			WoEOptimized:   true,
			ShowEnemyOnly:  false,
			ShowAlliedOnly: false,
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
			Enabled          bool    `json:"enabled"`
			Color            string  `json:"color"`
			Opacity          float32 `json:"opacity"`
			Thickness        float32 `json:"thickness"`
			Segments         int     `json:"segments"`
			WoEOptimized     bool    `json:"woe_optimized"`
			ShowOnlyInCombat bool    `json:"show_only_in_combat"`
		}{
			Enabled:          true,
			Color:            "0x7F00FF00", // Green with 50% opacity
			Opacity:          0.7,          // More opaque for WoE
			Thickness:        3.0,          // Thicker for better visibility
			Segments:         64,           // Smoother circle
			WoEOptimized:     true,
			ShowOnlyInCombat: false,
		},
	}
}

// LoadWOEConfig loads WoE configuration from file
func LoadWOEConfig(filename string) (*WOEConfig, error) {
	// Check if config file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// Create default WoE config file
		config := DefaultWOEConfig()
		if err := config.SaveWOEConfig(filename); err != nil {
			return nil, fmt.Errorf("error creating default WoE config: %v", err)
		}
		return config, nil
	}

	// Read config file
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading WoE config file: %v", err)
	}

	// Parse JSON
	var config WOEConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error parsing WoE config file: %v", err)
	}

	return &config, nil
}

// SaveWOEConfig saves WoE configuration to file
func (c *WOEConfig) SaveWOEConfig(filename string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("error creating config directory: %v", err)
	}

	// Marshal to JSON with indentation
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling WoE config: %v", err)
	}

	// Write to file
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("error writing WoE config file: %v", err)
	}

	return nil
}

// GetWOEConfigPath returns the default WoE config file path
func GetWOEConfigPath() string {
	// Get the directory where the program is executed
	execDir, err := os.Getwd()
	if err != nil {
		return "dtap-woe-config.json"
	}
	return filepath.Join(execDir, "dtap-woe-config.json")
}

// ConvertToStandardConfig converts WoE config to standard config for compatibility
func (c *WOEConfig) ConvertToStandardConfig() *Config {
	return &Config{
		AutoPot: struct {
			HPPotKey   string `json:"hp_pot_key"`
			SPPotKey   string `json:"sp_pot_key"`
			HPPercent  int    `json:"hp_percent"`
			SPPercent  int    `json:"sp_percent"`
			CheckDelay int    `json:"check_delay_ms"`
			Enabled    bool   `json:"enabled"`
		}{
			HPPotKey:   c.AutoPot.HPPotKey,
			SPPotKey:   c.AutoPot.SPPotKey,
			HPPercent:  c.AutoPot.HPPercent,
			SPPercent:  c.AutoPot.SPPercent,
			CheckDelay: c.AutoPot.CheckDelay,
			Enabled:    c.AutoPot.Enabled,
		},
		AutoBuff: struct {
			Enabled bool `json:"enabled"`
			Delay   int  `json:"delay_ms"`
		}{
			Enabled: c.WOE.AutoBuff,
			Delay:   c.WOE.BuffDelay,
		},
		Spammer: struct {
			Enabled     bool     `json:"enabled"`
			Keys        []string `json:"keys"`
			Delay       int      `json:"delay_ms"`
			SkillCancel bool     `json:"skill_cancel"`
			DoubleClick bool     `json:"double_click"`
		}{
			Enabled:     c.Spammer.Enabled,
			Keys:        c.Spammer.Keys,
			Delay:       c.Spammer.Delay,
			SkillCancel: c.Spammer.SkillCancel,
			DoubleClick: c.Spammer.DoubleClick,
		},
		UI: struct {
			ShowDebug bool   `json:"show_debug"`
			LogLevel  string `json:"log_level"`
		}{
			ShowDebug: c.UI.ShowDebug,
			LogLevel:  c.UI.LogLevel,
		},
	}
}
