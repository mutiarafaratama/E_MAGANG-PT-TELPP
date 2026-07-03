import { ref } from 'vue'
import { getToken, onMessage } from 'firebase/messaging'
import { messaging } from '@/lib/firebase'
import api from '@/lib/api'

const fcmToken = ref(null)
const fcmPermission = ref(
  typeof Notification !== 'undefined' ? Notification.permission : 'default'
)

const VAPID_KEY = import.meta.env.VITE_FIREBASE_VAPID_KEY

export function useFCM() {
  async function requestPermission() {
    if (!('Notification' in window) || !('serviceWorker' in navigator)) return false
    const permission = await Notification.requestPermission()
    fcmPermission.value = permission
    return permission === 'granted'
  }

  async function initToken() {
    try {
      const swReg = await navigator.serviceWorker.ready
      const token = await getToken(messaging, {
        vapidKey: VAPID_KEY,
        serviceWorkerRegistration: swReg,
      })
      if (token) {
        fcmToken.value = token
        // Kirim token ke backend agar bisa push notif
        try { await api.post('/api/notifikasi/fcm-token', { token }) } catch {}
        return token
      }
    } catch (err) {
      console.warn('[FCM] Gagal ambil token:', err)
    }
    return null
  }

  // Tangkap notif saat app terbuka (foreground)
  function listenForeground(onNotif) {
    return onMessage(messaging, (payload) => {
      const { title, body } = payload.notification || {}
      if (onNotif) {
        onNotif({ title, body, data: payload.data })
      } else if (title && Notification.permission === 'granted') {
        new Notification(title, { body: body || '', icon: '/logo_emagang.png' })
      }
    })
  }

  // Panggil ini sekali saat user login (di dashboard masing-masing)
  async function init(onNotif) {
    const granted = await requestPermission()
    if (!granted) return null
    await initToken()
    return listenForeground(onNotif)
  }

  return {
    fcmToken,
    fcmPermission,
    requestPermission,
    initToken,
    listenForeground,
    init,
  }
}
