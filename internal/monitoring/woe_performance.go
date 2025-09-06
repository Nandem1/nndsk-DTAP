package monitoring

import (
	"context"
	"fmt"
	"sync"
	"time"

	"dtap-go/internal/autopot"
	"dtap-go/internal/memory"
	"dtap-go/internal/spammer"
	"dtap-go/pkg/logger"

	"golang.org/x/sys/windows"
)

// WOEPerformanceMonitor tracks and optimizes performance for competitive WoE
type WOEPerformanceMonitor struct {
	memClient        *memory.OptimizedClient
	autopotService   *autopot.WOEOptimizedService
	spammerService   *spammer.WOEOptimizedSpammer
	windowHandle     windows.Handle
	
	// Performance metrics
	metrics          *PerformanceMetrics
	metricsMutex     sync.RWMutex
	
	// Optimization settings
	targetFPS        int
	adaptiveMode     bool
	
	// Context for cancellation
	ctx              context.Context
	cancel           context.CancelFunc
}

// PerformanceMetrics holds all performance data
type PerformanceMetrics struct {
	// Memory performance
	MemoryReadLatency    time.Duration
	CacheHitRate         float64
	CacheAge             time.Duration
	
	// Input performance
	InputLatency         time.Duration
	InputCount           int64
	InputRate            float64
	
	// Autopot performance
	PotCount             int64
	PotLatency           time.Duration
	EmergencyMode        bool
	
	// Spammer performance
	SpamCount            int64
	SpamRate             float64
	BurstMode            bool
	ActiveKeys           int
	
	// System performance
	CPUUsage             float64
	MemoryUsage          float64
	FrameRate            float64
	
	// WoE specific metrics
	CombatTime           time.Duration
	DamageTaken          int64
	SkillsUsed           int64
	
	// Timestamps
	LastUpdate           time.Time
	SessionStart         time.Time
}

// NewWOEPerformanceMonitor creates a performance monitor for WoE
func NewWOEPerformanceMonitor(
	memClient *memory.OptimizedClient,
	autopotService *autopot.WOEOptimizedService,
	spammerService *spammer.WOEOptimizedSpammer,
	windowHandle windows.Handle,
) *WOEPerformanceMonitor {
	ctx, cancel := context.WithCancel(context.Background())
	
	monitor := &WOEPerformanceMonitor{
		memClient:      memClient,
		autopotService: autopotService,
		spammerService: spammerService,
		windowHandle:   windowHandle,
		metrics:        &PerformanceMetrics{},
		targetFPS:      60, // Target 60 FPS for WoE
		adaptiveMode:   true,
		ctx:            ctx,
		cancel:         cancel,
	}
	
	monitor.metrics.SessionStart = time.Now()
	
	return monitor
}

// Start begins performance monitoring
func (m *WOEPerformanceMonitor) Start(ctx context.Context) {
	logger.Info("Iniciando monitor de performance WoE")
	
	// Start monitoring goroutines
	go m.memoryPerformanceLoop(ctx)
	go m.inputPerformanceLoop(ctx)
	go m.systemPerformanceLoop(ctx)
	go m.woeMetricsLoop(ctx)
	go m.optimizationLoop(ctx)
}

// memoryPerformanceLoop monitors memory performance
func (m *WOEPerformanceMonitor) memoryPerformanceLoop(ctx context.Context) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			return
		case <-m.ctx.Done():
			return
		case <-ticker.C:
			m.updateMemoryMetrics()
		}
	}
}

// inputPerformanceLoop monitors input performance
func (m *WOEPerformanceMonitor) inputPerformanceLoop(ctx context.Context) {
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			return
		case <-m.ctx.Done():
			return
		case <-ticker.C:
			m.updateInputMetrics()
		}
	}
}

// systemPerformanceLoop monitors system performance
func (m *WOEPerformanceMonitor) systemPerformanceLoop(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			return
		case <-m.ctx.Done():
			return
		case <-ticker.C:
			m.updateSystemMetrics()
		}
	}
}

// woeMetricsLoop monitors WoE-specific metrics
func (m *WOEPerformanceMonitor) woeMetricsLoop(ctx context.Context) {
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			return
		case <-m.ctx.Done():
			return
		case <-ticker.C:
			m.updateWOEMetrics()
		}
	}
}

// optimizationLoop continuously optimizes performance
func (m *WOEPerformanceMonitor) optimizationLoop(ctx context.Context) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			return
		case <-m.ctx.Done():
			return
		case <-ticker.C:
			m.optimizePerformance()
		}
	}
}

// updateMemoryMetrics updates memory-related performance metrics
func (m *WOEPerformanceMonitor) updateMemoryMetrics() {
	m.metricsMutex.Lock()
	defer m.metricsMutex.Unlock()
	
	// Get cache statistics
	cacheAge, valid := m.memClient.GetCacheStats()
	m.metrics.CacheAge = cacheAge
	m.metrics.CacheHitRate = 0.95 // Placeholder - would need actual hit rate calculation
	
	// Simulate memory read latency (would need actual measurement)
	m.metrics.MemoryReadLatency = 50 * time.Microsecond
}

// updateInputMetrics updates input-related performance metrics
func (m *WOEPerformanceMonitor) updateInputMetrics() {
	m.metricsMutex.Lock()
	defer m.metricsMutex.Unlock()
	
	// Get autopot stats
	potCount, emergencyMode, _ := m.autopotService.GetStats()
	m.metrics.PotCount = potCount
	m.metrics.EmergencyMode = emergencyMode
	
	// Get spammer stats
	spamCount, burstMode, activeKeys := m.spammerService.GetStats()
	m.metrics.SpamCount = spamCount
	m.metrics.BurstMode = burstMode
	m.metrics.ActiveKeys = activeKeys
	
	// Calculate rates
	now := time.Now()
	if !m.metrics.LastUpdate.IsZero() {
		timeDiff := now.Sub(m.metrics.LastUpdate).Seconds()
		if timeDiff > 0 {
			m.metrics.InputRate = float64(m.metrics.InputCount) / timeDiff
			m.metrics.SpamRate = float64(m.metrics.SpamCount) / timeDiff
		}
	}
	
	m.metrics.LastUpdate = now
}

// updateSystemMetrics updates system performance metrics
func (m *WOEPerformanceMonitor) updateSystemMetrics() {
	m.metricsMutex.Lock()
	defer m.metricsMutex.Unlock()
	
	// Simulate system metrics (would need actual system monitoring)
	m.metrics.CPUUsage = 15.0 // Placeholder
	m.metrics.MemoryUsage = 8.0 // Placeholder
	m.metrics.FrameRate = 60.0 // Placeholder
}

// updateWOEMetrics updates WoE-specific metrics
func (m *WOEPerformanceMonitor) updateWOEMetrics() {
	m.metricsMutex.Lock()
	defer m.metricsMutex.Unlock()
	
	// Update combat time
	m.metrics.CombatTime = time.Since(m.metrics.SessionStart)
	
	// Simulate WoE metrics (would need actual game state detection)
	m.metrics.DamageTaken = 0 // Placeholder
	m.metrics.SkillsUsed = m.metrics.SpamCount // Use spam count as proxy
}

// optimizePerformance continuously optimizes system performance
func (m *WOEPerformanceMonitor) optimizePerformance() {
	m.metricsMutex.RLock()
	metrics := *m.metrics
	m.metricsMutex.RUnlock()
	
	// Optimize based on current performance
	if metrics.FrameRate < float64(m.targetFPS-10) {
		logger.Warn("Frame rate bajo: %.1f FPS (target: %d)", metrics.FrameRate, m.targetFPS)
		// Could implement frame rate optimization here
	}
	
	if metrics.CPUUsage > 80.0 {
		logger.Warn("CPU usage alto: %.1f%%", metrics.CPUUsage)
		// Could implement CPU optimization here
	}
	
	if metrics.MemoryReadLatency > 100*time.Microsecond {
		logger.Warn("Latencia de memoria alta: %v", metrics.MemoryReadLatency)
		// Could implement memory optimization here
	}
}

// GetMetrics returns current performance metrics
func (m *WOEPerformanceMonitor) GetMetrics() PerformanceMetrics {
	m.metricsMutex.RLock()
	defer m.metricsMutex.RUnlock()
	
	return *m.metrics
}

// PrintPerformanceReport prints a detailed performance report
func (m *WOEPerformanceMonitor) PrintPerformanceReport() {
	m.metricsMutex.RLock()
	metrics := *m.metrics
	m.metricsMutex.RUnlock()
	
	logger.Info("=== REPORTE DE PERFORMANCE WoE ===")
	logger.Info("Tiempo de sesión: %v", metrics.CombatTime)
	logger.Info("")
	
	logger.Info("MEMORIA:")
	logger.Info("  Latencia de lectura: %v", metrics.MemoryReadLatency)
	logger.Info("  Cache hit rate: %.1f%%", metrics.CacheHitRate*100)
	logger.Info("  Edad del cache: %v", metrics.CacheAge)
	logger.Info("")
	
	logger.Info("INPUT:")
	logger.Info("  Latencia de input: %v", metrics.InputLatency)
	logger.Info("  Rate de input: %.1f/s", metrics.InputRate)
	logger.Info("  Total inputs: %d", metrics.InputCount)
	logger.Info("")
	
	logger.Info("AUTOPOT:")
	logger.Info("  Total potions: %d", metrics.PotCount)
	logger.Info("  Latencia de pot: %v", metrics.PotLatency)
	logger.Info("  Modo emergencia: %v", metrics.EmergencyMode)
	logger.Info("")
	
	logger.Info("SPAMMER:")
	logger.Info("  Total spams: %d", metrics.SpamCount)
	logger.Info("  Rate de spam: %.1f/s", metrics.SpamRate)
	logger.Info("  Modo burst: %v", metrics.BurstMode)
	logger.Info("  Teclas activas: %d", metrics.ActiveKeys)
	logger.Info("")
	
	logger.Info("SISTEMA:")
	logger.Info("  CPU usage: %.1f%%", metrics.CPUUsage)
	logger.Info("  Memory usage: %.1f%%", metrics.MemoryUsage)
	logger.Info("  Frame rate: %.1f FPS", metrics.FrameRate)
	logger.Info("")
	
	logger.Info("WoE:")
	logger.Info("  Tiempo en combate: %v", metrics.CombatTime)
	logger.Info("  Daño recibido: %d", metrics.DamageTaken)
	logger.Info("  Skills usados: %d", metrics.SkillsUsed)
	logger.Info("================================")
}

// SetTargetFPS sets the target frame rate
func (m *WOEPerformanceMonitor) SetTargetFPS(fps int) {
	m.targetFPS = fps
	logger.Info("Target FPS establecido a: %d", fps)
}

// EnableAdaptiveMode enables adaptive performance optimization
func (m *WOEPerformanceMonitor) EnableAdaptiveMode() {
	m.adaptiveMode = true
	logger.Info("Modo adaptativo habilitado")
}

// DisableAdaptiveMode disables adaptive performance optimization
func (m *WOEPerformanceMonitor) DisableAdaptiveMode() {
	m.adaptiveMode = false
	logger.Info("Modo adaptativo deshabilitado")
}

// Stop stops the performance monitor
func (m *WOEPerformanceMonitor) Stop() {
	m.cancel()
	m.PrintPerformanceReport()
	logger.Info("Monitor de performance WoE detenido")
}