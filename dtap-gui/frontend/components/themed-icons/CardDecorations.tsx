import React from 'react';
import { useTheme } from '../../contexts/ThemeContext';
import { IconCat, IconHeart } from '@tabler/icons-react';

interface CardDecorationsProps {
  isActive?: boolean;
  className?: string;
}

export const CardDecorations: React.FC<CardDecorationsProps> = ({ 
  isActive = false,
  className = ''
}) => {
  const { theme } = useTheme();

  if (theme === 'dark') {
    return (
      <div className={`absolute top-3 right-3 ${className}`}>
        <IconCat 
          className={`w-5 h-5 text-gray-400 transition-all duration-300 ${
            isActive ? 'animate-pulse text-gray-300' : ''
          }`}
          stroke={2}
          fill="currentColor"
        />
      </div>
    );
  }

  // Light mode - Corazón
  return (
    <div className={`absolute top-3 right-3 ${className}`}>
      <IconHeart 
        className={`w-5 h-5 text-pink-300 transition-all duration-300 ${
          isActive ? 'animate-pulse text-pink-400' : ''
        }`}
        stroke={2}
        fill="currentColor"
      />
    </div>
  );
};
