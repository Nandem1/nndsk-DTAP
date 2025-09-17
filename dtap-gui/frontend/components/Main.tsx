import { motion } from 'framer-motion'
import { useState } from 'react'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { useTransparency } from '../contexts/TransparencyContext'
import { CardDecorations } from './themed-icons'
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
  isPaused: boolean
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
      GoMemoryTotal: number
      GoMemoryPercent: number
      ProcessMemoryUsage: number
      ProcessCPUUsage: number
      GoroutineCount: number
      CPUUsage: number
    }
    autopot_prediction?: {
      prediction_mode: boolean
      continuous_spam_mode: boolean
      damage_rate: number
      high_damage_threshold: number
      spam_remaining_time: number
      spam_total_time: number
      spam_elapsed_time: number
      prediction_duration: number
    }
  }
  message: string
  onConnect: () => void
  onStartStop: () => void
}

export default function Main({
  isRunning,
  isPaused,
  isConnected,
  stats,
  message,
  onConnect,
  onStartStop
}: MainProps) {
  const { transparencyMode } = useTransparency()
  
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
      className={`h-full flex flex-col overflow-hidden ${transparencyMode === 'hyprland' ? 'hyprland-background' : 'glass-background'}`}
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
          <Card className={`relative ${transparencyMode === 'hyprland' ? 'hyprland-card' : 'glass-card'}`}>
            <CardDecorations isActive={isRunning} />
            <CardHeader className="pb-4">
              <CardTitle className="flex items-center gap-2 pr-16">
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
                      <Sword className="h-4 w-4 job-icon" />
                      <span className="text-xs text-muted-foreground">{stats.jobName}</span>
                    </div>
                  )}
                  {stats.mapName && (
                    <div className="flex items-center gap-2">
                      <MapPin className="h-4 w-4 map-icon" />
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
                    className="hp-bar h-2 rounded-full"
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
                    className="sp-bar h-2 rounded-full"
                    initial={{ width: 0 }}
                    animate={{ width: `${spPercent}%` }}
                    transition={{ duration: 0.3 }}
                  />
                </div>
              </div>

              {/* AutoPot Prediction Status */}
              {stats.autopot_prediction && (
                <div className="pt-4 border-t border-border">
                  <div className="space-y-2">
                    <div className="flex items-center justify-between">
                      <span className="text-sm font-medium">AutoPot Status</span>
                      <div className="flex items-center gap-2">
                        {stats.autopot_prediction.continuous_spam_mode ? (
                          <motion.div
                            className="flex items-center gap-1 px-2 py-1 spam-mode-badge rounded-md text-xs"
                            animate={{ opacity: [1, 0.5, 1] }}
                            transition={{ duration: 1, repeat: Infinity }}
                          >
                            <div className="w-2 h-2 rounded-full" style={{ backgroundColor: 'hsl(var(--spam-mode))' }}></div>
                            SPAM CONTINUO
                          </motion.div>
                        ) : stats.autopot_prediction.prediction_mode ? (
                          <motion.div
                            className="flex items-center gap-1 px-2 py-1 prediction-mode-badge rounded-md text-xs"
                            animate={{ opacity: [1, 0.7, 1] }}
                            transition={{ duration: 2, repeat: Infinity }}
                          >
                            <div className="w-2 h-2 rounded-full" style={{ backgroundColor: 'hsl(var(--prediction-mode))' }}></div>
                            PREDICCIÓN
                          </motion.div>
                        ) : (
                          <div className="flex items-center gap-1 px-2 py-1 normal-mode-badge rounded-md text-xs">
                            <div className="w-2 h-2 rounded-full" style={{ backgroundColor: 'hsl(var(--normal-mode))' }}></div>
                            NORMAL
                          </div>
                        )}
                      </div>
                    </div>
                    
                    {stats.autopot_prediction.continuous_spam_mode && (
                      <div className="text-xs text-muted-foreground">
                        <div>Daño: {stats.autopot_prediction.damage_rate.toFixed(1)}%/seg</div>
                        <div>Tiempo restante: {(stats.autopot_prediction.spam_remaining_time / 1000).toFixed(1)}s</div>
                      </div>
                    )}
                    
                    {stats.autopot_prediction.prediction_mode && !stats.autopot_prediction.continuous_spam_mode && (
                      <div className="text-xs text-muted-foreground">
                        <div>Predicción activa: {(stats.autopot_prediction.prediction_duration / 1000).toFixed(1)}s</div>
                      </div>
                    )}
                  </div>
                </div>
              )}

              {/* Main Control */}
              <div className="pt-4 space-y-2">
                <Button
                  onClick={onStartStop}
                  disabled={!isConnected}
                  className={`w-full h-12 ${
                    !isRunning 
                      ? 'start-button' 
                      : isPaused
                      ? 'pause-button'
                      : 'stop-button'
                  }`}
                >
                  {!isRunning ? (
                    <>
                      <Play className="h-4 w-4 mr-2" />
                      Start Services
                    </>
                  ) : isPaused ? (
                    <>
                      <Play className="h-4 w-4 mr-2" />
                      Resume Services
                    </>
                  ) : (
                    <>
                      <Pause className="h-4 w-4 mr-2" />
                      Pause Services
                    </>
                  )}
                </Button>
                
              </div>
            </CardContent>
          </Card>

          {/* System Performance */}
          <Card className={`relative ${transparencyMode === 'hyprland' ? 'hyprland-card' : 'glass-card'}`}>
            <CardDecorations isActive={isConnected} />
            <CardHeader className="pb-4">
              <CardTitle className="flex items-center gap-2 pr-16">
                <Clock className="h-5 w-5" />
                System Performance
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                <div className="grid grid-cols-2 gap-3">
                  <div className="metric-card">
                    <div className="text-2xl font-bold">
                      {stats.system ? `${stats.system.CPUUsage.toFixed(1)}%` : '0%'}
                    </div>
                    <div className="text-sm text-muted-foreground">Process CPU</div>
                  </div>
                  <div className="metric-card">
                    <div className="text-2xl font-bold">
                      {stats.system ? formatBytes(stats.system.ProcessMemoryUsage) : '0 B'}
                    </div>
                    <div className="text-sm text-muted-foreground">Process RAM</div>
                  </div>
                </div>
                <div className="metric-card">
                  <div className="text-xl font-bold">{formatUptime(stats.uptime)}</div>
                  <div className="text-xs text-muted-foreground">Uptime</div>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Message Display */}
        {message && (
          <Card className={`relative ${transparencyMode === 'hyprland' ? 'hyprland-card' : 'glass-card'}`}>
            <CardDecorations isActive={true} />
            <CardContent className="pt-6 pr-16">
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
