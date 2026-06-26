// Firebase Cloud Messaging Service Worker
// Mendukung notifikasi push di background (browser/PWA)
importScripts('https://www.gstatic.com/firebasejs/10.12.2/firebase-app-compat.js');
importScripts('https://www.gstatic.com/firebasejs/10.12.2/firebase-messaging-compat.js');

// Konfigurasi Firebase (public config — aman di service worker)
const FIREBASE_CONFIG = {
  apiKey: 'AIzaSyCR3jFTC6yKekEt41cmOWJA7h2EXIgXuLI',
  authDomain: 'e-magangtelpp.firebaseapp.com',
  projectId: 'e-magangtelpp',
  storageBucket: 'e-magangtelpp.firebasestorage.app',
  messagingSenderId: '443196985548',
  appId: '1:443196985548:web:b39c97d2115d307c53f269',
};

function initFirebase(config) {
  try {
    if (!firebase.apps.length) {
      firebase.initializeApp(config);
    }
    const messaging = firebase.messaging();

    messaging.onBackgroundMessage((payload) => {
      const { title = 'Notifikasi Baru', body = '' } = payload.notification ?? {};
      const route = payload.data?.route ?? '/';
      self.registration.showNotification(title, {
        body,
        icon: '/logotel.png',
        badge: '/logotel.png',
        tag: payload.data?.tipe ?? 'emagang-notif',
        renotify: true,
        data: { route, ...payload.data },
      });
    });
  } catch (e) {
    console.warn('[FCM SW] init error:', e);
  }
}

// Inisialisasi Firebase saat SW aktif pertama kali
// (diperlukan agar background notif bekerja walau app tertutup)
initFirebase(FIREBASE_CONFIG);

// Terima konfigurasi override dari main app jika berbeda
self.addEventListener('message', (event) => {
  if (event.data?.type === 'FIREBASE_CONFIG' && event.data.config?.apiKey) {
    // Reinit hanya jika config berbeda (misal env lokal)
    if (event.data.config.projectId !== FIREBASE_CONFIG.projectId) {
      initFirebase(event.data.config);
    }
  }
});

// Saat user klik notifikasi → buka/fokus tab dan navigasi ke route
self.addEventListener('notificationclick', (event) => {
  event.notification.close();
  const route = event.notification.data?.route ?? '/';
  event.waitUntil(
    clients.matchAll({ type: 'window', includeUncontrolled: true }).then((windowClients) => {
      // Cari tab yang sudah terbuka di domain yang sama
      for (const client of windowClients) {
        const url = new URL(client.url);
        if (url.pathname !== '/login' && 'focus' in client) {
          client.focus();
          client.postMessage({ type: 'NAVIGATE', route });
          return;
        }
      }
      // Tidak ada tab — buka baru
      if (clients.openWindow) {
        return clients.openWindow(route);
      }
    })
  );
});
