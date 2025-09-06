package main

import "C"
import (
	"syscall"
	"time"
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
	user32      *syscall.LazyDLL
	gdi32       *syscall.LazyDLL
	gameWindow  syscall.Handle
	dc          syscall.Handle
	pen         syscall.Handle
	brush       syscall.Handle
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

		// Initialize Windows API
		user32 = syscall.NewLazyDLL("user32.dll")
		gdi32 = syscall.NewLazyDLL("gdi32.dll")

		// Start rendering loop
		go startRenderingLoop()

		initialized = true
	case 0: // DLL_PROCESS_DETACH
		// Cleanup when DLL is unloaded
		cleanupResources()
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

// startRenderingLoop starts the main rendering loop
func startRenderingLoop() {
	// Find the game window
	gameWindow = findGameWindow()
	if gameWindow == 0 {
		return
	}

	// Get device context
	dc = getDeviceContext(gameWindow)
	if dc == 0 {
		return
	}

	// Create drawing objects
	pen = createPen(2, 0x00FF00)  // Green pen, 2px width
	brush = createBrush(0x00FF00) // Green brush

	// Main rendering loop
	ticker := time.NewTicker(16 * time.Millisecond) // ~60 FPS
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			renderFrame()
		}
	}
}

// findGameWindow finds the Ragnarok Online window
func findGameWindow() syscall.Handle {
	findWindow := user32.NewProc("FindWindowW")
	ret, _, _ := findWindow.Call(
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("Ragnarok"))),
		0,
	)
	return syscall.Handle(ret)
}

// getDeviceContext gets the device context for the window
func getDeviceContext(hwnd syscall.Handle) syscall.Handle {
	getDC := user32.NewProc("GetDC")
	ret, _, _ := getDC.Call(uintptr(hwnd))
	return syscall.Handle(ret)
}

// createPen creates a pen for drawing
func createPen(width int, color uint32) syscall.Handle {
	createPen := gdi32.NewProc("CreatePen")
	ret, _, _ := createPen.Call(
		uintptr(0), // PS_SOLID
		uintptr(width),
		uintptr(color),
	)
	return syscall.Handle(ret)
}

// createBrush creates a brush for drawing
func createBrush(color uint32) syscall.Handle {
	createSolidBrush := gdi32.NewProc("CreateSolidBrush")
	ret, _, _ := createSolidBrush.Call(uintptr(color))
	return syscall.Handle(ret)
}

// renderFrame renders a single frame
func renderFrame() {
	if sharedData == nil || !initialized {
		return
	}

	// Draw range display if enabled
	if sharedData.RangeDisplayEnabled {
		drawRangeDisplay()
	}

	// Draw skill markers if enabled
	if sharedData.SkillMarkersEnabled {
		drawSkillMarkers()
	}
}

// drawRangeDisplay draws the character's attack range
func drawRangeDisplay() {
	if dc == 0 || pen == 0 {
		return
	}

	// Select pen
	selectObject := gdi32.NewProc("SelectObject")
	selectObject.Call(uintptr(dc), uintptr(pen))

	// Draw circle (simplified - center of screen)
	ellipse := gdi32.NewProc("Ellipse")
	ellipse.Call(
		uintptr(dc),
		300, 200, // Left, Top
		500, 400, // Right, Bottom
	)
}

// drawSkillMarkers draws skill markers
func drawSkillMarkers() {
	if dc == 0 || brush == 0 {
		return
	}

	// Select brush
	selectObject := gdi32.NewProc("SelectObject")
	selectObject.Call(uintptr(dc), uintptr(brush))

	// Draw some test markers
	ellipse := gdi32.NewProc("Ellipse")
	ellipse.Call(
		uintptr(dc),
		100, 100, // Left, Top
		120, 120, // Right, Bottom
	)
}

// cleanupResources cleans up Windows resources
func cleanupResources() {
	if pen != 0 {
		deleteObject := gdi32.NewProc("DeleteObject")
		deleteObject.Call(uintptr(pen))
	}
	if brush != 0 {
		deleteObject := gdi32.NewProc("DeleteObject")
		deleteObject.Call(uintptr(brush))
	}
	if dc != 0 {
		releaseDC := user32.NewProc("ReleaseDC")
		releaseDC.Call(uintptr(gameWindow), uintptr(dc))
	}
}

func main() {
	// This is required for Go DLLs but won't be called
}
