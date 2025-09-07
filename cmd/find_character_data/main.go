package main

import (
	"fmt"
	"unsafe"

	"dtap-go/internal/process"
	"dtap-go/pkg/constants"

	"golang.org/x/sys/windows"
)

type CharacterData struct {
	Name      string
	Job       string
	JobID     uint32
	BaseLevel uint32
	JobLevel  uint32
	HP        uint32
	MaxHP     uint32
	SP        uint32
	MaxSP     uint32
	Weight    uint32
	MaxWeight uint32
	Zeny      uint32
}

func main() {
	fmt.Println("DTAP-Go - Character Data Finder")
	fmt.Println("===============================")

	// Buscar proceso de RO
	fmt.Println("Buscando procesos de Ragnarok Online...")
	roProcess, err := process.FindROProcess()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		fmt.Println("Asegúrate de que Ragnarok Online esté ejecutándose.")
		fmt.Println("Presiona Enter para salir...")
		fmt.Scanln()
		return
	}

	fmt.Printf("✓ Encontrado: %s (PID: %d)\n", roProcess.Name, roProcess.PID)

	// Abrir proceso
	handle, err := windows.OpenProcess(
		constants.ProcessAccessFlags,
		false,
		roProcess.PID,
	)
	if err != nil {
		fmt.Printf("Error abriendo proceso: %v\n", err)
		fmt.Println("Presiona Enter para salir...")
		fmt.Scanln()
		return
	}
	defer windows.CloseHandle(handle)

	fmt.Println("Proceso abierto exitosamente")
	fmt.Println()

	// Datos de los personajes
	characters := []CharacterData{
		{
			Name:      "Sokoninaru",
			Job:       "Champion",
			JobID:     4016,
			BaseLevel: 99,
			JobLevel:  70,
			HP:        11558,
			MaxHP:     11558,
			SP:        736,
			MaxSP:     736,
			Weight:    0,
			MaxWeight: 4130,
			Zeny:      9223,
		},
		{
			Name:      "draw",
			Job:       "Biochemist",
			JobID:     18,
			BaseLevel: 99,
			JobLevel:  70,
			HP:        11112,
			MaxHP:     11112,
			SP:        964,
			MaxSP:     964,
			Weight:    1467,
			MaxWeight: 5030,
			Zeny:      9223,
		},
	}

	// Buscar cada valor en memoria
	for i, char := range characters {
		fmt.Printf("=== PERSONAJE %d: %s (%s) ===\n", i+1, char.Name, char.Job)
		fmt.Printf("Job ID: %d, Base Level: %d, Job Level: %d\n", char.JobID, char.BaseLevel, char.JobLevel)
		fmt.Printf("HP: %d/%d, SP: %d/%d\n", char.HP, char.MaxHP, char.SP, char.MaxSP)
		fmt.Printf("Weight: %d/%d, Zeny: %d\n", char.Weight, char.MaxWeight, char.Zeny)
		fmt.Println()

		// Buscar cada valor
		searchValues := map[string]uint32{
			"Job ID":     char.JobID,
			"Base Level": char.BaseLevel,
			"Job Level":  char.JobLevel,
			"HP":         char.HP,
			"Max HP":     char.MaxHP,
			"SP":         char.SP,
			"Max SP":     char.MaxSP,
			"Weight":     char.Weight,
			"Max Weight": char.MaxWeight,
			"Zeny":       char.Zeny,
		}

		for valueName, value := range searchValues {
			if value == 0 {
				continue // Saltar valores 0
			}
			
			fmt.Printf("Buscando %s: %d\n", valueName, value)
			addresses := findValueInMemory(handle, value)
			
			if len(addresses) > 0 {
				fmt.Printf("  ✓ Encontrado en %d direcciones:\n", len(addresses))
				// Mostrar solo las primeras 10 direcciones
				for j, addr := range addresses {
					if j >= 10 {
						fmt.Printf("  ... y %d direcciones más\n", len(addresses)-10)
						break
					}
					fmt.Printf("    0x%X\n", addr)
				}
			} else {
				fmt.Printf("  ❌ No encontrado\n")
			}
			fmt.Println()
		}

		fmt.Println("----------------------------------------")
		fmt.Println()
	}

	// Buscar patrones de estructuras
	fmt.Println("=== ANÁLISIS DE PATRONES ===")
	analyzePatterns(handle, characters)

	fmt.Println("Presiona Enter para salir...")
	fmt.Scanln()
}

func findValueInMemory(handle windows.Handle, value uint32) []uintptr {
	var addresses []uintptr

	// Buscar en el heap
	heapStart := uintptr(0x10000000)
	heapEnd := uintptr(0x20000000)
	
	// Buscar en chunks de 32MB para ser más eficiente
	for chunkStart := heapStart; chunkStart < heapEnd; chunkStart += 0x2000000 {
		chunkEnd := chunkStart + 0x2000000
		if chunkEnd > heapEnd {
			chunkEnd = heapEnd
		}
		
		for addr := chunkStart; addr < chunkEnd; addr += 4 {
			readValue, err := readUint32(handle, addr)
			if err == nil && readValue == value {
				addresses = append(addresses, addr)
			}
		}
	}

	return addresses
}

func analyzePatterns(handle windows.Handle, characters []CharacterData) {
	fmt.Println("Buscando patrones de estructuras de datos...")
	fmt.Println()

	// Buscar direcciones que contengan múltiples valores del mismo personaje
	for _, char := range characters {
		fmt.Printf("Analizando patrones para %s:\n", char.Name)
		
		// Buscar direcciones que contengan HP y MaxHP cerca
		hpAddresses := findValueInMemory(handle, char.HP)
		maxHPAddresses := findValueInMemory(handle, char.MaxHP)
		
		fmt.Printf("  HP (%d): %d direcciones\n", char.HP, len(hpAddresses))
		fmt.Printf("  MaxHP (%d): %d direcciones\n", char.MaxHP, len(maxHPAddresses))
		
		// Buscar direcciones cercanas
		nearbyPairs := findNearbyAddresses(hpAddresses, maxHPAddresses, 64) // Dentro de 64 bytes
		if len(nearbyPairs) > 0 {
			fmt.Printf("  ✓ Encontradas %d parejas HP/MaxHP cercanas:\n", len(nearbyPairs))
			for j, pair := range nearbyPairs {
				if j >= 5 {
					fmt.Printf("  ... y %d parejas más\n", len(nearbyPairs)-5)
					break
				}
				fmt.Printf("    HP: 0x%X, MaxHP: 0x%X (diferencia: %d bytes)\n", 
					pair.addr1, pair.addr2, pair.distance)
			}
		}
		
		// Buscar direcciones que contengan SP y MaxSP cerca
		spAddresses := findValueInMemory(handle, char.SP)
		maxSPAddresses := findValueInMemory(handle, char.MaxSP)
		
		fmt.Printf("  SP (%d): %d direcciones\n", char.SP, len(spAddresses))
		fmt.Printf("  MaxSP (%d): %d direcciones\n", char.MaxSP, len(maxSPAddresses))
		
		nearbySPPairs := findNearbyAddresses(spAddresses, maxSPAddresses, 64)
		if len(nearbySPPairs) > 0 {
			fmt.Printf("  ✓ Encontradas %d parejas SP/MaxSP cercanas:\n", len(nearbySPPairs))
			for j, pair := range nearbySPPairs {
				if j >= 5 {
					fmt.Printf("  ... y %d parejas más\n", len(nearbySPPairs)-5)
					break
				}
				fmt.Printf("    SP: 0x%X, MaxSP: 0x%X (diferencia: %d bytes)\n", 
					pair.addr1, pair.addr2, pair.distance)
			}
		}
		
		// Buscar direcciones que contengan Job ID cerca de otros stats
		jobIDAddresses := findValueInMemory(handle, char.JobID)
		fmt.Printf("  Job ID (%d): %d direcciones\n", char.JobID, len(jobIDAddresses))
		
		// Buscar Job ID cerca de HP
		nearbyJobHP := findNearbyAddresses(jobIDAddresses, hpAddresses, 128)
		if len(nearbyJobHP) > 0 {
			fmt.Printf("  ✓ Encontradas %d parejas Job ID/HP cercanas:\n", len(nearbyJobHP))
			for j, pair := range nearbyJobHP {
				if j >= 3 {
					fmt.Printf("  ... y %d parejas más\n", len(nearbyJobHP)-3)
					break
				}
				fmt.Printf("    Job ID: 0x%X, HP: 0x%X (diferencia: %d bytes)\n", 
					pair.addr1, pair.addr2, pair.distance)
			}
		}
		
		fmt.Println()
	}
}

type AddressPair struct {
	addr1    uintptr
	addr2    uintptr
	distance int
}

func findNearbyAddresses(addresses1, addresses2 []uintptr, maxDistance int) []AddressPair {
	var pairs []AddressPair
	
	for _, addr1 := range addresses1 {
		for _, addr2 := range addresses2 {
			distance := int(addr2) - int(addr1)
			if distance < 0 {
				distance = -distance
			}
			if distance <= maxDistance {
				pairs = append(pairs, AddressPair{
					addr1:    addr1,
					addr2:    addr2,
					distance: distance,
				})
			}
		}
	}
	
	return pairs
}

// readUint32 lee un entero de 32 bits desde memoria
func readUint32(handle windows.Handle, address uintptr) (uint32, error) {
	var buffer [4]byte
	var bytesRead uintptr

	err := windows.ReadProcessMemory(
		handle,
		address,
		&buffer[0],
		4,
		&bytesRead,
	)
	if err != nil {
		return 0, err
	}

	if bytesRead != 4 {
		return 0, fmt.Errorf("solo se leyeron %d bytes de 4", bytesRead)
	}

	return *(*uint32)(unsafe.Pointer(&buffer[0])), nil
}
