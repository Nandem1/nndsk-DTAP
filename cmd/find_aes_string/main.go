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
	logger.Info("=== FIND AES STRING - Búsqueda de Attack Speed como String ===")
	logger.Info("Buscando strings AES en memoria de Ragnarok Online...")
	logger.Info("=============================================================")

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

	// Búsqueda de strings AES
	searchAESStrings(memClient)

	logger.Info("")
	logger.Info("Presiona Enter para salir...")
	fmt.Scanln()
}

// searchAESStrings busca strings relacionados con AES en memoria
func searchAESStrings(memClient *memory.Client) {
	// Strings a buscar relacionados con Attack Speed
	aesStrings := []string{
		"AES",
		"aes",
		"Attack Speed",
		"attack speed",
		"ATTACK SPEED",
		"ASPD",
		"aspd",
		"ASPD:",
		"aspd:",
		"Attack Speed:",
		"attack speed:",
		"Velocidad de Ataque",
		"velocidad de ataque",
		"VELOCIDAD DE ATAQUE",
		"Velocidad de Ataque:",
		"velocidad de ataque:",
		"Attack",
		"attack",
		"Speed",
		"speed",
		"SPEED",
		"Velocidad",
		"velocidad",
		"VELOCIDAD",
		"Ataque",
		"ataque",
		"ATAQUE",
	}

	totalFound := 0

	logger.Info("🔍 INICIANDO BÚSQUEDA DE STRINGS AES")
	logger.Info("Buscando %d strings relacionados con Attack Speed", len(aesStrings))
	logger.Info("")

	// === RANGO 1: ALREDEDOR DEL PERSONAJE ===
	logger.Info("📍 RANGO 1: Alrededor del personaje (0x10df5d8)")
	baseAddr := uintptr(0x10df5d8)
	found1 := searchStringsInRange(memClient, baseAddr, -1000, 2000, aesStrings, "personaje")
	totalFound += found1

	// === RANGO 2: ALREDEDOR DE HP/SP ===
	logger.Info("")
	logger.Info("📍 RANGO 2: Alrededor de HP/SP (0x10DCE10)")
	hpAddr := uintptr(0x10DCE10)
	found2 := searchStringsInRange(memClient, hpAddr, -2000, 4000, aesStrings, "HP/SP")
	totalFound += found2

	// === RANGO 3: RANGOS CONOCIDOS DE RO ===
	logger.Info("")
	logger.Info("📍 RANGO 3: Rangos conocidos de RO")
	found3 := searchStringsInKnownRanges(memClient, aesStrings)
	totalFound += found3

	// === RANGO 4: BÚSQUEDA EN HEAP ===
	logger.Info("")
	logger.Info("📍 RANGO 4: Heap (0x10000000 - 0x20000000)")
	found4 := searchStringsInHeap(memClient, aesStrings)
	totalFound += found4

	// === RANGO 5: BÚSQUEDA EN STACK ===
	logger.Info("")
	logger.Info("📍 RANGO 5: Stack (0x20000000 - 0x30000000)")
	found5 := searchStringsInStack(memClient, aesStrings)
	totalFound += found5

	// === RESUMEN FINAL ===
	logger.Info("")
	logger.Info("🎯 RESUMEN DE BÚSQUEDA:")
	logger.Info("   Rango 1 (Personaje): %d encontrados", found1)
	logger.Info("   Rango 2 (HP/SP):     %d encontrados", found2)
	logger.Info("   Rango 3 (Conocidos): %d encontrados", found3)
	logger.Info("   Rango 4 (Heap):      %d encontrados", found4)
	logger.Info("   Rango 5 (Stack):     %d encontrados", found5)
	logger.Info("   TOTAL:               %d encontrados", totalFound)

	if totalFound == 0 {
		logger.Info("")
		logger.Info("❌ No se encontraron strings AES en ningún rango.")
		logger.Info("   Posibles causas:")
		logger.Info("   - El personaje no está en el juego")
		logger.Info("   - Los strings están en un formato diferente")
		logger.Info("   - La dirección está en un rango no explorado")
		logger.Info("   - Los strings están encriptados o comprimidos")
	} else {
		logger.Info("")
		logger.Info("✅ Se encontraron %d strings relacionados con AES", totalFound)
		logger.Info("   Revisa las direcciones encontradas arriba.")
	}
}

// searchStringsInRange busca strings en un rango específico
func searchStringsInRange(memClient *memory.Client, baseAddr uintptr, startOffset, endOffset int, targetStrings []string, rangeName string) int {
	found := 0
	logger.Info("   Buscando strings AES en %s (offset %d a %d)...", rangeName, startOffset, endOffset)

	for i := startOffset; i <= endOffset; i++ {
		addr := baseAddr + uintptr(i)
		// Leer 256 bytes para buscar strings
		data, err := memClient.ReadBytes(addr, 256)
		if err == nil {
			dataStr := string(data)
			for _, targetStr := range targetStrings {
				if strings.Contains(dataStr, targetStr) {
					found++
					logger.Info("   ✅ String '%s' encontrado en: 0x%x (offset: %d)", targetStr, addr, i)
					// Mostrar contexto del string
					showStringContext(dataStr, targetStr, addr)
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

// searchStringsInKnownRanges busca strings en rangos conocidos de RO
func searchStringsInKnownRanges(memClient *memory.Client, targetStrings []string) int {
	found := 0
	logger.Info("   Buscando strings AES en rangos conocidos de RO...")

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
	}

	for _, r := range knownRanges {
		rangeFound := 0
		for addr := r.start; addr < r.end; addr += 4 {
			// Leer 128 bytes para buscar strings
			data, err := memClient.ReadBytes(addr, 128)
			if err == nil {
				dataStr := string(data)
				for _, targetStr := range targetStrings {
					if strings.Contains(dataStr, targetStr) {
						found++
						rangeFound++
						if rangeFound <= 2 {
							logger.Info("   ✅ String '%s' en %s: 0x%x", targetStr, r.name, addr)
							showStringContext(dataStr, targetStr, addr)
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

// searchStringsInHeap busca strings en el heap
func searchStringsInHeap(memClient *memory.Client, targetStrings []string) int {
	found := 0
	logger.Info("   Buscando strings AES en heap...")

	// Buscar en bloques de 1000 direcciones
	for block := 0; block < 100; block++ { // Reducido para evitar detección
		blockFound := 0
		startAddr := uintptr(0x10000000 + uintptr(block*4000))

		for i := 0; i < 1000; i++ {
			addr := startAddr + uintptr(i)
			// Leer 64 bytes para buscar strings
			data, err := memClient.ReadBytes(addr, 64)
			if err == nil {
				dataStr := string(data)
				for _, targetStr := range targetStrings {
					if strings.Contains(dataStr, targetStr) {
						found++
						blockFound++
						if blockFound <= 2 {
							logger.Info("   ✅ String '%s' en Heap: 0x%x", targetStr, addr)
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

// searchStringsInStack busca strings en el stack
func searchStringsInStack(memClient *memory.Client, targetStrings []string) int {
	found := 0
	logger.Info("   Buscando strings AES en stack...")

	// Buscar en bloques de 1000 direcciones
	for block := 0; block < 100; block++ { // Reducido para evitar detección
		blockFound := 0
		startAddr := uintptr(0x20000000 + uintptr(block*4000))

		for i := 0; i < 1000; i++ {
			addr := startAddr + uintptr(i)
			// Leer 64 bytes para buscar strings
			data, err := memClient.ReadBytes(addr, 64)
			if err == nil {
				dataStr := string(data)
				for _, targetStr := range targetStrings {
					if strings.Contains(dataStr, targetStr) {
						found++
						blockFound++
						if blockFound <= 2 {
							logger.Info("   ✅ String '%s' en Stack: 0x%x", targetStr, addr)
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

// showStringContext muestra el contexto alrededor del string encontrado
func showStringContext(dataStr, targetStr string, addr uintptr) {
	// Encontrar la posición del string
	pos := strings.Index(dataStr, targetStr)
	if pos == -1 {
		return
	}

	// Mostrar contexto (32 caracteres antes y después)
	start := pos - 32
	if start < 0 {
		start = 0
	}
	end := pos + len(targetStr) + 32
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
