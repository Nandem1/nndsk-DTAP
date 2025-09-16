# Mejoras del Spammer - Control de Ratio Teclas/Clicks

## 🎯 Problema Solucionado

El spammer anteriormente enviaba un click del mouse por cada tecla presionada, lo que causaba:
- Movimiento involuntario del personaje
- Clicks excesivos que interferían con el gameplay
- Falta de control sobre la frecuencia de clicks

## ✨ Nuevas Funcionalidades

### 1. **Control de Ratio Teclas/Clicks**
- **Configuración**: `key_click_ratio` (por defecto: 3)
- **Funcionamiento**: Envía 1 click por cada N teclas presionadas
- **Ejemplo**: Con ratio 3, se envían 3 teclas antes de 1 click

### 2. **Delay Adicional para Clicks**
- **Configuración**: `click_delay_ms` (por defecto: 50ms)
- **Funcionamiento**: Tiempo mínimo entre clicks del mouse
- **Beneficio**: Evita clicks demasiado rápidos

## 🔧 Configuración

### Archivo de Configuración
```json
{
  "spammer": {
    "enabled": true,
    "keys": ["F1", "F2", "F3"],
    "delay_ms": 0.1,
    "mouse_click": true,
    "mode": "simple",
    "key_click_ratio": 3,    // 3 teclas por 1 click
    "click_delay_ms": 50     // 50ms entre clicks
  }
}
```

### Interfaz de Usuario
En la sección "Spammer" del panel de Cheats:
- **Key/Click Ratio**: Controla cuántas teclas se envían por click
- **Click Delay**: Delay adicional entre clicks en milisegundos

## 📊 Ejemplos de Configuración

### Configuración Conservadora (Menos Clicks)
```json
{
  "key_click_ratio": 5,     // 5 teclas por 1 click
  "click_delay_ms": 100     // 100ms entre clicks
}
```

### Configuración Agresiva (Más Clicks)
```json
{
  "key_click_ratio": 2,     // 2 teclas por 1 click
  "click_delay_ms": 25      // 25ms entre clicks
```

### Solo Teclas (Sin Clicks)
```json
{
  "mouse_click": false      // Deshabilitar clicks completamente
}
```

## 🎮 Casos de Uso

### 1. **Farming/Leveling**
- **Ratio**: 4-5 teclas por click
- **Delay**: 75-100ms
- **Beneficio**: Menos movimiento, más estabilidad

### 2. **PvP/Combat**
- **Ratio**: 2-3 teclas por click
- **Delay**: 25-50ms
- **Beneficio**: Respuesta más rápida

### 3. **Skills con Cast Time**
- **Ratio**: 1-2 teclas por click
- **Delay**: 10-25ms
- **Beneficio**: Máxima velocidad

## 🔄 Lógica de Funcionamiento

```go
// Pseudocódigo de la lógica
func shouldSendClick() bool {
    // Verificar ratio de teclas
    if keyPressCount < keyClickRatio {
        return false
    }
    
    // Verificar delay mínimo
    if timeSinceLastClick < clickDelay {
        return false
    }
    
    return true
}
```

## 📈 Beneficios

1. **Control Preciso**: Ajusta exactamente cuántos clicks quieres
2. **Estabilidad**: Evita movimiento involuntario del personaje
3. **Flexibilidad**: Diferentes configuraciones para diferentes situaciones
4. **Compatibilidad**: Funciona con todas las configuraciones existentes

## 🚀 Uso Recomendado

1. **Empezar con valores por defecto** (ratio: 3, delay: 50ms)
2. **Ajustar según necesidad**:
   - Si el personaje se mueve mucho → Aumentar ratio
   - Si necesitas más velocidad → Disminuir ratio y delay
3. **Probar en diferentes situaciones** de juego
4. **Guardar configuraciones** para diferentes actividades

## 🔧 Archivos Modificados

- `pkg/config/config.go` - Estructura de configuración
- `internal/normal/spammer/ultra_fast_spammer.go` - Lógica del spammer
- `dtap-gui/frontend/components/Cheats.tsx` - Interfaz de usuario
- `dtap-gui/frontend/components/App.tsx` - Manejo de estado
- `dtap-config.json` - Configuración por defecto

## 🎯 Resultado Final

Ahora tienes control total sobre la frecuencia de clicks del spammer, permitiendo:
- ✅ Jugar sin movimiento involuntario
- ✅ Ajustar velocidad según necesidad
- ✅ Configurar diferentes perfiles para diferentes actividades
- ✅ Mantener la funcionalidad existente intacta
