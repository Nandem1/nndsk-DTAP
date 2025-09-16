import { motion } from 'framer-motion'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { useTheme } from '../contexts/ThemeContext'
import { Palette, Moon, Sun } from 'lucide-react'

export default function ThemeDemo() {
  const { theme, toggleTheme } = useTheme()

  const colorSwatches = [
    { name: 'Primary', class: 'bg-primary text-primary-foreground' },
    { name: 'Secondary', class: 'bg-secondary text-secondary-foreground' },
    { name: 'Accent', class: 'bg-accent text-accent-foreground' },
    { name: 'Muted', class: 'bg-muted text-muted-foreground' },
    { name: 'Card', class: 'bg-card text-card-foreground border border-border' },
    { name: 'Background', class: 'bg-background text-foreground border border-border' }
  ]

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.3 }}
      className="space-y-6"
    >
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Palette className="h-5 w-5 text-primary" />
            Theme System Demo
          </CardTitle>
        </CardHeader>
        <CardContent className="space-y-6">
          {/* Current Theme Info */}
          <div className="flex items-center justify-between p-4 bg-secondary rounded-lg">
            <div className="flex items-center gap-3">
              {theme === 'dark' ? (
                <Moon className="h-5 w-5 text-primary" />
              ) : (
                <Sun className="h-5 w-5 text-primary" />
              )}
              <div>
                <p className="font-medium">Current Theme</p>
                <p className="text-sm text-muted-foreground capitalize">{theme} Mode</p>
              </div>
            </div>
            <Button onClick={toggleTheme} variant="outline" size="sm">
              Switch to {theme === 'dark' ? 'Light' : 'Dark'}
            </Button>
          </div>

          {/* Color Palette */}
          <div>
            <h3 className="text-lg font-semibold mb-4">Color Palette</h3>
            <div className="grid grid-cols-2 md:grid-cols-3 gap-4">
              {colorSwatches.map((swatch, index) => (
                <motion.div
                  key={swatch.name}
                  initial={{ opacity: 0, scale: 0.9 }}
                  animate={{ opacity: 1, scale: 1 }}
                  transition={{ delay: index * 0.1 }}
                  className={`p-4 rounded-lg ${swatch.class} transition-all duration-300`}
                >
                  <p className="font-medium">{swatch.name}</p>
                  <p className="text-xs opacity-75 mt-1">
                    {swatch.class.split(' ')[0].replace('bg-', '')}
                  </p>
                </motion.div>
              ))}
            </div>
          </div>

          {/* Theme Features */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <Card>
              <CardContent className="p-4">
                <h4 className="font-semibold mb-2">✨ Features</h4>
                <ul className="text-sm space-y-1 text-muted-foreground">
                  <li>• Smooth transitions (0.3s)</li>
                  <li>• Persistent theme storage</li>
                  <li>• Coquette amethyst palette</li>
                  <li>• Animated toggle button</li>
                </ul>
              </CardContent>
            </Card>
            
            <Card>
              <CardContent className="p-4">
                <h4 className="font-semibold mb-2">🎨 Palette</h4>
                <ul className="text-sm space-y-1 text-muted-foreground">
                  <li>• HSL color system</li>
                  <li>• CSS custom properties</li>
                  <li>• Tailwind integration</li>
                  <li>• Dark/Light variants</li>
                </ul>
              </CardContent>
            </Card>
          </div>
        </CardContent>
      </Card>
    </motion.div>
  )
}
