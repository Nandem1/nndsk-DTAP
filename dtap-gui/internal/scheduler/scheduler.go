package scheduler

import (
	"context"
	"sync"
	"time"

	"dtap-gui/internal/metrics"
	"dtap-gui/pkg/config"
	"dtap-gui/pkg/constants"
	"dtap-gui/pkg/logger"
)

// Priority define la prioridad de una tarea
type Priority int

const (
	PriorityCritical Priority = iota // AutoPot - supervivencia
	PriorityHigh                     // Spammer - combate
	PriorityNormal                   // AutoBuff - soporte
	PriorityLow                      // Tareas de mantenimiento
)

// Task representa una tarea programable
type Task struct {
	ID          string
	Priority    Priority
	Interval    time.Duration
	LastRun     time.Time
	NextRun     time.Time
	Handler     func(ctx context.Context) error
	Enabled     bool
	Description string
}

// Scheduler maneja la ejecución de tareas con prioridades
type Scheduler struct {
	mu      sync.RWMutex
	tasks   map[string]*Task
	ctx     context.Context
	cancel  context.CancelFunc
	running bool
	paused  bool // Nuevo: estado de pausa
	metrics *metrics.Metrics

	// Configuración
	baseInterval time.Duration
}

// NewScheduler crea un nuevo scheduler
func NewScheduler() *Scheduler {
	ctx, cancel := context.WithCancel(context.Background())
	return &Scheduler{
		tasks:        make(map[string]*Task),
		ctx:          ctx,
		cancel:       cancel,
		baseInterval: constants.SchedulerBaseInterval,
	}
}

// AddTask agrega una tarea al scheduler
func (s *Scheduler) AddTask(task *Task) {
	s.mu.Lock()
	defer s.mu.Unlock()

	task.NextRun = time.Now().Add(task.Interval)
	s.tasks[task.ID] = task

	logger.Info("📅 Tarea agregada: %s (prioridad: %d, intervalo: %v)",
		task.Description, task.Priority, task.Interval)
}

// RemoveTask remueve una tarea del scheduler
func (s *Scheduler) RemoveTask(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if task, exists := s.tasks[id]; exists {
		delete(s.tasks, id)
		logger.Info("📅 Tarea removida: %s", task.Description)
	}
}

// EnableTask habilita una tarea
func (s *Scheduler) EnableTask(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if task, exists := s.tasks[id]; exists {
		task.Enabled = true
		logger.Info("📅 Tarea habilitada: %s", task.Description)
	}
}

// DisableTask deshabilita una tarea
func (s *Scheduler) DisableTask(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if task, exists := s.tasks[id]; exists {
		task.Enabled = false
		logger.Info("📅 Tarea deshabilitada: %s", task.Description)
	}
}

// Start inicia el scheduler
func (s *Scheduler) Start() {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return
	}
	s.running = true
	s.mu.Unlock()

	logger.Info("🚀 Scheduler iniciado con %d tareas", len(s.tasks))

	go s.run()
}

// Pause pausa el scheduler sin detenerlo completamente
func (s *Scheduler) Pause() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running || s.paused {
		return
	}

	s.paused = true
	logger.Info("⏸️ Scheduler pausado")
}

// Resume reanuda el scheduler pausado
func (s *Scheduler) Resume() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running || !s.paused {
		return
	}

	s.paused = false
	logger.Info("▶️ Scheduler reanudado")
}

// Stop detiene el scheduler
func (s *Scheduler) Stop() {
	s.mu.Lock()
	if !s.running {
		s.mu.Unlock()
		return
	}
	s.running = false
	s.paused = false
	s.mu.Unlock()

	s.cancel()

	// Crear nuevo contexto para permitir reinicio
	s.ctx, s.cancel = context.WithCancel(context.Background())

	logger.Info("🛑 Scheduler detenido")
}

// UpdateConfig actualiza la configuración del scheduler
func (s *Scheduler) UpdateConfig(cfg *config.Config) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Actualizar tareas según la nueva configuración
	if autopotTask, exists := s.tasks["autopot"]; exists {
		autopotTask.Enabled = cfg.AutoPot.Enabled
	}

	if autobuffTask, exists := s.tasks["autobuff"]; exists {
		autobuffTask.Enabled = cfg.AutoBuff.Enabled
	}

	if spammerTask, exists := s.tasks["spammer"]; exists {
		spammerTask.Enabled = cfg.Spammer.Enabled
	}

	logger.Info("🔄 Scheduler configurado: autopot=%v, autobuff=%v, spammer=%v",
		cfg.AutoPot.Enabled, cfg.AutoBuff.Enabled, cfg.Spammer.Enabled)
}

// run es el loop principal del scheduler
func (s *Scheduler) run() {
	ticker := time.NewTicker(s.baseInterval)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			// Verificar si está pausado
			s.mu.RLock()
			paused := s.paused
			s.mu.RUnlock()

			if !paused {
				s.processTasks()
			}
		}
	}
}

// processTasks procesa todas las tareas pendientes
func (s *Scheduler) processTasks() {
	s.mu.RLock()
	now := time.Now()

	// Crear lista de tareas ordenadas por prioridad
	var tasksToRun []*Task
	for _, task := range s.tasks {
		if task.Enabled && now.After(task.NextRun) {
			tasksToRun = append(tasksToRun, task)
		}
	}
	s.mu.RUnlock()

	// Ordenar por prioridad (menor número = mayor prioridad)
	for i := 0; i < len(tasksToRun)-1; i++ {
		for j := i + 1; j < len(tasksToRun); j++ {
			if tasksToRun[i].Priority > tasksToRun[j].Priority {
				tasksToRun[i], tasksToRun[j] = tasksToRun[j], tasksToRun[i]
			}
		}
	}

	// Ejecutar tareas en orden de prioridad
	for _, task := range tasksToRun {
		s.executeTask(task)
	}
}

// executeTask ejecuta una tarea individual
func (s *Scheduler) executeTask(task *Task) {
	start := time.Now()

	// Ejecutar la tarea
	if err := task.Handler(s.ctx); err != nil {
		logger.Error("Error ejecutando tarea %s: %v", task.ID, err)
	}

	// Actualizar tiempos
	task.LastRun = start
	task.NextRun = start.Add(task.Interval)

	// Reportar métricas si están disponibles
	if s.metrics != nil {
		s.metrics.IncrementSpammerLoops() // Esto se puede refinar por tipo de tarea
	}
}

// GetStats retorna estadísticas del scheduler
func (s *Scheduler) GetStats() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	stats := map[string]interface{}{
		"running": s.running,
		"tasks":   len(s.tasks),
	}

	// Contar tareas por prioridad
	priorityCount := make(map[Priority]int)
	for _, task := range s.tasks {
		priorityCount[task.Priority]++
	}

	stats["priority_counts"] = map[string]int{
		"critical": priorityCount[PriorityCritical],
		"high":     priorityCount[PriorityHigh],
		"normal":   priorityCount[PriorityNormal],
		"low":      priorityCount[PriorityLow],
	}

	return stats
}

// SetMetrics establece la instancia de métricas
func (s *Scheduler) SetMetrics(m *metrics.Metrics) {
	s.metrics = m
}

// GetTask retorna una tarea por ID
func (s *Scheduler) GetTask(id string) (*Task, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, exists := s.tasks[id]
	return task, exists
}

// ListTasks retorna todas las tareas
func (s *Scheduler) ListTasks() []*Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var tasks []*Task
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}

	return tasks
}
