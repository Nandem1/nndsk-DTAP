package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
	"unsafe"

	"dtap-go/internal/process"
	"dtap-go/pkg/constants"

	"golang.org/x/sys/windows"
)

// getJobName convierte el ID del job a un nombre legible
func getJobName(jobID uint32) string {
	jobNames := map[uint32]string{
		0:    "Novice",
		1:    "Swordman",
		2:    "Magician",
		3:    "Archer",
		4:    "Acolyte",
		5:    "Merchant",
		6:    "Thief",
		7:    "Knight",
		8:    "Priest",
		9:    "Wizard",
		10:   "Blacksmith",
		11:   "Hunter",
		12:   "Assassin",
		14:   "Crusader",
		15:   "Monk",
		16:   "Sage",
		17:   "Rogue",
		18:   "Alchemist",
		19:   "Bard",
		20:   "Dancer",
		23:   "Super Novice",
		24:   "Gunslinger",
		25:   "Ninja",
		4001: "High Novice",
		4002: "High Swordman",
		4003: "High Magician",
		4004: "High Archer",
		4005: "High Acolyte",
		4006: "High Merchant",
		4007: "High Thief",
		4008: "Lord Knight",
		4009: "High Priest",
		4010: "High Wizard",
		4011: "Whitesmith",
		4012: "Sniper",
		4013: "Assassin Cross",
		4015: "Paladin",
		4016: "Champion",
		4017: "Professor",
		4018: "Stalker",
		4019: "Creator",
		4020: "Clown",
		4021: "Gypsy",
		4023: "Baby Novice",
		4024: "Baby Swordman",
		4025: "Baby Magician",
		4026: "Baby Archer",
		4027: "Baby Acolyte",
		4028: "Baby Merchant",
		4029: "Baby Thief",
		4030: "Baby Knight",
		4031: "Baby Priest",
		4032: "Baby Wizard",
		4033: "Baby Blacksmith",
		4034: "Baby Hunter",
		4035: "Baby Assassin",
		4037: "Baby Crusader",
		4038: "Baby Monk",
		4039: "Baby Sage",
		4040: "Baby Rogue",
		4041: "Baby Alchemist",
		4042: "Baby Bard",
		4043: "Baby Dancer",
		4045: "Super Baby",
		4046: "Taekwon",
		4047: "Star Gladiator",
		4049: "Soul Linker",
		4050: "Gangsi",
		4051: "Death Knight",
		4052: "Dark Collector",
		4054: "Rune Knight",
		4055: "Warlock",
		4056: "Ranger",
		4057: "Arch Bishop",
		4058: "Mechanic",
		4059: "Guillotine Cross",
		4060: "Rune Knight2",
		4061: "Warlock2",
		4062: "Ranger2",
		4063: "Arch Bishop2",
		4064: "Mechanic2",
		4065: "Guillotine2",
		4066: "Royal Guard",
		4067: "Sorcerer",
		4068: "Minstrel",
		4069: "Wanderer",
		4070: "Sura",
		4071: "Genetic",
		4072: "Shadow Chaser",
		4073: "Royal Guard2",
		4074: "Sorcerer2",
		4075: "Minstrel2",
		4076: "Wanderer2",
		4077: "Sura2",
		4078: "Genetic2",
		4079: "Shadow Chaser2",
		4096: "Baby Rune",
		4097: "Baby Warlock",
		4098: "Baby Ranger",
		4099: "Baby Bishop",
		4100: "Baby Mechanic",
		4101: "Baby Cross",
		4102: "Baby Guard",
		4103: "Baby Sorcerer",
		4104: "Baby Minstrel",
		4105: "Baby Wanderer",
		4106: "Baby Sura",
		4107: "Baby Genetic",
		4108: "Baby Chaser",
		4190: "Super Novice E",
		4191: "Super Baby E",
		4211: "Kagerou",
		4212: "Oboro",
	}

	if name, exists := jobNames[jobID]; exists {
		return name
	}
	return fmt.Sprintf("Unknown Job (%d)", jobID)
}

type JobMonitor struct {
	handle    windows.Handle
	lastJobID uint32
	addresses map[uintptr]uint32
}

func NewJobMonitor(handle windows.Handle) *JobMonitor {
	return &JobMonitor{
		handle:    handle,
		lastJobID: 0,
		addresses: make(map[uintptr]uint32),
	}
}

func (jm *JobMonitor) FindCurrentJobID() (uint32, []uintptr, error) {
	var currentJobID uint32
	var foundAddresses []uintptr

	// Buscar en el heap (rango más probable)
	heapStart := uintptr(0x10000000)
	heapEnd := uintptr(0x20000000)

	// Buscar en chunks de 32MB para ser más eficiente
	for chunkStart := heapStart; chunkStart < heapEnd; chunkStart += 0x2000000 {
		chunkEnd := chunkStart + 0x2000000
		if chunkEnd > heapEnd {
			chunkEnd = heapEnd
		}

		for addr := chunkStart; addr < chunkEnd; addr += 4 {
			value, err := readUint32(jm.handle, addr)
			if err == nil {
				// Verificar si es un job ID válido
				if isValidJobID(value) {
					if currentJobID == 0 {
						currentJobID = value
					}
					if value == currentJobID {
						foundAddresses = append(foundAddresses, addr)
					}
				}
			}
		}
	}

	return currentJobID, foundAddresses, nil
}

func (jm *JobMonitor) MonitorJobChanges() {
	fmt.Println("🔍 Monitor Dinámico de Job ID")
	fmt.Println("=============================")
	fmt.Println("Monitoreando cambios de job en tiempo real...")
	fmt.Println("Presiona Ctrl+C para salir")
	fmt.Println()

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			currentJobID, addresses, err := jm.FindCurrentJobID()
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}

			if currentJobID != 0 {
				// Verificar si el job cambió
				if jm.lastJobID != 0 && jm.lastJobID != currentJobID {
					fmt.Printf("🔄 CAMBIO DETECTADO!\n")
					fmt.Printf("   Job anterior: %d (%s)\n", jm.lastJobID, getJobName(jm.lastJobID))
					fmt.Printf("   Job actual: %d (%s)\n", currentJobID, getJobName(currentJobID))
					fmt.Printf("   Direcciones encontradas: %d\n", len(addresses))

					// Mostrar las primeras 5 direcciones
					for i, addr := range addresses {
						if i >= 5 {
							fmt.Printf("   ... y %d direcciones más\n", len(addresses)-5)
							break
						}
						fmt.Printf("   Dirección %d: 0x%X\n", i+1, addr)
					}
					fmt.Println()
				} else if jm.lastJobID == 0 {
					fmt.Printf("✅ Job inicial detectado: %d (%s)\n", currentJobID, getJobName(currentJobID))
					fmt.Printf("   Direcciones encontradas: %d\n", len(addresses))
					fmt.Println()
				}

				jm.lastJobID = currentJobID
				jm.addresses = make(map[uintptr]uint32)
				for _, addr := range addresses {
					jm.addresses[addr] = currentJobID
				}
			} else {
				fmt.Printf("⏳ Esperando detección de job...\n")
			}
		}
	}
}

func isValidJobID(value uint32) bool {
	// Lista de job IDs válidos conocidos
	validJobIDs := map[uint32]bool{
		0: true, 1: true, 2: true, 3: true, 4: true, 5: true, 6: true, 7: true, 8: true, 9: true,
		10: true, 11: true, 12: true, 14: true, 15: true, 16: true, 17: true, 18: true, 19: true, 20: true,
		23: true, 24: true, 25: true,
		4001: true, 4002: true, 4003: true, 4004: true, 4005: true, 4006: true, 4007: true,
		4008: true, 4009: true, 4010: true, 4011: true, 4012: true, 4013: true, 4015: true, 4016: true,
		4017: true, 4018: true, 4019: true, 4020: true, 4021: true, 4023: true, 4024: true, 4025: true,
		4026: true, 4027: true, 4028: true, 4029: true, 4030: true, 4031: true, 4032: true, 4033: true,
		4034: true, 4035: true, 4037: true, 4038: true, 4039: true, 4040: true, 4041: true, 4042: true,
		4043: true, 4045: true, 4046: true, 4047: true, 4049: true, 4050: true, 4051: true, 4052: true,
		4054: true, 4055: true, 4056: true, 4057: true, 4058: true, 4059: true, 4060: true, 4061: true,
		4062: true, 4063: true, 4064: true, 4065: true, 4066: true, 4067: true, 4068: true, 4069: true,
		4070: true, 4071: true, 4072: true, 4073: true, 4074: true, 4075: true, 4076: true, 4077: true,
		4078: true, 4079: true, 4096: true, 4097: true, 4098: true, 4099: true, 4100: true, 4101: true,
		4102: true, 4103: true, 4104: true, 4105: true, 4106: true, 4107: true, 4108: true,
		4190: true, 4191: true, 4211: true, 4212: true,
	}

	return validJobIDs[value]
}

func main() {
	fmt.Println("DTAP-Go - Dynamic Job Monitor")
	fmt.Println("=============================")

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

	// Crear monitor
	monitor := NewJobMonitor(handle)

	// Verificar si se quiere hacer una búsqueda específica
	if len(os.Args) > 1 {
		jobID, err := strconv.ParseUint(os.Args[1], 10, 32)
		if err == nil {
			fmt.Printf("Buscando job ID específico: %d (%s)\n", jobID, getJobName(uint32(jobID)))
			fmt.Println("=====================================")

			// Buscar el job ID específico
			heapStart := uintptr(0x10000000)
			heapEnd := uintptr(0x20000000)
			foundCount := 0

			for chunkStart := heapStart; chunkStart < heapEnd; chunkStart += 0x2000000 {
				chunkEnd := chunkStart + 0x2000000
				if chunkEnd > heapEnd {
					chunkEnd = heapEnd
				}

				for addr := chunkStart; addr < chunkEnd; addr += 4 {
					value, err := readUint32(handle, addr)
					if err == nil && value == uint32(jobID) {
						fmt.Printf("✓ ENCONTRADO: 0x%X (valor: %d)\n", addr, value)
						foundCount++
					}
				}
			}

			if foundCount == 0 {
				fmt.Println("❌ No se encontró el job ID especificado")
			} else {
				fmt.Printf("✅ Se encontraron %d direcciones con el job ID %d\n", foundCount, jobID)
			}

			fmt.Println("\nPresiona Enter para salir...")
			fmt.Scanln()
			return
		}
	}

	// Iniciar monitoreo dinámico
	monitor.MonitorJobChanges()
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
