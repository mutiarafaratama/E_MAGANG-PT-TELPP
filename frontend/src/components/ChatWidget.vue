<template>
  <div class="cw">
    <!-- Header -->
    <div class="cw__header">
      <div class="cw__header-info">
        <div class="cw__avatar">
          <img src="/logo_chatbot.png" alt="Bot" class="cw__avatar-img" />
          <span class="cw__online-dot"></span>
        </div>
        <div>
          <div class="cw__name">Pusat Bantuan</div>
          <div class="cw__status" v-if="!activeTiket">Siap Membantu</div>
          <div class="cw__status" v-else-if="activeTiket.status === 'menunggu'">Menunggu HRD...</div>
          <div class="cw__status cw__status--active" v-else-if="activeTiket.status === 'diproses'">HRD Sedang Membalas</div>
          <div class="cw__status" v-else>Tiket {{ activeTiket.status }}</div>
        </div>
      </div>
      <div class="cw__timer" v-if="activeTiket && activeTiket.expires_at && activeTiket.status === 'menunggu'">
        <svg width="12" height="12" viewBox="0 0 24 24" fill="none"><circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2"/><polyline points="12 6 12 12 16 14" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>
        {{ timerLabel }}
      </div>
    </div>

    <!-- View: Chatbot NLP (belum ada tiket aktif) -->
    <template v-if="!activeTiket || activeTiket.status === 'hangus'">
      <div class="cw__messages" ref="msgContainerBot">
        <div v-for="(msg, i) in botMessages" :key="i" class="cw__msg" :class="'cw__msg--' + msg.from">
          <div class="cw__bubble">
            <span v-html="formatMsg(msg.text)"></span>
          </div>
          <div class="cw__time">{{ msg.time }}</div>
        </div>
        <div v-if="botTyping" class="cw__msg cw__msg--bot">
          <div class="cw__bubble cw__bubble--typing"><span></span><span></span><span></span></div>
        </div>
      </div>

      <!-- CTA Buat Tiket -->
      <div v-if="showTiketCta" class="cw__cta">
        <p>Ingin melanjutkan ke HRD?</p>
        <button class="cw__cta-btn" @click="showTiketForm = true">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none"><path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z" stroke="currentColor" stroke-width="2"/></svg>
          Hubungi HRD Langsung
        </button>
      </div>

      <!-- Form Buat Tiket -->
      <div v-if="showTiketForm" class="cw__form">
        <div class="cw__form-title">Buat Tiket ke HRD</div>
        <input v-model="tiketSubjek" class="cw__form-input" placeholder="Subjek pertanyaan..." maxlength="100" />
        <textarea v-model="tiketFormPesan" class="cw__form-textarea" placeholder="Jelaskan pertanyaan atau masalah Anda..." rows="3" maxlength="1000"></textarea>
        <select v-model="tiketKategori" class="cw__form-select">
          <option value="umum">Umum</option>
          <option value="absensi">Absensi</option>
          <option value="dokumen">Dokumen</option>
          <option value="sertifikat">Sertifikat</option>
        </select>
        <div class="cw__form-actions">
          <button class="cw__form-cancel" @click="showTiketForm = false">Batal</button>
          <button class="cw__form-submit" @click="buatTiket" :disabled="!tiketSubjek.trim() || !tiketFormPesan.trim() || loading">
            {{ loading ? 'Mengirim...' : 'Kirim ke HRD' }}
          </button>
        </div>
      </div>

      <div v-if="!showTiketForm" class="cw__input-area">
        <input
          v-model="botInput"
          class="cw__input"
          placeholder="Ketik pertanyaan..."
          @keydown.enter.prevent="kirimBot"
          :disabled="botTyping"
          maxlength="500"
        />
        <button class="cw__send" @click="kirimBot" :disabled="!botInput.trim() || botTyping">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none">
            <line x1="22" y1="2" x2="11" y2="13" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
            <polygon points="22 2 15 22 11 13 2 9 22 2" stroke="currentColor" stroke-width="2" stroke-linejoin="round"/>
          </svg>
        </button>
      </div>
    </template>

    <!-- View: Chat Tiket Aktif -->
    <template v-else>
      <div class="cw__tiket-info">
        <span class="cw__tiket-num">{{ activeTiket.nomor_tiket }}</span>
        <span class="cw__tiket-kategori">{{ activeTiket.kategori }}</span>
        <span class="cw__tiket-status" :class="'cw__tiket-status--' + activeTiket.status">{{ labelStatus(activeTiket.status) }}</span>
      </div>

      <div class="cw__messages" ref="msgContainerTiket">
        <div v-for="(msg, i) in tiketPesan" :key="i" class="cw__msg" :class="msg.sender_id === userId ? 'cw__msg--user' : 'cw__msg--bot'">
          <div class="cw__bubble">
            <div class="cw__sender" v-if="msg.sender_id !== userId">{{ msg.sender_nama }}</div>
            {{ msg.pesan }}
          </div>
          <div class="cw__time">{{ formatTime(msg.created_at) }}</div>
        </div>
      </div>

      <div class="cw__input-area" v-if="activeTiket.status === 'diproses'">
        <input
          v-model="tiketInput"
          class="cw__input"
          placeholder="Balas pesan..."
          @keydown.enter.prevent="kirimPesan"
          :disabled="loading"
          maxlength="1000"
        />
        <button class="cw__send" @click="kirimPesan" :disabled="!tiketInput.trim() || loading">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none">
            <line x1="22" y1="2" x2="11" y2="13" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
            <polygon points="22 2 15 22 11 13 2 9 22 2" stroke="currentColor" stroke-width="2" stroke-linejoin="round"/>
          </svg>
        </button>
      </div>

      <div class="cw__waiting" v-else-if="activeTiket.status === 'menunggu'">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none"><circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2"/><polyline points="12 6 12 12 16 14" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>
        Menunggu HRD mengambil tiket...
      </div>

      <div class="cw__selesai" v-else-if="activeTiket.status === 'selesai'">
        <div>✅ Tiket telah diselesaikan</div>
        <button class="cw__cta-btn" @click="resetBot" style="margin-top:8px">Tanya Lagi</button>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onBeforeUnmount, watch, nextTick } from 'vue'
import axios from 'axios'

const props = defineProps({
  wsHub: { type: Object, default: null }
})

const emit = defineEmits(['unread'])

const user = JSON.parse(localStorage.getItem('user') || '{}')
const userId = user.id || ''
const token = localStorage.getItem('access_token') || ''

// Bot state
const botMessages = ref([{
  from: 'bot',
  text: 'Halo, ' + (user.nama_lengkap?.split(' ')[0] || 'Kak') + '! 👋\n\nSaya asisten otomatis e-Magang. Tanyakan apa saja seputar program magang!',
  time: nowTime()
}])
const botInput = ref('')
const botTyping = ref(false)
const showTiketCta = ref(false)
const showTiketForm = ref(false)
const tiketSubjek = ref('')
const tiketFormPesan = ref('')
const tiketKategori = ref('umum')

// Tiket state
const activeTiket = ref(null)
const tiketPesan_ = ref([])
const tiketInput = ref('')
const loading = ref(false)

// Timer hangus
const timerLabel = ref('')
let timerInterval = null

const msgContainerBot = ref(null)
const msgContainerTiket = ref(null)

function nowTime() {
  return new Date().toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' })
}

function formatMsg(text) {
  return text.replace(/\n/g, '<br>').replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')
}

function formatTime(ts) {
  if (!ts) return ''
  return new Date(ts).toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' })
}

function labelStatus(s) {
  return { menunggu: 'Menunggu', diproses: 'Diproses', selesai: 'Selesai', hangus: 'Kedaluwarsa' }[s] || s
}

async function scrollBot() {
  await nextTick()
  if (msgContainerBot.value) msgContainerBot.value.scrollTop = msgContainerBot.value.scrollHeight
}

async function scrollTiket() {
  await nextTick()
  if (msgContainerTiket.value) msgContainerTiket.value.scrollTop = msgContainerTiket.value.scrollHeight
}

// NLP Bot
async function kirimBot() {
  const teks = botInput.value.trim()
  if (!teks || botTyping.value) return

  botMessages.value.push({ from: 'user', text: teks, time: nowTime() })
  botInput.value = ''
  showTiketCta.value = false
  showTiketForm.value = false
  await scrollBot()

  botTyping.value = true
  await scrollBot()

  try {
    const { data } = await axios.post('/api/chatbot/query', { pesan: teks })
    const result = data.data
    await new Promise(r => setTimeout(r, 600))
    botTyping.value = false

    if (result.terjawab) {
      botMessages.value.push({ from: 'bot', text: result.jawaban, time: nowTime() })
      showTiketCta.value = false
    } else {
      botMessages.value.push({
        from: 'bot',
        text: 'Maaf, saya tidak bisa menjawab pertanyaan ini secara otomatis. 🙏\n\nKlik tombol di bawah untuk menghubungi HRD secara langsung.',
        time: nowTime()
      })
      showTiketCta.value = true
      tiketSubjek.value = teks.substring(0, 80)
      tiketFormPesan.value = teks
    }
  } catch {
    botTyping.value = false
    botMessages.value.push({ from: 'bot', text: 'Maaf, terjadi gangguan. Silakan coba lagi.', time: nowTime() })
  }

  await scrollBot()
}

// Buat tiket
async function buatTiket() {
  if (!tiketSubjek.value.trim() || !tiketFormPesan.value.trim()) return
  loading.value = true
  try {
    const { data } = await axios.post('/api/chat/tiket', {
      subjek: tiketSubjek.value,
      pesan: tiketFormPesan.value,
      kategori: tiketKategori.value
    }, { headers: { Authorization: 'Bearer ' + token } })

    activeTiket.value = data.data
    showTiketForm.value = false
    showTiketCta.value = false
    await loadPesan()
    startTimer()
  } catch (e) {
    alert(e?.response?.data?.message || 'Gagal membuat tiket')
  } finally {
    loading.value = false
  }
}

async function kirimPesan() {
  if (!tiketInput.value.trim() || !activeTiket.value) return
  loading.value = true
  try {
    await axios.post(`/api/chat/tiket/${activeTiket.value.id}/pesan`,
      { pesan: tiketInput.value },
      { headers: { Authorization: 'Bearer ' + token } }
    )
    tiketInput.value = ''
    await loadPesan()
  } catch (e) {
    alert(e?.response?.data?.message || 'Gagal kirim pesan')
  } finally {
    loading.value = false
  }
}

async function loadPesan() {
  if (!activeTiket.value) return
  try {
    const { data } = await axios.get(`/api/chat/tiket/${activeTiket.value.id}/pesan`,
      { headers: { Authorization: 'Bearer ' + token } })
    tiketPesan_.value = data.data || []
    await scrollTiket()
  } catch {}
}

// Expose tiketPesan sebagai computed
const tiketPesan = computed(() => tiketPesan_.value)

async function loadActiveTiket() {
  try {
    const { data } = await axios.get('/api/chat/tiket/saya',
      { headers: { Authorization: 'Bearer ' + token } })
    const list = data.data || []
    const aktif = list.find(t => t.status === 'menunggu' || t.status === 'diproses')
    activeTiket.value = aktif || null
    if (activeTiket.value) {
      await loadPesan()
      startTimer()
    }
  } catch {}
}

function startTimer() {
  if (timerInterval) clearInterval(timerInterval)
  updateTimer()
  timerInterval = setInterval(updateTimer, 1000)
}

function updateTimer() {
  if (!activeTiket.value?.expires_at || activeTiket.value.status !== 'menunggu') {
    timerLabel.value = ''
    return
  }
  const diff = new Date(activeTiket.value.expires_at) - new Date()
  if (diff <= 0) {
    timerLabel.value = 'Kedaluwarsa'
    clearInterval(timerInterval)
    return
  }
  const h = Math.floor(diff / 3600000)
  const m = Math.floor((diff % 3600000) / 60000)
  const s = Math.floor((diff % 60000) / 1000)
  timerLabel.value = h > 0 ? `${h}j ${m}m` : m > 0 ? `${m}m ${s}d` : `${s}d`
}

function resetBot() {
  activeTiket.value = null
  tiketPesan_.value = []
  showTiketCta.value = false
  showTiketForm.value = false
}

// WebSocket events
function handleWsMessage(event) {
  try {
    const msg = JSON.parse(event.data)
    if (msg.type === 'chat_message' && activeTiket.value && msg.data.tiket_id === activeTiket.value.id) {
      loadPesan()
      if (msg.data.sender_id !== userId) emit('unread', 1)
    }
    if (msg.type === 'tiket_update' && activeTiket.value && msg.data.tiket_id === activeTiket.value.id) {
      activeTiket.value = { ...activeTiket.value, status: msg.data.status }
    }
  } catch {}
}

onMounted(() => {
  loadActiveTiket()
  if (props.wsHub) {
    props.wsHub.addEventListener('message', handleWsMessage)
  }
})

onBeforeUnmount(() => {
  if (timerInterval) clearInterval(timerInterval)
  if (props.wsHub) {
    props.wsHub.removeEventListener('message', handleWsMessage)
  }
})

watch(activeTiket, () => {
  if (!activeTiket.value) {
    if (timerInterval) clearInterval(timerInterval)
  }
})
</script>

<style scoped>
.cw {
  display: flex;
  flex-direction: column;
  height: 100%;
  min-height: 0;
  font-family: 'Poppins', sans-serif;
  background: #fff;
  overflow: hidden;
}

.cw__header {
  background: linear-gradient(135deg, #16a34a, #15803d);
  padding: 14px 16px;
  color: white;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.cw__header-info { display: flex; align-items: center; gap: 10px; }
.cw__avatar {
  width: 42px; height: 42px;
  position: relative;
  flex-shrink: 0;
}
.cw__avatar-img {
  width: 42px; height: 42px;
  border-radius: 50%;
  object-fit: cover;
  background: rgba(255,255,255,0.15);
  border: 2px solid rgba(255,255,255,0.3);
}
.cw__online-dot {
  position: absolute;
  bottom: 1px; right: 1px;
  width: 10px; height: 10px;
  background: #4ade80;
  border-radius: 50%;
  border: 2px solid #15803d;
}
.cw__name { font-weight: 600; font-size: 14px; }
.cw__status { font-size: 11px; opacity: 0.8; margin-top: 1px; }
.cw__status--active { opacity: 1; color: #bbf7d0; }

.cw__timer {
  display: flex; align-items: center; gap: 4px;
  font-size: 11px; background: rgba(255,255,255,0.15);
  padding: 3px 8px; border-radius: 10px;
}

.cw__tiket-info {
  padding: 8px 14px;
  background: #f0fdf4;
  border-bottom: 1px solid #bbf7d0;
  display: flex; align-items: center; gap: 8px;
  flex-shrink: 0;
}

.cw__tiket-num { font-size: 12px; font-weight: 600; color: #15803d; }
.cw__tiket-kategori { font-size: 11px; background: #dcfce7; color: #166534; padding: 2px 8px; border-radius: 10px; text-transform: capitalize; }
.cw__tiket-status { font-size: 11px; padding: 2px 8px; border-radius: 10px; margin-left: auto; }
.cw__tiket-status--menunggu { background: #fef9c3; color: #a16207; }
.cw__tiket-status--diproses { background: #dbeafe; color: #1d4ed8; }
.cw__tiket-status--selesai { background: #dcfce7; color: #166534; }

.cw__messages {
  flex: 1;
  overflow-y: auto;
  padding: 12px 14px;
  display: flex;
  flex-direction: column;
  gap: 8px;
  background: #f8fafb;
  min-height: 0;
}

.cw__msg { display: flex; flex-direction: column; max-width: 85%; }
.cw__msg--user { align-self: flex-end; align-items: flex-end; }
.cw__msg--bot { align-self: flex-start; align-items: flex-start; }

.cw__bubble {
  padding: 9px 13px;
  border-radius: 14px;
  font-size: 13px;
  line-height: 1.5;
  word-break: break-word;
}

.cw__msg--user .cw__bubble { background: #16a34a; color: white; border-bottom-right-radius: 4px; }
.cw__msg--bot .cw__bubble {
  background: white; color: #1e293b;
  border: 1px solid #e2e8f0; border-bottom-left-radius: 4px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.06);
}
.cw__sender { font-size: 10px; color: #64748b; margin-bottom: 3px; font-weight: 600; }
.cw__time { font-size: 10px; color: #94a3b8; margin-top: 3px; padding: 0 4px; }

.cw__bubble--typing { display: flex; gap: 4px; align-items: center; padding: 12px 16px; }
.cw__bubble--typing span {
  width: 6px; height: 6px; border-radius: 50%; background: #94a3b8;
  animation: typing 1.2s infinite;
}
.cw__bubble--typing span:nth-child(2) { animation-delay: 0.2s; }
.cw__bubble--typing span:nth-child(3) { animation-delay: 0.4s; }
@keyframes typing {
  0%, 60%, 100% { transform: translateY(0); opacity: 0.4; }
  30% { transform: translateY(-4px); opacity: 1; }
}

.cw__cta, .cw__form {
  background: #f0fdf4;
  border-top: 1px solid #bbf7d0;
  padding: 10px 14px;
  flex-shrink: 0;
}
.cw__cta p { font-size: 12px; color: #15803d; margin: 0 0 8px; }
.cw__cta-btn {
  display: inline-flex; align-items: center; gap: 6px;
  background: #16a34a; color: white; text-decoration: none;
  border-radius: 8px; padding: 7px 14px; font-size: 12px; font-weight: 600;
  transition: background 0.2s; width: 100%; justify-content: center;
  border: none; cursor: pointer;
}
.cw__cta-btn:hover { background: #15803d; }

.cw__form-title { font-size: 13px; font-weight: 600; color: #166534; margin-bottom: 8px; }
.cw__form-input, .cw__form-textarea, .cw__form-select {
  width: 100%; border: 1px solid #bbf7d0; border-radius: 8px;
  padding: 7px 10px; font-size: 12px; font-family: inherit;
  outline: none; margin-bottom: 6px; box-sizing: border-box;
  background: white;
}
.cw__form-input:focus, .cw__form-textarea:focus, .cw__form-select:focus { border-color: #16a34a; }
.cw__form-textarea { resize: none; }
.cw__form-actions { display: flex; gap: 6px; justify-content: flex-end; }
.cw__form-cancel {
  padding: 6px 12px; border: 1px solid #e2e8f0; border-radius: 6px;
  background: white; cursor: pointer; font-size: 12px;
}
.cw__form-submit {
  padding: 6px 14px; background: #16a34a; color: white; border: none;
  border-radius: 6px; cursor: pointer; font-size: 12px; font-weight: 600;
}
.cw__form-submit:disabled { opacity: 0.6; cursor: not-allowed; }

.cw__waiting {
  display: flex; align-items: center; gap: 8px;
  padding: 12px 14px; color: #94a3b8; font-size: 12px;
  border-top: 1px solid #e2e8f0; flex-shrink: 0;
  animation: pulse 2s infinite;
}
@keyframes pulse { 0%,100% { opacity: 1; } 50% { opacity: 0.5; } }

.cw__selesai {
  padding: 12px 14px; border-top: 1px solid #e2e8f0;
  text-align: center; font-size: 12px; color: #64748b; flex-shrink: 0;
}

.cw__input-area {
  display: flex; gap: 8px; padding: 10px 12px;
  border-top: 1px solid #e2e8f0; flex-shrink: 0; background: white;
}
.cw__input {
  flex: 1; border: 1px solid #e2e8f0; border-radius: 20px;
  padding: 8px 14px; font-size: 13px; outline: none;
  font-family: inherit; transition: border-color 0.2s; background: #f8fafb;
}
.cw__input:focus { border-color: #16a34a; background: white; }
.cw__input:disabled { opacity: 0.6; }
.cw__send {
  width: 36px; height: 36px; border-radius: 50%;
  background: #16a34a; color: white; border: none;
  cursor: pointer; display: flex; align-items: center; justify-content: center;
  transition: background 0.2s, transform 0.15s;
}
.cw__send:hover:not(:disabled) { background: #15803d; transform: scale(1.05); }
.cw__send:disabled { background: #94a3b8; cursor: not-allowed; }
</style>
