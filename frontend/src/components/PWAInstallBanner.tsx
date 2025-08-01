'use client';

import { useState, useEffect } from 'react';
import { X, Download, Smartphone, Check } from 'lucide-react';
import { usePWA } from '@/hooks/usePWA';
import { motion, AnimatePresence } from 'framer-motion';

export default function PWAInstallBanner() {
  const { isInstallable, isInstalled, installPWA } = usePWA();
  const [showBanner, setShowBanner] = useState(false);
  const [isInstalling, setIsInstalling] = useState(false);

  useEffect(() => {
    // Check if user has dismissed the banner before
    const dismissed = localStorage.getItem('pwa-install-dismissed');
    const dismissedTime = dismissed ? parseInt(dismissed) : 0;
    const daysSinceDismissed = (Date.now() - dismissedTime) / (1000 * 60 * 60 * 24);

    // Show banner if installable and not recently dismissed (7 days)
    if (isInstallable && !isInstalled && daysSinceDismissed > 7) {
      // Delay showing banner for better UX
      const timer = setTimeout(() => setShowBanner(true), 3000);
      return () => clearTimeout(timer);
    }
  }, [isInstallable, isInstalled]);

  const handleInstall = async () => {
    setIsInstalling(true);
    const installed = await installPWA();
    
    if (installed) {
      setShowBanner(false);
    }
    setIsInstalling(false);
  };

  const handleDismiss = () => {
    setShowBanner(false);
    localStorage.setItem('pwa-install-dismissed', Date.now().toString());
  };

  if (!showBanner) return null;

  return (
    <AnimatePresence>
      <motion.div
        initial={{ y: 100, opacity: 0 }}
        animate={{ y: 0, opacity: 1 }}
        exit={{ y: 100, opacity: 0 }}
        transition={{ type: 'spring', stiffness: 300, damping: 30 }}
        className="fixed bottom-0 left-0 right-0 z-50 p-4 md:p-6"
      >
        <div className="max-w-4xl mx-auto">
          <div className="bg-white rounded-lg shadow-xl border border-gray-200 overflow-hidden">
            <div className="p-4 md:p-6">
              <div className="flex items-start justify-between">
                <div className="flex-1">
                  <div className="flex items-center mb-3">
                    <div className="bg-blue-100 p-2 rounded-lg mr-3">
                      <Smartphone className="h-6 w-6 text-blue-600" />
                    </div>
                    <h3 className="text-lg font-semibold text-gray-900">
                      安裝 FastenMind 應用程式
                    </h3>
                  </div>
                  
                  <p className="text-gray-600 mb-4">
                    將 FastenMind 安裝到您的裝置，享受更快速、更便捷的使用體驗
                  </p>
                  
                  <div className="grid grid-cols-1 md:grid-cols-3 gap-3 mb-5">
                    <div className="flex items-center text-sm text-gray-600">
                      <Check className="h-4 w-4 text-green-500 mr-2 flex-shrink-0" />
                      <span>離線使用</span>
                    </div>
                    <div className="flex items-center text-sm text-gray-600">
                      <Check className="h-4 w-4 text-green-500 mr-2 flex-shrink-0" />
                      <span>快速啟動</span>
                    </div>
                    <div className="flex items-center text-sm text-gray-600">
                      <Check className="h-4 w-4 text-green-500 mr-2 flex-shrink-0" />
                      <span>即時通知</span>
                    </div>
                  </div>
                  
                  <div className="flex flex-wrap gap-3">
                    <button
                      onClick={handleInstall}
                      disabled={isInstalling}
                      className="inline-flex items-center px-4 py-2 bg-blue-600 text-white font-medium rounded-lg hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                    >
                      <Download className="h-4 w-4 mr-2" />
                      {isInstalling ? '安裝中...' : '立即安裝'}
                    </button>
                    
                    <button
                      onClick={handleDismiss}
                      className="inline-flex items-center px-4 py-2 bg-gray-100 text-gray-700 font-medium rounded-lg hover:bg-gray-200 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-500 transition-colors"
                    >
                      稍後再說
                    </button>
                  </div>
                </div>
                
                <button
                  onClick={handleDismiss}
                  className="ml-4 text-gray-400 hover:text-gray-600 focus:outline-none"
                >
                  <X className="h-5 w-5" />
                </button>
              </div>
            </div>
            
            {/* Platform-specific instructions */}
            <div className="bg-gray-50 px-4 py-3 border-t border-gray-200">
              <p className="text-xs text-gray-600">
                <span className="font-medium">提示：</span>
                {navigator.userAgent.includes('iPhone') || navigator.userAgent.includes('iPad')
                  ? ' 在 Safari 中，點擊分享按鈕然後選擇「加入主畫面」'
                  : ' 點擊瀏覽器網址欄右側的安裝按鈕'}
              </p>
            </div>
          </div>
        </div>
      </motion.div>
    </AnimatePresence>
  );
}
