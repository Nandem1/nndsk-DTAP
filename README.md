# DTAP - Ragnarok Online Assistant

Un asistente avanzado para Ragnarok Online con dos versiones especializadas: **Normal** y **WoE Optimizada**.

## 🚀 Características

### Versión Normal
- **Autopot**: Uso automático de pociones HP/SP
- **Spammer**: Spam de teclas con click izquierdo automático
- **Skill Cancel**: Cancelación de habilidades
- **Double Click**: Doble click para cancelar
- **Interfaz simple**: Fácil de usar y configurar

### Versión WoE (War of Emperium)
- **Ultra-fast Autopot**: Autopot optimizado para competencias
- **Smart Spammer**: Spammer de 1ms sin rate limiting
- **Performance Monitoring**: Monitoreo en tiempo real
- **Memory Optimization**: Cache de memoria para máxima velocidad
- **Burst Mode**: Modo de ráfaga para situaciones críticas
- **Predictive Mode**: Predicción de daño para autopot inteligente

## 📁 Estructura del Proyecto

```
dtap-go/
├── cmd/                    # Puntos de entrada
│   ├── dtap/              # Versión normal
│   │   └── main.go
│   └── dtap-woe/          # Versión WoE
│       └── main.go
├── internal/              # Código interno
│   ├── normal/           # Código para versión normal
│   │   ├── autopot/
│   │   ├── spammer/
│   │   └── monitoring/
│   ├── woe/              # Código para versión WoE
│   │   ├── autopot/
│   │   ├── spammer/
│   │   └── monitoring/
│   ├── input/            # Sistema de input compartido
│   ├── memory/           # Sistema de memoria compartido
│   ├── process/          # Gestión de procesos
│   ├── scanner/          # Escaneo de patrones
│   └── finder/           # Búsqueda de direcciones
├── pkg/                  # Código compartido
│   ├── config/           # Configuración
│   ├── constants/        # Constantes
│   └── logger/           # Sistema de logging
├── configs/              # Archivos de configuración
│   ├── dtap-config.json
│   └── dtap-woe-config.json
└── cmd/                  # Herramientas de línea de comandos
    ├── analyze_character/
    ├── find_addresses/
    ├── find_biochemist/
    ├── find_character/
    ├── inspect_character/
    └── parse_character/
```

## 🛠️ Compilación

### Usando Makefile (Recomendado)

```bash
# Compilar ambas versiones
make build

# Compilar solo versión normal
make build-normal

# Compilar solo versión WoE
make build-woe

# Limpiar archivos compilados
make clean

# Ver todas las opciones
make help
```

### Compilación Manual

```bash
# Versión normal
go build -o dtap.exe cmd/dtap/main.go

# Versión WoE
go build -o dtap-woe.exe cmd/dtap-woe/main.go
```

## 🎮 Uso

### Versión Normal
```bash
./dtap.exe
```

### Versión WoE
```bash
./dtap-woe.exe
```

## ⚙️ Configuración

### Versión Normal
- Archivo: `configs/dtap-config.json`
- Configuración básica de autopot y spammer

### Versión WoE
- Archivo: `configs/dtap-woe-config.json`
- Configuración avanzada con optimizaciones para WoE

## 🔧 Herramientas de Línea de Comandos

```bash
# Analizar personaje
./analyze_character.exe

# Encontrar direcciones de memoria
./find_addresses.exe

# Encontrar bioquímico
./find_biochemist.exe

# Encontrar personaje
./find_character.exe

# Inspeccionar personaje
./inspect_character.exe

# Parsear personaje
./parse_character.exe
```

## 📊 Rendimiento

### Versión Normal
- **Autopot**: ~100ms de latencia
- **Spammer**: Delay configurable (default: 25ms)
- **Memoria**: Lectura directa sin cache

### Versión WoE
- **Autopot**: ~50µs de latencia
- **Spammer**: 1ms loop (122K spams/segundo)
- **Memoria**: Cache optimizado con 95% hit rate
- **Monitoring**: Reportes de performance en tiempo real

## 🚨 Requisitos

- Windows 10/11
- Ragnarok Online ejecutándose
- Permisos de administrador (para acceso a memoria)

## 📝 Notas

- La versión WoE está optimizada para competencias y puede usar más recursos
- Ambas versiones son compatibles con la misma configuración base
- El sistema de logging permite debuggear problemas fácilmente
- Los archivos de configuración se crean automáticamente si no existen

## 🤝 Contribución

1. Fork el proyecto
2. Crea una rama para tu feature (`git checkout -b feature/AmazingFeature`)
3. Commit tus cambios (`git commit -m 'Add some AmazingFeature'`)
4. Push a la rama (`git push origin feature/AmazingFeature`)
5. Abre un Pull Request

## 📄 Licencia

Este proyecto está bajo la Licencia MIT - ver el archivo [LICENSE](LICENSE) para detalles.