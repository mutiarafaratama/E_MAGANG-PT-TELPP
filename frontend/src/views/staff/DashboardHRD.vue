<template>
  <DashboardLayout :nav-groups="navGroups" role-name="Staff HRD" default-tab="beranda" ref="layout"
    @tab-change="onTabChange">
    <template #default>

      <!-- BERANDA ──────────────────────────────────────────────── -->
      <template v-if="activeTab === 'beranda'">
        <div class="welcome-banner">
          <div>
            <div class="welcome-banner__greeting">Halo, {{ firstName }}! 👋</div>
            <div class="welcome-banner__sub">Kelola seleksi dan pelaksanaan magang dari sini.</div>
          </div>
          <div class="welcome-banner__icon">
            <svg width="36" height="36" viewBox="0 0 24 24" fill="none">
              <path d="M17 21v-2a4 4 0 00-4-4H5a4 4 0 00-4 4v2" stroke="rgba(255,255,255,0.4)" stroke-width="1.5"/>
              <circle cx="9" cy="7" r="4" stroke="rgba(255,255,255,0.4)" stroke-width="1.5"/>
              <path d="M23 21v-2a4 4 0 00-3-3.87" stroke="#86efac" stroke-width="1.5" stroke-linecap="round"/>
              <path d="M16 3.13a4 4 0 010 7.75" stroke="#86efac" stroke-width="1.5" stroke-linecap="round"/>
            </svg>
          </div>
        </div>

        <div class="stats-row">
          <div class="stat-card">
            <div class="stat-card__icon" style="background:#fff7ed;color:#ea580c">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none">
                <path d="M17 21v-2a4 4 0 00-4-4H5a4 4 0 00-4 4v2" stroke="currentColor" stroke-width="2"/>
                <circle cx="9" cy="7" r="4" stroke="currentColor" stroke-width="2"/>
              </svg>
            </div>
            <div>
              <div class="stat-card__label">Total Pengajuan</div>
              <div class="stat-card__value">{{ statsLoading ? '…' : totalAll }}</div>
            </div>
          </div>
          <div class="stat-card">
            <div class="stat-card__icon" style="background:#fefce8;color:#ca8a04">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none">
                <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2"/>
                <polyline points="12 6 12 12 16 14" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
              </svg>
            </div>
            <div>
              <div class="stat-card__label">Menunggu Verifikasi</div>
              <div class="stat-card__value">{{ statsLoading ? '…' : totalPending }}</div>
            </div>
          </div>
          <div class="stat-card">
            <div class="stat-card__icon" style="background:#f0fdf4;color:#16a34a">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none">
                <polyline points="20 6 9 17 4 12" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              </svg>
            </div>
            <div>
              <div class="stat-card__label">Diterima</div>
              <div class="stat-card__value">{{ statsLoading ? '…' : totalDiterima }}</div>
            </div>
          </div>
          <div class="stat-card">
            <div class="stat-card__icon" style="background:#fef2f2;color:#dc2626">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none">
                <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2"/>
                <line x1="15" y1="9" x2="9" y2="15" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
                <line x1="9" y1="9" x2="15" y2="15" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
              </svg>
            </div>
            <div>
              <div class="stat-card__label">Ditolak</div>
              <div class="stat-card__value">{{ statsLoading ? '…' : totalDitolak }}</div>
            </div>
          </div>
        </div>

        <!-- Perlu ditindaklanjuti -->
        <div class="card">
          <div class="card-header">
            <h3 class="card-title">Pengajuan Perlu Ditindaklanjuti</h3>
            <span class="badge badge--yellow">{{ totalPending }} Pending</span>
          </div>
          <div v-if="statsLoading" class="empty-state"><div class="spinner"></div></div>
          <div v-else-if="pendingItems.length === 0" class="empty-state">
            <div class="empty-state__icon">
              <svg width="36" height="36" viewBox="0 0 24 24" fill="none">
                <path d="M17 21v-2a4 4 0 00-4-4H5a4 4 0 00-4 4v2" stroke="#d1d5db" stroke-width="1.5"/>
                <circle cx="9" cy="7" r="4" stroke="#d1d5db" stroke-width="1.5"/>
              </svg>
            </div>
            <p>Tidak ada pengajuan yang perlu diproses saat ini.</p>
          </div>
          <div v-else class="table-wrap">
            <table class="data-table">
              <thead>
                <tr>
                  <th>Nama</th><th>Institusi</th><th>Kategori</th>
                  <th>Tanggal</th><th>Status</th><th>Aksi</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="p in pendingItems.slice(0, 5)" :key="p.id">
                  <td class="td-name">{{ p.nama_lengkap }}</td>
                  <td>{{ p.asal_institusi }}</td>
                  <td>{{ formatKategori(p.kategori_magang) }}</td>
                  <td>{{ formatDate(p.created_at) }}</td>
                  <td><span :class="statusBadgeClass(p.status)">{{ formatStatus(p.status) }}</span></td>
                  <td><button class="btn-detail" @click="goToVerifikasi">Tinjau</button></td>
                </tr>
              </tbody>
            </table>
          </div>
          <div v-if="pendingItems.length > 5" class="card-footer">
            <button class="btn-link" @click="goToVerifikasi">Lihat semua {{ pendingItems.length }} pengajuan →</button>
          </div>
        </div>
      </template>

      <!-- VERIFIKASI BERKAS ────────────────────────────────────── -->
      <template v-else-if="activeTab === 'verifikasi'">
        <HRDVerifikasi @updated="fetchPengajuan" />
      </template>

      <!-- MENUNGGU PENEMPATAN ──────────────────────────────────── -->
      <template v-else-if="activeTab === 'peserta-penempatan'">
        <HRDPenempatanMenunggu />
      </template>

      <!-- SEDANG BERLANGSUNG ───────────────────────────────────── -->
      <template v-else-if="activeTab === 'peserta-berlangsung'">
        <HRDBerlangsung />
      </template>

      <!-- PLACEHOLDER TABS ─────────────────────────────────────── -->
      <template v-else-if="activeTab === 'absensi'">
        <HRDRekapAbsen />
      </template>

      <!-- LAPORAN PESERTA ──────────────────────────────────────── -->
      <template v-else-if="activeTab === 'laporan-peserta'">
        <HRDLaporanPeserta />
      </template>

      <template v-else-if="activeTab === 'penilaian'">
        <HRDPenilaian />
      </template>

      <template v-else-if="activeTab === 'sertifikat'">
        <HRDSertifikat />
      </template>

      <template v-else-if="activeTab === 'perpanjangan'">
        <PerpanjanganHRDView />
      </template>

      <template v-else-if="activeTab === 'profil'">
        <StaffProfil />
      </template>

    </template>
  </DashboardLayout>

  <!-- Chat FAB — Kelola Tiket Peserta -->
  <ChatFAB v-model="chatOpen" :unread-count="badge.chat_unread">
    <ChatWidgetHRD />
  </ChatFAB>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from "vue";
import { useAuth } from "@/hooks/useAuth";
import api from "@/lib/api";
import { useAppWS } from "@/composables/useAppWS";
import DashboardLayout from "@/layouts/DashboardLayout.vue";
import ChatFAB from "@/components/ChatFAB.vue";
import ChatWidget from "@/components/ChatWidget.vue";
import HRDVerifikasi from "./HRDVerifikasi.vue";
import HRDPenempatanMenunggu from "./HRDPenempatanMenunggu.vue";
import HRDBerlangsung from "./HRDBerlangsung.vue";
import HRDRekapAbsen from "./HRDRekapAbsen.vue";
import HRDLaporanPeserta from "./HRDLaporanPeserta.vue";
import HRDPenilaian from "./HRDPenilaian.vue";
import HRDSertifikat from "./HRDSertifikat.vue";
import StaffProfil from "./StaffProfil.vue";
import PerpanjanganHRDView from "./PerpanjanganHRDView.vue";
import ChatWidgetHRD from "@/components/ChatWidgetHRD.vue";

// ── auth ────────────────────────────────────────────────────────────
const { user } = useAuth();
const layout    = ref(null);
const activeTab = ref("beranda");
const chatOpen  = ref(false);
const firstName = computed(() => (user.value?.nama_lengkap ?? "HRD").split(" ")[0]);

// ── sidebar badge counts ─────────────────────────────────────────────
interface HRDBadge {
  verifikasi_berkas: number;
  menunggu_penempatan: number;
  izin_sakit_pending: number;
  laporan_pending: number;
  penilaian_siap: number;
  chat_unread: number;
  perpanjangan_pending: number;
}
const badge = ref<HRDBadge>({
  verifikasi_berkas: 0, menunggu_penempatan: 0,
  izin_sakit_pending: 0, laporan_pending: 0,
  penilaian_siap: 0, chat_unread: 0,
  perpanjangan_pending: 0,
});

// seenBadge: count terakhir yang sudah "dilihat" user, persist ke localStorage
const LS_HRD = "hrd_seen_badge";
const _defaultSeen: HRDBadge = { verifikasi_berkas:0, menunggu_penempatan:0, izin_sakit_pending:0, laporan_pending:0, penilaian_siap:0, chat_unread:0, perpanjangan_pending:0 };
const seenBadge = ref<HRDBadge>((() => {
  try { return { ..._defaultSeen, ...JSON.parse(localStorage.getItem(LS_HRD) ?? "{}") }; } catch { return { ..._defaultSeen }; }
})());
function saveSeen() { localStorage.setItem(LS_HRD, JSON.stringify(seenBadge.value)); }

// Mapping tab → field badge yang relevan
const TAB_BADGE_HRD: Record<string, keyof HRDBadge> = {
  verifikasi:          "verifikasi_berkas",
  peserta:             "menunggu_penempatan",
  "peserta-penempatan":"menunggu_penempatan",
  absensi:             "izin_sakit_pending",
  "laporan-peserta":   "laporan_pending",
  penilaian:           "penilaian_siap",
  perpanjangan:        "perpanjangan_pending",
};

async function fetchBadge() {
  try {
    const r = await api.get("/api/badge/hrd");
    badge.value = r.data?.data ?? badge.value;
  } catch {}
}

// WebSocket — singleton, dipakai bersama komponen lain dalam sesi yang sama
const { connect: wsConnect, disconnect: wsDisconnect, subscribe: wsSubscribe } = useAppWS();
let wsUnsub: (() => void) | null = null;

// ── beranda stats state ─────────────────────────────────────────────
interface Pengajuan {
  id: string; nama_lengkap: string; asal_institusi: string;
  kategori_magang: string; status: string; created_at: string;
}
const allPengajuan  = ref<Pengajuan[]>([]);
const statsLoading  = ref(false);

const totalAll      = computed(() => allPengajuan.value.length);
const totalPending  = computed(() => allPengajuan.value.filter(p => ["diajukan","menunggu_verifikasi","diproses"].includes(p.status)).length);
const totalDiterima = computed(() => allPengajuan.value.filter(p => p.status === "diterima").length);
const totalDitolak  = computed(() => allPengajuan.value.filter(p => p.status === "ditolak").length);
const pendingItems  = computed(() => allPengajuan.value.filter(p => ["diajukan","menunggu_verifikasi","diproses"].includes(p.status)));

async function fetchPengajuan() {
  statsLoading.value = true;
  try {
    const res = await api.get("/api/pengajuan?page=1&limit=200");
    allPengajuan.value = Array.isArray(res.data.data) ? res.data.data : [];
  } catch {}
  finally { statsLoading.value = false; }
}

// ── nav ─────────────────────────────────────────────────────────────
const ICON = {
  home:     `<svg width="16" height="16" viewBox="0 0 24 24" fill="none"><path d="M3 9l9-7 9 7v11a2 2 0 01-2 2H5a2 2 0 01-2-2z" stroke="currentColor" stroke-width="2"/><polyline points="9 22 9 12 15 12 15 22" stroke="currentColor" stroke-width="2"/></svg>`,
  verify:   `<svg width="16" height="16" viewBox="0 0 24 24" fill="none"><path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z" stroke="currentColor" stroke-width="2"/><polyline points="14 2 14 8 20 8" stroke="currentColor" stroke-width="2"/><polyline points="9 15 11 17 15 13" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>`,
  users:    `<svg width="16" height="16" viewBox="0 0 24 24" fill="none"><path d="M17 21v-2a4 4 0 00-4-4H5a4 4 0 00-4 4v2" stroke="currentColor" stroke-width="2"/><circle cx="9" cy="7" r="4" stroke="currentColor" stroke-width="2"/></svg>`,
  calendar: `<svg width="16" height="16" viewBox="0 0 24 24" fill="none"><rect x="3" y="4" width="18" height="18" rx="2" ry="2" stroke="currentColor" stroke-width="2"/><line x1="3" y1="10" x2="21" y2="10" stroke="currentColor" stroke-width="2"/></svg>`,
  star:     `<svg width="16" height="16" viewBox="0 0 24 24" fill="none"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2" stroke="currentColor" stroke-width="2"/></svg>`,
  award:    `<svg width="16" height="16" viewBox="0 0 24 24" fill="none"><circle cx="12" cy="8" r="6" stroke="currentColor" stroke-width="2"/><path d="M15.477 12.89L17 22l-5-3-5 3 1.523-9.11" stroke="currentColor" stroke-width="2"/></svg>`,
  chat:     `<svg width="16" height="16" viewBox="0 0 24 24" fill="none"><path d="M21 15a2 2 0 01-2 2H7l-4 4V5a2 2 0 012-2h14a2 2 0 012 2z" stroke="currentColor" stroke-width="2"/></svg>`,
  laporan:  `<svg width="16" height="16" viewBox="0 0 24 24" fill="none"><path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z" stroke="currentColor" stroke-width="2" stroke-linejoin="round"/><polyline points="14 2 14 8 20 8" stroke="currentColor" stroke-width="2" stroke-linejoin="round"/><line x1="16" y1="13" x2="8" y2="13" stroke="currentColor" stroke-width="2" stroke-linecap="round"/><line x1="16" y1="17" x2="8" y2="17" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>`,
};

// dot helper: dot muncul jika count saat ini lebih besar dari yang terakhir dilihat
function d(cur: number, seen: number) { return cur > seen; }

const navGroups = computed(() => [
  { items: [{ key:"beranda", label:"Beranda", icon: ICON.home }] },
  {
    label: "Seleksi & Verifikasi",
    items: [
      {
        key:"verifikasi", label:"Verifikasi Berkas", icon: ICON.verify,
        dot: badge.value.verifikasi_berkas > 0,
      },
      {
        key:"peserta", label:"Peserta Aktif", icon: ICON.users,
        dot: d(badge.value.menunggu_penempatan, seenBadge.value.menunggu_penempatan),
        children: [
          { key:"peserta-penempatan", label:"Menunggu Penempatan", icon:"", dot: d(badge.value.menunggu_penempatan, seenBadge.value.menunggu_penempatan) },
          { key:"peserta-berlangsung", label:"Sedang Berlangsung", icon:"" },
        ],
      },
    ],
  },
  {
    label: "Pelaksanaan",
    items: [
      { key:"absensi",         label:"Rekap Absensi",     icon: ICON.calendar, dot: d(badge.value.izin_sakit_pending,    seenBadge.value.izin_sakit_pending) },
      { key:"laporan-peserta", label:"Laporan Peserta",   icon: ICON.laporan,  dot: d(badge.value.laporan_pending,       seenBadge.value.laporan_pending) },
      { key:"penilaian",       label:"Penilaian Peserta", icon: ICON.star,     dot: d(badge.value.penilaian_siap,        seenBadge.value.penilaian_siap) },
      { key:"sertifikat",      label:"Sertifikat",        icon: ICON.award },
      { key:"perpanjangan",    label:"Perpanjangan",      icon: ICON.calendar, dot: d(badge.value.perpanjangan_pending, seenBadge.value.perpanjangan_pending) },
    ],
  },
]);

// ── tab change — simpan count yang sudah dilihat ke localStorage ────
function markSeen(tab: string) {
  const key = TAB_BADGE_HRD[tab];
  if (key) {
    seenBadge.value = { ...seenBadge.value, [key]: badge.value[key] };
    saveSeen();
  }
}

function onTabChange(tab: string) {
  activeTab.value = tab;
  markSeen(tab);
  fetchBadge();
  if (tab === "beranda") fetchPengajuan();
}

function goToVerifikasi() {
  activeTab.value = "verifikasi";
  markSeen("verifikasi");
  if (layout.value) (layout.value as any).activeTab.value = "verifikasi";
}

// ── format helpers (beranda table) ──────────────────────────────────
function formatDate(d: string) {
  return new Date(d).toLocaleDateString("id-ID", { day:"2-digit", month:"short", year:"numeric" });
}
function formatStatus(s: string) {
  return ({ diajukan:"Diajukan", menunggu_verifikasi:"Menunggu Verifikasi", diproses:"Diproses", diterima:"Diterima", ditolak:"Ditolak" } as Record<string,string>)[s] ?? s;
}
function statusBadgeClass(s: string) {
  const base = "status-badge ";
  if (s === "diterima") return base + "status-badge--green";
  if (s === "ditolak")  return base + "status-badge--red";
  if (s === "diproses") return base + "status-badge--blue";
  return base + "status-badge--yellow";
}
function formatKategori(k: string) {
  return ({ siswa_smk:"Siswa SMK", mahasiswa_d3:"Mahasiswa D3", mahasiswa_s1:"Mahasiswa S1", penelitian:"Penelitian" } as Record<string,string>)[k] ?? k;
}

onMounted(() => {
  fetchPengajuan();
  fetchBadge();
  wsConnect();
  wsUnsub = wsSubscribe((msg) => {
    if (msg.type === "badge_update" || msg.type === "notifikasi" || msg.type === "chat_message" || msg.type === "tiket_update") {
      fetchBadge();
    }
  });
  setInterval(fetchBadge, 30_000);
});

onUnmounted(() => {
  wsUnsub?.();
  wsDisconnect();
});
</script>

<style scoped>
/* ── welcome ── */
.welcome-banner { background:linear-gradient(135deg,#0d2818 0%,#1a5c20 100%); border-radius:14px; padding:24px 28px; color:#fff; display:flex; align-items:center; justify-content:space-between; gap:16px; }
.welcome-banner__greeting { font-size:17px; font-weight:700; }
.welcome-banner__sub { font-size:12.5px; color:rgba(255,255,255,0.65); margin-top:4px; }
.welcome-banner__icon { opacity:0.8; }
/* ── stats ── */
.stats-row { display:grid; grid-template-columns:repeat(4,1fr); gap:14px; }
.stat-card { background:#fff; border-radius:12px; padding:16px; display:flex; align-items:center; gap:12px; border:1px solid #e9f5e9; box-shadow:0 1px 3px rgba(13,40,24,0.05); }
.stat-card__icon { width:38px; height:38px; border-radius:10px; display:flex; align-items:center; justify-content:center; flex-shrink:0; }
.stat-card__label { font-size:11px; color:#6b7280; font-weight:500; }
.stat-card__value { font-size:20px; font-weight:700; color:#111827; margin-top:1px; }
/* ── card ── */
.card { background:#fff; border-radius:14px; border:1px solid #e9f5e9; box-shadow:0 1px 3px rgba(13,40,24,0.05); overflow:hidden; }
.card-header { display:flex; align-items:center; justify-content:space-between; padding:16px 20px; border-bottom:1px solid #f0faf0; gap:12px; flex-wrap:wrap; }
.card-title { font-size:13.5px; font-weight:700; color:#111827; margin:0; }
.card-footer { padding:12px 20px; border-top:1px solid #f3f4f6; text-align:center; }
.btn-link { background:none; border:none; color:#48AF4A; font-size:12.5px; font-weight:600; cursor:pointer; font-family:inherit; }
.btn-link:hover { text-decoration:underline; }
/* ── table ── */
.table-wrap { overflow-x:auto; }
.data-table { width:100%; border-collapse:collapse; font-size:13px; }
.data-table th { padding:11px 16px; text-align:left; font-size:10.5px; font-weight:600; color:#6b7280; background:#f9fafb; border-bottom:1px solid #f1f5f9; text-transform:uppercase; letter-spacing:0.04em; white-space:nowrap; }
.data-table td { padding:13px 16px; border-bottom:1px solid #f9fafb; color:#374151; vertical-align:middle; }
.td-name { min-width:140px; }
/* ── btn ── */
.btn-detail { background:#f3f4f6; color:#374151; border:none; border-radius:7px; padding:5px 11px; font-size:11.5px; font-weight:600; font-family:inherit; cursor:pointer; }
.btn-detail:hover { background:#e5e7eb; }
.btn-green-sm { background:#48AF4A; color:#fff; border:none; border-radius:8px; padding:6px 14px; font-size:12px; font-weight:600; font-family:inherit; cursor:pointer; white-space:nowrap; display:flex; align-items:center; gap:5px; }
.btn-green-sm:hover { background:#48AF4A; }
/* ── badges ── */
.badge { padding:3px 10px; border-radius:100px; font-size:11px; font-weight:600; }
.badge--yellow { background:#fefce8; color:#ca8a04; }
.badge--gray   { background:#f3f4f6; color:#6b7280; }
.status-badge { padding:3px 10px; border-radius:100px; font-size:11px; font-weight:600; white-space:nowrap; }
.status-badge--yellow { background:#fefce8; color:#ca8a04; }
.status-badge--green  { background:#f0fdf4; color:#16a34a; }
.status-badge--red    { background:#fef2f2; color:#dc2626; }
.status-badge--blue   { background:#f0fdf4; color:#1a5c20; }
/* ── empty / spinner ── */
.empty-state { display:flex; flex-direction:column; align-items:center; padding:44px 24px; gap:12px; text-align:center; }
.empty-state__icon { width:72px; height:72px; background:#f9fafb; border-radius:50%; display:flex; align-items:center; justify-content:center; }
.empty-state p { font-size:13px; color:#9ca3af; line-height:1.7; margin:0; }
.spinner { width:24px; height:24px; border:2.5px solid #e5e7eb; border-top-color:#48AF4A; border-radius:50%; animation:spin 0.7s linear infinite; }
@keyframes spin { to { transform:rotate(360deg); } }
@media (max-width:700px) { .stats-row { grid-template-columns:1fr 1fr; } }
@media (max-width:420px) { .stats-row { grid-template-columns:1fr; } }
</style>
