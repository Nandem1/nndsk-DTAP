import { motion } from 'framer-motion'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Switch } from '@/components/ui/switch'
import { Label } from '@/components/ui/label'
import { 
  Zap, 
  Shield, 
  Target,
  Save
} from 'lucide-react'

interface CheatsProps {
  autoPotEnabled: boolean
  autoBuffEnabled: boolean
  spammerEnabled: boolean
  hpThreshold: number
  spThreshold: number
  spammerDelay: number
  gloomKey: string
  aspdKey: string
  hpPotKey: string
  spPotKey: string
  spammerKeys: string[]
  // Nuevas props para control de ratio teclas/clicks
  keyClickRatio: number
  clickDelay: number
  // Nueva prop para predicción
  predictionMode: boolean
  // Nuevos campos configurables
  autopotDelay: number
  autobuffDelay: number
  mouseClick: boolean
  spammerMode: string
  onAutoPotChange: (enabled: boolean) => void
  onAutoBuffChange: (enabled: boolean) => void
  onSpammerChange: (enabled: boolean) => void
  onHpThresholdChange: (value: number) => void
  onSpThresholdChange: (value: number) => void
  onSpammerDelayChange: (value: number) => void
  onGloomKeyChange: (key: string) => void
  onAspdKeyChange: (key: string) => void
  onHpPotKeyChange: (key: string) => void
  onSpPotKeyChange: (key: string) => void
  onSpammerKeysChange: (keys: string[]) => void
  onKeyClickRatioChange: (ratio: number) => void
  onClickDelayChange: (delay: number) => void
  onPredictionModeChange: (enabled: boolean) => void
  onAutopotDelayChange: (delay: number) => void
  onAutobuffDelayChange: (delay: number) => void
  onMouseClickChange: (enabled: boolean) => void
  onSpammerModeChange: (mode: string) => void
  onSaveConfig: () => void
}

export default function Cheats({
  autoPotEnabled,
  autoBuffEnabled,
  spammerEnabled,
  hpThreshold,
  spThreshold,
  spammerDelay,
  gloomKey,
  aspdKey,
  hpPotKey,
  spPotKey,
  spammerKeys,
  keyClickRatio,
  clickDelay,
  predictionMode,
  autopotDelay,
  autobuffDelay,
  mouseClick,
  spammerMode,
  onAutoPotChange,
  onAutoBuffChange,
  onSpammerChange,
  onHpThresholdChange,
  onSpThresholdChange,
  onSpammerDelayChange,
  onGloomKeyChange,
  onAspdKeyChange,
  onHpPotKeyChange,
  onSpPotKeyChange,
  onSpammerKeysChange,
  onKeyClickRatioChange,
  onClickDelayChange,
  onPredictionModeChange,
  onAutopotDelayChange,
  onAutobuffDelayChange,
  onMouseClickChange,
  onSpammerModeChange,
  onSaveConfig
}: CheatsProps) {
  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      exit={{ opacity: 0, y: -20 }}
      transition={{ duration: 0.3 }}
      className="h-full flex flex-col"
    >
      <div className="flex-1 p-4 space-y-4">
        {/* Cheats Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4">
          {/* Auto Pot */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2 text-sm">
                <Zap className="h-4 w-4 text-yellow-500" />
                Auto Pot
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-3">
              <div className="flex items-center justify-between">
                <Label>Enable Auto Pot</Label>
                <Switch
                  checked={autoPotEnabled}
                  onCheckedChange={onAutoPotChange}
                />
              </div>
              
              <div className="flex items-center justify-between">
                <div>
                  <Label>AutoPot Predictivo</Label>
                  <p className="text-xs text-muted-foreground">Predice daño y spamea por 5 segundos</p>
                </div>
                <Switch
                  checked={predictionMode}
                  onCheckedChange={onPredictionModeChange}
                  disabled={!autoPotEnabled}
                />
              </div>
              
              <div className="space-y-2">
                <div>
                  <Label className="text-xs">HP Threshold (%)</Label>
                  <input
                    type="number"
                    value={hpThreshold}
                    onChange={(e) => onHpThresholdChange(Number(e.target.value))}
                    className="w-full mt-1 px-2 py-1 bg-background border border-border rounded text-xs"
                    min="1"
                    max="100"
                  />
                </div>
                
                <div>
                  <Label className="text-xs">SP Threshold (%)</Label>
                  <input
                    type="number"
                    value={spThreshold}
                    onChange={(e) => onSpThresholdChange(Number(e.target.value))}
                    className="w-full mt-1 px-2 py-1 bg-background border border-border rounded text-xs"
                    min="1"
                    max="100"
                  />
                </div>
                
                <div>
                  <Label className="text-xs">HP Pot Key</Label>
                  <input
                    type="text"
                    value={hpPotKey}
                    onChange={(e) => onHpPotKeyChange(e.target.value)}
                    className="w-full mt-1 px-2 py-1 bg-background border border-border rounded text-xs"
                    placeholder="F1"
                  />
                </div>
                
                <div>
                  <Label className="text-xs">SP Pot Key</Label>
                  <input
                    type="text"
                    value={spPotKey}
                    onChange={(e) => onSpPotKeyChange(e.target.value)}
                    className="w-full mt-1 px-2 py-1 bg-background border border-border rounded text-xs"
                    placeholder="F2"
                  />
                </div>
                
                <div>
                  <Label className="text-xs">Check Delay (ms)</Label>
                  <input
                    type="number"
                    value={autopotDelay}
                    onChange={(e) => onAutopotDelayChange(Number(e.target.value))}
                    className="w-full mt-1 px-2 py-1 bg-background border border-border rounded text-xs"
                    min="1"
                    max="1000"
                  />
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Auto Buff */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2 text-sm">
                <Shield className="h-4 w-4 text-blue-500" />
                Auto Buff
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-3">
              <div className="flex items-center justify-between">
                <Label>Enable Auto Buff</Label>
                <Switch
                  checked={autoBuffEnabled}
                  onCheckedChange={onAutoBuffChange}
                />
              </div>
              
              <div className="space-y-2">
                <div>
                  <Label className="text-xs">Gloom Key</Label>
                  <input
                    type="text"
                    value={gloomKey}
                    onChange={(e) => onGloomKeyChange(e.target.value)}
                    className="w-full mt-1 px-2 py-1 bg-background border border-border rounded text-xs"
                    placeholder="F3"
                  />
                </div>
                
                <div>
                  <Label className="text-xs">ASPD Key</Label>
                  <input
                    type="text"
                    value={aspdKey}
                    onChange={(e) => onAspdKeyChange(e.target.value)}
                    className="w-full mt-1 px-2 py-1 bg-background border border-border rounded text-xs"
                    placeholder="F4"
                  />
                </div>
                
                <div>
                  <Label className="text-xs">Buff Delay (ms)</Label>
                  <input
                    type="number"
                    value={autobuffDelay}
                    onChange={(e) => onAutobuffDelayChange(Number(e.target.value))}
                    className="w-full mt-1 px-2 py-1 bg-background border border-border rounded text-xs"
                    min="1000"
                    max="30000"
                  />
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Spammer */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2 text-sm">
                <Target className="h-4 w-4 text-red-500" />
                Spammer
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-3">
              <div className="flex items-center justify-between">
                <Label>Enable Spammer</Label>
                <Switch
                  checked={spammerEnabled}
                  onCheckedChange={onSpammerChange}
                />
              </div>
              
              <div>
                <Label className="text-sm">Delay (ms)</Label>
                <input
                  type="number"
                  step="0.1"
                  value={spammerDelay}
                  onChange={(e) => onSpammerDelayChange(Number(e.target.value))}
                  className="w-full mt-1 px-3 py-2 bg-background border border-border rounded-md text-sm"
                  min="0.1"
                  max="1000"
                />
              </div>
              
              <div>
                <Label className="text-sm">Spam Keys (comma separated)</Label>
                <input
                  type="text"
                  value={spammerKeys.join(', ')}
                  onChange={(e) => onSpammerKeysChange(e.target.value.split(',').map(k => k.trim()).filter(k => k))}
                  className="w-full mt-1 px-3 py-2 bg-background border border-border rounded-md text-sm"
                  placeholder="F5, F6, F7"
                />
              </div>
              
              <div className="space-y-3 pt-2 border-t border-border">
                <div className="text-sm font-medium text-foreground">Click Control</div>
                
                <div>
                  <Label className="text-xs">Key/Click Ratio</Label>
                  <input
                    type="number"
                    value={keyClickRatio}
                    onChange={(e) => onKeyClickRatioChange(Number(e.target.value))}
                    className="w-full mt-1 px-2 py-1 bg-background border border-border rounded text-xs"
                    min="1"
                    max="10"
                    placeholder="3"
                  />
                  <div className="text-xs text-muted-foreground mt-1">
                    {keyClickRatio} tecla(s) por 1 click
                  </div>
                </div>
                
                <div>
                  <Label className="text-xs">Click Delay (ms)</Label>
                  <input
                    type="number"
                    value={clickDelay}
                    onChange={(e) => onClickDelayChange(Number(e.target.value))}
                    className="w-full mt-1 px-2 py-1 bg-background border border-border rounded text-xs"
                    min="0"
                    max="1000"
                    placeholder="50"
                  />
                  <div className="text-xs text-muted-foreground mt-1">
                    Delay adicional entre clicks
                  </div>
                </div>
                
                <div className="flex items-center justify-between">
                  <div>
                    <Label className="text-xs">Mouse Click</Label>
                    <div className="text-xs text-muted-foreground">Habilitar clicks del mouse</div>
                  </div>
                  <Switch
                    checked={mouseClick}
                    onCheckedChange={onMouseClickChange}
                  />
                </div>
                
                <div>
                  <Label className="text-xs">Spammer Mode</Label>
                  <select
                    value={spammerMode}
                    onChange={(e) => onSpammerModeChange(e.target.value)}
                    className="w-full mt-1 px-2 py-1 bg-background border border-border rounded text-xs"
                    disabled
                  >
                    <option value="ultra_fast">Ultra Fast</option>
                  </select>
                  <div className="text-xs text-muted-foreground mt-1">
                    Único modo disponible
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>

      </div>
      
      {/* Fixed Save Button */}
      <div className="border-t bg-background p-4">
        <div className="flex justify-end">
          <Button onClick={onSaveConfig} className="px-6">
            <Save className="h-4 w-4 mr-2" />
            Save Configuration
          </Button>
        </div>
      </div>
    </motion.div>
  )
}
