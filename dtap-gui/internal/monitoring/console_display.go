package monitoring

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

// ConsoleDisplay maneja la visualización en consola
type ConsoleDisplay struct {
	lastUpdate time.Time
	updateFreq time.Duration
}

// NewConsoleDisplay crea un nuevo display de consola
func NewConsoleDisplay(updateFreq time.Duration) *ConsoleDisplay {
	return &ConsoleDisplay{
		updateFreq: updateFreq,
	}
}

// DisplayStats muestra las estadísticas en consola de forma optimizada
func (cd *ConsoleDisplay) DisplayStats(
	hp, maxHP uint32, hpPercent int, sp, maxSP uint32, spPercent int,
	serviceStats map[string]interface{},
	systemStats *SystemStats,
) {
	now := time.Now()

	// Solo actualizar si ha pasado suficiente tiempo
	if now.Sub(cd.lastUpdate) < cd.updateFreq {
		return
	}

	cd.lastUpdate = now

	// Limpiar línea anterior - compatible con Windows
	fmt.Print("\r") // Solo return, sin códigos ANSI

	// Construir línea de estadísticas
	var parts []string

	// HP/SP Status
	parts = append(parts, fmt.Sprintf("HP: %d/%d (%d%%)", hp, maxHP, hpPercent))
	parts = append(parts, fmt.Sprintf("SP: %d/%d (%d%%)", sp, maxSP, spPercent))

	// Service Status
	if uptime, ok := serviceStats["uptime"].(time.Duration); ok {
		parts = append(parts, fmt.Sprintf("Uptime: %s", formatDuration(uptime)))
	}

	// System Stats
	if systemStats != nil {
		parts = append(parts, fmt.Sprintf("RAM: %s/%s (%s)",
			FormatBytes(systemStats.MemoryUsage),
			FormatBytes(systemStats.MemoryTotal),
			FormatPercent(systemStats.MemoryPercent)))

		parts = append(parts, fmt.Sprintf("GoRAM: %s", FormatBytes(systemStats.GoMemoryUsage)))
		parts = append(parts, fmt.Sprintf("Goroutines: %d", systemStats.GoroutineCount))
	}

	// Service Stats
	if spammerLoops, ok := serviceStats["spammer_loops"].(int64); ok {
		parts = append(parts, fmt.Sprintf("Spam: %d", spammerLoops))
	}

	if autopotPots, ok := serviceStats["autopot_potions"].(int64); ok {
		parts = append(parts, fmt.Sprintf("Pots: %d", autopotPots))
	}

	if autobuffBuffs, ok := serviceStats["autobuff_buffs"].(int64); ok {
		parts = append(parts, fmt.Sprintf("Buffs: %d", autobuffBuffs))
	}

	// Construir línea completa
	line := strings.Join(parts, " | ")

	// Asegurar que la línea tenga al menos 120 caracteres para limpiar
	if len(line) < 120 {
		line += strings.Repeat(" ", 120-len(line))
	}

	// Imprimir línea completa
	fmt.Print(line)

	// Forzar flush para asegurar que se muestre
	os.Stdout.Sync()
}

// DisplayDetailedStats muestra estadísticas detalladas
func (cd *ConsoleDisplay) DisplayDetailedStats(
	hp, maxHP uint32, hpPercent int, sp, maxSP uint32, spPercent int,
	serviceStats map[string]interface{},
	systemStats *SystemStats,
) {
	// Limpiar pantalla
	fmt.Print("\033[2J\033[H") // Clear screen and move cursor to top

	fmt.Println("=== DTAP-GO STATUS ===")
	fmt.Printf("Timestamp: %s\n", time.Now().Format("15:04:05"))
	fmt.Println()

	// HP/SP Status
	fmt.Println("--- Character Status ---")
	fmt.Printf("HP: %d/%d (%d%%)\n", hp, maxHP, hpPercent)
	fmt.Printf("SP: %d/%d (%d%%)\n", sp, maxSP, spPercent)
	fmt.Println()

	// Service Status
	fmt.Println("--- Service Status ---")
	if uptime, ok := serviceStats["uptime"].(time.Duration); ok {
		fmt.Printf("Uptime: %s\n", formatDuration(uptime))
	}
	if started, ok := serviceStats["started"].(bool); ok {
		fmt.Printf("Running: %t\n", started)
	}
	fmt.Println()

	// Service Statistics
	fmt.Println("--- Service Statistics ---")
	if spammerLoops, ok := serviceStats["spammer_loops"].(int64); ok {
		fmt.Printf("Spammer Loops: %d\n", spammerLoops)
	}
	if autopotLoops, ok := serviceStats["autopot_loops"].(int64); ok {
		fmt.Printf("AutoPot Loops: %d\n", autopotLoops)
	}
	if autobuffLoops, ok := serviceStats["autobuff_loops"].(int64); ok {
		fmt.Printf("AutoBuff Loops: %d\n", autobuffLoops)
	}
	if autopotPots, ok := serviceStats["autopot_potions"].(int64); ok {
		fmt.Printf("Potions Used: %d\n", autopotPots)
	}
	if autobuffBuffs, ok := serviceStats["autobuff_buffs"].(int64); ok {
		fmt.Printf("Buffs Used: %d\n", autobuffBuffs)
	}
	fmt.Println()

	// System Statistics
	if systemStats != nil {
		fmt.Println("--- System Statistics ---")
		fmt.Printf("System RAM: %s/%s (%s)\n",
			FormatBytes(systemStats.MemoryUsage),
			FormatBytes(systemStats.MemoryTotal),
			FormatPercent(systemStats.MemoryPercent))
		fmt.Printf("Go RAM: %s/%s (%s)\n",
			FormatBytes(systemStats.GoMemoryUsage),
			FormatBytes(systemStats.GoMemoryTotal),
			FormatPercent(systemStats.GoMemoryPercent))
		fmt.Printf("Process RAM: %s\n", FormatBytes(systemStats.ProcessMemoryUsage))
		fmt.Printf("Goroutines: %d\n", systemStats.GoroutineCount)
		fmt.Printf("CPU Usage: %s\n", FormatPercent(systemStats.ProcessCPUUsage))
		fmt.Println()
	}

	// Go Runtime Stats
	fmt.Println("--- Go Runtime ---")
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("GC Cycles: %d\n", m.NumGC)
	fmt.Printf("GC Pause: %s\n", time.Duration(m.PauseTotalNs).String())
	fmt.Printf("Heap Objects: %d\n", m.HeapObjects)
	fmt.Printf("Stack Inuse: %s\n", FormatBytes(m.StackInuse))
}

// formatDuration formatea una duración de forma legible
func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%.0fs", d.Seconds())
	} else if d < time.Hour {
		return fmt.Sprintf("%.0fm", d.Minutes())
	} else {
		hours := int(d.Hours())
		minutes := int(d.Minutes()) % 60
		return fmt.Sprintf("%dh%dm", hours, minutes)
	}
}
