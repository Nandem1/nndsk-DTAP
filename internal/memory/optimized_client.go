package memory

import (
	"context"
	"fmt"
	"sync"
	"time"
	"unsafe"

	"dtap-go/pkg/constants"
	"dtap-go/pkg/logger"

	"golang.org/x/sys/windows"
)

// OptimizedClient provides high-performance memory operations for WoE
type OptimizedClient struct {
	processHandle windows.Handle
	pid           uint32
	
	// Cached addresses
	hpBaseAddr    uintptr
	spBaseAddr    uintptr
	statusAddr    uintptr
	
	// Memory cache for frequent reads
	cache         *MemoryCache
	cacheMutex    sync.RWMutex
	
	// Background refresh
	refreshTicker *time.Ticker
	ctx           context.Context
	cancel        context.CancelFunc
}

// MemoryCache holds frequently accessed memory values
type MemoryCache struct {
	HP        uint32
	MaxHP     uint32
	SP        uint32
	MaxSP     uint32
	LastUpdate time.Time
	Valid     bool
}

// NewOptimizedClient creates a high-performance memory client
func NewOptimizedClient(pid uint32) (*OptimizedClient, error) {
	handle, err := windows.OpenProcess(
		constants.ProcessAccessFlags,
		false,
		pid,
	)
	if err != nil {
		return nil, fmt.Errorf("no se pudo abrir el proceso: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	
	client := &OptimizedClient{
		processHandle: handle,
		pid:           pid,
		hpBaseAddr:    constants.HPBaseAddress,
		spBaseAddr:    constants.SPBaseAddress,
		statusAddr:    constants.StatusBufferAddr,
		cache:         &MemoryCache{},
		ctx:           ctx,
		cancel:        cancel,
	}

	// Start background refresh at 50ms intervals (20 FPS)
	client.refreshTicker = time.NewTicker(50 * time.Millisecond)
	go client.backgroundRefresh()

	logger.Info("Cliente de memoria optimizado inicializado")
	return client, nil
}

// backgroundRefresh continuously updates the memory cache
func (c *OptimizedClient) backgroundRefresh() {
	for {
		select {
		case <-c.ctx.Done():
			return
		case <-c.refreshTicker.C:
			c.updateCache()
		}
	}
}

// updateCache reads and caches HP/SP values
func (c *OptimizedClient) updateCache() {
	c.cacheMutex.Lock()
	defer c.cacheMutex.Unlock()

	// Read HP values
	if hp, err := c.readUint32Fast(c.hpBaseAddr); err == nil {
		c.cache.HP = hp
	}
	if maxHP, err := c.readUint32Fast(c.hpBaseAddr + 4); err == nil {
		c.cache.MaxHP = maxHP
	}

	// Read SP values
	if sp, err := c.readUint32Fast(c.spBaseAddr); err == nil {
		c.cache.SP = sp
	}
	if maxSP, err := c.readUint32Fast(c.spBaseAddr + 4); err == nil {
		c.cache.MaxSP = maxSP
	}

	c.cache.LastUpdate = time.Now()
	c.cache.Valid = true
}

// readUint32Fast performs optimized memory read
func (c *OptimizedClient) readUint32Fast(address uintptr) (uint32, error) {
	var value uint32
	var bytesRead uintptr

	err := windows.ReadProcessMemory(
		c.processHandle,
		address,
		(*byte)(unsafe.Pointer(&value)),
		4,
		&bytesRead,
	)
	
	if err != nil || bytesRead != 4 {
		return 0, fmt.Errorf("error leyendo memoria en 0x%X: %v", address, err)
	}

	return value, nil
}

// GetHP returns cached HP values (ultra-fast)
func (c *OptimizedClient) GetHP() (current, max uint32, valid bool) {
	c.cacheMutex.RLock()
	defer c.cacheMutex.RUnlock()
	
	if !c.cache.Valid {
		return 0, 0, false
	}
	
	// Check if cache is fresh (less than 100ms old)
	if time.Since(c.cache.LastUpdate) > 100*time.Millisecond {
		return 0, 0, false
	}
	
	return c.cache.HP, c.cache.MaxHP, true
}

// GetSP returns cached SP values (ultra-fast)
func (c *OptimizedClient) GetSP() (current, max uint32, valid bool) {
	c.cacheMutex.RLock()
	defer c.cacheMutex.RUnlock()
	
	if !c.cache.Valid {
		return 0, 0, false
	}
	
	// Check if cache is fresh (less than 100ms old)
	if time.Since(c.cache.LastUpdate) > 100*time.Millisecond {
		return 0, 0, false
	}
	
	return c.cache.SP, c.cache.MaxSP, true
}

// ReadHP provides backward compatibility
func (c *OptimizedClient) ReadHP() (current, max uint32, err error) {
	if current, max, valid := c.GetHP(); valid {
		return current, max, nil
	}
	
	// Fallback to direct read
	current, err = c.readUint32Fast(c.hpBaseAddr)
	if err != nil {
		return 0, 0, err
	}
	
	max, err = c.readUint32Fast(c.hpBaseAddr + 4)
	if err != nil {
		return 0, 0, err
	}
	
	return current, max, nil
}

// ReadSP provides backward compatibility
func (c *OptimizedClient) ReadSP() (current, max uint32, err error) {
	if current, max, valid := c.GetSP(); valid {
		return current, max, nil
	}
	
	// Fallback to direct read
	current, err = c.readUint32Fast(c.spBaseAddr)
	if err != nil {
		return 0, 0, err
	}
	
	max, err = c.readUint32Fast(c.spBaseAddr + 4)
	if err != nil {
		return 0, 0, err
	}
	
	return current, max, nil
}

// Close shuts down the optimized client
func (c *OptimizedClient) Close() error {
	c.cancel()
	c.refreshTicker.Stop()
	return windows.CloseHandle(c.processHandle)
}

// GetCacheStats returns cache performance statistics
func (c *OptimizedClient) GetCacheStats() (age time.Duration, valid bool) {
	c.cacheMutex.RLock()
	defer c.cacheMutex.RUnlock()
	
	return time.Since(c.cache.LastUpdate), c.cache.Valid
}