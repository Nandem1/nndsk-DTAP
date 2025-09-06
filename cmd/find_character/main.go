package main

import (
	"fmt"
	"unsafe"

	"dtap-go/internal/process"
	"dtap-go/pkg/constants"

	"golang.org/x/sys/windows"
)

type CharacterFinder struct {
	processHandle windows.Handle
}

func NewCharacterFinder(handle windows.Handle) *CharacterFinder {
	return &CharacterFinder{
		processHandle: handle,
	}
}

func (cf *CharacterFinder) FindCharacterInfo() ([]CharacterInfo, error) {
	var results []CharacterInfo

	fmt.Println("Buscando personaje 'Sokoninaru' clase 'Champion'...")
	fmt.Println("(Nombres más específicos y únicos)")

	// Buscar el nuevo nombre del personaje
	nameAddresses, err := cf.findStringSimple("Sokoninaru")
	if err != nil {
		return nil, fmt.Errorf("error buscando 'Sokoninaru': %v", err)
	}

	// Buscar la nueva clase del personaje
	classAddresses, err := cf.findStringSimple("Champion")
	if err != nil {
		return nil, fmt.Errorf("error buscando 'Champion': %v", err)
	}

	fmt.Printf("Encontradas %d direcciones con 'Sokoninaru'\n", len(nameAddresses))
	fmt.Printf("Encontradas %d direcciones con 'Champion'\n", len(classAddresses))
	fmt.Println()

	// Mostrar TODOS los resultados
	allAddresses := append(nameAddresses, classAddresses...)

	for i, addr := range allAddresses {
		context, _ := cf.ReadContext(addr)

		results = append(results, CharacterInfo{
			NameAddress:  addr,
			ClassAddress: 0,
			Distance:     0,
			NameContext:  context,
			ClassContext: "",
		})

		fmt.Printf("--- Dirección %d ---\n", i+1)
		fmt.Printf("Address: 0x%x\n", addr)
		fmt.Printf("Contexto: %s\n", context)
		fmt.Println()
	}

	return results, nil
}

func (cf *CharacterFinder) findStringSimple(target string) ([]uintptr, error) {
	var addresses []uintptr
	targetBytes := []byte(target)

	fmt.Printf("Buscando '%s' usando la MISMA estrategia que find_biochemist...\n", target)

	// COPIAR EXACTAMENTE la lógica del find_biochemist que SÍ funciona
	var mbi windows.MemoryBasicInformation
	address := uintptr(0)
	regionsChecked := 0
	maxRegions := 2000
	regionsSearched := 0

	for regionsChecked < maxRegions {
		err := windows.VirtualQueryEx(cf.processHandle, address, &mbi, unsafe.Sizeof(mbi))
		if err != nil {
			address += uintptr(mbi.RegionSize)
			regionsChecked++
			continue
		}

		// Mostrar progreso cada 50 regiones
		if regionsChecked%50 == 0 {
			fmt.Printf("Progreso: %d regiones revisadas, %d regiones buscadas...\n", regionsChecked, regionsSearched)
		}

		// MISMA lógica que find_biochemist
		if mbi.State == windows.MEM_COMMIT &&
			(mbi.Protect&windows.PAGE_READONLY != 0 ||
				mbi.Protect&windows.PAGE_READWRITE != 0 ||
				mbi.Protect&windows.PAGE_EXECUTE_READ != 0 ||
				mbi.Protect&windows.PAGE_EXECUTE_READWRITE != 0) {

			// Leer región de memoria (optimizado para regiones grandes)
			regionSize := mbi.RegionSize
			if regionSize > 50*1024*1024 { // Si es muy grande, leer solo los primeros 50MB
				regionSize = 50 * 1024 * 1024
			}

			buffer := make([]byte, regionSize)
			var bytesRead uintptr

			err := windows.ReadProcessMemory(cf.processHandle, uintptr(mbi.BaseAddress), &buffer[0], uintptr(len(buffer)), &bytesRead)
			if err == nil && bytesRead > 0 {
				regionsSearched++
				// Buscar string en el buffer
				foundAddresses := cf.searchInBuffer(buffer, targetBytes, uintptr(mbi.BaseAddress))

				// NO filtrar - mostrar TODOS los resultados
				for _, addr := range foundAddresses {
					addresses = append(addresses, addr)
					fmt.Printf("✓ Encontrada ocurrencia en región 0x%x (tamaño: %d bytes)\n", mbi.BaseAddress, mbi.RegionSize)
				}
			}
		}

		address += uintptr(mbi.RegionSize)
		regionsChecked++
	}

	fmt.Printf("Búsqueda completada. Revisadas %d regiones, buscadas %d regiones.\n", regionsChecked, regionsSearched)
	return addresses, nil
}

func (cf *CharacterFinder) stringToHex(s string) uint32 {
	// Convertir string a valor hexadecimal de 32 bits
	// "draw" -> 0x64726177 (little-endian)
	// "Biochemist" -> 0x42696F6368656D697374 (truncado a 32 bits)

	if len(s) == 0 {
		return 0
	}

	// Tomar los primeros 4 caracteres para un uint32
	bytes := []byte(s)
	if len(bytes) > 4 {
		bytes = bytes[:4]
	}

	// Convertir a uint32 (little-endian)
	result := uint32(0)
	for i, b := range bytes {
		result |= uint32(b) << (i * 8)
	}

	return result
}

func (cf *CharacterFinder) findHexValue(targetValue uint32) ([]uintptr, error) {
	var addresses []uintptr

	fmt.Printf("Buscando valor 0x%x en TODA la memoria...\n", targetValue)

	// Enumerar TODAS las regiones de memoria
	var mbi windows.MemoryBasicInformation
	address := uintptr(0)
	regionsChecked := 0
	maxRegions := 50000
	regionsSearched := 0

	for regionsChecked < maxRegions {
		err := windows.VirtualQueryEx(cf.processHandle, address, &mbi, unsafe.Sizeof(mbi))
		if err != nil {
			address += uintptr(mbi.RegionSize)
			regionsChecked++
			continue
		}

		// Mostrar progreso cada 100 regiones
		if regionsChecked%100 == 0 {
			fmt.Printf("Progreso: %d regiones revisadas, %d regiones buscadas...\n", regionsChecked, regionsSearched)
		}

		// Buscar en TODAS las regiones (sin filtros)
		if mbi.State == windows.MEM_COMMIT {
			// Leer región de memoria
			regionSize := mbi.RegionSize
			if regionSize > 100*1024*1024 { // Límite de 100MB por región
				regionSize = 100 * 1024 * 1024
			}

			buffer := make([]byte, regionSize)
			var bytesRead uintptr

			err := windows.ReadProcessMemory(cf.processHandle, uintptr(mbi.BaseAddress), &buffer[0], uintptr(len(buffer)), &bytesRead)
			if err == nil && bytesRead > 0 {
				regionsSearched++
				// Buscar valor de 32 bits en el buffer
				foundAddresses := cf.searchUint32Value(buffer, targetValue, uintptr(mbi.BaseAddress))
				addresses = append(addresses, foundAddresses...)

				if len(foundAddresses) > 0 {
					fmt.Printf("✓ Encontradas %d ocurrencias del valor 0x%x en región 0x%x\n", len(foundAddresses), targetValue, mbi.BaseAddress)
				}
			}
		}

		address += uintptr(mbi.RegionSize)
		regionsChecked++
	}

	fmt.Printf("Búsqueda de valor 0x%x completada. Revisadas %d regiones, buscadas %d regiones.\n", targetValue, regionsChecked, regionsSearched)
	return addresses, nil
}

func (cf *CharacterFinder) searchUint32Value(buffer []byte, targetValue uint32, baseAddress uintptr) []uintptr {
	var addresses []uintptr

	// Buscar el valor de 32 bits en el buffer
	for i := 0; i <= len(buffer)-4; i++ {
		// Leer 4 bytes como uint32 (little-endian)
		value := uint32(buffer[i]) |
			uint32(buffer[i+1])<<8 |
			uint32(buffer[i+2])<<16 |
			uint32(buffer[i+3])<<24

		if value == targetValue {
			addresses = append(addresses, baseAddress+uintptr(i))
		}
	}

	return addresses
}

func (cf *CharacterFinder) findByteSequence(target string) ([]uintptr, error) {
	var addresses []uintptr
	targetBytes := []byte(target)

	fmt.Printf("Convirtiendo '%s' a bytes: %v\n", target, targetBytes)
	fmt.Printf("Buscando secuencia de bytes en TODA la memoria...\n")

	// Enumerar TODAS las regiones de memoria
	var mbi windows.MemoryBasicInformation
	address := uintptr(0)
	regionsChecked := 0
	maxRegions := 50000
	regionsSearched := 0

	for regionsChecked < maxRegions {
		err := windows.VirtualQueryEx(cf.processHandle, address, &mbi, unsafe.Sizeof(mbi))
		if err != nil {
			address += uintptr(mbi.RegionSize)
			regionsChecked++
			continue
		}

		// Mostrar progreso cada 100 regiones
		if regionsChecked%100 == 0 {
			fmt.Printf("Progreso: %d regiones revisadas, %d regiones buscadas...\n", regionsChecked, regionsSearched)
		}

		// Buscar en TODAS las regiones (sin filtros)
		if mbi.State == windows.MEM_COMMIT {
			// Leer región de memoria
			regionSize := mbi.RegionSize
			if regionSize > 100*1024*1024 { // Límite de 100MB por región
				regionSize = 100 * 1024 * 1024
			}

			buffer := make([]byte, regionSize)
			var bytesRead uintptr

			err := windows.ReadProcessMemory(cf.processHandle, uintptr(mbi.BaseAddress), &buffer[0], uintptr(len(buffer)), &bytesRead)
			if err == nil && bytesRead > 0 {
				regionsSearched++
				// Buscar secuencia de bytes en el buffer
				foundAddresses := cf.searchByteSequence(buffer, targetBytes, uintptr(mbi.BaseAddress))
				addresses = append(addresses, foundAddresses...)

				if len(foundAddresses) > 0 {
					fmt.Printf("✓ Encontradas %d ocurrencias de secuencia en región 0x%x (tamaño: %d bytes)\n", len(foundAddresses), mbi.BaseAddress, mbi.RegionSize)
				}
			}
		}

		address += uintptr(mbi.RegionSize)
		regionsChecked++
	}

	fmt.Printf("Búsqueda de secuencia de bytes completada. Revisadas %d regiones, buscadas %d regiones.\n", regionsChecked, regionsSearched)
	return addresses, nil
}

func (cf *CharacterFinder) findStringUnrestricted(target string) ([]uintptr, error) {
	var addresses []uintptr
	targetBytes := []byte(target)

	fmt.Printf("Iniciando búsqueda SIN RESTRICCIONES de '%s'...\n", target)

	// Enumerar TODAS las regiones de memoria (sin límites ni restricciones)
	var mbi windows.MemoryBasicInformation
	address := uintptr(0)
	regionsChecked := 0
	maxRegions := 50000 // Límite muy alto
	regionsSearched := 0

	for regionsChecked < maxRegions {
		err := windows.VirtualQueryEx(cf.processHandle, address, &mbi, unsafe.Sizeof(mbi))
		if err != nil {
			address += uintptr(mbi.RegionSize)
			regionsChecked++
			continue
		}

		// Mostrar progreso cada 100 regiones
		if regionsChecked%100 == 0 {
			fmt.Printf("Progreso: %d regiones revisadas, %d regiones buscadas...\n", regionsChecked, regionsSearched)
		}

		// Buscar en TODAS las regiones (sin filtros de protección)
		if mbi.State == windows.MEM_COMMIT {
			// Leer región de memoria (sin límite de tamaño)
			regionSize := mbi.RegionSize
			// Sin límite de tamaño - intentar leer toda la región

			buffer := make([]byte, regionSize)
			var bytesRead uintptr

			err := windows.ReadProcessMemory(cf.processHandle, uintptr(mbi.BaseAddress), &buffer[0], uintptr(len(buffer)), &bytesRead)
			if err == nil && bytesRead > 0 {
				regionsSearched++
				// Buscar string en el buffer
				foundAddresses := cf.searchInBuffer(buffer, targetBytes, uintptr(mbi.BaseAddress))
				addresses = append(addresses, foundAddresses...)

				if len(foundAddresses) > 0 {
					fmt.Printf("✓ Encontradas %d ocurrencias en región 0x%x (tamaño: %d bytes)\n", len(foundAddresses), mbi.BaseAddress, mbi.RegionSize)
				}
			}
		}

		address += uintptr(mbi.RegionSize)
		regionsChecked++
	}

	fmt.Printf("Búsqueda SIN RESTRICCIONES completada. Revisadas %d regiones, buscadas %d regiones.\n", regionsChecked, regionsSearched)
	return addresses, nil
}

func (cf *CharacterFinder) findStringNearAddress(target string, baseAddress uintptr, radius uintptr) ([]uintptr, error) {
	var addresses []uintptr
	targetBytes := []byte(target)

	// Calcular rango de búsqueda
	startAddress := baseAddress - radius
	endAddress := baseAddress + radius

	fmt.Printf("Buscando '%s' en rango: 0x%x - 0x%x (radio: %d bytes)\n", target, startAddress, endAddress, radius)

	// Enumerar regiones de memoria en el rango específico
	var mbi windows.MemoryBasicInformation
	address := startAddress
	regionsChecked := 0
	maxRegions := 1000
	regionsSearched := 0

	for regionsChecked < maxRegions && address < endAddress {
		err := windows.VirtualQueryEx(cf.processHandle, address, &mbi, unsafe.Sizeof(mbi))
		if err != nil {
			address += uintptr(mbi.RegionSize)
			regionsChecked++
			continue
		}

		// Verificar si la región está en nuestro rango de interés
		regionStart := uintptr(mbi.BaseAddress)
		regionEnd := regionStart + uintptr(mbi.RegionSize)

		if regionStart >= endAddress {
			break // Fuera del rango
		}

		if regionEnd <= startAddress {
			address += uintptr(mbi.RegionSize)
			regionsChecked++
			continue
		}

		// Buscar en regiones accesibles
		if mbi.State == windows.MEM_COMMIT &&
			(mbi.Protect&windows.PAGE_READONLY != 0 ||
				mbi.Protect&windows.PAGE_READWRITE != 0 ||
				mbi.Protect&windows.PAGE_EXECUTE_READ != 0 ||
				mbi.Protect&windows.PAGE_EXECUTE_READWRITE != 0) {

			// Leer región de memoria
			regionSize := mbi.RegionSize
			if regionSize > 10*1024*1024 { // Límite más pequeño para búsqueda localizada
				regionSize = 10 * 1024 * 1024
			}

			buffer := make([]byte, regionSize)
			var bytesRead uintptr

			err := windows.ReadProcessMemory(cf.processHandle, uintptr(mbi.BaseAddress), &buffer[0], uintptr(len(buffer)), &bytesRead)
			if err == nil && bytesRead > 0 {
				regionsSearched++
				// Buscar string en el buffer
				foundAddresses := cf.searchInBuffer(buffer, targetBytes, uintptr(mbi.BaseAddress))
				addresses = append(addresses, foundAddresses...)

				if len(foundAddresses) > 0 {
					fmt.Printf("✓ Encontradas %d ocurrencias en región 0x%x\n", len(foundAddresses), mbi.BaseAddress)
				}
			}
		}

		address += uintptr(mbi.RegionSize)
		regionsChecked++
	}

	fmt.Printf("Búsqueda localizada completada. Revisadas %d regiones, buscadas %d regiones.\n", regionsChecked, regionsSearched)
	return addresses, nil
}

func (cf *CharacterFinder) findString(target string) ([]uintptr, error) {
	var addresses []uintptr
	targetBytes := []byte(target)

	// Enumerar regiones de memoria (copiado del find_biochemist que SÍ funciona)
	var mbi windows.MemoryBasicInformation
	address := uintptr(0)
	regionsChecked := 0
	maxRegions := 2000
	regionsSearched := 0

	for regionsChecked < maxRegions {
		err := windows.VirtualQueryEx(cf.processHandle, address, &mbi, unsafe.Sizeof(mbi))
		if err != nil {
			address += uintptr(mbi.RegionSize)
			regionsChecked++
			continue
		}

		// Mostrar progreso cada 50 regiones
		if regionsChecked%50 == 0 {
			fmt.Printf("Progreso: %d regiones revisadas, %d regiones buscadas...\n", regionsChecked, regionsSearched)
		}

		// Buscar en todas las regiones accesibles del juego (igual que find_biochemist)
		if mbi.State == windows.MEM_COMMIT &&
			(mbi.Protect&windows.PAGE_READONLY != 0 ||
				mbi.Protect&windows.PAGE_READWRITE != 0 ||
				mbi.Protect&windows.PAGE_EXECUTE_READ != 0 ||
				mbi.Protect&windows.PAGE_EXECUTE_READWRITE != 0) {

			// Leer región de memoria (optimizado para regiones grandes)
			regionSize := mbi.RegionSize
			if regionSize > 50*1024*1024 { // Si es muy grande, leer solo los primeros 50MB
				regionSize = 50 * 1024 * 1024
			}

			buffer := make([]byte, regionSize)
			var bytesRead uintptr

			err := windows.ReadProcessMemory(cf.processHandle, uintptr(mbi.BaseAddress), &buffer[0], uintptr(len(buffer)), &bytesRead)
			if err == nil && bytesRead > 0 {
				regionsSearched++
				// Buscar string en el buffer
				foundAddresses := cf.searchInBuffer(buffer, targetBytes, uintptr(mbi.BaseAddress))
				addresses = append(addresses, foundAddresses...)
			}
		}

		address += uintptr(mbi.RegionSize)
		regionsChecked++
	}

	fmt.Printf("Búsqueda completada. Revisadas %d regiones, buscadas %d regiones.\n", regionsChecked, regionsSearched)
	return addresses, nil
}

func (cf *CharacterFinder) searchByteSequence(buffer []byte, target []byte, baseAddress uintptr) []uintptr {
	var addresses []uintptr

	for i := 0; i <= len(buffer)-len(target); i++ {
		if bytesEqual(buffer[i:i+len(target)], target) {
			addresses = append(addresses, baseAddress+uintptr(i))
		}
	}

	return addresses
}

func (cf *CharacterFinder) searchInBuffer(buffer []byte, target []byte, baseAddress uintptr) []uintptr {
	var addresses []uintptr

	for i := 0; i <= len(buffer)-len(target); i++ {
		if bytesEqual(buffer[i:i+len(target)], target) {
			addresses = append(addresses, baseAddress+uintptr(i))
		}
	}

	return addresses
}

func (cf *CharacterFinder) ReadContext(address uintptr) (string, error) {
	buffer := make([]byte, 64)
	var bytesRead uintptr

	err := windows.ReadProcessMemory(cf.processHandle, address, &buffer[0], uintptr(len(buffer)), &bytesRead)
	if err != nil {
		return "", err
	}

	// Convertir a string, reemplazando caracteres no imprimibles
	result := ""
	for _, b := range buffer {
		if b >= 32 && b <= 126 {
			result += string(b)
		} else {
			result += "."
		}
	}

	return result, nil
}

func (cf *CharacterFinder) isSelfProgramContext(context string) bool {
	selfPhrases := []string{
		"DTAP-Go",
		"Character Finder",
		"Buscando string",
		"Encontradas",
		"Dirección",
		"Contexto:",
		"=== RESULTADOS ===",
		"=== INSTRUCCIONES ===",
	}

	for _, phrase := range selfPhrases {
		if contains(context, phrase) {
			return true
		}
	}

	return false
}

type CharacterInfo struct {
	NameAddress  uintptr
	ClassAddress uintptr
	Distance     uintptr
	NameContext  string
	ClassContext string
}

func bytesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr
}

func main() {
	fmt.Println("DTAP-Go - Character Info Finder")
	fmt.Println("===============================")

	// Buscar proceso de Ragnarok Online
	roProcess, err := process.FindROProcess()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		fmt.Println("Asegúrate de que Ragnarok Online esté ejecutándose.")
		fmt.Scanln()
		return
	}

	fmt.Printf("Proceso encontrado: %s (PID: %d)\n", roProcess.Name, roProcess.PID)

	// Abrir proceso con los MISMOS permisos que find_biochemist
	handle, err := windows.OpenProcess(constants.ProcessAccessFlags, false, roProcess.PID)
	if err != nil {
		fmt.Printf("Error abriendo proceso: %v\n", err)
		fmt.Scanln()
		return
	}
	defer windows.CloseHandle(handle)

	fmt.Println("Proceso abierto exitosamente")
	fmt.Println("Buscando personaje 'Sokoninaru' clase 'Champion'...")
	fmt.Println("(Nombres más específicos y únicos)")
	fmt.Println("(Búsqueda simple de strings sin complicaciones)")

	// Buscar información del personaje
	characterFinder := NewCharacterFinder(handle)
	results, err := characterFinder.FindCharacterInfo()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		fmt.Scanln()
		return
	}

	if len(results) == 0 {
		fmt.Println("No se encontró información del personaje")
		fmt.Println("Asegúrate de que tu personaje sea 'Sokoninaru' y clase 'Champion'")
		fmt.Scanln()
		return
	}

	fmt.Printf("\n=== RESULTADOS ===\n")
	fmt.Printf("Encontradas %d coincidencias de información del personaje:\n\n", len(results))

	for i, result := range results {
		fmt.Printf("--- Información %d ---\n", i+1)
		fmt.Printf("Nombre:  0x%x (distancia: %d bytes)\n", result.NameAddress, result.Distance)
		fmt.Printf("Clase:   0x%x\n", result.ClassAddress)
		fmt.Printf("Contexto Nombre:  %s\n", result.NameContext)
		fmt.Printf("Contexto Clase:   %s\n", result.ClassContext)
		fmt.Println()
	}

	fmt.Println("=== RECOMENDACIÓN ===")
	if len(results) > 0 {
		// Buscar la mejor dirección (no cero)
		var bestAddress uintptr
		for _, result := range results {
			if result.NameAddress != 0 {
				bestAddress = result.NameAddress
				break
			}
		}

		if bestAddress != 0 {
			fmt.Printf("Usar dirección: 0x%x\n", bestAddress)
			fmt.Printf("Esta dirección contiene información del personaje\n")
			fmt.Printf("Ideal para implementar AutoBuff\n")
		} else {
			fmt.Println("No se encontraron direcciones válidas para AutoBuff")
		}
	} else {
		fmt.Println("No se encontraron resultados para AutoBuff")
	}

	fmt.Println("\nPresiona Enter para salir...")
	fmt.Scanln()
}
