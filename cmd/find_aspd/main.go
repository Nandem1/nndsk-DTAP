package main

import (
	"fmt"

	"dtap-go/internal/memory"
	"dtap-go/internal/process"
	"dtap-go/pkg/logger"
)

func main() {
	// Configurar logger
	logger.Global.SetLevel(logger.INFO)

	// Mostrar información de la aplicación
	logger.Info("=== FIND AES - Búsqueda de Attack Speed ===")
	logger.Info("Buscando AES (Attack Speed) en memoria de Ragnarok Online...")
	logger.Info("==================================================")

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

	// Búsqueda ampliada de AES
	searchAESComprehensive(memClient)

	logger.Info("")
	logger.Info("Presiona Enter para salir...")
	fmt.Scanln()
}

// searchAESComprehensive realiza una búsqueda exhaustiva de AES
func searchAESComprehensive(memClient *memory.Client) {
	// Buscar múltiples valores posibles de AES
	aesValues := []uint32{166, 158, 150, 140, 130, 120, 110, 100, 90, 80, 70, 60, 50, 40, 30, 20, 10}
	totalFound := 0

	logger.Info("🔍 INICIANDO BÚSQUEDA EXHAUSTIVA DE AES (Attack Speed)")
	logger.Info("Buscando valores comunes de AES: %v", aesValues)
	logger.Info("")

	// === RANGO 1: ALREDEDOR DEL PERSONAJE ===
	logger.Info("📍 RANGO 1: Alrededor del personaje (0x10df5d8)")
	baseAddr := uintptr(0x10df5d8)
	found1 := searchMultipleValuesInRange(memClient, baseAddr, -500, 1000, aesValues, "personaje")
	totalFound += found1

	// === RANGO 2: ALREDEDOR DE HP/SP ===
	logger.Info("")
	logger.Info("📍 RANGO 2: Alrededor de HP/SP (0x10DCE10)")
	hpAddr := uintptr(0x10DCE10)
	found2 := searchMultipleValuesInRange(memClient, hpAddr, -1000, 2000, aesValues, "HP/SP")
	totalFound += found2

	// === RANGO 3: HEAP COMPLETO ===
	logger.Info("")
	logger.Info("📍 RANGO 3: Heap completo (0x10000000 - 0x20000000)")
	found3 := searchMultipleValuesInHeap(memClient, aesValues)
	totalFound += found3

	// === RANGO 4: STACK ===
	logger.Info("")
	logger.Info("📍 RANGO 4: Stack (0x20000000 - 0x30000000)")
	found4 := searchMultipleValuesInStack(memClient, aesValues)
	totalFound += found4

	// === RANGO 5: BÚSQUEDA EN RANGOS CONOCIDOS ===
	logger.Info("")
	logger.Info("📍 RANGO 5: Rangos conocidos de RO")
	found5 := searchMultipleValuesInKnownRanges(memClient, aesValues)
	totalFound += found5

	// === RANGO 6: BÚSQUEDA ALEATORIA EN MEMORIA ===
	logger.Info("")
	logger.Info("📍 RANGO 6: Búsqueda aleatoria en memoria")
	found6 := searchMultipleValuesInRandomMemory(memClient, aesValues)
	totalFound += found6

	// === RESUMEN FINAL ===
	logger.Info("")
	logger.Info("🎯 RESUMEN DE BÚSQUEDA:")
	logger.Info("   Rango 1 (Personaje): %d encontrados", found1)
	logger.Info("   Rango 2 (HP/SP):     %d encontrados", found2)
	logger.Info("   Rango 3 (Heap):      %d encontrados", found3)
	logger.Info("   Rango 4 (Stack):     %d encontrados", found4)
	logger.Info("   Rango 5 (Conocidos): %d encontrados", found5)
	logger.Info("   Rango 6 (Aleatorio): %d encontrados", found6)
	logger.Info("   TOTAL:               %d encontrados", totalFound)

	if totalFound == 0 {
		logger.Info("")
		logger.Info("❌ No se encontró AES en ningún rango.")
		logger.Info("   Posibles causas:")
		logger.Info("   - El personaje no está en el juego")
		logger.Info("   - El AES tiene un valor no común")
		logger.Info("   - La dirección está en un rango no explorado")
		logger.Info("   - El valor se almacena en un formato diferente")
	} else {
		logger.Info("")
		logger.Info("✅ Se encontraron %d direcciones con valores de AES", totalFound)
		logger.Info("   Revisa las direcciones encontradas arriba.")
	}
}

// searchMultipleValuesInRange busca múltiples valores en un rango específico
func searchMultipleValuesInRange(memClient *memory.Client, baseAddr uintptr, startOffset, endOffset int, targetValues []uint32, rangeName string) int {
	found := 0
	logger.Info("   Buscando valores AES en %s (offset %d a %d)...", rangeName, startOffset, endOffset)

	for i := startOffset; i <= endOffset; i++ {
		addr := baseAddr + uintptr(i*4)
		value, err := memClient.ReadUint32(addr)
		if err == nil {
			for _, targetValue := range targetValues {
				if value == targetValue {
					found++
					logger.Info("   ✅ AES %d encontrado en: 0x%x (offset: %d)", targetValue, addr, i*4)
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

// searchMultipleValuesInHeap busca múltiples valores en el heap de memoria
func searchMultipleValuesInHeap(memClient *memory.Client, targetValues []uint32) int {
	found := 0
	logger.Info("   Buscando valores AES en heap...")

	// Buscar en bloques de 1000 direcciones para no saturar
	for block := 0; block < 1000; block++ {
		blockFound := 0
		startAddr := uintptr(0x10000000 + uintptr(block*4000)) // 1000 * 4 bytes

		for i := 0; i < 1000; i++ {
			addr := startAddr + uintptr(i*4)
			value, err := memClient.ReadUint32(addr)
			if err == nil {
				for _, targetValue := range targetValues {
					if value == targetValue {
						found++
						blockFound++
						if blockFound <= 3 { // Mostrar solo los primeros 3 por bloque
							logger.Info("   ✅ AES %d en Heap: 0x%x", targetValue, addr)
						}
						break
					}
				}
			}
		}

		if blockFound > 3 {
			logger.Info("   ... y %d más en bloque 0x%x", blockFound-3, startAddr)
		}
	}

	if found == 0 {
		logger.Info("   ❌ No encontrado en heap")
	}
	return found
}

// searchMultipleValuesInStack busca múltiples valores en el stack de memoria
func searchMultipleValuesInStack(memClient *memory.Client, targetValues []uint32) int {
	found := 0
	logger.Info("   Buscando valores AES en stack...")

	// Buscar en bloques de 1000 direcciones
	for block := 0; block < 1000; block++ {
		blockFound := 0
		startAddr := uintptr(0x20000000 + uintptr(block*4000))

		for i := 0; i < 1000; i++ {
			addr := startAddr + uintptr(i*4)
			value, err := memClient.ReadUint32(addr)
			if err == nil {
				for _, targetValue := range targetValues {
					if value == targetValue {
						found++
						blockFound++
						if blockFound <= 3 {
							logger.Info("   ✅ AES %d en Stack: 0x%x", targetValue, addr)
						}
						break
					}
				}
			}
		}

		if blockFound > 3 {
			logger.Info("   ... y %d más en bloque 0x%x", blockFound-3, startAddr)
		}
	}

	if found == 0 {
		logger.Info("   ❌ No encontrado en stack")
	}
	return found
}

// searchMultipleValuesInKnownRanges busca múltiples valores en rangos conocidos de RO
func searchMultipleValuesInKnownRanges(memClient *memory.Client, targetValues []uint32) int {
	found := 0
	logger.Info("   Buscando valores AES en rangos conocidos de RO...")

	// Rangos conocidos de RO (basados en experiencia previa)
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
	}

	for _, r := range knownRanges {
		rangeFound := 0
		for addr := r.start; addr < r.end; addr += 4 {
			value, err := memClient.ReadUint32(addr)
			if err == nil {
				for _, targetValue := range targetValues {
					if value == targetValue {
						found++
						rangeFound++
						if rangeFound <= 2 {
							logger.Info("   ✅ AES %d en %s: 0x%x", targetValue, r.name, addr)
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

// searchMultipleValuesInRandomMemory busca múltiples valores en memoria de forma aleatoria
func searchMultipleValuesInRandomMemory(memClient *memory.Client, targetValues []uint32) int {
	found := 0
	logger.Info("   Búsqueda aleatoria de valores AES en memoria...")

	// Buscar en múltiples rangos aleatorios
	randomRanges := []struct {
		start uintptr
		end   uintptr
	}{
		{0x400000, 0x500000},  // Código
		{0x500000, 0x600000},  // Datos
		{0x600000, 0x700000},  // Más datos
		{0x700000, 0x800000},  // Más datos
		{0x800000, 0x900000},  // Más datos
		{0x900000, 0xA00000},  // Más datos
		{0xA00000, 0xB00000},  // Más datos
		{0xB00000, 0xC00000},  // Más datos
		{0xC00000, 0xD00000},  // Más datos
		{0xD00000, 0xE00000},  // Más datos
		{0xE00000, 0xF00000},  // Más datos
		{0xF00000, 0x1000000}, // Más datos
	}

	for _, r := range randomRanges {
		rangeFound := 0
		for addr := r.start; addr < r.end; addr += 4 {
			value, err := memClient.ReadUint32(addr)
			if err == nil {
				for _, targetValue := range targetValues {
					if value == targetValue {
						found++
						rangeFound++
						if rangeFound <= 1 {
							logger.Info("   ✅ AES %d en Random: 0x%x", targetValue, addr)
						}
						break
					}
				}
			}
		}
		if rangeFound > 1 {
			logger.Info("   ... y %d más en rango 0x%x-0x%x", rangeFound-1, r.start, r.end)
		}
	}

	if found == 0 {
		logger.Info("   ❌ No encontrado en búsqueda aleatoria")
	}
	return found
}
