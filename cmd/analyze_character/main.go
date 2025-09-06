package main

import (
	"fmt"

	"dtap-go/internal/process"
	"dtap-go/pkg/constants"

	"golang.org/x/sys/windows"
)

type CharacterAnalyzer struct {
	processHandle windows.Handle
}

func NewCharacterAnalyzer(handle windows.Handle) *CharacterAnalyzer {
	return &CharacterAnalyzer{
		processHandle: handle,
	}
}

func (ca *CharacterAnalyzer) AnalyzeCharacterStructure(baseAddress uintptr) error {
	fmt.Printf("=== ANÁLISIS ESTRUCTURAL AVANZADO ===\n")
	fmt.Printf("Dirección base: 0x%x\n", baseAddress)
	fmt.Println()

	// Leer un bloque grande de memoria
	buffer := make([]byte, 1024)
	var bytesRead uintptr

	err := windows.ReadProcessMemory(ca.processHandle, baseAddress, &buffer[0], uintptr(len(buffer)), &bytesRead)
	if err != nil {
		return fmt.Errorf("error leyendo memoria: %v", err)
	}

	fmt.Printf("Analizando %d bytes desde 0x%x\n", bytesRead, baseAddress)
	fmt.Println()

	// Buscar estructuras específicas
	ca.findJobID(buffer[:bytesRead], baseAddress)
	ca.findClassID(buffer[:bytesRead], baseAddress)
	ca.findCharacterID(buffer[:bytesRead], baseAddress)
	ca.findLevels(buffer[:bytesRead], baseAddress)
	ca.findStats(buffer[:bytesRead], baseAddress)
	ca.findHP_SP(buffer[:bytesRead], baseAddress)
	ca.findPosition(buffer[:bytesRead], baseAddress)

	return nil
}

func (ca *CharacterAnalyzer) findJobID(data []byte, baseAddress uintptr) {
	fmt.Println("--- JOB ID / CLASS ID ---")

	// Buscar patrones que podrían ser Job ID
	// En RO, los Job IDs suelen ser valores específicos
	jobIDs := map[uint32]string{
		0: "Novice",
		1: "Swordman", 2: "Thief", 3: "Acolyte", 4: "Mage", 5: "Archer", 6: "Merchant",
		7: "Knight", 8: "Assassin", 9: "Priest", 10: "Wizard", 11: "Hunter", 12: "Blacksmith",
		13: "Crusader", 14: "Rogue", 15: "Monk", 16: "Sage", 17: "Bard", 18: "Dancer", 19: "Alchemist",
		20: "Super Novice", 21: "Gunslinger", 22: "Ninja", 23: "Taekwon", 24: "Star Gladiator",
		25: "Soul Linker", 26: "Gangsi", 27: "Death Knight", 28: "Dark Collector",
		4008: "High Novice", 4009: "High Swordman", 4010: "High Thief", 4011: "High Acolyte",
		4012: "High Mage", 4013: "High Archer", 4014: "High Merchant",
		4015: "Lord Knight", 4016: "Assassin Cross", 4017: "High Priest", 4018: "High Wizard",
		4019: "Sniper", 4020: "Whitesmith", 4021: "Paladin", 4022: "Stalker", 4023: "Champion",
		4024: "Professor", 4025: "Clown", 4026: "Gypsy", 4027: "Creator",
	}

	// Buscar enteros de 32 bits que coincidan con Job IDs conocidos
	for i := 0; i <= len(data)-4; i += 4 {
		value := uint32(data[i]) |
			uint32(data[i+1])<<8 |
			uint32(data[i+2])<<16 |
			uint32(data[i+3])<<24

		if jobName, exists := jobIDs[value]; exists {
			fmt.Printf("✓ Job ID encontrado: %d (%s) en offset +%d (0x%x)\n",
				value, jobName, i, baseAddress+uintptr(i))
		}
	}
	fmt.Println()
}

func (ca *CharacterAnalyzer) findClassID(data []byte, baseAddress uintptr) {
	fmt.Println("--- CLASS ID / RACE ID ---")

	// Buscar Race IDs (Human, Doram, etc.)
	raceIDs := map[uint32]string{
		0: "Human",
		1: "Doram",
	}

	// Buscar Class IDs específicos
	classIDs := map[uint32]string{
		0: "Normal",
		1: "Baby",
		2: "Third Class",
		3: "Extended",
	}

	for i := 0; i <= len(data)-4; i += 4 {
		value := uint32(data[i]) |
			uint32(data[i+1])<<8 |
			uint32(data[i+2])<<16 |
			uint32(data[i+3])<<24

		if raceName, exists := raceIDs[value]; exists {
			fmt.Printf("✓ Race ID encontrado: %d (%s) en offset +%d (0x%x)\n",
				value, raceName, i, baseAddress+uintptr(i))
		}

		if className, exists := classIDs[value]; exists {
			fmt.Printf("✓ Class ID encontrado: %d (%s) en offset +%d (0x%x)\n",
				value, className, i, baseAddress+uintptr(i))
		}
	}
	fmt.Println()
}

func (ca *CharacterAnalyzer) findCharacterID(data []byte, baseAddress uintptr) {
	fmt.Println("--- CHARACTER ID ---")

	// Buscar valores que podrían ser Character ID
	// Los Character IDs suelen ser números grandes y únicos
	for i := 0; i <= len(data)-4; i += 4 {
		value := uint32(data[i]) |
			uint32(data[i+1])<<8 |
			uint32(data[i+2])<<16 |
			uint32(data[i+3])<<24

		// Character IDs típicamente están en rangos específicos
		if value > 1000 && value < 1000000 {
			fmt.Printf("Posible Character ID: %d en offset +%d (0x%x)\n",
				value, i, baseAddress+uintptr(i))
		}
	}
	fmt.Println()
}

func (ca *CharacterAnalyzer) findLevels(data []byte, baseAddress uintptr) {
	fmt.Println("--- LEVELS ---")

	// Buscar valores que podrían ser niveles
	for i := 0; i <= len(data)-4; i += 4 {
		value := uint32(data[i]) |
			uint32(data[i+1])<<8 |
			uint32(data[i+2])<<16 |
			uint32(data[i+3])<<24

		// Niveles típicamente están entre 1 y 175
		if value >= 1 && value <= 175 {
			fmt.Printf("Posible Level: %d en offset +%d (0x%x)\n",
				value, i, baseAddress+uintptr(i))
		}
	}
	fmt.Println()
}

func (ca *CharacterAnalyzer) findStats(data []byte, baseAddress uintptr) {
	fmt.Println("--- STATS (STR, AGI, VIT, INT, DEX, LUK) ---")

	// Buscar valores que podrían ser stats
	for i := 0; i <= len(data)-4; i += 4 {
		value := uint32(data[i]) |
			uint32(data[i+1])<<8 |
			uint32(data[i+2])<<16 |
			uint32(data[i+3])<<24

		// Stats típicamente están entre 1 y 120
		if value >= 1 && value <= 120 {
			fmt.Printf("Posible Stat: %d en offset +%d (0x%x)\n",
				value, i, baseAddress+uintptr(i))
		}
	}
	fmt.Println()
}

func (ca *CharacterAnalyzer) findHP_SP(data []byte, baseAddress uintptr) {
	fmt.Println("--- HP / SP ---")

	// Buscar valores que podrían ser HP/SP
	for i := 0; i <= len(data)-4; i += 4 {
		value := uint32(data[i]) |
			uint32(data[i+1])<<8 |
			uint32(data[i+2])<<16 |
			uint32(data[i+3])<<24

		// HP/SP pueden ser valores grandes
		if value > 100 && value < 100000 {
			fmt.Printf("Posible HP/SP: %d en offset +%d (0x%x)\n",
				value, i, baseAddress+uintptr(i))
		}
	}
	fmt.Println()
}

func (ca *CharacterAnalyzer) findPosition(data []byte, baseAddress uintptr) {
	fmt.Println("--- POSITION (X, Y, Z) ---")

	// Buscar coordenadas (suelen ser valores grandes)
	for i := 0; i <= len(data)-4; i += 4 {
		value := uint32(data[i]) |
			uint32(data[i+1])<<8 |
			uint32(data[i+2])<<16 |
			uint32(data[i+3])<<24

		// Coordenadas típicamente están en rangos específicos
		if value > 1000 && value < 1000000 {
			fmt.Printf("Posible Coordenada: %d en offset +%d (0x%x)\n",
				value, i, baseAddress+uintptr(i))
		}
	}
	fmt.Println()
}

func main() {
	fmt.Println("DTAP-Go - Character Structure Analyzer")
	fmt.Println("======================================")

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

	// Crear analizador
	analyzer := NewCharacterAnalyzer(handle)

	// Analizar estructura del personaje
	err = analyzer.AnalyzeCharacterStructure(baseAddress)
	if err != nil {
		fmt.Printf("Error analizando estructura: %v\n", err)
		fmt.Scanln()
		return
	}

	fmt.Println("=== ANÁLISIS COMPLETADO ===")
	fmt.Println("Presiona Enter para salir...")
	fmt.Scanln()
}
