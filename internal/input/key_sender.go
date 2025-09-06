package input

import (
	"syscall"
	"time"

	"dtap-go/pkg/constants"

	"golang.org/x/sys/windows"
)

// KeySender handles sending keys to windows
type KeySender struct {
	user32      *syscall.LazyDLL
	postMessage *syscall.LazyProc
}

// NewKeySender creates a new key sender
func NewKeySender() *KeySender {
	user32 := syscall.NewLazyDLL("user32.dll")
	return &KeySender{
		user32:      user32,
		postMessage: user32.NewProc("PostMessageW"),
	}
}

// SendKey sends a key to the specified window
func (ks *KeySender) SendKey(hwnd windows.Handle, key string) error {
	keyCode := getKeyCode(key)
	if keyCode == 0 {
		return nil // Tecla no soportada, ignorar
	}

	// Verificar que el handle sea válido
	if hwnd == 0 {
		return nil // Handle inválido
	}

	// Enviar keydown
	ret1, _, _ := ks.postMessage.Call(
		uintptr(hwnd),
		constants.WM_KEYDOWN,
		uintptr(keyCode),
		0,
	)

	// Delay mínimo para skill canceling
	time.Sleep(1 * time.Millisecond)

	// Enviar keyup
	ret2, _, _ := ks.postMessage.Call(
		uintptr(hwnd),
		constants.WM_KEYUP,
		uintptr(keyCode),
		0,
	)

	// Verificar que ambos mensajes se enviaron correctamente
	if ret1 == 0 || ret2 == 0 {
		return nil // Error en el envío, pero no es crítico
	}

	return nil
}
