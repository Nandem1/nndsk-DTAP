import React from 'react';
import { useTheme } from '../../contexts/ThemeContext';
import { IconPaw, IconCloud } from '@tabler/icons-react';

interface NavigationDecorationsProps {
  className?: string;
}

export const NavigationDecorations: React.FC<NavigationDecorationsProps> = ({ 
  className = ''
}) => {
  const { theme } = useTheme();

  if (theme === 'dark') {
    return (
      <div className={`flex items-center space-x-1 ${className}`}>
        <IconPaw 
          className="w-3 h-3 text-gray-300 animate-bounce"
          stroke={2}
          fill="currentColor"
        />
        <IconPaw 
          className="w-3 h-3 text-gray-300 animate-bounce"
          stroke={2}
          fill="currentColor"
          style={{ animationDelay: '0.2s' }}
        />
        <IconPaw 
          className="w-3 h-3 text-gray-300 animate-bounce"
          stroke={2}
          fill="currentColor"
          style={{ animationDelay: '0.4s' }}
        />
      </div>
    );
  }

  // Light mode - Nubes saltando
  return (
    <div className={`flex items-center space-x-1 ${className}`}>
      <IconCloud 
        className="w-3 h-3 text-blue-300 animate-bounce"
        stroke={2}
        fill="currentColor"
      />
      <IconCloud 
        className="w-3 h-3 text-blue-300 animate-bounce"
        stroke={2}
        fill="currentColor"
        style={{ animationDelay: '0.2s' }}
      />
      <IconCloud 
        className="w-3 h-3 text-blue-300 animate-bounce"
        stroke={2}
        fill="currentColor"
        style={{ animationDelay: '0.4s' }}
      />
    </div>
  );
};
