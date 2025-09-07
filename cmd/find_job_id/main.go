package main

import (
	"fmt"
	"os"
	"strconv"
	"unsafe"

	"dtap-go/internal/process"
	"dtap-go/pkg/constants"

	"golang.org/x/sys/windows"
)

func main() {
	fmt.Println("DTAP-Go - Job ID Finder")
	fmt.Println("=======================")

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

	// Obtener job ID a buscar
	var jobID uint32
	if len(os.Args) > 1 {
		// Job ID desde argumento
		id, err := strconv.ParseUint(os.Args[1], 10, 32)
		if err != nil {
			fmt.Printf("Error parseando job ID: %v\n", err)
			fmt.Println("Presiona Enter para salir...")
			fmt.Scanln()
			return
		}
		jobID = uint32(id)
	} else {
		// Job ID desde input
		fmt.Print("Ingresa el job ID a buscar (ej: 4019 para Creator): ")
		var input string
		fmt.Scanln(&input)

		id, err := strconv.ParseUint(input, 10, 32)
		if err != nil {
			fmt.Printf("Error parseando job ID: %v\n", err)
			fmt.Println("Presiona Enter para salir...")
			fmt.Scanln()
			return
		}
		jobID = uint32(id)
	}

	fmt.Printf("Buscando job ID: %d (%s)\n", jobID, getJobName(jobID))
	fmt.Println("=====================================")

	// Buscar en diferentes rangos - BÚSQUEDA COMPLETA SIN LÍMITES
	found := false
	foundCount := 0

	// Rango 1: Alrededor de la dirección del personaje (rango ampliado)
	fmt.Println("Buscando alrededor de la dirección del personaje (rango ampliado)...")
	baseAddr := uintptr(0x10df5d8)
	for i := -1000; i <= 2000; i++ {
		addr := baseAddr + uintptr(i*4)
		value, err := readUint32(handle, addr)
		if err == nil && value == jobID {
			fmt.Printf("✓ [PERSONAJE] Offset %d: 0x%X -> Valor: %d (0x%X) - %s\n",
				i*4, addr, value, value, getJobName(value))
			found = true
			foundCount++
		}
	}

	// Rango 2: Alrededor de las direcciones de HP/SP (rango ampliado)
	fmt.Println("Buscando alrededor de las direcciones de HP/SP (rango ampliado)...")
	hpAddr := uintptr(0x10DCE10)
	for i := -1000; i <= 1000; i++ {
		addr := hpAddr + uintptr(i*4)
		value, err := readUint32(handle, addr)
		if err == nil && value == jobID {
			fmt.Printf("✓ [HP/SP] Offset %d: 0x%X -> Valor: %d (0x%X) - %s\n",
				i*4, addr, value, value, getJobName(value))
			found = true
			foundCount++
		}
	}

	// Rango 3: Búsqueda en heap completo (sin límites)
	fmt.Println("Buscando en heap completo (sin límites)...")
	fmt.Println("ADVERTENCIA: Esta búsqueda puede tomar varios minutos...")

	// Buscar en rangos de 64MB cada vez para no sobrecargar
	heapStart := uintptr(0x10000000)
	heapEnd := uintptr(0x20000000) // 256MB total

	for chunkStart := heapStart; chunkStart < heapEnd; chunkStart += 0x4000000 { // 64MB chunks
		chunkEnd := chunkStart + 0x4000000
		if chunkEnd > heapEnd {
			chunkEnd = heapEnd
		}

		fmt.Printf("Escaneando chunk: 0x%X - 0x%X\n", chunkStart, chunkEnd)

		for addr := chunkStart; addr < chunkEnd; addr += 4 {
			value, err := readUint32(handle, addr)
			if err == nil && value == jobID {
				fmt.Printf("✓ [HEAP] 0x%X -> Valor: %d (0x%X) - %s\n",
					addr, value, value, getJobName(value))
				found = true
				foundCount++
			}
		}
	}

	// Rango 4: Búsqueda en stack y otras áreas
	fmt.Println("Buscando en stack y otras áreas de memoria...")
	stackStart := uintptr(0x7FF00000)
	stackEnd := uintptr(0x7FFFFFFF)

	for addr := stackStart; addr < stackEnd; addr += 4 {
		value, err := readUint32(handle, addr)
		if err == nil && value == jobID {
			fmt.Printf("✓ [STACK] 0x%X -> Valor: %d (0x%X) - %s\n",
				addr, value, value, getJobName(value))
			found = true
			foundCount++
		}
	}

	if !found {
		fmt.Println("❌ No se encontró el job ID en los rangos buscados")
		fmt.Println("Intenta con un rango más amplio o verifica que el job ID sea correcto")
	} else {
		fmt.Printf("✅ Búsqueda completada. Se encontraron %d direcciones con el job ID %d (0x%X)\n",
			foundCount, jobID, jobID)
	}

	fmt.Println("\nPresiona Enter para salir...")
	fmt.Scanln()
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
