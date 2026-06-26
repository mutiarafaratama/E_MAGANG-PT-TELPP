<template>
  <div class="faq-page">
    <!-- NAVBAR -->
    <nav class="navbar" :class="{ 'navbar--scrolled': navScrolled }">
      <div class="container navbar__inner">
        <a href="/" class="navbar__brand">
          <img src="/logotel.png" alt="PT TELPP" class="navbar__logo" />
          <span>e-Magang <strong>PT TELPP</strong></span>
        </a>
        <div class="navbar__links">
          <a href="/#alur">Alur</a>
          <a href="/#jadwal">Jadwal</a>
          <a href="/#syarat">Syarat</a>
          <a href="/faq" class="active">FAQ</a>
        </div>
        <div class="navbar__actions">
          <template v-if="currentUser">
            <router-link :to="`/dashboard/${currentUser.role}`" class="btn-ghost navbar__profile navbar__action--desktop">
              <div class="profile-avatar">{{ currentUser.nama_lengkap?.[0]?.toUpperCase() ?? 'U' }}</div>
              <span>{{ firstName }}</span>
            </router-link>
            <router-link :to="`/dashboard/${currentUser.role}`" class="btn-primary navbar__action--desktop">Dashboard →</router-link>
          </template>
          <template v-else>
            <router-link to="/login" class="btn-ghost navbar__action--desktop">Masuk</router-link>
            <router-link to="/daftar" class="btn-primary navbar__action--desktop">Daftar Magang</router-link>
          </template>
          <!-- Hamburger — mobile only -->
          <button class="navbar__hamburger" :class="{ 'navbar__hamburger--open': mobileMenuOpen }" @click.stop="toggleMobileMenu" aria-label="Menu">
            <span></span><span></span><span></span>
          </button>
        </div>
      </div>

      <!-- Mobile dropdown menu -->
      <Transition name="mobile-menu">
        <div v-if="mobileMenuOpen" class="navbar__mobile-menu">
          <template v-if="currentUser">
            <div class="navbar__mobile-user">
              <div class="navbar__mobile-avatar">{{ currentUser.nama_lengkap?.[0]?.toUpperCase() ?? 'U' }}</div>
              <div class="navbar__mobile-user-info">
                <div class="navbar__mobile-user-name">{{ firstName }}</div>
                <div class="navbar__mobile-user-role">{{ currentUser.role === 'peserta' ? 'Peserta Magang' : currentUser.role === 'hrd' ? 'Staff HRD' : 'Administrator' }}</div>
              </div>
            </div>
            <div class="navbar__mobile-divider"></div>
          </template>
          <a href="/#alur"   @click="closeMobileMenu">Alur Pendaftaran</a>
          <a href="/#jadwal" @click="closeMobileMenu">Jadwal</a>
          <a href="/#syarat" @click="closeMobileMenu">Persyaratan</a>
          <a href="/faq"     @click="closeMobileMenu" class="active">FAQ</a>
          <div class="navbar__mobile-divider"></div>
          <template v-if="currentUser">
            <router-link :to="`/dashboard/${currentUser.role}`" class="navbar__mobile-cta" @click="closeMobileMenu">Buka Dashboard →</router-link>
          </template>
          <template v-else>
            <router-link to="/login" class="navbar__mobile-cta" @click="closeMobileMenu">Masuk ke Portal</router-link>
          </template>
        </div>
      </Transition>
    </nav>

    <!-- HERO -->
    <section class="faq-hero">
      <div class="container">
        <div class="breadcrumb">
          <router-link to="/">Beranda</router-link>
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none"><path d="M9 18l6-6-6-6" stroke="#9ca3af" stroke-width="2" stroke-linecap="round"/></svg>
          <span>FAQ</span>
        </div>
        <h1>Pertanyaan yang Sering<br><span class="hero-accent">Diajukan</span></h1>
        <p class="hero-sub">Temukan jawaban dari pertanyaan umum seputar program magang PT TanjungEnim Lestari Pulp and Paper.</p>

        <!-- Search -->
        <div class="search-wrap">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" class="search-icon">
            <circle cx="11" cy="11" r="8" stroke="#9ca3af" stroke-width="2"/>
            <path d="M21 21l-4.35-4.35" stroke="#9ca3af" stroke-width="2" stroke-linecap="round"/>
          </svg>
          <input
            v-model="query"
            type="text"
            class="search-input"
            placeholder="Cari pertanyaan..."
            @input="openFaq = -1"
          />
          <button v-if="query" class="search-clear" @click="query = ''; openFaq = -1">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none"><path d="M18 6L6 18M6 6l12 12" stroke="#6b7280" stroke-width="2" stroke-linecap="round"/></svg>
          </button>
        </div>
      </div>
    </section>

    <!-- FAQ CONTENT -->
    <section class="faq-body">
      <div class="container container--wide">

        <!-- Loading -->
        <div v-if="loading" class="state-box">
          <div class="spinner"></div>
          <span>Memuat FAQ...</span>
        </div>

        <!-- Empty search -->
        <div v-else-if="filtered.length === 0" class="state-box">
          <svg width="48" height="48" viewBox="0 0 24 24" fill="none">
            <circle cx="11" cy="11" r="8" stroke="#d1d5db" stroke-width="1.5"/>
            <path d="M21 21l-4.35-4.35" stroke="#d1d5db" stroke-width="1.5" stroke-linecap="round"/>
          </svg>
          <div class="state-title">Tidak ada hasil untuk "{{ query }}"</div>
          <div class="state-sub">Coba kata kunci lain atau hubungi tim HRD melalui Chat Helpdesk.</div>
        </div>

        <!-- FAQ List -->
        <div v-else>
          <p class="result-info" v-if="query">
            {{ filtered.length }} pertanyaan ditemukan untuk "<strong>{{ query }}</strong>"
          </p>
          <div class="faq-grid">
            <div
              v-for="(f, i) in filtered"
              :key="(f as any).id ?? i"
              class="faq-item"
              :class="{ 'faq-item--open': openFaq === i }"
              @click="toggleFaq(i)"
            >
              <div class="faq-item__q">
                <span class="faq-num">{{ String(i + 1).padStart(2, '0') }}</span>
                <span class="faq-question">{{ (f as any).pertanyaan }}</span>
                <svg width="18" height="18" viewBox="0 0 24 24" fill="none" class="faq-chevron" :class="{ 'faq-chevron--open': openFaq === i }">
                  <path d="M19 9l-7 7-7-7" stroke="#48AF4A" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                </svg>
              </div>
              <div class="faq-item__a" :class="{ 'faq-item__a--open': openFaq === i }">
                <div class="faq-item__a-inner">{{ (f as any).jawaban }}</div>
              </div>
            </div>
          </div>
        </div>

        <!-- Still not found? -->
        <div class="helpdesk-card">
          <div class="helpdesk-card__icon">
            <svg width="22" height="22" viewBox="0 0 24 24" fill="none">
              <path d="M21 15a2 2 0 01-2 2H7l-4 4V5a2 2 0 012-2h14a2 2 0 012 2z" stroke="#48AF4A" stroke-width="2" stroke-linejoin="round"/>
            </svg>
          </div>
          <div class="helpdesk-card__body">
            <div class="helpdesk-card__title">Masih belum menemukan jawaban?</div>
            <div class="helpdesk-card__sub">Hubungi tim HRD langsung melalui Chat Helpdesk. Kami siap membantu.</div>
          </div>
          <template v-if="currentUser">
            <router-link :to="`/dashboard/${currentUser.role}`" class="helpdesk-card__btn">Buka Helpdesk →</router-link>
          </template>
          <template v-else>
            <router-link to="/login" class="helpdesk-card__btn">Masuk Sekarang →</router-link>
          </template>
        </div>
      </div>
    </section>

    <!-- FOOTER -->
    <footer class="footer">
      <div class="container footer__inner">
        <a href="/" class="footer__brand">
          <img src="/logotel.png" alt="PT TELPP" class="footer__logo" />
          <span>e-Magang <strong>PT TELPP</strong></span>
        </a>
        <p class="footer__copy">© 2026 e-Magang PT TELPP · PT TanjungEnim Lestari Pulp and Paper</p>
        <div class="footer__links">
          <router-link to="/#alur">Alur</router-link>
          <router-link to="/#jadwal">Jadwal</router-link>
          <router-link to="/#syarat">Syarat</router-link>
          <router-link to="/faq">FAQ</router-link>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';

const query = ref('');
const openFaq = ref(-1);
const loading = ref(true);
const navScrolled = ref(false);
const mobileMenuOpen = ref(false);
function toggleMobileMenu() { mobileMenuOpen.value = !mobileMenuOpen.value; }
function closeMobileMenu() { mobileMenuOpen.value = false; }

const DEFAULT_FAQS = [
  { pertanyaan: 'Kapan periode pendaftaran magang dibuka?',     jawaban: 'Pendaftaran magang dibuka setiap periode (biasanya awal semester). Pantau halaman Jadwal di portal untuk informasi terbaru.' },
  { pertanyaan: 'Apakah ada biaya untuk mendaftar magang?',    jawaban: 'Tidak ada biaya pendaftaran. Program magang di PT TELPP sepenuhnya gratis untuk peserta.' },
  { pertanyaan: 'Berapa lama durasi magang?',                  jawaban: 'Durasi magang minimal 1 bulan dan maksimal 6 bulan, tergantung kebutuhan divisi dan kesepakatan.' },
  { pertanyaan: 'Apakah peserta magang mendapat uang saku?',   jawaban: 'Ketentuan benefit akan disampaikan secara resmi saat proses penerimaan. Silakan tanyakan kepada tim HRD.' },
  { pertanyaan: 'Berapa lama proses verifikasi berkas?',       jawaban: 'Proses verifikasi berkas memakan waktu 3–5 hari kerja setelah semua dokumen lengkap diunggah.' },
  { pertanyaan: 'Divisi apa saja yang tersedia?',              jawaban: 'Tersedia lebih dari 30 divisi: Produksi, IT, Keuangan, HRD, Lingkungan, Teknik, dan banyak lagi. Penempatan disesuaikan latar belakang pendidikan.' },
  { pertanyaan: 'Bagaimana cara memantau status pengajuan?',   jawaban: 'Login ke dashboard portal e-Magang. Status pengajuan ditampilkan secara real-time beserta riwayat perubahan.' },
  { pertanyaan: 'Apakah bisa mendaftar ulang jika ditolak?',   jawaban: 'Bisa. Jika pengajuan sebelumnya ditolak atau sudah ditutup, Anda dapat mengajukan kembali di periode berikutnya.' },
  { pertanyaan: 'Dokumen apa saja yang perlu disiapkan?',      jawaban: 'Umumnya: surat pengantar dari institusi, CV, transkrip nilai, dan KTP/KTM. Detail lengkap ada di halaman Syarat.' },
  { pertanyaan: 'Apakah tersedia sertifikat setelah magang?',  jawaban: 'Ya. Peserta yang menyelesaikan magang dan mendapat penilaian akan mendapat sertifikat resmi yang dapat diunduh melalui portal.' },
];

const faqs = ref<any[]>([]);

const filtered = computed(() => {
  const q = query.value.trim().toLowerCase();
  if (!q) return faqs.value;
  return faqs.value.filter((f: any) =>
    (f.pertanyaan ?? '').toLowerCase().includes(q) ||
    (f.jawaban ?? '').toLowerCase().includes(q)
  );
});

const currentUser = computed<{ nama_lengkap: string; role: string } | null>(() => {
  try {
    const s = localStorage.getItem('user');
    return s ? JSON.parse(s) : null;
  } catch { return null; }
});
const firstName = computed(() => currentUser.value?.nama_lengkap?.split(' ')[0] ?? '');

function toggleFaq(i: number) {
  openFaq.value = openFaq.value === i ? -1 : i;
}

onMounted(async () => {
  window.addEventListener('scroll', () => {
    navScrolled.value = window.scrollY > 30;
  }, { passive: true });

  try {
    const res = await fetch('/api/landing/faq');
    const data = await res.json();
    if (Array.isArray(data.data) && data.data.length > 0) {
      faqs.value = data.data.filter((f: any) => f.pertanyaan?.trim().length > 1);
    } else {
      faqs.value = DEFAULT_FAQS;
    }
  } catch {
    faqs.value = DEFAULT_FAQS;
  } finally {
    loading.value = false;
  }
});
</script>

<style scoped>
/* ── Layout ── */
.faq-page { min-height: 100vh; display: flex; flex-direction: column; background: #f9fafb; font-family: 'Inter', system-ui, sans-serif; }
.container { max-width: 1200px; margin: 0 auto; padding: 0 64px; }
.container--narrow { max-width: 760px; }
.container--wide   { max-width: 1000px; }

/* ── Navbar ── */
.navbar {
  position: sticky; top: 0; z-index: 50;
  background: rgba(249,250,251,0.25);
  backdrop-filter: blur(0px);
  border-bottom: 1px solid transparent;
  transition: background .35s ease, backdrop-filter .35s ease, border-color .35s ease, box-shadow .35s ease;
}
.navbar--scrolled {
  background: rgba(249,250,251,0.92);
  backdrop-filter: blur(18px);
  border-bottom-color: rgba(229,231,235,0.7);
  box-shadow: 0 1px 24px rgba(0,0,0,0.06);
}
.navbar__inner { height: 64px; display: flex; align-items: center; justify-content: space-between; }
.navbar__brand { display: flex; align-items: center; gap: 8px; font-size: 16px; font-weight: 700; color: #0b1c30; text-decoration: none; }
.navbar__logo { height: 28px; width: auto; }
.navbar__brand strong { color: #48AF4A; }
.navbar__links { display: flex; gap: 28px; }
.navbar__links a { font-size: 14px; font-weight: 500; color: #6b7280; text-decoration: none; transition: color 0.15s; }
.navbar__links a:hover, .navbar__links a.active { color: #48AF4A; }
.navbar__actions { display: flex; gap: 10px; align-items: center; }
.navbar__profile { display: flex; align-items: center; gap: 6px; }

/* Hamburger — hidden on desktop */
.navbar__hamburger {
  display: none; flex-direction: column; justify-content: center; gap: 6px;
  width: 36px; height: 36px; background: none; border: none; cursor: pointer; padding: 4px; border-radius: 8px;
}
.navbar__hamburger span {
  display: block; height: 2px; background: #0b1c30; border-radius: 2px;
  transition: transform 0.22s ease, opacity 0.15s ease;
}
.navbar__hamburger--open span:nth-child(1) { transform: translateY(8px) rotate(45deg); }
.navbar__hamburger--open span:nth-child(2) { opacity: 0; transform: scaleX(0); }
.navbar__hamburger--open span:nth-child(3) { transform: translateY(-8px) rotate(-45deg); }

/* Mobile dropdown */
.navbar__mobile-menu {
  background: #fff; border-top: 1px solid #f1f5f9;
  box-shadow: 0 8px 32px rgba(0,0,0,0.10);
  display: flex; flex-direction: column; padding: 8px 0 16px;
}
.navbar__mobile-menu a {
  padding: 13px 24px; font-size: 14px; font-weight: 500; color: #374151;
  text-decoration: none; transition: background 0.15s, color 0.15s;
}
.navbar__mobile-menu a.active { color: #48AF4A; font-weight: 600; }
.navbar__mobile-menu a:hover { background: #f0fdf4; color: #48AF4A; }
.navbar__mobile-divider { height: 1px; background: #f1f5f9; margin: 8px 20px; }
.navbar__mobile-user {
  display: flex; align-items: center; gap: 12px; padding: 14px 24px 10px;
}
.navbar__mobile-avatar {
  width: 38px; height: 38px; border-radius: 50%; background: #48AF4A;
  color: #fff; font-size: 14px; font-weight: 700;
  display: flex; align-items: center; justify-content: center; flex-shrink: 0;
}
.navbar__mobile-user-info { display: flex; flex-direction: column; gap: 2px; min-width: 0; }
.navbar__mobile-user-name { font-size: 14px; font-weight: 700; color: #0b1c30; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.navbar__mobile-user-role { font-size: 11px; color: #6b7280; font-weight: 500; }
.navbar__mobile-cta {
  margin: 6px 16px 0; padding: 12px 20px !important; background: #48AF4A !important;
  color: #fff !important; border-radius: 12px; font-weight: 600 !important;
  font-size: 14px !important; text-align: center; text-decoration: none;
}
.navbar__mobile-cta:hover { background: #3d9e3f !important; }

/* Transition */
.mobile-menu-enter-active { transition: opacity 0.2s ease, transform 0.22s ease; }
.mobile-menu-leave-active { transition: opacity 0.15s ease, transform 0.18s ease; }
.mobile-menu-enter-from   { opacity: 0; transform: translateY(-8px); }
.mobile-menu-leave-to     { opacity: 0; transform: translateY(-4px); }
.profile-avatar { width: 28px; height: 28px; border-radius: 50%; background: #48AF4A; color: #fff; font-size: 12px; font-weight: 700; display: flex; align-items: center; justify-content: center; flex-shrink: 0; }
.btn-primary {
  background: #48AF4A; color: #fff; border: none; border-radius: 100px;
  padding: 10px 22px; font-size: 14px; font-weight: 600; cursor: pointer;
  font-family: "Poppins", sans-serif; transition: opacity 0.15s;
  text-decoration: none; display: inline-flex; align-items: center;
}
.btn-primary:hover { opacity: 0.9; }
.btn-ghost {
  background: none; border: none; color: #48AF4A; font-size: 14px;
  font-weight: 600; cursor: pointer; font-family: "Poppins", sans-serif; padding: 8px 12px;
  text-decoration: none; display: inline-flex; align-items: center;
}

/* ── Hero ── */
.faq-hero {
  background: linear-gradient(160deg, #0b1c30 0%, #1a3a1f 60%, #0d2b10 100%);
  padding: 64px 0 72px;
  color: #fff;
}
.breadcrumb { display: flex; align-items: center; gap: 6px; margin-bottom: 20px; font-size: 13px; color: #9ca3af; }
.breadcrumb a { color: #9ca3af; text-decoration: none; transition: color .15s; }
.breadcrumb a:hover { color: #86efac; }
.faq-hero h1 { font-size: clamp(28px, 5vw, 44px); font-weight: 800; line-height: 1.2; margin: 0 0 14px; }
.hero-accent { color: #48AF4A; }
.hero-sub { font-size: 15px; color: #d1fae5; max-width: 540px; line-height: 1.7; margin: 0 0 36px; }

/* ── Search ── */
.search-wrap {
  position: relative; max-width: 560px;
  background: #fff; border-radius: 12px;
  display: flex; align-items: center; gap: 10px;
  padding: 0 14px; box-shadow: 0 4px 24px rgba(0,0,0,0.18);
}
.search-icon { flex-shrink: 0; }
.search-input {
  flex: 1; border: none; outline: none; padding: 14px 0;
  font-size: 14.5px; color: #111827; background: transparent;
  font-family: inherit;
}
.search-input::placeholder { color: #9ca3af; }
.search-clear {
  background: none; border: none; cursor: pointer; padding: 4px;
  display: flex; align-items: center; flex-shrink: 0;
}

/* ── Body ── */
.faq-body { flex: 1; padding: 52px 0 72px; }
.result-info { font-size: 13px; color: #6b7280; margin-bottom: 16px; }
.result-info strong { color: #111827; }

/* ── FAQ Grid (2 kolom) ── */
.faq-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
  align-items: start;
}
.faq-item {
  background: #fff; border: 1.5px solid #e5e7eb; border-radius: 12px;
  cursor: pointer; transition: border-color .2s, box-shadow .2s; overflow: hidden;
}
.faq-item:hover { border-color: #a7d9a8; background: #fafffe; }
.faq-item--open { border-color: #48AF4A; box-shadow: 0 4px 20px rgba(72,175,74,0.11); }
.faq-item__q { display: flex; align-items: center; gap: 14px; padding: 16px 20px; }
.faq-num {
  width: 28px; height: 28px; border-radius: 7px; flex-shrink: 0;
  background: #f0fdf4; color: #16a34a; font-size: 11px; font-weight: 700;
  display: flex; align-items: center; justify-content: center; transition: all .2s;
}
.faq-item--open .faq-num { background: #48AF4A; color: #fff; }
.faq-question { flex: 1; font-size: 14px; font-weight: 600; color: #0b1c30; line-height: 1.5; }
.faq-chevron { transition: transform .3s; flex-shrink: 0; }
.faq-chevron--open { transform: rotate(180deg); }
.faq-item__a {
  max-height: 0; overflow: hidden; opacity: 0;
  transition: max-height .35s ease, opacity .25s ease;
}
.faq-item__a--open { max-height: 400px; opacity: 1; }
.faq-item__a-inner { padding: 0 20px 18px 62px; font-size: 13.5px; color: #4b5563; line-height: 1.75; }

/* ── State ── */
.state-box { display: flex; flex-direction: column; align-items: center; gap: 12px; padding: 64px 24px; text-align: center; color: #6b7280; }
.state-title { font-size: 15px; font-weight: 700; color: #374151; }
.state-sub   { font-size: 13px; color: #9ca3af; max-width: 360px; line-height: 1.6; }
.spinner { width: 26px; height: 26px; border: 3px solid #e5e7eb; border-top-color: #48AF4A; border-radius: 50%; animation: spin .7s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }

/* ── Helpdesk card ── */
.helpdesk-card {
  display: flex; align-items: center; gap: 16px; flex-wrap: wrap;
  margin-top: 40px; padding: 20px 24px; border-radius: 14px;
  background: linear-gradient(135deg, #f0fdf4, #dcfce7);
  border: 1.5px solid #86efac;
}
.helpdesk-card__icon {
  width: 44px; height: 44px; border-radius: 11px; background: #fff;
  display: flex; align-items: center; justify-content: center; flex-shrink: 0;
  box-shadow: 0 2px 8px rgba(72,175,74,0.12);
}
.helpdesk-card__body { flex: 1; min-width: 0; }
.helpdesk-card__title { font-size: 14px; font-weight: 700; color: #14532d; margin-bottom: 2px; }
.helpdesk-card__sub   { font-size: 13px; color: #166534; }
.helpdesk-card__btn {
  padding: 9px 22px; border-radius: 8px; background: #48AF4A;
  color: #fff; font-size: 13.5px; font-weight: 600; text-decoration: none;
  white-space: nowrap; flex-shrink: 0; transition: background .15s;
}
.helpdesk-card__btn:hover { background: #3a9e3c; }

/* ── Footer ── */
.footer { background: #0b1c30; padding: 28px 0; margin-top: auto; }
.footer__inner {
  display: flex; align-items: center; justify-content: space-between;
  gap: 16px; flex-wrap: wrap;
}
.footer__brand { display: flex; align-items: center; gap: 8px; color: #fff; text-decoration: none; font-size: 14px; font-weight: 600; }
.footer__logo  { height: 28px; width: 28px; object-fit: contain; }
.footer__copy  { font-size: 12px; color: #6b7280; margin: 0; }
.footer__links { display: flex; gap: 20px; }
.footer__links a { font-size: 13px; color: #9ca3af; text-decoration: none; transition: color .15s; }
.footer__links a:hover { color: #86efac; }

@media (max-width: 768px) {
  .faq-grid { grid-template-columns: 1fr; }
}
@media (max-width: 640px) {
  .container { padding: 0 20px; }
  .navbar__links { display: none; }
  .navbar__action--desktop { display: none !important; }
  .navbar__hamburger { display: flex; }
  .faq-hero { padding: 48px 0 56px; }
  .footer__inner { flex-direction: column; text-align: center; }
  .helpdesk-card { flex-direction: column; text-align: center; }
}
</style>
