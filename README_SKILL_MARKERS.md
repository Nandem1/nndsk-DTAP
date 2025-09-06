# 🎯 DTAP Skill Markers - Sistema de Marcado de Habilidades

## 📋 Descripción

El sistema de **Skill Markers** (M2E - MiniMiniEffect) es una funcionalidad avanzada que marca habilidades en el piso del juego con colores personalizables, proporcionando una ventaja táctica significativa en combate.

## ✨ Características

### 🎨 **Marcado Visual**
- **Colores personalizables** por tipo de habilidad
- **Transparencia ajustable** (0.0 - 1.0)
- **Tamaño configurable** de las marcas
- **Duración automática** basada en el tipo de skill

### 🏗️ **Habilidades Soportadas**
- **Stone Wall (SW)** - Muros de piedra (Gris)
- **Fire Wall (FW)** - Muros de fuego (Rojo)
- **Sanctuary** - Santuario (Cian)
- **Safety Wall (SG)** - Muro de seguridad (Púrpura)
- **Fire Pillar (FP)** - Pilares de fuego (Amarillo)
- **Magnus Exorcismus (ME)** - Exorcismo (Cian)
- **Pneuma** - Protección contra proyectiles (Cian)
- **Ice Wall (IW)** - Muros de hielo (Púrpura)
- **Quagmire (QM)** - Pantano (Verde)
- **Warp Portal** - Portal de teletransporte (Blanco)

### ⚙️ **Configuración**

#### Versión Normal
```json
{
  "skill_markers": {
    "enabled": true,
    "opacity": 0.5,
    "size": 4.8,
    "skills": {
      "stone_wall": {
        "color": "0x7F008888",
        "enabled": true
      }
    }
  }
}
```

#### Versión WoE (Optimizada)
```json
{
  "skill_markers": {
    "enabled": true,
    "opacity": 0.7,
    "size": 4.8,
    "woe_optimized": true,
    "show_enemy_only": false,
    "show_allied_only": false,
    "skills": {
      "stone_wall": {
        "color": "0x7F008888",
        "enabled": true
      }
    }
  }
}
```

## 🎮 **Uso**

### 1. **Habilitar Skill Markers**
```bash
# Versión Normal
./dtap.exe

# Versión WoE
./dtap-woe.exe
```

### 2. **Configurar Colores**
Los colores se configuran en formato hexadecimal: `0xAARRGGBB`
- **AA**: Alpha (transparencia) - 00-FF
- **RR**: Red (rojo) - 00-FF  
- **GG**: Green (verde) - 00-FF
- **BB**: Blue (azul) - 00-FF

### 3. **Ejemplos de Colores**
```json
"0x7F008888"  // Gris semi-transparente (50% opacidad)
"0x7F880000"  // Rojo semi-transparente
"0x7F00FFFF"  // Cian semi-transparente
"0x7F880088"  // Púrpura semi-transparente
"0x7F888800"  // Amarillo semi-transparente
"0x7FFFFFFF"  // Blanco semi-transparente
```

## 🚀 **Ventajas Tácticas**

### ⚔️ **En Combate Normal**
- **Identificar amenazas** - Ver muros enemigos
- **Evitar trampas** - Detectar habilidades hostiles
- **Posicionamiento** - Coordinar movimientos

### 🏰 **En WoE (War of Emperium)**
- **Ventaja competitiva** - Información en tiempo real
- **Coordinación de equipo** - Marcar objetivos
- **Timing perfecto** - Duración exacta de habilidades
- **Estrategia avanzada** - Planificación táctica

## 🔧 **Configuración Avanzada**

### **Parámetros Principales**
- `enabled`: Habilita/deshabilita el sistema
- `opacity`: Transparencia (0.0 = invisible, 1.0 = opaco)
- `size`: Tamaño de las marcas en unidades de juego

### **Parámetros WoE**
- `woe_optimized`: Optimizaciones específicas para WoE
- `show_enemy_only`: Mostrar solo habilidades enemigas
- `show_allied_only`: Mostrar solo habilidades aliadas

### **Configuración por Skill**
- `color`: Color hexadecimal de la marca
- `enabled`: Habilita/deshabilita el skill específico

## 📊 **Monitoreo**

### **Versión Normal**
```
HP: 1500/2000 (75%) | SP: 800/1000 (80%) | Markers: 3
```

### **Versión WoE**
```
[WoE] HP:1500/2000(75%) SP:800/1000(80%) | Pots:15 | Spams:45 | Keys:3 | Markers:5 | Cache:true(2.1ms)
```

## 🛠️ **Implementación Técnica**

### **Arquitectura**
```
pkg/visual/
├── skill_marker_config.go    # Configuración de colores
├── skill_marker.go          # Renderizador de marcas
├── overlay.go               # Manager del overlay
└── config_adapter.go        # Adaptador de configuraciones
```

### **Características Técnicas**
- **Tiempo real** - Actualización cada 100ms
- **Memoria eficiente** - Gestión automática de marcas expiradas
- **Thread-safe** - Operaciones concurrentes seguras
- **Configurable** - Sin necesidad de recompilación

## 🎯 **Casos de Uso**

### **1. PvP (Player vs Player)**
- Detectar muros enemigos
- Evitar áreas de efecto
- Coordinar ataques

### **2. WoE (War of Emperium)**
- Marcar defensas enemigas
- Identificar puntos de entrada
- Coordinar asaltos

### **3. PvM (Player vs Monster)**
- Evitar habilidades de jefes
- Posicionamiento estratégico
- Optimización de rutas

## 🔒 **Seguridad**

- **Solo visual** - No modifica el juego
- **No detectable** - Overlay externo
- **Configurable** - Se puede deshabilitar
- **Ligero** - Mínimo impacto en performance

## 📈 **Performance**

### **Versión Normal**
- **CPU**: < 1% adicional
- **Memoria**: < 5MB adicionales
- **FPS**: Sin impacto notable

### **Versión WoE**
- **CPU**: < 0.5% adicional (optimizado)
- **Memoria**: < 3MB adicionales (cache)
- **FPS**: Sin impacto (60 FPS target)

## 🎉 **Conclusión**

El sistema de **Skill Markers** transforma DTAP en una herramienta táctica avanzada, proporcionando información visual crucial para el éxito en combate. Con configuración flexible y optimizaciones específicas para WoE, es la funcionalidad perfecta para jugadores competitivos.

---

*¡Domina el campo de batalla con información visual en tiempo real!* 🎯⚔️
