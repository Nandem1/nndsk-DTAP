package main

import (
	"fmt"
	"unsafe"

	"dtap-go/internal/process"
	"dtap-go/pkg/constants"

	"golang.org/x/sys/windows"
)

type CharacterInspector struct {
	processHandle windows.Handle
}

func NewCharacterInspector(handle windows.Handle) *CharacterInspector {
	return &CharacterInspector{
		processHandle: handle,
	}
}

func (ci *CharacterInspector) InspectCharacterInfo(baseAddress uintptr) error {
	fmt.Printf("=== INSPECCIÓN COMPLETA DEL PERSONAJE ===\n")
	fmt.Printf("Dirección base: 0x%x\n", baseAddress)
	fmt.Println()

	// Leer diferentes rangos de memoria alrededor de la dirección
	ranges := []struct {
		name   string
		offset int64
		size   uintptr
	}{
		{"Información básica (-64 bytes)", -64, 128},
		{"Datos del personaje (-32 bytes)", -32, 64},
		{"Información central (0 bytes)", 0, 64},
		{"Datos extendidos (+32 bytes)", 32, 64},
		{"Información adicional (+64 bytes)", 64, 64},
		{"Datos completos (-128 bytes)", -128, 256},
	}

	for _, r := range ranges {
		fmt.Printf("--- %s ---\n", r.name)
		address := uintptr(int64(baseAddress) + r.offset)

		// Leer memoria
		buffer := make([]byte, r.size)
		var bytesRead uintptr

		err := windows.ReadProcessMemory(ci.processHandle, address, &buffer[0], r.size, &bytesRead)
		if err != nil {
			fmt.Printf("Error leyendo memoria en 0x%x: %v\n", address, err)
			continue
		}

		// Mostrar como bytes hex
		fmt.Printf("Dirección: 0x%x (tamaño: %d bytes)\n", address, bytesRead)
		fmt.Printf("Bytes: ")
		for i, b := range buffer[:bytesRead] {
			if i%16 == 0 && i > 0 {
				fmt.Printf("\n       ")
			}
			fmt.Printf("%02x ", b)
		}
		fmt.Println()

		// Mostrar como strings
		fmt.Printf("Strings: ")
		ci.extractStrings(buffer[:bytesRead])
		fmt.Println()

		// Mostrar como enteros
		fmt.Printf("Enteros: ")
		ci.extractIntegers(buffer[:bytesRead])
		fmt.Println()

		// Mostrar como floats
		fmt.Printf("Floats: ")
		ci.extractFloats(buffer[:bytesRead])
		fmt.Println()

		fmt.Println()
	}

	return nil
}

func (ci *CharacterInspector) extractStrings(data []byte) {
	var current string
	for _, b := range data {
		if b >= 32 && b <= 126 { // Caracteres imprimibles
			current += string(b)
		} else {
			if len(current) >= 3 { // Solo strings de 3+ caracteres
				fmt.Printf("'%s' ", current)
			}
			current = ""
		}
	}
	if len(current) >= 3 {
		fmt.Printf("'%s' ", current)
	}
}

func (ci *CharacterInspector) extractIntegers(data []byte) {
	// Mostrar enteros de 32 bits (little-endian)
	for i := 0; i <= len(data)-4; i += 4 {
		value := uint32(data[i]) |
			uint32(data[i+1])<<8 |
			uint32(data[i+2])<<16 |
			uint32(data[i+3])<<24

		// Solo mostrar valores que parezcan relevantes (no 0, no muy grandes)
		if value > 0 && value < 1000000 {
			fmt.Printf("%d ", value)
		}
	}
}

func (ci *CharacterInspector) extractFloats(data []byte) {
	// Mostrar floats de 32 bits
	for i := 0; i <= len(data)-4; i += 4 {
		// Convertir bytes a float32
		bits := uint32(data[i]) |
			uint32(data[i+1])<<8 |
			uint32(data[i+2])<<16 |
			uint32(data[i+3])<<24

		value := *(*float32)(unsafe.Pointer(&bits))

		// Solo mostrar valores que parezcan relevantes
		if value > 0 && value < 10000 {
			fmt.Printf("%.2f ", value)
		}
	}
}

func (ci *CharacterInspector) AnalyzeCharacterStructure(baseAddress uintptr) error {
	fmt.Printf("=== ANÁLISIS ESTRUCTURAL DEL PERSONAJE ===\n")
	fmt.Printf("Dirección base: 0x%x\n", baseAddress)
	fmt.Println()

	// Leer un bloque grande de memoria
	buffer := make([]byte, 512)
	var bytesRead uintptr

	err := windows.ReadProcessMemory(ci.processHandle, baseAddress, &buffer[0], uintptr(len(buffer)), &bytesRead)
	if err != nil {
		return fmt.Errorf("error leyendo memoria: %v", err)
	}

	fmt.Printf("Leyendo %d bytes desde 0x%x\n", bytesRead, baseAddress)
	fmt.Println()

	// Buscar patrones conocidos
	ci.findKnownPatterns(buffer[:bytesRead], baseAddress)

	// Buscar strings de interés
	ci.findInterestingStrings(buffer[:bytesRead], baseAddress)

	// Buscar valores numéricos interesantes
	ci.findInterestingNumbers(buffer[:bytesRead], baseAddress)

	return nil
}

func (ci *CharacterInspector) findKnownPatterns(data []byte, baseAddress uintptr) {
	fmt.Println("--- Patrones Conocidos ---")

	// Buscar "draw"
	if pos := ci.findString(data, "draw"); pos != -1 {
		fmt.Printf("'draw' encontrado en offset +%d (0x%x)\n", pos, baseAddress+uintptr(pos))
	}

	// Buscar "Biochemist"
	if pos := ci.findString(data, "Biochemist"); pos != -1 {
		fmt.Printf("'Biochemist' encontrado en offset +%d (0x%x)\n", pos, baseAddress+uintptr(pos))
	}

	// Buscar otros strings comunes
	commonStrings := []string{"Jobs:", "Level:", "HP:", "SP:", "STR:", "AGI:", "VIT:", "INT:", "DEX:", "LUK:"}
	for _, s := range commonStrings {
		if pos := ci.findString(data, s); pos != -1 {
			fmt.Printf("'%s' encontrado en offset +%d (0x%x)\n", s, pos, baseAddress+uintptr(pos))
		}
	}
	fmt.Println()
}

func (ci *CharacterInspector) findInterestingStrings(data []byte, baseAddress uintptr) {
	fmt.Println("--- Strings Interesantes ---")

	var current string
	for i, b := range data {
		if b >= 32 && b <= 126 {
			current += string(b)
		} else {
			if len(current) >= 4 { // Strings de 4+ caracteres
				fmt.Printf("String '%s' en offset +%d (0x%x)\n", current, i-len(current), baseAddress+uintptr(i-len(current)))
			}
			current = ""
		}
	}
	if len(current) >= 4 {
		fmt.Printf("String '%s' en offset +%d (0x%x)\n", current, len(data)-len(current), baseAddress+uintptr(len(data)-len(current)))
	}
	fmt.Println()
}

func (ci *CharacterInspector) findInterestingNumbers(data []byte, baseAddress uintptr) {
	fmt.Println("--- Números Interesantes ---")

	// Buscar enteros de 32 bits que podrían ser stats
	for i := 0; i <= len(data)-4; i += 4 {
		value := uint32(data[i]) |
			uint32(data[i+1])<<8 |
			uint32(data[i+2])<<16 |
			uint32(data[i+3])<<24

		// Filtrar valores que podrían ser stats de personaje
		if value >= 1 && value <= 999 && (value == 1 || value == 10 || value == 50 || value == 99 || value == 100) {
			fmt.Printf("Posible stat %d en offset +%d (0x%x)\n", value, i, baseAddress+uintptr(i))
		}
	}
	fmt.Println()
}

func (ci *CharacterInspector) findString(data []byte, target string) int {
	targetBytes := []byte(target)
	for i := 0; i <= len(data)-len(targetBytes); i++ {
		match := true
		for j, b := range targetBytes {
			if data[i+j] != b {
				match = false
				break
			}
		}
		if match {
			return i
		}
	}
	return -1
}

func main() {
	fmt.Println("DTAP-Go - Character Inspector")
	fmt.Println("=============================")

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

	// Crear inspector
	inspector := NewCharacterInspector(handle)

	// Inspeccionar información del personaje
	err = inspector.InspectCharacterInfo(baseAddress)
	if err != nil {
		fmt.Printf("Error inspeccionando personaje: %v\n", err)
		fmt.Scanln()
		return
	}

	// Analizar estructura
	err = inspector.AnalyzeCharacterStructure(baseAddress)
	if err != nil {
		fmt.Printf("Error analizando estructura: %v\n", err)
		fmt.Scanln()
		return
	}

	fmt.Println("=== INSPECCIÓN COMPLETADA ===")
	fmt.Println("Presiona Enter para salir...")
	fmt.Scanln()
}
