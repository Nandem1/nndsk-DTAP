package visual

import (
	"fmt"
	"sync"
	"time"

	"dtap-go/pkg/logger"
)

// SkillMarker represents a skill marker on the ground
type SkillMarker struct {
	SkillID   int           `json:"skill_id"`
	Position  Position      `json:"position"`
	Color     string        `json:"color"`
	StartTime time.Time     `json:"start_time"`
	Duration  time.Duration `json:"duration"`
	Size      float32       `json:"size"`
}

// Position represents a 3D position in the game world
type Position struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
	Z float32 `json:"z"`
}

// SkillMarkerRenderer handles rendering of skill markers
type SkillMarkerRenderer struct {
	config     *SkillMarkerConfig
	markers    map[string]*SkillMarker
	markersMux sync.RWMutex
	logger     *logger.Logger
	enabled    bool
}

// NewSkillMarkerRenderer creates a new skill marker renderer
func NewSkillMarkerRenderer(config *SkillMarkerConfig, logger *logger.Logger) *SkillMarkerRenderer {
	return &SkillMarkerRenderer{
		config:  config,
		markers: make(map[string]*SkillMarker),
		logger:  logger,
		enabled: config.Enabled,
	}
}

// AddMarker adds a new skill marker
func (r *SkillMarkerRenderer) AddMarker(skillID int, position Position, duration time.Duration) {
	if !r.config.IsSkillEnabled(skillID) {
		return
	}

	skillConfig, exists := r.config.GetSkillConfig(skillID)
	if !exists {
		return
	}

	markerID := fmt.Sprintf("%d_%.2f_%.2f_%.2f", skillID, position.X, position.Y, position.Z)

	r.markersMux.Lock()
	defer r.markersMux.Unlock()

	marker := &SkillMarker{
		SkillID:   skillID,
		Position:  position,
		Color:     skillConfig.Color,
		StartTime: time.Now(),
		Duration:  duration,
		Size:      r.config.Size,
	}

	r.markers[markerID] = marker

	skillName, _ := SkillIDMapping[skillID]
	r.logger.Debug("Added skill marker",
		"skill_id", skillID,
		"skill_name", skillName,
		"position", fmt.Sprintf("(%.2f, %.2f, %.2f)", position.X, position.Y, position.Z),
		"duration", duration,
		"color", skillConfig.Color)
}

// RemoveMarker removes a skill marker
func (r *SkillMarkerRenderer) RemoveMarker(markerID string) {
	r.markersMux.Lock()
	defer r.markersMux.Unlock()

	if marker, exists := r.markers[markerID]; exists {
		skillName, _ := SkillIDMapping[marker.SkillID]
		r.logger.Debug("Removed skill marker",
			"skill_id", marker.SkillID,
			"skill_name", skillName,
			"marker_id", markerID)
		delete(r.markers, markerID)
	}
}

// UpdateMarkers updates and cleans up expired markers
func (r *SkillMarkerRenderer) UpdateMarkers() {
	r.markersMux.Lock()
	defer r.markersMux.Unlock()

	now := time.Now()
	var expiredMarkers []string

	for markerID, marker := range r.markers {
		if now.Sub(marker.StartTime) > marker.Duration {
			expiredMarkers = append(expiredMarkers, markerID)
		}
	}

	for _, markerID := range expiredMarkers {
		if marker, exists := r.markers[markerID]; exists {
			skillName, _ := SkillIDMapping[marker.SkillID]
			r.logger.Debug("Expired skill marker",
				"skill_id", marker.SkillID,
				"skill_name", skillName,
				"marker_id", markerID)
			delete(r.markers, markerID)
		}
	}
}

// GetActiveMarkers returns all active markers
func (r *SkillMarkerRenderer) GetActiveMarkers() []*SkillMarker {
	r.markersMux.RLock()
	defer r.markersMux.RUnlock()

	markers := make([]*SkillMarker, 0, len(r.markers))
	for _, marker := range r.markers {
		markers = append(markers, marker)
	}

	return markers
}

// GetMarkerCount returns the number of active markers
func (r *SkillMarkerRenderer) GetMarkerCount() int {
	r.markersMux.RLock()
	defer r.markersMux.RUnlock()

	return len(r.markers)
}

// SetEnabled enables or disables the skill marker renderer
func (r *SkillMarkerRenderer) SetEnabled(enabled bool) {
	r.enabled = enabled
	r.logger.Info("Skill marker renderer", "enabled", enabled)
}

// IsEnabled returns whether the renderer is enabled
func (r *SkillMarkerRenderer) IsEnabled() bool {
	return r.enabled && r.config.Enabled
}

// UpdateConfig updates the configuration
func (r *SkillMarkerRenderer) UpdateConfig(config *SkillMarkerConfig) {
	r.config = config
	r.logger.Info("Updated skill marker configuration", "enabled", config.Enabled)
}

// ClearAllMarkers removes all active markers
func (r *SkillMarkerRenderer) ClearAllMarkers() {
	r.markersMux.Lock()
	defer r.markersMux.Unlock()

	count := len(r.markers)
	r.markers = make(map[string]*SkillMarker)
	r.logger.Info("Cleared all skill markers", "count", count)
}

// GetStats returns statistics about the skill marker renderer
func (r *SkillMarkerRenderer) GetStats() map[string]interface{} {
	r.markersMux.RLock()
	defer r.markersMux.RUnlock()

	skillCounts := make(map[string]int)
	for _, marker := range r.markers {
		skillName, _ := SkillIDMapping[marker.SkillID]
		skillCounts[skillName]++
	}

	return map[string]interface{}{
		"enabled":        r.enabled,
		"total_markers":  len(r.markers),
		"skill_counts":   skillCounts,
		"config_enabled": r.config.Enabled,
		"opacity":        r.config.Opacity,
		"size":           r.config.Size,
	}
}
