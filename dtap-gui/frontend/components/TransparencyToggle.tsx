import { motion } from 'framer-motion'
import { useTransparency } from '../contexts/TransparencyContext'
import { Sparkles, Zap } from 'lucide-react'

export default function TransparencyToggle() {
  const { transparencyMode, toggleTransparencyMode } = useTransparency()

  return (
    <motion.button
      onClick={toggleTransparencyMode}
      className="relative flex items-center justify-center w-10 h-10 rounded-lg bg-secondary/50 hover:bg-secondary transition-colors"
      whileHover={{ scale: 1.05 }}
      whileTap={{ scale: 0.95 }}
      title={`Switch to ${transparencyMode === 'glass' ? 'Hyprland' : 'Glass'} mode`}
    >
      <motion.div
        className="relative"
        animate={{ rotate: transparencyMode === 'hyprland' ? 180 : 0 }}
        transition={{ duration: 0.3 }}
      >
        {transparencyMode === 'glass' ? (
          <Sparkles className="h-4 w-4 text-primary" />
        ) : (
          <Zap className="h-4 w-4 text-warning" />
        )}
      </motion.div>
      
      {/* Indicator dot */}
      <motion.div
        className="absolute -top-1 -right-1 w-3 h-3 rounded-full"
        animate={{
          backgroundColor: transparencyMode === 'glass' ? 'hsl(var(--primary))' : 'hsl(var(--warning))'
        }}
        transition={{ duration: 0.3 }}
      />
    </motion.button>
  )
}
