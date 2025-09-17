# Sistema de Temas DTAP

## Descripción
El sistema de temas de DTAP permite alternar entre modo oscuro (estilo Vercel) y modo claro con una paleta de colores azul elegante sobre crema suave.

## Características

### 🎨 Paleta de Colores
- **Modo Oscuro**: Estilo Vercel con tonos neutros y profesionales
- **Modo Claro**: Crema suave con azul elegante - fondo cálido y acogedor
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

### Modo Oscuro (por defecto) - Estilo Vercel
- `--background`: 0 0% 3.9% (fondo negro puro)
- `--primary`: 0 0% 98% (blanco puro para texto principal)
- `--secondary`: 0 0% 14.9% (gris oscuro neutro)

### Modo Claro - Crema Suave con Azul Elegante
- `--background`: 200 15% 96% (crema suave con tinte azul sutil)
- `--primary`: 200 70% 60% (azul elegante y profesional)
- `--secondary`: 200 15% 90% (gris muy claro con tinte azul mínimo)
- `--foreground`: 200 20% 20% (texto oscuro con mejor contraste)
- `--card`: 200 20% 94% (tarjetas ligeramente más oscuras para contraste)

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
