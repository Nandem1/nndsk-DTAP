# 🔧 DTAP Injection System - Sistema de Inyección de DLL

## 📋 Descripción

El sistema de inyección de DTAP está diseñado para proporcionar funcionalidades avanzadas como **skill markers** y **range display** mediante la inyección de una DLL en el proceso de Ragnarok Online, similar al enfoque utilizado por SimpleROHook.

## 🏗️ Arquitectura

### **Componentes Principales**

```
pkg/injection/
├── minimal_injector.go      # Inyector básico (simulación)
├── injection_manager.go     # Manager principal del sistema
└── config.go               # Configuración del sistema
```

### **Flujo de Inyección**

1. **Detección del Proceso** - Encuentra el proceso de Ragnarok Online
2. **Inicialización** - Crea el sistema de inyección
3. **Configuración** - Establece parámetros de la DLL
4. **Inyección** - Inyecta la DLL en el proceso objetivo
5. **Comunicación** - Establece comunicación con la DLL inyectada

## ⚙️ Configuración

### **Configuración Básica**
```json
{
  "injection": {
    "enabled": true,
    "dll_path": "dtap-hook.dll",
    "shared_memory_name": "DTAPSharedMemory",
    "hook_type": 5,
    "auto_inject": true,
    "enable_m2e": true,
    "enable_bbe": true,
    "ground_z_bias": 0,
    "alpha_level": 0,
    "enable_dead_cell": true,
    "enable_chat_scope": true,
    "fix_window_mode": true,
    "show_frame_rate": true,
    "enable_object_info": false,
    "enable_audio_mode": false,
    "cpu_cooler_level": 0
  }
}
```

### **Parámetros de Configuración**

| Parámetro | Descripción | Valores |
|-----------|-------------|---------|
| `enabled` | Habilita el sistema de inyección | `true`/`false` |
| `dll_path` | Ruta a la DLL a inyectar | `"dtap-hook.dll"` |
| `shared_memory_name` | Nombre de la memoria compartida | `"DTAPSharedMemory"` |
| `hook_type` | Tipo de hook de Windows | `5` (WH_CBT) |
| `auto_inject` | Inyección automática al iniciar | `true`/`false` |
| `enable_m2e` | Habilita skill markers | `true`/`false` |
| `enable_bbe` | Habilita BBE | `true`/`false` |
| `ground_z_bias` | Sesgo Z del suelo | `0-100` |
| `alpha_level` | Nivel de transparencia | `0-255` |

## 🎯 Funcionalidades

### **1. Skill Markers (M2E)**
- **Marcado de habilidades** en el suelo del juego
- **Colores personalizables** para cada habilidad
- **Duración configurable** de las marcas
- **Detección automática** de habilidades

### **2. Range Display**
- **Círculo de rango** alrededor del personaje
- **Actualización dinámica** según el arma/clase
- **Colores y transparencia** configurables
- **Optimización para WoE**

### **3. Otras Funcionalidades**
- **BBE (Bubble Effect)** - Efectos de burbujas
- **Dead Cell** - Marcado de celdas muertas
- **Chat Scope** - Alcance del chat
- **Window Mode Fix** - Corrección de modo ventana
- **Frame Rate Display** - Mostrar FPS
- **Object Information** - Información de objetos
- **Audio Mode** - Modo de audio
- **CPU Cooler** - Enfriamiento de CPU

## 🔧 Implementación Técnica

### **Método de Inyección**
```go
// 1. Abrir proceso objetivo
processHandle := windows.OpenProcess(permissions, false, processID)

// 2. Obtener dirección de LoadLibraryW
loadLibraryAddr := windows.GetProcAddress(kernel32, "LoadLibraryW")

// 3. Asignar memoria en proceso objetivo
remoteMemory := windows.VirtualAllocEx(processHandle, ...)

// 4. Escribir ruta de DLL en memoria remota
windows.WriteProcessMemory(processHandle, remoteMemory, dllPath, ...)

// 5. Crear thread remoto para cargar DLL
threadHandle := windows.CreateRemoteThread(processHandle, ...)

// 6. Esperar completación
windows.WaitForSingleObject(threadHandle, windows.INFINITE)
```

### **Comunicación con DLL**
- **Memoria Compartida** - Para intercambio de datos
- **Hooks de Windows** - Para interceptar eventos
- **Configuración Dinámica** - Cambios en tiempo real

## 🚀 Uso

### **1. Preparación**
```bash
# Asegúrate de tener la DLL compilada
dtap-hook.dll  # DLL a inyectar
```

### **2. Configuración**
```json
{
  "injection": {
    "enabled": true,
    "dll_path": "dtap-hook.dll",
    "auto_inject": true
  }
}
```

### **3. Ejecución**
```bash
# Versión Normal
./dtap.exe

# Versión WoE
./dtap-woe.exe
```

### **4. Verificación**
```
[INFO] Sistema de inyección inicializado
[INFO] DLL inyectada exitosamente
[INFO] M2E (skill markers): enabled=true
[INFO] BBE: enabled=true
```

## 🔒 Seguridad

### **Consideraciones de Seguridad**
- **Solo lectura** - No modifica archivos del juego
- **Proceso local** - Solo afecta al proceso actual
- **Configurable** - Se puede deshabilitar fácilmente
- **Logging** - Registro completo de actividades

### **Detección Anti-Cheat**
- **Riesgo bajo** - Solo visual, no modifica gameplay
- **Configurable** - Se puede deshabilitar
- **Opcional** - No es requerido para funcionalidad básica

## 📊 Monitoreo

### **Logs del Sistema**
```
[INFO] Injection system initialized successfully
[INFO] Simulating DLL injection dll_path=dtap-hook.dll process_id=1234
[INFO] Note: Real DLL injection requires a compiled DLL file
[INFO] DLL injected successfully dll_path=C:\path\to\dtap-hook.dll
[INFO] M2E (skill markers) enabled=true
[INFO] BBE enabled=true
```

### **Estado del Sistema**
- **Inyectado**: `true`/`false`
- **Proceso Válido**: `true`/`false`
- **Configuración**: Valores actuales
- **Errores**: Logs de errores

## 🛠️ Desarrollo

### **Crear DLL Personalizada**
1. **Compilar DLL** en C/C++
2. **Implementar funciones** de hook
3. **Configurar memoria compartida**
4. **Probar inyección**

### **Extender Funcionalidades**
```go
// Agregar nueva funcionalidad
func (im *InjectionManager) EnableNewFeature(enabled bool) error {
    im.logger.Info("New Feature", "enabled", enabled)
    return nil
}
```

## 🎯 Próximos Pasos

### **Implementación Completa**
1. **Compilar DLL real** - Crear la DLL de hook
2. **Implementar hooks** - Interceptar funciones del juego
3. **Memoria compartida** - Comunicación bidireccional
4. **Testing** - Pruebas exhaustivas

### **Funcionalidades Avanzadas**
- **Packet Interception** - Interceptar paquetes
- **Memory Scanning** - Escaneo de memoria
- **Advanced Hooks** - Hooks más complejos
- **Performance Monitoring** - Monitoreo de rendimiento

## 🎉 Conclusión

El sistema de inyección de DTAP proporciona una base sólida para implementar funcionalidades avanzadas como skill markers y range display. Aunque actualmente es una simulación, la arquitectura está preparada para una implementación completa con DLL real.

### **Ventajas**
- ✅ **Arquitectura modular** - Fácil de extender
- ✅ **Configuración flexible** - Parámetros ajustables
- ✅ **Logging completo** - Debugging fácil
- ✅ **Seguridad** - Riesgo mínimo
- ✅ **Compatibilidad** - Funciona con ambas versiones

---

*¡El sistema de inyección está listo para la implementación completa!* 🔧⚡
