package input

import (
	"syscall"
)

// KeyDetector handles key state detection
type KeyDetector struct {
	user32      *syscall.LazyDLL
	getKeyState *syscall.LazyProc
}

// NewKeyDetector creates a new key detector
func NewKeyDetector() *KeyDetector {
	user32 := syscall.NewLazyDLL("user32.dll")
	return &KeyDetector{
		user32:      user32,
		getKeyState: user32.NewProc("GetAsyncKeyState"),
	}
}

// IsKeyPressed checks if a key is currently pressed
func (kd *KeyDetector) IsKeyPressed(key string) bool {
	keyCode := getKeyCode(key)
	if keyCode == 0 {
		return false
	}

	// GetAsyncKeyState returns:
	// - High bit set if key is currently pressed
	// - Low bit set if key was pressed since last call
	ret, _, _ := kd.getKeyState.Call(uintptr(keyCode))

	// Check if high bit is set (key is currently pressed)
	return (ret & 0x8000) != 0
}

// IsModifierKey checks if a key is a modifier key (Shift, Ctrl, Alt)
func (kd *KeyDetector) IsModifierKey(key string) bool {
	modifierKeys := map[string]bool{
		"SHIFT": true,
		"CTRL":  true,
		"ALT":   true,
	}
	return modifierKeys[key]
}

// IsAnyNonModifierKeyPressed checks if any non-modifier key is pressed
func (kd *KeyDetector) IsAnyNonModifierKeyPressed(keys []string) bool {
	for _, key := range keys {
		if !kd.IsModifierKey(key) && kd.IsKeyPressed(key) {
			return true
		}
	}
	return false
}

// GetPressedKeys returns a list of currently pressed keys from the given list
func (kd *KeyDetector) GetPressedKeys(keys []string) []string {
	var pressed []string

	for _, key := range keys {
		if kd.IsKeyPressed(key) {
			pressed = append(pressed, key)
		}
	}

	return pressed
}

// IsAnyKeyPressed checks if any key from the list is currently pressed
func (kd *KeyDetector) IsAnyKeyPressed(keys []string) bool {
	for _, key := range keys {
		if kd.IsKeyPressed(key) {
			return true
		}
	}
	return false
}
