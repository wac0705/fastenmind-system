import { useEffect, useState, useCallback } from 'react';
import { toast } from 'react-hot-toast';

interface BeforeInstallPromptEvent extends Event {
  prompt: () => Promise<void>;
  userChoice: Promise<{ outcome: 'accepted' | 'dismissed' }>;
}

interface PWAState {
  isInstallable: boolean;
  isInstalled: boolean;
  isOffline: boolean;
  deferredPrompt: BeforeInstallPromptEvent | null;
  swRegistration: ServiceWorkerRegistration | null;
}

export function usePWA() {
  const [state, setState] = useState<PWAState>({
    isInstallable: false,
    isInstalled: false,
    isOffline: typeof navigator !== 'undefined' ? !navigator.onLine : false,
    deferredPrompt: null,
    swRegistration: null,
  });

  // Check if app is installed
  useEffect(() => {
    if (typeof navigator !== 'undefined' && 'getInstalledRelatedApps' in navigator) {
      (navigator as any).getInstalledRelatedApps().then((apps: any[]) => {
        setState(prev => ({ ...prev, isInstalled: apps.length > 0 }));
      });
    }

    // Check if running as PWA
    if (typeof window !== 'undefined') {
      const isStandalone = window.matchMedia('(display-mode: standalone)').matches ||
                          (window.navigator as any).standalone ||
                          document.referrer.includes('android-app://');
      
      setState(prev => ({ ...prev, isInstalled: isStandalone }));
    }
  }, []);

  // Handle beforeinstallprompt event
  useEffect(() => {
    const handleBeforeInstallPrompt = (e: Event) => {
      e.preventDefault();
      const promptEvent = e as BeforeInstallPromptEvent;
      setState(prev => ({
        ...prev,
        deferredPrompt: promptEvent,
        isInstallable: true,
      }));
    };

    window.addEventListener('beforeinstallprompt', handleBeforeInstallPrompt);

    return () => {
      window.removeEventListener('beforeinstallprompt', handleBeforeInstallPrompt);
    };
  }, []);

  // Handle online/offline status
  useEffect(() => {
    const handleOnline = () => {
      setState(prev => ({ ...prev, isOffline: false }));
      toast.success('網路已連線');
    };

    const handleOffline = () => {
      setState(prev => ({ ...prev, isOffline: true }));
      toast.error('您目前處於離線模式');
    };

    window.addEventListener('online', handleOnline);
    window.addEventListener('offline', handleOffline);

    return () => {
      window.removeEventListener('online', handleOnline);
      window.removeEventListener('offline', handleOffline);
    };
  }, []);

  // Register service worker
  useEffect(() => {
    if (typeof navigator !== 'undefined' && 'serviceWorker' in navigator && process.env.NODE_ENV === 'production') {
      navigator.serviceWorker
        .register('/sw.js')
        .then((registration) => {
          setState(prev => ({ ...prev, swRegistration: registration }));

          // Check for updates
          registration.addEventListener('updatefound', () => {
            const newWorker = registration.installing;
            if (newWorker) {
              newWorker.addEventListener('statechange', () => {
                if (newWorker.state === 'installed' && typeof navigator !== 'undefined' && navigator.serviceWorker.controller) {
                  toast(
                    (t) => (
                      <div>
                        <p>New version available</p>
                        <button
                          className="ml-2 bg-blue-500 text-white px-3 py-1 rounded"
                          onClick={() => {
                            newWorker.postMessage({ type: 'SKIP_WAITING' });
                            window.location.reload();
                            toast.dismiss(t.id);
                          }}
                        >
                          Update
                        </button>
                      </div>
                    ),
                    { duration: Infinity }
                  );
                }
              });
            }
          });
        })
        .catch((error) => {
          console.error('Service Worker registration failed:', error);
        });
    }
  }, []);

  // Install PWA
  const installPWA = useCallback(async () => {
    if (!state.deferredPrompt) {
      console.log('No deferred prompt available');
      return false;
    }

    try {
      await state.deferredPrompt.prompt();
      const { outcome } = await state.deferredPrompt.userChoice;
      
      if (outcome === 'accepted') {
        setState(prev => ({
          ...prev,
          isInstalled: true,
          isInstallable: false,
          deferredPrompt: null,
        }));
        toast.success('FastenMind 已安裝到您的裝置');
        return true;
      }
    } catch (error) {
      console.error('Error installing PWA:', error);
      toast.error('安裝失敗，請稍後再試');
    }

    setState(prev => ({
      ...prev,
      deferredPrompt: null,
      isInstallable: false,
    }));
    return false;
  }, [state.deferredPrompt]);

  // Subscribe to push notifications
  const subscribeToPush = useCallback(async () => {
    if (!state.swRegistration) {
      console.error('No service worker registration');
      return false;
    }

    try {
      const permission = await Notification.requestPermission();
      
      if (permission !== 'granted') {
        toast.error('請允許通知權限以接收重要訊息');
        return false;
      }

      const subscription = await state.swRegistration.pushManager.subscribe({
        userVisibleOnly: true,
        applicationServerKey: urlBase64ToUint8Array(process.env.NEXT_PUBLIC_VAPID_KEY || '') as BufferSource,
      });

      // Send subscription to server
      await fetch('/api/push/subscribe', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(subscription),
      });

      toast.success('已開啟通知功能');
      return true;
    } catch (error) {
      console.error('Error subscribing to push:', error);
      toast.error('訂閱通知失敗');
      return false;
    }
  }, [state.swRegistration]);

  // Unsubscribe from push notifications
  const unsubscribeFromPush = useCallback(async () => {
    if (!state.swRegistration) {
      return false;
    }

    try {
      const subscription = await state.swRegistration.pushManager.getSubscription();
      if (subscription) {
        await subscription.unsubscribe();
        
        // Notify server
        await fetch('/api/push/unsubscribe', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({ endpoint: subscription.endpoint }),
        });

        toast.success('已關閉通知功能');
        return true;
      }
    } catch (error) {
      console.error('Error unsubscribing from push:', error);
      toast.error('取消訂閱失敗');
    }
    return false;
  }, [state.swRegistration]);

  // Check if push notifications are supported and enabled
  const checkPushStatus = useCallback(async () => {
    if (!state.swRegistration || !('PushManager' in window)) {
      return { supported: false, enabled: false };
    }

    try {
      const subscription = await state.swRegistration.pushManager.getSubscription();
      return {
        supported: true,
        enabled: !!subscription,
      };
    } catch (error) {
      console.error('Error checking push status:', error);
      return { supported: false, enabled: false };
    }
  }, [state.swRegistration]);

  // Cache resources for offline use
  const cacheResources = useCallback(async (urls: string[]) => {
    if (!state.swRegistration) {
      return false;
    }

    try {
      const messageChannel = new MessageChannel();
      
      return new Promise<boolean>((resolve) => {
        messageChannel.port1.onmessage = (event) => {
          resolve(event.data.success);
        };

        if (typeof navigator !== 'undefined') {
          navigator.serviceWorker.controller?.postMessage(
            {
              type: 'CACHE_RESOURCES',
              resources: urls,
            },
            [messageChannel.port2]
          );
        }

        // Timeout after 5 seconds
        setTimeout(() => resolve(false), 5000);
      });
    } catch (error) {
      console.error('Error caching resources:', error);
      return false;
    }
  }, [state.swRegistration]);

  return {
    ...state,
    installPWA,
    subscribeToPush,
    unsubscribeFromPush,
    checkPushStatus,
    cacheResources,
  };
}

// Helper function to convert VAPID key
function urlBase64ToUint8Array(base64String: string): Uint8Array {
  const padding = '='.repeat((4 - (base64String.length % 4)) % 4);
  const base64 = (base64String + padding)
    .replace(/\-/g, '+')
    .replace(/_/g, '/');

  const rawData = window.atob(base64);
  const buffer = new ArrayBuffer(rawData.length);
  const outputArray = new Uint8Array(buffer);

  for (let i = 0; i < rawData.length; ++i) {
    outputArray[i] = rawData.charCodeAt(i);
  }
  return outputArray;
}
