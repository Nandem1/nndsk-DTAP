package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
	"unsafe"

	"dtap-go/internal/autopot"
	"dtap-go/internal/memory"
	"dtap-go/internal/monitoring"
	"dtap-go/internal/process"
	"dtap-go/internal/spammer"
	"dtap-go/pkg/config"
	"dtap-go/pkg/constants"
	"dtap-go/pkg/logger"

	"golang.org/x/sys/windows"
)

func main() {
	// Cargar configuración WoE optimizada
	woeConfig, err := config.LoadWOEConfig(config.GetWOEConfigPath())
	if err != nil {
		logger.Error("Error cargando configuración WoE: %v", err)
		os.Exit(1)
	}

	// Configurar logger
	logger.Global.SetLevel(logger.ParseLogLevel(woeConfig.UI.LogLevel))

	// Mostrar información de la aplicación
	logger.Info("%s v%s - WoE Optimized Edition", constants.AppName, constants.AppVersion)
	logger.Info("%s", constants.AppDesc)
	logger.Info("===================================")
	logger.Info("🚀 MODO COMPETITIVO WoE ACTIVADO 🚀")
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

	// Crear cliente de memoria optimizado
	logger.Info("Conectando a la memoria del proceso (modo optimizado)...")
	memClient, err := memory.NewOptimizedClient(roProcess.PID)
	if err != nil {
		logger.Error("Error conectando a memoria: %v", err)
		logger.Info("Presiona Enter para salir...")
		fmt.Scanln()
		return
	}
	defer memClient.Close()

	logger.Info("Conectado a la memoria del proceso (modo optimizado)")

	// Inicializar servicios WoE optimizados
	autopotService := autopot.NewWOEOptimizedService(woeConfig.ConvertToStandardConfig(), memClient, roHandle)
	spammerService := spammer.NewWOEOptimizedSpammer(woeConfig.ConvertToStandardConfig(), roHandle)
	
	// Inicializar monitor de performance
	performanceMonitor := monitoring.NewWOEPerformanceMonitor(memClient, autopotService, spammerService, roHandle)

	// Mostrar configuración WoE
	logger.Info("Configuración WoE:")
	logger.Info("- HP Pot: %s (cuando HP < %d%%, emergencia < %d%%)", 
		woeConfig.AutoPot.HPPotKey, woeConfig.AutoPot.HPPercent, woeConfig.AutoPot.EmergencyHP)
	logger.Info("- SP Pot: %s (cuando SP < %d%%, emergencia < %d%%)", 
		woeConfig.AutoPot.SPPotKey, woeConfig.AutoPot.SPPercent, woeConfig.AutoPot.EmergencySP)
	logger.Info("- Check Delay: %dms (emergencia: %dms)", 
		woeConfig.AutoPot.CheckDelay, woeConfig.AutoPot.EmergencyDelay)
	logger.Info("- Smart Mode: %v, Predictive Mode: %v", 
		woeConfig.AutoPot.SmartMode, woeConfig.AutoPot.PredictiveMode)

	if woeConfig.Spammer.Enabled {
		logger.Info("- Spammer: HABILITADO (%d teclas, %dms delay, burst: %dms)", 
			len(woeConfig.Spammer.Keys), woeConfig.Spammer.Delay, woeConfig.Spammer.BurstDelay)
		logger.Info("- Teclas spammer: %v", woeConfig.Spammer.Keys)
		logger.Info("- Skill Cancel: %v, Double Click: %v", 
			woeConfig.Spammer.SkillCancel, woeConfig.Spammer.DoubleClick)
		logger.Info("- Burst Mode: %v, Ultra Fast: %v", 
			woeConfig.Spammer.BurstMode, woeConfig.Spammer.UltraFastMode)
	} else {
		logger.Info("- Spammer: DESHABILITADO")
	}

	logger.Info("- Performance: Target FPS: %d, Adaptive: %v", 
		woeConfig.Performance.TargetFPS, woeConfig.Performance.AdaptiveMode)
	logger.Info("- Monitoring: %v, Report Interval: %ds", 
		woeConfig.Monitoring.Enabled, woeConfig.Monitoring.ReportInterval)

	logger.Info("Presiona Ctrl+C para salir...")

	// Canal para manejar interrupciones
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	// Contexto para cancelar goroutines
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Iniciar servicios WoE optimizados
	autopotService.Start(ctx)
	spammerService.Start(ctx)
	
	if woeConfig.Monitoring.Enabled {
		performanceMonitor.Start(ctx)
	}

	// Loop principal optimizado
	ticker := time.NewTicker(500 * time.Millisecond) // Mostrar estado cada 500ms
	defer ticker.Stop()

	// Performance report ticker
	var reportTicker *time.Ticker
	if woeConfig.Monitoring.Enabled && woeConfig.Monitoring.ReportInterval > 0 {
		reportTicker = time.NewTicker(time.Duration(woeConfig.Monitoring.ReportInterval) * time.Second)
		defer reportTicker.Stop()
	}

	for {
		select {
		case <-interrupt:
			logger.Info("Saliendo...")
			cancel() // Cancelar todas las goroutines
			
			// Mostrar estadísticas finales
			if woeConfig.Monitoring.Enabled {
				performanceMonitor.Stop()
			}
			
			return
			
		case <-ticker.C:
			// Mostrar estado optimizado
			if woeConfig.UI.ShowPerformance {
				displayOptimizedStatus(memClient, autopotService, spammerService)
			}
			
		case <-reportTicker.C:
			// Mostrar reporte de performance
			if woeConfig.Monitoring.Enabled {
				performanceMonitor.PrintPerformanceReport()
			}
		}
	}
}

// displayOptimizedStatus muestra el estado optimizado del sistema
func displayOptimizedStatus(memClient *memory.OptimizedClient, autopotService *autopot.WOEOptimizedService, spammerService *spammer.WOEOptimizedSpammer) {
	// Obtener valores de memoria (ultra-rápido con cache)
	hp, maxHP, hpValid := memClient.GetHP()
	sp, maxSP, spValid := memClient.GetSP()
	
	if !hpValid || !spValid {
		fmt.Printf("\r[WoE] Memoria no disponible...")
		return
	}

	// Calcular porcentajes
	hpPercent := int((float64(hp) / float64(maxHP)) * 100)
	spPercent := int((float64(sp) / float64(maxSP)) * 100)

	// Obtener estadísticas de servicios
	potCount, emergencyMode, _ := autopotService.GetStats()
	spamCount, burstMode, activeKeys := spammerService.GetStats()
	
	// Obtener estadísticas de cache
	cacheAge, cacheValid := memClient.GetCacheStats()
	
	// Mostrar estado compacto
	status := fmt.Sprintf("\r[WoE] HP:%d/%d(%d%%) SP:%d/%d(%d%%) | Pots:%d | Spams:%d | Keys:%d | Cache:%v(%.1fms)",
		hp, maxHP, hpPercent, sp, maxSP, spPercent, potCount, spamCount, activeKeys, cacheValid, float64(cacheAge.Nanoseconds())/1000000)
	
	if emergencyMode {
		status += " [EMERGENCY]"
	}
	if burstMode {
		status += " [BURST]"
	}
	
	fmt.Print(status)
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