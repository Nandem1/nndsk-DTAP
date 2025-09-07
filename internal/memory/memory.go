package memory

import (
	"fmt"
	"unsafe"

	"dtap-go/pkg/constants"
	"dtap-go/pkg/logger"

	"golang.org/x/sys/windows"
)

type Client struct {
	processHandle windows.Handle
	pid           uint32
	hpBaseAddr    uintptr
	spBaseAddr    uintptr
	statusAddr    uintptr
}

// NewClient crea un nuevo cliente de memoria
func NewClient(pid uint32) (*Client, error) {
	handle, err := windows.OpenProcess(
		constants.ProcessAccessFlags,
		false,
		pid,
	)
	if err != nil {
		return nil, fmt.Errorf("no se pudo abrir el proceso: %v", err)
	}

	client := &Client{
		processHandle: handle,
		pid:           pid,
	}

	// Usar direcciones de constantes
	logger.Info("Inicializando cliente de memoria con direcciones predefinidas")

	client.hpBaseAddr = constants.HPBaseAddress
	client.spBaseAddr = constants.SPBaseAddress
	client.statusAddr = constants.StatusBufferAddr

	logger.Info("HP base: 0x%X", client.hpBaseAddr)
	logger.Info("SP base: 0x%X", client.spBaseAddr)
	logger.Info("Status buffer: 0x%X", client.statusAddr)

	return client, nil
}

// Close cierra el handle del proceso
func (c *Client) Close() error {
	return windows.CloseHandle(c.processHandle)
}

// ReadMemory lee memoria del proceso
func (c *Client) ReadMemory(address uintptr, size uintptr) ([]byte, error) {
	buffer := make([]byte, size)
	var bytesRead uintptr

	err := windows.ReadProcessMemory(
		c.processHandle,
		address,
		&buffer[0],
		size,
		&bytesRead,
	)
	if err != nil {
		return nil, fmt.Errorf("error leyendo memoria: %v", err)
	}

	if bytesRead != size {
		return nil, fmt.Errorf("solo se leyeron %d bytes de %d", bytesRead, size)
	}

	return buffer, nil
}

// ReadUint32 lee un entero de 32 bits
func (c *Client) ReadUint32(address uintptr) (uint32, error) {
	data, err := c.ReadMemory(address, 4)
	if err != nil {
		return 0, err
	}

	return *(*uint32)(unsafe.Pointer(&data[0])), nil
}

// ReadHP lee HP actual y máximo
func (c *Client) ReadHP() (current, max uint32, err error) {
	current, err = c.ReadUint32(c.hpBaseAddr)
	if err != nil {
		return 0, 0, err
	}

	max, err = c.ReadUint32(c.hpBaseAddr + 4)
	if err != nil {
		return 0, 0, err
	}

	return current, max, nil
}

// ReadSP lee SP actual y máximo
func (c *Client) ReadSP() (current, max uint32, err error) {
	current, err = c.ReadUint32(c.spBaseAddr)
	if err != nil {
		return 0, 0, err
	}

	max, err = c.ReadUint32(c.spBaseAddr + 4)
	if err != nil {
		return 0, 0, err
	}

	return current, max, nil
}

// DebugAddresses muestra información de debugging sobre las direcciones
func (c *Client) DebugAddresses() {
	fmt.Printf("[DEBUG] Direcciones encontradas:\n")
	fmt.Printf("  HP Base: 0x%X\n", c.hpBaseAddr)
	fmt.Printf("  SP Base: 0x%X\n", c.spBaseAddr)
	fmt.Printf("  Status: 0x%X\n", c.statusAddr)

	// Probar diferentes offsets para HP
	fmt.Printf("\n[DEBUG] Probando diferentes offsets para HP:\n")
	for i := 0; i < 5; i++ {
		addr := c.hpBaseAddr + uintptr(i*4)
		val, err := c.ReadUint32(addr)
		if err == nil {
			fmt.Printf("  HP+%d (0x%X): 0x%X (%d)\n", i*4, addr, val, val)
		}
	}

	// Probar diferentes offsets para SP
	fmt.Printf("\n[DEBUG] Probando diferentes offsets para SP:\n")
	for i := 0; i < 5; i++ {
		addr := c.spBaseAddr + uintptr(i*4)
		val, err := c.ReadUint32(addr)
		if err == nil {
			fmt.Printf("  SP+%d (0x%X): 0x%X (%d)\n", i*4, addr, val, val)
		}
	}
}

// GetStatusBufferAddress retorna la dirección del status buffer
func (c *Client) GetStatusBufferAddress() (uintptr, error) {
	return c.statusAddr, nil
}
