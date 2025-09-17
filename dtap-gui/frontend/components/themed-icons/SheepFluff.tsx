import React from 'react';
import { IconCloud, IconCloudRain } from '@tabler/icons-react';

interface SheepFluffProps {
  count?: number;
  size?: 'sm' | 'md' | 'lg';
  className?: string;
  animated?: boolean;
}

export const SheepFluff: React.FC<SheepFluffProps> = ({ 
  count = 2, 
  size = 'sm', 
  className = '',
  animated = false
}) => {
  const sizeClasses = {
    sm: 'w-4 h-4',
    md: 'w-5 h-5',
    lg: 'w-6 h-6'
  };

  const animationClass = animated ? 'animate-bounce' : '';

  return (
    <div className={`flex items-center space-x-1 ${className}`}>
      {Array.from({ length: count }).map((_, index) => (
        <div key={index} className="relative">
          <IconCloud
            className={`${sizeClasses[size]} text-white ${animationClass}`}
            style={{ 
              animationDelay: `${index * 0.3}s`,
              filter: 'drop-shadow(0 0 3px rgba(255,255,255,0.4))'
            }}
            stroke={2}
            fill="currentColor"
          />
          <IconCloudRain
            className={`${sizeClasses[size]} text-blue-200 absolute inset-0 ${animationClass}`}
            style={{ 
              animationDelay: `${index * 0.3}s`,
              opacity: 0.6
            }}
            stroke={1.5}
            fill="currentColor"
          />
        </div>
      ))}
    </div>
  );
};
