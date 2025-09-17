import React from 'react';
import { IconCat, IconEye } from '@tabler/icons-react';

interface CatEyesProps {
  isActive: boolean;
  size?: 'sm' | 'md' | 'lg';
  className?: string;
}

export const CatEyes: React.FC<CatEyesProps> = ({ 
  isActive, 
  size = 'sm', 
  className = '' 
}) => {
  const sizeClasses = {
    sm: 'w-4 h-4',
    md: 'w-5 h-5',
    lg: 'w-6 h-6'
  };

  const eyeColor = isActive 
    ? 'text-emerald-400' 
    : 'text-yellow-400';

  return (
    <div className={`flex items-center space-x-1 ${className}`}>
      <div className="relative">
        <IconCat 
          className={`${sizeClasses[size]} ${eyeColor} transition-all duration-300 ${
            isActive ? 'animate-pulse' : ''
          }`}
          stroke={2}
        />
        <div className="absolute inset-0 flex items-center justify-center">
          <IconEye 
            className={`${sizeClasses[size === 'sm' ? 'sm' : 'md']} ${eyeColor} transition-colors duration-300`}
            stroke={1.5}
            fill="currentColor"
          />
        </div>
      </div>
    </div>
  );
};
