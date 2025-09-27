package spammer

import (
	"context"
	"syscall"
	"time"

	"dtap-gui/internal/metrics"
	"dtap-gui/pkg/config"
	"dtap-gui/pkg/constants"
	"dtap-gui/pkg/logger"

	"golang.org/x/sys/windows"
)

// Spammer - Implementación basada en lógica C# con performance Go
type Spammer struct {
	config       *config.Config
	windowHandle windows.Handle

	// APIs directas (como C#)
	user32           *syscall.LazyDLL
	postMessage      *syscall.LazyProc
	getAsyncKeyState *syscall.LazyProc

	// Estado simple (como C#)
	keyCodes []uint32
	delay    time.Duration
	running  bool

	// Métricas
	loopCount   int64
	lastLogTime time.Time
}

// NewSpammer crea un nuevo spammer basado en lógica C#
func NewSpammer(cfg *config.Config, windowHandle windows.Handle) *Spammer {
	user32 := syscall.NewLazyDLL("user32.dll")

	// Pre-convertir teclas a códigos
	keyCodes := make([]uint32, 0, len(cfg.Spammer.Keys))
	for _, key := range cfg.Spammer.Keys {
		keyCodes = append(keyCodes, getKeyCode(key))
	}

	// Calcular delay (como C#) con validación estricta
	delay := time.Duration(cfg.Spammer.Delay) * time.Millisecond
	if delay <= 0 {
		delay = constants.SpammerMinDelay // Default como C#
		logger.Warn("⚠️ Spammer delay configurado como 0, usando 1ms por defecto")
	}

	return &Spammer{
		config:           cfg,
		windowHandle:     windowHandle,
		user32:           user32,
		postMessage:      user32.NewProc("PostMessageW"),
		getAsyncKeyState: user32.NewProc("GetAsyncKeyState"),
		keyCodes:         keyCodes,
		delay:            delay,
		lastLogTime:      time.Now(),
	}
}

// Start inicia el spammer con lógica C#
func (s *Spammer) Start(ctx context.Context) {
	if !s.config.Spammer.Enabled {
		logger.Info("Spammer está deshabilitado")
		return
	}

	if len(s.config.Spammer.Keys) == 0 {
		logger.Warn("Spammer habilitado pero no hay teclas configuradas")
		return
	}

	logger.Info("🚀 INICIANDO SPAMMER (Lógica C# + Performance Go)")
	logger.Info("Teclas: %v", s.config.Spammer.Keys)
	logger.Info("Delay: %v (configurado: %vms)", s.delay, s.config.Spammer.Delay)
	logger.Info("Mouse Click: %v", s.config.Spammer.MouseClick)

	s.running = true

	// Goroutine principal (como C# pero no bloqueante)
	go s.runSpammerLoop(ctx)
}

// runSpammerLoop - Loop principal con lógica C#
func (s *Spammer) runSpammerLoop(ctx context.Context) {
	for s.running {
		select {
		case <-ctx.Done():
			logger.Info("Spammer detenido")
			s.running = false
			return
		default:
			s.processKeys() // Lógica exacta del C#
		}
	}
}

// ProcessSpammer - Método para el scheduler (compatible con pausar/reanudar)
func (s *Spammer) ProcessSpammer(ctx context.Context) error {
	// Lógica exacta del C#: una tecla a la vez, timing síncrono
	for _, keyCode := range s.keyCodes {
		if s.isKeyPressed(keyCode) {
			// Debug: Log cada vez que se detecta tecla presionada
			logger.Debug("🎯 Spammer: Tecla detectada, enviando...")

			// 1. Enviar tecla (como C#)
			s.sendKeyDirect(keyCode)

			// 2. Click inmediato después (como C#)
			if s.config.Spammer.MouseClick {
				s.sendClickDirect()
			}

			// 3. Delay fijo (como C#) - CRÍTICO: Aplicar delay configurado
			time.Sleep(s.delay)

			// 4. Métricas
			s.loopCount++
			metrics.GetMetrics().IncrementSpammerLoops()

			// Debug: Log después del delay
			logger.Debug("✅ Spammer: Tecla enviada, delay aplicado: %v", s.delay)

			return nil // Una tecla a la vez, como C#
		}
	}
	return nil
}

// processKeys - Lógica exacta del C# (estabilidad garantizada) - DEPRECATED
func (s *Spammer) processKeys() {
	// Lógica exacta del C#: una tecla a la vez, timing síncrono
	for _, keyCode := range s.keyCodes {
		if s.isKeyPressed(keyCode) {
			// 1. Enviar tecla (como C#)
			s.sendKeyDirect(keyCode)

			// 2. Click inmediato después (como C#)
			if s.config.Spammer.MouseClick {
				s.sendClickDirect()
			}

			// 3. Delay fijo (como C#)
			time.Sleep(s.delay)

			// 4. Métricas
			s.loopCount++
			metrics.GetMetrics().IncrementSpammerLoops()

			return // Una tecla a la vez, como C#
		}
	}
}

// sendKeyDirect - Envío directo de tecla (como C#)
func (s *Spammer) sendKeyDirect(keyCode uint32) {
	// PostMessage directo, como C#
	s.postMessage.Call(uintptr(s.windowHandle), uintptr(constants.WM_KEYDOWN), uintptr(keyCode), 0)
	s.postMessage.Call(uintptr(s.windowHandle), uintptr(constants.WM_KEYUP), uintptr(keyCode), 0)

	// Reportar métricas
	metrics.GetMetrics().IncrementSpammerEvents()
}

// sendClickDirect - Click directo (como C#)
func (s *Spammer) sendClickDirect() {
	// Click directo, como C#
	s.postMessage.Call(uintptr(s.windowHandle), uintptr(constants.WM_LBUTTONDOWN), 0, 0)
	time.Sleep(constants.SpammerMinDelay) // Delay mínimo como C#
	s.postMessage.Call(uintptr(s.windowHandle), uintptr(constants.WM_LBUTTONUP), 0, 0)

	// Reportar métricas
	metrics.GetMetrics().IncrementSpammerEvents()
}

// isKeyPressed - Verificación directa (como C#)
func (s *Spammer) isKeyPressed(keyCode uint32) bool {
	// Verificación directa, como C#
	state, _, _ := s.getAsyncKeyState.Call(uintptr(keyCode))
	return state&0x8000 != 0
}

// Stop detiene el spammer
func (s *Spammer) Stop() {
	s.running = false
}

// GetStats retorna estadísticas
func (s *Spammer) GetStats() (int64, time.Time) {
	return s.loopCount, s.lastLogTime
}

// getKeyCode - Conversión de tecla a código (optimizada)
func getKeyCode(key string) uint32 {
	switch key {
	case "F1":
		return constants.KeyF1
	case "F2":
		return constants.KeyF2
	case "F3":
		return constants.KeyF3
	case "F4":
		return constants.KeyF4
	case "F5":
		return constants.KeyF5
	case "F6":
		return constants.KeyF6
	case "F7":
		return constants.KeyF7
	case "F8":
		return constants.KeyF8
	case "F9":
		return constants.KeyF9
	case "F10":
		return constants.KeyF10
	case "F11":
		return constants.KeyF11
	case "F12":
		return constants.KeyF12
	case "1":
		return constants.Key1
	case "2":
		return constants.Key2
	case "3":
		return constants.Key3
	case "4":
		return constants.Key4
	case "5":
		return constants.Key5
	case "6":
		return constants.Key6
	case "7":
		return constants.Key7
	case "8":
		return constants.Key8
	case "9":
		return constants.Key9
	case "0":
		return constants.Key0
	case "A":
		return constants.KeyA
	case "B":
		return constants.KeyB
	case "C":
		return constants.KeyC
	case "D":
		return constants.KeyD
	case "E":
		return constants.KeyE
	case "F":
		return constants.KeyF
	case "G":
		return constants.KeyG
	case "H":
		return constants.KeyH
	case "I":
		return constants.KeyI
	case "J":
		return constants.KeyJ
	case "K":
		return constants.KeyK
	case "L":
		return constants.KeyL
	case "M":
		return constants.KeyM
	case "N":
		return constants.KeyN
	case "O":
		return constants.KeyO
	case "P":
		return constants.KeyP
	case "Q":
		return constants.KeyQ
	case "R":
		return constants.KeyR
	case "S":
		return constants.KeyS
	case "T":
		return constants.KeyT
	case "U":
		return constants.KeyU
	case "V":
		return constants.KeyV
	case "W":
		return constants.KeyW
	case "X":
		return constants.KeyX
	case "Y":
		return constants.KeyY
	case "Z":
		return constants.KeyZ
	default:
		return 0
	}
}
