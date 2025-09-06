package finder

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

// AddressFinder busca direcciones de memoria por valores específicos
type AddressFinder struct {
	processHandle windows.Handle
}

// NewAddressFinder crea un nuevo finder
func NewAddressFinder(processHandle windows.Handle) *AddressFinder {
	return &AddressFinder{
		processHandle: processHandle,
	}
}

// FindValue busca un valor específico en la memoria del proceso
func (af *AddressFinder) FindValue(value uint32) ([]uintptr, error) {
	var addresses []uintptr

	// Enumerar regiones de memoria
	regions, err := af.enumerateMemoryRegions()
	if err != nil {
		return nil, fmt.Errorf("error enumerando regiones: %v", err)
	}

	fmt.Printf("Buscando valor %d (0x%X) en %d regiones de memoria...\n", value, value, len(regions))

	for i, region := range regions {
		if i%10 == 0 {
			fmt.Printf("Escaneando región %d/%d (0x%X - 0x%X)\n",
				i+1, len(regions), region.BaseAddress, region.BaseAddress+region.Size)
		}

		// Leer la región de memoria
		buffer := make([]byte, region.Size)
		var bytesRead uintptr

		err := windows.ReadProcessMemory(
			af.processHandle,
			region.BaseAddress,
			&buffer[0],
			region.Size,
			&bytesRead,
		)

		if err != nil {
			continue // Saltar regiones que no se pueden leer
		}

		// Buscar el valor en el buffer
		for offset := 0; offset <= len(buffer)-4; offset += 4 {
			// Leer 4 bytes como uint32
			val := *(*uint32)(unsafe.Pointer(&buffer[offset]))
			if val == value {
				addr := region.BaseAddress + uintptr(offset)
				addresses = append(addresses, addr)

				// Limitar resultados para no saturar
				if len(addresses) >= 100 {
					fmt.Printf("Encontradas 100+ direcciones, limitando resultados...\n")
					return addresses, nil
				}
			}
		}
	}

	return addresses, nil
}

// FindHPAndSP busca direcciones de HP y SP que estén cerca
func (af *AddressFinder) FindHPAndSP(hpCurrent, hpMax, spCurrent, spMax uint32) ([]HPSPAddress, error) {
	fmt.Printf("Buscando HP: %d/%d y SP: %d/%d\n", hpCurrent, hpMax, spCurrent, spMax)

	// Buscar direcciones de HP actual
	hpCurrentAddrs, err := af.FindValue(hpCurrent)
	if err != nil {
		return nil, fmt.Errorf("error buscando HP actual: %v", err)
	}
	fmt.Printf("Encontradas %d direcciones para HP actual (%d)\n", len(hpCurrentAddrs), hpCurrent)

	// Buscar direcciones de HP máximo
	hpMaxAddrs, err := af.FindValue(hpMax)
	if err != nil {
		return nil, fmt.Errorf("error buscando HP máximo: %v", err)
	}
	fmt.Printf("Encontradas %d direcciones para HP máximo (%d)\n", len(hpMaxAddrs), hpMax)

	// Buscar direcciones de SP actual
	spCurrentAddrs, err := af.FindValue(spCurrent)
	if err != nil {
		return nil, fmt.Errorf("error buscando SP actual: %v", err)
	}
	fmt.Printf("Encontradas %d direcciones para SP actual (%d)\n", len(spCurrentAddrs), spCurrent)

	// Buscar direcciones de SP máximo
	spMaxAddrs, err := af.FindValue(spMax)
	if err != nil {
		return nil, fmt.Errorf("error buscando SP máximo: %v", err)
	}
	fmt.Printf("Encontradas %d direcciones para SP máximo (%d)\n", len(spMaxAddrs), spMax)

	// Buscar combinaciones que estén cerca (HP actual + 4 = HP máximo, etc.)
	var results []HPSPAddress

	for _, hpCurr := range hpCurrentAddrs {
		for _, hpMax := range hpMaxAddrs {
			// Verificar si HP máximo está 4 bytes después de HP actual
			if hpMax == hpCurr+4 {
				// Buscar SP cerca de esta dirección
				for _, spCurr := range spCurrentAddrs {
					for _, spMax := range spMaxAddrs {
						// Verificar si SP máximo está 4 bytes después de SP actual
						if spMax == spCurr+4 {
							// Verificar si SP está cerca de HP (típicamente 8-16 bytes de diferencia)
							diff := uintptr(0)
							if spCurr > hpCurr {
								diff = spCurr - hpCurr
							} else {
								diff = hpCurr - spCurr
							}

							if diff >= 8 && diff <= 32 { // Rango típico
								results = append(results, HPSPAddress{
									HPCurrent: hpCurr,
									HPMax:     hpMax,
									SPCurrent: spCurr,
									SPMax:     spMax,
									Distance:  diff,
								})
							}
						}
					}
				}
			}
		}
	}

	return results, nil
}

// HPSPAddress representa un conjunto de direcciones de HP/SP
type HPSPAddress struct {
	HPCurrent uintptr
	HPMax     uintptr
	SPCurrent uintptr
	SPMax     uintptr
	Distance  uintptr
}

// enumerateMemoryRegions enumera las regiones de memoria del proceso
func (af *AddressFinder) enumerateMemoryRegions() ([]MemoryRegion, error) {
	var regions []MemoryRegion
	var address uintptr

	for {
		var mbi windows.MemoryBasicInformation
		size := unsafe.Sizeof(mbi)

		err := windows.VirtualQueryEx(
			af.processHandle,
			address,
			&mbi,
			size,
		)

		if err != nil {
			break
		}

		// Solo regiones con MEM_COMMIT y que sean legibles
		if mbi.State == windows.MEM_COMMIT &&
			(mbi.Protect&windows.PAGE_READONLY != 0 ||
				mbi.Protect&windows.PAGE_READWRITE != 0 ||
				mbi.Protect&windows.PAGE_EXECUTE_READ != 0 ||
				mbi.Protect&windows.PAGE_EXECUTE_READWRITE != 0) {

			regions = append(regions, MemoryRegion{
				BaseAddress: uintptr(mbi.BaseAddress),
				Size:        uintptr(mbi.RegionSize),
				Protect:     mbi.Protect,
			})
		}

		address = uintptr(mbi.BaseAddress) + uintptr(mbi.RegionSize)
	}

	return regions, nil
}

// MemoryRegion representa una región de memoria
type MemoryRegion struct {
	BaseAddress uintptr
	Size        uintptr
	Protect     uint32
}
