'use client';

import { ReactNode, useState, useRef, useCallback } from 'react';
import { motion, useAnimation } from 'framer-motion';
import { RefreshCw } from 'lucide-react';

interface PullToRefreshProps {
  onRefresh: () => Promise<void>;
  children: ReactNode;
  threshold?: number;
}

export default function PullToRefresh({
  onRefresh,
  children,
  threshold = 80,
}: PullToRefreshProps) {
  const [isRefreshing, setIsRefreshing] = useState(false);
  const [pullDistance, setPullDistance] = useState(0);
  const [isDragging, setIsDragging] = useState(false);
  const controls = useAnimation();
  const containerRef = useRef<HTMLDivElement>(null);
  const startY = useRef(0);

  const handleTouchStart = useCallback((e: React.TouchEvent) => {
    if (containerRef.current?.scrollTop === 0) {
      startY.current = e.touches[0].clientY;
      setIsDragging(true);
    }
  }, []);

  const handleTouchMove = useCallback((e: React.TouchEvent) => {
    if (!isDragging || isRefreshing) return;

    const currentY = e.touches[0].clientY;
    const distance = currentY - startY.current;

    if (distance > 0 && containerRef.current?.scrollTop === 0) {
      e.preventDefault();
      // Apply resistance
      const adjustedDistance = Math.min(distance * 0.5, 150);
      setPullDistance(adjustedDistance);
      
      // Rotate the refresh icon based on pull distance
      const rotation = (adjustedDistance / threshold) * 360;
      controls.set({ rotate: rotation });
    }
  }, [isDragging, isRefreshing, threshold, controls]);

  const handleTouchEnd = useCallback(async () => {
    setIsDragging(false);

    if (pullDistance >= threshold && !isRefreshing) {
      setIsRefreshing(true);
      
      // Animate to loading position
      setPullDistance(60);
      controls.start({
        rotate: [0, 360],
        transition: {
          duration: 1,
          repeat: Infinity,
          ease: 'linear',
        },
      });

      try {
        await onRefresh();
      } catch (error) {
        console.error('Refresh failed:', error);
      } finally {
        setIsRefreshing(false);
        setPullDistance(0);
        controls.stop();
        controls.set({ rotate: 0 });
      }
    } else {
      // Snap back
      setPullDistance(0);
      controls.set({ rotate: 0 });
    }
  }, [pullDistance, threshold, isRefreshing, onRefresh, controls]);

  const getRefreshText = () => {
    if (isRefreshing) return '正在刷新...';
    if (pullDistance >= threshold) return '釋放以刷新';
    return '下拉刷新';
  };

  return (
    <div
      ref={containerRef}
      className="relative h-full overflow-y-auto"
      onTouchStart={handleTouchStart}
      onTouchMove={handleTouchMove}
      onTouchEnd={handleTouchEnd}
    >
      {/* Pull to refresh indicator */}
      <div
        className="absolute top-0 left-0 right-0 flex items-center justify-center transition-all duration-200 ease-out"
        style={{
          height: `${pullDistance}px`,
          marginTop: `-${pullDistance}px`,
        }}
      >
        <div className="flex flex-col items-center justify-center">
          <motion.div animate={controls}>
            <RefreshCw
              className={`h-6 w-6 mb-2 ${
                isRefreshing ? 'text-blue-600' : 'text-gray-400'
              }`}
            />
          </motion.div>
          <span className="text-sm text-gray-600">{getRefreshText()}</span>
        </div>
      </div>

      {/* Main content */}
      <div
        style={{
          transform: `translateY(${pullDistance}px)`,
          transition: isDragging ? 'none' : 'transform 0.2s ease-out',
        }}
      >
        {children}
      </div>
    </div>
  );
}
