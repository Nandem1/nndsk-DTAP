package input

import (
	"context"
	"syscall"
	"time"

	"dtap-go/internal/metrics"
	"dtap-go/pkg/constants"

	"golang.org/x/sys/windows"
)

// InputEvent representa un evento de entrada a enviar por PostMessageW
type InputEvent struct {
	Msg    uintptr
	WParam uintptr
	LParam uintptr
}

var (
	inputCh     chan InputEvent
	postMessage *syscall.LazyProc
	hwndGlobal  windows.Handle
)

// Start inicializa el dispatcher único de input.
// Debe llamarse una vez al inicio del ServiceManager.
func Start(ctx context.Context, hwnd windows.Handle) {
	if inputCh != nil {
		return
	}
	user32 := syscall.NewLazyDLL("user32.dll")
	postMessage = user32.NewProc("PostMessageW")
	hwndGlobal = hwnd
	inputCh = make(chan InputEvent, 256)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case ev := <-inputCh:
				// Llamada secuencial a PostMessageW
				postMessage.Call(
					uintptr(hwndGlobal),
					ev.Msg,
					ev.WParam,
					ev.LParam,
				)
				// Pequeño delay entre eventos para que el juego los procese
				time.Sleep(1 * time.Millisecond)

				// Reportar métricas
				metrics.GetMetrics().IncrementInputEvents()
			}
		}
	}()
}

// SendKey encola un KeyDown seguido de KeyUp para el keyCode dado.
func SendKey(keyCode uint32) {
	if inputCh == nil {
		return
	}
	inputCh <- InputEvent{Msg: uintptr(constants.WM_KEYDOWN), WParam: uintptr(keyCode), LParam: 0}
	inputCh <- InputEvent{Msg: uintptr(constants.WM_KEYUP), WParam: uintptr(keyCode), LParam: 0}
}

// SendLeftClick encola un click izquierdo completo (down+up).
func SendLeftClick() {
	if inputCh == nil {
		return
	}
	// Enviar mouse down y up con coordenadas (0,0) en lParam
	inputCh <- InputEvent{Msg: uintptr(constants.WM_LBUTTONDOWN), WParam: 0, LParam: 0}
	inputCh <- InputEvent{Msg: uintptr(constants.WM_LBUTTONUP), WParam: 0, LParam: 0}
}
