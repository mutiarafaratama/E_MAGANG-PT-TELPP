import { ref, onMounted, onUnmounted } from 'vue'

const DISMISS_KEY = 'install_prompt_dismissed_until'
const DISMISS_DAYS = 7

const deferredPrompt = ref(null)
const isInstallable = ref(false)
const isInstalled = ref(false)

function isStandalone() {
  return (
    window.matchMedia('(display-mode: standalone)').matches ||
    window.navigator.standalone === true
  )
}

function isDismissed() {
  const until = localStorage.getItem(DISMISS_KEY)
  if (!until) return false
  return Date.now() < parseInt(until, 10)
}

export function useInstallPrompt() {
  function setup() {
    if (isStandalone()) {
      isInstalled.value = true
      return
    }

    const handler = (e) => {
      e.preventDefault()
      deferredPrompt.value = e
      if (!isDismissed()) {
        isInstallable.value = true
      }
    }

    window.addEventListener('beforeinstallprompt', handler)
    window.addEventListener('appinstalled', () => {
      isInstallable.value = false
      isInstalled.value = true
      deferredPrompt.value = null
    })

    onUnmounted(() => {
      window.removeEventListener('beforeinstallprompt', handler)
    })
  }

  async function install() {
    if (!deferredPrompt.value) return
    deferredPrompt.value.prompt()
    const { outcome } = await deferredPrompt.value.userChoice
    if (outcome === 'accepted') {
      isInstallable.value = false
    }
    deferredPrompt.value = null
  }

  function dismiss() {
    isInstallable.value = false
    const until = Date.now() + DISMISS_DAYS * 24 * 60 * 60 * 1000
    localStorage.setItem(DISMISS_KEY, String(until))
  }

  onMounted(setup)

  return { isInstallable, isInstalled, install, dismiss }
}
