package constants

import "time"

// Timing constants for DTAP services
const (
	// AutoPot timing constants - Supervivencia (Priority 0)
	AutoPotCheckInterval     = 15 * time.Millisecond  // 66 Hz - crítico para supervivencia
	AutoPotSpamDuration      = 2 * time.Second        // Duración del spam continuo
	AutoPotPredictionSamples = 20                     // Muestras para predicción
	AutoPotPredictionWindow  = 200 * time.Millisecond // Ventana de predicción

	// ATK-DEF timing constants - Equipo (Priority 1)
	ATKDEFCheckInterval = 1 * time.Millisecond // 100 Hz - muy crítico para equipo
	ATKDEFClickDelay    = 1 * time.Millisecond // Delay para clicks

	// Spammer timing constants - Combate (Priority 2)
	SpammerMinDelay = 1 * time.Millisecond // Delay mínimo para estabilidad

	// AutoBuff timing constants - Soporte (Priority 3)
	AutoBuffCheckInterval  = 250 * time.Millisecond // 4 Hz - menos crítico
	AutoBuffSlowModeDelay  = 3 * time.Second        // Delay en modo lento
	AutoBuffSuccessTimeout = 5 * time.Second        // Timeout para buff exitoso

	// Scheduler timing constants
	SchedulerBaseInterval = 1 * time.Millisecond // Intervalo base del scheduler

	// Input timing constants
	InputMinDelay = 1 * time.Millisecond // Delay mínimo para input

	// Metrics timing constants
	MetricsUpdateInterval = 1 * time.Second // Intervalo de actualización de métricas

	// Service timeout constants
	ServiceStartTimeout = 5 * time.Second // Timeout para inicio de servicios

	// Death monitoring constants
	DeathCheckInterval  = 100 * time.Millisecond  // Intervalo para verificar muerte
	ResurrectionDelay   = 1500 * time.Millisecond // Delay después de resurrección
	MemoryErrorCooldown = 500 * time.Millisecond  // Cooldown después de error de memoria
)
