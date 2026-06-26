import { ref } from 'vue'
import { getToken, onMessage } from 'firebase/messaging'
import { getFirebaseMessaging } from '@/lib/firebase'
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
      const messaging = await getFirebaseMessaging()
      if (!messaging) return null

      const swReg = await navigator.serviceWorker.ready
      const token = await getToken(messaging, {
        vapidKey: VAPID_KEY,
        serviceWorkerRegistration: swReg,
      })
      if (token) {
        fcmToken.value = token
        try { await api.post('/api/notifikasi/fcm-token', { token }) } catch {}
        return token
      }
    } catch (err) {
      console.warn('[FCM] Gagal ambil token:', err)
    }
    return null
  }

  async function listenForeground(onNotif) {
    const messaging = await getFirebaseMessaging()
    if (!messaging) return () => {}
    return onMessage(messaging, (payload) => {
      const { title, body } = payload.notification || {}
      if (onNotif) {
        onNotif({ title, body, data: payload.data })
      } else if (title && Notification.permission === 'granted') {
        new Notification(title, { body: body || '', icon: '/logotel.png' })
      }
    })
  }

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
