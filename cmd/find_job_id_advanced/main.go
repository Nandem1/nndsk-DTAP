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

func main() {
	fmt.Println("DTAP-Go - Advanced Job ID Finder")
	fmt.Println("================================")

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
		fmt.Print("Ingresa el job ID a buscar (ej: 4016 para Champion): ")
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

	found := false
	foundCount := 0

	// 1. Buscar en la dirección sospechosa del inspect_character (0x10df6b4)
	fmt.Println("1. Buscando en dirección sospechosa del inspect_character...")
	suspiciousAddr := uintptr(0x10df6b4)
	value, err := readUint32(handle, suspiciousAddr)
	if err == nil {
		fmt.Printf("   Dirección 0x%X: %d (%s)\n", suspiciousAddr, value, getJobName(value))
		if value == jobID {
			fmt.Printf("   ✓ ENCONTRADO en 0x%X!\n", suspiciousAddr)
			found = true
			foundCount++
		}
	}

	// 2. Buscar alrededor de la dirección del personaje (más amplio)
	fmt.Println("2. Buscando alrededor de la dirección del personaje (rango amplio)...")
	baseAddr := uintptr(0x10df5d8)
	for i := -200; i <= 400; i++ {
		addr := baseAddr + uintptr(i*4)
		value, err := readUint32(handle, addr)
		if err == nil && value == jobID {
			fmt.Printf("   ✓ ENCONTRADO en offset %d: 0x%X (valor: %d)\n", i*4, addr, value)
			found = true
			foundCount++
		}
	}

	// 3. Buscar en diferentes formatos (byte, word, dword)
	fmt.Println("3. Buscando en diferentes formatos...")

	// Buscar como byte (8 bits)
	fmt.Println("   Buscando como BYTE (8 bits)...")
	for i := -100; i <= 200; i++ {
		addr := baseAddr + uintptr(i)
		value, err := readUint8(handle, addr)
		if err == nil && uint32(value) == jobID {
			fmt.Printf("   ✓ ENCONTRADO como BYTE en offset %d: 0x%X (valor: %d)\n", i, addr, value)
			found = true
			foundCount++
		}
	}

	// Buscar como word (16 bits)
	fmt.Println("   Buscando como WORD (16 bits)...")
	for i := -100; i <= 200; i++ {
		addr := baseAddr + uintptr(i*2)
		value, err := readUint16(handle, addr)
		if err == nil && uint32(value) == jobID {
			fmt.Printf("   ✓ ENCONTRADO como WORD en offset %d: 0x%X (valor: %d)\n", i*2, addr, value)
			found = true
			foundCount++
		}
	}

	// 4. Buscar en rangos específicos basados en el inspect_character
	fmt.Println("4. Buscando en rangos específicos del inspect_character...")

	// Rango donde encontró "Posible stat 10" (offset +220)
	statRange := baseAddr + 220
	for i := -20; i <= 20; i++ {
		addr := statRange + uintptr(i*4)
		value, err := readUint32(handle, addr)
		if err == nil && value == jobID {
			fmt.Printf("   ✓ ENCONTRADO en rango de stats offset %d: 0x%X (valor: %d)\n", 220+i*4, addr, value)
			found = true
			foundCount++
		}
	}

	// 5. Buscar en el heap con rango más amplio
	fmt.Println("5. Buscando en heap (rango amplio)...")
	heapStart := uintptr(0x10000000)
	heapEnd := uintptr(0x20000000)

	// Buscar en chunks de 16MB para ser más eficiente
	for chunkStart := heapStart; chunkStart < heapEnd; chunkStart += 0x1000000 {
		chunkEnd := chunkStart + 0x1000000
		if chunkEnd > heapEnd {
			chunkEnd = heapEnd
		}

		for addr := chunkStart; addr < chunkEnd; addr += 4 {
			value, err := readUint32(handle, addr)
			if err == nil && value == jobID {
				fmt.Printf("   ✓ ENCONTRADO en heap: 0x%X (valor: %d)\n", addr, value)
				found = true
				foundCount++
			}
		}
	}

	if !found {
		fmt.Println("❌ No se encontró el job ID en los rangos buscados")
		fmt.Println("Posibles causas:")
		fmt.Println("1. El job ID no está en memoria como entero de 32 bits")
		fmt.Println("2. Está en un formato diferente (byte, word, o con offset)")
		fmt.Println("3. Está en una región de memoria no accesible")
		fmt.Println("4. El job ID actual no es el esperado")
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

// readUint16 lee un entero de 16 bits desde memoria
func readUint16(handle windows.Handle, address uintptr) (uint16, error) {
	var buffer [2]byte
	var bytesRead uintptr

	err := windows.ReadProcessMemory(
		handle,
		address,
		&buffer[0],
		2,
		&bytesRead,
	)
	if err != nil {
		return 0, err
	}

	if bytesRead != 2 {
		return 0, fmt.Errorf("solo se leyeron %d bytes de 2", bytesRead)
	}

	return *(*uint16)(unsafe.Pointer(&buffer[0])), nil
}

// readUint8 lee un entero de 8 bits desde memoria
func readUint8(handle windows.Handle, address uintptr) (uint8, error) {
	var buffer [1]byte
	var bytesRead uintptr

	err := windows.ReadProcessMemory(
		handle,
		address,
		&buffer[0],
		1,
		&bytesRead,
	)
	if err != nil {
		return 0, err
	}

	if bytesRead != 1 {
		return 0, fmt.Errorf("solo se leyó %d byte de 1", bytesRead)
	}

	return buffer[0], nil
}
