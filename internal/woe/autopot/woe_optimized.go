package autopot

import (
	"context"
	"sync"
	"time"

	"dtap-go/internal/input"
	"dtap-go/internal/memory"
	"dtap-go/pkg/config"
	"dtap-go/pkg/logger"

	"golang.org/x/sys/windows"
)

// WOEOptimizedService provides ultra-fast autopot for competitive WoE
type WOEOptimizedService struct {
	config           *config.Config
	memClient        *memory.OptimizedClient
	ultraFastSender  *input.UltraFastSender
	windowHandle     windows.Handle
	
	// Performance tracking
	lastHPPercent    int
	lastSPPercent    int
	potCount         int64
	lastPotTime      time.Time
	
	// Smart potion logic
	emergencyMode    bool
	lastDamageTime   time.Time
	damageHistory    []time.Time
	damageMutex      sync.Mutex
	
	// Context for cancellation
	ctx              context.Context
	cancel           context.CancelFunc
}

// NewWOEOptimizedService creates an optimized autopot service for WoE
func NewWOEOptimizedService(cfg *config.Config, memClient *memory.OptimizedClient, windowHandle windows.Handle) *WOEOptimizedService {
	ctx, cancel := context.WithCancel(context.Background())
	
	service := &WOEOptimizedService{
		config:          cfg,
		memClient:       memClient,
		ultraFastSender: input.NewUltraFastSender(),
		windowHandle:    windowHandle,
		ctx:             ctx,
		cancel:          cancel,
		damageHistory:   make([]time.Time, 0, 10), // Track last 10 damage events
	}
	
	service.ultraFastSender.SetWindowHandle(windowHandle)
	
	return service
}

// Start begins the ultra-optimized autopot service
func (s *WOEOptimizedService) Start(ctx context.Context) {
	if !s.config.AutoPot.Enabled {
		logger.Info("AutoPot WoE está deshabilitado")
		return
	}

	logger.Info("Iniciando AutoPot WoE Optimizado - Modo Competitivo")
	
	// Start multiple goroutines for different optimization levels
	go s.ultraFastPotLoop(ctx)    // 10ms loop for emergency
	go s.smartPotLoop(ctx)        // 50ms loop for normal operation
	go s.damageDetectionLoop(ctx) // 20ms loop for damage tracking
}

// ultraFastPotLoop handles emergency potion usage (10ms intervals)
func (s *WOEOptimizedService) ultraFastPotLoop(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Millisecond) // 100 FPS
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			s.checkEmergencyPots()
		}
	}
}

// smartPotLoop handles intelligent potion usage (50ms intervals)
func (s *WOEOptimizedService) smartPotLoop(ctx context.Context) {
	ticker := time.NewTicker(50 * time.Millisecond) // 20 FPS
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			s.checkSmartPots()
		}
	}
}

// damageDetectionLoop tracks damage patterns for predictive potion usage
func (s *WOEOptimizedService) damageDetectionLoop(ctx context.Context) {
	ticker := time.NewTicker(20 * time.Millisecond) // 50 FPS
	defer ticker.Stop()

	var lastHP uint32 = 0

	for {
		select {
		case <-ctx.Done():
			return
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			if hp, _, valid := s.memClient.GetHP(); valid {
				if lastHP > 0 && hp < lastHP {
					// Damage detected
					s.recordDamage()
				}
				lastHP = hp
			}
		}
	}
}

// checkEmergencyPots handles critical HP/SP situations
func (s *WOEOptimizedService) checkEmergencyPots() {
	hp, maxHP, valid := s.memClient.GetHP()
	if !valid {
		return
	}
	
	sp, maxSP, valid := s.memClient.GetSP()
	if !valid {
		return
	}
	
	hpPercent := int((float64(hp) / float64(maxHP)) * 100)
	spPercent := int((float64(sp) / float64(maxSP)) * 100)
	
	// Emergency HP (below 20% or in emergency mode)
	if hpPercent < 20 || (s.emergencyMode && hpPercent < s.config.AutoPot.HPPercent) {
		s.useHPPotion()
		s.emergencyMode = true
	}
	
	// Emergency SP (below 15% or critical for skills)
	if spPercent < 15 || (s.emergencyMode && spPercent < s.config.AutoPot.SPPercent) {
		s.useSPPotion()
	}
	
	// Exit emergency mode if HP is above 80%
	if hpPercent > 80 {
		s.emergencyMode = false
	}
}

// checkSmartPots handles intelligent potion usage based on damage patterns
func (s *WOEOptimizedService) checkSmartPots() {
	hp, maxHP, valid := s.memClient.GetHP()
	if !valid {
		return
	}
	
	sp, maxSP, valid := s.memClient.GetSP()
	if !valid {
		return
	}
	
	hpPercent := int((float64(hp) / float64(maxHP)) * 100)
	spPercent := int((float64(sp) / float64(maxSP)) * 100)
	
	// Smart HP potion based on damage rate
	if s.shouldUseHPPotion(hpPercent) {
		s.useHPPotion()
	}
	
	// Smart SP potion based on usage patterns
	if s.shouldUseSPPotion(spPercent) {
		s.useSPPotion()
	}
}

// shouldUseHPPotion determines if HP potion should be used based on damage patterns
func (s *WOEOptimizedService) shouldUseHPPotion(currentPercent int) bool {
	// Always use if below threshold
	if currentPercent < s.config.AutoPot.HPPercent {
		return true
	}
	
	// Predictive usage based on damage rate
	if s.isTakingDamage() && currentPercent < s.config.AutoPot.HPPercent+20 {
		return true
	}
	
	return false
}

// shouldUseSPPotion determines if SP potion should be used
func (s *WOEOptimizedService) shouldUseSPPotion(currentPercent int) bool {
	// Always use if below threshold
	if currentPercent < s.config.AutoPot.SPPercent {
		return true
	}
	
	// Predictive usage if actively using skills
	if s.isActivelyUsingSkills() && currentPercent < s.config.AutoPot.SPPercent+15 {
		return true
	}
	
	return false
}

// useHPPotion uses HP potion with ultra-fast input
func (s *WOEOptimizedService) useHPPotion() {
	if time.Since(s.lastPotTime) < 50*time.Millisecond {
		return // Prevent potion spam
	}
	
	s.ultraFastSender.SendKeyUltraFast(s.config.AutoPot.HPPotKey)
	s.potCount++
	s.lastPotTime = time.Now()
	
	logger.Debug("HP Pot usado - Total: %d", s.potCount)
}

// useSPPotion uses SP potion with ultra-fast input
func (s *WOEOptimizedService) useSPPotion() {
	if time.Since(s.lastPotTime) < 50*time.Millisecond {
		return // Prevent potion spam
	}
	
	s.ultraFastSender.SendKeyUltraFast(s.config.AutoPot.SPPotKey)
	s.potCount++
	s.lastPotTime = time.Now()
	
	logger.Debug("SP Pot usado - Total: %d", s.potCount)
}

// recordDamage records damage events for pattern analysis
func (s *WOEOptimizedService) recordDamage() {
	s.damageMutex.Lock()
	defer s.damageMutex.Unlock()
	
	now := time.Now()
	s.damageHistory = append(s.damageHistory, now)
	s.lastDamageTime = now
	
	// Keep only last 10 damage events
	if len(s.damageHistory) > 10 {
		s.damageHistory = s.damageHistory[1:]
	}
}

// isTakingDamage checks if player is currently taking damage
func (s *WOEOptimizedService) isTakingDamage() bool {
	s.damageMutex.Lock()
	defer s.damageMutex.Unlock()
	
	// Check if damage occurred in last 2 seconds
	return time.Since(s.lastDamageTime) < 2*time.Second
}

// isActivelyUsingSkills checks if player is actively using skills (placeholder)
func (s *WOEOptimizedService) isActivelyUsingSkills() bool {
	// This could be enhanced with SP usage pattern detection
	// For now, assume active if SP is being consumed
	sp, maxSP, valid := s.memClient.GetSP()
	if !valid {
		return false
	}
	
	spPercent := int((float64(sp) / float64(maxSP)) * 100)
	return spPercent < 70 // Assume active skill usage if SP below 70%
}

// GetStats returns performance statistics
func (s *WOEOptimizedService) GetStats() (potCount int64, emergencyMode bool, lastPotTime time.Time) {
	return s.potCount, s.emergencyMode, s.lastPotTime
}

// UpdateConfig updates the service configuration
func (s *WOEOptimizedService) UpdateConfig(cfg *config.Config) {
	s.config = cfg
	logger.Info("Configuración del AutoPot WoE actualizada")
}

// Stop stops the service
func (s *WOEOptimizedService) Stop() {
	s.cancel()
	logger.Info("AutoPot WoE detenido - Total potions: %d", s.potCount)
}