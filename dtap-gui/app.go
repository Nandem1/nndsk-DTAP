package main

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"syscall"
	"unsafe"

	"dtap-gui/internal/memory"
	"dtap-gui/internal/monitoring"
	"dtap-gui/internal/process"
	"dtap-gui/internal/services"
	"dtap-gui/pkg/config"
	"dtap-gui/pkg/logger"

	"golang.org/x/sys/windows"
)

// App struct
type App struct {
	ctx            context.Context
	serviceManager *services.ServiceManager
	memClient      *memory.Client
	roProcess      *process.Process
	roHandle       windows.Handle
	config         *config.Config
	systemMonitor  *monitoring.SystemMonitor
	mu             sync.RWMutex
	isRunning      bool
	stats          map[string]interface{}
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		stats: make(map[string]interface{}),
	}
}

// startup is called at application startup
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Configurar logger
	logger.Global.SetLevel(logger.INFO)

	// Cargar configuración
	cfg, err := config.LoadConfig(config.GetConfigPath())
	if err != nil {
		logger.Error("Error cargando configuración: %v", err)
		return
	}
	a.config = cfg

	// Inicializar monitoreo del sistema
	a.systemMonitor = monitoring.NewSystemMonitor()

	logger.Info("DTAP GUI iniciado correctamente")
}

// domReady is called after front-end resources have been loaded
func (a *App) domReady(ctx context.Context) {
	// Add your action here
	logger.Info("Frontend cargado, listo para conectar con RO")
}

// beforeClose is called when the application is about to quit
func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	if a.serviceManager != nil && a.isRunning {
		logger.Info("Deteniendo servicios antes de cerrar...")
		a.serviceManager.Stop()
	}
	return false
}

// shutdown is called at application termination
func (a *App) shutdown(ctx context.Context) {
	if a.serviceManager != nil {
		a.serviceManager.Stop()
	}
	if a.memClient != nil {
		a.memClient.Close()
	}
	logger.Info("DTAP GUI cerrado correctamente")
}

// ConnectToRO connects to Ragnarok Online process
func (a *App) ConnectToRO() map[string]interface{} {
	a.mu.Lock()
	defer a.mu.Unlock()

	result := map[string]interface{}{
		"success": false,
		"message": "",
		"data":    nil,
	}

	// Buscar proceso de RO
	logger.Info("Buscando proceso de Ragnarok Online...")
	roProcess, err := process.FindROProcess()
	if err != nil {
		result["message"] = fmt.Sprintf("Error: %v. Asegúrate de que Ragnarok Online esté ejecutándose.", err)
		return result
	}

	a.roProcess = roProcess
	logger.Info("Proceso encontrado: %s (PID: %d)", roProcess.Name, roProcess.PID)

	// Obtener handle de la ventana
	logger.Info("Obteniendo handle de la ventana...")
	roHandle, err := a.getWindowHandle(roProcess.PID)
	if err != nil {
		result["message"] = fmt.Sprintf("Error obteniendo ventana: %v", err)
		return result
	}
	a.roHandle = roHandle
	logger.Info("Handle de ventana obtenido: 0x%x", roHandle)

	// Crear cliente de memoria
	logger.Info("Conectando a la memoria del proceso...")
	memClient, err := memory.NewClient(roProcess.PID)
	if err != nil {
		result["message"] = fmt.Sprintf("Error conectando a memoria: %v", err)
		return result
	}
	a.memClient = memClient
	logger.Info("Conectado a la memoria del proceso")

	// Inicializar ServiceManager
	a.serviceManager = services.NewServiceManager(a.config, memClient, roHandle)

	result["success"] = true
	result["message"] = "Conectado exitosamente a Ragnarok Online"
	result["data"] = map[string]interface{}{
		"process_name":  roProcess.Name,
		"pid":           roProcess.PID,
		"window_handle": fmt.Sprintf("0x%x", roHandle),
	}

	return result
}

// StartServices starts all DTAP services
func (a *App) StartServices() map[string]interface{} {
	a.mu.Lock()
	defer a.mu.Unlock()

	result := map[string]interface{}{
		"success": false,
		"message": "",
	}

	if a.serviceManager == nil {
		result["message"] = "No hay conexión con Ragnarok Online. Conecta primero."
		return result
	}

	if a.isRunning {
		result["message"] = "Los servicios ya están ejecutándose"
		return result
	}

	// Iniciar ServiceManager
	if err := a.serviceManager.Start(); err != nil {
		result["message"] = fmt.Sprintf("Error iniciando servicios: %v", err)
		return result
	}

	a.isRunning = true
	result["success"] = true
	result["message"] = "Servicios iniciados correctamente"

	return result
}

// StopServices stops all DTAP services
func (a *App) StopServices() map[string]interface{} {
	a.mu.Lock()
	defer a.mu.Unlock()

	result := map[string]interface{}{
		"success": false,
		"message": "",
	}

	if !a.isRunning {
		result["message"] = "Los servicios no están ejecutándose"
		return result
	}

	if a.serviceManager != nil {
		if err := a.serviceManager.Stop(); err != nil {
			result["message"] = fmt.Sprintf("Error deteniendo servicios: %v", err)
			return result
		}
	}

	a.isRunning = false
	result["success"] = true
	result["message"] = "Servicios detenidos correctamente"

	return result
}

// GetStats returns current statistics
func (a *App) GetStats() map[string]interface{} {
	a.mu.RLock()
	defer a.mu.RUnlock()

	stats := map[string]interface{}{
		"is_running": a.isRunning,
		"connected":  a.memClient != nil,
	}

	if a.memClient != nil {
		// Leer HP/SP
		hp, maxHP, err := a.memClient.ReadHP()
		if err == nil {
			stats["hp"] = hp
			stats["max_hp"] = maxHP
			stats["hp_percent"] = int((float64(hp) / float64(maxHP)) * 100)
		}

		sp, maxSP, err := a.memClient.ReadSP()
		if err == nil {
			stats["sp"] = sp
			stats["max_sp"] = maxSP
			stats["sp_percent"] = int((float64(sp) / float64(maxSP)) * 100)
		}
	}

	if a.serviceManager != nil {
		serviceStats := a.serviceManager.GetStats()
		for k, v := range serviceStats {
			stats[k] = v
		}
	}

	if a.systemMonitor != nil {
		systemStats := a.systemMonitor.GetSystemStats()
		stats["system"] = systemStats
		stats["uptime"] = a.systemMonitor.GetUptime().Seconds()
	}

	return stats
}

// GetConfig returns current configuration
func (a *App) GetConfig() map[string]interface{} {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if a.config == nil {
		return map[string]interface{}{
			"error": "Configuración no cargada",
		}
	}

	// Convertir configuración a map para JSON
	configData := map[string]interface{}{
		"autopot": map[string]interface{}{
			"enabled":     a.config.AutoPot.Enabled,
			"hp_pot_key":  a.config.AutoPot.HPPotKey,
			"sp_pot_key":  a.config.AutoPot.SPPotKey,
			"hp_percent":  a.config.AutoPot.HPPercent,
			"sp_percent":  a.config.AutoPot.SPPercent,
			"check_delay": a.config.AutoPot.CheckDelay,
		},
		"autobuff": map[string]interface{}{
			"enabled":   a.config.AutoBuff.Enabled,
			"delay":     a.config.AutoBuff.Delay,
			"gloom_key": a.config.AutoBuff.GloomKey,
			"aspd_key":  a.config.AutoBuff.AspdKey,
		},
		"spammer": map[string]interface{}{
			"enabled":      a.config.Spammer.Enabled,
			"keys":         a.config.Spammer.Keys,
			"delay":        a.config.Spammer.Delay,
			"skill_cancel": a.config.Spammer.SkillCancel,
			"double_click": a.config.Spammer.DoubleClick,
			"mouse_click":  a.config.Spammer.MouseClick,
			"mode":         a.config.Spammer.Mode,
		},
		"ui": map[string]interface{}{
			"show_debug": a.config.UI.ShowDebug,
			"log_level":  a.config.UI.LogLevel,
		},
	}

	return configData
}

// UpdateConfig updates the configuration
func (a *App) UpdateConfig(configJSON string) map[string]interface{} {
	a.mu.Lock()
	defer a.mu.Unlock()

	result := map[string]interface{}{
		"success": false,
		"message": "",
	}

	var newConfig config.Config
	if err := json.Unmarshal([]byte(configJSON), &newConfig); err != nil {
		result["message"] = fmt.Sprintf("Error parseando configuración: %v", err)
		return result
	}

	// Actualizar configuración en memoria (sin guardar JSON)
	a.config = &newConfig

	// Actualizar configuración en ServiceManager si está ejecutándose
	if a.serviceManager != nil {
		a.serviceManager.UpdateConfig(&newConfig)
	}

	result["success"] = true
	result["message"] = "Configuración actualizada correctamente"

	return result
}

// getWindowHandle obtiene el handle de la ventana del proceso
func (a *App) getWindowHandle(pid uint32) (windows.Handle, error) {
	var hwnd windows.Handle

	// Callback para enumerar ventanas
	enumProc := syscall.NewCallback(func(h syscall.Handle, lParam uintptr) uintptr {
		var processID uint32
		windows.GetWindowThreadProcessId(windows.HWND(h), &processID)

		if processID == pid {
			hwnd = windows.Handle(h)
			return 0 // Parar enumeración
		}
		return 1 // Continuar enumeración
	})

	windows.EnumWindows(enumProc, unsafe.Pointer(uintptr(0)))

	if hwnd == 0 {
		return 0, fmt.Errorf("no se encontró ventana para el proceso")
	}

	return hwnd, nil
}
