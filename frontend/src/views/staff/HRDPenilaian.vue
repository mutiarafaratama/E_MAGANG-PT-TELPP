<template>
  <div>
    <!-- Loading list -->
    <div v-if="loadingList" class="state-box"><div class="spinner"></div></div>

    <!-- List kosong -->
    <div v-else-if="list.length === 0" class="state-box state-box--empty">
      <svg width="44" height="44" viewBox="0 0 24 24" fill="none">
        <polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"
          stroke="#d1d5db" stroke-width="1.5"/>
      </svg>
      <div class="empty-title">Belum ada peserta yang perlu dinilai</div>
      <div class="empty-sub">Peserta akan muncul di sini setelah laporan magang disetujui.</div>
    </div>

    <!-- Tabel daftar peserta -->
    <div v-else-if="!selectedItem" class="card">
      <div class="card-header">
        <div>
          <h3 class="card-title">Penilaian Peserta</h3>
          <p class="card-sub">Input dan kelola nilai akhir magang dari seluruh peserta</p>
        </div>
      </div>
      <div class="table-wrap">
        <table class="data-table">
          <thead>
            <tr>
              <th>Peserta</th>
              <th>Divisi</th>
              <th>Periode</th>
              <th>Status</th>
              <th>Nilai Akhir</th>
              <th>Aksi</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="item in list" :key="item.pelaksanaan_id">
              <td>
                <div class="name-cell">
                  <div class="name-avatar">{{ item.nama_lengkap?.[0]?.toUpperCase() ?? '?' }}</div>
                  <div>
                    <div class="name-text">{{ item.nama_lengkap }}</div>
                    <div class="name-sub">{{ item.email }}</div>
                  </div>
                </div>
              </td>
              <td><span class="text-sm">{{ item.divisi || '—' }}</span></td>
              <td>
                <div class="text-sm">{{ fmtDate(item.tanggal_mulai) }}</div>
                <div class="text-xs text-gray">s/d {{ fmtDate(item.tanggal_selesai) }}</div>
              </td>
              <td>
                <span v-if="item.nilai_akhir != null" class="chip chip--green">Sudah dinilai</span>
                <span v-else class="chip chip--yellow">Belum dinilai</span>
              </td>
              <td>
                <span v-if="item.nilai_akhir != null" class="nilai-badge" :style="{ color: nilaiColor(item.nilai_akhir) }">
                  {{ item.nilai_akhir.toFixed(2) }}
                  <span class="grade-small">{{ nilaiGrade(item.nilai_akhir) }}</span>
                </span>
                <span v-else class="text-gray text-sm">—</span>
              </td>
              <td>
                <button class="btn-primary-sm" @click="openForm(item)">
                  {{ item.nilai_akhir != null ? 'Edit Nilai' : 'Input Nilai' }}
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- ═══ FORM PENILAIAN ═══ -->
    <div v-else class="form-wrap">
      <button class="btn-back" @click="closeForm">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none">
          <path d="M19 12H5M12 5l-7 7 7 7" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
        </svg>
        Kembali ke daftar
      </button>

      <!-- Info peserta -->
      <div class="peserta-info-card">
        <div class="peserta-avatar">{{ selectedItem.nama_lengkap?.[0]?.toUpperCase() }}</div>
        <div class="peserta-info-body">
          <div class="peserta-name">{{ selectedItem.nama_lengkap }}</div>
          <div class="peserta-meta">{{ selectedItem.divisi || 'Belum ada divisi' }} · {{ fmtDate(selectedItem.tanggal_mulai) }} – {{ fmtDate(selectedItem.tanggal_selesai) }}</div>
          <div class="peserta-meta">Pembimbing: <strong>{{ selectedItem.pembimbing || '—' }}</strong></div>
        </div>
        <div v-if="nilaiAkhir > 0" class="peserta-preview-score" :style="{ color: nilaiColor(nilaiAkhir) }">
          <div class="preview-score-num">{{ nilaiAkhir.toFixed(2) }}</div>
          <div class="preview-score-label">{{ nilaiGradeLabel(nilaiAkhir) }}</div>
        </div>
      </div>

      <!-- Import Excel Banner -->
      <div class="excel-banner">
        <div class="excel-banner__left">
          <div class="excel-icon">
            <svg width="22" height="22" viewBox="0 0 24 24" fill="none">
              <rect x="3" y="3" width="18" height="18" rx="3" stroke="#16a34a" stroke-width="1.5"/>
              <path d="M8 8l3 4-3 4M12 16h4" stroke="#16a34a" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
          </div>
          <div>
            <div class="excel-banner__title">Import dari Excel</div>
            <div class="excel-banner__sub">Upload file .xlsx template daftar penilaian — Bagian I, III, IV otomatis terisi</div>
          </div>
        </div>
        <label class="btn-excel" :class="{ 'btn-excel--loading': parsigExcel }">
          <span v-if="parsigExcel" class="spinner-sm"></span>
          <svg v-else width="15" height="15" viewBox="0 0 24 24" fill="none">
            <path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4M17 8l-5-5-5 5M12 3v12" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
          {{ parsigExcel ? 'Memproses...' : 'Pilih File Excel' }}
          <input type="file" accept=".xlsx,.xls" @change="onExcelUpload" :disabled="parsigExcel" style="display:none"/>
        </label>
      </div>
      <div v-if="excelMsg" class="excel-msg" :class="excelMsg.type === 'ok' ? 'excel-msg--ok' : 'excel-msg--err'">
        <svg v-if="excelMsg.type === 'ok'" width="14" height="14" viewBox="0 0 24 24" fill="none">
          <path d="M20 6L9 17l-5-5" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
        </svg>
        <svg v-else width="14" height="14" viewBox="0 0 24 24" fill="none">
          <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="1.5"/>
          <path d="M12 8v4M12 16h.01" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
        </svg>
        {{ excelMsg.text }}
      </div>

      <div v-if="saveError" class="alert alert--error">{{ saveError }}</div>
      <div v-if="saveSuccess" class="alert alert--success">{{ saveSuccess }}</div>

      <!-- ── I. KEPRIBADIAN / ETOS KERJA ── -->
      <div class="section-card">
        <div class="section-header">
          <div class="section-num">I</div>
          <div class="section-header__body">
            <div class="section-title">Kepribadian / Etos Kerja</div>
            <div class="section-sub">9 komponen — diisi dari Excel atau manual</div>
          </div>
          <div class="section-avg" v-if="avgI > 0">Rata-rata: <strong :style="{ color: nilaiColor(avgI) }">{{ avgI.toFixed(2) }}</strong></div>
        </div>
        <div class="komponen-grid">
          <div class="komponen-row" v-for="k in bagianI" :key="k.key">
            <div class="komponen-info">
              <div class="komponen-num">{{ k.no }}</div>
              <div class="komponen-name">{{ k.label }}</div>
            </div>
            <div class="komponen-input-wrap">
              <input type="number" class="nilai-input" min="0" max="100" step="1"
                v-model.number="(form as any)[k.key]"
                @input="recalc" placeholder="0–100"/>
              <div class="nilai-bar-wrap">
                <div class="nilai-bar" :style="{ width: ((form as any)[k.key] || 0) + '%', background: nilaiColor((form as any)[k.key] || 0) }"></div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- ── II. KEMAMPUAN KEJURUAN (Manual) ── -->
      <div class="section-card">
        <div class="section-header">
          <div class="section-num">II</div>
          <div class="section-header__body">
            <div class="section-title">Kemampuan Kejuruan</div>
            <div class="section-sub">Komponen disesuaikan per divisi — input manual HRD</div>
          </div>
          <div class="section-avg" v-if="avgII > 0">Rata-rata: <strong :style="{ color: nilaiColor(avgII) }">{{ avgII.toFixed(2) }}</strong></div>
        </div>
        <div class="komponen-grid" v-if="form.kejuruan.length > 0">
          <div class="komponen-row" v-for="(k, i) in form.kejuruan" :key="i">
            <div class="kejuruan-no">{{ i + 1 }}</div>
            <div class="kejuruan-name-wrap">
              <input type="text" class="inp-komponen" v-model="k.komponen" placeholder="Nama komponen, mis: Pemrograman Web"/>
            </div>
            <div class="komponen-input-wrap">
              <input type="number" class="nilai-input" min="0" max="100" step="1"
                v-model.number="k.nilai" @input="recalc" placeholder="0–100"/>
              <div class="nilai-bar-wrap">
                <div class="nilai-bar" :style="{ width: (k.nilai || 0) + '%', background: nilaiColor(k.nilai || 0) }"></div>
              </div>
            </div>
            <button class="btn-del-kejuruan" @click="removeKejuruan(i)" title="Hapus">
              <svg width="13" height="13" viewBox="0 0 24 24" fill="none">
                <path d="M18 6L6 18M6 6l12 12" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
              </svg>
            </button>
          </div>
        </div>
        <div v-else class="kejuruan-empty">
          <svg width="30" height="30" viewBox="0 0 24 24" fill="none"><circle cx="12" cy="12" r="9" stroke="#d1d5db" stroke-width="1.5"/><path d="M12 8v4M12 16h.01" stroke="#d1d5db" stroke-width="2" stroke-linecap="round"/></svg>
          <span>Belum ada komponen kejuruan. Tambahkan sesuai bidang peserta.</span>
        </div>
        <button class="btn-add-kejuruan" @click="addKejuruan">
          <svg width="13" height="13" viewBox="0 0 24 24" fill="none">
            <path d="M12 5v14M5 12h14" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
          </svg>
          Tambah Komponen Kejuruan
        </button>
      </div>

      <!-- ── III. K3 DAN PRODUKTIVITAS ── -->
      <div class="section-card">
        <div class="section-header">
          <div class="section-num">III</div>
          <div class="section-header__body">
            <div class="section-title">K3 dan Produktivitas</div>
            <div class="section-sub">4 komponen — diisi dari Excel atau manual</div>
          </div>
          <div class="section-avg" v-if="avgIII > 0">Rata-rata: <strong :style="{ color: nilaiColor(avgIII) }">{{ avgIII.toFixed(2) }}</strong></div>
        </div>
        <div class="komponen-grid">
          <div class="komponen-row" v-for="k in bagianIII" :key="k.key">
            <div class="komponen-info">
              <div class="komponen-num">{{ k.no }}</div>
              <div class="komponen-name">{{ k.label }}</div>
            </div>
            <div class="komponen-input-wrap">
              <input type="number" class="nilai-input" min="0" max="100" step="1"
                v-model.number="(form as any)[k.key]"
                @input="recalc" placeholder="0–100"/>
              <div class="nilai-bar-wrap">
                <div class="nilai-bar" :style="{ width: ((form as any)[k.key] || 0) + '%', background: nilaiColor((form as any)[k.key] || 0) }"></div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- ── IV. PRESENTASI ── -->
      <div class="section-card">
        <div class="section-header">
          <div class="section-num">IV</div>
          <div class="section-header__body">
            <div class="section-title">Presentasi</div>
            <div class="section-sub">4 komponen — khusus peserta perguruan tinggi (D1/D3/S1)</div>
          </div>
          <div class="section-avg" v-if="avgIV > 0">Rata-rata: <strong :style="{ color: nilaiColor(avgIV) }">{{ avgIV.toFixed(2) }}</strong></div>
        </div>
        <div class="komponen-grid">
          <div class="komponen-row" v-for="k in bagianIV" :key="k.key">
            <div class="komponen-info">
              <div class="komponen-num">{{ k.no }}</div>
              <div class="komponen-name">{{ k.label }}</div>
            </div>
            <div class="komponen-input-wrap">
              <input type="number" class="nilai-input" min="0" max="100" step="1"
                v-model.number="(form as any)[k.key]"
                @input="recalc" placeholder="0–100"/>
              <div class="nilai-bar-wrap">
                <div class="nilai-bar" :style="{ width: ((form as any)[k.key] || 0) + '%', background: nilaiColor((form as any)[k.key] || 0) }"></div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- ── PENANDATANGAN ── -->
      <div class="section-card">
        <div class="section-header" style="margin-bottom:14px">
          <div class="section-num">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none"><path d="M12 20h9M16.5 3.5a2.121 2.121 0 013 3L7 19l-4 1 1-4L16.5 3.5z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>
          </div>
          <div class="section-header__body">
            <div class="section-title">Penandatangan</div>
            <div class="section-sub">Pembimbing otomatis dari data peserta, Manager diisi manual</div>
          </div>
        </div>
        <div class="penandatangan-grid">
          <div class="penandatangan-item">
            <label class="field-label">Pembimbing Lapangan</label>
            <input type="text" class="inp-readonly" :value="selectedItem.pembimbing || '—'" readonly/>
            <div class="field-hint">Otomatis dari data peserta magang</div>
          </div>
          <div class="penandatangan-item">
            <label class="field-label">Manager (Dept. terkait) <span class="req">*</span></label>
            <input type="text" class="inp-field" v-model="form.manager_nama" placeholder="Nama Manager departemen terkait"/>
          </div>
        </div>
      </div>

      <!-- ── CATATAN ── -->
      <div class="section-card">
        <div class="section-title" style="margin-bottom:10px">Catatan Pembimbing</div>
        <textarea class="textarea-catatan" v-model="form.catatan" rows="3"
          placeholder="Tuliskan catatan, saran, atau komentar untuk peserta…"></textarea>
      </div>

      <!-- ── KETERANGAN NILAI ── -->
      <div class="keterangan-card">
        <div class="keterangan-title">Keterangan Nilai</div>
        <div class="keterangan-grid">
          <div class="keterangan-item keterangan-item--a"><span class="k-grade">A</span><span class="k-label">Baik Sekali</span><span class="k-range">85 – 100</span></div>
          <div class="keterangan-item keterangan-item--b"><span class="k-grade">B</span><span class="k-label">Baik</span><span class="k-range">75 – 84</span></div>
          <div class="keterangan-item keterangan-item--c"><span class="k-grade">C</span><span class="k-label">Cukup</span><span class="k-range">60 – 74</span></div>
          <div class="keterangan-item keterangan-item--d"><span class="k-grade">D</span><span class="k-label">Kurang</span><span class="k-range">&lt; 60</span></div>
        </div>
      </div>

      <!-- ── RINGKASAN & SIMPAN ── -->
      <div class="summary-card">
        <div class="summary-title">Ringkasan Nilai</div>
        <div class="summary-rows">
          <div class="summary-row">
            <span>I. Kepribadian / Etos Kerja (9 komponen)</span>
            <strong :style="{ color: avgI > 0 ? nilaiColor(avgI) : '#9ca3af' }">{{ avgI > 0 ? avgI.toFixed(2) : '—' }}</strong>
          </div>
          <div class="summary-row">
            <span>II. Kemampuan Kejuruan ({{ form.kejuruan.filter(k=>k.komponen).length }} komponen)</span>
            <strong :style="{ color: avgII > 0 ? nilaiColor(avgII) : '#9ca3af' }">{{ avgII > 0 ? avgII.toFixed(2) : '—' }}</strong>
          </div>
          <div class="summary-row">
            <span>III. K3 dan Produktivitas (4 komponen)</span>
            <strong :style="{ color: avgIII > 0 ? nilaiColor(avgIII) : '#9ca3af' }">{{ avgIII > 0 ? avgIII.toFixed(2) : '—' }}</strong>
          </div>
          <div class="summary-row">
            <span>IV. Presentasi (4 komponen)</span>
            <strong :style="{ color: avgIV > 0 ? nilaiColor(avgIV) : '#9ca3af' }">{{ avgIV > 0 ? avgIV.toFixed(2) : '—' }}</strong>
          </div>
          <div class="summary-divider"></div>
          <div class="summary-row summary-row--total">
            <span>Nilai Akhir <span class="total-formula">(Rata-rata {{ totalItems }} item)</span></span>
            <span class="summary-final" :style="{ color: nilaiAkhir > 0 ? nilaiColor(nilaiAkhir) : '#9ca3af' }">
              {{ nilaiAkhir > 0 ? nilaiAkhir.toFixed(2) : '—' }}
              <span v-if="nilaiAkhir > 0" class="summary-grade">{{ nilaiGradeLabel(nilaiAkhir) }}</span>
            </span>
          </div>
        </div>
        <button class="btn-simpan" :disabled="saving || nilaiAkhir <= 0" @click="simpan">
          <span v-if="saving" class="spinner-sm"></span>
          <svg v-else width="15" height="15" viewBox="0 0 24 24" fill="none">
            <path d="M19 21H5a2 2 0 01-2-2V5a2 2 0 012-2h11l5 5v11a2 2 0 01-2 2z" stroke="currentColor" stroke-width="2" stroke-linejoin="round"/>
            <polyline points="17 21 17 13 7 13 7 21" stroke="currentColor" stroke-width="2" stroke-linejoin="round"/>
            <polyline points="7 3 7 8 15 8" stroke="currentColor" stroke-width="2" stroke-linejoin="round"/>
          </svg>
          Simpan Penilaian
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, reactive } from "vue";
import * as XLSX from "xlsx";
import api from "@/lib/api";

interface ListItem {
  pelaksanaan_id: string;
  status: string;
  nama_lengkap: string;
  email: string;
  divisi: string;
  pembimbing: string;
  tanggal_mulai: string;
  tanggal_selesai: string;
  penilaian_id: string | null;
  nilai_akhir: number | null;
  dinilai_at: string | null;
}

interface KejuruanItem { komponen: string; nilai: number }

interface FormState {
  // I. Kepribadian
  nilai_motivasi: number;
  nilai_inisiatif: number;
  nilai_disiplin_waktu: number;
  nilai_kerajinan: number;
  nilai_kreativitas: number;
  nilai_tanggung_jawab: number;
  nilai_kerjasama: number;
  nilai_adaptasi: number;
  nilai_kehadiran: number;
  // II. Kejuruan
  kejuruan: KejuruanItem[];
  // III. K3
  nilai_k3_safety: number;
  nilai_k3_metode: number;
  nilai_k3_manajemen: number;
  nilai_k3_volume: number;
  // IV. Presentasi
  nilai_prs_proses: number;
  nilai_prs_teori: number;
  nilai_prs_judul: number;
  nilai_prs_data: number;
  // Penandatangan & catatan
  manager_nama: string;
  catatan: string;
}

const loadingList  = ref(false);
const saving       = ref(false);
const parsigExcel  = ref(false);
const saveError    = ref("");
const saveSuccess  = ref("");
const excelMsg     = ref<{ type: "ok" | "err"; text: string } | null>(null);
const list         = ref<ListItem[]>([]);
const selectedItem = ref<ListItem | null>(null);

const form = reactive<FormState>({
  nilai_motivasi: 0, nilai_inisiatif: 0, nilai_disiplin_waktu: 0,
  nilai_kerajinan: 0, nilai_kreativitas: 0, nilai_tanggung_jawab: 0,
  nilai_kerjasama: 0, nilai_adaptasi: 0, nilai_kehadiran: 0,
  kejuruan: [],
  nilai_k3_safety: 0, nilai_k3_metode: 0, nilai_k3_manajemen: 0, nilai_k3_volume: 0,
  nilai_prs_proses: 0, nilai_prs_teori: 0, nilai_prs_judul: 0, nilai_prs_data: 0,
  manager_nama: "", catatan: "",
});

// ── Definisi komponen per bagian ──
const bagianI = [
  { no: 1, key: "nilai_motivasi",       label: "Motivasi" },
  { no: 2, key: "nilai_inisiatif",      label: "Inisiatif" },
  { no: 3, key: "nilai_disiplin_waktu", label: "Disiplin Waktu" },
  { no: 4, key: "nilai_kerajinan",      label: "Kerajinan" },
  { no: 5, key: "nilai_kreativitas",    label: "Kreativitas" },
  { no: 6, key: "nilai_tanggung_jawab", label: "Tanggung Jawab" },
  { no: 7, key: "nilai_kerjasama",      label: "Kerjasama" },
  { no: 8, key: "nilai_adaptasi",       label: "Adaptasi dengan Lingkungan Kerja" },
  { no: 9, key: "nilai_kehadiran",      label: "Kehadiran" },
];

const bagianIII = [
  { no: 1, key: "nilai_k3_safety",    label: "Penggunaan dan Penerapan K3 (Safety Equipment)" },
  { no: 2, key: "nilai_k3_metode",    label: "Metode Kerja dalam Menggunakan Alat dan Bahan" },
  { no: 3, key: "nilai_k3_manajemen", label: "Manajemen Kerja" },
  { no: 4, key: "nilai_k3_volume",    label: "Volume Pekerjaan yang Diselesaikan" },
];

const bagianIV = [
  { no: 1, key: "nilai_prs_proses", label: "Pemahaman Proses Pulp PT TELPP" },
  { no: 2, key: "nilai_prs_teori",  label: "Pemahaman Teori" },
  { no: 3, key: "nilai_prs_judul",  label: "Pemahaman Judul Praktek Kerja Lapangan" },
  { no: 4, key: "nilai_prs_data",   label: "Kemampuan Mengelolah Data Dengan Baik dan Benar" },
];

// ── Rata-rata per bagian ──
const avgI = computed(() => {
  const vals = bagianI.map(k => (form as any)[k.key] as number).filter(v => v > 0);
  return vals.length === bagianI.length ? vals.reduce((a, b) => a + b, 0) / vals.length : 0;
});
const avgII = computed(() => {
  const vals = form.kejuruan.filter(k => k.komponen && k.nilai > 0).map(k => k.nilai);
  return vals.length > 0 ? vals.reduce((a, b) => a + b, 0) / vals.length : 0;
});
const avgIII = computed(() => {
  const vals = bagianIII.map(k => (form as any)[k.key] as number).filter(v => v > 0);
  return vals.length === bagianIII.length ? vals.reduce((a, b) => a + b, 0) / vals.length : 0;
});
const avgIV = computed(() => {
  const vals = bagianIV.map(k => (form as any)[k.key] as number).filter(v => v > 0);
  return vals.length === bagianIV.length ? vals.reduce((a, b) => a + b, 0) / vals.length : 0;
});

// Total item untuk nilai akhir: 9 kepribadian + n kejuruan + 4 K3 + 4 presentasi
const totalItems = computed(() => {
  const kjCount = form.kejuruan.filter(k => k.komponen && k.nilai > 0).length;
  // Hitung hanya bagian yang sudah terisi semua
  let count = 0;
  if (avgI.value > 0) count += 9;
  count += kjCount;
  if (avgIII.value > 0) count += 4;
  if (avgIV.value > 0) count += 4;
  return count;
});

// Nilai akhir = jumlah semua nilai terisi / total item (sama dengan rumus Excel)
const nilaiAkhir = computed(() => {
  let total = 0;
  let count = 0;
  if (avgI.value > 0) {
    total += bagianI.reduce((s, k) => s + ((form as any)[k.key] as number), 0);
    count += 9;
  }
  const kjFilled = form.kejuruan.filter(k => k.komponen && k.nilai > 0);
  if (kjFilled.length > 0) {
    total += kjFilled.reduce((s, k) => s + k.nilai, 0);
    count += kjFilled.length;
  }
  if (avgIII.value > 0) {
    total += bagianIII.reduce((s, k) => s + ((form as any)[k.key] as number), 0);
    count += 4;
  }
  if (avgIV.value > 0) {
    total += bagianIV.reduce((s, k) => s + ((form as any)[k.key] as number), 0);
    count += 4;
  }
  if (count === 0) return 0;
  return Math.round((total / count) * 100) / 100;
});

function recalc() { /* nilai akhir dihitung via computed */ }

function addKejuruan() {
  form.kejuruan.push({ komponen: "", nilai: 0 });
}
function removeKejuruan(idx: number) {
  form.kejuruan.splice(idx, 1);
}

// ── Excel upload & parse (SheetJS — langsung di browser, tanpa kirim ke server) ──
function cellNum(sheet: XLSX.WorkSheet, addr: string): number {
  const cell = sheet[addr];
  if (!cell) return 0;
  // Gunakan cached value (.v) — works untuk formula maupun angka biasa
  const v = Number(cell.v);
  return isFinite(v) && v > 0 ? v : 0;
}
function cellStr(sheet: XLSX.WorkSheet, addr: string): string {
  const cell = sheet[addr];
  return cell ? String(cell.v ?? "").trim() : "";
}

function onExcelUpload(e: Event) {
  const input = e.target as HTMLInputElement;
  const file = input.files?.[0];
  if (!file) return;
  input.value = "";

  parsigExcel.value = true;
  excelMsg.value = null;

  const reader = new FileReader();
  reader.onload = (ev) => {
    try {
      const data = ev.target?.result;
      if (!data) throw new Error("File kosong");

      const wb = XLSX.read(data, { type: "array", cellFormula: false, cellNF: false });
      if (!wb.SheetNames.length) throw new Error("File tidak memiliki sheet");

      const ws = wb.Sheets[wb.SheetNames[0]];

      // I. KEPRIBADIAN / ETOS KERJA — kolom F, baris 16–24
      form.nilai_motivasi       = cellNum(ws, "F16");
      form.nilai_inisiatif      = cellNum(ws, "F17");
      form.nilai_disiplin_waktu = cellNum(ws, "F18");
      form.nilai_kerajinan      = cellNum(ws, "F19");
      form.nilai_kreativitas    = cellNum(ws, "F20");
      form.nilai_tanggung_jawab = cellNum(ws, "F21");
      form.nilai_kerjasama      = cellNum(ws, "F22");
      form.nilai_adaptasi       = cellNum(ws, "F23");
      form.nilai_kehadiran      = cellNum(ws, "F24");

      // II. KEMAMPUAN KEJURUAN — kolom D (nama) & F (nilai), baris 27–31
      const kejuruanBaru: KejuruanItem[] = [];
      for (let row = 27; row <= 31; row++) {
        const nama  = cellStr(ws, `D${row}`);
        const nilai = cellNum(ws, `F${row}`);
        if (nama && nilai > 0) kejuruanBaru.push({ komponen: nama, nilai });
      }
      if (kejuruanBaru.length > 0) form.kejuruan = kejuruanBaru;

      // III. K3 DAN PRODUKTIVITAS — kolom F, baris 34–37
      form.nilai_k3_safety    = cellNum(ws, "F34");
      form.nilai_k3_metode    = cellNum(ws, "F35");
      form.nilai_k3_manajemen = cellNum(ws, "F36");
      form.nilai_k3_volume    = cellNum(ws, "F37");

      // IV. PRESENTASI — kolom F, baris 40–43
      form.nilai_prs_proses = cellNum(ws, "F40");
      form.nilai_prs_teori  = cellNum(ws, "F41");
      form.nilai_prs_judul  = cellNum(ws, "F42");
      form.nilai_prs_data   = cellNum(ws, "F43");

      // Hitung preview nilai akhir dari data Excel
      const terisi = [
        form.nilai_motivasi, form.nilai_inisiatif, form.nilai_disiplin_waktu,
        form.nilai_kerajinan, form.nilai_kreativitas, form.nilai_tanggung_jawab,
        form.nilai_kerjasama, form.nilai_adaptasi, form.nilai_kehadiran,
        form.nilai_k3_safety, form.nilai_k3_metode, form.nilai_k3_manajemen, form.nilai_k3_volume,
        form.nilai_prs_proses, form.nilai_prs_teori, form.nilai_prs_judul, form.nilai_prs_data,
        ...form.kejuruan.map(k => k.nilai),
      ].filter(v => v > 0);

      const preview = terisi.length > 0
        ? (terisi.reduce((a, b) => a + b, 0) / terisi.length).toFixed(2)
        : "—";

      excelMsg.value = {
        type: "ok",
        text: `Excel berhasil dibaca (${terisi.length} item nilai). Preview rata-rata: ${preview}. Lengkapi Bagian II (Kejuruan) & nama Manager, lalu simpan.`,
      };
    } catch (err: any) {
      excelMsg.value = {
        type: "err",
        text: "Gagal membaca file Excel: " + (err?.message ?? "format tidak dikenali"),
      };
    } finally {
      parsigExcel.value = false;
    }
  };
  reader.onerror = () => {
    excelMsg.value = { type: "err", text: "Gagal membuka file. Coba lagi." };
    parsigExcel.value = false;
  };
  reader.readAsArrayBuffer(file);
}

function nilaiColor(n: number) {
  if (n >= 85) return "#16a34a";
  if (n >= 75) return "#1a5c20";
  if (n >= 60) return "#d97706";
  return "#dc2626";
}
function nilaiGrade(n: number) {
  if (n >= 85) return "A";
  if (n >= 75) return "B";
  if (n >= 60) return "C";
  return "D";
}
function nilaiGradeLabel(n: number) {
  if (n >= 85) return "A — Baik Sekali";
  if (n >= 75) return "B — Baik";
  if (n >= 60) return "C — Cukup";
  return "D — Kurang";
}
function fmtDate(d: string) {
  return new Date(d).toLocaleDateString("id-ID", { day: "2-digit", month: "short", year: "numeric" });
}

function resetForm() {
  Object.assign(form, {
    nilai_motivasi: 0, nilai_inisiatif: 0, nilai_disiplin_waktu: 0,
    nilai_kerajinan: 0, nilai_kreativitas: 0, nilai_tanggung_jawab: 0,
    nilai_kerjasama: 0, nilai_adaptasi: 0, nilai_kehadiran: 0,
    kejuruan: [],
    nilai_k3_safety: 0, nilai_k3_metode: 0, nilai_k3_manajemen: 0, nilai_k3_volume: 0,
    nilai_prs_proses: 0, nilai_prs_teori: 0, nilai_prs_judul: 0, nilai_prs_data: 0,
    manager_nama: "", catatan: "",
  });
}

async function openForm(item: ListItem) {
  selectedItem.value = item;
  saveError.value = ""; saveSuccess.value = ""; excelMsg.value = null;
  resetForm();
  if (item.penilaian_id) {
    try {
      const res = await api.get(`/api/penilaian/${item.pelaksanaan_id}`);
      const p = res.data;
      if (p) {
        form.nilai_motivasi       = p.nilai_motivasi       ?? 0;
        form.nilai_inisiatif      = p.nilai_inisiatif      ?? 0;
        form.nilai_disiplin_waktu = p.nilai_disiplin_waktu ?? 0;
        form.nilai_kerajinan      = p.nilai_kerajinan      ?? 0;
        form.nilai_kreativitas    = p.nilai_kreativitas    ?? 0;
        form.nilai_tanggung_jawab = p.nilai_tanggung_jawab ?? 0;
        form.nilai_kerjasama      = p.nilai_kerjasama      ?? 0;
        form.nilai_adaptasi       = p.nilai_adaptasi       ?? 0;
        form.nilai_kehadiran      = p.nilai_kehadiran      ?? 0;
        form.nilai_k3_safety      = p.nilai_k3_safety      ?? 0;
        form.nilai_k3_metode      = p.nilai_k3_metode      ?? 0;
        form.nilai_k3_manajemen   = p.nilai_k3_manajemen   ?? 0;
        form.nilai_k3_volume      = p.nilai_k3_volume      ?? 0;
        form.nilai_prs_proses     = p.nilai_prs_proses     ?? 0;
        form.nilai_prs_teori      = p.nilai_prs_teori      ?? 0;
        form.nilai_prs_judul      = p.nilai_prs_judul      ?? 0;
        form.nilai_prs_data       = p.nilai_prs_data       ?? 0;
        form.manager_nama         = p.manager_nama         ?? "";
        form.catatan              = p.catatan              ?? "";
        form.kejuruan = (p.kejuruan ?? []).map((k: any) => ({ komponen: k.komponen, nilai: k.nilai }));
      }
    } catch { /* biarkan form kosong */ }
  }
}

function closeForm() {
  selectedItem.value = null; saveError.value = ""; saveSuccess.value = ""; excelMsg.value = null;
}

async function simpan() {
  if (!selectedItem.value) return;
  saving.value = true; saveError.value = ""; saveSuccess.value = "";
  try {
    const payload = {
      pelaksanaan_id:       selectedItem.value.pelaksanaan_id,
      nilai_motivasi:       form.nilai_motivasi,
      nilai_inisiatif:      form.nilai_inisiatif,
      nilai_disiplin_waktu: form.nilai_disiplin_waktu,
      nilai_kerajinan:      form.nilai_kerajinan,
      nilai_kreativitas:    form.nilai_kreativitas,
      nilai_tanggung_jawab: form.nilai_tanggung_jawab,
      nilai_kerjasama:      form.nilai_kerjasama,
      nilai_adaptasi:       form.nilai_adaptasi,
      nilai_kehadiran:      form.nilai_kehadiran,
      kejuruan:             form.kejuruan.filter(k => k.komponen.trim() !== ""),
      nilai_k3_safety:      form.nilai_k3_safety,
      nilai_k3_metode:      form.nilai_k3_metode,
      nilai_k3_manajemen:   form.nilai_k3_manajemen,
      nilai_k3_volume:      form.nilai_k3_volume,
      nilai_prs_proses:     form.nilai_prs_proses,
      nilai_prs_teori:      form.nilai_prs_teori,
      nilai_prs_judul:      form.nilai_prs_judul,
      nilai_prs_data:       form.nilai_prs_data,
      manager_nama:         form.manager_nama,
      catatan:              form.catatan,
      nilai_akhir:          nilaiAkhir.value,
    };
    await api.post(`/api/penilaian/${selectedItem.value.pelaksanaan_id}`, payload);
    saveSuccess.value = "Penilaian berhasil disimpan!";
    const idx = list.value.findIndex(i => i.pelaksanaan_id === selectedItem.value!.pelaksanaan_id);
    if (idx >= 0) {
      list.value[idx].nilai_akhir = nilaiAkhir.value;
      list.value[idx].penilaian_id = "updated";
    }
    selectedItem.value = { ...selectedItem.value, nilai_akhir: nilaiAkhir.value, penilaian_id: "updated" };
  } catch (e: any) {
    saveError.value = e?.response?.data?.message ?? "Gagal menyimpan penilaian.";
  } finally {
    saving.value = false;
  }
}

onMounted(async () => {
  loadingList.value = true;
  try {
    const res = await api.get("/api/penilaian");
    list.value = res.data ?? [];
  } catch { list.value = []; }
  finally { loadingList.value = false; }
});
</script>

<style scoped>
.page-header { display:flex; align-items:center; margin-bottom:20px; gap:12px; }
.page-header__left { display:flex; align-items:center; gap:14px; }
.page-icon { width:48px; height:48px; border-radius:12px; background:#f0fdf4; color:#0d2818; border:1px solid #e9d5ff; display:flex; align-items:center; justify-content:center; flex-shrink:0; }
.page-title { font-size:18px; font-weight:700; color:#111827; margin:0 0 3px; }
.page-sub   { font-size:12.5px; color:#6b7280; margin:0; }

.state-box { display:flex; flex-direction:column; align-items:center; justify-content:center; gap:12px; padding:56px 24px; background:#fff; border-radius:14px; border:1px solid #e9f5e9; text-align:center; }
.spinner { width:28px; height:28px; border:3px solid #e5e7eb; border-top-color:#48AF4A; border-radius:50%; animation:spin .7s linear infinite; }
@keyframes spin { to { transform:rotate(360deg); } }
.empty-title { font-size:15px; font-weight:700; color:#374151; }
.empty-sub   { font-size:12.5px; color:#9ca3af; max-width:360px; line-height:1.6; }

/* Table */
.card { background:#fff; border:1px solid #e9f5e9; border-radius:14px; overflow:hidden; }
.card-header { padding:16px 20px; border-bottom:1px solid #f0faf0; }
.card-title  { font-size:13.5px; font-weight:700; color:#111827; margin:0 0 2px; }
.card-sub    { font-size:11.5px; color:#9ca3af; margin:0; }
.table-wrap { overflow-x:auto; }
.data-table { width:100%; border-collapse:collapse; }
.data-table th { font-size:11px; font-weight:700; color:#6b7280; text-transform:uppercase; letter-spacing:.04em; padding:10px 14px; background:#f9fafb; border-bottom:1px solid #f3f4f6; white-space:nowrap; }
.data-table td { padding:12px 14px; border-bottom:1px solid #f9fafb; font-size:13px; color:#374151; vertical-align:middle; }
.data-table tr:last-child td { border-bottom:none; }
.name-cell { display:flex; align-items:center; gap:10px; }
.name-avatar { width:34px; height:34px; border-radius:50%; background:#f0fdf4; color:#16a34a; font-size:13px; font-weight:700; display:flex; align-items:center; justify-content:center; flex-shrink:0; }
.name-text { font-size:13px; font-weight:600; color:#111827; }
.name-sub  { font-size:11px; color:#9ca3af; }
.text-sm   { font-size:12.5px; }
.text-xs   { font-size:11px; }
.text-gray { color:#9ca3af; }
.chip { font-size:11px; font-weight:700; padding:2px 9px; border-radius:100px; }
.chip--green  { background:#f0fdf4; color:#16a34a; }
.chip--yellow { background:#fffbeb; color:#b45309; }
.nilai-badge { font-size:14px; font-weight:800; display:flex; flex-direction:column; gap:1px; }
.grade-small { font-size:10px; font-weight:600; opacity:.7; }
.btn-primary-sm { background:#48AF4A; color:#fff; border:none; border-radius:8px; padding:6px 14px; font-size:12px; font-weight:600; cursor:pointer; font-family:inherit; }
.btn-primary-sm:hover { background:#3a9e3c; }

/* Form wrap */
.form-wrap { display:flex; flex-direction:column; gap:14px; }
.btn-back { display:inline-flex; align-items:center; gap:6px; background:transparent; border:1.5px solid #e5e7eb; border-radius:9px; padding:7px 14px; font-size:12.5px; font-weight:600; color:#6b7280; cursor:pointer; font-family:inherit; }
.btn-back:hover { border-color:#48AF4A; color:#0d2818; }

/* Peserta info */
.peserta-info-card { background:#fff; border:1px solid #e9f5e9; border-radius:14px; padding:18px 20px; display:flex; align-items:center; gap:16px; }
.peserta-avatar { width:48px; height:48px; border-radius:50%; background:#f0fdf4; color:#16a34a; font-size:20px; font-weight:700; display:flex; align-items:center; justify-content:center; flex-shrink:0; }
.peserta-info-body { flex:1; min-width:0; }
.peserta-name { font-size:15px; font-weight:700; color:#111827; }
.peserta-meta { font-size:12px; color:#6b7280; margin-top:3px; }
.peserta-preview-score { text-align:center; flex-shrink:0; }
.preview-score-num { font-size:30px; font-weight:800; line-height:1; }
.preview-score-label { font-size:11px; font-weight:600; margin-top:3px; }

/* Excel banner */
.excel-banner { background:linear-gradient(135deg, #f0fdf4 0%, #dcfce7 100%); border:1.5px solid #86efac; border-radius:14px; padding:16px 20px; display:flex; align-items:center; justify-content:space-between; gap:16px; }
.excel-banner__left { display:flex; align-items:center; gap:14px; }
.excel-icon { width:42px; height:42px; border-radius:10px; background:#fff; border:1px solid #86efac; display:flex; align-items:center; justify-content:center; flex-shrink:0; }
.excel-banner__title { font-size:14px; font-weight:700; color:#14532d; }
.excel-banner__sub   { font-size:12px; color:#166534; margin-top:2px; line-height:1.4; }
.btn-excel { display:inline-flex; align-items:center; gap:6px; background:#16a34a; color:#fff; border:none; border-radius:9px; padding:9px 18px; font-size:13px; font-weight:600; cursor:pointer; white-space:nowrap; flex-shrink:0; font-family:inherit; }
.btn-excel:hover { background:#16a34a; }
.btn-excel--loading { opacity:.7; cursor:not-allowed; }
.excel-msg { display:flex; align-items:flex-start; gap:8px; border-radius:10px; padding:10px 14px; font-size:12.5px; line-height:1.5; }
.excel-msg--ok  { background:#f0fdf4; color:#16a34a; border:1px solid #bbf7d0; }
.excel-msg--err { background:#fef2f2; color:#dc2626; border:1px solid #fecaca; }

/* Alert */
.alert { border-radius:10px; padding:11px 16px; font-size:13px; font-weight:500; }
.alert--error   { background:#fef2f2; color:#dc2626; border:1px solid #fecaca; }
.alert--success { background:#f0fdf4; color:#16a34a; border:1px solid #bbf7d0; }

/* Section cards */
.section-card { background:#fff; border:1px solid #e9f5e9; border-radius:14px; padding:20px; }
.section-header { display:flex; align-items:flex-start; gap:12px; margin-bottom:16px; }
.section-header__body { flex:1; min-width:0; }
.section-num { width:30px; height:30px; min-width:30px; background:#48AF4A; color:#fff; border-radius:50%; display:flex; align-items:center; justify-content:center; font-size:11px; font-weight:800; margin-top:1px; flex-shrink:0; }
.section-title { font-size:14px; font-weight:700; color:#111827; }
.section-sub   { font-size:12px; color:#6b7280; margin-top:2px; }
.section-avg   { margin-left:auto; font-size:12.5px; color:#6b7280; white-space:nowrap; flex-shrink:0; }

/* Komponen grid */
.komponen-grid { display:flex; flex-direction:column; gap:12px; }
.komponen-row { display:flex; align-items:center; gap:12px; }
.komponen-info { display:flex; align-items:center; gap:10px; flex:1; min-width:0; }
.komponen-num { width:22px; height:22px; border-radius:50%; background:#f3f4f6; color:#6b7280; font-size:11px; font-weight:700; display:flex; align-items:center; justify-content:center; flex-shrink:0; }
.komponen-name { font-size:13px; color:#374151; font-weight:500; line-height:1.3; }
.komponen-input-wrap { flex-shrink:0; width:150px; }
.nilai-input { width:100%; border:1.5px solid #e5e7eb; border-radius:8px; padding:7px 10px; font-size:14px; font-weight:600; color:#111827; font-family:inherit; outline:none; text-align:center; }
.nilai-input:focus { border-color:#48AF4A; box-shadow:0 0 0 3px rgba(72,175,74,.1); }
.nilai-bar-wrap { height:4px; background:#f3f4f6; border-radius:100px; overflow:hidden; margin-top:5px; }
.nilai-bar { height:100%; border-radius:100px; transition:width .3s; }

/* Kejuruan */
.kejuruan-no { width:22px; height:22px; border-radius:50%; background:#fef9c3; color:#854d0e; font-size:11px; font-weight:700; display:flex; align-items:center; justify-content:center; flex-shrink:0; }
.kejuruan-name-wrap { flex:1; min-width:0; }
.inp-komponen { width:100%; border:1.5px solid #e5e7eb; border-radius:8px; padding:7px 10px; font-size:13px; color:#374151; font-family:inherit; outline:none; }
.inp-komponen:focus { border-color:#48AF4A; box-shadow:0 0 0 3px rgba(72,175,74,.1); }
.btn-del-kejuruan { width:30px; height:30px; border-radius:7px; border:1px solid #fecaca; background:#fef2f2; color:#dc2626; display:flex; align-items:center; justify-content:center; cursor:pointer; flex-shrink:0; }
.kejuruan-empty { display:flex; align-items:center; gap:10px; padding:16px 0; color:#9ca3af; font-size:13px; }
.btn-add-kejuruan { display:inline-flex; align-items:center; gap:6px; margin-top:12px; background:#f0fdf4; color:#16a34a; border:1.5px dashed #86efac; border-radius:8px; padding:7px 16px; font-size:13px; font-weight:600; cursor:pointer; font-family:inherit; }
.btn-add-kejuruan:hover { background:#dcfce7; }

/* Penandatangan */
.penandatangan-grid { display:grid; grid-template-columns:1fr 1fr; gap:16px; }
@media (max-width: 600px) { .penandatangan-grid { grid-template-columns:1fr; } }
.penandatangan-item { display:flex; flex-direction:column; gap:6px; }
.field-label { font-size:12px; font-weight:600; color:#374151; }
.req { color:#dc2626; }
.inp-readonly { width:100%; border:1.5px solid #f3f4f6; background:#f9fafb; border-radius:8px; padding:8px 11px; font-size:13px; color:#6b7280; font-family:inherit; outline:none; cursor:default; }
.inp-field { width:100%; border:1.5px solid #e5e7eb; border-radius:8px; padding:8px 11px; font-size:13px; color:#111827; font-family:inherit; outline:none; }
.inp-field:focus { border-color:#48AF4A; box-shadow:0 0 0 3px rgba(72,175,74,.1); }
.field-hint { font-size:11px; color:#9ca3af; }

/* Keterangan */
.keterangan-card { background:#fff; border:1px solid #e9f5e9; border-radius:14px; padding:16px 20px; }
.keterangan-title { font-size:12px; font-weight:700; color:#6b7280; text-transform:uppercase; letter-spacing:.05em; margin-bottom:12px; }
.keterangan-grid { display:flex; gap:10px; flex-wrap:wrap; }
.keterangan-item { display:flex; align-items:center; gap:8px; border-radius:8px; padding:8px 14px; font-size:12.5px; flex:1; min-width:120px; }
.keterangan-item--a { background:#f0fdf4; }
.keterangan-item--b { background:#f0fdf4; }
.keterangan-item--c { background:#fffbeb; }
.keterangan-item--d { background:#fef2f2; }
.k-grade { font-size:15px; font-weight:800; }
.keterangan-item--a .k-grade { color:#16a34a; }
.keterangan-item--b .k-grade { color:#1a5c20; }
.keterangan-item--c .k-grade { color:#d97706; }
.keterangan-item--d .k-grade { color:#dc2626; }
.k-label { font-weight:600; color:#374151; }
.k-range { margin-left:auto; color:#6b7280; font-size:12px; }

/* Catatan */
.textarea-catatan { width:100%; border:1.5px solid #e5e7eb; border-radius:8px; padding:10px 12px; font-size:13px; color:#374151; font-family:inherit; outline:none; resize:vertical; }
.textarea-catatan:focus { border-color:#48AF4A; box-shadow:0 0 0 3px rgba(72,175,74,.1); }

/* Summary */
.summary-card { background:#fff; border:2px solid #48AF4A; border-radius:14px; padding:20px; }
.summary-title { font-size:13px; font-weight:700; color:#374151; margin-bottom:14px; text-transform:uppercase; letter-spacing:.04em; }
.summary-rows { display:flex; flex-direction:column; gap:8px; }
.summary-row { display:flex; align-items:center; justify-content:space-between; font-size:13px; color:#374151; }
.summary-row strong { font-weight:700; }
.summary-divider { border-top:2px solid #f3f4f6; margin:6px 0; }
.summary-row--total { font-size:15px; font-weight:700; color:#111827; }
.total-formula { font-size:11px; font-weight:400; color:#9ca3af; margin-left:4px; }
.summary-final { font-size:22px; font-weight:800; display:flex; align-items:baseline; gap:8px; }
.summary-grade { font-size:12px; font-weight:600; opacity:.8; }
.btn-simpan { margin-top:16px; width:100%; display:flex; align-items:center; justify-content:center; gap:8px; background:#48AF4A; color:#fff; border:none; border-radius:10px; padding:13px; font-size:14px; font-weight:700; cursor:pointer; font-family:inherit; }
.btn-simpan:hover { background:#3a9e3c; }
.btn-simpan:disabled { opacity:.5; cursor:not-allowed; }

.spinner-sm { width:16px; height:16px; border:2.5px solid rgba(255,255,255,.3); border-top-color:#fff; border-radius:50%; animation:spin .7s linear infinite; display:inline-block; }
</style>
