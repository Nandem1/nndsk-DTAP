package spammer

import (
	"context"
	"sync"
	"time"

	"dtap-go/internal/input"
	"dtap-go/pkg/config"
	"dtap-go/pkg/logger"

	"golang.org/x/sys/windows"
)

// WOEOptimizedSpammer provides ultra-fast spamming for competitive WoE
type WOEOptimizedSpammer struct {
	config           *config.Config
	ultraFastSender  *input.UltraFastSender
	keyDetector      *input.KeyDetector
	windowHandle     windows.Handle
	
	// Performance tracking
	spamCount        int64
	lastSpamTime     time.Time
	activeKeys       map[string]bool
	activeKeysMutex  sync.RWMutex
	
	// Smart spamming features
	burstMode        bool
	burstStartTime   time.Time
	skillCancelMode  bool
	comboMode        bool
	
	// Context for cancellation
	ctx              context.Context
	cancel           context.CancelFunc
}

// NewWOEOptimizedSpammer creates an optimized spammer for WoE
func NewWOEOptimizedSpammer(cfg *config.Config, windowHandle windows.Handle) *WOEOptimizedSpammer {
	ctx, cancel := context.WithCancel(context.Background())
	
	spammer := &WOEOptimizedSpammer{
		config:          cfg,
		ultraFastSender: input.NewUltraFastSender(),
		keyDetector:     input.NewKeyDetector(),
		windowHandle:    windowHandle,
		activeKeys:      make(map[string]bool),
		ctx:             ctx,
		cancel:          cancel,
	}
	
	spammer.ultraFastSender.SetWindowHandle(windowHandle)
	
	return spammer
}

// Start begins the ultra-optimized spammer service
func (s *WOEOptimizedSpammer) Start(ctx context.Context) {
	if !s.config.Spammer.Enabled {
		logger.Info("Spammer WoE está deshabilitado")
		return
	}

	if len(s.config.Spammer.Keys) == 0 {
		logger.Warn("Spammer WoE habilitado pero no hay teclas configuradas")
		return
	}

	logger.Info("Iniciando Spammer WoE Optimizado - Modo Competitivo")
	logger.Info("Teclas configuradas: %v", s.config.Spammer.Keys)
	logger.Info("Delay: %dms", s.config.Spammer.Delay)
	
	// Start multiple optimization goroutines
	go s.ultraFastSpamLoop(ctx)     // 5ms loop for key detection
	go s.smartSpamLoop(ctx)         // Variable delay based on mode
	go s.burstModeLoop(ctx)         // Burst mode management
}

// ultraFastSpamLoop handles ultra-fast key detection and spamming (5ms intervals)
func (s *WOEOptimizedSpammer) ultraFastSpamLoop(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Millisecond) // 200 FPS
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			s.detectAndSpamKeys()
		}
	}
}

// smartSpamLoop handles intelligent spamming based on game state
func (s *WOEOptimizedSpammer) smartSpamLoop(ctx context.Context) {
	// Dynamic delay based on mode
	delay := time.Duration(s.config.Spammer.Delay) * time.Millisecond
	if s.burstMode {
		delay = delay / 2 // Half delay in burst mode
	}
	
	ticker := time.NewTicker(delay)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			s.processSmartSpam()
			
			// Adjust delay dynamically
			newDelay := time.Duration(s.config.Spammer.Delay) * time.Millisecond
			if s.burstMode {
				newDelay = newDelay / 2
			}
			ticker.Reset(newDelay)
		}
	}
}

// burstModeLoop manages burst mode activation
func (s *WOEOptimizedSpammer) burstModeLoop(ctx context.Context) {
	ticker := time.NewTicker(100 * time.Millisecond) // Check every 100ms
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			s.manageBurstMode()
		}
	}
}

// detectAndSpamKeys detects pressed keys and spams them immediately
func (s *WOEOptimizedSpammer) detectAndSpamKeys() {
	s.activeKeysMutex.Lock()
	defer s.activeKeysMutex.Unlock()
	
	for _, key := range s.config.Spammer.Keys {
		if !s.keyDetector.IsModifierKey(key) && s.keyDetector.IsKeyPressed(key) {
			// Key is pressed, add to active keys
			if !s.activeKeys[key] {
				s.activeKeys[key] = true
				logger.Debug("Tecla activada: %s", key)
			}
			
			// Spam the key immediately
			s.spamKey(key)
		} else {
			// Key is not pressed, remove from active keys
			if s.activeKeys[key] {
				delete(s.activeKeys, key)
				logger.Debug("Tecla desactivada: %s", key)
			}
		}
	}
}

// processSmartSpam processes intelligent spamming
func (s *WOEOptimizedSpammer) processSmartSpam() {
	s.activeKeysMutex.RLock()
	activeKeys := make([]string, 0, len(s.activeKeys))
	for key := range s.activeKeys {
		activeKeys = append(activeKeys, key)
	}
	s.activeKeysMutex.RUnlock()
	
	// Process all active keys
	for _, key := range activeKeys {
		s.spamKey(key)
		
		// Send mouse click if configured
		if s.config.Spammer.SkillCancel {
			s.sendSkillCancelClick()
		}
	}
}

// spamKey sends a key with ultra-fast input
func (s *WOEOptimizedSpammer) spamKey(key string) {
	// Prevent spam flooding
	if time.Since(s.lastSpamTime) < 10*time.Millisecond {
		return
	}
	
	s.ultraFastSender.SendKeyUltraFast(key)
	s.spamCount++
	s.lastSpamTime = time.Now()
	
	// Enable burst mode if spamming rapidly
	if s.spamCount%10 == 0 && time.Since(s.burstStartTime) < 1*time.Second {
		s.burstMode = true
		s.burstStartTime = time.Now()
	}
}

// sendSkillCancelClick sends mouse click for skill canceling
func (s *WOEOptimizedSpammer) sendSkillCancelClick() {
	// Send click to center of window for skill canceling
	s.ultraFastSender.SendMouseClickUltraFast(400, 300)
	
	// Double click if enabled
	if s.config.Spammer.DoubleClick {
		time.Sleep(50 * time.Microsecond) // Ultra-minimal delay
		s.ultraFastSender.SendMouseClickUltraFast(400, 300)
	}
}

// manageBurstMode manages burst mode activation and deactivation
func (s *WOEOptimizedSpammer) manageBurstMode() {
	if s.burstMode {
		// Exit burst mode after 3 seconds
		if time.Since(s.burstStartTime) > 3*time.Second {
			s.burstMode = false
			logger.Debug("Burst mode desactivado")
		}
	}
}

// EnableBurstMode manually enables burst mode
func (s *WOEOptimizedSpammer) EnableBurstMode() {
	s.burstMode = true
	s.burstStartTime = time.Now()
	logger.Info("Burst mode activado manualmente")
}

// DisableBurstMode manually disables burst mode
func (s *WOEOptimizedSpammer) DisableBurstMode() {
	s.burstMode = false
	logger.Info("Burst mode desactivado manualmente")
}

// ToggleSkillCancelMode toggles skill cancel mode
func (s *WOEOptimizedSpammer) ToggleSkillCancelMode() {
	s.skillCancelMode = !s.skillCancelMode
	logger.Info("Skill cancel mode: %v", s.skillCancelMode)
}

// ToggleComboMode toggles combo mode for advanced sequences
func (s *WOEOptimizedSpammer) ToggleComboMode() {
	s.comboMode = !s.comboMode
	logger.Info("Combo mode: %v", s.comboMode)
}

// GetActiveKeys returns currently active keys
func (s *WOEOptimizedSpammer) GetActiveKeys() []string {
	s.activeKeysMutex.RLock()
	defer s.activeKeysMutex.RUnlock()
	
	keys := make([]string, 0, len(s.activeKeys))
	for key := range s.activeKeys {
		keys = append(keys, key)
	}
	
	return keys
}

// GetStats returns spammer performance statistics
func (s *WOEOptimizedSpammer) GetStats() (spamCount int64, burstMode bool, activeKeys int) {
	s.activeKeysMutex.RLock()
	defer s.activeKeysMutex.RUnlock()
	
	return s.spamCount, s.burstMode, len(s.activeKeys)
}

// UpdateConfig updates the spammer configuration
func (s *WOEOptimizedSpammer) UpdateConfig(cfg *config.Config) {
	s.config = cfg
	logger.Info("Configuración del Spammer WoE actualizada")
}

// Stop stops the spammer service
func (s *WOEOptimizedSpammer) Stop() {
	s.cancel()
	logger.Info("Spammer WoE detenido - Total spams: %d", s.spamCount)
}