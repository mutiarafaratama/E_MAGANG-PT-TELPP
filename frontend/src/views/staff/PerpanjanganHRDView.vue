<template>
  <div>
    <!-- Header + tombol langsung -->
    <div class="page-head">
      <div>
        <div class="page-title">Perpanjangan Magang</div>
        <div class="page-sub">Kelola pengajuan perpanjangan dari peserta dan perpanjang langsung jika diperlukan.</div>
      </div>
      <button class="btn-langsung" @click="showModalLangsung = true">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none">
          <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2"/>
          <line x1="12" y1="8" x2="12" y2="16" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
          <line x1="8" y1="12" x2="16" y2="12" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
        </svg>
        Perpanjang Langsung
      </button>
    </div>

    <!-- Card: Tab + Tabel -->
    <div class="card">
      <div class="tab-bar">
        <button v-for="t in tabs" :key="t.key" class="tab-btn" :class="{ 'tab-btn--active': activeTab === t.key }" @click="activeTab = t.key">
          {{ t.label }}
          <span v-if="t.key === 'menunggu' && pendingCount > 0" class="tab-badge">{{ pendingCount }}</span>
        </button>
      </div>

      <div v-if="loading" class="empty-state"><div class="spinner"></div></div>

      <div v-else-if="filtered.length === 0" class="empty-state">
        <svg width="36" height="36" viewBox="0 0 24 24" fill="none">
          <rect x="3" y="4" width="18" height="18" rx="2" stroke="#d1d5db" stroke-width="1.5"/>
          <line x1="3" y1="10" x2="21" y2="10" stroke="#d1d5db" stroke-width="1.5"/>
        </svg>
        <p>{{ emptyLabel }}</p>
      </div>

      <div v-else class="table-wrap">
        <table class="data-table">
          <thead>
            <tr>
              <th>Peserta</th>
              <th>Divisi</th>
              <th>Pengaju</th>
              <th>Durasi</th>
              <th>Tanggal Lama → Baru</th>
              <th>Alasan</th>
              <th>Status</th>
              <th>Aksi</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="p in filtered" :key="p.id">
              <td>
                <div class="name-cell">
                  <div class="name-avatar">{{ p.nama_peserta?.[0]?.toUpperCase() ?? '?' }}</div>
                  <div class="name-text">{{ p.nama_peserta ?? '–' }}</div>
                </div>
              </td>
              <td class="text-muted">{{ p.divisi ?? '–' }}</td>
              <td>
                <span class="role-chip" :class="p.role_pengaju === 'hrd' ? 'chip--blue' : 'chip--gray'">
                  {{ p.role_pengaju === 'hrd' ? 'HRD Langsung' : 'Peserta' }}
                </span>
              </td>
              <td><strong class="durasi-val">+{{ p.durasi_hari }} hari</strong></td>
              <td class="date-cell">
                <div>{{ fmtDate(p.tanggal_selesai_lama) }}</div>
                <div class="arrow">↓</div>
                <div class="date-baru">{{ fmtDate(p.tanggal_selesai_baru) }}</div>
              </td>
              <td class="alasan-cell">{{ p.alasan }}</td>
              <td>
                <span class="status-chip" :class="`chip-status--${p.status}`">{{ statusLabel(p.status) }}</span>
              </td>
              <td>
                <div class="action-cell">
                  <!-- Lihat surat jika ada -->
                  <a v-if="p.surat_path" :href="`/uploads/${p.surat_path}`" target="_blank" class="btn-icon" title="Lihat Surat">
                    <svg width="13" height="13" viewBox="0 0 24 24" fill="none">
                      <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z" stroke="currentColor" stroke-width="2"/>
                      <circle cx="12" cy="12" r="3" stroke="currentColor" stroke-width="2"/>
                    </svg>
                  </a>
                  <!-- Approve -->
                  <button v-if="p.status === 'menunggu'" class="btn-approve" @click="openApprove(p)" :disabled="actionLoading === p.id">
                    Setujui
                  </button>
                  <!-- Tolak -->
                  <button v-if="p.status === 'menunggu'" class="btn-tolak" @click="openTolak(p)" :disabled="actionLoading === p.id">
                    Tolak
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- MODAL: Approve (opsional upload surat) -->
    <div v-if="modalApprove" class="modal-backdrop" @click.self="modalApprove = null">
      <div class="modal">
        <div class="modal-head">
          <div class="modal-title">Setujui Perpanjangan</div>
          <button class="modal-close" @click="modalApprove = null">✕</button>
        </div>
        <div class="modal-body">
          <p class="modal-info">
            Perpanjangan <strong>{{ modalApprove.durasi_hari }} hari</strong> untuk
            <strong>{{ modalApprove.nama_peserta }}</strong>.
            Tanggal selesai baru: <strong>{{ fmtDate(modalApprove.tanggal_selesai_baru) }}</strong>.
          </p>
          <div class="field">
            <label class="label">Upload Surat Perpanjangan <span class="hint">(opsional)</span></label>
            <input type="file" accept=".pdf,.jpg,.jpeg,.png" @change="onApproveFile" class="file-input" />
          </div>
          <div class="field">
            <label class="label">Catatan (opsional)</label>
            <textarea v-model="approveCatatan" class="textarea" rows="2" placeholder="Catatan untuk peserta..."></textarea>
          </div>
          <div v-if="actionErr" class="alert alert--error">{{ actionErr }}</div>
        </div>
        <div class="modal-foot">
          <button class="btn-cancel" @click="modalApprove = null">Batal</button>
          <button class="btn-ok" :disabled="actionLoading === modalApprove?.id" @click="doApprove">
            <span v-if="actionLoading === modalApprove?.id" class="spinner-sm"></span>
            <span v-else>Setujui</span>
          </button>
        </div>
      </div>
    </div>

    <!-- MODAL: Tolak -->
    <div v-if="modalTolak" class="modal-backdrop" @click.self="modalTolak = null">
      <div class="modal">
        <div class="modal-head">
          <div class="modal-title">Tolak Pengajuan</div>
          <button class="modal-close" @click="modalTolak = null">✕</button>
        </div>
        <div class="modal-body">
          <p class="modal-info">Tolak perpanjangan <strong>{{ modalTolak.durasi_hari }} hari</strong> dari <strong>{{ modalTolak.nama_peserta }}</strong>?</p>
          <div class="field">
            <label class="label">Catatan Penolakan (opsional)</label>
            <textarea v-model="tolakCatatan" class="textarea" rows="2" placeholder="Alasan penolakan..."></textarea>
          </div>
          <div v-if="actionErr" class="alert alert--error">{{ actionErr }}</div>
        </div>
        <div class="modal-foot">
          <button class="btn-cancel" @click="modalTolak = null">Batal</button>
          <button class="btn-tolak-ok" :disabled="actionLoading === modalTolak?.id" @click="doTolak">
            <span v-if="actionLoading === modalTolak?.id" class="spinner-sm"></span>
            <span v-else>Tolak Pengajuan</span>
          </button>
        </div>
      </div>
    </div>

    <!-- MODAL: Perpanjang Langsung -->
    <div v-if="showModalLangsung" class="modal-backdrop" @click.self="showModalLangsung = false">
      <div class="modal modal--wide">
        <div class="modal-head">
          <div class="modal-title">Perpanjang Langsung</div>
          <button class="modal-close" @click="showModalLangsung = false">✕</button>
        </div>
        <div class="modal-body">
          <div class="modal-hint">Perpanjangan ini langsung diterapkan tanpa perlu persetujuan tambahan.</div>

          <div class="field">
            <label class="label">Pilih Peserta Aktif <span class="req">*</span></label>
            <div v-if="loadingPeserta" class="text-muted" style="font-size:12px">Memuat daftar peserta...</div>
            <select v-else v-model="langsung.pelaksanaan_id" class="select">
              <option value="">— Pilih peserta —</option>
              <option v-for="ps in pesertaAktif" :key="ps.id" :value="ps.id">
                {{ ps.nama_peserta }} — {{ ps.divisi ?? 'Tanpa divisi' }} (s/d {{ fmtDate(ps.tanggal_selesai) }})
              </option>
            </select>
          </div>

          <div class="field">
            <label class="label">Durasi Perpanjangan <span class="req">*</span></label>
            <div class="durasi-grid">
              <button v-for="d in durasiOptions" :key="d" class="durasi-btn" :class="{ 'durasi-btn--active': langsung.durasi_hari === d }" @click="langsung.durasi_hari = d">
                {{ d }} hari
              </button>
            </div>
          </div>

          <div class="field">
            <label class="label">Alasan Perpanjangan <span class="req">*</span></label>
            <textarea v-model="langsung.alasan" class="textarea" rows="3" placeholder="Alasan perpanjangan (min. 10 karakter)..."></textarea>
          </div>

          <div class="field">
            <label class="label">Upload Surat Perpanjangan <span class="hint">(opsional)</span></label>
            <input type="file" accept=".pdf,.jpg,.jpeg,.png" @change="onLangsungFile" class="file-input" />
          </div>

          <div v-if="langsungErr" class="alert alert--error">{{ langsungErr }}</div>
          <div v-if="langsungOk" class="alert alert--ok">{{ langsungOk }}</div>
        </div>
        <div class="modal-foot">
          <button class="btn-cancel" @click="showModalLangsung = false">Batal</button>
          <button class="btn-ok" :disabled="langsungLoading" @click="doLangsung">
            <span v-if="langsungLoading" class="spinner-sm"></span>
            <span v-else>Terapkan Perpanjangan</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import api from "@/lib/api";

const tabs = [
  { key: "semua",     label: "Semua" },
  { key: "menunggu",  label: "Menunggu" },
  { key: "disetujui", label: "Disetujui" },
  { key: "ditolak",   label: "Ditolak" },
];
const durasiOptions = [7, 14, 30, 60, 90];

const activeTab    = ref("semua");
const loading      = ref(true);
const list         = ref<any[]>([]);
const actionLoading = ref<string | null>(null);
const actionErr    = ref("");

const modalApprove  = ref<any>(null);
const modalTolak    = ref<any>(null);
const approveCatatan = ref("");
const tolakCatatan  = ref("");
const approveFile   = ref<File | null>(null);

const showModalLangsung = ref(false);
const loadingPeserta    = ref(false);
const pesertaAktif      = ref<any[]>([]);
const langsungLoading   = ref(false);
const langsungErr       = ref("");
const langsungOk        = ref("");
const langsungFile      = ref<File | null>(null);
const langsung = ref({ pelaksanaan_id: "", durasi_hari: 0, alasan: "" });

const filtered = computed(() => {
  if (activeTab.value === "semua") return list.value;
  return list.value.filter(p => p.status === activeTab.value);
});

const pendingCount = computed(() => list.value.filter(p => p.status === "menunggu").length);

const emptyLabel = computed(() => {
  const m: Record<string, string> = {
    semua: "Belum ada pengajuan perpanjangan",
    menunggu: "Tidak ada pengajuan yang menunggu",
    disetujui: "Belum ada yang disetujui",
    ditolak: "Belum ada yang ditolak",
  };
  return m[activeTab.value] ?? "Tidak ada data";
});

function statusLabel(s: string) {
  return { menunggu: "Menunggu", disetujui: "Disetujui", ditolak: "Ditolak" }[s] ?? s;
}

function fmtDate(d: string) {
  if (!d) return "—";
  return new Date(d).toLocaleDateString("id-ID", { day: "numeric", month: "short", year: "numeric" });
}

async function fetchList() {
  loading.value = true;
  try {
    const r = await api.get("/api/perpanjangan");
    list.value = Array.isArray(r.data?.data) ? r.data.data : [];
  } catch { list.value = []; }
  finally { loading.value = false; }
}

async function fetchPesertaAktif() {
  loadingPeserta.value = true;
  try {
    const r = await api.get("/api/pelaksanaan?status=aktif&limit=200");
    const all = Array.isArray(r.data?.data?.list) ? r.data.data.list : [];
    pesertaAktif.value = all.filter((p: any) => !p.sudah_diperpanjang);
  } catch { pesertaAktif.value = []; }
  finally { loadingPeserta.value = false; }
}

// Open modal approve
function openApprove(p: any) {
  modalApprove.value  = p;
  approveCatatan.value = "";
  approveFile.value    = null;
  actionErr.value      = "";
}
function openTolak(p: any) {
  modalTolak.value  = p;
  tolakCatatan.value = "";
  actionErr.value   = "";
}

function onApproveFile(e: Event) {
  approveFile.value = (e.target as HTMLInputElement).files?.[0] ?? null;
}
function onLangsungFile(e: Event) {
  langsungFile.value = (e.target as HTMLInputElement).files?.[0] ?? null;
}

async function doApprove() {
  if (!modalApprove.value) return;
  actionErr.value = "";
  actionLoading.value = modalApprove.value.id;
  try {
    const fd = new FormData();
    fd.append("catatan_hrd", approveCatatan.value);
    if (approveFile.value) fd.append("surat", approveFile.value);
    await api.patch(`/api/perpanjangan/${modalApprove.value.id}/approve`, fd, {
      headers: { "Content-Type": "multipart/form-data" },
    });
    modalApprove.value = null;
    await fetchList();
  } catch (e: any) {
    actionErr.value = e?.response?.data?.message ?? "Gagal menyetujui pengajuan.";
  } finally { actionLoading.value = null; }
}

async function doTolak() {
  if (!modalTolak.value) return;
  actionErr.value = "";
  actionLoading.value = modalTolak.value.id;
  try {
    await api.patch(`/api/perpanjangan/${modalTolak.value.id}/tolak`, { catatan_hrd: tolakCatatan.value });
    modalTolak.value = null;
    await fetchList();
  } catch (e: any) {
    actionErr.value = e?.response?.data?.message ?? "Gagal menolak pengajuan.";
  } finally { actionLoading.value = null; }
}

async function doLangsung() {
  langsungErr.value = "";
  langsungOk.value  = "";
  if (!langsung.value.pelaksanaan_id) { langsungErr.value = "Pilih peserta terlebih dahulu."; return; }
  if (!langsung.value.durasi_hari)    { langsungErr.value = "Pilih durasi perpanjangan."; return; }
  if (langsung.value.alasan.length < 10) { langsungErr.value = "Alasan minimal 10 karakter."; return; }

  langsungLoading.value = true;
  try {
    const fd = new FormData();
    fd.append("pelaksanaan_id", langsung.value.pelaksanaan_id);
    fd.append("durasi_hari",    String(langsung.value.durasi_hari));
    fd.append("alasan",         langsung.value.alasan);
    if (langsungFile.value) fd.append("surat", langsungFile.value);
    await api.post("/api/perpanjangan/langsung", fd, {
      headers: { "Content-Type": "multipart/form-data" },
    });
    langsungOk.value = `Perpanjangan ${langsung.value.durasi_hari} hari berhasil diterapkan!`;
    langsung.value = { pelaksanaan_id: "", durasi_hari: 0, alasan: "" };
    langsungFile.value = null;
    await fetchList();
    await fetchPesertaAktif();
    setTimeout(() => { showModalLangsung.value = false; langsungOk.value = ""; }, 2000);
  } catch (e: any) {
    langsungErr.value = e?.response?.data?.message ?? "Gagal menerapkan perpanjangan.";
  } finally { langsungLoading.value = false; }
}

onMounted(() => {
  fetchList();
  fetchPesertaAktif();
});
</script>

<style scoped>
.page-head { display: flex; align-items: flex-start; justify-content: space-between; gap: 12px; margin-bottom: 16px; }
.page-title { font-size: 16px; font-weight: 700; color: #111827; }
.page-sub   { font-size: 12.5px; color: #6b7280; margin-top: 2px; }
.btn-langsung { display: flex; align-items: center; gap: 6px; padding: 9px 16px; background: #16a34a; color: #fff; border: none; border-radius: 9px; font-size: 13px; font-weight: 600; cursor: pointer; white-space: nowrap; }
.btn-langsung:hover { background: #15803d; }

/* Tabs */
.tab-bar { display: flex; gap: 4px; padding: 14px 16px 0; border-bottom: 1px solid #f0faf0; margin-bottom: 0; }
.tab-btn { padding: 7px 16px; border-radius: 8px; border: 1.5px solid #e5e7eb; background: #fff; font-size: 12.5px; font-weight: 500; color: #6b7280; cursor: pointer; display: flex; align-items: center; gap: 6px; }
.tab-btn:hover { border-color: #d1d5db; color: #374151; }
.tab-btn--active { background: #f0fdf4; border-color: #86efac; color: #15803d; font-weight: 700; }
.tab-badge { background: #ef4444; color: #fff; font-size: 10px; font-weight: 700; padding: 1px 6px; border-radius: 100px; }

/* Card + table */
.card { background: #fff; border: 1px solid #e5e7eb; border-radius: 12px; overflow: hidden; }
.empty-state { display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 10px; padding: 48px; color: #9ca3af; font-size: 13px; }
.table-wrap { overflow-x: auto; }
.data-table { width: 100%; border-collapse: collapse; font-size: 13px; }
.data-table th { padding: 11px 14px; text-align: left; font-size: 10.5px; font-weight: 600; color: #6b7280; background: #f9fafb; border-bottom: 1px solid #f1f5f9; text-transform: uppercase; letter-spacing: .04em; white-space: nowrap; }
.data-table td { padding: 13px 14px; border-bottom: 1px solid #f9fafb; color: #374151; vertical-align: middle; }
.data-table tbody tr:last-child td { border-bottom: none; }
.text-muted { color: #9ca3af; }

.name-cell   { display: flex; align-items: center; gap: 8px; }
.name-avatar { width: 30px; height: 30px; border-radius: 50%; background: #dbeafe; color: #1d4ed8; font-size: 12px; font-weight: 700; display: flex; align-items: center; justify-content: center; flex-shrink: 0; }
.name-text   { font-weight: 500; font-size: 13px; color: #111827; }

.role-chip { padding: 2px 9px; border-radius: 100px; font-size: 10.5px; font-weight: 600; }
.chip--blue { background: #eff6ff; color: #1d4ed8; }
.chip--gray { background: #f3f4f6; color: #6b7280; }

.durasi-val { font-size: 14px; color: #15803d; }

.date-cell { font-size: 12px; color: #374151; line-height: 1.5; }
.arrow     { color: #9ca3af; font-size: 11px; }
.date-baru { font-weight: 600; color: #15803d; }

.alasan-cell { max-width: 180px; font-size: 12px; color: #4b5563; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }

.status-chip { padding: 3px 10px; border-radius: 100px; font-size: 11px; font-weight: 600; white-space: nowrap; }
.chip-status--menunggu  { background: #fefce8; color: #ca8a04; }
.chip-status--disetujui { background: #f0fdf4; color: #15803d; }
.chip-status--ditolak   { background: #fef2f2; color: #dc2626; }

.action-cell { display: flex; align-items: center; gap: 6px; }
.btn-icon { display: flex; align-items: center; justify-content: center; width: 28px; height: 28px; border: 1px solid #e5e7eb; border-radius: 7px; background: #fff; color: #374151; cursor: pointer; text-decoration: none; }
.btn-approve { padding: 5px 12px; background: #16a34a; color: #fff; border: none; border-radius: 7px; font-size: 12px; font-weight: 600; cursor: pointer; }
.btn-approve:hover { background: #15803d; }
.btn-approve:disabled { opacity: 0.55; cursor: not-allowed; }
.btn-tolak { padding: 5px 12px; background: #fff; color: #dc2626; border: 1.5px solid #fca5a5; border-radius: 7px; font-size: 12px; font-weight: 600; cursor: pointer; }
.btn-tolak:hover { background: #fef2f2; }
.btn-tolak:disabled { opacity: 0.55; cursor: not-allowed; }

/* Modal */
.modal-backdrop { position: fixed; inset: 0; background: rgba(0,0,0,.4); display: flex; align-items: center; justify-content: center; z-index: 200; }
.modal { background: #fff; border-radius: 14px; width: 440px; max-width: 95vw; box-shadow: 0 20px 40px rgba(0,0,0,.15); }
.modal--wide { width: 540px; }
.modal-head { display: flex; align-items: center; justify-content: space-between; padding: 18px 22px 14px; border-bottom: 1px solid #f1f5f9; }
.modal-title { font-size: 15px; font-weight: 700; color: #111827; }
.modal-close { background: none; border: none; font-size: 16px; color: #9ca3af; cursor: pointer; padding: 4px; }
.modal-close:hover { color: #374151; }
.modal-body { padding: 18px 22px; }
.modal-hint { font-size: 12.5px; color: #6b7280; background: #f8fafc; padding: 10px 14px; border-radius: 8px; margin-bottom: 14px; }
.modal-info { font-size: 13px; color: #374151; margin-bottom: 14px; line-height: 1.5; }
.modal-foot { display: flex; gap: 10px; justify-content: flex-end; padding: 14px 22px 18px; border-top: 1px solid #f1f5f9; }

.field { margin-bottom: 14px; }
.label { display: block; font-size: 12.5px; font-weight: 600; color: #374151; margin-bottom: 6px; }
.req   { color: #ef4444; }
.hint  { font-weight: 400; color: #9ca3af; margin-left: 4px; }
.select { width: 100%; padding: 9px 12px; border: 1.5px solid #e5e7eb; border-radius: 8px; font-size: 13px; color: #111827; background: #fff; }
.select:focus { outline: none; border-color: #16a34a; }
.textarea { width: 100%; padding: 9px 12px; border: 1.5px solid #e5e7eb; border-radius: 8px; font-size: 13px; color: #111827; resize: vertical; font-family: inherit; box-sizing: border-box; }
.textarea:focus { outline: none; border-color: #16a34a; }
.file-input { font-size: 12.5px; color: #374151; }

.durasi-grid { display: flex; gap: 8px; flex-wrap: wrap; }
.durasi-btn { padding: 7px 16px; border-radius: 8px; border: 1.5px solid #e5e7eb; background: #f9fafb; font-size: 13px; font-weight: 500; color: #374151; cursor: pointer; }
.durasi-btn:hover { border-color: #16a34a; color: #15803d; }
.durasi-btn--active { border-color: #16a34a; background: #f0fdf4; color: #15803d; font-weight: 700; }

.alert { padding: 10px 14px; border-radius: 8px; font-size: 12.5px; margin-top: 10px; }
.alert--error { background: #fef2f2; color: #b91c1c; border: 1px solid #fca5a5; }
.alert--ok    { background: #f0fdf4; color: #15803d; border: 1px solid #86efac; }

.btn-cancel { padding: 9px 18px; border: 1.5px solid #e5e7eb; background: #fff; border-radius: 9px; font-size: 13px; font-weight: 500; color: #374151; cursor: pointer; }
.btn-ok     { padding: 9px 22px; background: #16a34a; color: #fff; border: none; border-radius: 9px; font-size: 13px; font-weight: 600; cursor: pointer; display: flex; align-items: center; gap: 7px; }
.btn-ok:hover:not(:disabled) { background: #15803d; }
.btn-ok:disabled { opacity: 0.55; cursor: not-allowed; }
.btn-tolak-ok { padding: 9px 22px; background: #dc2626; color: #fff; border: none; border-radius: 9px; font-size: 13px; font-weight: 600; cursor: pointer; display: flex; align-items: center; gap: 7px; }
.btn-tolak-ok:hover:not(:disabled) { background: #b91c1c; }
.btn-tolak-ok:disabled { opacity: 0.55; cursor: not-allowed; }

.spinner { width: 32px; height: 32px; border: 3px solid #e5e7eb; border-top-color: #16a34a; border-radius: 50%; animation: spin .7s linear infinite; }
.spinner-sm { width: 15px; height: 15px; border: 2px solid rgba(255,255,255,0.4); border-top-color: #fff; border-radius: 50%; animation: spin .7s linear infinite; display: inline-block; }
@keyframes spin { to { transform: rotate(360deg); } }
</style>
