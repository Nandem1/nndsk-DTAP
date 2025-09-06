package injection

import (
	"fmt"
	"path/filepath"
	"sync"

	"dtap-go/pkg/logger"
)

// InjectionManager manages the DLL injection system
type InjectionManager struct {
	windowsInjector    *WindowsInjector
	simpleHookInjector *SimpleHookInjector
	gepardBypass       *GepardBypass
	realGepardBypass   *RealGepardBypass
	logger             *logger.Logger
	processID          uint32
	dllPath            string
	mu                 sync.RWMutex
	injected           bool
}

// NewInjectionManager creates a new injection manager
func NewInjectionManager(processID uint32, dllPath string, logger *logger.Logger) *InjectionManager {
	return &InjectionManager{
		processID: processID,
		dllPath:   dllPath,
		logger:    logger,
	}
}

// Initialize initializes the injection system
func (im *InjectionManager) Initialize() error {
	im.mu.Lock()
	defer im.mu.Unlock()

	if im.injected {
		return fmt.Errorf("injection system already initialized")
	}

	// Create Windows injector
	im.windowsInjector = NewWindowsInjector(im.processID, im.logger)

	// Create Simple Hook injector (like SimpleROHook)
	im.simpleHookInjector = NewSimpleHookInjector(im.processID, im.dllPath, im.logger)

	// Create Gepard bypass
	im.gepardBypass = NewGepardBypass(im.logger)

	// Create real Gepard bypass
	im.realGepardBypass = NewRealGepardBypass(im.logger)

	im.logger.Info("Injection system initialized successfully")
	return nil
}

// Inject injects the DLL into the target process
func (im *InjectionManager) Inject() error {
	im.mu.Lock()
	defer im.mu.Unlock()

	if im.injected {
		return fmt.Errorf("DLL already injected")
	}

	if im.windowsInjector == nil {
		return fmt.Errorf("injection system not initialized")
	}

	// Get absolute path to DLL
	absDllPath, err := filepath.Abs(im.dllPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute DLL path: %w", err)
	}

	// Start real Gepard bypass if needed
	if im.realGepardBypass.IsGepardPresent() {
		im.logger.Info("Gepard detected, starting real bypass...")
		if err := im.realGepardBypass.StartBypass(im.processID); err != nil {
			im.logger.Warn("Failed to start real Gepard bypass, trying basic bypass", "error", err)
			// Fallback to basic bypass
			if im.gepardBypass.IsGepardPresent() {
				if err := im.gepardBypass.StartBypass(); err != nil {
					im.logger.Warn("Failed to start basic Gepard bypass", "error", err)
				} else {
					im.logger.Info("Basic Gepard bypass started successfully")
				}
			}
		} else {
			im.logger.Info("Real Gepard bypass started successfully")
		}
	} else {
		im.logger.Info("No Gepard detected, proceeding with normal injection")
	}

	// Try SimpleROHook-style injection first
	im.logger.Info("Attempting SimpleROHook-style injection...")
	err = im.simpleHookInjector.InjectDLL()
	if err != nil {
		im.logger.Warn("SimpleROHook injection failed, trying Windows API injection", "error", err)
		// Fallback to Windows API injection
		err = im.windowsInjector.InjectDLL(absDllPath)
		if err != nil {
			return fmt.Errorf("all injection methods failed: %w", err)
		}
	} else {
		im.logger.Info("SimpleROHook-style injection successful")
	}

	im.injected = true
	im.logger.Info("DLL injected successfully", "dll_path", absDllPath)
	return nil
}

// IsInjected returns whether the DLL is injected
func (im *InjectionManager) IsInjected() bool {
	im.mu.RLock()
	defer im.mu.RUnlock()
	return im.injected
}

// IsGepardPresent returns whether Gepard is present in the process
func (im *InjectionManager) IsGepardPresent() bool {
	if im.gepardBypass == nil {
		return false
	}
	return im.gepardBypass.IsGepardPresent()
}

// IsGepardBypassActive returns whether the Gepard bypass is active
func (im *InjectionManager) IsGepardBypassActive() bool {
	if im.gepardBypass == nil {
		return false
	}
	return im.gepardBypass.IsBypassActive()
}

// SetROWindowHandle sets the Ragnarok Online window handle
func (im *InjectionManager) SetROWindowHandle(hwnd uintptr) error {
	im.logger.Debug("RO window handle set", "handle", hwnd)
	return nil
}

// EnableM2E enables skill markers (M2E)
func (im *InjectionManager) EnableM2E(enabled bool) error {
	im.logger.Info("M2E (skill markers)", "enabled", enabled)
	return nil
}

// Cleanup cleans up the injection system
func (im *InjectionManager) Cleanup() error {
	im.mu.Lock()
	defer im.mu.Unlock()

	var err error

	// Stop real Gepard bypass
	if im.realGepardBypass != nil {
		if bypassErr := im.realGepardBypass.StopBypass(); bypassErr != nil {
			im.logger.Error("Failed to stop real Gepard bypass: %v", bypassErr)
			if err == nil {
				err = bypassErr
			}
		}
		im.realGepardBypass = nil
	}

	// Stop basic Gepard bypass
	if im.gepardBypass != nil {
		if bypassErr := im.gepardBypass.StopBypass(); bypassErr != nil {
			im.logger.Error("Failed to stop Gepard bypass: %v", bypassErr)
			if err == nil {
				err = bypassErr
			}
		}
		im.gepardBypass = nil
	}

	// Cleanup simple hook injector
	if im.simpleHookInjector != nil {
		if hookErr := im.simpleHookInjector.Cleanup(); hookErr != nil {
			im.logger.Error("Failed to cleanup simple hook injector: %v", hookErr)
			if err == nil {
				err = hookErr
			}
		}
		im.simpleHookInjector = nil
	}

	// Windows injector doesn't need cleanup
	im.windowsInjector = nil

	im.injected = false
	im.logger.Info("Injection system cleaned up")
	return err
}
