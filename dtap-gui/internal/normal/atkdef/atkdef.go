package atkdef

import (
	"context"
	"fmt"
	"sync"
	"time"

	"dtap-gui/internal/input"
	"dtap-gui/pkg/config"
	"dtap-gui/pkg/constants"
	"dtap-gui/pkg/logger"

	"golang.org/x/sys/windows"
)

// ATKDEFType representa el tipo de modo (ATK o DEF)
type ATKDEFType int

const (
	ATK ATKDEFType = iota
	DEF
)

// ATKDEFMode implementa el modo ATK-DEF optimizado para Go
type ATKDEFMode struct {
	config       *config.Config
	windowHandle windows.Handle

	// Estados
	mu          sync.RWMutex
	isActive    bool
	isRunning   bool
	currentMode ATKDEFType

	// Timing optimizado
	switchDelay time.Duration
	spamDelay   time.Duration
	clickDelay  time.Duration

	// Estadísticas
	stats struct {
		TotalActivations int64
		TotalSwitches    int64
		LastActivation   time.Time
		LastSwitch       time.Time
	}

	// Control de goroutines
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// NewATKDEFMode crea una nueva instancia del modo ATK-DEF
func NewATKDEFMode(cfg *config.Config, windowHandle windows.Handle) *ATKDEFMode {
	ctx, cancel := context.WithCancel(context.Background())

	return &ATKDEFMode{
		config:       cfg,
		windowHandle: windowHandle,
		switchDelay:  time.Duration(cfg.ATKDEF.SwitchDelay) * time.Millisecond,
		spamDelay:    time.Duration(cfg.ATKDEF.SpamDelay) * time.Millisecond,
		clickDelay:   constants.ATKDEFClickDelay,
		ctx:          ctx,
		cancel:       cancel,
	}
}

// Start inicia el modo ATK-DEF
func (a *ATKDEFMode) Start() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.isRunning {
		return nil // Ya está ejecutándose
	}

	if !a.config.ATKDEF.Enabled {
		return fmt.Errorf("ATK-DEF mode no está habilitado")
	}

	logger.Info("🛡️ Iniciando ATK-DEF Mode - Prioridad EQUIPO (100Hz)")
	a.isRunning = true

	// Iniciar goroutine dedicada para chequeo crítico
	go a.runCriticalLoop(a.ctx)

	return nil
}

// runCriticalLoop - Goroutine dedicada para chequeo ultra rápido de ATK-DEF
func (a *ATKDEFMode) runCriticalLoop(ctx context.Context) {
	// Ticker para ATK-DEF: 10ms (100 Hz) - crítico para equipo
	ticker := time.NewTicker(constants.ATKDEFCheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			a.processATKDEFCritical() // ← Chequeo ultra rápido
		}
	}
}

// processATKDEFCritical - Procesamiento crítico para ATK-DEF
func (a *ATKDEFMode) processATKDEFCritical() {
	// Verificar si está habilitado
	if !a.config.ATKDEF.Enabled {
		return
	}

	// Verificar si la tecla de activación está presionada
	activationKey := a.config.ATKDEF.ActivationKey
	if activationKey == "" {
		return
	}

	// Usar el dispatcher de input para verificar el estado de la tecla
	isPressed := input.IsKeyPressed(activationKey)

	a.mu.Lock()
	wasActive := a.isActive
	a.mu.Unlock()

	if isPressed && !wasActive {
		// Activar modo ATK-DEF
		logger.Info("🎯 ATK-DEF: Activando modo - tecla %s presionada", activationKey)
		a.activateMode()
	} else if !isPressed && wasActive {
		// Desactivar modo ATK-DEF
		logger.Info("🎯 ATK-DEF: Desactivando modo - tecla %s soltada", activationKey)
		a.deactivateMode()
	}
}

// Stop detiene el modo ATK-DEF
func (a *ATKDEFMode) Stop() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if !a.isRunning {
		return nil // Ya está detenido
	}

	logger.Info("🛑 Deteniendo ATK-DEF Mode")
	a.isRunning = false

	return nil
}

// activateMode activa el modo ATK-DEF
func (a *ATKDEFMode) activateMode() {
	a.mu.Lock()
	a.isActive = true
	a.stats.TotalActivations++
	a.stats.LastActivation = time.Now()
	a.mu.Unlock()

	logger.Info("🎯 ATK-DEF Mode activado - Activation Key: %s", a.config.ATKDEF.ActivationKey)

	// Ejecutar secuencia ATK-DEF
	go a.executeSequence()
}

// deactivateMode desactiva el modo ATK-DEF
func (a *ATKDEFMode) deactivateMode() {
	a.mu.Lock()
	a.isActive = false
	a.mu.Unlock()

	logger.Info("⏹️ ATK-DEF Mode desactivado")
}

// executeSequence ejecuta la secuencia completa ATK-DEF
func (a *ATKDEFMode) executeSequence() {
	defer func() {
		if r := recover(); r != nil {
			logger.Error("❌ Panic en secuencia ATK-DEF: %v", r)
		}
	}()

	// 1. Equipar items ATK
	a.equipATKItems()

	// 2. Ejecutar spammer mientras la tecla esté presionada
	a.executeSpammer()

	// 3. Equipar items DEF
	a.equipDEFItems()
}

// equipATKItems equipa los items de ataque
func (a *ATKDEFMode) equipATKItems() {
	atkKeys := a.config.ATKDEF.ATKKeys
	if len(atkKeys) == 0 {
		logger.Warn("⚠️ No hay teclas ATK configuradas")
		return
	}

	logger.Info("⚔️ Equipando items ATK...")

	for name, key := range atkKeys {
		if key == "" {
			continue
		}

		// Enviar tecla
		input.SendKeyByName(a.windowHandle, key)
		logger.Info("🔑 Tecla ATK enviada: %s (%s)", name, key)

		// Delay entre switches
		time.Sleep(a.switchDelay)
	}

	a.mu.Lock()
	a.stats.TotalSwitches++
	a.stats.LastSwitch = time.Now()
	a.mu.Unlock()
}

// equipDEFItems equipa los items de defensa
func (a *ATKDEFMode) equipDEFItems() {
	defKeys := a.config.ATKDEF.DEFKeys
	if len(defKeys) == 0 {
		logger.Warn("⚠️ No hay teclas DEF configuradas")
		return
	}

	logger.Info("🛡️ Equipando items DEF...")

	for name, key := range defKeys {
		if key == "" {
			continue
		}

		// Enviar tecla
		input.SendKeyByName(a.windowHandle, key)
		logger.Info("🔑 Tecla DEF enviada: %s (%s)", name, key)

		// Delay entre switches
		time.Sleep(a.switchDelay)
	}

	a.mu.Lock()
	a.stats.TotalSwitches++
	a.stats.LastSwitch = time.Now()
	a.mu.Unlock()
}

// executeSpammer ejecuta el spammer mientras la tecla esté presionada
func (a *ATKDEFMode) executeSpammer() {
	activationKey := a.config.ATKDEF.ActivationKey
	if activationKey == "" {
		return
	}

	logger.Info("🎯 Iniciando spammer ATK-DEF...")

	// Loop mientras la tecla esté presionada
	for {
		// Verificar si la tecla sigue presionada
		if !input.IsKeyPressed(activationKey) {
			break
		}

		// Enviar tecla de activación (spammer)
		input.SendKeyByName(a.windowHandle, activationKey)

		// Click del mouse si está habilitado
		if a.config.ATKDEF.MouseClick {
			input.SendMouseClickByName(a.windowHandle)
			time.Sleep(a.clickDelay)
		}

		// Delay del spammer - usar delay mínimo si está en 0
		if a.spamDelay <= 0 {
			time.Sleep(constants.ATKDEFClickDelay)
		} else {
			time.Sleep(a.spamDelay)
		}
	}

	logger.Info("⏹️ Spammer ATK-DEF finalizado")
}

// UpdateConfig actualiza la configuración del modo ATK-DEF
func (a *ATKDEFMode) UpdateConfig(cfg *config.Config) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.config = cfg
	a.switchDelay = time.Duration(cfg.ATKDEF.SwitchDelay) * time.Millisecond
	a.spamDelay = time.Duration(cfg.ATKDEF.SpamDelay) * time.Millisecond

	logger.Info("🔄 Configuración ATK-DEF actualizada")
}

// GetStats retorna las estadísticas del modo ATK-DEF
func (a *ATKDEFMode) GetStats() (int64, int64, time.Time, time.Time) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.stats.TotalActivations, a.stats.TotalSwitches,
		a.stats.LastActivation, a.stats.LastSwitch
}

// IsActive retorna si el modo ATK-DEF está activo
func (a *ATKDEFMode) IsActive() bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.isActive
}

// IsRunning retorna si el modo ATK-DEF está ejecutándose
func (a *ATKDEFMode) IsRunning() bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.isRunning
}

// ProcessATKDEF es el handler principal para el scheduler
func (a *ATKDEFMode) ProcessATKDEF(ctx context.Context) error {
	// Verificar si está habilitado
	if !a.config.ATKDEF.Enabled {
		return nil
	}

	// Verificar si la tecla de activación está presionada
	activationKey := a.config.ATKDEF.ActivationKey
	if activationKey == "" {
		return nil
	}

	// Usar el dispatcher de input para verificar el estado de la tecla
	isPressed := input.IsKeyPressed(activationKey)

	a.mu.Lock()
	wasActive := a.isActive
	a.mu.Unlock()

	// Debug logging
	if isPressed {
		logger.Debug("🎯 ATK-DEF: Tecla %s presionada, wasActive: %v", activationKey, wasActive)
	}

	if isPressed && !wasActive {
		// Activar modo ATK-DEF
		logger.Info("🎯 ATK-DEF: Activando modo - tecla %s presionada", activationKey)
		a.activateMode()
	} else if !isPressed && wasActive {
		// Desactivar modo ATK-DEF
		logger.Info("🎯 ATK-DEF: Desactivando modo - tecla %s soltada", activationKey)
		a.deactivateMode()
	}

	return nil
}
