import { motion } from 'framer-motion'
import { useState } from 'react'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { 
  Activity, 
  Wifi, 
  WifiOff, 
  Play, 
  Pause,
  Clock,
  MousePointer,
  User,
  MapPin,
  Sword
} from 'lucide-react'

interface MainProps {
  isRunning: boolean
  isConnected: boolean
  stats: {
    hp: number
    maxHp: number
    sp: number
    maxSp: number
    uptime: string | number
    totalActions: number
    characterName?: string
    mapName?: string
    jobName?: string
    system?: {
      MemoryUsage: number
      MemoryTotal: number
      MemoryPercent: number
      GoMemoryUsage: number
      GoroutineCount: number
      CPUUsage: number
    }
  }
  message: string
  onConnect: () => void
  onStartStop: () => void
}

export default function Main({
  isRunning,
  isConnected,
  stats,
  message,
  onConnect,
  onStartStop
}: MainProps) {
  const [debugInfo, setDebugInfo] = useState<any>(null)
  const [showDebug, setShowDebug] = useState(false)
  
  const hpPercent = (stats.hp / stats.maxHp) * 100
  const spPercent = (stats.sp / stats.maxSp) * 100

  const formatUptime = (uptime: string | number) => {
    // Si el uptime viene como número (timestamp), convertirlo a formato legible
    if (uptime && !isNaN(Number(uptime))) {
      const seconds = Math.floor(Number(uptime)) // Redondear hacia abajo para eliminar decimales
      const hours = Math.floor(seconds / 3600)
      const minutes = Math.floor((seconds % 3600) / 60)
      const secs = seconds % 60
      return `${hours.toString().padStart(2, '0')}:${minutes.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`
    }
    // Si ya viene formateado, devolverlo tal como está
    return uptime || '00:00:00'
  }

  const formatBytes = (bytes: number) => {
    const units = ['B', 'KB', 'MB', 'GB']
    let size = bytes
    let unitIndex = 0
    
    while (size >= 1024 && unitIndex < units.length - 1) {
      size /= 1024
      unitIndex++
    }
    
    return `${size.toFixed(1)} ${units[unitIndex]}`
  }

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      exit={{ opacity: 0, y: -20 }}
      transition={{ duration: 0.3 }}
      className="h-full flex flex-col overflow-hidden"
    >
      <div className="flex-1 p-6 space-y-6 overflow-y-auto">
        {/* Header */}
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold">Dashboard</h1>
            <p className="text-muted-foreground">Monitor your Ragnarok Online session</p>
          </div>
          
          <div className="flex items-center gap-4">
            <div className="flex items-center gap-2">
              <div className={`status-indicator ${isConnected ? 'connected' : 'disconnected'}`}></div>
              <span className="text-sm font-medium">
                {isConnected ? 'Connected' : 'Disconnected'}
              </span>
            </div>
            
            <Button
              onClick={onConnect}
              disabled={isConnected}
              variant={isConnected ? "secondary" : "default"}
            >
              {isConnected ? (
                <>
                  <Wifi className="h-4 w-4 mr-2" />
                  Connected
                </>
              ) : (
                <>
                  <WifiOff className="h-4 w-4 mr-2" />
                  Connect
                </>
              )}
            </Button>
          </div>
        </div>

        {/* Main Grid */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {/* HP/SP Status */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Activity className="h-5 w-5" />
                Character Status
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              {/* Character Info */}
              {stats.characterName && (
                <div className="space-y-2 pb-4 border-b border-border">
                  <div className="flex items-center gap-2">
                    <User className="h-4 w-4 text-primary" />
                    <span className="text-sm font-medium">{stats.characterName}</span>
                  </div>
                  {stats.jobName && (
                    <div className="flex items-center gap-2">
                      <Sword className="h-4 w-4 text-yellow-500" />
                      <span className="text-xs text-muted-foreground">{stats.jobName}</span>
                    </div>
                  )}
                  {stats.mapName && (
                    <div className="flex items-center gap-2">
                      <MapPin className="h-4 w-4 text-green-500" />
                      <span className="text-xs text-muted-foreground">{stats.mapName}</span>
                    </div>
                  )}
                </div>
              )}

              {/* HP Bar */}
              <div>
                <div className="flex justify-between text-sm mb-2">
                  <span>HP</span>
                  <span className="font-mono">{stats.hp} / {stats.maxHp} ({Math.round(hpPercent)}%)</span>
                </div>
                <div className="w-full bg-secondary rounded-full h-2">
                  <motion.div 
                    className="bg-red-500 h-2 rounded-full"
                    initial={{ width: 0 }}
                    animate={{ width: `${hpPercent}%` }}
                    transition={{ duration: 0.3 }}
                  />
                </div>
              </div>

              {/* SP Bar */}
              <div>
                <div className="flex justify-between text-sm mb-2">
                  <span>SP</span>
                  <span className="font-mono">{stats.sp} / {stats.maxSp} ({Math.round(spPercent)}%)</span>
                </div>
                <div className="w-full bg-secondary rounded-full h-2">
                  <motion.div 
                    className="bg-blue-500 h-2 rounded-full"
                    initial={{ width: 0 }}
                    animate={{ width: `${spPercent}%` }}
                    transition={{ duration: 0.3 }}
                  />
                </div>
              </div>

              {/* Main Control */}
              <div className="pt-4 space-y-2">
                <Button
                  onClick={onStartStop}
                  disabled={!isConnected}
                  className={`w-full h-12 ${
                    isRunning 
                      ? 'bg-red-600 hover:bg-red-700' 
                      : 'bg-green-600 hover:bg-green-700'
                  }`}
                >
                  {isRunning ? (
                    <>
                      <Pause className="h-4 w-4 mr-2" />
                      Stop Services
                    </>
                  ) : (
                    <>
                      <Play className="h-4 w-4 mr-2" />
                      Start Services
                    </>
                  )}
                </Button>
                
                {/* Debug Buttons */}
                <div className="space-y-2">
                  <Button
                    onClick={async () => {
                      try {
                        const { DebugCharacterInfo } = await import('../wailsjs/wailsjs/go/main/App')
                        const result = await DebugCharacterInfo()
                        if (result.success) {
                          setDebugInfo(result.data)
                          setShowDebug(true)
                        } else {
                          alert(`Error: ${result.message}`)
                        }
                      } catch (error) {
                        alert(`Error: ${error}`)
                      }
                    }}
                    disabled={!isConnected}
                    variant="outline"
                    className="w-full"
                  >
                    🔍 Debug Job ID
                  </Button>
                  
                  <Button
                    onClick={async () => {
                      try {
                        const { SearchJobIDInMemory } = await import('../wailsjs/wailsjs/go/main/App')
                        const result = await SearchJobIDInMemory()
                        if (result.success) {
                          setDebugInfo(result.data)
                          setShowDebug(true)
                          alert(`Búsqueda completada: ${result.message}`)
                        } else {
                          alert(`Error: ${result.message}`)
                        }
                      } catch (error) {
                        alert(`Error: ${error}`)
                      }
                    }}
                    disabled={!isConnected}
                    variant="outline"
                    className="w-full bg-blue-50 hover:bg-blue-100"
                  >
                    🔎 Buscar Job ID en Memoria
                  </Button>
                  
                  <Button
                    onClick={async () => {
                      try {
                        const { MonitorJobAddresses } = await import('../wailsjs/wailsjs/go/main/App')
                        const result = await MonitorJobAddresses()
                        if (result.success) {
                          setDebugInfo(result.data)
                          setShowDebug(true)
                          alert(`Monitoreo completado: ${result.message}`)
                        } else {
                          alert(`Error: ${result.message}`)
                        }
                      } catch (error) {
                        alert(`Error: ${error}`)
                      }
                    }}
                    disabled={!isConnected}
                    variant="outline"
                    className="w-full bg-purple-50 hover:bg-purple-100"
                  >
                    📊 Monitorear Direcciones
                  </Button>
                </div>
              </div>
            </CardContent>
          </Card>

          {/* System Performance */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Clock className="h-5 w-5" />
                System Performance
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-2 gap-4">
                <div className="metric-card">
                  <div className="text-2xl font-bold">{formatUptime(stats.uptime)}</div>
                  <div className="text-sm text-muted-foreground">Uptime</div>
                </div>
                
                <div className="metric-card">
                  <div className="text-2xl font-bold">
                    {stats.system ? `${stats.system.CPUUsage.toFixed(1)}%` : '0%'}
                  </div>
                  <div className="text-sm text-muted-foreground">CPU Usage</div>
                </div>
                
                <div className="metric-card">
                  <div className="text-2xl font-bold">
                    {stats.system ? formatBytes(stats.system.MemoryUsage) : '0 B'}
                  </div>
                  <div className="text-sm text-muted-foreground">RAM Usage</div>
                </div>
                
                <div className="metric-card">
                  <div className="text-2xl font-bold">
                    {stats.system ? stats.system.GoroutineCount : 0}
                  </div>
                  <div className="text-sm text-muted-foreground">Goroutines</div>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Message Display */}
        {message && (
          <Card>
            <CardContent className="pt-6">
              <div className="flex items-center gap-2 text-sm">
                <MousePointer className="h-4 w-4" />
                <span>{message}</span>
              </div>
            </CardContent>
          </Card>
        )}

        {/* Debug Information */}
        {showDebug && debugInfo && (
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center justify-between">
                <span>🔍 Debug Information</span>
                <Button
                  onClick={() => setShowDebug(false)}
                  variant="outline"
                  size="sm"
                >
                  Cerrar
                </Button>
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                <div className="bg-muted p-4 rounded-lg">
                  <h4 className="font-medium mb-2">Character Name:</h4>
                  <p className="text-sm font-mono">{debugInfo.character_name || 'No encontrado'}</p>
                </div>
                
                <div className="bg-muted p-4 rounded-lg">
                  <h4 className="font-medium mb-2">Job ID Analysis:</h4>
                  <div className="space-y-1 max-h-96 overflow-y-auto border rounded p-2 bg-background">
                    {Object.entries(debugInfo)
                      .filter(([key]) => key.startsWith('offset_') || key.startsWith('near_') || key.startsWith('heap_') || key.startsWith('addr_'))
                      .sort(([a], [b]) => {
                        // Ordenar por tipo: addr_ (monitoreo) primero, luego otros
                        if (a.startsWith('addr_') && !b.startsWith('addr_')) return -1
                        if (!a.startsWith('addr_') && b.startsWith('addr_')) return 1
                        if (a.startsWith('near_character_') && !b.startsWith('near_character_')) return -1
                        if (!a.startsWith('near_character_') && b.startsWith('near_character_')) return 1
                        if (a.startsWith('near_hp_') && !b.startsWith('near_hp_')) return -1
                        if (!a.startsWith('near_hp_') && b.startsWith('near_hp_')) return 1
                        return a.localeCompare(b)
                      })
                      .map(([key, value]: [string, any]) => {
                        const isSearchResult = key.startsWith('near_') || key.startsWith('heap_')
                        const isMonitorResult = key.startsWith('addr_')
                        const isChanged = value.is_changed
                        
                        return (
                          <div key={key} className={`flex justify-between items-center text-xs border-b pb-1 hover:bg-muted/50 p-1 rounded ${
                            isMonitorResult ? (isChanged ? 'bg-red-50 border-red-200' : 'bg-purple-50 border-purple-200') : 
                            isSearchResult ? 'bg-green-50 border-green-200' : ''
                          }`}>
                            <div className="flex flex-col">
                              <span className={`font-mono ${
                                isMonitorResult ? (isChanged ? 'text-red-700 font-bold' : 'text-purple-700 font-bold') :
                                isSearchResult ? 'text-green-700 font-bold' : 'text-blue-600'
                              }`}>
                                {key}
                              </span>
                              {(isSearchResult || isMonitorResult) && value.address && (
                                <span className="text-xs text-gray-600">Address: {value.address}</span>
                              )}
                              {isMonitorResult && value.hex_value && (
                                <span className="text-xs text-gray-600">Hex: {value.hex_value}</span>
                              )}
                            </div>
                            <div className="text-right">
                              <div className="font-mono">Value: {value.value}</div>
                              <div className={`font-medium ${
                                isMonitorResult && isChanged ? 'text-red-600 font-bold' :
                                value.value === 0 ? 'text-green-600' : 
                                value.value === 4019 ? 'text-purple-600' : 
                                value.job_name.includes('Unknown') ? 'text-gray-500' : 'text-primary'
                              }`}>
                                {value.job_name}
                                {isMonitorResult && isChanged && ' ⚠️ CHANGED!'}
                              </div>
                              {isSearchResult && value.offset_from_base !== undefined && (
                                <div className="text-xs text-green-600">Offset: {value.offset_from_base}</div>
                              )}
                            </div>
                          </div>
                        )
                      })}
                  </div>
                </div>
                
                <div className="text-xs text-muted-foreground bg-yellow-50 p-3 rounded border">
                  <p><strong>Instrucciones:</strong></p>
                  <p>1. <strong>Haz scroll</strong> en el área de arriba para ver todos los offsets</p>
                  <p>2. <strong>Direcciones Monitoreadas (púrpura):</strong> Las 7 direcciones encontradas para Creator</p>
                  <p>3. <strong>Cambia de job</strong> en el juego y presiona &quot;📊 Monitorear Direcciones&quot;</p>
                  <p>4. <strong>Busca la dirección que cambió</strong> (aparecerá en rojo con &quot;⚠️ CHANGED!&quot;)</p>
                  <p>5. <strong>Esa es la dirección correcta</strong> del job ID</p>
                  <p>6. <strong>Envíame la dirección</strong> que cambió para actualizar el código</p>
                </div>
              </div>
            </CardContent>
          </Card>
        )}
      </div>
    </motion.div>
  )
}
