package autopot

import (
	"context"
	"syscall"
	"time"

	"dtap-gui/internal/memory"
	"dtap-gui/pkg/config"
	"dtap-gui/pkg/constants"
	"dtap-gui/pkg/logger"

	"golang.org/x/sys/windows"
)

// UltraFastAutoPot - Versión ultra optimizada del autopot con predicción
type UltraFastAutoPot struct {
	config       *config.Config
	memClient    *memory.Client
	windowHandle windows.Handle

	// Direct Windows API - sin abstracciones
	user32      *syscall.LazyDLL
	postMessage *syscall.LazyProc

	// Pre-allocated variables para evitar GC
	hpPotKeyCode uint32
	spPotKeyCode uint32
	loopCount    int64
	lastLogTime  time.Time

	// Performance tracking
	lastHPPercent int
	lastSPPercent int
	potionsUsed   int64

	// Sistema de predicción para anticipar el pot 200ms antes
	hpHistory      []int     // Historial de HP para detectar tendencias
	spHistory      []int     // Historial de SP para detectar tendencias
	lastHPCheck    time.Time // Última vez que se revisó HP
	lastSPCheck    time.Time // Última vez que se revisó SP
	predictionMode bool      // Si está en modo predicción activo
	damageRate     float64   // Tasa de daño por segundo detectada
}

// NewUltraFastAutoPot crea el autopot más rápido posible
func NewUltraFastAutoPot(cfg *config.Config, memClient *memory.Client, windowHandle windows.Handle) *UltraFastAutoPot {
	user32 := syscall.NewLazyDLL("user32.dll")

	// Pre-convertir teclas a códigos para evitar conversión en runtime
	hpPotKeyCode := getKeyCodeUltraFast(cfg.AutoPot.HPPotKey)
	spPotKeyCode := getKeyCodeUltraFast(cfg.AutoPot.SPPotKey)

	return &UltraFastAutoPot{
		config:         cfg,
		memClient:      memClient,
		windowHandle:   windowHandle,
		user32:         user32,
		postMessage:    user32.NewProc("PostMessageW"),
		hpPotKeyCode:   hpPotKeyCode,
		spPotKeyCode:   spPotKeyCode,
		lastLogTime:    time.Now(),
		lastHPPercent:  100,
		lastSPPercent:  100,
		hpHistory:      make([]int, 0, 10), // Historial de 10 muestras
		spHistory:      make([]int, 0, 10), // Historial de 10 muestras
		lastHPCheck:    time.Now(),
		lastSPCheck:    time.Now(),
		predictionMode: false,
		damageRate:     0.0,
	}
}

// Start inicia el autopot ultra rápido
func (s *UltraFastAutoPot) Start(ctx context.Context) {
	if !s.config.AutoPot.Enabled {
		logger.Info("UltraFast AutoPot está deshabilitado")
		return
	}

	logger.Info("🚀 UltraFast AutoPot iniciado - Modo PREDICCIÓN 200ms")
	go s.runUltraFastLoop(ctx)
}

// runUltraFastLoop - Loop principal ultra optimizado
func (s *UltraFastAutoPot) runUltraFastLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			s.loopCount++

			// Procesar autopot ultra rápido
			s.processAutoPotUltraFast()

			// Log deshabilitado para evitar interferencia con display
			// if s.loopCount%50000 == 0 {
			//	now := time.Now()
			//	elapsed := now.Sub(startTime)
			//	logger.Info("⚡ UltraFast AutoPot loop #%d: %v elapsed, %d potions used",
			//		s.loopCount, elapsed, s.potionsUsed)
			//	s.lastLogTime = now
			// }

			// Delay mínimo para evitar saturar el juego
			// 9µs para alta estabilidad (mismo que el spammer)
			time.Sleep(9 * time.Microsecond)
		}
	}
}

// processAutoPotUltraFast - Procesamiento ultra rápido con predicción
func (s *UltraFastAutoPot) processAutoPotUltraFast() {
	now := time.Now()

	// Leer HP/SP de forma ultra rápida
	hp, maxHP, err := s.memClient.ReadHP()
	if err != nil {
		return // Silenciar errores para máxima velocidad
	}

	sp, maxSP, err := s.memClient.ReadSP()
	if err != nil {
		return // Silenciar errores para máxima velocidad
	}

	// Calcular porcentajes de forma ultra rápida
	hpPercent := int((float64(hp) / float64(maxHP)) * 100)
	spPercent := int((float64(sp) / float64(maxSP)) * 100)

	// Actualizar historial de HP para predicción
	s.updateHPHistory(hpPercent, now)
	s.updateSPHistory(spPercent, now)

	// Sistema de predicción: anticipar el pot 200ms antes
	predictedHP := s.predictHPIn200ms(hpPercent, now)
	predictedSP := s.predictSPIn200ms(spPercent, now)

	// Auto pot HP con predicción - ultra rápido
	if predictedHP < s.config.AutoPot.HPPercent || hpPercent < s.config.AutoPot.HPPercent {
		s.sendKeyUltraFast(s.hpPotKeyCode)
		s.potionsUsed++
		s.lastHPPercent = hpPercent
	}

	// Auto pot SP con predicción - ultra rápido
	if predictedSP < s.config.AutoPot.SPPercent || spPercent < s.config.AutoPot.SPPercent {
		s.sendKeyUltraFast(s.spPotKeyCode)
		s.potionsUsed++
		s.lastSPPercent = spPercent
	}
}

// sendKeyUltraFast - Envío de tecla ultra rápido (mismo que spammer)
func (s *UltraFastAutoPot) sendKeyUltraFast(keyCode uint32) {
	if keyCode == 0 {
		return
	}

	// KeyDown
	s.postMessage.Call(
		uintptr(s.windowHandle),
		constants.WM_KEYDOWN,
		uintptr(keyCode),
		0,
	)

	// Delay mínimo para que RO procese el input
	time.Sleep(9 * time.Microsecond)

	// KeyUp
	s.postMessage.Call(
		uintptr(s.windowHandle),
		constants.WM_KEYUP,
		uintptr(keyCode),
		0,
	)
}

// getKeyCodeUltraFast - Conversión ultra rápida de tecla a código
func getKeyCodeUltraFast(key string) uint32 {
	switch key {
	case "F1":
		return 0x70
	case "F2":
		return 0x71
	case "F3":
		return 0x72
	case "F4":
		return 0x73
	case "F5":
		return 0x74
	case "F6":
		return 0x75
	case "F7":
		return 0x76
	case "F8":
		return 0x77
	case "F9":
		return 0x78
	case "F10":
		return 0x79
	case "F11":
		return 0x7A
	case "F12":
		return 0x7B
	case "1":
		return 0x31
	case "2":
		return 0x32
	case "3":
		return 0x33
	case "4":
		return 0x34
	case "5":
		return 0x35
	case "6":
		return 0x36
	case "7":
		return 0x37
	case "8":
		return 0x38
	case "9":
		return 0x39
	case "0":
		return 0x30
	default:
		return 0
	}
}

// GetStats retorna estadísticas del autopot
func (s *UltraFastAutoPot) GetStats() (int64, int64, time.Time) {
	return s.loopCount, s.potionsUsed, s.lastLogTime
}

// updateHPHistory - Actualiza el historial de HP para predicción
func (s *UltraFastAutoPot) updateHPHistory(hpPercent int, now time.Time) {
	// Agregar al historial
	s.hpHistory = append(s.hpHistory, hpPercent)

	// Mantener solo los últimos 10 valores
	if len(s.hpHistory) > 10 {
		s.hpHistory = s.hpHistory[1:]
	}

	s.lastHPCheck = now
}

// updateSPHistory - Actualiza el historial de SP para predicción
func (s *UltraFastAutoPot) updateSPHistory(spPercent int, now time.Time) {
	// Agregar al historial
	s.spHistory = append(s.spHistory, spPercent)

	// Mantener solo los últimos 10 valores
	if len(s.spHistory) > 10 {
		s.spHistory = s.spHistory[1:]
	}

	s.lastSPCheck = now
}

// predictHPIn200ms - Predice el HP en 200ms basado en la tendencia
func (s *UltraFastAutoPot) predictHPIn200ms(currentHP int, now time.Time) int {
	if len(s.hpHistory) < 3 {
		return currentHP // No hay suficiente historial
	}

	// Calcular la tasa de cambio de HP por segundo
	timeDiff := now.Sub(s.lastHPCheck).Seconds()
	if timeDiff <= 0 {
		return currentHP
	}

	// Calcular tendencia de los últimos 3 valores
	recent := s.hpHistory[len(s.hpHistory)-3:]
	hpChange := float64(recent[0] - recent[2])     // Cambio en los últimos 3 checks
	hpChangePerSecond := hpChange / (timeDiff * 3) // Cambio por segundo

	// Si está perdiendo HP, predecir 200ms hacia adelante
	if hpChangePerSecond < 0 {
		predictedChange := hpChangePerSecond * 0.2 // 200ms = 0.2 segundos
		predictedHP := currentHP + int(predictedChange)

		// Asegurar que no sea negativo
		if predictedHP < 0 {
			predictedHP = 0
		}

		return predictedHP
	}

	return currentHP
}

// predictSPIn200ms - Predice el SP en 200ms basado en la tendencia
func (s *UltraFastAutoPot) predictSPIn200ms(currentSP int, now time.Time) int {
	if len(s.spHistory) < 3 {
		return currentSP // No hay suficiente historial
	}

	// Calcular la tasa de cambio de SP por segundo
	timeDiff := now.Sub(s.lastSPCheck).Seconds()
	if timeDiff <= 0 {
		return currentSP
	}

	// Calcular tendencia de los últimos 3 valores
	recent := s.spHistory[len(s.spHistory)-3:]
	spChange := float64(recent[0] - recent[2])     // Cambio en los últimos 3 checks
	spChangePerSecond := spChange / (timeDiff * 3) // Cambio por segundo

	// Si está perdiendo SP, predecir 200ms hacia adelante
	if spChangePerSecond < 0 {
		predictedChange := spChangePerSecond * 0.2 // 200ms = 0.2 segundos
		predictedSP := currentSP + int(predictedChange)

		// Asegurar que no sea negativo
		if predictedSP < 0 {
			predictedSP = 0
		}

		return predictedSP
	}

	return currentSP
}

// UpdateConfig actualiza la configuración
func (s *UltraFastAutoPot) UpdateConfig(cfg *config.Config) {
	s.config = cfg
	s.hpPotKeyCode = getKeyCodeUltraFast(cfg.AutoPot.HPPotKey)
	s.spPotKeyCode = getKeyCodeUltraFast(cfg.AutoPot.SPPotKey)
	logger.Info("Configuración del UltraFast AutoPot actualizada")
}
