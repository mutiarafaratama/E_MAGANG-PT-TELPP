<template>
  <div class="ak">
    <!-- Top Bar -->
    <div class="ak__topbar">
      <div class="ak__topbar-left">
        <div class="ak__title">Knowledge Base</div>
        <div class="ak__sub">Data untuk chatbot NLP menjawab pertanyaan otomatis</div>
      </div>
      <button class="ak__btn-add" @click="openForm()">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none"><line x1="12" y1="5" x2="12" y2="19" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"/><line x1="5" y1="12" x2="19" y2="12" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"/></svg>
        Tambah
      </button>
    </div>

    <!-- Filter + Search -->
    <div class="ak__filter">
      <input v-model="search" class="ak__search" placeholder="Cari pertanyaan atau kata kunci...">
      <div class="ak__selects">
        <select v-model="filterKategori" class="ak__select">
          <option value="">Semua Kategori</option>
          <option value="umum">Umum</option>
          <option value="absensi">Absensi</option>
          <option value="dokumen">Dokumen</option>
          <option value="sertifikat">Sertifikat</option>
        </select>
        <select v-model="filterStatus" class="ak__select">
          <option value="">Semua Status</option>
          <option value="aktif">Aktif</option>
          <option value="nonaktif">Nonaktif</option>
        </select>
      </div>
    </div>

    <!-- NLP Tester -->
    <div class="ak__tester">
      <div class="ak__tester-label">Uji Chatbot NLP</div>
      <div class="ak__tester-row">
        <input
          v-model="nlpQuery"
          class="ak__tester-input"
          placeholder="Ketik pertanyaan untuk test..."
          @keydown.enter="testNLP"
        />
        <button class="ak__tester-btn" @click="testNLP" :disabled="!nlpQuery.trim() || nlpLoading">
          {{ nlpLoading ? 'Menguji...' : 'Uji' }}
        </button>
      </div>
      <div v-if="nlpResult" class="ak__tester-result" :class="nlpResult.terjawab ? 'ak__tester-result--ok' : 'ak__tester-result--fail'">
        <div class="ak__tester-result-head">
          <span class="ak__tester-status-dot" :class="nlpResult.terjawab ? 'dot--ok' : 'dot--fail'"></span>
          <span class="ak__tester-status-label">{{ nlpResult.terjawab ? 'Terjawab' : 'Tidak Terjawab' }}</span>
          <span class="ak__tester-skor">Skor: {{ (nlpResult.skor * 100).toFixed(0) }}% &middot; {{ nlpResult.sumber || '-' }}</span>
        </div>
        <div v-if="nlpResult.terjawab" class="ak__tester-answer">{{ nlpResult.jawaban }}</div>
        <div v-else class="ak__tester-answer ak__tester-answer--dim">Tidak ada knowledge yang cocok. Pertimbangkan menambah kata kunci baru.</div>
      </div>
    </div>

    <!-- List -->
    <div v-if="loading" class="ak__loading">
      <div class="ak__spinner"></div>
      <span>Memuat data...</span>
    </div>
    <div v-else-if="filtered.length === 0" class="ak__empty">
      <svg width="32" height="32" viewBox="0 0 24 24" fill="none"><path d="M9 9h6M9 12h6M9 15h4" stroke="#d1d5db" stroke-width="1.5" stroke-linecap="round"/><path d="M5 3h14a2 2 0 012 2v14a2 2 0 01-2 2H5a2 2 0 01-2-2V5a2 2 0 012-2z" stroke="#d1d5db" stroke-width="1.5"/></svg>
      <p>Tidak ada knowledge ditemukan.</p>
    </div>
    <div v-else class="ak__list">
      <div
        v-for="k in filtered"
        :key="k.id"
        class="ak__card"
        :class="{ 'ak__card--inactive': !k.is_active }"
      >
        <div class="ak__card-header">
          <div class="ak__card-pertanyaan">{{ k.pertanyaan }}</div>
          <div class="ak__card-actions">
            <span class="ak__badge" :class="'ak__badge--' + k.kategori">{{ k.kategori }}</span>
            <button
              class="ak__toggle"
              :class="k.is_active ? 'ak__toggle--on' : 'ak__toggle--off'"
              @click="toggle(k)"
              :title="k.is_active ? 'Nonaktifkan' : 'Aktifkan'"
            >{{ k.is_active ? 'Aktif' : 'Nonaktif' }}</button>
            <button class="ak__icon-btn" @click="openForm(k)" title="Edit">
              <svg width="13" height="13" viewBox="0 0 24 24" fill="none"><path d="M11 4H4a2 2 0 00-2 2v14a2 2 0 002 2h14a2 2 0 002-2v-7" stroke="currentColor" stroke-width="2"/><path d="M18.5 2.5a2.121 2.121 0 013 3L12 15l-4 1 1-4 9.5-9.5z" stroke="currentColor" stroke-width="2"/></svg>
            </button>
            <button class="ak__icon-btn ak__icon-btn--del" @click="hapus(k)" title="Hapus">
              <svg width="13" height="13" viewBox="0 0 24 24" fill="none"><polyline points="3 6 5 6 21 6" stroke="currentColor" stroke-width="2"/><path d="M19 6l-1 14a2 2 0 01-2 2H8a2 2 0 01-2-2L5 6" stroke="currentColor" stroke-width="2"/><path d="M10 11v6M14 11v6" stroke="currentColor" stroke-width="2"/></svg>
            </button>
          </div>
        </div>
        <div class="ak__card-keywords">
          <span v-for="kw in k.kata_kunci" :key="kw" class="ak__kw">{{ kw }}</span>
          <span v-if="k.kata_kunci.length === 0" class="ak__kw-empty">Belum ada kata kunci</span>
        </div>
        <div class="ak__card-jawaban">{{ k.jawaban.substring(0, 120) }}{{ k.jawaban.length > 120 ? '...' : '' }}</div>
      </div>
    </div>

    <!-- Modal Form -->
    <div v-if="showForm" class="ak__overlay" @click.self="closeForm">
      <div class="ak__modal">
        <div class="ak__modal-header">
          <div class="ak__modal-title">{{ editItem ? 'Edit Knowledge' : 'Tambah Knowledge' }}</div>
          <button class="ak__modal-close" @click="closeForm">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none"><line x1="18" y1="6" x2="6" y2="18" stroke="currentColor" stroke-width="2" stroke-linecap="round"/><line x1="6" y1="6" x2="18" y2="18" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>
          </button>
        </div>

        <label class="ak__label">Contoh Pertanyaan <span class="ak__req">*</span></label>
        <input v-model="form.pertanyaan" class="ak__input" placeholder="Contoh: Cara mendaftar magang" />

        <label class="ak__label">Kata Kunci <span class="ak__req">*</span></label>
        <div class="ak__kw-input">
          <input v-model="kwInput" class="ak__input" placeholder="Ketik kata kunci lalu Enter..." @keydown.enter.prevent="addKw" />
          <button class="ak__kw-add-btn" @click="addKw">Tambah</button>
        </div>
        <div class="ak__kw-tags">
          <span v-for="(kw, i) in form.kata_kunci" :key="i" class="ak__kw-tag">
            {{ kw }}
            <button @click="removeKw(i)" class="ak__kw-tag-remove">&times;</button>
          </span>
          <span v-if="form.kata_kunci.length === 0" class="ak__kw-empty">Belum ada kata kunci</span>
        </div>

        <label class="ak__label">Jawaban <span class="ak__req">*</span></label>
        <textarea v-model="form.jawaban" class="ak__textarea" placeholder="Tulis jawaban lengkap..." rows="4"></textarea>

        <label class="ak__label">Kategori</label>
        <select v-model="form.kategori" class="ak__input">
          <option value="umum">Umum</option>
          <option value="absensi">Absensi</option>
          <option value="dokumen">Dokumen</option>
          <option value="sertifikat">Sertifikat</option>
        </select>

        <label class="ak__label-check">
          <input type="checkbox" v-model="form.is_active" class="ak__checkbox" />
          Aktif (digunakan oleh chatbot)
        </label>

        <div class="ak__modal-actions">
          <button class="ak__btn-cancel" @click="closeForm">Batal</button>
          <button
            class="ak__btn-save"
            @click="simpan"
            :disabled="saving || !form.pertanyaan || !form.jawaban || form.kata_kunci.length === 0"
          >{{ saving ? 'Menyimpan...' : 'Simpan' }}</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import axios from 'axios'

const token = localStorage.getItem('access_token') || ''
const headers = { Authorization: 'Bearer ' + token }

const list = ref([])
const loading = ref(false)
const saving = ref(false)
const search = ref('')
const filterKategori = ref('')
const filterStatus = ref('')

const showForm = ref(false)
const editItem = ref(null)
const kwInput = ref('')

const form = ref({
  pertanyaan: '',
  kata_kunci: [],
  jawaban: '',
  kategori: 'umum',
  is_active: true
})

const nlpQuery = ref('')
const nlpResult = ref(null)
const nlpLoading = ref(false)

const filtered = computed(() => {
  return list.value.filter(k => {
    const s = search.value.toLowerCase()
    const matchSearch = !s || k.pertanyaan.toLowerCase().includes(s) ||
      k.jawaban.toLowerCase().includes(s) ||
      k.kata_kunci.some(kw => kw.toLowerCase().includes(s))
    const matchKat = !filterKategori.value || k.kategori === filterKategori.value
    const matchStatus = !filterStatus.value ||
      (filterStatus.value === 'aktif' ? k.is_active : !k.is_active)
    return matchSearch && matchKat && matchStatus
  })
})

async function load() {
  loading.value = true
  try {
    const { data } = await axios.get('/api/admin/knowledge', { headers })
    list.value = data.data || []
  } finally {
    loading.value = false
  }
}

function openForm(item) {
  editItem.value = item || null
  if (item) {
    form.value = {
      pertanyaan: item.pertanyaan,
      kata_kunci: [...item.kata_kunci],
      jawaban: item.jawaban,
      kategori: item.kategori,
      is_active: item.is_active
    }
  } else {
    form.value = { pertanyaan: '', kata_kunci: [], jawaban: '', kategori: 'umum', is_active: true }
  }
  kwInput.value = ''
  showForm.value = true
}

function closeForm() {
  showForm.value = false
  editItem.value = null
}

function addKw() {
  const kw = kwInput.value.trim().toLowerCase()
  if (kw && !form.value.kata_kunci.includes(kw)) {
    form.value.kata_kunci.push(kw)
  }
  kwInput.value = ''
}

function removeKw(i) {
  form.value.kata_kunci.splice(i, 1)
}

async function simpan() {
  if (!form.value.pertanyaan || !form.value.jawaban || form.value.kata_kunci.length === 0) return
  saving.value = true
  try {
    const payload = {
      pertanyaan: form.value.pertanyaan,
      kata_kunci: form.value.kata_kunci,
      jawaban: form.value.jawaban,
      kategori: form.value.kategori,
      is_active: form.value.is_active
    }
    if (editItem.value) {
      await axios.put(`/api/admin/knowledge/${editItem.value.id}`, payload, { headers })
    } else {
      await axios.post('/api/admin/knowledge', payload, { headers })
    }
    closeForm()
    await load()
  } catch (e) {
    alert(e?.response?.data?.message || 'Gagal menyimpan')
  } finally {
    saving.value = false
  }
}

async function hapus(k) {
  if (!confirm(`Hapus knowledge "${k.pertanyaan}"?`)) return
  await axios.delete(`/api/admin/knowledge/${k.id}`, { headers })
  await load()
}

async function toggle(k) {
  await axios.patch(`/api/admin/knowledge/${k.id}/toggle`, {}, { headers })
  await load()
}

async function testNLP() {
  if (!nlpQuery.value.trim()) return
  nlpLoading.value = true
  nlpResult.value = null
  try {
    const { data } = await axios.post('/api/chatbot/query', { pesan: nlpQuery.value })
    nlpResult.value = data.data
  } finally {
    nlpLoading.value = false
  }
}

onMounted(load)
</script>

<style scoped>
.ak {
  padding: 0;
}

/* ── Top Bar ── */
.ak__topbar {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 12px;
  margin-bottom: 16px;
  flex-wrap: wrap;
}
.ak__topbar-left { flex: 1; min-width: 0; }
.ak__title { font-size: 17px; font-weight: 700; color: #111827; }
.ak__sub { font-size: 12px; color: #6b7280; margin-top: 3px; }
.ak__btn-add {
  display: inline-flex; align-items: center; gap: 6px;
  background: #16a34a; color: white; border: none;
  border-radius: 8px; padding: 9px 16px;
  font-size: 13px; font-weight: 600; cursor: pointer;
  white-space: nowrap; font-family: inherit;
  transition: background 0.15s;
  flex-shrink: 0;
}
.ak__btn-add:hover { background: #15803d; }

/* ── Filter ── */
.ak__filter {
  display: flex;
  gap: 8px;
  margin-bottom: 14px;
  flex-wrap: wrap;
  align-items: center;
}
.ak__search {
  flex: 1;
  min-width: 180px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  padding: 8px 12px;
  font-size: 13px;
  outline: none;
  font-family: inherit;
  transition: border-color 0.15s;
}
.ak__search:focus { border-color: #16a34a; }
.ak__selects { display: flex; gap: 8px; flex-wrap: wrap; }
.ak__select {
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  padding: 8px 10px;
  font-size: 13px;
  outline: none;
  background: white;
  cursor: pointer;
  font-family: inherit;
  color: #374151;
}

/* ── NLP Tester ── */
.ak__tester {
  background: #f9fafb;
  border: 1px solid #e5e7eb;
  border-radius: 10px;
  padding: 14px 16px;
  margin-bottom: 16px;
}
.ak__tester-label {
  font-size: 12px;
  font-weight: 600;
  color: #374151;
  margin-bottom: 10px;
  text-transform: uppercase;
  letter-spacing: 0.04em;
}
.ak__tester-row { display: flex; gap: 8px; }
.ak__tester-input {
  flex: 1;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  padding: 8px 12px;
  font-size: 13px;
  outline: none;
  background: white;
  font-family: inherit;
  transition: border-color 0.15s;
}
.ak__tester-input:focus { border-color: #16a34a; }
.ak__tester-btn {
  background: #16a34a;
  color: white;
  border: none;
  border-radius: 8px;
  padding: 8px 18px;
  font-size: 13px;
  cursor: pointer;
  font-weight: 600;
  font-family: inherit;
  transition: background 0.15s;
  white-space: nowrap;
}
.ak__tester-btn:disabled { opacity: 0.6; cursor: not-allowed; }
.ak__tester-btn:not(:disabled):hover { background: #15803d; }

.ak__tester-result {
  margin-top: 10px;
  padding: 10px 12px;
  border-radius: 8px;
  font-size: 12.5px;
}
.ak__tester-result--ok  { background: #f0fdf4; border: 1px solid #bbf7d0; }
.ak__tester-result--fail { background: #fef2f2; border: 1px solid #fecaca; }
.ak__tester-result-head {
  display: flex;
  align-items: center;
  gap: 7px;
  margin-bottom: 6px;
}
.ak__tester-status-dot {
  width: 8px; height: 8px;
  border-radius: 50%;
  flex-shrink: 0;
}
.dot--ok   { background: #16a34a; }
.dot--fail { background: #ef4444; }
.ak__tester-status-label { font-weight: 600; color: #111827; }
.ak__tester-skor {
  font-size: 11px;
  color: #6b7280;
  margin-left: auto;
}
.ak__tester-answer { color: #374151; line-height: 1.55; }
.ak__tester-answer--dim { color: #6b7280; }

/* ── List ── */
.ak__loading {
  display: flex; align-items: center; gap: 10px;
  padding: 32px; justify-content: center; color: #6b7280; font-size: 13px;
}
.ak__spinner {
  width: 20px; height: 20px;
  border: 2px solid #e5e7eb;
  border-top-color: #16a34a;
  border-radius: 50%;
  animation: spin 0.7s linear infinite;
}
@keyframes spin { to { transform: rotate(360deg); } }
.ak__empty {
  display: flex; flex-direction: column; align-items: center;
  gap: 10px; padding: 40px; text-align: center;
  color: #9ca3af; font-size: 13px;
}

.ak__list { display: flex; flex-direction: column; gap: 8px; }

/* ── Card ── */
.ak__card {
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 10px;
  padding: 14px 16px;
  transition: box-shadow 0.15s, border-color 0.15s;
}
.ak__card:hover { box-shadow: 0 2px 8px rgba(0,0,0,0.07); border-color: #d1d5db; }
.ak__card--inactive { opacity: 0.5; }

.ak__card-header {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  margin-bottom: 8px;
  flex-wrap: wrap;
}
.ak__card-pertanyaan {
  flex: 1;
  min-width: 0;
  font-weight: 600;
  font-size: 13.5px;
  color: #111827;
  line-height: 1.4;
}
.ak__card-actions {
  display: flex;
  align-items: center;
  gap: 5px;
  flex-shrink: 0;
  flex-wrap: wrap;
}

/* Badge kategori */
.ak__badge {
  font-size: 10px;
  padding: 2px 8px;
  border-radius: 10px;
  font-weight: 600;
  text-transform: capitalize;
  white-space: nowrap;
}
.ak__badge--umum      { background: #f1f5f9; color: #475569; }
.ak__badge--absensi   { background: #eff6ff; color: #1d4ed8; }
.ak__badge--dokumen   { background: #f5f3ff; color: #7c3aed; }
.ak__badge--sertifikat { background: #fef9c3; color: #a16207; }

/* Toggle aktif/nonaktif */
.ak__toggle {
  font-size: 10.5px;
  padding: 2px 8px;
  border-radius: 6px;
  border: none;
  cursor: pointer;
  font-weight: 600;
  font-family: inherit;
  white-space: nowrap;
  transition: opacity 0.15s;
}
.ak__toggle:hover { opacity: 0.8; }
.ak__toggle--on  { background: #dcfce7; color: #166534; }
.ak__toggle--off { background: #f1f5f9; color: #64748b; }

/* Icon buttons */
.ak__icon-btn {
  width: 28px; height: 28px;
  border-radius: 6px;
  border: 1px solid #e5e7eb;
  background: white;
  cursor: pointer;
  display: flex; align-items: center; justify-content: center;
  color: #6b7280;
  transition: all 0.15s;
}
.ak__icon-btn:hover { background: #f0fdf4; border-color: #16a34a; color: #16a34a; }
.ak__icon-btn--del:hover { background: #fef2f2; border-color: #ef4444; color: #ef4444; }

/* Card keywords */
.ak__card-keywords {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
  margin-bottom: 7px;
}
.ak__kw {
  background: #f1f5f9;
  color: #475569;
  font-size: 11px;
  padding: 2px 8px;
  border-radius: 6px;
}
.ak__kw-empty { color: #9ca3af; font-size: 11px; font-style: italic; }

.ak__card-jawaban {
  font-size: 12px;
  color: #6b7280;
  line-height: 1.55;
}

/* ── Modal Overlay ── */
.ak__overlay {
  position: fixed; inset: 0;
  background: rgba(0,0,0,0.4);
  z-index: 1000;
  display: flex; align-items: center; justify-content: center;
  padding: 16px;
}
.ak__modal {
  background: white;
  border-radius: 14px;
  padding: 24px;
  width: 100%;
  max-width: 520px;
  max-height: 90vh;
  overflow-y: auto;
}
.ak__modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 18px;
}
.ak__modal-title { font-size: 15px; font-weight: 700; color: #111827; }
.ak__modal-close {
  background: none; border: none; cursor: pointer;
  color: #6b7280; padding: 4px; border-radius: 6px;
  display: flex; align-items: center; justify-content: center;
  transition: background 0.15s;
}
.ak__modal-close:hover { background: #f3f4f6; color: #111827; }

/* Form elements */
.ak__label {
  display: block;
  font-size: 12px;
  font-weight: 600;
  color: #374151;
  margin-bottom: 5px;
  margin-top: 14px;
}
.ak__label:first-of-type { margin-top: 0; }
.ak__req { color: #ef4444; margin-left: 2px; }
.ak__input {
  width: 100%;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  padding: 8px 12px;
  font-size: 13px;
  outline: none;
  box-sizing: border-box;
  font-family: inherit;
  transition: border-color 0.15s;
}
.ak__input:focus { border-color: #16a34a; }
.ak__textarea {
  width: 100%;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  padding: 8px 12px;
  font-size: 13px;
  outline: none;
  resize: vertical;
  font-family: inherit;
  box-sizing: border-box;
  transition: border-color 0.15s;
  line-height: 1.55;
}
.ak__textarea:focus { border-color: #16a34a; }

/* Keyword input dalam form */
.ak__kw-input { display: flex; gap: 6px; }
.ak__kw-add-btn {
  background: white;
  border: 1px solid #e5e7eb;
  color: #374151;
  border-radius: 8px;
  padding: 8px 14px;
  font-size: 12.5px;
  cursor: pointer;
  font-family: inherit;
  white-space: nowrap;
  font-weight: 500;
  transition: all 0.15s;
}
.ak__kw-add-btn:hover { border-color: #16a34a; color: #16a34a; background: #f0fdf4; }

.ak__kw-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 5px;
  margin-top: 8px;
  min-height: 28px;
}
.ak__kw-tag {
  background: #f1f5f9;
  color: #374151;
  font-size: 12px;
  padding: 3px 8px 3px 10px;
  border-radius: 6px;
  display: flex; align-items: center; gap: 5px;
  border: 1px solid #e5e7eb;
}
.ak__kw-tag-remove {
  background: none; border: none;
  cursor: pointer; color: #9ca3af;
  font-size: 15px; line-height: 1; padding: 0;
  display: flex; align-items: center;
  transition: color 0.15s;
}
.ak__kw-tag-remove:hover { color: #ef4444; }

.ak__label-check {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: #374151;
  margin-top: 14px;
  cursor: pointer;
}
.ak__checkbox { width: 15px; height: 15px; cursor: pointer; accent-color: #16a34a; }

.ak__modal-actions {
  display: flex;
  gap: 8px;
  justify-content: flex-end;
  margin-top: 22px;
  padding-top: 16px;
  border-top: 1px solid #f3f4f6;
}
.ak__btn-cancel {
  padding: 8px 16px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  background: white;
  cursor: pointer;
  font-size: 13px;
  font-family: inherit;
  color: #374151;
  transition: background 0.15s;
}
.ak__btn-cancel:hover { background: #f3f4f6; }
.ak__btn-save {
  padding: 8px 22px;
  background: #16a34a;
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 13px;
  font-weight: 600;
  font-family: inherit;
  transition: background 0.15s;
}
.ak__btn-save:hover:not(:disabled) { background: #15803d; }
.ak__btn-save:disabled { opacity: 0.55; cursor: not-allowed; }

/* ── Mobile Responsive ── */
@media (max-width: 600px) {
  .ak__topbar { flex-direction: column; align-items: stretch; }
  .ak__btn-add { width: 100%; justify-content: center; }

  .ak__filter { flex-direction: column; }
  .ak__search { width: 100%; }
  .ak__selects { width: 100%; }
  .ak__select { flex: 1; }

  .ak__tester-row { flex-direction: column; }
  .ak__tester-btn { width: 100%; }

  .ak__card-header { flex-direction: column; gap: 8px; }
  .ak__card-actions { justify-content: flex-start; }

  .ak__modal { padding: 18px; }
  .ak__kw-input { flex-direction: column; }
  .ak__kw-add-btn { width: 100%; }

  .ak__modal-actions { flex-direction: column-reverse; }
  .ak__btn-save, .ak__btn-cancel { width: 100%; text-align: center; justify-content: center; }
}
</style>
