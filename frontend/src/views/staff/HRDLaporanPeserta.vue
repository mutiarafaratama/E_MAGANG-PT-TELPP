<template>
  <div>
    <!-- Tab filter -->
    <div class="tab-bar">
      <button
        v-for="t in tabs" :key="t.key"
        class="tab-btn"
        :class="{ 'tab-btn--active': activeTab === t.key }"
        @click="switchTab(t.key)"
      >
        {{ t.label }}
        <span v-if="t.key === 'menunggu_review' && totalMenunggu > 0" class="tab-badge">{{ totalMenunggu }}</span>
      </button>
    </div>

    <!-- Tabel -->
    <div class="card">
      <div class="card-head">
        <div class="card-head__title">Laporan Peserta</div>
        <div class="card-head__sub">Review dan kelola laporan akhir magang dari seluruh peserta</div>
      </div>
      <div v-if="loading" class="empty-state"><div class="spinner"></div></div>

      <div v-else-if="list.length === 0" class="empty-state">
        <div class="empty-state__icon">
          <svg width="36" height="36" viewBox="0 0 24 24" fill="none">
            <path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z" stroke="#d1d5db" stroke-width="1.5"/>
            <polyline points="14 2 14 8 20 8" stroke="#d1d5db" stroke-width="1.5"/>
          </svg>
        </div>
        <p>{{ emptyLabel }}</p>
      </div>

      <div v-else class="table-wrap">
        <table class="data-table">
          <thead>
            <tr>
              <th>Peserta</th>
              <th>Institusi / Divisi</th>
              <th>File Laporan</th>
              <th>Versi</th>
              <th>Dikirim</th>
              <th>Status</th>
              <th>Aksi</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="lap in list" :key="lap.id">
              <td>
                <div class="name-cell">
                  <div class="name-avatar">{{ lap.nama_peserta?.[0]?.toUpperCase() ?? '?' }}</div>
                  <div>
                    <div class="name-text">{{ lap.nama_peserta ?? '–' }}</div>
                  </div>
                </div>
              </td>
              <td>
                <div class="cell-two">
                  <span>{{ lap.asal_institusi ?? '–' }}</span>
                  <span class="cell-sub">{{ lap.divisi ?? 'Belum ditempatkan' }}</span>
                </div>
              </td>
              <td>
                <div class="file-cell">
                  <svg width="13" height="13" viewBox="0 0 24 24" fill="none">
                    <path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z" stroke="#6b7280" stroke-width="2"/>
                    <polyline points="14 2 14 8 20 8" stroke="#6b7280" stroke-width="2"/>
                  </svg>
                  <span class="file-name">{{ lap.nama_file }}</span>
                  <!-- Preview -->
                  <button class="btn-icon" title="Preview" @click="openPreview(lap)">
                    <svg width="11" height="11" viewBox="0 0 24 24" fill="none">
                      <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z" stroke="currentColor" stroke-width="2"/>
                      <circle cx="12" cy="12" r="3" stroke="currentColor" stroke-width="2"/>
                    </svg>
                  </button>
                  <!-- Download -->
                  <a :href="`/api/laporan/${lap.id}/download`" target="_blank" class="btn-icon btn-icon--dl" title="Unduh">
                    <svg width="11" height="11" viewBox="0 0 24 24" fill="none">
                      <path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
                      <polyline points="7 10 12 15 17 10" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                      <line x1="12" y1="15" x2="12" y2="3" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
                    </svg>
                  </a>
                </div>
              </td>
              <td><span class="versi-badge">v{{ lap.versi }}</span></td>
              <td class="td-date">{{ fmtDate(lap.diupload_at) }}</td>
              <td><span class="status-chip" :class="`chip--${lap.status}`">{{ statusLabel(lap.status) }}</span></td>
              <td>
                <div class="aksi-cell">
                  <button
                    v-if="lap.status === 'menunggu_review'"
                    class="btn-review"
                    @click="openReview(lap)"
                  >Review</button>
                  <button
                    v-else
                    class="btn-lihat"
                    @click="openReview(lap)"
                  >Lihat</button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Pagination -->
      <div v-if="totalPages > 1" class="pagination">
        <button class="pg-btn" :disabled="page <= 1" @click="page--; fetchData()">‹ Prev</button>
        <span class="pg-info">Hal {{ page }} / {{ totalPages }}</span>
        <button class="pg-btn" :disabled="page >= totalPages" @click="page++; fetchData()">Next ›</button>
      </div>
    </div>

    <!-- Modal Review -->
    <Teleport to="body">
      <div v-if="reviewTarget" class="modal-backdrop" @click.self="closeReview">
        <div class="modal">
          <div class="modal-header">
            <div>
              <div class="modal-title">{{ reviewTarget.status === 'menunggu_review' ? 'Review Laporan' : 'Detail Laporan' }}</div>
              <div class="modal-sub">{{ reviewTarget.nama_peserta }} — v{{ reviewTarget.versi }}</div>
            </div>
            <button class="modal-close" @click="closeReview">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none"><path d="M18 6L6 18M6 6l12 12" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>
            </button>
          </div>

          <div class="modal-body">
            <!-- Info file -->
            <div class="review-file-info">
              <svg width="28" height="28" viewBox="0 0 24 24" fill="none">
                <path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z" stroke="#48AF4A" stroke-width="2" stroke-linejoin="round"/>
                <polyline points="14 2 14 8 20 8" stroke="#48AF4A" stroke-width="2" stroke-linejoin="round"/>
              </svg>
              <div class="review-file-body">
                <div class="review-file-name">{{ reviewTarget.nama_file }}</div>
                <div class="review-file-meta">
                  <span v-if="reviewTarget.ukuran_bytes">{{ fmtSize(reviewTarget.ukuran_bytes) }}</span>
                  <span>Dikirim {{ fmtDate(reviewTarget.diupload_at) }}</span>
                </div>
              </div>
              <div class="review-file-actions">
                <button class="btn-file-action" title="Preview dokumen" @click="openPreview(reviewTarget)">
                  <svg width="13" height="13" viewBox="0 0 24 24" fill="none">
                    <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z" stroke="currentColor" stroke-width="2"/>
                    <circle cx="12" cy="12" r="3" stroke="currentColor" stroke-width="2"/>
                  </svg>
                  Preview
                </button>
                <a :href="`/api/laporan/${reviewTarget.id}/download`" target="_blank" class="btn-file-action btn-file-action--dl">
                  <svg width="13" height="13" viewBox="0 0 24 24" fill="none">
                    <path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
                    <polyline points="7 10 12 15 17 10" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                    <line x1="12" y1="15" x2="12" y2="3" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
                  </svg>
                  Unduh
                </a>
              </div>
            </div>

            <!-- Status (jika sudah direview) -->
            <div v-if="reviewTarget.status !== 'menunggu_review'" class="reviewed-status-box" :class="`reviewed-status--${reviewTarget.status}`">
              <svg v-if="reviewTarget.status === 'disetujui'" width="16" height="16" viewBox="0 0 24 24" fill="none"><circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2"/><path d="M8 12l3 3 5-5" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>
              <svg v-else width="16" height="16" viewBox="0 0 24 24" fill="none"><circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2"/><path d="M23 4v6h-6" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/><path d="M20.49 15a9 9 0 11-2.12-9.36L23 10" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>
              <div>
                <div class="reviewed-status-title">{{ reviewTarget.status === 'disetujui' ? 'Laporan telah disetujui' : 'Laporan diminta revisi' }}</div>
                <div v-if="reviewTarget.catatan_hrd" class="reviewed-status-catatan">Catatan: {{ reviewTarget.catatan_hrd }}</div>
              </div>
            </div>

            <!-- Form catatan (hanya jika menunggu review) -->
            <div v-if="reviewTarget.status === 'menunggu_review'" class="field">
              <label class="field-label">Catatan untuk Peserta <span class="field-hint">(opsional — wajib diisi jika meminta revisi)</span></label>
              <textarea v-model="reviewCatatan" class="field-textarea" rows="3"
                placeholder="Tuliskan catatan atau alasan revisi…"></textarea>
            </div>

            <div v-if="reviewError" class="review-err">{{ reviewError }}</div>
          </div>

          <div class="modal-footer">
            <button class="btn-cancel" @click="closeReview" :disabled="reviewing">Tutup</button>
            <template v-if="reviewTarget.status === 'menunggu_review'">
              <button class="btn-revisi" :disabled="reviewing" @click="submitReview('revisi')">
                <div v-if="reviewing && reviewAction === 'revisi'" class="btn-spinner"></div>
                <svg v-else width="13" height="13" viewBox="0 0 24 24" fill="none"><path d="M23 4v6h-6" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/><path d="M20.49 15a9 9 0 11-2.12-9.36L23 10" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>
                Minta Revisi
              </button>
              <button class="btn-acc" :disabled="reviewing" @click="submitReview('disetujui')">
                <div v-if="reviewing && reviewAction === 'disetujui'" class="btn-spinner"></div>
                <svg v-else width="13" height="13" viewBox="0 0 24 24" fill="none"><polyline points="20 6 9 17 4 12" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>
                Setujui Laporan
              </button>
            </template>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- Preview Lightbox -->
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

interface LaporanDetail {
  id: string;
  pelaksanaan_id: string;
  versi: number;
  nama_file: string;
  path_file: string;
  ukuran_bytes: number | null;
  mime_type: string | null;
  status: string;
  catatan_hrd: string | null;
  direview_at: string | null;
  diupload_at: string;
  nama_peserta: string;
  asal_institusi: string;
  divisi: string | null;
}

const tabs = [
  { key: "menunggu_review", label: "Menunggu Review" },
  { key: "disetujui",       label: "Disetujui" },
  { key: "revisi",          label: "Perlu Revisi" },
  { key: "",                label: "Semua" },
];

const activeTab    = ref("menunggu_review");
const loading      = ref(false);
const list         = ref<LaporanDetail[]>([]);
const page         = ref(1);
const total        = ref(0);
const totalPages   = ref(1);
const totalMenunggu = ref(0);

const emptyLabel = computed(() => {
  const m: Record<string, string> = {
    menunggu_review: "Tidak ada laporan yang menunggu review.",
    disetujui:       "Belum ada laporan yang disetujui.",
    revisi:          "Tidak ada laporan yang perlu direvisi.",
    "":              "Belum ada laporan dari peserta.",
  };
  return m[activeTab.value] ?? "Tidak ada data.";
});

async function fetchData() {
  loading.value = true;
  try {
    const params = new URLSearchParams({ page: String(page.value), limit: "20" });
    if (activeTab.value) params.set("status", activeTab.value);
    const r = await api.get(`/api/laporan?${params}`);
    list.value       = Array.isArray(r.data?.data) ? r.data.data : [];
    total.value      = r.data?.total ?? 0;
    totalPages.value = r.data?.total_pages ?? 1;
  } catch {
    list.value = [];
  } finally {
    loading.value = false;
  }
}

async function fetchMenungguCount() {
  try {
    const r = await api.get("/api/laporan?status=menunggu_review&page=1&limit=1");
    totalMenunggu.value = r.data?.total ?? 0;
  } catch { totalMenunggu.value = 0; }
}

function switchTab(key: string) {
  activeTab.value = key;
  page.value = 1;
  fetchData();
}

// ── Review modal ──────────────────────────────────────────────
const reviewTarget  = ref<LaporanDetail | null>(null);
const reviewCatatan = ref("");
const reviewing     = ref(false);
const reviewAction  = ref<"disetujui" | "revisi" | "">("");
const reviewError   = ref("");

function openReview(lap: LaporanDetail) {
  reviewTarget.value  = lap;
  reviewCatatan.value = "";
  reviewError.value   = "";
  reviewAction.value  = "";
}

function closeReview() {
  reviewTarget.value = null;
}

async function submitReview(action: "disetujui" | "revisi") {
  if (!reviewTarget.value) return;
  if (action === "revisi" && !reviewCatatan.value.trim()) {
    reviewError.value = "Catatan wajib diisi saat meminta revisi.";
    return;
  }
  reviewError.value = "";
  reviewing.value   = true;
  reviewAction.value = action;
  try {
    await api.patch(`/api/laporan/${reviewTarget.value.id}/review`, {
      status: action,
      catatan_hrd: reviewCatatan.value.trim(),
    });
    closeReview();
    page.value = 1;
    await Promise.all([fetchData(), fetchMenungguCount()]);
  } catch (e: any) {
    reviewError.value = e.response?.data?.message ?? "Gagal mengirim review.";
  } finally {
    reviewing.value    = false;
    reviewAction.value = "";
  }
}

// ── Preview ───────────────────────────────────────────────────
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

async function openPreview(lap: LaporanDetail) {
  previewState.value = { show: true, loading: true, error: false, url: "", type: "other", nama: lap.nama_file, lapId: lap.id };
  try {
    const res = await api.get(`/api/laporan/${lap.id}/download`, { responseType: "blob" });
    const mime = lap.mime_type ?? (res.headers["content-type"] as string | undefined) ?? "application/octet-stream";
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

// ── Helpers ───────────────────────────────────────────────────
function statusLabel(s: string) {
  return ({
    menunggu_review: "Menunggu Review",
    disetujui:       "Disetujui",
    revisi:          "Perlu Revisi",
  } as Record<string, string>)[s] ?? s;
}

function fmtDate(iso: string) {
  if (!iso) return "–";
  return new Date(iso).toLocaleDateString("id-ID", { day: "numeric", month: "short", year: "numeric" });
}

function fmtSize(bytes: number | null) {
  if (!bytes) return "";
  if (bytes < 1024)    return `${bytes} B`;
  if (bytes < 1048576) return `${(bytes / 1024).toFixed(1)} KB`;
  return `${(bytes / 1048576).toFixed(1)} MB`;
}

onMounted(() => {
  fetchData();
  fetchMenungguCount();
});
</script>

<style scoped>
.card-head { padding:16px 20px; border-bottom:1px solid #f0faf0; }
.card-head__title { font-size:14px; font-weight:700; color:#111827; margin-bottom:2px; }
.card-head__sub   { font-size:12px; color:#9ca3af; }

/* ── Tab bar ── */
.tab-bar { display:flex; gap:4px; background:#fff; border:1px solid #e9f5e9; border-radius:12px; padding:4px; margin-bottom:16px; overflow-x:auto; }
.tab-btn { display:flex; align-items:center; gap:6px; padding:7px 14px; border-radius:9px; border:none; font-size:12.5px; font-weight:600; font-family:inherit; cursor:pointer; background:none; color:#6b7280; white-space:nowrap; transition:all .15s; }
.tab-btn:hover { background:#f0fdf4; color:#16a34a; }
.tab-btn--active { background:#48AF4A; color:#fff; }
.tab-badge { background:rgba(255,255,255,0.25); color:inherit; font-size:10px; font-weight:800; padding:1px 7px; border-radius:100px; }
.tab-btn--active .tab-badge { background:rgba(0,0,0,0.15); }
.tab-btn:not(.tab-btn--active) .tab-badge { background:#ef4444; color:#fff; }

/* ── Card / table ── */
.card { background:#fff; border-radius:14px; border:1px solid #e9f5e9; box-shadow:0 1px 3px rgba(13,40,24,.05); overflow:hidden; }
.table-wrap { overflow-x:auto; }
.data-table { width:100%; border-collapse:collapse; font-size:12.5px; }
.data-table th { padding:11px 16px; text-align:left; font-size:10.5px; font-weight:600; color:#6b7280; background:#f9fafb; border-bottom:1px solid #f1f5f9; text-transform:uppercase; letter-spacing:0.04em; white-space:nowrap; }
.data-table td { padding:12px 16px; border-bottom:1px solid #f9fafb; color:#374151; vertical-align:middle; }
.data-table tr:last-child td { border-bottom:none; }

.name-cell { display:flex; align-items:center; gap:10px; }
.name-avatar { width:32px; height:32px; border-radius:8px; background:linear-gradient(135deg,#48AF4A,#1a5c20); color:#fff; font-size:12px; font-weight:700; display:flex; align-items:center; justify-content:center; flex-shrink:0; }
.name-text { font-weight:600; color:#111827; font-size:12.5px; }
.cell-two { display:flex; flex-direction:column; gap:2px; }
.cell-sub { font-size:11px; color:#9ca3af; }

.file-cell { display:flex; align-items:center; gap:5px; min-width:0; }
.file-name { font-size:12px; color:#374151; white-space:nowrap; overflow:hidden; text-overflow:ellipsis; max-width:140px; }
.btn-icon { width:24px; height:24px; border-radius:6px; background:#f3f4f6; border:1px solid #e5e7eb; display:flex; align-items:center; justify-content:center; color:#6b7280; cursor:pointer; flex-shrink:0; transition:all .15s; }
.btn-icon:hover { background:#e5e7eb; color:#374151; }
.btn-icon--dl { text-decoration:none; }

.versi-badge { background:#f0fdf4; color:#1a5c20; font-size:10.5px; font-weight:700; padding:2px 8px; border-radius:100px; }
.td-date { white-space:nowrap; font-size:12px; color:#6b7280; }

.status-chip { font-size:10.5px; font-weight:700; padding:3px 10px; border-radius:100px; white-space:nowrap; }
.chip--menunggu_review { background:#fffbeb; color:#b45309; }
.chip--disetujui       { background:#f0fdf4; color:#16a34a; }
.chip--revisi          { background:#fef2f2; color:#dc2626; }

.aksi-cell { display:flex; align-items:center; gap:6px; }
.btn-review { border:none; border-radius:7px; padding:5px 11px; font-size:11.5px; font-weight:700; font-family:inherit; cursor:pointer; background:#dcfce7; color:#0d2818; white-space:nowrap; transition:background .15s; }
.btn-review:hover { background:#bbf7d0; }
.btn-lihat { border:none; border-radius:7px; padding:5px 11px; font-size:11.5px; font-weight:700; font-family:inherit; cursor:pointer; background:#f3f4f6; color:#6b7280; white-space:nowrap; transition:background .15s; }
.btn-lihat:hover { background:#e5e7eb; color:#374151; }

/* ── Pagination ── */
.pagination { display:flex; align-items:center; justify-content:center; gap:12px; padding:14px; border-top:1px solid #f3f4f6; }
.pg-btn { border:1px solid #e5e7eb; background:#fff; border-radius:8px; padding:6px 14px; font-size:12.5px; font-weight:600; font-family:inherit; cursor:pointer; color:#374151; }
.pg-btn:disabled { opacity:.4; cursor:not-allowed; }
.pg-btn:not(:disabled):hover { background:#f0fdf4; border-color:#86efac; }
.pg-info { font-size:12px; color:#6b7280; }

/* ── Empty / spinner ── */
.empty-state { display:flex; flex-direction:column; align-items:center; padding:48px 24px; gap:12px; text-align:center; }
.empty-state__icon { width:72px; height:72px; background:#f9fafb; border-radius:50%; display:flex; align-items:center; justify-content:center; }
.empty-state p { font-size:13px; color:#9ca3af; margin:0; }
.spinner { width:24px; height:24px; border:2.5px solid #e5e7eb; border-top-color:#48AF4A; border-radius:50%; animation:spin .7s linear infinite; }

/* ── Modal ── */
.modal-backdrop { position:fixed; inset:0; background:rgba(0,0,0,0.5); backdrop-filter:blur(3px); z-index:300; display:flex; align-items:center; justify-content:center; padding:20px; }
.modal { background:#fff; border-radius:16px; width:min(540px,100%); box-shadow:0 20px 60px rgba(0,0,0,0.2); display:flex; flex-direction:column; overflow:hidden; }
.modal-header { display:flex; align-items:flex-start; justify-content:space-between; padding:20px 24px 16px; border-bottom:1px solid #f0faf0; gap:12px; }
.modal-title { font-size:16px; font-weight:700; color:#111827; margin:0 0 2px; }
.modal-sub   { font-size:12px; color:#6b7280; }
.modal-close { background:#f3f4f6; border:none; border-radius:8px; width:32px; height:32px; display:flex; align-items:center; justify-content:center; cursor:pointer; color:#6b7280; flex-shrink:0; }
.modal-close:hover { background:#e5e7eb; color:#111827; }
.modal-body { padding:20px 24px; display:flex; flex-direction:column; gap:16px; }
.modal-footer { display:flex; gap:8px; padding:16px 24px; border-top:1px solid #f0faf0; }

/* ── Review file info ── */
.review-file-info { display:flex; align-items:center; gap:12px; padding:12px 14px; background:#f0fdf4; border:1px solid #bbf7d0; border-radius:12px; }
.review-file-body { flex:1; min-width:0; }
.review-file-name { font-size:13px; font-weight:600; color:#111827; white-space:nowrap; overflow:hidden; text-overflow:ellipsis; }
.review-file-meta { display:flex; gap:8px; font-size:11px; color:#6b7280; margin-top:2px; }
.review-file-actions { display:flex; align-items:center; gap:6px; flex-shrink:0; }
.btn-file-action { display:flex; align-items:center; gap:5px; background:#fff; border:1px solid #d1d5db; border-radius:8px; padding:5px 10px; font-size:12px; font-weight:600; color:#374151; cursor:pointer; white-space:nowrap; font-family:inherit; text-decoration:none; transition:all .15s; }
.btn-file-action:hover { background:#f0fdf4; border-color:#86efac; }
.btn-file-action--dl { color:#374151; }

/* ── Reviewed status ── */
.reviewed-status-box { display:flex; align-items:flex-start; gap:10px; padding:12px 14px; border-radius:10px; border:1px solid; font-size:12.5px; }
.reviewed-status--disetujui { background:#f0fdf4; border-color:#86efac; color:#16a34a; }
.reviewed-status--revisi    { background:#fef2f2; border-color:#fecaca; color:#991b1b; }
.reviewed-status-title { font-weight:700; font-size:13px; }
.reviewed-status-catatan { margin-top:3px; font-size:12px; opacity:.8; }

.field-label { font-size:12px; font-weight:600; color:#374151; display:block; margin-bottom:6px; }
.field-hint  { font-size:11px; color:#9ca3af; font-weight:400; }
.field-textarea { width:100%; border:1.5px solid #d1d5db; border-radius:9px; padding:9px 12px; font-size:12.5px; font-family:inherit; color:#111827; resize:vertical; outline:none; box-sizing:border-box; }
.field-textarea:focus { border-color:#48AF4A; box-shadow:0 0 0 3px rgba(72,175,74,.1); }
.review-err { font-size:12px; color:#dc2626; background:#fef2f2; border:1px solid #fecaca; border-radius:8px; padding:8px 12px; }

.btn-cancel { flex:1; background:#f3f4f6; color:#374151; border:none; border-radius:9px; padding:10px 14px; font-size:12.5px; font-weight:600; font-family:inherit; cursor:pointer; }
.btn-cancel:hover:not(:disabled) { background:#e5e7eb; }
.btn-cancel:disabled { opacity:.5; cursor:not-allowed; }
.btn-revisi { flex:1.5; display:flex; align-items:center; justify-content:center; gap:7px; border:none; border-radius:9px; padding:10px 14px; font-size:12.5px; font-weight:700; font-family:inherit; cursor:pointer; background:#fef2f2; color:#dc2626; }
.btn-revisi:hover:not(:disabled) { background:#fee2e2; }
.btn-revisi:disabled { opacity:.5; cursor:not-allowed; }
.btn-acc { flex:1.5; display:flex; align-items:center; justify-content:center; gap:7px; border:none; border-radius:9px; padding:10px 14px; font-size:12.5px; font-weight:700; font-family:inherit; cursor:pointer; background:#48AF4A; color:#fff; }
.btn-acc:hover:not(:disabled) { background:#48AF4A; }
.btn-acc:disabled { opacity:.5; cursor:not-allowed; }
.btn-spinner { width:14px; height:14px; border:2px solid rgba(255,255,255,.35); border-top-color:currentColor; border-radius:50%; animation:spin .7s linear infinite; }

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
