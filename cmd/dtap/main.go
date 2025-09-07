package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
	"unsafe"

	"dtap-go/internal/memory"
	"dtap-go/internal/monitoring"
	"dtap-go/internal/process"
	"dtap-go/internal/services"
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
		logger.Info("- Spammer: HABILITADO (%d teclas + click izquierdo, %dms delay)", len(cfg.Spammer.Keys), cfg.Spammer.Delay)
		logger.Info("- Teclas spammer: %v", cfg.Spammer.Keys)
		if cfg.Spammer.SkillCancel {
			logger.Info("- Skill Cancel: HABILITADO (cancela habilidades)")
		}
		if cfg.Spammer.DoubleClick {
			logger.Info("- Double Click: HABILITADO (doble click para cancelar)")
		}
	} else {
		logger.Info("- Spammer: DESHABILITADO")
	}

	logger.Info("Presiona Ctrl+C para salir...")

	// Canal para manejar interrupciones
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

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
