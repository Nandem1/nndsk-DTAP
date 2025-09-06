package injection

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"unsafe"

	"dtap-go/pkg/logger"

	"golang.org/x/sys/windows"
)

// WindowsInjector provides real Windows API DLL injection
type WindowsInjector struct {
	processID uint32
	logger    *logger.Logger
}

// NewWindowsInjector creates a new Windows injector
func NewWindowsInjector(processID uint32, logger *logger.Logger) *WindowsInjector {
	return &WindowsInjector{
		processID: processID,
		logger:    logger,
	}
}

// InjectDLL performs real DLL injection using Windows API
func (wi *WindowsInjector) InjectDLL(dllPath string) error {
	wi.logger.Info("Starting real Windows DLL injection", "dll_path", dllPath, "process_id", wi.processID)

	// Check if DLL file exists
	if _, err := os.Stat(dllPath); os.IsNotExist(err) {
		return fmt.Errorf("DLL file does not exist: %s", dllPath)
	}

	// Get absolute path
	absPath, err := filepath.Abs(dllPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Open the target process
	processHandle, err := windows.OpenProcess(
		windows.PROCESS_CREATE_THREAD|windows.PROCESS_QUERY_INFORMATION|
			windows.PROCESS_VM_OPERATION|windows.PROCESS_VM_WRITE|windows.PROCESS_VM_READ,
		false,
		wi.processID,
	)
	if err != nil {
		return fmt.Errorf("failed to open process: %w", err)
	}
	defer windows.CloseHandle(processHandle)

	// Get the address of LoadLibraryW
	kernel32 := windows.NewLazySystemDLL("kernel32.dll")
	loadLibraryW := kernel32.NewProc("LoadLibraryW")
	if loadLibraryW.Find() != nil {
		return fmt.Errorf("failed to find LoadLibraryW")
	}

	// Convert DLL path to UTF-16
	dllPathUTF16, err := syscall.UTF16PtrFromString(absPath)
	if err != nil {
		return fmt.Errorf("failed to convert DLL path to UTF-16: %w", err)
	}

	// Allocate memory in the target process
	remoteMemory, err := wi.virtualAllocEx(
		processHandle,
		nil,
		uintptr(len(absPath)*2), // UTF-16 takes 2 bytes per character
		windows.MEM_COMMIT|windows.MEM_RESERVE,
		windows.PAGE_READWRITE,
	)
	if err != nil {
		return fmt.Errorf("failed to allocate memory in target process: %w", err)
	}

	// Write the DLL path to the allocated memory
	var bytesWritten uintptr
	err = wi.writeProcessMemory(
		processHandle,
		remoteMemory,
		(*byte)(unsafe.Pointer(dllPathUTF16)),
		uintptr(len(absPath)*2),
		&bytesWritten,
	)
	if err != nil {
		wi.virtualFreeEx(processHandle, remoteMemory, 0, windows.MEM_RELEASE)
		return fmt.Errorf("failed to write DLL path to target process: %w", err)
	}

	// Create a remote thread to load the DLL
	threadHandle, err := wi.createRemoteThread(
		processHandle,
		nil,
		0,
		loadLibraryW.Addr(),
		remoteMemory,
		0,
		nil,
	)
	if err != nil {
		wi.virtualFreeEx(processHandle, remoteMemory, 0, windows.MEM_RELEASE)
		return fmt.Errorf("failed to create remote thread: %w", err)
	}
	defer windows.CloseHandle(threadHandle)

	// Wait for the thread to complete
	_, err = windows.WaitForSingleObject(threadHandle, windows.INFINITE)
	if err != nil {
		return fmt.Errorf("failed to wait for remote thread: %w", err)
	}

	// Get the exit code to check if LoadLibrary succeeded
	var exitCode uint32
	err = wi.getExitCodeThread(threadHandle, &exitCode)
	if err != nil {
		return fmt.Errorf("failed to get thread exit code: %w", err)
	}

	wi.logger.Info("LoadLibrary exit code", "exit_code", exitCode)

	if exitCode == 0 {
		// Try to get more detailed error information
		wi.logger.Error("LoadLibrary failed in target process - possible causes:")
		wi.logger.Error("1. DLL file not found or path incorrect")
		wi.logger.Error("2. DLL has missing dependencies")
		wi.logger.Error("3. Target process has anti-injection protection")
		wi.logger.Error("4. DLL is not compatible with target process architecture")
		return fmt.Errorf("LoadLibrary failed in target process (exit code: %d)", exitCode)
	}

	// Clean up allocated memory
	wi.virtualFreeEx(processHandle, remoteMemory, 0, windows.MEM_RELEASE)

	wi.logger.Info("DLL injected successfully using Windows API", "dll_path", absPath, "module_handle", exitCode)
	return nil
}

// Windows API functions using syscall
func (wi *WindowsInjector) virtualAllocEx(processHandle windows.Handle, address unsafe.Pointer, size uintptr, allocationType uint32, protect uint32) (unsafe.Pointer, error) {
	ret, _, err := syscall.Syscall6(
		syscall.NewLazyDLL("kernel32.dll").NewProc("VirtualAllocEx").Addr(),
		6,
		uintptr(processHandle),
		uintptr(address),
		size,
		uintptr(allocationType),
		uintptr(protect),
		0,
	)
	if ret == 0 {
		return nil, err
	}
	return unsafe.Pointer(ret), nil
}

func (wi *WindowsInjector) writeProcessMemory(processHandle windows.Handle, baseAddress unsafe.Pointer, buffer *byte, size uintptr, numberOfBytesWritten *uintptr) error {
	ret, _, err := syscall.Syscall6(
		syscall.NewLazyDLL("kernel32.dll").NewProc("WriteProcessMemory").Addr(),
		6,
		uintptr(processHandle),
		uintptr(baseAddress),
		uintptr(unsafe.Pointer(buffer)),
		size,
		uintptr(unsafe.Pointer(numberOfBytesWritten)),
		0,
	)
	if ret == 0 {
		return err
	}
	return nil
}

func (wi *WindowsInjector) createRemoteThread(processHandle windows.Handle, securityAttributes *windows.SecurityAttributes, stackSize uintptr, startAddress uintptr, parameter unsafe.Pointer, creationFlags uint32, threadID *uint32) (windows.Handle, error) {
	ret, _, err := syscall.Syscall9(
		syscall.NewLazyDLL("kernel32.dll").NewProc("CreateRemoteThread").Addr(),
		9,
		uintptr(processHandle),
		uintptr(unsafe.Pointer(securityAttributes)),
		stackSize,
		startAddress,
		uintptr(parameter),
		uintptr(creationFlags),
		uintptr(unsafe.Pointer(threadID)),
		0,
		0,
	)
	if ret == 0 {
		return 0, err
	}
	return windows.Handle(ret), nil
}

func (wi *WindowsInjector) getExitCodeThread(threadHandle windows.Handle, exitCode *uint32) error {
	ret, _, err := syscall.Syscall(
		syscall.NewLazyDLL("kernel32.dll").NewProc("GetExitCodeThread").Addr(),
		2,
		uintptr(threadHandle),
		uintptr(unsafe.Pointer(exitCode)),
		0,
	)
	if ret == 0 {
		return err
	}
	return nil
}

func (wi *WindowsInjector) virtualFreeEx(processHandle windows.Handle, address unsafe.Pointer, size uintptr, freeType uint32) error {
	ret, _, err := syscall.Syscall6(
		syscall.NewLazyDLL("kernel32.dll").NewProc("VirtualFreeEx").Addr(),
		4,
		uintptr(processHandle),
		uintptr(address),
		size,
		uintptr(freeType),
		0,
		0,
	)
	if ret == 0 {
		return err
	}
	return nil
}

// IsProcessValid checks if the target process is still valid
func (wi *WindowsInjector) IsProcessValid() bool {
	processHandle, err := windows.OpenProcess(windows.PROCESS_QUERY_INFORMATION, false, wi.processID)
	if err != nil {
		return false
	}
	defer windows.CloseHandle(processHandle)

	var exitCode uint32
	err = windows.GetExitCodeProcess(processHandle, &exitCode)
	if err != nil {
		return false
	}

	// STILL_ACTIVE is 259
	return exitCode == 259
}

// GetModuleHandle gets the handle of a loaded module in the target process
func (wi *WindowsInjector) GetModuleHandle(moduleName string) (uintptr, error) {
	processHandle, err := windows.OpenProcess(windows.PROCESS_QUERY_INFORMATION|windows.PROCESS_VM_READ, false, wi.processID)
	if err != nil {
		return 0, fmt.Errorf("failed to open process: %w", err)
	}
	defer windows.CloseHandle(processHandle)

	// This is a simplified version
	wi.logger.Debug("Getting module handle", "module_name", moduleName)
	return 0, fmt.Errorf("module handle retrieval not implemented")
}
