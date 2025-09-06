package input

import (
	"syscall"

	"golang.org/x/sys/windows"
)

var (
	user32          = syscall.NewLazyDLL("user32.dll")
	procPostMessage = user32.NewProc("PostMessageW")
)

// SendKey envía una tecla al proceso (función legacy para compatibilidad)
func SendKey(hwnd windows.Handle, key string) error {
	sender := NewKeySender()
	return sender.SendKey(hwnd, key)
}
