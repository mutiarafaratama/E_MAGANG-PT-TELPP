import { initializeApp, getApps } from 'firebase/app'
import { getMessaging, getToken, onMessage } from 'firebase/messaging'
import api from '@/lib/api'
import { playNotifSound } from '@/services/sound'

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
  if (!firebaseConfig.apiKey || !firebaseConfig.projectId) {
    console.warn('[FCM] ❌ VITE_FIREBASE_* env vars tidak terset — cek .env.local frontend')
    return null
  }
  try {
    const app = getApps().length ? getApps()[0] : initializeApp(firebaseConfig)
    if (!_messaging) _messaging = getMessaging(app)
    return _messaging
  } catch (e) {
    console.warn('[FCM] ❌ Inisialisasi Firebase gagal:', e)
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
    console.log('[FCM] Izin notifikasi:', result)
    return result === 'granted'
  } catch (e) {
    console.warn('[FCM] ❌ requestPermission gagal:', e)
    return false
  }
}

// Kirim config Firebase ke service worker dan simpan ke cache-nya
async function sendConfigToSW(reg) {
  const config = firebaseConfig
  if (!config.apiKey) return

  // Kirim ke semua SW instance yang bisa dijangkau
  const targets = [reg.installing, reg.waiting, reg.active].filter(Boolean)
  targets.forEach((sw) => {
    try { sw.postMessage({ type: 'FIREBASE_CONFIG', config }) } catch (_) {}
  })

  // Juga kirim via controller (SW yang sudah aktif mengontrol halaman)
  if (navigator.serviceWorker.controller) {
    try {
      navigator.serviceWorker.controller.postMessage({ type: 'FIREBASE_CONFIG', config })
    } catch (_) {}
  }

  console.log('[FCM] Config Firebase dikirim ke service worker — project:', config.projectId)
}

async function registerSW() {
  try {
    const reg = await navigator.serviceWorker.register('/firebase-messaging-sw.js', {
      scope: '/',
      updateViaCache: 'none',
    })
    console.log('[FCM] Service worker terdaftar — state:', reg.active?.state ?? 'installing')

    await sendConfigToSW(reg)

    // Tunggu SW siap menerima pesan
    await navigator.serviceWorker.ready

    // Kirim ulang config setelah SW ready (pastikan sudah diterima)
    const readyReg = await navigator.serviceWorker.getRegistration('/firebase-messaging-sw.js')
    if (readyReg) await sendConfigToSW(readyReg)

    return reg
  } catch (e) {
    console.warn('[FCM] ❌ Service worker registration gagal:', e)
    return null
  }
}

export async function getFCMToken() {
  const messaging = getFirebaseMessaging()
  if (!messaging) return null

  if (!VAPID_KEY) {
    console.warn('[FCM] ❌ VITE_FIREBASE_VAPID_KEY tidak terset')
    return null
  }

  try {
    const swReg = await registerSW()
    console.log('[FCM] Meminta FCM token dari Firebase...')
    const token = await getToken(messaging, {
      vapidKey: VAPID_KEY,
      serviceWorkerRegistration: swReg ?? undefined,
    })
    if (token) {
      console.log('[FCM] ✅ Token didapat:', token.substring(0, 20) + '...')
    } else {
      console.warn('[FCM] ⚠️  Token kosong — pastikan izin notifikasi sudah granted')
    }
    return token || null
  } catch (e) {
    console.warn('[FCM] ❌ Gagal get token:', e)
    return null
  }
}

// Tampilkan notifikasi foreground dengan dukungan klik → navigasi
function showForegroundNotif(title, body, route) {
  if (Notification.permission !== 'granted') return
  playNotifSound()
  try {
    if ('serviceWorker' in navigator) {
      navigator.serviceWorker.ready.then((reg) => {
        reg.showNotification(title, {
          body,
          icon:     '/logotel.png',
          badge:    '/logotel.png',
          tag:      'emagang-foreground',
          renotify: true,
          data:     { route: route || '/' },
        })
      }).catch(() => _fallbackNotif(title, body, route))
    } else {
      _fallbackNotif(title, body, route)
    }
  } catch (e) {
    _fallbackNotif(title, body, route)
  }
}

function _fallbackNotif(title, body, route) {
  try {
    const n = new Notification(title, { body, icon: '/logotel.png' })
    if (route) {
      n.onclick = () => { window.focus(); window.location.href = route; n.close() }
    }
  } catch (_) { /* silent */ }
}

export async function setupPushNotifications() {
  console.log('[FCM] setupPushNotifications dipanggil')

  if (!isNotificationSupported()) {
    console.warn('[FCM] ⚠️  Browser tidak support push notification')
    return
  }

  if (!firebaseConfig.apiKey) {
    console.error('[FCM] ❌ VITE_FIREBASE_API_KEY tidak terset — push notif tidak bisa jalan.' +
      '\n    Untuk lokal: isi VITE_FIREBASE_* di frontend/.env.local' +
      '\n    Untuk Replit: pastikan secrets VITE_FIREBASE_* sudah diset')
    return
  }

  const currentPermission = Notification.permission
  console.log('[FCM] Status izin notifikasi saat ini:', currentPermission)

  if (currentPermission === 'denied') {
    console.warn('[FCM] ⚠️  Izin notifikasi ditolak user — tidak bisa lanjut')
    return
  }

  let granted = currentPermission === 'granted'
  if (!granted) {
    console.log('[FCM] Meminta izin notifikasi ke user...')
    granted = await requestPermission()
  }
  if (!granted) {
    console.warn('[FCM] ⚠️  User menolak izin notifikasi')
    return
  }

  const token = await getFCMToken()
  if (!token) {
    console.warn('[FCM] ⚠️  Tidak bisa dapat token — push notif tidak aktif')
    return
  }

  try {
    await api.post('/api/notifikasi/fcm-token', { token })
    console.log('[FCM] ✅ Token berhasil disimpan ke server')
  } catch (e) {
    console.warn('[FCM] ❌ Gagal kirim token ke server:', e?.response?.status, e?.message)
  }

  // Handle pesan foreground — guard agar listener tidak terdaftar ganda
  const messaging = getFirebaseMessaging()
  if (messaging && !_messageHandlerRegistered) {
    _messageHandlerRegistered = true
    onMessage(messaging, (payload) => {
      console.log('[FCM] 🔔 Pesan foreground diterima:', payload)
      const { title = 'Notifikasi Baru', body = '' } = payload.notification ?? {}
      const route = payload.data?.route ?? null
      showForegroundNotif(title, body, route)
    })
    console.log('[FCM] ✅ Foreground message listener aktif')
  }
}
