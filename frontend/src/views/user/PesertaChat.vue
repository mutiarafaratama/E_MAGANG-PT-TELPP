<template>
  <div class="hd-wrap">

    <!-- ── PANEL: DAFTAR TIKET ── -->
    <div v-if="activePanel === 'list'" class="panel">
      <div class="panel-head">
        <div>
          <div class="panel-title">Tiket Helpdesk</div>
          <div class="panel-sub">Hubungi HRD untuk pertanyaan seputar magang kamu</div>
        </div>
        <button
          v-if="!tiketAktifExist"
          class="btn-primary"
          @click="showBuatTiket = true"
        >+ Buat Tiket</button>
        <div v-else class="tiket-aktif-note">
          <svg width="13" height="13" viewBox="0 0 24 24" fill="none"><circle cx="12" cy="12" r="10" stroke="#ca8a04" stroke-width="2"/><line x1="12" y1="8" x2="12" y2="12" stroke="#ca8a04" stroke-width="2" stroke-linecap="round"/><circle cx="12" cy="16" r="1" fill="#ca8a04"/></svg>
          Tiket aktif sedang berjalan
        </div>
      </div>

      <div v-if="loading" class="empty-state"><div class="spinner"></div></div>
      <div v-else-if="tiketList.length === 0" class="empty-state">
        <div class="es-icon">
          <svg width="36" height="36" viewBox="0 0 24 24" fill="none"><path d="M21 15a2 2 0 01-2 2H7l-4 4V5a2 2 0 012-2h14a2 2 0 012 2z" stroke="#d1d5db" stroke-width="1.5"/></svg>
        </div>
        <p>Belum ada tiket.<br/>Klik <strong>+ Buat Tiket</strong> untuk menghubungi HRD.</p>
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
            </div>
            <div class="tiket-subjek">{{ t.subjek || '(tanpa subjek)' }}</div>
            <div v-if="t.last_message" class="tiket-preview">{{ t.last_message }}</div>
          </div>
          <div class="tiket-row__right">
            <span :class="['status-badge', `status-badge--${t.status}`]">{{ labelStatus(t.status) }}</span>
            <div class="tiket-waktu">{{ formatWaktu(t.updated_at) }}</div>
            <span v-if="t.unread_count > 0" class="unread-dot">{{ t.unread_count }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- ── PANEL: CHAT ── -->
    <div v-if="activePanel === 'chat' && tiketAktif" class="panel panel--chat">
      <div class="chat-head">
        <button class="btn-back" @click="bukaList">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none"><path d="M19 12H5M12 5l-7 7 7 7" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>
          Kembali
        </button>
        <div class="chat-head__info">
          <div class="chat-head__nomor">{{ tiketAktif.nomor_tiket }}</div>
          <div class="chat-head__subjek">{{ tiketAktif.subjek }}</div>
          <div class="chat-head__meta">
            <span :class="['kat-badge', `kat-badge--${tiketAktif.kategori}`]">{{ labelKategori(tiketAktif.kategori) }}</span>
            <span :class="['status-badge', `status-badge--${tiketAktif.status}`]">{{ labelStatus(tiketAktif.status) }}</span>
            <span v-if="tiketAktif.assigned_nama" class="chat-head__hrd">
              <svg width="11" height="11" viewBox="0 0 24 24" fill="none"><path d="M20 21v-2a4 4 0 00-4-4H8a4 4 0 00-4 4v2" stroke="#6b7280" stroke-width="2"/><circle cx="12" cy="7" r="4" stroke="#6b7280" stroke-width="2"/></svg>
              {{ tiketAktif.assigned_nama }}
            </span>
            <span v-else class="chat-head__hrd chat-head__hrd--waiting">Menunggu HRD...</span>
          </div>
        </div>
      </div>

      <div class="chat-messages" ref="msgBox">
        <div v-if="pesanLoading" class="empty-state" style="padding:20px"><div class="spinner"></div></div>
        <div v-else-if="pesanList.length === 0" class="empty-state" style="padding:24px">
          <p>Belum ada pesan. Kirim pesan pertama di bawah.</p>
        </div>
        <template v-else>
          <div v-for="p in pesanList" :key="p.id" :class="['msg-row', p.sender_id === userId ? 'msg-row--me' : 'msg-row--them']">
            <div v-if="p.sender_id !== userId" class="msg-avatar">{{ p.sender_nama?.[0]?.toUpperCase() ?? 'H' }}</div>
            <div class="msg-col">
              <div v-if="p.sender_id !== userId" class="msg-sender">{{ p.sender_nama }}</div>
              <div class="msg-bubble">{{ p.pesan }}</div>
              <div class="msg-time">{{ formatWaktu(p.created_at) }}</div>
            </div>
          </div>
        </template>
      </div>

      <div v-if="tiketAktif.status === 'selesai'" class="chat-closed">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z" stroke="#6b7280" stroke-width="2"/></svg>
        Tiket ini sudah ditutup.
      </div>
      <div v-else class="chat-input-row">
        <textarea
          v-model="inputPesan"
          class="chat-input"
          placeholder="Tulis pesan..."
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

  <!-- Modal Buat Tiket -->
  <Teleport to="body">
    <div v-if="showBuatTiket" class="modal-backdrop" @click.self="showBuatTiket = false">
      <div class="modal-box">
        <div class="modal-title">Buat Tiket Helpdesk</div>
        <div class="form-field">
          <label class="form-label">Kategori <span class="req">*</span></label>
          <select v-model="newKategori" class="form-select">
            <option value="umum">Umum</option>
            <option value="absensi">Absensi</option>
            <option value="dokumen">Dokumen</option>
            <option value="sertifikat">Sertifikat</option>
          </select>
        </div>
        <div class="form-field" style="margin-top:12px">
          <label class="form-label">Subjek <span class="req">*</span></label>
          <input v-model="newSubjek" class="form-input" placeholder="Topik pertanyaan Anda..." />
        </div>
        <div class="form-field" style="margin-top:12px">
          <label class="form-label">Pesan Awal <span class="req">*</span></label>
          <textarea v-model="newPesan" class="form-textarea" rows="4" placeholder="Tuliskan pertanyaan atau kendala Anda..."></textarea>
        </div>
        <div v-if="buatError" class="alert-err">{{ buatError }}</div>
        <div class="modal-actions">
          <button class="btn-cancel" @click="showBuatTiket = false">Batal</button>
          <button class="btn-confirm" @click="buatTiket" :disabled="buatLoading">{{ buatLoading ? 'Mengirim...' : 'Kirim Tiket' }}</button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, nextTick, computed } from "vue";
import { useAppWS } from "@/composables/useAppWS";
import { useAuth } from "@/hooks/useAuth";
import api from "@/lib/api";

const { user } = useAuth();
const userId = computed(() => user.value?.id ?? "");

const loading      = ref(false);
const tiketList    = ref<any[]>([]);
const tiketAktif   = ref<any>(null);
const pesanList    = ref<any[]>([]);
const pesanLoading = ref(false);
const inputPesan   = ref("");
const kirimLoading = ref(false);
const msgBox       = ref<HTMLElement | null>(null);
const inputEl      = ref<HTMLTextAreaElement | null>(null);
const activePanel  = ref<"list" | "chat">("list");

const showBuatTiket = ref(false);
const newKategori  = ref("umum");
const newSubjek    = ref("");
const newPesan     = ref("");
const buatLoading  = ref(false);
const buatError    = ref("");

// Apakah ada tiket aktif (status != selesai)
const tiketAktifExist = computed(() =>
  tiketList.value.some(t => t.status !== "selesai")
);

// ── WebSocket ─────────────────────────────────────────────────────────
const { connect: wsConnect, disconnect: wsDisconnect, subscribe: wsSubscribe } = useAppWS();
let wsUnsub: (() => void) | null = null;

// ── Helpers ────────────────────────────────────────────────────────────
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

// ── API Calls ─────────────────────────────────────────────────────────
async function fetchTiket() {
  loading.value = true;
  try {
    const res = await api.get("/api/chat/tiket/saya");
    tiketList.value = Array.isArray(res.data?.data) ? res.data.data : [];
    // Sync tiketAktif jika panel chat sedang terbuka
    if (tiketAktif.value) {
      const updated = tiketList.value.find(t => t.id === tiketAktif.value.id);
      if (updated) tiketAktif.value = updated;
    }
  } catch {
    tiketList.value = [];
  } finally {
    loading.value = false;
  }
}

function bukaList() {
  activePanel.value = "list";
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

async function buatTiket() {
  buatError.value = "";
  if (!newSubjek.value.trim() || !newPesan.value.trim()) {
    buatError.value = "Subjek dan pesan wajib diisi"; return;
  }
  buatLoading.value = true;
  try {
    const res = await api.post("/api/chat/tiket", {
      subjek: newSubjek.value.trim(),
      pesan: newPesan.value.trim(),
      kategori: newKategori.value,
    });
    showBuatTiket.value = false;
    newSubjek.value = ""; newPesan.value = ""; newKategori.value = "umum";
    await fetchTiket();
    if (res.data?.data) bukaChat(res.data.data);
  } catch (e: any) {
    buatError.value = e.response?.data?.message ?? "Gagal membuat tiket";
  } finally {
    buatLoading.value = false;
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
      // Update tiketAktif in-place jika sedang dibuka
      if (tiketAktif.value && tiketAktif.value.id === p.tiket_id) {
        tiketAktif.value = {
          ...tiketAktif.value,
          status: p.status,
          assigned_to: p.assigned_to ?? tiketAktif.value.assigned_to,
          assigned_nama: p.assigned_nama ?? tiketAktif.value.assigned_nama,
        };
      }
      // Refresh daftar tiket agar preview row juga update
      fetchTiket();
    } else if (msg.type === "notifikasi") {
      // Refresh tiket saat ada notifikasi chat (e.g. tiket dibuat, balasan baru)
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

/* Panels */
.panel {
  background: #fff; border-radius: 14px; border: 1px solid #e9f5e9;
  box-shadow: 0 1px 3px rgba(13,40,24,0.05); overflow: hidden;
}
.panel--chat { display: flex; flex-direction: column; }

.panel-head {
  display: flex; align-items: center; justify-content: space-between;
  padding: 16px 20px; border-bottom: 1px solid #f0faf0;
}
.panel-title  { font-size: 14px; font-weight: 700; color: #111827; }
.panel-sub    { font-size: 12px; color: #9ca3af; margin-top: 2px; }

.tiket-aktif-note {
  display: flex; align-items: center; gap: 5px;
  font-size: 11.5px; color: #ca8a04; font-weight: 600;
  background: #fefce8; border: 1px solid #fde68a;
  padding: 5px 11px; border-radius: 8px;
}

/* Tiket List */
.tiket-list { padding: 4px 0; }
.tiket-row {
  display: flex; align-items: flex-start; justify-content: space-between;
  gap: 12px; padding: 14px 20px;
  border-bottom: 1px solid #f9fafb; cursor: pointer;
  transition: background .1s;
}
.tiket-row:hover { background: #f9fafb; }
.tiket-row:last-child { border-bottom: none; }
.tiket-row__left { flex: 1; min-width: 0; }
.tiket-row__top  { display: flex; align-items: center; gap: 8px; margin-bottom: 4px; }
.tiket-nomor  { font-size: 11px; font-weight: 700; color: #48AF4A; font-family: monospace; }
.tiket-subjek { font-size: 13px; font-weight: 600; color: #111827; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.tiket-preview { font-size: 12px; color: #9ca3af; margin-top: 3px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.tiket-row__right { display: flex; flex-direction: column; align-items: flex-end; gap: 5px; flex-shrink: 0; }
.tiket-waktu  { font-size: 11px; color: #9ca3af; }
.unread-dot {
  background: #ef4444; color: #fff;
  font-size: 9.5px; font-weight: 700;
  padding: 1px 5px; border-radius: 100px;
}

/* Kategori Badge */
.kat-badge { font-size: 10.5px; font-weight: 700; padding: 2px 8px; border-radius: 100px; }
.kat-badge--umum       { background: #f3f4f6; color: #374151; }
.kat-badge--absensi    { background: #eff6ff; color: #1d4ed8; }
.kat-badge--dokumen    { background: #f5f3ff; color: #7c3aed; }
.kat-badge--sertifikat { background: #fef9c3; color: #a16207; }

/* Status Badge */
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
  margin-bottom: 10px;
}
.btn-back:hover { background: #e9f5e9; color: #1a5c20; }

/* Chat Head */
.chat-head {
  padding: 14px 20px; border-bottom: 1px solid #f0faf0;
  background: #fafffe;
}
.chat-head__nomor  { font-size: 11px; font-weight: 700; color: #48AF4A; font-family: monospace; }
.chat-head__subjek { font-size: 13.5px; font-weight: 700; color: #111827; margin: 2px 0 6px; }
.chat-head__meta   { display: flex; align-items: center; gap: 8px; flex-wrap: wrap; }
.chat-head__hrd    { display: flex; align-items: center; gap: 4px; font-size: 11.5px; color: #6b7280; }
.chat-head__hrd--waiting { color: #9ca3af; font-style: italic; }

/* Chat Messages */
.chat-messages {
  flex: 1; min-height: 260px; max-height: 420px;
  overflow-y: auto; padding: 16px 20px;
  display: flex; flex-direction: column; gap: 12px;
}
.msg-row {
  display: flex; align-items: flex-end; gap: 8px;
}
.msg-row--me { flex-direction: row-reverse; }
.msg-avatar {
  width: 28px; height: 28px; border-radius: 50%;
  background: linear-gradient(135deg, #48AF4A, #1a5c20);
  color: #fff; font-size: 11px; font-weight: 700;
  display: flex; align-items: center; justify-content: center;
  flex-shrink: 0;
}
.msg-col { display: flex; flex-direction: column; max-width: 70%; }
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
  padding: 10px 14px; border-top: 1px solid #f0faf0;
}
.chat-input {
  flex: 1; border: 1.5px solid #e5e7eb; border-radius: 10px;
  padding: 9px 13px; font-size: 13px; font-family: "Poppins", sans-serif;
  resize: none; outline: none; overflow: hidden;
  min-height: 38px; max-height: 120px;
  transition: border-color .15s; line-height: 1.5;
}
.chat-input:focus { border-color: #48AF4A; }
.btn-kirim {
  background: #48AF4A; color: #fff; border: none; border-radius: 10px;
  width: 38px; height: 38px; display: flex; align-items: center; justify-content: center;
  cursor: pointer; flex-shrink: 0; transition: background .15s;
}
.btn-kirim:hover { background: #3d9e3f; }
.btn-kirim:disabled { background: #d1d5db; cursor: default; }

/* Empty & Spinner */
.empty-state {
  display: flex; flex-direction: column; align-items: center;
  padding: 44px 24px; gap: 12px; text-align: center;
}
.es-icon {
  width: 72px; height: 72px; background: #f9fafb; border-radius: 50%;
  display: flex; align-items: center; justify-content: center;
}
.empty-state p { font-size: 13px; color: #9ca3af; line-height: 1.7; margin: 0; }
.spinner {
  width: 34px; height: 34px; border: 3px solid #e5e7eb;
  border-top-color: #48AF4A; border-radius: 50%;
  animation: spin .8s linear infinite;
}
@keyframes spin { to { transform: rotate(360deg); } }

/* Primary Button */
.btn-primary {
  background: #48AF4A; color: #fff; border: none; border-radius: 9px;
  padding: 7px 16px; font-size: 12.5px; font-weight: 600;
  font-family: "Poppins", sans-serif; cursor: pointer; transition: background .15s;
}
.btn-primary:hover { background: #3d9e3f; }

/* Modal */
.modal-backdrop {
  position: fixed; inset: 0; background: rgba(0,0,0,0.4);
  display: flex; align-items: center; justify-content: center; z-index: 1000;
}
.modal-box {
  background: #fff; border-radius: 16px; padding: 28px 24px;
  width: 100%; max-width: 420px; box-shadow: 0 20px 60px rgba(0,0,0,0.15);
}
.modal-title { font-size: 16px; font-weight: 700; color: #111827; margin-bottom: 16px; }
.modal-actions { display: flex; gap: 10px; justify-content: flex-end; margin-top: 18px; }
.form-field { display: flex; flex-direction: column; gap: 5px; }
.form-label { font-size: 12px; font-weight: 600; color: #374151; }
.req { color: #ef4444; }
.form-input, .form-select {
  border: 1.5px solid #e5e7eb; border-radius: 9px;
  padding: 9px 12px; font-size: 13px; font-family: "Poppins", sans-serif;
  outline: none; transition: border-color .15s;
}
.form-input:focus, .form-select:focus { border-color: #48AF4A; }
.form-textarea {
  width: 100%; border: 1.5px solid #e5e7eb; border-radius: 9px;
  padding: 9px 12px; font-size: 13px; font-family: "Poppins", sans-serif;
  resize: vertical; outline: none; transition: border-color .15s; box-sizing: border-box;
}
.form-textarea:focus { border-color: #48AF4A; }
.alert-err {
  background: #fff1f2; color: #be123c; border: 1px solid #fecdd3;
  padding: 8px 12px; border-radius: 8px; font-size: 12.5px; margin-top: 10px;
}
.btn-cancel {
  background: #f3f4f6; color: #374151; border: none; border-radius: 9px;
  padding: 9px 20px; font-size: 13px; font-weight: 600; font-family: "Poppins",sans-serif; cursor: pointer;
}
.btn-confirm {
  background: #48AF4A; color: #fff; border: none; border-radius: 9px;
  padding: 9px 20px; font-size: 13px; font-weight: 600; font-family: "Poppins",sans-serif; cursor: pointer;
}
.btn-confirm:disabled { background: #d1d5db; cursor: default; }
</style>
