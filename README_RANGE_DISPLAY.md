# 🎯 DTAP Range Display - Sistema de Visualización de Rango

## 📋 Descripción

El sistema de **Range Display** muestra un **círculo visual** alrededor del personaje que representa su **rango de ataque**, proporcionando información táctica crucial para el posicionamiento y combate.

## ✨ Características

### 🎨 **Visualización del Rango**
- **Círculo dinámico** que muestra el rango de ataque del personaje
- **Colores personalizables** en formato hexadecimal
- **Transparencia ajustable** (0.0 - 1.0)
- **Grosor de línea configurable**
- **Círculo suave** con segmentos ajustables

### ⚙️ **Configuración**

#### Versión Normal
```json
{
  "range_display": {
    "enabled": true,
    "color": "0x7F00FF00",
    "opacity": 0.5,
    "thickness": 2.0,
    "segments": 32
  }
}
```

#### Versión WoE (Optimizada)
```json
{
  "range_display": {
    "enabled": true,
    "color": "0x7F00FF00",
    "opacity": 0.7,
    "thickness": 3.0,
    "segments": 64,
    "woe_optimized": true,
    "show_only_in_combat": false
  }
}
```

## 🎮 **Parámetros de Configuración**

### **Parámetros Básicos**
- `enabled`: Habilita/deshabilita el range display
- `color`: Color del círculo en formato "0xAARRGGBB"
- `opacity`: Transparencia (0.0 = invisible, 1.0 = opaco)
- `thickness`: Grosor de la línea del círculo
- `segments`: Número de segmentos (más = círculo más suave)

### **Parámetros WoE**
- `woe_optimized`: Optimizaciones específicas para WoE
- `show_only_in_combat`: Mostrar solo durante combate

## 🎯 **Colores Recomendados**

```json
"0x7F00FF00"  // Verde semi-transparente (recomendado)
"0x7F0000FF"  // Azul semi-transparente
"0x7FFF0000"  // Rojo semi-transparente
"0x7FFFFFFF"  // Blanco semi-transparente
"0x7FFFFF00"  // Amarillo semi-transparente
"0x7FFF00FF"  // Magenta semi-transparente
```

## 🚀 **Ventajas Tácticas**

### ⚔️ **En Combate Normal**
- **Posicionamiento preciso** - Saber exactamente dónde puedes atacar
- **Evitar errores** - No intentar atacar fuera de rango
- **Optimización de movimiento** - Moverte a la distancia perfecta

### 🏰 **En WoE (War of Emperium)**
- **Ventaja competitiva** - Información visual en tiempo real
- **Coordinación de equipo** - Todos ven el mismo rango
- **Estrategia avanzada** - Posicionamiento táctico perfecto
- **Timing de habilidades** - Usar skills en el momento exacto

## 🔧 **Implementación Técnica**

### **Arquitectura**
```
pkg/visual/
├── range_display.go        # Sistema de visualización de rango
├── skill_marker.go         # Sistema de marcado de habilidades
├── overlay.go              # Manager del overlay
└── config_adapter.go       # Adaptador de configuraciones
```

### **Características Técnicas**
- **Tiempo real** - Actualización cada 50ms (20 FPS)
- **Dinámico** - Se adapta al rango del personaje
- **Eficiente** - Mínimo impacto en performance
- **Configurable** - Sin necesidad de recompilación

## 📊 **Monitoreo**

### **Versión Normal**
```
HP: 1500/2000 (75%) | SP: 800/1000 (80%) | Markers: 3 | Range: 2.5
```

### **Versión WoE**
```
[WoE] HP:1500/2000(75%) SP:800/1000(80%) | Pots:15 | Spams:45 | Keys:3 | Markers:5 | Range:2.5 | Cache:true(2.1ms)
```

## 🎯 **Casos de Uso**

### **1. PvP (Player vs Player)**
- Posicionamiento perfecto para ataques
- Evitar estar fuera de rango
- Coordinar ataques de equipo

### **2. WoE (War of Emperium)**
- Asaltos coordinados
- Defensa estratégica
- Control de territorio

### **3. PvM (Player vs Monster)**
- Optimización de farming
- Evitar daño innecesario
- Eficiencia de combate

## 🔒 **Seguridad**

- **Solo visual** - No modifica el juego
- **No detectable** - Overlay externo
- **Configurable** - Se puede deshabilitar
- **Ligero** - Mínimo impacto en performance

## 📈 **Performance**

### **Versión Normal**
- **CPU**: < 0.5% adicional
- **Memoria**: < 2MB adicionales
- **FPS**: Sin impacto notable

### **Versión WoE**
- **CPU**: < 0.3% adicional (optimizado)
- **Memoria**: < 1MB adicionales (cache)
- **FPS**: Sin impacto (60 FPS target)

## 🎉 **Ejemplo de Uso**

### **1. Habilitar Range Display**
```bash
# Editar configuración
{
  "range_display": {
    "enabled": true,
    "color": "0x7F00FF00",
    "opacity": 0.5,
    "thickness": 2.0,
    "segments": 32
  }
}
```

### **2. Ejecutar DTAP**
```bash
# Versión Normal
./dtap.exe

# Versión WoE (recomendado)
./dtap-woe.exe
```

### **3. Ver el Círculo**
- El círculo verde aparecerá alrededor de tu personaje
- Se actualizará dinámicamente con tu rango de ataque
- Cambiará de tamaño según tu arma/clase

## 🎯 **Diferencias entre Versiones**

### **Versión Normal**
- Círculo básico con 32 segmentos
- Opacidad del 50%
- Grosor de línea 2.0
- Actualización cada 50ms

### **Versión WoE**
- Círculo optimizado con 64 segmentos
- Opacidad del 70% (más visible)
- Grosor de línea 3.0 (más grueso)
- Actualización cada 50ms (optimizada)
- Opción de mostrar solo en combate

## 🎉 **Conclusión**

El sistema de **Range Display** complementa perfectamente los **Skill Markers**, proporcionando información visual completa para el combate. Con el círculo de rango, ahora puedes:

- ✅ **Ver tu rango de ataque** en tiempo real
- ✅ **Posicionarte perfectamente** para combate
- ✅ **Coordinar ataques** de equipo
- ✅ **Optimizar tu gameplay** táctico

---

*¡Domina el campo de batalla con información visual completa!* 🎯⚔️
