<template>
  <div class="aj-wrap">
    <div class="page-card">
    <!-- HEADER -->
    <div class="aj-header">
      <div>
        <h3 class="aj-title">Kelola Periode Magang</h3>
        <p class="aj-sub">Tambah, edit, aktifkan, atau hapus periode pendaftaran magang</p>
      </div>
      <div class="aj-header-actions">
        <a href="/#jadwal" target="_blank" class="btn-preview">Lihat Landing ↗</a>
        <button class="btn-add" @click="openModal(null)">+ Tambah Periode</button>
      </div>
    </div>

    <!-- PERIODE LIST -->
    <div class="aj-section">
      <div v-if="loading" class="aj-loading"><div class="spinner"></div></div>
      <div v-else-if="list.length === 0" class="aj-empty">
        <svg width="36" height="36" viewBox="0 0 24 24" fill="none"><rect x="3" y="4" width="18" height="18" rx="2" stroke="#d1d5db" stroke-width="1.5"/><line x1="3" y1="10" x2="21" y2="10" stroke="#d1d5db" stroke-width="1.5"/></svg>
        <p>Belum ada periode. Klik <strong>+ Tambah Periode</strong> untuk mulai.</p>
      </div>
      <div v-else class="aj-table-wrap">
        <table class="aj-table">
          <thead>
            <tr>
              <th>Nama</th>
              <th>Tanggal Buka</th>
              <th>Tanggal Tutup</th>
              <th>Kuota</th>
              <th>Status</th>
              <th>Aksi</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="p in list" :key="p.id">
              <td class="aj-td-name">{{ p.nama }}</td>
              <td>{{ formatDate(p.tanggal_buka) }}</td>
              <td>{{ formatDate(p.tanggal_tutup) }}</td>
              <td>{{ p.kuota ?? '∞' }}</td>
              <td>
                <span class="status-pill" :class="periodeStatusClass(p)">
                  {{ periodeStatusLabel(p) }}
                </span>
              </td>
              <td>
                <div class="aj-td-actions">
                  <button class="btn-sm btn-sm--edit" @click="openModal(p)">Edit</button>
                  <button
                    class="btn-sm"
                    :class="p.is_active ? 'btn-sm--orange' : 'btn-sm--green'"
                    @click="toggleAktif(p)"
                  >{{ p.is_active ? 'Nonaktifkan' : 'Aktifkan' }}</button>
                  <button class="btn-sm btn-sm--red" @click="confirmDelete(p)">Hapus</button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- ALERTS -->
    <div v-if="error" class="aj-alert aj-alert--error">{{ error }}</div>

    <!-- MODAL TAMBAH/EDIT -->
    <Teleport to="body">
      <div v-if="showModal" class="modal-backdrop" @click.self="closeModal">
        <div class="modal-box">
          <div class="modal-header">
            <h4>{{ editTarget ? 'Edit Periode' : 'Tambah Periode Baru' }}</h4>
            <button class="modal-close" @click="closeModal">×</button>
          </div>
          <div class="modal-body">
            <div class="mf-group">
              <label>Nama Periode</label>
              <input v-model="modalForm.nama" type="text" placeholder="contoh: Periode I 2025" class="mf-input" />
            </div>
            <div class="mf-row">
              <div class="mf-group">
                <label>Tanggal Buka</label>
                <input v-model="modalForm.tanggal_buka" type="date" class="mf-input" />
              </div>
              <div class="mf-group">
                <label>Tanggal Tutup</label>
                <input v-model="modalForm.tanggal_tutup" type="date" class="mf-input" />
              </div>
            </div>
            <div class="mf-group">
              <label>Kuota <span class="mf-hint">(kosongkan = tidak terbatas)</span></label>
              <input v-model="modalForm.kuota" type="number" min="1" placeholder="mis. 50" class="mf-input mf-input--half" />
            </div>
            <div v-if="modalError" class="aj-alert aj-alert--error">{{ modalError }}</div>
          </div>
          <div class="modal-footer">
            <button class="btn-cancel" @click="closeModal">Batal</button>
            <button class="btn-save" @click="savePeriode" :disabled="saving">
              <span v-if="saving" class="spinner-sm"></span>
              {{ saving ? 'Menyimpan…' : 'Simpan' }}
            </button>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- MODAL KONFIRMASI HAPUS -->
    <Teleport to="body">
      <div v-if="showDeleteModal" class="modal-backdrop" @click.self="showDeleteModal = false">
        <div class="modal-box modal-box--sm">
          <div class="modal-header">
            <h4>Hapus Periode</h4>
            <button class="modal-close" @click="showDeleteModal = false">×</button>
          </div>
          <div class="modal-body">
            <p class="delete-confirm-text">
              Yakin ingin menghapus periode <strong>{{ deleteTarget?.nama }}</strong>?
              <br/>
              <span class="delete-warn">Tindakan ini tidak bisa dibatalkan.</span>
            </p>
          </div>
          <div class="modal-footer">
            <button class="btn-cancel" @click="showDeleteModal = false">Batal</button>
            <button class="btn-delete" @click="doDelete" :disabled="deleting">
              <span v-if="deleting" class="spinner-sm"></span>
              {{ deleting ? 'Menghapus…' : 'Ya, Hapus' }}
            </button>
          </div>
        </div>
      </div>
    </Teleport>

    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import api from '@/lib/api';

interface Periode {
  id: string;
  nama: string;
  tanggal_buka: string;
  tanggal_tutup: string;
  kuota: number | null;
  is_active: boolean;
  created_at: string;
}

const list        = ref<Periode[]>([]);
const loading     = ref(true);
const error       = ref('');
const saving      = ref(false);
const deleting    = ref(false);
const showModal   = ref(false);
const showDeleteModal = ref(false);
const editTarget  = ref<Periode | null>(null);
const deleteTarget = ref<Periode | null>(null);
const modalError  = ref('');

const modalForm = ref({ nama: '', tanggal_buka: '', tanggal_tutup: '', kuota: '' });

function formatDate(d: string) {
  return new Date(d).toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' });
}

function periodeStatusLabel(p: Periode) {
  const now  = new Date();
  const buka = new Date(p.tanggal_buka);
  const tup  = new Date(p.tanggal_tutup);
  if (!p.is_active) return 'Nonaktif';
  if (now < buka)   return 'Akan Datang';
  if (now > tup)    return 'Selesai';
  return 'Aktif';
}

function periodeStatusClass(p: Periode) {
  const lbl = periodeStatusLabel(p);
  return {
    'pill--green':  lbl === 'Aktif',
    'pill--yellow': lbl === 'Akan Datang',
    'pill--gray':   lbl === 'Selesai' || lbl === 'Nonaktif',
  };
}

async function fetchAll() {
  loading.value = true;
  error.value = '';
  try {
    const res = await api.get('/api/admin/periode');
    list.value = Array.isArray(res.data?.data) ? res.data.data : [];
  } catch { error.value = 'Gagal memuat data periode.'; }
  finally { loading.value = false; }
}

async function toggleAktif(p: Periode) {
  error.value = '';
  try {
    await api.patch(`/api/admin/periode/${p.id}/aktif`, { is_active: !p.is_active });
    await fetchAll();
  } catch { error.value = 'Gagal mengubah status periode.'; }
}

function openModal(p: Periode | null) {
  editTarget.value = p;
  modalError.value = '';
  if (p) {
    modalForm.value = {
      nama: p.nama,
      tanggal_buka:  p.tanggal_buka.split('T')[0],
      tanggal_tutup: p.tanggal_tutup.split('T')[0],
      kuota: p.kuota !== null ? String(p.kuota) : '',
    };
  } else {
    modalForm.value = { nama: '', tanggal_buka: '', tanggal_tutup: '', kuota: '' };
  }
  showModal.value = true;
}

function closeModal() { showModal.value = false; }

async function savePeriode() {
  if (!modalForm.value.nama || !modalForm.value.tanggal_buka || !modalForm.value.tanggal_tutup) {
    modalError.value = 'Nama, tanggal buka, dan tanggal tutup wajib diisi.';
    return;
  }
  saving.value = true;
  modalError.value = '';
  try {
    const payload = {
      nama:          modalForm.value.nama,
      tanggal_buka:  modalForm.value.tanggal_buka,
      tanggal_tutup: modalForm.value.tanggal_tutup,
      kuota:         modalForm.value.kuota ? parseInt(modalForm.value.kuota) : null,
    };
    if (editTarget.value) {
      await api.put(`/api/admin/periode/${editTarget.value.id}`, payload);
    } else {
      await api.post('/api/admin/periode', payload);
    }
    await fetchAll();
    closeModal();
  } catch (e: any) {
    modalError.value = e.response?.data?.message ?? 'Gagal menyimpan periode.';
  } finally { saving.value = false; }
}

function confirmDelete(p: Periode) {
  deleteTarget.value = p;
  showDeleteModal.value = true;
}

async function doDelete() {
  if (!deleteTarget.value) return;
  deleting.value = true;
  error.value = '';
  try {
    await api.delete(`/api/admin/periode/${deleteTarget.value.id}`);
    showDeleteModal.value = false;
    deleteTarget.value = null;
    await fetchAll();
  } catch (e: any) {
    error.value = e.response?.data?.message ?? 'Gagal menghapus periode.';
    showDeleteModal.value = false;
  } finally { deleting.value = false; }
}

onMounted(fetchAll);
</script>

<style scoped>
.aj-wrap  { display: flex; flex-direction: column; gap: 20px; }
.page-card { background:#fff; border:1px solid #e9f5e9; border-radius:14px; padding:20px; box-shadow:0 1px 3px rgba(13,40,24,.04); display:flex; flex-direction:column; gap:20px; }

.aj-header {
  display: flex; justify-content: space-between; align-items: flex-start;
  gap: 16px; flex-wrap: wrap;
}
.aj-title { font-size: 15px; font-weight: 700; color: #111827; margin: 0 0 4px; }
.aj-sub   { font-size: 12.5px; color: #6b7280; margin: 0; }
.aj-header-actions { display: flex; gap: 8px; align-items: center; }

.btn-preview {
  background: none; border: 1.5px solid #d1d5db; border-radius: 8px;
  padding: 7px 14px; font-size: 12px; font-weight: 600; color: #374151;
  cursor: pointer; text-decoration: none; font-family: "Poppins", sans-serif;
  transition: border-color .15s;
}
.btn-preview:hover { border-color: #48AF4A; color: #48AF4A; }
.btn-add {
  background: #48AF4A; color: #fff; border: none; border-radius: 8px;
  padding: 8px 18px; font-size: 13px; font-weight: 600; cursor: pointer;
  font-family: "Poppins", sans-serif; transition: opacity .15s;
}
.btn-add:hover { opacity: .88; }

/* ── Section ── */
.aj-section {
  background: #fff; border-radius: 12px; border: 1px solid #e9f5e9;
  padding: 4px 0; display: flex; flex-direction: column;
  box-shadow: 0 1px 3px rgba(13,40,24,.04);
}

/* ── Table ── */
.aj-table-wrap { overflow-x: auto; }
.aj-table {
  width: 100%; border-collapse: collapse; font-size: 13px;
}
.aj-table th {
  padding: 10px 14px; text-align: left; font-size: 11px; font-weight: 700;
  color: #6b7280; border-bottom: 1.5px solid #f1f5f9; background: #fafafa;
  text-transform: uppercase; letter-spacing: .5px;
}
.aj-table td { padding: 12px 14px; border-bottom: 1px solid #f1f5f9; color: #374151; }
.aj-table tbody tr:last-child td { border-bottom: none; }
.aj-table tbody tr:hover { background: #fafff9; }
.aj-td-name   { font-weight: 600; color: #111827; }
.aj-td-actions { display: flex; gap: 6px; align-items: center; }

.status-pill {
  display: inline-flex; align-items: center;
  padding: 3px 10px; border-radius: 100px; font-size: 11px; font-weight: 700;
}
.pill--green  { background: #dcfce7; color: #16a34a; }
.pill--yellow { background: #fef3c7; color: #92400e; }
.pill--gray   { background: #f3f4f6; color: #6b7280; }

.btn-sm {
  padding: 5px 11px; border-radius: 6px; font-size: 11.5px; font-weight: 600;
  cursor: pointer; border: none; font-family: "Poppins", sans-serif; transition: opacity .15s;
  white-space: nowrap;
}
.btn-sm:hover      { opacity: .82; }
.btn-sm--edit      { background: #f1f5f9; color: #374151; }
.btn-sm--green     { background: #dcfce7; color: #16a34a; }
.btn-sm--orange    { background: #fef3c7; color: #92400e; }
.btn-sm--red       { background: #fee2e2; color: #dc2626; }

/* ── Alerts ── */
.aj-alert { padding: 10px 14px; border-radius: 8px; font-size: 12.5px; }
.aj-alert--error { background: #fef2f2; color: #dc2626; border: 1px solid #fecaca; }

.aj-loading { display: flex; justify-content: center; padding: 40px; }
.aj-empty   { display: flex; flex-direction: column; align-items: center; gap: 10px; padding: 40px 0; color: #9ca3af; font-size: 13px; text-align: center; }

/* ── Modal ── */
.modal-backdrop {
  position: fixed; inset: 0; background: rgba(0,0,0,.45);
  z-index: 9999; display: flex; align-items: center; justify-content: center;
}
.modal-box {
  background: #fff; border-radius: 16px; width: 480px; max-width: 95vw;
  box-shadow: 0 20px 60px rgba(0,0,0,.2); overflow: hidden;
}
.modal-box--sm { width: 400px; }
.modal-header {
  display: flex; justify-content: space-between; align-items: center;
  padding: 18px 22px; border-bottom: 1px solid #f1f5f9;
}
.modal-header h4 { font-size: 14px; font-weight: 700; margin: 0; color: #111827; }
.modal-close { background: none; border: none; font-size: 20px; color: #9ca3af; cursor: pointer; line-height: 1; padding: 0 4px; }
.modal-body   { padding: 20px 22px; display: flex; flex-direction: column; gap: 14px; }
.modal-footer { padding: 16px 22px; border-top: 1px solid #f1f5f9; display: flex; justify-content: flex-end; gap: 10px; }

.mf-group { display: flex; flex-direction: column; gap: 5px; }
.mf-row   { display: grid; grid-template-columns: 1fr 1fr; gap: 12px; }
.mf-group label { font-size: 12px; font-weight: 600; color: #374151; }
.mf-hint  { font-size: 10.5px; color: #9ca3af; font-weight: 400; margin-left: 4px; }
.mf-input {
  border: 1.5px solid #e5e7eb; border-radius: 8px; padding: 9px 12px;
  font-size: 13px; font-family: "Poppins", sans-serif; color: #111827;
  outline: none; transition: border-color .15s; width: 100%; box-sizing: border-box;
}
.mf-input:focus  { border-color: #48AF4A; }
.mf-input--half  { max-width: 160px; }

.btn-cancel {
  background: #f1f5f9; color: #374151; border: none; border-radius: 8px;
  padding: 8px 20px; font-size: 13px; font-weight: 600; cursor: pointer;
  font-family: "Poppins", sans-serif;
}
.btn-save {
  background: #48AF4A; color: #fff; border: none; border-radius: 8px;
  padding: 8px 20px; font-size: 13px; font-weight: 600; cursor: pointer;
  font-family: "Poppins", sans-serif; display: inline-flex; align-items: center; gap: 7px;
  transition: opacity .15s;
}
.btn-save:disabled { opacity: .65; cursor: not-allowed; }
.btn-delete {
  background: #dc2626; color: #fff; border: none; border-radius: 8px;
  padding: 8px 20px; font-size: 13px; font-weight: 600; cursor: pointer;
  font-family: "Poppins", sans-serif; display: inline-flex; align-items: center; gap: 7px;
  transition: opacity .15s;
}
.btn-delete:hover    { background: #b91c1c; }
.btn-delete:disabled { opacity: .65; cursor: not-allowed; }

.delete-confirm-text { font-size: 13.5px; color: #374151; line-height: 1.7; margin: 0; }
.delete-warn { font-size: 12px; color: #dc2626; font-weight: 600; }

.spinner    { width: 28px; height: 28px; border: 3px solid #e5e7eb; border-top-color: #48AF4A; border-radius: 50%; animation: spin .7s linear infinite; }
.spinner-sm { width: 14px; height: 14px; border: 2px solid rgba(255,255,255,.4); border-top-color: #fff; border-radius: 50%; animation: spin .7s linear infinite; display: inline-block; }
@keyframes spin { to { transform: rotate(360deg); } }

@media (max-width: 600px) {
  .mf-row { grid-template-columns: 1fr; }
  .aj-td-actions { flex-wrap: wrap; }
}
</style>
