# DTAP Hook Native - Compilación y Uso

## Resumen

Se ha compilado exitosamente el hook nativo `dtap-hook-native.dll` usando MinGW MSYS. Este hook proporciona funcionalidades de comunicación entre DTAP y el cliente de Ragnarok Online.

## Archivos Generados

- `dtap-hook-native.dll` (46,987 bytes) - DLL principal compilada
- `dtap-hook-native.lib` (4,852 bytes) - Librería de importación
- `dtap-hook-native.def` (91 bytes) - Archivo de definición de exportaciones
- `dtap-hook-native.c` (2,976 bytes) - Código fuente C

## Funciones Exportadas

La DLL exporta las siguientes funciones:

1. **EnableBBE(BOOL enabled)** - Habilita/deshabilita el display de rango
2. **EnableM2E(BOOL enabled)** - Habilita/deshabilita los marcadores de habilidades
3. **IsInitialized()** - Verifica si la DLL está inicializada
4. **SetAlphaLevel(INT32 level)** - Establece el nivel de transparencia
5. **SetGroundZBias(INT32 bias)** - Establece el sesgo Z del suelo
6. **TestFunction()** - Función de prueba que retorna 42

## Características Técnicas

- **Arquitectura**: x86-64 (64-bit)
- **Formato**: PE32+ (Windows)
- **Compilador**: GCC 15.1.0 (MinGW-W64)
- **Memoria compartida**: Utiliza `CreateFileMappingA` para comunicación
- **Inicialización**: Se ejecuta automáticamente al cargar la DLL

## Uso desde Go

```go
package main

import (
    "fmt"
    "syscall"
)

var (
    hookDLL = syscall.NewLazyDLL("dtap-hook-native.dll")
    enableBBE = hookDLL.NewProc("EnableBBE")
    enableM2E = hookDLL.NewProc("EnableM2E")
    // ... otras funciones
)

func main() {
    // Habilitar display de rango
    enableBBE.Call(1)
    
    // Habilitar marcadores de habilidades
    enableM2E.Call(1)
    
    // Establecer nivel de transparencia
    setAlphaLevel.Call(128)
}
```

## Comando de Compilación

```bash
gcc -shared -o dtap-hook-native.dll dtap-hook-native.c \
    "-Wl,--enable-stdcall-fixup" \
    "-Wl,--enable-auto-import" \
    "-Wl,--enable-runtime-pseudo-reloc" \
    "-Wl,--out-implib,dtap-hook-native.lib"
```

## Estado de las Pruebas

✅ **Compilación exitosa** - La DLL se compila sin errores
✅ **Exportaciones correctas** - Todas las funciones están disponibles
✅ **Carga desde Go** - Se puede cargar y usar desde código Go
✅ **Funciones operativas** - Todas las funciones responden correctamente
✅ **Memoria compartida** - Sistema de comunicación funcional

## Próximos Pasos

1. Integrar el hook nativo con el sistema de inyección de DTAP
2. Implementar la comunicación bidireccional con el cliente
3. Agregar funcionalidades de renderizado visual
4. Optimizar el rendimiento para uso en tiempo real

## Notas Técnicas

- La DLL utiliza memoria compartida para comunicación con DTAP
- Se inicializa automáticamente al ser cargada
- Compatible con Windows 64-bit
- Requiere MinGW MSYS para compilación
- Las funciones void retornan códigos de estado de Windows
