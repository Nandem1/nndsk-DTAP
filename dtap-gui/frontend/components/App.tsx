import { useState, useEffect } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import Navigation from './Navigation'
import Main from './Main'
import Cheats from './Cheats'
import Toast from './Toast'
import { 
  ConnectToRO, 
  StartServices, 
  StopServices, 
  GetStats, 
  GetConfig, 
  UpdateConfig 
} from '../wailsjs/wailsjs/go/main/App'

type Route = 'main' | 'cheats'

export default function App() {
  const [currentRoute, setCurrentRoute] = useState<Route>('main')
  const [isRunning, setIsRunning] = useState(false)
  const [isConnected, setIsConnected] = useState(false)
  const [autoPotEnabled, setAutoPotEnabled] = useState(true)
  const [autoBuffEnabled, setAutoBuffEnabled] = useState(false)
  const [spammerEnabled, setSpammerEnabled] = useState(false)
  const [hpThreshold, setHpThreshold] = useState(50)
  const [spThreshold, setSpThreshold] = useState(30)
  const [spammerDelay, setSpammerDelay] = useState(1.0)
  const [gloomKey, setGloomKey] = useState('F3')
  const [aspdKey, setAspdKey] = useState('F4')
  const [hpPotKey, setHpPotKey] = useState('F1')
  const [spPotKey, setSpPotKey] = useState('F2')
  const [spammerKeys, setSpammerKeys] = useState<string[]>(['F5', 'F6', 'F7'])
  const [toast, setToast] = useState<{
    isVisible: boolean
    message: string
    type: 'success' | 'error' | 'warning'
  }>({
    isVisible: false,
    message: '',
    type: 'success'
  })
  const [stats, setStats] = useState({
    hp: 100,
    maxHp: 100,
    sp: 100,
    maxSp: 100,
    uptime: '00:00:00',
    totalActions: 0,
    system: {
      MemoryUsage: 0,
      MemoryTotal: 0,
      MemoryPercent: 0,
      GoMemoryUsage: 0,
      GoroutineCount: 0,
      CPUUsage: 0
    }
  })
  const [message, setMessage] = useState('')

  // Función para mostrar toast
  const showToast = (message: string, type: 'success' | 'error' | 'warning') => {
    setToast({ isVisible: true, message, type })
    setTimeout(() => {
      setToast(prev => ({ ...prev, isVisible: false }))
    }, 3000)
  }

  // Cargar configuración al iniciar
  useEffect(() => {
    loadConfig()
  }, [])

  // Actualizar estadísticas en tiempo real
  useEffect(() => {
    const interval = setInterval(async () => {
      if (isConnected) {
        try {
          const statsData = await GetStats()
          setStats(prev => ({
            ...prev,
            hp: statsData.hp || prev.hp,
            maxHp: statsData.max_hp || prev.maxHp,
            sp: statsData.sp || prev.sp,
            maxSp: statsData.max_sp || prev.maxSp,
            uptime: statsData.uptime || prev.uptime,
            totalActions: statsData.total_actions || prev.totalActions,
            system: statsData.system || prev.system
          }))
          setIsRunning(statsData.is_running || false)
        } catch (error) {
          console.error('Error obteniendo estadísticas:', error)
        }
      }
    }, 16) // ~60 FPS (1000ms / 60 = 16.67ms)
    return () => clearInterval(interval)
  }, [isConnected])

  const loadConfig = async () => {
    try {
      const config = await GetConfig()
      if (config.autopot) {
        setAutoPotEnabled(config.autopot.enabled)
        setHpThreshold(config.autopot.hp_percent)
        setSpThreshold(config.autopot.sp_percent)
        setHpPotKey(config.autopot.hp_pot_key || 'F1')
        setSpPotKey(config.autopot.sp_pot_key || 'F2')
      }
      if (config.autobuff) {
        setAutoBuffEnabled(config.autobuff.enabled)
        setGloomKey(config.autobuff.gloom_key || 'F3')
        setAspdKey(config.autobuff.aspd_key || 'F4')
      }
      if (config.spammer) {
        setSpammerEnabled(config.spammer.enabled)
        setSpammerDelay(config.spammer.delay)
        setSpammerKeys(config.spammer.keys || ['F5', 'F6', 'F7'])
      }
    } catch (error) {
      console.error('Error cargando configuración:', error)
    }
  }

  const handleConnect = async () => {
    try {
      const result = await ConnectToRO()
      if (result.success) {
        setIsConnected(true)
        setMessage(result.message)
      } else {
        setMessage(result.message)
      }
    } catch (error) {
      setMessage(`Error conectando: ${error}`)
    }
  }

  const handleStartStop = async () => {
    try {
      if (isRunning) {
        const result = await StopServices()
        setMessage(result.message)
      } else {
        const result = await StartServices()
        setMessage(result.message)
      }
    } catch (error) {
      setMessage(`Error: ${error}`)
    }
  }

  const handleConfigUpdate = async () => {
    try {
      const configData = {
        autopot: {
          enabled: autoPotEnabled,
          hp_pot_key: hpPotKey,
          sp_pot_key: spPotKey,
          hp_percent: hpThreshold,
          sp_percent: spThreshold,
          check_delay: 100
        },
        autobuff: {
          enabled: autoBuffEnabled,
          delay: 5000,
          gloom_key: gloomKey,
          aspd_key: aspdKey
        },
        spammer: {
          enabled: spammerEnabled,
          keys: spammerKeys,
          delay: spammerDelay,
          skill_cancel: false,
          double_click: false,
          mouse_click: true,
          mode: "simple"
        },
        ui: {
          show_debug: false,
          log_level: "info"
        }
      }
      
      const result = await UpdateConfig(JSON.stringify(configData))
      if (result.success) {
        showToast('Configuración guardada correctamente', 'success')
      } else {
        showToast(result.message || 'Error guardando configuración', 'error')
      }
    } catch (error) {
      showToast(`Error actualizando configuración: ${error}`, 'error')
    }
  }

  return (
    <div className="h-screen w-screen overflow-hidden flex flex-col">
      {/* Navigation */}
      <Navigation 
        currentRoute={currentRoute} 
        onRouteChange={setCurrentRoute} 
      />

      {/* Main Content */}
      <div className="flex-1 overflow-hidden">
        <AnimatePresence mode="wait">
          {currentRoute === 'main' && (
            <Main
              key="main"
              isRunning={isRunning}
              isConnected={isConnected}
              stats={stats}
              message={message}
              onConnect={handleConnect}
              onStartStop={handleStartStop}
            />
          )}
          
          {currentRoute === 'cheats' && (
            <Cheats
              key="cheats"
              autoPotEnabled={autoPotEnabled}
              autoBuffEnabled={autoBuffEnabled}
              spammerEnabled={spammerEnabled}
              hpThreshold={hpThreshold}
              spThreshold={spThreshold}
              spammerDelay={spammerDelay}
              gloomKey={gloomKey}
              aspdKey={aspdKey}
              hpPotKey={hpPotKey}
              spPotKey={spPotKey}
              spammerKeys={spammerKeys}
              onAutoPotChange={setAutoPotEnabled}
              onAutoBuffChange={setAutoBuffEnabled}
              onSpammerChange={setSpammerEnabled}
              onHpThresholdChange={setHpThreshold}
              onSpThresholdChange={setSpThreshold}
              onSpammerDelayChange={setSpammerDelay}
              onGloomKeyChange={setGloomKey}
              onAspdKeyChange={setAspdKey}
              onHpPotKeyChange={setHpPotKey}
              onSpPotKeyChange={setSpPotKey}
              onSpammerKeysChange={setSpammerKeys}
              onSaveConfig={handleConfigUpdate}
            />
          )}
          
        </AnimatePresence>
      </div>
      
      {/* Toast */}
      <Toast
        isVisible={toast.isVisible}
        message={toast.message}
        type={toast.type}
        onClose={() => setToast(prev => ({ ...prev, isVisible: false }))}
      />
    </div>
  )
}