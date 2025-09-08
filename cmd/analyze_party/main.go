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
	logger.Info("=== ANALYZE PARTY - Análisis de Datos de Party ===")
	logger.Info("Analizando el área donde se encontró 'seamaeplayo'...")
	logger.Info("=================================================")

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

	// Análisis específico del área de party
	analyzePartyArea(memClient)

	logger.Info("")
	logger.Info("Presiona Enter para salir...")
	fmt.Scanln()
}

// analyzePartyArea analiza específicamente el área donde se encontró seamaeplayo
func analyzePartyArea(memClient *memory.Client) {
	logger.Info("🔍 ANÁLISIS ESPECÍFICO DEL ÁREA DE PARTY")
	logger.Info("")

	// Área específica donde se encontró seamaeplayo
	// Base: 0x10DCE10 (HP/SP), offset -342 = 0x10dccba
	baseAddr := uintptr(0x10dccba)

	logger.Info("📍 ANÁLISIS DETALLADO DEL ÁREA 0x10dccba")
	logger.Info("   (Área donde se encontró 'seamaeplayo' cerca de coordenadas X)")
	logger.Info("")

	// Leer un bloque grande de memoria alrededor de esta área
	blockSize := 2048 // 2KB de datos
	logger.Info("Leyendo %d bytes desde 0x%x...", blockSize, baseAddr)

	data, err := memClient.ReadBytes(baseAddr, blockSize)
	if err != nil {
		logger.Error("Error leyendo memoria: %v", err)
		return
	}

	// Convertir a string para análisis
	dataStr := string(data)

	// Buscar strings específicos de party
	partyStrings := []string{
		"seamaeplayo",
		"party",
		"Party",
		"PARTY",
		"member",
		"Member",
		"MEMBER",
		"leader",
		"Leader",
		"LEADER",
		"name",
		"Name",
		"NAME",
		"char",
		"Char",
		"CHAR",
		"player",
		"Player",
		"PLAYER",
	}

	logger.Info("🔍 BUSCANDO STRINGS DE PARTY EN EL ÁREA:")
	logger.Info("")

	foundStrings := make(map[string][]int)

	for _, searchStr := range partyStrings {
		positions := findStringPositions(dataStr, searchStr)
		if len(positions) > 0 {
			foundStrings[searchStr] = positions
			logger.Info("✅ String '%s' encontrado en %d posiciones:", searchStr, len(positions))
			for i, pos := range positions {
				if i < 3 { // Mostrar solo las primeras 3 posiciones
					addr := baseAddr + uintptr(pos)
					context := getStringContext(dataStr, pos, 64)
					logger.Info("   Posición %d: 0x%x - Contexto: '%s'", pos, addr, context)
				}
			}
			if len(positions) > 3 {
				logger.Info("   ... y %d posiciones más", len(positions)-3)
			}
			logger.Info("")
		}
	}

	// Análisis de coordenadas
	logger.Info("📍 ANÁLISIS DE COORDENADAS:")
	logger.Info("")

	// Buscar patrones de coordenadas (números)
	coordinatePatterns := []string{"X", "Y", "Z", "x", "y", "z"}

	for _, coord := range coordinatePatterns {
		positions := findStringPositions(dataStr, coord)
		if len(positions) > 0 {
			logger.Info("✅ Coordenada '%s' encontrada en %d posiciones:", coord, len(positions))
			for i, pos := range positions {
				if i < 2 { // Mostrar solo las primeras 2 posiciones
					addr := baseAddr + uintptr(pos)
					context := getStringContext(dataStr, pos, 32)
					logger.Info("   Posición %d: 0x%x - Contexto: '%s'", pos, addr, context)
				}
			}
			if len(positions) > 2 {
				logger.Info("   ... y %d posiciones más", len(positions)-2)
			}
			logger.Info("")
		}
	}

	// Análisis de valores numéricos
	logger.Info("🔢 ANÁLISIS DE VALORES NUMÉRICOS:")
	logger.Info("")

	// Leer valores como enteros de 32 bits
	logger.Info("Leyendo valores enteros de 32 bits en el área:")
	for i := 0; i < blockSize/4; i++ {
		addr := baseAddr + uintptr(i*4)
		value, err := memClient.ReadUint32(addr)
		if err == nil {
			// Mostrar solo valores interesantes (coordenadas típicas, IDs, etc.)
			if value > 0 && (value < 1000000 || value > 0x10000000) {
				logger.Info("   Offset %3d (0x%x): %10d (0x%08X)", i*4, addr, value, value)
			}
		}
	}

	// Análisis de estructura
	logger.Info("")
	logger.Info("🏗️ ANÁLISIS DE ESTRUCTURA:")
	logger.Info("")

	if len(foundStrings) > 0 {
		logger.Info("✅ ESTRUCTURA DETECTADA:")
		logger.Info("   - Se encontraron %d tipos de strings diferentes", len(foundStrings))
		logger.Info("   - El área contiene datos de party/coordenadas")
		logger.Info("   - Posible estructura de datos del personaje o party")

		// Buscar el string más común
		maxCount := 0
		mostCommon := ""
		for str, positions := range foundStrings {
			if len(positions) > maxCount {
				maxCount = len(positions)
				mostCommon = str
			}
		}

		if mostCommon != "" {
			logger.Info("   - String más común: '%s' (%d ocurrencias)", mostCommon, maxCount)
		}
	} else {
		logger.Info("❌ No se detectó estructura clara en esta área")
	}

	// Recomendaciones
	logger.Info("")
	logger.Info("💡 RECOMENDACIONES:")
	logger.Info("   1. El área 0x10dccba contiene datos de party/coordenadas")
	logger.Info("   2. 'seamaeplayo' es un nombre de miembro de party")
	logger.Info("   3. Las coordenadas X están cerca de los nombres de party")
	logger.Info("   4. Esta área puede contener la estructura completa del personaje")
	logger.Info("   5. Buscar en offsets específicos para encontrar ASPD")
	logger.Info("   6. Analizar los valores numéricos para encontrar coordenadas reales")
}

// findStringPositions encuentra todas las posiciones de un string en los datos
func findStringPositions(data, searchStr string) []int {
	var positions []int
	start := 0

	for {
		pos := strings.Index(data[start:], searchStr)
		if pos == -1 {
			break
		}
		positions = append(positions, start+pos)
		start += pos + 1
	}

	return positions
}

// getStringContext obtiene el contexto alrededor de una posición
func getStringContext(data string, pos, contextSize int) string {
	start := pos - contextSize/2
	if start < 0 {
		start = 0
	}
	end := pos + contextSize/2
	if end > len(data) {
		end = len(data)
	}

	context := data[start:end]
	// Limpiar caracteres no imprimibles
	cleanContext := ""
	for _, c := range context {
		if c >= 32 && c <= 126 {
			cleanContext += string(c)
		} else {
			cleanContext += "."
		}
	}

	return cleanContext
}

