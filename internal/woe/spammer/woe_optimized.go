package spammer

import (
	"context"
	"time"

	"dtap-go/internal/input"
	"dtap-go/pkg/config"
	"dtap-go/pkg/logger"

	"golang.org/x/sys/windows"
)

// WOEOptimizedSpammer provides ultra-fast spamming for competitive WoE
type WOEOptimizedSpammer struct {
	config       *config.Config
	keyDetector  *input.KeyDetector
	keySender    *input.KeySender
	mouseSender  *input.MouseSender
	windowHandle windows.Handle

	// Performance tracking
	spamCount int64

	// Context for cancellation
	ctx    context.Context
	cancel context.CancelFunc
}

// NewWOEOptimizedSpammer creates an optimized spammer for WoE
func NewWOEOptimizedSpammer(cfg *config.Config, windowHandle windows.Handle) *WOEOptimizedSpammer {
	ctx, cancel := context.WithCancel(context.Background())

	spammer := &WOEOptimizedSpammer{
		config:       cfg,
		keyDetector:  input.NewKeyDetector(),
		keySender:    input.NewKeySender(),
		mouseSender:  input.NewMouseSender(),
		windowHandle: windowHandle,
		ctx:          ctx,
		cancel:       cancel,
	}

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

	// Start optimized spammer loop (same as normal but faster)
	go s.runOptimizedSpammerLoop(ctx)
	go s.keepAliveLoop(ctx) // Keep-alive to prevent stopping
}

// runOptimizedSpammerLoop is the main optimized spammer loop (same as normal but faster)
func (s *WOEOptimizedSpammer) runOptimizedSpammerLoop(ctx context.Context) {
	// Use much faster ticker for WoE - 1ms instead of normal delay
	ticker := time.NewTicker(1 * time.Millisecond) // Ultra-fast for WoE
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Info("Spammer WoE detenido por contexto externo")
			return
		case <-s.ctx.Done():
			logger.Info("Spammer WoE detenido por contexto interno")
			return
		case <-ticker.C:
			s.processSpammerKeys()
		}
	}
}

// keepAliveLoop ensures the spammer stays active and provides status updates
func (s *WOEOptimizedSpammer) keepAliveLoop(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second) // Status every 5 seconds
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Info("Keep-alive loop detenido por contexto externo")
			return
		case <-s.ctx.Done():
			logger.Info("Keep-alive loop detenido por contexto interno")
			return
		case <-ticker.C:
			spamCount := s.spamCount
			logger.Info("Spammer WoE Status - Total spams: %d", spamCount)
		}
	}
}

// processSpammerKeys processes all configured spammer keys (same as normal version)
func (s *WOEOptimizedSpammer) processSpammerKeys() {
	// Procesar teclas (ignorando modifier keys) - same logic as normal version
	for _, key := range s.config.Spammer.Keys {
		if !s.keyDetector.IsModifierKey(key) && s.keyDetector.IsKeyPressed(key) {
			s.sendKey(key)
			s.sendLeftClick()

			// Doble click si está habilitado
			if s.config.Spammer.DoubleClick {
				s.sendLeftClick()
			}

			// Skill cancel click if configured
			if s.config.Spammer.SkillCancel {
				s.sendSkillCancelClick()
			}

			s.spamCount++
		}
	}

	// Log activity every 1000 iterations to track performance
	if s.spamCount%1000 == 0 && s.spamCount > 0 {
		logger.Debug("Spammer WoE activo - Total spams: %d", s.spamCount)
	}
}

// sendKey sends a key to the game window (same as normal version)
func (s *WOEOptimizedSpammer) sendKey(key string) {
	err := s.keySender.SendKey(s.windowHandle, key)
	if err != nil {
		logger.Error("Error enviando tecla %s: %v", key, err)
	}
}

// sendLeftClick sends a left mouse click to the center of the game window (same as normal version)
func (s *WOEOptimizedSpammer) sendLeftClick() {
	// Click en el centro de la ventana (coordenadas típicas)
	err := s.mouseSender.SendLeftClick(s.windowHandle, 400, 300)
	if err != nil {
		logger.Error("Error enviando click izquierdo: %v", err)
	}
}

// sendSkillCancelClick sends mouse click for skill canceling
func (s *WOEOptimizedSpammer) sendSkillCancelClick() {
	// Send click to center of window for skill canceling
	err := s.mouseSender.SendLeftClick(s.windowHandle, 400, 300)
	if err != nil {
		logger.Error("Error enviando skill cancel click: %v", err)
	}
}

// GetStats returns spammer performance statistics
func (s *WOEOptimizedSpammer) GetStats() (spamCount int64, burstMode bool, activeKeys int) {
	return s.spamCount, false, 0 // Simplified for WoE
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
