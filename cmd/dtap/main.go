package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
	"unsafe"

	"dtap-go/internal/memory"
	"dtap-go/internal/normal/autopot"
	"dtap-go/internal/normal/spammer"
	"dtap-go/internal/process"
	"dtap-go/pkg/config"
	"dtap-go/pkg/constants"
	"dtap-go/pkg/injection"
	"dtap-go/pkg/logger"
	"dtap-go/pkg/visual"

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

	// Inicializar servicios
	spammerService := spammer.NewSpammerService(cfg, roHandle)
	autopotService := autopot.NewAutoPotService(cfg, memClient, roHandle)

	// Inicializar skill markers y range display
	configAdapter := visual.NewConfigAdapter()
	skillMarkerConfig := configAdapter.ConvertFromNormalConfig(cfg)
	rangeDisplayConfig := configAdapter.ConvertRangeDisplayFromNormalConfig(cfg)
	overlayManager := visual.NewOverlayManager(skillMarkerConfig, rangeDisplayConfig, logger.Global)

	// Inicializar sistema de inyección
	injectionConfig := injection.DefaultInjectionConfig()
	injectionManager := injection.NewInjectionManager(roProcess.PID, injectionConfig.DLLPath, logger.Global)

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

	if cfg.SkillMarkers.Enabled {
		logger.Info("- Skill Markers: HABILITADO (marcado de habilidades en el piso)")
		logger.Info("- Opacidad: %.1f, Tamaño: %.1f", cfg.SkillMarkers.Opacity, cfg.SkillMarkers.Size)
	} else {
		logger.Info("- Skill Markers: DESHABILITADO")
	}

	if cfg.RangeDisplay.Enabled {
		logger.Info("- Range Display: HABILITADO (círculo de rango del personaje)")
		logger.Info("- Color: %s, Opacidad: %.1f, Grosor: %.1f", cfg.RangeDisplay.Color, cfg.RangeDisplay.Opacity, cfg.RangeDisplay.Thickness)
	} else {
		logger.Info("- Range Display: DESHABILITADO")
	}

	if injectionConfig.Enabled {
		logger.Info("- Inyección DLL: HABILITADO (funcionalidades avanzadas)")
		logger.Info("- DLL: %s, Memoria Compartida: %s", injectionConfig.DLLPath, injectionConfig.SharedMemoryName)

		if injectionConfig.EnableGepardBypass {
			logger.Info("- Gepard Bypass: HABILITADO (protección anti-inyección)")
			logger.Info("- Auto-detección: %v, Bypass CRC: %v, Restaurar funciones: %v",
				injectionConfig.AutoDetectGepard, injectionConfig.BypassCRC, injectionConfig.RestoreFunctions)
		} else {
			logger.Info("- Gepard Bypass: DESHABILITADO")
		}
	} else {
		logger.Info("- Inyección DLL: DESHABILITADO")
	}

	logger.Info("Presiona Ctrl+C para salir...")

	// Canal para manejar interrupciones
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	// Contexto para cancelar goroutines
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Iniciar servicios
	spammerService.Start(ctx)
	autopotService.Start(ctx)

	// Iniciar overlay manager
	if err := overlayManager.Start(); err != nil {
		logger.Error("Error iniciando overlay manager: %v", err)
	} else {
		logger.Info("Overlay manager iniciado")
	}
	defer overlayManager.Stop()

	// Inicializar sistema de inyección
	if injectionConfig.Enabled {
		if err := injectionManager.Initialize(); err != nil {
			logger.Error("Error inicializando sistema de inyección: %v", err)
		} else {
			logger.Info("Sistema de inyección inicializado")

			// Configurar ventana de RO
			if err := injectionManager.SetROWindowHandle(uintptr(roHandle)); err != nil {
				logger.Error("Error configurando ventana de RO: %v", err)
			}

			// Configurar M2E (skill markers)
			if err := injectionManager.EnableM2E(cfg.SkillMarkers.Enabled); err != nil {
				logger.Error("Error configurando M2E: %v", err)
			}

			// Intentar inyectar DLL
			if injectionConfig.AutoInject {
				if err := injectionManager.Inject(); err != nil {
					logger.Error("Error inyectando DLL: %v", err)
					logger.Info("Nota: La DLL 'dtap-hook.dll' debe estar en el directorio del ejecutable")
				} else {
					logger.Info("DLL inyectada exitosamente")
				}
			}
		}
		defer injectionManager.Cleanup()
	}

	// Loop principal - solo para mostrar estado
	ticker := time.NewTicker(1 * time.Second) // Mostrar estado cada segundo
	defer ticker.Stop()

	// Ticker para simular skill markers (solo para testing)
	skillTicker := time.NewTicker(5 * time.Second)
	defer skillTicker.Stop()

	// Simular skill markers una vez al inicio si están habilitados
	if cfg.SkillMarkers.Enabled {
		overlayManager.SimulateSkillDetection()
	}

	for {
		select {
		case <-interrupt:
			logger.Info("Saliendo...")
			cancel() // Cancelar todas las goroutines
			return
		case <-ticker.C:
			// Solo mostrar estado, el autopot está en su propio servicio
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

			// Mostrar estado
			markerCount := overlayManager.GetActiveMarkers()
			fmt.Printf("\rHP: %d/%d (%d%%) | SP: %d/%d (%d%%) | Markers: %d",
				hp, maxHP, hpPercent, sp, maxSP, spPercent, len(markerCount))
		case <-skillTicker.C:
			// Simular detección de skills cada 5 segundos (solo para testing)
			if cfg.SkillMarkers.Enabled {
				overlayManager.SimulateSkillDetection()
			}
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
