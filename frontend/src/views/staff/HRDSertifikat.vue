<template>
  <div>
    <!-- Loading -->
    <div v-if="loading" class="card state-box">
      <div class="spinner"></div>
      <span>Memuat daftar peserta…</span>
    </div>

    <!-- Empty -->
    <div v-else-if="!list.length" class="card state-box state-box--empty">
      <svg width="48" height="48" viewBox="0 0 24 24" fill="none">
        <circle cx="12" cy="8" r="6" stroke="#d1d5db" stroke-width="1.5"/>
        <path d="M15.477 12.89L17 22l-5-3-5 3 1.523-9.11" stroke="#d1d5db" stroke-width="1.5"/>
      </svg>
      <div class="empty-title">Belum Ada Peserta Siap Sertifikat</div>
      <div class="empty-sub">Peserta dengan nilai akhir akan muncul di sini untuk diterbitkan sertifikatnya.</div>
    </div>

    <template v-else>
      <!-- ── Card utama ── -->
      <div class="card">

        <!-- Card header: info + stats -->
        <div class="card-header">
          <div class="card-header__left">
            <div class="card-icon">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none">
                <circle cx="12" cy="8" r="6" stroke="currentColor" stroke-width="2"/>
                <path d="M15.477 12.89L17 22l-5-3-5 3 1.523-9.11" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              </svg>
            </div>
            <div>
              <h3 class="card-title">Sertifikat Peserta</h3>
              <p class="card-sub">Upload sertifikat PDF — otomatis tersedia di portal peserta setelah di-upload.</p>
            </div>
          </div>
          <div class="card-stats">
            <div class="stat-chip">
              <span class="stat-chip__num">{{ list.length }}</span>
              <span class="stat-chip__lbl">Total</span>
            </div>
            <div class="stat-chip stat-chip--done">
              <span class="stat-chip__num">{{ sudahCount }}</span>
              <span class="stat-chip__lbl">Diterbitkan</span>
            </div>
            <div class="stat-chip stat-chip--pending">
              <span class="stat-chip__num">{{ belumCount }}</span>
              <span class="stat-chip__lbl">Belum</span>
            </div>
          </div>
        </div>

        <!-- Filter + Search -->
        <div class="filter-bar">
          <button
            v-for="f in filters" :key="f.key"
            :class="['filter-btn', activeFilter === f.key ? 'filter-btn--active' : '']"
            @click="activeFilter = f.key">
            {{ f.label }}
          </button>
          <div class="filter-spacer"></div>
          <input v-model="search" class="search-input" placeholder="Cari nama peserta…" />
        </div>

        <!-- Cards list -->
        <div class="cert-list">
          <div v-for="item in filteredList" :key="String(item.pelaksanaan_id)" class="cert-card">
            <!-- Kiri: info peserta -->
            <div class="cert-card__info">
              <div class="cert-card__avatar">{{ initials(item.nama_lengkap) }}</div>
              <div class="cert-card__detail">
                <div class="cert-card__nama">{{ item.nama_lengkap }}</div>
                <div class="cert-card__meta">{{ item.email }}</div>
                <div class="cert-card__tags">
                  <span v-if="item.institusi" class="tag tag--institusi">{{ item.institusi }}</span>
                  <span class="tag tag--gray">{{ item.divisi || '—' }}</span>
                  <span class="tag tag--gray">{{ fmtPeriode(item.tanggal_mulai, item.tanggal_selesai) }}</span>
                  <span class="tag tag--nilai">Nilai: {{ item.nilai?.toFixed(1) ?? '—' }}</span>
                </div>
              </div>
            </div>

            <!-- Kanan: status + aksi -->
            <div class="cert-card__action">
              <div v-if="item.sertifikat_generated" class="status-done">
                <svg width="13" height="13" viewBox="0 0 24 24" fill="none">
                  <path d="M5 13l4 4L19 7" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"/>
                </svg>
                Sudah Diterbitkan
              </div>
              <div v-else class="status-pending">Belum Diterbitkan</div>

              <div v-if="item.sertifikat_generated" class="cert-card__date">
                {{ fmtDate(item.sertifikat_generated_at) }}
              </div>

              <!-- Baris tombol aksi -->
              <div class="cert-card__btns">
                <!-- Preview & Unduh — hanya jika sudah ada PDF -->
                <template v-if="item.sertifikat_generated">
                  <button class="btn-action btn-action--preview" @click="previewSertifikat(item)" title="Preview PDF">
                    <svg width="13" height="13" viewBox="0 0 24 24" fill="none">
                      <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z" stroke="currentColor" stroke-width="2"/>
                      <circle cx="12" cy="12" r="3" stroke="currentColor" stroke-width="2"/>
                    </svg>
                    Preview
                  </button>
                  <button class="btn-action btn-action--dl"
                    :disabled="downloading === String(item.pelaksanaan_id)"
                    @click="downloadSertifikat(item)" title="Unduh PDF">
                    <div v-if="downloading === String(item.pelaksanaan_id)" class="btn-spinner btn-spinner--dark"></div>
                    <svg v-else width="13" height="13" viewBox="0 0 24 24" fill="none">
                      <path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
                      <polyline points="7 10 12 15 17 10" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                      <line x1="12" y1="15" x2="12" y2="3" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
                    </svg>
                    Unduh
                  </button>
                </template>

                <!-- Upload / Ganti -->
                <label :for="`upload-${item.pelaksanaan_id}`"
                  class="btn-action"
                  :class="item.sertifikat_generated ? 'btn-action--replace' : 'btn-action--upload'"
                  :title="item.sertifikat_generated ? 'Ganti PDF sertifikat' : 'Upload sertifikat PDF'">
                  <div v-if="uploading === String(item.pelaksanaan_id)" class="btn-spinner"></div>
                  <svg v-else width="13" height="13" viewBox="0 0 24 24" fill="none">
                    <path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
                    <polyline points="17 8 12 3 7 8" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                    <line x1="12" y1="3" x2="12" y2="15" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
                  </svg>
                  {{ item.sertifikat_generated ? 'Ganti PDF' : 'Upload' }}
                </label>
                <input
                  :id="`upload-${item.pelaksanaan_id}`"
                  type="file" accept=".pdf" class="file-hidden"
                  @change="(e) => handleUpload(e, item)" />
              </div>
            </div>
          </div>
        </div>

      </div><!-- end card -->
    </template>

    <!-- PDF Preview Modal -->
    <Teleport to="body">
      <div v-if="pdfModal.show" class="pdf-backdrop" @click.self="closePdfModal">
        <div class="pdf-modal">
          <div class="pdf-modal__head">
            <span>{{ pdfModal.nama }} — Sertifikat Magang</span>
            <div class="pdf-modal__acts">
              <button class="pdf-dl-btn"
                :disabled="downloading === pdfModal.id"
                @click="downloadSertifikatById(pdfModal.id!, pdfModal.nama!)">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none">
                  <path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
                  <polyline points="7 10 12 15 17 10" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                  <line x1="12" y1="15" x2="12" y2="3" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
                </svg>
                Unduh
              </button>
              <button class="pdf-close-btn" @click="closePdfModal">✕</button>
            </div>
          </div>
          <div class="pdf-modal__body">
            <div v-if="!pdfModal.src" class="pdf-loading">
              <div class="spinner"></div><span>Memuat PDF…</span>
            </div>
            <iframe v-else :src="pdfModal.src" class="pdf-iframe" />
          </div>
        </div>
      </div>
    </Teleport>

    <!-- Toast -->
    <Teleport to="body">
      <div v-if="toast.show" :class="['toast', `toast--${toast.type}`]">
        <svg v-if="toast.type === 'success'" width="16" height="16" viewBox="0 0 24 24" fill="none">
          <path d="M5 13l4 4L19 7" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"/>
        </svg>
        <svg v-else width="16" height="16" viewBox="0 0 24 24" fill="none">
          <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2"/>
          <line x1="12" y1="8" x2="12" y2="12" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
          <line x1="12" y1="16" x2="12.01" y2="16" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
        </svg>
        {{ toast.message }}
      </div>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import api from "@/lib/api";

interface SertifikatItem {
  pelaksanaan_id: string;
  status: string;
  nama_lengkap: string;
  email: string;
  institusi: string;
  divisi: string;
  pembimbing: string;
  tanggal_mulai: string;
  tanggal_selesai: string;
  nilai: number | null;
  sertifikat_generated: boolean;
  sertifikat_generated_at: string | null;
}

const loading    = ref(true);
const list       = ref<SertifikatItem[]>([]);
const search     = ref("");
const activeFilter = ref<"semua"|"belum"|"sudah">("semua");
const uploading  = ref<string | null>(null);
const downloading = ref<string | null>(null);

const filters = [
  { key: "semua",  label: "Semua" },
  { key: "belum",  label: "Belum Diterbitkan" },
  { key: "sudah",  label: "Sudah Diterbitkan" },
] as const;

const toast = ref({ show: false, type: "success", message: "" });
const pdfModal = ref<{ show: boolean; src: string; id: string | null; nama: string | null }>({
  show: false, src: "", id: null, nama: null,
});

const sudahCount = computed(() => list.value.filter(i => i.sertifikat_generated).length);
const belumCount = computed(() => list.value.filter(i => !i.sertifikat_generated).length);

const filteredList = computed(() => {
  let res = list.value;
  if (activeFilter.value === "belum") res = res.filter(i => !i.sertifikat_generated);
  if (activeFilter.value === "sudah") res = res.filter(i => i.sertifikat_generated);
  if (search.value.trim()) {
    const q = search.value.toLowerCase();
    res = res.filter(i =>
      i.nama_lengkap.toLowerCase().includes(q) ||
      i.email.toLowerCase().includes(q) ||
      (i.institusi ?? "").toLowerCase().includes(q)
    );
  }
  return res;
});

function initials(nama: string) {
  return nama.split(" ").slice(0, 2).map(w => w[0]).join("").toUpperCase();
}
function fmtDate(d: string | null) {
  if (!d) return "—";
  return new Date(d).toLocaleDateString("id-ID", { day: "2-digit", month: "short", year: "numeric" });
}
function fmtPeriode(mulai: string, selesai: string) {
  const fmt = (s: string) =>
    new Date(s).toLocaleDateString("id-ID", { day: "2-digit", month: "short", year: "numeric" });
  return `${fmt(mulai)} – ${fmt(selesai)}`;
}
function showToast(message: string, type: "success"|"error" = "success") {
  toast.value = { show: true, type, message };
  setTimeout(() => { toast.value.show = false; }, 4000);
}

async function loadList() {
  loading.value = true;
  try {
    const res = await api.get("/api/sertifikat");
    list.value = res.data?.data ?? [];
  } catch {
    list.value = [];
  } finally {
    loading.value = false;
  }
}

async function previewSertifikat(item: SertifikatItem) {
  pdfModal.value = { show: true, src: "", id: String(item.pelaksanaan_id), nama: item.nama_lengkap };
  try {
    const r = await api.get(`/api/pelaksanaan/${item.pelaksanaan_id}/sertifikat/download`, { responseType: "blob" });
    pdfModal.value.src = URL.createObjectURL(new Blob([r.data], { type: "application/pdf" }));
  } catch {
    showToast("Gagal memuat preview PDF", "error");
    pdfModal.value.show = false;
  }
}

function closePdfModal() {
  if (pdfModal.value.src) URL.revokeObjectURL(pdfModal.value.src);
  pdfModal.value = { show: false, src: "", id: null, nama: null };
}

async function downloadSertifikat(item: SertifikatItem) {
  await downloadSertifikatById(String(item.pelaksanaan_id), item.nama_lengkap);
}

async function downloadSertifikatById(id: string, nama: string) {
  downloading.value = id;
  try {
    const r = await api.get(`/api/pelaksanaan/${id}/sertifikat/download`, { responseType: "blob" });
    const url = URL.createObjectURL(new Blob([r.data], { type: "application/pdf" }));
    const a   = document.createElement("a");
    a.href     = url;
    a.download = `sertifikat_${nama.replace(/\s+/g, "_")}.pdf`;
    a.click();
    URL.revokeObjectURL(url);
  } catch {
    showToast("Gagal mengunduh sertifikat", "error");
  } finally {
    downloading.value = null;
  }
}

async function handleUpload(event: Event, item: SertifikatItem) {
  const input = event.target as HTMLInputElement;
  const file  = input.files?.[0];
  if (!file) return;
  if (file.type !== "application/pdf") {
    showToast("Hanya file PDF yang diperbolehkan", "error");
    input.value = "";
    return;
  }
  uploading.value = String(item.pelaksanaan_id);
  const formData = new FormData();
  formData.append("file", file);
  try {
    await api.post(`/api/pelaksanaan/${item.pelaksanaan_id}/sertifikat`, formData, {
      headers: { "Content-Type": "multipart/form-data" },
    });
    showToast(`Sertifikat ${item.nama_lengkap} berhasil diterbitkan`);
    await loadList();
  } catch (e: any) {
    showToast(e.response?.data?.message ?? "Gagal mengunggah sertifikat", "error");
  } finally {
    uploading.value = null;
    input.value = "";
  }
}

onMounted(loadList);
</script>

<style scoped>
/* ── State ── */
.state-box { display:flex; flex-direction:column; align-items:center; gap:14px; padding:52px 24px; text-align:center; color:#6b7280; font-size:13px; }
.empty-title { font-size:16px; font-weight:700; color:#374151; }
.empty-sub   { font-size:12.5px; color:#9ca3af; line-height:1.7; max-width:360px; }
.spinner { width:26px; height:26px; border:3px solid #e5e7eb; border-top-color:#48AF4A; border-radius:50%; animation:spin .7s linear infinite; flex-shrink:0; }
@keyframes spin { to { transform:rotate(360deg); } }

/* ── Card wrapper ── */
.card { background:#fff; border:1px solid #e9f5e9; border-radius:14px; box-shadow:0 1px 3px rgba(13,40,24,.05); overflow:hidden; }

/* ── Card header ── */
.card-header {
  display:flex; align-items:center; justify-content:space-between;
  gap:16px; flex-wrap:wrap;
  padding:16px 20px; border-bottom:1px solid #f0faf0;
}
.card-header__left { display:flex; align-items:center; gap:12px; flex:1; min-width:0; }
.card-icon {
  width:36px; height:36px; border-radius:9px; background:#f0fdf4;
  color:#16a34a; display:flex; align-items:center; justify-content:center; flex-shrink:0;
}
.card-title { font-size:13.5px; font-weight:700; color:#111827; margin:0 0 2px; }
.card-sub   { font-size:11.5px; color:#9ca3af; margin:0; }
.card-stats { display:flex; gap:8px; flex-shrink:0; }

/* ── Stat chips ── */
.stat-chip { display:flex; flex-direction:column; align-items:center; padding:8px 16px; background:#f9fafb; border:1px solid #e5e7eb; border-radius:8px; min-width:72px; }
.stat-chip--done    { border-color:#86efac; background:#f0fdf4; }
.stat-chip--pending { border-color:#fca5a5; background:#fef2f2; }
.stat-chip__num { font-size:18px; font-weight:800; color:#111827; line-height:1.2; }
.stat-chip--done .stat-chip__num    { color:#16a34a; }
.stat-chip--pending .stat-chip__num { color:#dc2626; }
.stat-chip__lbl { font-size:10.5px; color:#9ca3af; font-weight:500; }

/* ── Filter ── */
.filter-bar { display:flex; gap:8px; align-items:center; flex-wrap:wrap; padding:12px 20px; background:#f9fafb; border-bottom:1px solid #f0faf0; }
.filter-btn { padding:6px 14px; border-radius:20px; border:1.5px solid #e5e7eb; background:#fff; font-size:12.5px; font-weight:500; color:#6b7280; cursor:pointer; font-family:inherit; transition:all .15s; }
.filter-btn:hover { background:#f9fafb; }
.filter-btn--active { background:#48AF4A; color:#fff; border-color:#48AF4A; font-weight:600; }
.filter-spacer { flex:1; }
.search-input { padding:7px 12px; border:1.5px solid #e5e7eb; border-radius:8px; font-size:13px; font-family:inherit; outline:none; width:220px; color:#111827; background:#fff; }
.search-input:focus { border-color:#48AF4A; }

/* ── Card list ── */
.cert-list { display:flex; flex-direction:column; gap:0; padding:12px 16px; gap:8px; }
.cert-card {
  display:flex; align-items:center; justify-content:space-between; gap:16px;
  background:#fff; border:1px solid #e5e7eb; border-radius:10px;
  padding:14px 18px; transition:border-color .15s; flex-wrap:wrap;
}
.cert-card:hover { border-color:#86efac; }

.cert-card__info { display:flex; align-items:flex-start; gap:14px; flex:1; min-width:0; }
.cert-card__avatar { width:42px; height:42px; border-radius:50%; background:#0d2818; color:#fff; font-size:14px; font-weight:700; display:flex; align-items:center; justify-content:center; flex-shrink:0; }
.cert-card__detail { min-width:0; }
.cert-card__nama { font-size:14px; font-weight:700; color:#111827; margin-bottom:2px; }
.cert-card__meta { font-size:12px; color:#9ca3af; margin-bottom:6px; }
.cert-card__tags { display:flex; gap:6px; flex-wrap:wrap; }
.tag { font-size:11px; font-weight:600; padding:2px 8px; border-radius:4px; }
.tag--gray     { background:#f3f4f6; color:#6b7280; }
.tag--nilai    { background:#f0fdf4; color:#16a34a; }
.tag--institusi { background:#f0fdf4; color:#0d2818; }

.cert-card__action { display:flex; flex-direction:column; align-items:flex-end; gap:6px; flex-shrink:0; }
.status-done    { display:flex; align-items:center; gap:5px; font-size:12px; font-weight:600; color:#16a34a; }
.status-pending { font-size:12px; font-weight:600; color:#9ca3af; }
.cert-card__date { font-size:11px; color:#9ca3af; }
.cert-card__btns { display:flex; gap:6px; flex-wrap:wrap; align-items:center; }

/* ── Tombol aksi ── */
.file-hidden { display:none; }
.btn-action {
  display:inline-flex; align-items:center; gap:5px; cursor:pointer;
  padding:7px 12px; border-radius:7px; font-size:12px; font-weight:600;
  border:1.5px solid transparent; font-family:inherit;
  transition:background .15s, border-color .15s; user-select:none;
}
.btn-action:disabled { opacity:.5; cursor:not-allowed; }

.btn-action--preview  { background:#f0fdf4; color:#0d2818; border-color:#bbf7d0; }
.btn-action--preview:hover { background:#dcfce7; }

.btn-action--dl  { background:#f0fdf4; color:#16a34a; border-color:#86efac; }
.btn-action--dl:hover { background:#dcfce7; }

.btn-action--upload  { background:#48AF4A; color:#fff; border-color:#48AF4A; }
.btn-action--upload:hover { background:#3d9e3f; }

.btn-action--replace { background:#f9fafb; color:#374151; border-color:#e5e7eb; }
.btn-action--replace:hover { background:#f3f4f6; }

.btn-spinner { width:12px; height:12px; border:2px solid rgba(255,255,255,.35); border-top-color:#fff; border-radius:50%; animation:spin .7s linear infinite; }
.btn-spinner--dark { border:2px solid #e5e7eb; border-top-color:#6b7280; }

/* ── PDF Modal ── */
.pdf-backdrop { position:fixed; inset:0; background:rgba(0,0,0,.55); z-index:9999; display:flex; align-items:center; justify-content:center; padding:16px; }
.pdf-modal { background:#fff; border-radius:14px; display:flex; flex-direction:column; width:100%; max-width:860px; height:90vh; overflow:hidden; box-shadow:0 24px 64px rgba(0,0,0,.25); }
.pdf-modal__head { display:flex; align-items:center; justify-content:space-between; padding:14px 18px; border-bottom:1px solid #e5e7eb; font-weight:700; color:#111827; font-size:13.5px; }
.pdf-modal__acts { display:flex; gap:8px; align-items:center; }
.pdf-dl-btn { display:flex; align-items:center; gap:6px; background:#48AF4A; color:#fff; border:none; border-radius:7px; padding:6px 14px; font-size:12px; font-weight:600; cursor:pointer; font-family:inherit; }
.pdf-dl-btn:hover { background:#3d9e3f; }
.pdf-dl-btn:disabled { opacity:.5; cursor:not-allowed; }
.pdf-close-btn { background:#f9fafb; border:1.5px solid #e5e7eb; border-radius:7px; width:30px; height:30px; cursor:pointer; font-size:14px; color:#6b7280; display:flex; align-items:center; justify-content:center; }
.pdf-close-btn:hover { background:#f3f4f6; color:#374151; }
.pdf-modal__body { flex:1; display:flex; overflow:hidden; }
.pdf-loading { display:flex; flex-direction:column; align-items:center; justify-content:center; gap:10px; flex:1; color:#6b7280; font-size:13px; }
.pdf-iframe { width:100%; height:100%; border:none; }

/* ── Toast ── */
.toast { position:fixed; bottom:24px; right:24px; display:flex; align-items:center; gap:10px; padding:12px 18px; border-radius:10px; font-size:13px; font-weight:600; box-shadow:0 8px 24px rgba(0,0,0,.15); z-index:9999; animation:slideUp .25s ease; }
.toast--success { background:#0d2818; color:#fff; }
.toast--error   { background:#dc2626; color:#fff; }
@keyframes slideUp { from { opacity:0; transform:translateY(10px); } to { opacity:1; transform:translateY(0); } }
</style>
