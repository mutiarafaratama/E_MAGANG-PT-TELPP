<template>
  <div class="card">
    <div class="card-header">
      <div>
        <h3 class="card-title">Sedang Berlangsung</h3>
        <p class="card-subtitle">Peserta yang sudah memiliki jadwal magang aktif</p>
      </div>
      <div class="card-header-actions">
        <span v-if="!loading" class="count-badge">{{ rows.length }} peserta</span>
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
    <div v-else-if="rows.length === 0" class="empty-state">
      <div class="empty-state__icon">
        <svg width="36" height="36" viewBox="0 0 24 24" fill="none">
          <path d="M17 21v-2a4 4 0 00-4-4H5a4 4 0 00-4 4v2" stroke="#d1d5db" stroke-width="1.5"/>
          <circle cx="9" cy="7" r="4" stroke="#d1d5db" stroke-width="1.5"/>
          <path d="M23 21v-2a4 4 0 00-3-3.87" stroke="#d1d5db" stroke-width="1.5" stroke-linecap="round"/>
          <path d="M16 3.13a4 4 0 010 7.75" stroke="#d1d5db" stroke-width="1.5" stroke-linecap="round"/>
        </svg>
      </div>
      <p>Belum ada peserta dengan jadwal magang aktif.</p>
    </div>
    <div v-else class="table-wrap">
      <table class="data-table">
        <thead>
          <tr>
            <th>Peserta</th><th>Divisi</th><th>Periode</th>
            <th>Status</th><th>Sisa Hari</th><th>Aksi</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="pl in rows" :key="pl.id">
            <td>
              <div class="name-cell">
                <div class="name-avatar">{{ pl.nama_lengkap[0] }}</div>
                <div>
                  <div class="name-text">{{ pl.nama_lengkap }}</div>
                  <div class="name-sub">{{ pl.asal_institusi }}</div>
                  <div class="name-sub">{{ formatKategori(pl.kategori_magang) }}</div>
                </div>
              </div>
            </td>
            <td>
              <span v-if="pl.divisi" class="tag">{{ pl.divisi }}</span>
              <span v-else class="name-sub">–</span>
            </td>
            <td>
              <div class="name-text" style="font-size:12px;white-space:nowrap">{{ formatDate(pl.tanggal_mulai) }}</div>
              <div class="name-sub" style="white-space:nowrap">s/d {{ formatDate(pl.tanggal_selesai) }}</div>
            </td>
            <td><span :class="statusClass(pl.status)">{{ formatStatus(pl.status) }}</span></td>
            <td>
              <span v-if="sisaHari(pl.tanggal_selesai) < 0" class="sisa-hari sisa-hari--over">Lewat {{ Math.abs(sisaHari(pl.tanggal_selesai)) }}h</span>
              <span v-else-if="sisaHari(pl.tanggal_selesai) <= 7" class="sisa-hari sisa-hari--warn">{{ sisaHari(pl.tanggal_selesai) }} hari</span>
              <span v-else class="sisa-hari">{{ sisaHari(pl.tanggal_selesai) }} hari</span>
            </td>
            <td>
              <div class="aksi-cell">
                <button
                  v-if="pl.status === 'menunggu_mulai'"
                  class="btn-aksi btn-aksi--green"
                  :disabled="updatingId === pl.id"
                  @click="updateStatus(pl.id, 'aktif')"
                >{{ updatingId === pl.id ? '...' : 'Aktifkan' }}</button>
                <button
                  v-else-if="pl.status === 'upload_laporan'"
                  class="btn-aksi btn-aksi--orange"
                  :disabled="updatingId === pl.id"
                  @click="updateStatus(pl.id, 'penilaian')"
                >{{ updatingId === pl.id ? '...' : 'Ke Penilaian' }}</button>
                <button class="btn-aksi btn-aksi--ghost" @click="openDetail(pl)">Detail</button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <div v-if="updateError" class="update-error">{{ updateError }}</div>
  </div>

  <!-- ── Drawer: Detail Peserta ─────────────────────────────── -->
  <Teleport to="body">
    <div v-if="showDetail" class="drawer-overlay" @click.self="showDetail = false">
      <div class="drawer">
        <div class="drawer-header">
          <div>
            <h2 class="drawer-title">Detail Peserta</h2>
            <p class="drawer-sub" v-if="selected">{{ selected.nama_lengkap }}</p>
          </div>
          <button class="drawer-close" @click="showDetail = false">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none"><path d="M18 6L6 18M6 6l12 12" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>
          </button>
        </div>

        <div v-if="detailLoading" class="drawer-body drawer-loading"><div class="spinner"></div></div>

        <div v-else-if="selected" class="drawer-body">

          <!-- Status banner -->
          <div class="detail-status-banner" :class="`status-banner--${selected.status}`">
            <span :class="statusClass(selected.status)">{{ formatStatus(selected.status) }}</span>
            <span class="banner-sisa">
              <span v-if="sisaHari(selected.tanggal_selesai) < 0">Lewat {{ Math.abs(sisaHari(selected.tanggal_selesai)) }} hari</span>
              <span v-else>{{ sisaHari(selected.tanggal_selesai) }} hari lagi</span>
            </span>
          </div>

          <!-- Data Pribadi -->
          <div v-if="pengajuanDetail" class="detail-section">
            <div class="detail-section__title">Data Pribadi</div>
            <div class="info-grid">
              <div class="info-item info-item--full">
                <span class="info-label">Nama Lengkap</span>
                <span class="info-val">{{ pengajuanDetail.nama_lengkap }}</span>
              </div>
              <div class="info-item">
                <span class="info-label">Email</span>
                <span class="info-val">{{ pengajuanDetail.email ?? '–' }}</span>
              </div>
              <div class="info-item">
                <span class="info-label">No. HP</span>
                <span class="info-val">{{ pengajuanDetail.no_hp ?? '–' }}</span>
              </div>
              <div class="info-item">
                <span class="info-label">Jenis Kelamin</span>
                <span class="info-val">{{ formatJK(pengajuanDetail.jenis_kelamin) }}</span>
              </div>
              <div class="info-item">
                <span class="info-label">Tempat Lahir</span>
                <span class="info-val">{{ pengajuanDetail.tempat_lahir ?? '–' }}</span>
              </div>
              <div class="info-item">
                <span class="info-label">Tanggal Lahir</span>
                <span class="info-val">{{ formatDate(pengajuanDetail.tanggal_lahir) }}</span>
              </div>
              <div class="info-item info-item--full">
                <span class="info-label">Alamat</span>
                <span class="info-val">{{ pengajuanDetail.alamat ?? '–' }}</span>
              </div>
            </div>
          </div>

          <!-- Data Akademik -->
          <div v-if="pengajuanDetail" class="detail-section">
            <div class="detail-section__title">Data Akademik</div>
            <div class="info-grid">
              <div class="info-item">
                <span class="info-label">Asal Institusi</span>
                <span class="info-val">{{ pengajuanDetail.asal_institusi ?? '–' }}</span>
              </div>
              <div class="info-item">
                <span class="info-label">Jurusan</span>
                <span class="info-val">{{ pengajuanDetail.jurusan ?? '–' }}</span>
              </div>
              <div class="info-item">
                <span class="info-label">Kelas / Semester</span>
                <span class="info-val">{{ pengajuanDetail.kelas_semester ?? '–' }}</span>
              </div>
              <div class="info-item">
                <span class="info-label">Nomor Induk</span>
                <span class="info-val">{{ pengajuanDetail.nomor_induk ?? '–' }}</span>
              </div>
              <div class="info-item">
                <span class="info-label">Kategori Magang</span>
                <span class="info-val">{{ formatKategori(pengajuanDetail.kategori_magang) }}</span>
              </div>
              <div class="info-item">
                <span class="info-label">Tanggal Daftar</span>
                <span class="info-val">{{ formatDate(pengajuanDetail.created_at) }}</span>
              </div>
            </div>
          </div>

          <!-- Informasi Magang -->
          <div class="detail-section">
            <div class="detail-section__title">Informasi Magang</div>
            <div class="info-grid">
              <div class="info-item info-item--full">
                <span class="info-label">Divisi / Unit Kerja</span>
                <span class="info-val">{{ selected.divisi ?? '–' }}</span>
              </div>
              <div class="info-item">
                <span class="info-label">Tanggal Mulai</span>
                <span class="info-val">{{ formatDate(selected.tanggal_mulai) }}</span>
              </div>
              <div class="info-item">
                <span class="info-label">Tanggal Selesai</span>
                <span class="info-val">{{ formatDate(selected.tanggal_selesai) }}</span>
              </div>
              <div class="info-item">
                <span class="info-label">Status</span>
                <span class="info-val"><span :class="statusClass(selected.status)">{{ formatStatus(selected.status) }}</span></span>
              </div>
              <div class="info-item">
                <span class="info-label">Sisa Hari</span>
                <span class="info-val">
                  <span v-if="sisaHari(selected.tanggal_selesai) < 0" class="sisa-hari sisa-hari--over">Lewat {{ Math.abs(sisaHari(selected.tanggal_selesai)) }} hari</span>
                  <span v-else-if="sisaHari(selected.tanggal_selesai) <= 7" class="sisa-hari sisa-hari--warn">{{ sisaHari(selected.tanggal_selesai) }} hari</span>
                  <span v-else class="sisa-hari">{{ sisaHari(selected.tanggal_selesai) }} hari</span>
                </span>
              </div>
            </div>
          </div>

          <!-- Pembimbing -->
          <div class="detail-section">
            <div class="detail-section__title-row">
              <span class="detail-section__title">Pembimbing Magang</span>
              <button v-if="!editPembimbing" class="edit-inline-btn" @click="startEditPembimbing">Ganti</button>
            </div>

            <!-- Mode Edit -->
            <div v-if="editPembimbing" class="pembimbing-edit">
              <input v-model="editPembimbingText" type="text" class="pembimbing-input"
                placeholder="Nama pembimbing, kosongkan untuk menghapus"
                :disabled="savingPembimbing" @keyup.enter="savePembimbing" />
              <div v-if="pembimbingError" class="pembimbing-err">{{ pembimbingError }}</div>
              <div class="pembimbing-edit-btns">
                <button class="btn-cancel-sm" @click="editPembimbing = false; pembimbingError = ''"
                  :disabled="savingPembimbing">Batal</button>
                <button class="btn-save-sm" @click="savePembimbing" :disabled="savingPembimbing">
                  {{ savingPembimbing ? 'Menyimpan…' : 'Simpan' }}
                </button>
              </div>
            </div>

            <!-- Mode Tampil -->
            <template v-else>
              <div v-if="selected.pembimbing_nama" class="pembimbing-box">
                <div class="pembimbing-avatar">{{ selected.pembimbing_nama[0] }}</div>
                <div>
                  <div class="pembimbing-name">{{ selected.pembimbing_nama }}</div>
                  <div class="pembimbing-role">Pembimbing Magang</div>
                </div>
              </div>
              <div v-else class="pembimbing-empty">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none"><path d="M20 21v-2a4 4 0 00-4-4H8a4 4 0 00-4 4v2" stroke="#d1d5db" stroke-width="1.5"/><circle cx="12" cy="7" r="4" stroke="#d1d5db" stroke-width="1.5"/></svg>
                Belum ditugaskan
              </div>
            </template>
          </div>

          <!-- Laporan Magang (tampil jika status upload_laporan atau penilaian) -->
          <div v-if="selected && ['upload_laporan','penilaian','selesai'].includes(selected.status)" class="detail-section laporan-section">
            <div class="detail-section__title-row">
              <span class="detail-section__title">Laporan Magang</span>
              <span v-if="laporanList.length" class="laporan-count">{{ laporanList.length }} file</span>
            </div>

            <div v-if="laporanLoading" class="laporan-loading"><div class="spinner spinner--sm"></div></div>

            <div v-else-if="laporanList.length === 0" class="laporan-empty">
              <svg width="26" height="26" viewBox="0 0 24 24" fill="none"><path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z" stroke="#d1d5db" stroke-width="1.5"/><polyline points="14 2 14 8 20 8" stroke="#d1d5db" stroke-width="1.5"/></svg>
              Peserta belum mengunggah laporan.
            </div>

            <div v-else class="laporan-list">
              <div v-for="lap in laporanList" :key="lap.id" class="laporan-item" :class="{ 'laporan-item--latest': lap.id === laporanList[0].id }">
                <div class="laporan-item__icon">
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none"><path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z" stroke="currentColor" stroke-width="2"/><polyline points="14 2 14 8 20 8" stroke="currentColor" stroke-width="2"/></svg>
                </div>
                <div class="laporan-item__info">
                  <div class="laporan-item__name">{{ lap.nama_file }}</div>
                  <div class="laporan-item__meta">
                    <span class="lap-versi">v{{ lap.versi }}</span>
                    <span v-if="lap.ukuran_bytes">{{ fmtSize(lap.ukuran_bytes) }}</span>
                    <span>{{ formatDate(lap.diupload_at) }}</span>
                  </div>
                </div>
                <div class="laporan-item__right">
                  <span class="lap-status" :class="`lap-status--${lap.status}`">{{ formatStatusLaporan(lap.status) }}</span>
                  <a :href="`/api/laporan/${lap.id}/download`" target="_blank" class="btn-dl-lap" title="Unduh">
                    <svg width="13" height="13" viewBox="0 0 24 24" fill="none"><path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4" stroke="currentColor" stroke-width="2" stroke-linecap="round"/><polyline points="7 10 12 15 17 10" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/><line x1="12" y1="15" x2="12" y2="3" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>
                  </a>
                </div>
              </div>
            </div>

            <!-- Action review (hanya untuk laporan terbaru yang menunggu) -->
            <div v-if="latestLaporan && latestLaporan.status === 'menunggu_review'" class="review-action-box">
              <div class="review-action-box__label">Tindakan Review</div>
              <textarea v-model="reviewCatatan" class="review-catatan-input" rows="2"
                placeholder="Catatan untuk peserta (opsional, wajib jika revisi)…"></textarea>
              <div v-if="reviewError" class="review-err">{{ reviewError }}</div>
              <div class="review-btns">
                <button class="btn-revisi" :disabled="reviewing" @click="submitReview('revisi')">
                  <div v-if="reviewing && reviewAction === 'revisi'" class="spinner spinner--sm spinner--white"></div>
                  <svg v-else width="14" height="14" viewBox="0 0 24 24" fill="none"><path d="M23 4v6h-6" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/><path d="M20.49 15a9 9 0 11-2.12-9.36L23 10" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>
                  Minta Revisi
                </button>
                <button class="btn-acc" :disabled="reviewing" @click="submitReview('disetujui')">
                  <div v-if="reviewing && reviewAction === 'disetujui'" class="spinner spinner--sm spinner--white"></div>
                  <svg v-else width="14" height="14" viewBox="0 0 24 24" fill="none"><polyline points="20 6 9 17 4 12" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>
                  Setujui Laporan
                </button>
              </div>
            </div>

            <!-- Sudah disetujui info -->
            <div v-else-if="latestLaporan && latestLaporan.status === 'disetujui'" class="review-done-box">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none"><circle cx="12" cy="12" r="10" fill="#16a34a"/><path d="M8 12l3 3 5-5" stroke="#fff" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>
              Laporan disetujui — form penilaian bisa dibuka
            </div>
          </div>

          <!-- Nilai (jika sudah ada) -->
          <div v-if="selected.nilai !== null && selected.nilai !== undefined" class="detail-section">
            <div class="detail-section__title">Penilaian Akhir</div>
            <div class="nilai-display">
              <div class="nilai-score">{{ selected.nilai }}</div>
              <div class="nilai-label">/ 100</div>
            </div>
          </div>

        </div>

        <!-- Footer: workflow action -->
        <div v-if="selected && ['menunggu_mulai','upload_laporan'].includes(selected.status)" class="drawer-footer">
          <button
            v-if="selected.status === 'menunggu_mulai'"
            class="btn-full btn-full--green"
            :disabled="updatingId === selected.id"
            @click="updateStatus(selected.id, 'aktif')"
          >{{ updatingId === selected.id ? 'Memproses...' : 'Aktifkan Magang' }}</button>
          <button
            v-else-if="selected.status === 'upload_laporan'"
            class="btn-full btn-full--orange"
            :disabled="updatingId === selected.id"
            @click="updateStatus(selected.id, 'penilaian')"
          >{{ updatingId === selected.id ? 'Memproses...' : 'Pindah ke Tahap Penilaian' }}</button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import api from "@/lib/api";

interface Pengajuan {
  id: string;
  nama_lengkap: string;
  asal_institusi: string;
  jurusan: string;
  kategori_magang: string;
}
interface Pelaksanaan {
  id: string;
  pengajuan_id: string;
  tanggal_mulai: string;
  tanggal_selesai: string;
  divisi: string | null;
  pembimbing_id: string | null;
  pembimbing_nama: string | null;
  status: string;
  nilai: number | null;
}
interface Row extends Pelaksanaan {
  nama_lengkap: string;
  asal_institusi: string;
  jurusan: string;
  kategori_magang: string;
}
interface PengajuanDetail {
  id: string;
  nama_lengkap: string;
  email: string;
  no_hp: string;
  jenis_kelamin: string;
  tempat_lahir: string;
  tanggal_lahir: string;
  alamat: string;
  asal_institusi: string;
  jurusan: string;
  kelas_semester: string;
  nomor_induk: string;
  kategori_magang: string;
  created_at: string;
}
const allPelaksanaan = ref<Pelaksanaan[]>([]);
const allPengajuan   = ref<Pengajuan[]>([]);
const loading        = ref(false);
const error          = ref<string | null>(null);
const updatingId     = ref<string | null>(null);
const updateError    = ref('');

const showDetail      = ref(false);
const selected        = ref<Row | null>(null);
const detailLoading   = ref(false);
const pengajuanDetail = ref<PengajuanDetail | null>(null);

const editPembimbing    = ref(false);
const editPembimbingText = ref('');
const savingPembimbing  = ref(false);
const pembimbingError   = ref('');

// ── Laporan Magang ────────────────────────────────────────────
const laporanList    = ref<any[]>([]);
const laporanLoading = ref(false);
const latestLaporan  = computed(() => laporanList.value[0] ?? null);
const reviewCatatan  = ref('');
const reviewing      = ref(false);
const reviewAction   = ref<'disetujui'|'revisi'|''>('');
const reviewError    = ref('');

function fmtSize(bytes: number | null) {
  if (!bytes) return '';
  if (bytes < 1024)    return `${bytes} B`;
  if (bytes < 1048576) return `${(bytes / 1024).toFixed(1)} KB`;
  return `${(bytes / 1048576).toFixed(1)} MB`;
}

function formatStatusLaporan(s: string) {
  const m: Record<string, string> = {
    menunggu_review: 'Menunggu Review',
    disetujui:       'Disetujui',
    revisi:          'Perlu Revisi',
  };
  return m[s] ?? s;
}

async function fetchLaporan(pelaksanaanId: string) {
  laporanLoading.value = true;
  laporanList.value = [];
  reviewCatatan.value = '';
  reviewError.value = '';
  reviewAction.value = '';
  try {
    const r = await api.get(`/api/laporan/pelaksanaan/${pelaksanaanId}`);
    laporanList.value = Array.isArray(r.data?.data) ? r.data.data : [];
  } catch {
    laporanList.value = [];
  } finally {
    laporanLoading.value = false;
  }
}

async function submitReview(action: 'disetujui' | 'revisi') {
  if (!latestLaporan.value) return;
  if (action === 'revisi' && !reviewCatatan.value.trim()) {
    reviewError.value = 'Catatan wajib diisi saat meminta revisi.';
    return;
  }
  reviewError.value = '';
  reviewing.value = true;
  reviewAction.value = action;
  try {
    await api.patch(`/api/laporan/${latestLaporan.value.id}/review`, {
      status: action,
      catatan_hrd: reviewCatatan.value.trim(),
    });
    // Refresh laporan dan data pelaksanaan
    if (selected.value) {
      await fetchLaporan(selected.value.id);
      await fetchData();
      // Update status di selected
      const updated = allPelaksanaan.value.find(p => p.id === selected.value?.id);
      if (updated) {
        selected.value = {
          ...selected.value,
          status: updated.status,
        };
      }
    }
  } catch (e: any) {
    reviewError.value = e.response?.data?.message ?? 'Gagal mengirim review.';
  } finally {
    reviewing.value = false;
    reviewAction.value = '';
  }
}

const rows = computed<Row[]>(() => {
  const pMap = new Map(allPengajuan.value.map(p => [p.id, p]));
  return allPelaksanaan.value
    .filter(pl => pl.status !== "selesai")
    .map(pl => {
      const pj = pMap.get(pl.pengajuan_id);
      return {
        ...pl,
        nama_lengkap:    pj?.nama_lengkap    ?? "–",
        asal_institusi:  pj?.asal_institusi  ?? "–",
        jurusan:         pj?.jurusan         ?? "–",
        kategori_magang: pj?.kategori_magang ?? "–",
      };
    });
});

const pembimbingNama = computed<string>(() => selected.value?.pembimbing_nama ?? "");

async function fetchData() {
  loading.value = true; error.value = null;
  try {
    const [rPl, rPj] = await Promise.all([
      api.get("/api/pelaksanaan?page=1&limit=200"),
      api.get("/api/pengajuan?page=1&limit=200"),
    ]);
    allPelaksanaan.value = Array.isArray(rPl.data.data?.list) ? rPl.data.data.list : [];
    allPengajuan.value   = Array.isArray(rPj.data.data)       ? rPj.data.data      : [];
  } catch (e: any) {
    error.value = e.response?.data?.message ?? "Gagal memuat data.";
  } finally { loading.value = false; }
}

async function openDetail(row: Row) {
  selected.value = row;
  pengajuanDetail.value = null;
  showDetail.value = true;
  detailLoading.value = true;
  editPembimbing.value = false;
  editPembimbingText.value = '';
  pembimbingError.value = '';
  laporanList.value = [];
  try {
    const promises: Promise<any>[] = [api.get(`/api/pengajuan/${row.pengajuan_id}`)];
    if (['upload_laporan','penilaian','selesai'].includes(row.status)) {
      promises.push(fetchLaporan(row.id));
    }
    const [rPj] = await Promise.all(promises);
    pengajuanDetail.value = rPj.data?.data ?? null;
  } catch {
    // non-fatal
  } finally {
    detailLoading.value = false;
  }
}

function startEditPembimbing() {
  editPembimbingText.value = selected.value?.pembimbing_nama ?? '';
  pembimbingError.value = '';
  editPembimbing.value = true;
}

async function savePembimbing() {
  if (!selected.value) return;
  savingPembimbing.value = true;
  pembimbingError.value = '';
  try {
    await api.patch(`/api/pelaksanaan/${selected.value.id}/pembimbing`, {
      pembimbing: editPembimbingText.value.trim(),
    });
    const newNama = editPembimbingText.value.trim() || null;
    selected.value = { ...selected.value, pembimbing_nama: newNama };
    const idx = allPelaksanaan.value.findIndex(p => p.id === selected.value!.id);
    if (idx !== -1) allPelaksanaan.value[idx] = { ...allPelaksanaan.value[idx], pembimbing_nama: newNama };
    editPembimbing.value = false;
  } catch (e: any) {
    pembimbingError.value = e.response?.data?.message ?? 'Gagal menyimpan pembimbing';
  } finally {
    savingPembimbing.value = false;
  }
}

async function updateStatus(id: string, status: string) {
  updatingId.value = id;
  updateError.value = '';
  try {
    await api.patch(`/api/pelaksanaan/${id}/status`, { status });
    const idx = allPelaksanaan.value.findIndex(p => p.id === id);
    if (idx !== -1) allPelaksanaan.value[idx].status = status;
    if (selected.value?.id === id) selected.value = { ...selected.value, status };
  } catch (e: any) {
    updateError.value = e.response?.data?.message ?? 'Gagal memperbarui status';
    setTimeout(() => { updateError.value = ''; }, 4000);
  } finally {
    updatingId.value = null;
  }
}

function sisaHari(selesai: string): number {
  const s = new Date(selesai); s.setHours(0,0,0,0);
  const t = new Date(); t.setHours(0,0,0,0);
  return Math.ceil((s.getTime() - t.getTime()) / 86400000);
}
function formatStatus(s: string) {
  return ({ menunggu_mulai:"Belum Mulai", aktif:"Aktif", upload_laporan:"Upload Laporan", penilaian:"Penilaian", selesai:"Selesai" } as Record<string,string>)[s] ?? s;
}
function statusClass(s: string) {
  if (s === "aktif")          return "sp-badge sp-badge--green";
  if (s === "menunggu_mulai") return "sp-badge sp-badge--gray";
  if (s === "upload_laporan") return "sp-badge sp-badge--blue";
  if (s === "penilaian")      return "sp-badge sp-badge--orange";
  return "sp-badge sp-badge--gray";
}
function formatDate(d: string) {
  if (!d) return '–';
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
function formatJK(v: string) {
  return v === 'laki_laki' ? 'Laki-laki' : v === 'perempuan' ? 'Perempuan' : '–';
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
@keyframes spin { to { transform:rotate(360deg); } }

.sp-badge { display:inline-flex; align-items:center; font-size:11px; font-weight:700; padding:3px 9px; border-radius:100px; white-space:nowrap; }
.sp-badge--green  { background:#dcfce7; color:#16a34a; }
.sp-badge--gray   { background:#f3f4f6; color:#6b7280; }
.sp-badge--blue   { background:#dbeafe; color:#2563eb; }
.sp-badge--orange { background:#ffedd5; color:#ea580c; }

.sisa-hari       { font-size:12px; font-weight:700; color:#374151; }
.sisa-hari--warn { color:#ea580c; }
.sisa-hari--over { color:#dc2626; }

.aksi-cell { display:flex; align-items:center; gap:6px; flex-wrap:wrap; }
.btn-aksi { border:none; border-radius:7px; padding:5px 11px; font-size:11.5px; font-weight:700; font-family:inherit; cursor:pointer; white-space:nowrap; transition:opacity .15s; }
.btn-aksi:disabled { opacity:0.5; cursor:default; }
.btn-aksi--green  { background:#dcfce7; color:#15803d; }
.btn-aksi--green:hover:not(:disabled)  { background:#bbf7d0; }
.btn-aksi--orange { background:#ffedd5; color:#c2410c; }
.btn-aksi--orange:hover:not(:disabled) { background:#fed7aa; }
.btn-aksi--ghost  { background:#f3f4f6; color:#374151; }
.btn-aksi--ghost:hover { background:#e5e7eb; }

.update-error { margin:12px 16px; background:#fff1f2; border:1px solid #fecdd3; color:#be123c; font-size:12.5px; padding:8px 14px; border-radius:8px; }

/* ── Drawer ────────────────────────────────────────────────── */
.drawer-overlay { position:fixed; inset:0; background:rgba(0,0,0,0.45); backdrop-filter:blur(2px); z-index:200; display:flex; justify-content:flex-end; }
.drawer { width:min(560px,100vw); height:100vh; background:#fff; display:flex; flex-direction:column; box-shadow:-4px 0 24px rgba(0,0,0,0.12); overflow:hidden; }
.drawer-header { display:flex; align-items:flex-start; justify-content:space-between; padding:20px 24px 16px; border-bottom:1px solid #f0faf0; flex-shrink:0; }
.drawer-title { font-size:16px; font-weight:700; color:#111827; margin:0 0 2px; }
.drawer-sub { font-size:12.5px; color:#6b7280; margin:0; }
.drawer-close { background:#f3f4f6; border:none; border-radius:8px; width:32px; height:32px; display:flex; align-items:center; justify-content:center; cursor:pointer; color:#6b7280; flex-shrink:0; margin-top:2px; }
.drawer-close:hover { background:#e5e7eb; color:#111827; }
.drawer-body { flex:1; overflow-y:auto; padding:20px 24px; display:flex; flex-direction:column; gap:20px; }
.drawer-loading { align-items:center; justify-content:center; }
.drawer-footer { flex-shrink:0; padding:16px 24px; border-top:1px solid #f0faf0; background:#fff; display:flex; flex-direction:column; gap:10px; }

/* ── Detail sections (sama persis dgn HRDVerifikasi) ── */
.detail-status-banner { display:flex; align-items:center; gap:10px; padding:10px 14px; background:#f9fafb; border-radius:10px; }
.banner-sisa { font-size:12.5px; color:#6b7280; margin-left:auto; font-weight:500; }
.status-banner--aktif          { background:#f0fdf4; }
.status-banner--menunggu_mulai { background:#f9fafb; }
.status-banner--upload_laporan { background:#eff6ff; }
.status-banner--penilaian      { background:#fff7ed; }
.status-banner--selesai        { background:#f0fdf4; }

.detail-section { display:flex; flex-direction:column; gap:10px; }
.detail-section__title { font-size:11px; font-weight:700; color:#9ca3af; text-transform:uppercase; letter-spacing:0.08em; }
.detail-section__title-row { display:flex; align-items:center; gap:8px; }
.info-grid { display:grid; grid-template-columns:1fr 1fr; gap:10px; }
.info-item { display:flex; flex-direction:column; gap:2px; }
.info-item--full { grid-column:1/-1; }
.info-label { font-size:11px; color:#9ca3af; font-weight:500; }
.info-val { font-size:13px; color:#111827; font-weight:500; }

/* ── Pembimbing ── */
.pembimbing-box { display:flex; align-items:center; gap:12px; background:#f9fafb; border:1px solid #e9f5e9; border-radius:10px; padding:12px 16px; }
.pembimbing-avatar { width:36px; height:36px; border-radius:50%; background:linear-gradient(135deg,#48AF4A,#2d8f30); color:#fff; font-size:14px; font-weight:700; display:flex; align-items:center; justify-content:center; flex-shrink:0; }
.pembimbing-name { font-size:13px; font-weight:600; color:#111827; }
.pembimbing-role { font-size:11px; color:#9ca3af; margin-top:2px; }
.pembimbing-empty { display:flex; align-items:center; gap:8px; font-size:13px; color:#9ca3af; background:#f9fafb; border-radius:8px; padding:12px 14px; }
.edit-inline-btn { font-size:11px; font-weight:600; color:#16a34a; background:none; border:1px solid #bbf7d0; border-radius:6px; padding:2px 9px; cursor:pointer; font-family:inherit; transition:background .15s; margin-left:auto; }
.edit-inline-btn:hover { background:#f0fdf4; }
.pembimbing-edit { display:flex; flex-direction:column; gap:8px; }
.pembimbing-input { border:1.5px solid #d1d5db; border-radius:9px; padding:9px 12px; font-size:13px; font-family:inherit; color:#111827; outline:none; width:100%; box-sizing:border-box; transition:border-color .15s; }
.pembimbing-input:focus { border-color:#48AF4A; box-shadow:0 0 0 3px rgba(72,175,74,.12); }
.pembimbing-input:disabled { background:#f9fafb; color:#9ca3af; cursor:not-allowed; }
.pembimbing-err { font-size:12px; color:#dc2626; padding:2px 0; }
.pembimbing-edit-btns { display:flex; gap:8px; }
.btn-cancel-sm { flex:1; background:#f3f4f6; color:#374151; border:none; border-radius:8px; padding:8px 12px; font-size:12.5px; font-weight:600; font-family:inherit; cursor:pointer; }
.btn-cancel-sm:hover:not(:disabled) { background:#e5e7eb; }
.btn-cancel-sm:disabled { opacity:.5; cursor:not-allowed; }
.btn-save-sm { flex:1; background:#48AF4A; color:#fff; border:none; border-radius:8px; padding:8px 12px; font-size:12.5px; font-weight:600; font-family:inherit; cursor:pointer; }
.btn-save-sm:hover:not(:disabled) { background:#3d9e3f; }
.btn-save-sm:disabled { opacity:.5; cursor:not-allowed; }

/* ── Laporan Magang ── */
.laporan-section { background:#fafafa; border:1px solid #f1f5f9; border-radius:12px; padding:14px; }
.laporan-count { margin-left:auto; background:#e0f2fe; color:#0369a1; font-size:10.5px; font-weight:700; padding:1px 8px; border-radius:100px; }
.laporan-loading { display:flex; justify-content:center; padding:16px 0; }
.laporan-empty { display:flex; align-items:center; gap:10px; font-size:12.5px; color:#9ca3af; padding:10px 0; }
.laporan-list { display:flex; flex-direction:column; gap:6px; margin-bottom:12px; }
.laporan-item { display:flex; align-items:center; gap:10px; padding:10px 12px; border-radius:10px; background:#fff; border:1px solid #e5e7eb; }
.laporan-item--latest { border-color:#86efac; background:#f0fdf4; }
.laporan-item__icon { width:32px; height:32px; border-radius:8px; background:#f3f4f6; display:flex; align-items:center; justify-content:center; color:#6b7280; flex-shrink:0; }
.laporan-item--latest .laporan-item__icon { background:#dcfce7; color:#16a34a; }
.laporan-item__info { flex:1; min-width:0; }
.laporan-item__name { font-size:12px; font-weight:600; color:#111827; white-space:nowrap; overflow:hidden; text-overflow:ellipsis; }
.laporan-item__meta { display:flex; align-items:center; gap:6px; margin-top:2px; font-size:11px; color:#9ca3af; flex-wrap:wrap; }
.lap-versi { background:#eff6ff; color:#2563eb; font-size:10px; font-weight:700; padding:1px 6px; border-radius:100px; }
.laporan-item__right { display:flex; align-items:center; gap:6px; flex-shrink:0; }
.lap-status { font-size:10.5px; font-weight:600; padding:2px 8px; border-radius:100px; white-space:nowrap; }
.lap-status--menunggu_review { background:#fffbeb; color:#b45309; }
.lap-status--disetujui       { background:#f0fdf4; color:#16a34a; }
.lap-status--revisi          { background:#fef2f2; color:#dc2626; }
.btn-dl-lap { width:28px; height:28px; border-radius:7px; background:#f3f4f6; border:1px solid #e5e7eb; display:flex; align-items:center; justify-content:center; color:#6b7280; text-decoration:none; }
.btn-dl-lap:hover { background:#e5e7eb; color:#374151; }

.review-action-box { background:#fff; border:1.5px solid #fde68a; border-radius:10px; padding:12px 14px; display:flex; flex-direction:column; gap:8px; }
.review-action-box__label { font-size:11px; font-weight:700; color:#9ca3af; text-transform:uppercase; letter-spacing:0.06em; }
.review-catatan-input { border:1.5px solid #d1d5db; border-radius:8px; padding:8px 11px; font-size:12.5px; font-family:inherit; color:#111827; resize:vertical; outline:none; width:100%; box-sizing:border-box; }
.review-catatan-input:focus { border-color:#48AF4A; box-shadow:0 0 0 3px rgba(72,175,74,.1); }
.review-err { font-size:11.5px; color:#dc2626; }
.review-btns { display:flex; gap:8px; }
.btn-revisi { flex:1; display:flex; align-items:center; justify-content:center; gap:6px; border:none; border-radius:8px; padding:9px 12px; font-size:12.5px; font-weight:700; font-family:inherit; cursor:pointer; background:#fef2f2; color:#dc2626; }
.btn-revisi:hover:not(:disabled) { background:#fee2e2; }
.btn-revisi:disabled { opacity:.5; cursor:not-allowed; }
.btn-acc { flex:1; display:flex; align-items:center; justify-content:center; gap:6px; border:none; border-radius:8px; padding:9px 12px; font-size:12.5px; font-weight:700; font-family:inherit; cursor:pointer; background:#48AF4A; color:#fff; }
.btn-acc:hover:not(:disabled) { background:#3d9e3f; }
.btn-acc:disabled { opacity:.5; cursor:not-allowed; }

.review-done-box { display:flex; align-items:center; gap:8px; font-size:12.5px; font-weight:600; color:#16a34a; background:#f0fdf4; border:1px solid #bbf7d0; border-radius:8px; padding:9px 12px; }

.spinner--sm { width:14px; height:14px; border-width:2px; }
.spinner--white { border-color:rgba(255,255,255,.35); border-top-color:#fff; }

/* ── Nilai ── */
.nilai-display { display:flex; align-items:baseline; gap:6px; }
.nilai-score { font-size:36px; font-weight:800; color:#48AF4A; }
.nilai-label { font-size:14px; color:#9ca3af; }

/* ── Footer action buttons ── */
.btn-full { width:100%; border:none; border-radius:10px; padding:12px; font-size:13.5px; font-weight:700; font-family:inherit; cursor:pointer; display:flex; align-items:center; justify-content:center; gap:8px; }
.btn-full:disabled { opacity:0.5; cursor:default; }
.btn-full--green  { background:#48AF4A; color:#fff; }
.btn-full--green:hover:not(:disabled) { background:#3d9e3f; }
.btn-full--orange { background:#ea580c; color:#fff; }
.btn-full--orange:hover:not(:disabled) { background:#c2410c; }

@media (max-width:600px) { .info-grid { grid-template-columns:1fr; } .drawer { width:100vw; } }
</style>
