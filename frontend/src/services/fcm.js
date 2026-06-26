import { initializeApp, getApps } from 'firebase/app'
import { getMessaging, getToken, onMessage } from 'firebase/messaging'
import api from '@/lib/api'

const firebaseConfig = {
  apiKey:            import.meta.env.VITE_FIREBASE_API_KEY,
  authDomain:        import.meta.env.VITE_FIREBASE_AUTH_DOMAIN,
  projectId:         import.meta.env.VITE_FIREBASE_PROJECT_ID,
  storageBucket:     import.meta.env.VITE_FIREBASE_STORAGE_BUCKET,
  messagingSenderId: import.meta.env.VITE_FIREBASE_MESSAGING_SENDER_ID,
  appId:             import.meta.env.VITE_FIREBASE_APP_ID,
}

const VAPID_KEY = import.meta.env.VITE_FIREBASE_VAPID_KEY

let _messaging = null
let _messageHandlerRegistered = false

function getFirebaseMessaging() {
  if (!firebaseConfig.apiKey || !firebaseConfig.projectId) return null
  try {
    const app = getApps().length ? getApps()[0] : initializeApp(firebaseConfig)
    if (!_messaging) _messaging = getMessaging(app)
    return _messaging
  } catch (e) {
    console.warn('[FCM] Inisialisasi gagal:', e)
    return null
  }
}

export function isNotificationSupported() {
  return (
    typeof window !== 'undefined' &&
    'Notification' in window &&
    'serviceWorker' in navigator &&
    'PushManager' in window
  )
}

export async function requestPermission() {
  if (!isNotificationSupported()) return false
  try {
    const result = await Notification.requestPermission()
    return result === 'granted'
  } catch (e) {
    return false
  }
}

async function registerSW() {
  try {
    const reg = await navigator.serviceWorker.register('/firebase-messaging-sw.js', {
      scope: '/',
    })
    // Kirim konfigurasi Firebase ke service worker (untuk override jika perlu)
    const target = reg.installing || reg.waiting || reg.active
    if (target) {
      target.postMessage({ type: 'FIREBASE_CONFIG', config: firebaseConfig })
    }
    await navigator.serviceWorker.ready
    return reg
  } catch (e) {
    console.warn('[FCM] Service worker registration gagal:', e)
    return null
  }
}

export async function getFCMToken() {
  const messaging = getFirebaseMessaging()
  if (!messaging || !VAPID_KEY) return null
  try {
    const swReg = await registerSW()
    const token = await getToken(messaging, {
      vapidKey: VAPID_KEY,
      serviceWorkerRegistration: swReg ?? undefined,
    })
    return token || null
  } catch (e) {
    console.warn('[FCM] Gagal get token:', e)
    return null
  }
}

// Tampilkan notifikasi foreground dengan dukungan klik → navigasi
function showForegroundNotif(title, body, route) {
  if (Notification.permission !== 'granted') return
  try {
    if ('serviceWorker' in navigator) {
      navigator.serviceWorker.ready.then((reg) => {
        reg.showNotification(title, {
          body,
          icon: '/logotel.png',
          badge: '/logotel.png',
          tag: 'emagang-foreground',
          renotify: true,
          data: { route: route || '/' },
        })
      }).catch(() => {
        // Fallback ke Notification API
        _fallbackNotif(title, body, route)
      })
    } else {
      _fallbackNotif(title, body, route)
    }
  } catch (e) {
    _fallbackNotif(title, body, route)
  }
}

function _fallbackNotif(title, body, route) {
  try {
    const n = new Notification(title, {
      body,
      icon: '/logotel.png',
    })
    if (route) {
      n.onclick = () => {
        window.focus()
        window.location.href = route
        n.close()
      }
    }
  } catch (e) { /* silent */ }
}

export async function setupPushNotifications() {
  if (!isNotificationSupported()) return
  if (!firebaseConfig.apiKey) return

  const currentPermission = Notification.permission
  if (currentPermission === 'denied') return

  let granted = currentPermission === 'granted'
  if (!granted) {
    granted = await requestPermission()
  }
  if (!granted) return

  const token = await getFCMToken()
  if (!token) return

  try {
    await api.post('/api/notifikasi/fcm-token', { token })
    console.log('[FCM] Token berhasil disimpan ke server')
  } catch (e) {
    console.warn('[FCM] Gagal kirim token ke server:', e)
  }

  // Handle pesan foreground — guard agar listener tidak terdaftar ganda
  // saat DashboardLayout remount (navigasi antar halaman)
  const messaging = getFirebaseMessaging()
  if (messaging && !_messageHandlerRegistered) {
    _messageHandlerRegistered = true
    onMessage(messaging, (payload) => {
      const { title = 'Notifikasi Baru', body = '' } = payload.notification ?? {}
      const route = payload.data?.route ?? null
      showForegroundNotif(title, body, route)
    })
  }
}
