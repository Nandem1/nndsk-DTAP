import { motion } from 'framer-motion'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { 
  Activity, 
  Wifi, 
  WifiOff, 
  Play, 
  Pause,
  Clock,
  MousePointer
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
      className="h-full flex flex-col"
    >
      <div className="flex-1 p-6 space-y-6">
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
              <div className="pt-4">
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
      </div>
    </motion.div>
  )
}
