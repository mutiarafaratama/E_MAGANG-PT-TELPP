<template>
  <div>
    <!-- Header + filter bar -->
    <div class="card" style="margin-bottom:0">
      <div class="card-header">
        <div>
          <h3 class="card-title">Kelola Dokumen</h3>
          <div class="card-sub">Semua file yang diupload peserta maupun HRD ke sistem</div>
        </div>
        <div class="header-actions">
          <span v-if="total > 0" class="badge-count">{{ total }} dokumen</span>
          <button class="btn-refresh" @click="fetchDokumen" :disabled="loading" title="Muat Ulang">
            <svg width="13" height="13" viewBox="0 0 24 24" fill="none" :class="{ spinning: loading }"><polyline points="1 4 1 10 7 10" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/><path d="M3.51 15a9 9 0 102.13-9.36L1 10" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>
            Muat Ulang
          </button>
        </div>
      </div>

      <!-- Filter -->
      <div class="filter-bar">
        <div class="filter-group">
          <label class="filter-label">Filter Jenis</label>
          <select v-model="filterJenis" class="filter-select" @change="goPage(1)">
            <option value="">Semua Jenis</option>
            <option v-for="j in JENIS_OPTIONS" :key="j.value" :value="j.value">{{ j.label }}</option>
          </select>
        </div>
        <div class="filter-group">
          <label class="filter-label">Cari Pemilik</label>
          <div class="search-input-wrap">
            <svg class="search-icon" width="13" height="13" viewBox="0 0 24 24" fill="none"><circle cx="11" cy="11" r="8" stroke="currentColor" stroke-width="2"/><line x1="21" y1="21" x2="16.65" y2="16.65" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>
            <input
              v-model="filterSearch"
              type="text"
              class="filter-search"
              placeholder="Nama atau email pemilik…"
              @input="onSearchInput"
            />
            <button v-if="filterSearch" class="search-clear" @click="filterSearch=''; goPage(1)" title="Hapus">×</button>
          </div>
        </div>
        <div class="filter-group" style="margin-left:auto">
          <label class="filter-label">Per halaman</label>
          <select v-model.number="limitVal" class="filter-select filter-select--sm" @change="goPage(1)">
            <option :value="20">20</option>
            <option :value="50">50</option>
            <option :value="100">100</option>
          </select>
        </div>
      </div>
    </div>

    <!-- Table -->
    <div class="card" style="border-top-left-radius:0;border-top-right-radius:0;border-top:none">
      <div v-if="loading" class="empty-state" style="padding:52px">
        <div class="spinner"></div>
      </div>
      <div v-else-if="fetchError" class="empty-state">
        <div class="empty-state__icon">
          <svg width="36" height="36" viewBox="0 0 24 24" fill="none"><circle cx="12" cy="12" r="10" stroke="#f87171" stroke-width="1.5"/><line x1="12" y1="8" x2="12" y2="12" stroke="#f87171" stroke-width="2" stroke-linecap="round"/><line x1="12" y1="16" x2="12.01" y2="16" stroke="#f87171" stroke-width="2" stroke-linecap="round"/></svg>
        </div>
        <p style="color:#dc2626">{{ fetchError }}</p>
        <button class="btn-retry" @click="fetchDokumen">Coba Lagi</button>
      </div>
      <div v-else-if="list.length === 0" class="empty-state">
        <div class="empty-state__icon">
          <svg width="36" height="36" viewBox="0 0 24 24" fill="none"><path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z" stroke="#d1d5db" stroke-width="1.5"/><polyline points="14 2 14 8 20 8" stroke="#d1d5db" stroke-width="1.5"/></svg>
        </div>
        <p>Belum ada dokumen yang tersimpan.</p>
      </div>
      <div v-else class="table-wrap">
        <table class="data-table">
          <thead>
            <tr>
              <th>#</th>
              <th>Nama File</th>
              <th>Jenis</th>
              <th>Pemilik</th>
              <th>Ukuran</th>
              <th>Tanggal Upload</th>
              <th>Aksi</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(d, i) in list" :key="d.id">
              <td class="td-no">{{ (currentPage - 1) * limitVal + i + 1 }}</td>
              <td class="td-file">
                <div class="file-icon" :class="iconClass(d.mime_type)">
                  <svg v-if="isPDF(d.mime_type)" width="12" height="12" viewBox="0 0 24 24" fill="none"><path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z" stroke="currentColor" stroke-width="2"/><polyline points="14 2 14 8 20 8" stroke="currentColor" stroke-width="2"/><line x1="8" y1="13" x2="16" y2="13" stroke="currentColor" stroke-width="2"/><line x1="8" y1="17" x2="16" y2="17" stroke="currentColor" stroke-width="2"/></svg>
                  <svg v-else-if="isImage(d.mime_type)" width="12" height="12" viewBox="0 0 24 24" fill="none"><rect x="3" y="3" width="18" height="18" rx="2" stroke="currentColor" stroke-width="2"/><circle cx="8.5" cy="8.5" r="1.5" stroke="currentColor" stroke-width="2"/><polyline points="21 15 16 10 5 21" stroke="currentColor" stroke-width="2"/></svg>
                  <svg v-else width="12" height="12" viewBox="0 0 24 24" fill="none"><path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z" stroke="currentColor" stroke-width="2"/><polyline points="14 2 14 8 20 8" stroke="currentColor" stroke-width="2"/></svg>
                </div>
                <span class="td-file__name" :title="d.nama_file">{{ d.nama_file }}</span>
              </td>
              <td><span class="jenis-badge" :class="`jenis-badge--${d.jenis}`">{{ JENIS_LABEL[d.jenis] ?? d.jenis }}</span></td>
              <td class="td-pemilik">
                <span v-if="d.nama_pemilik">{{ d.nama_pemilik }}</span>
                <span v-else class="td-anon">— Publik —</span>
              </td>
              <td class="td-size">{{ formatSize(d.ukuran_bytes) }}</td>
              <td class="td-date">{{ formatDate(d.uploaded_at) }}</td>
              <td>
                <div class="aksi-cell">
                  <button class="btn-aksi btn-aksi--green" @click="preview(d)" title="Lihat">
                    <svg width="12" height="12" viewBox="0 0 24 24" fill="none"><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z" stroke="currentColor" stroke-width="2"/><circle cx="12" cy="12" r="3" stroke="currentColor" stroke-width="2"/></svg>
                    Lihat
                  </button>
                  <button class="btn-aksi btn-aksi--ghost" title="Download" @click="downloadFile(d)">
                    <svg width="12" height="12" viewBox="0 0 24 24" fill="none"><path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4" stroke="currentColor" stroke-width="2" stroke-linecap="round"/><polyline points="7 10 12 15 17 10" stroke="currentColor" stroke-width="2" stroke-linecap="round"/><line x1="12" y1="15" x2="12" y2="3" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>
                    Download
                  </button>
                  <button class="btn-aksi btn-aksi--red" @click="confirmHapus(d)" title="Hapus">
                    <svg width="12" height="12" viewBox="0 0 24 24" fill="none"><polyline points="3 6 5 6 21 6" stroke="currentColor" stroke-width="2"/><path d="M19 6l-1 14a2 2 0 01-2 2H8a2 2 0 01-2-2L5 6" stroke="currentColor" stroke-width="2"/><path d="M10 11v6M14 11v6" stroke="currentColor" stroke-width="2" stroke-linecap="round"/><path d="M9 6V4a1 1 0 011-1h4a1 1 0 011 1v2" stroke="currentColor" stroke-width="2"/></svg>
                    Hapus
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Pagination -->
      <div v-if="totalPages > 1" class="pagination">
        <button class="pg-btn" :disabled="currentPage <= 1" @click="goPage(currentPage - 1)">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none"><polyline points="15 18 9 12 15 6" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>
        </button>
        <span class="pg-info">Halaman {{ currentPage }} / {{ totalPages }}</span>
        <button class="pg-btn" :disabled="currentPage >= totalPages" @click="goPage(currentPage + 1)">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none"><polyline points="9 18 15 12 9 6" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>
        </button>
      </div>
    </div>

    <!-- Toast -->
    <Transition name="toast-fade">
      <div v-if="toast" class="toast" :class="`toast--${toastType}`">{{ toast }}</div>
    </Transition>
  </div>

  <!-- Preview Modal -->
  <Teleport to="body">
    <Transition name="modal-fade">
      <div v-if="previewModal.show" class="modal-backdrop" @click.self="closePreview">
        <div class="modal-box modal-box--preview">
          <div class="modal-header">
            <div>
              <div class="modal-title">{{ previewModal.nama }}</div>
              <div class="modal-sub">{{ JENIS_LABEL[previewModal.jenis ?? ''] ?? previewModal.jenis }}</div>
            </div>
            <div style="display:flex;gap:8px;align-items:center">
              <a v-if="previewModal.blobUrl" :href="previewModal.blobUrl"
                 :download="previewModal.nama" class="btn-aksi btn-aksi--green" style="font-size:12px;padding:7px 14px;text-decoration:none">
                <svg width="12" height="12" viewBox="0 0 24 24" fill="none"><path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4" stroke="currentColor" stroke-width="2" stroke-linecap="round"/><polyline points="7 10 12 15 17 10" stroke="currentColor" stroke-width="2" stroke-linecap="round"/><line x1="12" y1="15" x2="12" y2="3" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>
                Download
              </a>
              <button class="modal-close" @click="closePreview">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none"><line x1="18" y1="6" x2="6" y2="18" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"/><line x1="6" y1="6" x2="18" y2="18" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"/></svg>
              </button>
            </div>
          </div>
          <div class="modal-body">
            <div v-if="previewModal.loading" class="preview-loading">
              <div class="spinner"></div>
              <span>Memuat file…</span>
            </div>
            <img v-else-if="previewModal.blobUrl && previewModal.type === 'image'"
                 :src="previewModal.blobUrl" class="preview-img" alt="preview" />
            <iframe v-else-if="previewModal.blobUrl && previewModal.type === 'pdf'"
                    :src="previewModal.blobUrl" class="preview-iframe"></iframe>
            <div v-else-if="previewModal.blobUrl" class="preview-other">
              <svg width="44" height="44" viewBox="0 0 24 24" fill="none"><path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z" stroke="#9ca3af" stroke-width="1.5"/><polyline points="14 2 14 8 20 8" stroke="#9ca3af" stroke-width="1.5"/></svg>
              <p>Format file tidak dapat ditampilkan di browser.</p>
              <a :href="previewModal.blobUrl" :download="previewModal.nama"
                 class="btn-aksi btn-aksi--green" style="text-decoration:none;margin-top:8px">Download File</a>
            </div>
            <div v-else-if="!previewModal.loading" class="preview-other">
              <p style="color:#ef4444">Gagal memuat file.</p>
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>

  <!-- Hapus Confirm Modal -->
  <Teleport to="body">
    <Transition name="modal-fade">
      <div v-if="hapusModal.show" class="modal-backdrop" @click.self="hapusModal.show = false">
        <div class="modal-box modal-box--confirm">
          <div class="confirm-icon">
            <svg width="26" height="26" viewBox="0 0 24 24" fill="none"><path d="M12 9v4M12 17h.01" stroke="#ef4444" stroke-width="2.5" stroke-linecap="round"/><path d="M10.29 3.86L1.82 18a2 2 0 001.71 3h16.94a2 2 0 001.71-3L13.71 3.86a2 2 0 00-3.42 0z" stroke="#ef4444" stroke-width="2"/></svg>
          </div>
          <div class="confirm-title">Hapus Dokumen?</div>
          <div class="confirm-body">
            File <strong>{{ hapusModal.nama }}</strong> akan dihapus permanen dari server dan database. Aksi ini tidak dapat dibatalkan.
          </div>
          <div class="confirm-actions">
            <button class="btn-cancel" @click="hapusModal.show = false" :disabled="hapusModal.loading">Batal</button>
            <button class="btn-red" @click="doHapus" :disabled="hapusModal.loading">
              <div v-if="hapusModal.loading" class="spinner-sm" style="border-top-color:#fff"></div>
              {{ hapusModal.loading ? 'Menghapus…' : 'Ya, Hapus' }}
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue';
import api from '@/lib/api';

interface DokumenItem {
  id: string;
  pengajuan_id: string | null;
  user_id: string | null;
  jenis: string;
  nama_file: string;
  path_file: string;
  ukuran_bytes: number | null;
  mime_type: string | null;
  uploaded_at: string;
  nama_pemilik: string | null;
  email_pemilik: string | null;
}

const JENIS_OPTIONS = [
  { value: 'proposal_magang',  label: 'Proposal Magang' },
  { value: 'surat_pengantar',  label: 'Surat Pengantar' },
  { value: 'ktp',              label: 'KTP' },
  { value: 'ktm',              label: 'KTM' },
  { value: 'pasfoto',          label: 'Pas Foto' },
  { value: 'bpjs_kis',         label: 'BPJS / KIS' },
  { value: 'surat_balasan',    label: 'Surat Balasan' },
  { value: 'laporan_magang',   label: 'Laporan Magang' },
  { value: 'sertifikat',       label: 'Sertifikat' },
];

const JENIS_LABEL: Record<string, string> = Object.fromEntries(
  JENIS_OPTIONS.map(j => [j.value, j.label])
);

const list         = ref<DokumenItem[]>([]);
const loading      = ref(false);
const fetchError   = ref('');
const total        = ref(0);
const currentPage  = ref(1);
const limitVal     = ref(20);
const filterJenis  = ref('');
const filterSearch = ref('');
const toast        = ref('');
const toastType    = ref<'ok' | 'err'>('ok');
let toastTimer: ReturnType<typeof setTimeout> | null = null;
let searchTimer: ReturnType<typeof setTimeout> | null = null;

const totalPages = computed(() => Math.max(1, Math.ceil(total.value / limitVal.value)));

function onSearchInput() {
  if (searchTimer) clearTimeout(searchTimer);
  searchTimer = setTimeout(() => goPage(1), 400);
}

async function fetchDokumen() {
  loading.value = true;
  fetchError.value = '';
  try {
    const params: Record<string, string | number> = {
      page: currentPage.value,
      limit: limitVal.value,
    };
    if (filterJenis.value)  params.jenis  = filterJenis.value;
    if (filterSearch.value) params.search = filterSearch.value;
    const r = await api.get('/api/admin/dokumen', { params });
    list.value  = Array.isArray(r.data?.data) ? r.data.data : [];
    total.value = r.data?.total ?? 0;
  } catch (e: any) {
    list.value  = [];
    total.value = 0;
    fetchError.value = e?.code === 'ERR_NETWORK' || e?.message?.includes('ECONNREFUSED')
      ? 'Koneksi ke server gagal. Klik Muat Ulang untuk mencoba lagi.'
      : (e?.response?.data?.message ?? 'Gagal memuat data dokumen.');
  } finally {
    loading.value = false;
  }
}

function goPage(p: number) {
  if (p < 1 || p > totalPages.value) return;
  currentPage.value = p;
  fetchDokumen();
}

function showToast(msg: string, type: 'ok' | 'err' = 'ok') {
  if (toastTimer) clearTimeout(toastTimer);
  toast.value = msg;
  toastType.value = type;
  toastTimer = setTimeout(() => { toast.value = ''; }, 3000);
}

// ── Preview ──────────────────────────────────────────────────
const previewModal = ref<{
  show: boolean; loading: boolean;
  blobUrl: string; type: 'image' | 'pdf' | 'other';
  nama: string; jenis: string | null;
}>({ show: false, loading: false, blobUrl: '', type: 'other', nama: '', jenis: null });

async function preview(d: DokumenItem) {
  if (previewModal.value.blobUrl) URL.revokeObjectURL(previewModal.value.blobUrl);
  previewModal.value = { show: true, loading: true, blobUrl: '', type: 'other', nama: d.nama_file, jenis: d.jenis };
  try {
    const r = await api.get(`/api/admin/dokumen/${d.id}/download`, { responseType: 'blob' });
    const mimeType = String(r.headers['content-type'] || d.mime_type || '');
    const blob = new Blob([r.data], { type: mimeType });
    const url  = URL.createObjectURL(blob);
    const type = isPDF(mimeType) ? 'pdf' : isImage(mimeType) ? 'image' : 'other';
    previewModal.value = { show: true, loading: false, blobUrl: url, type, nama: d.nama_file, jenis: d.jenis };
  } catch {
    previewModal.value.loading = false;
    showToast('Gagal memuat file.', 'err');
  }
}

function closePreview() {
  if (previewModal.value.blobUrl) URL.revokeObjectURL(previewModal.value.blobUrl);
  previewModal.value.show = false;
  previewModal.value.blobUrl = '';
}

// ── Download ─────────────────────────────────────────────────
async function downloadFile(d: DokumenItem) {
  try {
    const r = await api.get(`/api/admin/dokumen/${d.id}/download`, { responseType: 'blob' });
    const mimeType = String(r.headers['content-type'] || d.mime_type || 'application/octet-stream');
    const blob = new Blob([r.data], { type: mimeType });
    const url  = URL.createObjectURL(blob);
    const a    = document.createElement('a');
    a.href     = url;
    a.download = d.nama_file;
    a.click();
    setTimeout(() => URL.revokeObjectURL(url), 5000);
  } catch {
    showToast('Gagal mendownload file.', 'err');
  }
}

// ── Hapus ─────────────────────────────────────────────────────
const hapusModal = ref<{ show: boolean; loading: boolean; id: string; nama: string }>(
  { show: false, loading: false, id: '', nama: '' }
);

function confirmHapus(d: DokumenItem) {
  hapusModal.value = { show: true, loading: false, id: d.id, nama: d.nama_file };
}

async function doHapus() {
  hapusModal.value.loading = true;
  try {
    await api.delete(`/api/admin/dokumen/${hapusModal.value.id}`);
    hapusModal.value.show = false;
    showToast('Dokumen berhasil dihapus.');
    await fetchDokumen();
  } catch (e: any) {
    showToast(e.response?.data?.message ?? 'Gagal menghapus dokumen.', 'err');
  } finally {
    hapusModal.value.loading = false;
  }
}

// ── Helpers ──────────────────────────────────────────────────
function isPDF(mime: string | null) { return !!mime && mime.includes('pdf'); }
function isImage(mime: string | null) { return !!mime && mime.startsWith('image/'); }
function iconClass(mime: string | null) {
  if (isPDF(mime)) return 'file-icon--pdf';
  if (isImage(mime)) return 'file-icon--img';
  return 'file-icon--doc';
}
function formatSize(bytes: number | null) {
  if (!bytes) return '–';
  if (bytes < 1024) return bytes + ' B';
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB';
  return (bytes / 1024 / 1024).toFixed(2) + ' MB';
}
function formatDate(iso: string) {
  if (!iso) return '–';
  return new Date(iso).toLocaleDateString('id-ID', { day: '2-digit', month: 'short', year: 'numeric', hour: '2-digit', minute: '2-digit' });
}

onMounted(() => { fetchDokumen(); });
onUnmounted(() => {
  if (previewModal.value.blobUrl) URL.revokeObjectURL(previewModal.value.blobUrl);
  if (toastTimer) clearTimeout(toastTimer);
});
</script>

<style scoped>
.card { background: #fff; border-radius: 14px; border: 1px solid #e9f5e9; box-shadow: 0 1px 3px rgba(13,40,24,0.05); overflow: hidden; }
.card-header { display: flex; align-items: center; justify-content: space-between; padding: 16px 20px; border-bottom: 1px solid #f0faf0; gap: 12px; }
.card-title  { font-size: 13.5px; font-weight: 700; color: #111827; margin: 0 0 2px; }
.card-sub    { font-size: 11.5px; color: #9ca3af; }
.header-actions { display: flex; align-items: center; gap: 8px; }
.badge-count { background: #f0fdf4; color: #16a34a; border: 1px solid #bbf7d0; border-radius: 100px; font-size: 11px; font-weight: 700; padding: 3px 10px; }
.btn-refresh { display: inline-flex; align-items: center; gap: 5px; background: #f3f4f6; border: 1px solid #e5e7eb; border-radius: 7px; padding: 5px 10px; font-size: 11px; font-weight: 600; color: #374151; cursor: pointer; font-family: inherit; transition: background .15s; }
.btn-refresh:hover:not(:disabled) { background: #e5e7eb; }
.btn-refresh:disabled { opacity: .5; cursor: default; }
.btn-refresh svg.spinning { animation: spin .8s linear infinite; }
.btn-retry { background: #dc2626; color: #fff; border: none; border-radius: 7px; padding: 7px 18px; font-size: 12px; font-weight: 700; font-family: inherit; cursor: pointer; transition: background .15s; }
.btn-retry:hover { background: #b91c1c; }
@keyframes spin { to { transform: rotate(360deg); } }

.filter-bar { display: flex; align-items: flex-end; gap: 16px; padding: 14px 20px; background: #f9fafb; border-bottom: 1px solid #f0faf0; flex-wrap: wrap; }
.filter-group { display: flex; flex-direction: column; gap: 4px; }
.filter-label { font-size: 11px; font-weight: 600; color: #6b7280; text-transform: uppercase; letter-spacing: .04em; }
.filter-select { border: 1.5px solid #e5e7eb; border-radius: 8px; padding: 7px 10px; font-size: 12.5px; font-family: inherit; color: #111827; background: #fff; outline: none; cursor: pointer; }
.filter-select:focus { border-color: #48AF4A; }
.filter-select--sm { min-width: 72px; }
.search-input-wrap { position: relative; display: flex; align-items: center; }
.search-icon { position: absolute; left: 9px; color: #9ca3af; pointer-events: none; }
.filter-search {
  border: 1.5px solid #e5e7eb; border-radius: 8px; padding: 7px 28px 7px 28px;
  font-size: 12.5px; font-family: inherit; color: #111827; background: #fff;
  outline: none; min-width: 210px; transition: border-color .15s;
}
.filter-search:focus { border-color: #48AF4A; }
.filter-search::placeholder { color: #b4b9c1; }
.search-clear {
  position: absolute; right: 8px; background: none; border: none; cursor: pointer;
  color: #9ca3af; font-size: 16px; line-height: 1; padding: 0 2px;
}
.search-clear:hover { color: #374151; }

.empty-state { display: flex; flex-direction: column; align-items: center; padding: 52px 24px; gap: 12px; text-align: center; }
.empty-state__icon { width: 72px; height: 72px; background: #f9fafb; border-radius: 50%; display: flex; align-items: center; justify-content: center; }
.empty-state p { font-size: 13px; color: #9ca3af; margin: 0; }

.table-wrap { overflow-x: auto; }
.data-table { width: 100%; border-collapse: collapse; font-size: 12.5px; }
.data-table th { padding: 10px 14px; text-align: left; font-size: 10.5px; font-weight: 700; color: #6b7280; background: #f9fafb; border-bottom: 1px solid #f1f5f9; text-transform: uppercase; letter-spacing: .04em; white-space: nowrap; }
.data-table td { padding: 12px 14px; border-bottom: 1px solid #f9fafb; color: #374151; vertical-align: middle; }
.data-table tbody tr:last-child td { border-bottom: none; }
.data-table tbody tr:hover td { background: #fafff9; }

.td-no   { color: #9ca3af; font-size: 12px; width: 40px; }
.td-file { display: flex; align-items: center; gap: 9px; min-width: 180px; max-width: 260px; }
.td-file__name { overflow: hidden; text-overflow: ellipsis; white-space: nowrap; font-weight: 500; color: #111827; font-size: 12.5px; }
.td-pemilik { font-size: 12.5px; }
.td-anon  { color: #9ca3af; font-style: italic; }
.td-size  { color: #6b7280; font-size: 12px; white-space: nowrap; }
.td-date  { color: #6b7280; font-size: 12px; white-space: nowrap; }

.file-icon { width: 26px; height: 26px; border-radius: 7px; display: flex; align-items: center; justify-content: center; flex-shrink: 0; }
.file-icon--pdf { background: #fef2f2; color: #ef4444; }
.file-icon--img { background: #f0fdf4; color: #3b82f6; }
.file-icon--doc { background: #f9fafb; color: #6b7280; }

.jenis-badge { display: inline-block; font-size: 10.5px; font-weight: 700; padding: 3px 9px; border-radius: 100px; white-space: nowrap; }
.jenis-badge--proposal_magang  { background: #f0fdf4; color: #0d2818; }
.jenis-badge--surat_pengantar  { background: #f0fdf4; color: #0d2818; }
.jenis-badge--ktp              { background: #f0fdf4; color: #9333ea; }
.jenis-badge--ktm              { background: #f0fdf4; color: #9333ea; }
.jenis-badge--pasfoto          { background: #f0fdf4; color: #0284c7; }
.jenis-badge--bpjs_kis         { background: #f0fdf4; color: #16a34a; }
.jenis-badge--surat_balasan    { background: #fefce8; color: #ca8a04; }
.jenis-badge--laporan_magang   { background: #fff7ed; color: #ea580c; }
.jenis-badge--sertifikat       { background: #f0fdf4; color: #7c3aed; }

.aksi-cell { display: flex; align-items: center; gap: 6px; flex-wrap: nowrap; }
.btn-aksi { border: none; border-radius: 7px; padding: 5px 10px; font-size: 11px; font-weight: 700; font-family: inherit; cursor: pointer; white-space: nowrap; transition: opacity .15s; display: inline-flex; align-items: center; gap: 4px; }
.btn-aksi--ghost { background: #f3f4f6; color: #374151; }
.btn-aksi--ghost:hover { background: #e5e7eb; }
.btn-aksi--green { background: #dcfce7; color: #16a34a; }
.btn-aksi--green:hover { background: #bbf7d0; }
.btn-aksi--red   { background: #fee2e2; color: #dc2626; }
.btn-aksi--red:hover { background: #fecaca; }

.pagination { display: flex; align-items: center; justify-content: center; gap: 14px; padding: 16px; border-top: 1px solid #f0faf0; }
.pg-btn { background: #f3f4f6; border: none; border-radius: 8px; width: 32px; height: 32px; display: flex; align-items: center; justify-content: center; cursor: pointer; color: #374151; }
.pg-btn:hover:not(:disabled) { background: #e5e7eb; }
.pg-btn:disabled { opacity: .4; cursor: default; }
.pg-info { font-size: 12.5px; color: #6b7280; font-weight: 600; }

/* Spinner */
.spinner { width: 36px; height: 36px; border: 3px solid #e5e7eb; border-top-color: #48AF4A; border-radius: 50%; animation: spin .8s linear infinite; margin: 0 auto; }
.spinner-sm { width: 13px; height: 13px; border: 2px solid rgba(0,0,0,.15); border-top-color: #ef4444; border-radius: 50%; animation: spin .7s linear infinite; display: inline-block; }
@keyframes spin { to { transform: rotate(360deg); } }

/* Toast */
.toast { position: fixed; bottom: 24px; left: 50%; transform: translateX(-50%); padding: 11px 22px; border-radius: 10px; font-size: 13px; font-weight: 600; z-index: 9999; box-shadow: 0 4px 20px rgba(0,0,0,.15); }
.toast--ok  { background: #111827; color: #fff; }
.toast--err { background: #ef4444; color: #fff; }
.toast-fade-enter-active, .toast-fade-leave-active { transition: opacity .3s, transform .3s; }
.toast-fade-enter-from, .toast-fade-leave-to { opacity: 0; transform: translateX(-50%) translateY(8px); }

/* Modal shared */
.modal-backdrop { position: fixed; inset: 0; background: rgba(0,0,0,.5); backdrop-filter: blur(3px); z-index: 300; display: flex; align-items: center; justify-content: center; padding: 16px; }
.modal-fade-enter-active, .modal-fade-leave-active { transition: opacity .2s; }
.modal-fade-enter-from, .modal-fade-leave-to { opacity: 0; }

/* Preview modal */
.modal-box--preview { background: #fff; border-radius: 18px; width: min(900px, 96vw); max-height: 90vh; display: flex; flex-direction: column; box-shadow: 0 24px 80px rgba(0,0,0,.22); overflow: hidden; }
.modal-header { display: flex; align-items: flex-start; justify-content: space-between; padding: 16px 18px; border-bottom: 1px solid #f0f0f0; gap: 12px; flex-shrink: 0; }
.modal-title  { font-size: 14px; font-weight: 700; color: #111827; word-break: break-all; }
.modal-sub    { font-size: 11.5px; color: #9ca3af; margin-top: 2px; }
.modal-close  { background: #f3f4f6; border: none; border-radius: 8px; width: 30px; height: 30px; flex-shrink: 0; display: flex; align-items: center; justify-content: center; cursor: pointer; color: #6b7280; }
.modal-close:hover { background: #e5e7eb; }
.modal-body   { flex: 1; overflow: hidden; background: #f3f4f6; min-height: 400px; }
.preview-loading { display: flex; flex-direction: column; align-items: center; justify-content: center; height: 100%; gap: 14px; color: #6b7280; font-size: 13px; padding: 60px; }
.preview-img    { display: block; max-width: 100%; max-height: 74vh; object-fit: contain; margin: auto; padding: 16px; }
.preview-iframe { width: 100%; height: 74vh; border: none; display: block; }
.preview-other  { display: flex; flex-direction: column; align-items: center; justify-content: center; height: 100%; gap: 10px; padding: 60px; color: #9ca3af; font-size: 13px; }

/* Confirm modal */
.modal-box--confirm { background: #fff; border-radius: 18px; width: min(400px, 94vw); padding: 28px 24px; text-align: center; box-shadow: 0 24px 80px rgba(0,0,0,.22); }
.confirm-icon  { width: 60px; height: 60px; background: #fef2f2; border-radius: 50%; display: flex; align-items: center; justify-content: center; margin: 0 auto 16px; }
.confirm-title { font-size: 16px; font-weight: 700; color: #111827; margin-bottom: 10px; }
.confirm-body  { font-size: 13px; color: #6b7280; line-height: 1.65; margin-bottom: 24px; }
.confirm-actions { display: flex; gap: 10px; }
.btn-cancel { flex: 1; background: #f3f4f6; color: #374151; border: none; border-radius: 9px; padding: 10px 16px; font-size: 13px; font-weight: 600; font-family: inherit; cursor: pointer; }
.btn-cancel:hover:not(:disabled) { background: #e5e7eb; }
.btn-cancel:disabled { opacity: .5; cursor: default; }
.btn-red    { flex: 1; background: #ef4444; color: #fff; border: none; border-radius: 9px; padding: 10px 22px; font-size: 13px; font-weight: 600; font-family: "Poppins",sans-serif; cursor: pointer; display: flex; align-items: center; justify-content: center; gap: 7px; }
.btn-red:hover:not(:disabled) { background: #dc2626; }
.btn-red:disabled { opacity: .5; cursor: not-allowed; }
</style>
