package metrics

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof" // Importar pprof endpoints
	"runtime"
	"sync"
	"time"

	"dtap-gui/pkg/logger"
)

// Metrics recopila métricas de rendimiento del sistema
type Metrics struct {
	mu sync.RWMutex

	// Métricas de CPU y memoria
	CPUUsage    float64
	MemoryUsage uint64
	GoRoutines  int

	// Métricas de servicios
	SpammerLoops   int64
	AutoPotLoops   int64
	AutoBuffLoops  int64
	SpammerEvents  int64
	AutoPotEvents  int64
	AutoBuffEvents int64

	// Métricas de timing
	LastUpdate time.Time
	Uptime     time.Duration
	StartTime  time.Time

	// Métricas de input dispatcher
	InputEventsSent    int64
	InputEventsQueued  int64
	InputEventsDropped int64

	// Configuración
	Enabled bool
	Port    int
}

var (
	globalMetrics *Metrics
	once          sync.Once
)

// GetMetrics retorna la instancia global de métricas
func GetMetrics() *Metrics {
	once.Do(func() {
		globalMetrics = &Metrics{
			StartTime: time.Now(),
			Enabled:   false,
			Port:      6060,
		}
	})
	return globalMetrics
}

// Enable habilita las métricas y el servidor pprof
func (m *Metrics) Enable(port int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.Enabled {
		return nil // Ya está habilitado
	}

	m.Enabled = true
	m.Port = port

	// Iniciar servidor HTTP para pprof
	go m.startPprofServer()

	logger.Info("📊 Métricas habilitadas en puerto %d", port)
	logger.Info("🔍 pprof disponible en: http://localhost:%d/debug/pprof/", port)
	return nil
}

// Disable deshabilita las métricas
func (m *Metrics) Disable() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Enabled = false
	logger.Info("📊 Métricas deshabilitadas")
}

// startPprofServer inicia el servidor HTTP para pprof
func (m *Metrics) startPprofServer() {
	// Usar DefaultServeMux que ya incluye pprof
	http.HandleFunc("/metrics", m.handleMetrics)

	logger.Info("🚀 Servidor de métricas iniciado en puerto %d", m.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", m.Port), nil); err != nil && err != http.ErrServerClosed {
		logger.Error("Error en servidor de métricas: %v", err)
	}
}

// handleMetrics maneja el endpoint /metrics
func (m *Metrics) handleMetrics(w http.ResponseWriter, r *http.Request) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	w.Header().Set("Content-Type", "text/plain")

	fmt.Fprintf(w, "# DTAP Metrics\n")
	fmt.Fprintf(w, "dtap_uptime_seconds %f\n", m.Uptime.Seconds())
	fmt.Fprintf(w, "dtap_cpu_usage_percent %f\n", m.CPUUsage)
	fmt.Fprintf(w, "dtap_memory_usage_bytes %d\n", m.MemoryUsage)
	fmt.Fprintf(w, "dtap_goroutines %d\n", m.GoRoutines)

	fmt.Fprintf(w, "dtap_spammer_loops %d\n", m.SpammerLoops)
	fmt.Fprintf(w, "dtap_autopot_loops %d\n", m.AutoPotLoops)
	fmt.Fprintf(w, "dtap_autobuff_loops %d\n", m.AutoBuffLoops)

	fmt.Fprintf(w, "dtap_spammer_events %d\n", m.SpammerEvents)
	fmt.Fprintf(w, "dtap_autopot_events %d\n", m.AutoPotEvents)
	fmt.Fprintf(w, "dtap_autobuff_events %d\n", m.AutoBuffEvents)

	fmt.Fprintf(w, "dtap_input_events_sent %d\n", m.InputEventsSent)
	fmt.Fprintf(w, "dtap_input_events_queued %d\n", m.InputEventsQueued)
	fmt.Fprintf(w, "dtap_input_events_dropped %d\n", m.InputEventsDropped)
}

// UpdateSystemMetrics actualiza las métricas del sistema
func (m *Metrics) UpdateSystemMetrics() {
	if !m.Enabled {
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Actualizar tiempo de actividad
	m.Uptime = time.Since(m.StartTime)
	m.LastUpdate = time.Now()

	// Actualizar métricas de Go
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	m.MemoryUsage = memStats.Alloc
	m.GoRoutines = runtime.NumGoroutine()

	// CPU usage sería más complejo de implementar, por ahora lo dejamos en 0
	// En una implementación real, usarías gopsutil o similar
}

// IncrementSpammerLoops incrementa el contador de loops del spammer
func (m *Metrics) IncrementSpammerLoops() {
	if !m.Enabled {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.SpammerLoops++
}

// IncrementAutoPotLoops incrementa el contador de loops del autopot
func (m *Metrics) IncrementAutoPotLoops() {
	if !m.Enabled {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.AutoPotLoops++
}

// IncrementAutoBuffLoops incrementa el contador de loops del autobuff
func (m *Metrics) IncrementAutoBuffLoops() {
	if !m.Enabled {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.AutoBuffLoops++
}

// IncrementSpammerEvents incrementa el contador de eventos del spammer
func (m *Metrics) IncrementSpammerEvents() {
	if !m.Enabled {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.SpammerEvents++
}

// IncrementAutoPotEvents incrementa el contador de eventos del autopot
func (m *Metrics) IncrementAutoPotEvents() {
	if !m.Enabled {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.AutoPotEvents++
}

// IncrementAutoBuffEvents incrementa el contador de eventos del autobuff
func (m *Metrics) IncrementAutoBuffEvents() {
	if !m.Enabled {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.AutoBuffEvents++
}

// IncrementInputEvents incrementa el contador de eventos de input
func (m *Metrics) IncrementInputEvents() {
	if !m.Enabled {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.InputEventsSent++
}

// SetInputEventsQueued establece el número de eventos en cola
func (m *Metrics) SetInputEventsQueued(count int64) {
	if !m.Enabled {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.InputEventsQueued = count
}

// IncrementInputEventsDropped incrementa el contador de eventos descartados
func (m *Metrics) IncrementInputEventsDropped() {
	if !m.Enabled {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.InputEventsDropped++
}

// GetStats retorna las métricas actuales
func (m *Metrics) GetStats() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return map[string]interface{}{
		"uptime":               m.Uptime.Seconds(),
		"cpu_usage":            m.CPUUsage,
		"memory_usage":         m.MemoryUsage,
		"goroutines":           m.GoRoutines,
		"spammer_loops":        m.SpammerLoops,
		"autopot_loops":        m.AutoPotLoops,
		"autobuff_loops":       m.AutoBuffLoops,
		"spammer_events":       m.SpammerEvents,
		"autopot_events":       m.AutoPotEvents,
		"autobuff_events":      m.AutoBuffEvents,
		"input_events_sent":    m.InputEventsSent,
		"input_events_queued":  m.InputEventsQueued,
		"input_events_dropped": m.InputEventsDropped,
		"enabled":              m.Enabled,
	}
}

// StartMetricsUpdater inicia el actualizador de métricas
func (m *Metrics) StartMetricsUpdater(ctx context.Context) {
	if !m.Enabled {
		return
	}

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			m.UpdateSystemMetrics()
		}
	}
}
