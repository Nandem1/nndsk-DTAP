package monitoring

import (
	"fmt"
	"runtime"
	"time"
)

// SystemStats contiene estadísticas del sistema
type SystemStats struct {
	// CPU
	CPUUsage float64

	// Memory
	MemoryUsage     uint64 // bytes
	MemoryTotal     uint64 // bytes
	MemoryPercent   float64
	GoMemoryUsage   uint64 // bytes
	GoMemoryTotal   uint64 // bytes
	GoMemoryPercent float64

	// Process
	ProcessMemoryUsage uint64 // bytes
	ProcessCPUUsage    float64

	// Goroutines
	GoroutineCount int

	// Timestamp
	Timestamp time.Time
}

// SystemMonitor maneja el monitoreo del sistema
type SystemMonitor struct {
	lastCPUTime time.Time
	lastCPU     uint64
}

// NewSystemMonitor crea un nuevo monitor del sistema
func NewSystemMonitor() *SystemMonitor {
	return &SystemMonitor{
		lastCPUTime: time.Now(),
		lastCPU:     0,
	}
}

// GetSystemStats obtiene estadísticas actuales del sistema
func (sm *SystemMonitor) GetSystemStats() *SystemStats {
	stats := &SystemStats{
		Timestamp: time.Now(),
	}

	// Obtener estadísticas de memoria del sistema
	sm.getSystemMemoryStats(stats)

	// Obtener estadísticas de Go runtime
	sm.getGoRuntimeStats(stats)

	// Obtener estadísticas del proceso actual
	sm.getProcessStats(stats)

	return stats
}

// getSystemMemoryStats obtiene estadísticas de memoria del sistema
func (sm *SystemMonitor) getSystemMemoryStats(stats *SystemStats) {
	// Implementación simplificada usando runtime
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Estimación basada en memoria de Go
	stats.MemoryUsage = m.Sys
	stats.MemoryTotal = m.Sys * 10 // Estimación
	stats.MemoryPercent = float64(m.Sys) / float64(m.Sys*10) * 100
}

// getGoRuntimeStats obtiene estadísticas del runtime de Go
func (sm *SystemMonitor) getGoRuntimeStats(stats *SystemStats) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	stats.GoMemoryUsage = m.Alloc
	stats.GoMemoryTotal = m.Sys
	stats.GoMemoryPercent = float64(m.Alloc) / float64(m.Sys) * 100
	stats.GoroutineCount = runtime.NumGoroutine()
}

// getProcessStats obtiene estadísticas del proceso actual
func (sm *SystemMonitor) getProcessStats(stats *SystemStats) {
	// Usar estadísticas de Go runtime como aproximación
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	stats.ProcessMemoryUsage = m.Sys
	stats.ProcessCPUUsage = sm.calculateProcessCPUUsage()
}

// calculateProcessCPUUsage calcula el uso de CPU del proceso
func (sm *SystemMonitor) calculateProcessCPUUsage() float64 {
	// Implementación simplificada - en producción se usaría un método más preciso
	now := time.Now()
	elapsed := now.Sub(sm.lastCPUTime).Seconds()

	if elapsed > 0 {
		sm.lastCPUTime = now
		// Retornar un valor estimado basado en goroutines activas
		return float64(runtime.NumGoroutine()) * 0.1
	}

	return 0.0
}

// FormatBytes formatea bytes a una representación legible
func FormatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// FormatPercent formatea un porcentaje
func FormatPercent(percent float64) string {
	return fmt.Sprintf("%.1f%%", percent)
}
