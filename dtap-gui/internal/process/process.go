package process

import (
	"fmt"
	"strings"
	"unsafe"

	"dtap-gui/pkg/constants"

	"golang.org/x/sys/windows"
)

type Process struct {
	PID  uint32
	Name string
}

// FindROProcess busca el proceso de Ragnarok Online
func FindROProcess() (*Process, error) {
	// Usar nombres de constantes
	roNames := constants.ROProcessNames

	fmt.Println("Buscando procesos de Ragnarok Online...")

	for _, name := range roNames {
		pid, err := findProcessByName(name)
		if err == nil {
			fmt.Printf("✓ Encontrado: %s (PID: %d)\n", name, pid)
			return &Process{
				PID:  pid,
				Name: name,
			}, nil
		}
	}

	// Si no encuentra ninguno, mostrar algunos procesos para debugging
	fmt.Println("No se encontró ningún proceso de RO. Algunos procesos disponibles:")
	listSomeProcesses()

	return nil, fmt.Errorf("no se encontró ningún proceso de Ragnarok Online")
}

// listSomeProcesses lista algunos procesos que podrían ser RO
func listSomeProcesses() {
	snapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return
	}
	defer windows.CloseHandle(snapshot)

	var pe32 windows.ProcessEntry32
	pe32.Size = uint32(unsafe.Sizeof(pe32))

	err = windows.Process32First(snapshot, &pe32)
	if err != nil {
		return
	}

	count := 0
	for {
		processName := windows.UTF16PtrToString(&pe32.ExeFile[0])
		if processName != "" {
			// Mostrar solo procesos que podrían ser juegos
			lowerName := strings.ToLower(processName)
			if strings.Contains(lowerName, "ragnarok") ||
				strings.Contains(lowerName, "ro") ||
				strings.Contains(lowerName, "game") ||
				strings.Contains(lowerName, "client") ||
				strings.Contains(lowerName, "gravity") {
				fmt.Printf("  - %s (PID: %d)\n", processName, pe32.ProcessID)
				count++
			}
		}

		err = windows.Process32Next(snapshot, &pe32)
		if err != nil {
			break
		}
	}

	if count == 0 {
		fmt.Println("  No se encontraron procesos que parezcan ser juegos.")
	}
}

// listAllProcesses lista todos los procesos para debugging
func listAllProcesses() {
	snapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return
	}
	defer windows.CloseHandle(snapshot)

	var pe32 windows.ProcessEntry32
	pe32.Size = uint32(unsafe.Sizeof(pe32))

	err = windows.Process32First(snapshot, &pe32)
	if err != nil {
		return
	}

	count := 0
	for {
		processName := windows.UTF16PtrToString(&pe32.ExeFile[0])
		if processName != "" {
			// Mostrar todos los procesos, pero limitar a 50 para no saturar
			if count < 50 {
				fmt.Printf("  - %s (PID: %d)\n", processName, pe32.ProcessID)
			}
			count++
		}

		err = windows.Process32Next(snapshot, &pe32)
		if err != nil {
			break
		}
	}

	if count > 50 {
		fmt.Printf("  ... y %d procesos más\n", count-50)
	}
}

// findProcessByName busca un proceso por nombre
func findProcessByName(name string) (uint32, error) {
	snapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return 0, err
	}
	defer windows.CloseHandle(snapshot)

	var pe32 windows.ProcessEntry32
	pe32.Size = uint32(unsafe.Sizeof(pe32))

	err = windows.Process32First(snapshot, &pe32)
	if err != nil {
		return 0, err
	}

	for {
		processName := windows.UTF16PtrToString(&pe32.ExeFile[0])
		if strings.EqualFold(processName, name) {
			return pe32.ProcessID, nil
		}

		err = windows.Process32Next(snapshot, &pe32)
		if err != nil {
			break
		}
	}

	return 0, fmt.Errorf("proceso no encontrado")
}

// OpenProcess abre un handle al proceso
func (p *Process) OpenHandle() (windows.Handle, error) {
	handle, err := windows.OpenProcess(
		windows.PROCESS_VM_READ|windows.PROCESS_VM_WRITE|windows.PROCESS_VM_OPERATION,
		false,
		p.PID,
	)
	if err != nil {
		return 0, fmt.Errorf("no se pudo abrir el proceso: %v", err)
	}

	return handle, nil
}
