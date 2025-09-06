package main

import (
	"fmt"

	"dtap-go/internal/process"
	"dtap-go/pkg/constants"

	"golang.org/x/sys/windows"
)

type CharacterData struct {
	Name        string
	Level       int
	JobLevel    int
	BaseLevel   int
	HP          int
	SP          int
	STR         int
	AGI         int
	VIT         int
	INT         int
	DEX         int
	LUK         int
	MapName     string
	CharacterID int
}

func main() {
	fmt.Println("DTAP-Go - Character Parser")
	fmt.Println("==========================")

	// Buscar proceso de Ragnarok Online
	roProcess, err := process.FindROProcess()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		fmt.Println("Asegúrate de que Ragnarok Online esté ejecutándose.")
		fmt.Scanln()
		return
	}

	fmt.Printf("Proceso encontrado: %s (PID: %d)\n", roProcess.Name, roProcess.PID)

	// Abrir proceso
	handle, err := windows.OpenProcess(constants.ProcessAccessFlags, false, roProcess.PID)
	if err != nil {
		fmt.Printf("Error abriendo proceso: %v\n", err)
		fmt.Scanln()
		return
	}
	defer windows.CloseHandle(handle)

	fmt.Println("Proceso abierto exitosamente")
	fmt.Println()

	// Pedir dirección al usuario
	fmt.Print("Ingresa la dirección base del personaje (ej: 0x10df5d8): ")
	var addressStr string
	fmt.Scanln(&addressStr)

	var baseAddress uintptr
	_, err = fmt.Sscanf(addressStr, "0x%x", &baseAddress)
	if err != nil {
		_, err = fmt.Sscanf(addressStr, "%x", &baseAddress)
		if err != nil {
			fmt.Printf("Error parseando dirección: %v\n", err)
			fmt.Scanln()
			return
		}
	}

	// Leer datos del personaje
	characterData, err := parseCharacterData(handle, baseAddress)
	if err != nil {
		fmt.Printf("Error parseando datos: %v\n", err)
		fmt.Scanln()
		return
	}

	// Mostrar información del personaje
	fmt.Println("=== INFORMACIÓN DEL PERSONAJE ===")
	fmt.Printf("Nombre: %s\n", characterData.Name)
	fmt.Printf("Nivel Base: %d\n", characterData.BaseLevel)
	fmt.Printf("Nivel Job: %d\n", characterData.JobLevel)
	fmt.Printf("Nivel Total: %d\n", characterData.Level)
	fmt.Printf("HP: %d\n", characterData.HP)
	fmt.Printf("SP: %d\n", characterData.SP)
	fmt.Println()
	fmt.Printf("=== STATS ===\n")
	fmt.Printf("STR: %d\n", characterData.STR)
	fmt.Printf("AGI: %d\n", characterData.AGI)
	fmt.Printf("VIT: %d\n", characterData.VIT)
	fmt.Printf("INT: %d\n", characterData.INT)
	fmt.Printf("DEX: %d\n", characterData.DEX)
	fmt.Printf("LUK: %d\n", characterData.LUK)
	fmt.Println()
	fmt.Printf("=== INFORMACIÓN ADICIONAL ===\n")
	fmt.Printf("ID del Personaje: %d\n", characterData.CharacterID)
	fmt.Printf("Mapa: %s\n", characterData.MapName)
	fmt.Println()

	fmt.Println("Presiona Enter para salir...")
	fmt.Scanln()
}

func parseCharacterData(handle windows.Handle, baseAddress uintptr) (*CharacterData, error) {
	// Leer un bloque grande de memoria
	buffer := make([]byte, 256)
	var bytesRead uintptr

	err := windows.ReadProcessMemory(handle, baseAddress, &buffer[0], uintptr(len(buffer)), &bytesRead)
	if err != nil {
		return nil, fmt.Errorf("error leyendo memoria: %v", err)
	}

	character := &CharacterData{}

	// Extraer nombre (primeros 20 bytes)
	nameBytes := make([]byte, 0)
	for i := 0; i < 20 && i < len(buffer); i++ {
		if buffer[i] == 0 {
			break
		}
		if buffer[i] >= 32 && buffer[i] <= 126 {
			nameBytes = append(nameBytes, buffer[i])
		}
	}
	character.Name = string(nameBytes)

	// Extraer números que podrían ser stats
	// Basado en el análisis anterior, los números están en posiciones específicas
	numbers := extractNumbers(buffer)

	if len(numbers) >= 6 {
		// Interpretar los números basado en patrones comunes de RO
		character.BaseLevel = numbers[0]          // 11
		character.JobLevel = numbers[1]           // 15
		character.Level = numbers[0] + numbers[1] // Nivel total
		character.STR = numbers[2]                // 1
		character.CharacterID = numbers[3]        // 30066
		character.HP = numbers[4]                 // 119
		character.SP = numbers[5]                 // Otro valor
	}

	// Extraer nombre del mapa
	mapName := extractMapName(buffer)
	character.MapName = mapName

	// Valores por defecto si no se encontraron
	if character.Name == "" {
		character.Name = "Desconocido"
	}
	if character.BaseLevel == 0 {
		character.BaseLevel = 1
	}
	if character.JobLevel == 0 {
		character.JobLevel = 1
	}
	if character.Level == 0 {
		character.Level = character.BaseLevel + character.JobLevel
	}

	return character, nil
}

func extractNumbers(data []byte) []int {
	var numbers []int

	// Buscar enteros de 32 bits (little-endian)
	for i := 0; i <= len(data)-4; i += 4 {
		value := uint32(data[i]) |
			uint32(data[i+1])<<8 |
			uint32(data[i+2])<<16 |
			uint32(data[i+3])<<24

		// Filtrar valores que podrían ser stats de personaje
		if value > 0 && value < 100000 {
			numbers = append(numbers, int(value))
		}
	}

	return numbers
}

func extractMapName(data []byte) string {
	// Buscar strings que terminen en .rsw (archivos de mapa de RO)
	for i := 0; i < len(data)-4; i++ {
		if data[i] == '.' &&
			i+3 < len(data) &&
			data[i+1] == 'r' &&
			data[i+2] == 's' &&
			data[i+3] == 'w' {

			// Encontrar el inicio del string
			start := i
			for start > 0 && data[start-1] != 0 {
				start--
			}

			// Extraer el string
			var mapName []byte
			for j := start; j < i+4 && j < len(data); j++ {
				if data[j] >= 32 && data[j] <= 126 {
					mapName = append(mapName, data[j])
				}
			}

			if len(mapName) > 0 {
				return string(mapName)
			}
		}
	}

	return "Desconocido"
}
