# DTAP - Ragnarok Online Assistant

Un asistente moderno para Ragnarok Online con **sistema de predicción inteligente**, **ATK-DEF Mode** y interfaz gráfica intuitiva.

## 🙏 Créditos y Proyectos Base

Este proyecto ha tomado inspiración y código base de los siguientes proyectos de la comunidad:

- **[k1ngJ/dtAP](https://github.com/k1ngJ/dtAP)** - Proyecto original en C# que sirvió como base fundamental
- **[SmookyzX/SmookyzAP](https://github.com/SmookyzX/SmookyzAP)** - Autopotion avanzado con High Ping Mode
- **[4RTools/4RTools](https://github.com/4RTools/4RTools)** - Suite completa de herramientas para Ragnarok Online

Agradecemos a todos los desarrolladores de la comunidad por su trabajo y contribuciones al ecosistema de herramientas para Ragnarok Online.

## 🚀 Características

- **🧠 Predicción Inteligente**: Anticipa el daño 200ms antes usando regresión lineal
- **⚡ AutoPot Ultra Rápido**: Autopot con predicción de daño (15ms interval)
- **🎮 Spammer Optimizado**: Spammer de alta velocidad (1ms interval)
- **🔄 AutoBuff Inteligente**: Sistema de buffs automático (250ms interval)
- **⚔️ ATK-DEF Mode**: Cambio automático entre equipos de ataque y defensa
- **🖥️ Interfaz Moderna**: GUI React/TypeScript con configuración visual
- **📊 Monitoreo en Vivo**: Estadísticas actualizadas a 60 FPS
- **🎨 Temas Personalizables**: Múltiples temas con transparencia
- **⚙️ Configuración en Tiempo Real**: Ajustes sin reiniciar
- **🔧 Scheduler Ultra-Optimizado**: Sistema de tareas con prioridades

## 🖼️ Interfaz de Usuario

### Modo Oscuro
![DTAP Dark Mode](dtap-gui/images/dark-mode-1.png)
![DTAP Dark Mode](dtap-gui/images/dark-mode-2.png)

### Modo Claro
![DTAP Light Mode](dtap-gui/images/light-mode-1.png)
![DTAP Light Mode](dtap-gui/images/light-mode-2.png)

## 📁 Estructura del Proyecto

```
dtap-gui/
├── frontend/             # Interfaz React/TypeScript
│   ├── components/       # Componentes de la UI
│   │   ├── ATKDEF.tsx    # Configuración ATK-DEF Mode
│   │   ├── Cheats.tsx    # Configuración Spammer
│   │   └── Dashboard.tsx # Panel principal
│   ├── contexts/         # Contextos de React
│   ├── styles/           # Estilos CSS
│   └── pages/            # Páginas principales
├── internal/             # Lógica del backend
│   ├── normal/           # Servicios principales
│   │   ├── autopot/      # Sistema de autopot
│   │   ├── spammer/      # Sistema de spammer
│   │   ├── autobuff/     # Sistema de autobuff
│   │   └── atkdef/       # Sistema ATK-DEF Mode
│   ├── scheduler/        # Scheduler ultra-optimizado
│   ├── memory/           # Acceso a memoria
│   ├── prediction/       # Sistema de predicción
│   ├── input/            # Sistema de input
│   └── services/         # Gestión de servicios
├── pkg/                  # Paquetes compartidos
│   ├── config/          # Configuración
│   ├── constants/       # Constantes del sistema
│   └── logger/           # Sistema de logging
├── build/                # Archivos de build
└── main.go               # Punto de entrada
```

## 🛠️ Instalación

### Prerrequisitos
- **Go 1.21+**
- **Node.js 16+**
- **Wails v2**
- **Windows 10/11**

### Instalación
```bash
# Instalar Wails
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Clonar el repositorio
git clone <tu-repo>
cd dtap-gui

# Instalar dependencias
npm install

# Compilar
wails build

# Ejecutar en desarrollo
wails dev
```


## 🎮 Uso

1. **Ejecutar la aplicación**
   ```bash
   wails dev
   ```

2. **Conectar a Ragnarok Online**
   - Asegúrate de que RO esté ejecutándose
   - Haz clic en "Conectar" en la GUI

3. **Configurar servicios**
   - **AutoPot**: Configura teclas y porcentajes de HP/SP
   - **Spammer**: Configura teclas y delay (1ms ultra-fast)
   - **AutoBuff**: Configura buffs automáticos
   - **ATK-DEF Mode**: Configura cambio automático de equipos

4. **Iniciar servicios**
   - Haz clic en "Iniciar" para activar los servicios
   - Los servicios se ejecutarán automáticamente con scheduler optimizado

## ⚙️ Configuración

La configuración se realiza directamente desde la interfaz gráfica:

- **AutoPot**: Configura teclas de HP/SP y porcentajes
- **Spammer**: Configura teclas, delay (1ms ultra-fast) y modo
- **AutoBuff**: Configura buffs y delays
- **ATK-DEF Mode**: Configura tecla de activación, equipos ATK/DEF y delays
- **Temas**: Cambia entre temas claro/oscuro
- **Transparencia**: Ajusta la transparencia de la ventana

Todos los cambios se aplican en tiempo real sin necesidad de reiniciar.


## 📊 Características Técnicas

- **Predicción de Daño**: Anticipa el daño 200ms antes usando regresión lineal
- **Alto Rendimiento**: AutoPot a 15ms, Spammer a 1ms, AutoBuff a 250ms, ATK-DEF a 1ms
- **Scheduler Ultra-Optimizado**: Sistema de tareas con prioridades (Critical > High > Normal > Low)
- **Constantes Configurables**: Cero configuraciones hardcodeadas, todo configurable
- **Optimizaciones de Sistema**: timeBeginPeriod(1) para máxima precisión
- **Interfaz Fluida**: Actualizaciones a 60 FPS sin lag
- **Acceso a Memoria**: Lectura directa de HP/SP del proceso de RO
- **ATK-DEF Mode**: Cambio automático entre equipos con tecla de activación

## 🆕 Últimas Mejoras

### ⚡ Optimizaciones del Scheduler
- **Sistema de Prioridades**: Critical (AutoPot) > High (Spammer/ATK-DEF) > Normal (AutoBuff) > Low (Mantenimiento)
- **Constantes Configurables**: Eliminadas todas las configuraciones hardcodeadas
- **Delays Optimizados**: Spammer y ATK-DEF a 1ms ultra-fast
- **Logging Mejorado**: Logs contextuales con emojis para fácil identificación

### ⚔️ ATK-DEF Mode
- **Tecla de Activación**: Configurable (F1-F12, letras, números)
- **Equipos Dinámicos**: Configuración de equipos ATK y DEF
- **Spam Integrado**: Combina spammer con cambio de equipos
- **Delays Configurables**: Switch delay y spam delay independientes

### 🔧 Arquitectura Mejorada
- **Constantes Centralizadas**: `pkg/constants/timing.go`
- **Scheduler Optimizado**: Base interval 1ms para máxima precisión
- **Logging Balanceado**: 108 logs distribuidos en 10 archivos
- **Código Mantenible**: Estructura escalable y documentada

## 📈 Rendimiento

### ⚡ Velocidades de Ejecución
- **AutoPot**: 15ms (66 Hz) - Crítico para supervivencia
- **Spammer**: 1ms (1000 Hz) - Ultra-fast para combate
- **ATK-DEF**: 1ms (1000 Hz) - Ultra-fast para cambio de equipos
- **AutoBuff**: 250ms (4 Hz) - Optimizado para soporte

### 🎯 Prioridades del Scheduler
1. **Critical**: AutoPot (supervivencia)
2. **High**: Spammer, ATK-DEF (combate)
3. **Normal**: AutoBuff (soporte)
4. **Low**: Mantenimiento

### 📊 Métricas en Tiempo Real
- **Uptime**: Tiempo de ejecución
- **Loops**: Contadores de ejecución por servicio
- **Performance**: Estadísticas de rendimiento
- **Memory**: Uso de memoria del sistema

## 🚨 Requisitos

- **Windows 10/11** (64-bit)
- **Ragnarok Online** ejecutándose
- **Permisos de administrador** (para acceso a memoria)
- **Go 1.21+** y **Node.js 16+** (para compilación)

## 📄 Licencia

Este proyecto está bajo la Licencia MIT.

---

**DTAP** - *Ragnarok Online Assistant* 🚀