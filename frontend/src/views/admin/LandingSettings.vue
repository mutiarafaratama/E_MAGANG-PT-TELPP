<template>
  <div class="admin-layout">
    <!-- SIDEBAR -->
    <aside class="sidebar">
      <div class="sidebar__brand">
        <svg width="22" height="22" viewBox="0 0 24 24" fill="none">
          <path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z" fill="#86efac"/>
        </svg>
        <div>
          <div class="sidebar__brand-name">e-Magang</div>
          <div class="sidebar__brand-sub">PT TELPP</div>
        </div>
      </div>

      <nav class="sidebar__nav">
        <div class="sidebar__section-label">Konten</div>
        <a href="#" class="sidebar__link sidebar__link--active">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none">
            <path d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
          </svg>
          Landing Page
        </a>
      </nav>

      <div class="sidebar__footer">
        <div class="sidebar__user">
          <div class="sidebar__avatar">{{ userInitial }}</div>
          <div>
            <div class="sidebar__user-name">{{ user?.nama_lengkap }}</div>
            <div class="sidebar__user-role">{{ user?.role }}</div>
          </div>
        </div>
        <button class="btn-logout" @click="logout">Keluar</button>
      </div>
    </aside>

    <!-- MAIN -->
    <main class="admin-main">
      <div class="admin-header">
        <div>
          <h1>Pengaturan Landing Page</h1>
          <p>Kelola konten dan gambar hero yang tampil di halaman publik.</p>
        </div>
        <div class="header-actions">
          <a href="/" target="_blank" class="btn-preview">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none">
              <path d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
            Lihat Landing Page
          </a>
        </div>
      </div>

      <div v-if="pageLoading" class="loading-state">
        <div class="spinner-lg"></div>
        <p>Memuat data...</p>
      </div>

      <div v-else class="settings-grid">
        <!-- HERO IMAGE CARD -->
        <div class="card card--full">
          <div class="card-header">
            <div class="card-title">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none">
                <path d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" stroke="#48AF4A" stroke-width="2" stroke-linecap="round"/>
              </svg>
              Gambar Hero
            </div>
            <span class="card-badge">Background halaman login & register</span>
          </div>

          <div class="hero-editor">
            <!-- Preview -->
            <div class="hero-preview" :style="heroPreviewBg">
              <div class="hero-preview__overlay"></div>
              <div class="hero-preview__label">Preview Hero Background</div>
              <div v-if="!heroImageUrl" class="hero-preview__empty">
                <svg width="32" height="32" viewBox="0 0 24 24" fill="none">
                  <path d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" stroke="rgba(255,255,255,0.4)" stroke-width="1.5" stroke-linecap="round"/>
                </svg>
                <span>Belum ada gambar — menggunakan gradient default</span>
              </div>
            </div>

            <!-- Upload & URL -->
            <div class="hero-controls">
              <div class="control-section">
                <div class="control-label">Upload Gambar Baru</div>
                <div class="upload-area" @click="triggerFileInput" @dragover.prevent @drop.prevent="onDrop">
                  <input ref="fileInput" type="file" accept="image/*" class="hidden-input" @change="onFileChange" />
                  <svg width="24" height="24" viewBox="0 0 24 24" fill="none">
                    <path d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" stroke="#48AF4A" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                  </svg>
                  <div class="upload-text">
                    <strong>Klik untuk pilih file</strong> atau drag & drop
                    <span>PNG, JPG, WEBP — maks 10MB</span>
                  </div>
                </div>
                <div v-if="selectedFile" class="file-selected">
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none">
                    <path d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" stroke="#48AF4A" stroke-width="2" stroke-linecap="round"/>
                  </svg>
                  {{ selectedFile.name }}
                  <button @click="clearFile" class="clear-file">×</button>
                </div>
                <button
                  v-if="selectedFile"
                  class="btn-upload"
                  :disabled="uploading"
                  @click="uploadHeroImage"
                >
                  <span v-if="uploading" class="spinner-sm"></span>
                  {{ uploading ? "Mengupload..." : "Upload Gambar" }}
                </button>
              </div>

              <div class="divider-or">atau</div>

              <div class="control-section">
                <div class="control-label">Gunakan URL Gambar</div>
                <div class="url-input-wrap">
                  <input
                    v-model="heroUrlInput"
                    type="url"
                    placeholder="https://example.com/gambar.jpg"
                    class="url-input"
                  />
                  <button class="btn-use-url" @click="applyHeroUrl" :disabled="!heroUrlInput">
                    Terapkan
                  </button>
                </div>
              </div>

              <div class="control-section">
                <div class="control-label">URL Gambar Saat Ini</div>
                <div class="current-url">
                  <span v-if="heroImageUrl" class="url-text">{{ heroImageUrl }}</span>
                  <span v-else class="url-empty">Tidak ada gambar — gradient default aktif</span>
                  <button v-if="heroImageUrl" class="btn-remove-img" @click="removeHeroImage">
                    Hapus Gambar
                  </button>
                </div>
              </div>

              <div v-if="heroSaveMsg" class="save-msg" :class="heroSaveMsg.type">
                {{ heroSaveMsg.text }}
              </div>
            </div>
          </div>
        </div>

        <!-- TEKS HERO CARD -->
        <div class="card">
          <div class="card-header">
            <div class="card-title">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none">
                <path d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" stroke="#48AF4A" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              </svg>
              Teks Hero
            </div>
          </div>
          <div class="form-fields">
            <div class="form-group">
              <label>Judul Hero</label>
              <input v-model="content.hero_title" type="text" placeholder="Program Magang PT TELPP..." />
            </div>
            <div class="form-group">
              <label>Subjudul Hero</label>
              <textarea v-model="content.hero_subtitle" rows="3" placeholder="Deskripsi singkat..."></textarea>
            </div>
            <button class="btn-save" @click="saveField('hero_title'); saveField('hero_subtitle')" :disabled="saving.hero_title">
              <span v-if="saving.hero_title" class="spinner-sm"></span>
              Simpan Teks Hero
            </button>
          </div>
        </div>

        <!-- KONTEN LAINNYA -->
        <div class="card">
          <div class="card-header">
            <div class="card-title">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none">
                <path d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" stroke="#48AF4A" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              </svg>
              Konten Lainnya
            </div>
          </div>
          <div class="form-fields">
            <div class="form-group">
              <label>Syarat Umum</label>
              <textarea v-model="content.syarat_umum" rows="5" placeholder="1. ...&#10;2. ..."></textarea>
              <button class="btn-save-small" @click="saveField('syarat_umum')" :disabled="saving.syarat_umum">
                {{ saving.syarat_umum ? "..." : "Simpan" }}
              </button>
            </div>
            <div class="form-group">
              <label>Ketentuan Magang</label>
              <textarea v-model="content.ketentuan_magang" rows="5" placeholder="1. ...&#10;2. ..."></textarea>
              <button class="btn-save-small" @click="saveField('ketentuan_magang')" :disabled="saving.ketentuan_magang">
                {{ saving.ketentuan_magang ? "..." : "Simpan" }}
              </button>
            </div>
            <div class="form-group">
              <label>Kontak HRD</label>
              <input v-model="content.kontak_hrd" type="text" placeholder="email@telpp.co.id | (0734) xxx-xxx" />
              <button class="btn-save-small" @click="saveField('kontak_hrd')" :disabled="saving.kontak_hrd">
                {{ saving.kontak_hrd ? "..." : "Simpan" }}
              </button>
            </div>
          </div>
        </div>

        <!-- FAQ CARD -->
        <div class="card card--full">
          <div class="card-header">
            <div class="card-title">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none">
                <path d="M8.228 9c.549-1.165 2.03-2 3.772-2 2.21 0 4 1.343 4 3 0 1.4-1.278 2.575-3.006 2.907-.542.104-.994.54-.994 1.093m0 3h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" stroke="#48AF4A" stroke-width="2" stroke-linecap="round"/>
              </svg>
              Kelola FAQ
            </div>
          </div>

          <div class="faq-editor">
            <!-- Add new FAQ -->
            <div class="faq-new">
              <div class="faq-new__title">Tambah FAQ Baru</div>
              <div class="form-group">
                <label>Pertanyaan</label>
                <input v-model="newFaq.pertanyaan" type="text" placeholder="Apa persyaratan magang?" />
              </div>
              <div class="form-group">
                <label>Jawaban</label>
                <textarea v-model="newFaq.jawaban" rows="3" placeholder="Penjelasan jawaban..."></textarea>
              </div>
              <button class="btn-save" @click="addFaq" :disabled="savingFaq || !newFaq.pertanyaan || !newFaq.jawaban">
                <span v-if="savingFaq" class="spinner-sm"></span>
                Tambah FAQ
              </button>
            </div>

            <!-- FAQ List -->
            <div class="faq-list">
              <div class="faq-list__title">Daftar FAQ ({{ faqs.length }})</div>
              <div v-if="faqs.length === 0" class="faq-empty">Belum ada FAQ.</div>
              <div v-for="faq in faqs" :key="faq.id" class="faq-item">
                <div class="faq-item__header">
                  <span class="faq-item__num">{{ faq.urutan }}</span>
                  <div class="faq-item__body">
                    <div class="faq-item__q">{{ faq.pertanyaan }}</div>
                    <div class="faq-item__a">{{ faq.jawaban }}</div>
                  </div>
                  <div class="faq-item__actions">
                    <span class="faq-badge" :class="faq.is_active ? 'faq-badge--active' : 'faq-badge--inactive'">
                      {{ faq.is_active ? "Aktif" : "Nonaktif" }}
                    </span>
                    <button class="btn-delete" @click="deleteFaq(faq.id)">
                      <svg width="14" height="14" viewBox="0 0 24 24" fill="none">
                        <path d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                      </svg>
                      Hapus
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import api from "@/lib/api";
import { useAuth } from "@/hooks/useAuth";

const { user, logout } = useAuth();

const userInitial = computed(() =>
  user.value?.nama_lengkap?.charAt(0).toUpperCase() || "A"
);

interface FAQ {
  id: string;
  pertanyaan: string;
  jawaban: string;
  urutan: number;
  is_active: boolean;
}

const pageLoading = ref(true);
const content = ref<Record<string, string>>({
  hero_title: "",
  hero_subtitle: "",
  syarat_umum: "",
  ketentuan_magang: "",
  kontak_hrd: "",
});
const heroImageUrl = ref<string | null>(null);
const heroUrlInput = ref("");
const faqs = ref<FAQ[]>([]);
const saving = ref<Record<string, boolean>>({});
const savingFaq = ref(false);
const uploading = ref(false);
const selectedFile = ref<File | null>(null);
const fileInput = ref<HTMLInputElement | null>(null);
const heroSaveMsg = ref<{ type: string; text: string } | null>(null);
const newFaq = ref({ pertanyaan: "", jawaban: "" });

const heroPreviewBg = computed(() => {
  if (heroImageUrl.value) {
    return {
      background: `url('${heroImageUrl.value}') center/cover no-repeat`,
    };
  }
  return {
    background: "linear-gradient(135deg, #0b1c30 0%, #1a3a1f 55%, #0d2b10 100%)",
  };
});

async function loadContent() {
  try {
    const [contentRes, faqRes] = await Promise.all([
      api.get("/api/landing/content"),
      api.get("/api/landing/faq"),
    ]);
    const data = contentRes.data.data as Record<string, string>;
    Object.keys(content.value).forEach((key) => {
      if (data[key]) content.value[key] = data[key];
    });
    if (data.hero_image) {
      heroImageUrl.value = data.hero_image;
    }
    faqs.value = faqRes.data.data || [];
  } catch (err) {
    console.error("Gagal memuat data:", err);
  } finally {
    pageLoading.value = false;
  }
}

async function saveField(kunci: string) {
  saving.value[kunci] = true;
  try {
    await api.put("/api/landing/content", {
      kunci,
      nilai: content.value[kunci],
      tipe: "text",
    });
  } catch (err) {
    console.error("Gagal simpan:", err);
  } finally {
    saving.value[kunci] = false;
  }
}

async function saveHeroImage(url: string, tipe = "image") {
  heroSaveMsg.value = null;
  try {
    await api.put("/api/landing/content", {
      kunci: "hero_image",
      nilai: url,
      tipe,
    });
    heroImageUrl.value = url;
    heroSaveMsg.value = { type: "success", text: "Gambar hero berhasil disimpan!" };
    setTimeout(() => (heroSaveMsg.value = null), 3000);
  } catch (err) {
    heroSaveMsg.value = { type: "error", text: "Gagal menyimpan gambar. Coba lagi." };
  }
}

async function removeHeroImage() {
  heroSaveMsg.value = null;
  try {
    await api.put("/api/landing/content", {
      kunci: "hero_image",
      nilai: "",
      tipe: "image",
    });
    heroImageUrl.value = null;
    heroUrlInput.value = "";
    heroSaveMsg.value = { type: "success", text: "Gambar hero dihapus. Gradient default aktif." };
    setTimeout(() => (heroSaveMsg.value = null), 3000);
  } catch (err) {
    heroSaveMsg.value = { type: "error", text: "Gagal menghapus gambar." };
  }
}

function triggerFileInput() {
  fileInput.value?.click();
}

function onFileChange(e: Event) {
  const input = e.target as HTMLInputElement;
  if (input.files?.[0]) {
    selectedFile.value = input.files[0];
  }
}

function onDrop(e: DragEvent) {
  if (e.dataTransfer?.files?.[0]) {
    selectedFile.value = e.dataTransfer.files[0];
  }
}

function clearFile() {
  selectedFile.value = null;
  if (fileInput.value) fileInput.value.value = "";
}

async function uploadHeroImage() {
  if (!selectedFile.value) return;
  uploading.value = true;
  heroSaveMsg.value = null;
  try {
    const formData = new FormData();
    formData.append("file", selectedFile.value);
    formData.append("jenis", "hero_image");

    const res = await api.post("/api/dokumen/upload", formData, {
      headers: { "Content-Type": "multipart/form-data" },
    });
    const doc = res.data.data;
    const downloadUrl = `${import.meta.env.VITE_API_BASE_URL || "http://localhost:8080"}/api/dokumen/${doc.id}/download`;
    await saveHeroImage(downloadUrl);
    clearFile();
  } catch (err: any) {
    heroSaveMsg.value = {
      type: "error",
      text: err.response?.data?.message || "Upload gagal. Coba gunakan URL gambar.",
    };
  } finally {
    uploading.value = false;
  }
}

async function applyHeroUrl() {
  if (!heroUrlInput.value) return;
  await saveHeroImage(heroUrlInput.value, "text");
  heroUrlInput.value = "";
}

async function addFaq() {
  if (!newFaq.value.pertanyaan || !newFaq.value.jawaban) return;
  savingFaq.value = true;
  try {
    await api.post("/api/landing/faq", {
      pertanyaan: newFaq.value.pertanyaan,
      jawaban: newFaq.value.jawaban,
      urutan: faqs.value.length + 1,
      is_active: true,
    });
    newFaq.value = { pertanyaan: "", jawaban: "" };
    await loadContent();
  } catch (err) {
    console.error("Gagal tambah FAQ:", err);
  } finally {
    savingFaq.value = false;
  }
}

async function deleteFaq(id: string) {
  if (!confirm("Hapus FAQ ini?")) return;
  try {
    await api.delete(`/api/landing/faq/${id}`);
    faqs.value = faqs.value.filter((f) => f.id !== id);
  } catch (err) {
    console.error("Gagal hapus FAQ:", err);
  }
}

onMounted(loadContent);
</script>

<style scoped>
* { box-sizing: border-box; }

.admin-layout {
  display: flex;
  min-height: 100vh;
  font-family: "Poppins", sans-serif;
  background: #f0f4f8;
}

/* ── SIDEBAR ── */
.sidebar {
  width: 240px;
  flex-shrink: 0;
  background: #0b1c30;
  display: flex;
  flex-direction: column;
  padding: 24px 0;
}

.sidebar__brand {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 0 20px 24px;
  border-bottom: 1px solid rgba(255,255,255,0.08);
}

.sidebar__brand-name {
  font-size: 14px;
  font-weight: 700;
  color: #fff;
  line-height: 1.2;
}

.sidebar__brand-sub {
  font-size: 11px;
  color: #86efac;
  font-weight: 600;
}

.sidebar__section-label {
  font-size: 10px;
  font-weight: 700;
  color: rgba(255,255,255,0.35);
  letter-spacing: 1.5px;
  text-transform: uppercase;
  padding: 20px 20px 8px;
}

.sidebar__nav { flex: 1; }

.sidebar__link {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 20px;
  font-size: 13px;
  font-weight: 500;
  color: rgba(255,255,255,0.6);
  text-decoration: none;
  transition: all 0.15s;
  border-left: 3px solid transparent;
}

.sidebar__link:hover { color: #fff; background: rgba(255,255,255,0.05); }

.sidebar__link--active {
  color: #86efac;
  background: rgba(72,175,74,0.1);
  border-left-color: #48AF4A;
}

.sidebar__footer {
  padding: 20px;
  border-top: 1px solid rgba(255,255,255,0.08);
}

.sidebar__user {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 14px;
}

.sidebar__avatar {
  width: 34px;
  height: 34px;
  border-radius: 50%;
  background: #48AF4A;
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  font-size: 14px;
  flex-shrink: 0;
}

.sidebar__user-name {
  font-size: 13px;
  font-weight: 600;
  color: #fff;
  line-height: 1.2;
  max-width: 140px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.sidebar__user-role {
  font-size: 11px;
  color: rgba(255,255,255,0.45);
  text-transform: capitalize;
}

.btn-logout {
  width: 100%;
  padding: 8px;
  background: rgba(255,255,255,0.07);
  border: 1px solid rgba(255,255,255,0.12);
  border-radius: 8px;
  color: rgba(255,255,255,0.65);
  font-size: 12px;
  font-weight: 600;
  cursor: pointer;
  font-family: "Poppins", sans-serif;
  transition: all 0.15s;
}

.btn-logout:hover { background: rgba(255,255,255,0.12); color: #fff; }

/* ── MAIN ── */
.admin-main {
  flex: 1;
  overflow-y: auto;
  padding: 36px 40px;
}

.admin-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  margin-bottom: 32px;
}

.admin-header h1 {
  font-size: 24px;
  font-weight: 800;
  color: #0b1c30;
  margin-bottom: 4px;
}

.admin-header p {
  font-size: 14px;
  color: #64748b;
}

.btn-preview {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 9px 18px;
  border: 1.5px solid #48AF4A;
  color: #48AF4A;
  border-radius: 8px;
  font-size: 13px;
  font-weight: 600;
  text-decoration: none;
  transition: all 0.15s;
  font-family: "Poppins", sans-serif;
}

.btn-preview:hover { background: #e8f5e9; }

/* ── LOADING ── */
.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px 0;
  gap: 16px;
  color: #64748b;
  font-size: 14px;
}

.spinner-lg {
  width: 36px;
  height: 36px;
  border: 3px solid #e5e7eb;
  border-top-color: #48AF4A;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

.spinner-sm {
  width: 14px;
  height: 14px;
  border: 2px solid rgba(255,255,255,0.4);
  border-top-color: #fff;
  border-radius: 50%;
  animation: spin 0.7s linear infinite;
  display: inline-block;
}

@keyframes spin { to { transform: rotate(360deg); } }

/* ── GRID ── */
.settings-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 24px;
}

.card--full { grid-column: 1 / -1; }

/* ── CARD ── */
.card {
  background: #fff;
  border-radius: 16px;
  border: 1px solid #e5e7eb;
  overflow: hidden;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 24px;
  border-bottom: 1px solid #f0f0f0;
}

.card-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 15px;
  font-weight: 700;
  color: #0b1c30;
}

.card-badge {
  font-size: 11px;
  background: #e8f5e9;
  color: #48AF4A;
  padding: 3px 10px;
  border-radius: 100px;
  font-weight: 600;
}

/* ── HERO EDITOR ── */
.hero-editor {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 0;
}

.hero-preview {
  position: relative;
  min-height: 280px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.hero-preview__overlay {
  position: absolute;
  inset: 0;
  background: rgba(0,0,0,0.35);
}

.hero-preview__label {
  position: absolute;
  top: 14px;
  left: 14px;
  font-size: 11px;
  font-weight: 600;
  color: rgba(255,255,255,0.7);
  background: rgba(0,0,0,0.4);
  padding: 4px 10px;
  border-radius: 6px;
  z-index: 1;
}

.hero-preview__empty {
  position: relative;
  z-index: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
  color: rgba(255,255,255,0.45);
  font-size: 12px;
  text-align: center;
  padding: 20px;
}

.hero-controls {
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 20px;
  border-left: 1px solid #f0f0f0;
}

.control-section { display: flex; flex-direction: column; gap: 10px; }

.control-label {
  font-size: 12px;
  font-weight: 700;
  color: #374151;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.upload-area {
  border: 2px dashed #d1d5db;
  border-radius: 10px;
  padding: 20px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
  cursor: pointer;
  transition: border-color 0.15s;
}

.upload-area:hover { border-color: #48AF4A; }

.hidden-input { display: none; }

.upload-text {
  text-align: center;
  font-size: 13px;
  color: #64748b;
  line-height: 1.6;
}

.upload-text strong { color: #48AF4A; }

.upload-text span {
  display: block;
  font-size: 11px;
  color: #94a3b8;
  margin-top: 2px;
}

.file-selected {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: #48AF4A;
  background: #e8f5e9;
  padding: 6px 12px;
  border-radius: 6px;
}

.clear-file {
  margin-left: auto;
  background: none;
  border: none;
  cursor: pointer;
  color: #64748b;
  font-size: 16px;
  line-height: 1;
  padding: 0;
}

.btn-upload {
  padding: 9px 18px;
  background: #48AF4A;
  color: #fff;
  border: none;
  border-radius: 8px;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  font-family: "Poppins", sans-serif;
  display: flex;
  align-items: center;
  gap: 6px;
  transition: opacity 0.15s;
}

.btn-upload:disabled { opacity: 0.6; cursor: not-allowed; }

.divider-or {
  text-align: center;
  font-size: 12px;
  color: #94a3b8;
  position: relative;
}

.divider-or::before,
.divider-or::after {
  content: "";
  position: absolute;
  top: 50%;
  width: 40%;
  height: 1px;
  background: #e5e7eb;
}

.divider-or::before { left: 0; }
.divider-or::after { right: 0; }

.url-input-wrap {
  display: flex;
  gap: 8px;
}

.url-input {
  flex: 1;
  padding: 9px 12px;
  border: 1.5px solid #e5e7eb;
  border-radius: 8px;
  font-size: 13px;
  font-family: "Poppins", sans-serif;
  outline: none;
  color: #0b1c30;
}

.url-input:focus { border-color: #48AF4A; }

.btn-use-url {
  padding: 9px 16px;
  background: #0b1c30;
  color: #fff;
  border: none;
  border-radius: 8px;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  font-family: "Poppins", sans-serif;
  white-space: nowrap;
  transition: opacity 0.15s;
}

.btn-use-url:disabled { opacity: 0.4; cursor: not-allowed; }

.current-url {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  background: #f9fafb;
  border-radius: 8px;
  padding: 10px 14px;
}

.url-text {
  font-size: 12px;
  color: #374151;
  word-break: break-all;
  flex: 1;
  line-height: 1.5;
}

.url-empty {
  font-size: 12px;
  color: #94a3b8;
  font-style: italic;
}

.btn-remove-img {
  flex-shrink: 0;
  padding: 4px 10px;
  background: #fef2f2;
  border: 1px solid #fca5a5;
  color: #dc2626;
  border-radius: 6px;
  font-size: 11px;
  font-weight: 600;
  cursor: pointer;
  font-family: "Poppins", sans-serif;
  transition: all 0.15s;
}

.btn-remove-img:hover { background: #fee2e2; }

.save-msg {
  padding: 10px 14px;
  border-radius: 8px;
  font-size: 13px;
  font-weight: 500;
}

.save-msg.success { background: #e8f5e9; color: #2e7d32; border: 1px solid #a5d6a7; }
.save-msg.error { background: #fef2f2; color: #dc2626; border: 1px solid #fca5a5; }

/* ── FORM FIELDS ── */
.form-fields {
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 7px;
}

.form-group label {
  font-size: 12px;
  font-weight: 600;
  color: #374151;
  text-transform: uppercase;
  letter-spacing: 0.3px;
}

.form-group input,
.form-group textarea {
  padding: 10px 14px;
  border: 1.5px solid #e5e7eb;
  border-radius: 8px;
  font-size: 14px;
  font-family: "Poppins", sans-serif;
  color: #0b1c30;
  outline: none;
  transition: border-color 0.15s;
  background: #fafafa;
  resize: vertical;
}

.form-group input:focus,
.form-group textarea:focus {
  border-color: #48AF4A;
  background: #fff;
}

.btn-save {
  padding: 11px 24px;
  background: #48AF4A;
  color: #fff;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 700;
  cursor: pointer;
  font-family: "Poppins", sans-serif;
  display: flex;
  align-items: center;
  gap: 6px;
  width: fit-content;
  transition: opacity 0.15s;
}

.btn-save:disabled { opacity: 0.6; cursor: not-allowed; }

.btn-save-small {
  padding: 7px 16px;
  background: #48AF4A;
  color: #fff;
  border: none;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 600;
  cursor: pointer;
  font-family: "Poppins", sans-serif;
  width: fit-content;
  transition: opacity 0.15s;
}

.btn-save-small:disabled { opacity: 0.6; cursor: not-allowed; }

/* ── FAQ EDITOR ── */
.faq-editor {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 0;
}

.faq-new {
  padding: 24px;
  border-right: 1px solid #f0f0f0;
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.faq-new__title {
  font-size: 14px;
  font-weight: 700;
  color: #0b1c30;
}

.faq-list {
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 12px;
  max-height: 480px;
  overflow-y: auto;
}

.faq-list__title {
  font-size: 14px;
  font-weight: 700;
  color: #0b1c30;
  margin-bottom: 4px;
}

.faq-empty {
  font-size: 13px;
  color: #94a3b8;
  text-align: center;
  padding: 30px 0;
}

.faq-item {
  background: #f9fafb;
  border-radius: 10px;
  padding: 14px;
  border: 1px solid #e5e7eb;
}

.faq-item__header {
  display: flex;
  align-items: flex-start;
  gap: 10px;
}

.faq-item__num {
  flex-shrink: 0;
  width: 24px;
  height: 24px;
  border-radius: 50%;
  background: #48AF4A;
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 11px;
  font-weight: 700;
  margin-top: 2px;
}

.faq-item__body { flex: 1; }

.faq-item__q {
  font-size: 13px;
  font-weight: 600;
  color: #0b1c30;
  margin-bottom: 4px;
  line-height: 1.4;
}

.faq-item__a {
  font-size: 12px;
  color: #64748b;
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.faq-item__actions {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 6px;
  flex-shrink: 0;
}

.faq-badge {
  font-size: 10px;
  font-weight: 700;
  padding: 2px 8px;
  border-radius: 100px;
}

.faq-badge--active { background: #e8f5e9; color: #2e7d32; }
.faq-badge--inactive { background: #f1f5f9; color: #94a3b8; }

.btn-delete {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 5px 10px;
  background: #fef2f2;
  border: 1px solid #fca5a5;
  color: #dc2626;
  border-radius: 6px;
  font-size: 11px;
  font-weight: 600;
  cursor: pointer;
  font-family: "Poppins", sans-serif;
  transition: all 0.15s;
}

.btn-delete:hover { background: #fee2e2; }
</style>
