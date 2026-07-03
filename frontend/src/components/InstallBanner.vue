<template>
  <Transition name="install-slide">
    <div v-if="isInstallable" class="install-banner" role="complementary" aria-label="Install aplikasi">
      <div class="install-banner__icon">
        <img src="/logo_emagang.png" alt="e-Magang" />
      </div>

      <div class="install-banner__body">
        <p class="install-banner__title">Install e-Magang</p>
        <p class="install-banner__sub">Akses penuh, notifikasi, & bisa dibuka tanpa browser</p>
      </div>

      <button class="install-banner__btn" @click="install">Install</button>

      <button class="install-banner__close" @click="dismiss" aria-label="Tutup">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none">
          <path d="M18 6L6 18M6 6l12 12" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"/>
        </svg>
      </button>
    </div>
  </Transition>
</template>

<script setup>
import { useInstallPrompt } from '@/composables/useInstallPrompt'

const { isInstallable, install, dismiss } = useInstallPrompt()
</script>

<style scoped>
.install-banner {
  position: fixed;
  bottom: 20px;
  left: 50%;
  transform: translateX(-50%);
  z-index: 9999;
  display: flex;
  align-items: center;
  gap: 12px;
  background: #fff;
  border: 1px solid #e5e7eb;
  border-radius: 16px;
  padding: 12px 16px;
  box-shadow: 0 8px 32px rgba(0,0,0,0.14), 0 2px 8px rgba(0,0,0,0.08);
  min-width: 300px;
  max-width: calc(100vw - 40px);
}

.install-banner__icon {
  flex-shrink: 0;
  width: 40px;
  height: 40px;
  border-radius: 10px;
  overflow: hidden;
  background: #f0fdf4;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px solid #bbf7d0;
}

.install-banner__icon img {
  width: 32px;
  height: 32px;
  object-fit: contain;
}

.install-banner__body {
  flex: 1;
  min-width: 0;
}

.install-banner__title {
  font-size: 13.5px;
  font-weight: 700;
  color: #111827;
  margin: 0 0 2px;
  white-space: nowrap;
}

.install-banner__sub {
  font-size: 11.5px;
  color: #6b7280;
  margin: 0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.install-banner__btn {
  flex-shrink: 0;
  background: #16a34a;
  color: #fff;
  font-size: 12.5px;
  font-weight: 700;
  border: none;
  border-radius: 8px;
  padding: 7px 16px;
  cursor: pointer;
  transition: background 0.15s;
  white-space: nowrap;
}

.install-banner__btn:hover {
  background: #15803d;
}

.install-banner__close {
  flex-shrink: 0;
  background: none;
  border: none;
  color: #9ca3af;
  cursor: pointer;
  padding: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 6px;
  transition: color 0.15s, background 0.15s;
}

.install-banner__close:hover {
  color: #374151;
  background: #f3f4f6;
}

/* Animasi masuk dari bawah */
.install-slide-enter-active,
.install-slide-leave-active {
  transition: opacity 0.3s ease, transform 0.3s ease;
}

.install-slide-enter-from,
.install-slide-leave-to {
  opacity: 0;
  transform: translateX(-50%) translateY(24px);
}

/* Mobile: full width hampir */
@media (max-width: 480px) {
  .install-banner {
    bottom: 16px;
    left: 16px;
    right: 16px;
    transform: none;
    min-width: unset;
    max-width: unset;
  }

  .install-slide-enter-from,
  .install-slide-leave-to {
    opacity: 0;
    transform: translateY(24px);
  }
}
</style>
