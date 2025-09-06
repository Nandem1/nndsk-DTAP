package main

import (
	"fmt"
	"log"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	kernel32               = syscall.NewLazyDLL("kernel32.dll")
	procVirtualAllocEx     = kernel32.NewProc("VirtualAllocEx")
	procVirtualFreeEx      = kernel32.NewProc("VirtualFreeEx")
	procWriteProcessMemory = kernel32.NewProc("WriteProcessMemory")
	procCreateRemoteThread = kernel32.NewProc("CreateRemoteThread")
	procGetExitCodeThread  = kernel32.NewProc("GetExitCodeThread")
	procLoadLibraryW       = kernel32.NewProc("LoadLibraryW")
)

func main() {
	// Find the Ragnarok process
	processID := uint32(8904) // HoneyRO.exe PID from your log

	fmt.Printf("Testing DLL injection into process %d...\n", processID)

	// Open process
	processHandle, err := windows.OpenProcess(0x1F0FFF, false, processID)
	if err != nil {
		log.Fatalf("Failed to open process: %v", err)
	}
	defer windows.CloseHandle(processHandle)

	// Prepare DLL path
	dllPath := "test-minimal.dll"
	dllPathUTF16, err := syscall.UTF16PtrFromString(dllPath)
	if err != nil {
		log.Fatalf("Failed to convert DLL path: %v", err)
	}

	// Allocate memory in target process
	remoteMemory, _, _ := procVirtualAllocEx.Call(
		uintptr(processHandle),
		0,
		uintptr(len(dllPath)*2),
		0x1000|0x2000, // MEM_COMMIT|MEM_RESERVE
		0x04,          // PAGE_READWRITE
	)
	if remoteMemory == 0 {
		log.Fatalf("Failed to allocate memory")
	}
	defer procVirtualFreeEx.Call(uintptr(processHandle), remoteMemory, 0, 0x8000) // MEM_RELEASE

	// Write DLL path to target process
	var bytesWritten uintptr
	ret, _, _ := procWriteProcessMemory.Call(
		uintptr(processHandle),
		remoteMemory,
		uintptr(unsafe.Pointer(dllPathUTF16)),
		uintptr(len(dllPath)*2),
		uintptr(unsafe.Pointer(&bytesWritten)),
	)
	if ret == 0 {
		log.Fatalf("Failed to write DLL path")
	}

	fmt.Printf("DLL path written to target process: %d bytes\n", bytesWritten)

	// Create remote thread
	threadHandle, _, _ := procCreateRemoteThread.Call(
		uintptr(processHandle),
		0,
		0,
		procLoadLibraryW.Addr(),
		remoteMemory,
		0,
		0,
	)
	if threadHandle == 0 {
		log.Fatalf("Failed to create remote thread")
	}
	defer windows.CloseHandle(windows.Handle(threadHandle))

	// Wait for thread completion
	_, err = windows.WaitForSingleObject(windows.Handle(threadHandle), windows.INFINITE)
	if err != nil {
		log.Fatalf("Failed to wait for thread: %v", err)
	}

	// Get exit code
	var exitCode uint32
	ret, _, _ = procGetExitCodeThread.Call(threadHandle, uintptr(unsafe.Pointer(&exitCode)))
	if ret == 0 {
		log.Fatalf("Failed to get exit code")
	}

	fmt.Printf("LoadLibrary exit code: %d\n", exitCode)

	if exitCode == 0 {
		fmt.Println("❌ LoadLibrary failed - DLL could not be loaded")
		fmt.Println("Possible causes:")
		fmt.Println("1. DLL file not found")
		fmt.Println("2. Missing dependencies")
		fmt.Println("3. Anti-injection protection")
		fmt.Println("4. Architecture mismatch")
	} else {
		fmt.Printf("✅ LoadLibrary succeeded! Module handle: 0x%x\n", exitCode)
	}
}
