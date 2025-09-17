import React, { createContext, useContext, useEffect, useState } from 'react'

export type TransparencyMode = 'glass' | 'hyprland'

interface TransparencyContextType {
  transparencyMode: TransparencyMode
  setTransparencyMode: (mode: TransparencyMode) => void
  toggleTransparencyMode: () => void
}

const TransparencyContext = createContext<TransparencyContextType | undefined>(undefined)

interface TransparencyProviderProps {
  children: React.ReactNode
}

export function TransparencyProvider({ children }: TransparencyProviderProps) {
  const [transparencyMode, setTransparencyModeState] = useState<TransparencyMode>('glass')

  // Cargar modo de transparencia desde localStorage al inicializar
  useEffect(() => {
    const savedMode = localStorage.getItem('dtap-transparency-mode') as TransparencyMode
    if (savedMode && (savedMode === 'glass' || savedMode === 'hyprland')) {
      setTransparencyModeState(savedMode)
    }
  }, [])

  // Aplicar modo de transparencia al documento
  useEffect(() => {
    const root = document.documentElement
    root.classList.remove('glass-mode', 'hyprland-mode')
    root.classList.add(`${transparencyMode}-mode`)
    
    // Guardar en localStorage
    localStorage.setItem('dtap-transparency-mode', transparencyMode)
  }, [transparencyMode])

  const setTransparencyMode = (newMode: TransparencyMode) => {
    setTransparencyModeState(newMode)
  }

  const toggleTransparencyMode = () => {
    setTransparencyModeState(prev => prev === 'glass' ? 'hyprland' : 'glass')
  }

  return (
    <TransparencyContext.Provider value={{ transparencyMode, setTransparencyMode, toggleTransparencyMode }}>
      {children}
    </TransparencyContext.Provider>
  )
}

export function useTransparency() {
  const context = useContext(TransparencyContext)
  if (context === undefined) {
    throw new Error('useTransparency must be used within a TransparencyProvider')
  }
  return context
}
