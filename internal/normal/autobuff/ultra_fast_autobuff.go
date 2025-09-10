package autobuff

import (
	"context"
	"syscall"
	"time"

	"dtap-go/internal/memory"
	"dtap-go/internal/metrics"
	"dtap-go/pkg/config"
	"dtap-go/pkg/constants"
	"dtap-go/pkg/logger"

	"golang.org/x/sys/windows"
)

// BuffTypes - Tipos de buffs en Ragnarok Online
type BuffTypes uint32

const (
	ST_PROVOKE               BuffTypes = 0
	ST_IMPROVE_CONCENTRATION BuffTypes = 3 // Gloom
	ST_QUAGMIRE              BuffTypes = 8
	ST_ASPD_A                BuffTypes = 37 // ASPD
	ST_ASPD_B                BuffTypes = 38
	ST_ASPD_C                BuffTypes = 39
	ST_BUFFER_END            BuffTypes = 0xFFFFFFFF
)

// UltraFastAutoBuff - Versión ultra optimizada del autobuff
type UltraFastAutoBuff struct {
	config       *config.Config
	memClient    *memory.Client
	windowHandle windows.Handle

	// Direct Windows API - sin abstracciones
	user32      *syscall.LazyDLL
	postMessage *syscall.LazyProc

	// Pre-allocated variables para evitar GC
	gloomKeyCode uint32
	aspdKeyCode  uint32
	loopCount    int64
	lastLogTime  time.Time

	// Performance tracking
	buffsUsed    int64
	lastBuffTime time.Time

	// Buff status tracking
	foundGloom    bool
	foundAspd     bool
	foundQuagmire bool
	statusBuffer  []uint32 // Buffer para leer status

	// Timing control
	lastGloomUse time.Time
	lastAspdUse  time.Time

	// Smart speed control for missing consumables
	lastSuccessfulBuff  time.Time
	consecutiveFailures int
	slowMode            bool
	slowModeDelay       time.Duration
}

// NewUltraFastAutoBuff crea el autobuff más rápido posible
func NewUltraFastAutoBuff(cfg *config.Config, memClient *memory.Client, windowHandle windows.Handle) *UltraFastAutoBuff {
	user32 := syscall.NewLazyDLL("user32.dll")

	// Pre-convertir teclas a códigos para evitar conversión en runtime
	gloomKeyCode := getKeyCodeUltraFast(cfg.AutoBuff.GloomKey)
	aspdKeyCode := getKeyCodeUltraFast(cfg.AutoBuff.AspdKey)

	return &UltraFastAutoBuff{
		config:       cfg,
		memClient:    memClient,
		windowHandle: windowHandle,
		user32:       user32,
		postMessage:  user32.NewProc("PostMessageW"),
		gloomKeyCode: gloomKeyCode,
		aspdKeyCode:  aspdKeyCode,
		lastLogTime:  time.Now(),
		lastBuffTime: time.Now(),
		statusBuffer: make([]uint32, 100), // Buffer de 100 elementos
		lastGloomUse: time.Now(),
		lastAspdUse:  time.Now(),
		// Smart speed control initialization
		lastSuccessfulBuff:  time.Now(),
		consecutiveFailures: 0,
		slowMode:            false,
		slowModeDelay:       3 * time.Second,
	}
}

// Start inicia el autobuff ultra rápido
func (s *UltraFastAutoBuff) Start(ctx context.Context) {
	if !s.config.AutoBuff.Enabled {
		logger.Info("UltraFast AutoBuff está deshabilitado")
		return
	}

	logger.Info("🚀 UltraFast AutoBuff iniciado - Modo ULTRA RÁPIDO")
	go s.runUltraFastLoop(ctx)
}

// ProcessAutoBuff procesa el autobuff (para scheduler)
func (s *UltraFastAutoBuff) ProcessAutoBuff(ctx context.Context) error {
	s.processAutoBuffUltraFast()
	s.loopCount++

	// Sistema inteligente de velocidad
	s.adjustSpeedBasedOnSuccess()

	// Reportar métricas
	metrics.GetMetrics().IncrementAutoBuffLoops()

	return nil
}

// runUltraFastLoop - Loop optimizado con Ticker
func (s *UltraFastAutoBuff) runUltraFastLoop(ctx context.Context) {
	// Ticker para AutoBuff: 250ms (4 Hz) - menos crítico, optimizado para eficiencia
	ticker := time.NewTicker(250 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.loopCount++

			// Procesar autobuff con timing controlado
			s.processAutoBuffUltraFast()

			// Sistema inteligente de velocidad
			s.adjustSpeedBasedOnSuccess()

			// Reportar métricas
			metrics.GetMetrics().IncrementAutoBuffLoops()
		}
	}
}

// processAutoBuffUltraFast - Procesamiento ultra rápido de buffs
func (s *UltraFastAutoBuff) processAutoBuffUltraFast() {
	now := time.Now()

	// Leer status buffer de forma ultra rápida
	err := s.readStatusBuffer()
	if err != nil {
		return // Silenciar errores para máxima velocidad
	}

	// Analizar buffs de forma ultra rápida
	s.analyzeBuffs()

	// Usar buffs si es necesario
	s.useBuffsIfNeeded(now)
}

// readStatusBuffer - Lee el status buffer de memoria
func (s *UltraFastAutoBuff) readStatusBuffer() error {
	// Leer el status buffer completo de una vez
	statusBufferAddr, err := s.memClient.GetStatusBufferAddress()
	if err != nil {
		return err
	}

	// Leer 100 elementos del buffer (400 bytes)
	for i := 0; i < 100; i++ {
		value, err := s.memClient.ReadUint32(statusBufferAddr + uintptr(i*4))
		if err != nil {
			return err
		}
		s.statusBuffer[i] = value
	}

	return nil
}

// analyzeBuffs - Analiza qué buffs están activos
func (s *UltraFastAutoBuff) analyzeBuffs() {
	s.foundGloom = false
	s.foundAspd = false
	s.foundQuagmire = false

	// Analizar buffer de forma ultra rápida
	for i := 0; i < 100; i++ {
		buffValue := s.statusBuffer[i]

		switch BuffTypes(buffValue) {
		case ST_IMPROVE_CONCENTRATION:
			s.foundGloom = true
		case ST_ASPD_A, ST_ASPD_B, ST_ASPD_C:
			s.foundAspd = true
		case ST_QUAGMIRE:
			s.foundQuagmire = true
		case ST_BUFFER_END:
			return // Fin del buffer
		}

		// Si ya encontramos todos los buffs necesarios, salir
		if s.foundGloom && s.foundAspd {
			break
		}
	}
}

// useBuffsIfNeeded - Usa buffs si es necesario
func (s *UltraFastAutoBuff) useBuffsIfNeeded(now time.Time) {
	// Usar Gloom si no está activo y no hay Quagmire
	if !s.foundGloom && !s.foundQuagmire {
		// Verificar delay mínimo entre usos - ultra rápido (1ms)
		if now.Sub(s.lastGloomUse) >= 1*time.Millisecond {
			s.sendKeyUltraFast(s.gloomKeyCode)
			s.buffsUsed++
			s.lastGloomUse = now
			s.lastBuffTime = now
			s.lastSuccessfulBuff = now

			// Reportar métricas
			metrics.GetMetrics().IncrementAutoBuffEvents()
		}
	}

	// Usar ASPD si no está activo
	if !s.foundAspd {
		// Verificar delay mínimo entre usos - ultra rápido (1ms)
		if now.Sub(s.lastAspdUse) >= 1*time.Millisecond {
			s.sendKeyUltraFast(s.aspdKeyCode)
			s.buffsUsed++
			s.lastAspdUse = now
			s.lastBuffTime = now
			s.lastSuccessfulBuff = now

			// Reportar métricas
			metrics.GetMetrics().IncrementAutoBuffEvents()
		}
	}
}

// adjustSpeedBasedOnSuccess - Sistema inteligente de velocidad
func (s *UltraFastAutoBuff) adjustSpeedBasedOnSuccess() {
	now := time.Now()

	// Solo activar modo lento si realmente no hay consumibles disponibles
	// y han pasado más de 5 segundos sin buffs exitosos
	if now.Sub(s.lastSuccessfulBuff) > 5*time.Second && !s.slowMode {
		// Verificar si realmente no hay consumibles disponibles
		if !s.foundGloom && !s.foundAspd {
			s.slowMode = true
			logger.Info("🔄 AutoBuff: Sin consumibles detectado, activando modo lento (3s delay)")
		}
	}

	// Solo volver al modo rápido si realmente se usaron consumibles
	// y tenemos buffs activos
	if s.slowMode && (s.foundGloom || s.foundAspd) {
		s.slowMode = false
		s.consecutiveFailures = 0
		logger.Info("⚡ AutoBuff: Consumibles disponibles, volviendo a modo ultra rápido")
	}
}

// sendKeyUltraFast - Envío de tecla ultra rápido (mismo que spammer y autopot)
func (s *UltraFastAutoBuff) sendKeyUltraFast(keyCode uint32) {
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

	// Delay mínimo para que RO procese el input - ultra rápido
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

// GetStats retorna estadísticas del autobuff
func (s *UltraFastAutoBuff) GetStats() (int64, int64, time.Time) {
	return s.loopCount, s.buffsUsed, s.lastLogTime
}

// UpdateConfig actualiza la configuración
func (s *UltraFastAutoBuff) UpdateConfig(cfg *config.Config) {
	s.config = cfg
	s.gloomKeyCode = getKeyCodeUltraFast(cfg.AutoBuff.GloomKey)
	s.aspdKeyCode = getKeyCodeUltraFast(cfg.AutoBuff.AspdKey)
	logger.Info("Configuración del UltraFast AutoBuff actualizada")
}
