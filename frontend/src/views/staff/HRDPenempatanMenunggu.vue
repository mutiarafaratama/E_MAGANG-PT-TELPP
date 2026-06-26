<template>
  <div class="card">
    <div class="card-header">
      <div>
        <h3 class="card-title">Menunggu Penempatan</h3>
        <p class="card-subtitle">Pengajuan diterima yang belum memiliki jadwal magang</p>
      </div>
      <div class="card-header-actions">
        <span v-if="!loading" class="count-badge">{{ penerimaanTanpaJadwal.length }} peserta</span>
        <button class="btn-green-sm" @click="fetchData">
          <svg width="13" height="13" viewBox="0 0 24 24" fill="none"><path d="M23 4v6h-6" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/><path d="M1 20v-6h6" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/><path d="M3.51 9a9 9 0 0114.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0020.49 15" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>
          Refresh
        </button>
      </div>
    </div>

    <div v-if="loading" class="empty-state"><div class="spinner"></div></div>
    <div v-else-if="error" class="empty-state">
      <p class="text-red">{{ error }}</p>
      <button class="btn-green-sm" @click="fetchData">Coba lagi</button>
    </div>
    <div v-else-if="penerimaanTanpaJadwal.length === 0" class="empty-state">
      <div class="empty-state__icon">
        <svg width="36" height="36" viewBox="0 0 24 24" fill="none"><polyline points="20 6 9 17 4 12" stroke="#d1d5db" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/></svg>
      </div>
      <p>Semua pengajuan yang diterima sudah dijadwalkan.</p>
    </div>
    <div v-else class="table-wrap">
      <table class="data-table">
        <thead>
          <tr>
            <th>Peserta</th><th>Institusi / Jurusan</th><th>Kategori</th><th>Tanggal Diterima</th><th>Aksi</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="p in penerimaanTanpaJadwal" :key="p.id">
            <td>
              <div class="name-cell">
                <div class="name-avatar">{{ p.nama_lengkap[0] }}</div>
                <div>
                  <div class="name-text">{{ p.nama_lengkap }}</div>
                  <div class="name-sub">{{ p.no_hp }}</div>
                </div>
              </div>
            </td>
            <td>
              <div class="name-text" style="font-size:12.5px">{{ p.asal_institusi }}</div>
              <div class="name-sub">{{ p.jurusan }}</div>
            </td>
            <td><span class="tag">{{ formatKategori(p.kategori_magang) }}</span></td>
            <td>{{ formatDate(p.created_at) }}</td>
            <td>
              <button class="btn-set-jadwal" @click="openJadwalModal(p)">
                <svg width="13" height="13" viewBox="0 0 24 24" fill="none"><rect x="3" y="4" width="18" height="18" rx="2" stroke="currentColor" stroke-width="2"/><line x1="3" y1="10" x2="21" y2="10" stroke="currentColor" stroke-width="2"/><line x1="8" y1="2" x2="8" y2="6" stroke="currentColor" stroke-width="2" stroke-linecap="round"/><line x1="16" y1="2" x2="16" y2="6" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>
                Set Jadwal
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>

  <!-- ── MODAL SET JADWAL ──────────────────────────────────── -->
  <Teleport to="body">
    <div v-if="showModal" class="modal-overlay" @click.self="closeModal">
      <div class="jadwal-modal">
        <div class="jadwal-modal__header">
          <div>
            <h3 class="jadwal-modal__title">Set Jadwal Magang</h3>
            <p class="jadwal-modal__sub" v-if="modalTarget">{{ modalTarget.nama_lengkap }} · {{ modalTarget.asal_institusi }}</p>
          </div>
          <button class="drawer-close" @click="closeModal">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none"><path d="M18 6L6 18M6 6l12 12" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>
          </button>
        </div>
        <div class="jadwal-modal__body">
          <div class="jform-row">
            <div class="jform-group">
              <label class="jform-label">Tanggal Mulai <span class="jform-req">*</span></label>
              <input v-model="form.tanggal_mulai" type="date" class="jform-input" :disabled="modalLoading"/>
            </div>
            <div class="jform-group">
              <label class="jform-label">Tanggal Selesai <span class="jform-req">*</span></label>
              <input v-model="form.tanggal_selesai" type="date" class="jform-input" :disabled="modalLoading" :min="form.tanggal_mulai"/>
            </div>
          </div>
          <div v-if="form.tanggal_mulai && form.tanggal_selesai" class="jform-duration">
            <svg width="13" height="13" viewBox="0 0 24 24" fill="none"><circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2"/><polyline points="12 6 12 12 16 14" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>
            Durasi: <strong>{{ Math.max(0, Math.round((new Date(form.tanggal_selesai).getTime() - new Date(form.tanggal_mulai).getTime()) / 86400000)) }} hari</strong>
          </div>
          <div class="jform-group">
            <label class="jform-label">Divisi / Unit Kerja <span class="jform-req">*</span></label>
            <select v-model="form.divisi" class="jform-input" :disabled="modalLoading || divisiLoading">
              <option value="" disabled>{{ divisiLoading ? 'Memuat divisi…' : '— Pilih divisi —' }}</option>
              <option v-for="d in divisiOptions" :key="d.id" :value="d.nama">{{ d.nama }}</option>
            </select>
          </div>
          <div class="jform-group">
            <label class="jform-label">Pembimbing <span class="jform-opt">(opsional)</span></label>
            <input v-model="form.pembimbing" type="text" class="jform-input" :disabled="modalLoading"
              placeholder="Nama pembimbing, kosongkan jika belum ada" />
          </div>
          <div v-if="modalError" class="jform-error">{{ modalError }}</div>
          <div v-if="modalSuccess" class="jform-success">
            <svg width="15" height="15" viewBox="0 0 24 24" fill="none"><polyline points="20 6 9 17 4 12" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"/></svg>
            Jadwal berhasil disimpan!
          </div>
        </div>
        <div class="jadwal-modal__footer">
          <button class="btn-cancel-modal" @click="closeModal" :disabled="modalLoading">Batal</button>
          <button class="btn-jadwal-simpan"
            :disabled="modalLoading || !form.tanggal_mulai || !form.tanggal_selesai || !form.divisi"
            @click="submitJadwal">
            <div v-if="modalLoading" class="spinner spinner--sm spinner--white"></div>
            <svg v-else width="14" height="14" viewBox="0 0 24 24" fill="none"><polyline points="20 6 9 17 4 12" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"/></svg>
            {{ modalLoading ? 'Menyimpan…' : 'Simpan Jadwal' }}
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import api from "@/lib/api";

interface Pengajuan {
  id: string; nama_lengkap: string; no_hp: string;
  asal_institusi: string; jurusan: string; kategori_magang: string;
  status: string; created_at: string;
}
interface Pelaksanaan { id: string; pengajuan_id: string; }

const allPengajuan  = ref<Pengajuan[]>([]);
const allPelaksanaan = ref<Pelaksanaan[]>([]);
const loading       = ref(false);
const error         = ref<string | null>(null);

// divisi dropdown
interface DivisiOption { id: string; nama: string; }
const divisiOptions = ref<DivisiOption[]>([]);
const divisiLoading = ref(false);


async function fetchDivisi() {
  divisiLoading.value = true;
  try {
    const r = await api.get('/api/divisi');
    divisiOptions.value = Array.isArray(r.data?.data) ? r.data.data : [];
  } catch { divisiOptions.value = []; }
  finally { divisiLoading.value = false; }
}

const penerimaanTanpaJadwal = computed(() => {
  const scheduled = new Set(allPelaksanaan.value.map(p => p.pengajuan_id));
  return allPengajuan.value.filter(p => p.status === "diterima" && !scheduled.has(p.id));
});

// modal
const showModal    = ref(false);
const modalTarget  = ref<Pengajuan | null>(null);
const form         = ref({ tanggal_mulai: "", tanggal_selesai: "", divisi: "", pembimbing: "" });
const modalLoading = ref(false);
const modalError   = ref<string | null>(null);
const modalSuccess = ref(false);

async function fetchData() {
  loading.value = true; error.value = null;
  try {
    const [rPengajuan, rPelaksanaan] = await Promise.all([
      api.get("/api/pengajuan?page=1&limit=200"),
      api.get("/api/pelaksanaan?page=1&limit=200"),
    ]);
    allPengajuan.value  = Array.isArray(rPengajuan.data.data)  ? rPengajuan.data.data  : [];
    allPelaksanaan.value = Array.isArray(rPelaksanaan.data.data?.list) ? rPelaksanaan.data.data.list : [];
  } catch (e: any) {
    error.value = e.response?.data?.message ?? "Gagal memuat data.";
  } finally { loading.value = false; }
}

function openJadwalModal(p: Pengajuan) {
  modalTarget.value = p;
  form.value = { tanggal_mulai: "", tanggal_selesai: "", divisi: "", pembimbing: "" };
  modalError.value = null; modalSuccess.value = false; showModal.value = true;
  if (divisiOptions.value.length === 0) fetchDivisi();
}
function closeModal() { showModal.value = false; modalTarget.value = null; modalError.value = null; }

async function submitJadwal() {
  if (!modalTarget.value) return;
  modalLoading.value = true; modalError.value = null; modalSuccess.value = false;
  try {
    await api.post(`/api/pelaksanaan/pengajuan/${modalTarget.value.id}`, {
      tanggal_mulai:   form.value.tanggal_mulai,
      tanggal_selesai: form.value.tanggal_selesai,
      divisi:          form.value.divisi,
      pembimbing:      form.value.pembimbing.trim() || undefined,
    });
    modalSuccess.value = true;
    const [rPelaksanaan, rPengajuan] = await Promise.all([
      api.get("/api/pelaksanaan?page=1&limit=200"),
      api.get("/api/pengajuan?page=1&limit=200"),
    ]);
    allPelaksanaan.value = Array.isArray(rPelaksanaan.data.data?.list) ? rPelaksanaan.data.data.list : [];
    allPengajuan.value   = Array.isArray(rPengajuan.data.data)         ? rPengajuan.data.data        : [];
    setTimeout(() => closeModal(), 1200);
  } catch (e: any) {
    modalError.value = e.response?.data?.message ?? "Gagal menyimpan jadwal.";
  } finally { modalLoading.value = false; }
}

function formatDate(d: string | null) {
  if (!d) return "-";
  return new Date(d).toLocaleDateString("id-ID", { day:"2-digit", month:"short", year:"numeric" });
}
function formatKategori(k: string) {
  return ({
    siswa_smk:"Siswa SMK", mahasiswa_d3:"Mahasiswa D3", mahasiswa_s1:"Mahasiswa S1",
    smk:"SMK / Sederajat", d3_s1_s2:"D3 / S1 / S2",
    d3:"D3", s1:"S1", s2:"S2",
    penelitian:"Penelitian", lainnya:"Lainnya",
  } as Record<string,string>)[k] ?? k;
}

onMounted(fetchData);
</script>

<style scoped>
.card { background:#fff; border-radius:14px; border:1px solid #e9f5e9; box-shadow:0 1px 3px rgba(13,40,24,0.05); overflow:hidden; }
.card-header { display:flex; align-items:center; justify-content:space-between; padding:16px 20px; border-bottom:1px solid #f0faf0; gap:12px; flex-wrap:wrap; }
.card-title { font-size:13.5px; font-weight:700; color:#111827; margin:0; }
.card-subtitle { font-size:12px; color:#9ca3af; margin:2px 0 0; }
.card-header-actions { display:flex; align-items:center; gap:8px; }
.count-badge { background:#f0fdf4; border:1px solid #bbf7d0; color:#16a34a; font-size:11px; font-weight:700; padding:4px 12px; border-radius:100px; }
.btn-green-sm { background:#48AF4A; color:#fff; border:none; border-radius:8px; padding:6px 14px; font-size:12px; font-weight:600; font-family:inherit; cursor:pointer; white-space:nowrap; display:flex; align-items:center; gap:5px; }
.btn-green-sm:hover { background:#3d9e3f; }
.table-wrap { overflow-x:auto; }
.data-table { width:100%; border-collapse:collapse; font-size:13px; }
.data-table th { padding:11px 16px; text-align:left; font-size:10.5px; font-weight:600; color:#6b7280; background:#f9fafb; border-bottom:1px solid #f1f5f9; text-transform:uppercase; letter-spacing:0.04em; white-space:nowrap; }
.data-table td { padding:13px 16px; border-bottom:1px solid #f9fafb; color:#374151; vertical-align:middle; }
.name-cell { display:flex; align-items:center; gap:10px; }
.name-avatar { width:32px; height:32px; border-radius:8px; background:linear-gradient(135deg,#48AF4A,#2d8f30); color:#fff; font-size:13px; font-weight:700; display:flex; align-items:center; justify-content:center; flex-shrink:0; }
.name-text { font-weight:600; color:#111827; font-size:12.5px; }
.name-sub { font-size:11px; color:#9ca3af; }
.tag { background:#eff6ff; color:#1d4ed8; border-radius:6px; padding:2px 8px; font-size:11px; font-weight:600; white-space:nowrap; }
.empty-state { display:flex; flex-direction:column; align-items:center; padding:44px 24px; gap:12px; text-align:center; }
.empty-state__icon { width:72px; height:72px; background:#f9fafb; border-radius:50%; display:flex; align-items:center; justify-content:center; }
.empty-state p { font-size:13px; color:#9ca3af; line-height:1.7; margin:0; }
.text-red { color:#dc2626; font-size:13px; }
.spinner { width:24px; height:24px; border:2.5px solid #e5e7eb; border-top-color:#48AF4A; border-radius:50%; animation:spin 0.7s linear infinite; }
.spinner--sm { width:14px; height:14px; border-width:2px; }
.spinner--white { border-color:rgba(255,255,255,0.3); border-top-color:#fff; }
@keyframes spin { to { transform:rotate(360deg); } }
.btn-set-jadwal { display:inline-flex; align-items:center; gap:5px; background:#0d2818; color:#fff; border:none; border-radius:8px; padding:6px 13px; font-size:12px; font-weight:600; font-family:inherit; cursor:pointer; white-space:nowrap; transition:background 0.15s; }
.btn-set-jadwal:hover { background:#1a5c20; }
/* modal */
.modal-overlay { position:fixed; inset:0; background:rgba(0,0,0,0.5); backdrop-filter:blur(3px); z-index:200; display:flex; align-items:center; justify-content:center; padding:20px; }
.jadwal-modal { background:#fff; border-radius:18px; width:min(520px,100%); box-shadow:0 24px 80px rgba(0,0,0,0.22); display:flex; flex-direction:column; overflow:hidden; }
.jadwal-modal__header { display:flex; align-items:flex-start; justify-content:space-between; padding:22px 24px 16px; border-bottom:1px solid #f0faf0; gap:12px; }
.jadwal-modal__title { font-size:16px; font-weight:700; color:#111827; margin:0 0 3px; }
.jadwal-modal__sub { font-size:12.5px; color:#6b7280; margin:0; }
.jadwal-modal__body { padding:20px 24px; display:flex; flex-direction:column; gap:16px; }
.jadwal-modal__footer { display:flex; gap:10px; padding:16px 24px; border-top:1px solid #f0faf0; }
.drawer-close { background:#f3f4f6; border:none; border-radius:8px; width:32px; height:32px; display:flex; align-items:center; justify-content:center; cursor:pointer; color:#6b7280; flex-shrink:0; }
.drawer-close:hover { background:#e5e7eb; color:#111827; }
.jform-row { display:grid; grid-template-columns:1fr 1fr; gap:12px; }
.jform-group { display:flex; flex-direction:column; gap:5px; }
.jform-label { font-size:12px; font-weight:600; color:#374151; }
.jform-req { color:#dc2626; }
.jform-opt { color:#9ca3af; font-weight:400; }
.jform-input { border:1px solid #e5e7eb; border-radius:9px; padding:9px 12px; font-size:13px; font-family:inherit; color:#111827; outline:none; transition:border-color 0.15s; background:#fff; }
.jform-input:focus { border-color:#48AF4A; box-shadow:0 0 0 3px rgba(72,175,74,0.12); }
.jform-input:disabled { background:#f9fafb; color:#9ca3af; cursor:not-allowed; }
.jform-duration { display:flex; align-items:center; gap:6px; font-size:12.5px; color:#374151; background:#f0fdf4; border:1px solid #bbf7d0; border-radius:8px; padding:8px 12px; }
.jform-error { font-size:12.5px; color:#dc2626; background:#fef2f2; border:1px solid #fecaca; border-radius:8px; padding:9px 13px; }
.jform-success { display:flex; align-items:center; gap:7px; font-size:12.5px; color:#16a34a; background:#f0fdf4; border:1px solid #bbf7d0; border-radius:8px; padding:9px 13px; font-weight:600; }
.btn-cancel-modal { flex:1; background:#f3f4f6; color:#374151; border:none; border-radius:9px; padding:10px 16px; font-size:13px; font-weight:600; font-family:inherit; cursor:pointer; }
.btn-cancel-modal:hover:not(:disabled) { background:#e5e7eb; }
.btn-cancel-modal:disabled { opacity:0.5; cursor:default; }
.btn-jadwal-simpan { flex:1; background:#48AF4A; color:#fff; border:none; border-radius:10px; padding:11px 16px; font-size:13px; font-weight:600; font-family:inherit; cursor:pointer; display:flex; align-items:center; justify-content:center; gap:7px; transition:background 0.15s; }
.btn-jadwal-simpan:hover:not(:disabled) { background:#3d9e3f; }
.btn-jadwal-simpan:disabled { opacity:0.5; cursor:not-allowed; }
</style>
