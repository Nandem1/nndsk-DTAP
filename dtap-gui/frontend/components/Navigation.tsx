import { motion } from 'framer-motion'
import { Button } from '@/components/ui/button'
import { Home, Zap, Settings } from 'lucide-react'
import ThemeToggle from './ThemeToggle'
import TransparencyToggle from './TransparencyToggle'
import { useTransparency } from '../contexts/TransparencyContext'
import { NavigationDecorations } from './themed-icons'

type Route = 'main' | 'cheats'

interface NavigationProps {
  currentRoute: Route
  onRouteChange: (route: Route) => void
}

export default function Navigation({ currentRoute, onRouteChange }: NavigationProps) {
  const { transparencyMode } = useTransparency()
  const routes = [
    { id: 'main' as Route, label: 'Main', icon: Home },
    { id: 'cheats' as Route, label: 'Cheats', icon: Zap }
  ]

  const handleDoubleClick = () => {
    if (window.runtime) {
      window.runtime.WindowMaximise()
    }
  }

  return (
    <nav 
      onDoubleClick={handleDoubleClick}
      className={`${transparencyMode === 'hyprland' ? 'hyprland-nav' : 'glass-nav'} titlebar`}
    >
      <div className="flex h-16 items-center px-6">
        <div className="flex items-center space-x-8 flex-1">
          <div className="flex items-center space-x-2">
            <div className="h-8 w-8 rounded bg-primary"></div>
            <span className="text-xl font-bold">DTAP</span>
          </div>
          
          <div className="flex items-center space-x-1 no-drag">
            {routes.map((route) => {
              const Icon = route.icon
              return (
                <Button
                  key={route.id}
                  variant="ghost"
                  onClick={() => onRouteChange(route.id)}
                  className={`nav-item ${currentRoute === route.id ? 'active' : ''}`}
                >
                  <Icon className="h-4 w-4 mr-2" />
                  {route.label}
                </Button>
              )
            })}
          </div>
        </div>
        
        {/* Navigation Decorations between Cheats and Toggles */}
        <NavigationDecorations className="mr-4" />
        
        {/* Theme and Transparency Toggles */}
        <div className="flex items-center space-x-2 no-drag">
          <TransparencyToggle />
          <ThemeToggle />
        </div>
        
        {/* Window Controls */}
        <div className="flex items-center space-x-2 no-drag">
          <Button
            variant="ghost"
            size="sm"
            className="h-8 w-8 p-0"
            onClick={() => window.runtime?.WindowMinimise()}
          >
            <span className="text-xs">−</span>
          </Button>
          <Button
            variant="ghost"
            size="sm"
            className="h-8 w-8 p-0"
            onClick={() => window.runtime?.WindowMaximise()}
          >
            <span className="text-xs">□</span>
          </Button>
          <Button
            variant="ghost"
            size="sm"
            className="h-8 w-8 p-0 hover:bg-destructive hover:text-destructive-foreground"
            onClick={() => window.runtime?.Quit()}
          >
            <span className="text-xs">×</span>
          </Button>
        </div>
      </div>
    </nav>
  )
}
