// Service Worker e-Magang PT TELPP — PWA + FCM

importScripts('https://www.gstatic.com/firebasejs/11.0.0/firebase-app-compat.js')
importScripts('https://www.gstatic.com/firebasejs/11.0.0/firebase-messaging-compat.js')

const CACHE_NAME = 'emagang-v2'
const STATIC_ASSETS = ['/', '/manifest.json']

// ── Firebase init ──────────────────────────────────────────────────────────────
firebase.initializeApp({
  apiKey: 'AIzaSyCR3jFTC6yKekEt41cmOWJA7h2EXIgXuLI',
  authDomain: 'e-magangtelpp.firebaseapp.com',
  projectId: 'e-magangtelpp',
  storageBucket: 'e-magangtelpp.firebasestorage.app',
  messagingSenderId: '443196985548',
  appId: '1:443196985548:web:b39c97d2115d307c53f269',
})

const messaging = firebase.messaging()

// Notifikasi background (app tertutup / di tab lain)
messaging.onBackgroundMessage((payload) => {
  const title = payload.notification?.title || 'e-Magang PT TELPP'
  const body = payload.notification?.body || ''
  const route = payload.data?.route || '/'

  self.registration.showNotification(title, {
    body,
    icon: '/logotel.png',
    badge: '/logotel.png',
    data: { route },
    vibrate: [200, 100, 200],
    tag: 'emagang-notif',
    renotify: true,
  })
})

// ── Install ────────────────────────────────────────────────────────────────────
self.addEventListener('install', (event) => {
  event.waitUntil(
    caches.open(CACHE_NAME).then((cache) => cache.addAll(STATIC_ASSETS).catch(() => {}))
  )
  self.skipWaiting()
})

// ── Activate — bersihkan cache lama ───────────────────────────────────────────
self.addEventListener('activate', (event) => {
  event.waitUntil(
    caches.keys().then((keys) =>
      Promise.all(keys.filter((k) => k !== CACHE_NAME).map((k) => caches.delete(k)))
    )
  )
  self.clients.claim()
})

// ── Fetch — network first, fallback cache untuk navigasi ──────────────────────
self.addEventListener('fetch', (event) => {
  const url = new URL(event.request.url)

  // Skip API + uploads — selalu ke network
  if (url.pathname.startsWith('/api/') || url.pathname.startsWith('/uploads/')) return

  // Navigasi: network first, fallback ke cache
  if (event.request.mode === 'navigate') {
    event.respondWith(
      fetch(event.request).catch(() =>
        caches.match('/').then((r) => r || fetch(event.request))
      )
    )
  }
})

// ── Notification click — buka route yang relevan ──────────────────────────────
self.addEventListener('notificationclick', (event) => {
  event.notification.close()
  const route = event.notification.data?.route || '/'

  event.waitUntil(
    clients.matchAll({ type: 'window', includeUncontrolled: true }).then((clientList) => {
      for (const client of clientList) {
        if (client.url.includes(self.location.origin) && 'focus' in client) {
          client.focus()
          client.postMessage({ type: 'navigate', route })
          return
        }
      }
      if (clients.openWindow) return clients.openWindow(route)
    })
  )
})
