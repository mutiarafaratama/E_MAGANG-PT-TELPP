<template>
  <div class="perp-wrap">
    <!-- Loading -->
    <div v-if="loading" class="state-box">
      <div class="spinner"></div>
    </div>

    <!-- Sudah diperpanjang / ada pengajuan -->
    <template v-else>
      <!-- Banner status pengajuan jika sudah ada -->
      <div v-if="data" class="status-banner" :class="`banner--${data.status}`">
        <div class="banner-icon">
          <svg v-if="data.status === 'menunggu'" width="22" height="22" viewBox="0 0 24 24" fill="none">
            <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2"/>
            <polyline points="12 6 12 12 16 14" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
          </svg>
          <svg v-else-if="data.status === 'disetujui'" width="22" height="22" viewBox="0 0 24 24" fill="none">
            <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2"/>
            <path d="M8 12l3 3 5-5" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
          <svg v-else width="22" height="22" viewBox="0 0 24 24" fill="none">
            <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2"/>
            <line x1="12" y1="8" x2="12" y2="12" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
            <circle cx="12" cy="16" r="1" fill="currentColor"/>
          </svg>
        </div>
        <div class="banner-body">
          <div class="banner-title">{{ bannerTitle }}</div>
          <div class="banner-sub">{{ bannerSub }}</div>
          <div v-if="data.status === 'disetujui'" class="banner-detail">
            Tanggal selesai baru:
            <strong>{{ fmtDate(data.tanggal_selesai_baru) }}</strong>
            (diperpanjang {{ data.durasi_hari }} hari)
          </div>
          <div v-if="data.status === 'ditolak' && data.catatan_hrd" class="banner-catatan">
            <svg width="12" height="12" viewBox="0 0 24 24" fill="none"><path d="M21 15a2 2 0 01-2 2H7l-4 4V5a2 2 0 012-2h14a2 2 0 012 2z" stroke="currentColor" stroke-width="2"/></svg>
            Catatan HRD: <em>{{ data.catatan_hrd }}</em>
          </div>
        </div>
        <div class="banner-durasi">+{{ data.durasi_hari }} hari</div>
      </div>

      <!-- Sertifikat perpanjangan dari HRD (jika ada) -->
      <div v-if="data?.surat_path" class="surat-card">
        <div class="surat-card__label">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none">
            <path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z" stroke="currentColor" stroke-width="2"/>
            <polyline points="14 2 14 8 20 8" stroke="currentColor" stroke-width="2"/>
          </svg>
          Surat Perpanjangan dari HRD
        </div>
        <a :href="`/uploads/${data.surat_path}`" target="_blank" class="btn-dl">
          <svg width="13" height="13" viewBox="0 0 24 24" fill="none">
            <path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
            <polyline points="7 10 12 15 17 10" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
            <line x1="12" y1="15" x2="12" y2="3" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
          </svg>
          Unduh Surat
        </a>
      </div>

      <!-- Form pengajuan (tampil hanya jika belum ada pengajuan & status aktif) -->
      <div v-if="!data && pelaksanaanAktif" class="form-card">
        <div class="form-card__head">
          <div class="form-card__title">Ajukan Perpanjangan Magang</div>
          <div class="form-card__sub">Pengajuan akan ditinjau oleh HRD sebelum diterapkan.</div>
        </div>

        <div class="field">
          <label class="label">Durasi Perpanjangan</label>
          <div class="durasi-grid">
            <button
              v-for="d in durasiOptions" :key="d"
              class="durasi-btn"
              :class="{ 'durasi-btn--active': form.durasi_hari === d }"
              @click="form.durasi_hari = d"
            >
              {{ d }} hari
            </button>
          </div>
        </div>

        <div class="field">
          <label class="label">Alasan Perpanjangan <span class="req">*</span></label>
          <textarea
            v-model="form.alasan"
            class="textarea"
            rows="4"
            placeholder="Jelaskan alasan perpanjangan secara singkat (min. 10 karakter)..."
          ></textarea>
          <div class="field-hint">{{ form.alasan.length }} karakter</div>
        </div>

        <div v-if="errMsg" class="alert alert--error">{{ errMsg }}</div>
        <div v-if="okMsg" class="alert alert--ok">{{ okMsg }}</div>

        <button class="btn-submit" :disabled="submitting || form.durasi_hari === 0 || form.alasan.length < 10" @click="submit">
          <span v-if="submitting" class="spinner-sm"></span>
          <span v-else>Kirim Pengajuan</span>
        </button>
      </div>

      <!-- Info: tidak bisa ajukan karena status bukan aktif -->
      <div v-if="!data && !pelaksanaanAktif" class="state-box state-box--locked">
        <div class="locked-icon">
          <svg width="40" height="40" viewBox="0 0 24 24" fill="none">
            <rect x="3" y="11" width="18" height="11" rx="2" stroke="#f59e0b" stroke-width="2"/>
            <path d="M7 11V7a5 5 0 0110 0v4" stroke="#f59e0b" stroke-width="2" stroke-linecap="round"/>
          </svg>
        </div>
        <div class="locked-title">Pengajuan Belum Tersedia</div>
        <div class="locked-sub">Perpanjangan hanya bisa diajukan saat magang berstatus <strong>Aktif</strong>.</div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from "vue";
import api from "@/lib/api";

const props = defineProps<{ pelaksanaan: any }>();

const loading   = ref(true);
const submitting = ref(false);
const data      = ref<any>(null);
const errMsg    = ref("");
const okMsg     = ref("");

const durasiOptions = [7, 14, 30, 60, 90];

const form = ref({ durasi_hari: 0, alasan: "" });

// Sisa hari hingga tanggal_selesai (negatif jika sudah lewat)
const sisaHariSelesai = computed(() => {
  if (!props.pelaksanaan?.tanggal_selesai) return 999;
  const selesai = new Date(props.pelaksanaan.tanggal_selesai);
  const today = new Date();
  today.setHours(0, 0, 0, 0);
  selesai.setHours(0, 0, 0, 0);
  return Math.floor((selesai.getTime() - today.getTime()) / (1000 * 60 * 60 * 24));
});

// Boleh ajukan: status aktif DAN dalam 5 hari terakhir sebelum selesai
const pelaksanaanAktif = computed(() =>
  props.pelaksanaan?.status === 'aktif' && sisaHariSelesai.value <= 5
);

const bannerTitle = computed(() => {
  if (data.value?.status === "menunggu")   return "Pengajuan Perpanjangan Sedang Ditinjau";
  if (data.value?.status === "disetujui")  return "Perpanjangan Magang Disetujui";
  if (data.value?.status === "ditolak")    return "Pengajuan Perpanjangan Ditolak";
  return "";
});

const bannerSub = computed(() => {
  if (data.value?.status === "menunggu")
    return `Anda mengajukan perpanjangan ${data.value.durasi_hari} hari. Menunggu keputusan HRD.`;
  if (data.value?.status === "disetujui")
    return `Durasi magang Anda telah diperpanjang ${data.value.durasi_hari} hari oleh HRD.`;
  if (data.value?.status === "ditolak")
    return "Pengajuan tidak disetujui. Hubungi HRD untuk info lebih lanjut.";
  return "";
});

function fmtDate(d: string) {
  if (!d) return "—";
  return new Date(d).toLocaleDateString("id-ID", { day: "numeric", month: "long", year: "numeric" });
}

async function fetchData() {
  loading.value = true;
  try {
    const r = await api.get("/api/perpanjangan/saya");
    data.value = r.data?.data ?? null;
  } catch {
    data.value = null;
  } finally {
    loading.value = false;
  }
}

async function submit() {
  errMsg.value = "";
  okMsg.value  = "";
  if (form.value.durasi_hari === 0) { errMsg.value = "Pilih durasi perpanjangan."; return; }
  if (form.value.alasan.length < 10) { errMsg.value = "Alasan minimal 10 karakter."; return; }

  submitting.value = true;
  try {
    await api.post("/api/perpanjangan", { durasi_hari: form.value.durasi_hari, alasan: form.value.alasan });
    okMsg.value = "Pengajuan berhasil dikirim. Menunggu persetujuan HRD.";
    await fetchData();
  } catch (e: any) {
    errMsg.value = e?.response?.data?.message ?? "Gagal mengirim pengajuan.";
  } finally {
    submitting.value = false;
  }
}

watch(
  () => props.pelaksanaan?.id,
  async (pelID) => {
    if (!pelID) return;
    await fetchData();
  },
  { immediate: true }
);
</script>

<style scoped>
.perp-wrap { display: flex; flex-direction: column; gap: 16px; }

/* State box */
.state-box { display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 12px; padding: 48px 24px; text-align: center; }
.state-box--locked { background: #fffbeb; border: 1px solid #fde68a; border-radius: 12px; }
.locked-icon { margin-bottom: 4px; }
.locked-title { font-size: 15px; font-weight: 600; color: #92400e; }
.locked-sub   { font-size: 13px; color: #78350f; line-height: 1.5; max-width: 340px; }

/* Banner status */
.status-banner { display: flex; align-items: flex-start; gap: 14px; padding: 18px 20px; border-radius: 12px; border: 1.5px solid; }
.banner--menunggu  { background: #fefce8; border-color: #fde68a; color: #78350f; }
.banner--disetujui { background: #f0fdf4; border-color: #86efac; color: #14532d; }
.banner--ditolak   { background: #fef2f2; border-color: #fca5a5; color: #7f1d1d; }
.banner-icon  { flex-shrink: 0; margin-top: 1px; }
.banner-body  { flex: 1; }
.banner-title { font-size: 14px; font-weight: 700; margin-bottom: 4px; }
.banner-sub   { font-size: 12.5px; line-height: 1.5; opacity: 0.85; }
.banner-detail { margin-top: 6px; font-size: 12.5px; }
.banner-catatan { margin-top: 6px; font-size: 12px; opacity: 0.85; display: flex; align-items: center; gap: 5px; }
.banner-durasi { flex-shrink: 0; font-size: 20px; font-weight: 800; opacity: 0.5; white-space: nowrap; }

/* Surat card */
.surat-card { display: flex; align-items: center; justify-content: space-between; gap: 12px; padding: 14px 18px; background: #f8fafc; border: 1px solid #e2e8f0; border-radius: 10px; }
.surat-card__label { display: flex; align-items: center; gap: 7px; font-size: 13px; font-weight: 500; color: #374151; }
.btn-dl { display: flex; align-items: center; gap: 5px; padding: 6px 14px; background: #1d4ed8; color: #fff; border-radius: 7px; font-size: 12px; font-weight: 600; text-decoration: none; transition: background .15s; }
.btn-dl:hover { background: #1e40af; }

/* Form card */
.form-card { background: #fff; border: 1px solid #e5e7eb; border-radius: 12px; overflow: hidden; }
.form-card__head { padding: 20px 24px 16px; border-bottom: 1px solid #f1f5f9; }
.form-card__title { font-size: 15px; font-weight: 700; color: #111827; }
.form-card__sub   { font-size: 12.5px; color: #6b7280; margin-top: 3px; }

.field { padding: 0 24px; margin-top: 18px; }
.label { display: block; font-size: 12.5px; font-weight: 600; color: #374151; margin-bottom: 8px; }
.req { color: #ef4444; }

/* Durasi grid */
.durasi-grid { display: flex; gap: 8px; flex-wrap: wrap; }
.durasi-btn { padding: 8px 18px; border-radius: 8px; border: 1.5px solid #e5e7eb; background: #f9fafb; font-size: 13px; font-weight: 500; color: #374151; cursor: pointer; transition: all .15s; }
.durasi-btn:hover { border-color: #16a34a; color: #15803d; background: #f0fdf4; }
.durasi-btn--active { border-color: #16a34a; background: #f0fdf4; color: #15803d; font-weight: 700; }

.textarea { width: 100%; padding: 10px 13px; border: 1.5px solid #e5e7eb; border-radius: 8px; font-size: 13px; color: #111827; resize: vertical; font-family: inherit; box-sizing: border-box; }
.textarea:focus { outline: none; border-color: #16a34a; }
.field-hint { font-size: 11px; color: #9ca3af; margin-top: 4px; text-align: right; }

.alert { margin: 12px 24px 0; padding: 10px 14px; border-radius: 8px; font-size: 12.5px; }
.alert--error { background: #fef2f2; color: #b91c1c; border: 1px solid #fca5a5; }
.alert--ok    { background: #f0fdf4; color: #15803d; border: 1px solid #86efac; }

.btn-submit { display: block; width: calc(100% - 48px); margin: 18px 24px 24px; padding: 12px; background: #16a34a; color: #fff; border: none; border-radius: 10px; font-size: 14px; font-weight: 600; cursor: pointer; transition: background .15s; }
.btn-submit:hover:not(:disabled) { background: #15803d; }
.btn-submit:disabled { opacity: 0.55; cursor: not-allowed; }

.spinner { width: 32px; height: 32px; border: 3px solid #e5e7eb; border-top-color: #16a34a; border-radius: 50%; animation: spin .7s linear infinite; }
.spinner-sm { width: 16px; height: 16px; border: 2px solid rgba(255,255,255,0.4); border-top-color: #fff; border-radius: 50%; animation: spin .7s linear infinite; display: inline-block; }
@keyframes spin { to { transform: rotate(360deg); } }
</style>
