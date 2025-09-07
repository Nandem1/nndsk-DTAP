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
      <div className="flex-1 p-6 space-y-6">
        {/* Header */}
        <div>
          <h1 className="text-3xl font-bold">Cheats</h1>
          <p className="text-muted-foreground">Configure your automation settings</p>
        </div>

        {/* Cheats Grid */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Auto Pot */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Zap className="h-5 w-5 text-yellow-500" />
                Auto Pot
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="flex items-center justify-between">
                <Label>Enable Auto Pot</Label>
                <Switch
                  checked={autoPotEnabled}
                  onCheckedChange={onAutoPotChange}
                />
              </div>
              
              <div className="space-y-3">
                <div>
                  <Label className="text-sm">HP Threshold (%)</Label>
                  <input
                    type="number"
                    value={hpThreshold}
                    onChange={(e) => onHpThresholdChange(Number(e.target.value))}
                    className="w-full mt-1 px-3 py-2 bg-background border border-border rounded-md text-sm"
                    min="1"
                    max="100"
                  />
                </div>
                
                <div>
                  <Label className="text-sm">SP Threshold (%)</Label>
                  <input
                    type="number"
                    value={spThreshold}
                    onChange={(e) => onSpThresholdChange(Number(e.target.value))}
                    className="w-full mt-1 px-3 py-2 bg-background border border-border rounded-md text-sm"
                    min="1"
                    max="100"
                  />
                </div>
                
                <div>
                  <Label className="text-sm">HP Pot Key</Label>
                  <input
                    type="text"
                    value={hpPotKey}
                    onChange={(e) => onHpPotKeyChange(e.target.value)}
                    className="w-full mt-1 px-3 py-2 bg-background border border-border rounded-md text-sm"
                    placeholder="F1"
                  />
                </div>
                
                <div>
                  <Label className="text-sm">SP Pot Key</Label>
                  <input
                    type="text"
                    value={spPotKey}
                    onChange={(e) => onSpPotKeyChange(e.target.value)}
                    className="w-full mt-1 px-3 py-2 bg-background border border-border rounded-md text-sm"
                    placeholder="F2"
                  />
                </div>
                
                <div className="text-sm text-muted-foreground">
                  Check Delay: 100ms (hardcoded)
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Auto Buff */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Shield className="h-5 w-5 text-blue-500" />
                Auto Buff
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="flex items-center justify-between">
                <Label>Enable Auto Buff</Label>
                <Switch
                  checked={autoBuffEnabled}
                  onCheckedChange={onAutoBuffChange}
                />
              </div>
              
              <div className="space-y-3">
                <div>
                  <Label className="text-sm">Gloom Key</Label>
                  <input
                    type="text"
                    value={gloomKey}
                    onChange={(e) => onGloomKeyChange(e.target.value)}
                    className="w-full mt-1 px-3 py-2 bg-background border border-border rounded-md text-sm"
                    placeholder="F3"
                  />
                </div>
                
                <div>
                  <Label className="text-sm">ASPD Key</Label>
                  <input
                    type="text"
                    value={aspdKey}
                    onChange={(e) => onAspdKeyChange(e.target.value)}
                    className="w-full mt-1 px-3 py-2 bg-background border border-border rounded-md text-sm"
                    placeholder="F4"
                  />
                </div>
                
                <div className="text-sm text-muted-foreground">
                  Delay: 5000ms (hardcoded)
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Spammer */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Target className="h-5 w-5 text-red-500" />
                Spammer
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
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
            </CardContent>
          </Card>
        </div>

        {/* Save Button */}
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
