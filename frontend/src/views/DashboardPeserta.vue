<template>
  <DashboardLayout :nav-groups="navGroups" role-name="Peserta Magang" default-tab="beranda" ref="layout">
    <template #default>

      <!-- ── BERANDA ── -->
      <template v-if="activeTab === 'beranda'">
        <div class="welcome-banner">
          <div class="welcome-banner__left">
            <div class="welcome-banner__greeting">Selamat datang, {{ firstName }}!</div>
            <div class="welcome-banner__sub">Pantau perjalanan magang Anda dari sini.</div>
          </div>
          <div class="welcome-banner__date">
            <div class="welcome-date__day">{{ todayDay }}</div>
            <div class="welcome-date__info">{{ todayDate }}</div>
          </div>
        </div>

        <div class="stats-row">
          <div class="stat-card">
            <div class="stat-card__icon" style="background:#f0fdf4;color:#16a34a">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none"><path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z" stroke="currentColor" stroke-width="2"/><polyline points="14 2 14 8 20 8" stroke="currentColor" stroke-width="2"/></svg>
            </div>
            <div>
              <div class="stat-card__label">Status Pengajuan</div>
              <div class="stat-card__value">{{ statusPengajuanLabel }}</div>
            </div>
          </div>
          <div class="stat-card">
            <div class="stat-card__icon" style="background:#eff6ff;color:#3b82f6">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none"><rect x="3" y="4" width="18" height="18" rx="2" ry="2" stroke="currentColor" stroke-width="2"/><line x1="3" y1="10" x2="21" y2="10" stroke="currentColor" stroke-width="2"/></svg>
            </div>
            <div><div class="stat-card__label">Total Absensi</div><div class="stat-card__value">0 Hari</div></div>
          </div>
          <div class="stat-card">
            <div class="stat-card__icon" style="background:#fefce8;color:#ca8a04">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2" stroke="currentColor" stroke-width="2"/></svg>
            </div>
            <div><div class="stat-card__label">Nilai Akhir</div><div class="stat-card__value">—</div></div>
          </div>
          <div class="stat-card">
            <div class="stat-card__icon" style="background:#fdf4ff;color:#9333ea">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none"><path d="M21 15a2 2 0 01-2 2H7l-4 4V5a2 2 0 012-2h14a2 2 0 012 2z" stroke="currentColor" stroke-width="2"/></svg>
            </div>
            <div><div class="stat-card__label">Tiket Helpdesk</div><div class="stat-card__value">0 Aktif</div></div>
          </div>
        </div>

        <div class="card">
          <div class="card-header">
            <h3 class="card-title">Panduan Langkah Magang</h3>
            <span class="badge badge--gray">Step 1 of 4</span>
          </div>
          <div class="steps">
            <div class="step step--active">
              <div class="step__num">1</div>
              <div class="step__body">
                <div class="step__title">Buat Pengajuan Magang</div>
                <div class="step__desc">Lengkapi data diri dan informasi akademis Anda</div>
              </div>
              <button class="btn-green-sm" @click="setTab('pengajuan')">Mulai →</button>
            </div>
            <div class="step">
              <div class="step__num">2</div>
              <div class="step__body">
                <div class="step__title">Upload Dokumen Persyaratan</div>
                <div class="step__desc">KTP, KTM, proposal magang, pasfoto, BPJS/KIS</div>
              </div>
            </div>
            <div class="step">
              <div class="step__num">3</div>
              <div class="step__body">
                <div class="step__title">Tunggu Verifikasi HRD</div>
                <div class="step__desc">Tim HRD akan mereview berkas Anda dalam 3–5 hari kerja</div>
              </div>
            </div>
            <div class="step">
              <div class="step__num">4</div>
              <div class="step__body">
                <div class="step__title">Mulai Magang dan Absensi Harian</div>
                <div class="step__desc">Lakukan absensi setiap hari selama masa magang</div>
              </div>
            </div>
          </div>
        </div>
      </template>

      <!-- ── PENGAJUAN ── -->
      <template v-else-if="activeTab === 'pengajuan'">

        <!-- Sudah tersubmit — tampilkan status -->
        <template v-if="pengajuanSaya">
          <div class="card">
            <div class="card-header">
              <h3 class="card-title">Status Pengajuan Magang</h3>
              <span :class="['badge', statusBadgeClass(pengajuanSaya.status)]">{{ statusLabel(pengajuanSaya.status) }}</span>
            </div>
            <div class="status-detail">
              <div class="status-detail__row">
                <div class="status-detail__item">
                  <div class="status-detail__label">Nama Lengkap</div>
                  <div class="status-detail__value">{{ pengajuanSaya.nama_lengkap }}</div>
                </div>
                <div class="status-detail__item">
                  <div class="status-detail__label">Institusi</div>
                  <div class="status-detail__value">{{ pengajuanSaya.asal_institusi }}</div>
                </div>
                <div class="status-detail__item">
                  <div class="status-detail__label">Jurusan</div>
                  <div class="status-detail__value">{{ pengajuanSaya.jurusan }}</div>
                </div>
                <div class="status-detail__item">
                  <div class="status-detail__label">Kategori</div>
                  <div class="status-detail__value">{{ kategoriLabel(pengajuanSaya.kategori_magang) }}</div>
                </div>
              </div>
              <div v-if="pengajuanSaya.catatan_hrd" class="status-catatan">
                <div class="status-catatan__title">Catatan dari HRD:</div>
                <div class="status-catatan__body">{{ pengajuanSaya.catatan_hrd }}</div>
              </div>
              <div class="status-timeline">
                <div class="status-timeline__item status-timeline__item--done">
                  <div class="status-timeline__dot"></div>
                  <div class="status-timeline__text">Pengajuan dikirim</div>
                  <div class="status-timeline__date">{{ formatTanggal(pengajuanSaya.created_at) }}</div>
                </div>
                <div :class="['status-timeline__item', ['diproses','diterima'].includes(pengajuanSaya.status) ? 'status-timeline__item--done' : '']">
                  <div class="status-timeline__dot"></div>
                  <div class="status-timeline__text">Verifikasi berkas HRD</div>
                </div>
                <div :class="['status-timeline__item', pengajuanSaya.status === 'diterima' ? 'status-timeline__item--done' : '']">
                  <div class="status-timeline__dot"></div>
                  <div class="status-timeline__text">Pengajuan diterima</div>
                </div>
              </div>
            </div>
          </div>
        </template>

        <!-- Belum submit — tampilkan form multi-step -->
        <template v-else>
          <div class="form-card">
            <!-- Header -->
            <div class="form-header">
              <div class="form-header__title">Formulir Pengajuan Magang</div>
              <div class="form-header__sub">PT TanjungEnim Lestari Pulp and Paper</div>
            </div>

            <!-- Stepper -->
            <div class="stepper">
              <div v-for="(s, i) in STEPS" :key="s.id" class="stepper-item">
                <div class="stepper-col">
                  <div :class="['stepper-dot', formStep === s.id ? 'stepper-dot--active' : formStep > s.id ? 'stepper-dot--done' : '']">
                    <template v-if="formStep > s.id">
                      <svg width="12" height="12" viewBox="0 0 24 24" fill="none"><path d="M5 13l4 4L19 7" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"/></svg>
                    </template>
                    <template v-else>{{ s.id }}</template>
                  </div>
                  <div :class="['stepper-label', formStep === s.id ? 'stepper-label--active' : formStep > s.id ? 'stepper-label--done' : '']">{{ s.label }}</div>
                </div>
                <div v-if="i < STEPS.length - 1" :class="['stepper-line', formStep > s.id ? 'stepper-line--done' : '']"></div>
              </div>
            </div>

            <!-- Step 1: Data Diri -->
            <div v-if="formStep === 1" class="step-body">
              <div class="step-head">
                <div class="step-title">Data Diri</div>
                <div class="step-sub">Lengkapi informasi pribadi sesuai KTP</div>
              </div>
              <div class="form-grid-2">
                <div class="field">
                  <label class="field-label">Nama Lengkap <span class="req">*</span></label>
                  <input class="inp" v-model="diri.nama_lengkap" placeholder="Sesuai KTP" />
                </div>
                <div class="field">
                  <label class="field-label">Jenis Kelamin <span class="req">*</span></label>
                  <select class="inp" v-model="diri.jenis_kelamin">
                    <option value="">-- Pilih --</option>
                    <option value="laki_laki">Laki-laki</option>
                    <option value="perempuan">Perempuan</option>
                  </select>
                </div>
                <div class="field">
                  <label class="field-label">Tempat Lahir <span class="req">*</span></label>
                  <input class="inp" v-model="diri.tempat_lahir" placeholder="Kota kelahiran" />
                </div>
                <div class="field">
                  <label class="field-label">Tanggal Lahir <span class="req">*</span></label>
                  <input class="inp" type="date" v-model="diri.tanggal_lahir" />
                </div>
                <div class="field">
                  <label class="field-label">Nomor HP / WhatsApp <span class="req">*</span></label>
                  <input class="inp" v-model="diri.no_hp" placeholder="08xxxxxxxxxx" />
                </div>
                <div class="field">
                  <label class="field-label">Email Aktif <span class="req">*</span></label>
                  <input class="inp" type="email" v-model="diri.email" placeholder="nama@email.com" />
                </div>
              </div>
              <div class="field">
                <label class="field-label">Alamat Lengkap <span class="req">*</span></label>
                <textarea class="inp" v-model="diri.alamat" rows="3" placeholder="Jalan, nomor, kelurahan, kecamatan, kota, provinsi"></textarea>
              </div>
            </div>

            <!-- Step 2: Akademis -->
            <div v-else-if="formStep === 2" class="step-body">
              <div class="step-head">
                <div class="step-title">Informasi Akademis</div>
                <div class="step-sub">Data institusi dan program studi Anda</div>
              </div>
              <div class="field">
                <label class="field-label">Kategori Magang <span class="req">*</span></label>
                <div class="kategori-grid">
                  <button
                    v-for="k in KATEGORI"
                    :key="k.value"
                    type="button"
                    :class="['kategori-btn', akademis.kategori_magang === k.value ? 'kategori-btn--active' : '']"
                    @click="akademis.kategori_magang = k.value"
                  >
                    <span class="kategori-radio"></span>
                    <span>{{ k.label }}</span>
                  </button>
                </div>
              </div>
              <div class="form-grid-2">
                <div class="field">
                  <label class="field-label">Asal Institusi / Sekolah <span class="req">*</span></label>
                  <input class="inp" v-model="akademis.asal_institusi" placeholder="Nama sekolah / universitas" />
                </div>
                <div class="field">
                  <label class="field-label">Jurusan / Program Studi <span class="req">*</span></label>
                  <input class="inp" v-model="akademis.jurusan" placeholder="Teknik Informatika" />
                </div>
                <div class="field">
                  <label class="field-label">{{ akademis.kategori_magang === 'smk' ? 'Kelas' : 'Semester' }} <span class="req">*</span></label>
                  <input class="inp" v-model="akademis.kelas_semester" :placeholder="akademis.kategori_magang === 'smk' ? 'XI / XII' : 'Semester 5'" />
                </div>
                <div class="field">
                  <label class="field-label">{{ akademis.kategori_magang === 'smk' ? 'NIS / NISN' : 'NIM' }} <span class="req">*</span></label>
                  <input class="inp" v-model="akademis.nomor_induk" placeholder="Nomor induk" />
                </div>
              </div>
              <div class="info-box">
                Kategori magang menentukan dokumen yang perlu diunggah pada langkah berikutnya.
              </div>
            </div>

            <!-- Step 3: Dokumen -->
            <div v-else-if="formStep === 3" class="step-body">
              <div class="step-head">
                <div class="step-title">Upload Dokumen</div>
                <div class="step-sub">Unggah semua berkas yang diperlukan. Maks 10 MB per file.</div>
              </div>
              <div class="dok-grid">
                <div v-for="dok in DOKUMEN_LIST" :key="dok.key" class="dok-item">
                  <div class="dok-item__header">
                    <span class="dok-item__label">{{ dok.label }}</span>
                    <span v-if="dok.required" class="badge-required">Wajib</span>
                    <span v-else class="badge-optional">Opsional</span>
                  </div>
                  <div
                    :class="['dropzone', dragKey === dok.key ? 'dropzone--drag' : '', files[dok.key] ? 'dropzone--done' : '']"
                    @dragover.prevent="dragKey = dok.key"
                    @dragleave="dragKey = ''"
                    @drop.prevent="onDrop($event, dok.key)"
                    @click="triggerInput(dok.key)"
                  >
                    <input
                      :ref="el => { if (el) fileInputs[dok.key] = el as HTMLInputElement }"
                      type="file"
                      :accept="dok.accept"
                      style="display:none"
                      @change="onFileChange($event, dok.key)"
                    />
                    <template v-if="files[dok.key]">
                      <div class="dropzone-done">
                        <svg width="16" height="16" viewBox="0 0 24 24" fill="none"><circle cx="12" cy="12" r="10" fill="#48AF4A"/><path d="M8 12l3 3 5-6" stroke="#fff" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>
                        <div class="dropzone-done__info">
                          <div class="dropzone-done__name">{{ files[dok.key]!.name }}</div>
                          <div class="dropzone-done__size">{{ fmtSize(files[dok.key]!.size) }}</div>
                        </div>
                        <button class="dropzone-done__remove" type="button" @click.stop="files[dok.key] = null">
                          <svg width="14" height="14" viewBox="0 0 24 24" fill="none"><path d="M18 6L6 18M6 6l12 12" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>
                        </button>
                      </div>
                    </template>
                    <template v-else>
                      <div class="dropzone-empty">
                        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" style="flex-shrink:0;color:#9ca3af"><path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4" stroke="currentColor" stroke-width="1.8" stroke-linecap="round"/><polyline points="17 8 12 3 7 8" stroke="currentColor" stroke-width="1.8" stroke-linecap="round"/><line x1="12" y1="3" x2="12" y2="15" stroke="currentColor" stroke-width="1.8" stroke-linecap="round"/></svg>
                        <div>
                          <div class="dropzone-empty__text">Klik atau drag file ke sini</div>
                          <div class="dropzone-empty__hint">{{ dok.hint }}</div>
                        </div>
                      </div>
                    </template>
                  </div>
                </div>
              </div>
            </div>

            <!-- Step 4: Review & Submit -->
            <div v-else-if="formStep === 4" class="step-body">
              <div class="step-head">
                <div class="step-title">Review dan Kirim</div>
                <div class="step-sub">Periksa kembali data sebelum mengirim pengajuan</div>
              </div>

              <div class="review-section">
                <div class="review-section__title">Data Diri</div>
                <div class="review-grid">
                  <div class="review-row"><div class="review-row__label">Nama Lengkap</div><div class="review-row__value">{{ diri.nama_lengkap || '—' }}</div></div>
                  <div class="review-row"><div class="review-row__label">Jenis Kelamin</div><div class="review-row__value">{{ diri.jenis_kelamin === 'laki_laki' ? 'Laki-laki' : diri.jenis_kelamin === 'perempuan' ? 'Perempuan' : '—' }}</div></div>
                  <div class="review-row"><div class="review-row__label">Tempat Lahir</div><div class="review-row__value">{{ diri.tempat_lahir || '—' }}</div></div>
                  <div class="review-row"><div class="review-row__label">Tanggal Lahir</div><div class="review-row__value">{{ diri.tanggal_lahir || '—' }}</div></div>
                  <div class="review-row"><div class="review-row__label">No. HP</div><div class="review-row__value">{{ diri.no_hp || '—' }}</div></div>
                  <div class="review-row"><div class="review-row__label">Email</div><div class="review-row__value">{{ diri.email || '—' }}</div></div>
                  <div class="review-row review-row--full"><div class="review-row__label">Alamat</div><div class="review-row__value">{{ diri.alamat || '—' }}</div></div>
                </div>
              </div>

              <div class="review-section">
                <div class="review-section__title">Data Akademis</div>
                <div class="review-grid">
                  <div class="review-row"><div class="review-row__label">Kategori</div><div class="review-row__value">{{ kategoriLabel(akademis.kategori_magang) }}</div></div>
                  <div class="review-row"><div class="review-row__label">Institusi</div><div class="review-row__value">{{ akademis.asal_institusi || '—' }}</div></div>
                  <div class="review-row"><div class="review-row__label">Jurusan</div><div class="review-row__value">{{ akademis.jurusan || '—' }}</div></div>
                  <div class="review-row"><div class="review-row__label">Kelas / Semester</div><div class="review-row__value">{{ akademis.kelas_semester || '—' }}</div></div>
                  <div class="review-row"><div class="review-row__label">Nomor Induk</div><div class="review-row__value">{{ akademis.nomor_induk || '—' }}</div></div>
                </div>
              </div>

              <div class="review-section">
                <div class="review-section__title">Dokumen</div>
                <div class="dok-status-grid">
                  <div v-for="d in DOKUMEN_LIST" :key="d.key" :class="['dok-status', files[d.key] ? 'dok-status--ok' : d.required ? 'dok-status--miss' : '']">
                    <svg width="14" height="14" viewBox="0 0 24 24" fill="none">
                      <template v-if="files[d.key]"><circle cx="12" cy="12" r="10" fill="#16a34a"/><path d="M8 12l3 3 5-6" stroke="#fff" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></template>
                      <template v-else-if="d.required"><circle cx="12" cy="12" r="10" fill="#fecaca"/><path d="M15 9l-6 6M9 9l6 6" stroke="#dc2626" stroke-width="2" stroke-linecap="round"/></template>
                      <template v-else><circle cx="12" cy="12" r="10" stroke="#d1d5db" stroke-width="1.5" fill="none"/></template>
                    </svg>
                    <span>{{ d.label }}</span>
                    <span v-if="d.required && !files[d.key]" class="dok-missing">Belum upload</span>
                  </div>
                </div>
              </div>

              <div v-if="missingRequired.length > 0" class="alert-warn">
                Masih ada {{ missingRequired.length }} dokumen wajib yang belum diunggah. Kembali ke langkah Dokumen.
              </div>

              <div class="review-note">
                Pastikan semua data sudah benar. Pengajuan yang sudah dikirim tidak dapat diedit. Tim HRD akan memverifikasi berkas dalam <strong>3–5 hari kerja</strong>.
              </div>

              <div v-if="submitError" class="alert-warn" style="margin-top:0">{{ submitError }}</div>
            </div>

            <!-- Footer navigasi -->
            <div class="form-footer">
              <div>
                <button v-if="formStep > 1" class="btn-back" type="button" @click="formStep--" :disabled="submitting">Kembali</button>
              </div>
              <div class="step-counter">Langkah {{ formStep }} dari {{ STEPS.length }}</div>
              <button
                v-if="formStep < 4"
                class="btn-next"
                type="button"
                @click="nextStep"
              >Lanjut</button>
              <button
                v-else
                class="btn-next btn-submit"
                type="button"
                @click="submitPengajuan"
                :disabled="submitting || missingRequired.length > 0"
              >
                <template v-if="submitting">
                  <span class="spinner"></span> Mengirim...
                </template>
                <template v-else>Kirim Pengajuan</template>
              </button>
            </div>
          </div>
        </template>
      </template>

      <!-- ── DOKUMEN ── -->
      <template v-else-if="activeTab === 'dokumen'">
        <div class="card">
          <div class="card-header"><h3 class="card-title">Dokumen Saya</h3></div>
          <div class="empty-state">
            <div class="empty-state__icon"><svg width="40" height="40" viewBox="0 0 24 24" fill="none"><path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4" stroke="#d1d5db" stroke-width="1.5"/><polyline points="17 8 12 3 7 8" stroke="#d1d5db" stroke-width="1.5"/><line x1="12" y1="3" x2="12" y2="15" stroke="#d1d5db" stroke-width="1.5"/></svg></div>
            <p>Belum ada dokumen terupload.<br/>Upload dokumen setelah mengisi pengajuan.</p>
          </div>
        </div>
      </template>

      <!-- ── ABSENSI ── -->
      <template v-else-if="activeTab === 'absensi'">
        <div class="card">
          <div class="card-header"><h3 class="card-title">Absensi Harian</h3></div>
          <div class="empty-state">
            <div class="empty-state__icon"><svg width="40" height="40" viewBox="0 0 24 24" fill="none"><rect x="3" y="4" width="18" height="18" rx="2" ry="2" stroke="#d1d5db" stroke-width="1.5"/><line x1="3" y1="10" x2="21" y2="10" stroke="#d1d5db" stroke-width="1.5"/></svg></div>
            <p>Absensi tersedia setelah magang dimulai.</p>
          </div>
        </div>
      </template>

      <!-- ── NOTIFIKASI ── -->
      <template v-else-if="activeTab === 'notifikasi'">
        <div class="card">
          <div class="card-header"><h3 class="card-title">Notifikasi</h3></div>
          <div class="empty-state">
            <div class="empty-state__icon"><svg width="40" height="40" viewBox="0 0 24 24" fill="none"><path d="M18 8A6 6 0 006 8c0 7-3 9-3 9h18s-3-2-3-9" stroke="#d1d5db" stroke-width="1.5"/><path d="M13.73 21a2 2 0 01-3.46 0" stroke="#d1d5db" stroke-width="1.5"/></svg></div>
            <p>Belum ada notifikasi.</p>
          </div>
        </div>
      </template>

      <!-- ── CHAT ── -->
      <template v-else-if="activeTab === 'chat'">
        <div class="card">
          <div class="card-header">
            <h3 class="card-title">Chat Helpdesk</h3>
            <button class="btn-green-sm">+ Buat Tiket</button>
          </div>
          <div class="empty-state">
            <div class="empty-state__icon"><svg width="40" height="40" viewBox="0 0 24 24" fill="none"><path d="M21 15a2 2 0 01-2 2H7l-4 4V5a2 2 0 012-2h14a2 2 0 012 2z" stroke="#d1d5db" stroke-width="1.5"/></svg></div>
            <p>Belum ada tiket helpdesk.</p>
          </div>
        </div>
      </template>

    </template>
  </DashboardLayout>
</template>

<script setup lang="ts">
import { ref, computed, reactive, onMounted } from "vue";
import DashboardLayout from "@/layouts/DashboardLayout.vue";
import { useAuth } from "@/hooks/useAuth";
import api from "@/lib/api";

const { user } = useAuth();
const layout = ref<InstanceType<typeof DashboardLayout> | null>(null);
const activeTab = computed(() => layout.value?.activeTab ?? "beranda");

const firstName = computed(() => user.value?.nama_lengkap?.split(" ")[0] ?? "");

const now = new Date();
const todayDay = now.toLocaleDateString("id-ID", { weekday: "long" });
const todayDate = now.toLocaleDateString("id-ID", { day: "numeric", month: "long", year: "numeric" });

function setTab(tab: string) {
  if (layout.value) layout.value.activeTab = tab;
}

// ── Pengajuan state ──────────────────────────────────────────
const pengajuanSaya = ref<any>(null);
const loadingPengajuan = ref(false);

async function fetchPengajuanSaya() {
  try {
    loadingPengajuan.value = true;
    const res = await api.get("/api/peserta/pengajuan/saya");
    pengajuanSaya.value = res.data?.data ?? null;
  } catch (e: any) {
    if (e?.response?.status === 404 || e?.response?.status === 204) {
      pengajuanSaya.value = null;
    }
  } finally {
    loadingPengajuan.value = false;
  }
}

onMounted(fetchPengajuanSaya);

const statusPengajuanLabel = computed(() => {
  if (!pengajuanSaya.value) return "Belum Diajukan";
  return statusLabel(pengajuanSaya.value.status);
});

function statusLabel(status: string) {
  const m: Record<string, string> = {
    diajukan:            "Baru Diajukan",
    menunggu_verifikasi: "Menunggu Verifikasi",
    diproses:            "Sedang Diproses",
    diterima:            "Diterima",
    ditolak:             "Ditolak",
  };
  return m[status] ?? status;
}

function statusBadgeClass(status: string) {
  const m: Record<string, string> = {
    diajukan:            "badge badge--yellow",
    menunggu_verifikasi: "badge badge--yellow",
    diproses:            "badge badge--blue",
    diterima:            "badge badge--green",
    ditolak:             "badge badge--red",
  };
  return m[status] ?? "badge badge--gray";
}

function kategoriLabel(val: string) {
  const m: Record<string, string> = {
    smk:        "SMK / Sekolah Menengah Kejuruan",
    d3_s1_s2:  "D3 / S1 / S2 (Perguruan Tinggi)",
    penelitian: "Penelitian / Riset",
  };
  return m[val] ?? val;
}

function formatTanggal(iso: string) {
  if (!iso) return "—";
  return new Date(iso).toLocaleDateString("id-ID", { day: "numeric", month: "long", year: "numeric" });
}

// ── Form multi-step ──────────────────────────────────────────
const STEPS = [
  { id: 1, label: "Data Diri" },
  { id: 2, label: "Akademis" },
  { id: 3, label: "Dokumen" },
  { id: 4, label: "Review & Kirim" },
];

const KATEGORI = [
  { value: "smk",        label: "SMK / Sekolah Menengah Kejuruan" },
  { value: "d3_s1_s2",  label: "D3 / S1 / S2 (Perguruan Tinggi)" },
  { value: "penelitian", label: "Penelitian / Riset" },
];

const DOKUMEN_LIST = [
  { key: "proposal_magang", label: "Proposal Magang",       required: true,  accept: ".pdf,.doc,.docx", hint: "PDF atau Word, maks 10 MB" },
  { key: "ktp",             label: "KTP / Kartu Identitas", required: true,  accept: "image/*,.pdf",    hint: "JPG, PNG, atau PDF" },
  { key: "ktm",             label: "KTM / Kartu Pelajar",   required: false, accept: "image/*,.pdf",    hint: "JPG, PNG, atau PDF" },
  { key: "pasfoto",         label: "Pasfoto 3x4",           required: true,  accept: "image/*",         hint: "JPG atau PNG, latar merah/biru" },
  { key: "bpjs_kis",        label: "BPJS / KIS",            required: true,  accept: "image/*,.pdf",    hint: "Foto kartu atau PDF" },
];

const formStep = ref(1);
const submitting = ref(false);
const submitError = ref("");
const dragKey = ref("");

const diri = reactive({
  nama_lengkap: user.value?.nama_lengkap ?? "",
  jenis_kelamin: "",
  tempat_lahir: "",
  tanggal_lahir: "",
  no_hp: "",
  email: user.value?.email ?? "",
  alamat: "",
});

const akademis = reactive({
  kategori_magang: "d3_s1_s2",
  asal_institusi: "",
  jurusan: "",
  kelas_semester: "",
  nomor_induk: "",
});

const files = reactive<Record<string, File | null>>({
  proposal_magang: null,
  ktp: null,
  ktm: null,
  pasfoto: null,
  bpjs_kis: null,
});

const fileInputs: Record<string, HTMLInputElement> = {};

function triggerInput(key: string) {
  fileInputs[key]?.click();
}

function onFileChange(e: Event, key: string) {
  const f = (e.target as HTMLInputElement).files?.[0];
  if (f) files[key] = f;
}

function onDrop(e: DragEvent, key: string) {
  dragKey.value = "";
  const f = e.dataTransfer?.files?.[0];
  if (f) files[key] = f;
}

function fmtSize(b: number) {
  return b > 1048576 ? `${(b / 1048576).toFixed(1)} MB` : `${(b / 1024).toFixed(0)} KB`;
}

const missingRequired = computed(() =>
  DOKUMEN_LIST.filter(d => d.required && !files[d.key])
);

function validateStep(step: number): string {
  if (step === 1) {
    if (!diri.nama_lengkap.trim()) return "Nama lengkap wajib diisi";
    if (!diri.jenis_kelamin)       return "Jenis kelamin wajib dipilih";
    if (!diri.tempat_lahir.trim()) return "Tempat lahir wajib diisi";
    if (!diri.tanggal_lahir)       return "Tanggal lahir wajib diisi";
    if (!diri.no_hp.trim())        return "Nomor HP wajib diisi";
    if (!diri.email.trim())        return "Email wajib diisi";
    if (!diri.alamat.trim())       return "Alamat wajib diisi";
  }
  if (step === 2) {
    if (!akademis.asal_institusi.trim()) return "Asal institusi wajib diisi";
    if (!akademis.jurusan.trim())        return "Jurusan wajib diisi";
    if (!akademis.kelas_semester.trim()) return "Kelas/Semester wajib diisi";
    if (!akademis.nomor_induk.trim())    return "Nomor induk wajib diisi";
  }
  return "";
}

function nextStep() {
  const err = validateStep(formStep.value);
  if (err) {
    submitError.value = err;
    setTimeout(() => { submitError.value = ""; }, 3000);
    return;
  }
  submitError.value = "";
  formStep.value++;
}

async function submitPengajuan() {
  if (missingRequired.value.length > 0) return;
  submitting.value = true;
  submitError.value = "";

  try {
    // 1. Submit data pengajuan
    const res = await api.post("/api/peserta/pengajuan", {
      step1: {
        nama_lengkap:  diri.nama_lengkap,
        tempat_lahir:  diri.tempat_lahir,
        tanggal_lahir: diri.tanggal_lahir,
        jenis_kelamin: diri.jenis_kelamin,
        alamat:        diri.alamat,
        no_hp:         diri.no_hp,
        email:         diri.email,
      },
      step2: {
        kategori_magang: akademis.kategori_magang,
        nomor_induk:     akademis.nomor_induk,
        asal_institusi:  akademis.asal_institusi,
        jurusan:         akademis.jurusan,
        kelas_semester:  akademis.kelas_semester,
      },
    });

    const pengajuanId: string = res.data?.data?.id ?? "";

    // 2. Upload setiap dokumen
    const uploads = DOKUMEN_LIST.filter(d => files[d.key]);
    for (const dok of uploads) {
      const fd = new FormData();
      fd.append("file", files[dok.key] as File);
      fd.append("jenis", dok.key);
      if (pengajuanId) fd.append("pengajuan_id", pengajuanId);
      await api.post("/api/peserta/dokumen/upload", fd, {
        headers: { "Content-Type": "multipart/form-data" },
      });
    }

    // 3. Refresh status
    await fetchPengajuanSaya();

  } catch (e: any) {
    submitError.value = e?.response?.data?.message ?? "Terjadi kesalahan, coba lagi.";
  } finally {
    submitting.value = false;
  }
}

// ── Nav ──────────────────────────────────────────────────────
const ICON = {
  home:     `<svg width="16" height="16" viewBox="0 0 24 24" fill="none"><path d="M3 9l9-7 9 7v11a2 2 0 01-2 2H5a2 2 0 01-2-2z" stroke="currentColor" stroke-width="2"/><polyline points="9 22 9 12 15 12 15 22" stroke="currentColor" stroke-width="2"/></svg>`,
  file:     `<svg width="16" height="16" viewBox="0 0 24 24" fill="none"><path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z" stroke="currentColor" stroke-width="2"/><polyline points="14 2 14 8 20 8" stroke="currentColor" stroke-width="2"/></svg>`,
  upload:   `<svg width="16" height="16" viewBox="0 0 24 24" fill="none"><path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4" stroke="currentColor" stroke-width="2"/><polyline points="17 8 12 3 7 8" stroke="currentColor" stroke-width="2"/><line x1="12" y1="3" x2="12" y2="15" stroke="currentColor" stroke-width="2"/></svg>`,
  calendar: `<svg width="16" height="16" viewBox="0 0 24 24" fill="none"><rect x="3" y="4" width="18" height="18" rx="2" ry="2" stroke="currentColor" stroke-width="2"/><line x1="16" y1="2" x2="16" y2="6" stroke="currentColor" stroke-width="2"/><line x1="8" y1="2" x2="8" y2="6" stroke="currentColor" stroke-width="2"/><line x1="3" y1="10" x2="21" y2="10" stroke="currentColor" stroke-width="2"/></svg>`,
  bell:     `<svg width="16" height="16" viewBox="0 0 24 24" fill="none"><path d="M18 8A6 6 0 006 8c0 7-3 9-3 9h18s-3-2-3-9" stroke="currentColor" stroke-width="2"/><path d="M13.73 21a2 2 0 01-3.46 0" stroke="currentColor" stroke-width="2"/></svg>`,
  chat:     `<svg width="16" height="16" viewBox="0 0 24 24" fill="none"><path d="M21 15a2 2 0 01-2 2H7l-4 4V5a2 2 0 012-2h14a2 2 0 012 2z" stroke="currentColor" stroke-width="2"/></svg>`,
};

const navGroups = [
  {
    items: [
      { key: "beranda", label: "Beranda", icon: ICON.home },
    ],
  },
  {
    label: "Magang Saya",
    items: [
      { key: "pengajuan", label: "Pengajuan Magang", icon: ICON.file },
      { key: "dokumen",   label: "Dokumen Saya",     icon: ICON.upload },
      { key: "absensi",   label: "Absensi Harian",   icon: ICON.calendar },
    ],
  },
  {
    label: "Komunikasi",
    items: [
      { key: "notifikasi", label: "Notifikasi",    icon: ICON.bell },
      { key: "chat",       label: "Chat Helpdesk", icon: ICON.chat },
    ],
  },
];
</script>

<style scoped>
/* ── Banner & Stats ─────────────────────────────────────────── */
.welcome-banner {
  background: linear-gradient(135deg, #0d2818 0%, #1a5c20 100%);
  border-radius: 14px;
  padding: 24px 28px;
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}
.welcome-banner__greeting { font-size: 17px; font-weight: 700; }
.welcome-banner__sub { font-size: 12.5px; color: rgba(255,255,255,0.65); margin-top: 4px; }
.welcome-date__day { font-size: 22px; font-weight: 800; color: #86efac; text-align: right; }
.welcome-date__info { font-size: 11px; color: rgba(255,255,255,0.55); text-align: right; }

.stats-row { display: grid; grid-template-columns: repeat(4,1fr); gap: 14px; }
.stat-card {
  background: #fff; border-radius: 12px; padding: 16px;
  display: flex; align-items: center; gap: 12px;
  border: 1px solid #e9f5e9; box-shadow: 0 1px 3px rgba(13,40,24,0.05);
}
.stat-card__icon { width: 38px; height: 38px; border-radius: 10px; display: flex; align-items: center; justify-content: center; flex-shrink: 0; }
.stat-card__label { font-size: 11px; color: #6b7280; font-weight: 500; }
.stat-card__value { font-size: 13.5px; font-weight: 700; color: #111827; margin-top: 2px; }

/* ── Shared card ────────────────────────────────────────────── */
.card { background: #fff; border-radius: 14px; border: 1px solid #e9f5e9; box-shadow: 0 1px 3px rgba(13,40,24,0.05); overflow: hidden; }
.card-header { display: flex; align-items: center; justify-content: space-between; padding: 16px 20px; border-bottom: 1px solid #f0faf0; }
.card-title { font-size: 13.5px; font-weight: 700; color: #111827; margin: 0; }

.badge { padding: 3px 10px; border-radius: 100px; font-size: 11px; font-weight: 600; }
.badge--gray   { background: #f3f4f6; color: #6b7280; }
.badge--green  { background: #dcfce7; color: #16a34a; }
.badge--yellow { background: #fef9c3; color: #92400e; }
.badge--red    { background: #fef2f2; color: #dc2626; }
.badge--orange { background: #fff7ed; color: #ea580c; }
.badge--blue   { background: #eff6ff; color: #2563eb; }

/* ── Steps (panduan) ────────────────────────────────────────── */
.steps { padding: 8px 20px 16px; display: flex; flex-direction: column; gap: 0; }
.step { display: flex; align-items: center; gap: 14px; padding: 14px 0; border-bottom: 1px solid #f8faf8; }
.step:last-child { border-bottom: none; }
.step__num { width: 28px; height: 28px; border-radius: 50%; background: #e5e7eb; color: #9ca3af; font-size: 12px; font-weight: 700; display: flex; align-items: center; justify-content: center; flex-shrink: 0; }
.step--active .step__num { background: #48AF4A; color: #fff; }
.step__body { flex: 1; }
.step__title { font-size: 13px; font-weight: 600; color: #111827; }
.step__desc { font-size: 11.5px; color: #9ca3af; margin-top: 2px; }

/* ── Empty state ────────────────────────────────────────────── */
.empty-state { display: flex; flex-direction: column; align-items: center; padding: 44px 24px; gap: 12px; text-align: center; }
.empty-state__icon { width: 72px; height: 72px; background: #f9fafb; border-radius: 50%; display: flex; align-items: center; justify-content: center; }
.empty-state p { font-size: 13px; color: #9ca3af; line-height: 1.7; margin: 0; }

/* ── Buttons ────────────────────────────────────────────────── */
.btn-green { background: #48AF4A; color: #fff; border: none; border-radius: 9px; padding: 10px 22px; font-size: 13px; font-weight: 600; font-family: "Poppins", sans-serif; cursor: pointer; transition: background 0.15s; }
.btn-green:hover { background: #3d9e3f; }
.btn-green-sm { background: #48AF4A; color: #fff; border: none; border-radius: 8px; padding: 6px 14px; font-size: 12px; font-weight: 600; font-family: "Poppins", sans-serif; cursor: pointer; white-space: nowrap; }
.btn-green-sm:hover { background: #3d9e3f; }

/* ── Pengajuan status ───────────────────────────────────────── */
.status-detail { padding: 20px; display: flex; flex-direction: column; gap: 18px; }
.status-detail__row { display: grid; grid-template-columns: 1fr 1fr; gap: 14px; }
.status-detail__item {}
.status-detail__label { font-size: 11px; color: #9ca3af; margin-bottom: 3px; }
.status-detail__value { font-size: 13.5px; font-weight: 600; color: #111827; }
.status-catatan { background: #fffbeb; border: 1px solid #fde68a; border-radius: 9px; padding: 12px 14px; }
.status-catatan__title { font-size: 11.5px; font-weight: 700; color: #92400e; margin-bottom: 5px; }
.status-catatan__body { font-size: 13px; color: #78350f; line-height: 1.6; }
.status-timeline { display: flex; flex-direction: column; gap: 10px; padding-top: 4px; }
.status-timeline__item { display: flex; align-items: center; gap: 10px; font-size: 12.5px; color: #9ca3af; }
.status-timeline__item--done { color: #16a34a; font-weight: 600; }
.status-timeline__dot { width: 10px; height: 10px; border-radius: 50%; background: #e5e7eb; flex-shrink: 0; }
.status-timeline__item--done .status-timeline__dot { background: #16a34a; }
.status-timeline__date { margin-left: auto; font-size: 11px; color: #9ca3af; font-weight: 400; }

/* ── Form card ──────────────────────────────────────────────── */
.form-card { background: #fff; border-radius: 14px; border: 1px solid #e9f5e9; box-shadow: 0 1px 3px rgba(13,40,24,0.05); overflow: hidden; }
.form-header { background: linear-gradient(135deg, #0d2818 0%, #1a5c20 100%); padding: 18px 24px; }
.form-header__title { font-size: 15px; font-weight: 800; color: #fff; }
.form-header__sub { font-size: 11.5px; color: rgba(255,255,255,0.5); margin-top: 3px; }

/* ── Stepper ────────────────────────────────────────────────── */
.stepper { display: flex; align-items: flex-start; padding: 18px 24px 14px; border-bottom: 1px solid #f0faf0; gap: 0; }
.stepper-item { flex: 1; display: flex; align-items: flex-start; }
.stepper-col { display: flex; flex-direction: column; align-items: center; gap: 6px; }
.stepper-dot { width: 28px; height: 28px; border-radius: 50%; background: #f3f4f6; color: #9ca3af; font-size: 12px; font-weight: 700; display: flex; align-items: center; justify-content: center; flex-shrink: 0; transition: all 0.2s; border: 2px solid transparent; }
.stepper-dot--active { background: #48AF4A; color: #fff; border-color: #48AF4A; box-shadow: 0 0 0 4px rgba(72,175,74,0.15); }
.stepper-dot--done   { background: #0d2818; color: #86efac; border-color: #0d2818; }
.stepper-label { font-size: 10px; font-weight: 600; color: #d1d5db; white-space: nowrap; }
.stepper-label--active { color: #48AF4A; }
.stepper-label--done   { color: #374151; }
.stepper-line { flex: 1; height: 2px; background: #e5e7eb; margin: 13px 6px 0; transition: background 0.2s; }
.stepper-line--done { background: #48AF4A; }

/* ── Step body ──────────────────────────────────────────────── */
.step-body { padding: 20px 24px; display: flex; flex-direction: column; gap: 14px; }
.step-head { margin-bottom: 2px; }
.step-title { font-size: 14px; font-weight: 800; color: #0d2818; }
.step-sub   { font-size: 11.5px; color: #9ca3af; margin-top: 2px; }

/* ── Fields ─────────────────────────────────────────────────── */
.form-grid-2 { display: grid; grid-template-columns: 1fr 1fr; gap: 12px; }
.field { display: flex; flex-direction: column; gap: 5px; }
.field-label { font-size: 11.5px; font-weight: 600; color: #374151; }
.req { color: #e11d48; }
.inp { border: 1.5px solid #e5e7eb; border-radius: 9px; padding: 9px 12px; font-size: 13px; font-family: "Poppins", sans-serif; color: #111827; outline: none; width: 100%; background: #fff; transition: border-color 0.15s; appearance: none; }
.inp:focus { border-color: #48AF4A; box-shadow: 0 0 0 3px rgba(72,175,74,0.10); }
textarea.inp { resize: vertical; }

/* ── Kategori picker ────────────────────────────────────────── */
.kategori-grid { display: flex; flex-direction: column; gap: 7px; }
.kategori-btn { display: flex; align-items: center; gap: 11px; padding: 11px 14px; border-radius: 9px; border: 1.5px solid #e5e7eb; background: #fff; font-size: 13px; font-family: "Poppins", sans-serif; color: #374151; cursor: pointer; text-align: left; transition: all 0.15s; }
.kategori-btn--active { border-color: #48AF4A; background: #f0fdf4; color: #0d2818; font-weight: 600; }
.kategori-radio { width: 14px; height: 14px; border-radius: 50%; border: 2px solid #d1d5db; flex-shrink: 0; transition: all 0.15s; }
.kategori-btn--active .kategori-radio { border-color: #48AF4A; background: #48AF4A; box-shadow: inset 0 0 0 3px #fff; }

/* ── Info box ───────────────────────────────────────────────── */
.info-box { background: #eff6ff; border: 1px solid #bfdbfe; border-radius: 9px; padding: 11px 14px; font-size: 12px; color: #1d4ed8; line-height: 1.6; }

/* ── Dokumen upload ─────────────────────────────────────────── */
.dok-grid { display: flex; flex-direction: column; gap: 12px; }
.dok-item { display: flex; flex-direction: column; gap: 5px; }
.dok-item__header { display: flex; align-items: center; gap: 7px; }
.dok-item__label { font-size: 12px; font-weight: 600; color: #374151; }
.badge-required { background: #fef2f2; color: #dc2626; border-radius: 100px; padding: 2px 8px; font-size: 10px; font-weight: 700; }
.badge-optional { background: #f3f4f6; color: #6b7280;  border-radius: 100px; padding: 2px 8px; font-size: 10px; font-weight: 600; }

.dropzone { border: 2px dashed #d1d5db; border-radius: 9px; padding: 12px 14px; cursor: pointer; transition: all 0.15s; background: #fafafa; }
.dropzone:hover, .dropzone--drag { border-color: #48AF4A; background: #f0fdf4; }
.dropzone--done { border-style: solid; border-color: #48AF4A; background: #f0fdf4; }
.dropzone-empty { display: flex; align-items: center; gap: 10px; }
.dropzone-empty__text { font-size: 12.5px; color: #374151; font-weight: 500; }
.dropzone-empty__hint { font-size: 11px; color: #9ca3af; margin-top: 1px; }
.dropzone-done { display: flex; align-items: center; gap: 10px; }
.dropzone-done__info { flex: 1; }
.dropzone-done__name { font-size: 12.5px; font-weight: 600; color: #0d2818; }
.dropzone-done__size { font-size: 11px; color: #6b7280; }
.dropzone-done__remove { background: none; border: none; color: #9ca3af; cursor: pointer; padding: 4px; border-radius: 4px; display: flex; align-items: center; }
.dropzone-done__remove:hover { background: #fee2e2; color: #dc2626; }

/* ── Review ─────────────────────────────────────────────────── */
.review-section { background: #fafafa; border-radius: 11px; padding: 14px 16px; display: flex; flex-direction: column; gap: 10px; }
.review-section__title { font-size: 11px; font-weight: 700; color: #0d2818; text-transform: uppercase; letter-spacing: 0.06em; padding-bottom: 8px; border-bottom: 1px solid #e5e7eb; }
.review-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 8px; }
.review-row { display: flex; flex-direction: column; gap: 2px; }
.review-row--full { grid-column: 1 / -1; }
.review-row__label { font-size: 10.5px; color: #9ca3af; }
.review-row__value { font-size: 12.5px; font-weight: 600; color: #111827; }
.dok-status-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 6px; }
.dok-status { display: flex; align-items: center; gap: 7px; font-size: 12px; color: #6b7280; }
.dok-status--ok   { color: #16a34a; font-weight: 600; }
.dok-status--miss { color: #dc2626; }
.dok-missing { font-size: 10px; color: #dc2626; font-weight: 600; margin-left: auto; }

.review-note { background: #fffbeb; border: 1px solid #fde68a; border-radius: 9px; padding: 11px 14px; font-size: 12px; color: #78350f; line-height: 1.6; }
.alert-warn  { background: #fef2f2; border: 1px solid #fecaca; border-radius: 9px; padding: 11px 14px; font-size: 12px; color: #dc2626; }

/* ── Form footer ────────────────────────────────────────────── */
.form-footer { padding: 14px 24px; border-top: 1px solid #f0faf0; display: flex; align-items: center; justify-content: space-between; background: #fafafa; }
.step-counter { font-size: 11px; color: #9ca3af; }
.btn-back { background: #fff; border: 1.5px solid #e5e7eb; color: #374151; border-radius: 9px; padding: 9px 20px; font-size: 13px; font-weight: 600; font-family: "Poppins", sans-serif; cursor: pointer; transition: all 0.15s; }
.btn-back:hover { border-color: #48AF4A; color: #0d2818; }
.btn-next { background: #48AF4A; color: #fff; border: none; border-radius: 9px; padding: 9px 24px; font-size: 13px; font-weight: 700; font-family: "Poppins", sans-serif; cursor: pointer; transition: background 0.15s; display: flex; align-items: center; gap: 7px; }
.btn-next:hover:not(:disabled) { background: #3d9e3f; }
.btn-next:disabled { opacity: 0.6; cursor: not-allowed; }
.btn-submit { background: #0d2818; }
.btn-submit:hover:not(:disabled) { background: #1a5c20; }

.spinner { width: 14px; height: 14px; border: 2px solid rgba(255,255,255,0.3); border-top-color: #fff; border-radius: 50%; animation: spin 0.7s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }

/* ── Responsive ─────────────────────────────────────────────── */
@media (max-width: 700px) {
  .stats-row { grid-template-columns: 1fr 1fr; }
  .welcome-date__day, .welcome-date__info { display: none; }
  .form-grid-2 { grid-template-columns: 1fr; }
  .review-grid  { grid-template-columns: 1fr; }
  .dok-status-grid { grid-template-columns: 1fr; }
  .status-detail__row { grid-template-columns: 1fr; }
}
@media (max-width: 420px) {
  .stats-row { grid-template-columns: 1fr; }
}
</style>
