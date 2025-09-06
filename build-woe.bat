@echo off
echo ========================================
echo DTAP-Go WoE Optimized Build Script
echo ========================================
echo.

echo [1/4] Cleaning previous builds...
if exist dtap-woe.exe del /q dtap-woe.exe
if exist dtap.exe del /q dtap.exe
if exist find_addresses.exe del /q find_addresses.exe
echo ✅ Clean complete

echo.
echo [2/4] Installing dependencies...
go mod tidy
go mod download
echo ✅ Dependencies installed

echo.
echo [3/4] Building WoE Optimized version...
go build -ldflags "-s -w" -o dtap-woe.exe main_woe.go
if %errorlevel% neq 0 (
    echo ❌ Build failed!
    pause
    exit /b 1
)
echo ✅ WoE build complete

echo.
echo [4/4] Building standard version (backup)...
go build -ldflags "-s -w" -o dtap.exe main.go
if %errorlevel% neq 0 (
    echo ⚠️  Standard build failed, but WoE build succeeded
) else (
    echo ✅ Standard build complete
)

echo.
echo ========================================
echo 🚀 BUILD COMPLETE! 🚀
echo ========================================
echo.
echo Files created:
if exist dtap-woe.exe echo ✅ dtap-woe.exe (WoE Optimized - RECOMMENDED)
if exist dtap.exe echo ✅ dtap.exe (Standard version)
echo.
echo To run WoE optimized version:
echo   dtap-woe.exe
echo.
echo To run standard version:
echo   dtap.exe
echo.
echo ========================================
pause