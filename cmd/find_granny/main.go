package main

import (
	"fmt"
	"strings"

	"dtap-go/internal/memory"
	"dtap-go/internal/process"
	"dtap-go/pkg/logger"
)

func main() {
	// Configurar logger
	logger.Global.SetLevel(logger.INFO)

	// Mostrar información de la aplicación
	logger.Info("=== FIND GRANNY - Búsqueda de Funciones de Granny ===")
	logger.Info("Buscando funciones de Granny relacionadas con velocidad de ataque...")
	logger.Info("=====================================================")

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
	logger.Info("")

	// Búsqueda de funciones de Granny
	searchGrannyFunctions(memClient)

	logger.Info("")
	logger.Info("Presiona Enter para salir...")
	fmt.Scanln()
}

// searchGrannyFunctions busca funciones de Granny relacionadas con velocidad
func searchGrannyFunctions(memClient *memory.Client) {
	// Funciones de Granny relacionadas con velocidad y animación
	grannyFunctions := []string{
		// Funciones de control de velocidad
		"_GrannyGetControlSpeed@4",
		"_GrannySetControlSpeed@8",
		"_GrannyGetControlLocalDuration@4",
		"_GrannySetControlLocalDuration@8",

		// Funciones de loop y estado
		"_GrannyGetControlLoopState@12",
		"_GrannySetControlLoopState@12",
		"_GrannyGetControlLoopCount@4",
		"_GrannyGetControlLoopIndex@4",

		// Funciones de reloj y tiempo
		"_GrannyGetControlRawLocalClock@4",
		"_GrannySetControlRawLocalClock@8",
		"_GrannyGetControlWeight@4",
		"_GrannySetControlWeight@8",

		// Funciones de animación
		"_GrannyGetControlUserDataArray@4",
		"_GrannySetControlUserDataArray@8",
		"_GrannyGetLocalPoseBoneCount@4",
		"_GrannyGetLocalPoseBoneArray@4",

		// Funciones de archivo y datos
		"_GrannyGetDataTreeFromFile@4",
		"_GrannyGetFileInfo@4",
		"_GrannyGetFileTypeTag@4",
		"_GrannyGetIGCInfo@4",

		// Funciones de contador
		"_GrannyGetCounterCount@0",
		"_GrannyGetCounterResults@8",

		// Funciones de árbol de datos
		"_GrannyGetConvertedTreeSize@12",
		"_GrannyGetGRNSectionArray@4",

		// Funciones de cabeza de Granny
		"_GrannyGetGrannyHeadBezier@4",
		"_GrannyGetGrannyHeadBezierCount@0",
		"_GrannyGetGrannyHeadWidthOverHeight@0",

		// Funciones de callback
		"_GrannyGetDefaultFileReaderOpenCallback@0",

		// Strings más generales
		"Granny",
		"granny",
		"GRANNY",
		"ControlSpeed",
		"controlSpeed",
		"CONTROL_SPEED",
		"LocalDuration",
		"localDuration",
		"LOCAL_DURATION",
		"LoopState",
		"loopState",
		"LOOP_STATE",
		"RawLocalClock",
		"rawLocalClock",
		"RAW_LOCAL_CLOCK",
		"UserDataArray",
		"userDataArray",
		"USER_DATA_ARRAY",
		"PoseBone",
		"poseBone",
		"POSE_BONE",
		"DataTree",
		"dataTree",
		"DATA_TREE",
		"FileInfo",
		"fileInfo",
		"FILE_INFO",
		"Counter",
		"counter",
		"COUNTER",
		"Bezier",
		"bezier",
		"BEZIER",
	}

	totalFound := 0

	logger.Info("🔍 INICIANDO BÚSQUEDA DE FUNCIONES GRANNY")
	logger.Info("Buscando %d funciones y strings de Granny", len(grannyFunctions))
	logger.Info("")

	// === RANGO 1: ALREDEDOR DEL PERSONAJE ===
	logger.Info("📍 RANGO 1: Alrededor del personaje (0x10df5d8)")
	baseAddr := uintptr(0x10df5d8)
	found1 := searchGrannyInRange(memClient, baseAddr, -2000, 4000, grannyFunctions, "personaje")
	totalFound += found1

	// === RANGO 2: ALREDEDOR DE HP/SP ===
	logger.Info("")
	logger.Info("📍 RANGO 2: Alrededor de HP/SP (0x10DCE10)")
	hpAddr := uintptr(0x10DCE10)
	found2 := searchGrannyInRange(memClient, hpAddr, -4000, 8000, grannyFunctions, "HP/SP")
	totalFound += found2

	// === RANGO 3: RANGOS CONOCIDOS DE RO ===
	logger.Info("")
	logger.Info("📍 RANGO 3: Rangos conocidos de RO")
	found3 := searchGrannyInKnownRanges(memClient, grannyFunctions)
	totalFound += found3

	// === RANGO 4: BÚSQUEDA EN HEAP ===
	logger.Info("")
	logger.Info("📍 RANGO 4: Heap (0x10000000 - 0x20000000)")
	found4 := searchGrannyInHeap(memClient, grannyFunctions)
	totalFound += found4

	// === RANGO 5: BÚSQUEDA EN STACK ===
	logger.Info("")
	logger.Info("📍 RANGO 5: Stack (0x20000000 - 0x30000000)")
	found5 := searchGrannyInStack(memClient, grannyFunctions)
	totalFound += found5

	// === RANGO 6: BÚSQUEDA EN CÓDIGO ===
	logger.Info("")
	logger.Info("📍 RANGO 6: Código (0x400000 - 0x1000000)")
	found6 := searchGrannyInCode(memClient, grannyFunctions)
	totalFound += found6

	// === RESUMEN FINAL ===
	logger.Info("")
	logger.Info("🎯 RESUMEN DE BÚSQUEDA:")
	logger.Info("   Rango 1 (Personaje): %d encontrados", found1)
	logger.Info("   Rango 2 (HP/SP):     %d encontrados", found2)
	logger.Info("   Rango 3 (Conocidos): %d encontrados", found3)
	logger.Info("   Rango 4 (Heap):      %d encontrados", found4)
	logger.Info("   Rango 5 (Stack):     %d encontrados", found5)
	logger.Info("   Rango 6 (Código):    %d encontrados", found6)
	logger.Info("   TOTAL:               %d encontrados", totalFound)

	if totalFound == 0 {
		logger.Info("")
		logger.Info("❌ No se encontraron funciones de Granny en ningún rango.")
		logger.Info("   Posibles causas:")
		logger.Info("   - El personaje no está en el juego")
		logger.Info("   - Las funciones están en un formato diferente")
		logger.Info("   - La dirección está en un rango no explorado")
		logger.Info("   - Las funciones están encriptadas o comprimidas")
	} else {
		logger.Info("")
		logger.Info("✅ Se encontraron %d funciones/strings de Granny", totalFound)
		logger.Info("   Revisa las direcciones encontradas arriba.")
		logger.Info("")
		logger.Info("💡 ANÁLISIS:")
		logger.Info("   - Las funciones de Granny manejan animaciones y velocidad")
		logger.Info("   - '_GrannyGetControlSpeed@4' puede estar relacionada con ASPD")
		logger.Info("   - '_GrannyGetControlLocalDuration@4' puede controlar duración de ataques")
		logger.Info("   - '_GrannyGetControlLoopState@12' puede manejar estados de ataque")
	}
}

// searchGrannyInRange busca funciones de Granny en un rango específico
func searchGrannyInRange(memClient *memory.Client, baseAddr uintptr, startOffset, endOffset int, targetFunctions []string, rangeName string) int {
	found := 0
	logger.Info("   Buscando funciones Granny en %s (offset %d a %d)...", rangeName, startOffset, endOffset)

	for i := startOffset; i <= endOffset; i++ {
		addr := baseAddr + uintptr(i)
		// Leer 512 bytes para buscar funciones
		data, err := memClient.ReadBytes(addr, 512)
		if err == nil {
			dataStr := string(data)
			for _, targetFunc := range targetFunctions {
				if strings.Contains(dataStr, targetFunc) {
					found++
					logger.Info("   ✅ Función '%s' encontrada en: 0x%x (offset: %d)", targetFunc, addr, i)
					// Mostrar contexto de la función
					showGrannyContext(dataStr, targetFunc, addr)
					break // Solo mostrar una vez por dirección
				}
			}
		}
	}

	if found == 0 {
		logger.Info("   ❌ No encontrado en %s", rangeName)
	}
	return found
}

// searchGrannyInKnownRanges busca funciones de Granny en rangos conocidos de RO
func searchGrannyInKnownRanges(memClient *memory.Client, targetFunctions []string) int {
	found := 0
	logger.Info("   Buscando funciones Granny en rangos conocidos de RO...")

	// Rangos conocidos de RO
	knownRanges := []struct {
		name  string
		start uintptr
		end   uintptr
	}{
		{"Status Buffer", 0x10DD000, 0x10DE000},
		{"Character Data", 0x10DF000, 0x10E0000},
		{"Map Data", 0x10E0000, 0x10E1000},
		{"Inventory", 0x10E1000, 0x10E2000},
		{"Skills", 0x10E2000, 0x10E3000},
		{"Stats", 0x10E3000, 0x10E4000},
		{"Equipment", 0x10E4000, 0x10E5000},
		{"Party", 0x10E5000, 0x10E6000},
		{"Guild", 0x10E6000, 0x10E7000},
		{"Chat", 0x10E7000, 0x10E8000},
		{"UI Text", 0x10E8000, 0x10E9000},
		{"Menu Text", 0x10E9000, 0x10EA000},
		{"Error Messages", 0x10EA000, 0x10EB000},
		{"System Messages", 0x10EB000, 0x10EC000},
		{"Item Names", 0x10EC000, 0x10ED000},
		{"Animation Data", 0x10ED000, 0x10EE000},
		{"Model Data", 0x10EE000, 0x10EF000},
		{"Texture Data", 0x10EF000, 0x10F0000},
	}

	for _, r := range knownRanges {
		rangeFound := 0
		for addr := r.start; addr < r.end; addr += 4 {
			// Leer 256 bytes para buscar funciones
			data, err := memClient.ReadBytes(addr, 256)
			if err == nil {
				dataStr := string(data)
				for _, targetFunc := range targetFunctions {
					if strings.Contains(dataStr, targetFunc) {
						found++
						rangeFound++
						if rangeFound <= 2 {
							logger.Info("   ✅ Función '%s' en %s: 0x%x", targetFunc, r.name, addr)
							showGrannyContext(dataStr, targetFunc, addr)
						}
						break
					}
				}
			}
		}
		if rangeFound > 2 {
			logger.Info("   ... y %d más en %s", rangeFound-2, r.name)
		}
	}

	if found == 0 {
		logger.Info("   ❌ No encontrado en rangos conocidos")
	}
	return found
}

// searchGrannyInHeap busca funciones de Granny en el heap
func searchGrannyInHeap(memClient *memory.Client, targetFunctions []string) int {
	found := 0
	logger.Info("   Buscando funciones Granny en heap...")

	// Buscar en bloques de 1000 direcciones
	for block := 0; block < 200; block++ { // Aumentado para Granny
		blockFound := 0
		startAddr := uintptr(0x10000000 + uintptr(block*4000))

		for i := 0; i < 1000; i++ {
			addr := startAddr + uintptr(i)
			// Leer 128 bytes para buscar funciones
			data, err := memClient.ReadBytes(addr, 128)
			if err == nil {
				dataStr := string(data)
				for _, targetFunc := range targetFunctions {
					if strings.Contains(dataStr, targetFunc) {
						found++
						blockFound++
						if blockFound <= 2 {
							logger.Info("   ✅ Función '%s' en Heap: 0x%x", targetFunc, addr)
						}
						break
					}
				}
			}
		}

		if blockFound > 2 {
			logger.Info("   ... y %d más en bloque 0x%x", blockFound-2, startAddr)
		}
	}

	if found == 0 {
		logger.Info("   ❌ No encontrado en heap")
	}
	return found
}

// searchGrannyInStack busca funciones de Granny en el stack
func searchGrannyInStack(memClient *memory.Client, targetFunctions []string) int {
	found := 0
	logger.Info("   Buscando funciones Granny en stack...")

	// Buscar en bloques de 1000 direcciones
	for block := 0; block < 200; block++ { // Aumentado para Granny
		blockFound := 0
		startAddr := uintptr(0x20000000 + uintptr(block*4000))

		for i := 0; i < 1000; i++ {
			addr := startAddr + uintptr(i)
			// Leer 128 bytes para buscar funciones
			data, err := memClient.ReadBytes(addr, 128)
			if err == nil {
				dataStr := string(data)
				for _, targetFunc := range targetFunctions {
					if strings.Contains(dataStr, targetFunc) {
						found++
						blockFound++
						if blockFound <= 2 {
							logger.Info("   ✅ Función '%s' en Stack: 0x%x", targetFunc, addr)
						}
						break
					}
				}
			}
		}

		if blockFound > 2 {
			logger.Info("   ... y %d más en bloque 0x%x", blockFound-2, startAddr)
		}
	}

	if found == 0 {
		logger.Info("   ❌ No encontrado en stack")
	}
	return found
}

// searchGrannyInCode busca funciones de Granny en el código
func searchGrannyInCode(memClient *memory.Client, targetFunctions []string) int {
	found := 0
	logger.Info("   Buscando funciones Granny en código...")

	// Buscar en bloques de código
	for block := 0; block < 100; block++ {
		blockFound := 0
		startAddr := uintptr(0x400000 + uintptr(block*4000))

		for i := 0; i < 1000; i++ {
			addr := startAddr + uintptr(i)
			// Leer 64 bytes para buscar funciones
			data, err := memClient.ReadBytes(addr, 64)
			if err == nil {
				dataStr := string(data)
				for _, targetFunc := range targetFunctions {
					if strings.Contains(dataStr, targetFunc) {
						found++
						blockFound++
						if blockFound <= 1 {
							logger.Info("   ✅ Función '%s' en Código: 0x%x", targetFunc, addr)
						}
						break
					}
				}
			}
		}

		if blockFound > 1 {
			logger.Info("   ... y %d más en bloque 0x%x", blockFound-1, startAddr)
		}
	}

	if found == 0 {
		logger.Info("   ❌ No encontrado en código")
	}
	return found
}

// showGrannyContext muestra el contexto alrededor de la función de Granny encontrada
func showGrannyContext(dataStr, targetFunc string, addr uintptr) {
	// Encontrar la posición de la función
	pos := strings.Index(dataStr, targetFunc)
	if pos == -1 {
		return
	}

	// Mostrar contexto (64 caracteres antes y después)
	start := pos - 64
	if start < 0 {
		start = 0
	}
	end := pos + len(targetFunc) + 64
	if end > len(dataStr) {
		end = len(dataStr)
	}

	context := dataStr[start:end]
	// Limpiar caracteres no imprimibles
	cleanContext := ""
	for _, c := range context {
		if c >= 32 && c <= 126 {
			cleanContext += string(c)
		} else {
			cleanContext += "."
		}
	}

	logger.Info("      Contexto: '%s'", cleanContext)
}
