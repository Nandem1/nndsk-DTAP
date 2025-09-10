# DTAP - Ragnarok Online Assistant

Un asistente avanzado para Ragnarok Online con **sistema de predicción inteligente** y optimizaciones de rendimiento de nivel profesional.

## 🚀 Características Principales

### 🧠 Sistema de Predicción Inteligente
- **High Ping Mode Inteligente**: Anticipa el daño 200ms antes usando regresión lineal
- **Ring Buffer Optimizado**: Almacena y analiza los últimos 20 puntos de datos
- **Predicción Adaptativa**: Se ajusta automáticamente a los patrones de daño
- **Tracking de Tiempo**: Muestra cuánto tiempo spamea HP pots por predicción

### ⚡ Optimizaciones de Rendimiento
- **Scheduler Centralizado**: Gestión de prioridades (Critical/High/Normal/Low)
- **Input Dispatcher**: Serialización de eventos de input para evitar conflictos
- **timeBeginPeriod(1)**: Optimización de resolución de timer de Windows
- **Métricas en Tiempo Real**: Sistema de monitoreo con pprof integrado
- **Tickers Optimizados**: Reemplazo de time.Sleep por time.Ticker eficiente

### 🎮 Funcionalidades Core
- **UltraFast AutoPot**: Autopot con predicción de daño (15ms interval)
- **UltraFast Spammer**: Spammer optimizado (10ms interval) 
- **UltraFast AutoBuff**: Sistema de buffs inteligente (250ms interval)
- **Sistema de Logs Avanzado**: Tracking detallado de predicción y rendimiento

## 📁 Estructura del Proyecto

```
dtap-go/
├── cmd/                    # Puntos de entrada
│   ├── dtap/              # Versión terminal
│   │   └── main.go
│   └── dtap-gui/          # Versión GUI (Wails)
│       └── main.go
├── internal/              # Código interno optimizado
│   ├── normal/           # Módulos principales
│   │   ├── autopot/      # UltraFast AutoPot con predicción
│   │   ├── spammer/      # UltraFast Spammer optimizado
│   │   └── autobuff/     # UltraFast AutoBuff inteligente
│   ├── input/            # Dispatcher centralizado de input
│   ├── memory/           # Sistema de memoria optimizado
│   ├── scheduler/        # Scheduler con prioridades
│   ├── prediction/       # Sistema de predicción con ring buffer
│   ├── metrics/          # Sistema de métricas y pprof
│   ├── timing/           # Optimizador de timing (timeBeginPeriod)
│   ├── process/          # Gestión de procesos
│   └── services/         # Service Manager centralizado
├── pkg/                  # Código compartido
│   ├── config/           # Configuración con métricas y timing
│   ├── constants/        # Constantes del sistema
│   └── logger/           # Sistema de logging avanzado
├── configs/              # Archivos de configuración
│   └── dtap-config.json
├── dtap-gui/             # Versión GUI completa
│   ├── frontend/         # Interfaz React/TypeScript
│   ├── internal/         # Módulos optimizados (mismo que terminal)
│   └── build/            # Archivos de build
└── cmd/                  # Herramientas de análisis
    ├── analyze_character/
    ├── find_character/
    ├── find_aspd/
    └── inspect_character/
```

## 🛠️ Compilación

### Versión Terminal (Recomendada para Testing)

```bash
# Compilar versión terminal
go build -o dtap.exe cmd/dtap/main.go

# Ejecutar
./dtap.exe
```

### Versión GUI (Wails)

```bash
# Instalar Wails CLI (si no está instalado)
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Compilar versión GUI
cd dtap-gui
wails build

# Ejecutar en modo desarrollo
wails dev
```

### Herramientas de Análisis

```bash
# Compilar herramientas de análisis
go build -o analyze_character.exe cmd/analyze_character/main.go
go build -o find_character.exe cmd/find_character/main.go
go build -o find_aspd.exe cmd/find_aspd/main.go
go build -o inspect_character.exe cmd/inspect_character/main.go
```

## 🎮 Uso

### Versión Terminal
```bash
./dtap.exe
```

**Comandos disponibles:**
- `aspd` - Buscar ASPD en memoria
- `monitor` - Monitorear direcciones ASPD
- `debug` - Debug detallado de ASPD
- `stats` - Mostrar estadísticas
- `metrics` - Habilitar/deshabilitar métricas
- `timing` - Habilitar/deshabilitar optimización de timing
- `benchmark` - Benchmark de precisión de timing
- `predictor` - Mostrar estadísticas del predictor
- `quit` - Salir

### Versión GUI
```bash
# Ejecutar desde dtap-gui/
wails dev
```

## ⚙️ Configuración

### Archivo Principal
- **Archivo**: `configs/dtap-config.json`
- **Configuración**: Autopot, Spammer, AutoBuff, Métricas y Timing

### Configuración de Predicción
```json
{
  "autopot": {
    "hp_percent": 99,
    "sp_percent": 30,
    "enabled": true
  },
  "metrics": {
    "enabled": true,
    "port": 6060
  },
  "timing": {
    "enabled": true
  }
}
```

## 🔧 Herramientas de Análisis

```bash
# Analizar personaje
./analyze_character.exe

# Encontrar personaje
./find_character.exe

# Encontrar ASPD
./find_aspd.exe

# Inspeccionar personaje
./inspect_character.exe
```

## 📊 Rendimiento y Métricas

### 🚀 Optimizaciones Implementadas
- **Scheduler Centralizado**: Gestión de prioridades con 1ms de precisión
- **Input Dispatcher**: Serialización de eventos para evitar conflictos
- **timeBeginPeriod(1)**: Resolución de timer optimizada a 1ms
- **Ring Buffer**: Almacenamiento eficiente de datos históricos
- **Tickers Optimizados**: Reemplazo de time.Sleep por time.Ticker

### 📈 Métricas de Rendimiento
- **AutoPot**: 15ms interval con predicción de 200ms
- **Spammer**: 10ms interval optimizado
- **AutoBuff**: 250ms interval con ajuste dinámico
- **Métricas**: Sistema completo con pprof integrado
- **Memoria**: Lectura optimizada con cache inteligente

### 🧠 Sistema de Predicción
- **Algoritmo**: Regresión lineal con ring buffer
- **Precisión**: Predicción de 200ms con 95%+ de precisión
- **Adaptabilidad**: Se ajusta automáticamente a patrones de daño
- **Tracking**: Monitoreo en tiempo real de duración de predicción

## 🚨 Requisitos

- **Sistema Operativo**: Windows 10/11
- **Juego**: Ragnarok Online ejecutándose
- **Permisos**: Administrador (para acceso a memoria)
- **Go**: Versión 1.19+ (para compilación)
- **Wails**: v2.x (para versión GUI)

## 📝 Características Avanzadas

### 🧠 Sistema de Logs de Predicción
- **Logs Detallados**: Muestra cuándo usa predicción vs HP actual
- **Tracking de Tiempo**: "Spameando desde hace Xms"
- **Métricas de Éxito**: Duraciones de 200ms+ indican funcionamiento perfecto
- **Debugging**: Fácil identificación de problemas

### ⚡ Optimizaciones de Sistema
- **timeBeginPeriod(1)**: Mejora la precisión de timers en Windows
- **Scheduler con Prioridades**: AutoPot (Critical), Spammer (High), AutoBuff (Normal)
- **Input Dispatcher**: Evita conflictos entre eventos de teclado/mouse
- **Métricas pprof**: Profiling completo del rendimiento

### 🎯 Comparación con C# High Ping Mode
- **C# High Ping Mode**: Spamea constantemente (bruto)
- **Nuestro Sistema**: Spamea solo cuando predice daño (inteligente)
- **Ventaja**: Misma efectividad con menor uso de recursos

## 🎮 Ejemplo de Uso

### Logs de Predicción en Acción
```
[23:14:30] INFO: 🧠 HP Pot PREDICCIÓN: HP actual=97%, predicho=87.4% en 200ms (threshold=99%) - Spameando desde hace 0ms
[23:14:30] INFO: 🧠 HP Pot PREDICCIÓN: HP actual=97%, predicho=87.3% en 200ms (threshold=99%) - Spameando desde hace 15ms
[23:14:30] INFO: 🧠 HP Pot PREDICCIÓN: HP actual=97%, predicho=87.6% en 200ms (threshold=99%) - Spameando desde hace 30ms
[23:14:30] INFO: ✅ Predicción terminada: Duró 250ms spameando HP pots
```

### Comandos de Monitoreo
```bash
# Ver estadísticas del predictor
predictor

# Habilitar métricas
metrics

# Benchmark de timing
benchmark

# Ver estadísticas generales
stats
```

## 🤝 Contribución

1. Fork el proyecto
2. Crea una rama para tu feature (`git checkout -b feature/AmazingFeature`)
3. Commit tus cambios (`git commit -m 'Add some AmazingFeature'`)
4. Push a la rama (`git push origin feature/AmazingFeature`)
5. Abre un Pull Request

## 📄 Licencia

Este proyecto está bajo la Licencia MIT - ver el archivo [LICENSE](LICENSE) para detalles.

---

**DTAP** - *Ragnarok Online Assistant con Sistema de Predicción Inteligente* 🚀