<template>
  <div>
    <!-- Loading -->
    <div v-if="loading" class="state-box"><div class="spinner"></div></div>

    <!-- Belum ada nilai -->
    <div v-else-if="!penilaian" class="state-box state-box--empty">
      <svg width="52" height="52" viewBox="0 0 24 24" fill="none">
        <polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2" stroke="#d1d5db" stroke-width="1.5"/>
      </svg>
      <div class="empty-title">Nilai Belum Tersedia</div>
      <div class="empty-sub">Nilai akhir magang akan diberikan oleh HRD setelah seluruh tahapan selesai.</div>
      <div class="status-steps">
        <div v-for="step in statusSteps" :key="step.key"
          :class="['step-item', step.done ? 'step-item--done' : step.active ? 'step-item--active' : '']">
          <div class="step-dot">
            <svg v-if="step.done" width="10" height="10" viewBox="0 0 24 24" fill="none">
              <path d="M5 13l4 4L19 7" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
          </div>
          <span>{{ step.label }}</span>
        </div>
      </div>
    </div>

    <!-- Ada nilai -->
    <template v-else>
      <div class="nilai-card">
        <!-- Topbar -->
        <div class="nilai-topbar no-print">
          <span class="nilai-topbar__title">Lembar Penilaian Magang</span>
          <button class="btn-print" @click="cetakPDF">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none">
              <polyline points="6 9 6 2 18 2 18 9" stroke="currentColor" stroke-width="2" stroke-linejoin="round"/>
              <path d="M6 18H4a2 2 0 01-2-2v-5a2 2 0 012-2h16a2 2 0 012 2v5a2 2 0 01-2 2h-2" stroke="currentColor" stroke-width="2" stroke-linejoin="round"/>
              <rect x="6" y="14" width="12" height="8" stroke="currentColor" stroke-width="2" stroke-linejoin="round"/>
            </svg>
            Cetak / Unduh PDF
          </button>
        </div>

      <!-- ═══ LEMBAR DAFTAR PENILAIAN ═══ -->
      <div id="lembar-penilaian" class="lembar">

        <!-- ── KOP SURAT ── -->
        <div class="kop">
          <img src="/logotel.png" alt="PT TELPP" class="kop-logo" />
          <div class="kop-text">
            <div class="kop-nama">PT TANJUNGENIM LESTARI PULP AND PAPER</div>
            <div class="kop-sub">Tanjung Enim, Sumatera Selatan</div>
          </div>
        </div>
        <div class="kop-divider-wrap">
          <div class="kop-divider-thick"></div>
          <div class="kop-divider-thin"></div>
        </div>
        <div class="kop-judul">DAFTAR PENILAIAN PESERTA MAGANG</div>

        <!-- ── DATA DIRI PESERTA ── -->
        <!-- SMK -->
        <template v-if="isSMK">
          <table class="info-table">
            <tr>
              <td class="it-label">Nama Siswa</td>
              <td class="it-sep">:</td>
              <td class="it-val">{{ pengajuan?.nama_lengkap ?? pelaksanaan?.nama_peserta ?? '—' }}</td>
              <td class="it-label">Unit Kerja / Divisi</td>
              <td class="it-sep">:</td>
              <td class="it-val">{{ pelaksanaan?.divisi ?? '—' }}</td>
            </tr>
            <tr>
              <td class="it-label">NISN</td>
              <td class="it-sep">:</td>
              <td class="it-val">{{ pengajuan?.nomor_induk ?? '—' }}</td>
              <td class="it-label">Nama Pembimbing</td>
              <td class="it-sep">:</td>
              <td class="it-val">{{ pelaksanaan?.pembimbing_nama ?? '—' }}</td>
            </tr>
            <tr>
              <td class="it-label">Kelas</td>
              <td class="it-sep">:</td>
              <td class="it-val">{{ pengajuan?.kelas_semester ?? '—' }}</td>
              <td class="it-label">Periode Magang</td>
              <td class="it-sep">:</td>
              <td class="it-val">{{ fmtDate(pelaksanaan?.tanggal_mulai) }} – {{ fmtDate(pelaksanaan?.tanggal_selesai) }}</td>
            </tr>
            <tr>
              <td class="it-label">Kompetensi Keahlian</td>
              <td class="it-sep">:</td>
              <td class="it-val">{{ pengajuan?.jurusan ?? '—' }}</td>
              <td class="it-label">Manager Dept.</td>
              <td class="it-sep">:</td>
              <td class="it-val">{{ penilaian.manager_nama ?? '—' }}</td>
            </tr>
            <tr>
              <td class="it-label">Asal Sekolah</td>
              <td class="it-sep">:</td>
              <td class="it-val" colspan="4">{{ pengajuan?.asal_institusi ?? '—' }}</td>
            </tr>
          </table>
        </template>

        <!-- Mahasiswa (D3/S1/S2) -->
        <template v-else-if="isMahasiswa">
          <table class="info-table">
            <tr>
              <td class="it-label">Nama Mahasiswa</td>
              <td class="it-sep">:</td>
              <td class="it-val">{{ pengajuan?.nama_lengkap ?? pelaksanaan?.nama_peserta ?? '—' }}</td>
              <td class="it-label">Unit Kerja / Divisi</td>
              <td class="it-sep">:</td>
              <td class="it-val">{{ pelaksanaan?.divisi ?? '—' }}</td>
            </tr>
            <tr>
              <td class="it-label">NIM</td>
              <td class="it-sep">:</td>
              <td class="it-val">{{ pengajuan?.nomor_induk ?? '—' }}</td>
              <td class="it-label">Nama Pembimbing</td>
              <td class="it-sep">:</td>
              <td class="it-val">{{ pelaksanaan?.pembimbing_nama ?? '—' }}</td>
            </tr>
            <tr>
              <td class="it-label">Program Studi</td>
              <td class="it-sep">:</td>
              <td class="it-val">{{ pengajuan?.jurusan ?? '—' }}</td>
              <td class="it-label">Periode Magang</td>
              <td class="it-sep">:</td>
              <td class="it-val">{{ fmtDate(pelaksanaan?.tanggal_mulai) }} – {{ fmtDate(pelaksanaan?.tanggal_selesai) }}</td>
            </tr>
            <tr>
              <td class="it-label">Semester</td>
              <td class="it-sep">:</td>
              <td class="it-val">{{ pengajuan?.kelas_semester ?? '—' }}</td>
              <td class="it-label">Manager Dept.</td>
              <td class="it-sep">:</td>
              <td class="it-val">{{ penilaian.manager_nama ?? '—' }}</td>
            </tr>
            <tr>
              <td class="it-label">Asal Perguruan Tinggi</td>
              <td class="it-sep">:</td>
              <td class="it-val" colspan="4">{{ pengajuan?.asal_institusi ?? '—' }}</td>
            </tr>
          </table>
        </template>

        <!-- Penelitian -->
        <template v-else-if="isPenelitian">
          <table class="info-table">
            <tr>
              <td class="it-label">Nama Peneliti</td>
              <td class="it-sep">:</td>
              <td class="it-val">{{ pengajuan?.nama_lengkap ?? pelaksanaan?.nama_peserta ?? '—' }}</td>
              <td class="it-label">Unit Kerja / Divisi</td>
              <td class="it-sep">:</td>
              <td class="it-val">{{ pelaksanaan?.divisi ?? '—' }}</td>
            </tr>
            <tr>
              <td class="it-label">NIM / NRP</td>
              <td class="it-sep">:</td>
              <td class="it-val">{{ pengajuan?.nomor_induk ?? '—' }}</td>
              <td class="it-label">Pembimbing Lapangan</td>
              <td class="it-sep">:</td>
              <td class="it-val">{{ pelaksanaan?.pembimbing_nama ?? '—' }}</td>
            </tr>
            <tr>
              <td class="it-label">Bidang / Jurusan</td>
              <td class="it-sep">:</td>
              <td class="it-val">{{ pengajuan?.jurusan ?? '—' }}</td>
              <td class="it-label">Periode Penelitian</td>
              <td class="it-sep">:</td>
              <td class="it-val">{{ fmtDate(pelaksanaan?.tanggal_mulai) }} – {{ fmtDate(pelaksanaan?.tanggal_selesai) }}</td>
            </tr>
            <tr>
              <td class="it-label">Semester / Angkatan</td>
              <td class="it-sep">:</td>
              <td class="it-val">{{ pengajuan?.kelas_semester ?? '—' }}</td>
              <td class="it-label">Manager Dept.</td>
              <td class="it-sep">:</td>
              <td class="it-val">{{ penilaian.manager_nama ?? '—' }}</td>
            </tr>
            <tr>
              <td class="it-label">Asal Institusi</td>
              <td class="it-sep">:</td>
              <td class="it-val" colspan="4">{{ pengajuan?.asal_institusi ?? '—' }}</td>
            </tr>
          </table>
        </template>

        <!-- Default / Lainnya -->
        <template v-else>
          <table class="info-table">
            <tr>
              <td class="it-label">Nama Peserta</td>
              <td class="it-sep">:</td>
              <td class="it-val">{{ pengajuan?.nama_lengkap ?? pelaksanaan?.nama_peserta ?? '—' }}</td>
              <td class="it-label">Unit Kerja / Divisi</td>
              <td class="it-sep">:</td>
              <td class="it-val">{{ pelaksanaan?.divisi ?? '—' }}</td>
            </tr>
            <tr>
              <td class="it-label">Nomor Induk</td>
              <td class="it-sep">:</td>
              <td class="it-val">{{ pengajuan?.nomor_induk ?? '—' }}</td>
              <td class="it-label">Nama Pembimbing</td>
              <td class="it-sep">:</td>
              <td class="it-val">{{ pelaksanaan?.pembimbing_nama ?? '—' }}</td>
            </tr>
            <tr>
              <td class="it-label">Jurusan</td>
              <td class="it-sep">:</td>
              <td class="it-val">{{ pengajuan?.jurusan ?? '—' }}</td>
              <td class="it-label">Periode</td>
              <td class="it-sep">:</td>
              <td class="it-val">{{ fmtDate(pelaksanaan?.tanggal_mulai) }} – {{ fmtDate(pelaksanaan?.tanggal_selesai) }}</td>
            </tr>
            <tr>
              <td class="it-label">Kelas / Semester</td>
              <td class="it-sep">:</td>
              <td class="it-val">{{ pengajuan?.kelas_semester ?? '—' }}</td>
              <td class="it-label">Manager Dept.</td>
              <td class="it-sep">:</td>
              <td class="it-val">{{ penilaian.manager_nama ?? '—' }}</td>
            </tr>
            <tr>
              <td class="it-label">Asal Institusi</td>
              <td class="it-sep">:</td>
              <td class="it-val" colspan="4">{{ pengajuan?.asal_institusi ?? '—' }}</td>
            </tr>
          </table>
        </template>

        <!-- ── TABEL PENILAIAN (satu tabel, colgroup fixed) ── -->
        <table class="nilai-table" style="margin-top:14px">
          <colgroup>
            <col class="col-no" />
            <col class="col-unsur" />
            <col class="col-angka" />
            <col class="col-ket" />
          </colgroup>
          <thead>
            <tr>
              <th class="th-no">No</th>
              <th>Unsur yang Dinilai</th>
              <th class="th-angka">Nilai (Angka)</th>
              <th class="th-ket">Keterangan</th>
            </tr>
          </thead>
          <tbody>
            <!-- I. KEPRIBADIAN -->
            <tr class="tr-section">
              <td colspan="4" class="td-section">I. KEPRIBADIAN / ETOS KERJA</td>
            </tr>
            <tr v-for="(k, i) in kepribadianRows" :key="k.key">
              <td class="td-no">{{ i + 1 }}</td>
              <td>{{ k.label }}</td>
              <td class="td-angka">{{ fmt((penilaian as any)[k.key]) }}</td>
              <td class="td-ket">{{ keteranganNilai((penilaian as any)[k.key]) }}</td>
            </tr>

            <!-- II. KEJURUAN -->
            <tr class="tr-section">
              <td colspan="4" class="td-section">II. KEMAMPUAN KEJURUAN</td>
            </tr>
            <tr v-for="(k, i) in penilaian.kejuruan ?? []" :key="'kj-'+i">
              <td class="td-no">{{ i + 1 }}</td>
              <td>{{ k.komponen }}</td>
              <td class="td-angka">{{ k.nilai > 0 ? k.nilai.toFixed(1) : '—' }}</td>
              <td class="td-ket">{{ keteranganNilai(k.nilai) }}</td>
            </tr>
            <tr v-if="!penilaian.kejuruan?.length">
              <td class="td-no">—</td>
              <td colspan="3" class="td-empty">Tidak ada data kejuruan</td>
            </tr>

            <!-- III. K3 -->
            <tr class="tr-section">
              <td colspan="4" class="td-section">III. K3 DAN PRODUKTIVITAS</td>
            </tr>
            <tr v-for="(k, i) in k3Rows" :key="k.key">
              <td class="td-no">{{ i + 1 }}</td>
              <td>{{ k.label }}</td>
              <td class="td-angka">{{ fmt((penilaian as any)[k.key]) }}</td>
              <td class="td-ket">{{ keteranganNilai((penilaian as any)[k.key]) }}</td>
            </tr>

            <!-- IV. PRESENTASI (hanya mahasiswa / penelitian) -->
            <template v-if="isMahasiswa || isPenelitian">
              <tr class="tr-section">
                <td colspan="4" class="td-section">IV. PRESENTASI *)</td>
              </tr>
              <tr v-for="(k, i) in presentasiRows" :key="k.key">
                <td class="td-no">{{ i + 1 }}</td>
                <td>{{ k.label }}</td>
                <td class="td-angka">{{ fmt((penilaian as any)[k.key]) }}</td>
                <td class="td-ket">{{ keteranganNilai((penilaian as any)[k.key]) }}</td>
              </tr>
            </template>

            <!-- JUMLAH & RATA-RATA -->
            <tr class="tr-jumlah">
              <td class="td-no"></td>
              <td><strong>JUMLAH</strong></td>
              <td class="td-angka"><strong>{{ totalNilai.toFixed(1) }}</strong></td>
              <td></td>
            </tr>
            <tr class="tr-final">
              <td class="td-no"></td>
              <td><strong>NILAI RATA-RATA</strong></td>
              <td class="td-angka td-final">
                <strong>{{ penilaian.nilai_akhir?.toFixed(2) ?? '—' }}</strong>
              </td>
              <td class="td-ket td-grade">
                {{ nilaiGrade(penilaian.nilai_akhir ?? 0) }}
              </td>
            </tr>
          </tbody>
        </table>

        <!-- KETERANGAN NILAI -->
        <div class="keterangan-box">
          <div class="ket-row">
            <span class="ket-title">KETERANGAN:</span>
            <span class="ket-item">"A" Baik Sekali : 85 – 100</span>
            <span class="ket-sep">|</span>
            <span class="ket-item">"B" Baik : 75 – 84</span>
            <span class="ket-sep">|</span>
            <span class="ket-item">"C" Cukup : 60 – 74</span>
            <span class="ket-sep">|</span>
            <span class="ket-item">"D" Kurang : &lt; 60</span>
          </div>
          <div v-if="isMahasiswa || isPenelitian" class="ket-note">
            *) Kolom Presentasi diisi hanya untuk peserta Perguruan Tinggi (D3/S1/S2)
          </div>
        </div>

        <!-- CATATAN PEMBIMBING -->
        <div v-if="penilaian.catatan" class="catatan-box">
          <span class="catatan-label">Catatan Pembimbing: </span>
          <em>{{ penilaian.catatan }}</em>
        </div>

        <!-- TTD -->
        <div class="ttd-row">
          <div class="ttd-box">
            <div class="ttd-title">{{ isSMK ? 'Manager (Dept. terkait)' : 'Manager (Dept. terkait)' }}</div>
            <div class="ttd-space"></div>
            <div class="ttd-name">{{ penilaian.manager_nama || '…………………………' }}</div>
          </div>
          <div class="ttd-spacer"></div>
          <div class="ttd-box">
            <div class="ttd-title">Pembimbing Lapangan,</div>
            <div class="ttd-space"></div>
            <div class="ttd-name">{{ pelaksanaan?.pembimbing_nama ?? '…………………………' }}</div>
          </div>
        </div>

        <div class="footer-note">
          Dicetak pada {{ new Date().toLocaleDateString('id-ID', { day:'2-digit', month:'long', year:'numeric' }) }}
          · e-Magang PT TanjungEnim Lestari Pulp and Paper
        </div>
      </div><!-- end lembar -->
      </div><!-- end nilai-card -->
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from "vue";
import api from "@/lib/api";

const props = defineProps<{
  pelaksanaan: any;
  pengajuan?: any;
}>();

interface Kejuruan { komponen: string; nilai: number }
interface Penilaian {
  id: string; pelaksanaan_id: string;
  nilai_motivasi?: number; nilai_inisiatif?: number; nilai_disiplin_waktu?: number;
  nilai_kerajinan?: number; nilai_kreativitas?: number; nilai_tanggung_jawab?: number;
  nilai_kerjasama?: number; nilai_adaptasi?: number; nilai_kehadiran?: number;
  kejuruan?: Kejuruan[];
  nilai_k3_safety?: number; nilai_k3_metode?: number;
  nilai_k3_manajemen?: number; nilai_k3_volume?: number;
  nilai_prs_proses?: number; nilai_prs_teori?: number;
  nilai_prs_judul?: number; nilai_prs_data?: number;
  manager_nama?: string;
  catatan?: string; nilai_akhir?: number; dinilai_at?: string;
}

const loading    = ref(false);
const penilaian  = ref<Penilaian | null>(null);

// ── Deteksi kategori ──
const kategori = computed<string>(() =>
  props.pengajuan?.kategori_magang ?? props.pelaksanaan?.kategori_magang ?? ""
);
const isSMK       = computed(() => kategori.value === "smk");
const isMahasiswa = computed(() =>
  ["d3", "s1", "s2", "d3_s1_s2"].includes(kategori.value)
);
const isPenelitian = computed(() => kategori.value === "penelitian");

// ── Baris komponen per bagian ──
const kepribadianRows = [
  { key: "nilai_motivasi",       label: "Motivasi" },
  { key: "nilai_inisiatif",      label: "Inisiatif" },
  { key: "nilai_disiplin_waktu", label: "Disiplin Waktu" },
  { key: "nilai_kerajinan",      label: "Kerajinan" },
  { key: "nilai_kreativitas",    label: "Kreativitas" },
  { key: "nilai_tanggung_jawab", label: "Tanggung Jawab" },
  { key: "nilai_kerjasama",      label: "Kerjasama" },
  { key: "nilai_adaptasi",       label: "Adaptasi dengan Lingkungan Kerja" },
  { key: "nilai_kehadiran",      label: "Kehadiran" },
];
const k3Rows = [
  { key: "nilai_k3_safety",    label: "Penggunaan dan Penerapan K3 (Safety Equipment)" },
  { key: "nilai_k3_metode",    label: "Metode Kerja dalam Menggunakan Alat dan Bahan" },
  { key: "nilai_k3_manajemen", label: "Manajemen Kerja" },
  { key: "nilai_k3_volume",    label: "Volume Pekerjaan yang Diselesaikan" },
];
const presentasiRows = [
  { key: "nilai_prs_proses", label: "Pemahaman Proses Pulp PT TELPP" },
  { key: "nilai_prs_teori",  label: "Pemahaman Teori" },
  { key: "nilai_prs_judul",  label: "Pemahaman Judul Praktek Kerja Lapangan" },
  { key: "nilai_prs_data",   label: "Kemampuan Mengelolah Data Dengan Baik dan Benar" },
];

const totalNilai = computed(() => {
  const p = penilaian.value;
  if (!p) return 0;
  const kepSum = kepribadianRows.reduce((s, k) => s + ((p as any)[k.key] ?? 0), 0);
  const kjSum  = (p.kejuruan ?? []).reduce((s, k) => s + (k.nilai ?? 0), 0);
  const k3Sum  = k3Rows.reduce((s, k) => s + ((p as any)[k.key] ?? 0), 0);
  const prsSum = (isMahasiswa.value || isPenelitian.value)
    ? presentasiRows.reduce((s, k) => s + ((p as any)[k.key] ?? 0), 0)
    : 0;
  return kepSum + kjSum + k3Sum + prsSum;
});

function fmt(v?: number) { return v != null && v > 0 ? v.toFixed(1) : '—'; }
function fmtDate(d?: string) {
  if (!d) return '—';
  return new Date(d).toLocaleDateString("id-ID", { day: "2-digit", month: "long", year: "numeric" });
}
function keteranganNilai(v?: number) {
  if (v == null || v === 0) return '—';
  if (v >= 85) return 'A — Baik Sekali';
  if (v >= 75) return 'B — Baik';
  if (v >= 60) return 'C — Cukup';
  return 'D — Kurang';
}
function nilaiColor(n: number) {
  if (n >= 85) return "#16a34a";
  if (n >= 75) return "#1a5c20";
  if (n >= 60) return "#d97706";
  return "#dc2626";
}
function nilaiGrade(n: number) {
  if (n >= 85) return "A — Baik Sekali";
  if (n >= 75) return "B — Baik";
  if (n >= 60) return "C — Cukup";
  return "D — Kurang";
}

const statusSteps = computed(() => {
  const status = props.pelaksanaan?.status ?? "";
  const order  = ["aktif", "upload_laporan", "penilaian", "selesai"];
  const labels: Record<string, string> = {
    aktif: "Magang Aktif", upload_laporan: "Upload Laporan",
    penilaian: "Proses Penilaian", selesai: "Nilai Keluar",
  };
  const currentIdx = order.indexOf(status);
  return order.map((key, i) => ({
    key, label: labels[key],
    done: i < currentIdx, active: i === currentIdx,
  }));
});

function cetakPDF() { window.print(); }

watch(
  () => props.pelaksanaan?.id,
  async (pelID) => {
    if (!pelID) return;
    loading.value = true;
    try {
      const res = await api.get(`/api/penilaian/${pelID}`);
      penilaian.value = res.data ?? null;
    } catch { penilaian.value = null; }
    finally { loading.value = false; }
  },
  { immediate: true }
);
</script>

<style scoped>
/* ── State ── */
.state-box { display:flex; flex-direction:column; align-items:center; justify-content:center; gap:12px; padding:48px 24px; background:#fff; border-radius:14px; border:1px solid #e9f5e9; text-align:center; }
.spinner { width:28px; height:28px; border:3px solid #e5e7eb; border-top-color:#48AF4A; border-radius:50%; animation:spin .7s linear infinite; }
@keyframes spin { to { transform:rotate(360deg); } }
.empty-title { font-size:16px; font-weight:700; color:#374151; }
.empty-sub   { font-size:12.5px; color:#9ca3af; line-height:1.7; max-width:380px; }
.status-steps { display:flex; align-items:center; flex-wrap:wrap; justify-content:center; gap:0; margin-top:12px; }
.step-item { display:flex; align-items:center; gap:6px; font-size:11.5px; color:#9ca3af; font-weight:500; }
.step-item + .step-item::before { content:"›"; margin:0 6px; color:#d1d5db; }
.step-item--done   { color:#16a34a; }
.step-item--active { color:#0d2818; font-weight:700; }
.step-dot { width:18px; height:18px; border-radius:50%; border:2px solid currentColor; display:flex; align-items:center; justify-content:center; flex-shrink:0; }
.step-item--done .step-dot   { background:#f0fdf4; border-color:#16a34a; color:#16a34a; }
.step-item--active .step-dot { background:#f0fdf4; border-color:#0d2818; }

/* ── Kartu nilai wrapper ── */
.nilai-card { background:#fff; border:1.5px solid #e9f5e9; border-radius:14px; overflow:hidden; box-shadow:0 1px 3px rgba(13,40,24,.04); }
.nilai-topbar { display:flex; align-items:center; justify-content:space-between; padding:12px 16px; border-bottom:1px solid #f0fdf4; background:#fafffe; gap:10px; }
.nilai-topbar__title { font-size:13.5px; font-weight:700; color:#111827; }
.btn-print { display:inline-flex; align-items:center; gap:7px; background:#48AF4A; color:#fff; border:none; border-radius:8px; padding:7px 16px; font-size:12.5px; font-weight:600; cursor:pointer; font-family:inherit; flex-shrink:0; }
.btn-print:hover { background:#3d9e3f; }

/* ── Lembar utama ── */
.lembar { background:#fff; padding:24px 28px; font-size:13px; color:#111827; line-height:1.5; }

/* ── Kop surat ── */
.kop { display:flex; align-items:center; gap:14px; margin-bottom:8px; }
.kop-logo { height:52px; width:auto; object-fit:contain; flex-shrink:0; }
.kop-text { }
.kop-nama { font-size:15px; font-weight:800; color:#0d2818; letter-spacing:.02em; text-transform:uppercase; }
.kop-sub  { font-size:11.5px; color:#6b7280; margin-top:2px; }
.kop-divider-wrap { margin-bottom:6px; }
.kop-divider-thick { border-top:3px solid #0d2818; }
.kop-divider-thin  { border-top:1px solid #0d2818; margin-top:2px; }
.kop-judul { font-size:13px; font-weight:800; text-align:center; letter-spacing:.08em; text-transform:uppercase; color:#0d2818; margin:8px 0 14px; text-decoration:underline; text-underline-offset:3px; }

/* ── Tabel info peserta ── */
.info-table { width:100%; border-collapse:collapse; margin-bottom:4px; font-size:12.5px; }
.info-table tr td { padding:2.5px 4px; vertical-align:top; }
.it-label { width:150px; color:#374151; font-weight:500; white-space:nowrap; }
.it-sep   { width:14px; color:#374151; }
.it-val   { color:#111827; font-weight:600; padding-right:20px; }

/* ── Tabel nilai ── */
.nilai-table { width:100%; border-collapse:collapse; font-size:12.5px; table-layout:fixed; }
.nilai-table th { background:#f9fafb; border:1px solid #d1d5db; padding:6px 8px; font-size:11px; font-weight:700; color:#374151; text-align:left; }
.nilai-table td { border:1px solid #d1d5db; padding:5px 8px; vertical-align:middle; overflow:hidden; }

/* colgroup: lebar kolom fixed */
.col-no    { width:38px; }
.col-angka { width:110px; }
.col-ket   { width:160px; }
.col-unsur { /* sisa lebar otomatis */ }

.th-no    { text-align:center; }
.th-angka { text-align:center; }
.th-ket   { text-align:center; }
.td-no    { text-align:center; color:#6b7280; font-size:11.5px; }
.td-angka { text-align:center; font-weight:700; font-size:14px; }
.td-ket   { font-size:11.5px; color:#6b7280; }
.td-empty { text-align:center; color:#9ca3af; font-style:italic; }

/* Baris judul seksi (pengganti div.section-label) */
.tr-section .td-section {
  font-size:11.5px; font-weight:800; letter-spacing:.04em;
  color:#0d2818; background:#f0fdf4;
  border-left:3px solid #48AF4A;
  padding:4px 10px; text-transform:uppercase;
}

.tr-jumlah td { background:#f9fafb; }
.tr-final  td { background:#f3f4f6; }
.td-final  { font-size:16px !important; }
.td-grade  { font-weight:700; color:#111827; }

/* ── Keterangan ── */
.keterangan-box { margin-top:12px; padding:8px 12px; background:#f9fafb; border:1px solid #d1d5db; border-radius:4px; }
.ket-row { display:flex; align-items:center; gap:8px; flex-wrap:wrap; }
.ket-title { font-size:11px; font-weight:700; color:#374151; white-space:nowrap; margin-right:4px; }
.ket-item { font-size:11.5px; font-weight:500; color:#374151; white-space:nowrap; }
.ket-sep { font-size:11.5px; color:#9ca3af; }
.ket-note { font-size:11px; color:#6b7280; font-style:italic; margin-top:4px; }

/* ── Catatan ── */
.catatan-box { margin-top:10px; font-size:12.5px; padding:8px 12px; background:#f9fafb; border:1px solid #d1d5db; border-radius:4px; color:#374151; }
.catatan-label { font-weight:700; }

/* ── TTD ── */
.ttd-row    { display:flex; justify-content:space-between; margin:24px 0 12px; }
.ttd-spacer { flex:1; }
.ttd-box    { text-align:center; width:210px; }
.ttd-title  { font-size:12px; font-weight:600; color:#374151; margin-bottom:4px; }
.ttd-space  { height:64px; border-bottom:1px solid #374151; margin:0 10px; }
.ttd-name   { font-size:11.5px; font-weight:600; margin-top:6px; color:#374151; }

/* ── Footer ── */
.footer-note { font-size:10.5px; color:#9ca3af; text-align:center; margin-top:12px; border-top:1px solid #f3f4f6; padding-top:10px; }

/* ── Responsive: info-table 1 kolom di mobile ── */
@media (max-width: 640px) {
  .info-table,
  .info-table tbody { display: block; }

  .info-table tr {
    display: grid;
    grid-template-columns: auto auto 1fr;
    margin-bottom: 1px;
  }

  /* kolom kiri (1-3) → baris grid 1 */
  .info-table td:nth-child(1) { grid-column: 1; grid-row: 1; }
  .info-table td:nth-child(2) { grid-column: 2; grid-row: 1; }
  .info-table td:nth-child(3) { grid-column: 3; grid-row: 1; padding-right: 0; }

  /* kolom kanan (4-6) → baris grid 2 (stacked di bawah) */
  .info-table td:nth-child(4) { grid-column: 1; grid-row: 2; }
  .info-table td:nth-child(5) { grid-column: 2; grid-row: 2; }
  .info-table td:nth-child(6) { grid-column: 3; grid-row: 2; padding-right: 0; }

  /* label kiri sedikit lebih lebar di mobile */
  .it-label { width: auto; min-width: 110px; }
}

/* ── PRINT ── */
@media print {
  .no-print { display:none !important; }
  .nilai-card { border:none !important; box-shadow:none !important; border-radius:0 !important; }
  .lembar { padding:16px; }
  body, html { background:#fff !important; }
  .section-label { -webkit-print-color-adjust:exact; print-color-adjust:exact; }
  .tr-jumlah td, .tr-final td { -webkit-print-color-adjust:exact; print-color-adjust:exact; }
}
</style>

<!-- CSS global (non-scoped) untuk sembunyikan layout saat cetak -->
<style>
@media print {
  .topbar,
  .sidebar,
  .sidebar-overlay,
  .topbar-hamburger { display: none !important; }

  .main-area {
    margin-left: 0 !important;
    padding: 0 !important;
  }

  .app-shell {
    display: block !important;
  }
}
</style>
