<template>
  <DashboardLayout :nav-groups="navGroups" role-name="Super Admin" default-tab="beranda" ref="layout">
    <template #default>

      <!-- ── BERANDA ── -->
      <template v-if="activeTab === 'beranda'">
        <div class="welcome-banner">
          <div>
            <div class="welcome-banner__greeting">Selamat datang, {{ firstName }}! 🛡️</div>
            <div class="welcome-banner__sub">Kelola seluruh sistem e-Magang PT TELPP dari sini.</div>
          </div>
          <div class="welcome-banner__badge">ADMIN</div>
        </div>

        <div class="stats-row">
          <div class="stat-card">
            <div class="stat-card__icon" style="background:#eff6ff;color:#3b82f6">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none"><path d="M17 21v-2a4 4 0 00-4-4H5a4 4 0 00-4 4v2" stroke="currentColor" stroke-width="2"/><circle cx="9" cy="7" r="4" stroke="currentColor" stroke-width="2"/></svg>
            </div>
            <div><div class="stat-card__label">Total User</div><div class="stat-card__value">2</div><div class="stat-card__sub">Admin & HRD</div></div>
          </div>
          <div class="stat-card">
            <div class="stat-card__icon" style="background:#f0fdf4;color:#16a34a">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none"><rect x="3" y="4" width="18" height="18" rx="2" ry="2" stroke="currentColor" stroke-width="2"/><line x1="3" y1="10" x2="21" y2="10" stroke="currentColor" stroke-width="2"/></svg>
            </div>
            <div><div class="stat-card__label">Periode Aktif</div><div class="stat-card__value">0</div></div>
          </div>
          <div class="stat-card">
            <div class="stat-card__icon" style="background:#fefce8;color:#ca8a04">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none"><circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2"/><line x1="2" y1="12" x2="22" y2="12" stroke="currentColor" stroke-width="2"/><path d="M12 2a15.3 15.3 0 014 10 15.3 15.3 0 01-4 10 15.3 15.3 0 01-4-10 15.3 15.3 0 014-10z" stroke="currentColor" stroke-width="2"/></svg>
            </div>
            <div><div class="stat-card__label">Landing Content</div><div class="stat-card__value">6</div><div class="stat-card__sub">item tersimpan</div></div>
          </div>
          <div class="stat-card">
            <div class="stat-card__icon" style="background:#fff1f2;color:#e11d48">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none"><circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2"/><line x1="12" y1="8" x2="12" y2="12" stroke="currentColor" stroke-width="2" stroke-linecap="round"/><line x1="12" y1="16" x2="12.01" y2="16" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>
            </div>
            <div><div class="stat-card__label">FAQ Aktif</div><div class="stat-card__value">6</div></div>
          </div>
        </div>

        <div class="quick-grid">
          <button class="quick-card" @click="setTab('users')">
            <div class="quick-card__icon" style="background:#eff6ff;color:#3b82f6">
              <svg width="22" height="22" viewBox="0 0 24 24" fill="none"><path d="M17 21v-2a4 4 0 00-4-4H5a4 4 0 00-4 4v2" stroke="currentColor" stroke-width="2"/><circle cx="9" cy="7" r="4" stroke="currentColor" stroke-width="2"/></svg>
            </div>
            <div class="quick-card__label">Manajemen User</div>
          </button>
          <button class="quick-card" @click="setTab('periode')">
            <div class="quick-card__icon" style="background:#f0fdf4;color:#16a34a">
              <svg width="22" height="22" viewBox="0 0 24 24" fill="none"><rect x="3" y="4" width="18" height="18" rx="2" ry="2" stroke="currentColor" stroke-width="2"/><line x1="3" y1="10" x2="21" y2="10" stroke="currentColor" stroke-width="2"/></svg>
            </div>
            <div class="quick-card__label">Kelola Periode</div>
          </button>
          <button class="quick-card" @click="setTab('landing')">
            <div class="quick-card__icon" style="background:#fefce8;color:#ca8a04">
              <svg width="22" height="22" viewBox="0 0 24 24" fill="none"><circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2"/><line x1="2" y1="12" x2="22" y2="12" stroke="currentColor" stroke-width="2"/></svg>
            </div>
            <div class="quick-card__label">Landing Page</div>
          </button>
          <button class="quick-card" @click="setTab('statistik')">
            <div class="quick-card__icon" style="background:#fdf4ff;color:#9333ea">
              <svg width="22" height="22" viewBox="0 0 24 24" fill="none"><line x1="18" y1="20" x2="18" y2="10" stroke="currentColor" stroke-width="2" stroke-linecap="round"/><line x1="12" y1="20" x2="12" y2="4" stroke="currentColor" stroke-width="2" stroke-linecap="round"/><line x1="6" y1="20" x2="6" y2="14" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>
            </div>
            <div class="quick-card__label">Statistik</div>
          </button>
        </div>
      </template>

      <!-- ── MANAJEMEN USER ── -->
      <template v-else-if="activeTab === 'users'">
        <div class="card">
          <div class="card-header">
            <h3 class="card-title">Manajemen User</h3>
            <button class="btn-green-sm">+ Tambah User</button>
          </div>
          <div class="table-wrap">
            <table class="data-table">
              <thead><tr><th>Nama</th><th>Email</th><th>Role</th><th>Status</th><th>Aksi</th></tr></thead>
              <tbody>
                <tr>
                  <td>Super Admin</td><td>admin@telpp.co.id</td>
                  <td><span class="role-chip role-chip--admin">Admin</span></td>
                  <td><span class="status-chip status-chip--active">Aktif</span></td>
                  <td><button class="action-btn">Edit</button></td>
                </tr>
                <tr>
                  <td>Staff HRD</td><td>hrd@telpp.co.id</td>
                  <td><span class="role-chip role-chip--hrd">HRD</span></td>
                  <td><span class="status-chip status-chip--active">Aktif</span></td>
                  <td><button class="action-btn">Edit</button></td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </template>

      <!-- ── PERIODE ── -->
      <template v-else-if="activeTab === 'periode'">
        <div class="card">
          <div class="card-header">
            <h3 class="card-title">Periode Magang</h3>
            <button class="btn-green-sm">+ Tambah Periode</button>
          </div>
          <div class="empty-state">
            <div class="empty-state__icon"><svg width="36" height="36" viewBox="0 0 24 24" fill="none"><rect x="3" y="4" width="18" height="18" rx="2" ry="2" stroke="#d1d5db" stroke-width="1.5"/><line x1="3" y1="10" x2="21" y2="10" stroke="#d1d5db" stroke-width="1.5"/></svg></div>
            <p>Belum ada periode magang.<br/>Tambahkan periode baru untuk membuka pendaftaran.</p>
            <button class="btn-green">+ Buat Periode Baru</button>
          </div>
        </div>
      </template>

      <!-- ── LANDING ── -->
      <template v-else-if="activeTab === 'landing'">
        <div class="card">
          <div class="card-header">
            <h3 class="card-title">Landing Page Settings</h3>
            <button class="btn-green-sm" @click="goToLandingSettings">Buka Editor →</button>
          </div>
          <div class="empty-state">
            <div class="empty-state__icon"><svg width="36" height="36" viewBox="0 0 24 24" fill="none"><circle cx="12" cy="12" r="10" stroke="#d1d5db" stroke-width="1.5"/><line x1="2" y1="12" x2="22" y2="12" stroke="#d1d5db" stroke-width="1.5"/></svg></div>
            <p>Kelola konten hero, syarat, FAQ, dan kontak dari halaman editor.</p>
            <button class="btn-green" @click="goToLandingSettings">Buka Editor Landing Page</button>
          </div>
        </div>
      </template>

      <!-- ── FAQ ── -->
      <template v-else-if="activeTab === 'faq'">
        <div class="card">
          <div class="card-header">
            <h3 class="card-title">Kelola FAQ</h3>
            <button class="btn-green-sm">+ Tambah FAQ</button>
          </div>
          <div class="empty-state">
            <div class="empty-state__icon"><svg width="36" height="36" viewBox="0 0 24 24" fill="none"><circle cx="12" cy="12" r="10" stroke="#d1d5db" stroke-width="1.5"/><path d="M9.09 9a3 3 0 015.83 1c0 2-3 3-3 3" stroke="#d1d5db" stroke-width="1.5" stroke-linecap="round"/><line x1="12" y1="17" x2="12.01" y2="17" stroke="#d1d5db" stroke-width="2" stroke-linecap="round"/></svg></div>
            <p>Kelola daftar pertanyaan yang sering diajukan.</p>
          </div>
        </div>
      </template>

      <!-- ── STATISTIK ── -->
      <template v-else-if="activeTab === 'statistik'">
        <div class="stats-row">
          <div class="stat-card">
            <div class="stat-card__icon" style="background:#eff6ff;color:#3b82f6"><svg width="18" height="18" viewBox="0 0 24 24" fill="none"><path d="M17 21v-2a4 4 0 00-4-4H5a4 4 0 00-4 4v2" stroke="currentColor" stroke-width="2"/><circle cx="9" cy="7" r="4" stroke="currentColor" stroke-width="2"/></svg></div>
            <div><div class="stat-card__label">Total Peserta</div><div class="stat-card__value">0</div></div>
          </div>
          <div class="stat-card">
            <div class="stat-card__icon" style="background:#f0fdf4;color:#16a34a"><svg width="18" height="18" viewBox="0 0 24 24" fill="none"><polyline points="20 6 9 17 4 12" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg></div>
            <div><div class="stat-card__label">Selesai Magang</div><div class="stat-card__value">0</div></div>
          </div>
          <div class="stat-card">
            <div class="stat-card__icon" style="background:#fefce8;color:#ca8a04"><svg width="18" height="18" viewBox="0 0 24 24" fill="none"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2" stroke="currentColor" stroke-width="2"/></svg></div>
            <div><div class="stat-card__label">Rata-rata Nilai</div><div class="stat-card__value">—</div></div>
          </div>
          <div class="stat-card">
            <div class="stat-card__icon" style="background:#fff1f2;color:#e11d48"><svg width="18" height="18" viewBox="0 0 24 24" fill="none"><path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z" stroke="currentColor" stroke-width="2"/></svg></div>
            <div><div class="stat-card__label">Sertifikat Diterbitkan</div><div class="stat-card__value">0</div></div>
          </div>
        </div>
        <div class="card">
          <div class="card-header"><h3 class="card-title">Grafik Pendaftar per Bulan</h3></div>
          <div class="empty-state">
            <div class="empty-state__icon"><svg width="36" height="36" viewBox="0 0 24 24" fill="none"><line x1="18" y1="20" x2="18" y2="10" stroke="#d1d5db" stroke-width="1.5" stroke-linecap="round"/><line x1="12" y1="20" x2="12" y2="4" stroke="#d1d5db" stroke-width="1.5" stroke-linecap="round"/><line x1="6" y1="20" x2="6" y2="14" stroke="#d1d5db" stroke-width="1.5" stroke-linecap="round"/></svg></div>
            <p>Belum ada data statistik.</p>
          </div>
        </div>
      </template>

    </template>
  </DashboardLayout>
</template>

<script setup lang="ts">
import { ref, computed } from "vue";
import { useRouter } from "vue-router";
import DashboardLayout from "@/layouts/DashboardLayout.vue";
import { useAuth } from "@/hooks/useAuth";

const { user } = useAuth();
const router = useRouter();
const layout = ref<InstanceType<typeof DashboardLayout> | null>(null);
const activeTab = computed(() => layout.value?.activeTab ?? "beranda");
const firstName = computed(() => user.value?.nama_lengkap?.split(" ")[0] ?? "");

function setTab(tab: string) {
  if (layout.value) layout.value.activeTab = tab;
}
function goToLandingSettings() {
  router.push("/admin/landing-settings");
}

const ICON = {
  home:    `<svg width="16" height="16" viewBox="0 0 24 24" fill="none"><path d="M3 9l9-7 9 7v11a2 2 0 01-2 2H5a2 2 0 01-2-2z" stroke="currentColor" stroke-width="2"/><polyline points="9 22 9 12 15 12 15 22" stroke="currentColor" stroke-width="2"/></svg>`,
  users:   `<svg width="16" height="16" viewBox="0 0 24 24" fill="none"><path d="M17 21v-2a4 4 0 00-4-4H5a4 4 0 00-4 4v2" stroke="currentColor" stroke-width="2"/><circle cx="9" cy="7" r="4" stroke="currentColor" stroke-width="2"/></svg>`,
  period:  `<svg width="16" height="16" viewBox="0 0 24 24" fill="none"><rect x="3" y="4" width="18" height="18" rx="2" ry="2" stroke="currentColor" stroke-width="2"/><line x1="3" y1="10" x2="21" y2="10" stroke="currentColor" stroke-width="2"/></svg>`,
  globe:   `<svg width="16" height="16" viewBox="0 0 24 24" fill="none"><circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2"/><line x1="2" y1="12" x2="22" y2="12" stroke="currentColor" stroke-width="2"/><path d="M12 2a15.3 15.3 0 014 10 15.3 15.3 0 01-4 10 15.3 15.3 0 01-4-10 15.3 15.3 0 014-10z" stroke="currentColor" stroke-width="2"/></svg>`,
  faq:     `<svg width="16" height="16" viewBox="0 0 24 24" fill="none"><circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2"/><path d="M9.09 9a3 3 0 015.83 1c0 2-3 3-3 3" stroke="currentColor" stroke-width="2" stroke-linecap="round"/><line x1="12" y1="17" x2="12.01" y2="17" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>`,
  stats:   `<svg width="16" height="16" viewBox="0 0 24 24" fill="none"><line x1="18" y1="20" x2="18" y2="10" stroke="currentColor" stroke-width="2" stroke-linecap="round"/><line x1="12" y1="20" x2="12" y2="4" stroke="currentColor" stroke-width="2" stroke-linecap="round"/><line x1="6" y1="20" x2="6" y2="14" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>`,
};

const navGroups = [
  {
    items: [
      { key: "beranda", label: "Beranda", icon: ICON.home },
    ],
  },
  {
    label: "Manajemen",
    items: [
      { key: "users",   label: "Manajemen User",  icon: ICON.users },
      { key: "periode", label: "Periode Magang",   icon: ICON.period },
    ],
  },
  {
    label: "Konten Web",
    items: [
      { key: "landing", label: "Landing Page", icon: ICON.globe },
      { key: "faq",     label: "Kelola FAQ",   icon: ICON.faq },
    ],
  },
  {
    label: "Laporan",
    items: [
      { key: "statistik", label: "Statistik", icon: ICON.stats },
    ],
  },
];
</script>

<style scoped>
.welcome-banner { background: linear-gradient(135deg, #0d2818 0%, #1a5c20 100%); border-radius: 14px; padding: 24px 28px; color: #fff; display: flex; align-items: center; justify-content: space-between; gap: 16px; }
.welcome-banner__greeting { font-size: 17px; font-weight: 700; }
.welcome-banner__sub { font-size: 12.5px; color: rgba(255,255,255,0.65); margin-top: 4px; }
.welcome-banner__badge { background: rgba(72,175,74,0.25); border: 1px solid rgba(134,239,172,0.4); color: #86efac; font-size: 11px; font-weight: 800; letter-spacing: 0.12em; padding: 6px 14px; border-radius: 100px; }

.stats-row { display: grid; grid-template-columns: repeat(4,1fr); gap: 14px; }
.stat-card { background: #fff; border-radius: 12px; padding: 16px; display: flex; align-items: center; gap: 12px; border: 1px solid #e9f5e9; box-shadow: 0 1px 3px rgba(13,40,24,0.05); }
.stat-card__icon { width: 38px; height: 38px; border-radius: 10px; display: flex; align-items: center; justify-content: center; flex-shrink: 0; }
.stat-card__label { font-size: 11px; color: #6b7280; font-weight: 500; }
.stat-card__value { font-size: 20px; font-weight: 700; color: #111827; margin-top: 1px; }
.stat-card__sub { font-size: 10.5px; color: #9ca3af; }

.quick-grid { display: grid; grid-template-columns: repeat(4,1fr); gap: 14px; }
.quick-card { background: #fff; border: 1.5px solid #e9f5e9; border-radius: 14px; padding: 20px 14px; display: flex; flex-direction: column; align-items: center; gap: 12px; cursor: pointer; font-family: "Poppins",sans-serif; transition: border-color 0.15s, box-shadow 0.15s; }
.quick-card:hover { border-color: #48AF4A; box-shadow: 0 4px 14px rgba(72,175,74,0.12); }
.quick-card__icon { width: 48px; height: 48px; border-radius: 12px; display: flex; align-items: center; justify-content: center; }
.quick-card__label { font-size: 12px; font-weight: 600; color: #374151; text-align: center; }

.card { background: #fff; border-radius: 14px; border: 1px solid #e9f5e9; box-shadow: 0 1px 3px rgba(13,40,24,0.05); overflow: hidden; }
.card-header { display: flex; align-items: center; justify-content: space-between; padding: 16px 20px; border-bottom: 1px solid #f0faf0; }
.card-title { font-size: 13.5px; font-weight: 700; color: #111827; margin: 0; }

.empty-state { display: flex; flex-direction: column; align-items: center; padding: 44px 24px; gap: 12px; text-align: center; }
.empty-state__icon { width: 72px; height: 72px; background: #f9fafb; border-radius: 50%; display: flex; align-items: center; justify-content: center; }
.empty-state p { font-size: 13px; color: #9ca3af; line-height: 1.7; margin: 0; }

.table-wrap { overflow-x: auto; }
.data-table { width: 100%; border-collapse: collapse; font-size: 13px; }
.data-table th { padding: 11px 16px; text-align: left; font-size: 10.5px; font-weight: 600; color: #6b7280; background: #f9fafb; border-bottom: 1px solid #f1f5f9; text-transform: uppercase; letter-spacing: 0.04em; }
.data-table td { padding: 13px 16px; border-bottom: 1px solid #f9fafb; color: #374151; }

.role-chip { padding: 3px 10px; border-radius: 100px; font-size: 11px; font-weight: 600; }
.role-chip--admin { background: #f5f3ff; color: #7c3aed; }
.role-chip--hrd { background: #eff6ff; color: #3b82f6; }
.status-chip { padding: 3px 10px; border-radius: 100px; font-size: 11px; font-weight: 600; }
.status-chip--active { background: #f0fdf4; color: #16a34a; }
.action-btn { background: none; border: 1px solid #e5e7eb; border-radius: 6px; padding: 4px 12px; font-size: 12px; font-family: "Poppins",sans-serif; cursor: pointer; color: #374151; transition: all 0.15s; }
.action-btn:hover { border-color: #48AF4A; color: #48AF4A; }

.btn-green { background: #48AF4A; color: #fff; border: none; border-radius: 9px; padding: 10px 22px; font-size: 13px; font-weight: 600; font-family: "Poppins",sans-serif; cursor: pointer; }
.btn-green-sm { background: #48AF4A; color: #fff; border: none; border-radius: 8px; padding: 6px 14px; font-size: 12px; font-weight: 600; font-family: "Poppins",sans-serif; cursor: pointer; white-space: nowrap; }

@media (max-width: 700px) { .stats-row, .quick-grid { grid-template-columns: 1fr 1fr; } }
@media (max-width: 420px) { .stats-row, .quick-grid { grid-template-columns: 1fr; } }
</style>
