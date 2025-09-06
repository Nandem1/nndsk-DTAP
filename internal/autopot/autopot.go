package autopot

import (
	"context"
	"time"

	"dtap-go/internal/input"
	"dtap-go/internal/memory"
	"dtap-go/pkg/config"
	"dtap-go/pkg/logger"

	"golang.org/x/sys/windows"
)

// AutoPotService handles automatic potion usage
type AutoPotService struct {
	config       *config.Config
	memClient    *memory.Client
	keySender    *input.KeySender
	windowHandle windows.Handle
}

// NewAutoPotService creates a new autopot service
func NewAutoPotService(cfg *config.Config, memClient *memory.Client, windowHandle windows.Handle) *AutoPotService {
	return &AutoPotService{
		config:       cfg,
		memClient:    memClient,
		keySender:    input.NewKeySender(),
		windowHandle: windowHandle,
	}
}

// Start begins the autopot service
func (s *AutoPotService) Start(ctx context.Context) {
	if !s.config.AutoPot.Enabled {
		logger.Info("AutoPot está deshabilitado")
		return
	}

	// AutoPot iniciado silenciosamente

	go s.runAutoPotLoop(ctx)
}

// runAutoPotLoop is the main autopot loop
func (s *AutoPotService) runAutoPotLoop(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(s.config.AutoPot.CheckDelay) * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.checkAndUsePotions()
		}
	}
}

// checkAndUsePotions checks HP/SP and uses potions if needed
func (s *AutoPotService) checkAndUsePotions() {
	// Leer HP/SP
	hp, maxHP, err := s.memClient.ReadHP()
	if err != nil {
		logger.Error("Error leyendo HP: %v", err)
		return
	}

	sp, maxSP, err := s.memClient.ReadSP()
	if err != nil {
		logger.Error("Error leyendo SP: %v", err)
		return
	}

	// Calcular porcentajes
	hpPercent := int((float64(hp) / float64(maxHP)) * 100)
	spPercent := int((float64(sp) / float64(maxSP)) * 100)

	// Solo mostrar estado cuando hay cambios significativos o se usan potions
	// El estado se muestra en el main loop cada segundo

	// Auto pot HP
	if hpPercent < s.config.AutoPot.HPPercent {
		s.keySender.SendKey(s.windowHandle, s.config.AutoPot.HPPotKey)
	}

	// Auto pot SP
	if spPercent < s.config.AutoPot.SPPercent {
		s.keySender.SendKey(s.windowHandle, s.config.AutoPot.SPPotKey)
	}
}

// UpdateConfig updates the autopot configuration
func (s *AutoPotService) UpdateConfig(cfg *config.Config) {
	s.config = cfg
	logger.Info("Configuración del AutoPot actualizada")
}
