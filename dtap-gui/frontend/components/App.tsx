import { useState, useEffect } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import Navigation from './Navigation'
import Main from './Main'
import Cheats from './Cheats'
import Toast from './Toast'
import { ThemeProvider } from '../contexts/ThemeContext'
import { TransparencyProvider } from '../contexts/TransparencyContext'
import { 
  ConnectToRO, 
  StartServices, 
  StopServices, 
  PauseServices,
  ResumeServices,
  GetStats, 
  GetConfig, 
  UpdateConfig 
} from '../wailsjs/wailsjs/go/main/App'

type Route = 'main' | 'cheats'

export default function App() {
  const [currentRoute, setCurrentRoute] = useState<Route>('main')
  const [isRunning, setIsRunning] = useState(false)
  const [isPaused, setIsPaused] = useState(false)
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
  const [keyClickRatio, setKeyClickRatio] = useState(3)
  const [clickDelay, setClickDelay] = useState(50)
  const [predictionMode, setPredictionMode] = useState(true)
  // Nuevos campos configurables
  const [autopotDelay, setAutopotDelay] = useState(100)
  const [autobuffDelay, setAutobuffDelay] = useState(5000)
  const [mouseClick, setMouseClick] = useState(true)
  const [spammerMode, setSpammerMode] = useState('ultra_fast')
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
    characterName: '',
    mapName: '',
    jobName: '',
    system: {
      MemoryUsage: 0,
      MemoryTotal: 0,
      MemoryPercent: 0,
      GoMemoryUsage: 0,
      GoroutineCount: 0,
      CPUUsage: 0
    },
    autopot_prediction: {
      prediction_mode: false,
      continuous_spam_mode: false,
      damage_rate: 0,
      high_damage_threshold: -15,
      spam_remaining_time: 0,
      spam_total_time: 5000,
      spam_elapsed_time: 0,
      prediction_duration: 0
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
            characterName: statsData.character_name || prev.characterName,
            mapName: statsData.map_name || prev.mapName,
            jobName: statsData.job_name || prev.jobName,
            system: statsData.system || prev.system,
            autopot_prediction: statsData.autopot_prediction || prev.autopot_prediction
          }))
          setIsRunning(statsData.is_running || false)
          setIsPaused(statsData.is_paused || false)
        } catch (error) {
          // Error obteniendo estadísticas - silenciado
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
        setPredictionMode(config.autopot.prediction_mode !== false) // Por defecto true
        setAutopotDelay(config.autopot.check_delay_ms || 100)
      }
      if (config.autobuff) {
        setAutoBuffEnabled(config.autobuff.enabled)
        setGloomKey(config.autobuff.gloom_key || 'F3')
        setAspdKey(config.autobuff.aspd_key || 'F4')
        setAutobuffDelay(config.autobuff.delay_ms || 5000)
      }
      if (config.spammer) {
        setSpammerEnabled(config.spammer.enabled)
        setSpammerDelay(config.spammer.delay)
        setSpammerKeys(config.spammer.keys || ['F5', 'F6', 'F7'])
        setKeyClickRatio(config.spammer.key_click_ratio || 3)
        setClickDelay(config.spammer.click_delay_ms || 50)
        setMouseClick(config.spammer.mouse_click !== false) // Por defecto true
        setSpammerMode(config.spammer.mode || 'ultra_fast')
      }
    } catch (error) {
      // Error cargando configuración - silenciado
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
      if (!isRunning) {
        // Iniciar servicios
        const result = await StartServices()
        setMessage(result.message)
      } else if (isPaused) {
        // Reanudar servicios
        const result = await ResumeServices()
        setMessage(result.message)
      } else {
        // Pausar servicios
        const result = await PauseServices()
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
          check_delay_ms: autopotDelay,
          prediction_mode: predictionMode
        },
        autobuff: {
          enabled: autoBuffEnabled,
          delay_ms: autobuffDelay,
          gloom_key: gloomKey,
          aspd_key: aspdKey
        },
        spammer: {
          enabled: spammerEnabled,
          keys: spammerKeys,
          delay: spammerDelay,
          mouse_click: mouseClick,
          mode: spammerMode,
          key_click_ratio: keyClickRatio,
          click_delay_ms: clickDelay
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
    <ThemeProvider>
      <TransparencyProvider>
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
                isPaused={isPaused}
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
              keyClickRatio={keyClickRatio}
              clickDelay={clickDelay}
              predictionMode={predictionMode}
              autopotDelay={autopotDelay}
              autobuffDelay={autobuffDelay}
              mouseClick={mouseClick}
              spammerMode={spammerMode}
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
              onKeyClickRatioChange={setKeyClickRatio}
              onClickDelayChange={setClickDelay}
              onPredictionModeChange={setPredictionMode}
              onAutopotDelayChange={setAutopotDelay}
              onAutobuffDelayChange={setAutobuffDelay}
              onMouseClickChange={setMouseClick}
              onSpammerModeChange={setSpammerMode}
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
      </TransparencyProvider>
    </ThemeProvider>
  )
}