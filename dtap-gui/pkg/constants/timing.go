package constants

import "time"

// Timing constants for DTAP services
const (
	// AutoPot timing constants
	AutoPotCheckInterval     = 15 * time.Millisecond  // 66 Hz - crítico para supervivencia
	AutoPotSpamDuration      = 2 * time.Second        // Duración del spam continuo
	AutoPotPredictionSamples = 20                     // Muestras para predicción
	AutoPotPredictionWindow  = 200 * time.Millisecond // Ventana de predicción

	// AutoBuff timing constants
	AutoBuffCheckInterval  = 250 * time.Millisecond // 4 Hz - menos crítico
	AutoBuffSlowModeDelay  = 3 * time.Second        // Delay en modo lento
	AutoBuffSuccessTimeout = 5 * time.Second        // Timeout para buff exitoso

	// Spammer timing constants
	SpammerMinDelay = 1 * time.Millisecond // Delay mínimo para estabilidad

	// ATK-DEF timing constants
	ATKDEFClickDelay = 1 * time.Millisecond // Delay para clicks

	// Scheduler timing constants
	SchedulerBaseInterval = 1 * time.Millisecond // Intervalo base del scheduler

	// Input timing constants
	InputMinDelay = 1 * time.Millisecond // Delay mínimo para input

	// Metrics timing constants
	MetricsUpdateInterval = 1 * time.Second // Intervalo de actualización de métricas

	// Service timeout constants
	ServiceStartTimeout = 5 * time.Second // Timeout para inicio de servicios
)
