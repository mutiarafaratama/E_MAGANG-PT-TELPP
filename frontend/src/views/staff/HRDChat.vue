<template>
  <div class="hd-wrap">

    <!-- ── PANEL: INBOX ── -->
    <div v-if="activePanel === 'inbox'" class="panel">
      <div class="panel-head">
        <div>
          <div class="panel-title">Inbox Helpdesk</div>
          <div class="panel-sub">Tiket masuk dari seluruh peserta magang</div>
        </div>
        <div class="panel-head__right">
          <span v-if="totalMenunggu > 0" class="badge-menunggu">{{ totalMenunggu }} menunggu</span>
          <span class="panel-count">{{ totalTiket }} tiket</span>
        </div>
      </div>

      <!-- Filter Bar -->
      <div class="filter-bar">
        <select v-model="filterStatus" class="filter-select" @change="fetchTiket">
          <option value="">Semua Status</option>
          <option value="menunggu">Menunggu</option>
          <option value="diproses">Diproses</option>
          <option value="selesai">Selesai</option>
        </select>
        <select v-model="filterKategori" class="filter-select" @change="fetchTiket">
          <option value="">Semua Kategori</option>
          <option value="umum">Umum</option>
          <option value="absensi">Absensi</option>
          <option value="dokumen">Dokumen</option>
          <option value="sertifikat">Sertifikat</option>
        </select>
      </div>

      <div v-if="loading" class="empty-state"><div class="spinner"></div></div>
      <div v-else-if="tiketList.length === 0" class="empty-state">
        <div class="es-icon">
          <svg width="36" height="36" viewBox="0 0 24 24" fill="none"><path d="M21 15a2 2 0 01-2 2H7l-4 4V5a2 2 0 012-2h14a2 2 0 012 2z" stroke="#d1d5db" stroke-width="1.5"/></svg>
        </div>
        <p>Tidak ada tiket masuk saat ini.</p>
      </div>
      <div v-else class="tiket-list">
        <div
          v-for="t in tiketList" :key="t.id"
          class="tiket-row"
          @click="bukaChat(t)"
        >
          <div class="tiket-row__left">
            <div class="tiket-row__top">
              <span class="tiket-nomor">{{ t.nomor_tiket }}</span>
              <span :class="['kat-badge', `kat-badge--${t.kategori}`]">{{ labelKategori(t.kategori) }}</span>
              <span :class="['status-badge', `status-badge--${t.status}`]">{{ labelStatus(t.status) }}</span>
            </div>
            <div class="tiket-subjek">{{ t.subjek || '(tanpa subjek)' }}</div>
            <div class="tiket-meta">
              <span class="tiket-user">{{ t.user_nama }}</span>
              <span v-if="t.assigned_nama" class="tiket-assigned">· {{ t.assigned_nama }}</span>
              <span v-else class="tiket-unassigned">· Belum diambil</span>
            </div>
            <div v-if="t.last_message" class="tiket-preview">{{ t.last_message }}</div>
          </div>
          <div class="tiket-row__right">
            <div class="tiket-waktu">{{ formatWaktu(t.updated_at) }}</div>
            <span v-if="t.unread_count > 0" class="unread-dot">{{ t.unread_count }}</span>
          </div>
        </div>
      </div>

      <!-- Pagination -->
      <div v-if="totalPages > 1" class="pagination">
        <button class="page-btn" :disabled="page <= 1" @click="changePage(page - 1)">‹</button>
        <span class="page-info">{{ page }} / {{ totalPages }}</span>
        <button class="page-btn" :disabled="page >= totalPages" @click="changePage(page + 1)">›</button>
      </div>
    </div>

    <!-- ── PANEL: CHAT ── -->
    <div v-if="activePanel === 'chat' && tiketAktif" class="panel panel--chat">
      <div class="chat-head">
        <button class="btn-back" @click="activePanel = 'inbox'">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none"><path d="M19 12H5M12 5l-7 7 7 7" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>
          Kembali
        </button>
        <div class="chat-head__body">
          <div class="chat-head__left">
            <div class="chat-head__nomor">{{ tiketAktif.nomor_tiket }}</div>
            <div class="chat-head__subjek">{{ tiketAktif.subjek }}</div>
            <div class="chat-head__meta">
              <span :class="['kat-badge', `kat-badge--${tiketAktif.kategori}`]">{{ labelKategori(tiketAktif.kategori) }}</span>
              <span :class="['status-badge', `status-badge--${tiketAktif.status}`]">{{ labelStatus(tiketAktif.status) }}</span>
              <span class="chat-head__user">{{ tiketAktif.user_nama }}</span>
            </div>
          </div>
          <div class="chat-head__actions">
            <span v-if="tiketAktif.assigned_nama" class="assigned-tag">{{ tiketAktif.assigned_nama }}</span>
            <button
              v-if="tiketAktif.status === 'menunggu'"
              class="btn-ambil"
              @click="ambilTiket"
              :disabled="ambilLoading"
            >{{ ambilLoading ? 'Memproses...' : 'Ambil Tiket' }}</button>
            <button
              v-if="tiketAktif.status === 'diproses'"
              class="btn-selesai"
              @click="selesaikanTiket"
              :disabled="selesaiLoading"
            >{{ selesaiLoading ? '...' : 'Tandai Selesai' }}</button>
          </div>
        </div>
      </div>

      <!-- Chat Messages -->
      <div class="chat-messages" ref="msgBox">
        <div v-if="pesanLoading" class="empty-state" style="padding:20px"><div class="spinner"></div></div>
        <div v-else-if="pesanList.length === 0" class="empty-state" style="padding:24px">
          <p>Belum ada pesan dalam tiket ini.</p>
        </div>
        <template v-else>
          <div v-for="p in pesanList" :key="p.id" :class="['msg-row', p.sender_id === userId ? 'msg-row--me' : 'msg-row--them']">
            <div v-if="p.sender_id !== userId" class="msg-avatar">
              {{ p.sender_nama?.[0]?.toUpperCase() ?? 'U' }}
            </div>
            <div class="msg-col">
              <div v-if="p.sender_id !== userId" class="msg-sender">{{ p.sender_nama }}</div>
              <div class="msg-bubble">{{ p.pesan }}</div>
              <div class="msg-time">{{ formatWaktu(p.created_at) }}</div>
            </div>
          </div>
        </template>
      </div>

      <!-- Input -->
      <div v-if="tiketAktif.status === 'selesai'" class="chat-closed">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z" stroke="#6b7280" stroke-width="2"/></svg>
        Tiket ini sudah ditutup.
      </div>
      <div v-else-if="tiketAktif.status === 'menunggu'" class="chat-closed">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z" stroke="#6b7280" stroke-width="2"/></svg>
        Klik <strong style="color:#374151">Ambil Tiket</strong> terlebih dahulu untuk membalas.
      </div>
      <div v-else class="chat-input-row">
        <textarea
          v-model="inputPesan"
          class="chat-input"
          placeholder="Tulis pesan... (Enter untuk kirim)"
          rows="1"
          @keydown.enter.exact.prevent="kirimPesan"
          @input="autoResize"
          ref="inputEl"
        ></textarea>
        <button class="btn-kirim" @click="kirimPesan" :disabled="!inputPesan.trim() || kirimLoading">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none"><line x1="22" y1="2" x2="11" y2="13" stroke="currentColor" stroke-width="2" stroke-linecap="round"/><polygon points="22 2 15 22 11 13 2 9 22 2" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, nextTick } from "vue";
import { useAppWS } from "@/composables/useAppWS";
import { useAuth } from "@/hooks/useAuth";
import api from "@/lib/api";

const { user } = useAuth();
const userId = computed(() => user.value?.id ?? "");

const loading       = ref(false);
const tiketList     = ref<any[]>([]);
const tiketAktif    = ref<any>(null);
const pesanList     = ref<any[]>([]);
const pesanLoading  = ref(false);
const inputPesan    = ref("");
const kirimLoading  = ref(false);
const ambilLoading  = ref(false);
const selesaiLoading = ref(false);
const msgBox        = ref<HTMLElement | null>(null);
const inputEl       = ref<HTMLTextAreaElement | null>(null);
const activePanel   = ref<"inbox" | "chat">("inbox");

const filterStatus   = ref("");
const filterKategori = ref("");
const page           = ref(1);
const totalTiket     = ref(0);
const limit          = 20;
const totalPages     = computed(() => Math.max(1, Math.ceil(totalTiket.value / limit)));
const totalMenunggu  = computed(() => tiketList.value.filter(t => t.status === "menunggu").length);

const { connect: wsConnect, disconnect: wsDisconnect, subscribe: wsSubscribe } = useAppWS();
let wsUnsub: (() => void) | null = null;

function labelStatus(s: string) {
  return ({ menunggu: "Menunggu", diproses: "Diproses", selesai: "Selesai" } as Record<string,string>)[s] ?? s;
}
function labelKategori(k: string) {
  return ({ umum: "Umum", absensi: "Absensi", dokumen: "Dokumen", sertifikat: "Sertifikat" } as Record<string,string>)[k] ?? k;
}
function formatWaktu(iso: string) {
  if (!iso) return "–";
  const d = new Date(iso);
  const diff = (Date.now() - d.getTime()) / 1000;
  if (diff < 60)    return "Baru saja";
  if (diff < 3600)  return `${Math.floor(diff / 60)} mnt lalu`;
  if (diff < 86400) return `${Math.floor(diff / 3600)} jam lalu`;
  return d.toLocaleDateString("id-ID", { day: "2-digit", month: "short", year: "numeric" });
}
function scrollBottom() {
  if (msgBox.value) msgBox.value.scrollTop = msgBox.value.scrollHeight;
}
function autoResize(e: Event) {
  const el = e.target as HTMLTextAreaElement;
  el.style.height = "auto";
  el.style.height = Math.min(el.scrollHeight, 120) + "px";
}

async function fetchTiket() {
  loading.value = true;
  try {
    const params = new URLSearchParams({ page: String(page.value), limit: String(limit) });
    if (filterStatus.value)   params.set("status", filterStatus.value);
    if (filterKategori.value) params.set("kategori", filterKategori.value);
    const res = await api.get(`/api/chat/tiket?${params}`);
    tiketList.value  = Array.isArray(res.data?.data) ? res.data.data : [];
    totalTiket.value = res.data?.total ?? 0;
    if (tiketAktif.value) {
      const updated = tiketList.value.find(t => t.id === tiketAktif.value.id);
      if (updated) tiketAktif.value = updated;
    }
  } finally {
    loading.value = false;
  }
}

function changePage(p: number) {
  page.value = p;
  fetchTiket();
}

async function bukaChat(t: any) {
  tiketAktif.value = t;
  activePanel.value = "chat";
  pesanLoading.value = true;
  pesanList.value = [];
  try {
    const res = await api.get(`/api/chat/tiket/${t.id}/pesan`);
    pesanList.value = Array.isArray(res.data?.data) ? res.data.data : [];
  } finally {
    pesanLoading.value = false;
    await nextTick();
    scrollBottom();
  }
}

async function kirimPesan() {
  if (!inputPesan.value.trim() || !tiketAktif.value) return;
  const teks = inputPesan.value.trim();
  inputPesan.value = "";
  if (inputEl.value) inputEl.value.style.height = "auto";
  kirimLoading.value = true;
  try {
    const res = await api.post(`/api/chat/tiket/${tiketAktif.value.id}/pesan`, { pesan: teks });
    if (res.data?.data) {
      pesanList.value.push(res.data.data);
      await nextTick();
      scrollBottom();
    }
  } catch { /* silent */ } finally {
    kirimLoading.value = false;
  }
}

async function ambilTiket() {
  if (!tiketAktif.value) return;
  ambilLoading.value = true;
  try {
    const res = await api.patch(`/api/chat/tiket/${tiketAktif.value.id}/ambil`);
    if (res.data?.data) {
      tiketAktif.value = res.data.data;
      fetchTiket();
    }
  } catch { /* silent */ } finally {
    ambilLoading.value = false;
  }
}

async function selesaikanTiket() {
  if (!tiketAktif.value) return;
  selesaiLoading.value = true;
  try {
    await api.patch(`/api/chat/tiket/${tiketAktif.value.id}/status`, { status: "selesai" });
    tiketAktif.value = { ...tiketAktif.value, status: "selesai" };
    fetchTiket();
  } catch { /* silent */ } finally {
    selesaiLoading.value = false;
  }
}

onMounted(() => {
  fetchTiket();
  wsConnect();
  wsUnsub = wsSubscribe((msg) => {
    if (msg.type === "chat_message") {
      const p = msg.payload;
      if (tiketAktif.value && p.tiket_id === tiketAktif.value.id) {
        pesanList.value.push({
          id: crypto.randomUUID(),
          tiket_id: p.tiket_id,
          sender_id: p.sender_id,
          sender_nama: p.sender_nama,
          sender_role: p.sender_role,
          pesan: p.pesan,
          is_read: false,
          created_at: p.created_at,
        });
        nextTick(() => scrollBottom());
      }
      fetchTiket();
    } else if (msg.type === "tiket_update") {
      const p = msg.payload;
      // Update tiketAktif in-place agar header status/tombol langsung berubah
      if (tiketAktif.value && tiketAktif.value.id === p.tiket_id) {
        tiketAktif.value = {
          ...tiketAktif.value,
          status: p.status,
          assigned_to: p.assigned_to ?? tiketAktif.value.assigned_to,
          assigned_nama: p.assigned_nama ?? tiketAktif.value.assigned_nama,
        };
      }
      // Refresh inbox agar row status/badge unread sinkron
      fetchTiket();
    }
  });
});

onUnmounted(() => {
  wsUnsub?.();
  wsDisconnect();
});
</script>

<style scoped>
.hd-wrap { display: flex; flex-direction: column; gap: 0; }

/* Panel */
.panel {
  background: #fff; border-radius: 14px; border: 1px solid #e9f5e9;
  box-shadow: 0 1px 3px rgba(13,40,24,0.05); overflow: hidden;
}
.panel--chat { display: flex; flex-direction: column; }

/* Panel Head */
.panel-head {
  display: flex; align-items: center; justify-content: space-between;
  padding: 16px 20px; border-bottom: 1px solid #f0faf0; gap: 12px;
}
.panel-title { font-size: 14px; font-weight: 700; color: #111827; }
.panel-sub   { font-size: 12px; color: #9ca3af; margin-top: 2px; }
.panel-head__right { display: flex; align-items: center; gap: 8px; flex-shrink: 0; }
.panel-count { font-size: 12px; color: #9ca3af; }
.badge-menunggu {
  background: #fefce8; color: #ca8a04; border: 1px solid #fde68a;
  font-size: 11px; font-weight: 700; padding: 2px 8px; border-radius: 100px;
}

/* Filter Bar */
.filter-bar {
  display: flex; align-items: center; gap: 8px;
  padding: 10px 20px; border-bottom: 1px solid #f0faf0;
  background: #fafffe;
}
.filter-select {
  border: 1.5px solid #e5e7eb; border-radius: 8px;
  padding: 5px 10px; font-size: 12px; font-family: "Poppins", sans-serif;
  color: #374151; outline: none; cursor: pointer; background: #fff;
  transition: border-color .15s;
}
.filter-select:focus { border-color: #48AF4A; }

/* Tiket List */
.tiket-list { padding: 4px 0; }
.tiket-row {
  display: flex; align-items: flex-start; justify-content: space-between;
  gap: 12px; padding: 13px 20px;
  border-bottom: 1px solid #f9fafb; cursor: pointer;
  transition: background .1s;
}
.tiket-row:hover { background: #f9fafb; }
.tiket-row:last-child { border-bottom: none; }
.tiket-row__left { flex: 1; min-width: 0; }
.tiket-row__top  { display: flex; align-items: center; gap: 6px; margin-bottom: 4px; flex-wrap: wrap; }
.tiket-nomor  { font-size: 11px; font-weight: 700; color: #48AF4A; font-family: monospace; }
.tiket-subjek { font-size: 13px; font-weight: 600; color: #111827; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.tiket-meta   { display: flex; align-items: center; gap: 4px; margin-top: 3px; flex-wrap: wrap; }
.tiket-user   { font-size: 11.5px; color: #6b7280; }
.tiket-assigned  { font-size: 11.5px; color: #374151; font-weight: 500; }
.tiket-unassigned { font-size: 11px; color: #9ca3af; }
.tiket-preview { font-size: 12px; color: #9ca3af; margin-top: 3px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.tiket-row__right { display: flex; flex-direction: column; align-items: flex-end; gap: 5px; flex-shrink: 0; }
.tiket-waktu  { font-size: 11px; color: #9ca3af; white-space: nowrap; }
.unread-dot {
  background: #ef4444; color: #fff;
  font-size: 9.5px; font-weight: 700;
  padding: 1px 6px; border-radius: 100px;
}

/* Pagination */
.pagination {
  display: flex; align-items: center; justify-content: center;
  gap: 10px; padding: 12px;
  border-top: 1px solid #f0faf0;
}
.page-btn {
  background: #f3f4f6; border: none; border-radius: 7px;
  width: 30px; height: 30px; font-size: 14px; cursor: pointer;
  display: flex; align-items: center; justify-content: center;
  font-family: "Poppins", sans-serif; color: #374151;
}
.page-btn:disabled { opacity: .4; cursor: default; }
.page-btn:not(:disabled):hover { background: #e9f5e9; }
.page-info { font-size: 12.5px; color: #6b7280; }

/* Kategori & Status Badge */
.kat-badge { font-size: 10.5px; font-weight: 700; padding: 2px 8px; border-radius: 100px; }
.kat-badge--umum       { background: #f3f4f6; color: #374151; }
.kat-badge--absensi    { background: #eff6ff; color: #1d4ed8; }
.kat-badge--dokumen    { background: #f5f3ff; color: #7c3aed; }
.kat-badge--sertifikat { background: #fef9c3; color: #a16207; }

.status-badge { font-size: 10.5px; font-weight: 700; padding: 2px 8px; border-radius: 100px; }
.status-badge--menunggu { background: #fefce8; color: #ca8a04; }
.status-badge--diproses  { background: #f0fdf4; color: #1a5c20; }
.status-badge--selesai   { background: #f3f4f6; color: #6b7280; }

/* Back button */
.btn-back {
  display: inline-flex; align-items: center; gap: 6px;
  background: #f3f4f6; border: none; border-radius: 8px;
  padding: 6px 12px; font-size: 12px; font-weight: 600;
  color: #374151; cursor: pointer; font-family: "Poppins", sans-serif;
  margin-bottom: 10px; flex-shrink: 0;
}
.btn-back:hover { background: #e9f5e9; color: #1a5c20; }

/* Chat Head */
.chat-head {
  padding: 14px 18px; border-bottom: 1px solid #f0faf0;
  background: #fafffe;
}
.chat-head__body {
  display: flex; align-items: flex-start; justify-content: space-between;
  gap: 12px; flex-wrap: wrap;
}
.chat-head__left { flex: 1; min-width: 0; }
.chat-head__nomor  { font-size: 11px; font-weight: 700; color: #48AF4A; font-family: monospace; }
.chat-head__subjek { font-size: 13.5px; font-weight: 700; color: #111827; margin: 2px 0 6px; }
.chat-head__meta   { display: flex; align-items: center; gap: 8px; flex-wrap: wrap; }
.chat-head__user   { font-size: 11.5px; color: #6b7280; }
.chat-head__actions { display: flex; align-items: center; gap: 8px; flex-shrink: 0; flex-wrap: wrap; }

.btn-ambil {
  background: #48AF4A; color: #fff; border: none; border-radius: 8px;
  padding: 7px 14px; font-size: 12px; font-weight: 600;
  font-family: "Poppins", sans-serif; cursor: pointer; transition: background .15s;
}
.btn-ambil:hover { background: #3d9e3f; }
.btn-ambil:disabled { background: #9ca3af; cursor: default; }

.btn-selesai {
  background: #f9fafb; color: #374151; border: 1.5px solid #e5e7eb; border-radius: 8px;
  padding: 6px 13px; font-size: 12px; font-weight: 600;
  font-family: "Poppins", sans-serif; cursor: pointer; transition: all .15s;
}
.btn-selesai:hover { background: #f3f4f6; }
.btn-selesai:disabled { opacity: .5; cursor: default; }

.assigned-tag {
  font-size: 11.5px; color: #6b7280; font-weight: 500;
  background: #f3f4f6; padding: 4px 10px; border-radius: 7px;
}

/* Chat Messages */
.chat-messages {
  flex: 1; min-height: 260px; max-height: 400px;
  overflow-y: auto; padding: 16px 20px;
  display: flex; flex-direction: column; gap: 12px;
}
.msg-row { display: flex; align-items: flex-end; gap: 8px; }
.msg-row--me { flex-direction: row-reverse; }
.msg-avatar {
  width: 28px; height: 28px; border-radius: 50%;
  background: linear-gradient(135deg, #94a3b8, #64748b);
  color: #fff; font-size: 11px; font-weight: 700;
  display: flex; align-items: center; justify-content: center; flex-shrink: 0;
}
.msg-col { display: flex; flex-direction: column; max-width: 72%; }
.msg-row--me .msg-col { align-items: flex-end; }
.msg-sender { font-size: 10.5px; font-weight: 600; color: #6b7280; margin-bottom: 3px; }
.msg-bubble {
  background: #f3f4f6; border-radius: 12px 12px 12px 3px;
  padding: 9px 13px; font-size: 13px; color: #111827; line-height: 1.55;
  word-break: break-word; white-space: pre-wrap;
}
.msg-row--me .msg-bubble {
  background: #f0fdf4; color: #0d2818;
  border-radius: 12px 12px 3px 12px;
}
.msg-time { font-size: 10.5px; color: #9ca3af; margin-top: 3px; }

/* Chat Input */
.chat-closed {
  display: flex; align-items: center; justify-content: center; gap: 6px;
  font-size: 12px; color: #9ca3af;
  padding: 10px 16px 14px;
  border-top: 1px solid #f0faf0;
}
.chat-input-row {
  display: flex; align-items: flex-end; gap: 8px;
  padding: 10px 16px 14px; border-top: 1px solid #f0faf0;
}
.chat-input {
  flex: 1; border: 1.5px solid #e5e7eb; border-radius: 10px;
  padding: 9px 12px; font-size: 13px; font-family: "Poppins", sans-serif;
  color: #111827; outline: none; resize: none; line-height: 1.5;
  min-height: 38px; max-height: 120px; overflow-y: auto;
  transition: border-color .15s;
}
.chat-input:focus { border-color: #48AF4A; }
.btn-kirim {
  width: 36px; height: 36px; border-radius: 9px; flex-shrink: 0;
  background: #48AF4A; color: #fff; border: none; cursor: pointer;
  display: flex; align-items: center; justify-content: center;
  transition: background .15s;
}
.btn-kirim:hover { background: #3d9e3f; }
.btn-kirim:disabled { background: #d1d5db; cursor: default; }

/* Empty + Spinner */
.empty-state {
  display: flex; flex-direction: column; align-items: center;
  justify-content: center; gap: 10px; padding: 32px;
  font-size: 13px; color: #9ca3af; text-align: center;
}
.es-icon { opacity: .6; }
.spinner {
  width: 24px; height: 24px; border: 2.5px solid #e5e7eb;
  border-top-color: #48AF4A; border-radius: 50%;
  animation: spin .7s linear infinite;
}
@keyframes spin { to { transform: rotate(360deg); } }
</style>
