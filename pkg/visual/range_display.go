package visual

import (
	"context"
	"sync"
	"time"

	"dtap-go/pkg/logger"
)

// RangeDisplay represents the character's attack range visualization
type RangeDisplay struct {
	enabled    bool
	rangeValue float32
	position   Position
	color      string
	opacity    float32
	thickness  float32
	segments   int
	logger     *logger.Logger
	ctx        context.Context
	cancel     context.CancelFunc
	wg         sync.WaitGroup
}

// RangeDisplayConfig represents configuration for range display
type RangeDisplayConfig struct {
	Enabled   bool    `json:"enabled"`
	Color     string  `json:"color"`     // Hex color in format "0xAARRGGBB"
	Opacity   float32 `json:"opacity"`   // 0.0 to 1.0
	Thickness float32 `json:"thickness"` // Line thickness
	Segments  int     `json:"segments"`  // Number of segments in the circle
}

// DefaultRangeDisplayConfig returns default range display configuration
func DefaultRangeDisplayConfig() *RangeDisplayConfig {
	return &RangeDisplayConfig{
		Enabled:   true,
		Color:     "0x7F00FF00", // Green with 50% opacity
		Opacity:   0.5,
		Thickness: 2.0,
		Segments:  32, // Smooth circle
	}
}

// NewRangeDisplay creates a new range display
func NewRangeDisplay(config *RangeDisplayConfig, logger *logger.Logger) *RangeDisplay {
	ctx, cancel := context.WithCancel(context.Background())

	return &RangeDisplay{
		enabled:    config.Enabled,
		rangeValue: 1.5, // Default range (1.5 cells)
		position:   Position{X: 0, Y: 0, Z: 0},
		color:      config.Color,
		opacity:    config.Opacity,
		thickness:  config.Thickness,
		segments:   config.Segments,
		logger:     logger,
		ctx:        ctx,
		cancel:     cancel,
	}
}

// Start starts the range display
func (rd *RangeDisplay) Start() error {
	if !rd.enabled {
		rd.logger.Info("Range display disabled")
		return nil
	}

	rd.logger.Info("Starting range display")

	// Start the update loop
	rd.wg.Add(1)
	go rd.updateLoop()

	rd.logger.Info("Range display started successfully")
	return nil
}

// Stop stops the range display
func (rd *RangeDisplay) Stop() {
	rd.logger.Info("Stopping range display")

	rd.cancel()
	rd.wg.Wait()

	rd.logger.Info("Range display stopped")
}

// updateLoop runs the main update loop
func (rd *RangeDisplay) updateLoop() {
	defer rd.wg.Done()

	ticker := time.NewTicker(50 * time.Millisecond) // 20 FPS for smooth circle
	defer ticker.Stop()

	for {
		select {
		case <-rd.ctx.Done():
			return
		case <-ticker.C:
			rd.update()
		}
	}
}

// update performs a single update cycle
func (rd *RangeDisplay) update() {
	if !rd.enabled {
		return
	}

	// Here we would integrate with the game's rendering system
	// For now, we just log the current state
	rd.logger.Debug("Range display active",
		"range", rd.rangeValue,
		"position", rd.position,
		"color", rd.color)
}

// SetRange sets the attack range value
func (rd *RangeDisplay) SetRange(rangeValue float32) {
	rd.rangeValue = rangeValue
	rd.logger.Debug("Range updated", "range", rangeValue)
}

// SetPosition sets the character position
func (rd *RangeDisplay) SetPosition(x, y, z float32) {
	rd.position = Position{X: x, Y: y, Z: z}
}

// SetEnabled enables or disables the range display
func (rd *RangeDisplay) SetEnabled(enabled bool) {
	rd.enabled = enabled
	rd.logger.Info("Range display", "enabled", enabled)
}

// IsEnabled returns whether the range display is enabled
func (rd *RangeDisplay) IsEnabled() bool {
	return rd.enabled
}

// UpdateConfig updates the range display configuration
func (rd *RangeDisplay) UpdateConfig(config *RangeDisplayConfig) {
	rd.enabled = config.Enabled
	rd.color = config.Color
	rd.opacity = config.Opacity
	rd.thickness = config.Thickness
	rd.segments = config.Segments
	rd.logger.Info("Updated range display configuration", "enabled", config.Enabled)
}

// GetStats returns statistics about the range display
func (rd *RangeDisplay) GetStats() map[string]interface{} {
	return map[string]interface{}{
		"enabled":   rd.enabled,
		"range":     rd.rangeValue,
		"position":  rd.position,
		"color":     rd.color,
		"opacity":   rd.opacity,
		"thickness": rd.thickness,
		"segments":  rd.segments,
	}
}

// SimulateRangeDetection simulates detecting character range (for testing)
func (rd *RangeDisplay) SimulateRangeDetection() {
	// Simulate different range values based on character class/weapon
	testRanges := []float32{1.5, 2.0, 3.0, 4.0, 5.0}

	// Cycle through different ranges for testing
	currentTime := time.Now().Unix()
	rangeIndex := int(currentTime) % len(testRanges)

	rd.SetRange(testRanges[rangeIndex])
	rd.SetPosition(100, 0, 100) // Simulate character position

	rd.logger.Debug("Simulated range detection",
		"range", rd.rangeValue,
		"position", rd.position)
}
