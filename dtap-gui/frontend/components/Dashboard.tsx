import { motion } from 'framer-motion'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Switch } from '@/components/ui/switch'
import { Label } from '@/components/ui/label'
import { 
  Play, 
  Pause, 
  Zap, 
  Shield, 
  Target,
  Activity,
  Wifi,
  WifiOff,
  User,
  MapPin,
  Sword
} from 'lucide-react'

interface DashboardProps {
  isRunning: boolean
  isConnected: boolean
  autoPotEnabled: boolean
  autoBuffEnabled: boolean
  spammerEnabled: boolean
  stats: {
    hp: number
    maxHp: number
    sp: number
    maxSp: number
    uptime: string
    totalActions: number
    characterName?: string
    mapName?: string
    jobName?: string
  }
  message: string
  onConnect: () => void
  onStartStop: () => void
  onAutoPotChange: (enabled: boolean) => void
  onAutoBuffChange: (enabled: boolean) => void
  onSpammerChange: (enabled: boolean) => void
  onOpenSettings: () => void
}

export default function Dashboard({
  isRunning,
  isConnected,
  autoPotEnabled,
  autoBuffEnabled,
  spammerEnabled,
  stats,
  message,
  onConnect,
  onStartStop,
  onAutoPotChange,
  onAutoBuffChange,
  onSpammerChange,
  onOpenSettings
}: DashboardProps) {
  const hpPercent = (stats.hp / stats.maxHp) * 100
  const spPercent = (stats.sp / stats.maxSp) * 100

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      exit={{ opacity: 0, y: -20 }}
      transition={{ duration: 0.3 }}
      className="h-full flex flex-col"
    >
      {/* Connection Status */}
      <motion.div
        initial={{ opacity: 0, y: -10 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.1 }}
        className="mb-8"
      >
        <div className="vercel-card p-6">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-4">
              <div className={`status-indicator ${isConnected ? 'connected' : 'disconnected'}`}></div>
              <div>
                <span className="text-foreground text-sm font-medium">
                  {isConnected ? 'Connected' : 'Disconnected'}
                </span>
                <div className="text-xs text-muted-foreground">
                  {isConnected ? 'Ragnarok Online' : 'No Process Found'}
                </div>
              </div>
            </div>
            <Button
              onClick={onConnect}
              disabled={isConnected}
              className="vercel-button"
            >
              {isConnected ? 'Connected' : 'Connect'}
            </Button>
          </div>
          {message && (
            <div className="mt-4 p-3 bg-muted/50 rounded-lg border-l-2 border-primary">
              <div className="text-xs text-foreground">
                {message}
              </div>
            </div>
          )}
        </div>
      </motion.div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6 flex-1">
        {/* Status Panel */}
        <motion.div
          initial={{ opacity: 0, x: -20 }}
          animate={{ opacity: 1, x: 0 }}
          transition={{ delay: 0.2 }}
        >
          <div className="vercel-card p-6 h-full">
            <div className="flex items-center gap-3 mb-6">
              <Activity className="h-5 w-5 text-primary" />
              <h3 className="text-foreground text-sm font-medium">Status</h3>
            </div>
            
            <div className="space-y-6">
              {/* Character Info */}
              {stats.characterName && (
                <div className="space-y-3 pb-4 border-b border-border">
                  <div className="flex items-center gap-2">
                    <User className="h-4 w-4 text-primary" />
                    <span className="text-sm font-medium text-foreground">{stats.characterName}</span>
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
                <div className="flex justify-between text-xs text-muted-foreground mb-2">
                  <span>HP</span>
                  <span className="text-foreground font-medium">{Math.round(hpPercent)}%</span>
                </div>
                <div className="w-full bg-muted rounded-full h-2">
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
                <div className="flex justify-between text-xs text-muted-foreground mb-2">
                  <span>SP</span>
                  <span className="text-foreground font-medium">{Math.round(spPercent)}%</span>
                </div>
                <div className="w-full bg-muted rounded-full h-2">
                  <motion.div 
                    className="sp-bar h-2 rounded-full"
                    initial={{ width: 0 }}
                    animate={{ width: `${spPercent}%` }}
                    transition={{ duration: 0.3 }}
                  />
                </div>
              </div>

              {/* Stats */}
              <div className="grid grid-cols-2 gap-4 pt-4 border-t border-border">
                <div className="text-center">
                  <div className="text-2xl font-bold text-foreground">{stats.uptime}</div>
                  <div className="text-xs text-muted-foreground">Uptime</div>
                </div>
                <div className="text-center">
                  <div className="text-2xl font-bold text-foreground">{stats.totalActions}</div>
                  <div className="text-xs text-muted-foreground">Actions</div>
                </div>
              </div>
            </div>
          </div>
        </motion.div>

        {/* Controls Panel */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.3 }}
        >
          <div className="vercel-card p-6 h-full">
            <div className="flex items-center gap-3 mb-6">
              <Play className="h-5 w-5 text-primary" />
              <h3 className="text-foreground text-sm font-medium">Controls</h3>
            </div>
            
            <div className="space-y-6">
              {/* Main Control */}
              <div className="text-center">
                <Button
                  onClick={onStartStop}
                  disabled={!isConnected}
                  className={`w-full h-12 text-sm font-medium transition-all duration-200 ${
                    isRunning 
                      ? 'stop-button' 
                      : 'start-button'
                  } ${!isConnected ? 'opacity-50 cursor-not-allowed' : ''}`}
                >
                  {isRunning ? (
                    <>
                      <Pause className="h-4 w-4 mr-2" />
                      Stop
                    </>
                  ) : (
                    <>
                      <Play className="h-4 w-4 mr-2" />
                      Start
                    </>
                  )}
                </Button>
              </div>

              {/* Auto Pot */}
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-3">
                  <Zap className="h-4 w-4 text-warning" />
                  <Label className="text-foreground text-sm">Auto Pot</Label>
                </div>
                <Switch
                  checked={autoPotEnabled}
                  onCheckedChange={onAutoPotChange}
                  className="vercel-switch"
                />
              </div>

              {/* Auto Buff */}
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-3">
                  <Shield className="h-4 w-4 text-info" />
                  <Label className="text-foreground text-sm">Auto Buff</Label>
                </div>
                <Switch
                  checked={autoBuffEnabled}
                  onCheckedChange={onAutoBuffChange}
                  className="vercel-switch"
                />
              </div>

              {/* Spammer */}
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-3">
                  <Target className="h-4 w-4 text-error" />
                  <Label className="text-foreground text-sm">Spammer</Label>
                </div>
                <Switch
                  checked={spammerEnabled}
                  onCheckedChange={onSpammerChange}
                  className="vercel-switch"
                />
              </div>
            </div>
          </div>
        </motion.div>

        {/* Settings Panel */}
        <motion.div
          initial={{ opacity: 0, x: 20 }}
          animate={{ opacity: 1, x: 0 }}
          transition={{ delay: 0.4 }}
        >
          <div className="vercel-card p-6 h-full">
            <div className="flex items-center gap-3 mb-6">
              <Activity className="h-5 w-5 text-primary" />
              <h3 className="text-foreground text-sm font-medium">Quick Settings</h3>
            </div>
            
            <div className="space-y-4">
              <div className="text-center">
                <Button
                  onClick={onOpenSettings}
                  className="vercel-button w-full"
                >
                  Advanced Settings
                </Button>
              </div>
              
              <div className="text-center text-xs text-muted-foreground">
                Configure thresholds, delays, and advanced options
              </div>
            </div>
          </div>
        </motion.div>
      </div>
    </motion.div>
  )
}
