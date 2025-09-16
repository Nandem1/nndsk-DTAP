# Sistema de Temas DTAP

## Descripción
El sistema de temas de DTAP permite alternar entre modo oscuro y modo claro con una paleta de colores coquette amatista suave.

## Características

### 🎨 Paleta de Colores
- **Modo Oscuro**: Tonos púrpura/amatista suaves con fondo oscuro
- **Modo Claro**: Tonos púrpura/amatista suaves con fondo claro
- **Transiciones suaves**: Animaciones de 0.3s entre cambios de tema

### 🔧 Funcionalidades
- **Persistencia**: El tema seleccionado se guarda en localStorage
- **Toggle animado**: Botón con iconos de sol/luna y animaciones
- **Context API**: Sistema de contexto React para manejo global del tema
- **CSS Variables**: Variables CSS dinámicas para fácil personalización

## Uso

### Hook useTheme
```tsx
import { useTheme } from '../contexts/ThemeContext'

function MyComponent() {
  const { theme, setTheme, toggleTheme } = useTheme()
  
  return (
    <button onClick={toggleTheme}>
      Current theme: {theme}
    </button>
  )
}
```

### Componente ThemeToggle
```tsx
import ThemeToggle from './ThemeToggle'

function Navigation() {
  return (
    <nav>
      {/* ... otros elementos ... */}
      <ThemeToggle />
    </nav>
  )
}
```

## Estructura de Archivos

```
frontend/
├── contexts/
│   └── ThemeContext.tsx      # Contexto y provider del tema
├── components/
│   └── ThemeToggle.tsx       # Componente toggle del tema
├── types/
│   └── theme.d.ts            # Tipos TypeScript
└── styles/
    └── globals.css           # Variables CSS de temas
```

## Variables CSS

### Modo Oscuro (por defecto)
- `--background`: 240 10% 3.9% (fondo muy oscuro con tinte púrpura)
- `--primary`: 270 100% 70% (púrpura vibrante)
- `--secondary`: 240 3.7% 15.9% (gris oscuro con tinte púrpura)

### Modo Claro
- `--background`: 240 10% 98% (fondo claro con tinte púrpura)
- `--primary`: 270 100% 60% (púrpura más suave)
- `--secondary`: 240 4.8% 95.9% (gris claro con tinte púrpura)

## Personalización

Para modificar los colores, edita las variables CSS en `styles/globals.css`:

```css
:root {
  /* Modo oscuro */
  --primary: 270 100% 70%; /* Cambia estos valores */
}

.light {
  /* Modo claro */
  --primary: 270 100% 60%; /* Cambia estos valores */
}
```

## Integración con Tailwind

El sistema usa las clases de Tailwind CSS con variables CSS personalizadas:

```tsx
<div className="bg-background text-foreground">
  <button className="bg-primary text-primary-foreground">
    Botón con tema
  </button>
</div>
```

## Animaciones

- **Transición de tema**: 0.3s ease
- **Toggle button**: Scale animation en hover/tap
- **Iconos**: Rotación suave entre sol/luna
