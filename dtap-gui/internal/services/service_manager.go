package services

import (
	"context"
	"sync"
	"time"

	"dtap-gui/internal/memory"
	"dtap-gui/internal/normal/autobuff"
	"dtap-gui/internal/normal/autopot"
	"dtap-gui/internal/normal/spammer"
	"dtap-gui/pkg/config"
	"dtap-gui/pkg/logger"

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

// Start inicia todos los servicios de forma coordinada
func (sm *ServiceManager) Start() error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if sm.started {
		return nil // Ya está iniciado
	}

	logger.Info("🚀 Iniciando ServiceManager con mejores prácticas")

	// Inicializar servicios
	sm.spammerService = spammer.NewUltraFastSpammer(sm.config, sm.windowHandle)
	sm.autopotService = autopot.NewUltraFastAutoPot(sm.config, sm.memClient, sm.windowHandle)
	sm.autobuffService = autobuff.NewUltraFastAutoBuff(sm.config, sm.memClient, sm.windowHandle)

	// Iniciar servicios con control de goroutines
	sm.startService("Spammer", sm.spammerService.Start)
	sm.startService("AutoPot", sm.autopotService.Start)
	sm.startService("AutoBuff", sm.autobuffService.Start)

	sm.started = true
	logger.Info("✅ ServiceManager iniciado correctamente")

	return nil
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
