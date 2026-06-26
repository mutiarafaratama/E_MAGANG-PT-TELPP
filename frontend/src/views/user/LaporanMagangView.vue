<template>
  <div>
    <!-- Loading -->
    <div v-if="loading" class="state-box">
      <div class="spinner"></div>
    </div>

    <!-- Belum waktunya -->
    <div v-else-if="!canUpload && laporanList.length === 0" class="state-box state-box--locked">
      <div class="locked-icon">
        <svg width="44" height="44" viewBox="0 0 24 24" fill="none">
          <rect x="3" y="11" width="18" height="11" rx="2" stroke="#f59e0b" stroke-width="2"/>
          <path d="M7 11V7a5 5 0 0110 0v4" stroke="#f59e0b" stroke-width="2" stroke-linecap="round"/>
          <circle cx="12" cy="16" r="1.5" fill="#f59e0b"/>
        </svg>
      </div>
      <div class="locked-title">Pengumpulan Laporan Belum Dibuka</div>
      <div class="locked-sub">
        Fitur upload laporan tersedia saat status magang Anda <strong>"Aktif"</strong>
        atau <strong>"Upload Laporan"</strong>.
        Jadwal magang perlu dikonfirmasi oleh HRD terlebih dahulu.
      </div>
      <div class="status-chip" :class="`chip--${statusCls}`">
        Status saat ini: <strong>{{ statusLabel }}</strong>
      </div>
      <button class="btn-refresh-status" @click="emit('refresh')">
        <svg width="13" height="13" viewBox="0 0 24 24" fill="none"><path d="M23 4v6h-6" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/><path d="M1 20v-6h6" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/><path d="M3.51 9a9 9 0 0114.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0020.49 15" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>
        Cek Status Terbaru
      </button>
    </div>

    <template v-else>
      <!-- Status review banner (tampil jika ada laporan) -->
      <div v-if="latest" class="review-banner" :class="`review-banner--${latest.status}`">
        <div class="review-banner__icon">
          <svg v-if="latest.status === 'menunggu_review'" width="22" height="22" viewBox="0 0 24 24" fill="none">
            <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2"/>
            <polyline points="12 6 12 12 16 14" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
          </svg>
          <svg v-else-if="latest.status === 'disetujui'" width="22" height="22" viewBox="0 0 24 24" fill="none">
            <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2"/>
            <path d="M8 12l3 3 5-5" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
          <svg v-else width="22" height="22" viewBox="0 0 24 24" fill="none">
            <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2"/>
            <line x1="12" y1="8" x2="12" y2="12" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
            <circle cx="12" cy="16" r="1" fill="currentColor"/>
          </svg>
        </div>
        <div class="review-banner__body">
          <div class="review-banner__title">{{ reviewBannerTitle }}</div>
          <div class="review-banner__sub">{{ reviewBannerSub }}</div>
          <div v-if="latest.status === 'revisi' && latest.catatan_hrd" class="review-catatan">
            <svg width="12" height="12" viewBox="0 0 24 24" fill="none"><path d="M21 15a2 2 0 01-2 2H7l-4 4V5a2 2 0 012-2h14a2 2 0 012 2z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>
            Catatan HRD: <em>{{ latest.catatan_hrd }}</em>
          </div>
        </div>
        <span class="review-versi">v{{ latest.versi }}</span>
      </div>

      <!-- Riwayat semua upload -->
      <div v-if="laporanList.length > 0" class="history-card">
        <div class="history-card__header">
          <span class="history-card__title">Riwayat Pengiriman</span>
          <span class="history-count">{{ laporanList.length }} file</span>
        </div>
        <div class="history-list">
          <div v-for="lap in laporanList" :key="lap.id" class="history-item">
            <div class="history-item__icon">
              <svg width="15" height="15" viewBox="0 0 24 24" fill="none">
                <path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z" stroke="currentColor" stroke-width="2"/>
                <polyline points="14 2 14 8 20 8" stroke="currentColor" stroke-width="2"/>
              </svg>
            </div>
            <div class="history-item__info">
              <div class="history-item__name">{{ lap.nama_file }}</div>
              <div class="history-item__meta">
                <span class="versi-badge">v{{ lap.versi }}</span>
                <span>{{ fmtSize(lap.ukuran_bytes) }}</span>
                <span>{{ fmtDate(lap.diupload_at) }}</span>
              </div>
            </div>
            <div class="history-item__right">
              <span class="status-chip-sm" :class="`chip-sm--${lap.status}`">{{ statusLabelLap(lap.status) }}</span>
              <!-- Preview -->
              <button class="btn-action" title="Preview" @click="openPreview(lap)">
                <svg width="13" height="13" viewBox="0 0 24 24" fill="none">
                  <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z" stroke="currentColor" stroke-width="2"/>
                  <circle cx="12" cy="12" r="3" stroke="currentColor" stroke-width="2"/>
                </svg>
              </button>
              <!-- Download -->
              <a :href="`/api/laporan/${lap.id}/download`" target="_blank" class="btn-action" title="Unduh">
                <svg width="13" height="13" viewBox="0 0 24 24" fill="none">
                  <path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
                  <polyline points="7 10 12 15 17 10" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                  <line x1="12" y1="15" x2="12" y2="3" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
                </svg>
              </a>
            </div>
          </div>
        </div>
      </div>

      <!-- Form Upload -->
      <div v-if="canUpload" class="upload-card">
        <div class="upload-card__topbar">
          <span class="upload-card__title">
            {{ laporanList.length === 0 ? 'Upload Laporan Akhir Magang' : 'Upload Ulang / Revisi Laporan' }}
          </span>
          <button class="btn-panduan" @click="showStepInfo = true" title="Panduan upload laporan">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none">
              <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2"/>
              <line x1="12" y1="16" x2="12" y2="12" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
              <circle cx="12" cy="8" r="1" fill="currentColor"/>
            </svg>
            Panduan
          </button>
        </div>
        <div class="upload-card__body">
        <div class="upload-card__sub">Format: PDF, DOC, DOCX. Maksimal 20 MB.</div>

        <!-- Drop zone -->
        <label class="drop-zone" :class="{ 'drop-zone--dragover': dragging, 'drop-zone--has-file': selectedFile }">
          <input type="file" accept=".pdf,.doc,.docx" @change="onFileChange"
            @dragenter.prevent="dragging=true" @dragleave="dragging=false"
            style="display:none" ref="fileInput"/>
          <template v-if="!selectedFile">
            <svg width="36" height="36" viewBox="0 0 24 24" fill="none">
              <path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4" stroke="#9ca3af" stroke-width="1.8" stroke-linecap="round"/>
              <polyline points="17 8 12 3 7 8" stroke="#9ca3af" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"/>
              <line x1="12" y1="3" x2="12" y2="15" stroke="#9ca3af" stroke-width="1.8" stroke-linecap="round"/>
            </svg>
            <div class="dz-label">Seret file ke sini atau <span class="dz-link" @click.prevent="fileInput?.click()">klik untuk pilih</span></div>
            <div class="dz-sub">PDF, DOC, DOCX — maks 20 MB</div>
          </template>
          <template v-else>
            <svg width="28" height="28" viewBox="0 0 24 24" fill="none">
              <path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z" stroke="#16a34a" stroke-width="2" stroke-linejoin="round"/>
              <polyline points="14 2 14 8 20 8" stroke="#16a34a" stroke-width="2"/>
            </svg>
            <div class="dz-file-name">{{ selectedFile.name }}</div>
            <div class="dz-file-size">{{ fmtSize(selectedFile.size) }}</div>
            <button class="dz-remove" @click.prevent="selectedFile = null">Ganti file</button>
          </template>
        </label>

        <div v-if="errMsg" class="upload-err">{{ errMsg }}</div>

        <button class="btn-submit" :disabled="!selectedFile || uploading" @click="doUpload">
          <svg v-if="!uploading" width="15" height="15" viewBox="0 0 24 24" fill="none">
            <path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
            <polyline points="17 8 12 3 7 8" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            <line x1="12" y1="3" x2="12" y2="15" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
          </svg>
          <div v-else class="btn-spinner"></div>
          {{ uploading ? "Mengunggah…" : (laporanList.length === 0 ? "Kirim Laporan" : "Kirim Ulang Laporan") }}
        </button>

        <div v-if="uploadOk" class="upload-ok">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none"><path d="M20 6L9 17l-5-5" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"/></svg>
          Laporan berhasil dikirim! Menunggu review dari HRD.
        </div>
        </div><!-- end upload-card__body -->
      </div>

      <!-- Laporan sudah disetujui -->
      <div v-else-if="latest && latest.status === 'disetujui'" class="info-box info-box--green">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none"><circle cx="12" cy="12" r="10" fill="#16a34a"/><path d="M8 12l3 3 5-5" stroke="#fff" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>
        Laporan Anda telah <strong>disetujui oleh HRD</strong>. Form penilaian akan segera dibuka oleh HRD.
      </div>
    </template>

    <!-- ── Modal Panduan Upload ── -->
    <Teleport to="body">
      <div v-if="showStepInfo" class="modal-backdrop" @click.self="showStepInfo = false">
        <div class="modal modal--info">
          <div class="modal-header">
            <div class="modal-header__left">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" style="color:#16a34a;flex-shrink:0">
                <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2"/>
                <line x1="12" y1="16" x2="12" y2="12" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
                <circle cx="12" cy="8" r="1" fill="currentColor"/>
              </svg>
              <div class="modal-title">Panduan Upload Laporan</div>
            </div>
            <button class="modal-close" @click="showStepInfo = false">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none"><path d="M18 6L6 18M6 6l12 12" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>
            </button>
          </div>
          <div class="modal-body">
            <div class="step-list">
              <div v-for="(s, i) in steps" :key="i" class="step-item">
                <div class="step-item__num">{{ i + 1 }}</div>
                <div class="step-item__body">
                  <div class="step-item__title">{{ s.short }}</div>
                  <ul class="step-item__details">
                    <li v-for="(d, di) in s.details" :key="di">{{ d }}</li>
                  </ul>
                </div>
              </div>
            </div>
          </div>
          <div class="modal-footer modal-footer--single">
            <button class="btn-ok" @click="showStepInfo = false">Mengerti</button>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- ── Modal Preview Laporan ── -->
    <Teleport to="body">
      <div v-if="previewState.show" class="preview-overlay" @click.self="closePreview">
        <div class="preview-box">
          <div class="preview-header">
            <span class="preview-filename">{{ previewState.nama }}</span>
            <div class="preview-header-actions">
              <a v-if="previewState.lapId" :href="`/api/laporan/${previewState.lapId}/download`" target="_blank" class="preview-btn-dl" title="Download">
                <svg width="15" height="15" viewBox="0 0 24 24" fill="none"><path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4" stroke="currentColor" stroke-width="2" stroke-linecap="round"/><polyline points="7 10 12 15 17 10" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/><line x1="12" y1="15" x2="12" y2="3" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>
              </a>
              <button class="preview-close" @click="closePreview">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none"><path d="M18 6L6 18M6 6l12 12" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>
              </button>
            </div>
          </div>
          <div class="preview-body">
            <div v-if="previewState.loading" class="preview-loading"><div class="spinner"></div></div>
            <div v-else-if="previewState.error" class="preview-error">
              <svg width="32" height="32" viewBox="0 0 24 24" fill="none"><circle cx="12" cy="12" r="10" stroke="#9ca3af" stroke-width="1.5"/><line x1="15" y1="9" x2="9" y2="15" stroke="#9ca3af" stroke-width="1.5" stroke-linecap="round"/><line x1="9" y1="9" x2="15" y2="15" stroke="#9ca3af" stroke-width="1.5" stroke-linecap="round"/></svg>
              <p>Tidak dapat menampilkan file ini.</p>
            </div>
            <iframe v-else-if="previewState.type === 'pdf'" :src="previewState.url" class="preview-iframe"/>
            <div v-else class="preview-error">
              <svg width="40" height="40" viewBox="0 0 24 24" fill="none"><path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z" stroke="#9ca3af" stroke-width="1.5"/><polyline points="14 2 14 8 20 8" stroke="#9ca3af" stroke-width="1.5"/></svg>
              <p>Preview hanya tersedia untuk file PDF.</p>
              <a v-if="previewState.lapId" :href="`/api/laporan/${previewState.lapId}/download`" target="_blank" class="btn-green-sm">Download File</a>
            </div>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import api from "@/lib/api";

const props = defineProps<{ pelaksanaan: any }>();
const emit  = defineEmits<{ (e: 'refresh'): void }>();

const loading      = ref(true);
const uploading    = ref(false);
const dragging     = ref(false);
const errMsg       = ref("");
const uploadOk     = ref(false);
const laporanList  = ref<any[]>([]);
const selectedFile = ref<File | null>(null);
const fileInput    = ref<HTMLInputElement | null>(null);

const canUpload = computed(() =>
  ["aktif", "upload_laporan"].includes(props.pelaksanaan?.status ?? "")
);

const latest = computed(() => laporanList.value[0] ?? null);

const statusLabel = computed(() => {
  const m: Record<string, string> = {
    menunggu_mulai: "Menunggu Dimulai", aktif: "Aktif",
    upload_laporan: "Upload Laporan",   penilaian: "Sedang Dinilai",
    selesai: "Selesai",
  };
  return m[props.pelaksanaan?.status ?? ""] ?? "–";
});

const statusCls = computed(() => {
  const s = props.pelaksanaan?.status ?? "";
  if (s === "aktif")          return "green";
  if (s === "upload_laporan") return "yellow";
  if (s === "penilaian" || s === "selesai") return "blue";
  return "gray";
});

const reviewBannerTitle = computed(() => {
  if (!latest.value) return "";
  const m: Record<string, string> = {
    menunggu_review: "Laporan sedang ditinjau oleh HRD",
    disetujui:       "Laporan disetujui oleh HRD",
    revisi:          "HRD meminta revisi laporan Anda",
  };
  return m[latest.value.status] ?? "";
});

const reviewBannerSub = computed(() => {
  if (!latest.value) return "";
  const m: Record<string, string> = {
    menunggu_review: "Mohon tunggu, HRD akan meninjau dan memberikan keputusan segera.",
    disetujui:       "Selamat! Laporan Anda telah diterima. Form penilaian akan segera dibuka.",
    revisi:          "Silakan upload ulang laporan sesuai catatan HRD di bawah ini.",
  };
  return m[latest.value.status] ?? "";
});

function statusLabelLap(s: string) {
  return ({ menunggu_review: "Menunggu Review", disetujui: "Disetujui", revisi: "Perlu Revisi" } as Record<string,string>)[s] ?? s;
}

function fmtDate(iso: string) {
  if (!iso) return "";
  return new Date(iso).toLocaleDateString("id-ID", { day: "numeric", month: "short", year: "numeric" });
}

function fmtSize(bytes: number | null) {
  if (!bytes) return "";
  if (bytes < 1024)    return `${bytes} B`;
  if (bytes < 1048576) return `${(bytes / 1024).toFixed(1)} KB`;
  return `${(bytes / 1048576).toFixed(1)} MB`;
}

// ── Step info ─────────────────────────────────────────────────────────
const steps = [
  {
    short: "Pastikan laporan magang sudah selesai ditulis dan disetujui pembimbing lapangan.",
    details: [
      "Laporan harus mencakup seluruh kegiatan selama periode magang berlangsung.",
      "Pastikan laporan sudah mendapat persetujuan atau tanda tangan dari pembimbing lapangan Anda.",
      "Periksa kembali isi laporan: latar belakang, deskripsi kegiatan, hasil, dan kesimpulan.",
      "Jika ada format khusus dari kampus, pastikan sudah sesuai sebelum diupload.",
    ],
  },
  {
    short: "Format file: PDF (disarankan), DOC, atau DOCX. Ukuran maksimal 20 MB.",
    details: [
      "Format PDF sangat disarankan karena tidak berubah tampilannya di berbagai perangkat.",
      "Untuk mengkonversi DOC/DOCX ke PDF: buka di Microsoft Word → File → Simpan Sebagai → PDF.",
      "Ukuran file maksimal adalah 20 MB. Kompres gambar di dalam laporan jika ukuran terlalu besar.",
      "Nama file akan disimpan seperti aslinya, gunakan nama yang deskriptif (contoh: Laporan_Magang_NamaAnda.pdf).",
    ],
  },
  {
    short: "Upload file dan tunggu konfirmasi review dari tim HRD.",
    details: [
      "Setelah upload berhasil, status laporan akan menjadi 'Menunggu Review'.",
      "Tim HRD akan meninjau laporan Anda dalam waktu beberapa hari kerja.",
      "Anda akan mendapat notifikasi ketika HRD memberikan keputusan.",
      "Jika laporan perlu revisi, Anda bisa upload ulang setelah membaca catatan dari HRD.",
      "Pastikan koneksi internet stabil saat proses upload berlangsung.",
    ],
  },
];

const showStepInfo = ref(false);

// ── Preview ───────────────────────────────────────────────────────────
interface PreviewState {
  show: boolean;
  loading: boolean;
  error: boolean;
  url: string;
  type: "pdf" | "other";
  nama: string;
  lapId: string;
}

const previewState = ref<PreviewState>({
  show: false, loading: false, error: false, url: "", type: "other", nama: "", lapId: "",
});

async function openPreview(lap: any) {
  previewState.value = { show: true, loading: true, error: false, url: "", type: "other", nama: lap.nama_file, lapId: lap.id };
  try {
    const res = await api.get(`/api/laporan/${lap.id}/download`, { responseType: "blob" });
    const mime = lap.mime_type ?? res.headers["content-type"] ?? "application/octet-stream";
    const blob = new Blob([res.data], { type: mime });
    const url = URL.createObjectURL(blob);
    const type = mime === "application/pdf" ? "pdf" : "other";
    previewState.value = { ...previewState.value, loading: false, url, type };
  } catch {
    previewState.value = { ...previewState.value, loading: false, error: true };
  }
}

function closePreview() {
  if (previewState.value.url) URL.revokeObjectURL(previewState.value.url);
  previewState.value = { show: false, loading: false, error: false, url: "", type: "other", nama: "", lapId: "" };
}

// ── Upload ────────────────────────────────────────────────────────────
function onFileChange(e: Event) {
  const f = (e.target as HTMLInputElement).files?.[0];
  if (!f) return;
  if (f.size > 20 * 1024 * 1024) { errMsg.value = "Ukuran file melebihi 20 MB."; return; }
  errMsg.value = "";
  selectedFile.value = f;
}

async function doUpload() {
  if (!selectedFile.value) return;
  uploading.value = true; errMsg.value = ""; uploadOk.value = false;
  try {
    const fd = new FormData();
    fd.append("file", selectedFile.value);
    await api.post("/api/laporan/upload", fd, { headers: { "Content-Type": "multipart/form-data" } });
    selectedFile.value = null;
    uploadOk.value = true;
    await fetchLaporan();
  } catch (e: any) {
    errMsg.value = e.response?.data?.message ?? "Gagal mengunggah laporan. Coba lagi.";
  } finally {
    uploading.value = false;
  }
}

async function fetchLaporan() {
  try {
    const r = await api.get("/api/laporan/saya");
    laporanList.value = Array.isArray(r.data?.data) ? r.data.data : [];
  } catch {
    laporanList.value = [];
  }
}

onMounted(async () => {
  await fetchLaporan();
  loading.value = false;
});
</script>

<style scoped>
.upload-card__titlerow { display:flex; align-items:center; justify-content:space-between; gap:10px; margin-bottom:4px; }
.btn-panduan {
  display:flex; align-items:center; gap:5px; flex-shrink:0;
  background:#f0fdf4; border:1.5px solid #bbf7d0; border-radius:8px;
  padding:5px 11px; font-size:12px; font-weight:600; color:#16a34a;
  cursor:pointer; font-family:inherit; transition:all .15s;
}
.btn-panduan:hover { background:#dcfce7; border-color:#48AF4A; }

/* ── State locked ── */
.state-box { display:flex; flex-direction:column; align-items:center; justify-content:center; gap:12px; padding:56px 24px; background:#fff; border-radius:14px; border:1px solid #e9f5e9; text-align:center; }
.locked-icon { margin-bottom:4px; }
.locked-title { font-size:16px; font-weight:700; color:#374151; }
.btn-refresh-status { display:inline-flex; align-items:center; gap:6px; margin-top:8px; background:#f0fdf4; border:1px solid #bbf7d0; color:#16a34a; border-radius:8px; padding:7px 16px; font-size:12.5px; font-weight:600; font-family:inherit; cursor:pointer; }
.btn-refresh-status:hover { background:#dcfce7; }
.locked-sub { font-size:12.5px; color:#9ca3af; line-height:1.7; max-width:380px; }
.status-chip { margin-top:4px; font-size:12px; font-weight:600; padding:5px 14px; border-radius:100px; }
.chip--green  { background:#f0fdf4; color:#16a34a; border:1px solid #86efac; }
.chip--yellow { background:#fffbeb; color:#b45309; border:1px solid #fde68a; }
.chip--blue   { background:#f0fdf4; color:#0d2818; border:1px solid #bbf7d0; }
.chip--gray   { background:#f9fafb; color:#6b7280; border:1px solid #e5e7eb; }

/* ── Review banner ── */
.review-banner { display:flex; align-items:flex-start; gap:14px; border-radius:12px; padding:16px 20px; margin-bottom:16px; border:1.5px solid; }
.review-banner--menunggu_review { background:#fffbeb; border-color:#fde68a; color:#92400e; }
.review-banner--disetujui       { background:#f0fdf4; border-color:#86efac; color:#16a34a; }
.review-banner--revisi          { background:#fef2f2; border-color:#fecaca; color:#991b1b; }
.review-banner__icon { flex-shrink:0; margin-top:1px; }
.review-banner__body { flex:1; min-width:0; }
.review-banner__title { font-size:13.5px; font-weight:700; margin-bottom:3px; }
.review-banner__sub   { font-size:12px; line-height:1.6; }
.review-catatan { margin-top:8px; font-size:12px; display:flex; align-items:flex-start; gap:5px; padding:8px 12px; background:rgba(0,0,0,0.04); border-radius:8px; }
.review-catatan em { font-style:italic; }
.review-versi { font-size:10.5px; font-weight:700; padding:2px 8px; background:rgba(0,0,0,0.08); border-radius:100px; white-space:nowrap; align-self:flex-start; }

/* ── History ── */
.history-card { background:#fff; border:1.5px solid #e9f5e9; border-radius:14px; overflow:hidden; margin-bottom:16px; box-shadow:0 1px 3px rgba(13,40,24,.04); }
.history-card__header { display:flex; align-items:center; justify-content:space-between; padding:12px 16px; border-bottom:1px solid #f0fdf4; background:#fafffe; }
.history-card__title { font-size:13.5px; font-weight:700; color:#111827; }
.history-count { background:#f3f4f6; color:#6b7280; font-size:11px; font-weight:700; padding:3px 12px; border-radius:100px; }
.history-list { display:flex; flex-direction:column; }
.history-item { display:flex; align-items:center; gap:12px; padding:12px 16px; border-bottom:1px solid #f9fafb; }
.history-item:last-child { border-bottom:none; }
.history-item__icon { width:34px; height:34px; border-radius:8px; background:#f3f4f6; color:#6b7280; display:flex; align-items:center; justify-content:center; flex-shrink:0; }
.history-item__info { flex:1; min-width:0; }
.history-item__name { font-size:12.5px; font-weight:600; color:#111827; white-space:nowrap; overflow:hidden; text-overflow:ellipsis; }
.history-item__meta { display:flex; align-items:center; gap:6px; margin-top:2px; font-size:11px; color:#9ca3af; flex-wrap:wrap; }
.versi-badge { background:#f0fdf4; color:#1a5c20; font-size:10px; font-weight:700; padding:1px 6px; border-radius:100px; }
.history-item__right { display:flex; align-items:center; gap:6px; flex-shrink:0; }
.status-chip-sm { font-size:10.5px; font-weight:600; padding:2px 9px; border-radius:100px; white-space:nowrap; }
.chip-sm--menunggu_review { background:#fffbeb; color:#b45309; }
.chip-sm--disetujui       { background:#f0fdf4; color:#16a34a; }
.chip-sm--revisi          { background:#fef2f2; color:#dc2626; }
.btn-action { width:30px; height:30px; border-radius:7px; background:#f3f4f6; border:1px solid #e5e7eb; display:flex; align-items:center; justify-content:center; color:#6b7280; text-decoration:none; transition:all .15s; cursor:pointer; }
.btn-action:hover { background:#e5e7eb; color:#374151; }

/* ── Upload form ── */
.upload-card { background:#fff; border:1.5px solid #e9f5e9; border-radius:14px; overflow:hidden; box-shadow:0 1px 3px rgba(13,40,24,.04); }
.upload-card__topbar { display:flex; align-items:center; justify-content:space-between; padding:12px 16px; border-bottom:1px solid #f0fdf4; background:#fafffe; gap:10px; }
.upload-card__title { font-size:13.5px; font-weight:700; color:#111827; }
.upload-card__body  { padding:20px 20px 24px; }
.upload-card__sub   { font-size:12px; color:#6b7280; margin-bottom:18px; }


.drop-zone {
  display:flex; flex-direction:column; align-items:center; justify-content:center; gap:8px;
  border:2px dashed #d1d5db; border-radius:12px; padding:40px 24px; cursor:pointer;
  transition:border-color .15s,background .15s; text-align:center;
}
.drop-zone:hover,.drop-zone--dragover { border-color:#48AF4A; background:#f0fdf4; }
.drop-zone--has-file { border-color:#86efac; background:#f0fdf4; border-style:solid; }
.dz-label { font-size:13px; font-weight:600; color:#374151; }
.dz-link  { color:#48AF4A; cursor:pointer; text-decoration:underline; }
.dz-sub   { font-size:11.5px; color:#9ca3af; }
.dz-file-name { font-size:13.5px; font-weight:700; color:#16a34a; }
.dz-file-size { font-size:11.5px; color:#6b7280; }
.dz-remove { background:none; border:1px solid #d1d5db; border-radius:6px; padding:4px 12px; font-size:11.5px; color:#6b7280; cursor:pointer; margin-top:4px; }
.dz-remove:hover { background:#f3f4f6; }

.upload-err { font-size:12px; color:#be123c; background:#fff1f2; border:1px solid #fecdd3; border-radius:8px; padding:8px 12px; margin-top:12px; }
.upload-ok  { display:flex; align-items:center; gap:7px; font-size:12.5px; font-weight:600; color:#16a34a; background:#f0fdf4; border:1px solid #bbf7d0; border-radius:8px; padding:10px 14px; margin-top:12px; }

.btn-submit {
  display:flex; align-items:center; gap:8px; justify-content:center;
  margin-top:16px; width:100%; padding:12px; border-radius:10px;
  background:#48AF4A; color:#fff; font-size:13px; font-weight:700;
  font-family:inherit; border:none; cursor:pointer; transition:background .15s;
}
.btn-submit:hover:not(:disabled) { background:#3a9e3c; }
.btn-submit:disabled { opacity:.5; cursor:not-allowed; }
.btn-spinner { width:16px; height:16px; border:2px solid rgba(255,255,255,.4); border-top-color:#fff; border-radius:50%; animation:spin .7s linear infinite; }

/* ── Info box ── */
.info-box { display:flex; align-items:center; gap:10px; padding:14px 16px; border-radius:12px; font-size:13px; margin-top:16px; }
.info-box--green { background:#f0fdf4; border:1.5px solid #86efac; color:#16a34a; }

.spinner { width:28px; height:28px; border:3px solid #e5e7eb; border-top-color:#48AF4A; border-radius:50%; animation:spin .7s linear infinite; }

/* ── Modal Panduan ── */
.modal-backdrop { position:fixed; inset:0; background:rgba(0,0,0,0.5); backdrop-filter:blur(3px); z-index:300; display:flex; align-items:center; justify-content:center; padding:20px; }
.modal { background:#fff; border-radius:16px; width:min(500px,100%); box-shadow:0 20px 60px rgba(0,0,0,0.2); display:flex; flex-direction:column; overflow:hidden; max-height:90vh; }
.modal--info {}
.modal-header { display:flex; align-items:center; justify-content:space-between; padding:16px 20px; border-bottom:1px solid #bbf7d0; gap:12px; background:#f0fdf4; flex-shrink:0; }
.modal-header__left { display:flex; align-items:center; gap:8px; flex:1; min-width:0; }
.modal-title { font-size:14px; font-weight:700; color:#0d2818; }
.modal-close { background:rgba(0,0,0,0.06); border:none; border-radius:8px; width:30px; height:30px; display:flex; align-items:center; justify-content:center; cursor:pointer; color:#0d2818; flex-shrink:0; }
.modal-close:hover { background:rgba(0,0,0,0.12); }
.modal-body { padding:20px; overflow-y:auto; }
.step-list { display:flex; flex-direction:column; gap:16px; }
.step-item { display:flex; gap:14px; align-items:flex-start; }
.step-item__num { width:26px; height:26px; min-width:26px; background:#48AF4A; color:#fff; border-radius:50%; display:flex; align-items:center; justify-content:center; font-size:11px; font-weight:800; margin-top:1px; }
.step-item__body { flex:1; min-width:0; }
.step-item__title { font-size:13px; font-weight:700; color:#0d2818; margin-bottom:6px; line-height:1.4; }
.step-item__details { margin:0; padding-left:16px; display:flex; flex-direction:column; gap:4px; }
.step-item__details li { font-size:12.5px; color:#6b7280; line-height:1.6; }
.modal-footer { display:flex; gap:8px; padding:14px 20px; border-top:1px solid #f3f4f6; flex-shrink:0; }
.modal-footer--single { justify-content:flex-end; }
.btn-ok { background:#48AF4A; color:#fff; border:none; border-radius:9px; padding:9px 22px; font-size:13px; font-weight:700; font-family:inherit; cursor:pointer; }
.btn-ok:hover { background:#3a9e3c; }

/* ── Preview lightbox ── */
.preview-overlay { position:fixed; inset:0; background:rgba(0,0,0,0.75); backdrop-filter:blur(4px); z-index:400; display:flex; align-items:center; justify-content:center; padding:20px; }
.preview-box { background:#fff; border-radius:14px; width:min(860px,100%); max-height:90vh; display:flex; flex-direction:column; overflow:hidden; box-shadow:0 24px 80px rgba(0,0,0,0.3); }
.preview-header { display:flex; align-items:center; justify-content:space-between; padding:14px 18px; border-bottom:1px solid #f3f4f6; gap:12px; flex-shrink:0; }
.preview-filename { font-size:13px; font-weight:600; color:#111827; white-space:nowrap; overflow:hidden; text-overflow:ellipsis; min-width:0; }
.preview-header-actions { display:flex; align-items:center; gap:8px; flex-shrink:0; }
.preview-btn-dl { width:34px; height:34px; border-radius:8px; background:#f3f4f6; border:1px solid #e5e7eb; display:flex; align-items:center; justify-content:center; color:#374151; text-decoration:none; transition:all .15s; }
.preview-btn-dl:hover { background:#e5e7eb; }
.preview-close { width:34px; height:34px; border-radius:8px; background:#f3f4f6; border:none; display:flex; align-items:center; justify-content:center; cursor:pointer; color:#374151; transition:all .15s; }
.preview-close:hover { background:#e5e7eb; }
.preview-body { flex:1; overflow:hidden; min-height:0; }
.preview-loading { display:flex; align-items:center; justify-content:center; height:400px; }
.preview-error { display:flex; flex-direction:column; align-items:center; justify-content:center; gap:12px; height:300px; color:#9ca3af; font-size:13px; }
.preview-iframe { width:100%; height:100%; min-height:500px; border:none; }
.btn-green-sm { background:#48AF4A; color:#fff; border:none; border-radius:8px; padding:8px 16px; font-size:12.5px; font-weight:600; font-family:inherit; cursor:pointer; text-decoration:none; }
.btn-green-sm:hover { background:#3a9e3c; }

@keyframes spin { to { transform:rotate(360deg); } }
</style>
