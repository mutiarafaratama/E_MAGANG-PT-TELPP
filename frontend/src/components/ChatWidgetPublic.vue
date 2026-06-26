<template>
  <div class="cwp">
    <!-- Header -->
    <div class="cwp__header">
      <div class="cwp__header-info">
        <div class="cwp__avatar">
          <img src="/logo_chatbot.png" alt="Bot" class="cwp__avatar-img" />
          <span class="cwp__online-dot"></span>
        </div>
        <div>
          <div class="cwp__name">Asisten e-Magang</div>
          <div class="cwp__status">
            <span class="cwp__status-dot"></span>
            Online
          </div>
        </div>
      </div>
    </div>

    <!-- Messages -->
    <div class="cwp__messages" ref="msgContainer">
      <div v-for="(msg, i) in messages" :key="i" class="cwp__msg" :class="'cwp__msg--' + msg.from">
        <div class="cwp__msg-row">
          <div v-if="msg.from === 'bot'" class="cwp__bot-avatar">
            <img src="/logo_chatbot.png" alt="Bot" />
          </div>
          <div>
            <div class="cwp__bubble" :class="'cwp__bubble--' + msg.from">
              <span v-html="formatMsg(msg.text)"></span>
            </div>
            <div class="cwp__time">{{ msg.time }}</div>
          </div>
        </div>
      </div>

      <!-- Typing indicator -->
      <div v-if="isTyping" class="cwp__msg cwp__msg--bot">
        <div class="cwp__msg-row">
          <div class="cwp__bot-avatar">
            <img src="/logo_chatbot.png" alt="Bot" />
          </div>
          <div>
            <div class="cwp__bubble cwp__bubble--bot cwp__bubble--typing">
              <span></span><span></span><span></span>
            </div>
            <div class="cwp__time">mengetik...</div>
          </div>
        </div>
      </div>
    </div>

    <!-- CTA Login jika tidak terjawab -->
    <div v-if="showLoginCta" class="cwp__cta">
      <div class="cwp__cta-icon">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none"><path d="M21 15a2 2 0 01-2 2H7l-4 4V5a2 2 0 012-2h14a2 2 0 012 2z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>
      </div>
      <p>Pertanyaan Anda memerlukan bantuan langsung dari tim HRD kami.</p>
      <router-link to="/login" class="cwp__cta-btn">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2" stroke="currentColor" stroke-width="2"/><circle cx="12" cy="7" r="4" stroke="currentColor" stroke-width="2"/></svg>
        Login untuk Hubungi HRD
      </router-link>
    </div>

    <!-- Input -->
    <div class="cwp__input-area">
      <input
        v-model="inputText"
        class="cwp__input"
        placeholder="Ketik pertanyaan..."
        @keydown.enter.prevent="kirim"
        :disabled="isTyping"
        maxlength="500"
        autocomplete="off"
      />
      <button class="cwp__send" @click="kirim" :disabled="!inputText.trim() || isTyping">
        <svg width="17" height="17" viewBox="0 0 24 24" fill="none">
          <line x1="22" y1="2" x2="11" y2="13" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
          <polygon points="22 2 15 22 11 13 2 9 22 2" fill="currentColor"/>
        </svg>
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, nextTick, onMounted } from 'vue'
import axios from 'axios'

const messages = ref([
  {
    from: 'bot',
    text: 'Selamat datang di Asisten e-Magang PT TELPP.\n\nSaya siap membantu menjawab pertanyaan seputar program magang. Silakan ketik pertanyaan Anda.',
    time: nowTime()
  }
])
const inputText = ref('')
const isTyping = ref(false)
const showLoginCta = ref(false)
const msgContainer = ref(null)

function nowTime() {
  return new Date().toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' })
}

function formatMsg(text) {
  return text
    .replace(/\n/g, '<br>')
    .replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')
}

async function scrollBottom() {
  await nextTick()
  if (msgContainer.value) {
    msgContainer.value.scrollTop = msgContainer.value.scrollHeight
  }
}

async function kirim() {
  const teks = inputText.value.trim()
  if (!teks || isTyping.value) return

  messages.value.push({ from: 'user', text: teks, time: nowTime() })
  inputText.value = ''
  showLoginCta.value = false
  await scrollBottom()

  isTyping.value = true
  await scrollBottom()

  try {
    const { data } = await axios.post('/api/chatbot/query', { pesan: teks })
    const result = data.data

    await new Promise(r => setTimeout(r, 800))
    isTyping.value = false

    if (result.terjawab) {
      messages.value.push({ from: 'bot', text: result.jawaban, time: nowTime() })
      showLoginCta.value = false
    } else {
      messages.value.push({
        from: 'bot',
        text: 'Mohon maaf, saya belum memiliki informasi untuk pertanyaan tersebut.\n\nSilakan login dan hubungi tim HRD secara langsung melalui fitur tiket chat.',
        time: nowTime()
      })
      showLoginCta.value = true
    }
  } catch {
    isTyping.value = false
    messages.value.push({
      from: 'bot',
      text: 'Maaf, terjadi gangguan koneksi. Silakan coba lagi.',
      time: nowTime()
    })
  }

  await scrollBottom()
}

onMounted(() => scrollBottom())
</script>

<style scoped>
.cwp {
  display: flex;
  flex-direction: column;
  height: 100%;
  max-height: 560px;
  font-family: 'Poppins', sans-serif;
  background: #fff;
  border-radius: 18px;
  overflow: hidden;
}

/* Header */
.cwp__header {
  background: linear-gradient(135deg, #16a34a 0%, #15803d 100%);
  padding: 14px 16px;
  color: white;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.cwp__header-info {
  display: flex;
  align-items: center;
  gap: 10px;
}

.cwp__avatar {
  width: 42px;
  height: 42px;
  position: relative;
  flex-shrink: 0;
}

.cwp__avatar-img {
  width: 42px;
  height: 42px;
  border-radius: 50%;
  object-fit: cover;
  background: rgba(255,255,255,0.15);
  border: 2px solid rgba(255,255,255,0.3);
}

.cwp__online-dot {
  position: absolute;
  bottom: 1px;
  right: 1px;
  width: 10px;
  height: 10px;
  background: #4ade80;
  border-radius: 50%;
  border: 2px solid #15803d;
}

.cwp__name {
  font-weight: 700;
  font-size: 14px;
  letter-spacing: 0.01em;
}

.cwp__status {
  font-size: 11px;
  opacity: 0.9;
  margin-top: 2px;
  display: flex;
  align-items: center;
  gap: 5px;
}

.cwp__status-dot {
  width: 6px;
  height: 6px;
  background: #4ade80;
  border-radius: 50%;
  animation: pulse-dot 2s infinite;
}

@keyframes pulse-dot {
  0%, 100% { opacity: 1; transform: scale(1); }
  50% { opacity: 0.6; transform: scale(0.85); }
}

.cwp__header-badge {
  font-size: 10px;
  font-weight: 700;
  background: rgba(255,255,255,0.2);
  border: 1px solid rgba(255,255,255,0.3);
  padding: 3px 8px;
  border-radius: 10px;
  letter-spacing: 0.05em;
}

/* Messages */
.cwp__messages {
  flex: 1;
  overflow-y: auto;
  padding: 14px 14px 8px;
  display: flex;
  flex-direction: column;
  gap: 10px;
  background: #f0fdf4;
  min-height: 0;
  scroll-behavior: smooth;
}

.cwp__messages::-webkit-scrollbar { width: 4px; }
.cwp__messages::-webkit-scrollbar-track { background: transparent; }
.cwp__messages::-webkit-scrollbar-thumb { background: #bbf7d0; border-radius: 2px; }

.cwp__msg { display: flex; flex-direction: column; }
.cwp__msg--user { align-items: flex-end; }
.cwp__msg--bot { align-items: flex-start; }

.cwp__msg-row {
  display: flex;
  align-items: flex-end;
  gap: 7px;
}

.cwp__msg--user .cwp__msg-row { flex-direction: row-reverse; }

.cwp__bot-avatar {
  width: 28px;
  height: 28px;
  flex-shrink: 0;
  margin-bottom: 2px;
}

.cwp__bot-avatar img {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  object-fit: cover;
  border: 1px solid #bbf7d0;
}

.cwp__bubble {
  padding: 9px 13px;
  border-radius: 16px;
  font-size: 13px;
  line-height: 1.55;
  word-break: break-word;
  max-width: 240px;
}

.cwp__bubble--user {
  background: linear-gradient(135deg, #16a34a, #15803d);
  color: white;
  border-bottom-right-radius: 4px;
  box-shadow: 0 2px 8px rgba(22,163,74,0.25);
}

.cwp__bubble--bot {
  background: white;
  color: #1e293b;
  border: 1px solid #dcfce7;
  border-bottom-left-radius: 4px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.06);
}

.cwp__time {
  font-size: 10px;
  color: #94a3b8;
  margin-top: 4px;
  padding: 0 4px;
}

/* Typing indicator */
.cwp__bubble--typing {
  display: flex;
  gap: 5px;
  align-items: center;
  padding: 12px 16px;
  min-width: 56px;
}

.cwp__bubble--typing span {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: #86efac;
  animation: typing-bounce 1.3s infinite;
}

.cwp__bubble--typing span:nth-child(2) { animation-delay: 0.2s; }
.cwp__bubble--typing span:nth-child(3) { animation-delay: 0.4s; }

@keyframes typing-bounce {
  0%, 60%, 100% { transform: translateY(0); opacity: 0.4; }
  30% { transform: translateY(-5px); opacity: 1; }
}

/* CTA Login */
.cwp__cta {
  background: #f0fdf4;
  border-top: 1px solid #bbf7d0;
  padding: 12px 14px;
  flex-shrink: 0;
  text-align: center;
}

.cwp__cta-icon {
  font-size: 18px;
  margin-bottom: 4px;
}

.cwp__cta p {
  font-size: 12px;
  color: #166534;
  margin: 0 0 10px;
  line-height: 1.5;
}

.cwp__cta-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  background: linear-gradient(135deg, #16a34a, #15803d);
  color: white;
  text-decoration: none;
  border-radius: 10px;
  padding: 8px 16px;
  font-size: 12px;
  font-weight: 600;
  transition: opacity 0.2s, transform 0.15s;
  width: 100%;
  justify-content: center;
  box-shadow: 0 2px 8px rgba(22,163,74,0.3);
}

.cwp__cta-btn:hover { opacity: 0.92; transform: translateY(-1px); }

/* Input */
.cwp__input-area {
  display: flex;
  gap: 8px;
  padding: 10px 12px;
  border-top: 1px solid #dcfce7;
  flex-shrink: 0;
  background: white;
  align-items: center;
}

.cwp__input {
  flex: 1;
  border: 1.5px solid #dcfce7;
  border-radius: 22px;
  padding: 9px 15px;
  font-size: 13px;
  outline: none;
  font-family: inherit;
  transition: border-color 0.2s, box-shadow 0.2s;
  background: #f8fafb;
}

.cwp__input:focus {
  border-color: #16a34a;
  background: white;
  box-shadow: 0 0 0 3px rgba(22,163,74,0.08);
}

.cwp__input:disabled { opacity: 0.6; }

.cwp__send {
  width: 38px;
  height: 38px;
  border-radius: 50%;
  background: linear-gradient(135deg, #16a34a, #15803d);
  color: white;
  border: none;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  transition: opacity 0.2s, transform 0.15s;
  box-shadow: 0 2px 8px rgba(22,163,74,0.3);
}

.cwp__send:hover:not(:disabled) { opacity: 0.9; transform: scale(1.05); }
.cwp__send:disabled { background: #cbd5e1; box-shadow: none; cursor: not-allowed; }
</style>
