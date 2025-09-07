import { motion } from 'framer-motion'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { ArrowLeft, Settings as SettingsIcon } from 'lucide-react'

interface SettingsProps {
  hpThreshold: number
  spThreshold: number
  spammerDelay: number
  onHpThresholdChange: (value: number) => void
  onSpThresholdChange: (value: number) => void
  onSpammerDelayChange: (value: number) => void
  onConfigUpdate: () => void
  onBack: () => void
}

export default function Settings({
  hpThreshold,
  spThreshold,
  spammerDelay,
  onHpThresholdChange,
  onSpThresholdChange,
  onSpammerDelayChange,
  onConfigUpdate,
  onBack
}: SettingsProps) {
  return (
    <motion.div
      initial={{ opacity: 0, x: 20 }}
      animate={{ opacity: 1, x: 0 }}
      exit={{ opacity: 0, x: -20 }}
      transition={{ duration: 0.3 }}
      className="h-full flex flex-col"
    >
      {/* Header */}
      <motion.div
        initial={{ opacity: 0, y: -10 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.1 }}
        className="mb-8"
      >
        <div className="flex items-center gap-4 mb-4">
          <Button
            onClick={onBack}
            variant="ghost"
            size="sm"
            className="p-2"
          >
            <ArrowLeft className="h-4 w-4" />
          </Button>
          <div>
            <h2 className="text-2xl font-bold text-foreground">Settings</h2>
            <p className="text-muted-foreground text-sm">Configure your DTAP preferences</p>
          </div>
        </div>
      </motion.div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 flex-1">
        {/* Auto Pot Settings */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.2 }}
        >
          <div className="vercel-card p-6 h-full">
            <div className="flex items-center gap-3 mb-6">
              <SettingsIcon className="h-5 w-5 text-primary" />
              <h3 className="text-foreground text-sm font-medium">Auto Pot</h3>
            </div>
            
            <div className="space-y-4">
              <div>
                <Label className="text-foreground text-sm mb-2 block">HP Threshold (%)</Label>
                <Input
                  type="number"
                  value={hpThreshold}
                  onChange={(e) => onHpThresholdChange(Number(e.target.value))}
                  className="vercel-input"
                  min="1"
                  max="100"
                />
              </div>

              <div>
                <Label className="text-foreground text-sm mb-2 block">SP Threshold (%)</Label>
                <Input
                  type="number"
                  value={spThreshold}
                  onChange={(e) => onSpThresholdChange(Number(e.target.value))}
                  className="vercel-input"
                  min="1"
                  max="100"
                />
              </div>
            </div>
          </div>
        </motion.div>

        {/* Spammer Settings */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.3 }}
        >
          <div className="vercel-card p-6 h-full">
            <div className="flex items-center gap-3 mb-6">
              <SettingsIcon className="h-5 w-5 text-primary" />
              <h3 className="text-foreground text-sm font-medium">Spammer</h3>
            </div>
            
            <div className="space-y-4">
              <div>
                <Label className="text-foreground text-sm mb-2 block">Delay (ms)</Label>
                <Input
                  type="number"
                  step="0.1"
                  value={spammerDelay}
                  onChange={(e) => onSpammerDelayChange(Number(e.target.value))}
                  className="vercel-input"
                  min="0.1"
                  max="1000"
                />
              </div>
            </div>
          </div>
        </motion.div>
      </div>

      {/* Save Button */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.4 }}
        className="mt-8"
      >
        <div className="vercel-card p-6">
          <div className="flex items-center justify-between">
            <div>
              <h4 className="text-foreground text-sm font-medium">Save Configuration</h4>
              <p className="text-muted-foreground text-xs">Apply your changes to the system</p>
            </div>
            <Button
              onClick={onConfigUpdate}
              className="vercel-button"
            >
              <SettingsIcon className="h-4 w-4 mr-2" />
              Save Changes
            </Button>
          </div>
        </div>
      </motion.div>
    </motion.div>
  )
}
