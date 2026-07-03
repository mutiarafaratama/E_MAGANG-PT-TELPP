// Firebase Cloud Messaging Service Worker
// Mendukung notifikasi push di background (browser/PWA)
importScripts('https://www.gstatic.com/firebasejs/10.12.2/firebase-app-compat.js');
importScripts('https://www.gstatic.com/firebasejs/10.12.2/firebase-messaging-compat.js');

// ── Config default (public config — aman di-commit, bukan secret) ──────────
// Dipakai untuk background push saat app tertutup.
// Bisa di-override oleh app via postMessage jika berbeda.
const DEFAULT_CONFIG = {
  apiKey:            'AIzaSyCHxKRMV5CR-kLk_XLj9xImFfS1pzXxg94',
  authDomain:        'e-magang-pt-telpp.firebaseapp.com',
  projectId:         'e-magang-pt-telpp',
  storageBucket:     'e-magang-pt-telpp.firebasestorage.app',
  messagingSenderId: '531598128459',
  appId:             '1:531598128459:web:1676f7ca1bea26f7bb64ce',
};

const CACHE_NAME = 'fcm-config-v1';
const CACHE_KEY  = '/fcm-config-data';

// ── Inisialisasi Firebase di awal (wajib agar push handler terdaftar tepat waktu) ──
// Firebase messaging HARUS di-init saat evaluasi awal SW — bukan di dalam async handler.
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
        icon:     '/logo_emagang.png',
        badge:    '/logo_emagang.png',
        tag:      payload.data?.tipe ?? 'emagang-notif',
        renotify: true,
        data:     { route, ...payload.data },
      });
    });

    console.log('[FCM SW] ✅ Firebase aktif — project:', config.projectId);
  } catch (e) {
    console.warn('[FCM SW] ❌ Init error:', e);
  }
}

// Init langsung saat SW dievaluasi (sinkron) agar push/notificationclick handler
// terdaftar pada initial evaluation — ini yang mencegah SW warning di browser.
initFirebase(DEFAULT_CONFIG);

// ── Simpan config baru ke Cache API ───────────────────────────────────────
async function saveConfig(config) {
  try {
    const cache = await caches.open(CACHE_NAME);
    await cache.put(CACHE_KEY, new Response(JSON.stringify(config), {
      headers: { 'Content-Type': 'application/json' },
    }));
  } catch (_) { /* silent */ }
}

// ── Lifecycle ──────────────────────────────────────────────────────────────
self.addEventListener('install', (event) => {
  console.log('[FCM SW] install');
  self.skipWaiting();
});

self.addEventListener('activate', (event) => {
  console.log('[FCM SW] activate');
  event.waitUntil(self.clients.claim());
});

// ── Terima config dari main app (untuk update jika config berubah) ─────────
self.addEventListener('message', async (event) => {
  if (event.data?.type !== 'FIREBASE_CONFIG') return;
  const config = event.data.config;
  if (!config?.apiKey) return;

  // Simpan ke cache untuk referensi sesi berikutnya
  await saveConfig(config);
  console.log('[FCM SW] Config diterima dan disimpan — project:', config.projectId);
});

// ── Klik notifikasi → buka/fokus tab + navigasi ────────────────────────────
// Hanya fokus tab yang sudah berada di halaman app (dashboard/staff/admin) —
// halaman publik (landing, login, dll) tidak punya SW message listener,
// sehingga tidak bisa menerima NAVIGATE dan navigasi akan hilang.
const APP_PREFIXES = ['/dashboard', '/staff', '/admin'];

self.addEventListener('notificationclick', (event) => {
  event.notification.close();
  const route = event.notification.data?.route ?? '/';
  const fullUrl = self.location.origin + route;

  event.waitUntil(
    clients.matchAll({ type: 'window', includeUncontrolled: true }).then((list) => {
      const appClient = list.find((client) => {
        try {
          const pathname = new URL(client.url).pathname;
          return APP_PREFIXES.some(p => pathname.startsWith(p)) && 'focus' in client;
        } catch (_) { return false; }
      });

      if (appClient) {
        appClient.focus();
        appClient.postMessage({ type: 'NAVIGATE', route });
        return;
      }

      if (clients.openWindow) {
        return clients.openWindow(fullUrl);
      }
    })
  );
});
