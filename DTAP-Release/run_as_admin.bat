@echo off
echo ========================================
echo    DTAP - Ejecutar como Administrador
echo ========================================
echo.

echo Verificando permisos de administrador...
net session >nul 2>&1
if %errorLevel% == 0 (
    echo ✓ Ejecutando como administrador
    echo.
    echo Iniciando DTAP...
    dtap.exe
) else (
    echo ❌ Se requieren permisos de administrador
    echo Solicitando elevacion de privilegios...
    powershell -Command "Start-Process '%~f0' -Verb RunAs"
    exit /b
)

echo.
echo DTAP ha terminado.
pause
