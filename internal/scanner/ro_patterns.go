package scanner

import (
	"golang.org/x/sys/windows"
)

// Patrones comunes para Ragnarok Online
// Estos patrones buscan instrucciones que acceden a las direcciones de HP/SP

// HPPatterns - Patrones para encontrar la dirección base de HP
var HPPatterns = []string{
	// Patrones comunes de acceso a HP
	"8B 0D ?? ?? ?? ?? 85 C9 74 ??", // MOV ECX, [HP_ADDRESS]; TEST ECX, ECX; JZ
	"A1 ?? ?? ?? ?? 85 C0 74 ??",    // MOV EAX, [HP_ADDRESS]; TEST EAX, EAX; JZ
	"8B 15 ?? ?? ?? ?? 85 D2 74 ??", // MOV EDX, [HP_ADDRESS]; TEST EDX, EDX; JZ
	"8B 35 ?? ?? ?? ?? 85 F6 74 ??", // MOV ESI, [HP_ADDRESS]; TEST ESI, ESI; JZ
	"8B 3D ?? ?? ?? ?? 85 FF 74 ??", // MOV EDI, [HP_ADDRESS]; TEST EDI, EDI; JZ

	// Patrones de comparación de HP
	"3B 0D ?? ?? ?? ?? 7C ??", // CMP ECX, [HP_ADDRESS]; JL
	"3B 05 ?? ?? ?? ?? 7C ??", // CMP EAX, [HP_ADDRESS]; JL
	"3B 15 ?? ?? ?? ?? 7C ??", // CMP EDX, [HP_ADDRESS]; JL

	// Patrones de asignación de HP
	"89 0D ?? ?? ?? ??", // MOV [HP_ADDRESS], ECX
	"89 05 ?? ?? ?? ??", // MOV [HP_ADDRESS], EAX
	"89 15 ?? ?? ?? ??", // MOV [HP_ADDRESS], EDX
}

// SPPatterns - Patrones para encontrar la dirección base de SP
var SPPatterns = []string{
	// Patrones comunes de acceso a SP
	"8B 0D ?? ?? ?? ?? 85 C9 75 ??", // MOV ECX, [SP_ADDRESS]; TEST ECX, ECX; JNZ
	"A1 ?? ?? ?? ?? 85 C0 75 ??",    // MOV EAX, [SP_ADDRESS]; TEST EAX, EAX; JNZ
	"8B 15 ?? ?? ?? ?? 85 D2 75 ??", // MOV EDX, [SP_ADDRESS]; TEST EDX, EDX; JNZ

	// Patrones de comparación de SP
	"3B 0D ?? ?? ?? ?? 7D ??", // CMP ECX, [SP_ADDRESS]; JGE
	"3B 05 ?? ?? ?? ?? 7D ??", // CMP EAX, [SP_ADDRESS]; JGE
	"3B 15 ?? ?? ?? ?? 7D ??", // CMP EDX, [SP_ADDRESS]; JGE

	// Patrones de asignación de SP
	"89 0D ?? ?? ?? ??", // MOV [SP_ADDRESS], ECX
	"89 05 ?? ?? ?? ??", // MOV [SP_ADDRESS], EAX
	"89 15 ?? ?? ?? ??", // MOV [SP_ADDRESS], EDX
}

// StatusPatterns - Patrones para encontrar el buffer de estados
var StatusPatterns = []string{
	// Patrones de acceso al buffer de estados
	"8B 0D ?? ?? ?? ?? 85 C9 0F 84", // MOV ECX, [STATUS_ADDRESS]; TEST ECX, ECX; JZ
	"A1 ?? ?? ?? ?? 85 C0 0F 84",    // MOV EAX, [STATUS_ADDRESS]; TEST EAX, EAX; JZ
	"8B 15 ?? ?? ?? ?? 85 D2 0F 84", // MOV EDX, [STATUS_ADDRESS]; TEST EDX, EDX; JZ
}

// FindHPBaseAddress busca la dirección base de HP usando pattern scanning
func FindHPBaseAddress(processHandle windows.Handle) (uintptr, error) {
	return ScanMultiplePatterns(processHandle, HPPatterns)
}

// FindSPBaseAddress busca la dirección base de SP usando pattern scanning
func FindSPBaseAddress(processHandle windows.Handle) (uintptr, error) {
	return ScanMultiplePatterns(processHandle, SPPatterns)
}

// FindStatusBufferAddress busca la dirección del buffer de estados
func FindStatusBufferAddress(processHandle windows.Handle) (uintptr, error) {
	return ScanMultiplePatterns(processHandle, StatusPatterns)
}
