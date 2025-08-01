const CACHE_NAME = 'fastenmind-v1.0.0';
const STATIC_CACHE_NAME = 'fastenmind-static-v1.0.0';
const DYNAMIC_CACHE_NAME = 'fastenmind-dynamic-v1.0.0';

// 靜態資源快取列表
const STATIC_FILES = [
  '/',
  '/dashboard',
  '/inquiries',
  '/quotes',
  '/orders',
  '/inventory',
  '/offline',
  '/manifest.json',
  // 添加重要的靜態資源
  '/_next/static/css/',
  '/_next/static/js/',
  '/icons/icon-192x192.png',
  '/icons/icon-512x512.png'
];

// API 端點需要網路優先的策略
const API_ENDPOINTS = [
  '/api/inquiries',
  '/api/quotes',
  '/api/orders',
  '/api/inventory',
  '/api/mobile'
];

// 安裝 Service Worker
self.addEventListener('install', (event) => {
  console.log('[SW] Installing Service Worker');
  
  event.waitUntil(
    Promise.all([
      // 快取靜態資源
      caches.open(STATIC_CACHE_NAME).then((cache) => {
        console.log('[SW] Caching static files');
        return cache.addAll(STATIC_FILES.filter(url => !url.includes('_next')));
      }),
      // 立即啟用新的 Service Worker
      self.skipWaiting()
    ])
  );
});

// 啟用 Service Worker
self.addEventListener('activate', (event) => {
  console.log('[SW] Activating Service Worker');
  
  event.waitUntil(
    Promise.all([
      // 清理舊快取
      caches.keys().then((cacheNames) => {
        return Promise.all(
          cacheNames.map((cacheName) => {
            if (cacheName !== STATIC_CACHE_NAME && 
                cacheName !== DYNAMIC_CACHE_NAME &&
                cacheName.startsWith('fastenmind-')) {
              console.log('[SW] Deleting old cache:', cacheName);
              return caches.delete(cacheName);
            }
          })
        );
      }),
      // 立即控制所有客戶端
      self.clients.claim()
    ])
  );
});

// 攔截網路請求
self.addEventListener('fetch', (event) => {
  const requestUrl = new URL(event.request.url);
  
  // 跳過 Chrome 擴展和非 HTTP(S) 請求
  if (!event.request.url.startsWith('http')) {
    return;
  }
  
  // API 請求使用網路優先策略
  if (isApiRequest(event.request.url)) {
    event.respondWith(networkFirstStrategy(event.request));
    return;
  }
  
  // 靜態資源使用快取優先策略
  if (isStaticAsset(event.request.url)) {
    event.respondWith(cacheFirstStrategy(event.request));
    return;
  }
  
  // 頁面請求使用網路優先，快取後備策略
  if (event.request.mode === 'navigate') {
    event.respondWith(navigationStrategy(event.request));
    return;
  }
  
  // 其他請求使用快取優先策略
  event.respondWith(cacheFirstStrategy(event.request));
});

// 處理推播通知
self.addEventListener('push', (event) => {
  console.log('[SW] Push notification received');
  
  let notificationData = {
    title: 'FastenMind',
    body: '您有新的訊息',
    icon: '/icons/icon-192x192.png',
    badge: '/icons/badge-72x72.png',
    tag: 'fastenmind-notification',
    requireInteraction: false,
    actions: []
  };
  
  if (event.data) {
    try {
      const data = event.data.json();
      notificationData = { ...notificationData, ...data };
      
      // 添加動作按鈕
      if (data.type === 'inquiry') {
        notificationData.actions = [
          {
            action: 'view',
            title: '查看詢價',
            icon: '/icons/action-view.png'
          },
          {
            action: 'reply',
            title: '快速回覆',
            icon: '/icons/action-reply.png'
          }
        ];
      } else if (data.type === 'quote') {
        notificationData.actions = [
          {
            action: 'approve',
            title: '核准',
            icon: '/icons/action-approve.png'
          },
          {
            action: 'review',
            title: '審核',
            icon: '/icons/action-review.png'
          }
        ];
      }
      
    } catch (error) {
      console.error('[SW] Error parsing push data:', error);
    }
  }
  
  event.waitUntil(
    self.registration.showNotification(notificationData.title, notificationData)
  );
});

// 處理通知點擊
self.addEventListener('notificationclick', (event) => {
  console.log('[SW] Notification clicked:', event.notification.tag);
  
  event.notification.close();
  
  const notificationData = event.notification.data || {};
  let targetUrl = '/dashboard';
  
  // 根據通知類型決定目標頁面
  if (event.action) {
    switch (event.action) {
      case 'view':
        if (notificationData.resourceType === 'inquiry') {
          targetUrl = `/inquiries/${notificationData.resourceId}`;
        } else if (notificationData.resourceType === 'quote') {
          targetUrl = `/quotes/${notificationData.resourceId}`;
        }
        break;
      case 'reply':
        targetUrl = `/inquiries/${notificationData.resourceId}/reply`;
        break;
      case 'approve':
        targetUrl = `/quotes/${notificationData.resourceId}/approve`;
        break;
      case 'review':
        targetUrl = `/quotes/${notificationData.resourceId}`;
        break;
      default:
        targetUrl = notificationData.actionUrl || '/dashboard';
    }
  } else if (notificationData.actionUrl) {
    targetUrl = notificationData.actionUrl;
  }
  
  event.waitUntil(
    clients.matchAll({ type: 'window', includeUncontrolled: true })
      .then((clientList) => {
        // 嘗試找到現有的窗口
        for (const client of clientList) {
          if (client.url.includes(self.location.origin) && 'focus' in client) {
            client.postMessage({
              type: 'NAVIGATE',
              url: targetUrl
            });
            return client.focus();
          }
        }
        
        // 沒有現有窗口，打開新窗口
        if (clients.openWindow) {
          return clients.openWindow(targetUrl);
        }
      })
  );
  
  // 發送點擊事件到分析系統
  if (notificationData.notificationId) {
    fetch('/api/mobile/notifications/' + notificationData.notificationId + '/clicked', {
      method: 'POST'
    }).catch(() => {
      // 忽略錯誤，不影響用戶體驗
    });
  }
});

// 處理通知關閉
self.addEventListener('notificationclose', (event) => {
  console.log('[SW] Notification closed:', event.notification.tag);
  
  // 可以在這裡記錄通知被關閉的分析數據
});

// 處理後台同步
self.addEventListener('sync', (event) => {
  console.log('[SW] Background sync:', event.tag);
  
  if (event.tag === 'offline-data-sync') {
    event.waitUntil(syncOfflineData());
  }
});

// 處理來自主線程的消息
self.addEventListener('message', (event) => {
  console.log('[SW] Message received:', event.data);
  
  if (event.data && event.data.type === 'SKIP_WAITING') {
    self.skipWaiting();
  } else if (event.data && event.data.type === 'CACHE_RESOURCES') {
    event.waitUntil(cacheResources(event.data.resources));
  }
});

// 快取策略實現

// 網路優先策略 (適用於 API 請求)
async function networkFirstStrategy(request) {
  try {
    const response = await fetch(request);
    
    if (response.ok) {
      const cache = await caches.open(DYNAMIC_CACHE_NAME);
      cache.put(request, response.clone());
    }
    
    return response;
  } catch (error) {
    console.log('[SW] Network failed, trying cache:', request.url);
    const cachedResponse = await caches.match(request);
    
    if (cachedResponse) {
      return cachedResponse;
    }
    
    // 返回離線頁面或錯誤響應
    if (request.mode === 'navigate') {
      return caches.match('/offline');
    }
    
    return new Response(
      JSON.stringify({ error: 'Network unavailable', offline: true }),
      {
        status: 503,
        statusText: 'Service Unavailable',
        headers: { 'Content-Type': 'application/json' }
      }
    );
  }
}

// 快取優先策略 (適用於靜態資源)
async function cacheFirstStrategy(request) {
  const cachedResponse = await caches.match(request);
  
  if (cachedResponse) {
    return cachedResponse;
  }
  
  try {
    const response = await fetch(request);
    
    if (response.ok) {
      const cache = await caches.open(DYNAMIC_CACHE_NAME);
      cache.put(request, response.clone());
    }
    
    return response;
  } catch (error) {
    console.log('[SW] Cache and network failed:', request.url);
    
    // 返回預設的離線資源
    if (request.destination === 'image') {
      return caches.match('/icons/offline-image.png');
    }
    
    return new Response('Resource unavailable offline', {
      status: 503,
      statusText: 'Service Unavailable'
    });
  }
}

// 導航策略 (適用於頁面請求)
async function navigationStrategy(request) {
  try {
    const response = await fetch(request);
    
    if (response.ok) {
      const cache = await caches.open(DYNAMIC_CACHE_NAME);
      cache.put(request, response.clone());
    }
    
    return response;
  } catch (error) {
    console.log('[SW] Navigation failed, trying cache:', request.url);
    
    // 嘗試從快取獲取
    const cachedResponse = await caches.match(request);
    if (cachedResponse) {
      return cachedResponse;
    }
    
    // 嘗試獲取通用的離線頁面
    const offlinePage = await caches.match('/offline');
    if (offlinePage) {
      return offlinePage;
    }
    
    // 返回基本的離線響應
    return new Response(`
      <!DOCTYPE html>
      <html>
        <head>
          <title>FastenMind - 離線</title>
          <meta name="viewport" content="width=device-width, initial-scale=1">
          <style>
            body { font-family: Arial, sans-serif; text-align: center; padding: 50px; }
            .offline-message { color: #666; margin: 20px 0; }
          </style>
        </head>
        <body>
          <h1>FastenMind</h1>
          <div class="offline-message">
            <p>目前無法連接到網路</p>
            <p>請檢查您的網路連接後重試</p>
          </div>
          <button onclick="window.location.reload()">重試</button>
        </body>
      </html>
    `, {
      headers: { 'Content-Type': 'text/html' }
    });
  }
}

// 同步離線資料
async function syncOfflineData() {
  try {
    console.log('[SW] Syncing offline data');
    
    // 獲取離線存儲的資料
    const offlineData = await getOfflineData();
    
    for (const data of offlineData) {
      try {
        const response = await fetch('/api/mobile/offline-data/sync', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify(data)
        });
        
        if (response.ok) {
          await removeOfflineData(data.id);
          console.log('[SW] Synced offline data:', data.id);
        }
      } catch (error) {
        console.error('[SW] Failed to sync data:', data.id, error);
      }
    }
  } catch (error) {
    console.error('[SW] Offline data sync failed:', error);
  }
}

// 快取額外資源
async function cacheResources(resources) {
  const cache = await caches.open(DYNAMIC_CACHE_NAME);
  
  try {
    await cache.addAll(resources);
    console.log('[SW] Cached additional resources:', resources.length);
  } catch (error) {
    console.error('[SW] Failed to cache resources:', error);
  }
}

// 獲取離線存儲的資料
async function getOfflineData() {
  // 這裡應該從 IndexedDB 或其他本地存儲獲取離線資料
  // 暫時返回空數組
  return [];
}

// 移除已同步的離線資料
async function removeOfflineData(id) {
  // 這裡應該從本地存儲移除已同步的資料
  console.log('[SW] Removing synced offline data:', id);
}

// 工具函數

function isApiRequest(url) {
  return API_ENDPOINTS.some(endpoint => url.includes(endpoint));
}

function isStaticAsset(url) {
  const staticExtensions = ['.js', '.css', '.png', '.jpg', '.jpeg', '.gif', '.svg', '.woff', '.woff2'];
  return staticExtensions.some(ext => url.includes(ext)) || url.includes('_next/static');
}

// PWA 安裝提示處理
self.addEventListener('beforeinstallprompt', (event) => {
  console.log('[SW] Before install prompt');
  event.preventDefault();
  
  // 發送消息給主線程，讓應用決定何時顯示安裝提示
  self.clients.matchAll().then((clients) => {
    clients.forEach((client) => {
      client.postMessage({
        type: 'INSTALL_PROMPT_AVAILABLE',
        event: event
      });
    });
  });
});

// 定期清理快取
setInterval(() => {
  caches.keys().then((cacheNames) => {
    cacheNames.forEach((cacheName) => {
      if (cacheName.startsWith('fastenmind-dynamic-')) {
        caches.open(cacheName).then((cache) => {
          // 清理超過一週的動態快取
          // 這裡可以實現更複雜的清理邏輯
        });
      }
    });
  });
}, 24 * 60 * 60 * 1000); // 每24小時執行一次