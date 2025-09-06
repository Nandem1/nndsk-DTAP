#include <windows.h>
#include <iostream>
#include <string>

// Shared memory structure for communication with DTAP
struct DTAPSharedData {
    bool skillMarkersEnabled;
    bool rangeDisplayEnabled;
    int groundZBias;
    int alphaLevel;
    char sharedMemoryName[64];
};

// Global variables
DTAPSharedData* g_sharedData = nullptr;
HANDLE g_sharedMemoryHandle = nullptr;
bool g_initialized = false;

// Function to initialize shared memory
bool InitializeSharedMemory() {
    // Create or open shared memory
    g_sharedMemoryHandle = CreateFileMappingA(
        INVALID_HANDLE_VALUE,
        NULL,
        PAGE_READWRITE,
        0,
        sizeof(DTAPSharedData),
        "DTAPSharedMemory"
    );

    if (g_sharedMemoryHandle == nullptr) {
        return false;
    }

    // Map the shared memory
    g_sharedData = (DTAPSharedData*)MapViewOfFile(
        g_sharedMemoryHandle,
        FILE_MAP_ALL_ACCESS,
        0,
        0,
        sizeof(DTAPSharedData)
    );

    if (g_sharedData == nullptr) {
        CloseHandle(g_sharedMemoryHandle);
        return false;
    }

    // Initialize default values
    g_sharedData->skillMarkersEnabled = false;
    g_sharedData->rangeDisplayEnabled = false;
    g_sharedData->groundZBias = 0;
    g_sharedData->alphaLevel = 0;
    strcpy_s(g_sharedData->sharedMemoryName, "DTAPSharedMemory");

    return true;
}

// Function to enable skill markers
void EnableSkillMarkers(bool enabled) {
    if (g_sharedData) {
        g_sharedData->skillMarkersEnabled = enabled;
    }
}

// Function to enable range display
void EnableRangeDisplay(bool enabled) {
    if (g_sharedData) {
        g_sharedData->rangeDisplayEnabled = enabled;
    }
}

// Function to set ground Z bias
void SetGroundZBias(int bias) {
    if (g_sharedData) {
        g_sharedData->groundZBias = bias;
    }
}

// Function to set alpha level
void SetAlphaLevel(int level) {
    if (g_sharedData) {
        g_sharedData->alphaLevel = level;
    }
}

// DLL Entry Point
BOOL APIENTRY DllMain(HMODULE hModule, DWORD ul_reason_for_call, LPVOID lpReserved) {
    switch (ul_reason_for_call) {
    case DLL_PROCESS_ATTACH:
        // Initialize when DLL is loaded
        if (InitializeSharedMemory()) {
            g_initialized = true;
            // Enable skill markers and range display by default
            EnableSkillMarkers(true);
            EnableRangeDisplay(true);
        }
        break;
    case DLL_THREAD_ATTACH:
        break;
    case DLL_THREAD_DETACH:
        break;
    case DLL_PROCESS_DETACH:
        // Cleanup when DLL is unloaded
        if (g_sharedData) {
            UnmapViewOfFile(g_sharedData);
        }
        if (g_sharedMemoryHandle) {
            CloseHandle(g_sharedMemoryHandle);
        }
        break;
    }
    return TRUE;
}

// Export functions for external access
extern "C" {
    __declspec(dllexport) void EnableM2E(bool enabled) {
        EnableSkillMarkers(enabled);
    }

    __declspec(dllexport) void EnableBBE(bool enabled) {
        EnableRangeDisplay(enabled);
    }

    __declspec(dllexport) void SetGroundZBias(int bias) {
        SetGroundZBias(bias);
    }

    __declspec(dllexport) void SetAlphaLevel(int level) {
        SetAlphaLevel(level);
    }

    __declspec(dllexport) bool IsInitialized() {
        return g_initialized;
    }
}
