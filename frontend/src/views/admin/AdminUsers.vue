<template>
  <div>
    <div class="card">
      <div class="card-header">
        <div>
          <h3 class="card-title">Manajemen User</h3>
          <p class="card-sub">Kelola akun Admin, HRD, dan Peserta sistem e-Magang</p>
        </div>
        <div style="display:flex;gap:8px;align-items:center">
          <select v-model="filterRole" class="filter-select" @change="fetchUsers">
            <option value="">Semua Role</option>
            <option value="admin">Admin</option>
            <option value="hrd">HRD</option>
            <option value="peserta">Peserta</option>
          </select>
          <button class="btn-green-sm" @click="openModal(null)">+ Tambah User</button>
        </div>
      </div>

      <div v-if="loading" class="empty-state"><div class="spinner"></div></div>
      <div v-else-if="users.length === 0" class="empty-state">
        <div class="empty-state__icon">
          <svg width="36" height="36" viewBox="0 0 24 24" fill="none"><path d="M17 21v-2a4 4 0 00-4-4H5a4 4 0 00-4 4v2" stroke="#d1d5db" stroke-width="1.5"/><circle cx="9" cy="7" r="4" stroke="#d1d5db" stroke-width="1.5"/></svg>
        </div>
        <p>Belum ada user. Klik <strong>+ Tambah User</strong> untuk membuat akun baru.</p>
      </div>
      <div v-else class="table-wrap">
        <table class="data-table">
          <thead>
            <tr><th>Nama</th><th>Email</th><th>Role</th><th>Status Akun</th><th>Status Magang</th><th>Password</th><th>Dibuat</th><th>Aksi</th></tr>
          </thead>
          <tbody>
            <tr v-for="u in users" :key="u.id">
              <td class="td-nama">{{ u.nama_lengkap }}</td>
              <td class="td-email">{{ u.email }}</td>
              <td>
                <span :class="roleChipClass(u.role)">{{ roleLabel(u.role) }}</span>
              </td>
              <td>
                <span :class="u.is_active ? 'status-pill status-pill--green' : 'status-pill status-pill--gray'">
                  {{ u.is_active ? 'Aktif' : 'Nonaktif' }}
                </span>
              </td>
              <td>
                <span v-if="u.role === 'peserta'" :class="statusMagangClass(u.status_magang)">
                  {{ statusMagangLabel(u.status_magang) }}
                </span>
                <span v-else class="status-pill status-pill--gray" style="font-size:10px;">—</span>
              </td>
              <td>
                <span :class="u.password_changed ? 'pw-pill pw-pill--changed' : 'pw-pill pw-pill--default'">
                  {{ u.password_changed ? 'Sudah diganti' : 'Default' }}
                </span>
              </td>
              <td class="td-date">{{ formatDate(u.created_at) }}</td>
              <td>
                <div class="aksi-cell">
                  <button class="btn-aksi btn-aksi--ghost" @click="openModal(u)">Edit</button>
                  <button
                    :class="u.is_active ? 'btn-aksi btn-aksi--warn' : 'btn-aksi btn-aksi--green'"
                    @click="toggleUser(u)"
                    :disabled="toggling === u.id"
                  >{{ u.is_active ? 'Nonaktifkan' : 'Aktifkan' }}</button>
                  <button
                    v-if="u.role === 'peserta' && u.status_magang === 'selesai'"
                    class="btn-aksi btn-aksi--danger"
                    @click="openHapus(u)"
                    :disabled="hapusing === u.id"
                  >Hapus</button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div v-if="total > limit" class="pagination">
        <button class="page-btn" :disabled="page <= 1" @click="changePage(page - 1)">‹ Prev</button>
        <span class="page-info">{{ page }} / {{ Math.ceil(total / limit) }}</span>
        <button class="page-btn" :disabled="page >= Math.ceil(total / limit)" @click="changePage(page + 1)">Next ›</button>
      </div>
    </div>
  </div>

  <!-- Modal Tambah / Edit -->
  <Teleport to="body">
    <div v-if="showModal" class="modal-overlay" @click.self="closeModal">
      <div class="user-modal">
        <div class="modal-header">
          <h3 class="modal-title">{{ editTarget ? 'Edit User' : 'Tambah User Baru' }}</h3>
          <button class="modal-close" @click="closeModal">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none"><path d="M18 6L6 18M6 6l12 12" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"/></svg>
          </button>
        </div>
        <div class="modal-body">
          <div class="dform-group">
            <label class="dform-label">Nama Lengkap <span class="dform-req">*</span></label>
            <input v-model="form.nama_lengkap" type="text" class="dform-input" placeholder="Nama lengkap user" :disabled="saving" />
          </div>
          <div v-if="!editTarget" class="dform-group">
            <label class="dform-label">Email <span class="dform-req">*</span></label>
            <input v-model="form.email" type="email" class="dform-input" placeholder="email@telpp.co.id" :disabled="saving" />
          </div>
          <div v-if="editTarget" class="dform-group">
            <label class="dform-label">Email</label>
            <div class="dform-static">{{ editTarget.email }}</div>
          </div>
          <div class="dform-group">
            <label class="dform-label">Role <span v-if="!editTarget || editTarget.role !== 'peserta'" class="dform-req">*</span></label>
            <div v-if="editTarget && editTarget.role === 'peserta'" class="dform-static">
              <span class="role-chip role-chip--peserta">Peserta</span>
              <span class="role-note">Role peserta tidak dapat diubah dari sini</span>
            </div>
            <div v-else class="role-radio-group">
              <label class="role-radio" :class="{ 'role-radio--active': form.role === 'hrd' }">
                <input type="radio" v-model="form.role" value="hrd" :disabled="saving" />
                <div class="role-radio__icon" style="background:#f0fdf4;color:#3b82f6">
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none"><path d="M17 21v-2a4 4 0 00-4-4H5a4 4 0 00-4 4v2" stroke="currentColor" stroke-width="2"/><circle cx="9" cy="7" r="4" stroke="currentColor" stroke-width="2"/></svg>
                </div>
                <div>
                  <div class="role-radio__name">HRD</div>
                  <div class="role-radio__desc">Kelola pengajuan & peserta magang</div>
                </div>
              </label>
              <label class="role-radio" :class="{ 'role-radio--active': form.role === 'admin' }">
                <input type="radio" v-model="form.role" value="admin" :disabled="saving" />
                <div class="role-radio__icon" style="background:#f5f3ff;color:#7c3aed">
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>
                </div>
                <div>
                  <div class="role-radio__name">Admin</div>
                  <div class="role-radio__desc">Akses penuh ke seluruh sistem</div>
                </div>
              </label>
            </div>
          </div>
          <div v-if="!editTarget" class="dform-group">
            <label class="dform-label">Password <span class="dform-req">*</span></label>
            <div class="pw-input-wrap">
              <input
                v-model="form.password"
                :type="showPw ? 'text' : 'password'"
                class="dform-input dform-input--pw"
                placeholder="Min. 8 karakter (huruf + angka + spesial)"
                :disabled="saving"
              />
              <button class="pw-toggle" type="button" @click="showPw = !showPw" tabindex="-1">
                <svg v-if="showPw" width="14" height="14" viewBox="0 0 24 24" fill="none"><path d="M17.94 17.94A10.07 10.07 0 0112 20c-7 0-11-8-11-8a18.45 18.45 0 015.06-5.94M9.9 4.24A9.12 9.12 0 0112 4c7 0 11 8 11 8a18.5 18.5 0 01-2.16 3.19m-6.72-1.07a3 3 0 11-4.24-4.24" stroke="currentColor" stroke-width="2" stroke-linecap="round"/><line x1="1" y1="1" x2="23" y2="23" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>
                <svg v-else width="14" height="14" viewBox="0 0 24 24" fill="none"><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z" stroke="currentColor" stroke-width="2"/><circle cx="12" cy="12" r="3" stroke="currentColor" stroke-width="2"/></svg>
              </button>
            </div>
            <div class="dform-hint">Minimal 8 karakter, harus mengandung huruf, angka, dan karakter spesial</div>
          </div>
          <div v-if="formError" class="dform-error">{{ formError }}</div>
        </div>
        <div class="modal-footer">
          <button class="btn-cancel" @click="closeModal" :disabled="saving">Batal</button>
          <button class="btn-green" @click="saveUser" :disabled="saving">
            <div v-if="saving" class="spinner-sm"></div>
            {{ saving ? 'Menyimpan…' : (editTarget ? 'Simpan Perubahan' : 'Buat Akun') }}
          </button>
        </div>
      </div>
    </div>
  </Teleport>

  <!-- Modal Konfirmasi Hapus -->
  <Teleport to="body">
    <div v-if="showHapusModal" class="modal-overlay" @click.self="closeHapus">
      <div class="hapus-modal">
        <div class="hapus-modal__icon">
          <svg width="28" height="28" viewBox="0 0 24 24" fill="none">
            <path d="M3 6h18M8 6V4h8v2M19 6l-1 14H6L5 6" stroke="#dc2626" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
        </div>
        <h3 class="hapus-modal__title">Hapus Akun Peserta</h3>
        <p class="hapus-modal__desc">
          Anda akan menghapus akun <strong>{{ hapusTarget?.nama_lengkap }}</strong>.
        </p>
        <div class="hapus-modal__info">
          <p style="margin:0 0 6px;font-size:11px;font-weight:700;color:#9ca3af;text-transform:uppercase;letter-spacing:.06em;">Yang akan terjadi:</p>
          <ul style="margin:0;padding-left:16px;font-size:12.5px;color:#374151;line-height:1.9;">
            <li>Sistem mengirim <strong>email backup</strong> ke peserta (termasuk sertifikat PDF jika ada)</li>
            <li>Akun, data absensi, laporan, dan seluruh riwayat magang <strong>dihapus permanen</strong></li>
          </ul>
        </div>
        <p class="hapus-modal__warn">Tindakan ini tidak dapat dibatalkan.</p>
        <div class="hapus-modal__footer">
          <button class="btn-cancel" @click="closeHapus" :disabled="hapusing === hapusTarget?.id">Batal</button>
          <button class="btn-danger" @click="konfirmasiHapus" :disabled="hapusing === hapusTarget?.id">
            <div v-if="hapusing === hapusTarget?.id" class="spinner-sm"></div>
            {{ hapusing === hapusTarget?.id ? 'Menghapus…' : 'Ya, Hapus Akun' }}
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import api from "@/lib/api";

interface UserItem {
  id: string;
  nama_lengkap: string;
  email: string;
  role: "admin" | "hrd" | "peserta";
  is_active: boolean;
  password_changed: boolean;
  created_at: string;
  status_magang?: string | null;
}

const users   = ref<UserItem[]>([]);
const loading = ref(false);
const total   = ref(0);
const page    = ref(1);
const limit   = ref(20);
const filterRole = ref("");
const toggling = ref<string | null>(null);
const hapusing = ref<string | null>(null);

const showModal  = ref(false);
const editTarget = ref<UserItem | null>(null);
const saving     = ref(false);
const formError  = ref("");
const showPw     = ref(false);

const showHapusModal = ref(false);
const hapusTarget    = ref<UserItem | null>(null);

const form = ref({ nama_lengkap: "", email: "", role: "hrd" as "admin" | "hrd" | "peserta", password: "" });

async function fetchUsers() {
  loading.value = true;
  try {
    const params: Record<string, string | number> = { page: page.value, limit: limit.value };
    if (filterRole.value) params.role = filterRole.value;
    const r = await api.get("/api/admin/users", { params });
    users.value = r.data?.data ?? [];
    total.value = r.data?.total ?? 0;
  } catch { users.value = []; }
  finally { loading.value = false; }
}

function openModal(u: UserItem | null) {
  formError.value = "";
  showPw.value = false;
  editTarget.value = u;
  if (u) {
    form.value = { nama_lengkap: u.nama_lengkap, email: u.email, role: u.role, password: "" };
  } else {
    form.value = { nama_lengkap: "", email: "", role: "hrd", password: "" };
  }
  showModal.value = true;
}

function closeModal() {
  showModal.value = false;
  editTarget.value = null;
  formError.value = "";
}

async function saveUser() {
  formError.value = "";
  if (!form.value.nama_lengkap.trim()) { formError.value = "Nama lengkap wajib diisi"; return; }
  if (!editTarget.value) {
    if (!form.value.email.trim()) { formError.value = "Email wajib diisi"; return; }
    if (!form.value.password.trim()) { formError.value = "Password wajib diisi"; return; }
  }
  saving.value = true;
  try {
    if (editTarget.value) {
      await api.patch(`/api/admin/users/${editTarget.value.id}`, {
        nama_lengkap: form.value.nama_lengkap.trim(),
        role: form.value.role,
      });
    } else {
      await api.post("/api/admin/users", {
        nama_lengkap: form.value.nama_lengkap.trim(),
        email: form.value.email.trim(),
        password: form.value.password,
        role: form.value.role,
      });
    }
    closeModal();
    await fetchUsers();
  } catch (e: any) {
    formError.value = e.response?.data?.message ?? "Gagal menyimpan user";
  } finally { saving.value = false; }
}

async function toggleUser(u: UserItem) {
  if (toggling.value) return;
  toggling.value = u.id;
  try {
    await api.patch(`/api/admin/users/${u.id}/toggle`, { is_active: !u.is_active });
    await fetchUsers();
  } catch { }
  finally { toggling.value = null; }
}

function openHapus(u: UserItem) {
  hapusTarget.value = u;
  showHapusModal.value = true;
}

function closeHapus() {
  if (hapusing.value) return;
  showHapusModal.value = false;
  hapusTarget.value = null;
}

async function konfirmasiHapus() {
  if (!hapusTarget.value || hapusing.value) return;
  hapusing.value = hapusTarget.value.id;
  try {
    await api.delete(`/api/admin/users/${hapusTarget.value.id}`);
    showHapusModal.value = false;
    hapusTarget.value = null;
    await fetchUsers();
  } catch (e: any) {
    alert(e.response?.data?.message ?? "Gagal menghapus akun");
  } finally { hapusing.value = null; }
}

function changePage(p: number) {
  page.value = p;
  fetchUsers();
}

function roleLabel(role: string) {
  if (role === "admin") return "Admin";
  if (role === "hrd")   return "HRD";
  return "Peserta";
}

function roleChipClass(role: string) {
  if (role === "admin")   return "role-chip role-chip--admin";
  if (role === "hrd")     return "role-chip role-chip--hrd";
  return "role-chip role-chip--peserta";
}

function statusMagangLabel(s?: string | null): string {
  const map: Record<string, string> = {
    belum_daftar:          "Belum Daftar",
    diajukan:              "Diajukan",
    menunggu_verifikasi:   "Menunggu Verifikasi",
    diproses:              "Diproses",
    diterima:              "Diterima",
    ditolak:               "Ditolak",
    menunggu_mulai:        "Menunggu Mulai",
    aktif:                 "Magang Aktif",
    upload_laporan:        "Upload Laporan",
    penilaian:             "Penilaian",
    selesai:               "Selesai",
  };
  return s ? (map[s] ?? s) : "Belum Daftar";
}

function statusMagangClass(s?: string | null): string {
  if (!s || s === "belum_daftar") return "status-pill status-pill--gray";
  if (s === "selesai")            return "status-pill status-pill--teal";
  if (s === "aktif")              return "status-pill status-pill--green";
  if (s === "ditolak")            return "status-pill status-pill--red";
  if (s === "upload_laporan" || s === "penilaian") return "status-pill status-pill--orange";
  return "status-pill status-pill--blue";
}

function formatDate(s: string) {
  return new Date(s).toLocaleDateString("id-ID", { day: "2-digit", month: "short", year: "numeric" });
}

onMounted(fetchUsers);
</script>

<style scoped>
.card { background: #fff; border-radius: 14px; border: 1px solid #e9f5e9; box-shadow: 0 1px 3px rgba(13,40,24,0.05); overflow: hidden; }
.card-header { display: flex; align-items: center; justify-content: space-between; padding: 16px 20px; border-bottom: 1px solid #f0faf0; gap: 12px; }
.card-title  { font-size: 13.5px; font-weight: 700; color: #111827; margin: 0; }
.card-sub    { font-size: 11.5px; color: #9ca3af; margin: 2px 0 0; }

.filter-select {
  border: 1.5px solid #e5e7eb; border-radius: 8px; padding: 5px 10px;
  font-size: 12px; font-family: "Poppins", sans-serif; color: #374151;
  outline: none; cursor: pointer; background: #fff;
}
.filter-select:focus { border-color: #48AF4A; }

.btn-green-sm { background: #48AF4A; color: #fff; border: none; border-radius: 8px; padding: 6px 14px; font-size: 12px; font-weight: 600; font-family: "Poppins",sans-serif; cursor: pointer; white-space: nowrap; }

.empty-state { display: flex; flex-direction: column; align-items: center; padding: 44px 24px; gap: 12px; text-align: center; }
.empty-state__icon { width: 72px; height: 72px; background: #f9fafb; border-radius: 50%; display: flex; align-items: center; justify-content: center; }
.empty-state p { font-size: 13px; color: #9ca3af; line-height: 1.7; margin: 0; }

.table-wrap { overflow-x: auto; }
.data-table { width: 100%; border-collapse: collapse; font-size: 13px; }
.data-table th { padding: 11px 16px; text-align: left; font-size: 10.5px; font-weight: 600; color: #6b7280; background: #f9fafb; border-bottom: 1px solid #f1f5f9; text-transform: uppercase; letter-spacing: .04em; white-space: nowrap; }
.data-table td { padding: 13px 16px; border-bottom: 1px solid #f9fafb; color: #374151; vertical-align: middle; }
.data-table tr:last-child td { border-bottom: none; }
.td-nama  { font-weight: 600; color: #111827; }
.td-email { font-size: 12px; color: #6b7280; }
.td-date  { font-size: 12px; color: #9ca3af; white-space: nowrap; }

.role-chip          { padding: 3px 10px; border-radius: 100px; font-size: 11px; font-weight: 700; }
.role-chip--admin   { background: #f5f3ff; color: #7c3aed; }
.role-chip--hrd     { background: #eff6ff; color: #3b82f6; }
.role-chip--peserta { background: #f0fdf4; color: #16a34a; }

.status-pill         { display: inline-flex; align-items: center; font-size: 11px; font-weight: 700; padding: 3px 10px; border-radius: 100px; white-space: nowrap; }
.status-pill--green  { background: #dcfce7; color: #16a34a; }
.status-pill--teal   { background: #ccfbf1; color: #0f766e; }
.status-pill--gray   { background: #f3f4f6; color: #6b7280; }
.status-pill--blue   { background: #dbeafe; color: #1d4ed8; }
.status-pill--orange { background: #fef3c7; color: #b45309; }
.status-pill--red    { background: #fee2e2; color: #dc2626; }

.pw-pill           { font-size: 11px; font-weight: 600; padding: 3px 9px; border-radius: 100px; }
.pw-pill--changed  { background: #f0fdf4; color: #16a34a; }
.pw-pill--default  { background: #fef9c3; color: #92400e; }

.aksi-cell { display: flex; align-items: center; gap: 6px; flex-wrap: nowrap; }
.btn-aksi { border: none; border-radius: 7px; padding: 5px 11px; font-size: 11.5px; font-weight: 700; font-family: inherit; cursor: pointer; white-space: nowrap; transition: opacity .15s; }
.btn-aksi:disabled { opacity: .4; cursor: default; }
.btn-aksi--ghost  { background: #f3f4f6; color: #374151; }
.btn-aksi--ghost:hover { background: #e5e7eb; }
.btn-aksi--green  { background: #dcfce7; color: #16a34a; }
.btn-aksi--green:hover { background: #bbf7d0; }
.btn-aksi--warn   { background: #fef9c3; color: #92400e; }
.btn-aksi--warn:hover { background: #fde68a; }
.btn-aksi--danger { background: #fee2e2; color: #dc2626; }
.btn-aksi--danger:hover { background: #fecaca; }

.pagination { display: flex; align-items: center; justify-content: center; gap: 14px; padding: 14px; border-top: 1px solid #f0faf0; }
.page-btn { background: #f3f4f6; border: none; border-radius: 7px; padding: 6px 14px; font-size: 12.5px; font-family: inherit; cursor: pointer; color: #374151; }
.page-btn:hover:not(:disabled) { background: #e5e7eb; }
.page-btn:disabled { opacity: .4; cursor: default; }
.page-info { font-size: 12.5px; color: #6b7280; }

/* ── Modal Tambah/Edit ── */
.modal-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.5); backdrop-filter: blur(3px); z-index: 300; display: flex; align-items: center; justify-content: center; padding: 20px; }
.user-modal { background: #fff; border-radius: 18px; width: min(480px, 100%); box-shadow: 0 24px 80px rgba(0,0,0,0.22); display: flex; flex-direction: column; overflow: hidden; }
.modal-header { display: flex; align-items: center; justify-content: space-between; padding: 20px 24px 14px; border-bottom: 1px solid #f0faf0; }
.modal-title  { font-size: 15px; font-weight: 700; color: #111827; margin: 0; }
.modal-close  { background: #f3f4f6; border: none; border-radius: 8px; width: 30px; height: 30px; display: flex; align-items: center; justify-content: center; cursor: pointer; color: #6b7280; }
.modal-close:hover { background: #e5e7eb; }
.modal-body   { padding: 20px 24px; display: flex; flex-direction: column; gap: 16px; }
.modal-footer { display: flex; gap: 10px; padding: 14px 24px; border-top: 1px solid #f0faf0; }

.dform-group  { display: flex; flex-direction: column; gap: 5px; }
.dform-label  { font-size: 12px; font-weight: 600; color: #374151; }
.dform-req    { color: #dc2626; }
.dform-hint   { font-size: 11px; color: #9ca3af; }
.dform-error  { font-size: 12.5px; color: #dc2626; background: #fef2f2; border: 1px solid #fecaca; border-radius: 8px; padding: 9px 13px; }
.dform-input  { border: 1.5px solid #e5e7eb; border-radius: 9px; padding: 9px 12px; font-size: 13px; font-family: inherit; color: #111827; outline: none; transition: border-color .15s; }
.dform-input:focus { border-color: #48AF4A; box-shadow: 0 0 0 3px rgba(72,175,74,.12); }
.dform-input:disabled { background: #f9fafb; color: #9ca3af; cursor: not-allowed; }
.dform-static { font-size: 13px; color: #6b7280; padding: 9px 12px; background: #f9fafb; border-radius: 9px; border: 1.5px solid #e5e7eb; display: flex; align-items: center; gap: 10px; }
.role-note { font-size: 11px; color: #9ca3af; }

.dform-input--pw { padding-right: 40px; width: 100%; box-sizing: border-box; }
.pw-input-wrap { position: relative; }
.pw-toggle {
  position: absolute; right: 10px; top: 50%; transform: translateY(-50%);
  background: none; border: none; cursor: pointer; color: #9ca3af; padding: 4px;
  display: flex; align-items: center;
}
.pw-toggle:hover { color: #374151; }

.role-radio-group { display: flex; flex-direction: column; gap: 8px; }
.role-radio {
  display: flex; align-items: center; gap: 12px;
  border: 1.5px solid #e5e7eb; border-radius: 10px; padding: 12px 14px;
  cursor: pointer; transition: border-color .15s, background .15s;
}
.role-radio input[type=radio] { display: none; }
.role-radio--active { border-color: #48AF4A; background: #f0fdf4; }
.role-radio:hover:not(.role-radio--active) { border-color: #a7d9a8; background: #fafffe; }
.role-radio__icon { width: 32px; height: 32px; border-radius: 8px; display: flex; align-items: center; justify-content: center; flex-shrink: 0; }
.role-radio__name { font-size: 13px; font-weight: 700; color: #111827; }
.role-radio__desc { font-size: 11px; color: #9ca3af; margin-top: 1px; }

.btn-cancel { flex: 1; background: #f3f4f6; color: #374151; border: none; border-radius: 9px; padding: 10px 16px; font-size: 13px; font-weight: 600; font-family: inherit; cursor: pointer; }
.btn-cancel:hover:not(:disabled) { background: #e5e7eb; }
.btn-cancel:disabled { opacity: .5; cursor: default; }
.btn-green { flex: 1; background: #48AF4A; color: #fff; border: none; border-radius: 9px; padding: 10px 22px; font-size: 13px; font-weight: 600; font-family: "Poppins",sans-serif; cursor: pointer; display: flex; align-items: center; justify-content: center; gap: 7px; }
.btn-green:hover:not(:disabled) { background: #3d9e3f; }
.btn-green:disabled { opacity: .5; cursor: not-allowed; }

/* ── Modal Hapus ── */
.hapus-modal {
  background: #fff; border-radius: 18px; width: min(440px, 100%);
  box-shadow: 0 24px 80px rgba(0,0,0,0.22);
  padding: 32px 28px 24px; display: flex; flex-direction: column; align-items: center; gap: 12px; text-align: center;
}
.hapus-modal__icon {
  width: 60px; height: 60px; background: #fee2e2; border-radius: 50%;
  display: flex; align-items: center; justify-content: center; margin-bottom: 4px;
}
.hapus-modal__title { font-size: 16px; font-weight: 700; color: #111827; margin: 0; }
.hapus-modal__desc  { font-size: 13.5px; color: #374151; margin: 0; line-height: 1.7; }
.hapus-modal__info  {
  background: #f9fafb; border: 1px solid #e5e7eb; border-radius: 10px;
  padding: 14px 16px; text-align: left; width: 100%; box-sizing: border-box;
}
.hapus-modal__warn  { font-size: 12px; color: #dc2626; font-weight: 600; margin: 0; }
.hapus-modal__footer { display: flex; gap: 10px; width: 100%; margin-top: 4px; }

.btn-danger { flex: 1; background: #dc2626; color: #fff; border: none; border-radius: 9px; padding: 10px 22px; font-size: 13px; font-weight: 600; font-family: "Poppins",sans-serif; cursor: pointer; display: flex; align-items: center; justify-content: center; gap: 7px; }
.btn-danger:hover:not(:disabled) { background: #b91c1c; }
.btn-danger:disabled { opacity: .5; cursor: not-allowed; }

.spinner { width: 36px; height: 36px; border: 3px solid #e5e7eb; border-top-color: #48AF4A; border-radius: 50%; animation: spin .8s linear infinite; }
.spinner-sm { width: 14px; height: 14px; border: 2px solid rgba(255,255,255,0.35); border-top-color: #fff; border-radius: 50%; animation: spin .7s linear infinite; display: inline-block; }
@keyframes spin { to { transform: rotate(360deg); } }
</style>
