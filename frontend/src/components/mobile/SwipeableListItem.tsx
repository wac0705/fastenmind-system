'use client';

import { ReactNode, useRef } from 'react';
import { motion, useAnimation, PanInfo } from 'framer-motion';
import { Trash2, Edit, Archive } from 'lucide-react';

interface SwipeAction {
  icon: typeof Trash2;
  label: string;
  color: string;
  bgColor: string;
  action: () => void;
}

interface SwipeableListItemProps {
  children: ReactNode;
  onDelete?: () => void;
  onEdit?: () => void;
  onArchive?: () => void;
  threshold?: number;
}

export default function SwipeableListItem({
  children,
  onDelete,
  onEdit,
  onArchive,
  threshold = 100,
}: SwipeableListItemProps) {
  const controls = useAnimation();
  const containerRef = useRef<HTMLDivElement>(null);

  const leftActions: SwipeAction[] = [];
  const rightActions: SwipeAction[] = [];

  if (onEdit) {
    leftActions.push({
      icon: Edit,
      label: '編輯',
      color: 'text-white',
      bgColor: 'bg-blue-500',
      action: onEdit,
    });
  }

  if (onArchive) {
    rightActions.push({
      icon: Archive,
      label: '封存',
      color: 'text-white',
      bgColor: 'bg-gray-500',
      action: onArchive,
    });
  }

  if (onDelete) {
    rightActions.push({
      icon: Trash2,
      label: '刪除',
      color: 'text-white',
      bgColor: 'bg-red-500',
      action: onDelete,
    });
  }

  const handleDragEnd = async (_: any, info: PanInfo) => {
    const offset = info.offset.x;
    const velocity = info.velocity.x;
    const width = containerRef.current?.offsetWidth || 0;

    // Determine if we should open the actions
    if (Math.abs(offset) > threshold || Math.abs(velocity) > 500) {
      // Swiped far enough or fast enough
      if (offset > 0 && leftActions.length > 0) {
        // Swiped right - show left actions
        await controls.start({ x: width * 0.3 });
      } else if (offset < 0 && rightActions.length > 0) {
        // Swiped left - show right actions
        await controls.start({ x: -width * 0.3 });
      } else {
        // Snap back
        await controls.start({ x: 0 });
      }
    } else {
      // Snap back
      await controls.start({ x: 0 });
    }
  };

  const handleActionClick = async (action: SwipeAction) => {
    // Animate closing
    await controls.start({ x: 0 });
    // Execute action
    action.action();
  };

  return (
    <div ref={containerRef} className="relative overflow-hidden">
      {/* Left Actions */}
      {leftActions.length > 0 && (
        <div className="absolute inset-y-0 left-0 flex">
          {leftActions.map((action, index) => (
            <button
              key={index}
              onClick={() => handleActionClick(action)}
              className={`flex items-center justify-center px-6 ${action.bgColor} ${action.color}`}
            >
              <div className="flex flex-col items-center">
                <action.icon className="h-5 w-5 mb-1" />
                <span className="text-xs">{action.label}</span>
              </div>
            </button>
          ))}
        </div>
      )}

      {/* Right Actions */}
      {rightActions.length > 0 && (
        <div className="absolute inset-y-0 right-0 flex">
          {rightActions.map((action, index) => (
            <button
              key={index}
              onClick={() => handleActionClick(action)}
              className={`flex items-center justify-center px-6 ${action.bgColor} ${action.color}`}
            >
              <div className="flex flex-col items-center">
                <action.icon className="h-5 w-5 mb-1" />
                <span className="text-xs">{action.label}</span>
              </div>
            </button>
          ))}
        </div>
      )}

      {/* Swipeable Content */}
      <motion.div
        drag="x"
        dragElastic={0.2}
        dragConstraints={{ left: -200, right: 200 }}
        onDragEnd={handleDragEnd}
        animate={controls}
        className="relative bg-white z-10"
      >
        {children}
      </motion.div>
    </div>
  );
}
