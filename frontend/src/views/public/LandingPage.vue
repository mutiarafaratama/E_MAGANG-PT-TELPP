<template>
  <div class="landing">

    <!-- Chat FAB Landing (anonymous) -->
    <ChatFAB v-model="chatOpen" :unread-count="0">
      <ChatWidgetPublic />
    </ChatFAB>

    <!-- NAVBAR -->
    <nav class="navbar" :class="{ 'navbar--scrolled': navScrolled }">
      <div class="container navbar__inner">
        <a href="/" class="navbar__brand">
          <img src="/logotel.png" alt="PT TELPP" class="navbar__logo" />
          <span>e-Magang <strong>PT TELPP</strong></span>
        </a>

        <div class="navbar__links">
          <a href="#alur" @click="closeMobileMenu">Alur</a>
          <a href="#jadwal" @click="closeMobileMenu">Jadwal</a>
          <a href="#syarat" @click="closeMobileMenu">Syarat</a>
          <a href="#faq" @click="closeMobileMenu">FAQ</a>
        </div>

        <div class="navbar__actions">
          <template v-if="currentUser">
            <!-- Desktop: tampilkan avatar + nama + tombol dashboard -->
            <router-link :to="dashboardPath" class="btn-ghost navbar__profile navbar__action--desktop">
              <div class="profile-avatar">{{ currentUser.nama_lengkap?.[0]?.toUpperCase() ?? 'U' }}</div>
              <span>{{ firstName }}</span>
            </router-link>
            <router-link :to="dashboardPath" class="btn-primary navbar__action--desktop">Dashboard →</router-link>
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
          <!-- Info akun jika sudah login -->
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

          <a href="#alur"   @click="closeMobileMenu">Alur Pendaftaran</a>
          <a href="#jadwal" @click="closeMobileMenu">Jadwal</a>
          <a href="#syarat" @click="closeMobileMenu">Persyaratan</a>
          <a href="#faq"    @click="closeMobileMenu">FAQ</a>
          <div class="navbar__mobile-divider"></div>

          <template v-if="currentUser">
            <router-link :to="dashboardPath" class="navbar__mobile-cta" @click="closeMobileMenu">Buka Dashboard →</router-link>
          </template>
          <template v-else>
            <router-link to="/login" class="navbar__mobile-cta" @click="closeMobileMenu">Masuk ke Portal</router-link>
          </template>
        </div>
      </Transition>
    </nav>

    <!-- HERO -->
    <section class="hero" :style="heroStyle">
      <div class="hero__overlay" :style="overlayStyle"></div>
      <div class="hero__decor hero__decor--1"></div>
      <div class="hero__decor hero__decor--2"></div>
      <div class="container hero__container">
        <div class="hero__content">
          <h1 class="hero__title hero-anim hero-anim--1">
            {{ hero.judul }}<br />
            <span class="text-green">{{ hero.judul_aksen }}</span>
          </h1>
          <div class="hero__stats hero-anim hero-anim--2" :style="statCardStyle">
            <template v-for="(s, i) in hero.stats" :key="s.num">
              <div class="stat-item">
                <div class="stat__num">{{ s.num }}</div>
                <div class="stat__label">{{ s.label }}</div>
              </div>
              <div v-if="i < hero.stats.length - 1" class="stat-divider"></div>
            </template>
          </div>
          <div class="hero__cta hero-anim hero-anim--3">
            <template v-if="!currentUser">
              <router-link to="/daftar" class="btn-primary btn-lg">Daftar Sekarang</router-link>
              <button class="btn-outline-white btn-lg">Panduan Sistem</button>
            </template>
            <template v-else>
              <router-link :to="dashboardPath" class="btn-primary btn-lg">Buka Dashboard →</router-link>
            </template>
          </div>
        </div>
      </div>
    </section>

    <!-- ALUR PENDAFTARAN -->
    <section id="alur" class="section section--white">
      <div class="container">
        <div v-if="alurItems.length > 0">
          <!-- Header full-width — di atas grid, tidak ikut dalam kolom -->
          <div class="alur-header anim-fade-up">
            <h2>Alur Pendaftaran</h2>
            <p class="alur-header__sub">Tahapan mudah bergabung bersama e-magang</p>
          </div>

          <!-- Grid 2 kolom: kiri = tab + teks, kanan = gambar -->
          <div
            class="alur-main-grid"
            @touchstart.passive="onAlurTouchStart"
            @touchend.passive="onAlurTouchEnd"
          >
            <!-- KIRI: tab (desktop only) + konten teks animasi -->
            <div class="alur-left-col">
              <!-- Tab buttons — hanya desktop, 1 baris -->
              <div class="alur-tabs alur-tabs--desktop anim-fade-up" data-delay="80">
                <button
                  v-for="(item, i) in alurItems"
                  :key="item.id || i"
                  class="alur-tab"
                  :class="{ 'alur-tab--active': i === alurIdx }"
                  @click="alurIdx = i"
                >Tahap {{ i + 1 }}</button>
              </div>
              <!-- Konten teks animasi -->
              <Transition name="alur-slide" mode="out-in">
                <div class="alur-text-content" :key="alurIdx">
                  <div class="alur-zig-watermark">{{ String(alurIdx + 1).padStart(2, '0') }}</div>
                  <h3 class="alur-zig-title">{{ alurItems[alurIdx]?.judul }}</h3>
                  <p class="alur-zig-desc">{{ alurItems[alurIdx]?.paragraf }}</p>
                </div>
              </Transition>
            </div>

            <!-- KANAN: gambar animasi -->
            <div class="alur-right-col">
              <Transition name="alur-slide" mode="out-in">
                <div class="alur-zig-img" :key="alurIdx">
                  <img
                    v-if="alurItems[alurIdx]?.gambar_url"
                    :src="alurItems[alurIdx].gambar_url"
                    :alt="alurItems[alurIdx].judul"
                  />
                  <div v-else class="alur-zig-empty">
                    <span>{{ String(alurIdx + 1).padStart(2, '0') }}</span>
                    <p>{{ alurItems[alurIdx]?.judul }}</p>
                  </div>
                </div>
              </Transition>
            </div>
          </div>

          <!-- Dot indicators + swipe hint — mobile only -->
          <div class="alur-mobile-nav">
            <div class="alur-dots">
              <button
                v-for="(_, i) in alurItems"
                :key="i"
                class="alur-dot"
                :class="{ 'alur-dot--active': i === alurIdx }"
                @click="alurIdx = i"
                :aria-label="`Tahap ${i + 1}`"
              ></button>
            </div>
            <div class="alur-swipe-hint">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none"><path d="M15 18l-6-6 6-6" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"/></svg>
              Geser untuk navigasi
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none"><path d="M9 18l6-6-6-6" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"/></svg>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- JADWAL -->
    <section id="jadwal" class="section section--surface" :style="jadwalSectionStyle">
      <div v-if="jadwalBg.type === 'image'" class="jadwal-bg-overlay" :style="jadwalOverlayStyle"></div>
      <div class="container">
        <div class="section-header anim-fade-up">
          <h2>Jadwal Penerimaan</h2>
          <p>Kuota terbatas setiap periodenya. Daftar sebelum penuh.</p>
        </div>

        <!-- ── PERIOD-GRID (new default design) ── -->
        <div v-if="jadwal_style === 'period-grid'" class="jpg-wrap anim-fade-up" data-delay="100">
          <div class="jpg__track" :style="(periodes.length > 0 ? periodes.length : periodsFallback.length) <= 3 ? 'justify-content:center' : ''">
            <!-- Data dari DB -->
            <template v-if="periodes.length > 0">
              <template v-for="(p, i) in periodes" :key="p.id">
                <div class="jpg__period">
                  <div class="jpg__circle" :class="periodeIsActive(p) ? 'jpg__circle--active' : ''">
                    <span class="jpg__circle-lbl">PERIODE</span>
                    <span class="jpg__circle-num">{{ String(i + 1).padStart(2, '0') }}</span>
                  </div>
                  <div class="jpg__badge" :class="periodeIsActive(p) ? 'jpg__badge--active' : ''">
                    {{ p.nama }}
                  </div>
                  <div class="jpg__card" :class="periodeIsActive(p) ? 'jpg__card--active' : ''">
                    <div class="jpg__card-open">Pendaftaran</div>
                    <div class="jpg__card-range">{{ fmtDate(p.tanggal_buka) }} – {{ fmtDate(p.tanggal_tutup) }}</div>
                    <span class="jpg__pill" :class="periodeStatusClass(p)">{{ periodeStatusLabel(p) }}</span>
                    <div v-if="p.kuota" class="jpg__kuota">Kuota: {{ p.kuota }}</div>
                  </div>
                </div>
                <div v-if="i < periodes.length - 1" class="jpg__connector">
                  <div class="jpg__connector-line"></div>
                </div>
              </template>
            </template>

            <!-- Fallback: tidak ada data dari DB -->
            <template v-else>
              <template v-for="(p, i) in periodsFallback" :key="p.label">
                <div class="jpg__period">
                  <div class="jpg__circle" :class="p.active ? 'jpg__circle--active' : ''">
                    <span class="jpg__circle-lbl">PERIODE</span>
                    <span class="jpg__circle-num">{{ String(i + 1).padStart(2, '0') }}</span>
                  </div>
                  <div class="jpg__badge" :class="p.active ? 'jpg__badge--active' : ''">
                    {{ p.label }}
                  </div>
                  <div class="jpg__card" :class="p.active ? 'jpg__card--active' : ''">
                    <div class="jpg__card-open">Pendaftaran</div>
                    <div class="jpg__card-range">{{ p.range }}</div>
                    <span class="jpg__pill" :class="p.active ? 'status--green' : 'status--gray'">{{ p.status }}</span>
                  </div>
                </div>
                <div v-if="i < periodsFallback.length - 1" class="jpg__connector">
                  <div class="jpg__connector-line"></div>
                </div>
              </template>
            </template>
          </div>
        </div>

        <!-- Style: timeline-vertical (backward compat CMS) -->
        <div v-else-if="jadwal_style === 'timeline-vertical'" class="jtv">
          <div v-for="(p, i) in periodes" :key="p.id" class="jtv__item anim-fade-up" :data-delay="i * 80">
            <div class="jtv__left">
              <div class="jtv__dot" :class="periodeStatusClass(p)"></div>
              <div class="jtv__line" v-if="i < periodes.length - 1"></div>
            </div>
            <div class="jtv__card" :class="periodeIsActive(p) ? 'jtv__card--active' : ''">
              <div class="jtv__tag" v-if="periodeIsActive(p)">Sedang Buka</div>
              <div class="jtv__name">{{ p.nama }}</div>
              <div class="jtv__range">{{ fmtDate(p.tanggal_buka) }} – {{ fmtDate(p.tanggal_tutup) }}</div>
              <span class="jtv__pill" :class="periodeStatusClass(p)">{{ periodeStatusLabel(p) }}</span>
              <div v-if="p.kuota" class="jtv__kuota">Kuota: {{ p.kuota }}</div>
            </div>
          </div>
        </div>

        <!-- Style: timeline-horizontal -->
        <div v-else-if="jadwal_style === 'timeline-horizontal'" class="jth">
          <div v-for="(p, i) in periodes" :key="p.id" class="jth__step anim-fade-up" :data-delay="i * 80">
            <div class="jth__head">
              <div class="jth__line jth__line--pre" v-if="i > 0"></div>
              <div class="jth__circle" :class="periodeStatusClass(p)">{{ i + 1 }}</div>
              <div class="jth__line jth__line--post" v-if="i < periodes.length - 1"></div>
            </div>
            <div class="jth__body">
              <div class="jth__name">{{ p.nama }}</div>
              <div class="jth__range">{{ fmtMonthYear(p.tanggal_buka) }} – {{ fmtMonthYear(p.tanggal_tutup) }}</div>
              <span class="jth__pill" :class="periodeStatusClass(p)">{{ periodeStatusLabel(p) }}</span>
            </div>
          </div>
        </div>

        <!-- Style: timeline-cards -->
        <div v-else-if="jadwal_style === 'timeline-cards'" class="jadwal-cards-grid">
          <div v-for="(p, i) in periodes" :key="p.id"
            class="jcard anim-fade-up"
            :class="periodeIsActive(p) ? 'jcard--active' : ''"
            :data-delay="i * 80"
          >
            <div class="jcard__tag" v-if="periodeIsActive(p)">Sedang Buka</div>
            <div class="jcard__name">{{ p.nama }}</div>
            <div class="jcard__range">{{ fmtDate(p.tanggal_buka) }} – {{ fmtDate(p.tanggal_tutup) }}</div>
            <span class="jcard__pill" :class="periodeStatusClass(p)">{{ periodeStatusLabel(p) }}</span>
            <div v-if="p.kuota" class="jcard__kuota">Kuota: {{ p.kuota }}</div>
          </div>
        </div>

        <!-- Style: fallback compact -->
        <div v-else class="jcmp">
          <div v-for="(p, i) in periodes" :key="p.id" class="jcmp__row anim-fade-up" :data-delay="i * 80">
            <div class="jcmp__dot" :class="periodeStatusClass(p)"></div>
            <div class="jcmp__name">{{ p.nama }}</div>
            <div class="jcmp__range">{{ fmtDate(p.tanggal_buka) }} – {{ fmtDate(p.tanggal_tutup) }}</div>
            <span class="jcmp__pill" :class="periodeStatusClass(p)">{{ periodeStatusLabel(p) }}</span>
          </div>
        </div>

      </div>
    </section>

    <!-- PERSYARATAN -->
    <section id="syarat" class="section section--white">
      <div class="container">
        <div class="section-header anim-fade-up">
          <h2>Persyaratan Dokumen</h2>
          <p>Format PDF/JPG, maksimal 10MB per file. Pastikan dokumen hasil scan.</p>
        </div>
        <!-- Desktop: 3-col grid -->
        <div class="docs-grid docs-desktop">
          <div
            v-for="(d, i) in dokumen"
            :key="d.title"
            class="docs-card anim-fade-up"
            :data-delay="i * 100"
          >
            <div class="docs-card__icon"><component :is="d.icon" /></div>
            <h3>{{ d.title }}</h3>
            <ul>
              <li v-for="item in d.items" :key="item">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none"><path d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" stroke="#48AF4A" stroke-width="2" stroke-linecap="round"/></svg>
                {{ item }}
              </li>
            </ul>
          </div>
        </div>

        <!-- Mobile: carousel satu kartu -->
        <div class="docs-carousel docs-mobile"
          @touchstart.passive="onDocsTouchStart"
          @touchend.passive="onDocsTouchEnd"
        >
          <Transition name="alur-slide" mode="out-in">
            <div class="docs-card docs-card--full" :key="docsIdx">
              <div class="docs-card__icon"><component :is="dokumen[docsIdx].icon" /></div>
              <h3>{{ dokumen[docsIdx].title }}</h3>
              <ul>
                <li v-for="item in dokumen[docsIdx].items" :key="item">
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none"><path d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" stroke="#48AF4A" stroke-width="2" stroke-linecap="round"/></svg>
                  {{ item }}
                </li>
              </ul>
            </div>
          </Transition>
          <div class="alur-mobile-nav">
            <div class="alur-dots">
              <button
                v-for="(_, i) in dokumen"
                :key="i"
                class="alur-dot"
                :class="{ 'alur-dot--active': i === docsIdx }"
                @click="docsIdx = i"
              ></button>
            </div>
            <div class="alur-swipe-hint">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none"><path d="M15 18l-6-6 6-6" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"/></svg>
              Geser untuk navigasi
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none"><path d="M9 18l6-6-6-6" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"/></svg>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- K3 & BUDAYA KERJA -->
    <section class="section section--dark">
      <div class="container">
        <div class="k3-grid">
          <div class="k3-left anim-fade-up">
            <h2 class="text-white">Budaya Kerja & K3</h2>
            <p class="text-white-muted">
              Keselamatan adalah prioritas utama. Seluruh peserta wajib mengikuti standar
              operasional dan K3L di lingkungan PT TELPP.
            </p>
            <div class="k3-rules">
              <div v-for="(r, i) in k3Rules" :key="r" class="k3-rule anim-fade-up" :data-delay="i * 130">
                <span class="k3-rule__num">0{{ i + 1 }}</span>
                <span>{{ r }}</span>
              </div>
            </div>
          </div>
          <div class="k3-items">
            <div
              v-for="(item, i) in k3Items"
              :key="item.label"
              class="k3-item anim-fade-up"
              :data-delay="i * 80"
            >
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
        <div class="section-header anim-fade-up">
          <h2>Pertanyaan yang Sering Diajukan</h2>
          <p>Belum menemukan jawaban? Kunjungi halaman FAQ lengkap atau hubungi tim HRD.</p>
        </div>
        <div class="faq-list">
          <div
            v-for="(f, i) in faqsPreview"
            :key="(f as any).id ?? i"
            class="faq-item"
            :class="{ 'faq-item--open': openFaq === i, 'anim-fade-up': true, 'is-visible': faqVisible.has(i) }"
            :data-faq-idx="String(i)"
            :data-delay="i * 55"
            @click.stop="toggleFaq(i)"
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
        <div class="faq-more-wrap">
          <router-link to="/faq" class="faq-more-btn">
            Lihat Semua Pertanyaan
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none">
              <path d="M5 12h14M12 5l7 7-7 7" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
          </router-link>
        </div>
      </div>
    </section>

    <!-- CTA BANNER -->
    <section class="cta-banner anim-fade-up">
      <div class="container cta-banner__inner">
        <h2>Siap Memulai Perjalananmu?</h2>
        <p>Bergabunglah bersama ratusan mahasiswa dan siswa SMK yang telah menjalani pengalaman magang berharga di PT TELPP.</p>
        <div class="cta-banner__btns">
          <router-link v-if="!currentUser" to="/daftar" class="btn-white">Ajukan Magang Sekarang</router-link>
          <router-link v-else :to="dashboardPath" class="btn-white">Buka Dashboard →</router-link>
          <button class="btn-outline-white">Lihat Panduan</button>
        </div>
      </div>
    </section>

    <!-- FOOTER -->
    <footer class="footer">
      <div class="container footer__grid">
        <div class="footer__brand-col">
          <a href="/" class="footer__brand-link">
            <img src="/logotel.png" alt="PT TELPP" class="footer__logo" />
            <span>e-Magang <strong>PT TELPP</strong></span>
          </a>
          <p class="footer__desc">Portal resmi manajemen magang terpadu untuk talenta muda Indonesia.</p>
          <p class="footer__address">PT TanjungEnim Lestari Pulp and Paper<br />Muara Enim, Sumatera Selatan</p>
        </div>
        <div>
          <div class="footer__col-title">Navigasi</div>
          <ul>
            <li><a href="#alur">Alur Pendaftaran</a></li>
            <li><a href="#jadwal">Jadwal Penerimaan</a></li>
            <li><a href="#syarat">Persyaratan</a></li>
            <li><a href="#faq">FAQ</a></li>
          </ul>
        </div>
        <div>
          <div class="footer__col-title">Akun</div>
          <ul>
            <li><router-link to="/login">Masuk</router-link></li>
            <li><router-link to="/daftar">Daftar Magang</router-link></li>
            <li><router-link to="/dashboard">Dashboard</router-link></li>
          </ul>
        </div>
        <div>
          <div class="footer__col-title">Kontak HRD</div>
          <ul>
            <li>hrd@telpp.co.id</li>
            <li>+62 734 123-456</li>
            <li>Senin–Jumat, 08.00–16.00 WIB</li>
          </ul>
        </div>
      </div>
      <div class="footer__bottom">
        <div class="container footer__bottom-inner">
          <span>© 2026 e-Magang PT TELPP. Hak cipta dilindungi.</span>
          <span>PT TanjungEnim Lestari Pulp and Paper</span>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, defineComponent, h, onMounted, onUnmounted, reactive } from "vue";
import ChatFAB from "@/components/ChatFAB.vue";
import ChatWidgetPublic from "@/components/ChatWidgetPublic.vue";

const chatOpen = ref(false);

const openFaq = ref(-1);
const faqVisible = reactive(new Set<number>());
const faqsPreview = computed(() =>
  faqs.value.filter((f: any) => f.pertanyaan && f.pertanyaan.trim().length > 1).slice(0, 4)
);
function toggleFaq(i: number) {
  openFaq.value = openFaq.value === i ? -1 : i;
}
const navScrolled = ref(false);
const mobileMenuOpen = ref(false);
function toggleMobileMenu() { mobileMenuOpen.value = !mobileMenuOpen.value; }
function closeMobileMenu() { mobileMenuOpen.value = false; }
let outsideClickHandler: ((e: MouseEvent) => void) | null = null;

// ── Hero content ──────────────────────────────────────────────
const hero = ref({
  judul:           'Mulai Karir',
  judul_aksen:     'PT TELPP',
  bg_url:          '',
  bg_opacity:      70,
  stat_card_opacity: 20,
  stats: [
    { num: '500+',    label: 'Alumni Magang' },
    { num: '30+',     label: 'Divisi Terbuka' },
    { num: '95%',     label: 'Kepuasan Peserta' },
    { num: '1–6 Bln', label: 'Durasi' },
  ],
});

const statCardStyle = computed(() => {
  const op = (hero.value.stat_card_opacity ?? 20) / 100;
  const border = Math.min(op + 0.15, 0.45);
  return {
    background: `rgba(255,255,255,${op.toFixed(2)})`,
    backdropFilter: 'blur(10px)',
    WebkitBackdropFilter: 'blur(10px)',
    boxShadow: '0 8px 32px rgba(0,0,0,0.22), 0 2px 8px rgba(0,0,0,0.12)',
    border: `1px solid rgba(255,255,255,${border.toFixed(2)})`,
  };
});

const jadwal_style = ref('period-grid');
const jadwalBg = ref({ type: 'none', color: '#f0fdf4', url: '', opacity: 50 });

const jadwalSectionStyle = computed(() => {
  if (jadwalBg.value.type === 'color') {
    return { backgroundColor: jadwalBg.value.color };
  }
  if (jadwalBg.value.type === 'image' && jadwalBg.value.url) {
    return {
      backgroundImage: `url(${jadwalBg.value.url})`,
      backgroundSize: 'cover',
      backgroundPosition: 'center',
    };
  }
  return {};
});

const jadwalOverlayStyle = computed(() => {
  if (jadwalBg.value.type === 'image') {
    const alpha = parseFloat((Math.max(0, 1 - jadwalBg.value.opacity / 100) * 0.85).toFixed(2));
    return { background: `rgba(0,0,0,${alpha})` };
  }
  return {};
});

interface Periode {
  id: string; nama: string;
  tanggal_buka: string; tanggal_tutup: string;
  kuota: number | null; is_active: boolean;
}
const periodes = ref<Periode[]>([]);

const heroStyle = computed(() => {
  if (hero.value.bg_url) {
    return {
      backgroundImage: `url(${hero.value.bg_url})`,
      backgroundSize: 'cover',
      backgroundPosition: 'center',
    };
  }
  return {};
});

const overlayStyle = computed(() => {
  const val = hero.value.bg_opacity;
  const alpha = hero.value.bg_url
    ? Math.max(0.30, 0.92 - (val / 100) * 0.62)
    : 0;
  return {
    background: `linear-gradient(135deg, rgba(11,28,48,${alpha}) 0%, rgba(26,58,31,${(alpha - 0.08).toFixed(2)}) 55%, rgba(13,43,16,${alpha}) 100%)`,
  };
});

// ── Alur carousel ────────────────────────────────────────────
interface AlurSlide { id?: string; judul: string; paragraf: string; gambar_url: string; urutan?: number; }
const alurItems    = ref<AlurSlide[]>([]);
const alurIdx      = ref(0);
const alurTouchX   = ref(0);

function onAlurTouchStart(e: TouchEvent) {
  alurTouchX.value = e.changedTouches[0].clientX;
}
function onAlurTouchEnd(e: TouchEvent) {
  const dx = e.changedTouches[0].clientX - alurTouchX.value;
  if (Math.abs(dx) < 40) return;
  if (dx < 0 && alurIdx.value < alurItems.value.length - 1) alurIdx.value++;
  else if (dx > 0 && alurIdx.value > 0) alurIdx.value--;
}

// ── Docs carousel (mobile) ──────────────────────────────────────
const docsIdx    = ref(0);
const docsTouchX = ref(0);
function onDocsTouchStart(e: TouchEvent) {
  docsTouchX.value = e.changedTouches[0].clientX;
}
function onDocsTouchEnd(e: TouchEvent) {
  const dx = e.changedTouches[0].clientX - docsTouchX.value;
  if (Math.abs(dx) < 40) return;
  if (dx < 0 && docsIdx.value < dokumen.length - 1) docsIdx.value++;
  else if (dx > 0 && docsIdx.value > 0) docsIdx.value--;
}

// ── Intersection Observer for scroll animations (bidirectional) ────────────
let scrollObserver: IntersectionObserver | null = null;
const animDelayTimers: number[] = [];

function initScrollAnimations() {
  scrollObserver?.disconnect();
  scrollObserver = new IntersectionObserver(
    (entries) => {
      entries.forEach((entry) => {
        const el = entry.target as HTMLElement;
        const delay = el.dataset.delay ? parseInt(el.dataset.delay) : 0;
        const faqIdx = el.dataset.faqIdx;
        if (entry.isIntersecting) {
          // Scroll into view → animate in
          const tid = setTimeout(() => {
            if (faqIdx !== undefined) faqVisible.add(parseInt(faqIdx));
            else el.classList.add('is-visible');
          }, delay) as unknown as number;
          animDelayTimers.push(tid);
        } else {
          // Scroll out of view → animate back out
          if (faqIdx !== undefined) faqVisible.delete(parseInt(faqIdx));
          else el.classList.remove('is-visible');
        }
      });
    },
    { threshold: 0.12, rootMargin: '0px 0px -32px 0px' }
  );
  document.querySelectorAll('.anim-fade-up').forEach((el) => scrollObserver!.observe(el));
}

onMounted(async () => {
  window.addEventListener('scroll', () => { navScrolled.value = window.scrollY > 30; }, { passive: true });

  // Fetch hero content
  try {
    const res = await fetch('/api/landing/content');
    const data = await res.json();
    const d = data.data ?? {};
    if (d.hero_judul)           hero.value.judul             = d.hero_judul;
    if (d.hero_judul_aksen)     hero.value.judul_aksen       = d.hero_judul_aksen;
    if (d.hero_bg_url)          hero.value.bg_url            = d.hero_bg_url;
    if (d.hero_bg_opacity)      hero.value.bg_opacity        = parseInt(d.hero_bg_opacity) || 70;
    if (d.hero_stat_card_opacity) hero.value.stat_card_opacity = parseInt(d.hero_stat_card_opacity) || 20;
    if (d.jadwal_style)        jadwal_style.value        = d.jadwal_style;
    if (d.jadwal_bg_type)      jadwalBg.value.type       = d.jadwal_bg_type;
    if (d.jadwal_bg_color)     jadwalBg.value.color      = d.jadwal_bg_color;
    if (d.jadwal_bg_url)       jadwalBg.value.url        = d.jadwal_bg_url;
    if (d.jadwal_bg_opacity)   jadwalBg.value.opacity    = parseInt(d.jadwal_bg_opacity) || 50;
    for (let i = 1; i <= 4; i++) {
      if (d[`hero_stat_${i}_num`])   hero.value.stats[i - 1].num   = d[`hero_stat_${i}_num`];
      if (d[`hero_stat_${i}_label`]) hero.value.stats[i - 1].label = d[`hero_stat_${i}_label`];
    }
  } catch { /* pakai default */ }

  // Fetch semua periodes (untuk jadwal section)
  try {
    const res = await fetch('/api/landing/periodes');
    const data = await res.json();
    if (Array.isArray(data.data) && data.data.length > 0) {
      periodes.value = data.data;
    }
  } catch { /* pakai fallback kosong */ }

  // Fetch alur
  try {
    const res = await fetch('/api/landing/alur');
    const data = await res.json();
    if (Array.isArray(data.data) && data.data.length > 0) {
      alurItems.value = data.data;
    } else {
      alurItems.value = steps.map(s => ({ judul: s.title, paragraf: s.desc, gambar_url: '' }));
    }
  } catch {
    alurItems.value = steps.map(s => ({ judul: s.title, paragraf: s.desc, gambar_url: '' }));
  }

  // Fetch FAQ dari API
  try {
    const res = await fetch('/api/landing/faq');
    const data = await res.json();
    if (Array.isArray(data.data) && data.data.length > 0) faqs.value = data.data;
  } catch { /* pakai fallback default */ }

  // Init animations after DOM is ready
  setTimeout(initScrollAnimations, 50);

  // Close mobile menu on outside click
  outsideClickHandler = (e: MouseEvent) => {
    const nav = document.querySelector('.navbar');
    if (nav && !nav.contains(e.target as Node)) closeMobileMenu();
  };
  document.addEventListener('click', outsideClickHandler);
});

onUnmounted(() => {
  scrollObserver?.disconnect();
  animDelayTimers.forEach(clearTimeout);
  if (outsideClickHandler) document.removeEventListener('click', outsideClickHandler);
});

const currentUser = computed<{ nama_lengkap: string; role: string } | null>(() => {
  try {
    const s = localStorage.getItem("user");
    return s ? JSON.parse(s) : null;
  } catch {
    return null;
  }
});

const firstName = computed(() => currentUser.value?.nama_lengkap?.split(" ")[0] ?? "");

const dashboardPath = computed(() => {
  const role = currentUser.value?.role;
  if (role === 'admin') return '/admin';
  if (role === 'hrd') return '/staff';
  return '/dashboard';
});

const steps = [
  { n: "01", title: "Isi Formulir",    desc: "Lengkapi data diri dan akademik melalui formulir online." },
  { n: "02", title: "Verifikasi HRD",  desc: "Tim HRD memeriksa kelayakan berkas dalam 3–5 hari kerja." },
  { n: "03", title: "Terima Akun",     desc: "Akun login e-Magang dikirim ke email Anda setelah diterima." },
  { n: "04", title: "Upload Berkas",   desc: "Login dan unggah dokumen persyaratan yang diperlukan." },
  { n: "05", title: "Mulai Magang",    desc: "Mulai magang dan pantau progress via dashboard." },
];

const periodsFallback = [
  { label: "Periode I",   range: "Jan – Mar 2025", status: "Tutup",       color: "#9ca3af", active: false },
  { label: "Periode II",  range: "Apr – Jun 2025", status: "Segera Buka", color: "#d97706", active: false },
  { label: "Periode III", range: "Jul – Sep 2025", status: "Buka",        color: "#48AF4A", active: true  },
  { label: "Periode IV",  range: "Okt – Des 2025", status: "Akan Datang", color: "#60a5fa", active: false },
];

function fmtDate(d: string) {
  return new Date(d).toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' });
}
function fmtMonthYear(d: string) {
  return new Date(d).toLocaleDateString('id-ID', { month: 'short', year: 'numeric' });
}
function periodeIsActive(p: Periode) {
  const now = new Date(); const b = new Date(p.tanggal_buka); const t = new Date(p.tanggal_tutup);
  return p.is_active && now >= b && now <= t;
}
function periodeStatusLabel(p: Periode) {
  const now = new Date(); const b = new Date(p.tanggal_buka); const t = new Date(p.tanggal_tutup);
  if (now > t)  return 'Selesai';
  if (now < b)  return 'Akan Datang';
  if (p.is_active) return 'Buka';
  return 'Ditutup';
}
function periodeStatusClass(p: Periode) {
  const lbl = periodeStatusLabel(p);
  if (lbl === 'Buka')         return 'status--green';
  if (lbl === 'Akan Datang')  return 'status--yellow';
  return 'status--gray';
}

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
  { label: "Helm Safety",    icon: makeIcon("M12 2a7 7 0 00-7 7v3H5v2h14v-2h0V9a7 7 0 00-7-7zM8 19a4 4 0 008 0H8z") },
  { label: "Seragam Sopan",  icon: makeIcon("M20 21v-2a4 4 0 00-4-4H8a4 4 0 00-4 4v2M12 11a4 4 0 100-8 4 4 0 000 8z") },
  { label: "Disiplin Waktu", icon: makeIcon("M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z") },
  { label: "Absensi Digital",icon: makeIcon("M12 18h.01M8 21h8a2 2 0 002-2V5a2 2 0 00-2-2H8a2 2 0 00-2 2v14a2 2 0 002 2z") },
  { label: "Laporan Berkala",icon: makeIcon("M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z") },
  { label: "Kerjasama Tim",  icon: makeIcon("M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z") },
];

const faqs = ref([
  { pertanyaan: "Berapa lama proses verifikasi berkas?",    jawaban: "Proses verifikasi berkas memakan waktu 3–5 hari kerja setelah semua dokumen lengkap diunggah ke sistem." },
  { pertanyaan: "Apakah ada kompensasi atau uang saku?",    jawaban: "Ketentuan benefit akan disampaikan secara resmi saat proses penerimaan. Silakan tanyakan langsung kepada tim HRD." },
  { pertanyaan: "Divisi apa saja yang tersedia?",           jawaban: "Tersedia lebih dari 30 divisi: Produksi, IT, Keuangan, HRD, Lingkungan, Teknik, dan banyak lagi. Penempatan disesuaikan dengan latar belakang pendidikan." },
  { pertanyaan: "Bagaimana cara memantau status pengajuan?",jawaban: "Login ke dashboard portal e-Magang. Status pengajuan ditampilkan secara real-time beserta riwayat setiap perubahan status." },
  { pertanyaan: "Apakah bisa mendaftar ulang jika ditolak?",jawaban: "Bisa. Jika pengajuan sebelumnya ditolak atau sudah ditutup, Anda dapat mengajukan kembali di periode berikutnya." },
]);
</script>

<style scoped>
/* ── LAYOUT ── */
.container { max-width: 1200px; margin: 0 auto; padding: 0 64px; }
.container--narrow { max-width: 760px; }

/* ── NAVBAR ── */
.navbar {
  position: sticky; top: 0; z-index: 50;
  background: rgba(249, 250, 251, 0.25);
  backdrop-filter: blur(0px);
  border-bottom: 1px solid transparent;
  transition: background 0.35s ease, backdrop-filter 0.35s ease, border-color 0.35s ease, box-shadow 0.35s ease;
}
.navbar--scrolled {
  background: rgba(249, 250, 251, 0.95);
  backdrop-filter: blur(18px);
  border-bottom-color: rgba(229,231,235,0.7);
  box-shadow: 0 1px 24px rgba(0,0,0,0.06);
}
.navbar__inner {
  height: 64px; display: flex; align-items: center; justify-content: space-between;
}
.navbar__brand {
  display: flex; align-items: center; gap: 8px;
  font-size: 16px; font-weight: 700; color: #0b1c30; text-decoration: none;
}
.navbar__logo { height: 28px; width: auto; }
.navbar__brand strong { color: #48AF4A; }
.navbar__links { display: flex; gap: 28px; }
.navbar__links a {
  font-size: 14px; font-weight: 500; color: #6b7280;
  text-decoration: none; transition: color 0.15s;
}
.navbar__links a:hover { color: #48AF4A; }
.navbar__actions { display: flex; gap: 10px; align-items: center; }
.navbar__profile { display: flex; align-items: center; gap: 6px; }
.profile-avatar { width: 28px; height: 28px; border-radius: 50%; background: #48AF4A; color: #fff; font-size: 12px; font-weight: 700; display: flex; align-items: center; justify-content: center; flex-shrink: 0; }

/* Hamburger button — hidden on desktop */
.navbar__hamburger {
  display: none; flex-direction: column; justify-content: space-between;
  width: 28px; height: 20px; background: none; border: none; cursor: pointer; padding: 0;
  flex-shrink: 0;
}
.navbar__hamburger span {
  display: block; height: 2px; border-radius: 2px;
  background: #0b1c30;
  transition: transform 0.28s ease, opacity 0.2s ease, width 0.2s ease;
  transform-origin: center;
}
.navbar__hamburger--open span:nth-child(1) { transform: translateY(9px) rotate(45deg); }
.navbar__hamburger--open span:nth-child(2) { opacity: 0; transform: scaleX(0); }
.navbar__hamburger--open span:nth-child(3) { transform: translateY(-9px) rotate(-45deg); }

/* Mobile dropdown */
.navbar__mobile-menu {
  position: absolute; top: 64px; left: 0; right: 0;
  background: #fff;
  border-bottom: 1px solid #e5e7eb;
  box-shadow: 0 8px 32px rgba(0,0,0,0.10);
  display: flex; flex-direction: column; padding: 8px 0 16px;
  z-index: 49;
}
.navbar__mobile-menu a {
  display: block; padding: 13px 28px;
  font-size: 15px; font-weight: 500; color: #374151;
  text-decoration: none; transition: background 0.12s, color 0.12s;
}
.navbar__mobile-menu a:hover { background: #f0fdf4; color: #48AF4A; }
.navbar__mobile-divider { height: 1px; background: #f1f5f9; margin: 8px 20px; }

/* User info block di mobile menu */
.navbar__mobile-user {
  display: flex; align-items: center; gap: 12px;
  padding: 14px 24px 12px;
}
.navbar__mobile-avatar {
  width: 38px; height: 38px; border-radius: 50%; flex-shrink: 0;
  background: #48AF4A; color: #fff; font-size: 15px; font-weight: 700;
  display: flex; align-items: center; justify-content: center;
}
.navbar__mobile-user-info { display: flex; flex-direction: column; gap: 2px; min-width: 0; }
.navbar__mobile-user-name { font-size: 14px; font-weight: 700; color: #0b1c30; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.navbar__mobile-user-role { font-size: 11px; color: #6b7280; font-weight: 500; }
.navbar__mobile-cta {
  margin: 4px 20px 0; padding: 12px 20px !important;
  background: #48AF4A !important; color: #fff !important;
  border-radius: 100px; font-weight: 700 !important; text-align: center;
  font-size: 14px !important;
}
.navbar__mobile-cta:hover { background: #3d9e3f !important; color: #fff !important; }

/* Mobile menu transition */
.mobile-menu-enter-active { transition: opacity 0.2s ease, transform 0.22s ease; }
.mobile-menu-leave-active { transition: opacity 0.15s ease, transform 0.18s ease; }
.mobile-menu-enter-from   { opacity: 0; transform: translateY(-8px); }
.mobile-menu-leave-to     { opacity: 0; transform: translateY(-4px); }

/* ── BUTTONS ── */
.btn-primary {
  background: #48AF4A; color: #fff; border: none; border-radius: 100px;
  padding: 10px 22px; font-size: 14px; font-weight: 600; cursor: pointer;
  font-family: "Poppins", sans-serif; transition: opacity 0.15s;
  text-decoration: none; display: inline-flex; align-items: center;
}
.btn-primary:hover { opacity: 0.9; }
.btn-primary.btn-lg { padding: 13px 32px; font-size: 15px; }
.btn-ghost {
  background: none; border: none; color: #48AF4A; font-size: 14px;
  font-weight: 600; cursor: pointer; font-family: "Poppins", sans-serif; padding: 8px 12px;
  text-decoration: none; display: inline-flex; align-items: center;
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
  transition: background 0.15s;
}
.btn-outline-white:hover { background: rgba(255,255,255,0.2); }
.btn-white {
  background: #fff; color: #48AF4A; border: none; border-radius: 100px;
  padding: 13px 32px; font-size: 15px; font-weight: 700; cursor: pointer;
  font-family: "Poppins", sans-serif; text-decoration: none;
  display: inline-flex; align-items: center;
}

/* ── SCROLL ANIMATIONS ── */
.anim-fade-up {
  opacity: 0;
  transform: translateY(28px);
  transition: opacity 0.65s ease, transform 0.65s ease;
}
.anim-fade-up.is-visible {
  opacity: 1;
  transform: translateY(0);
}

/* ── HERO ANIMATIONS (on-load, no observer) ── */
@keyframes heroFadeUp {
  from { opacity: 0; transform: translateY(24px); }
  to   { opacity: 1; transform: translateY(0); }
}
.hero-anim {
  animation: heroFadeUp 0.7s ease both;
}
.hero-anim--1 { animation-delay: 0.05s; }
.hero-anim--2 { animation-delay: 0.18s; }
.hero-anim--3 { animation-delay: 0.30s; }
.hero-anim--4 { animation-delay: 0.42s; }
.hero-anim--5 { animation-delay: 0.54s; }

/* ── HERO ── */
.hero {
  position: relative;
  min-height: calc(100vh - 64px);
  background: linear-gradient(135deg, #0b1c30 0%, #1a3a1f 55%, #0d2b10 100%);
  display: flex; align-items: center; overflow: hidden;
  padding-top: 30px;
  padding-bottom: 80px;
}
.hero__overlay {
  position: absolute; inset: 0; pointer-events: none; z-index: 0;
  background: linear-gradient(135deg, rgba(11,28,48,0.92) 0%, rgba(26,58,31,0.84) 55%, rgba(13,43,16,0.92) 100%);
}
.hero__decor {
  position: absolute; border-radius: 50%;
  background: radial-gradient(circle, rgba(72,175,74,0.12) 0%, transparent 70%);
  pointer-events: none;
}
.hero__decor--1 { width: 600px; height: 600px; top: -150px; right: -100px; }
.hero__decor--2 { width: 400px; height: 400px; bottom: -100px; left: -80px; opacity: 0.6; }
/* hero__container adalah wrapper full-width di dalam section hero */
.hero__container {
  display: flex; align-items: center; justify-content: center;
}
.hero__content {
  position: relative; z-index: 1;
  display: flex; flex-direction: column; gap: 24px;
  max-width: 1060px; width: 100%; margin: 0 auto;
  text-align: center; align-items: center;
}
.hero__title {
  font-size: clamp(26px, 3.6vw, 46px); font-weight: 900;
  color: #fff; line-height: 1.15; letter-spacing: -0.5px; margin: 0;
  white-space: nowrap;
}
.text-green { color: #86efac; }
.hero__title .text-green { font-size: 0.72em; display: block; white-space: nowrap; }
.hero__stats {
  display: inline-flex; flex-wrap: wrap; gap: 0;
  justify-content: center; align-items: center;
  padding: 20px 32px; border-radius: 20px;
}
.stat-item  { padding: 0 28px; text-align: center; }
.stat-divider {
  width: 1px; height: 40px; flex-shrink: 0;
  background: rgba(255,255,255,0.25);
}
.stat__num   { font-size: 28px; font-weight: 800; color: #fff; line-height: 1; }
.stat__label { font-size: 11.5px; color: rgba(255,255,255,0.75); font-weight: 500; margin-top: 4px; }
.hero__cta   { display: flex; gap: 14px; flex-wrap: wrap; align-items: center; justify-content: center; }

/* ── SECTIONS ── */
.section { padding: 96px 0; }
.section--white   { background: #edf1f8; }
.section--surface { background: #f8fafc; }
/* jadwal section supports custom bg */
#jadwal { position: relative; }
.jadwal-bg-overlay {
  position: absolute; inset: 0; z-index: 0; pointer-events: none;
}
#jadwal .container { position: relative; z-index: 1; }
.section--dark    { background: #0b1c30; }
.section-header { text-align: center; margin-bottom: 56px; }
.section-header h2 { font-size: 36px; font-weight: 800; color: #0b1c30; margin: 8px 0 12px; }
.section-header p  { font-size: 15px; color: #6b7280; max-width: 480px; margin: 0 auto; line-height: 1.6; }
.section-label {
  display: inline-block;
  font-size: 11px; font-weight: 700; letter-spacing: 1.5px; text-transform: uppercase;
  color: #48AF4A; background: #e8f5e9; padding: 4px 12px; border-radius: 100px;
}
.section-label--light { background: rgba(72,175,74,0.15); color: #86efac; }
.text-white       { color: #fff !important; }
.text-white-muted { color: rgba(255,255,255,0.65) !important; }

/* ── ALUR ── */
.alur-header { text-align: left; margin-bottom: 24px; }
.alur-header h2 { font-size: 36px; font-weight: 800; color: #0b1c30; margin: 0 0 6px; }
.alur-header__sub { font-size: 15px; color: #6b7280; margin: 0; line-height: 1.6; }

/* 2-column grid layout */
.alur-main-grid {
  display: grid; grid-template-columns: 1fr 1fr;
  gap: 60px; align-items: start; touch-action: pan-y; user-select: none;
}
.alur-left-col { display: flex; flex-direction: column; }
.alur-right-col { }
.alur-text-content { position: relative; margin-top: 24px; }

/* Tabs — 1 baris, tidak wrap */
.alur-tabs   { display: flex; gap: 6px; flex-wrap: nowrap; margin: 0 0 0; overflow: hidden; }
.alur-tab {
  background: #f1f5f9; border: 1.5px solid #e5e7eb; border-radius: 100px;
  padding: 9px 18px; font-size: 13px; font-weight: 600; color: #6b7280;
  cursor: pointer; font-family: "Poppins", sans-serif; transition: all 0.22s;
  white-space: nowrap; flex-shrink: 0;
}
.alur-tab:hover { background: #e8f5e9; border-color: #86efac; color: #48AF4A; }
.alur-tab--active { background: #48AF4A; color: #fff; border-color: #48AF4A; box-shadow: 0 4px 14px rgba(72,175,74,0.25); }

/* Transition saat ganti tahap */
.alur-slide-enter-active { transition: opacity .18s ease, transform .18s ease; }
.alur-slide-leave-active { transition: opacity .12s ease; }
.alur-slide-enter-from   { opacity: 0; transform: translateY(10px); }
.alur-slide-leave-to     { opacity: 0; }

/* Dot nav + swipe hint — hidden on desktop, shown on mobile */
.alur-mobile-nav { display: none; flex-direction: column; align-items: center; gap: 10px; margin-top: 24px; }
.alur-dots       { display: flex; gap: 8px; align-items: center; }
.alur-dot {
  width: 8px; height: 8px; border-radius: 50%; border: none; padding: 0;
  background: #d1d5db; cursor: pointer; transition: all .2s;
}
.alur-dot--active { background: #48AF4A; width: 22px; border-radius: 4px; }
.alur-swipe-hint {
  display: flex; align-items: center; gap: 5px;
  font-size: 11px; color: #9ca3af; font-weight: 500;
}
.alur-zigzag {
  display: grid; grid-template-columns: 1fr 1fr; gap: 60px; align-items: center;
}
.alur-zigzag--flip .alur-zig-text { order: 2; }
.alur-zigzag--flip .alur-zig-img  { order: 1; }
.alur-zig-text { position: relative; }
.alur-zig-watermark {
  font-size: 100px; font-weight: 900; color: #f0fdf4;
  position: absolute; top: -24px; left: -16px; line-height: 1; z-index: 0; pointer-events: none;
}
.alur-zig-step  { font-size: 12px; font-weight: 700; color: #48AF4A; text-transform: uppercase; letter-spacing: 1px; position: relative; z-index: 1; }
.alur-zig-title { font-size: 28px; font-weight: 800; color: #0b1c30; margin: 8px 0 12px; position: relative; z-index: 1; }
.alur-zig-desc  { font-size: 15px; color: #6b7280; line-height: 1.7; position: relative; z-index: 1; }
.alur-zig-img img { width: 100%; border-radius: 16px; aspect-ratio: 4/3; object-fit: cover; box-shadow: 0 20px 60px rgba(0,0,0,0.22), 0 4px 16px rgba(0,0,0,0.12); }
.alur-zig-empty {
  width: 100%; aspect-ratio: 4/3; background: linear-gradient(135deg, #e8f5e9, #f0fdf4);
  border-radius: 16px; display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 8px;
}
.alur-zig-empty span { font-size: 64px; font-weight: 900; color: #bbf7d0; }
.alur-zig-empty p    { font-size: 14px; font-weight: 600; color: #48AF4A; margin: 0; }

/* ── JADWAL — shared status pills ── */
.status--green  { background: #dcfce7; color: #16a34a; }
.status--yellow { background: #fef3c7; color: #92400e; }
.status--gray   { background: #f3f4f6; color: #6b7280; }

/* ── JADWAL — period-grid (new default) ── */
@keyframes lineFlow {
  0%   { background-position: 100% center; }
  100% { background-position: -100% center; }
}
.jpg-wrap {
  background: linear-gradient(135deg, #e8f5e9 0%, #f0fdf4 100%);
  border-radius: 20px; padding: 40px 36px 36px;
  border: 1px solid #c8e6c9;
  overflow-x: auto; -webkit-overflow-scrolling: touch;
}
.jpg__track {
  display: flex; align-items: flex-start;
  min-width: max-content; gap: 0;
}
.jpg__period {
  display: flex; flex-direction: column; align-items: center;
  min-width: 140px; max-width: 180px; flex: 1;
}
.jpg__connector {
  flex: 0.5; min-width: 28px; max-width: 80px;
  display: flex; align-items: flex-start;
  padding-top: 36px;
}
.jpg__connector-line {
  width: 100%; height: 2.5px; border-radius: 2px;
  background: linear-gradient(90deg,
    #48AF4A 0%, #86efac 30%, #48AF4A 60%, #86efac 90%, #48AF4A 100%);
  background-size: 200% 100%;
  animation: lineFlow 2.5s linear infinite;
}
.jpg__circle {
  width: 72px; height: 72px; border-radius: 50%;
  border: 2.5px solid #48AF4A; background: #fff;
  display: flex; flex-direction: column; align-items: center; justify-content: center;
  box-shadow: 0 4px 16px rgba(72,175,74,0.18);
  transition: transform 0.25s ease, box-shadow 0.25s ease;
  position: relative; z-index: 1; cursor: default;
}
.jpg__circle:hover { transform: translateY(-4px) scale(1.04); box-shadow: 0 10px 28px rgba(72,175,74,0.30); }
.jpg__circle--active { border-color: #2d8a30; background: #f0fdf4; box-shadow: 0 6px 22px rgba(72,175,74,0.35); }
.jpg__circle-lbl { font-size: 7px; font-weight: 800; color: #48AF4A; letter-spacing: 2px; text-transform: uppercase; line-height: 1; }
.jpg__circle-num { font-size: 26px; font-weight: 900; color: #0b1c30; line-height: 1; margin-top: 2px; }
.jpg__circle--active .jpg__circle-num { color: #2d8a30; }
.jpg__badge {
  margin-top: 14px; padding: 6px 12px; border-radius: 5px;
  background: #0b1c30; color: #fff;
  font-size: 9px; font-weight: 800; letter-spacing: 0.8px; text-transform: uppercase;
  text-align: center; max-width: 100%; white-space: nowrap;
  overflow: hidden; text-overflow: ellipsis;
}
.jpg__badge--active { background: #48AF4A; }
.jpg__card {
  margin-top: 8px; width: 100%;
  background: #fff; border-radius: 12px;
  padding: 14px 10px; text-align: center;
  border: 1.5px solid #e5e7eb;
  box-shadow: 0 2px 10px rgba(0,0,0,0.05);
  display: flex; flex-direction: column; align-items: center; gap: 4px;
  transition: border-color 0.22s, box-shadow 0.22s, transform 0.22s;
}
.jpg__card:hover { border-color: #86efac; transform: translateY(-2px); box-shadow: 0 6px 18px rgba(72,175,74,0.12); }
.jpg__card--active { border-color: #48AF4A; background: #f8fff8; box-shadow: 0 4px 16px rgba(72,175,74,0.14); }
.jpg__card-open  { font-size: 10px; font-weight: 700; color: #9ca3af; text-transform: uppercase; letter-spacing: 1px; }
.jpg__card-range { font-size: 12px; font-weight: 600; color: #0b1c30; line-height: 1.45; margin-top: 2px; }
.jpg__pill { font-size: 10px; font-weight: 700; padding: 3px 10px; border-radius: 100px; margin-top: 5px; }
.jpg__kuota { font-size: 10px; color: #9ca3af; margin-top: 2px; }

/* ── JADWAL — fallback + cards grid ── */
.jadwal-cards-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 20px; }
.jcard {
  border: 1.5px solid #e5e7eb; border-radius: 16px; padding: 24px;
  background: #fff; transition: transform .2s, box-shadow .2s;
  display: flex; flex-direction: column; gap: 6px;
}
.jcard:hover        { transform: translateY(-2px); box-shadow: 0 8px 24px rgba(0,0,0,.08); }
.jcard--active      { border-color: #48AF4A; background: #f0fdf4; }
.jcard__tag         { font-size: 10px; font-weight: 700; color: #48AF4A; background: #dcfce7; padding: 3px 10px; border-radius: 100px; display: inline-block; align-self: flex-start; }
.jcard__name        { font-size: 13px; font-weight: 700; color: #0b1c30; }
.jcard__range       { font-size: 12px; color: #6b7280; }
.jcard__status      { font-size: 13px; font-weight: 700; margin-top: 6px; }
.jcard__kuota       { font-size: 11px; color: #6b7280; margin-top: 4px; }
.jcard__pill        { font-size: 11px; font-weight: 700; padding: 3px 10px; border-radius: 100px; align-self: flex-start; margin-top: 8px; }

/* ── JADWAL — timeline vertical ── */
.jtv            { display: flex; flex-direction: column; gap: 0; max-width: 680px; margin: 0 auto; }
.jtv__item      { display: flex; gap: 16px; align-items: flex-start; }
.jtv__left      { display: flex; flex-direction: column; align-items: center; flex-shrink: 0; width: 24px; }
.jtv__dot       { width: 16px; height: 16px; border-radius: 50%; flex-shrink: 0; margin-top: 16px; border: 2px solid #fff; box-shadow: 0 0 0 2px #e5e7eb; }
.jtv__dot.status--green  { background: #48AF4A; box-shadow: 0 0 0 3px rgba(72,175,74,.2); }
.jtv__dot.status--yellow { background: #f59e0b; box-shadow: 0 0 0 3px rgba(245,158,11,.2); }
.jtv__dot.status--gray   { background: #9ca3af; }
.jtv__line      { width: 2px; background: #e5e7eb; flex: 1; min-height: 24px; }
.jtv__card      { flex: 1; border: 1.5px solid #e5e7eb; border-radius: 14px; padding: 18px 20px; background: #fff; margin: 8px 0 16px; display: flex; flex-direction: column; gap: 5px; transition: border-color .2s; }
.jtv__card:hover { border-color: #86efac; }
.jtv__card--active { border-color: #48AF4A; background: #f0fdf4; }
.jtv__tag       { font-size: 10px; font-weight: 700; color: #48AF4A; background: #dcfce7; padding: 3px 9px; border-radius: 100px; align-self: flex-start; }
.jtv__name      { font-size: 14px; font-weight: 700; color: #0b1c30; }
.jtv__range     { font-size: 12px; color: #6b7280; }
.jtv__pill      { font-size: 11px; font-weight: 700; padding: 3px 10px; border-radius: 100px; align-self: flex-start; margin-top: 4px; }
.jtv__kuota     { font-size: 11px; color: #6b7280; }

/* ── JADWAL — timeline horizontal ── */
.jth { display: flex; align-items: flex-start; justify-content: center; gap: 0; flex-wrap: wrap; }
.jth__step { display: flex; flex-direction: column; align-items: center; flex: 1; min-width: 100px; max-width: 220px; }
.jth__head { display: flex; align-items: center; width: 100%; }
.jth__circle {
  width: 40px; height: 40px; border-radius: 50%; display: flex; align-items: center; justify-content: center;
  font-weight: 800; font-size: 14px; color: #fff; flex-shrink: 0; z-index: 1;
}
.jth__circle.status--green  { background: #48AF4A; }
.jth__circle.status--yellow { background: #f59e0b; }
.jth__circle.status--gray   { background: #9ca3af; }
.jth__line     { flex: 1; height: 2px; background: #e5e7eb; }
.jth__body     { text-align: center; padding: 12px 8px 0; display: flex; flex-direction: column; align-items: center; gap: 4px; }
.jth__name     { font-size: 12.5px; font-weight: 700; color: #0b1c30; }
.jth__range    { font-size: 11px; color: #6b7280; }
.jth__pill     { font-size: 10px; font-weight: 700; padding: 2px 9px; border-radius: 100px; margin-top: 4px; }

/* ── JADWAL — compact list ── */
.jcmp { display: flex; flex-direction: column; gap: 8px; max-width: 760px; margin: 0 auto; }
.jcmp__row {
  display: flex; align-items: center; gap: 14px;
  padding: 14px 18px; background: #fff; border-radius: 10px;
  border: 1.5px solid #e5e7eb; transition: border-color .15s;
}
.jcmp__row:hover { border-color: #86efac; }
.jcmp__dot { width: 12px; height: 12px; border-radius: 50%; flex-shrink: 0; }
.jcmp__dot.status--green  { background: #48AF4A; }
.jcmp__dot.status--yellow { background: #f59e0b; }
.jcmp__dot.status--gray   { background: #9ca3af; }
.jcmp__name  { font-size: 13px; font-weight: 700; color: #0b1c30; flex: 1; }
.jcmp__range { font-size: 12px; color: #6b7280; white-space: nowrap; }
.jcmp__pill  { font-size: 11px; font-weight: 700; padding: 3px 12px; border-radius: 100px; flex-shrink: 0; }

/* ── PERSYARATAN ── */
.docs-desktop { display: grid; grid-template-columns: repeat(3, 1fr); gap: 24px; }
.docs-mobile  { display: none; }
.docs-carousel { touch-action: pan-y; user-select: none; }
.docs-card {
  border: 1.5px solid #e5e7eb; border-radius: 16px; padding: 28px;
  background: #fff; transition: border-color 0.2s, transform 0.2s, box-shadow 0.2s;
  box-shadow: 0 4px 16px rgba(0,0,0,0.07);
}
.docs-card:hover { border-color: #48AF4A; transform: translateY(-3px); box-shadow: 0 8px 28px rgba(72,175,74,0.13); }
.docs-card__icon { width: 48px; height: 48px; background: #e8f5e9; border-radius: 12px; display: flex; align-items: center; justify-content: center; margin-bottom: 16px; }
.docs-card h3    { font-size: 15px; font-weight: 700; color: #0b1c30; margin: 0 0 12px; }
.docs-card ul    { list-style: none; padding: 0; margin: 0; display: flex; flex-direction: column; gap: 8px; }
.docs-card li    { display: flex; align-items: center; gap: 8px; font-size: 13px; color: #374151; }
.docs-card--full { width: 100%; }

/* ── K3 ── */
.k3-grid  { display: grid; grid-template-columns: 1fr 1fr; gap: 64px; align-items: center; }
.k3-left  { display: flex; flex-direction: column; gap: 24px; }
.k3-left h2 { font-size: 34px; font-weight: 800; margin: 0; line-height: 1.15; }
.k3-left p  { font-size: 15px; line-height: 1.7; margin: 0; opacity: 0.72; }
.k3-rules   { display: flex; flex-direction: column; gap: 14px; }
.k3-rule    { display: flex; align-items: center; gap: 16px; }
.k3-rule__num {
  width: 38px; height: 38px; border-radius: 10px; flex-shrink: 0;
  background: rgba(72,175,74,0.15); border: 1px solid rgba(72,175,74,0.35);
  display: flex; align-items: center; justify-content: center;
  font-size: 12px; font-weight: 800; color: #48AF4A; letter-spacing: 0.5px;
}
.k3-rule span:last-child { font-size: 14px; color: rgba(255,255,255,0.85); line-height: 1.5; }
.k3-items   { display: grid; grid-template-columns: repeat(3, 1fr); gap: 14px; }
.k3-item {
  background: rgba(255,255,255,0.06); border: 1px solid rgba(255,255,255,0.1);
  border-radius: 12px; padding: 22px 14px; display: flex; flex-direction: column;
  align-items: center; gap: 10px; text-align: center;
  transition: background 0.2s, transform 0.2s, border-color 0.2s;
}
.k3-item:hover { background: rgba(72,175,74,0.15); transform: translateY(-3px); border-color: rgba(72,175,74,0.3); }
.k3-item__icon { width: 40px; height: 40px; background: rgba(72,175,74,0.15); border-radius: 10px; display: flex; align-items: center; justify-content: center; }
.k3-item span  { font-size: 12px; font-weight: 600; color: rgba(255,255,255,0.75); line-height: 1.4; }

/* ── FAQ ── */
.faq-list { display: flex; flex-direction: column; gap: 10px; }
.faq-item {
  border: 1.5px solid #e2e8f0; border-radius: 14px; cursor: pointer;
  transition: border-color 0.2s, box-shadow 0.2s, background 0.15s;
  background: #fff; box-shadow: 0 1px 4px rgba(0,0,0,0.04); overflow: hidden;
}
.faq-item:hover { border-color: #a7d9a8; background: #fafffe; }
.faq-item--open { border-color: #48AF4A; box-shadow: 0 4px 20px rgba(72,175,74,0.11); }
.faq-item__q { display: flex; align-items: center; gap: 14px; padding: 16px 20px; }
.faq-num {
  flex-shrink: 0; width: 30px; height: 30px; border-radius: 8px;
  background: #f0fdf4; color: #48AF4A; font-size: 11px; font-weight: 800;
  display: flex; align-items: center; justify-content: center; letter-spacing: 0.5px;
  transition: background 0.2s, color 0.2s;
}
.faq-item--open .faq-num { background: #48AF4A; color: #fff; }
.faq-question { flex: 1; font-size: 14px; font-weight: 600; color: #0b1c30; line-height: 1.5; }
.faq-chevron { transition: transform 0.3s ease; flex-shrink: 0; }
.faq-chevron--open { transform: rotate(180deg); }
.faq-item__a {
  max-height: 0; overflow: hidden;
  transition: max-height 0.35s ease, opacity 0.25s ease;
  opacity: 0;
}
.faq-item__a--open { max-height: 400px; opacity: 1; }
.faq-item__a-inner { padding: 0 20px 18px 64px; font-size: 13.5px; color: #4b5563; line-height: 1.75; }

/* ── FAQ more link ── */
.faq-more-wrap { display: flex; justify-content: center; margin-top: 28px; }
.faq-more-btn {
  display: inline-flex; align-items: center; gap: 8px;
  padding: 11px 28px; border-radius: 999px;
  border: 2px solid #48AF4A; color: #48AF4A;
  font-size: 14px; font-weight: 600; text-decoration: none;
  background: transparent; transition: all .2s;
}
.faq-more-btn:hover { background: #48AF4A; color: #fff; }

/* ── CTA BANNER ── */
.cta-banner {
  background: linear-gradient(135deg, #0b1c30 0%, #1a3a1f 60%, #0d2b10 100%);
  padding: 80px 0; text-align: center;
}
.cta-banner h2 { font-size: 36px; font-weight: 800; color: #fff; margin: 0 0 12px; }
.cta-banner p  { font-size: 15px; color: rgba(255,255,255,0.65); max-width: 540px; margin: 0 auto 32px; line-height: 1.7; }
.cta-banner__btns { display: flex; gap: 14px; justify-content: center; flex-wrap: wrap; }

/* ── FOOTER ── */
.footer { background: #0b1c30; padding: 64px 0 0; }
.footer__grid { display: grid; grid-template-columns: 2fr 1fr 1fr 1fr; gap: 48px; padding-bottom: 48px; }
.footer__brand-col { display: flex; flex-direction: column; gap: 14px; }
.footer__brand-link { display: inline-flex; align-items: center; gap: 8px; font-size: 15px; font-weight: 700; color: #fff; text-decoration: none; }
.footer__brand-link strong { color: #48AF4A; }
.footer__logo  { height: 36px; width: auto; object-fit: contain; background: #fff; padding: 4px 8px; border-radius: 8px; }
.footer__desc  { font-size: 13px; color: rgba(255,255,255,0.45); line-height: 1.65; margin: 0; }
.footer__address { font-size: 12px; color: rgba(255,255,255,0.28); line-height: 1.7; margin: 0; }
.footer__col-title { font-size: 11px; font-weight: 700; color: rgba(255,255,255,0.88); text-transform: uppercase; letter-spacing: 1.2px; margin-bottom: 16px; }
.footer ul   { list-style: none; padding: 0; margin: 0; display: flex; flex-direction: column; gap: 10px; }
.footer li   { font-size: 13px; color: rgba(255,255,255,0.45); }
.footer li a { color: rgba(255,255,255,0.5); text-decoration: none; transition: color 0.15s; }
.footer li a:hover { color: #86efac; }
.footer__bottom { border-top: 1px solid rgba(255,255,255,0.07); }
.footer__bottom-inner {
  display: flex; justify-content: space-between; align-items: center;
  padding: 20px 0; font-size: 12px; color: rgba(255,255,255,0.28);
}

/* ── RESPONSIVE ── */
@media (max-width: 1024px) {
  .jadwal-cards-grid { grid-template-columns: repeat(2, 1fr); }
  .jth { flex-wrap: wrap; }
  .jth__step { min-width: 140px; }
  .k3-grid     { grid-template-columns: 1fr; gap: 40px; }
  .alur-zigzag { grid-template-columns: 1fr; gap: 32px; }
  .alur-zigzag--flip .alur-zig-text { order: 1; }
  .alur-zigzag--flip .alur-zig-img  { order: 2; }
}

@media (max-width: 768px) {
  .container { padding: 0 20px; }

  /* Navbar mobile */
  .navbar__links        { display: none; }
  .navbar__action--desktop { display: none !important; }
  .navbar__hamburger    { display: flex; }
  .navbar__inner        { position: relative; }

  /* Hero */
  .hero { padding-top: 20px; padding-bottom: 60px; }
  .hero__title { font-size: 22px; white-space: normal; }
  .hero__title .text-green { white-space: normal; }
  .hero__sub    { font-size: 15px; }
  .hero__stats  { padding: 16px 20px; gap: 0; }
  .stat-item    { padding: 0 16px; }
  .stat__num    { font-size: 20px; }
  .stat__label  { font-size: 10px; }
  .stat-divider { height: 32px; }
  .hero__cta    { gap: 10px; }

  /* Sections */
  .section      { padding: 60px 0; }
  .section-header h2 { font-size: 26px; }
  .section-header { margin-bottom: 40px; }
  .alur-header h2 { font-size: 26px; }

  /* Docs */
  .docs-grid  { grid-template-columns: 1fr 1fr; }
  .docs-card  { padding: 18px; }

  /* K3 */
  .k3-items   { grid-template-columns: repeat(3, 1fr); }

  /* Footer */
  .footer__grid { grid-template-columns: 1fr 1fr; gap: 32px; }
  .footer__bottom-inner { flex-direction: column; gap: 6px; text-align: center; }

  /* Alur: 2-col → 1-col di tablet/mobile, gambar di bawah teks */
  .alur-main-grid     { grid-template-columns: 1fr; gap: 28px; }
  .alur-right-col     { } /* gambar ikut flow normal = di bawah teks */
  .alur-mobile-nav    { display: flex; }
  .alur-tabs--desktop { display: none; } /* sembunyikan tabs di mobile, pakai swipe */
  .alur-zig-watermark { font-size: 72px; }
  .alur-zig-title     { font-size: 22px; }

  /* Docs: carousel mode di mobile */
  .docs-desktop { display: none; }
  .docs-mobile  { display: block; }

  /* Jadwal period-grid */
  .jpg-wrap { padding: 28px 20px 28px; }
  .jpg__period { min-width: 120px; max-width: 150px; }
  .jpg__connector { min-width: 18px; max-width: 48px; }
  .jpg__circle { width: 62px; height: 62px; }
  .jpg__circle-num { font-size: 22px; }
  .jpg__connector { padding-top: 31px; }
  .jpg__badge { font-size: 8px; padding: 5px 8px; }

  /* Jadwal other */
  .jtv { padding: 0 4px; }
  .jth__step { min-width: 120px; }
}

@media (max-width: 480px) {
  .container { padding: 0 16px; }

  /* Hero */
  .hero__title { font-size: 18px; white-space: normal; }
  .hero__stats { flex-wrap: wrap; padding: 14px 16px; gap: 12px; }
  .stat-divider { display: none; }
  .stat-item    { padding: 0 10px; }
  .hero__cta    { flex-direction: column; }
  .hero__cta .btn-primary, .hero__cta .btn-outline-white { width: 100%; justify-content: center; }

  /* Jadwal */
  .jadwal-cards-grid { grid-template-columns: 1fr; }
  .jcmp__range { display: none; }
  .jth { gap: 0; }
  .jth__step { min-width: 80px; }
  .jth__circle { width: 32px; height: 32px; font-size: 12px; }

  /* Docs + K3 */
  .k3-items    { grid-template-columns: repeat(2, 1fr); }
  .k3-rule     { font-size: 12px; }

  /* CTA + Footer */
  .cta-banner h2 { font-size: 24px; }
  .cta-banner__btns { flex-direction: column; align-items: center; }
  .cta-banner__btns .btn-white, .cta-banner__btns .btn-outline-white { width: 100%; justify-content: center; }
  .footer__grid { grid-template-columns: 1fr; }
  .footer { padding: 48px 0 0; }

  /* Alur */
  .alur-tab { padding: 7px 14px; font-size: 12px; }
  .alur-zig-watermark { font-size: 56px; }
  .alur-zig-title { font-size: 20px; }
  .alur-zig-desc  { font-size: 14px; }

  /* Jadwal period-grid mobile kecil */
  .jpg__period { min-width: 100px; max-width: 130px; }
  .jpg__circle { width: 56px; height: 56px; }
  .jpg__circle-num { font-size: 20px; }
  .jpg__connector { padding-top: 28px; }
  .jpg__card { padding: 10px 8px; }
  .jpg__card-range { font-size: 11px; }

  /* Mobile menu wider padding */
  .navbar__mobile-menu a { padding: 14px 20px; }
}
</style>
