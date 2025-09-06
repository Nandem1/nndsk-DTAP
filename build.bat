@echo off
REM DTAP Build Script for Windows

echo Building DTAP - Normal and WoE Versions
echo ========================================

REM Build normal version
echo.
echo Building normal version...
go build -o dtap.exe cmd/dtap/main.go
if %errorlevel% neq 0 (
    echo ERROR: Failed to build normal version
    pause
    exit /b 1
)
echo Normal version built successfully!

REM Build WoE version
echo.
echo Building WoE version...
go build -o dtap-woe.exe cmd/dtap-woe/main.go
if %errorlevel% neq 0 (
    echo ERROR: Failed to build WoE version
    pause
    exit /b 1
)
echo WoE version built successfully!

echo.
echo ========================================
echo Build completed successfully!
echo.
echo Available executables:
echo - dtap.exe (Normal version)
echo - dtap-woe.exe (WoE version)
echo.
pause
