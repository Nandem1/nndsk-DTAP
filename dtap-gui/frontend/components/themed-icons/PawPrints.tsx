import React from 'react';
import { IconPaw } from '@tabler/icons-react';

interface PawPrintsProps {
  count?: number;
  size?: 'sm' | 'md' | 'lg';
  className?: string;
  animated?: boolean;
}

export const PawPrints: React.FC<PawPrintsProps> = ({ 
  count = 3, 
  size = 'sm', 
  className = '',
  animated = false
}) => {
  const sizeClasses = {
    sm: 'w-3 h-3',
    md: 'w-4 h-4',
    lg: 'w-5 h-5'
  };

  const animationClass = animated ? 'animate-bounce' : '';

  return (
    <div className={`flex items-center space-x-1 ${className}`}>
      {Array.from({ length: count }).map((_, index) => (
        <IconPaw
          key={index}
          className={`${sizeClasses[size]} text-gray-300 ${animationClass}`}
          style={{ 
            animationDelay: `${index * 0.2}s`,
          }}
          stroke={2}
          fill="currentColor"
        />
      ))}
    </div>
  );
};
