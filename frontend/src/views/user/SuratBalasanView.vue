<template>
  <div>

    <!-- Loading -->
    <div v-if="loading" class="sb-state">
      <div class="spin"></div>
      <span class="sb-state__txt">Memuat dokumen…</span>
    </div>

    <!-- Empty -->
    <div v-else-if="suratList.length === 0" class="sb-empty">
      <div class="sb-empty__ring">
        <svg width="40" height="40" viewBox="0 0 24 24" fill="none">
          <path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"
            stroke="#d1d5db" stroke-width="1.5"/>
          <polyline points="22,6 12,13 2,6" stroke="#d1d5db" stroke-width="1.5"/>
        </svg>
      </div>
      <p class="sb-empty__title">Belum ada surat balasan</p>
      <p class="sb-empty__sub">
        HRD akan mengunggah surat balasan setelah pengajuan magang kamu diverifikasi dan disetujui.
      </p>
      <div v-if="pengajuanStatus" class="sb-status-chip" :class="`chip--${statusCls}`">
        Status pengajuan: <strong>{{ statusLabel }}</strong>
      </div>
    </div>

    <!-- Tabel surat -->
    <div v-else class="sb-table-wrap">
      <div class="sb-table-topbar">
        <span class="sb-table-title">Surat Balasan Magang</span>
        <span class="sb-table-count">{{ suratList.length }} dokumen</span>
      </div>
      <table class="sb-table">
        <thead>
          <tr>
            <th class="col-jenis">Jenis Surat</th>
            <th class="col-nama">Nama File</th>
            <th class="col-ukuran">Ukuran</th>
            <th class="col-tgl">Tanggal</th>
            <th class="col-aksi">Aksi</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="dok in suratList" :key="dok.id" class="sb-row" @click="openPreview(dok)">
            <td class="col-jenis">
              <div class="sb-jenis-cell">
                <div class="sb-file-icon" :class="isPdf(dok) ? 'icon--pdf' : isImage(dok) ? 'icon--img' : 'icon--file'">
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
                <span class="sb-jenis-label">{{ labelFor(dok.jenis) }}</span>
              </div>
            </td>
            <td class="col-nama">
              <span class="sb-nama-txt">{{ dok.nama_file }}</span>
            </td>
            <td class="col-ukuran">{{ fmtSize(dok.ukuran_bytes) }}</td>
            <td class="col-tgl">{{ fmtDate(dok.created_at) }}</td>
            <td class="col-aksi" @click.stop>
              <div class="aksi-wrap">
                <button class="act-btn act-btn--preview" @click="openPreview(dok)">
                  <svg width="13" height="13" viewBox="0 0 24 24" fill="none">
                    <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z" stroke="currentColor" stroke-width="2"/>
                    <circle cx="12" cy="12" r="3" stroke="currentColor" stroke-width="2"/>
                  </svg>
                  Preview
                </button>
                <button class="act-btn act-btn--dl" @click="downloadDok(dok)"
                  :class="{ 'act-btn--loading': dlId === dok.id }">
                  <div v-if="dlId === dok.id" class="spin spin--sm"></div>
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

    <!-- Info -->
    <div class="sb-info">
      <svg width="15" height="15" viewBox="0 0 24 24" fill="none">
        <circle cx="12" cy="12" r="10" stroke="#16a34a" stroke-width="2"/>
        <line x1="12" y1="8" x2="12" y2="12" stroke="#16a34a" stroke-width="2" stroke-linecap="round"/>
        <line x1="12" y1="16" x2="12.01" y2="16" stroke="#16a34a" stroke-width="2" stroke-linecap="round"/>
      </svg>
      <span>Surat balasan diterbitkan oleh HRD dan berfungsi sebagai bukti penerimaan magang resmi dari PT TanjungEnim Lestari Pulp and Paper.</span>
    </div>

    <!-- ── MODAL PREVIEW ── -->
    <Teleport to="body">
      <div v-if="preview.show" class="pv-backdrop" @click.self="closePreview">
        <div class="pv-box">
          <div class="pv-head">
            <div class="pv-head__left">
              <div class="pv-head__jenis">{{ labelFor(preview.jenis) }}</div>
              <div class="pv-head__nama">{{ preview.nama }}</div>
            </div>
            <div class="pv-head__right">
              <button class="pv-act pv-act--dl" @click="downloadDok(preview.docRef!)" title="Download">
                <svg width="15" height="15" viewBox="0 0 24 24" fill="none">
                  <path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
                  <polyline points="7 10 12 15 17 10" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                  <line x1="12" y1="15" x2="12" y2="3" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
                </svg>
                Unduh
              </button>
              <button class="pv-close" @click="closePreview" title="Tutup">
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
import { ref, computed, onMounted } from "vue";
import api from "@/lib/api";

const props = defineProps<{ pelaksanaan?: any }>();

interface Dokumen {
  id: string; jenis: string; nama_file: string;
  mime_type: string; ukuran_bytes: number; created_at: string;
}

const loading  = ref(false);
const dokList  = ref<Dokumen[]>([]);
const dlId     = ref<string | null>(null);
const pengajuanStatus = ref<string>("");

const SURAT_TYPES = ["surat_balasan", "surat_keterangan", "surat_penerimaan"];

const suratList = computed(() =>
  dokList.value.filter(d => SURAT_TYPES.includes(d.jenis))
);

const statusLabel = computed(() => {
  const m: Record<string, string> = {
    aktif: "Aktif", selesai: "Selesai",
    upload_laporan: "Upload Laporan", penilaian: "Penilaian",
    diterima: "Diterima", diproses: "Diproses",
  };
  return m[pengajuanStatus.value] ?? "Menunggu";
});

const statusCls = computed(() => {
  if (["aktif", "diterima"].includes(pengajuanStatus.value)) return "green";
  if (pengajuanStatus.value === "selesai") return "blue";
  return "gray";
});

function labelFor(jenis: string) {
  const m: Record<string, string> = {
    surat_balasan:    "Surat Balasan Magang",
    surat_keterangan: "Surat Keterangan Magang",
    surat_penerimaan: "Surat Penerimaan Peserta",
  };
  return m[jenis] ?? jenis.replace(/_/g, " ");
}

function fmtDate(iso: string) {
  if (!iso) return "–";
  return new Date(iso).toLocaleDateString("id-ID", { day: "numeric", month: "short", year: "numeric" });
}

function fmtSize(bytes: number) {
  if (!bytes) return "–";
  if (bytes >= 1048576) return `${(bytes / 1048576).toFixed(1)} MB`;
  return `${Math.round(bytes / 1024)} KB`;
}

function isImage(d: Dokumen) { return d.mime_type?.startsWith("image/"); }
function isPdf(d: Dokumen)   { return d.mime_type === "application/pdf"; }

const preview = ref({
  show: false, loading: false, error: false,
  url: "", type: "other" as "image" | "pdf" | "other",
  nama: "", jenis: "", docRef: null as Dokumen | null,
});

async function openPreview(dok: Dokumen) {
  preview.value = { show: true, loading: true, error: false, url: "", type: "other", nama: dok.nama_file, jenis: dok.jenis, docRef: dok };
  try {
    const res = await api.get(`/api/dokumen/${dok.id}/download`, { responseType: "blob" });
    const mime = dok.mime_type ?? "application/octet-stream";
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
    const url = URL.createObjectURL(new Blob([res.data], { type: dok.mime_type ?? "application/octet-stream" }));
    const a   = Object.assign(document.createElement("a"), { href: url, download: dok.nama_file });
    document.body.appendChild(a); a.click();
    setTimeout(() => { URL.revokeObjectURL(url); a.remove(); }, 500);
  } catch { /* silent */ }
  finally { dlId.value = null; }
}

async function fetchDokumen() {
  loading.value = true;
  try {
    let pengajuanId = props.pelaksanaan?.pengajuan_id;
    if (!pengajuanId) {
      const pr = await api.get("/api/pengajuan/saya");
      const p  = pr.data?.data ?? pr.data;
      pengajuanId = p?.id;
      pengajuanStatus.value = p?.status ?? "";
    } else {
      pengajuanStatus.value = props.pelaksanaan?.status ?? "";
    }
    if (!pengajuanId) { dokList.value = []; return; }

    const r = await api.get(`/api/dokumen/pengajuan/${pengajuanId}`);
    dokList.value = Array.isArray(r.data?.data) ? r.data.data : [];
  } catch {
    dokList.value = [];
  } finally {
    loading.value = false;
  }
}

onMounted(fetchDokumen);
</script>

<style scoped>
/* ── States ─────────────────────────────────────────────────── */
.sb-state {
  display: flex; flex-direction: column; align-items: center; gap: 10px;
  padding: 56px 24px; background: #fff; border-radius: 14px;
  border: 1px solid #e9f5e9; color: #9ca3af; font-size: 13px;
}
.sb-state__txt { color: #9ca3af; font-size: 13px; }
.sb-empty {
  display: flex; flex-direction: column; align-items: center; gap: 10px;
  padding: 56px 24px; background: #fff; border-radius: 14px;
  border: 1px solid #e9f5e9; text-align: center;
}
.sb-empty__ring {
  width: 80px; height: 80px; border-radius: 50%;
  background: #f9fafb; border: 2px solid #f3f4f6;
  display: flex; align-items: center; justify-content: center; margin-bottom: 6px;
}
.sb-empty__title { font-size: 15px; font-weight: 700; color: #374151; margin: 0; }
.sb-empty__sub   { font-size: 12.5px; color: #9ca3af; max-width: 340px; line-height: 1.65; margin: 0; }
.sb-status-chip {
  margin-top: 4px; font-size: 12px; font-weight: 600;
  padding: 5px 14px; border-radius: 100px;
}
.chip--green { background: #f0fdf4; color: #15803d; border: 1px solid #86efac; }
.chip--blue  { background: #eff6ff; color: #1d4ed8; border: 1px solid #bfdbfe; }
.chip--gray  { background: #f9fafb; color: #6b7280; border: 1px solid #e5e7eb; }

/* ── Table ──────────────────────────────────────────────────── */
.sb-table-wrap {
  background: #fff; border: 1.5px solid #e9f5e9; border-radius: 14px;
  overflow: hidden; box-shadow: 0 1px 3px rgba(13,40,24,.04);
}
.sb-table-topbar {
  display: flex; align-items: center; justify-content: space-between;
  padding: 12px 16px; border-bottom: 1px solid #f0fdf4;
  background: #fafffe;
}
.sb-table-title {
  font-size: 13.5px; font-weight: 700; color: #111827;
}
.sb-table-count {
  font-size: 12px; font-weight: 700; color: #6b7280;
  background: #f3f4f6; border-radius: 100px; padding: 3px 12px;
}
.sb-table {
  width: 100%; border-collapse: collapse; font-size: 13px;
}
.sb-table thead tr {
  border-bottom: 1.5px solid #f0fdf4; background: #fafffe;
}
.sb-table th {
  padding: 11px 14px; text-align: left;
  font-size: 11px; font-weight: 700; color: #9ca3af;
  text-transform: uppercase; letter-spacing: .06em; white-space: nowrap;
}
.sb-row {
  border-bottom: 1px solid #f9fafb; cursor: pointer; transition: background .1s;
}
.sb-row:last-child { border-bottom: none; }
.sb-row:hover { background: #fafffe; }
.sb-table td { padding: 12px 14px; vertical-align: middle; }

.col-jenis { width: 220px; }
.sb-jenis-cell { display: flex; align-items: center; gap: 10px; }
.sb-file-icon {
  width: 30px; height: 30px; border-radius: 7px; flex-shrink: 0;
  display: flex; align-items: center; justify-content: center;
}
.icon--img  { background: #f0fdf4; color: #16a34a; }
.icon--pdf  { background: #fff5f5; color: #ef4444; }
.icon--file { background: #f9fafb; color: #6b7280; }
.sb-jenis-label { font-size: 12.5px; font-weight: 600; color: #374151; }

.col-nama { max-width: 220px; }
.sb-nama-txt {
  font-size: 12.5px; color: #6b7280;
  white-space: nowrap; overflow: hidden; text-overflow: ellipsis;
  display: block; max-width: 200px;
}

.col-ukuran, .col-tgl { font-size: 12px; color: #9ca3af; white-space: nowrap; }

.col-aksi { width: 160px; }
.aksi-wrap { display: flex; gap: 6px; }

.act-btn {
  display: flex; align-items: center; gap: 5px;
  font-size: 11.5px; font-weight: 600; padding: 5px 11px; border-radius: 7px;
  cursor: pointer; font-family: inherit; border: 1.5px solid transparent;
  transition: all .12s; white-space: nowrap;
}
.act-btn--preview { background: #f9fafb; color: #374151; border-color: #e5e7eb; }
.act-btn--preview:hover { background: #f0fdf4; color: #16a34a; border-color: #86efac; }
.act-btn--dl { background: #f0fdf4; color: #16a34a; border-color: #bbf7d0; }
.act-btn--dl:hover { background: #dcfce7; border-color: #86efac; }
.act-btn--loading { opacity: .6; pointer-events: none; }

/* ── Info box ───────────────────────────────────────────────── */
.sb-info {
  display: flex; align-items: flex-start; gap: 10px;
  background: #f0fdf4; border: 1px solid #bbf7d0; border-radius: 10px;
  padding: 12px 16px; margin-top: 16px;
  font-size: 12px; color: #15803d; line-height: 1.65;
}

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
  border-radius: 8px; cursor: pointer; font-family: inherit; border: none;
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
.spin--sm { width: 14px; height: 14px; border-width: 2px; }
@keyframes rot { to { transform: rotate(360deg); } }

@media (max-width: 640px) {
  .col-ukuran, .col-tgl { display: none; }
  .col-jenis { width: auto; }
}
</style>
