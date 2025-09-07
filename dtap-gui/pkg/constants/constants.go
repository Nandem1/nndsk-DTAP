package constants

// Memory addresses for HoneyRO client
const (
	// HP/SP addresses found by address finder
	HPBaseAddress    = 0x10DCE10 // HP Current
	HPMaxAddress     = 0x10DCE14 // HP Maximum
	SPBaseAddress    = 0x10DCE18 // SP Current
	SPMaxAddress     = 0x10DCE1C // SP Maximum
	StatusBufferAddr = 0x10DD284 // Status buffer (HP + 0x474)
)

// Process names to search for
var ROProcessNames = []string{
	"honeyro.exe", // HoneyRO
	"ragnarok.exe",
	"ragnarokonline.exe",
	"ro.exe",
	"kro.exe",
	"gravity.exe",
	"ragnarok_online.exe",
	"ragnarok_online_64.exe",
	"ragnarok_online_32.exe",
	"rag.exe",
	"ragnarokclient.exe",
	"ragnarok_client.exe",
	"roclient.exe",
	"ro_client.exe",
	"iRO.exe",
	"kRO.exe",
	"jRO.exe",
	"tRO.exe",
	"bRO.exe",
	"idRO.exe",
	"thRO.exe",
	"phRO.exe",
	"myRO.exe",
	"sgRO.exe",
	"twRO.exe",
	"hkRO.exe",
	"ro2.exe",
	"ragnarok2.exe",
}

// Default configuration values
const (
	DefaultHPPotKey   = "F1"
	DefaultSPPotKey   = "F2"
	DefaultHPPercent  = 50
	DefaultSPPercent  = 30
	DefaultCheckDelay = 100 // milliseconds
)

// Key codes for input
const (
	// Function keys
	KeyF1  = 0x70
	KeyF2  = 0x71
	KeyF3  = 0x72
	KeyF4  = 0x73
	KeyF5  = 0x74
	KeyF6  = 0x75
	KeyF7  = 0x76
	KeyF8  = 0x77
	KeyF9  = 0x78
	KeyF10 = 0x79
	KeyF11 = 0x7A
	KeyF12 = 0x7B

	// Mouse Messages
	WM_LBUTTONDOWN = 0x0201
	WM_LBUTTONUP   = 0x0202
	WM_RBUTTONDOWN = 0x0204
	WM_RBUTTONUP   = 0x0205
	WM_MBUTTONDOWN = 0x0207
	WM_MBUTTONUP   = 0x0208

	// Modifier Keys
	KeyShift = 0x10
	KeyCtrl  = 0x11
	KeyAlt   = 0x12

	// Number keys
	Key1 = 0x31
	Key2 = 0x32
	Key3 = 0x33
	Key4 = 0x34
	Key5 = 0x35
	Key6 = 0x36
	Key7 = 0x37
	Key8 = 0x38
	Key9 = 0x39
	Key0 = 0x30

	// Letter keys
	KeyA = 0x41
	KeyB = 0x42
	KeyC = 0x43
	KeyD = 0x44
	KeyE = 0x45
	KeyF = 0x46
	KeyG = 0x47
	KeyH = 0x48
	KeyI = 0x49
	KeyJ = 0x4A
	KeyK = 0x4B
	KeyL = 0x4C
	KeyM = 0x4D
	KeyN = 0x4E
	KeyO = 0x4F
	KeyP = 0x50
	KeyQ = 0x51
	KeyR = 0x52
	KeyS = 0x53
	KeyT = 0x54
	KeyU = 0x55
	KeyV = 0x56
	KeyW = 0x57
	KeyX = 0x58
	KeyY = 0x59
	KeyZ = 0x5A
)

// Windows message constants
const (
	WM_KEYDOWN = 0x0100
	WM_KEYUP   = 0x0101
)

// Memory access constants
const (
	ProcessAccessFlags = 0x001F0FFF // PROCESS_VM_READ | PROCESS_VM_WRITE | PROCESS_VM_OPERATION
)

// Application constants
const (
	AppName    = "DTAP-Go"
	AppVersion = "1.0.0"
	AppDesc    = "Ragnarok Online Auto Tool"
)

// Address finder constants - valores de ejemplo para testing
const (
	// Valores de HP/SP para address finder (cambiar según tu personaje)
	FinderHPCurrent = 6914
	FinderHPMax     = 6914
	FinderSPCurrent = 2511
	FinderSPMax     = 2511
)
