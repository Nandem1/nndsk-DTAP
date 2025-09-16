package autopot

import (
	"context"
	"syscall"
	"time"

	"dtap-go/internal/memory"
	"dtap-go/internal/metrics"
	"dtap-go/internal/prediction"
	"dtap-go/pkg/config"
	"dtap-go/pkg/constants"
	"dtap-go/pkg/logger"

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

	// Sistema de predicción optimizado con ring buffer
	predictor      *prediction.Predictor // Predictor con ring buffer
	lastHPCheck    time.Time             // Última vez que se revisó HP
	lastSPCheck    time.Time             // Última vez que se revisó SP
	predictionMode bool                  // Si está en modo predicción activo
	damageRate     float64               // Tasa de daño por segundo detectada

	// Tracking de predicción para logs
	predictionStartTime time.Time     // Cuando empezó la predicción
	predictionDuration  time.Duration // Duración total de la predicción
	lastHPAt100Percent  time.Time     // Última vez que HP estuvo al 100%

	// Sistema de spam continuo por daño alto
	continuousSpamMode     bool          // Si está en modo spam continuo
	continuousSpamStart    time.Time     // Cuando empezó el spam continuo
	continuousSpamDuration time.Duration // Duración del spam continuo (5 segundos)
	highDamageThreshold    float64       // Umbral de daño alto (ej: -20% por segundo)
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
		predictor:      prediction.NewPredictor(20, 200*time.Millisecond), // 20 muestras, predicción 200ms
		lastHPCheck:    time.Now(),
		lastSPCheck:    time.Now(),
		predictionMode: false,
		damageRate:     0.0,
		// Inicializar sistema de spam continuo
		continuousSpamMode:     false,
		continuousSpamDuration: 5 * time.Second, // 5 segundos de spam continuo
		highDamageThreshold:    -15.0,           // -15% HP por segundo = daño alto
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

// ProcessAutoPot procesa el autopot (para scheduler)
func (s *UltraFastAutoPot) ProcessAutoPot(ctx context.Context) error {
	s.processAutoPotUltraFast()
	s.loopCount++

	// Reportar métricas
	metrics.GetMetrics().IncrementAutoPotLoops()

	return nil
}

// runUltraFastLoop - Loop optimizado con Ticker
func (s *UltraFastAutoPot) runUltraFastLoop(ctx context.Context) {
	// Ticker para AutoPot: 15ms (66 Hz) - crítico para supervivencia
	ticker := time.NewTicker(15 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.loopCount++

			// Procesar autopot con timing controlado
			s.processAutoPotUltraFast()

			// Reportar métricas
			metrics.GetMetrics().IncrementAutoPotLoops()
		}
	}
}

// processAutoPotUltraFast - Procesamiento ultra rápido con predicción
func (s *UltraFastAutoPot) processAutoPotUltraFast() {
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

	// Actualizar predictor con nuevos datos
	s.predictor.AddHP(float64(hpPercent))
	s.predictor.AddSP(float64(spPercent))

	// Sistema de predicción: anticipar el pot 200ms antes
	predictedHP, hpPredictionValid := s.predictor.PredictHP(200 * time.Millisecond)
	predictedSP, spPredictionValid := s.predictor.PredictSP(200 * time.Millisecond)

	// Calcular tendencia de daño (cambio de HP por segundo)
	damageRate := s.calculateDamageRate()

	// Auto pot HP con predicción y spam continuo - ultra rápido
	shouldUseHPPot := hpPercent < s.config.AutoPot.HPPercent
	usedPrediction := false
	usedContinuousSpam := false
	now := time.Now()

	// Verificar si estamos en modo spam continuo
	if s.continuousSpamMode {
		timeInSpam := now.Sub(s.continuousSpamStart)
		if timeInSpam < s.continuousSpamDuration {
			// Seguir spameando por 5 segundos sin importar el HP
			shouldUseHPPot = true
			usedContinuousSpam = true
		} else {
			// Terminar modo spam continuo
			s.continuousSpamMode = false
			logger.Info("🛑 Spam continuo terminado después de %v", timeInSpam.Round(time.Millisecond))
		}
	} else if s.config.AutoPot.PredictionMode && hpPredictionValid {
		// Solo usar predicción si está habilitada en la configuración
		// Detectar daño alto y activar spam continuo
		if damageRate < s.highDamageThreshold {
			// Daño alto detectado, activar spam continuo por 5 segundos
			s.continuousSpamMode = true
			s.continuousSpamStart = now
			shouldUseHPPot = true
			usedContinuousSpam = true
			logger.Info("🚨 DAÑO ALTO DETECTADO: %.1f%%/seg - Activando spam continuo por 5 segundos", damageRate)
		} else {
			// Predicción normal
			predictionTriggered := predictedHP < float64(s.config.AutoPot.HPPercent)
			if predictionTriggered {
				shouldUseHPPot = true
				usedPrediction = true

				// Trackear tiempo de predicción (independientemente del HP actual)
				if s.predictionStartTime.IsZero() {
					s.predictionStartTime = now
					s.lastHPAt100Percent = now
				}
			}
		}
	}

	if shouldUseHPPot {
		s.sendKeyUltraFast(s.hpPotKeyCode)
		s.potionsUsed++
		s.lastHPPercent = hpPercent

		// Log detallado del motivo del potting
		if usedContinuousSpam {
			timeInSpam := now.Sub(s.continuousSpamStart)
			logger.Info("🚨 HP Pot SPAM CONTINUO: HP actual=%d%% - Spameando desde hace %v (daño: %.1f%%/seg)",
				hpPercent, timeInSpam.Round(time.Millisecond), damageRate)
		} else if usedPrediction {
			// Calcular duración de predicción
			var predictionDuration time.Duration
			if !s.predictionStartTime.IsZero() {
				predictionDuration = now.Sub(s.predictionStartTime)
			}

			if hpPercent >= s.config.AutoPot.HPPercent {
				logger.Info("🧠 HP Pot PREDICCIÓN: HP actual=%d%%, predicho=%.1f%% en 200ms (threshold=%d%%) - Spameando desde hace %v",
					hpPercent, predictedHP, s.config.AutoPot.HPPercent, predictionDuration.Round(time.Millisecond))
			} else {
				logger.Info("🧠 HP Pot PREDICCIÓN: HP actual=%d%%, predicho=%.1f%% en 200ms (threshold=%d%%) - Spameando desde hace %v",
					hpPercent, predictedHP, s.config.AutoPot.HPPercent, predictionDuration.Round(time.Millisecond))
			}
		} else {
			logger.Info("💊 HP Pot NORMAL: HP actual=%d%% (threshold=%d%%)",
				hpPercent, s.config.AutoPot.HPPercent)
		}

		// Reportar métricas
		metrics.GetMetrics().IncrementAutoPotEvents()
	}

	// Auto pot SP con predicción - ultra rápido
	shouldUseSPPot := spPercent < s.config.AutoPot.SPPercent
	usedSPPrediction := false

	if spPredictionValid {
		spPredictionTriggered := predictedSP < float64(s.config.AutoPot.SPPercent)
		// Usar predicción si predice que bajará por debajo del threshold
		// Esto permite spamear SP pots aunque esté al 100% (como High Ping Mode inteligente)
		if spPredictionTriggered {
			shouldUseSPPot = true
			usedSPPrediction = true
		}
	}

	if shouldUseSPPot {
		s.sendKeyUltraFast(s.spPotKeyCode)
		s.potionsUsed++
		s.lastSPPercent = spPercent

		// Log detallado del motivo del potting
		if usedSPPrediction {
			logger.Info("🧠 SP Pot PREDICCIÓN: SP actual=%d%%, predicho=%.1f%% en 200ms (threshold=%d%%)",
				spPercent, predictedSP, s.config.AutoPot.SPPercent)
		} else {
			logger.Info("💊 SP Pot NORMAL: SP actual=%d%% (threshold=%d%%)",
				spPercent, s.config.AutoPot.SPPercent)
		}

		// Reportar métricas
		metrics.GetMetrics().IncrementAutoPotEvents()
	} else {
		// Resetear tracking de predicción si no se está usando
		if !usedPrediction && !s.predictionStartTime.IsZero() {
			// La predicción terminó, mostrar duración total
			totalDuration := now.Sub(s.predictionStartTime)
			if totalDuration > 50*time.Millisecond { // Solo log si duró más de 50ms
				logger.Info("✅ Predicción terminada: Duró %v spameando HP pots al 100%%",
					totalDuration.Round(time.Millisecond))
			}
			s.predictionStartTime = time.Time{} // Reset
		}

		// Log de información del predictor cuando no se usan pots
		if hpPredictionValid && s.loopCount%100 == 0 { // Cada 100 loops para no spamear
			logger.Info("📊 HP Status: actual=%d%%, predicho=%.1f%% en 200ms, threshold=%d%%",
				hpPercent, predictedHP, s.config.AutoPot.HPPercent)
		}
	}

	if !shouldUseSPPot && spPredictionValid && s.loopCount%100 == 0 {
		logger.Info("📊 SP Status: actual=%d%%, predicho=%.1f%% en 200ms, threshold=%d%%",
			spPercent, predictedSP, s.config.AutoPot.SPPercent)
	}
}

// calculateDamageRate calcula la tasa de daño por segundo basada en el historial
func (s *UltraFastAutoPot) calculateDamageRate() float64 {
	// Obtener historial de HP de los últimos 3 segundos
	history := s.predictor.GetHPHistoryByDuration(3 * time.Second)
	if len(history) < 2 {
		return 0.0
	}

	// Calcular la pendiente (tendencia) de los datos
	// Usar regresión lineal simple para obtener la tasa de cambio
	first := history[0]
	last := history[len(history)-1]

	timeDiff := last.Timestamp.Sub(first.Timestamp)
	if timeDiff <= 0 {
		return 0.0
	}

	// Calcular cambio de HP por segundo
	hpChange := last.Value - first.Value
	damageRate := hpChange / timeDiff.Seconds()

	// Actualizar la tasa de daño almacenada
	s.damageRate = damageRate

	return damageRate
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

// GetPredictionStatus retorna el estado actual de la predicción
func (s *UltraFastAutoPot) GetPredictionStatus() map[string]interface{} {
	now := time.Now()

	status := map[string]interface{}{
		"prediction_mode":       s.predictionMode,
		"continuous_spam_mode":  s.continuousSpamMode,
		"damage_rate":           s.damageRate,
		"high_damage_threshold": s.highDamageThreshold,
	}

	if s.continuousSpamMode {
		timeInSpam := now.Sub(s.continuousSpamStart)
		remainingTime := s.continuousSpamDuration - timeInSpam
		if remainingTime < 0 {
			remainingTime = 0
		}

		status["spam_remaining_time"] = remainingTime.Milliseconds()
		status["spam_total_time"] = s.continuousSpamDuration.Milliseconds()
		status["spam_elapsed_time"] = timeInSpam.Milliseconds()
	}

	if !s.predictionStartTime.IsZero() {
		predictionDuration := now.Sub(s.predictionStartTime)
		status["prediction_duration"] = predictionDuration.Milliseconds()
	}

	return status
}

// UpdateConfig actualiza la configuración
func (s *UltraFastAutoPot) UpdateConfig(cfg *config.Config) {
	s.config = cfg
	s.hpPotKeyCode = getKeyCodeUltraFast(cfg.AutoPot.HPPotKey)
	s.spPotKeyCode = getKeyCodeUltraFast(cfg.AutoPot.SPPotKey)
	logger.Info("Configuración del UltraFast AutoPot actualizada")
}
