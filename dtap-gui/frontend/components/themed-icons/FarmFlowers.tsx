import React from 'react';
import { IconFlower, IconFlowerOff } from '@tabler/icons-react';

interface FarmFlowersProps {
  isActive: boolean;
  size?: 'sm' | 'md' | 'lg';
  className?: string;
}

export const FarmFlowers: React.FC<FarmFlowersProps> = ({ 
  isActive, 
  size = 'sm', 
  className = '' 
}) => {
  const sizeClasses = {
    sm: 'w-4 h-4',
    md: 'w-5 h-5',
    lg: 'w-6 h-6'
  };

  const flowerColor = isActive 
    ? 'text-green-400' 
    : 'text-yellow-300';

  return (
    <div className={`flex items-center space-x-1 ${className}`}>
      <div className="relative">
        <IconFlower 
          className={`${sizeClasses[size]} ${flowerColor} transition-all duration-300 ${
            isActive ? 'animate-pulse' : ''
          }`}
          stroke={2}
          fill="currentColor"
        />
        <IconFlowerOff 
          className={`${sizeClasses[size]} ${flowerColor} absolute inset-0 transition-all duration-300 ${
            isActive ? 'animate-pulse' : ''
          }`}
          stroke={1.5}
          fill="currentColor"
          style={{ opacity: 0.7 }}
        />
      </div>
      <div className="relative">
        <IconFlower 
          className={`${sizeClasses[size]} ${flowerColor} transition-all duration-300 ${
            isActive ? 'animate-pulse' : ''
          }`}
          stroke={2}
          fill="currentColor"
          style={{ animationDelay: '0.2s' }}
        />
        <IconFlowerOff 
          className={`${sizeClasses[size]} ${flowerColor} absolute inset-0 transition-all duration-300 ${
            isActive ? 'animate-pulse' : ''
          }`}
          stroke={1.5}
          fill="currentColor"
          style={{ opacity: 0.7, animationDelay: '0.2s' }}
        />
      </div>
    </div>
  );
};
