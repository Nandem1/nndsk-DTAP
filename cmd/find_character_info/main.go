package main

import (
	"fmt"
	"strings"

	"dtap-go/internal/memory"
	"dtap-go/internal/process"
	"dtap-go/pkg/logger"
)

func main() {
	// Configurar logger
	logger.Global.SetLevel(logger.INFO)

	// Mostrar información de la aplicación
	logger.Info("=== FIND CHARACTER INFO - Búsqueda de Información del Personaje ===")
	logger.Info("Buscando strings relacionados con información del personaje...")
	logger.Info("================================================================")

	// Buscar proceso de RO
	logger.Info("Buscando proceso de Ragnarok Online...")
	roProcess, err := process.FindROProcess()
	if err != nil {
		logger.Error("Error: %v", err)
		logger.Info("Asegúrate de que Ragnarok Online esté ejecutándose.")
		logger.Info("Presiona Enter para salir...")
		fmt.Scanln()
		return
	}

	logger.Info("Proceso encontrado: %s (PID: %d)", roProcess.Name, roProcess.PID)

	// Crear cliente de memoria
	logger.Info("Conectando a la memoria del proceso...")
	memClient, err := memory.NewClient(roProcess.PID)
	if err != nil {
		logger.Error("Error conectando a memoria: %v", err)
		logger.Info("Presiona Enter para salir...")
		fmt.Scanln()
		return
	}
	defer memClient.Close()

	logger.Info("Conectado a la memoria del proceso")
	logger.Info("")

	// Búsqueda de información del personaje
	searchCharacterInfo(memClient)

	logger.Info("")
	logger.Info("Presiona Enter para salir...")
	fmt.Scanln()
}

// searchCharacterInfo busca strings relacionados con información del personaje
func searchCharacterInfo(memClient *memory.Client) {
	// Strings relacionados con información del personaje
	characterStrings := []string{
		// Información básica del personaje
		"Character",
		"character",
		"CHARACTER",
		"Player",
		"player",
		"PLAYER",
		"Avatar",
		"avatar",
		"AVATAR",
		"Entity",
		"entity",
		"ENTITY",

		// Stats del personaje
		"Stats",
		"stats",
		"STATS",
		"Status",
		"status",
		"STATUS",
		"Stat",
		"stat",
		"STAT",
		"Level",
		"level",
		"LEVEL",
		"Job",
		"job",
		"JOB",
		"Class",
		"class",
		"CLASS",

		// Atributos
		"STR",
		"str",
		"AGI",
		"agi",
		"VIT",
		"vit",
		"INT",
		"int",
		"DEX",
		"dex",
		"LUK",
		"luk",
		"Strength",
		"strength",
		"Agility",
		"agility",
		"Vitality",
		"vitality",
		"Intelligence",
		"intelligence",
		"Dexterity",
		"dexterity",
		"Luck",
		"luck",

		// Velocidad de ataque
		"ASPD",
		"aspd",
		"AttackSpeed",
		"attackSpeed",
		"ATTACK_SPEED",
		"Attack Speed",
		"attack speed",
		"ATTACK SPEED",
		"Velocidad de Ataque",
		"velocidad de ataque",
		"VELOCIDAD DE ATAQUE",

		// Movimiento
		"MoveSpeed",
		"moveSpeed",
		"MOVE_SPEED",
		"Move Speed",
		"move speed",
		"MOVE SPEED",
		"Velocidad de Movimiento",
		"velocidad de movimiento",
		"VELOCIDAD DE MOVIMIENTO",

		// HP/SP
		"HP",
		"hp",
		"SP",
		"sp",
		"Health",
		"health",
		"HEALTH",
		"Spirit",
		"spirit",
		"SPIRIT",
		"MaxHP",
		"maxHP",
		"MAX_HP",
		"MaxSP",
		"maxSP",
		"MAX_SP",

		// Experiencia
		"EXP",
		"exp",
		"Experience",
		"experience",
		"EXPERIENCE",
		"BaseEXP",
		"baseEXP",
		"BASE_EXP",
		"JobEXP",
		"jobEXP",
		"JOB_EXP",

		// Posición
		"Position",
		"position",
		"POSITION",
		"Pos",
		"pos",
		"POS",
		"X",
		"x",
		"Y",
		"y",
		"Z",
		"z",
		"Coord",
		"coord",
		"COORD",
		"Coordinate",
		"coordinate",
		"COORDINATE",

		// Dirección
		"Direction",
		"direction",
		"DIRECTION",
		"Dir",
		"dir",
		"DIR",
		"Angle",
		"angle",
		"ANGLE",
		"Rotation",
		"rotation",
		"ROTATION",

		// ID y nombres
		"ID",
		"id",
		"Name",
		"name",
		"NAME",
		"CharName",
		"charName",
		"CHAR_NAME",
		"PlayerName",
		"playerName",
		"PLAYER_NAME",
		"AccountID",
		"accountID",
		"ACCOUNT_ID",
		"CharID",
		"charID",
		"CHAR_ID",

		// Mapas
		"Map",
		"map",
		"MAP",
		"MapName",
		"mapName",
		"MAP_NAME",
		"CurrentMap",
		"currentMap",
		"CURRENT_MAP",
		"MapID",
		"mapID",
		"MAP_ID",

		// Guild
		"Guild",
		"guild",
		"GUILD",
		"GuildName",
		"guildName",
		"GUILD_NAME",
		"GuildID",
		"guildID",
		"GUILD_ID",

		// Party
		"Party",
		"party",
		"PARTY",
		"PartyID",
		"partyID",
		"PARTY_ID",
		"PartyLeader",
		"partyLeader",
		"PARTY_LEADER",

		// Equipamiento
		"Equipment",
		"equipment",
		"EQUIPMENT",
		"Equip",
		"equip",
		"EQUIP",
		"Weapon",
		"weapon",
		"WEAPON",
		"Armor",
		"armor",
		"ARMOR",
		"Shield",
		"shield",
		"SHIELD",
		"Helmet",
		"helmet",
		"HELMET",
		"Shoes",
		"shoes",
		"SHOES",
		"Accessory",
		"accessory",
		"ACCESSORY",

		// Skills
		"Skill",
		"skill",
		"SKILL",
		"Skills",
		"skills",
		"SKILLS",
		"SkillPoint",
		"skillPoint",
		"SKILL_POINT",
		"SkillLevel",
		"skillLevel",
		"SKILL_LEVEL",

		// Items
		"Item",
		"item",
		"ITEM",
		"Items",
		"items",
		"ITEMS",
		"Inventory",
		"inventory",
		"INVENTORY",
		"Inv",
		"inv",
		"INV",

		// Money
		"Money",
		"money",
		"MONEY",
		"Zeny",
		"zeny",
		"ZENY",
		"Gold",
		"gold",
		"GOLD",
		"Cash",
		"cash",
		"CASH",

		// Estado
		"State",
		"state",
		"STATE",
		"Status",
		"status",
		"STATUS",
		"Condition",
		"condition",
		"CONDITION",
		"Buff",
		"buff",
		"BUFF",
		"Debuff",
		"debuff",
		"DEBUFF",

		// Comportamiento
		"Behavior",
		"behavior",
		"BEHAVIOR",
		"Action",
		"action",
		"ACTION",
		"Activity",
		"activity",
		"ACTIVITY",
		"Activity",
		"activity",
		"ACTIVITY",

		// Físico
		"Physics",
		"physics",
		"PHYSICS",
		"Body",
		"body",
		"BODY",
		"Collision",
		"collision",
		"COLLISION",
		"HitBox",
		"hitBox",
		"HIT_BOX",
		"BoundingBox",
		"boundingBox",
		"BOUNDING_BOX",

		// Renderizado
		"Render",
		"render",
		"RENDER",
		"Draw",
		"draw",
		"DRAW",
		"Model",
		"model",
		"MODEL",
		"Mesh",
		"mesh",
		"MESH",
		"Texture",
		"texture",
		"TEXTURE",
		"Animation",
		"animation",
		"ANIMATION",

		// Red
		"Network",
		"network",
		"NETWORK",
		"Packet",
		"packet",
		"PACKET",
		"Protocol",
		"protocol",
		"PROTOCOL",
		"Client",
		"client",
		"CLIENT",
		"Server",
		"server",
		"SERVER",

		// Base de datos
		"Database",
		"database",
		"DATABASE",
		"DB",
		"db",
		"Table",
		"table",
		"TABLE",
		"Record",
		"record",
		"RECORD",
		"Field",
		"field",
		"FIELD",

		// Configuración
		"Config",
		"config",
		"CONFIG",
		"Setting",
		"setting",
		"SETTING",
		"Option",
		"option",
		"OPTION",
		"Parameter",
		"parameter",
		"PARAMETER",

		// Memoria
		"Memory",
		"memory",
		"MEMORY",
		"Buffer",
		"buffer",
		"BUFFER",
		"Data",
		"data",
		"DATA",
		"Struct",
		"struct",
		"STRUCT",
		"Class",
		"class",
		"CLASS",
		"Object",
		"object",
		"OBJECT",
		"Instance",
		"instance",
		"INSTANCE",
	}

	totalFound := 0

	logger.Info("🔍 INICIANDO BÚSQUEDA DE INFORMACIÓN DEL PERSONAJE")
	logger.Info("Buscando %d strings relacionados con información del personaje", len(characterStrings))
	logger.Info("")

	// === RANGO 1: ALREDEDOR DEL PERSONAJE ===
	logger.Info("📍 RANGO 1: Alrededor del personaje (0x10df5d8)")
	baseAddr := uintptr(0x10df5d8)
	found1 := searchCharacterInRange(memClient, baseAddr, -2000, 4000, characterStrings, "personaje")
	totalFound += found1

	// === RANGO 2: ALREDEDOR DE HP/SP ===
	logger.Info("")
	logger.Info("📍 RANGO 2: Alrededor de HP/SP (0x10DCE10)")
	hpAddr := uintptr(0x10DCE10)
	found2 := searchCharacterInRange(memClient, hpAddr, -4000, 8000, characterStrings, "HP/SP")
	totalFound += found2

	// === RANGO 3: RANGOS CONOCIDOS DE RO ===
	logger.Info("")
	logger.Info("📍 RANGO 3: Rangos conocidos de RO")
	found3 := searchCharacterInKnownRanges(memClient, characterStrings)
	totalFound += found3

	// === RANGO 4: BÚSQUEDA EN HEAP ===
	logger.Info("")
	logger.Info("📍 RANGO 4: Heap (0x10000000 - 0x20000000)")
	found4 := searchCharacterInHeap(memClient, characterStrings)
	totalFound += found4

	// === RESUMEN FINAL ===
	logger.Info("")
	logger.Info("🎯 RESUMEN DE BÚSQUEDA:")
	logger.Info("   Rango 1 (Personaje): %d encontrados", found1)
	logger.Info("   Rango 2 (HP/SP):     %d encontrados", found2)
	logger.Info("   Rango 3 (Conocidos): %d encontrados", found3)
	logger.Info("   Rango 4 (Heap):      %d encontrados", found4)
	logger.Info("   TOTAL:               %d encontrados", totalFound)

	if totalFound == 0 {
		logger.Info("")
		logger.Info("❌ No se encontraron strings de información del personaje en ningún rango.")
		logger.Info("   Posibles causas:")
		logger.Info("   - El personaje no está en el juego")
		logger.Info("   - Los strings están en un formato diferente")
		logger.Info("   - La dirección está en un rango no explorado")
		logger.Info("   - Los strings están encriptados o comprimidos")
	} else {
		logger.Info("")
		logger.Info("✅ Se encontraron %d strings relacionados con información del personaje", totalFound)
		logger.Info("   Revisa las direcciones encontradas arriba.")
		logger.Info("")
		logger.Info("💡 ANÁLISIS:")
		logger.Info("   - Los strings encontrados pueden indicar estructuras de datos del personaje")
		logger.Info("   - Busca patrones como 'Character', 'Player', 'Stats', 'ASPD'")
		logger.Info("   - Las direcciones con múltiples strings pueden ser estructuras completas")
		logger.Info("   - Los strings de 'Position', 'Direction' pueden indicar datos de movimiento")
	}
}

// searchCharacterInRange busca strings de información del personaje en un rango específico
func searchCharacterInRange(memClient *memory.Client, baseAddr uintptr, startOffset, endOffset int, targetStrings []string, rangeName string) int {
	found := 0
	logger.Info("   Buscando información del personaje en %s (offset %d a %d)...", rangeName, startOffset, endOffset)

	for i := startOffset; i <= endOffset; i++ {
		addr := baseAddr + uintptr(i)
		// Leer 512 bytes para buscar strings
		data, err := memClient.ReadBytes(addr, 512)
		if err == nil {
			dataStr := string(data)
			for _, targetStr := range targetStrings {
				if strings.Contains(dataStr, targetStr) {
					found++
					logger.Info("   ✅ String '%s' encontrado en: 0x%x (offset: %d)", targetStr, addr, i)
					// Mostrar contexto del string
					showCharacterContext(dataStr, targetStr, addr)
					break // Solo mostrar una vez por dirección
				}
			}
		}
	}

	if found == 0 {
		logger.Info("   ❌ No encontrado en %s", rangeName)
	}
	return found
}

// searchCharacterInKnownRanges busca strings de información del personaje en rangos conocidos de RO
func searchCharacterInKnownRanges(memClient *memory.Client, targetStrings []string) int {
	found := 0
	logger.Info("   Buscando información del personaje en rangos conocidos de RO...")

	// Rangos conocidos de RO
	knownRanges := []struct {
		name  string
		start uintptr
		end   uintptr
	}{
		{"Status Buffer", 0x10DD000, 0x10DE000},
		{"Character Data", 0x10DF000, 0x10E0000},
		{"Map Data", 0x10E0000, 0x10E1000},
		{"Inventory", 0x10E1000, 0x10E2000},
		{"Skills", 0x10E2000, 0x10E3000},
		{"Stats", 0x10E3000, 0x10E4000},
		{"Equipment", 0x10E4000, 0x10E5000},
		{"Party", 0x10E5000, 0x10E6000},
		{"Guild", 0x10E6000, 0x10E7000},
		{"Chat", 0x10E7000, 0x10E8000},
		{"UI Text", 0x10E8000, 0x10E9000},
		{"Menu Text", 0x10E9000, 0x10EA000},
		{"Error Messages", 0x10EA000, 0x10EB000},
		{"System Messages", 0x10EB000, 0x10EC000},
		{"Item Names", 0x10EC000, 0x10ED000},
		{"Animation Data", 0x10ED000, 0x10EE000},
		{"Model Data", 0x10EE000, 0x10EF000},
		{"Texture Data", 0x10EF000, 0x10F0000},
		{"Player Data", 0x10F0000, 0x10F1000},
		{"Entity Data", 0x10F1000, 0x10F2000},
	}

	for _, r := range knownRanges {
		rangeFound := 0
		for addr := r.start; addr < r.end; addr += 4 {
			// Leer 256 bytes para buscar strings
			data, err := memClient.ReadBytes(addr, 256)
			if err == nil {
				dataStr := string(data)
				for _, targetStr := range targetStrings {
					if strings.Contains(dataStr, targetStr) {
						found++
						rangeFound++
						if rangeFound <= 2 {
							logger.Info("   ✅ String '%s' en %s: 0x%x", targetStr, r.name, addr)
							showCharacterContext(dataStr, targetStr, addr)
						}
						break
					}
				}
			}
		}
		if rangeFound > 2 {
			logger.Info("   ... y %d más en %s", rangeFound-2, r.name)
		}
	}

	if found == 0 {
		logger.Info("   ❌ No encontrado en rangos conocidos")
	}
	return found
}

// searchCharacterInHeap busca strings de información del personaje en el heap
func searchCharacterInHeap(memClient *memory.Client, targetStrings []string) int {
	found := 0
	logger.Info("   Buscando información del personaje en heap...")

	// Buscar en bloques de 1000 direcciones
	for block := 0; block < 100; block++ { // Reducido para evitar detección
		blockFound := 0
		startAddr := uintptr(0x10000000 + uintptr(block*4000))

		for i := 0; i < 1000; i++ {
			addr := startAddr + uintptr(i)
			// Leer 128 bytes para buscar strings
			data, err := memClient.ReadBytes(addr, 128)
			if err == nil {
				dataStr := string(data)
				for _, targetStr := range targetStrings {
					if strings.Contains(dataStr, targetStr) {
						found++
						blockFound++
						if blockFound <= 2 {
							logger.Info("   ✅ String '%s' en Heap: 0x%x", targetStr, addr)
						}
						break
					}
				}
			}
		}

		if blockFound > 2 {
			logger.Info("   ... y %d más en bloque 0x%x", blockFound-2, startAddr)
		}
	}

	if found == 0 {
		logger.Info("   ❌ No encontrado en heap")
	}
	return found
}

// showCharacterContext muestra el contexto alrededor del string de información del personaje encontrado
func showCharacterContext(dataStr, targetStr string, addr uintptr) {
	// Encontrar la posición del string
	pos := strings.Index(dataStr, targetStr)
	if pos == -1 {
		return
	}

	// Mostrar contexto (64 caracteres antes y después)
	start := pos - 64
	if start < 0 {
		start = 0
	}
	end := pos + len(targetStr) + 64
	if end > len(dataStr) {
		end = len(dataStr)
	}

	context := dataStr[start:end]
	// Limpiar caracteres no imprimibles
	cleanContext := ""
	for _, c := range context {
		if c >= 32 && c <= 126 {
			cleanContext += string(c)
		} else {
			cleanContext += "."
		}
	}

	logger.Info("      Contexto: '%s'", cleanContext)
}
