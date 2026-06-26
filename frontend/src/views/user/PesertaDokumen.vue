<template>
  <div>

    <!-- Loading -->
    <div v-if="loading" class="state-center">
      <div class="spin"></div>
      <span class="state-txt">Memuat dokumen…</span>
    </div>

    <!-- Empty -->
    <div v-else-if="dokumenList.length === 0" class="state-empty">
      <div class="empty-ring">
        <svg width="40" height="40" viewBox="0 0 24 24" fill="none">
          <path d="M22 19a2 2 0 01-2 2H4a2 2 0 01-2-2V5a2 2 0 012-2h5l2 3h9a2 2 0 012 2z"
            stroke="#d1d5db" stroke-width="1.5" stroke-linejoin="round"/>
        </svg>
      </div>
      <p class="empty-title">Belum ada dokumen</p>
      <p class="empty-sub">Dokumen pendaftaran akan muncul di sini setelah pengajuan magang kamu diproses oleh HRD.</p>
    </div>

    <!-- Tabel dokumen -->
    <div v-else class="dok-table-wrap">
      <div class="dok-table-topbar">
        <span class="dok-table-title">Dokumen Pendaftaran</span>
        <span class="dok-table-count">{{ dokumenList.length }} dokumen</span>
      </div>
      <table class="dok-table">
        <thead>
          <tr>
            <th class="col-jenis">Jenis Dokumen</th>
            <th class="col-nama">Nama File</th>
            <th class="col-ukuran">Ukuran</th>
            <th class="col-tgl">Tanggal Upload</th>
            <th class="col-aksi">Aksi</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="dok in dokumenList" :key="dok.id" class="dok-row" @click="openPreview(dok)">
            <td class="col-jenis">
              <div class="dok-jenis-cell">
                <div class="dok-file-icon" :class="iconClass(dok)">
                  <svg v-if="isPdf(dok)" width="14" height="14" viewBox="0 0 24 24" fill="none">
                    <path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z" stroke="currentColor" stroke-width="2" stroke-linejoin="round"/>
                    <polyline points="14 2 14 8 20 8" stroke="currentColor" stroke-width="2" stroke-linejoin="round"/>
                  </svg>
                  <svg v-else-if="isImage(dok)" width="14" height="14" viewBox="0 0 24 24" fill="none">
                    <rect x="3" y="3" width="18" height="18" rx="2" stroke="currentColor" stroke-width="2"/>
                    <circle cx="8.5" cy="8.5" r="1.5" stroke="currentColor" stroke-width="1.5"/>
                    <polyline points="21 15 16 10 5 21" stroke="currentColor" stroke-width="2" stroke-linejoin="round"/>
                  </svg>
                  <svg v-else width="14" height="14" viewBox="0 0 24 24" fill="none">
                    <path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z" stroke="currentColor" stroke-width="2" stroke-linejoin="round"/>
                    <polyline points="14 2 14 8 20 8" stroke="currentColor" stroke-width="2" stroke-linejoin="round"/>
                  </svg>
                </div>
                <span class="dok-jenis-label">{{ labelJenis(dok.jenis) }}</span>
              </div>
            </td>
            <td class="col-nama">
              <span class="dok-nama-txt">{{ dok.nama_file }}</span>
            </td>
            <td class="col-ukuran">{{ formatUkuran(dok.ukuran_bytes) }}</td>
            <td class="col-tgl">{{ formatTgl(dok.uploaded_at) }}</td>
            <td class="col-aksi" @click.stop>
              <div class="aksi-wrap">
                <button class="btn-aksi btn-aksi--preview" @click="openPreview(dok)" title="Preview">
                  <svg width="13" height="13" viewBox="0 0 24 24" fill="none">
                    <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z" stroke="currentColor" stroke-width="2"/>
                    <circle cx="12" cy="12" r="3" stroke="currentColor" stroke-width="2"/>
                  </svg>
                  Preview
                </button>
                <button class="btn-aksi btn-aksi--dl" @click="downloadDok(dok)" title="Unduh"
                  :class="{ 'btn-aksi--loading': dlId === dok.id }">
                  <div v-if="dlId === dok.id" class="spin spin--xs"></div>
                  <svg v-else width="13" height="13" viewBox="0 0 24 24" fill="none">
                    <path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
                    <polyline points="7 10 12 15 17 10" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                    <line x1="12" y1="15" x2="12" y2="3" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
                  </svg>
                  Unduh
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- ── MODAL PREVIEW ── -->
    <Teleport to="body">
      <div v-if="preview.show" class="pv-backdrop" @click.self="closePreview">
        <div class="pv-box">
          <div class="pv-head">
            <div class="pv-head__left">
              <div class="pv-head__jenis">{{ labelJenis(preview.jenis) }}</div>
              <div class="pv-head__nama">{{ preview.nama }}</div>
            </div>
            <div class="pv-head__right">
              <button class="pv-act pv-act--dl" @click="downloadDok(preview.docRef!)" title="Download">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none">
                  <path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
                  <polyline points="7 10 12 15 17 10" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                  <line x1="12" y1="15" x2="12" y2="3" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
                </svg>
                Unduh
              </button>
              <button class="pv-close" @click="closePreview">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none">
                  <path d="M18 6L6 18M6 6l12 12" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"/>
                </svg>
              </button>
            </div>
          </div>
          <div class="pv-body">
            <div v-if="preview.loading" class="pv-state">
              <div class="spin"></div>
              <span>Memuat file…</span>
            </div>
            <div v-else-if="preview.error" class="pv-state">
              <svg width="36" height="36" viewBox="0 0 24 24" fill="none">
                <circle cx="12" cy="12" r="10" stroke="#9ca3af" stroke-width="1.5"/>
                <path d="M15 9l-6 6M9 9l6 6" stroke="#9ca3af" stroke-width="1.5" stroke-linecap="round"/>
              </svg>
              <p>File tidak dapat ditampilkan.</p>
              <button class="pv-dl-fallback" @click="downloadDok(preview.docRef!)">Download File</button>
            </div>
            <img v-else-if="preview.type === 'image'" :src="preview.url" class="pv-image" alt="preview"/>
            <iframe v-else-if="preview.type === 'pdf'" :src="preview.url" class="pv-iframe"/>
            <div v-else class="pv-state">
              <svg width="44" height="44" viewBox="0 0 24 24" fill="none">
                <path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z" stroke="#d1d5db" stroke-width="1.5"/>
                <polyline points="14 2 14 8 20 8" stroke="#d1d5db" stroke-width="1.5"/>
              </svg>
              <p>Preview tidak tersedia untuk tipe file ini.</p>
              <button class="pv-dl-fallback" @click="downloadDok(preview.docRef!)">Download File</button>
            </div>
          </div>
        </div>
      </div>
    </Teleport>

  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import api from "@/lib/api";

interface Dokumen {
  id: string; jenis: string; nama_file: string;
  mime_type: string | null; ukuran_bytes: number | null; uploaded_at: string;
}

const loading     = ref(false);
const dokumenList = ref<Dokumen[]>([]);
const dlId        = ref<string | null>(null);

const LABEL: Record<string, string> = {
  proposal_magang:      "Proposal Magang",
  surat_pengantar:      "Surat Pengantar Instansi",
  ktp:                  "KTP",
  ktm:                  "KTM / Kartu Mahasiswa",
  pasfoto:              "Pas Foto",
  bpjs_kis:             "BPJS / KIS",
  sertifikat_kesehatan: "Sertifikat Kesehatan",
  transkrip_nilai:      "Transkrip Nilai",
  surat_balasan:        "Surat Balasan HRD",
  laporan_magang:       "Laporan Magang",
  sertifikat:           "Sertifikat Magang",
};

function labelJenis(j: string) {
  return LABEL[j] ?? j.replace(/_/g, " ");
}

function mimeOf(d: Dokumen) {
  if (d.mime_type) return d.mime_type;
  const ext = d.nama_file?.split(".").pop()?.toLowerCase() ?? "";
  if (["jpg","jpeg","png","gif","webp","bmp"].includes(ext)) return "image/jpeg";
  if (ext === "pdf") return "application/pdf";
  return "application/octet-stream";
}

function isImage(d: Dokumen) { return mimeOf(d).startsWith("image/"); }
function isPdf(d: Dokumen)   { return mimeOf(d) === "application/pdf"; }

function iconClass(d: Dokumen) {
  if (isImage(d)) return "icon--img";
  if (isPdf(d))   return "icon--pdf";
  return "icon--file";
}

function formatUkuran(bytes: number | null) {
  if (!bytes) return "–";
  if (bytes >= 1048576) return `${(bytes / 1048576).toFixed(1)} MB`;
  return `${Math.round(bytes / 1024)} KB`;
}

function formatTgl(iso: string) {
  if (!iso) return "–";
  return new Date(iso).toLocaleDateString("id-ID", { day: "numeric", month: "short", year: "numeric" });
}

const preview = ref({
  show: false, loading: false, error: false,
  url: "", type: "other" as "image" | "pdf" | "other",
  nama: "", jenis: "", docRef: null as Dokumen | null,
});

async function openPreview(dok: Dokumen) {
  preview.value = { show: true, loading: true, error: false, url: "", type: "other", nama: dok.nama_file, jenis: dok.jenis, docRef: dok };
  try {
    const res = await api.get(`/api/dokumen/${dok.id}/download`, { responseType: "blob" });
    const mime = mimeOf(dok);
    const url  = URL.createObjectURL(new Blob([res.data], { type: mime }));
    let type: "image" | "pdf" | "other" = "other";
    if (mime.startsWith("image/"))       type = "image";
    else if (mime === "application/pdf") type = "pdf";
    preview.value = { ...preview.value, loading: false, url, type };
  } catch {
    preview.value = { ...preview.value, loading: false, error: true };
  }
}

function closePreview() {
  if (preview.value.url) URL.revokeObjectURL(preview.value.url);
  preview.value = { show: false, loading: false, error: false, url: "", type: "other", nama: "", jenis: "", docRef: null };
}

async function downloadDok(dok: Dokumen) {
  dlId.value = dok.id;
  try {
    const res = await api.get(`/api/dokumen/${dok.id}/download`, { responseType: "blob" });
    const url = URL.createObjectURL(new Blob([res.data], { type: mimeOf(dok) }));
    const a   = Object.assign(document.createElement("a"), { href: url, download: dok.nama_file });
    document.body.appendChild(a); a.click();
    setTimeout(() => { URL.revokeObjectURL(url); a.remove(); }, 500);
  } catch { /* silent */ }
  finally { dlId.value = null; }
}

async function fetchDokumen() {
  loading.value = true;
  try {
    const pr = await api.get("/api/pengajuan/saya");
    const pengajuanId = pr.data?.data?.id ?? pr.data?.id;
    if (!pengajuanId) { dokumenList.value = []; return; }

    const res = await api.get(`/api/dokumen/pengajuan/${pengajuanId}`);
    const all = Array.isArray(res.data?.data) ? res.data.data : [];
    dokumenList.value = all.filter((d: any) => d.jenis !== 'surat_balasan');
  } catch {
    dokumenList.value = [];
  } finally {
    loading.value = false;
  }
}

onMounted(fetchDokumen);
</script>

<style scoped>
/* ── States ─────────────────────────────────────────────────── */
.state-center {
  display: flex; flex-direction: column; align-items: center; gap: 10px;
  padding: 56px 24px; background: #fff; border-radius: 14px;
  border: 1px solid #e9f5e9;
}
.state-txt { font-size: 13px; color: #9ca3af; }
.state-empty {
  display: flex; flex-direction: column; align-items: center; gap: 10px;
  padding: 56px 24px; background: #fff; border-radius: 14px;
  border: 1px solid #e9f5e9; text-align: center;
}
.empty-ring {
  width: 80px; height: 80px; border-radius: 50%;
  background: #f9fafb; border: 2px solid #f3f4f6;
  display: flex; align-items: center; justify-content: center; margin-bottom: 6px;
}
.empty-title { font-size: 15px; font-weight: 700; color: #374151; margin: 0; }
.empty-sub   { font-size: 12.5px; color: #9ca3af; max-width: 340px; line-height: 1.65; margin: 0; }

/* ── Table ──────────────────────────────────────────────────── */
.dok-table-wrap {
  background: #fff; border: 1.5px solid #e9f5e9; border-radius: 14px;
  overflow: hidden; box-shadow: 0 1px 3px rgba(13,40,24,.04);
}
.dok-table-topbar {
  display: flex; align-items: center; justify-content: space-between;
  padding: 12px 16px; border-bottom: 1px solid #f0fdf4;
  background: #fafffe;
}
.dok-table-title {
  font-size: 13.5px; font-weight: 700; color: #111827;
}
.dok-table-count {
  font-size: 12px; font-weight: 700; color: #6b7280;
  background: #f3f4f6; border-radius: 100px; padding: 3px 12px;
}
.dok-table {
  width: 100%; border-collapse: collapse; font-size: 13px;
}
.dok-table thead tr {
  border-bottom: 1.5px solid #f0fdf4; background: #fafffe;
}
.dok-table th {
  padding: 11px 14px; text-align: left;
  font-size: 11px; font-weight: 700; color: #9ca3af;
  text-transform: uppercase; letter-spacing: .06em; white-space: nowrap;
}
.dok-row {
  border-bottom: 1px solid #f9fafb; cursor: pointer;
  transition: background .1s;
}
.dok-row:last-child { border-bottom: none; }
.dok-row:hover { background: #fafffe; }
.dok-table td { padding: 12px 14px; vertical-align: middle; }

.col-jenis { width: 200px; }
.dok-jenis-cell { display: flex; align-items: center; gap: 10px; }
.dok-file-icon {
  width: 30px; height: 30px; border-radius: 7px; flex-shrink: 0;
  display: flex; align-items: center; justify-content: center;
}
.icon--img  { background: #f0fdf4; color: #16a34a; }
.icon--pdf  { background: #fff5f5; color: #ef4444; }
.icon--file { background: #f9fafb; color: #6b7280; }
.dok-jenis-label { font-size: 12.5px; font-weight: 600; color: #374151; }

.col-nama { max-width: 220px; }
.dok-nama-txt {
  font-size: 12.5px; color: #6b7280;
  white-space: nowrap; overflow: hidden; text-overflow: ellipsis;
  display: block; max-width: 200px;
}

.col-ukuran, .col-tgl {
  font-size: 12px; color: #9ca3af; white-space: nowrap;
}

.col-aksi { width: 160px; }
.aksi-wrap { display: flex; gap: 6px; }

.btn-aksi {
  display: flex; align-items: center; gap: 5px;
  font-size: 11.5px; font-weight: 600; padding: 5px 11px; border-radius: 7px;
  cursor: pointer; font-family: inherit; border: 1.5px solid transparent;
  transition: all .12s; white-space: nowrap;
}
.btn-aksi--preview {
  background: #f9fafb; color: #374151; border-color: #e5e7eb;
}
.btn-aksi--preview:hover { background: #f0fdf4; color: #16a34a; border-color: #86efac; }
.btn-aksi--dl {
  background: #f0fdf4; color: #16a34a; border-color: #bbf7d0;
}
.btn-aksi--dl:hover { background: #dcfce7; border-color: #86efac; }
.btn-aksi--loading { opacity: .5; pointer-events: none; }

/* ── Modal Preview ──────────────────────────────────────────── */
.pv-backdrop {
  position: fixed; inset: 0; background: rgba(0,0,0,.6);
  z-index: 9999; display: flex; align-items: center; justify-content: center;
  padding: 16px; backdrop-filter: blur(2px);
}
.pv-box {
  background: #fff; border-radius: 16px; overflow: hidden;
  display: flex; flex-direction: column;
  width: 100%; max-width: 880px; height: 90vh;
  box-shadow: 0 24px 64px rgba(0,0,0,.28);
}
.pv-head {
  display: flex; align-items: center; justify-content: space-between;
  padding: 14px 20px; border-bottom: 1px solid #f3f4f6; gap: 12px; flex-shrink: 0;
}
.pv-head__left { min-width: 0; }
.pv-head__jenis {
  font-size: 10px; font-weight: 700; color: #48AF4A;
  text-transform: uppercase; letter-spacing: .08em; margin-bottom: 2px;
}
.pv-head__nama {
  font-size: 14px; font-weight: 700; color: #111827;
  white-space: nowrap; overflow: hidden; text-overflow: ellipsis;
}
.pv-head__right { display: flex; align-items: center; gap: 8px; flex-shrink: 0; }
.pv-act {
  display: flex; align-items: center; gap: 6px;
  font-size: 12.5px; font-weight: 600; padding: 7px 16px;
  border-radius: 8px; cursor: pointer; font-family: inherit;
}
.pv-act--dl { background: #f0fdf4; color: #16a34a; border: 1.5px solid #86efac; }
.pv-act--dl:hover { background: #dcfce7; }
.pv-close {
  width: 32px; height: 32px; border-radius: 8px;
  background: #f9fafb; border: 1.5px solid #e5e7eb; color: #6b7280;
  display: flex; align-items: center; justify-content: center; cursor: pointer;
  transition: all .13s;
}
.pv-close:hover { background: #fee2e2; border-color: #fca5a5; color: #dc2626; }
.pv-body { flex: 1; display: flex; overflow: hidden; }
.pv-state {
  flex: 1; display: flex; flex-direction: column; align-items: center;
  justify-content: center; gap: 12px; color: #9ca3af; font-size: 13.5px;
}
.pv-state p { font-size: 13px; color: #9ca3af; margin: 0; }
.pv-dl-fallback {
  background: #f0fdf4; color: #16a34a; border: 1.5px solid #86efac;
  border-radius: 8px; padding: 8px 20px; font-size: 13px; font-weight: 600;
  cursor: pointer; font-family: inherit;
}
.pv-dl-fallback:hover { background: #dcfce7; }
.pv-image  { width: 100%; height: 100%; object-fit: contain; padding: 16px; }
.pv-iframe { width: 100%; height: 100%; border: none; }

/* ── Spinner ────────────────────────────────────────────────── */
.spin {
  width: 28px; height: 28px; border: 3px solid #e5e7eb;
  border-top-color: #48AF4A; border-radius: 50%; animation: rot .7s linear infinite;
}
.spin--xs { width: 12px; height: 12px; border-width: 2px; }
@keyframes rot { to { transform: rotate(360deg); } }

@media (max-width: 640px) {
  .col-ukuran, .col-tgl { display: none; }
  .col-jenis { width: auto; }
}
</style>