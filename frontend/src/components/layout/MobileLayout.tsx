'use client';

import { ReactNode } from 'react';
import { usePathname } from 'next/navigation';
import Link from 'next/link';
import {
  Home,
  FileText,
  Package,
  Users,
  Menu,
  Plus,
  Bell,
  User,
} from 'lucide-react';
import { cn } from '@/lib/utils';
import { motion } from 'framer-motion';

interface MobileLayoutProps {
  children: ReactNode;
}

const navigationItems = [
  { href: '/dashboard', icon: Home, label: '首頁' },
  { href: '/inquiries', icon: FileText, label: '詢價' },
  { href: '/quotes', icon: Package, label: '報價' },
  { href: '/customers', icon: Users, label: '客戶' },
  { href: '/profile', icon: User, label: '我的' },
];

export default function MobileLayout({ children }: MobileLayoutProps) {
  const pathname = usePathname();

  return (
    <div className="flex flex-col h-screen bg-gray-50">
      {/* Mobile Header */}
      <header className="bg-white border-b border-gray-200 px-4 py-3">
        <div className="flex items-center justify-between">
          <button className="p-2 -ml-2 rounded-lg hover:bg-gray-100">
            <Menu className="h-6 w-6 text-gray-600" />
          </button>
          
          <h1 className="text-lg font-semibold text-gray-900">
            FastenMind
          </h1>
          
          <button className="p-2 -mr-2 rounded-lg hover:bg-gray-100 relative">
            <Bell className="h-6 w-6 text-gray-600" />
            <span className="absolute top-1 right-1 h-2 w-2 bg-red-500 rounded-full" />
          </button>
        </div>
      </header>

      {/* Main Content */}
      <main className="flex-1 overflow-y-auto pb-20">
        <div className="px-4 py-4">
          {children}
        </div>
      </main>

      {/* Floating Action Button */}
      <motion.div
        className="fixed right-4 bottom-24 z-40"
        initial={{ scale: 0 }}
        animate={{ scale: 1 }}
        transition={{ type: 'spring', stiffness: 260, damping: 20 }}
      >
        <button className="bg-blue-600 text-white rounded-full p-4 shadow-lg hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500">
          <Plus className="h-6 w-6" />
        </button>
      </motion.div>

      {/* Bottom Navigation */}
      <nav className="fixed bottom-0 left-0 right-0 bg-white border-t border-gray-200">
        <div className="flex items-center justify-around py-2">
          {navigationItems.map((item) => {
            const isActive = pathname === item.href;
            return (
              <Link
                key={item.href}
                href={item.href}
                className="flex flex-col items-center justify-center py-2 px-3 rounded-lg transition-colors"
              >
                <item.icon
                  className={cn(
                    'h-6 w-6 mb-1',
                    isActive ? 'text-blue-600' : 'text-gray-400'
                  )}
                />
                <span
                  className={cn(
                    'text-xs',
                    isActive ? 'text-blue-600 font-medium' : 'text-gray-400'
                  )}
                >
                  {item.label}
                </span>
                {isActive && (
                  <motion.div
                    className="absolute -bottom-0.5 w-12 h-1 bg-blue-600 rounded-full"
                    layoutId="activeTab"
                  />
                )}
              </Link>
            );
          })}
        </div>
      </nav>
    </div>
  );
}
