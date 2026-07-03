<template>
  <Teleport to="body">
    <div class="notif-toasts" aria-live="polite">
      <TransitionGroup name="toast" tag="div" class="notif-toasts__list">
        <div
          v-for="t in toasts"
          :key="t.id"
          class="notif-toast"
          :class="{ 'notif-toast--clickable': !!t.route }"
          role="alert"
          @click="handleClick(t)"
        >
          <div class="notif-toast__icon">
            <img src="/logo_emagang.png" alt="e-Magang" />
          </div>
          <div class="notif-toast__body">
            <div class="notif-toast__title">{{ t.title }}</div>
            <div v-if="t.body" class="notif-toast__msg">{{ t.body }}</div>
          </div>
          <button
            class="notif-toast__close"
            @click.stop="$emit('dismiss', t.id)"
            aria-label="Tutup"
          >
            <svg width="12" height="12" viewBox="0 0 24 24" fill="none">
              <path d="M18 6L6 18M6 6l12 12" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"/>
            </svg>
          </button>
        </div>
      </TransitionGroup>
    </div>
  </Teleport>
</template>

<script setup>
const props = defineProps({
  toasts: {
    type: Array,
    default: () => [],
  },
})
const emit = defineEmits(['dismiss', 'navigate'])

function handleClick(t) {
  if (t.route) emit('navigate', t.route)
  emit('dismiss', t.id)
}
</script>

<style scoped>
.notif-toasts {
  position: fixed;
  top: 16px;
  right: 16px;
  z-index: 9999;
  display: flex;
  flex-direction: column;
  gap: 10px;
  pointer-events: none;
  max-width: min(360px, calc(100vw - 32px));
}

.notif-toasts__list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.notif-toast {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  padding: 12px 14px;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15), 0 1px 4px rgba(0, 0, 0, 0.08);
  border-left: 3px solid #48AF4A;
  pointer-events: all;
  cursor: default;
  min-width: 260px;
  max-width: 100%;
  transition: transform 0.15s, box-shadow 0.15s;
}

.notif-toast--clickable {
  cursor: pointer;
}

.notif-toast--clickable:hover {
  transform: translateY(-1px);
  box-shadow: 0 6px 24px rgba(0, 0, 0, 0.18), 0 2px 6px rgba(0, 0, 0, 0.1);
}

.notif-toast__icon {
  flex-shrink: 0;
  width: 32px;
  height: 32px;
  border-radius: 8px;
  overflow: hidden;
  background: #f0fdf4;
  display: flex;
  align-items: center;
  justify-content: center;
}

.notif-toast__icon img {
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.notif-toast__body {
  flex: 1;
  min-width: 0;
}

.notif-toast__title {
  font-weight: 600;
  font-size: 13px;
  color: #111827;
  line-height: 1.3;
  margin-bottom: 2px;
  word-break: break-word;
}

.notif-toast__msg {
  font-size: 12px;
  color: #6b7280;
  line-height: 1.4;
  word-break: break-word;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.notif-toast__close {
  flex-shrink: 0;
  background: none;
  border: none;
  cursor: pointer;
  color: #9ca3af;
  padding: 2px;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: color 0.15s, background 0.15s;
  margin-top: -2px;
}

.notif-toast__close:hover {
  color: #374151;
  background: #f3f4f6;
}

/* Animasi masuk/keluar */
.toast-enter-active {
  animation: toast-in 0.28s ease;
}
.toast-leave-active {
  animation: toast-out 0.22s ease forwards;
}

@keyframes toast-in {
  from { opacity: 0; transform: translateX(40px) scale(0.95); }
  to   { opacity: 1; transform: translateX(0) scale(1); }
}

@keyframes toast-out {
  from { opacity: 1; transform: translateX(0) scale(1); }
  to   { opacity: 0; transform: translateX(40px) scale(0.9); }
}

/* ── Mobile: pindah ke bawah, full-width ──────────────────────── */
@media (max-width: 600px) {
  .notif-toasts {
    top: auto;
    bottom: 80px;   /* hindari nav bar bawah mobile */
    right: 12px;
    left: 12px;
    max-width: 100%;
  }

  .notif-toast {
    min-width: unset;
    border-radius: 10px;
  }

  @keyframes toast-in {
    from { opacity: 0; transform: translateY(20px) scale(0.97); }
    to   { opacity: 1; transform: translateY(0) scale(1); }
  }

  @keyframes toast-out {
    from { opacity: 1; transform: translateY(0) scale(1); }
    to   { opacity: 0; transform: translateY(20px) scale(0.95); }
  }
}
</style>
