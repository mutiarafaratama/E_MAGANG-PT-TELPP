<template>
  <!-- Backdrop: klik di luar widget untuk tutup -->
  <teleport to="body">
    <div
      v-if="isOpen"
      class="chat-fab-backdrop"
      @click="closeWidget"
      @touchstart.passive="closeWidget"
    />
  </teleport>

  <div
    ref="fabRef"
    class="chat-fab"
    :class="{ 'chat-fab--open': isOpen, 'chat-fab--dragging': isDragging }"
    :style="fabStyle"
    @mousedown="onDragStart"
    @touchstart.passive="onTouchStart"
  >
    <!-- FAB Button: disembunyikan di mobile saat widget terbuka -->
    <button
      v-show="!isOpen || !isMobile"
      class="chat-fab__btn"
      @click.stop="toggleChat"
      :title="isOpen ? 'Tutup Chat' : 'Chat Bantuan'"
    >
      <span class="chat-fab__badge" v-if="unreadCount > 0">{{ unreadCount > 99 ? '99+' : unreadCount }}</span>
      <img v-if="!isOpen" src="/logo_chatbot.png" alt="Chat" class="chat-fab__icon-img" />
      <svg v-else width="22" height="22" viewBox="0 0 24 24" fill="none">
        <line x1="18" y1="6" x2="6" y2="18" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"/>
        <line x1="6" y1="6" x2="18" y2="18" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"/>
      </svg>
    </button>

    <!-- Chat Widget -->
    <div v-if="isOpen" class="chat-fab__widget" @click.stop @mousedown.stop @touchstart.stop>
      <slot />
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onBeforeUnmount, onMounted } from 'vue'

const props = defineProps({
  unreadCount: { type: Number, default: 0 },
  modelValue: { type: Boolean, default: false }
})

const emit = defineEmits(['update:modelValue'])

const fabRef = ref(null)
const isOpen = computed({
  get: () => props.modelValue,
  set: (v) => emit('update:modelValue', v)
})

const isDragging = ref(false)
const STORAGE_KEY = 'chat_fab_pos'

// Deteksi mobile
const isMobile = ref(false)
function checkMobile() { isMobile.value = window.innerWidth < 640 }
onMounted(() => {
  checkMobile()
  window.addEventListener('resize', checkMobile)
})
onBeforeUnmount(() => {
  window.removeEventListener('resize', checkMobile)
})

const defaultPos = () => ({ right: 24, bottom: 24, left: null, top: null })
const pos = ref(loadPos())

function loadPos() {
  try {
    const saved = localStorage.getItem(STORAGE_KEY)
    if (saved) return JSON.parse(saved)
  } catch {}
  return defaultPos()
}

function savePos() {
  localStorage.setItem(STORAGE_KEY, JSON.stringify(pos.value))
}

const fabStyle = computed(() => {
  const s = { position: 'fixed', zIndex: 9999 }
  if (pos.value.left !== null) s.left = pos.value.left + 'px'
  else s.right = (pos.value.right ?? 24) + 'px'
  if (pos.value.top !== null) s.top = pos.value.top + 'px'
  else s.bottom = (pos.value.bottom ?? 24) + 'px'
  return s
})

let startX = 0, startY = 0, startLeft = 0, startTop = 0, hasMoved = false

function getAbsolutePos() {
  const el = fabRef.value
  if (!el) return { left: window.innerWidth - 80, top: window.innerHeight - 80 }
  const rect = el.getBoundingClientRect()
  return { left: rect.left, top: rect.top }
}

function onDragStart(e) {
  if (e.target.closest('.chat-fab__widget')) return
  if (e.target.closest('.chat-fab__btn') && !isDragging.value) return
  isDragging.value = false
  hasMoved = false
  const absPos = getAbsolutePos()
  startX = e.clientX; startY = e.clientY
  startLeft = absPos.left; startTop = absPos.top
  document.addEventListener('mousemove', onDragMove)
  document.addEventListener('mouseup', onDragEnd)
}

function onTouchStart(e) {
  if (e.target.closest('.chat-fab__widget')) return
  const touch = e.touches[0]
  hasMoved = false
  const absPos = getAbsolutePos()
  startX = touch.clientX; startY = touch.clientY
  startLeft = absPos.left; startTop = absPos.top
  document.addEventListener('touchmove', onTouchMove, { passive: false })
  document.addEventListener('touchend', onTouchEnd)
}

function onDragMove(e) {
  const dx = e.clientX - startX, dy = e.clientY - startY
  if (Math.abs(dx) > 4 || Math.abs(dy) > 4) { isDragging.value = true; hasMoved = true }
  if (!isDragging.value) return
  applyDrag(startLeft + dx, startTop + dy)
}

function onTouchMove(e) {
  e.preventDefault()
  const touch = e.touches[0]
  const dx = touch.clientX - startX, dy = touch.clientY - startY
  if (Math.abs(dx) > 4 || Math.abs(dy) > 4) { isDragging.value = true; hasMoved = true }
  if (!isDragging.value) return
  applyDrag(startLeft + dx, startTop + dy)
}

function applyDrag(left, top) {
  const el = fabRef.value
  const w = el ? el.offsetWidth : 56, h = el ? el.offsetHeight : 56
  left = Math.max(8, Math.min(left, window.innerWidth - w - 8))
  top = Math.max(8, Math.min(top, window.innerHeight - h - 8))
  pos.value = { left, top, right: null, bottom: null }
}

function snapToEdge() {
  const el = fabRef.value
  if (!el) return
  const rect = el.getBoundingClientRect()
  const centerX = rect.left + rect.width / 2
  const snapRight = centerX > window.innerWidth / 2
  if (window.innerWidth < 640) {
    const top = Math.max(8, Math.min(rect.top, window.innerHeight - rect.height - 8))
    pos.value = snapRight
      ? { left: null, right: 8, top, bottom: null }
      : { left: 8, right: null, top, bottom: null }
  }
  savePos()
}

function onDragEnd() {
  document.removeEventListener('mousemove', onDragMove)
  document.removeEventListener('mouseup', onDragEnd)
  if (isDragging.value) snapToEdge()
  setTimeout(() => { isDragging.value = false }, 50)
}

function onTouchEnd() {
  document.removeEventListener('touchmove', onTouchMove)
  document.removeEventListener('touchend', onTouchEnd)
  if (isDragging.value) snapToEdge()
  setTimeout(() => { isDragging.value = false }, 50)
}

function toggleChat() {
  if (!isDragging.value && !hasMoved) isOpen.value = !isOpen.value
}

function closeWidget() {
  isOpen.value = false
}

onBeforeUnmount(() => {
  document.removeEventListener('mousemove', onDragMove)
  document.removeEventListener('mouseup', onDragEnd)
  document.removeEventListener('touchmove', onTouchMove)
  document.removeEventListener('touchend', onTouchEnd)
})
</script>

<style scoped>
/* Backdrop global untuk klik di luar */
:global(.chat-fab-backdrop) {
  position: fixed;
  inset: 0;
  z-index: 9990;
  background: transparent;
  cursor: default;
}

.chat-fab {
  user-select: none;
  touch-action: none;
}

.chat-fab__btn {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  background: linear-gradient(135deg, #16a34a, #15803d);
  color: white;
  border: none;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4px 24px rgba(22,163,74,0.45), 0 1px 4px rgba(0,0,0,0.12);
  position: relative;
  transition: background 0.2s, transform 0.2s, box-shadow 0.2s;
  -webkit-tap-highlight-color: transparent;
  z-index: 9999;
}

.chat-fab--dragging .chat-fab__btn {
  cursor: grabbing;
  transform: scale(1.08);
  box-shadow: 0 8px 32px rgba(22,163,74,0.55);
}

.chat-fab__btn:hover {
  background: linear-gradient(135deg, #15803d, #166534);
  transform: scale(1.06);
  box-shadow: 0 6px 28px rgba(22,163,74,0.5);
}

.chat-fab__icon-img {
  width: 34px;
  height: 34px;
  object-fit: contain;
  filter: drop-shadow(0 1px 2px rgba(0,0,0,0.15));
  border-radius: 50%;
}

.chat-fab__badge {
  position: absolute;
  top: -4px;
  right: -4px;
  background: #ef4444;
  color: white;
  font-size: 10px;
  font-weight: 700;
  min-width: 18px;
  height: 18px;
  border-radius: 9px;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0 4px;
  border: 2px solid white;
  line-height: 1;
  animation: badge-pop 0.3s ease;
}

@keyframes badge-pop {
  0% { transform: scale(0); }
  70% { transform: scale(1.2); }
  100% { transform: scale(1); }
}

.chat-fab__widget {
  position: fixed;
  bottom: 88px;
  right: 24px;
  width: 368px;
  height: min(560px, calc(100dvh - 100px));
  background: white;
  border-radius: 18px;
  box-shadow: 0 12px 48px rgba(0,0,0,0.18), 0 2px 8px rgba(0,0,0,0.08);
  overflow: hidden;
  display: flex;
  flex-direction: column;
  animation: widget-in 0.28s cubic-bezier(0.34,1.56,0.64,1);
  transform-origin: bottom right;
  z-index: 9998;
}

@keyframes widget-in {
  from { opacity: 0; transform: scale(0.88) translateY(12px); }
  to { opacity: 1; transform: scale(1) translateY(0); }
}

@media (max-width: 480px) {
  .chat-fab__widget {
    width: calc(100vw - 16px);
    right: 8px;
    left: 8px;
    bottom: 80px;
    height: min(480px, calc(100dvh - 110px));
    border-radius: 14px;
  }

  @keyframes widget-in {
    from { opacity: 0; transform: scale(0.94) translateY(10px); }
    to { opacity: 1; transform: scale(1) translateY(0); }
  }
}
</style>
