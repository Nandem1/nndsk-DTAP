import { motion } from 'framer-motion'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Switch } from '@/components/ui/switch'
import { Label } from '@/components/ui/label'
import { Input } from '@/components/ui/input'
import { useTransparency } from '../contexts/TransparencyContext'
import { CardDecorations } from './themed-icons'
import { 
  Sword, 
  Shield, 
  Zap,
  Settings,
  Plus,
  Trash2,
  Save
} from 'lucide-react'
import { useState } from 'react'

interface ATKDEFProps {
  enabled: boolean
  activationKey: string
  switchDelay: number
  spamDelay: number
  mouseClick: boolean
  atkKeys: Record<string, string>
  defKeys: Record<string, string>
  onEnabledChange: (enabled: boolean) => void
  onActivationKeyChange: (key: string) => void
  onSwitchDelayChange: (delay: number) => void
  onSpamDelayChange: (delay: number) => void
  onMouseClickChange: (enabled: boolean) => void
  onATKKeyAdd: (name: string, key: string) => void
  onATKKeyRemove: (name: string) => void
  onDEFKeyAdd: (name: string, key: string) => void
  onDEFKeyRemove: (name: string) => void
  onSaveConfig: () => void
}

export default function ATKDEF({
  enabled,
  activationKey,
  switchDelay,
  spamDelay,
  mouseClick,
  atkKeys,
  defKeys,
  onEnabledChange,
  onActivationKeyChange,
  onSwitchDelayChange,
  onSpamDelayChange,
  onMouseClickChange,
  onATKKeyAdd,
  onATKKeyRemove,
  onDEFKeyAdd,
  onDEFKeyRemove,
  onSaveConfig
}: ATKDEFProps) {
  const { transparencyMode } = useTransparency()
  const [newATKName, setNewATKName] = useState('')
  const [newATKKey, setNewATKKey] = useState('')
  const [newDEFName, setNewDEFName] = useState('')
  const [newDEFKey, setNewDEFKey] = useState('')

  const handleATKKeyAdd = () => {
    if (newATKName.trim() && newATKKey.trim()) {
      onATKKeyAdd(newATKName.trim(), newATKKey.trim())
      setNewATKName('')
      setNewATKKey('')
    }
  }

  const handleDEFKeyAdd = () => {
    if (newDEFName.trim() && newDEFKey.trim()) {
      onDEFKeyAdd(newDEFName.trim(), newDEFKey.trim())
      setNewDEFName('')
      setNewDEFKey('')
    }
  }

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      exit={{ opacity: 0, y: -20 }}
      transition={{ duration: 0.3 }}
      className={`h-full flex flex-col ${transparencyMode === 'hyprland' ? 'hyprland-background' : 'glass-background'}`}
    >
      <div className="flex-1 p-4 space-y-4">
        {/* Header */}
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-2xl font-bold flex items-center gap-2">
              <Sword className="h-6 w-6 text-primary" />
              ATK-DEF Mode
            </h1>
            <p className="text-muted-foreground">Configuración de cambio automático de equipamiento</p>
          </div>
        </div>

        {/* Main Configuration */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
          {/* General Settings */}
          <Card className={`relative ${transparencyMode === 'hyprland' ? 'hyprland-card' : 'glass-card'}`}>
            <CardDecorations isActive={enabled} />
            <CardHeader className="pb-3">
              <CardTitle className="flex items-center gap-2 text-sm pr-16">
                <Settings className="h-4 w-4 text-primary" />
                Configuración General
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="flex items-center justify-between">
                <Label>Habilitar ATK-DEF Mode</Label>
                <Switch
                  checked={enabled}
                  onCheckedChange={onEnabledChange}
                />
              </div>
              
              <div>
                <Label className="text-sm">Tecla de Activación</Label>
                <Input
                  value={activationKey}
                  onChange={(e) => onActivationKeyChange(e.target.value)}
                  placeholder="SPACE"
                  className="mt-1"
                />
                <p className="text-xs text-muted-foreground mt-1">
                  Tecla que activa el modo ATK-DEF mientras se mantiene presionada
                </p>
              </div>
              
              <div className="grid grid-cols-2 gap-3">
                <div>
                  <Label className="text-sm">Switch Delay (ms)</Label>
                  <Input
                    type="number"
                    value={switchDelay}
                    onChange={(e) => onSwitchDelayChange(Number(e.target.value))}
                    min="1"
                    max="1000"
                    className="mt-1"
                  />
                </div>
                <div>
                  <Label className="text-sm">Spam Delay (ms)</Label>
                  <Input
                    type="number"
                    step="0.1"
                    value={spamDelay}
                    onChange={(e) => onSpamDelayChange(Number(e.target.value))}
                    min="0.1"
                    max="100"
                    className="mt-1"
                  />
                </div>
              </div>
              
              <div className="flex items-center justify-between">
                <div>
                  <Label className="text-sm">Mouse Click</Label>
                  <p className="text-xs text-muted-foreground">Habilitar clicks del mouse durante el spam</p>
                </div>
                <Switch
                  checked={mouseClick}
                  onCheckedChange={onMouseClickChange}
                />
              </div>
            </CardContent>
          </Card>

          {/* ATK Keys */}
          <Card className={`relative ${transparencyMode === 'hyprland' ? 'hyprland-card' : 'glass-card'}`}>
            <CardDecorations isActive={Object.keys(atkKeys).length > 0} />
            <CardHeader className="pb-3">
              <CardTitle className="flex items-center gap-2 text-sm pr-16">
                <Sword className="h-4 w-4 text-error" />
                Teclas ATK (Ataque)
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-3">
              {/* Add new ATK key */}
              <div className="space-y-2 p-3 border border-border rounded-md">
                <div className="grid grid-cols-2 gap-2">
                  <Input
                    placeholder="Nombre del item"
                    value={newATKName}
                    onChange={(e) => setNewATKName(e.target.value)}
                    className="text-xs"
                  />
                  <Input
                    placeholder="Tecla (ej: F1)"
                    value={newATKKey}
                    onChange={(e) => setNewATKKey(e.target.value)}
                    className="text-xs"
                  />
                </div>
                <Button
                  onClick={handleATKKeyAdd}
                  size="sm"
                  className="w-full"
                  disabled={!newATKName.trim() || !newATKKey.trim()}
                >
                  <Plus className="h-3 w-3 mr-1" />
                  Agregar
                </Button>
              </div>

              {/* ATK Keys List */}
              <div className="space-y-2 max-h-40 overflow-y-auto">
                {Object.entries(atkKeys).map(([name, key]) => (
                  <div key={name} className="flex items-center justify-between p-2 bg-secondary rounded-md">
                    <div className="flex-1">
                      <div className="text-sm font-medium">{name}</div>
                      <div className="text-xs text-muted-foreground">{key}</div>
                    </div>
                    <Button
                      onClick={() => onATKKeyRemove(name)}
                      size="sm"
                      variant="ghost"
                      className="h-6 w-6 p-0 text-destructive hover:text-destructive"
                    >
                      <Trash2 className="h-3 w-3" />
                    </Button>
                  </div>
                ))}
                {Object.keys(atkKeys).length === 0 && (
                  <div className="text-center text-muted-foreground text-sm py-4">
                    No hay teclas ATK configuradas
                  </div>
                )}
              </div>
            </CardContent>
          </Card>
        </div>

        {/* DEF Keys */}
        <Card className={`relative ${transparencyMode === 'hyprland' ? 'hyprland-card' : 'glass-card'}`}>
          <CardDecorations isActive={Object.keys(defKeys).length > 0} />
          <CardHeader className="pb-3">
            <CardTitle className="flex items-center gap-2 text-sm pr-16">
              <Shield className="h-4 w-4 text-info" />
              Teclas DEF (Defensa)
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-3">
            {/* Add new DEF key */}
            <div className="space-y-2 p-3 border border-border rounded-md">
              <div className="grid grid-cols-2 gap-2">
                <Input
                  placeholder="Nombre del item"
                  value={newDEFName}
                  onChange={(e) => setNewDEFName(e.target.value)}
                  className="text-xs"
                />
                <Input
                  placeholder="Tecla (ej: F2)"
                  value={newDEFKey}
                  onChange={(e) => setNewDEFKey(e.target.value)}
                  className="text-xs"
                />
              </div>
              <Button
                onClick={handleDEFKeyAdd}
                size="sm"
                className="w-full"
                disabled={!newDEFName.trim() || !newDEFKey.trim()}
              >
                <Plus className="h-3 w-3 mr-1" />
                Agregar
              </Button>
            </div>

            {/* DEF Keys List */}
            <div className="space-y-2 max-h-40 overflow-y-auto">
              {Object.entries(defKeys).map(([name, key]) => (
                <div key={name} className="flex items-center justify-between p-2 bg-secondary rounded-md">
                  <div className="flex-1">
                    <div className="text-sm font-medium">{name}</div>
                    <div className="text-xs text-muted-foreground">{key}</div>
                  </div>
                  <Button
                    onClick={() => onDEFKeyRemove(name)}
                    size="sm"
                    variant="ghost"
                    className="h-6 w-6 p-0 text-destructive hover:text-destructive"
                  >
                    <Trash2 className="h-3 w-3" />
                  </Button>
                </div>
              ))}
              {Object.keys(defKeys).length === 0 && (
                <div className="text-center text-muted-foreground text-sm py-4">
                  No hay teclas DEF configuradas
                </div>
              )}
            </div>
          </CardContent>
        </Card>
      </div>
      
      {/* Fixed Save Button */}
      <div className="border-t bg-background p-4">
        <div className="flex justify-end">
          <Button onClick={onSaveConfig} className="px-6">
            <Save className="h-4 w-4 mr-2" />
            Guardar Configuración
          </Button>
        </div>
      </div>
    </motion.div>
  )
}
