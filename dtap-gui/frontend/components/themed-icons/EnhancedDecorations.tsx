import React from 'react';
import { useTheme } from '../../contexts/ThemeContext';
import { 
  IconCat, 
  IconPaw, 
  IconCloud, 
  IconFlower, 
  IconSun,
  IconMoon,
  IconStar,
  IconHeart
} from '@tabler/icons-react';

interface EnhancedDecorationsProps {
  isActive?: boolean;
  position?: 'top-left' | 'top-right' | 'bottom-left' | 'bottom-right';
  className?: string;
  intensity?: 'subtle' | 'medium' | 'bold';
}

export const EnhancedDecorations: React.FC<EnhancedDecorationsProps> = ({ 
  isActive = false,
  position = 'top-right',
  intensity = 'medium',
  className = ''
}) => {
  const { theme } = useTheme();

  const positionClasses = {
    'top-left': 'top-3 left-3',
    'top-right': 'top-3 right-3',
    'bottom-left': 'bottom-3 left-3',
    'bottom-right': 'bottom-3 right-3'
  };

  const intensityClasses = {
    subtle: 'opacity-40',
    medium: 'opacity-70',
    bold: 'opacity-100'
  };

  if (theme === 'dark') {
    return (
      <div className={`absolute ${positionClasses[position]} ${intensityClasses[intensity]} ${className}`}>
        <div className="flex flex-col items-center space-y-2">
          {/* Gatito principal */}
          <div className="relative">
            <IconCat 
              className={`w-6 h-6 ${isActive ? 'text-emerald-400' : 'text-gray-400'} transition-all duration-500 ${
                isActive ? 'animate-pulse' : ''
              }`}
              stroke={2}
              fill="currentColor"
            />
            {isActive && (
              <IconStar 
                className="w-3 h-3 text-yellow-400 absolute -top-1 -right-1 animate-ping"
                stroke={2}
                fill="currentColor"
              />
            )}
          </div>
          
          {/* Huellitas */}
          <div className="flex space-x-1">
            <IconPaw 
              className={`w-3 h-3 text-gray-300 transition-all duration-300 ${
                isActive ? 'animate-bounce' : ''
              }`}
              stroke={2}
              fill="currentColor"
            />
            <IconPaw 
              className={`w-3 h-3 text-gray-300 transition-all duration-300 ${
                isActive ? 'animate-bounce' : ''
              }`}
              stroke={2}
              fill="currentColor"
              style={{ animationDelay: '0.2s' }}
            />
          </div>
          
          {/* Luna para el tema oscuro */}
          <IconMoon 
            className={`w-4 h-4 text-blue-300 transition-all duration-500 ${
              isActive ? 'animate-pulse' : ''
            }`}
            stroke={2}
            fill="currentColor"
          />
        </div>
      </div>
    );
  }

  // Light theme - Ovejitas y elementos de granja
  return (
    <div className={`absolute ${positionClasses[position]} ${intensityClasses[intensity]} ${className}`}>
      <div className="flex flex-col items-center space-y-2">
        {/* Nube principal (ovejita) */}
        <div className="relative">
          <IconCloud 
            className={`w-6 h-6 ${isActive ? 'text-blue-300' : 'text-gray-300'} transition-all duration-500 ${
              isActive ? 'animate-bounce' : ''
            }`}
            stroke={2}
            fill="currentColor"
          />
          {isActive && (
            <IconSun 
              className="w-3 h-3 text-yellow-400 absolute -top-1 -right-1 animate-spin"
              stroke={2}
              fill="currentColor"
            />
          )}
        </div>
        
        {/* Flores */}
        <div className="flex space-x-1">
          <IconFlower 
            className={`w-3 h-3 ${isActive ? 'text-green-400' : 'text-yellow-300'} transition-all duration-300 ${
              isActive ? 'animate-pulse' : ''
            }`}
            stroke={2}
            fill="currentColor"
          />
          <IconFlower 
            className={`w-3 h-3 ${isActive ? 'text-green-400' : 'text-yellow-300'} transition-all duration-300 ${
              isActive ? 'animate-pulse' : ''
            }`}
            stroke={2}
            fill="currentColor"
            style={{ animationDelay: '0.2s' }}
          />
        </div>
        
        {/* Corazón para el tema claro */}
        <IconHeart 
          className={`w-4 h-4 text-pink-300 transition-all duration-500 ${
            isActive ? 'animate-pulse' : ''
          }`}
          stroke={2}
          fill="currentColor"
        />
      </div>
    </div>
  );
};
