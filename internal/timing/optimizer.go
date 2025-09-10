package timing

import (
	"sync"
	"syscall"
	"time"

	"dtap-go/pkg/logger"
)

// Optimizer maneja la optimización de timing de Windows
type Optimizer struct {
	mu              sync.RWMutex
	enabled         bool
	winmm           *syscall.LazyDLL
	timeBeginPeriod *syscall.LazyProc
	timeEndPeriod   *syscall.LazyProc
}

var (
	instance *Optimizer
	once     sync.Once
)

// GetOptimizer retorna la instancia singleton del optimizer
func GetOptimizer() *Optimizer {
	once.Do(func() {
		instance = &Optimizer{
			winmm: syscall.NewLazyDLL("winmm.dll"),
		}
		instance.timeBeginPeriod = instance.winmm.NewProc("timeBeginPeriod")
		instance.timeEndPeriod = instance.winmm.NewProc("timeEndPeriod")
	})
	return instance
}

// Enable activa la optimización de timing de Windows
func (o *Optimizer) Enable() error {
	o.mu.Lock()
	defer o.mu.Unlock()

	if o.enabled {
		return nil // Ya está habilitado
	}

	// Llamar a timeBeginPeriod(1) para máxima precisión
	ret, _, err := o.timeBeginPeriod.Call(1)
	if ret != 0 {
		logger.Error("Error al activar timeBeginPeriod: %v", err)
		return err
	}

	o.enabled = true
	logger.Info("⚡ Optimización de timing de Windows activada (timeBeginPeriod=1ms)")

	return nil
}

// Disable desactiva la optimización de timing de Windows
func (o *Optimizer) Disable() error {
	o.mu.Lock()
	defer o.mu.Unlock()

	if !o.enabled {
		return nil // Ya está deshabilitado
	}

	// Llamar a timeEndPeriod(1) para restaurar configuración por defecto
	ret, _, err := o.timeEndPeriod.Call(1)
	if ret != 0 {
		logger.Error("Error al desactivar timeEndPeriod: %v", err)
		return err
	}

	o.enabled = false
	logger.Info("🔄 Optimización de timing de Windows desactivada")

	return nil
}

// IsEnabled retorna si la optimización está activa
func (o *Optimizer) IsEnabled() bool {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.enabled
}

// GetStats retorna estadísticas del optimizer
func (o *Optimizer) GetStats() map[string]interface{} {
	o.mu.RLock()
	defer o.mu.RUnlock()

	return map[string]interface{}{
		"enabled": o.enabled,
		"period":  1, // 1ms cuando está habilitado
	}
}

// BenchmarkTiming realiza un benchmark de precisión de timing
func (o *Optimizer) BenchmarkTiming() map[string]interface{} {
	const iterations = 1000
	const targetDelay = 1 * time.Millisecond

	var totalError time.Duration
	var maxError time.Duration
	var minError time.Duration = time.Hour

	for i := 0; i < iterations; i++ {
		start := time.Now()
		time.Sleep(targetDelay)
		actual := time.Since(start)

		error := actual - targetDelay
		if error < 0 {
			error = -error
		}

		totalError += error
		if error > maxError {
			maxError = error
		}
		if error < minError {
			minError = error
		}
	}

	avgError := totalError / iterations

	return map[string]interface{}{
		"iterations":   iterations,
		"target_delay": targetDelay.String(),
		"avg_error":    avgError.String(),
		"max_error":    maxError.String(),
		"min_error":    minError.String(),
		"avg_error_us": avgError.Microseconds(),
		"max_error_us": maxError.Microseconds(),
		"min_error_us": minError.Microseconds(),
	}
}

// AutoDisable desactiva automáticamente la optimización al finalizar
func (o *Optimizer) AutoDisable() {
	if o.IsEnabled() {
		o.Disable()
	}
}
