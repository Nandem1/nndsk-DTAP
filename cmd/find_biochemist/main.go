package main

import (
	"fmt"
	"unsafe"

	"dtap-go/internal/process"
	"dtap-go/pkg/constants"

	"golang.org/x/sys/windows"
)

func main() {
	fmt.Println("DTAP-Go - Biochemist ID Finder")
	fmt.Println("==============================")
	fmt.Println("Buscando proceso de Ragnarok Online...")

	// Buscar proceso de RO
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
	fmt.Println("Buscando string 'Biochemist' en TODA la memoria del juego...")
	fmt.Println("(Búsqueda completa - puede tomar más tiempo)")

	// Buscar string "Biochemist" (optimizado y filtrado)
	biochemistFinder := NewBiochemistFinder(handle)
	addresses, err := biochemistFinder.FindBiochemistFiltered()
	if err != nil {
		fmt.Printf("Error buscando Biochemist: %v\n", err)
		fmt.Scanln()
		return
	}

	if len(addresses) == 0 {
		fmt.Println("No se encontró la string 'Biochemist' en memoria")
		fmt.Println("Asegúrate de que tu personaje sea Biochemist y esté en el juego")
		fmt.Scanln()
		return
	}

	fmt.Printf("\n=== RESULTADOS ===\n")
	fmt.Printf("Encontradas %d direcciones con 'Biochemist':\n", len(addresses))

	for i, addr := range addresses {
		fmt.Printf("--- Dirección %d ---\n", i+1)
		fmt.Printf("Address: 0x%x\n", addr)

		// Leer contexto alrededor de la dirección
		context, err := biochemistFinder.ReadContext(addr)
		if err == nil {
			fmt.Printf("Contexto: %s\n", context)
		}
		fmt.Println()
	}

	fmt.Println("=== INSTRUCCIONES ===")
	fmt.Println("1. Busca la dirección que contenga información relevante del personaje")
	fmt.Println("2. Anota la dirección que parezca más prometedora")
	fmt.Println("3. Esta dirección se puede usar para AutoBuff")
	fmt.Println()
	fmt.Println("Presiona Enter para salir...")
	fmt.Scanln()
}

// BiochemistFinder busca la string "Biochemist" en memoria
type BiochemistFinder struct {
	processHandle windows.Handle
}

// NewBiochemistFinder crea un nuevo finder
func NewBiochemistFinder(handle windows.Handle) *BiochemistFinder {
	return &BiochemistFinder{
		processHandle: handle,
	}
}

// FindBiochemistFiltered busca "Biochemist" filtrando memoria del propio programa
func (bf *BiochemistFinder) FindBiochemistFiltered() ([]uintptr, error) {
	var addresses []uintptr
	targetString := "Biochemist"
	targetBytes := []byte(targetString)

	// Enumerar regiones de memoria (búsqueda completa)
	var mbi windows.MemoryBasicInformation
	address := uintptr(0)
	regionsChecked := 0
	maxRegions := 2000 // Aumentar límite
	regionsSearched := 0

	for regionsChecked < maxRegions {
		err := windows.VirtualQueryEx(bf.processHandle, address, &mbi, unsafe.Sizeof(mbi))
		if err != nil {
			address += uintptr(mbi.RegionSize)
			regionsChecked++
			continue
		}

		// Buscar en todas las regiones accesibles del juego
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

			err := windows.ReadProcessMemory(bf.processHandle, uintptr(mbi.BaseAddress), &buffer[0], uintptr(len(buffer)), &bytesRead)
			if err == nil && bytesRead > 0 {
				regionsSearched++
				// Buscar string en el buffer
				foundAddresses := bf.searchInBuffer(buffer, targetBytes, uintptr(mbi.BaseAddress))

				// Filtrar resultados que parezcan del propio programa
				for _, addr := range foundAddresses {
					context, err := bf.ReadContext(addr)
					if err == nil && !bf.isSelfProgramContext(context) {
						addresses = append(addresses, addr)
						fmt.Printf("✓ Encontrada ocurrencia en región 0x%x (tamaño: %d bytes)\n", mbi.BaseAddress, mbi.RegionSize)
					}
				}
			}
		}

		address += uintptr(mbi.RegionSize)
		regionsChecked++

		// Mostrar progreso cada 50 regiones
		if regionsChecked%50 == 0 {
			fmt.Printf("Progreso: %d regiones revisadas, %d regiones buscadas...\n", regionsChecked, regionsSearched)
		}
	}

	fmt.Printf("Búsqueda completada. Revisadas %d regiones, buscadas %d regiones.\n", regionsChecked, regionsSearched)
	return addresses, nil
}

// isSelfProgramContext verifica si el contexto parece ser del propio programa
func (bf *BiochemistFinder) isSelfProgramContext(context string) bool {
	// Filtrar contextos que contengan texto del propio programa
	selfProgramTexts := []string{
		"DTAP-Go",
		"Biochemist ID Finder",
		"Buscando string",
		"Proceso abierto exitosamente",
		"Búsqueda completada",
	}

	for _, text := range selfProgramTexts {
		if len(context) > 0 && len(text) > 0 {
			// Verificar si el contexto contiene texto del programa
			for i := 0; i <= len(context)-len(text); i++ {
				if context[i:i+len(text)] == text {
					return true
				}
			}
		}
	}
	return false
}

// FindBiochemistOptimized busca "Biochemist" de forma optimizada
func (bf *BiochemistFinder) FindBiochemistOptimized() ([]uintptr, error) {
	var addresses []uintptr
	targetString := "Biochemist"
	targetBytes := []byte(targetString)

	// Enumerar regiones de memoria (optimizado pero más amplio)
	var mbi windows.MemoryBasicInformation
	address := uintptr(0)
	regionsChecked := 0
	maxRegions := 500 // Aumentar a 500 regiones
	regionsSearched := 0

	for regionsChecked < maxRegions {
		err := windows.VirtualQueryEx(bf.processHandle, address, &mbi, unsafe.Sizeof(mbi))
		if err != nil {
			address += uintptr(mbi.RegionSize)
			regionsChecked++
			continue
		}

		// Buscar en regiones accesibles (más amplio pero aún optimizado)
		if mbi.State == windows.MEM_COMMIT &&
			mbi.RegionSize < 10*1024*1024 && // Menos de 10MB (más amplio)
			(mbi.Protect&windows.PAGE_READONLY != 0 ||
				mbi.Protect&windows.PAGE_READWRITE != 0 ||
				mbi.Protect&windows.PAGE_EXECUTE_READ != 0) {

			// Leer región de memoria
			buffer := make([]byte, mbi.RegionSize)
			var bytesRead uintptr

			err := windows.ReadProcessMemory(bf.processHandle, uintptr(mbi.BaseAddress), &buffer[0], uintptr(len(buffer)), &bytesRead)
			if err == nil && bytesRead > 0 {
				regionsSearched++
				// Buscar string en el buffer
				foundAddresses := bf.searchInBuffer(buffer, targetBytes, uintptr(mbi.BaseAddress))
				addresses = append(addresses, foundAddresses...)

				// Si encontramos algo, mostrar progreso
				if len(foundAddresses) > 0 {
					fmt.Printf("✓ Encontrada ocurrencia en región 0x%x (tamaño: %d bytes)\n", mbi.BaseAddress, mbi.RegionSize)
				}
			}
		}

		address += uintptr(mbi.RegionSize)
		regionsChecked++

		// Mostrar progreso cada 50 regiones
		if regionsChecked%50 == 0 {
			fmt.Printf("Progreso: %d regiones revisadas, %d regiones buscadas...\n", regionsChecked, regionsSearched)
		}
	}

	fmt.Printf("Búsqueda completada. Revisadas %d regiones, buscadas %d regiones.\n", regionsChecked, regionsSearched)
	return addresses, nil
}

// FindBiochemist busca todas las ocurrencias de "Biochemist" (versión completa)
func (bf *BiochemistFinder) FindBiochemist() ([]uintptr, error) {
	var addresses []uintptr
	targetString := "Biochemist"
	targetBytes := []byte(targetString)

	// Enumerar regiones de memoria
	var mbi windows.MemoryBasicInformation
	address := uintptr(0)

	for {
		err := windows.VirtualQueryEx(bf.processHandle, address, &mbi, unsafe.Sizeof(mbi))
		if err != nil {
			address += uintptr(mbi.RegionSize)
			continue
		}

		// Solo buscar en regiones de memoria válidas y accesibles
		if mbi.State == windows.MEM_COMMIT &&
			(mbi.Protect&windows.PAGE_READONLY != 0 ||
				mbi.Protect&windows.PAGE_READWRITE != 0 ||
				mbi.Protect&windows.PAGE_EXECUTE_READ != 0 ||
				mbi.Protect&windows.PAGE_EXECUTE_READWRITE != 0) {

			// Leer región de memoria
			buffer := make([]byte, mbi.RegionSize)
			var bytesRead uintptr

			err := windows.ReadProcessMemory(bf.processHandle, uintptr(mbi.BaseAddress), &buffer[0], uintptr(len(buffer)), &bytesRead)
			if err == nil && bytesRead > 0 {
				// Buscar string en el buffer
				foundAddresses := bf.searchInBuffer(buffer, targetBytes, uintptr(mbi.BaseAddress))
				addresses = append(addresses, foundAddresses...)
			}
		}

		address += uintptr(mbi.RegionSize)
	}

	return addresses, nil
}

// searchInBuffer busca la string en un buffer de memoria
func (bf *BiochemistFinder) searchInBuffer(buffer []byte, target []byte, baseAddress uintptr) []uintptr {
	var addresses []uintptr

	for i := 0; i <= len(buffer)-len(target); i++ {
		match := true
		for j := 0; j < len(target); j++ {
			if buffer[i+j] != target[j] {
				match = false
				break
			}
		}

		if match {
			addresses = append(addresses, baseAddress+uintptr(i))
		}
	}

	return addresses
}

// ReadContext lee contexto alrededor de una dirección
func (bf *BiochemistFinder) ReadContext(address uintptr) (string, error) {
	// Leer 64 bytes alrededor de la dirección
	buffer := make([]byte, 64)
	var bytesRead uintptr

	// Ajustar para leer 32 bytes antes y 32 bytes después
	readAddress := address - 32
	err := windows.ReadProcessMemory(bf.processHandle, readAddress, &buffer[0], uintptr(len(buffer)), &bytesRead)
	if err != nil {
		return "", err
	}

	// Convertir a string, filtrando caracteres no imprimibles
	result := ""
	for _, b := range buffer {
		if b >= 32 && b <= 126 { // Caracteres imprimibles ASCII
			result += string(b)
		} else {
			result += "."
		}
	}

	return result, nil
}
