package scanner

import (
	"fmt"
	"strings"
	"unsafe"

	"golang.org/x/sys/windows"
)

// Pattern representa un patrón de bytes a buscar
type Pattern struct {
	Bytes []byte
	Mask  []bool // true = byte exacto, false = wildcard
}

// ParsePattern convierte un string de patrón a Pattern
// Ejemplo: "8B 0D ?? ?? ?? ?? 85 C9"
func ParsePattern(patternStr string) Pattern {
	parts := strings.Fields(patternStr)
	bytes := make([]byte, 0, len(parts))
	mask := make([]bool, 0, len(parts))

	for _, part := range parts {
		if part == "??" {
			bytes = append(bytes, 0)
			mask = append(mask, false) // wildcard
		} else {
			var b byte
			fmt.Sscanf(part, "%02X", &b)
			bytes = append(bytes, b)
			mask = append(mask, true) // byte exacto
		}
	}

	return Pattern{Bytes: bytes, Mask: mask}
}

// MemoryRegion representa una región de memoria
type MemoryRegion struct {
	BaseAddress uintptr
	Size        uintptr
	Protect     uint32
}

// EnumerateMemoryRegions enumera las regiones de memoria del proceso
func EnumerateMemoryRegions(processHandle windows.Handle) ([]MemoryRegion, error) {
	var regions []MemoryRegion
	var address uintptr

	for {
		var mbi windows.MemoryBasicInformation
		size := unsafe.Sizeof(mbi)

		err := windows.VirtualQueryEx(
			processHandle,
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

// ScanPattern busca un patrón en la memoria del proceso
func ScanPattern(processHandle windows.Handle, pattern Pattern) (uintptr, error) {
	regions, err := EnumerateMemoryRegions(processHandle)
	if err != nil {
		return 0, fmt.Errorf("error enumerando regiones: %v", err)
	}

	for _, region := range regions {
		// Leer la región de memoria
		buffer := make([]byte, region.Size)
		var bytesRead uintptr

		err := windows.ReadProcessMemory(
			processHandle,
			region.BaseAddress,
			&buffer[0],
			region.Size,
			&bytesRead,
		)

		if err != nil {
			continue // Saltar regiones que no se pueden leer
		}

		// Buscar el patrón en el buffer
		for i := 0; i <= len(buffer)-len(pattern.Bytes); i++ {
			if matchPattern(buffer[i:], pattern) {
				return region.BaseAddress + uintptr(i), nil
			}
		}
	}

	return 0, fmt.Errorf("patrón no encontrado")
}

// matchPattern verifica si un buffer coincide con un patrón
func matchPattern(buffer []byte, pattern Pattern) bool {
	if len(buffer) < len(pattern.Bytes) {
		return false
	}

	for i, patternByte := range pattern.Bytes {
		if pattern.Mask[i] && buffer[i] != patternByte {
			return false
		}
	}

	return true
}

// ScanMultiplePatterns busca múltiples patrones y retorna el primero encontrado
func ScanMultiplePatterns(processHandle windows.Handle, patterns []string) (uintptr, error) {
	for _, patternStr := range patterns {
		pattern := ParsePattern(patternStr)
		addr, err := ScanPattern(processHandle, pattern)
		if err == nil {
			return addr, nil
		}
	}

	return 0, fmt.Errorf("ningún patrón encontrado")
}
