<template>
  <DashboardLayout :nav-groups="navGroups" role-name="Peserta Magang" default-tab="beranda" ref="layout"
    @tab-change="onTabChange">
    <template #default>

      <!-- ══════════════════════════════════════════════════════════ -->
      <!--  BERANDA                                                   -->
      <!-- ══════════════════════════════════════════════════════════ -->
      <template v-if="activeTab === 'beranda'">

        <!-- ── Welcome Banner ────────────────────────────────────── -->
        <div class="welcome-banner">
          <div class="wb-left">
            <div class="wb-greeting">Selamat datang, {{ firstName }}!</div>
            <div class="wb-sub">{{ todayFull }}</div>
          </div>
          <div v-if="pelaksanaan && pelaksanaan.status === 'aktif'" class="wb-right">
            <template v-if="sisaHari > 0">
              <div class="wb-num">{{ sisaHari }}</div>
              <div class="wb-lbl">hari&nbsp;lagi</div>
            </template>
            <template v-else>
              <div class="wb-num">🎉</div>
              <div class="wb-lbl">Selesai hari ini!</div>
            </template>
            <div class="wb-progress-wrap">
              <div class="wb-progress-bar" :style="{ width: progressMagang + '%' }"></div>
            </div>
            <div class="wb-hari-ke">Hari ke-{{ hariKe }} dari {{ totalHariMagang }}</div>
          </div>
        </div>

        <!-- ── Panduan Pengguna ───────────────────────────────────── -->
        <div class="panduan-card">
          <div class="panduan-card__head">
            <div class="panduan-card__icon">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none"><circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2"/><line x1="12" y1="8" x2="12" y2="12" stroke="currentColor" stroke-width="2" stroke-linecap="round"/><circle cx="12" cy="16" r="1" fill="currentColor"/></svg>
            </div>
            <span class="panduan-card__title">Panduan Penggunaan Sistem</span>
          </div>
          <div class="panduan-steps">
            <div class="panduan-step">
              <div class="panduan-step__num">1</div>
              <div class="panduan-step__body">
                <div class="panduan-step__label">Upload Dokumen Pendaftaran</div>
                <div class="panduan-step__desc">Unggah semua berkas yang diminta (proposal, surat pengantar, KTP, dll.) melalui menu <strong>Dokumen</strong>.</div>
              </div>
            </div>
            <div class="panduan-step">
              <div class="panduan-step__num">2</div>
              <div class="panduan-step__body">
                <div class="panduan-step__label">Tunggu Verifikasi HRD</div>
                <div class="panduan-step__desc">HRD akan memverifikasi berkas dan mengunggah surat balasan. Cek notifikasi secara berkala.</div>
              </div>
            </div>
            <div class="panduan-step">
              <div class="panduan-step__num">3</div>
              <div class="panduan-step__body">
                <div class="panduan-step__label">Absensi Harian</div>
                <div class="panduan-step__desc">Saat magang aktif, lakukan absen masuk dan pulang setiap hari melalui menu <strong>Absensi</strong>.</div>
              </div>
            </div>
            <div class="panduan-step">
              <div class="panduan-step__num">4</div>
              <div class="panduan-step__body">
                <div class="panduan-step__label">Upload Laporan & Penilaian</div>
                <div class="panduan-step__desc">Di akhir periode, upload laporan akhir magang. Nilai dan sertifikat akan diterbitkan oleh HRD.</div>
              </div>
            </div>
          </div>
        </div>

      </template>

      <!-- ── DOKUMEN ── -->
      <template v-else-if="activeTab === 'dokumen'">
        <PesertaDokumen />
      </template>

      <!-- ── ABSENSI ── -->
      <template v-else-if="activeTab === 'absensi'">
        <PesertaAbsensi />
      </template>

      <!-- ── DOKUMEN PENDAFTARAN ── -->
      <template v-else-if="activeTab === 'dokumen_pendaftaran'">
        <PesertaDokumen />
      </template>

      <!-- ── SURAT BALASAN ── -->
      <template v-else-if="activeTab === 'surat_balasan'">
        <SuratBalasanView :pelaksanaan="pelaksanaan" />
      </template>

      <!-- ── LAPORAN MAGANG ── -->
      <template v-else-if="activeTab === 'laporan_magang'">
        <LaporanMagangView :pelaksanaan="pelaksanaan" @refresh="fetchPelaksanaan" />
      </template>

      <!-- ── PERPANJANGAN ── -->
      <template v-else-if="activeTab === 'perpanjangan'">
        <PerpanjanganView :pelaksanaan="pelaksanaan" />
      </template>

      <!-- ── NILAI ── -->
      <template v-else-if="activeTab === 'nilai'">
        <NilaiView :pelaksanaan="pelaksanaan" :pengajuan="pengajuan" />
      </template>

      <!-- ── SERTIFIKAT ── -->
      <template v-else-if="activeTab === 'sertifikat'">
        <SertifikatView :pelaksanaan="pelaksanaan" :pengajuan="pengajuan" />
      </template>

      <!-- ── PROFIL ── -->
      <template v-else-if="activeTab === 'profil'">
        <PesertaProfil />
      </template>

    </template>
  </DashboardLayout>

  <!-- Chat FAB — Komunikasi HRD -->
  <ChatFAB v-model="chatOpen" :unread-count="badgeP.chat_unread">
    <ChatWidget />
  </ChatFAB>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted } from "vue";
import DashboardLayout from "@/layouts/DashboardLayout.vue";
import { useAuth } from "@/hooks/useAuth";
import api from "@/lib/api";
import { useAppWS } from "@/composables/useAppWS";
import ChatFAB from "@/components/ChatFAB.vue";
import ChatWidget from "@/components/ChatWidget.vue";
import PesertaDokumen    from "./PesertaDokumen.vue";
import PesertaAbsensi   from "./PesertaAbsensi.vue";
import PesertaProfil    from "./PesertaProfil.vue";
import SuratBalasanView from "./SuratBalasanView.vue";
import LaporanMagangView from "./LaporanMagangView.vue";
import PerpanjanganView from "./PerpanjanganView.vue";
import NilaiView        from "./NilaiView.vue";
import SertifikatView   from "./SertifikatView.vue";

// ── Auth + layout ref ─────────────────────────────────────────
const { user } = useAuth();
const layout    = ref<InstanceType<typeof DashboardLayout> | null>(null);
const activeTab = ref("beranda");
const chatOpen  = ref(false);
const firstName = computed(() => user.value?.nama_lengkap?.split(" ")[0] ?? "");

// ── Today ─────────────────────────────────────────────────────
const now = new Date();
const todayISO  = now.toISOString().slice(0, 10);
const todayFull = now.toLocaleDateString("id-ID", { weekday: "long", day: "numeric", month: "long", year: "numeric" });
const bulanIniLabel = now.toLocaleDateString("id-ID", { month: "long", year: "numeric" });

// ── Data refs ──────────────────────────────────────────────────
const pelaksanaan   = ref<any>(null);
const pengajuan     = ref<any>(null);
const absensiList   = ref<any[]>([]);
const absensiLoading = ref(false);
const dokumenList   = ref<any[]>([]);
const recentNotif   = ref<any[]>([]);
const notifLoading  = ref(false);
const aksiLoading   = ref(false);
const aksiMsg       = ref<{ ok: boolean; text: string } | null>(null);

// ── Status magang ──────────────────────────────────────────────
const statusLabel = computed(() => {
  const m: Record<string, string> = {
    aktif: "Aktif", selesai: "Selesai", upload_laporan: "Upload Laporan",
    penilaian: "Penilaian", none: "Belum Dimulai",
  };
  return m[pelaksanaan.value?.status ?? "none"] ?? (pelaksanaan.value?.status ?? "—");
});

// ── Countdown ─────────────────────────────────────────────────
const sisaHari = computed(() => {
  if (!pelaksanaan.value?.tanggal_selesai) return 0;
  const selesai = new Date(pelaksanaan.value.tanggal_selesai);
  return Math.max(0, Math.ceil((selesai.getTime() - now.getTime()) / 86400000));
});
const hariKe = computed(() => {
  if (!pelaksanaan.value?.tanggal_mulai) return 0;
  const mulai = new Date(pelaksanaan.value.tanggal_mulai);
  return Math.max(1, Math.ceil((now.getTime() - mulai.getTime()) / 86400000));
});
const totalHariMagang = computed(() => {
  if (!pelaksanaan.value?.tanggal_mulai || !pelaksanaan.value?.tanggal_selesai) return 0;
  const mulai   = new Date(pelaksanaan.value.tanggal_mulai);
  const selesai = new Date(pelaksanaan.value.tanggal_selesai);
  return Math.ceil((selesai.getTime() - mulai.getTime()) / 86400000);
});
const progressMagang = computed(() => {
  if (!totalHariMagang.value) return 0;
  return Math.min(100, Math.round((hariKe.value / totalHariMagang.value) * 100));
});

// ── Absensi hari ini ──────────────────────────────────────────
const todayAbsensi = computed(() =>
  absensiList.value.find(a => a.tanggal?.slice(0, 10) === todayISO) ?? null
);

// ── Rekap bulan ini ───────────────────────────────────────────
const thisYear  = now.getFullYear();
const thisMonth = now.getMonth(); // 0-indexed

const rekapBulan = computed(() => {
  const list = absensiList.value.filter(a => {
    const d = new Date(a.tanggal);
    return d.getFullYear() === thisYear && d.getMonth() === thisMonth;
  });
  return {
    hadir: list.filter(a => a.keterangan === "hadir").length,
    izin:  list.filter(a => a.keterangan === "izin").length,
    sakit: list.filter(a => a.keterangan === "sakit").length,
    alpha: list.filter(a => a.keterangan === "alpha").length,
  };
});

const totalHariKerja = computed(() => {
  const r = rekapBulan.value;
  return r.hadir + r.izin + r.sakit + r.alpha;
});

const kehadiranPct = computed(() => {
  const r = rekapBulan.value;
  const total = r.hadir + r.izin + r.sakit + r.alpha;
  if (!total) return 0;
  return Math.round((r.hadir / total) * 100);
});

// ── Mini Kalender ─────────────────────────────────────────────
const calendarCells = computed(() => {
  const year  = thisYear;
  const month = thisMonth;

  const absensiMap: Record<string, string> = {};
  for (const a of absensiList.value) {
    const d = a.tanggal?.slice(0, 10);
    if (d) absensiMap[d] = a.keterangan;
  }

  const firstDay  = new Date(year, month, 1);
  const daysInMonth = new Date(year, month + 1, 0).getDate();
  // Monday-based: Monday=0, ..., Friday=4
  let offset = firstDay.getDay() - 1; // getDay(): Sun=0,Mon=1,...
  if (offset < 0) offset = 6; // Sunday → treat as end of prev week

  const cells: any[] = [];
  // Leading empty cells (only up to 4, for Mon–Fri)
  for (let i = 0; i < offset && i < 5; i++) {
    cells.push({ key: `e-${i}`, d: "", cls: "mc-cell--empty", title: "" });
  }

  for (let d = 1; d <= daysInMonth; d++) {
    const date = new Date(year, month, d);
    const dow  = date.getDay(); // 0=Sun, 6=Sat
    if (dow === 0 || dow === 6) continue; // skip weekends

    const iso   = `${year}-${String(month + 1).padStart(2, "0")}-${String(d).padStart(2, "0")}`;
    const ket   = absensiMap[iso];
    const isFuture = iso > todayISO;
    const isToday  = iso === todayISO;

    let cls = "mc-cell--future";
    let title = iso;
    if (isToday && !ket) { cls = "mc-cell--today"; title = "Hari ini"; }
    else if (ket === "hadir")  { cls = "mc-cell--hadir";  title = `${iso}: Hadir`; }
    else if (ket === "izin")   { cls = "mc-cell--izin";   title = `${iso}: Izin`; }
    else if (ket === "sakit")  { cls = "mc-cell--sakit";  title = `${iso}: Sakit`; }
    else if (ket === "alpha")  { cls = "mc-cell--alpha";  title = `${iso}: Alpha`; }
    else if (!isFuture)        { cls = "mc-cell--kosong"; title = `${iso}: Tidak ada catatan`; }

    cells.push({ key: iso, d, cls, title });
  }

  return cells;
});

// ── Dokumen status ────────────────────────────────────────────
const DOK_ITEMS = [
  { jenis: "proposal_magang",  label: "Proposal Magang" },
  { jenis: "surat_pengantar",  label: "Surat Pengantar" },
  { jenis: "ktp",              label: "KTP" },
  { jenis: "ktm",              label: "KTM" },
  { jenis: "pasfoto",          label: "Pas Foto" },
  { jenis: "bpjs_kis",         label: "BPJS / KIS" },
];

const dokumenStatus = computed(() =>
  DOK_ITEMS.map(item => ({
    ...item,
    uploaded: dokumenList.value.some(d => d.jenis === item.jenis),
  }))
);

const dokUploaded = computed(() => dokumenStatus.value.filter(d => d.uploaded).length);

// ── Data fetch ────────────────────────────────────────────────
async function fetchPelaksanaan() {
  try {
    const r = await api.get("/api/pelaksanaan/saya");
    pelaksanaan.value = r.data?.data ?? null;
  } catch { pelaksanaan.value = null; }
}

watch(activeTab, (tab) => {
  if (['laporan_magang', 'perpanjangan', 'nilai', 'sertifikat', 'surat_balasan'].includes(tab)) fetchPelaksanaan();
});

async function fetchBeranda() {
  absensiLoading.value = true;
  notifLoading.value   = true;

  const [r1, r2, r3, r4] = await Promise.allSettled([
    api.get("/api/pelaksanaan/saya"),
    api.get("/api/pengajuan/saya"),
    api.get("/api/absensi/saya"),
    api.get("/api/notifikasi"),
  ]);

  pelaksanaan.value = r1.status === "fulfilled" ? (r1.value.data?.data ?? null) : null;
  pengajuan.value   = r2.status === "fulfilled" ? (r2.value.data?.data ?? null) : null;
  absensiList.value = r3.status === "fulfilled"
    ? (Array.isArray(r3.value.data?.data) ? r3.value.data.data : []) : [];
  recentNotif.value = r4.status === "fulfilled"
    ? (Array.isArray(r4.value.data?.data) ? r4.value.data.data.slice(0, 5) : []) : [];

  absensiLoading.value = false;
  notifLoading.value   = false;

  // Fetch dokumen setelah dapat pengajuan_id
  const pengajuanID = pelaksanaan.value?.pengajuan_id ?? pengajuan.value?.id;
  if (pengajuanID) {
    try {
      const rd = await api.get(`/api/dokumen/pengajuan/${pengajuanID}`);
      dokumenList.value = Array.isArray(rd.data?.data) ? rd.data.data : [];
    } catch { dokumenList.value = []; }
  }
}

// ── Check-In / Check-Out ──────────────────────────────────────
async function doCheckIn() {
  aksiLoading.value = true; aksiMsg.value = null;
  try {
    await api.post("/api/absensi/checkin", {});
    aksiMsg.value = { ok: true, text: "Check-in berhasil! Selamat bekerja." };
    await fetchBeranda();
  } catch (e: any) {
    aksiMsg.value = { ok: false, text: e.response?.data?.message ?? "Gagal check-in" };
  } finally {
    aksiLoading.value = false;
    setTimeout(() => { aksiMsg.value = null; }, 4000);
  }
}

async function doCheckOut() {
  aksiLoading.value = true; aksiMsg.value = null;
  try {
    await api.patch("/api/absensi/checkout", {});
    aksiMsg.value = { ok: true, text: "Check-out berhasil! Sampai jumpa besok." };
    await fetchBeranda();
  } catch (e: any) {
    aksiMsg.value = { ok: false, text: e.response?.data?.message ?? "Gagal check-out" };
  } finally {
    aksiLoading.value = false;
    setTimeout(() => { aksiMsg.value = null; }, 4000);
  }
}

// ── Tab switch ────────────────────────────────────────────────
function switchTab(tab: string) {
  activeTab.value = tab;
  if (layout.value) layout.value.activeTab = tab;
}

function onTabChange(tab: string) {
  activeTab.value = tab;
  // Simpan count yang sudah dilihat ke localStorage
  const key = TAB_BADGE_PESERTA[tab];
  if (key) {
    seenBadgeP.value = { ...seenBadgeP.value, [key]: badgeP.value[key] };
    saveSeenP();
  }
  fetchBadgePeserta();
  if (tab === "beranda") fetchBeranda();
}

// ── Helpers ───────────────────────────────────────────────────
function fmtDate(iso: string | null | undefined) {
  if (!iso) return "—";
  return new Date(iso).toLocaleDateString("id-ID", { day: "numeric", month: "short", year: "numeric" });
}

function fmtRelative(iso: string) {
  if (!iso) return "";
  const diff = Math.floor((Date.now() - new Date(iso).getTime()) / 1000);
  if (diff < 60) return "Baru saja";
  if (diff < 3600) return `${Math.floor(diff / 60)} menit lalu`;
  if (diff < 86400) return `${Math.floor(diff / 3600)} jam lalu`;
  return `${Math.floor(diff / 86400)} hari lalu`;
}

// ── Nav ───────────────────────────────────────────────────────
const ICON = {
  home:        `<svg width="16" height="16" viewBox="0 0 24 24" fill="none"><path d="M3 9l9-7 9 7v11a2 2 0 01-2 2H5a2 2 0 01-2-2z" stroke="currentColor" stroke-width="2" stroke-linejoin="round"/><polyline points="9 22 9 12 15 12 15 22" stroke="currentColor" stroke-width="2" stroke-linejoin="round"/></svg>`,
  calendar:    `<svg width="16" height="16" viewBox="0 0 24 24" fill="none"><rect x="3" y="4" width="18" height="18" rx="2" stroke="currentColor" stroke-width="2"/><line x1="16" y1="2" x2="16" y2="6" stroke="currentColor" stroke-width="2"/><line x1="8" y1="2" x2="8" y2="6" stroke="currentColor" stroke-width="2"/><line x1="3" y1="10" x2="21" y2="10" stroke="currentColor" stroke-width="2"/></svg>`,
  folder:      `<svg width="16" height="16" viewBox="0 0 24 24" fill="none"><path d="M22 19a2 2 0 01-2 2H4a2 2 0 01-2-2V5a2 2 0 012-2h5l2 3h9a2 2 0 012 2z" stroke="currentColor" stroke-width="2" stroke-linejoin="round"/></svg>`,
  award:       `<svg width="16" height="16" viewBox="0 0 24 24" fill="none"><circle cx="12" cy="8" r="6" stroke="currentColor" stroke-width="2"/><path d="M8.21 13.89L7 23l5-3 5 3-1.21-9.12" stroke="currentColor" stroke-width="2" stroke-linejoin="round"/></svg>`,
  certificate: `<svg width="16" height="16" viewBox="0 0 24 24" fill="none"><rect x="2" y="3" width="20" height="14" rx="2" stroke="currentColor" stroke-width="2"/><path d="M8 21h8M12 17v4" stroke="currentColor" stroke-width="2" stroke-linecap="round"/><path d="M7 8h5M7 11h3" stroke="currentColor" stroke-width="2" stroke-linecap="round"/><circle cx="16.5" cy="9.5" r="2.5" stroke="currentColor" stroke-width="1.5"/></svg>`,
  chat:        `<svg width="16" height="16" viewBox="0 0 24 24" fill="none"><path d="M21 15a2 2 0 01-2 2H7l-4 4V5a2 2 0 012-2h14a2 2 0 012 2z" stroke="currentColor" stroke-width="2"/></svg>`,
  profil:      `<svg width="16" height="16" viewBox="0 0 24 24" fill="none"><path d="M20 21v-2a4 4 0 00-4-4H8a4 4 0 00-4 4v2" stroke="currentColor" stroke-width="2"/><circle cx="12" cy="7" r="4" stroke="currentColor" stroke-width="2"/></svg>`,
};

// ── sidebar badge counts ─────────────────────────────────────────────
interface PesertaBadge {
  laporan_revisi: number;
  nilai_baru: number;
  sertifikat_baru: number;
  chat_unread: number;
}
const badgeP = ref<PesertaBadge>({ laporan_revisi: 0, nilai_baru: 0, sertifikat_baru: 0, chat_unread: 0 });

// seenBadge: count terakhir yang sudah "dilihat" user, persist ke localStorage
const LS_PESERTA = "peserta_seen_badge";
const _defaultSeenP: PesertaBadge = { laporan_revisi:0, nilai_baru:0, sertifikat_baru:0, chat_unread:0 };
const seenBadgeP = ref<PesertaBadge>((() => {
  try { return { ..._defaultSeenP, ...JSON.parse(localStorage.getItem(LS_PESERTA) ?? "{}") }; } catch { return { ..._defaultSeenP }; }
})());
function saveSeenP() { localStorage.setItem(LS_PESERTA, JSON.stringify(seenBadgeP.value)); }

// Mapping tab → field badge yang relevan
const TAB_BADGE_PESERTA: Record<string, keyof PesertaBadge> = {
  dokumen_saya:  "laporan_revisi",
  laporan_magang:"laporan_revisi",
  nilai:         "nilai_baru",
  sertifikat:    "sertifikat_baru",
};

async function fetchBadgePeserta() {
  try {
    const r = await api.get("/api/badge/peserta");
    badgeP.value = r.data?.data ?? badgeP.value;
  } catch {}
}

// WebSocket — singleton, dipakai bersama komponen lain dalam sesi yang sama
const { connect: wsConnect, disconnect: wsDisconnect, subscribe: wsSubscribe } = useAppWS();
let wsUnsub: (() => void) | null = null;

// dot helper: dot muncul jika count saat ini lebih besar dari yang terakhir dilihat
function dp(cur: number, seen: number) { return cur > seen; }

const navGroups = computed(() => [
  {
    items: [
      { key: "beranda", label: "Beranda", icon: ICON.home },
    ],
  },
  {
    label: "Pelaksanaan Magang",
    items: [
      { key: "absensi", label: "Absensi Harian", icon: ICON.calendar },
      {
        key: "dokumen_saya",
        label: "Dokumen Saya",
        icon: ICON.folder,
        dot: dp(badgeP.value.laporan_revisi, seenBadgeP.value.laporan_revisi),
        children: [
          { key: "dokumen_pendaftaran", label: "Dokumen Pendaftaran", icon: "" },
          { key: "surat_balasan",       label: "Surat Balasan HRD",   icon: "" },
          { key: "laporan_magang",      label: "Laporan Magang",      icon: "", dot: dp(badgeP.value.laporan_revisi, seenBadgeP.value.laporan_revisi) },
          { key: "perpanjangan",       label: "Perpanjangan Magang", icon: "" },
        ],
      },
    ],
  },
  {
    label: "Hasil Magang",
    items: [
      { key: "nilai",      label: "Nilai Akhir",  icon: ICON.award,        dot: dp(badgeP.value.nilai_baru,      seenBadgeP.value.nilai_baru) },
      { key: "sertifikat", label: "Sertifikat",   icon: ICON.certificate,  dot: dp(badgeP.value.sertifikat_baru, seenBadgeP.value.sertifikat_baru) },
    ],
  },
]);

onMounted(() => {
  fetchBeranda();
  fetchBadgePeserta();
  wsConnect();
  wsUnsub = wsSubscribe((msg) => {
    if (msg.type === "badge_update" || msg.type === "notifikasi" || msg.type === "chat_message" || msg.type === "tiket_update") {
      fetchBadgePeserta();
    }
  });
  setInterval(fetchBadgePeserta, 30_000);
});

onUnmounted(() => {
  wsUnsub?.();
  wsDisconnect();
});
</script>

<style scoped>
/* ── Welcome Banner ─────────────────────────────────────────── */
/* ── Panduan Card ─────────────────────────────────────────────── */
.panduan-card {
  margin-top: 14px;
  background: #fff; border: 1.5px solid #e9f5e9; border-radius: 14px;
  overflow: hidden; box-shadow: 0 1px 3px rgba(13,40,24,.04);
}
.panduan-card__head {
  display: flex; align-items: center; gap: 10px;
  padding: 14px 18px; border-bottom: 1px solid #f0fdf4;
  background: #fafffe;
}
.panduan-card__icon {
  width: 30px; height: 30px; border-radius: 8px;
  background: #f0fdf4; color: #0d2818;
  display: flex; align-items: center; justify-content: center; flex-shrink: 0;
}
.panduan-card__title {
  font-size: 13.5px; font-weight: 700; color: #111827;
}
.panduan-steps {
  padding: 14px 18px; display: flex; flex-direction: column; gap: 14px;
}
.panduan-step {
  display: flex; align-items: flex-start; gap: 14px;
}
.panduan-step__num {
  width: 26px; height: 26px; border-radius: 50%; flex-shrink: 0;
  background: #0d2818; color: #fff;
  font-size: 12px; font-weight: 800;
  display: flex; align-items: center; justify-content: center;
  margin-top: 1px;
}
.panduan-step__body { flex: 1; }
.panduan-step__label {
  font-size: 13px; font-weight: 700; color: #111827; margin-bottom: 3px;
}
.panduan-step__desc {
  font-size: 12px; color: #6b7280; line-height: 1.6;
}
.panduan-step__desc strong { color: #374151; }

.welcome-banner {
  background: linear-gradient(135deg, #0d2818 0%, #1a5c20 60%, #2d7a2e 100%);
  border-radius: 14px; padding: 22px 28px; color: #fff;
  display: flex; align-items: center; justify-content: space-between; gap: 16px;
}
.wb-left {}
.wb-greeting { font-size: 18px; font-weight: 700; }
.wb-sub { font-size: 12px; color: rgba(255,255,255,.55); margin-top: 4px; }
.wb-right { text-align: right; min-width: 130px; }
.wb-num { font-size: 36px; font-weight: 800; color: #86efac; line-height: 1; }
.wb-lbl { font-size: 11px; color: rgba(255,255,255,.65); margin-top: 2px; }
.wb-progress-wrap {
  height: 4px; background: rgba(255,255,255,.15); border-radius: 100px; margin: 8px 0 4px; overflow: hidden;
}
.wb-progress-bar { height: 100%; background: #86efac; border-radius: 100px; transition: width .6s; }
.wb-hari-ke { font-size: 10.5px; color: rgba(255,255,255,.5); }

/* ── Top row ─────────────────────────────────────────────────── */
.top-row { display: grid; grid-template-columns: 1fr 1fr; gap: 14px; }

/* ── Card base ───────────────────────────────────────────────── */
.card { background: #fff; border-radius: 14px; border: 1px solid #e9f5e9; box-shadow: 0 1px 3px rgba(13,40,24,.05); overflow: hidden; }
.card-header { display: flex; align-items: center; justify-content: space-between; padding: 14px 18px; border-bottom: 1px solid #f0faf0; }
.card-title { font-size: 13px; font-weight: 700; color: #111827; margin: 0; }

/* ── Status Card ─────────────────────────────────────────────── */
.status-card { padding: 18px; display: flex; flex-direction: column; gap: 14px; }
.sc-top { display: flex; align-items: center; justify-content: space-between; gap: 8px; }
.sc-title { font-size: 11px; font-weight: 600; color: #6b7280; text-transform: uppercase; letter-spacing: .05em; }
.sc-badge { font-size: 11px; font-weight: 700; padding: 3px 10px; border-radius: 100px; }
.sc-badge--aktif         { background: #f0fdf4; color: #16a34a; border: 1px solid #86efac; }
.sc-badge--selesai       { background: #f0fdf4; color: #0d2818; border: 1px solid #bbf7d0; }
.sc-badge--upload_laporan { background: #fefce8; color: #a16207; border: 1px solid #fde68a; }
.sc-badge--penilaian     { background: #f0fdf4; color: #0d2818; border: 1px solid #e9d5ff; }
.sc-badge--none          { background: #f9fafb; color: #6b7280; border: 1px solid #e5e7eb; }
.sc-items { display: flex; flex-direction: column; gap: 9px; }
.sc-item  { display: flex; justify-content: space-between; align-items: center; padding-bottom: 9px; border-bottom: 1px solid #f9fafb; }
.sc-item:last-child { border-bottom: none; padding-bottom: 0; }
.sc-lbl { font-size: 11.5px; color: #9ca3af; font-weight: 500; }
.sc-val { font-size: 12.5px; font-weight: 600; color: #111827; }
.sc-empty { font-size: 12px; color: #9ca3af; padding: 8px 0; }

/* ── Aksi Cepat ──────────────────────────────────────────────── */
.aksi-card { padding: 18px; display: flex; flex-direction: column; gap: 12px; }
.aksi-grid { display: flex; flex-direction: column; gap: 8px; }
.aksi-btn {
  display: flex; align-items: center; gap: 8px;
  font-size: 12.5px; font-weight: 600; padding: 10px 16px;
  border-radius: 10px; cursor: pointer; font-family: inherit;
  border: 1.5px solid transparent; text-align: left; width: 100%;
  transition: all .15s;
}
.aksi-btn:disabled { opacity: .55; cursor: not-allowed; }
.aksi-btn--green  { background: #f0fdf4; color: #16a34a; border-color: #86efac; }
.aksi-btn--green:hover:not(:disabled)  { background: #dcfce7; }
.aksi-btn--blue   { background: #f0fdf4; color: #0d2818; border-color: #bbf7d0; }
.aksi-btn--blue:hover:not(:disabled)   { background: #dcfce7; }
.aksi-btn--yellow { background: #fffbeb; color: #b45309; border-color: #fde68a; }
.aksi-btn--yellow:hover:not(:disabled) { background: #fef3c7; }
.aksi-btn--purple { background: #f0fdf4; color: #0d2818; border-color: #e9d5ff; }
.aksi-btn--purple:hover:not(:disabled) { background: #f3e8ff; }
.aksi-btn--gray   { background: #f9fafb; color: #374151; border-color: #e5e7eb; }
.aksi-btn--gray:hover:not(:disabled)   { background: #f3f4f6; }
.aksi-done { display: flex; align-items: center; gap: 8px; font-size: 12.5px; font-weight: 600; color: #16a34a; padding: 10px 0; }
.aksi-flash { font-size: 12px; font-weight: 500; padding: 8px 12px; border-radius: 8px; margin-top: 4px; }
.aksi-flash--ok  { background: #f0fdf4; color: #16a34a; border: 1px solid #86efac; }
.aksi-flash--err { background: #fff1f2; color: #be123c; border: 1px solid #fecdd3; }

/* ── Rekap chips ─────────────────────────────────────────────── */
.rekap-row { display: flex; align-items: center; gap: 12px; padding: 14px 18px; flex-wrap: wrap; }
.rekap-chip { display: flex; flex-direction: column; align-items: center; padding: 10px 20px; border-radius: 12px; min-width: 64px; }
.rc-num { font-size: 26px; font-weight: 800; line-height: 1; }
.rc-lbl { font-size: 10.5px; font-weight: 500; margin-top: 3px; color: inherit; opacity: .75; }
.rekap-chip--green  { background: #f0fdf4; color: #16a34a; border: 1px solid #bbf7d0; }
.rekap-chip--yellow { background: #fffbeb; color: #b45309; border: 1px solid #fde68a; }
.rekap-chip--blue   { background: #f0fdf4; color: #0d2818; border: 1px solid #bbf7d0; }
.rekap-chip--red    { background: #fff1f2; color: #be123c; border: 1px solid #fecdd3; }
.rekap-progress { flex: 1; min-width: 160px; }
.rp-header { display: flex; justify-content: space-between; align-items: baseline; margin-bottom: 6px; }
.rp-label { font-size: 11.5px; color: #6b7280; font-weight: 500; }
.rp-pct   { font-size: 20px; font-weight: 800; }
.rp-track { height: 8px; background: #f3f4f6; border-radius: 100px; overflow: hidden; }
.rp-fill  { height: 100%; border-radius: 100px; transition: width .6s; }
.rp-sub   { font-size: 11px; color: #9ca3af; margin-top: 5px; }

/* ── Mini Kalender ───────────────────────────────────────────── */
.mini-cal { padding: 0 18px 16px; }
.mc-days-header {
  display: grid; grid-template-columns: repeat(5, 1fr);
  text-align: center; margin-bottom: 6px;
}
.mc-days-header span { font-size: 10px; font-weight: 700; color: #9ca3af; text-transform: uppercase; padding: 4px 0; }
.mc-grid {
  display: grid; grid-template-columns: repeat(5, 1fr); gap: 4px;
}
.mc-cell {
  aspect-ratio: 1; border-radius: 8px; display: flex; align-items: center; justify-content: center;
  font-size: 11.5px; font-weight: 600; cursor: default;
}
.mc-cell--empty   { visibility: hidden; }
.mc-cell--hadir   { background: #f0fdf4; color: #16a34a; border: 1.5px solid #86efac; }
.mc-cell--izin    { background: #fffbeb; color: #b45309; border: 1.5px solid #fde68a; }
.mc-cell--sakit   { background: #f0fdf4; color: #0d2818; border: 1.5px solid #bbf7d0; }
.mc-cell--alpha   { background: #fff1f2; color: #be123c; border: 1.5px solid #fecdd3; }
.mc-cell--future  { background: #f9fafb; color: #d1d5db; border: 1px solid #f3f4f6; }
.mc-cell--kosong  { background: #f3f4f6; color: #9ca3af; border: 1px dashed #e5e7eb; }
.mc-cell--today   { background: #ecfdf5; color: #059669; border: 2px solid #10b981; font-weight: 800; box-shadow: 0 0 0 2px #d1fae5; }
.mc-legend { display: flex; gap: 10px; flex-wrap: wrap; margin-top: 10px; }
.mc-leg { font-size: 10px; font-weight: 600; padding: 2px 8px; border-radius: 100px; }
.mc-leg--green  { background: #f0fdf4; color: #16a34a; border: 1px solid #bbf7d0; }
.mc-leg--yellow { background: #fffbeb; color: #b45309; border: 1px solid #fde68a; }
.mc-leg--blue   { background: #f0fdf4; color: #0d2818; border: 1px solid #bbf7d0; }
.mc-leg--red    { background: #fff1f2; color: #be123c; border: 1px solid #fecdd3; }
.mc-leg--future { background: #f9fafb; color: #9ca3af; border: 1px solid #e5e7eb; }

/* ── Bottom row ──────────────────────────────────────────────── */
.bottom-row { display: grid; grid-template-columns: 1fr 1fr; gap: 14px; }

/* ── Dokumen ─────────────────────────────────────────────────── */
.dok-count { font-size: 12px; font-weight: 700; color: #6b7280; background: #f3f4f6; padding: 2px 9px; border-radius: 100px; }
.dok-list { padding: 8px 18px 12px; display: flex; flex-direction: column; gap: 2px; }
.dok-item { display: flex; align-items: center; gap: 10px; padding: 7px 0; border-bottom: 1px solid #f9fafb; }
.dok-item:last-child { border-bottom: none; }
.dok-icon { width: 22px; height: 22px; border-radius: 50%; display: flex; align-items: center; justify-content: center; flex-shrink: 0; }
.dok-icon--ok    { background: #f0fdf4; color: #16a34a; border: 1.5px solid #86efac; }
.dok-icon--empty { background: #f9fafb; color: #d1d5db; border: 1.5px solid #e5e7eb; }
.dok-label { flex: 1; font-size: 12.5px; color: #374151; font-weight: 500; }
.dok-label--empty { color: #9ca3af; }
.dok-ok { font-size: 11px; font-weight: 700; color: #16a34a; }

/* ── Notifikasi ──────────────────────────────────────────────── */
.btn-text { font-size: 11.5px; font-weight: 600; color: #48AF4A; background: none; border: none; cursor: pointer; font-family: inherit; padding: 0; }
.notif-list { padding: 8px 18px 12px; display: flex; flex-direction: column; gap: 2px; }
.notif-item { display: flex; align-items: flex-start; gap: 10px; padding: 9px 0; border-bottom: 1px solid #f9fafb; }
.notif-item:last-child { border-bottom: none; }
.notif-item--unread { background: #f0fdf4; margin: 0 -18px; padding: 9px 18px; border-radius: 8px; border-bottom: 1px solid #e9f5e9; }
.notif-dot { width: 7px; height: 7px; border-radius: 50%; background: #48AF4A; flex-shrink: 0; margin-top: 5px; }
.notif-body { flex: 1; min-width: 0; }
.notif-title { font-size: 12.5px; font-weight: 600; color: #111827; }
.notif-msg   { font-size: 11.5px; color: #6b7280; margin-top: 2px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.notif-time  { font-size: 10.5px; color: #9ca3af; margin-top: 3px; }

/* ── Spinner + Empty ─────────────────────────────────────────── */
.spinner { border: 2px solid #e5e7eb; border-top-color: #48AF4A; border-radius: 50%; animation: spin .7s linear infinite; }
.spinner--sm { width: 16px; height: 16px; }
@keyframes spin { to { transform: rotate(360deg); } }
.empty-state { display: flex; align-items: center; justify-content: center; padding: 24px; }

/* ── Responsive ──────────────────────────────────────────────── */
@media (max-width: 860px) {
  .top-row    { grid-template-columns: 1fr; }
  .bottom-row { grid-template-columns: 1fr; }
}
@media (max-width: 600px) {
  .wb-num { font-size: 26px; }
  .rekap-chip { padding: 8px 12px; min-width: 52px; }
  .rc-num { font-size: 20px; }
}
</style>
