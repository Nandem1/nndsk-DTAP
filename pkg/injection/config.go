package injection

import "fmt"

// InjectionConfig represents configuration for the injection system
type InjectionConfig struct {
	Enabled          bool   `json:"enabled"`
	DLLPath          string `json:"dll_path"`
	SharedMemoryName string `json:"shared_memory_name"`
	HookType         int    `json:"hook_type"`
	AutoInject       bool   `json:"auto_inject"`
	EnableM2E        bool   `json:"enable_m2e"`
	EnableBBE        bool   `json:"enable_bbe"`
	GroundZBias      int32  `json:"ground_z_bias"`
	AlphaLevel       int32  `json:"alpha_level"`
	EnableDeadCell   bool   `json:"enable_dead_cell"`
	EnableChatScope  bool   `json:"enable_chat_scope"`
	FixWindowMode    bool   `json:"fix_window_mode"`
	ShowFrameRate    bool   `json:"show_frame_rate"`
	EnableObjectInfo bool   `json:"enable_object_info"`
	EnableAudioMode  bool   `json:"enable_audio_mode"`
	CPUCoolerLevel   int32  `json:"cpu_cooler_level"`

	// Gepard bypass settings
	EnableGepardBypass bool `json:"enable_gepard_bypass"`
	AutoDetectGepard   bool `json:"auto_detect_gepard"`
	BypassCRC          bool `json:"bypass_crc"`
	RestoreFunctions   bool `json:"restore_functions"`
}

// DefaultInjectionConfig returns the default injection configuration
func DefaultInjectionConfig() *InjectionConfig {
	return &InjectionConfig{
		Enabled:          true,
		DLLPath:          "dtap-hook-native.dll",
		SharedMemoryName: "DTAPSharedMemory",
		HookType:         5, // WH_CBT
		AutoInject:       true,
		EnableM2E:        true,
		EnableBBE:        true,
		GroundZBias:      0,
		AlphaLevel:       0,
		EnableDeadCell:   true,
		EnableChatScope:  true,
		FixWindowMode:    true,
		ShowFrameRate:    true,
		EnableObjectInfo: false,
		EnableAudioMode:  false,
		CPUCoolerLevel:   0,

		// Gepard bypass defaults
		EnableGepardBypass: true,
		AutoDetectGepard:   true,
		BypassCRC:          true,
		RestoreFunctions:   true,
	}
}

// Validate validates the injection configuration
func (ic *InjectionConfig) Validate() error {
	if ic.DLLPath == "" {
		return fmt.Errorf("DLL path cannot be empty")
	}

	if ic.SharedMemoryName == "" {
		return fmt.Errorf("shared memory name cannot be empty")
	}

	if ic.HookType < 0 || ic.HookType > 15 {
		return fmt.Errorf("invalid hook type: %d", ic.HookType)
	}

	return nil
}
