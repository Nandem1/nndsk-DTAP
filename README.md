# DTAP-Go - Ragnarok Online Auto Tool

Una herramienta de automatización para Ragnarok Online escrita en Go, con funcionalidades de auto pot, auto buff y spammer.

## 🚀 Características

- **Auto Pot** - Uso automático de pociones de HP/SP
- **Auto Buff** - Aplicación automática de buffs
- **Spammer** - Repetición automática de teclas
- **Address Finder** - Herramienta para encontrar direcciones de memoria
- **Pattern Scanning** - Búsqueda dinámica de direcciones
- **Configuración JSON** - Sistema de configuración flexible
- **Logging estructurado** - Sistema de logs profesional

## 📁 Estructura del Proyecto

```
dtap-go/
├── cmd/
│   └── find_addresses/     # Address finder tool
├── internal/
│   ├── memory/             # Memory operations
│   ├── process/            # Process management
│   ├── input/              # Input automation
│   ├── scanner/            # Pattern scanning
│   └── finder/             # Address finding
├── pkg/
│   ├── constants/          # All constants
│   ├── config/             # Configuration system
│   └── logger/             # Logging system
├── main.go                 # Main application
└── go.mod
```

## 🔧 Build Commands

### Build Principal (Auto Tool)
```bash
go build -o dtap.exe .
```

### Build Address Finder
```bash
go build -o find_addresses.exe ./cmd/find_addresses
```

## ⚙️ Configuración

### Direcciones de Memoria
Edita `pkg/constants/constants.go` para ajustar las direcciones de memoria:

```go
const (
    HPBaseAddress    = 0x10DCE10 // HP Current
    HPMaxAddress     = 0x10DCE14 // HP Maximum
    SPBaseAddress    = 0x10DCE18 // SP Current
    SPMaxAddress     = 0x10DCE1C // SP Maximum
    StatusBufferAddr = 0x10DD284 // Status buffer
)
```

### Valores para Address Finder
Edita las constantes en `pkg/constants/constants.go`:

```go
const (
    FinderHPCurrent = 6914  // Tu HP actual
    FinderHPMax     = 6914  // Tu HP máximo
    FinderSPCurrent = 2511  // Tu SP actual
    FinderSPMax     = 2511  // Tu SP máximo
)
```

### Configuración de la Aplicación
El archivo de configuración se crea automáticamente en el mismo directorio del programa como `dtap-config.json`:

```json
{
  "autopot": {
    "hp_pot_key": "F1",
    "sp_pot_key": "F2",
    "hp_percent": 50,
    "sp_percent": 30,
    "check_delay_ms": 100,
    "enabled": true
  },
  "autobuff": {
    "enabled": false,
    "delay_ms": 5000
  },
  "spammer": {
    "enabled": false,
    "keys": [],
    "delay_ms": 50
  },
  "ui": {
    "show_debug": false,
    "log_level": "info"
  }
}
```

## 🎯 Uso

### 1. Encontrar Direcciones de Memoria
```bash
# Ejecutar como administrador
.\find_addresses.exe
```

### 2. Ejecutar Auto Tool
```bash
# Ejecutar como administrador
.\dtap.exe
```

## 🔍 Proceso de Address Finding

1. **Asegúrate** de que tu personaje esté en standby (parado)
2. **Anota** los valores exactos de HP/SP del juego
3. **Edita** las constantes en `pkg/constants/constants.go`
4. **Rebuild** el address finder: `go build -o find_addresses.exe ./cmd/find_addresses`
5. **Ejecuta** como administrador: `.\find_addresses.exe`
6. **Copia** las direcciones encontradas a las constantes principales
7. **Rebuild** el programa principal: `go build -o dtap.exe .`

## 🛠️ Desarrollo

### Requisitos
- Go 1.21+
- Windows (para APIs de Windows)
- Permisos de administrador

### Dependencias
```bash
go mod tidy
```

### Estructura de Constantes
- `constants.go` - Todas las constantes del proyecto
- `config.go` - Sistema de configuración
- `logger.go` - Sistema de logging

## ⚙️ Configuración Rápida

### Cambiar Teclas y Delay del Auto Pot:

1. **Ejecuta** `.\dtap.exe` una vez (crea `dtap-config.json`)
2. **Edita** `dtap-config.json` en el mismo directorio:
   ```json
   {
     "autopot": {
       "hp_pot_key": "F1",        // Cambia por tu tecla de HP
       "sp_pot_key": "F2",        // Cambia por tu tecla de SP
       "hp_percent": 50,          // Porcentaje para usar HP pot
       "sp_percent": 30,          // Porcentaje para usar SP pot
       "check_delay_ms": 100,     // Delay en milisegundos
       "enabled": true
     }
   }
   ```
3. **Ejecuta** `.\dtap.exe` de nuevo (lee la nueva configuración)

## 📝 Notas Importantes

- **Siempre ejecutar como administrador** para acceso a memoria
- **Rebuild después de cambiar constantes** para que los cambios se apliquen
- **El address finder es específico** para cada cliente de RO
- **Las direcciones pueden cambiar** con updates del juego
- **El archivo de configuración se crea en el mismo directorio** del programa

## 🚨 Disclaimer

Esta herramienta es para uso educativo y personal. Úsala bajo tu propia responsabilidad y respetando los términos de servicio del juego.

## 📄 Licencia

Proyecto personal - Uso educativo únicamente.
