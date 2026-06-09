<template>
  <div class="landing">
    <!-- NAVBAR -->
    <nav class="navbar">
      <div class="container navbar__inner">
        <a href="/" class="navbar__brand">
          <svg width="22" height="22" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z" fill="#48AF4A"/>
          </svg>
          <span>e-Magang <strong>PT TELPP</strong></span>
        </a>

        <div class="navbar__links">
          <a href="#alur">Alur</a>
          <a href="#jadwal">Jadwal</a>
          <a href="#syarat">Syarat</a>
          <a href="#faq">FAQ</a>
        </div>

        <div class="navbar__actions">
          <button class="btn-ghost">Masuk</button>
          <button class="btn-primary">Daftar Sekarang</button>
        </div>
      </div>
    </nav>

    <!-- HERO -->
    <section class="hero">
      <div class="hero__bg"></div>
      <div class="container hero__content">
        <div class="hero__left">
          <div class="badge">
            <span class="badge__dot"></span>
            Pendaftaran Periode III 2025 Dibuka
          </div>
          <h1 class="hero__title">
            Mulai Karir<br />di <span class="text-green">PT TELPP</span>
          </h1>
          <p class="hero__sub">
            Platform resmi pendaftaran magang PT TanjungEnim Lestari Pulp and Paper.
            Proses transparan dari seleksi hingga sertifikasi dalam satu sistem terintegrasi.
          </p>
          <div class="hero__stats">
            <div v-for="s in stats" :key="s.num" class="stat">
              <div class="stat__num">{{ s.num }}</div>
              <div class="stat__label">{{ s.label }}</div>
            </div>
          </div>
          <div class="hero__cta">
            <button class="btn-primary btn-lg">Daftar Sekarang</button>
            <button class="btn-outline-white btn-lg">Panduan Sistem</button>
          </div>
        </div>

        <div class="hero__right">
          <div class="steps-card">
            <div class="steps-card__title">5 Langkah Mudah</div>
            <div v-for="(s, i) in steps" :key="s.n" class="step-item" :class="{ 'step-item--dim': i >= 2 }">
              <div class="step-item__num" :class="i === 0 ? 'num--primary' : i === 1 ? 'num--accent' : 'num--muted'">
                {{ s.n }}
              </div>
              <div>
                <div class="step-item__title">{{ s.title }}</div>
                <div class="step-item__desc">{{ s.desc }}</div>
              </div>
            </div>
            <div class="steps-card__note">Kuota terbatas — daftar segera</div>
          </div>
        </div>
      </div>
    </section>

    <!-- ALUR PENDAFTARAN -->
    <section id="alur" class="section section--white">
      <div class="container">
        <div class="section-header">
          <div class="section-label">Cara Daftar</div>
          <h2>Alur Pendaftaran</h2>
          <p>Ikuti prosedur digital kami untuk monitoring status pengajuan secara real-time.</p>
        </div>
        <div class="alur-grid">
          <div v-for="(s, i) in steps" :key="s.n" class="alur-step">
            <div class="alur-step__circle" :class="i < 2 ? 'circle--active' : ''">
              <span>{{ s.n }}</span>
            </div>
            <div class="alur-step__label">Langkah {{ s.n }}</div>
            <div class="alur-step__title">{{ s.title }}</div>
            <div class="alur-step__desc">{{ s.desc }}</div>
          </div>
        </div>
      </div>
    </section>

    <!-- JADWAL -->
    <section id="jadwal" class="section section--surface">
      <div class="container">
        <div class="section-header">
          <div class="section-label">Tahun 2025</div>
          <h2>Jadwal Penerimaan</h2>
          <p>Kuota terbatas setiap periodenya. Daftar sebelum penuh.</p>
        </div>
        <div class="jadwal-grid">
          <div v-for="p in periods" :key="p.label" class="jadwal-card" :class="p.active ? 'jadwal-card--active' : ''">
            <div class="jadwal-card__tag" v-if="p.active">Buka</div>
            <div class="jadwal-card__label">{{ p.label }}</div>
            <div class="jadwal-card__range">{{ p.range }}</div>
            <div class="jadwal-card__status" :style="{ color: p.color }">{{ p.status }}</div>
            <div v-if="p.active" class="jadwal-card__quota">
              <div class="quota-bar"><div class="quota-fill" style="width: 42%"></div></div>
              <span>58% kuota tersisa</span>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- PERSYARATAN -->
    <section id="syarat" class="section section--white">
      <div class="container">
        <div class="section-header">
          <div class="section-label">Dokumen Wajib</div>
          <h2>Persyaratan Dokumen</h2>
          <p>Format PDF/JPG, maksimal 100MB per file. Pastikan dokumen hasil scan.</p>
        </div>
        <div class="docs-grid">
          <div v-for="d in dokumen" :key="d.title" class="docs-card">
            <div class="docs-card__icon">
              <component :is="d.icon" />
            </div>
            <h3>{{ d.title }}</h3>
            <ul>
              <li v-for="item in d.items" :key="item">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none"><path d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" stroke="#48AF4A" stroke-width="2" stroke-linecap="round"/></svg>
                {{ item }}
              </li>
            </ul>
          </div>
        </div>
        <div class="docs-tip">
          <div class="docs-tip__icon">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none"><path d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" stroke="#48AF4A" stroke-width="2" stroke-linecap="round"/></svg>
          </div>
          <div>
            <strong>Tips Verifikasi Cepat</strong>
            <p>Unggah dokumen hasil scan (bukan foto smartphone) agar proses verifikasi lebih lancar.</p>
          </div>
          <button class="btn-outline">Panduan Unggah</button>
        </div>
      </div>
    </section>

    <!-- K3 & BUDAYA KERJA -->
    <section class="section section--dark">
      <div class="container">
        <div class="k3-grid">
          <div class="k3-left">
            <div class="section-label section-label--light">Standar TELPP</div>
            <h2 class="text-white">Budaya Kerja & K3</h2>
            <p class="text-white-muted">
              Keselamatan adalah prioritas utama. Seluruh peserta wajib mengikuti standar
              operasional dan K3L di lingkungan PT TELPP.
            </p>
            <div class="k3-rules">
              <div v-for="r in k3Rules" :key="r" class="k3-rule">
                <svg width="18" height="18" viewBox="0 0 24 24" fill="none"><path d="M5 13l4 4L19 7" stroke="#48AF4A" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"/></svg>
                <span>{{ r }}</span>
              </div>
            </div>
          </div>
          <div class="k3-items">
            <div v-for="item in k3Items" :key="item.label" class="k3-item">
              <div class="k3-item__icon">
                <component :is="item.icon" />
              </div>
              <span>{{ item.label }}</span>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- FAQ -->
    <section id="faq" class="section section--surface">
      <div class="container container--narrow">
        <div class="section-header">
          <div class="section-label">Pertanyaan Umum</div>
          <h2>FAQ</h2>
          <p>Belum menemukan jawaban? Hubungi HRD via Chat Helpdesk.</p>
        </div>
        <div class="faq-list">
          <div
            v-for="(f, i) in faqs"
            :key="i"
            class="faq-item"
            :class="{ 'faq-item--open': openFaq === i }"
            @click="openFaq = openFaq === i ? -1 : i"
          >
            <div class="faq-item__q">
              <span>{{ f.q }}</span>
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" class="faq-chevron" :style="openFaq === i ? 'transform:rotate(180deg)' : ''">
                <path d="M19 9l-7 7-7-7" stroke="#48AF4A" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              </svg>
            </div>
            <div v-if="openFaq === i" class="faq-item__a">{{ f.a }}</div>
          </div>
        </div>
        <div class="faq-cta">
          <button class="btn-outline">Buka Chat Helpdesk</button>
        </div>
      </div>
    </section>

    <!-- CTA BANNER -->
    <section class="cta-banner">
      <div class="container cta-banner__inner">
        <h2>Siap Memulai Perjalananmu?</h2>
        <p>Bergabunglah bersama ratusan mahasiswa dan siswa SMK yang telah menjalani pengalaman magang berharga di PT TELPP.</p>
        <div class="cta-banner__btns">
          <button class="btn-white">Daftar Gratis Sekarang</button>
          <button class="btn-outline-white">Lihat Panduan</button>
        </div>
      </div>
    </section>

    <!-- FOOTER -->
    <footer class="footer">
      <div class="container footer__grid">
        <div>
          <div class="footer__brand">e-Magang PT TELPP</div>
          <p class="footer__desc">Portal resmi manajemen magang terpadu untuk talenta muda Indonesia.</p>
          <div class="footer__address">PT TanjungEnim Lestari Pulp and Paper<br />Muara Enim, Sumatera Selatan</div>
        </div>
        <div v-for="col in footerCols" :key="col.title">
          <div class="footer__col-title">{{ col.title }}</div>
          <ul>
            <li v-for="l in col.links" :key="l">{{ l }}</li>
          </ul>
        </div>
      </div>
      <div class="footer__bottom">
        <span>© 2025 e-Magang PT TELPP. All rights reserved.</span>
        <span>Dibuat untuk Indonesia</span>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { ref, defineComponent, h } from "vue";

const openFaq = ref(-1);

const stats = [
  { num: "500+", label: "Alumni Magang" },
  { num: "30+", label: "Divisi Terbuka" },
  { num: "95%", label: "Kepuasan Peserta" },
  { num: "1–6 Bln", label: "Durasi" },
];

const steps = [
  { n: "01", title: "Buat Akun", desc: "Daftar dengan email institusi atau pribadi." },
  { n: "02", title: "Isi Formulir", desc: "Lengkapi data diri dan akademik." },
  { n: "03", title: "Upload Berkas", desc: "Unggah dokumen sesuai kategori." },
  { n: "04", title: "Verifikasi HRD", desc: "Tim HRD memeriksa kelayakan berkas." },
  { n: "05", title: "Mulai Magang", desc: "Terima surat konfirmasi dan mulai." },
];

const periods = [
  { label: "Periode I", range: "Jan – Mar 2025", status: "Tutup", color: "#9ca3af", active: false },
  { label: "Periode II", range: "Apr – Jun 2025", status: "Segera Buka", color: "#d97706", active: false },
  { label: "Periode III", range: "Jul – Sep 2025", status: "Buka", color: "#48AF4A", active: true },
  { label: "Periode IV", range: "Okt – Des 2025", status: "Akan Datang", color: "#60a5fa", active: false },
];

const dokumen = [
  {
    title: "Mahasiswa D3/S1/S2",
    icon: defineComponent({ render: () => h("svg", { width: 24, height: 24, viewBox: "0 0 24 24", fill: "none" }, [h("path", { d: "M12 14l9-5-9-5-9 5 9 5z", stroke: "#48AF4A", "stroke-width": "2", "stroke-linecap": "round", "stroke-linejoin": "round" }), h("path", { d: "M12 14l6.16-3.422a12.083 12.083 0 01.665 6.479A11.952 11.952 0 0012 20.055a11.952 11.952 0 00-6.824-2.998 12.078 12.078 0 01.665-6.479L12 14z", stroke: "#48AF4A", "stroke-width": "2", "stroke-linecap": "round", "stroke-linejoin": "round" })]) }),
    items: ["Proposal Magang", "Surat Pengantar Kampus", "KTM Aktif", "KTP", "Pasfoto 3×4", "BPJS/KIS"],
  },
  {
    title: "Siswa SMK (Prakerin)",
    icon: defineComponent({ render: () => h("svg", { width: 24, height: 24, viewBox: "0 0 24 24", fill: "none" }, [h("path", { d: "M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z", stroke: "#48AF4A", "stroke-width": "2" }), h("circle", { cx: "12", cy: "12", r: "3", stroke: "#48AF4A", "stroke-width": "2" })]) }),
    items: ["Surat Pengantar Sekolah", "Kartu Pelajar", "KTP/KK Orang Tua", "Pasfoto 3×4"],
  },
  {
    title: "Penelitian (S2/S3)",
    icon: defineComponent({ render: () => h("svg", { width: 24, height: 24, viewBox: "0 0 24 24", fill: "none" }, [h("path", { d: "M9 3H5a2 2 0 00-2 2v4m6-6h10a2 2 0 012 2v4M9 3v18m0 0h10a2 2 0 002-2V9M9 21H5a2 2 0 01-2-2V9m0 0h18", stroke: "#48AF4A", "stroke-width": "2", "stroke-linecap": "round", "stroke-linejoin": "round" })]) }),
    items: ["Proposal Penelitian", "Surat Institusi", "CV / Riwayat Hidup", "KTP", "BPJS/KIS"],
  },
];

const k3Rules = [
  "Helm dan sepatu safety wajib di area produksi",
  "Seragam hitam-putih atau almamater sopan",
  "Hadir tepat waktu sesuai jadwal yang ditetapkan",
];

function makeIcon(path: string) {
  return defineComponent({ render: () => h("svg", { width: 20, height: 20, viewBox: "0 0 24 24", fill: "none" }, [h("path", { d: path, stroke: "#48AF4A", "stroke-width": "2", "stroke-linecap": "round", "stroke-linejoin": "round" })]) });
}

const k3Items = [
  { label: "Helm Safety", icon: makeIcon("M12 2a7 7 0 00-7 7v3H5v2h14v-2h0V9a7 7 0 00-7-7zM8 19a4 4 0 008 0H8z") },
  { label: "Seragam Sopan", icon: makeIcon("M20 21v-2a4 4 0 00-4-4H8a4 4 0 00-4 4v2M12 11a4 4 0 100-8 4 4 0 000 8z") },
  { label: "Disiplin Waktu", icon: makeIcon("M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z") },
  { label: "Absensi Digital", icon: makeIcon("M12 18h.01M8 21h8a2 2 0 002-2V5a2 2 0 00-2-2H8a2 2 0 00-2 2v14a2 2 0 002 2z") },
  { label: "Laporan Berkala", icon: makeIcon("M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z") },
  { label: "Kerjasama Tim", icon: makeIcon("M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z") },
];

const faqs = [
  { q: "Berapa lama proses verifikasi berkas?", a: "Proses verifikasi berkas memakan waktu 3–5 hari kerja setelah semua dokumen lengkap diunggah ke sistem." },
  { q: "Apakah ada kompensasi atau uang saku?", a: "Ketentuan benefit akan disampaikan secara resmi saat proses penerimaan. Silakan tanyakan langsung kepada tim HRD." },
  { q: "Divisi apa saja yang tersedia untuk magang?", a: "Tersedia lebih dari 30 divisi: Produksi, IT, Keuangan, HRD, Lingkungan, Teknik, dan banyak lagi. Penempatan disesuaikan dengan latar belakang pendidikan." },
  { q: "Bagaimana cara memantau status pengajuan?", a: "Login ke dashboard portal e-Magang. Status pengajuan ditampilkan secara real-time beserta riwayat setiap perubahan status." },
  { q: "Apakah bisa mendaftar ulang jika ditolak?", a: "Bisa. Jika pengajuan sebelumnya ditolak atau sudah ditutup, Anda dapat mengajukan kembali di periode berikutnya." },
];

const footerCols = [
  { title: "Menu", links: ["Alur Pendaftaran", "Jadwal Penerimaan", "Persyaratan", "FAQ"] },
  { title: "Bantuan", links: ["Panduan Sistem", "Chat Helpdesk", "Unduh Panduan PDF", "Kebijakan Privasi"] },
  { title: "Kontak HRD", links: ["hrd@telpp.co.id", "+62 734 123-456", "Senin–Jumat 08.00–16.00 WIB"] },
];
</script>

<style scoped>
/* ── LAYOUT ── */
.container { max-width: 1200px; margin: 0 auto; padding: 0 64px; }
.container--narrow { max-width: 760px; }

/* ── NAVBAR ── */
.navbar {
  position: sticky; top: 0; z-index: 50;
  background: rgba(249, 250, 251, 0.95);
  backdrop-filter: blur(12px);
  border-bottom: 1px solid #e5e7eb;
}
.navbar__inner {
  height: 64px; display: flex; align-items: center; justify-content: space-between;
}
.navbar__brand {
  display: flex; align-items: center; gap: 8px;
  font-size: 16px; font-weight: 700; color: #0b1c30; text-decoration: none;
}
.navbar__brand strong { color: #48AF4A; }
.navbar__links { display: flex; gap: 28px; }
.navbar__links a {
  font-size: 14px; font-weight: 500; color: #64748b;
  text-decoration: none; transition: color 0.15s;
}
.navbar__links a:hover { color: #48AF4A; }
.navbar__actions { display: flex; gap: 10px; align-items: center; }

/* ── BUTTONS ── */
.btn-primary {
  background: #48AF4A; color: #fff; border: none; border-radius: 100px;
  padding: 10px 22px; font-size: 14px; font-weight: 600; cursor: pointer;
  font-family: "Poppins", sans-serif; transition: opacity 0.15s;
}
.btn-primary:hover { opacity: 0.9; }
.btn-primary.btn-lg { padding: 13px 32px; font-size: 15px; }
.btn-ghost {
  background: none; border: none; color: #48AF4A; font-size: 14px;
  font-weight: 600; cursor: pointer; font-family: "Poppins", sans-serif; padding: 8px 12px;
}
.btn-outline {
  background: none; border: 2px solid #48AF4A; color: #48AF4A;
  border-radius: 100px; padding: 10px 24px; font-size: 14px; font-weight: 600;
  cursor: pointer; font-family: "Poppins", sans-serif; transition: all 0.15s;
}
.btn-outline:hover { background: #e8f5e9; }
.btn-outline-white {
  background: rgba(255,255,255,0.12); border: 1.5px solid rgba(255,255,255,0.3);
  color: #fff; border-radius: 100px; padding: 13px 28px; font-size: 15px;
  font-weight: 600; cursor: pointer; font-family: "Poppins", sans-serif;
}
.btn-white {
  background: #fff; color: #48AF4A; border: none; border-radius: 100px;
  padding: 13px 32px; font-size: 15px; font-weight: 700; cursor: pointer;
  font-family: "Poppins", sans-serif;
}

/* ── HERO ── */
.hero {
  position: relative; min-height: 680px;
  background: linear-gradient(135deg, #0b1c30 0%, #1a3a1f 55%, #0d2b10 100%);
  display: flex; align-items: center; overflow: hidden;
}
.hero__bg {
  position: absolute; top: -100px; right: -100px;
  width: 500px; height: 500px; border-radius: 50%;
  background: rgba(72, 175, 74, 0.12); pointer-events: none;
}
.hero__content {
  display: grid; grid-template-columns: 1fr 1fr;
  gap: 64px; align-items: center; padding-top: 80px; padding-bottom: 80px;
  position: relative; z-index: 1;
}
.hero__left { display: flex; flex-direction: column; gap: 24px; }
.badge {
  display: inline-flex; align-items: center; gap: 8px;
  background: rgba(72,175,74,0.15); border: 1px solid rgba(72,175,74,0.3);
  border-radius: 100px; padding: 6px 14px; font-size: 12px; font-weight: 600;
  color: #86efac; width: fit-content;
}
.badge__dot {
  width: 7px; height: 7px; border-radius: 50%; background: #86efac;
  box-shadow: 0 0 0 3px rgba(134,239,172,0.25);
}
.hero__title { font-size: 50px; font-weight: 800; color: #fff; line-height: 1.15; }
.text-green { color: #86efac; }
.hero__sub { font-size: 16px; color: rgba(255,255,255,0.72); line-height: 1.7; max-width: 460px; }
.hero__stats { display: flex; gap: 32px; }
.stat__num { font-size: 24px; font-weight: 800; color: #86efac; }
.stat__label { font-size: 11px; color: rgba(255,255,255,0.5); font-weight: 500; margin-top: 2px; }
.hero__cta { display: flex; gap: 14px; }

/* ── STEPS CARD ── */
.steps-card {
  background: #fff; border-radius: 20px; padding: 32px 36px;
  box-shadow: 0 24px 60px rgba(0,0,0,0.25); max-width: 360px;
}
.steps-card__title {
  font-size: 16px; font-weight: 700; color: #0b1c30;
  padding-bottom: 16px; border-bottom: 1px solid #f0f0f0; margin-bottom: 20px;
}
.step-item {
  display: flex; align-items: center; gap: 14px; margin-bottom: 16px;
}
.step-item--dim { opacity: 0.35; }
.step-item__num {
  width: 36px; height: 36px; border-radius: 50%; display: flex; align-items: center;
  justify-content: center; font-weight: 800; font-size: 13px; flex-shrink: 0;
}
.num--primary { background: #48AF4A; color: #fff; }
.num--accent { background: #ffc857; color: #fff; }
.num--muted { background: #f0f0f0; color: #9ca3af; }
.step-item__title { font-weight: 600; font-size: 14px; color: #0b1c30; }
.step-item__desc { font-size: 11px; color: #94a3b8; line-height: 1.5; }
.steps-card__note {
  margin-top: 16px; padding: 10px 14px; background: #e8f5e9; border-radius: 8px;
  font-size: 12px; color: #48AF4A; font-weight: 600; text-align: center;
}

/* ── SECTION COMMON ── */
.section { padding: 80px 0; }
.section--white { background: #fff; }
.section--surface { background: #f0f4ff; }
.section--dark { background: #0b1c30; }
.section-header { text-align: center; margin-bottom: 52px; }
.section-header h2 { font-size: 34px; font-weight: 700; color: #0b1c30; margin-bottom: 10px; }
.section-header p { font-size: 16px; color: #64748b; }
.section-label {
  font-size: 11px; font-weight: 700; color: #48AF4A; letter-spacing: 2px;
  text-transform: uppercase; margin-bottom: 10px;
}
.section-label--light { color: #86efac; }

/* ── ALUR ── */
.alur-grid {
  display: grid; grid-template-columns: repeat(5, 1fr); gap: 20px;
  position: relative;
}
.alur-grid::before {
  content: ""; position: absolute; top: 38px; left: 12%; right: 12%;
  height: 2px; background: linear-gradient(to right, #48AF4A, #ffc857); z-index: 0;
}
.alur-step { text-align: center; position: relative; z-index: 1; }
.alur-step__circle {
  width: 76px; height: 76px; border-radius: 50%; margin: 0 auto 14px;
  background: #f0f0f0; border: 2px solid #e5e7eb;
  display: flex; align-items: center; justify-content: center;
  font-size: 17px; font-weight: 800; color: #9ca3af;
}
.alur-step__circle.circle--active {
  background: #48AF4A; border-color: #48AF4A; color: #fff;
  box-shadow: 0 6px 20px rgba(72,175,74,0.25);
}
.alur-step__label { font-size: 10px; font-weight: 700; color: #9ca3af; letter-spacing: 1px; text-transform: uppercase; margin-bottom: 4px; }
.alur-step__title { font-weight: 700; font-size: 14px; color: #0b1c30; margin-bottom: 4px; }
.alur-step__desc { font-size: 12px; color: #94a3b8; line-height: 1.5; }

/* ── JADWAL ── */
.jadwal-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 18px; }
.jadwal-card {
  background: #fff; border-radius: 14px; padding: 22px;
  border: 1.5px solid #e5e7eb; position: relative; overflow: hidden;
}
.jadwal-card--active {
  border-color: #48AF4A;
  box-shadow: 0 6px 24px rgba(72,175,74,0.12);
}
.jadwal-card__tag {
  position: absolute; top: 0; right: 0;
  background: #48AF4A; color: #fff; font-size: 10px; font-weight: 700;
  padding: 4px 12px; border-bottom-left-radius: 8px; letter-spacing: 0.5px;
}
.jadwal-card__label { font-size: 11px; font-weight: 700; color: #9ca3af; text-transform: uppercase; letter-spacing: 1px; margin-bottom: 6px; }
.jadwal-card__range { font-size: 17px; font-weight: 700; color: #0b1c30; margin-bottom: 10px; }
.jadwal-card__status { font-size: 12px; font-weight: 600; text-transform: capitalize; }
.jadwal-card__quota { margin-top: 14px; }
.quota-bar { height: 5px; background: #e8f5e9; border-radius: 3px; overflow: hidden; margin-bottom: 5px; }
.quota-fill { height: 100%; background: #48AF4A; border-radius: 3px; }
.jadwal-card__quota span { font-size: 11px; color: #48AF4A; font-weight: 600; }

/* ── DOKUMEN ── */
.docs-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 22px; margin-bottom: 24px; }
.docs-card {
  background: #f9fafb; border-radius: 16px; padding: 26px;
  border: 1px solid #e5e7eb;
}
.docs-card__icon {
  width: 44px; height: 44px; background: #e8f5e9; border-radius: 10px;
  display: flex; align-items: center; justify-content: center; margin-bottom: 14px;
}
.docs-card h3 { font-size: 16px; font-weight: 700; color: #0b1c30; margin-bottom: 16px; }
.docs-card ul { list-style: none; display: flex; flex-direction: column; gap: 8px; }
.docs-card li { display: flex; align-items: center; gap: 8px; font-size: 13px; color: #64748b; }
.docs-tip {
  display: flex; align-items: center; gap: 20px;
  background: linear-gradient(135deg, #e8f5e9, rgba(255,200,87,0.08));
  border: 1px solid rgba(72,175,74,0.2); border-radius: 14px; padding: 20px 24px;
}
.docs-tip__icon {
  width: 44px; height: 44px; background: #e8f5e9; border-radius: 10px;
  display: flex; align-items: center; justify-content: center; flex-shrink: 0;
}
.docs-tip strong { display: block; font-size: 14px; color: #48AF4A; margin-bottom: 4px; }
.docs-tip p { font-size: 13px; color: #64748b; }
.docs-tip .btn-outline { margin-left: auto; white-space: nowrap; flex-shrink: 0; }

/* ── K3 ── */
.k3-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 64px; align-items: center; }
.k3-left { display: flex; flex-direction: column; gap: 16px; }
.k3-left h2 { font-size: 34px; font-weight: 700; }
.text-white { color: #fff; }
.text-white-muted { color: rgba(255,255,255,0.6); font-size: 15px; line-height: 1.7; }
.k3-rules { display: flex; flex-direction: column; gap: 10px; margin-top: 8px; }
.k3-rule {
  display: flex; align-items: center; gap: 12px;
  background: rgba(255,255,255,0.07); border-radius: 10px; padding: 12px 16px;
}
.k3-rule span { color: rgba(255,255,255,0.85); font-size: 14px; }
.k3-items { display: grid; grid-template-columns: repeat(3, 1fr); gap: 14px; }
.k3-item {
  background: rgba(255,255,255,0.07); border: 1px solid rgba(255,255,255,0.1);
  border-radius: 14px; padding: 18px 12px; text-align: center;
  display: flex; flex-direction: column; align-items: center; gap: 8px;
}
.k3-item__icon {
  width: 40px; height: 40px; background: rgba(72,175,74,0.15); border-radius: 8px;
  display: flex; align-items: center; justify-content: center;
}
.k3-item span { font-size: 12px; font-weight: 600; color: rgba(255,255,255,0.8); }

/* ── FAQ ── */
.faq-list { display: flex; flex-direction: column; gap: 10px; margin-bottom: 28px; }
.faq-item {
  border: 1px solid #e5e7eb; border-radius: 12px; background: #fff;
  cursor: pointer; transition: border-color 0.15s;
  overflow: hidden;
}
.faq-item--open { border-color: #48AF4A; background: #f0fdf4; }
.faq-item__q {
  display: flex; justify-content: space-between; align-items: center;
  padding: 16px 18px; gap: 10px;
}
.faq-item__q span { font-size: 14px; font-weight: 600; color: #48AF4A; line-height: 1.5; }
.faq-chevron { flex-shrink: 0; transition: transform 0.2s; }
.faq-item__a { padding: 0 18px 16px; font-size: 13px; color: #64748b; line-height: 1.7; }
.faq-cta { text-align: center; }

/* ── CTA BANNER ── */
.cta-banner {
  padding: 72px 64px; text-align: center;
  background: linear-gradient(135deg, #48AF4A 0%, #2d8a30 100%);
}
.cta-banner__inner { max-width: 680px; margin: 0 auto; }
.cta-banner h2 { font-size: 32px; font-weight: 800; color: #fff; margin-bottom: 12px; }
.cta-banner p { font-size: 16px; color: rgba(255,255,255,0.82); margin-bottom: 32px; line-height: 1.6; }
.cta-banner__btns { display: flex; gap: 14px; justify-content: center; }

/* ── FOOTER ── */
.footer { background: #061810; padding: 56px 0 0; }
.footer__grid {
  display: grid; grid-template-columns: 2fr 1fr 1fr 1fr;
  gap: 48px; padding-bottom: 48px; border-bottom: 1px solid rgba(255,255,255,0.08);
}
.footer__brand { font-size: 17px; font-weight: 800; color: #86efac; margin-bottom: 10px; }
.footer__desc { font-size: 13px; color: rgba(255,255,255,0.5); line-height: 1.7; margin-bottom: 16px; }
.footer__address { font-size: 12px; color: rgba(255,255,255,0.3); line-height: 1.6; }
.footer__col-title { font-size: 13px; font-weight: 700; color: #fff; margin-bottom: 14px; }
.footer ul { list-style: none; display: flex; flex-direction: column; gap: 10px; }
.footer li { font-size: 13px; color: rgba(255,255,255,0.5); cursor: pointer; }
.footer li:hover { color: #86efac; }
.footer__bottom {
  display: flex; justify-content: space-between;
  padding: 18px 64px; font-size: 11px; color: rgba(255,255,255,0.25);
}
</style>
