package spammer

import (
	"context"
	"time"

	"dtap-go/internal/input"
	"dtap-go/pkg/config"
	"dtap-go/pkg/logger"

	"golang.org/x/sys/windows"
)

// SpammerService handles key spamming functionality
type SpammerService struct {
	config       *config.Config
	keyDetector  *input.KeyDetector
	keySender    *input.KeySender
	mouseSender  *input.MouseSender
	windowHandle windows.Handle
}

// NewSpammerService creates a new spammer service
func NewSpammerService(cfg *config.Config, windowHandle windows.Handle) *SpammerService {
	return &SpammerService{
		config:       cfg,
		keyDetector:  input.NewKeyDetector(),
		keySender:    input.NewKeySender(),
		mouseSender:  input.NewMouseSender(),
		windowHandle: windowHandle,
	}
}

// Start begins the spammer service
func (s *SpammerService) Start(ctx context.Context) {
	if !s.config.Spammer.Enabled {
		logger.Info("Spammer está deshabilitado")
		return
	}

	if len(s.config.Spammer.Keys) == 0 {
		logger.Warn("Spammer habilitado pero no hay teclas configuradas")
		return
	}

	logger.Info("Iniciando spammer con %d teclas: %v", len(s.config.Spammer.Keys), s.config.Spammer.Keys)
	logger.Info("Delay del spammer: %dms", s.config.Spammer.Delay)

	go s.runSpammerLoop(ctx)
}

// runSpammerLoop is the main spammer loop
func (s *SpammerService) runSpammerLoop(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(s.config.Spammer.Delay) * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Info("Spammer detenido")
			return
		case <-ticker.C:
			s.processSpammerKeys()
		}
	}
}

// processSpammerKeys processes all configured spammer keys
func (s *SpammerService) processSpammerKeys() {
	// Procesar teclas (ignorando modifier keys)
	for _, key := range s.config.Spammer.Keys {
		if !s.keyDetector.IsModifierKey(key) && s.keyDetector.IsKeyPressed(key) {
			s.sendKey(key)
			s.sendLeftClick()

			// Doble click si está habilitado
			if s.config.Spammer.DoubleClick {
				s.sendLeftClick()
			}
		}
	}
}

// sendKey sends a key to the game window
func (s *SpammerService) sendKey(key string) {
	err := s.keySender.SendKey(s.windowHandle, key)
	if err != nil {
		logger.Error("Error enviando tecla %s: %v", key, err)
	}
}

// sendLeftClick sends a left mouse click to the center of the game window
func (s *SpammerService) sendLeftClick() {
	// Click en el centro de la ventana (coordenadas típicas)
	err := s.mouseSender.SendLeftClick(s.windowHandle, 400, 300)
	if err != nil {
		logger.Error("Error enviando click izquierdo: %v", err)
	}
}

// UpdateConfig updates the spammer configuration
func (s *SpammerService) UpdateConfig(cfg *config.Config) {
	s.config = cfg
	logger.Info("Configuración del spammer actualizada")
}
