package monitoring

import (
	"fmt"
	"time"
)

// DisplayLogger maneja logs que no interfieren con el display
type DisplayLogger struct {
	lastLogTime time.Time
	logFreq     time.Duration
}

// NewDisplayLogger crea un nuevo logger para display
func NewDisplayLogger(logFreq time.Duration) *DisplayLogger {
	return &DisplayLogger{
		logFreq: logFreq,
	}
}

// LogInfo muestra un log sin interferir con el display
func (dl *DisplayLogger) LogInfo(format string, args ...interface{}) {
	now := time.Now()

	// Solo mostrar si ha pasado suficiente tiempo
	if now.Sub(dl.lastLogTime) < dl.logFreq {
		return
	}

	dl.lastLogTime = now

	// Limpiar línea actual y mostrar log
	fmt.Print("\r\033[2K") // Limpiar línea
	fmt.Printf("[%s] %s\n", time.Now().Format("15:04:05"), fmt.Sprintf(format, args...))
}

// LogError muestra un error sin interferir con el display
func (dl *DisplayLogger) LogError(format string, args ...interface{}) {
	// Los errores siempre se muestran
	fmt.Print("\r\033[2K") // Limpiar línea
	fmt.Printf("[%s] ERROR: %s\n", time.Now().Format("15:04:05"), fmt.Sprintf(format, args...))
}

// ClearLine limpia la línea actual
func (dl *DisplayLogger) ClearLine() {
	fmt.Print("\r\033[2K")
}

// NewLine agrega una nueva línea
func (dl *DisplayLogger) NewLine() {
	fmt.Println()
}
