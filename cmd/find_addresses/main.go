package main

import (
	"fmt"

	"dtap-go/internal/finder"
	"dtap-go/internal/process"
	"dtap-go/pkg/constants"

	"golang.org/x/sys/windows"
)

func main() {
	fmt.Println("DTAP-Go - Address Finder")
	fmt.Println("========================")

	// Buscar proceso de RO
	fmt.Println("Buscando proceso de Ragnarok Online...")
	roProcess, err := process.FindROProcess()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		fmt.Println("\nAsegúrate de que Ragnarok Online esté ejecutándose.")
		fmt.Println("Presiona Enter para salir...")
		fmt.Scanln()
		return
	}

	fmt.Printf("Proceso encontrado: %s (PID: %d)\n", roProcess.Name, roProcess.PID)

	// Abrir proceso
	handle, err := windows.OpenProcess(
		windows.PROCESS_VM_READ|windows.PROCESS_VM_WRITE|windows.PROCESS_VM_OPERATION,
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

	// Crear finder
	addressFinder := finder.NewAddressFinder(handle)

	// Valores de HP/SP desde constantes
	hpCurrent := uint32(constants.FinderHPCurrent)
	hpMax := uint32(constants.FinderHPMax)
	spCurrent := uint32(constants.FinderSPCurrent)
	spMax := uint32(constants.FinderSPMax)

	fmt.Printf("\nBuscando direcciones para:\n")
	fmt.Printf("HP: %d/%d\n", hpCurrent, hpMax)
	fmt.Printf("SP: %d/%d\n", spCurrent, spMax)
	fmt.Printf("\nNOTA: Si estos valores no coinciden con tu personaje,\n")
	fmt.Printf("edita las constantes en pkg/constants/constants.go:\n")
	fmt.Printf("- FinderHPCurrent = %d\n", constants.FinderHPCurrent)
	fmt.Printf("- FinderHPMax = %d\n", constants.FinderHPMax)
	fmt.Printf("- FinderSPCurrent = %d\n", constants.FinderSPCurrent)
	fmt.Printf("- FinderSPMax = %d\n", constants.FinderSPMax)
	fmt.Println()

	// Buscar direcciones
	results, err := addressFinder.FindHPAndSP(hpCurrent, hpMax, spCurrent, spMax)
	if err != nil {
		fmt.Printf("Error buscando direcciones: %v\n", err)
		fmt.Println("Presiona Enter para salir...")
		fmt.Scanln()
		return
	}

	// Mostrar resultados
	fmt.Printf("\n=== RESULTADOS ===\n")
	if len(results) == 0 {
		fmt.Println("No se encontraron direcciones válidas.")
		fmt.Println("Posibles causas:")
		fmt.Println("1. Los valores de HP/SP han cambiado")
		fmt.Println("2. El personaje no está en standby")
		fmt.Println("3. Los valores están en formato diferente")
	} else {
		fmt.Printf("Encontradas %d direcciones válidas:\n\n", len(results))

		for i, result := range results {
			fmt.Printf("--- Dirección %d ---\n", i+1)
			fmt.Printf("HP Current: 0x%X\n", result.HPCurrent)
			fmt.Printf("HP Max:     0x%X\n", result.HPMax)
			fmt.Printf("SP Current: 0x%X\n", result.SPCurrent)
			fmt.Printf("SP Max:     0x%X\n", result.SPMax)
			fmt.Printf("Distancia HP-SP: %d bytes\n", result.Distance)
			fmt.Println()
		}

		// Mostrar la primera dirección como recomendada
		if len(results) > 0 {
			best := results[0]
			fmt.Printf("=== DIRECCIÓN RECOMENDADA ===\n")
			fmt.Printf("HP Base: 0x%X\n", best.HPCurrent)
			fmt.Printf("SP Base: 0x%X\n", best.SPCurrent)
			fmt.Printf("Status:  0x%X (HP + 0x474)\n", best.HPCurrent+0x474)
		}
	}

	fmt.Println("\nPresiona Enter para salir...")
	fmt.Scanln()
}
