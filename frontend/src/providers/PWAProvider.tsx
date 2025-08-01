'use client';

import { createContext, useContext, useEffect, ReactNode } from 'react';
import { usePWA } from '@/hooks/usePWA';
import PWAInstallBanner from '@/components/PWAInstallBanner';
import { Toaster } from 'react-hot-toast';

interface PWAContextType {
  isInstallable: boolean;
  isInstalled: boolean;
  isOffline: boolean;
  installPWA: () => Promise<boolean>;
  subscribeToPush: () => Promise<boolean>;
  unsubscribeFromPush: () => Promise<boolean>;
  checkPushStatus: () => Promise<{ supported: boolean; enabled: boolean }>;
}

const PWAContext = createContext<PWAContextType | null>(null);

export function usePWAContext() {
  const context = useContext(PWAContext);
  if (!context) {
    throw new Error('usePWAContext must be used within PWAProvider');
  }
  return context;
}

interface PWAProviderProps {
  children: ReactNode;
}

export function PWAProvider({ children }: PWAProviderProps) {
  const pwa = usePWA();

  useEffect(() => {
    // Register service worker on mount
    if ('serviceWorker' in navigator && process.env.NODE_ENV === 'production') {
      window.addEventListener('load', () => {
        navigator.serviceWorker.register('/sw.js').catch((error) => {
          console.error('Service Worker registration failed:', error);
        });
      });
    }

    // Handle app installed event
    window.addEventListener('appinstalled', () => {
      console.log('PWA was installed');
    });

    // Handle beforeinstallprompt for browsers that support it
    let deferredPrompt: any;
    window.addEventListener('beforeinstallprompt', (e) => {
      e.preventDefault();
      deferredPrompt = e;
    });
  }, []);

  return (
    <PWAContext.Provider
      value={{
        isInstallable: pwa.isInstallable,
        isInstalled: pwa.isInstalled,
        isOffline: pwa.isOffline,
        installPWA: pwa.installPWA,
        subscribeToPush: pwa.subscribeToPush,
        unsubscribeFromPush: pwa.unsubscribeFromPush,
        checkPushStatus: pwa.checkPushStatus,
      }}
    >
      {children}
      <PWAInstallBanner />
      <Toaster position="top-center" />
    </PWAContext.Provider>
  );
}
