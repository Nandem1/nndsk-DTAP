package input

import (
	"dtap-go/pkg/constants"
)

// getKeyCode converts a string key to its Windows virtual key code
func getKeyCode(key string) uint32 {
	keyMap := map[string]uint32{
		// Function keys
		"F1":  constants.KeyF1,
		"F2":  constants.KeyF2,
		"F3":  constants.KeyF3,
		"F4":  constants.KeyF4,
		"F5":  constants.KeyF5,
		"F6":  constants.KeyF6,
		"F7":  constants.KeyF7,
		"F8":  constants.KeyF8,
		"F9":  constants.KeyF9,
		"F10": constants.KeyF10,
		"F11": constants.KeyF11,
		"F12": constants.KeyF12,

		// Modifier keys
		"SHIFT": constants.KeyShift,
		"CTRL":  constants.KeyCtrl,
		"ALT":   constants.KeyAlt,

		// Number keys
		"1": constants.Key1,
		"2": constants.Key2,
		"3": constants.Key3,
		"4": constants.Key4,
		"5": constants.Key5,
		"6": constants.Key6,
		"7": constants.Key7,
		"8": constants.Key8,
		"9": constants.Key9,
		"0": constants.Key0,

		// Letter keys
		"A": constants.KeyA,
		"B": constants.KeyB,
		"C": constants.KeyC,
		"D": constants.KeyD,
		"E": constants.KeyE,
		"F": constants.KeyF,
		"G": constants.KeyG,
		"H": constants.KeyH,
		"I": constants.KeyI,
		"J": constants.KeyJ,
		"K": constants.KeyK,
		"L": constants.KeyL,
		"M": constants.KeyM,
		"N": constants.KeyN,
		"O": constants.KeyO,
		"P": constants.KeyP,
		"Q": constants.KeyQ,
		"R": constants.KeyR,
		"S": constants.KeyS,
		"T": constants.KeyT,
		"U": constants.KeyU,
		"V": constants.KeyV,
		"W": constants.KeyW,
		"X": constants.KeyX,
		"Y": constants.KeyY,
		"Z": constants.KeyZ,
	}

	return keyMap[key]
}

// IsValidKey checks if a key string is valid
func IsValidKey(key string) bool {
	return getKeyCode(key) != 0
}

// GetSupportedKeys returns a list of all supported keys
func GetSupportedKeys() []string {
	return []string{
		// Function keys
		"F1", "F2", "F3", "F4", "F5", "F6", "F7", "F8", "F9", "F10", "F11", "F12",
		// Modifier keys
		"SHIFT", "CTRL", "ALT",
		// Number keys
		"1", "2", "3", "4", "5", "6", "7", "8", "9", "0",
		// Letter keys
		"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
	}
}
