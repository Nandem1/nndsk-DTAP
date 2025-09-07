import { motion } from 'framer-motion'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { 
  Settings, 
  Monitor,
  Keyboard,
  Mouse,
  Save
} from 'lucide-react'

interface ConfigProps {
  onSaveConfig: () => void
}

export default function Config({ onSaveConfig }: ConfigProps) {
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
          <h1 className="text-3xl font-bold">Configuración</h1>
          <p className="text-muted-foreground">Advanced settings and preferences</p>
        </div>

        {/* Config Grid */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {/* Display Settings */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Monitor className="h-5 w-5" />
                Display Settings
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div>
                <Label className="text-sm">Window Width</Label>
                <Input
                  type="number"
                  defaultValue="1200"
                  className="mt-1"
                />
              </div>
              
              <div>
                <Label className="text-sm">Window Height</Label>
                <Input
                  type="number"
                  defaultValue="800"
                  className="mt-1"
                />
              </div>
              
              <div>
                <Label className="text-sm">Theme</Label>
                <select className="w-full mt-1 px-3 py-2 bg-background border border-border rounded-md text-sm">
                  <option>Dark</option>
                  <option>Light</option>
                </select>
              </div>
            </CardContent>
          </Card>

          {/* Input Settings */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Keyboard className="h-5 w-5" />
                Input Settings
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div>
                <Label className="text-sm">HP Pot Key</Label>
                <Input
                  defaultValue="F1"
                  className="mt-1"
                />
              </div>
              
              <div>
                <Label className="text-sm">SP Pot Key</Label>
                <Input
                  defaultValue="F2"
                  className="mt-1"
                />
              </div>
              
              <div>
                <Label className="text-sm">Buff Key</Label>
                <Input
                  defaultValue="F3"
                  className="mt-1"
                />
              </div>
            </CardContent>
          </Card>

          {/* Mouse Settings */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Mouse className="h-5 w-5" />
                Mouse Settings
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div>
                <Label className="text-sm">Click Delay (ms)</Label>
                <Input
                  type="number"
                  defaultValue="100"
                  className="mt-1"
                />
              </div>
              
              <div>
                <Label className="text-sm">Double Click</Label>
                <select className="w-full mt-1 px-3 py-2 bg-background border border-border rounded-md text-sm">
                  <option>Enabled</option>
                  <option>Disabled</option>
                </select>
              </div>
            </CardContent>
          </Card>

          {/* System Settings */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Settings className="h-5 w-5" />
                System Settings
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div>
                <Label className="text-sm">Check Delay (ms)</Label>
                <Input
                  type="number"
                  defaultValue="100"
                  className="mt-1"
                />
              </div>
              
              <div>
                <Label className="text-sm">Log Level</Label>
                <select className="w-full mt-1 px-3 py-2 bg-background border border-border rounded-md text-sm">
                  <option>Info</option>
                  <option>Debug</option>
                  <option>Error</option>
                </select>
              </div>
              
              <div>
                <Label className="text-sm">Show Debug</Label>
                <select className="w-full mt-1 px-3 py-2 bg-background border border-border rounded-md text-sm">
                  <option>No</option>
                  <option>Yes</option>
                </select>
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Save Button */}
        <div className="flex justify-end">
          <Button onClick={onSaveConfig} className="px-6">
            <Save className="h-4 w-4 mr-2" />
            Save All Settings
          </Button>
        </div>
      </div>
    </motion.div>
  )
}
