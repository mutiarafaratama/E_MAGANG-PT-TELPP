<template>
  <!-- Bubble trigger -->
  <div class="chatbot-wrap">
    <button class="chatbot-bubble" @click="toggleChat" :class="{ 'chatbot-bubble--open': isOpen }" aria-label="Buka chat bantuan">
      <transition name="icon-swap">
        <svg v-if="!isOpen" key="chat" width="24" height="24" viewBox="0 0 24 24" fill="none">
          <path d="M21 15a2 2 0 01-2 2H7l-4 4V5a2 2 0 012-2h14a2 2 0 012 2z" stroke="#fff" stroke-width="2" stroke-linejoin="round"/>
        </svg>
        <svg v-else key="close" width="22" height="22" viewBox="0 0 24 24" fill="none">
          <path d="M18 6L6 18M6 6l12 12" stroke="#fff" stroke-width="2.5" stroke-linecap="round"/>
        </svg>
      </transition>
      <span v-if="unreadCount > 0 && !isOpen" class="chatbot-bubble__badge">{{ unreadCount }}</span>
    </button>

    <!-- Chat window -->
    <transition name="chat-window">
      <div v-if="isOpen" class="chatbot-window">
        <!-- Header -->
        <div class="chatbot-header">
          <div class="chatbot-header__avatar">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none">
              <path d="M21 15a2 2 0 01-2 2H7l-4 4V5a2 2 0 012-2h14a2 2 0 012 2z" stroke="#fff" stroke-width="2"/>
            </svg>
          </div>
          <div>
            <div class="chatbot-header__name">Asisten e-Magang</div>
            <div class="chatbot-header__status">
              <span class="status-dot"></span> Online — Siap membantu
            </div>
          </div>
        </div>

        <!-- Messages -->
        <div class="chatbot-messages" ref="messagesEl">
          <div v-for="(msg, i) in messages" :key="i" class="chatbot-msg" :class="`chatbot-msg--${msg.from}`">
            <div class="chatbot-msg__bubble">
              <span v-html="msg.text"></span>
            </div>
            <div class="chatbot-msg__time">{{ msg.time }}</div>
          </div>

          <!-- Quick replies -->
          <div v-if="showQuickReplies && !isTyping" class="chatbot-quickreplies">
            <button
              v-for="qr in quickReplies"
              :key="qr"
              class="chatbot-qr"
              @click="sendQuickReply(qr)"
            >{{ qr }}</button>
          </div>

          <!-- Typing indicator -->
          <div v-if="isTyping" class="chatbot-msg chatbot-msg--bot">
            <div class="chatbot-msg__bubble chatbot-typing">
              <span></span><span></span><span></span>
            </div>
          </div>
        </div>

        <!-- CTA ke login setelah tidak terjawab -->
        <div v-if="showLoginCta" class="chatbot-cta">
          <p>Pertanyaan Anda perlu dijawab langsung oleh tim HRD kami.</p>
          <div class="chatbot-cta__btns">
            <a href="/login" class="chatbot-cta__btn chatbot-cta__btn--primary">Masuk & Chat HRD</a>
            <a href="/daftar" class="chatbot-cta__btn">Daftar Dulu</a>
          </div>
        </div>

        <!-- Input -->
        <div class="chatbot-input-wrap">
          <input
            v-model="inputText"
            type="text"
            class="chatbot-input"
            placeholder="Ketik pertanyaan Anda..."
            @keyup.enter="sendMessage"
            :disabled="isTyping"
            maxlength="300"
          />
          <button class="chatbot-send" @click="sendMessage" :disabled="!inputText.trim() || isTyping">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none">
              <path d="M22 2L11 13M22 2l-7 20-4-9-9-4 20-7z" stroke="currentColor" stroke-width="2" stroke-linejoin="round"/>
            </svg>
          </button>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup lang="ts">
import { ref, nextTick, onMounted } from "vue";

interface KBItem {
  id: string;
  pertanyaan: string;
  jawaban: string;
}

interface ChatMsg {
  from: "bot" | "user";
  text: string;
  time: string;
}

const isOpen = ref(false);
const inputText = ref("");
const messages = ref<ChatMsg[]>([]);
const isTyping = ref(false);
const showQuickReplies = ref(true);
const showLoginCta = ref(false);
const unreadCount = ref(1);
const messagesEl = ref<HTMLElement | null>(null);
const kb = ref<KBItem[]>([]);

const quickReplies = [
  "Cara mendaftar magang",
  "Dokumen yang dibutuhkan",
  "Lama proses verifikasi",
  "Cara cek status pengajuan",
  "Siapa yang bisa mendaftar",
];

// ── NLP matching ─────────────────────────────────────────────────────────────
const STOP_WORDS = new Set([
  "yang","dan","ke","di","ini","itu","saya","kami","kamu","ada","tidak","bisa",
  "untuk","dari","atau","dengan","jika","apa","bagaimana","kapan","berapa",
  "dimana","siapa","apakah","sudah","belum","masih","lebih","akan","juga",
  "harus","cara","gimana","gmn","mohon","tolong","ingin","mau","mau","boleh",
  "yah","ya","dong","sih","nih","lah","deh","kah","nya","nya","nya",
]);

function tokenize(text: string): string[] {
  return text
    .toLowerCase()
    .replace(/[^a-z0-9\s]/g, " ")
    .split(/\s+/)
    .filter((w) => w.length > 2 && !STOP_WORDS.has(w));
}

function matchKB(query: string): KBItem | null {
  const qt = tokenize(query);
  if (qt.length === 0 || kb.value.length === 0) return null;

  let best: KBItem | null = null;
  let bestScore = 0;

  for (const item of kb.value) {
    const qt2 = tokenize(item.pertanyaan);
    const qa2 = tokenize(item.jawaban);
    const allTokens = [...qt2, ...qa2];

    const matches = qt.filter((t) =>
      allTokens.some((e) => e.includes(t) || t.includes(e))
    ).length;
    const baseScore = matches / qt.length;

    const qMatches = qt.filter((t) =>
      qt2.some((e) => e.includes(t) || t.includes(e))
    ).length;
    const score = baseScore + (qMatches / Math.max(qt.length, 1)) * 0.5;

    if (score > bestScore) {
      bestScore = score;
      best = item;
    }
  }

  return bestScore >= 0.28 ? best : null;
}

// ── Utils ─────────────────────────────────────────────────────────────────────
function nowTime(): string {
  return new Date().toLocaleTimeString("id-ID", { hour: "2-digit", minute: "2-digit" });
}

async function scrollBottom() {
  await nextTick();
  if (messagesEl.value) {
    messagesEl.value.scrollTop = messagesEl.value.scrollHeight;
  }
}

function addBotMsg(text: string) {
  messages.value.push({ from: "bot", text, time: nowTime() });
  scrollBottom();
}

function addUserMsg(text: string) {
  messages.value.push({ from: "user", text, time: nowTime() });
  scrollBottom();
}

async function typingThenReply(text: string, delay = 900) {
  isTyping.value = true;
  scrollBottom();
  await new Promise((r) => setTimeout(r, delay));
  isTyping.value = false;
  addBotMsg(text);
}

// ── Actions ───────────────────────────────────────────────────────────────────
function toggleChat() {
  isOpen.value = !isOpen.value;
  if (isOpen.value) {
    unreadCount.value = 0;
    scrollBottom();
  }
}

async function sendQuickReply(text: string) {
  showQuickReplies.value = false;
  addUserMsg(text);
  await replyFor(text);
}

async function sendMessage() {
  const text = inputText.value.trim();
  if (!text) return;
  inputText.value = "";
  showQuickReplies.value = false;
  showLoginCta.value = false;
  addUserMsg(text);
  await replyFor(text);
}

async function replyFor(query: string) {
  const match = matchKB(query);
  if (match) {
    await typingThenReply(match.jawaban);
    await new Promise((r) => setTimeout(r, 400));
    addBotMsg("Ada pertanyaan lain yang bisa saya bantu? 😊");
    showQuickReplies.value = false;
  } else {
    await typingThenReply(
      "Maaf, saya belum punya jawaban untuk pertanyaan tersebut. " +
      "Pertanyaan Anda akan lebih baik dijawab langsung oleh tim HRD kami.",
      1200
    );
    showLoginCta.value = true;
  }
}

// ── Init ──────────────────────────────────────────────────────────────────────
async function fetchKB() {
  try {
    const res = await fetch("/api/chat/knowledge");
    const data = await res.json();
    kb.value = data.data ?? [];
  } catch {
    kb.value = [];
  }
}

onMounted(async () => {
  await fetchKB();
  addBotMsg(
    "Halo! 👋 Selamat datang di e-Magang PT TELPP.<br>" +
    "Saya asisten virtual yang siap menjawab pertanyaan umum seputar pendaftaran magang.<br><br>" +
    "Pilih pertanyaan di bawah atau ketik langsung:"
  );
});
</script>

<style scoped>
.chatbot-wrap {
  position: fixed;
  bottom: 28px;
  right: 28px;
  z-index: 9999;
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 12px;
}

/* Bubble */
.chatbot-bubble {
  width: 56px; height: 56px;
  border-radius: 50%;
  background: linear-gradient(135deg, #48AF4A, #2d8a30);
  border: none; cursor: pointer;
  display: flex; align-items: center; justify-content: center;
  box-shadow: 0 4px 20px rgba(72,175,74,0.45);
  transition: transform 0.2s, box-shadow 0.2s;
  position: relative;
}
.chatbot-bubble:hover { transform: scale(1.08); box-shadow: 0 6px 24px rgba(72,175,74,0.55); }
.chatbot-bubble--open { background: linear-gradient(135deg, #374151, #1f2937); box-shadow: 0 4px 20px rgba(0,0,0,0.25); }
.chatbot-bubble__badge {
  position: absolute; top: -4px; right: -4px;
  background: #ef4444; color: #fff;
  font-size: 10px; font-weight: 700;
  width: 18px; height: 18px; border-radius: 50%;
  display: flex; align-items: center; justify-content: center;
  border: 2px solid #fff;
}

/* Window */
.chatbot-window {
  width: 340px;
  border-radius: 18px;
  background: #fff;
  box-shadow: 0 8px 40px rgba(0,0,0,0.18);
  display: flex; flex-direction: column;
  overflow: hidden;
  max-height: 520px;
}

/* Header */
.chatbot-header {
  background: linear-gradient(135deg, #48AF4A, #2d8a30);
  padding: 16px 18px;
  display: flex; align-items: center; gap: 12px;
}
.chatbot-header__avatar {
  width: 36px; height: 36px; border-radius: 50%;
  background: rgba(255,255,255,0.2);
  display: flex; align-items: center; justify-content: center;
  flex-shrink: 0;
}
.chatbot-header__name { font-size: 14px; font-weight: 700; color: #fff; }
.chatbot-header__status { font-size: 11px; color: rgba(255,255,255,0.8); display: flex; align-items: center; gap: 5px; margin-top: 2px; }
.status-dot { width: 7px; height: 7px; border-radius: 50%; background: #86efac; display: inline-block; }

/* Messages */
.chatbot-messages {
  flex: 1; overflow-y: auto; padding: 14px 14px 10px;
  display: flex; flex-direction: column; gap: 8px;
  scroll-behavior: smooth;
}
.chatbot-messages::-webkit-scrollbar { width: 4px; }
.chatbot-messages::-webkit-scrollbar-thumb { background: #e5e7eb; border-radius: 4px; }

.chatbot-msg { display: flex; flex-direction: column; max-width: 85%; }
.chatbot-msg--bot { align-self: flex-start; }
.chatbot-msg--user { align-self: flex-end; }

.chatbot-msg__bubble {
  padding: 9px 13px; border-radius: 14px;
  font-size: 13px; line-height: 1.55;
}
.chatbot-msg--bot .chatbot-msg__bubble {
  background: #f1f5f9; color: #1e293b;
  border-bottom-left-radius: 4px;
}
.chatbot-msg--user .chatbot-msg__bubble {
  background: linear-gradient(135deg, #48AF4A, #2d8a30);
  color: #fff;
  border-bottom-right-radius: 4px;
}
.chatbot-msg__time { font-size: 10px; color: #9ca3af; margin-top: 3px; }
.chatbot-msg--user .chatbot-msg__time { text-align: right; }

/* Typing */
.chatbot-typing { display: flex; align-items: center; gap: 5px; padding: 12px 16px; }
.chatbot-typing span {
  width: 7px; height: 7px; border-radius: 50%; background: #9ca3af;
  animation: blink 1.2s infinite;
}
.chatbot-typing span:nth-child(2) { animation-delay: 0.2s; }
.chatbot-typing span:nth-child(3) { animation-delay: 0.4s; }
@keyframes blink { 0%,80%,100% { opacity: 0.2; transform: scale(0.85); } 40% { opacity: 1; transform: scale(1.1); } }

/* Quick replies */
.chatbot-quickreplies {
  display: flex; flex-wrap: wrap; gap: 6px; padding-top: 4px;
}
.chatbot-qr {
  background: #f0fdf4; border: 1px solid #bbf7d0; color: #16a34a;
  border-radius: 100px; padding: 5px 12px; font-size: 12px; font-weight: 500;
  cursor: pointer; transition: all 0.15s; white-space: nowrap;
  font-family: inherit;
}
.chatbot-qr:hover { background: #dcfce7; border-color: #86efac; }

/* CTA */
.chatbot-cta {
  background: #f8fafc; border-top: 1px solid #e2e8f0;
  padding: 12px 14px;
}
.chatbot-cta p { font-size: 12px; color: #64748b; margin-bottom: 10px; line-height: 1.5; }
.chatbot-cta__btns { display: flex; gap: 8px; }
.chatbot-cta__btn {
  flex: 1; text-align: center; padding: 8px 10px; border-radius: 8px;
  font-size: 12px; font-weight: 600; cursor: pointer; text-decoration: none;
  border: 1.5px solid #d1d5db; color: #374151; background: #fff; transition: all 0.15s;
}
.chatbot-cta__btn:hover { background: #f1f5f9; }
.chatbot-cta__btn--primary {
  background: linear-gradient(135deg, #48AF4A, #2d8a30);
  color: #fff; border-color: transparent;
}
.chatbot-cta__btn--primary:hover { opacity: 0.9; }

/* Input */
.chatbot-input-wrap {
  display: flex; align-items: center; gap: 8px;
  padding: 10px 12px; border-top: 1px solid #f1f5f9;
  background: #fff;
}
.chatbot-input {
  flex: 1; border: 1.5px solid #e2e8f0; border-radius: 10px;
  padding: 8px 12px; font-size: 13px; outline: none;
  font-family: inherit; color: #1e293b; transition: border-color 0.15s;
}
.chatbot-input:focus { border-color: #48AF4A; }
.chatbot-input::placeholder { color: #9ca3af; }
.chatbot-send {
  width: 34px; height: 34px; border-radius: 50%;
  background: linear-gradient(135deg, #48AF4A, #2d8a30);
  color: #fff; border: none; cursor: pointer;
  display: flex; align-items: center; justify-content: center;
  transition: opacity 0.15s; flex-shrink: 0;
}
.chatbot-send:disabled { opacity: 0.45; cursor: not-allowed; }

/* Animations */
.chat-window-enter-active { animation: chatSlideIn 0.25s ease-out; }
.chat-window-leave-active { animation: chatSlideIn 0.2s ease-in reverse; }
@keyframes chatSlideIn {
  from { opacity: 0; transform: translateY(16px) scale(0.97); }
  to { opacity: 1; transform: translateY(0) scale(1); }
}
.icon-swap-enter-active, .icon-swap-leave-active { transition: opacity 0.15s; }
.icon-swap-enter-from, .icon-swap-leave-to { opacity: 0; }
</style>
