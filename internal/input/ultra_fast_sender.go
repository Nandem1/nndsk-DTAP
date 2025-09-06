package input

import (
	"sync"
	"time"

	"dtap-go/pkg/constants"

	"golang.org/x/sys/windows"
)

// UltraFastSender provides the fastest possible input for WoE
type UltraFastSender struct {
	user32      *windows.LazyDLL
	sendInput   *windows.LazyProc
	postMessage *windows.LazyProc
	keybdEvent  *windows.LazyProc

	// Input queue for batching
	inputQueue chan InputEvent
	queueMutex sync.Mutex

	// Performance tracking
	lastSendTime time.Time
	sendCount    int64
	windowHandle windows.Handle
}

// InputEvent represents a queued input event
type InputEvent struct {
	Type     InputType
	KeyCode  uint32
	Duration time.Duration
}

// InputType defines the type of input event
type InputType int

const (
	KeyDown InputType = iota
	KeyUp
	KeyPress
)

// NewUltraFastSender creates the fastest input sender
func NewUltraFastSender() *UltraFastSender {
	user32 := windows.NewLazySystemDLL("user32.dll")

	sender := &UltraFastSender{
		user32:      user32,
		sendInput:   user32.NewProc("SendInput"),
		postMessage: user32.NewProc("PostMessageW"),
		keybdEvent:  user32.NewProc("keybd_event"),
		inputQueue:  make(chan InputEvent, 100), // Buffer for 100 events
	}

	// Start input processing goroutine
	go sender.processInputQueue()

	return sender
}

// SetWindowHandle sets the target window
func (ufs *UltraFastSender) SetWindowHandle(hwnd windows.Handle) {
	ufs.windowHandle = hwnd
}

// SendKeyUltraFast sends a key with minimal latency
func (ufs *UltraFastSender) SendKeyUltraFast(key string) error {
	keyCode := getKeyCode(key)
	if keyCode == 0 {
		return nil
	}

	// Use keybd_event for maximum speed (bypasses message queue)
	ufs.keybdEvent.Call(
		uintptr(keyCode),
		0,
		0, // KEYEVENTF_KEYDOWN
		0,
	)

	// Minimal delay for skill canceling
	time.Sleep(100 * time.Microsecond) // 0.1ms instead of 1ms

	ufs.keybdEvent.Call(
		uintptr(keyCode),
		0,
		2, // KEYEVENTF_KEYUP
		0,
	)

	ufs.sendCount++
	ufs.lastSendTime = time.Now()

	return nil
}

// SendKeyBatch sends multiple keys in rapid succession
func (ufs *UltraFastSender) SendKeyBatch(keys []string) error {
	for _, key := range keys {
		if err := ufs.SendKeyUltraFast(key); err != nil {
			return err
		}
		// Ultra-minimal delay between keys
		time.Sleep(50 * time.Microsecond) // 0.05ms
	}
	return nil
}

// QueueKey queues a key for batch processing
func (ufs *UltraFastSender) QueueKey(key string, inputType InputType) {
	keyCode := getKeyCode(key)
	if keyCode == 0 {
		return
	}

	event := InputEvent{
		Type:    inputType,
		KeyCode: keyCode,
	}

	select {
	case ufs.inputQueue <- event:
	default:
		// Queue full, drop event
	}
}

// processInputQueue processes queued input events
func (ufs *UltraFastSender) processInputQueue() {
	for event := range ufs.inputQueue {
		switch event.Type {
		case KeyDown:
			ufs.keybdEvent.Call(
				uintptr(event.KeyCode),
				0,
				0, // KEYEVENTF_KEYDOWN
				0,
			)
		case KeyUp:
			ufs.keybdEvent.Call(
				uintptr(event.KeyCode),
				0,
				2, // KEYEVENTF_KEYUP
				0,
			)
		case KeyPress:
			ufs.keybdEvent.Call(
				uintptr(event.KeyCode),
				0,
				0, // KEYEVENTF_KEYDOWN
				0,
			)
			time.Sleep(100 * time.Microsecond)
			ufs.keybdEvent.Call(
				uintptr(event.KeyCode),
				0,
				2, // KEYEVENTF_KEYUP
				0,
			)
		}
	}
}

// SendMouseClickUltraFast sends mouse click with minimal latency
func (ufs *UltraFastSender) SendMouseClickUltraFast(x, y int32) error {
	if ufs.windowHandle == 0 {
		return nil
	}

	// Convert coordinates to lParam
	lParam := uintptr((y << 16) | (x & 0xFFFF))

	// Send mouse down
	ufs.postMessage.Call(
		uintptr(ufs.windowHandle),
		constants.WM_LBUTTONDOWN,
		1, // MK_LBUTTON
		lParam,
	)

	// Ultra-minimal delay
	time.Sleep(50 * time.Microsecond)

	// Send mouse up
	ufs.postMessage.Call(
		uintptr(ufs.windowHandle),
		constants.WM_LBUTTONUP,
		0,
		lParam,
	)

	return nil
}

// GetPerformanceStats returns input performance statistics
func (ufs *UltraFastSender) GetPerformanceStats() (count int64, lastSend time.Time) {
	return ufs.sendCount, ufs.lastSendTime
}

// FlushQueue processes all queued events immediately
func (ufs *UltraFastSender) FlushQueue() {
	// Process remaining events in queue
	for {
		select {
		case event := <-ufs.inputQueue:
			switch event.Type {
			case KeyDown:
				ufs.keybdEvent.Call(
					uintptr(event.KeyCode),
					0,
					0,
					0,
				)
			case KeyUp:
				ufs.keybdEvent.Call(
					uintptr(event.KeyCode),
					0,
					2,
					0,
				)
			case KeyPress:
				ufs.keybdEvent.Call(
					uintptr(event.KeyCode),
					0,
					0,
					0,
				)
				time.Sleep(100 * time.Microsecond)
				ufs.keybdEvent.Call(
					uintptr(event.KeyCode),
					0,
					2,
					0,
				)
			}
		default:
			return
		}
	}
}
