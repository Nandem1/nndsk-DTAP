package visual

import (
	"context"
	"sync"
	"time"

	"dtap-go/pkg/logger"
)

// OverlayManager manages the visual overlay system
type OverlayManager struct {
	skillRenderer *SkillMarkerRenderer
	rangeDisplay  *RangeDisplay
	logger        *logger.Logger
	ctx           context.Context
	cancel        context.CancelFunc
	wg            sync.WaitGroup
	enabled       bool
	updateTicker  *time.Ticker
}

// NewOverlayManager creates a new overlay manager
func NewOverlayManager(config *SkillMarkerConfig, rangeConfig *RangeDisplayConfig, logger *logger.Logger) *OverlayManager {
	ctx, cancel := context.WithCancel(context.Background())

	return &OverlayManager{
		skillRenderer: NewSkillMarkerRenderer(config, logger),
		rangeDisplay:  NewRangeDisplay(rangeConfig, logger),
		logger:        logger,
		ctx:           ctx,
		cancel:        cancel,
		enabled:       config.Enabled,
	}
}

// Start starts the overlay manager
func (om *OverlayManager) Start() error {
	if om.enabled {
		om.logger.Info("Starting overlay manager")

		// Start the update loop
		om.wg.Add(1)
		go om.updateLoop()

		// Start range display
		if err := om.rangeDisplay.Start(); err != nil {
			om.logger.Error("Error starting range display: %v", err)
		}

		om.logger.Info("Overlay manager started successfully")
	} else {
		om.logger.Info("Overlay manager disabled")
	}

	return nil
}

// Stop stops the overlay manager
func (om *OverlayManager) Stop() {
	om.logger.Info("Stopping overlay manager")

	om.cancel()

	if om.updateTicker != nil {
		om.updateTicker.Stop()
	}

	// Stop range display
	om.rangeDisplay.Stop()

	om.wg.Wait()
	om.logger.Info("Overlay manager stopped")
}

// updateLoop runs the main update loop
func (om *OverlayManager) updateLoop() {
	defer om.wg.Done()

	om.updateTicker = time.NewTicker(100 * time.Millisecond) // 10 FPS
	defer om.updateTicker.Stop()

	for {
		select {
		case <-om.ctx.Done():
			return
		case <-om.updateTicker.C:
			om.update()
		}
	}
}

// update performs a single update cycle
func (om *OverlayManager) update() {
	if !om.enabled {
		return
	}

	// Update skill markers (clean up expired ones)
	om.skillRenderer.UpdateMarkers()

	// Update range display
	om.rangeDisplay.SimulateRangeDetection()

	// Here we would integrate with the game's rendering system
	// For now, we just log the current state
	markerCount := om.skillRenderer.GetMarkerCount()
	if markerCount > 0 {
		om.logger.Debug("Active skill markers", "count", markerCount)
	}
}

// AddSkillMarker adds a skill marker to the overlay
func (om *OverlayManager) AddSkillMarker(skillID int, x, y, z float32, duration time.Duration) {
	if !om.enabled {
		return
	}

	position := Position{X: x, Y: y, Z: z}
	om.skillRenderer.AddMarker(skillID, position, duration)
}

// RemoveSkillMarker removes a skill marker from the overlay
func (om *OverlayManager) RemoveSkillMarker(markerID string) {
	om.skillRenderer.RemoveMarker(markerID)
}

// SetEnabled enables or disables the overlay
func (om *OverlayManager) SetEnabled(enabled bool) {
	om.enabled = enabled
	om.skillRenderer.SetEnabled(enabled)
	om.logger.Info("Overlay manager", "enabled", enabled)
}

// IsEnabled returns whether the overlay is enabled
func (om *OverlayManager) IsEnabled() bool {
	return om.enabled && om.skillRenderer.IsEnabled()
}

// GetActiveMarkers returns all active skill markers
func (om *OverlayManager) GetActiveMarkers() []*SkillMarker {
	return om.skillRenderer.GetActiveMarkers()
}

// ClearAllMarkers clears all skill markers
func (om *OverlayManager) ClearAllMarkers() {
	om.skillRenderer.ClearAllMarkers()
}

// UpdateConfig updates the overlay configuration
func (om *OverlayManager) UpdateConfig(config *SkillMarkerConfig) {
	om.skillRenderer.UpdateConfig(config)
	om.enabled = config.Enabled
	om.logger.Info("Updated overlay configuration", "enabled", config.Enabled)
}

// SetRangeDisplayEnabled enables or disables the range display
func (om *OverlayManager) SetRangeDisplayEnabled(enabled bool) {
	om.rangeDisplay.SetEnabled(enabled)
}

// IsRangeDisplayEnabled returns whether the range display is enabled
func (om *OverlayManager) IsRangeDisplayEnabled() bool {
	return om.rangeDisplay.IsEnabled()
}

// SetCharacterRange sets the character's attack range
func (om *OverlayManager) SetCharacterRange(rangeValue float32) {
	om.rangeDisplay.SetRange(rangeValue)
}

// SetCharacterPosition sets the character's position
func (om *OverlayManager) SetCharacterPosition(x, y, z float32) {
	om.rangeDisplay.SetPosition(x, y, z)
}

// GetStats returns statistics about the overlay system
func (om *OverlayManager) GetStats() map[string]interface{} {
	stats := om.skillRenderer.GetStats()
	stats["overlay_enabled"] = om.enabled
	stats["update_interval"] = "100ms"
	stats["range_display"] = om.rangeDisplay.GetStats()
	return stats
}

// SimulateSkillDetection simulates detecting skills in the game world
// This would be replaced with actual game memory scanning
func (om *OverlayManager) SimulateSkillDetection() {
	// Simulate some common skills for testing
	testSkills := []struct {
		skillID  int
		x, y, z  float32
		duration time.Duration
	}{
		{0x7E, 100, 0, 100, 30 * time.Second}, // Stone Wall
		{0x7F, 105, 0, 105, 20 * time.Second}, // Fire Wall
		{0x83, 110, 0, 110, 60 * time.Second}, // Sanctuary
		{0x86, 115, 0, 115, 30 * time.Second}, // Safety Wall
	}

	for _, skill := range testSkills {
		om.AddSkillMarker(skill.skillID, skill.x, skill.y, skill.z, skill.duration)
	}

	om.logger.Info("Simulated skill detection", "skills_added", len(testSkills))
}
