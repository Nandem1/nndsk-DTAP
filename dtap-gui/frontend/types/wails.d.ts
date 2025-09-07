declare global {
  interface Window {
    runtime?: {
      WindowMinimise: () => void
      WindowMaximise: () => void
      WindowUnmaximise: () => void
      WindowClose: () => void
      Quit: () => void
    }
  }
}

export {}
