package input

import (
	"syscall"
	"time"

	"dtap-go/pkg/constants"

	"golang.org/x/sys/windows"
)

// MouseSender handles sending mouse clicks to windows
type MouseSender struct {
	user32      *syscall.LazyDLL
	postMessage *syscall.LazyProc
}

// NewMouseSender creates a new mouse sender
func NewMouseSender() *MouseSender {
	user32 := syscall.NewLazyDLL("user32.dll")
	return &MouseSender{
		user32:      user32,
		postMessage: user32.NewProc("PostMessageW"),
	}
}

// SendLeftClick sends a left mouse click to the specified window
func (ms *MouseSender) SendLeftClick(hwnd windows.Handle, x, y int32) error {
	if hwnd == 0 {
		return nil // Handle inválido
	}

	// Crear lParam con coordenadas (y << 16 | x)
	lParam := uintptr((uint32(y) << 16) | uint32(x))

	// Enviar mousedown
	ret1, _, _ := ms.postMessage.Call(
		uintptr(hwnd),
		constants.WM_LBUTTONDOWN,
		0x0001, // MK_LBUTTON
		lParam,
	)

	// Delay mínimo para skill canceling
	time.Sleep(1 * time.Millisecond)

	// Enviar mouseup
	ret2, _, _ := ms.postMessage.Call(
		uintptr(hwnd),
		constants.WM_LBUTTONUP,
		0x0000, // No buttons pressed
		lParam,
	)

	// Verificar que ambos mensajes se enviaron correctamente
	if ret1 == 0 || ret2 == 0 {
		return nil // Error en el envío, pero no es crítico
	}

	return nil
}

// SendRightClick sends a right mouse click to the specified window
func (ms *MouseSender) SendRightClick(hwnd windows.Handle, x, y int32) error {
	if hwnd == 0 {
		return nil // Handle inválido
	}

	// Crear lParam con coordenadas (y << 16 | x)
	lParam := uintptr((uint32(y) << 16) | uint32(x))

	// Enviar mousedown
	ret1, _, _ := ms.postMessage.Call(
		uintptr(hwnd),
		constants.WM_RBUTTONDOWN,
		0x0002, // MK_RBUTTON
		lParam,
	)

	// Delay mínimo para skill canceling
	time.Sleep(1 * time.Millisecond)

	// Enviar mouseup
	ret2, _, _ := ms.postMessage.Call(
		uintptr(hwnd),
		constants.WM_RBUTTONUP,
		0x0000, // No buttons pressed
		lParam,
	)

	// Verificar que ambos mensajes se enviaron correctamente
	if ret1 == 0 || ret2 == 0 {
		return nil // Error en el envío, pero no es crítico
	}

	return nil
}
