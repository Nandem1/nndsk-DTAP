#include <windows.h>

// Minimal test DLL - just returns TRUE from DllMain
BOOL APIENTRY DllMain(HMODULE hModule, DWORD ul_reason_for_call, LPVOID lpReserved) {
    switch (ul_reason_for_call) {
    case DLL_PROCESS_ATTACH:
        // Just return TRUE to indicate successful loading
        return TRUE;
    case DLL_THREAD_ATTACH:
    case DLL_THREAD_DETACH:
    case DLL_PROCESS_DETACH:
        break;
    }
    return TRUE;
}

// Simple exported function for testing
__declspec(dllexport) int TestFunction() {
    return 42;
}

// Hook procedure for Windows Hook injection
LRESULT CALLBACK HookProc(int nCode, WPARAM wParam, LPARAM lParam) {
    return CallNextHookEx(NULL, nCode, wParam, lParam);
}
