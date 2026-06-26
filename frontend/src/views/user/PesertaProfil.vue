<template>
  <div class="profil-wrap">

    <!-- ── Hero ───────────────────────────────────────────────── -->
    <div class="hero">
      <div class="hero__bg-deco"></div>
      <div class="hero__content">
        <div class="hero__av">{{ profilInitials }}</div>
        <div class="hero__info">
          <div class="hero__name">{{ profileData?.nama_lengkap ?? user?.nama_lengkap ?? '—' }}</div>
          <div class="hero__email">{{ profileData?.email ?? user?.email ?? '' }}</div>
          <div class="hero__chips">
            <span class="chip chip--role">Peserta Magang</span>
            <span :class="['chip', profileData?.is_active !== false ? 'chip--aktif' : 'chip--off']">
              {{ profileData?.is_active !== false ? 'Akun Aktif' : 'Tidak Aktif' }}
            </span>
            <span v-if="pengajuanSaya?.status" :class="['chip', `chip--pj-${pengajuanSaya.status}`]">
              {{ labelStatusPengajuan(pengajuanSaya.status) }}
            </span>
          </div>
        </div>
        <div class="hero__actions">
          <button class="btn-pw" @click="showPasswordModal = true">
            <svg width="13" height="13" viewBox="0 0 24 24" fill="none">
              <rect x="3" y="11" width="18" height="11" rx="2" stroke="currentColor" stroke-width="2"/>
              <path d="M7 11V7a5 5 0 0110 0v4" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
            </svg>
            Ubah Password
          </button>
          <div class="hero__join">
            <span class="hero__join-lbl">Bergabung</span>
            <span class="hero__join-val">{{ profileData?.created_at ? fmtTgl(profileData.created_at) : '—' }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- ── Skeleton ───────────────────────────────────────────── -->
    <div v-if="loading" class="skeleton-grid">
      <div class="skel-card" v-for="i in 3" :key="i">
        <div class="skel-bar skel-bar--title"></div>
        <div class="skel-bar"></div>
        <div class="skel-bar skel-bar--short"></div>
        <div class="skel-bar"></div>
        <div class="skel-bar skel-bar--mid"></div>
      </div>
    </div>

    <!-- ── Tab Layout ─────────────────────────────────────────── -->
    <div v-else class="tab-card">

      <!-- Tab Nav -->
      <div class="tab-nav">
        <button
          v-for="tab in tabs" :key="tab.key"
          :class="['tab-btn', activeTab === tab.key ? 'tab-btn--active' : '']"
          @click="activeTab = tab.key"
        >
          <span class="tab-btn__icon" v-html="tab.icon"></span>
          {{ tab.label }}
        </button>
      </div>

      <!-- Tab: Data Pribadi -->
      <div v-if="activeTab === 'pribadi'" class="tab-pane">
        <div class="field-grid">
          <div class="field">
            <div class="field__label">Nama Lengkap</div>
            <div class="field__value">{{ pengajuanSaya?.nama_lengkap ?? profileData?.nama_lengkap ?? '—' }}</div>
          </div>
          <div class="field">
            <div class="field__label">Email</div>
            <div class="field__value">{{ profileData?.email ?? '—' }}</div>
          </div>
          <div class="field">
            <div class="field__label">No. HP / WhatsApp</div>
            <div class="field__value">{{ pengajuanSaya?.no_hp ?? '—' }}</div>
          </div>
          <div class="field">
            <div class="field__label">No. Identitas (KTP/KTM)</div>
            <div class="field__value">{{ pengajuanSaya?.nomor_induk ?? '—' }}</div>
          </div>
          <div class="field">
            <div class="field__label">Tempat Lahir</div>
            <div class="field__value">{{ pengajuanSaya?.tempat_lahir ?? '—' }}</div>
          </div>
          <div class="field">
            <div class="field__label">Tanggal Lahir</div>
            <div class="field__value">{{ pengajuanSaya?.tanggal_lahir ? fmtTgl(pengajuanSaya.tanggal_lahir) : '—' }}</div>
          </div>
          <div class="field">
            <div class="field__label">Jenis Kelamin</div>
            <div class="field__value">{{ fmtJenisKelamin(pengajuanSaya?.jenis_kelamin) }}</div>
          </div>
          <div class="field field--wide">
            <div class="field__label">Alamat</div>
            <div class="field__value field__value--addr">{{ pengajuanSaya?.alamat ?? '—' }}</div>
          </div>
        </div>
      </div>

      <!-- Tab: Data Akademik -->
      <div v-else-if="activeTab === 'akademik'" class="tab-pane">
        <template v-if="pengajuanSaya">
          <div class="field-grid">
            <div class="field">
              <div class="field__label">Asal Institusi</div>
              <div class="field__value">{{ pengajuanSaya.asal_institusi ?? '—' }}</div>
            </div>
            <div class="field">
              <div class="field__label">Jurusan / Program Studi</div>
              <div class="field__value">{{ pengajuanSaya.jurusan ?? '—' }}</div>
            </div>
            <div class="field">
              <div class="field__label">Nomor Induk (NIM/NISN)</div>
              <div class="field__value">{{ pengajuanSaya.nomor_induk ?? '—' }}</div>
            </div>
            <div class="field">
              <div class="field__label">Semester / Kelas</div>
              <div class="field__value">{{ pengajuanSaya.kelas_semester ?? '—' }}</div>
            </div>
            <div class="field">
              <div class="field__label">Kategori Magang</div>
              <div class="field__value">
                <span class="badge badge--kategori">{{ labelKategori(pengajuanSaya.kategori_magang) }}</span>
              </div>
            </div>
            <div class="field">
              <div class="field__label">Tanggal Pengajuan</div>
              <div class="field__value">{{ pengajuanSaya.created_at ? fmtTgl(pengajuanSaya.created_at) : '—' }}</div>
            </div>
          </div>
        </template>
        <div v-else class="section-empty">
          <svg width="32" height="32" viewBox="0 0 24 24" fill="none">
            <path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z" stroke="#d1d5db" stroke-width="1.5"/>
            <polyline points="14 2 14 8 20 8" stroke="#d1d5db" stroke-width="1.5"/>
          </svg>
          <p>Belum ada data pengajuan</p>
        </div>
      </div>

      <!-- Tab: Pelaksanaan -->
      <div v-else-if="activeTab === 'pelaksanaan'" class="tab-pane">
        <template v-if="pelaksanaan">
          <div class="field-grid">
            <div class="field">
              <div class="field__label">Status</div>
              <div class="field__value">
                <span :class="['badge', `badge--status-${pelaksanaan.status}`]">{{ labelStatus(pelaksanaan.status) }}</span>
              </div>
            </div>
            <div class="field">
              <div class="field__label">Divisi / Unit Kerja</div>
              <div class="field__value">{{ pelaksanaan.divisi ?? '—' }}</div>
            </div>
            <div class="field">
              <div class="field__label">Nama Pembimbing</div>
              <div class="field__value">{{ pelaksanaan.pembimbing_nama ?? '—' }}</div>
            </div>
            <div class="field">
              <div class="field__label">Tanggal Mulai</div>
              <div class="field__value">{{ pelaksanaan.tanggal_mulai ? fmtTgl(pelaksanaan.tanggal_mulai) : '—' }}</div>
            </div>
            <div class="field">
              <div class="field__label">Tanggal Selesai</div>
              <div class="field__value">{{ pelaksanaan.tanggal_selesai ? fmtTgl(pelaksanaan.tanggal_selesai) : '—' }}</div>
            </div>
            <div class="field">
              <div class="field__label">Nilai Akhir</div>
              <div class="field__value">
                <span v-if="pelaksanaan.nilai != null" class="badge badge--nilai">
                  {{ pelaksanaan.nilai.toFixed(1) }}
                </span>
                <span v-else>—</span>
              </div>
            </div>
          </div>

          <!-- Progress bar durasi -->
          <div v-if="pelaksanaan.status === 'aktif' && progressMagang > 0" class="magang-progress">
            <div class="magang-progress__header">
              <span class="magang-progress__lbl">Progres Magang</span>
              <span class="magang-progress__pct">{{ progressMagang }}%</span>
            </div>
            <div class="magang-progress__track">
              <div class="magang-progress__fill" :style="{ width: progressMagang + '%' }"></div>
            </div>
            <div class="magang-progress__sub">Hari ke-{{ hariKe }} dari {{ totalHari }} hari · {{ sisaHari }} hari lagi</div>
          </div>
        </template>
        <div v-else class="section-empty">
          <svg width="32" height="32" viewBox="0 0 24 24" fill="none">
            <rect x="3" y="4" width="18" height="18" rx="2" stroke="#d1d5db" stroke-width="1.5"/>
            <line x1="3" y1="10" x2="21" y2="10" stroke="#d1d5db" stroke-width="1.5"/>
          </svg>
          <p>Data pelaksanaan belum tersedia</p>
        </div>
      </div>

    </div>

  </div>

  <!-- ── Modal Ubah Password ──────────────────────────────────── -->
  <Teleport to="body">
    <div v-if="showPasswordModal" class="modal-backdrop" @click.self="closePasswordModal">
      <div class="modal-box">
        <div class="modal-head">
          <div class="modal-title">Ubah Password</div>
          <button class="modal-close" @click="closePasswordModal">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none">
              <path d="M18 6L6 18M6 6l12 12" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
            </svg>
          </button>
        </div>
        <div class="pw-form">
          <div class="pw-field">
            <label class="pw-label">Password Lama <span class="req">*</span></label>
            <input v-model="pwLama" type="password" class="pw-input"
              placeholder="Masukkan password saat ini" autocomplete="current-password" />
          </div>
          <div class="pw-field">
            <label class="pw-label">Password Baru <span class="req">*</span></label>
            <input v-model="pwBaru" type="password" class="pw-input"
              placeholder="Min 8 karakter" autocomplete="new-password" />
            <span class="pw-hint">Min 8 karakter · harus ada huruf, angka, dan karakter spesial</span>
          </div>
          <div class="pw-field">
            <label class="pw-label">Konfirmasi Password Baru <span class="req">*</span></label>
            <input v-model="pwKonfirmasi" type="password" class="pw-input"
              placeholder="Ulangi password baru" autocomplete="new-password"
              @keyup.enter="changePassword" />
          </div>
          <div v-if="pwError"   class="pw-alert pw-alert--err">{{ pwError }}</div>
          <div v-if="pwSuccess" class="pw-alert pw-alert--ok">{{ pwSuccess }}</div>
        </div>
        <div class="modal-actions">
          <button class="btn-cancel" @click="closePasswordModal">Batal</button>
          <button class="btn-confirm" @click="changePassword" :disabled="pwLoading">
            {{ pwLoading ? 'Menyimpan…' : 'Simpan Password' }}
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import { useAuth } from "@/hooks/useAuth";
import api from "@/lib/api";

const { user } = useAuth();

const profilInitials = computed(() => {
  const name = user.value?.nama_lengkap ?? "";
  return name.split(" ").map((n: string) => n[0]).slice(0, 2).join("").toUpperCase() || "?";
});

const loading       = ref(false);
const profileData   = ref<any>(null);
const pengajuanSaya = ref<any>(null);
const pelaksanaan   = ref<any>(null);

const activeTab = ref<'pribadi'|'akademik'|'pelaksanaan'>('pribadi');
const tabs: { key: 'pribadi'|'akademik'|'pelaksanaan'; label: string; icon: string }[] = [
  {
    key: 'pribadi',
    label: 'Data Pribadi',
    icon: '<svg width="13" height="13" viewBox="0 0 24 24" fill="none"><path d="M20 21v-2a4 4 0 00-4-4H8a4 4 0 00-4 4v2" stroke="currentColor" stroke-width="2"/><circle cx="12" cy="7" r="4" stroke="currentColor" stroke-width="2"/></svg>',
  },
  {
    key: 'akademik',
    label: 'Akademik',
    icon: '<svg width="13" height="13" viewBox="0 0 24 24" fill="none"><path d="M22 10v6M2 10l10-5 10 5-10 5-10-5z" stroke="currentColor" stroke-width="2"/><path d="M6 12v5c3.33 1.67 8.67 1.67 12 0v-5" stroke="currentColor" stroke-width="2"/></svg>',
  },
  {
    key: 'pelaksanaan',
    label: 'Pelaksanaan',
    icon: '<svg width="13" height="13" viewBox="0 0 24 24" fill="none"><rect x="3" y="4" width="18" height="18" rx="2" stroke="currentColor" stroke-width="2"/><line x1="16" y1="2" x2="16" y2="6" stroke="currentColor" stroke-width="2"/><line x1="8" y1="2" x2="8" y2="6" stroke="currentColor" stroke-width="2"/><line x1="3" y1="10" x2="21" y2="10" stroke="currentColor" stroke-width="2"/></svg>',
  },
];

const showPasswordModal = ref(false);
const pwLama       = ref("");
const pwBaru       = ref("");
const pwKonfirmasi = ref("");
const pwLoading    = ref(false);
const pwError      = ref("");
const pwSuccess    = ref("");

// ── Progress Magang ───────────────────────────────────────────
const hariKe = computed(() => {
  if (!pelaksanaan.value?.tanggal_mulai) return 0;
  const diff = Math.floor((Date.now() - new Date(pelaksanaan.value.tanggal_mulai).getTime()) / 86400000) + 1;
  return Math.max(0, diff);
});
const totalHari = computed(() => {
  if (!pelaksanaan.value?.tanggal_mulai || !pelaksanaan.value?.tanggal_selesai) return 0;
  return Math.ceil((new Date(pelaksanaan.value.tanggal_selesai).getTime() - new Date(pelaksanaan.value.tanggal_mulai).getTime()) / 86400000);
});
const sisaHari = computed(() => Math.max(0, totalHari.value - hariKe.value));
const progressMagang = computed(() => {
  if (!totalHari.value) return 0;
  return Math.min(100, Math.round((hariKe.value / totalHari.value) * 100));
});

// ── Label helpers ─────────────────────────────────────────────
function fmtTgl(iso: string) {
  if (!iso) return "—";
  return new Date(iso).toLocaleDateString("id-ID", { day: "numeric", month: "long", year: "numeric" });
}
function fmtJenisKelamin(val: string | undefined) {
  if (!val) return "—";
  const m: Record<string, string> = { l: "Laki-laki", p: "Perempuan", laki_laki: "Laki-laki", perempuan: "Perempuan" };
  return m[val.toLowerCase()] ?? val;
}
function labelKategori(k: string) {
  const m: Record<string, string> = {
    smk: "SMK / Sederajat", d3_s1_s2: "D3 / S1 / S2",
    d3: "D3", s1: "S1", s2: "S2",
    penelitian: "Penelitian", lainnya: "Lainnya",
  };
  return m[k] ?? k;
}
function labelStatus(s: string) {
  const m: Record<string, string> = {
    aktif: "Sedang Berjalan", selesai: "Selesai",
    upload_laporan: "Upload Laporan", penilaian: "Penilaian Akhir",
  };
  return m[s] ?? s;
}
function labelStatusPengajuan(s: string) {
  const m: Record<string, string> = {
    diajukan: "Baru Diajukan", menunggu_verifikasi: "Menunggu Verifikasi",
    diproses: "Sedang Diproses", diterima: "Diterima", ditolak: "Ditolak",
  };
  return m[s] ?? s;
}

// ── Fetch ─────────────────────────────────────────────────────
async function fetchData() {
  loading.value = true;
  try {
    const [r1, r2, r3] = await Promise.allSettled([
      api.get("/api/auth/me"),
      api.get("/api/pengajuan/saya"),
      api.get("/api/pelaksanaan/saya"),
    ]);
    if (r1.status === "fulfilled") profileData.value   = r1.value.data?.data ?? r1.value.data;
    if (r2.status === "fulfilled") pengajuanSaya.value = r2.value.data?.data ?? null;
    if (r3.status === "fulfilled") pelaksanaan.value   = r3.value.data?.data ?? null;
  } finally {
    loading.value = false;
  }
}

// ── Password change ───────────────────────────────────────────
function validatePw(s: string) {
  if (s.length < 8) return false;
  const special = "!@#$%^&*()_+=.,><?/";
  let hasLetter = false, hasDigit = false, hasSpecial = false;
  for (const c of s) {
    if (/[a-zA-Z]/.test(c)) hasLetter = true;
    else if (/[0-9]/.test(c)) hasDigit = true;
    else if (special.includes(c)) hasSpecial = true;
  }
  return hasLetter && hasDigit && hasSpecial;
}

async function changePassword() {
  pwError.value = ""; pwSuccess.value = "";
  if (!pwLama.value || !pwBaru.value || !pwKonfirmasi.value) {
    pwError.value = "Semua field wajib diisi"; return;
  }
  if (pwBaru.value !== pwKonfirmasi.value) {
    pwError.value = "Konfirmasi password baru tidak cocok"; return;
  }
  if (!validatePw(pwBaru.value)) {
    pwError.value = "Password baru harus min 8 karakter, mengandung huruf, angka, dan karakter spesial"; return;
  }
  pwLoading.value = true;
  try {
    await api.post("/api/auth/change-password", { old_password: pwLama.value, new_password: pwBaru.value });
    pwSuccess.value = "Password berhasil diubah!";
    pwLama.value = ""; pwBaru.value = ""; pwKonfirmasi.value = "";
    setTimeout(() => { showPasswordModal.value = false; pwSuccess.value = ""; }, 2000);
  } catch (e: any) {
    pwError.value = e.response?.data?.message ?? "Gagal mengubah password";
  } finally {
    pwLoading.value = false;
  }
}

function closePasswordModal() {
  showPasswordModal.value = false;
  pwLama.value = ""; pwBaru.value = ""; pwKonfirmasi.value = "";
  pwError.value = ""; pwSuccess.value = "";
}

onMounted(fetchData);
</script>

<style scoped>
.profil-wrap { display: flex; flex-direction: column; gap: 16px; }

/* ── Hero ───────────────────────────────────────────────────── */
.hero {
  position: relative; overflow: hidden;
  background: linear-gradient(135deg, #0b2016 0%, #1a5c20 55%, #2d7a2e 100%);
  border-radius: 16px; padding: 26px 28px;
}
.hero__bg-deco {
  position: absolute; top: -40px; right: -40px;
  width: 200px; height: 200px; border-radius: 50%;
  background: rgba(134,239,172,0.07); pointer-events: none;
}
.hero__content {
  position: relative; display: flex; align-items: center; gap: 20px; flex-wrap: wrap;
}
.hero__av {
  width: 68px; height: 68px; border-radius: 18px; flex-shrink: 0;
  background: rgba(134,239,172,0.15); border: 2px solid rgba(134,239,172,0.35);
  color: #86efac; font-size: 24px; font-weight: 800;
  display: flex; align-items: center; justify-content: center;
  letter-spacing: -.5px;
}
.hero__info { flex: 1; min-width: 0; }
.hero__name  { font-size: 18px; font-weight: 700; color: #fff; letter-spacing: -.2px; }
.hero__email { font-size: 12.5px; color: rgba(255,255,255,.5); margin-top: 3px; }
.hero__chips { display: flex; gap: 7px; flex-wrap: wrap; margin-top: 10px; }

.chip { font-size: 11px; font-weight: 600; padding: 3px 11px; border-radius: 100px; }
.chip--role  { background: rgba(255,255,255,.1); color: #d1fae5; border: 1px solid rgba(255,255,255,.15); }
.chip--aktif { background: rgba(134,239,172,.18); color: #86efac; border: 1px solid rgba(134,239,172,.35); }
.chip--off   { background: rgba(239,68,68,.18); color: #fca5a5; border: 1px solid rgba(239,68,68,.35); }
.chip--pj-diterima            { background: rgba(134,239,172,.18); color: #86efac; border: 1px solid rgba(134,239,172,.35); }
.chip--pj-diajukan            { background: rgba(255,255,255,.08); color: rgba(255,255,255,.55); border: 1px solid rgba(255,255,255,.15); }
.chip--pj-menunggu_verifikasi { background: rgba(250,204,21,.15); color: #fde68a; border: 1px solid rgba(250,204,21,.3); }
.chip--pj-diproses            { background: rgba(96,165,250,.15); color: #bfdbfe; border: 1px solid rgba(96,165,250,.3); }
.chip--pj-ditolak             { background: rgba(239,68,68,.18); color: #fca5a5; border: 1px solid rgba(239,68,68,.35); }

.hero__actions { display: flex; flex-direction: column; align-items: flex-end; gap: 10px; flex-shrink: 0; }
.btn-pw {
  display: flex; align-items: center; gap: 7px;
  background: rgba(255,255,255,.1); border: 1px solid rgba(255,255,255,.2);
  color: #fff; border-radius: 10px; padding: 9px 16px;
  font-size: 12.5px; font-weight: 600; font-family: inherit; cursor: pointer;
  transition: background .15s;
}
.btn-pw:hover { background: rgba(255,255,255,.18); }
.hero__join { text-align: right; }
.hero__join-lbl { font-size: 10.5px; color: rgba(255,255,255,.4); display: block; }
.hero__join-val { font-size: 12px; color: rgba(255,255,255,.6); font-weight: 500; }

/* ── Skeleton ────────────────────────────────────────────────── */
.skeleton-grid { display: grid; grid-template-columns: 1fr; gap: 14px; }
.skel-card { background: #fff; border-radius: 14px; border: 1px solid #e9f5e9; padding: 20px 22px; display: flex; flex-direction: column; gap: 12px; }
.skel-bar { height: 13px; background: linear-gradient(90deg,#f3f4f6 25%,#e9eaeb 50%,#f3f4f6 75%); border-radius: 6px; animation: shimmer 1.4s ease-in-out infinite; background-size: 200% 100%; }
.skel-bar--title { height: 16px; width: 40%; }
.skel-bar--short { width: 55%; }
.skel-bar--mid   { width: 75%; }
@keyframes shimmer { 0% { background-position: 200% 0; } 100% { background-position: -200% 0; } }

/* ── Tab Card ────────────────────────────────────────────────── */
.tab-card {
  background: #fff; border: 1.5px solid #edf5ed; border-radius: 16px;
  overflow: hidden; box-shadow: 0 1px 4px rgba(13,40,24,.04);
}

.tab-nav {
  display: flex; border-bottom: 1.5px solid #f0faf0;
  background: #f9fafb; padding: 0 6px;
}

.tab-btn {
  display: flex; align-items: center; gap: 6px;
  padding: 13px 16px; border: none; background: transparent;
  font-size: 12.5px; font-weight: 600; color: #6b7280;
  font-family: inherit; cursor: pointer; border-bottom: 2.5px solid transparent;
  margin-bottom: -1.5px; white-space: nowrap; transition: color .15s;
}
.tab-btn:hover { color: #374151; }
.tab-btn--active { color: #16a34a; border-bottom-color: #16a34a; }
.tab-btn__icon { display: flex; align-items: center; }

.tab-pane { padding: 20px 22px 22px; }

/* ── Sections (legacy, kept for section-empty) ───────────────── */
.sections { display: flex; flex-direction: column; gap: 14px; }

.section-card {
  background: #fff; border: 1.5px solid #edf5ed; border-radius: 16px;
  overflow: hidden; box-shadow: 0 1px 4px rgba(13,40,24,.04);
}

.section-head {
  display: flex; align-items: flex-start; gap: 12px;
  padding: 18px 22px 16px; border-bottom: 1.5px solid #f3faf3;
}
.section-head__icon {
  width: 36px; height: 36px; border-radius: 10px; flex-shrink: 0;
  display: flex; align-items: center; justify-content: center; margin-top: 1px;
}
.section-head__icon--green { background: #f0fdf4; color: #16a34a; border: 1px solid #bbf7d0; }
.section-head__icon--blue  { background: #eff6ff; color: #2563eb; border: 1px solid #bfdbfe; }
.section-head__icon--amber { background: #fffbeb; color: #b45309; border: 1px solid #fde68a; }
.section-head__title { font-size: 13.5px; font-weight: 700; color: #111827; }
.section-head__sub   { font-size: 11.5px; color: #9ca3af; margin-top: 2px; }

/* ── Field Grid ──────────────────────────────────────────────── */
.field-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 0;
}
.field {
  padding: 14px 22px 13px;
  border-bottom: 1px solid #f5faf5;
  border-right: 1px solid #f5faf5;
}
.field:last-child { border-bottom: none; }
.field--wide {
  grid-column: 1 / -1;
}
.field__label {
  font-size: 10.5px; font-weight: 600; color: #9ca3af;
  text-transform: uppercase; letter-spacing: .06em; margin-bottom: 5px;
}
.field__value {
  font-size: 13.5px; font-weight: 600; color: #1f2937; line-height: 1.4;
}
.field__value--addr {
  font-size: 13px; font-weight: 500; color: #374151;
  line-height: 1.6; white-space: pre-wrap;
}

/* ── Badges ──────────────────────────────────────────────────── */
.badge {
  display: inline-flex; font-size: 11.5px; font-weight: 700;
  padding: 3px 11px; border-radius: 100px;
}
.badge--kategori  { background: #eff6ff; color: #2563eb; border: 1px solid #bfdbfe; }
.badge--nilai     { background: #f0fdf4; color: #15803d; border: 1px solid #86efac; font-size: 14px; }
.badge--status-aktif          { background: #f0fdf4; color: #15803d; border: 1px solid #86efac; }
.badge--status-selesai        { background: #eff6ff; color: #1d4ed8; border: 1px solid #bfdbfe; }
.badge--status-upload_laporan { background: #fffbeb; color: #a16207; border: 1px solid #fde68a; }
.badge--status-penilaian      { background: #fdf4ff; color: #7e22ce; border: 1px solid #e9d5ff; }

/* ── Magang progress ─────────────────────────────────────────── */
.magang-progress {
  margin: 4px 22px 18px; padding: 14px 18px;
  background: #f9fffe; border: 1px solid #d1fae5; border-radius: 12px;
}
.magang-progress__header { display: flex; justify-content: space-between; align-items: baseline; margin-bottom: 8px; }
.magang-progress__lbl { font-size: 12px; color: #6b7280; font-weight: 500; }
.magang-progress__pct { font-size: 18px; font-weight: 800; color: #16a34a; }
.magang-progress__track { height: 7px; background: #e5e7eb; border-radius: 100px; overflow: hidden; }
.magang-progress__fill  { height: 100%; background: #48AF4A; border-radius: 100px; transition: width .6s; }
.magang-progress__sub   { font-size: 11px; color: #9ca3af; margin-top: 6px; }

/* ── Empty ───────────────────────────────────────────────────── */
.section-empty {
  display: flex; flex-direction: column; align-items: center; gap: 9px;
  padding: 32px 22px; text-align: center;
}
.section-empty p { font-size: 12.5px; color: #9ca3af; margin: 0; }

/* ── Modal Password ──────────────────────────────────────────── */
.modal-backdrop {
  position: fixed; inset: 0; background: rgba(0,0,0,.42);
  display: flex; align-items: center; justify-content: center;
  z-index: 1000; padding: 16px;
}
.modal-box {
  background: #fff; border-radius: 16px; padding: 0 24px 24px;
  width: 100%; max-width: 420px; box-shadow: 0 20px 60px rgba(0,0,0,.18);
}
.modal-head {
  display: flex; align-items: center; justify-content: space-between;
  padding: 20px 0 16px;
}
.modal-title { font-size: 16px; font-weight: 700; color: #111827; }
.modal-close {
  background: #f3f4f6; border: none; border-radius: 8px;
  width: 30px; height: 30px; display: flex; align-items: center; justify-content: center;
  cursor: pointer; color: #6b7280;
}
.modal-close:hover { background: #e5e7eb; }
.pw-form { display: flex; flex-direction: column; gap: 14px; }
.pw-field { display: flex; flex-direction: column; gap: 5px; }
.pw-label { font-size: 12px; font-weight: 600; color: #374151; }
.req { color: #ef4444; }
.pw-input {
  border: 1.5px solid #e5e7eb; border-radius: 9px; padding: 9px 12px;
  font-size: 13px; font-family: inherit; outline: none; transition: border-color .15s;
}
.pw-input:focus { border-color: #48AF4A; }
.pw-hint { font-size: 11px; color: #9ca3af; }
.pw-alert { padding: 8px 12px; border-radius: 8px; font-size: 12.5px; }
.pw-alert--err { background: #fff1f2; color: #be123c; border: 1px solid #fecdd3; }
.pw-alert--ok  { background: #f0fdf4; color: #16a34a; border: 1px solid #bbf7d0; }
.modal-actions { display: flex; gap: 10px; justify-content: flex-end; margin-top: 18px; }
.btn-cancel {
  background: #f3f4f6; color: #374151; border: none; border-radius: 9px;
  padding: 9px 20px; font-size: 13px; font-weight: 600; font-family: inherit; cursor: pointer;
}
.btn-confirm {
  background: #48AF4A; color: #fff; border: none; border-radius: 9px;
  padding: 9px 20px; font-size: 13px; font-weight: 600; font-family: inherit; cursor: pointer;
}
.btn-confirm:disabled { background: #d1d5db; cursor: default; }

/* ── Responsive ──────────────────────────────────────────────── */
@media (max-width: 700px) {
  .hero__content { gap: 14px; }
  .hero__actions { flex-direction: row; align-items: center; width: 100%; justify-content: space-between; }
  .field-grid { grid-template-columns: 1fr; }
}
</style>
