import React from 'react';
import { EnhancedDecorations } from './EnhancedDecorations';

interface ThemeDecorationsProps {
  isActive?: boolean;
  position?: 'top-left' | 'top-right' | 'bottom-left' | 'bottom-right';
  className?: string;
  intensity?: 'subtle' | 'medium' | 'bold';
}

export const ThemeDecorations: React.FC<ThemeDecorationsProps> = ({ 
  isActive = false,
  position = 'top-right',
  className = '',
  intensity = 'medium'
}) => {
  return (
    <EnhancedDecorations 
      isActive={isActive}
      position={position}
      className={className}
      intensity={intensity}
    />
  );
};
