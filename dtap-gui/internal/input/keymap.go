package input

import (
	"strings"

	"golang.org/x/sys/windows"
)

// KeyMap mapea nombres de teclas a códigos virtuales de Windows
var KeyMap = map[string]uint32{
	// Teclas de función
	"F1":  0x70,
	"F2":  0x71,
	"F3":  0x72,
	"F4":  0x73,
	"F5":  0x74,
	"F6":  0x75,
	"F7":  0x76,
	"F8":  0x77,
	"F9":  0x78,
	"F10": 0x79,
	"F11": 0x7A,
	"F12": 0x7B,

	// Teclas alfanuméricas
	"A": 0x41,
	"B": 0x42,
	"C": 0x43,
	"D": 0x44,
	"E": 0x45,
	"F": 0x46,
	"G": 0x47,
	"H": 0x48,
	"I": 0x49,
	"J": 0x4A,
	"K": 0x4B,
	"L": 0x4C,
	"M": 0x4D,
	"N": 0x4E,
	"O": 0x4F,
	"P": 0x50,
	"Q": 0x51,
	"R": 0x52,
	"S": 0x53,
	"T": 0x54,
	"U": 0x55,
	"V": 0x56,
	"W": 0x57,
	"X": 0x58,
	"Y": 0x59,
	"Z": 0x5A,

	// Números
	"0": 0x30,
	"1": 0x31,
	"2": 0x32,
	"3": 0x33,
	"4": 0x34,
	"5": 0x35,
	"6": 0x36,
	"7": 0x37,
	"8": 0x38,
	"9": 0x39,

	// Teclas especiales
	"SPACE":     0x20,
	"ENTER":     0x0D,
	"ESCAPE":    0x1B,
	"TAB":       0x09,
	"SHIFT":     0x10,
	"CTRL":      0x11,
	"ALT":       0x12,
	"BACKSPACE": 0x08,
	"DELETE":    0x2E,
	"INSERT":    0x2D,
	"HOME":      0x24,
	"END":       0x23,
	"PAGEUP":    0x21,
	"PAGEDOWN":  0x22,

	// Flechas
	"UP":    0x26,
	"DOWN":  0x28,
	"LEFT":  0x25,
	"RIGHT": 0x27,

	// Teclas del teclado numérico
	"NUMPAD0": 0x60,
	"NUMPAD1": 0x61,
	"NUMPAD2": 0x62,
	"NUMPAD3": 0x63,
	"NUMPAD4": 0x64,
	"NUMPAD5": 0x65,
	"NUMPAD6": 0x66,
	"NUMPAD7": 0x67,
	"NUMPAD8": 0x68,
	"NUMPAD9": 0x69,
	"NUMPAD*": 0x6A,
	"NUMPAD+": 0x6B,
	"NUMPAD-": 0x6D,
	"NUMPAD.": 0x6E,
	"NUMPAD/": 0x6F,
}

// GetKeyCode retorna el código virtual de una tecla
func GetKeyCode(keyName string) (uint32, bool) {
	keyName = strings.ToUpper(strings.TrimSpace(keyName))
	code, exists := KeyMap[keyName]
	return code, exists
}

// SendKeyByName envía una tecla por nombre
func SendKeyByName(hwnd windows.Handle, keyName string) {
	keyCode, exists := GetKeyCode(keyName)
	if !exists {
		return
	}
	SendKeyWithHandle(hwnd, keyCode)
}

// SendMouseClickByName envía un click del mouse
func SendMouseClickByName(hwnd windows.Handle) {
	SendMouseClickWithHandle(hwnd)
}
