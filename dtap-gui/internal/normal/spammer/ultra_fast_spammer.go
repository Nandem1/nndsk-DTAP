package spammer

import (
	"context"
	"runtime"
	"syscall"
	"time"

	"dtap-gui/pkg/config"
	"dtap-gui/pkg/constants"
	"dtap-gui/pkg/logger"

	"golang.org/x/sys/windows"
)

// UltraFastSpammer - Sin time.Sleep para delays ultra bajos
type UltraFastSpammer struct {
	config       *config.Config
	windowHandle windows.Handle

	// Direct Windows API
	user32           *syscall.LazyDLL
	postMessage      *syscall.LazyProc
	getAsyncKeyState *syscall.LazyProc

	// Pre-allocated variables
	keyCodes    []uint32
	loopCount   int64
	lastLogTime time.Time

	// Ultra-fast timing
	targetDelay time.Duration
	lastProcess time.Time
}

// NewUltraFastSpammer crea el spammer más rápido posible
func NewUltraFastSpammer(cfg *config.Config, windowHandle windows.Handle) *UltraFastSpammer {
	user32 := syscall.NewLazyDLL("user32.dll")

	// Pre-convertir teclas
	keyCodes := make([]uint32, 0, len(cfg.Spammer.Keys))
	for _, key := range cfg.Spammer.Keys {
		keyCodes = append(keyCodes, getKeyCodeUltraFast(key))
	}

	// Calcular delay objetivo - usar microsegundos para decimales
	targetDelay := time.Duration(cfg.Spammer.Delay*1000) * time.Microsecond
	if targetDelay <= 0 {
		targetDelay = 100 * time.Microsecond // 0.1ms por defecto
	}

	return &UltraFastSpammer{
		config:           cfg,
		windowHandle:     windowHandle,
		user32:           user32,
		postMessage:      user32.NewProc("PostMessageW"),
		getAsyncKeyState: user32.NewProc("GetAsyncKeyState"),
		keyCodes:         keyCodes,
		targetDelay:      targetDelay,
		lastLogTime:      time.Now(),
		lastProcess:      time.Now(),
	}
}

// Start inicia el spammer ultra rápido
func (s *UltraFastSpammer) Start(ctx context.Context) {
	if !s.config.Spammer.Enabled {
		logger.Info("UltraFastSpammer está deshabilitado")
		return
	}

	if len(s.config.Spammer.Keys) == 0 {
		logger.Warn("UltraFastSpammer habilitado pero no hay teclas configuradas")
		return
	}

	logger.Info("⚡ INICIANDO ULTRA FAST SPAMMER ⚡")
	logger.Info("Teclas: %v", s.config.Spammer.Keys)
	logger.Info("Delay objetivo: %v", s.targetDelay)

	// Configurar runtime para máximo rendimiento
	runtime.GOMAXPROCS(runtime.NumCPU())
	runtime.GC()

	go s.runUltraFastLoop(ctx)
}

// runUltraFastLoop - Loop ultra rápido con timing preciso
func (s *UltraFastSpammer) runUltraFastLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			logger.Info("UltraFastSpammer detenido")
			return
		default:
			// Procesar teclas - SIN timing para máxima velocidad
			s.processKeysUltraFast()
			s.loopCount++

			// Log deshabilitado para evitar interferencia con display
			// if s.loopCount%50000 == 0 {
			//	now = time.Now()
			//	elapsed = now.Sub(s.lastLogTime)
			//	logger.Info("⚡ UltraFast loop #%d: %v elapsed (ANIMATION CANCEL MODE)",
			//		s.loopCount, elapsed)
			//	s.lastLogTime = now
			// }

			// Delay mínimo para evitar saturar el juego
			// 9µs para alta estabilidad
			time.Sleep(9 * time.Microsecond)
		}
	}
}

// busyWait - Espera ultra rápida sin Sleep para delays < 1ms
func (s *UltraFastSpammer) busyWait(duration time.Duration) {
	if duration <= 0 {
		return
	}

	// Busy wait ultra optimizado - sin Gosched para máxima velocidad
	target := time.Now().Add(duration)
	for time.Now().Before(target) {
		// Sin yield para delays < 1ms - máxima velocidad
	}
}

// processKeysUltraFast - Procesamiento ultra optimizado
func (s *UltraFastSpammer) processKeysUltraFast() {
	// Verificar todas las teclas con early exit
	for _, keyCode := range s.keyCodes {
		if keyCode == 0 {
			continue
		}

		// Verificación ultra rápida
		state, _, _ := s.getAsyncKeyState.Call(uintptr(keyCode))
		if state&0x8000 != 0 { // Key is pressed
			// Log deshabilitado para evitar interferencia con display
			// if s.loopCount%1000 == 0 {
			//	logger.Info("🎯 Key pressed detected! Loop #%d, keyCode: %d", s.loopCount, keyCode)
			// }

			// Enviar tecla ultra rápido
			s.sendKeyUltraFast(keyCode)

			// Mouse click si está habilitado
			if s.config.Spammer.MouseClick {
				s.sendMouseClickUltraFast()
			}

			// Early exit - si encontramos una tecla presionada, no necesitamos verificar las demás
			return
		}
	}
}

// sendKeyUltraFast - Envío de tecla ultra optimizado
func (s *UltraFastSpammer) sendKeyUltraFast(keyCode uint32) {
	// PostMessage directo - sin verificaciones para máxima velocidad
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

// sendMouseClickUltraFast - Mouse click ultra optimizado
func (s *UltraFastSpammer) sendMouseClickUltraFast() {
	// MouseDown - sin verificaciones para máxima velocidad
	s.postMessage.Call(
		uintptr(s.windowHandle),
		constants.WM_LBUTTONDOWN,
		0,
		0,
	)

	// Delay mínimo para que RO procese el mouse
	time.Sleep(9 * time.Microsecond)

	// MouseUp
	s.postMessage.Call(
		uintptr(s.windowHandle),
		constants.WM_LBUTTONUP,
		0,
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
	case "A":
		return 0x41
	case "B":
		return 0x42
	case "C":
		return 0x43
	case "D":
		return 0x44
	case "E":
		return 0x45
	case "F":
		return 0x46
	case "G":
		return 0x47
	case "H":
		return 0x48
	case "I":
		return 0x49
	case "J":
		return 0x4A
	case "K":
		return 0x4B
	case "L":
		return 0x4C
	case "M":
		return 0x4D
	case "N":
		return 0x4E
	case "O":
		return 0x4F
	case "P":
		return 0x50
	case "Q":
		return 0x51
	case "R":
		return 0x52
	case "S":
		return 0x53
	case "T":
		return 0x54
	case "U":
		return 0x55
	case "V":
		return 0x56
	case "W":
		return 0x57
	case "X":
		return 0x58
	case "Y":
		return 0x59
	case "Z":
		return 0x5A
	default:
		return 0
	}
}

// GetStats retorna estadísticas
func (s *UltraFastSpammer) GetStats() (int64, time.Time) {
	return s.loopCount, s.lastLogTime
}
