<template>
  <div>
    <!-- Loading -->
    <div v-if="loading" class="state-box">
      <div class="spinner"></div>
      <span>Memeriksa sertifikat…</span>
    </div>

    <!-- Belum selesai -->
    <div v-else-if="!isSelesai" class="state-box state-box--empty">
      <div class="lock-circle">
        <svg width="36" height="36" viewBox="0 0 24 24" fill="none">
          <rect x="3" y="11" width="18" height="11" rx="2" stroke="#9ca3af" stroke-width="2"/>
          <path d="M7 11V7a5 5 0 0110 0v4" stroke="#9ca3af" stroke-width="2" stroke-linecap="round"/>
          <circle cx="12" cy="16" r="1.5" fill="#9ca3af"/>
        </svg>
      </div>
      <div class="empty-title">Sertifikat Belum Tersedia</div>
      <div class="empty-sub">
        Sertifikat magang akan diterbitkan oleh HRD setelah seluruh rangkaian magang Anda selesai dan nilai akhir telah ditetapkan.
      </div>
      <div class="progress-steps">
        <div v-for="(s, i) in steps" :key="i"
          :class="['pstep', s.done ? 'pstep--done' : s.current ? 'pstep--current' : '']">
          <div class="pstep__dot">
            <svg v-if="s.done" width="10" height="10" viewBox="0 0 24 24" fill="none">
              <path d="M5 13l4 4L19 7" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
            <span v-else-if="s.current" class="pulse"></span>
          </div>
          <span class="pstep__label">{{ s.label }}</span>
        </div>
      </div>
    </div>

    <!-- Sertifikat tersedia -->
    <template v-else>

      <!-- Card header -->
      <div class="cert-card">
        <div class="cert-card__top">
          <div class="cert-card__left">
            <div class="cert-icon">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
                <circle cx="12" cy="8" r="6" stroke="currentColor" stroke-width="2"/>
                <path d="M15.477 12.89L17 22l-5-3-5 3 1.523-9.11" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              </svg>
            </div>
            <div>
              <div class="cert-title">Sertifikat Magang Tersedia</div>
              <div class="cert-sub">Sertifikat resmi telah diterbitkan oleh HRD PT Tanjung Enim Lestari Pulp and Paper.</div>
            </div>
          </div>
          <div class="cert-badge">
            <svg width="12" height="12" viewBox="0 0 24 24" fill="none">
              <path d="M5 13l4 4L19 7" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
            Diterbitkan
          </div>
        </div>

        <div class="cert-actions">
          <button class="btn-primary" :disabled="pdfLoading" @click="downloadSertifikat">
            <div v-if="pdfLoading" class="btn-spinner"></div>
            <svg v-else width="15" height="15" viewBox="0 0 24 24" fill="none">
              <path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
              <polyline points="7 10 12 15 17 10" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              <line x1="12" y1="15" x2="12" y2="3" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
            </svg>
            {{ pdfLoading ? 'Memuat…' : 'Unduh Sertifikat PDF' }}
          </button>
          <button class="btn-ghost" @click="previewSertifikat">
            <svg width="15" height="15" viewBox="0 0 24 24" fill="none">
              <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z" stroke="currentColor" stroke-width="2"/>
              <circle cx="12" cy="12" r="3" stroke="currentColor" stroke-width="2"/>
            </svg>
            Preview
          </button>
        </div>
      </div>

      <!-- Bento grid info -->
      <div class="bento">
        <div class="bento-item">
          <div class="bento-label">Nama Peserta</div>
          <div class="bento-val">{{ namaUser }}</div>
        </div>
        <div class="bento-item">
          <div class="bento-label">Institusi</div>
          <div class="bento-val">{{ institusi }}</div>
        </div>
        <div class="bento-item">
          <div class="bento-label">Divisi</div>
          <div class="bento-val">{{ pelaksanaan?.divisi ?? '—' }}</div>
        </div>
        <div class="bento-item">
          <div class="bento-label">Periode</div>
          <div class="bento-val">{{ fmtPeriode }}</div>
        </div>
        <div class="bento-item">
          <div class="bento-label">Nilai Akhir</div>
          <div class="bento-val bento-val--score" :style="{ color: nilaiColor }">
            {{ pelaksanaan?.nilai?.toFixed(1) ?? '—' }}
          </div>
        </div>
        <div class="bento-item">
          <div class="bento-label">Predikat</div>
          <div class="bento-val bento-val--score" :style="{ color: nilaiColor }">{{ nilaiGrade }}</div>
        </div>
      </div>

    </template>

    <!-- PDF Preview Modal -->
    <Teleport to="body">
      <div v-if="pdfModal" class="pdf-backdrop" @click.self="pdfModal = false">
        <div class="pdf-modal">
          <div class="pdf-modal__head">
            <span>Sertifikat — {{ namaUser }}</span>
            <div class="pdf-modal__acts">
              <button class="pdf-dl-btn" @click="downloadSertifikat">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none">
                  <path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
                  <polyline points="7 10 12 15 17 10" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                  <line x1="12" y1="15" x2="12" y2="3" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
                </svg>
                Unduh
              </button>
              <button class="pdf-close-btn" @click="pdfModal = false">✕</button>
            </div>
          </div>
          <div class="pdf-modal__body">
            <div v-if="!pdfSrc" class="pdf-loading">
              <div class="spinner"></div><span>Memuat PDF…</span>
            </div>
            <iframe v-else :src="pdfSrc" class="pdf-iframe" />
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from "vue";
import api from "@/lib/api";
import { useAuth } from "@/hooks/useAuth";

const props = defineProps<{ pelaksanaan: any; pengajuan?: any }>();
const { user } = useAuth();
const loading    = ref(false);
const pdfLoading = ref(false);
const pdfModal   = ref(false);
const pdfSrc     = ref("");

const namaUser  = computed(() =>
  props.pengajuan?.nama_lengkap ?? user.value?.nama_lengkap ?? ""
);
const institusi = computed(() =>
  props.pengajuan?.asal_institusi ?? "—"
);
const isSelesai = computed(() => props.pelaksanaan?.sertifikat_generated === true);

const nilaiColor = computed(() => {
  const n = props.pelaksanaan?.nilai ?? 0;
  if (n >= 85) return "#16a34a";
  if (n >= 75) return "#1a5c20";
  if (n >= 65) return "#d97706";
  return "#dc2626";
});

const nilaiGrade = computed(() => {
  const n = props.pelaksanaan?.nilai ?? 0;
  if (n >= 85) return "Sangat Baik";
  if (n >= 75) return "Baik";
  if (n >= 65) return "Cukup";
  return "Kurang";
});

const fmtPeriode = computed(() => {
  const p = props.pelaksanaan;
  if (!p) return "—";
  const fmt = (iso: string) =>
    new Date(iso).toLocaleDateString("id-ID", { day: "numeric", month: "long", year: "numeric" });
  return `${fmt(p.tanggal_mulai)} – ${fmt(p.tanggal_selesai)}`;
});

const steps = computed(() => {
  const status = props.pelaksanaan?.status ?? "";
  const order  = ["aktif", "upload_laporan", "penilaian", "selesai"];
  const labels: Record<string, string> = {
    aktif: "Magang Aktif", upload_laporan: "Upload Laporan",
    penilaian: "Proses Penilaian", selesai: "Sertifikat Tersedia",
  };
  const ci = order.indexOf(status);
  return order.map((key, i) => ({ key, label: labels[key], done: i < ci, current: i === ci }));
});

async function downloadSertifikat() {
  const id = props.pelaksanaan?.id;
  if (!id) return;
  pdfLoading.value = true;
  try {
    const r = await api.get(`/api/pelaksanaan/${id}/sertifikat/download`, { responseType: "blob" });
    const url = URL.createObjectURL(new Blob([r.data], { type: "application/pdf" }));
    const a = document.createElement("a");
    a.href = url;
    a.download = `sertifikat_magang_${namaUser.value.replace(/\s+/g, "_")}.pdf`;
    a.click();
    URL.revokeObjectURL(url);
  } catch (e: any) {
    alert(e.response?.data?.message ?? "Gagal mengunduh sertifikat.");
  } finally {
    pdfLoading.value = false;
  }
}

async function previewSertifikat() {
  pdfModal.value = true;
  if (pdfSrc.value) return;
  const id = props.pelaksanaan?.id;
  if (!id) return;
  try {
    const r = await api.get(`/api/pelaksanaan/${id}/sertifikat/download`, { responseType: "blob" });
    pdfSrc.value = URL.createObjectURL(new Blob([r.data], { type: "application/pdf" }));
  } catch {
    pdfSrc.value = "";
  }
}
</script>

<style scoped>
/* ── State ── */
.state-box { display:flex; flex-direction:column; align-items:center; gap:14px; padding:52px 24px; background:#fff; border-radius:14px; border:1px solid #e9f5e9; text-align:center; }
.spinner   { width:28px; height:28px; border:3px solid #e5e7eb; border-top-color:#48AF4A; border-radius:50%; animation:spin .7s linear infinite; }
@keyframes spin { to { transform:rotate(360deg); } }
.lock-circle  { width:72px; height:72px; background:#f9fafb; border-radius:50%; display:flex; align-items:center; justify-content:center; border:2px solid #e5e7eb; margin:0 auto; }
.empty-title  { font-size:16px; font-weight:700; color:#374151; }
.empty-sub    { font-size:12.5px; color:#9ca3af; line-height:1.7; max-width:380px; }
.progress-steps { display:flex; align-items:center; flex-wrap:wrap; justify-content:center; margin-top:6px; }
.pstep { display:flex; align-items:center; gap:6px; font-size:11.5px; color:#d1d5db; font-weight:500; padding:4px 6px; }
.pstep + .pstep::before { content:"›"; margin-right:6px; color:#e5e7eb; }
.pstep--done    { color:#16a34a; }
.pstep--current { color:#48AF4A; font-weight:700; }
.pstep__dot { width:18px; height:18px; border-radius:50%; border:2px solid currentColor; display:flex; align-items:center; justify-content:center; flex-shrink:0; }
.pstep--done .pstep__dot { background:#f0fdf4; border-color:#16a34a; color:#16a34a; }
.pulse { width:6px; height:6px; background:#48AF4A; border-radius:50%; animation:pulse 1.4s infinite; }
@keyframes pulse { 0%,100%{opacity:1;transform:scale(1);}50%{opacity:.5;transform:scale(1.3);} }

/* ── Header card ── */
.cert-card {
  background:#fff;
  border:1.5px solid #e9f5e9;
  border-radius:14px;
  overflow:hidden;
  margin-bottom:12px;
  box-shadow:0 1px 3px rgba(13,40,24,.04);
}
.cert-card__top {
  display:flex; align-items:center; justify-content:space-between;
  gap:14px; padding:12px 16px;
  border-bottom:1px solid #f0fdf4; background:#fafffe;
  flex-wrap:wrap;
}
.cert-card__left { display:flex; align-items:center; gap:12px; flex:1; min-width:0; }
.cert-icon {
  width:38px; height:38px; border-radius:10px; flex-shrink:0;
  background:#0d2818; color:#fff;
  display:flex; align-items:center; justify-content:center;
}
.cert-title { font-size:13.5px; font-weight:700; color:#111827; margin-bottom:2px; }
.cert-sub   { font-size:11.5px; color:#6b7280; line-height:1.4; }
.cert-badge {
  display:inline-flex; align-items:center; gap:5px; flex-shrink:0;
  font-size:12px; font-weight:600; color:#16a34a;
  white-space:nowrap;
}

/* ── Tombol ── */
.cert-actions { display:flex; gap:10px; flex-wrap:wrap; padding:14px 16px; border-top:1px solid #f0fdf4; }
.btn-primary {
  display:inline-flex; align-items:center; gap:7px;
  background:#48AF4A; color:#fff; border:none;
  font-size:13px; font-weight:700; padding:9px 18px;
  border-radius:9px; cursor:pointer; font-family:inherit;
  transition:background .15s;
}
.btn-primary:hover:not(:disabled) { background:#3d9e3f; }
.btn-primary:disabled { opacity:.5; cursor:not-allowed; }
.btn-ghost {
  display:inline-flex; align-items:center; gap:7px;
  background:#fff; color:#374151;
  border:1.5px solid #d1d5db;
  font-size:13px; font-weight:600; padding:9px 18px;
  border-radius:9px; cursor:pointer; font-family:inherit;
  transition:background .15s;
}
.btn-ghost:hover { background:#f9fafb; }
.btn-spinner { width:14px; height:14px; border:2px solid rgba(255,255,255,.35); border-top-color:#fff; border-radius:50%; animation:spin .7s linear infinite; }

/* ── Bento grid ── */
.bento {
  display:grid;
  grid-template-columns:repeat(3, 1fr);
  gap:10px;
}
@media (max-width: 560px) {
  .bento { grid-template-columns:repeat(2, 1fr); }
}
.bento-item {
  background:#fff;
  border:1px solid #e9f5e9;
  border-radius:11px;
  padding:14px 16px;
}
.bento-label {
  font-size:10.5px; font-weight:700; color:#9ca3af;
  text-transform:uppercase; letter-spacing:.05em; margin-bottom:6px;
}
.bento-val { font-size:14px; font-weight:600; color:#111827; }
.bento-val--score { font-size:18px; font-weight:800; }

/* ── PDF Modal ── */
.pdf-backdrop { position:fixed; inset:0; background:rgba(0,0,0,.55); z-index:9999; display:flex; align-items:center; justify-content:center; padding:16px; }
.pdf-modal { background:#fff; border-radius:14px; display:flex; flex-direction:column; width:100%; max-width:860px; height:90vh; overflow:hidden; box-shadow:0 24px 64px rgba(0,0,0,.25); }
.pdf-modal__head { display:flex; align-items:center; justify-content:space-between; padding:14px 18px; border-bottom:1px solid #e5e7eb; font-weight:700; color:#111827; font-size:13.5px; background:#f9fafb; }
.pdf-modal__acts { display:flex; gap:8px; align-items:center; }
.pdf-dl-btn { display:flex; align-items:center; gap:6px; background:#48AF4A; color:#fff; border:none; border-radius:7px; padding:6px 14px; font-size:12px; font-weight:600; cursor:pointer; font-family:inherit; }
.pdf-dl-btn:hover { background:#3d9e3f; }
.pdf-close-btn { background:#fff; border:1.5px solid #e5e7eb; border-radius:7px; width:30px; height:30px; cursor:pointer; font-size:14px; color:#6b7280; display:flex; align-items:center; justify-content:center; }
.pdf-close-btn:hover { background:#f3f4f6; }
.pdf-modal__body { flex:1; display:flex; overflow:hidden; }
.pdf-loading { display:flex; flex-direction:column; align-items:center; justify-content:center; gap:10px; flex:1; color:#6b7280; font-size:13px; }
.pdf-iframe { width:100%; height:100%; border:none; }
</style>
