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

// SearchJobIDInMemory - Busca el job ID en toda la memoria del proceso
func (a *App) SearchJobIDInMemory() map[string]interface{} {
	result := map[string]interface{}{
		"success": false,
		"message": "",
		"data":    nil,
	}

	if a.memClient == nil {
		result["message"] = "No hay conexión con la memoria"
		return result
	}

	// Leer el job ID actual para buscarlo
	currentJobID, err := a.memClient.ReadJobID()
	if err != nil {
		result["message"] = "Error leyendo job ID actual"
		return result
	}

	// Buscar en diferentes rangos de memoria
	searchResults := make(map[string]interface{})

	// Rango 1: Alrededor de la dirección del personaje (más probable)
	baseAddr := uintptr(0x10df5d8)
	for i := -50; i <= 200; i++ {
		addr := baseAddr + uintptr(i*4)
		value, err := a.memClient.ReadUint32(addr)
		if err == nil && value == currentJobID {
			searchResults[fmt.Sprintf("near_character_0x%x", addr)] = map[string]interface{}{
				"address":          fmt.Sprintf("0x%x", addr),
				"value":            value,
				"job_name":         getJobName(value),
				"offset_from_base": i * 4,
			}
		}
	}

	// Rango 2: Alrededor de las direcciones de HP/SP
	hpAddr := uintptr(0x10DCE10)
	for i := -100; i <= 100; i++ {
		addr := hpAddr + uintptr(i*4)
		value, err := a.memClient.ReadUint32(addr)
		if err == nil && value == currentJobID {
			searchResults[fmt.Sprintf("near_hp_0x%x", addr)] = map[string]interface{}{
				"address":        fmt.Sprintf("0x%x", addr),
				"value":          value,
				"job_name":       getJobName(value),
				"offset_from_hp": i * 4,
			}
		}
	}

	// Rango 3: Búsqueda más amplia en el heap
	for i := 0; i < 1000; i++ {
		addr := uintptr(0x10000000 + uintptr(i*4)) // Rango del heap
		value, err := a.memClient.ReadUint32(addr)
		if err == nil && value == currentJobID {
			searchResults[fmt.Sprintf("heap_0x%x", addr)] = map[string]interface{}{
				"address":     fmt.Sprintf("0x%x", addr),
				"value":       value,
				"job_name":    getJobName(value),
				"heap_offset": i * 4,
			}
		}
	}

	result["success"] = true
	result["message"] = fmt.Sprintf("Búsqueda completada. Job ID actual: %d (%s)", currentJobID, getJobName(currentJobID))
	result["data"] = searchResults

	return result
}

// DebugCharacterInfo - Función para debuggear información del personaje
func (a *App) DebugCharacterInfo() map[string]interface{} {
	result := map[string]interface{}{
		"success": false,
		"message": "",
		"data":    nil,
	}

	if a.memClient == nil {
		result["message"] = "No hay conexión con la memoria"
		return result
	}

	// Leer diferentes offsets alrededor de la dirección del personaje
	baseAddr := uintptr(0x10df5d8)
	debugData := make(map[string]interface{})

	// Leer en diferentes offsets para encontrar el job ID
	for i := 0; i < 20; i++ {
		offset := i * 4 // Cada offset de 4 bytes
		addr := baseAddr + uintptr(offset)

		value, err := a.memClient.ReadUint32(addr)
		if err == nil {
			debugData[fmt.Sprintf("offset_%d_0x%x", offset, addr)] = map[string]interface{}{
				"value":    value,
				"job_name": getJobName(value),
			}
		}
	}

	// También leer el nombre del personaje para verificar
	characterName, _ := a.memClient.ReadCharacterName()
	debugData["character_name"] = characterName

	result["success"] = true
	result["message"] = "Debug info obtenida"
	result["data"] = debugData

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

		// Leer información del personaje
		characterName, err := a.memClient.ReadCharacterName()
		if err == nil && characterName != "" {
			stats["character_name"] = characterName
		}

		mapName, err := a.memClient.ReadMapName()
		if err == nil && mapName != "" {
			stats["map_name"] = mapName
		}

		jobID, err := a.memClient.ReadJobID()
		if err == nil {
			stats["job_id"] = jobID
			stats["job_name"] = getJobName(jobID)
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

// getJobName convierte el ID del job a un nombre legible
func getJobName(jobID uint32) string {
	jobNames := map[uint32]string{
		0:    "Novice",
		1:    "Swordman",
		2:    "Magician",
		3:    "Archer",
		4:    "Acolyte",
		5:    "Merchant",
		6:    "Thief",
		7:    "Knight",
		8:    "Priest",
		9:    "Wizard",
		10:   "Blacksmith",
		11:   "Hunter",
		12:   "Assassin",
		14:   "Crusader",
		15:   "Monk",
		16:   "Sage",
		17:   "Rogue",
		18:   "Alchemist",
		19:   "Bard",
		20:   "Dancer",
		23:   "Super Novice",
		24:   "Gunslinger",
		25:   "Ninja",
		4001: "High Novice",
		4002: "High Swordman",
		4003: "High Magician",
		4004: "High Archer",
		4005: "High Acolyte",
		4006: "High Merchant",
		4007: "High Thief",
		4008: "Lord Knight",
		4009: "High Priest",
		4010: "High Wizard",
		4011: "Whitesmith",
		4012: "Sniper",
		4013: "Assassin Cross",
		4015: "Paladin",
		4016: "Champion",
		4017: "Professor",
		4018: "Stalker",
		4019: "Creator",
		4020: "Clown",
		4021: "Gypsy",
		4023: "Baby Novice",
		4024: "Baby Swordman",
		4025: "Baby Magician",
		4026: "Baby Archer",
		4027: "Baby Acolyte",
		4028: "Baby Merchant",
		4029: "Baby Thief",
		4030: "Baby Knight",
		4031: "Baby Priest",
		4032: "Baby Wizard",
		4033: "Baby Blacksmith",
		4034: "Baby Hunter",
		4035: "Baby Assassin",
		4037: "Baby Crusader",
		4038: "Baby Monk",
		4039: "Baby Sage",
		4040: "Baby Rogue",
		4041: "Baby Alchemist",
		4042: "Baby Bard",
		4043: "Baby Dancer",
		4045: "Super Baby",
		4046: "Taekwon",
		4047: "Star Gladiator",
		4049: "Soul Linker",
		4050: "Gangsi",
		4051: "Death Knight",
		4052: "Dark Collector",
		4054: "Rune Knight",
		4055: "Warlock",
		4056: "Ranger",
		4057: "Arch Bishop",
		4058: "Mechanic",
		4059: "Guillotine Cross",
		4060: "Rune Knight2",
		4061: "Warlock2",
		4062: "Ranger2",
		4063: "Arch Bishop2",
		4064: "Mechanic2",
		4065: "Guillotine2",
		4066: "Royal Guard",
		4067: "Sorcerer",
		4068: "Minstrel",
		4069: "Wanderer",
		4070: "Sura",
		4071: "Genetic",
		4072: "Shadow Chaser",
		4073: "Royal Guard2",
		4074: "Sorcerer2",
		4075: "Minstrel2",
		4076: "Wanderer2",
		4077: "Sura2",
		4078: "Genetic2",
		4079: "Shadow Chaser2",
		4096: "Baby Rune",
		4097: "Baby Warlock",
		4098: "Baby Ranger",
		4099: "Baby Bishop",
		4100: "Baby Mechanic",
		4101: "Baby Cross",
		4102: "Baby Guard",
		4103: "Baby Sorcerer",
		4104: "Baby Minstrel",
		4105: "Baby Wanderer",
		4106: "Baby Sura",
		4107: "Baby Genetic",
		4108: "Baby Chaser",
		4190: "Super Novice E",
		4191: "Super Baby E",
		4211: "Kagerou",
		4212: "Oboro",
	}

	if name, exists := jobNames[jobID]; exists {
		return name
	}
	return fmt.Sprintf("Unknown Job (%d)", jobID)
}

// MonitorJobAddresses - Monitorea las direcciones específicas encontradas
func (a *App) MonitorJobAddresses() map[string]interface{} {
	result := map[string]interface{}{"success": false, "message": "", "data": nil}
	if a.memClient == nil {
		result["message"] = "No hay conexión con la memoria"
		return result
	}

	// Las 7 direcciones encontradas para Creator (4019)
	addresses := []uintptr{
		0x14A0208C,
		0x14A1A534,
		0x153F1FC0,
		0x153F1FC4,
		0x1FE5B16C,
		0x1FE5B170,
		0x1FE5B174,
	}

	monitorResults := make(map[string]interface{})

	for i, addr := range addresses {
		value, err := a.memClient.ReadUint32(addr)
		if err == nil {
			monitorResults[fmt.Sprintf("addr_%d_0x%x", i+1, addr)] = map[string]interface{}{
				"address":    fmt.Sprintf("0x%x", addr),
				"value":      value,
				"hex_value":  fmt.Sprintf("0x%X", value),
				"job_name":   getJobName(value),
				"is_changed": value != 4019, // Marcar si cambió de Creator
			}
		} else {
			monitorResults[fmt.Sprintf("addr_%d_0x%x", i+1, addr)] = map[string]interface{}{
				"address": fmt.Sprintf("0x%x", addr),
				"error":   "No se pudo leer",
			}
		}
	}

	result["success"] = true
	result["message"] = "Monitoreo de direcciones completado"
	result["data"] = monitorResults
	return result
}
