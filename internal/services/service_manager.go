package services

import (
	"context"
	"sync"
	"time"

	"dtap-go/internal/input"
	"dtap-go/internal/memory"
	"dtap-go/internal/normal/autobuff"
	"dtap-go/internal/normal/autopot"
	"dtap-go/internal/normal/spammer"
	"dtap-go/internal/scheduler"
	"dtap-go/pkg/config"
	"dtap-go/pkg/logger"

	"golang.org/x/sys/windows"
)

// ServiceManager maneja todos los servicios de forma coordinada
type ServiceManager struct {
	config       *config.Config
	memClient    *memory.Client
	windowHandle windows.Handle

	// Servicios
	spammerService  *spammer.UltraFastSpammer
	autopotService  *autopot.UltraFastAutoPot
	autobuffService *autobuff.UltraFastAutoBuff

	// Scheduler centralizado
	scheduler *scheduler.Scheduler

	// Control de goroutines
	wg      sync.WaitGroup
	ctx     context.Context
	cancel  context.CancelFunc
	started bool
	mu      sync.RWMutex

	// Estadísticas
	stats struct {
		StartTime    time.Time
		TotalLoops   int64
		TotalActions int64
	}
}

// NewServiceManager crea un nuevo manager de servicios
func NewServiceManager(cfg *config.Config, memClient *memory.Client, windowHandle windows.Handle) *ServiceManager {
	ctx, cancel := context.WithCancel(context.Background())

	return &ServiceManager{
		config:       cfg,
		memClient:    memClient,
		windowHandle: windowHandle,
		scheduler:    scheduler.NewScheduler(),
		ctx:          ctx,
		cancel:       cancel,
		stats: struct {
			StartTime    time.Time
			TotalLoops   int64
			TotalActions int64
		}{
			StartTime: time.Now(),
		},
	}
}

// Start inicia todos los servicios de forma coordinada con scheduler
func (sm *ServiceManager) Start() error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if sm.started {
		return nil // Ya está iniciado
	}

	logger.Info("🚀 Iniciando ServiceManager con scheduler centralizado")

	// Iniciar dispatcher de input (una sola vez)
	input.Start(sm.ctx, sm.windowHandle)

	// Inicializar servicios
	sm.spammerService = spammer.NewUltraFastSpammer(sm.config, sm.windowHandle)
	sm.autopotService = autopot.NewUltraFastAutoPot(sm.config, sm.memClient, sm.windowHandle)
	sm.autobuffService = autobuff.NewUltraFastAutoBuff(sm.config, sm.memClient, sm.windowHandle)

	// Configurar tareas del scheduler con prioridades
	sm.setupSchedulerTasks()

	// Iniciar scheduler
	sm.scheduler.Start()

	sm.started = true
	logger.Info("✅ ServiceManager iniciado con scheduler")

	return nil
}

// setupSchedulerTasks configura las tareas del scheduler con prioridades
func (sm *ServiceManager) setupSchedulerTasks() {
	// AutoPot - Prioridad crítica (supervivencia)
	if sm.config.AutoPot.Enabled {
		sm.scheduler.AddTask(&scheduler.Task{
			ID:          "autopot",
			Priority:    scheduler.PriorityCritical,
			Interval:    15 * time.Millisecond, // 66 Hz - crítico para supervivencia
			Handler:     sm.autopotService.ProcessAutoPot,
			Enabled:     true,
			Description: "AutoPot - Supervivencia crítica",
		})
	}

	// Spammer - Prioridad alta (combate)
	if sm.config.Spammer.Enabled {
		sm.scheduler.AddTask(&scheduler.Task{
			ID:          "spammer",
			Priority:    scheduler.PriorityHigh,
			Interval:    10 * time.Millisecond, // 100 Hz - balance velocidad/estabilidad
			Handler:     sm.spammerService.ProcessSpammer,
			Enabled:     true,
			Description: "Spammer - Combate prioritario",
		})
	}

	// AutoBuff - Prioridad normal (soporte)
	if sm.config.AutoBuff.Enabled {
		sm.scheduler.AddTask(&scheduler.Task{
			ID:          "autobuff",
			Priority:    scheduler.PriorityNormal,
			Interval:    250 * time.Millisecond, // 4 Hz - menos crítico
			Handler:     sm.autobuffService.ProcessAutoBuff,
			Enabled:     true,
			Description: "AutoBuff - Soporte",
		})
	}
}

// startService inicia un servicio con control de goroutines
func (sm *ServiceManager) startService(name string, startFunc func(context.Context)) {
	sm.wg.Add(1)
	go func() {
		defer sm.wg.Done()
		defer func() {
			if r := recover(); r != nil {
				logger.Error("❌ Panic en servicio %s: %v", name, r)
			}
		}()

		logger.Info("🔄 Iniciando servicio: %s", name)
		startFunc(sm.ctx)
		logger.Info("⏹️ Servicio %s terminado", name)
	}()
}

// Stop detiene todos los servicios de forma graceful
func (sm *ServiceManager) Stop() error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if !sm.started {
		return nil // Ya está detenido
	}

	logger.Info("🛑 Deteniendo ServiceManager...")

	// Detener scheduler
	sm.scheduler.Stop()

	// Cancelar contexto para todos los servicios
	sm.cancel()

	// Esperar a que todas las goroutines terminen (con timeout)
	done := make(chan struct{})
	go func() {
		sm.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		logger.Info("✅ ServiceManager detenido correctamente")
	case <-time.After(5 * time.Second):
		logger.Warn("⚠️ Timeout esperando que los servicios terminen")
	}

	sm.started = false
	return nil
}

// GetStats retorna estadísticas del ServiceManager
func (sm *ServiceManager) GetStats() map[string]interface{} {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	stats := map[string]interface{}{
		"started":     sm.started,
		"uptime":      time.Since(sm.stats.StartTime),
		"total_loops": sm.stats.TotalLoops,
	}

	// Agregar estadísticas de cada servicio si están disponibles
	if sm.spammerService != nil {
		spammerLoops, spammerTime := sm.spammerService.GetStats()
		stats["spammer_loops"] = spammerLoops
		stats["spammer_last_log"] = spammerTime
	}

	if sm.autopotService != nil {
		autopotLoops, autopotPots, autopotTime := sm.autopotService.GetStats()
		stats["autopot_loops"] = autopotLoops
		stats["autopot_potions"] = autopotPots
		stats["autopot_last_log"] = autopotTime
	}

	if sm.autobuffService != nil {
		autobuffLoops, autobuffBuffs, autobuffTime := sm.autobuffService.GetStats()
		stats["autobuff_loops"] = autobuffLoops
		stats["autobuff_buffs"] = autobuffBuffs
		stats["autobuff_last_log"] = autobuffTime
	}

	return stats
}

// UpdateConfig actualiza la configuración de todos los servicios
func (sm *ServiceManager) UpdateConfig(cfg *config.Config) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.config = cfg

	if sm.spammerService != nil {
		// El spammer no tiene UpdateConfig, se recrea
		sm.spammerService = spammer.NewUltraFastSpammer(cfg, sm.windowHandle)
	}

	if sm.autopotService != nil {
		sm.autopotService.UpdateConfig(cfg)
	}

	if sm.autobuffService != nil {
		sm.autobuffService.UpdateConfig(cfg)
	}

	logger.Info("🔄 Configuración actualizada en todos los servicios")
}

// IsRunning retorna si el ServiceManager está ejecutándose
func (sm *ServiceManager) IsRunning() bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.started
}
