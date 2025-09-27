package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
	"unsafe"

	"dtap-go/internal/memory"
	"dtap-go/internal/metrics"
	"dtap-go/internal/monitoring"
	"dtap-go/internal/process"
	"dtap-go/internal/services"
	"dtap-go/internal/timing"
	"dtap-go/pkg/config"
	"dtap-go/pkg/constants"
	"dtap-go/pkg/logger"

	"golang.org/x/sys/windows"
)

func main() {
	// Cargar configuración
	cfg, err := config.LoadConfig(config.GetConfigPath())
	if err != nil {
		logger.Error("Error cargando configuración: %v", err)
		os.Exit(1)
	}

	// Configurar logger
	logger.Global.SetLevel(logger.ParseLogLevel(cfg.UI.LogLevel))

	// Inicializar métricas si están habilitadas
	if cfg.Metrics.Enabled {
		metricsInstance := metrics.GetMetrics()
		if err := metricsInstance.Enable(cfg.Metrics.Port); err != nil {
			logger.Error("Error habilitando métricas: %v", err)
		}
	}

	// Inicializar timing optimizer si está habilitado
	if cfg.Timing.Enabled {
		timingOptimizer := timing.GetOptimizer()
		if err := timingOptimizer.Enable(); err != nil {
			logger.Error("Error habilitando timing optimizer: %v", err)
		} else {
			// Asegurar que se desactive al finalizar
			defer timingOptimizer.AutoDisable()
		}
	}

	// Mostrar información de la aplicación
	logger.Info("%s v%s - Versión Normal", constants.AppName, constants.AppVersion)
	logger.Info("%s", constants.AppDesc)
	logger.Info("===================================")

	// Buscar proceso de RO
	logger.Info("Buscando proceso de Ragnarok Online...")
	roProcess, err := process.FindROProcess()
	if err != nil {
		logger.Error("Error: %v", err)
		logger.Info("Asegúrate de que Ragnarok Online esté ejecutándose.")
		logger.Info("Presiona Enter para salir...")
		fmt.Scanln()
		return
	}

	logger.Info("Proceso encontrado: %s (PID: %d)", roProcess.Name, roProcess.PID)

	// Obtener handle de la ventana
	logger.Info("Obteniendo handle de la ventana...")
	roHandle, err := getWindowHandle(roProcess.PID)
	if err != nil {
		logger.Error("Error obteniendo ventana: %v", err)
		logger.Info("Presiona Enter para salir...")
		fmt.Scanln()
		return
	}
	logger.Info("Handle de ventana obtenido: 0x%x", roHandle)

	// Crear cliente de memoria
	logger.Info("Conectando a la memoria del proceso...")
	memClient, err := memory.NewClient(roProcess.PID)
	if err != nil {
		logger.Error("Error conectando a memoria: %v", err)
		logger.Info("Presiona Enter para salir...")
		fmt.Scanln()
		return
	}
	defer memClient.Close()

	logger.Info("Conectado a la memoria del proceso")

	// Inicializar ServiceManager con mejores prácticas
	serviceManager := services.NewServiceManager(cfg, memClient, roHandle)

	// Inicializar monitoreo del sistema
	systemMonitor := monitoring.NewSystemMonitor()
	consoleDisplay := monitoring.NewConsoleDisplay(500 * time.Millisecond) // Actualizar cada 500ms

	// Mostrar configuración
	logger.Info("Configuración:")
	logger.Info("- HP Pot: %s (cuando HP < %d%%)", cfg.AutoPot.HPPotKey, cfg.AutoPot.HPPercent)
	logger.Info("- SP Pot: %s (cuando SP < %d%%)", cfg.AutoPot.SPPotKey, cfg.AutoPot.SPPercent)
	logger.Info("- Delay: %dms", cfg.AutoPot.CheckDelay)

	if cfg.Spammer.Enabled {
		logger.Info("- Spammer: HABILITADO (%d teclas + click izquierdo, %.1fms delay)", len(cfg.Spammer.Keys), cfg.Spammer.Delay)
		logger.Info("- Teclas spammer: %v", cfg.Spammer.Keys)
		if cfg.Spammer.MouseClick {
			logger.Info("- Mouse Click: HABILITADO (click automático con teclas)")
		}
	} else {
		logger.Info("- Spammer: DESHABILITADO")
	}

	logger.Info("Presiona Ctrl+C para salir...")
	logger.Info("Comandos disponibles:")
	logger.Info("- 'aspd' - Buscar ASPD en memoria")
	logger.Info("- 'monitor' - Monitorear direcciones ASPD")
	logger.Info("- 'debug' - Debug detallado de ASPD")
	logger.Info("- 'stats' - Mostrar estadísticas")
	logger.Info("- 'metrics' - Habilitar/deshabilitar métricas")
	logger.Info("- 'timing' - Habilitar/deshabilitar optimización de timing")
	logger.Info("- 'benchmark' - Benchmark de precisión de timing")
	logger.Info("- 'predictor' - Mostrar estadísticas del predictor")
	logger.Info("- 'quit' - Salir")

	// Canal para manejar interrupciones
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	// Canal para comandos de usuario
	commandChan := make(chan string, 1)
	go handleUserCommands(commandChan)

	// El ServiceManager maneja su propio contexto

	// Iniciar ServiceManager
	if err := serviceManager.Start(); err != nil {
		logger.Error("Error iniciando ServiceManager: %v", err)
		return
	}
	defer serviceManager.Stop()

	// Loop principal - solo para mostrar estado
	ticker := time.NewTicker(1 * time.Second) // Mostrar estado cada segundo
	defer ticker.Stop()

	for {
		select {
		case <-interrupt:
			logger.Info("Saliendo...")
			// El ServiceManager se detiene automáticamente con defer
			return
		case command := <-commandChan:
			handleCommand(command, memClient, serviceManager, systemMonitor)
		case <-ticker.C:
			// Leer HP/SP
			hp, maxHP, err := memClient.ReadHP()
			if err != nil {
				continue
			}

			sp, maxSP, err := memClient.ReadSP()
			if err != nil {
				continue
			}

			// Calcular porcentajes
			hpPercent := int((float64(hp) / float64(maxHP)) * 100)
			spPercent := int((float64(sp) / float64(maxSP)) * 100)

			// Obtener estadísticas del ServiceManager y sistema
			serviceStats := serviceManager.GetStats()
			systemStats := systemMonitor.GetSystemStats()

			// Mostrar estadísticas optimizadas
			consoleDisplay.DisplayStats(
				hp, maxHP, hpPercent, sp, maxSP, spPercent,
				serviceStats, systemStats,
			)
		}
	}
}

// getWindowHandle obtiene el handle de la ventana del proceso
func getWindowHandle(pid uint32) (windows.Handle, error) {
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

// handleUserCommands maneja los comandos de entrada del usuario
func handleUserCommands(commandChan chan<- string) {
	for {
		var input string
		fmt.Scanln(&input)
		commandChan <- input
	}
}

// handleCommand procesa los comandos del usuario
func handleCommand(command string, memClient *memory.Client, serviceManager *services.ServiceManager, systemMonitor *monitoring.SystemMonitor) {
	switch command {
	case "aspd":
		searchASPDInMemory(memClient)
	case "monitor":
		monitorASPDAddresses(memClient)
	case "debug":
		debugASPDInfo(memClient)
	case "stats":
		showDetailedStats(memClient, serviceManager, systemMonitor)
	case "metrics":
		toggleMetrics()
	case "timing":
		toggleTiming()
	case "benchmark":
		runTimingBenchmark()
	case "predictor":
		showPredictorStats(serviceManager)
	case "quit":
		logger.Info("Saliendo...")
		os.Exit(0)
	default:
		logger.Info("Comando no reconocido: %s", command)
		logger.Info("Comandos disponibles: aspd, monitor, debug, stats, metrics, timing, benchmark, predictor, quit")
	}
}

// searchASPDInMemory busca el ASPD en toda la memoria del proceso
func searchASPDInMemory(memClient *memory.Client) {
	logger.Info("🔍 Buscando ASPD (158) en memoria...")

	targetASPD := uint32(158)
	searchResults := make(map[string]interface{})

	// Rango 1: Alrededor de la dirección del personaje
	baseAddr := uintptr(0x10df5d8)
	logger.Info("Buscando cerca del personaje (0x%x)...", baseAddr)

	for i := -100; i <= 300; i++ {
		addr := baseAddr + uintptr(i*4)
		value, err := memClient.ReadUint32(addr)
		if err == nil && value == targetASPD {
			searchResults[fmt.Sprintf("near_character_0x%x", addr)] = map[string]interface{}{
				"address":          fmt.Sprintf("0x%x", addr),
				"value":            value,
				"aspd":             value,
				"offset_from_base": i * 4,
			}
			logger.Info("✅ Encontrado ASPD en: 0x%x (offset: %d)", addr, i*4)
		}
	}

	// Rango 2: Alrededor de las direcciones de HP/SP
	hpAddr := uintptr(0x10DCE10)
	logger.Info("Buscando cerca de HP/SP (0x%x)...", hpAddr)

	for i := -200; i <= 200; i++ {
		addr := hpAddr + uintptr(i*4)
		value, err := memClient.ReadUint32(addr)
		if err == nil && value == targetASPD {
			searchResults[fmt.Sprintf("near_hp_0x%x", addr)] = map[string]interface{}{
				"address":        fmt.Sprintf("0x%x", addr),
				"value":          value,
				"aspd":           value,
				"offset_from_hp": i * 4,
			}
			logger.Info("✅ Encontrado ASPD en: 0x%x (offset desde HP: %d)", addr, i*4)
		}
	}

	// Rango 3: Búsqueda en heap
	logger.Info("Buscando en heap...")
	heapFound := 0
	for i := 0; i < 2000; i++ {
		addr := uintptr(0x10000000 + uintptr(i*4))
		value, err := memClient.ReadUint32(addr)
		if err == nil && value == targetASPD {
			searchResults[fmt.Sprintf("heap_0x%x", addr)] = map[string]interface{}{
				"address":     fmt.Sprintf("0x%x", addr),
				"value":       value,
				"aspd":        value,
				"heap_offset": i * 4,
			}
			heapFound++
			if heapFound <= 5 { // Mostrar solo los primeros 5
				logger.Info("✅ Encontrado ASPD en heap: 0x%x", addr)
			}
		}
	}
	if heapFound > 5 {
		logger.Info("... y %d más en heap", heapFound-5)
	}

	logger.Info("🎯 Búsqueda completada. Total encontrados: %d", len(searchResults))
}

// monitorASPDAddresses monitorea direcciones específicas de ASPD
func monitorASPDAddresses(memClient *memory.Client) {
	logger.Info("📊 Monitoreando direcciones ASPD...")

	addresses := []uintptr{
		0x10df5d8 + 0x10, // Cerca del personaje + offset
		0x10df5d8 + 0x20, // Cerca del personaje + offset
		0x10df5d8 + 0x30, // Cerca del personaje + offset
		0x10DCE10 + 0x50, // Cerca de HP + offset
		0x10DCE10 + 0x60, // Cerca de HP + offset
		0x10DCE10 + 0x70, // Cerca de HP + offset
	}

	for i, addr := range addresses {
		value, err := memClient.ReadUint32(addr)
		if err == nil {
			isASPD := value == 158
			status := "❌"
			if isASPD {
				status = "✅"
			}
			logger.Info("%s Addr %d: 0x%x = %d (0x%X) %s",
				status, i+1, addr, value, value,
				map[bool]string{true: "[ES ASPD!]", false: ""}[isASPD])
		} else {
			logger.Info("❌ Addr %d: 0x%x = ERROR", i+1, addr)
		}
	}
}

// debugASPDInfo muestra información detallada de debug
func debugASPDInfo(memClient *memory.Client) {
	logger.Info("🔧 Debug detallado de ASPD...")

	// Leer alrededor del personaje
	baseAddr := uintptr(0x10df5d8)
	logger.Info("=== OFFSETS DEL PERSONAJE ===")

	for i := 0; i < 20; i++ {
		offset := i * 4
		addr := baseAddr + uintptr(offset)
		value, err := memClient.ReadUint32(addr)
		if err == nil {
			isASPD := value == 158
			status := ""
			if isASPD {
				status = " ← ¡ASPD!"
			}
			logger.Info("Offset %2d: 0x%x = %4d (0x%X)%s",
				offset, addr, value, value, status)
		}
	}

	// Leer alrededor de HP/SP
	hpAddr := uintptr(0x10DCE10)
	logger.Info("=== OFFSETS DE HP/SP ===")

	for i := 0; i < 15; i++ {
		offset := i * 4
		addr := hpAddr + uintptr(offset)
		value, err := memClient.ReadUint32(addr)
		if err == nil {
			isASPD := value == 158
			status := ""
			if isASPD {
				status = " ← ¡ASPD!"
			}
			logger.Info("HP+%2d: 0x%x = %4d (0x%X)%s",
				offset, addr, value, value, status)
		}
	}
}

// showDetailedStats muestra estadísticas detalladas
func showDetailedStats(memClient *memory.Client, serviceManager *services.ServiceManager, systemMonitor *monitoring.SystemMonitor) {
	logger.Info("📈 Estadísticas detalladas:")

	// HP/SP
	hp, maxHP, err := memClient.ReadHP()
	if err == nil {
		hpPercent := int((float64(hp) / float64(maxHP)) * 100)
		logger.Info("HP: %d/%d (%d%%)", hp, maxHP, hpPercent)
	}

	sp, maxSP, err := memClient.ReadSP()
	if err == nil {
		spPercent := int((float64(sp) / float64(maxSP)) * 100)
		logger.Info("SP: %d/%d (%d%%)", sp, maxSP, spPercent)
	}

	// Job
	jobID, err := memClient.ReadJobID()
	if err == nil {
		logger.Info("Job ID: %d", jobID)
	}

	// ASPD
	aspd, err := memClient.ReadASPD()
	if err == nil {
		logger.Info("ASPD: %d", aspd)
	} else {
		logger.Info("ASPD: Error leyendo")
	}

	// Stats del servicio
	serviceStats := serviceManager.GetStats()
	logger.Info("Servicios: %v", serviceStats)
}

// toggleMetrics habilita/deshabilita las métricas
func toggleMetrics() {
	metricsInstance := metrics.GetMetrics()
	if metricsInstance.Enabled {
		metricsInstance.Disable()
		logger.Info("📊 Métricas deshabilitadas")
	} else {
		if err := metricsInstance.Enable(6060); err != nil {
			logger.Error("Error habilitando métricas: %v", err)
		} else {
			logger.Info("📊 Métricas habilitadas en puerto 6060")
			logger.Info("🔍 pprof disponible en: http://localhost:6060/debug/pprof/")
			logger.Info("📈 Métricas disponibles en: http://localhost:6060/metrics")
		}
	}
}

// toggleTiming habilita/deshabilita la optimización de timing
func toggleTiming() {
	timingOptimizer := timing.GetOptimizer()
	if timingOptimizer.IsEnabled() {
		if err := timingOptimizer.Disable(); err != nil {
			logger.Error("Error deshabilitando timing optimizer: %v", err)
		} else {
			logger.Info("⚡ Optimización de timing deshabilitada")
		}
	} else {
		if err := timingOptimizer.Enable(); err != nil {
			logger.Error("Error habilitando timing optimizer: %v", err)
		} else {
			logger.Info("⚡ Optimización de timing habilitada (timeBeginPeriod=1ms)")
			logger.Info("💡 Esto puede mejorar la precisión de los timers en Windows")
		}
	}
}

// runTimingBenchmark ejecuta un benchmark de precisión de timing
func runTimingBenchmark() {
	logger.Info("🏃 Ejecutando benchmark de timing...")

	timingOptimizer := timing.GetOptimizer()
	stats := timingOptimizer.BenchmarkTiming()

	logger.Info("📊 Resultados del benchmark:")
	logger.Info("- Iteraciones: %v", stats["iterations"])
	logger.Info("- Delay objetivo: %v", stats["target_delay"])
	logger.Info("- Error promedio: %v (%v μs)", stats["avg_error"], stats["avg_error_us"])
	logger.Info("- Error máximo: %v (%v μs)", stats["max_error"], stats["max_error_us"])
	logger.Info("- Error mínimo: %v (%v μs)", stats["min_error"], stats["min_error_us"])

	if timingOptimizer.IsEnabled() {
		logger.Info("✅ Timing optimizer está HABILITADO")
	} else {
		logger.Info("❌ Timing optimizer está DESHABILITADO")
	}
}

// showPredictorStats muestra estadísticas del predictor
func showPredictorStats(serviceManager *services.ServiceManager) {
	logger.Info("🧠 Estadísticas del Predictor:")

	// Obtener estadísticas del service manager
	stats := serviceManager.GetStats()
	logger.Info("📊 Servicios activos: %v", stats)

	// Mostrar información sobre el predictor
	logger.Info("🔍 El predictor analiza los últimos 20 puntos de datos")
	logger.Info("⏱️ Predice HP/SP con 200ms de anticipación")
	logger.Info("📈 Usa regresión lineal para calcular tendencias")
	logger.Info("💡 Busca en los logs los mensajes:")
	logger.Info("   - 🧠 HP/SP Pot PREDICCIÓN: Cuando usa predicción")
	logger.Info("   - 💊 HP/SP Pot NORMAL: Cuando usa HP/SP actual")
	logger.Info("   - 📊 HP/SP Status: Información periódica del predictor")
}
