package prediction

import (
	"sync"
	"time"
)

// DataPoint representa un punto de datos en el tiempo
type DataPoint struct {
	Timestamp time.Time
	Value     float64
}

// RingBuffer implementa un buffer circular para datos históricos
type RingBuffer struct {
	mu       sync.RWMutex
	data     []DataPoint
	capacity int
	head     int
	tail     int
	count    int
}

// NewRingBuffer crea un nuevo ring buffer
func NewRingBuffer(capacity int) *RingBuffer {
	return &RingBuffer{
		data:     make([]DataPoint, capacity),
		capacity: capacity,
		head:     0,
		tail:     0,
		count:    0,
	}
}

// Add agrega un nuevo punto de datos
func (rb *RingBuffer) Add(value float64) {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	rb.data[rb.tail] = DataPoint{
		Timestamp: time.Now(),
		Value:     value,
	}

	rb.tail = (rb.tail + 1) % rb.capacity

	if rb.count < rb.capacity {
		rb.count++
	} else {
		rb.head = (rb.head + 1) % rb.capacity
	}
}

// GetLatest retorna el valor más reciente
func (rb *RingBuffer) GetLatest() (float64, bool) {
	rb.mu.RLock()
	defer rb.mu.RUnlock()

	if rb.count == 0 {
		return 0, false
	}

	latestIndex := (rb.tail - 1 + rb.capacity) % rb.capacity
	return rb.data[latestIndex].Value, true
}

// GetHistory retorna los últimos N puntos de datos
func (rb *RingBuffer) GetHistory(n int) []DataPoint {
	rb.mu.RLock()
	defer rb.mu.RUnlock()

	if rb.count == 0 {
		return nil
	}

	if n > rb.count {
		n = rb.count
	}

	result := make([]DataPoint, n)
	start := (rb.tail - n + rb.capacity) % rb.capacity

	for i := 0; i < n; i++ {
		index := (start + i) % rb.capacity
		result[i] = rb.data[index]
	}

	return result
}

// GetCount retorna el número de elementos en el buffer
func (rb *RingBuffer) GetCount() int {
	rb.mu.RLock()
	defer rb.mu.RUnlock()
	return rb.count
}

// Predictor maneja la predicción de valores basada en historial
type Predictor struct {
	hpHistory *RingBuffer
	spHistory *RingBuffer
	mu        sync.RWMutex

	// Configuración
	predictionWindow time.Duration
	minDataPoints    int
}

// NewPredictor crea un nuevo predictor
func NewPredictor(capacity int, predictionWindow time.Duration) *Predictor {
	return &Predictor{
		hpHistory:        NewRingBuffer(capacity),
		spHistory:        NewRingBuffer(capacity),
		predictionWindow: predictionWindow,
		minDataPoints:    3, // Mínimo para predicción confiable
	}
}

// AddHP agrega un punto de datos de HP
func (p *Predictor) AddHP(hpPercent float64) {
	p.hpHistory.Add(hpPercent)
}

// AddSP agrega un punto de datos de SP
func (p *Predictor) AddSP(spPercent float64) {
	p.spHistory.Add(spPercent)
}

// PredictHP predice el HP en el futuro
func (p *Predictor) PredictHP(futureTime time.Duration) (float64, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	history := p.hpHistory.GetHistory(p.minDataPoints)
	if len(history) < p.minDataPoints {
		return 0, false
	}

	// Calcular tendencia simple (regresión lineal básica)
	trend := p.calculateTrend(history)

	// Predecir valor futuro
	latest, _ := p.hpHistory.GetLatest()
	predicted := latest + (trend * float64(futureTime.Milliseconds()) / 1000.0)

	// Limitar predicción a valores razonables
	if predicted < 0 {
		predicted = 0
	} else if predicted > 100 {
		predicted = 100
	}

	return predicted, true
}

// PredictSP predice el SP en el futuro
func (p *Predictor) PredictSP(futureTime time.Duration) (float64, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	history := p.spHistory.GetHistory(p.minDataPoints)
	if len(history) < p.minDataPoints {
		return 0, false
	}

	// Calcular tendencia simple
	trend := p.calculateTrend(history)

	// Predecir valor futuro
	latest, _ := p.spHistory.GetLatest()
	predicted := latest + (trend * float64(futureTime.Milliseconds()) / 1000.0)

	// Limitar predicción a valores razonables
	if predicted < 0 {
		predicted = 0
	} else if predicted > 100 {
		predicted = 100
	}

	return predicted, true
}

// calculateTrend calcula la tendencia de los datos históricos
func (p *Predictor) calculateTrend(history []DataPoint) float64 {
	if len(history) < 2 {
		return 0
	}

	// Regresión lineal simple: y = mx + b
	// m = (n*Σxy - Σx*Σy) / (n*Σx² - (Σx)²)

	n := float64(len(history))
	var sumX, sumY, sumXY, sumXX float64

	for i, point := range history {
		x := float64(i)
		y := point.Value

		sumX += x
		sumY += y
		sumXY += x * y
		sumXX += x * x
	}

	denominator := n*sumXX - sumX*sumX
	if denominator == 0 {
		return 0
	}

	slope := (n*sumXY - sumX*sumY) / denominator

	// Convertir slope a cambio por segundo
	// Asumiendo que los puntos están espaciados uniformemente
	if len(history) > 1 {
		timeDiff := history[len(history)-1].Timestamp.Sub(history[0].Timestamp)
		if timeDiff > 0 {
			slope = slope / timeDiff.Seconds()
		}
	}

	return slope
}

// GetHPStats retorna estadísticas del HP
func (p *Predictor) GetHPStats() map[string]interface{} {
	p.mu.RLock()
	defer p.mu.RUnlock()

	history := p.hpHistory.GetHistory(10) // Últimos 10 puntos

	stats := map[string]interface{}{
		"count":  p.hpHistory.GetCount(),
		"latest": 0.0,
		"trend":  0.0,
	}

	if len(history) > 0 {
		stats["latest"] = history[len(history)-1].Value

		if len(history) >= 2 {
			stats["trend"] = p.calculateTrend(history)
		}
	}

	return stats
}

// GetSPStats retorna estadísticas del SP
func (p *Predictor) GetSPStats() map[string]interface{} {
	p.mu.RLock()
	defer p.mu.RUnlock()

	history := p.spHistory.GetHistory(10) // Últimos 10 puntos

	stats := map[string]interface{}{
		"count":  p.spHistory.GetCount(),
		"latest": 0.0,
		"trend":  0.0,
	}

	if len(history) > 0 {
		stats["latest"] = history[len(history)-1].Value

		if len(history) >= 2 {
			stats["trend"] = p.calculateTrend(history)
		}
	}

	return stats
}

// Clear limpia todos los datos históricos
func (p *Predictor) Clear() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.hpHistory = NewRingBuffer(p.hpHistory.capacity)
	p.spHistory = NewRingBuffer(p.spHistory.capacity)
}
