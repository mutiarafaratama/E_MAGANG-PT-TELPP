<template>
  <div class="hero-admin">
    <div class="page-card">
    <!-- HEADER -->
    <div class="ha-header">
      <div>
        <h3 class="ha-title">Hero Section</h3>
        <p class="ha-sub">Customize tampilan utama landing page — judul, statistik card, dan background image</p>
      </div>
      <div class="ha-header-actions">
        <a href="/" target="_blank" class="btn-preview">Lihat Landing ↗</a>
        <button class="btn-save" @click="saveAll" :disabled="saving">
          <span v-if="saving" class="spinner-sm"></span>
          {{ saving ? 'Menyimpan…' : 'Simpan Semua' }}
        </button>
      </div>
    </div>

    <div v-if="loading" class="ha-loading"><div class="spinner"></div></div>

    <div v-else class="ha-grid">

      <!-- ── KOLOM KIRI ── -->
      <div class="ha-col">

        <!-- Background Image -->
        <div class="ha-card">
          <div class="ha-card-label">Background Image</div>
          <div class="ha-card-desc">Foto/gambar yang tampil di belakang hero. Resolusi direkomendasikan 1920×1080px.</div>
          <div class="bg-upload-area" @click="triggerBgInput" :class="{ 'bg-upload-area--uploading': bgUploading }">
            <input ref="bgFileInput" type="file" accept=".jpg,.jpeg,.png,.webp" style="display:none" @change="onBgFileChange" />
            <template v-if="bgPreview || form.hero_bg_url">
              <img :src="bgPreview || form.hero_bg_url" class="bg-preview-img" alt="Hero background" />
              <div class="bg-overlay-hint">Klik untuk ganti gambar</div>
            </template>
            <template v-else>
              <div class="bg-upload-icon">
                <svg width="32" height="32" viewBox="0 0 24 24" fill="none">
                  <rect x="3" y="3" width="18" height="18" rx="3" stroke="#94a3b8" stroke-width="1.5"/>
                  <path d="M3 16l5-5 4 4 3-3 4 4" stroke="#94a3b8" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
                  <circle cx="8.5" cy="8.5" r="1.5" fill="#94a3b8"/>
                </svg>
              </div>
              <div class="bg-upload-text">Klik untuk upload background</div>
              <div class="bg-upload-hint">JPG, PNG, WebP · Maks 10 MB · Rekomendasi 1920×1080</div>
            </template>
            <div v-if="bgUploading" class="bg-upload-spinner">
              <div class="spinner"></div>
              <span>Mengupload…</span>
            </div>
          </div>
          <div v-if="bgError" class="ha-error">{{ bgError }}</div>
          <div v-if="form.hero_bg_url" class="bg-remove">
            <button class="btn-remove-bg" @click.stop="removeBg">× Hapus background image</button>
          </div>
        </div>

        <!-- Opacity Slider -->
        <div class="ha-card" v-if="form.hero_bg_url || bgPreview">
          <div class="ha-card-label">Intensitas Background <span class="ha-tag-opacity">{{ form.hero_bg_opacity }}%</span></div>
          <div class="ha-card-desc">Seberapa kuat gambar background terlihat. Nilai tinggi = gambar lebih jelas, nilai rendah = lebih gelap.</div>
          <div class="opacity-slider-wrap">
            <svg class="opacity-icon" width="16" height="16" viewBox="0 0 24 24" fill="none"><path d="M21 12.79A9 9 0 1111.21 3 7 7 0 0021 12.79z" stroke="#94a3b8" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>
            <input type="range" min="0" max="100" step="5" v-model="form.hero_bg_opacity" class="opacity-slider" />
            <svg class="opacity-icon" width="16" height="16" viewBox="0 0 24 24" fill="none"><circle cx="12" cy="12" r="5" stroke="#94a3b8" stroke-width="2"/><line x1="12" y1="1" x2="12" y2="3" stroke="#94a3b8" stroke-width="2" stroke-linecap="round"/><line x1="12" y1="21" x2="12" y2="23" stroke="#94a3b8" stroke-width="2" stroke-linecap="round"/><line x1="4.22" y1="4.22" x2="5.64" y2="5.64" stroke="#94a3b8" stroke-width="2" stroke-linecap="round"/><line x1="18.36" y1="18.36" x2="19.78" y2="19.78" stroke="#94a3b8" stroke-width="2" stroke-linecap="round"/><line x1="1" y1="12" x2="3" y2="12" stroke="#94a3b8" stroke-width="2" stroke-linecap="round"/><line x1="21" y1="12" x2="23" y2="12" stroke="#94a3b8" stroke-width="2" stroke-linecap="round"/><line x1="4.22" y1="19.78" x2="5.64" y2="18.36" stroke="#94a3b8" stroke-width="2" stroke-linecap="round"/><line x1="18.36" y1="5.64" x2="19.78" y2="4.22" stroke="#94a3b8" stroke-width="2" stroke-linecap="round"/></svg>
          </div>
          <div class="opacity-preview-bar">
            <div class="opacity-preview-fill" :style="{ width: form.hero_bg_opacity + '%' }"></div>
          </div>
        </div>

        <!-- Judul -->
        <div class="ha-card">
          <div class="ha-card-label">Judul Utama</div>
          <div class="ha-card-desc">Baris pertama judul besar (warna putih).</div>
          <input v-model="form.hero_judul" type="text" class="ha-input" placeholder="Mulai Karir" />
        </div>

        <!-- Judul aksen -->
        <div class="ha-card">
          <div class="ha-card-label">Judul Aksen <span class="ha-tag-green">hijau</span></div>
          <div class="ha-card-desc">Baris kedua judul (warna hijau — nama perusahaan).</div>
          <input v-model="form.hero_judul_aksen" type="text" class="ha-input" placeholder="PT TELPP" />
        </div>

      </div>

      <!-- ── KOLOM KANAN ── -->
      <div class="ha-col">

        <!-- Statistik -->
        <div class="ha-card">
          <div class="ha-card-label">Statistik (4 item)</div>
          <div class="ha-card-desc">Angka dan label yang tampil sebagai kartu di hero.</div>
          <div class="stats-grid">
            <div v-for="i in 4" :key="i" class="stat-row">
              <div class="stat-row-num">
                <label class="stat-label">Angka</label>
                <input v-model="form[`hero_stat_${i}_num`]" type="text" class="ha-input ha-input-sm" :placeholder="defaultStats[i-1].num" />
              </div>
              <div class="stat-row-label">
                <label class="stat-label">Keterangan</label>
                <input v-model="form[`hero_stat_${i}_label`]" type="text" class="ha-input ha-input-sm" :placeholder="defaultStats[i-1].label" />
              </div>
            </div>
          </div>
        </div>

        <!-- Stat Card Opacity -->
        <div class="ha-card">
          <div class="ha-card-label">Transparansi Kartu Statistik <span class="ha-tag-opacity">{{ form.hero_stat_card_opacity }}%</span></div>
          <div class="ha-card-desc">Seberapa solid kartu statistik terlihat di atas background. 0% = transparan penuh, 100% = kartu putih solid.</div>
          <div class="opacity-slider-wrap">
            <svg class="opacity-icon" width="16" height="16" viewBox="0 0 24 24" fill="none"><path d="M21 12.79A9 9 0 1111.21 3 7 7 0 0021 12.79z" stroke="#94a3b8" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>
            <input type="range" min="0" max="100" step="5" v-model="form.hero_stat_card_opacity" class="opacity-slider" />
            <svg class="opacity-icon" width="16" height="16" viewBox="0 0 24 24" fill="none"><rect x="3" y="3" width="18" height="18" rx="3" stroke="#94a3b8" stroke-width="2"/></svg>
          </div>
          <!-- Preview kartu -->
          <div class="stat-card-preview">
            <div
              class="stat-card-preview-item"
              :style="{
                background: `rgba(255,255,255,${(parseInt(form.hero_stat_card_opacity)/100).toFixed(2)})`,
                boxShadow: '0 8px 32px rgba(0,0,0,0.22)',
                border: `1px solid rgba(255,255,255,${Math.min(parseInt(form.hero_stat_card_opacity)/100+0.15,0.45).toFixed(2)})`,
              }"
            >
              <div class="scp-num">500+</div>
              <div class="scp-label">Alumni Magang</div>
            </div>
            <div class="stat-card-preview-desc">Preview kartu statistik di atas background</div>
          </div>
        </div>

        <!-- Preview Live -->
        <div class="ha-card ha-card--preview">
          <div class="ha-card-label">Preview Mini</div>
          <div class="hero-preview" :style="previewStyle">
            <div class="hero-preview__overlay"></div>
            <div class="hero-preview__content">
              <div class="hp-title">
                {{ form.hero_judul || 'Mulai Karir' }}<br />
                <span class="hp-title-aksen">{{ form.hero_judul_aksen || 'PT TELPP' }}</span>
              </div>
              <div class="hp-stats">
                <div
                  v-for="i in 4" :key="i"
                  class="hp-stat-card"
                  :style="{
                    background: `rgba(255,255,255,${(parseInt(form.hero_stat_card_opacity)/100).toFixed(2)})`,
                    boxShadow: '0 4px 12px rgba(0,0,0,0.2)',
                  }"
                >
                  <div class="hp-stat-num">{{ form[`hero_stat_${i}_num`] || defaultStats[i-1].num }}</div>
                  <div class="hp-stat-label">{{ form[`hero_stat_${i}_label`] || defaultStats[i-1].label }}</div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Status / Notif -->
        <div v-if="saveError" class="ha-alert ha-alert--error">{{ saveError }}</div>
        <div v-if="saveSuccess" class="ha-alert ha-alert--success">{{ saveSuccess }}</div>

      </div>
    </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import api from '@/lib/api';

const loading    = ref(true);
const saving     = ref(false);
const saveError  = ref('');
const saveSuccess = ref('');
const bgUploading = ref(false);
const bgError    = ref('');
const bgPreview  = ref('');
const bgFileInput = ref<HTMLInputElement | null>(null);

const defaultStats = [
  { num: '500+', label: 'Alumni Magang' },
  { num: '30+',  label: 'Divisi Terbuka' },
  { num: '95%',  label: 'Kepuasan Peserta' },
  { num: '1–6 Bln', label: 'Durasi' },
];

const form = ref<Record<string, string>>({
  hero_judul:           'Mulai Karir',
  hero_judul_aksen:     'PT TELPP',
  hero_bg_url:          '',
  hero_bg_opacity:      '70',
  hero_stat_card_opacity: '20',
  hero_stat_1_num:   '500+',  hero_stat_1_label: 'Alumni Magang',
  hero_stat_2_num:   '30+',   hero_stat_2_label: 'Divisi Terbuka',
  hero_stat_3_num:   '95%',   hero_stat_3_label: 'Kepuasan Peserta',
  hero_stat_4_num:   '1–6 Bln', hero_stat_4_label: 'Durasi',
});

const previewStyle = computed(() => {
  const url = bgPreview.value || form.value.hero_bg_url;
  if (url) {
    return { backgroundImage: `url(${url})`, backgroundSize: 'cover', backgroundPosition: 'center' };
  }
  return {};
});

async function fetchContent() {
  loading.value = true;
  try {
    const r = await api.get('/api/landing/content');
    const d = r.data?.data ?? {};
    const keys = Object.keys(form.value);
    for (const k of keys) {
      if (d[k] !== undefined && d[k] !== null) form.value[k] = d[k];
    }
  } catch { /* pakai default */ }
  finally { loading.value = false; }
}

const HERO_KEYS = [
  'hero_judul','hero_judul_aksen','hero_bg_url','hero_bg_opacity','hero_stat_card_opacity',
  'hero_stat_1_num','hero_stat_1_label','hero_stat_2_num','hero_stat_2_label',
  'hero_stat_3_num','hero_stat_3_label','hero_stat_4_num','hero_stat_4_label',
];

async function saveAll() {
  saving.value = true; saveError.value = ''; saveSuccess.value = '';
  try {
    for (const k of HERO_KEYS) {
      await api.put('/api/landing/content', {
        kunci: k,
        nilai: form.value[k] ?? '',
        tipe: k === 'hero_bg_url' ? 'image' : 'text',
      });
    }
    saveSuccess.value = 'Hero section berhasil disimpan!';
    setTimeout(() => { saveSuccess.value = ''; }, 3500);
  } catch (e: any) {
    saveError.value = e.response?.data?.message ?? 'Gagal menyimpan. Coba lagi.';
  } finally { saving.value = false; }
}

function triggerBgInput() {
  if (!bgUploading.value) bgFileInput.value?.click();
}

async function onBgFileChange(e: Event) {
  const file = (e.target as HTMLInputElement).files?.[0];
  if (!file) return;
  bgError.value = '';
  bgPreview.value = URL.createObjectURL(file);
  bgUploading.value = true;
  try {
    const fd = new FormData();
    fd.append('gambar', file);
    const r = await api.post('/api/admin/hero/upload', fd, { headers: { 'Content-Type': 'multipart/form-data' } });
    form.value.hero_bg_url = r.data?.data?.url ?? '';
    bgPreview.value = form.value.hero_bg_url;
    saveSuccess.value = 'Background hero berhasil diupload!';
    setTimeout(() => { saveSuccess.value = ''; }, 3000);
  } catch (err: any) {
    bgError.value = err.response?.data?.message ?? 'Gagal upload gambar';
    bgPreview.value = '';
  } finally {
    bgUploading.value = false;
    (e.target as HTMLInputElement).value = '';
  }
}

async function removeBg() {
  if (!confirm('Hapus background image hero?')) return;
  form.value.hero_bg_url = '';
  bgPreview.value = '';
  try {
    await api.put('/api/landing/content', { kunci: 'hero_bg_url', nilai: '', tipe: 'image' });
    saveSuccess.value = 'Background dihapus.';
    setTimeout(() => { saveSuccess.value = ''; }, 2500);
  } catch { /* ignore */ }
}

onMounted(fetchContent);
</script>

<style scoped>
.hero-admin { display: flex; flex-direction: column; gap: 20px; }

.ha-header {
  display: flex; align-items: flex-start; justify-content: space-between; gap: 16px;
  flex-wrap: wrap;
}
.ha-title { font-size: 15px; font-weight: 700; color: #111827; margin: 0 0 4px; }
.ha-sub   { font-size: 12.5px; color: #6b7280; margin: 0; }
.ha-header-actions { display: flex; gap: 8px; align-items: center; flex-shrink: 0; }

.btn-preview {
  background: none; border: 1.5px solid #d1d5db; border-radius: 8px;
  padding: 7px 14px; font-size: 12px; font-weight: 600; color: #374151;
  cursor: pointer; text-decoration: none; font-family: "Poppins", sans-serif;
  transition: border-color 0.15s;
}
.btn-preview:hover { border-color: #48AF4A; color: #48AF4A; }
.btn-save {
  background: #48AF4A; color: #fff; border: none; border-radius: 8px;
  padding: 8px 20px; font-size: 13px; font-weight: 600;
  cursor: pointer; font-family: "Poppins", sans-serif;
  display: inline-flex; align-items: center; gap: 7px;
  transition: opacity 0.15s;
}
.btn-save:disabled { opacity: 0.65; cursor: not-allowed; }

.page-card { background:#fff; border:1px solid #e9f5e9; border-radius:14px; padding:20px; box-shadow:0 1px 3px rgba(13,40,24,.04); display:flex; flex-direction:column; gap:20px; }

.ha-loading { display: flex; justify-content: center; padding: 60px; }

.ha-grid {
  display: grid; grid-template-columns: 1fr 1fr; gap: 20px; align-items: start;
}

.ha-col { display: flex; flex-direction: column; gap: 16px; }

.ha-card {
  background: #fff; border-radius: 12px; border: 1px solid #e9f5e9;
  padding: 18px 20px; display: flex; flex-direction: column; gap: 10px;
  box-shadow: 0 1px 3px rgba(13,40,24,0.04);
}
.ha-card--preview { padding: 18px 20px; }

.ha-card-label { font-size: 12px; font-weight: 700; color: #111827; }
.ha-card-desc  { font-size: 11.5px; color: #6b7280; line-height: 1.5; margin-top: -4px; }
.ha-tag-green  { background: #dcfce7; color: #16a34a; font-size: 10px; padding: 2px 7px; border-radius: 4px; margin-left: 4px; font-weight: 600; }

.ha-input {
  border: 1.5px solid #e5e7eb; border-radius: 8px; padding: 9px 12px;
  font-size: 13px; font-family: "Poppins", sans-serif; color: #111827;
  outline: none; transition: border-color 0.15s; width: 100%; box-sizing: border-box;
}
.ha-input:focus { border-color: #48AF4A; }
.ha-input-sm { padding: 7px 10px; font-size: 12px; }

/* Background upload */
.bg-upload-area {
  border: 2px dashed #d1d5db; border-radius: 12px; padding: 0;
  cursor: pointer; position: relative; overflow: hidden;
  min-height: 180px; display: flex; flex-direction: column;
  align-items: center; justify-content: center; gap: 8px;
  background: #f9fafb; transition: border-color 0.15s;
}
.bg-upload-area:hover { border-color: #48AF4A; background: #f0faf0; }
.bg-upload-area--uploading { pointer-events: none; }
.bg-upload-icon { color: #94a3b8; }
.bg-upload-text { font-size: 13px; font-weight: 600; color: #374151; }
.bg-upload-hint { font-size: 11px; color: #9ca3af; }

.bg-preview-img {
  width: 100%; height: 180px; object-fit: cover;
  display: block; border-radius: 10px;
}
.bg-overlay-hint {
  position: absolute; inset: 0; background: rgba(0,0,0,0.35);
  display: flex; align-items: center; justify-content: center;
  font-size: 13px; font-weight: 600; color: #fff;
  border-radius: 10px; opacity: 0; transition: opacity 0.15s;
}
.bg-upload-area:hover .bg-overlay-hint { opacity: 1; }

.bg-upload-spinner {
  position: absolute; inset: 0; background: rgba(255,255,255,0.8);
  display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 8px;
  font-size: 12px; color: #48AF4A; font-weight: 600;
}

.btn-remove-bg {
  background: none; border: none; color: #ef4444; font-size: 11.5px;
  cursor: pointer; font-family: "Poppins", sans-serif; padding: 2px 0;
  text-align: left;
}
.btn-remove-bg:hover { text-decoration: underline; }

/* Stats grid */
.stats-grid { display: flex; flex-direction: column; gap: 12px; }
.stat-row { display: grid; grid-template-columns: 1fr 1.6fr; gap: 10px; }
.stat-row-num, .stat-row-label { display: flex; flex-direction: column; gap: 4px; }
.stat-label { font-size: 10.5px; font-weight: 600; color: #6b7280; }

/* Stat card opacity preview */
.stat-card-preview {
  display: flex; flex-direction: column; gap: 8px; margin-top: 4px;
}
.stat-card-preview-item {
  display: inline-flex; flex-direction: column; align-items: center;
  padding: 12px 20px; border-radius: 12px; width: fit-content;
  backdrop-filter: blur(10px);
  background-image: linear-gradient(135deg, #0b1c30 0%, #1a3a1f 100%);
  background-blend-mode: luminosity;
}
.stat-card-preview { background: linear-gradient(135deg, #0b1c30 0%, #1a3a1f 100%); border-radius: 10px; padding: 16px; }
.stat-card-preview-item { background-image: none !important; }
.scp-num   { font-size: 22px; font-weight: 800; color: #fff; line-height: 1; }
.scp-label { font-size: 10px; color: rgba(255,255,255,0.75); font-weight: 500; margin-top: 3px; }
.stat-card-preview-desc { font-size: 10.5px; color: #9ca3af; text-align: center; }

/* Mini Preview */
.hero-preview {
  border-radius: 10px; overflow: hidden;
  background: linear-gradient(135deg, #0b1c30 0%, #1a3a1f 55%, #0d2b10 100%);
  position: relative; min-height: 200px; padding: 24px 20px;
  display: flex; align-items: center; justify-content: center;
}
.hero-preview__overlay {
  position: absolute; inset: 0;
  background: linear-gradient(135deg, rgba(11,28,48,0.85) 0%, rgba(26,58,31,0.75) 55%, rgba(13,43,16,0.85) 100%);
  pointer-events: none;
}
.hero-preview__content {
  position: relative; z-index: 1;
  display: flex; flex-direction: column; gap: 12px; align-items: center; text-align: center; width: 100%;
}

.hp-title { font-size: 20px; font-weight: 800; color: #fff; line-height: 1.2; }
.hp-title-aksen { color: #86efac; }
.hp-stats { display: flex; gap: 8px; flex-wrap: wrap; justify-content: center; }
.hp-stat-card {
  padding: 8px 12px; border-radius: 8px; text-align: center;
  backdrop-filter: blur(8px); min-width: 56px;
  border: 1px solid rgba(255,255,255,0.2);
}
.hp-stat-num   { font-size: 13px; font-weight: 800; color: #fff; }
.hp-stat-label { font-size: 7.5px; color: rgba(255,255,255,0.7); font-weight: 500; margin-top: 2px; }

/* Alerts */
.ha-alert {
  padding: 11px 15px; border-radius: 8px; font-size: 12.5px; font-weight: 500;
}
.ha-alert--success { background: #f0fdf4; color: #16a34a; border: 1px solid #bbf7d0; }
.ha-alert--error   { background: #fef2f2; color: #dc2626; border: 1px solid #fecaca; }

.ha-error { font-size: 11.5px; color: #dc2626; }

/* Spinner */
.spinner    { width: 28px; height: 28px; border: 3px solid #e5e7eb; border-top-color: #48AF4A; border-radius: 50%; animation: spin 0.7s linear infinite; }
.spinner-sm { width: 14px; height: 14px; border: 2px solid rgba(255,255,255,0.4); border-top-color: #fff; border-radius: 50%; animation: spin 0.7s linear infinite; display: inline-block; }
@keyframes spin { to { transform: rotate(360deg); } }

/* Opacity slider */
.opacity-slider-wrap {
  display: flex; align-items: center; gap: 10px;
}
.opacity-icon { font-size: 16px; flex-shrink: 0; }
.opacity-slider {
  flex: 1; -webkit-appearance: none; appearance: none;
  height: 5px; border-radius: 3px; outline: none; cursor: pointer;
  background: linear-gradient(to right, #48AF4A 0%, #48AF4A v-bind("form.hero_bg_opacity + '%'"), #e5e7eb v-bind("form.hero_bg_opacity + '%'"), #e5e7eb 100%);
}
.opacity-slider::-webkit-slider-thumb {
  -webkit-appearance: none; width: 18px; height: 18px;
  border-radius: 50%; background: #48AF4A;
  border: 2px solid #fff; box-shadow: 0 1px 4px rgba(0,0,0,0.2);
  cursor: pointer;
}
.opacity-slider::-moz-range-thumb {
  width: 18px; height: 18px; border-radius: 50%; background: #48AF4A;
  border: 2px solid #fff; box-shadow: 0 1px 4px rgba(0,0,0,0.2); cursor: pointer;
}
.opacity-preview-bar {
  height: 4px; background: #f1f5f9; border-radius: 2px; overflow: hidden; margin-top: 6px;
}
.opacity-preview-fill {
  height: 100%; background: linear-gradient(to right, #e8f5e9, #48AF4A);
  border-radius: 2px; transition: width 0.1s;
}
.ha-tag-opacity {
  background: #f0fdf4; color: #16a34a; font-size: 10px;
  padding: 2px 7px; border-radius: 4px; margin-left: 4px; font-weight: 700;
}

/* Responsive */
@media (max-width: 900px) {
  .ha-grid { grid-template-columns: 1fr; }
}
</style>
