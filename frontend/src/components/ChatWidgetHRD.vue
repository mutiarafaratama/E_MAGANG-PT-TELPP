<template>
  <div class="hrd-cw">
    <!-- INBOX VIEW -->
    <div v-if="!selectedTiket" class="hrd-cw__inbox">
      <div class="hrd-cw__head">
        <div class="hrd-cw__head-title">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none"><path d="M21 15a2 2 0 01-2 2H7l-4 4V5a2 2 0 012-2h14a2 2 0 012 2z" stroke="currentColor" stroke-width="2"/></svg>
          Helpdesk HRD
        </div>
        <div class="hrd-cw__head-badge" v-if="totalMenunggu > 0">{{ totalMenunggu }}</div>
      </div>

      <!-- Filter -->
      <div class="hrd-cw__filter">
        <button
          v-for="f in filterOpts" :key="f.val"
          class="hrd-cw__filter-btn"
          :class="{ 'hrd-cw__filter-btn--active': filterStatus === f.val }"
          @click="filterStatus = f.val; fetchTiket()"
        >{{ f.label }}</button>
      </div>

      <!-- List -->
      <div class="hrd-cw__list" ref="listRef">
        <div v-if="loading" class="hrd-cw__empty"><div class="spinner-sm"></div></div>
        <div v-else-if="tiketList.length === 0" class="hrd-cw__empty">
          <svg width="28" height="28" viewBox="0 0 24 24" fill="none"><path d="M21 15a2 2 0 01-2 2H7l-4 4V5a2 2 0 012-2h14a2 2 0 012 2z" stroke="#d1d5db" stroke-width="1.5"/></svg>
          <p>Tidak ada tiket</p>
        </div>
        <div
          v-else
          v-for="t in tiketList"
          :key="t.id"
          class="hrd-cw__card"
          :class="{ 'hrd-cw__card--menunggu': t.status === 'menunggu' }"
          @click="bukaChat(t)"
        >
          <div class="hrd-cw__card-top">
            <span class="hrd-cw__nomor">{{ t.nomor_tiket }}</span>
            <span :class="['hrd-cw__status', `hrd-cw__status--${t.status}`]">{{ labelStatus(t.status) }}</span>
          </div>
          <div class="hrd-cw__subjek">{{ t.subjek || '(tanpa subjek)' }}</div>
          <div class="hrd-cw__meta">
            <span class="hrd-cw__user">{{ t.user_nama }}</span>
            <span v-if="t.last_message" class="hrd-cw__preview">· {{ t.last_message.slice(0, 40) }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- CHAT VIEW -->
    <div v-else class="hrd-cw__chat">
      <div class="hrd-cw__chat-head">
        <button class="hrd-cw__back" @click="selectedTiket = null">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none"><polyline points="15 18 9 12 15 6" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>
        </button>
        <div class="hrd-cw__chat-info">
          <div class="hrd-cw__chat-subjek">{{ selectedTiket.subjek }}</div>
          <div class="hrd-cw__chat-user">{{ selectedTiket.user_nama }}</div>
        </div>
        <span :class="['hrd-cw__status', `hrd-cw__status--${selectedTiket.status}`]">{{ labelStatus(selectedTiket.status) }}</span>
      </div>

      <!-- Actions -->
      <div class="hrd-cw__actions" v-if="selectedTiket.status !== 'selesai' && selectedTiket.status !== 'hangus'">
        <button
          v-if="selectedTiket.status === 'menunggu'"
          class="hrd-cw__btn hrd-cw__btn--claim"
          @click="ambilTiket"
          :disabled="actionLoading"
        >{{ actionLoading ? '...' : 'Ambil Tiket' }}</button>
        <button
          v-if="selectedTiket.status === 'diproses'"
          class="hrd-cw__btn hrd-cw__btn--done"
          @click="selesaikanTiket"
          :disabled="actionLoading"
        >{{ actionLoading ? '...' : 'Selesaikan' }}</button>
      </div>

      <!-- Messages -->
      <div class="hrd-cw__messages" ref="msgRef">
        <div v-if="msgLoading" class="hrd-cw__empty"><div class="spinner-sm"></div></div>
        <div v-for="(m, i) in pesanList" :key="i" class="hrd-cw__msg" :class="m.sender_role === 'peserta' ? 'hrd-cw__msg--peserta' : 'hrd-cw__msg--hrd'">
          <div class="hrd-cw__bubble">{{ m.pesan }}</div>
          <div class="hrd-cw__msg-meta">{{ m.sender_nama }} · {{ formatTime(m.created_at) }}</div>
        </div>
      </div>

      <!-- Reply -->
      <div class="hrd-cw__reply" v-if="selectedTiket.status === 'diproses'">
        <textarea
          v-model="replyText"
          class="hrd-cw__reply-input"
          placeholder="Tulis balasan... (Enter untuk kirim)"
          rows="2"
          @keydown.enter.exact.prevent="kirimBalasan"
          maxlength="1000"
        ></textarea>
        <button
          class="hrd-cw__reply-send"
          @click="kirimBalasan"
          :disabled="!replyText.trim() || sendLoading"
        >
          <svg v-if="!sendLoading" width="16" height="16" viewBox="0 0 24 24" fill="none"><line x1="22" y1="2" x2="11" y2="13" stroke="currentColor" stroke-width="2"/><polygon points="22 2 15 22 11 13 2 9 22 2" stroke="currentColor" stroke-width="2"/></svg>
          <div v-else class="spinner-sm"></div>
        </button>
      </div>
      <div class="hrd-cw__closed-note" v-else-if="selectedTiket.status === 'selesai' || selectedTiket.status === 'hangus'">
        Tiket ini sudah ditutup.
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, nextTick } from 'vue';
import api from '@/lib/api';
import { useAppWS } from '@/composables/useAppWS';

const filterOpts = [
  { val: 'menunggu', label: 'Menunggu' },
  { val: 'diproses', label: 'Diproses' },
  { val: '',         label: 'Semua' },
];

const tiketList     = ref<any[]>([]);
const loading       = ref(false);
const filterStatus  = ref('menunggu');
const totalMenunggu = ref(0);
const selectedTiket = ref<any>(null);
const pesanList     = ref<any[]>([]);
const msgLoading    = ref(false);
const replyText     = ref('');
const sendLoading   = ref(false);
const actionLoading = ref(false);
const listRef       = ref<HTMLElement|null>(null);
const msgRef        = ref<HTMLElement|null>(null);

function labelStatus(s: string) {
  return { menunggu: 'Menunggu', diproses: 'Diproses', selesai: 'Selesai', hangus: 'Hangus' }[s] ?? s;
}

function formatTime(iso: string) {
  if (!iso) return '–';
  const d = new Date(iso);
  if (isNaN(d.getTime())) return '–';
  return d.toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' });
}

async function fetchTiket() {
  loading.value = true;
  try {
    const params: any = { limit: 30 };
    if (filterStatus.value) params.status = filterStatus.value;
    const res = await api.get('/api/chat/tiket', { params });
    tiketList.value = Array.isArray(res.data?.data) ? res.data.data : [];
    // Hitung menunggu
    const res2 = await api.get('/api/chat/tiket', { params: { status: 'menunggu', limit: 1 } });
    totalMenunggu.value = res2.data?.total ?? 0;
  } catch { /* noop */ }
  finally { loading.value = false; }
}

async function bukaChat(t: any) {
  selectedTiket.value = t;
  msgLoading.value = true;
  replyText.value = '';
  try {
    const res = await api.get(`/api/chat/tiket/${t.id}/pesan`);
    pesanList.value = Array.isArray(res.data?.data) ? res.data.data : [];
  } catch { /* noop */ }
  finally {
    msgLoading.value = false;
    nextTick(() => { if (msgRef.value) msgRef.value.scrollTop = msgRef.value.scrollHeight; });
  }
}

async function ambilTiket() {
  if (!selectedTiket.value) return;
  actionLoading.value = true;
  try {
    const res = await api.patch(`/api/chat/tiket/${selectedTiket.value.id}/ambil`);
    selectedTiket.value = res.data?.data ?? { ...selectedTiket.value, status: 'diproses' };
    fetchTiket();
  } catch { /* noop */ }
  finally { actionLoading.value = false; }
}

async function selesaikanTiket() {
  if (!selectedTiket.value) return;
  actionLoading.value = true;
  try {
    await api.patch(`/api/chat/tiket/${selectedTiket.value.id}/status`, { status: 'selesai' });
    selectedTiket.value = { ...selectedTiket.value, status: 'selesai' };
    fetchTiket();
  } catch { /* noop */ }
  finally { actionLoading.value = false; }
}

async function kirimBalasan() {
  const pesan = replyText.value.trim();
  if (!pesan || !selectedTiket.value) return;
  sendLoading.value = true;
  try {
    const res = await api.post(`/api/chat/tiket/${selectedTiket.value.id}/pesan`, { pesan });
    if (res.data?.data) pesanList.value.push(res.data.data);
    replyText.value = '';
    nextTick(() => { if (msgRef.value) msgRef.value.scrollTop = msgRef.value.scrollHeight; });
  } catch { /* noop */ }
  finally { sendLoading.value = false; }
}

const { connect, disconnect, subscribe } = useAppWS();
let wsUnsub: (() => void) | null = null;

function handleWS(msg: any) {
  if (msg.type === 'chat_message') {
    const d = msg.payload;
    if (selectedTiket.value && d?.tiket_id === selectedTiket.value.id) {
      pesanList.value.push(d);
      nextTick(() => { if (msgRef.value) msgRef.value.scrollTop = msgRef.value.scrollHeight; });
    }
    fetchTiket();
  }
  if (msg.type === 'tiket_update') {
    if (selectedTiket.value && msg.payload?.tiket_id === selectedTiket.value.id) {
      selectedTiket.value = { ...selectedTiket.value, ...msg.payload };
    }
    fetchTiket();
  }
  if (msg.type === 'badge_update') {
    fetchTiket();
  }
}

onMounted(() => { fetchTiket(); connect(); wsUnsub = subscribe(handleWS); });
onUnmounted(() => { wsUnsub?.(); disconnect(); });
</script>

<style scoped>
.hrd-cw { display: flex; flex-direction: column; height: 100%; min-height: 0; overflow: hidden; font-family: "Poppins", sans-serif; }

/* ── Inbox ── */
.hrd-cw__inbox { display: flex; flex-direction: column; height: 100%; min-height: 0; overflow: hidden; }
.hrd-cw__head { display: flex; align-items: center; justify-content: space-between; padding: 14px 16px 10px; border-bottom: 1px solid #f0f0f0; }
.hrd-cw__head-title { display: flex; align-items: center; gap: 6px; font-size: 13px; font-weight: 700; color: #111827; }
.hrd-cw__head-badge { background: #ef4444; color: #fff; font-size: 11px; font-weight: 700; border-radius: 50%; width: 20px; height: 20px; display: flex; align-items: center; justify-content: center; }

.hrd-cw__filter { display: flex; gap: 4px; padding: 8px 12px; border-bottom: 1px solid #f0f0f0; }
.hrd-cw__filter-btn { border: 1px solid #e5e7eb; border-radius: 20px; padding: 3px 10px; font-size: 11px; font-family: "Poppins", sans-serif; cursor: pointer; background: #fff; color: #6b7280; transition: all .15s; }
.hrd-cw__filter-btn--active { background: #48AF4A; color: #fff; border-color: #48AF4A; }

.hrd-cw__list { flex: 1; overflow-y: auto; padding: 8px; display: flex; flex-direction: column; gap: 6px; }

.hrd-cw__card { background: #f9fafb; border: 1px solid #e9f5e9; border-radius: 10px; padding: 10px 12px; cursor: pointer; transition: all .15s; }
.hrd-cw__card:hover { background: #f0fdf4; border-color: #48AF4A; }
.hrd-cw__card--menunggu { border-left: 3px solid #f59e0b; }
.hrd-cw__card-top { display: flex; align-items: center; justify-content: space-between; margin-bottom: 3px; }
.hrd-cw__nomor { font-size: 11px; font-weight: 600; color: #6b7280; font-family: monospace; }
.hrd-cw__subjek { font-size: 12.5px; font-weight: 600; color: #111827; margin: 2px 0; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.hrd-cw__meta { font-size: 11px; color: #9ca3af; display: flex; gap: 4px; }
.hrd-cw__user { font-weight: 500; color: #6b7280; }
.hrd-cw__preview { white-space: nowrap; overflow: hidden; text-overflow: ellipsis; max-width: 120px; }

/* ── Status ── */
.hrd-cw__status { font-size: 10px; font-weight: 600; padding: 2px 7px; border-radius: 20px; }
.hrd-cw__status--menunggu { background: #fef3c7; color: #d97706; }
.hrd-cw__status--diproses { background: #dbeafe; color: #2563eb; }
.hrd-cw__status--selesai  { background: #d1fae5; color: #059669; }
.hrd-cw__status--hangus   { background: #f3f4f6; color: #9ca3af; }

/* ── Chat ── */
.hrd-cw__chat { display: flex; flex-direction: column; height: 100%; min-height: 0; overflow: hidden; }
.hrd-cw__chat-head { display: flex; align-items: center; gap: 10px; padding: 10px 12px; border-bottom: 1px solid #f0f0f0; }
.hrd-cw__back { border: none; background: none; cursor: pointer; display: flex; align-items: center; justify-content: center; padding: 4px; border-radius: 6px; color: #6b7280; }
.hrd-cw__back:hover { background: #f3f4f6; }
.hrd-cw__chat-info { flex: 1; min-width: 0; }
.hrd-cw__chat-subjek { font-size: 12.5px; font-weight: 600; color: #111827; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.hrd-cw__chat-user { font-size: 11px; color: #9ca3af; }

.hrd-cw__actions { display: flex; gap: 6px; padding: 8px 12px; border-bottom: 1px solid #f0f0f0; }
.hrd-cw__btn { border: none; border-radius: 6px; padding: 5px 12px; font-size: 11.5px; font-weight: 600; font-family: "Poppins", sans-serif; cursor: pointer; transition: opacity .15s; }
.hrd-cw__btn:disabled { opacity: .6; cursor: not-allowed; }
.hrd-cw__btn--claim { background: #3b82f6; color: #fff; }
.hrd-cw__btn--done  { background: #48AF4A; color: #fff; }

.hrd-cw__messages { flex: 1; overflow-y: auto; padding: 10px 12px; display: flex; flex-direction: column; gap: 8px; }
.hrd-cw__msg { display: flex; flex-direction: column; }
.hrd-cw__msg--peserta { align-items: flex-start; }
.hrd-cw__msg--hrd     { align-items: flex-end; }
.hrd-cw__bubble { max-width: 85%; padding: 7px 10px; border-radius: 10px; font-size: 12.5px; line-height: 1.5; word-break: break-word; }
.hrd-cw__msg--peserta .hrd-cw__bubble { background: #f3f4f6; color: #111827; border-bottom-left-radius: 3px; }
.hrd-cw__msg--hrd     .hrd-cw__bubble { background: #48AF4A; color: #fff; border-bottom-right-radius: 3px; }
.hrd-cw__msg-meta { font-size: 10px; color: #9ca3af; margin-top: 2px; }

.hrd-cw__reply { display: flex; gap: 6px; padding: 8px 10px; border-top: 1px solid #f0f0f0; align-items: flex-end; }
.hrd-cw__reply-input { flex: 1; border: 1px solid #e5e7eb; border-radius: 8px; padding: 7px 10px; font-size: 12px; font-family: "Poppins", sans-serif; resize: none; outline: none; }
.hrd-cw__reply-input:focus { border-color: #48AF4A; }
.hrd-cw__reply-send { background: #48AF4A; border: none; border-radius: 8px; padding: 8px 10px; cursor: pointer; display: flex; align-items: center; justify-content: center; color: #fff; transition: opacity .15s; }
.hrd-cw__reply-send:disabled { opacity: .5; cursor: not-allowed; }
.hrd-cw__closed-note { text-align: center; font-size: 11.5px; color: #9ca3af; padding: 10px; border-top: 1px solid #f0f0f0; }

.hrd-cw__empty { flex: 1; display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 8px; color: #9ca3af; font-size: 12.5px; padding: 16px; text-align: center; }
.spinner-sm { width: 18px; height: 18px; border: 2px solid #e5e7eb; border-top-color: #48AF4A; border-radius: 50%; animation: spin .7s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }
</style>
