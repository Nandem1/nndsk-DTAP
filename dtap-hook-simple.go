package main

import "C"
import (
	"unsafe"
)

// Shared memory structure for communication with DTAP
type DTAPSharedData struct {
	SkillMarkersEnabled bool
	RangeDisplayEnabled bool
	GroundZBias         int32
	AlphaLevel          int32
	SharedMemoryName    [64]byte
}

// Global variables
var (
	sharedData  *DTAPSharedData
	initialized bool
)

//export DllMain
func DllMain(hModule uintptr, ul_reason_for_call uint32, lpReserved uintptr) bool {
	switch ul_reason_for_call {
	case 1: // DLL_PROCESS_ATTACH
		// Initialize when DLL is loaded
		sharedData = &DTAPSharedData{
			SkillMarkersEnabled: true,
			RangeDisplayEnabled: true,
			GroundZBias:         0,
			AlphaLevel:          0,
		}
		copy(sharedData.SharedMemoryName[:], "DTAPSharedMemory")
		initialized = true
	case 0: // DLL_PROCESS_DETACH
		// Cleanup when DLL is unloaded
		sharedData = nil
		initialized = false
	}
	return true
}

//export EnableM2E
func EnableM2E(enabled bool) {
	if sharedData != nil {
		sharedData.SkillMarkersEnabled = enabled
	}
}

//export EnableBBE
func EnableBBE(enabled bool) {
	if sharedData != nil {
		sharedData.RangeDisplayEnabled = enabled
	}
}

//export SetGroundZBias
func SetGroundZBias(bias int32) {
	if sharedData != nil {
		sharedData.GroundZBias = bias
	}
}

//export SetAlphaLevel
func SetAlphaLevel(level int32) {
	if sharedData != nil {
		sharedData.AlphaLevel = level
	}
}

//export IsInitialized
func IsInitialized() bool {
	return initialized
}

//export GetSharedData
func GetSharedData() uintptr {
	if sharedData == nil {
		return 0
	}
	return uintptr(unsafe.Pointer(sharedData))
}

//export TestFunction
func TestFunction() int32 {
	return 42
}

func main() {
	// This is required for Go DLLs but won't be called
}
