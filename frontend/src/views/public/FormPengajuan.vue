<template>
  <div class="fp-page">

    <!-- ── SUCCESS SCREEN ── -->
    <div v-if="submitted" class="fp-success">
      <div class="fp-success__card">
        <div class="fp-success__check">
          <svg width="32" height="32" viewBox="0 0 24 24" fill="none"><path d="M5 12l5 5L20 7" stroke="#fff" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"/></svg>
        </div>
        <h2>Pengajuan Berhasil Dikirim</h2>
        <p>
          Terima kasih, <strong>{{ diri.nama_lengkap }}</strong>. Tim HRD akan meninjau berkas Anda dan
          mengirimkan akun login ke <strong>{{ diri.email }}</strong> dalam 3–5 hari kerja.
        </p>
        <div class="fp-success__meta">
          <div class="fp-success__meta-item">
            <svg width="15" height="15" viewBox="0 0 24 24" fill="none"><circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2"/><path d="M12 6v6l4 2" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>
            Proses verifikasi 3–5 hari kerja
          </div>
          <div class="fp-success__meta-item">
            <svg width="15" height="15" viewBox="0 0 24 24" fill="none"><path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z" stroke="currentColor" stroke-width="2"/><polyline points="22,6 12,13 2,6" stroke="currentColor" stroke-width="2"/></svg>
            Akun dikirim via email setelah diterima
          </div>
          <div class="fp-success__meta-item">
            <svg width="15" height="15" viewBox="0 0 24 24" fill="none"><path d="M9 12l2 2 4-4" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/><path d="M21 12c0 4.97-4.03 9-9 9s-9-4.03-9-9 4.03-9 9-9 9 4.03 9 9z" stroke="currentColor" stroke-width="2"/></svg>
            {{ docsUploaded }} dokumen berhasil diupload
          </div>
        </div>
        <router-link to="/" class="fp-success__btn">Kembali ke Beranda</router-link>
      </div>
    </div>

    <!-- ── MAIN FORM ── -->
    <template v-else>

    <!-- Navbar -->
    <nav class="fp-nav">
      <div class="fp-nav__inner">
        <router-link to="/" class="fp-nav__brand">
          <img src="/logotel.png" alt="TELPP" class="fp-nav__logo" />
          <span>e-Magang <strong>PT TELPP</strong></span>
        </router-link>
        <div class="fp-nav__right">
          <router-link to="/" class="fp-nav__link">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none"><path d="M19 12H5M12 5l-7 7 7 7" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>
            Beranda
          </router-link>
          <router-link to="/login" class="fp-nav__btn">Sudah punya akun? Masuk</router-link>
        </div>
      </div>
    </nav>

    <div class="fp-wrap">
      <!-- Header hijau -->
      <div class="fp-header">
        <div class="fp-header__inner">
          <div class="fp-header__title">Formulir Pengajuan Magang</div>
          <div class="fp-header__sub">PT TanjungEnim Lestari Pulp and Paper</div>
        </div>
      </div>

      <!-- Body -->
      <div class="fp-body">

        <!-- Stepper -->
        <div class="fp-stepper">
          <template v-for="(s, i) in STEPS" :key="i">
            <div class="fp-stepper__col">
              <div
                class="fp-stepper__circle"
                :class="{
                  'fp-stepper__circle--active': step === i,
                  'fp-stepper__circle--done': step > i,
                }"
              >
                <svg v-if="step > i" width="13" height="13" viewBox="0 0 24 24" fill="none"><path d="M5 12l5 5L20 7" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"/></svg>
                <span v-else>{{ i + 1 }}</span>
              </div>
              <div
                class="fp-stepper__label"
                :class="{
                  'fp-stepper__label--active': step === i,
                  'fp-stepper__label--done': step > i,
                }"
              >{{ s }}</div>
            </div>
            <div
              v-if="i < STEPS.length - 1"
              class="fp-stepper__line"
              :class="{ 'fp-stepper__line--done': step > i }"
            ></div>
          </template>
        </div>

        <!-- ── STEP 1: Data Diri ── -->
        <div v-if="step === 0" class="fp-section">
          <div class="fp-section__head">
            <div class="fp-section__title">Data Diri</div>
            <div class="fp-section__sub">Lengkapi informasi pribadi sesuai KTP</div>
          </div>
          <div class="fp-grid">
            <div class="fp-field">
              <label>Nama Lengkap <span class="req">*</span></label>
              <input v-model="diri.nama_lengkap" type="text" placeholder="Sesuai KTP" />
            </div>
            <div class="fp-field">
              <label>Jenis Kelamin <span class="req">*</span></label>
              <select v-model="diri.jenis_kelamin">
                <option value="">-- Pilih --</option>
                <option value="laki_laki">Laki-laki</option>
                <option value="perempuan">Perempuan</option>
              </select>
            </div>
            <div class="fp-field">
              <label>Tempat Lahir <span class="req">*</span></label>
              <input v-model="diri.tempat_lahir" type="text" placeholder="Kota kelahiran" />
            </div>
            <div class="fp-field">
              <label>Tanggal Lahir <span class="req">*</span></label>
              <input v-model="diri.tanggal_lahir" type="date" :max="todayMax" />
            </div>
            <div class="fp-field">
              <label>Nomor HP / WhatsApp <span class="req">*</span></label>
              <input v-model="diri.no_hp" type="tel" placeholder="08xxxxxxxxxx" @input="filterNoHp" inputmode="numeric" maxlength="15" />
            </div>
            <div class="fp-field">
              <label>Email Aktif <span class="req">*</span></label>
              <input v-model="diri.email" type="email" placeholder="nama@email.com" />
            </div>
          </div>
          <div class="fp-field fp-field--full">
            <label>Alamat Lengkap <span class="req">*</span></label>
            <textarea v-model="diri.alamat" rows="3" placeholder="Jalan, nomor, kelurahan, kecamatan, kota, provinsi"></textarea>
          </div>
        </div>

        <!-- ── STEP 2: Data Akademis ── -->
        <div v-else-if="step === 1" class="fp-section">
          <div class="fp-section__head">
            <div class="fp-section__title">Informasi Akademis</div>
            <div class="fp-section__sub">Data institusi dan program studi Anda</div>
          </div>

          <!-- Dropdown kategori -->
          <div class="fp-field fp-field--full" style="margin-bottom:20px">
            <label>Kategori Magang <span class="req">*</span></label>
            <select v-model="akademis.kategori_magang" @change="onKategoriChange">
              <option value="">-- Pilih Kategori --</option>
              <option value="smk">SMK / Sekolah Menengah Kejuruan</option>
              <option value="d3">D3 (Diploma Tiga)</option>
              <option value="s1">S1 (Sarjana)</option>
              <option value="s2">S2 (Magister)</option>
              <option value="penelitian">Penelitian / Riset</option>
              <option value="lainnya">Lainnya</option>
            </select>
          </div>

          <!-- Lainnya: input bebas -->
          <div v-if="isLainnya" class="fp-field fp-field--full" style="margin-bottom:20px">
            <label>Sebutkan Kategori Anda <span class="req">*</span></label>
            <input v-model="kategoriLainnya" type="text" placeholder="Contoh: Magang Industri, PKL Teknik, dll." />
          </div>

          <!-- Form dinamis sesuai kategori -->
          <div class="fp-grid" v-if="akademis.kategori_magang">

            <!-- Asal Institusi (semua) -->
            <div class="fp-field">
              <label>{{ labelInstitusi }} <span class="req">*</span></label>
              <input v-model="akademis.asal_institusi" type="text" :placeholder="placeholderInstitusi" />
            </div>

            <!-- Jurusan / Bidang (semua) -->
            <div class="fp-field">
              <label>{{ labelJurusan }} <span class="req">*</span></label>
              <input v-model="akademis.jurusan" type="text" :placeholder="placeholderJurusan" />
            </div>

            <!-- Kelas — hanya SMK -->
            <div v-if="isSmk" class="fp-field">
              <label>Kelas <span class="req">*</span></label>
              <select v-model="akademis.kelas_semester">
                <option value="">-- Pilih Kelas --</option>
                <option value="Kelas X">Kelas X</option>
                <option value="Kelas XI">Kelas XI</option>
                <option value="Kelas XII">Kelas XII</option>
              </select>
            </div>

            <!-- Semester — hanya D3/S1/S2 -->
            <div v-if="isPT" class="fp-field">
              <label>Semester <span class="req">*</span></label>
              <select v-model="akademis.kelas_semester">
                <option value="">-- Pilih Semester --</option>
                <option v-for="n in semesterMax" :key="n" :value="`Semester ${n}`">Semester {{ n }}</option>
              </select>
            </div>

            <!-- Nomor Induk — SMK / PT / Lainnya -->
            <div v-if="!isPenelitian" class="fp-field">
              <label>{{ labelNomorInduk }} <span class="req">*</span></label>
              <input v-model="akademis.nomor_induk" type="text" :placeholder="placeholderNomorInduk" />
            </div>

            <!-- Judul Penelitian — hanya Penelitian (full width) -->
            <div v-if="isPenelitian" class="fp-field fp-field--full">
              <label>Judul / Topik Penelitian <span class="req">*</span></label>
              <input v-model="akademis.kelas_semester" type="text" placeholder="Contoh: Analisis Efisiensi Produksi Kertas..." />
            </div>

            <!-- Nomor Identitas — hanya Penelitian -->
            <div v-if="isPenelitian" class="fp-field">
              <label>Nomor Identitas (NIK / NIM / NIP) <span class="req">*</span></label>
              <input v-model="akademis.nomor_induk" type="text" placeholder="Nomor identitas resmi" />
            </div>
          </div>

          <div class="fp-notice" v-if="akademis.kategori_magang">
            <svg width="15" height="15" viewBox="0 0 24 24" fill="none"><circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2"/><path d="M12 8v4M12 16h.01" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>
            Kategori magang menentukan dokumen yang perlu diunggah pada langkah berikutnya.
          </div>
        </div>

        <!-- ── STEP 3: Upload Dokumen ── -->
        <div v-else-if="step === 2" class="fp-section">
          <div class="fp-section__head">
            <div class="fp-section__title">Upload Dokumen</div>
            <div class="fp-section__sub">Unggah semua berkas yang diperlukan. Maks 10 MB per file.</div>
          </div>

          <div class="fp-doc-list">
            <!-- 1. Proposal Magang / Penelitian -->
            <div class="fp-doc">
              <div class="fp-doc__label">
                {{ labelProposal }} <span class="req">*</span>
              </div>
              <div
                class="fp-dropzone"
                :class="{ 'fp-dropzone--done': dokumenFiles.proposal_magang }"
                @click="triggerPick('proposal_magang')"
              >
                <template v-if="!dokumenFiles.proposal_magang">
                  <svg width="18" height="18" viewBox="0 0 24 24" fill="none" class="fp-dropzone__icon"><path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4" stroke="currentColor" stroke-width="1.8" stroke-linecap="round"/><polyline points="17 8 12 3 7 8" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"/><line x1="12" y1="3" x2="12" y2="15" stroke="currentColor" stroke-width="1.8" stroke-linecap="round"/></svg>
                  <span class="fp-dropzone__text">Klik atau drag file ke sini</span>
                  <span class="fp-dropzone__hint">PDF atau Word, maks 10 MB</span>
                </template>
                <template v-else>
                  <svg width="16" height="16" viewBox="0 0 24 24" fill="none"><path d="M5 12l5 5L20 7" stroke="#16a34a" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"/></svg>
                  <span class="fp-dropzone__filename">{{ dokumenFiles.proposal_magang.name }}</span>
                  <button class="fp-dropzone__remove" @click.stop="removeFile('proposal_magang')" type="button">Hapus</button>
                </template>
              </div>
              <input ref="pick_proposal_magang" type="file" accept=".pdf,.doc,.docx" style="display:none" @change="onFilePick('proposal_magang', $event)" />
            </div>

            <!-- 2. Surat Pengantar -->
            <div class="fp-doc">
              <div class="fp-doc__label">
                {{ labelSuratPengantar }} <span class="req">*</span>
              </div>
              <div
                class="fp-dropzone"
                :class="{ 'fp-dropzone--done': dokumenFiles.surat_pengantar }"
                @click="triggerPick('surat_pengantar')"
              >
                <template v-if="!dokumenFiles.surat_pengantar">
                  <svg width="18" height="18" viewBox="0 0 24 24" fill="none" class="fp-dropzone__icon"><path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4" stroke="currentColor" stroke-width="1.8" stroke-linecap="round"/><polyline points="17 8 12 3 7 8" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"/><line x1="12" y1="3" x2="12" y2="15" stroke="currentColor" stroke-width="1.8" stroke-linecap="round"/></svg>
                  <span class="fp-dropzone__text">Klik atau drag file ke sini</span>
                  <span class="fp-dropzone__hint">PDF atau Word, maks 10 MB</span>
                </template>
                <template v-else>
                  <svg width="16" height="16" viewBox="0 0 24 24" fill="none"><path d="M5 12l5 5L20 7" stroke="#16a34a" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"/></svg>
                  <span class="fp-dropzone__filename">{{ dokumenFiles.surat_pengantar.name }}</span>
                  <button class="fp-dropzone__remove" @click.stop="removeFile('surat_pengantar')" type="button">Hapus</button>
                </template>
              </div>
              <input ref="pick_surat_pengantar" type="file" accept=".pdf,.doc,.docx" style="display:none" @change="onFilePick('surat_pengantar', $event)" />
            </div>

            <!-- 3. KTP -->
            <div class="fp-doc">
              <div class="fp-doc__label">
                KTP / Kartu Identitas <span class="req">*</span>
              </div>
              <div
                class="fp-dropzone"
                :class="{ 'fp-dropzone--done': dokumenFiles.ktp }"
                @click="triggerPick('ktp')"
              >
                <template v-if="!dokumenFiles.ktp">
                  <svg width="18" height="18" viewBox="0 0 24 24" fill="none" class="fp-dropzone__icon"><path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4" stroke="currentColor" stroke-width="1.8" stroke-linecap="round"/><polyline points="17 8 12 3 7 8" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"/><line x1="12" y1="3" x2="12" y2="15" stroke="currentColor" stroke-width="1.8" stroke-linecap="round"/></svg>
                  <span class="fp-dropzone__text">Klik atau drag file ke sini</span>
                  <span class="fp-dropzone__hint">JPG, PNG, atau PDF</span>
                </template>
                <template v-else>
                  <svg width="16" height="16" viewBox="0 0 24 24" fill="none"><path d="M5 12l5 5L20 7" stroke="#16a34a" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"/></svg>
                  <span class="fp-dropzone__filename">{{ dokumenFiles.ktp.name }}</span>
                  <button class="fp-dropzone__remove" @click.stop="removeFile('ktp')" type="button">Hapus</button>
                </template>
              </div>
              <input ref="pick_ktp" type="file" accept=".jpg,.jpeg,.png,.pdf" style="display:none" @change="onFilePick('ktp', $event)" />
            </div>

            <!-- 4. KTM / Kartu Pelajar (wajib SMK & PT, opsional lainnya) -->
            <div class="fp-doc">
              <div class="fp-doc__label">
                {{ labelKtm }}
                <span v-if="isKtmRequired" class="req">*</span>
                <span v-else class="fp-badge fp-badge--optional">Opsional</span>
              </div>
              <div
                class="fp-dropzone"
                :class="{ 'fp-dropzone--done': dokumenFiles.ktm }"
                @click="triggerPick('ktm')"
              >
                <template v-if="!dokumenFiles.ktm">
                  <svg width="18" height="18" viewBox="0 0 24 24" fill="none" class="fp-dropzone__icon"><path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4" stroke="currentColor" stroke-width="1.8" stroke-linecap="round"/><polyline points="17 8 12 3 7 8" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"/><line x1="12" y1="3" x2="12" y2="15" stroke="currentColor" stroke-width="1.8" stroke-linecap="round"/></svg>
                  <span class="fp-dropzone__text">Klik atau drag file ke sini</span>
                  <span class="fp-dropzone__hint">JPG, PNG, atau PDF</span>
                </template>
                <template v-else>
                  <svg width="16" height="16" viewBox="0 0 24 24" fill="none"><path d="M5 12l5 5L20 7" stroke="#16a34a" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"/></svg>
                  <span class="fp-dropzone__filename">{{ dokumenFiles.ktm.name }}</span>
                  <button class="fp-dropzone__remove" @click.stop="removeFile('ktm')" type="button">Hapus</button>
                </template>
              </div>
              <input ref="pick_ktm" type="file" accept=".jpg,.jpeg,.png,.pdf" style="display:none" @change="onFilePick('ktm', $event)" />
            </div>

            <!-- 5. Pasfoto -->
            <div class="fp-doc">
              <div class="fp-doc__label">
                Pasfoto 3x4 <span class="req">*</span>
              </div>
              <div
                class="fp-dropzone"
                :class="{ 'fp-dropzone--done': dokumenFiles.pasfoto }"
                @click="triggerPick('pasfoto')"
              >
                <template v-if="!dokumenFiles.pasfoto">
                  <svg width="18" height="18" viewBox="0 0 24 24" fill="none" class="fp-dropzone__icon"><path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4" stroke="currentColor" stroke-width="1.8" stroke-linecap="round"/><polyline points="17 8 12 3 7 8" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"/><line x1="12" y1="3" x2="12" y2="15" stroke="currentColor" stroke-width="1.8" stroke-linecap="round"/></svg>
                  <span class="fp-dropzone__text">Klik atau drag file ke sini</span>
                  <span class="fp-dropzone__hint">JPG atau PNG, latar merah/biru</span>
                </template>
                <template v-else>
                  <svg width="16" height="16" viewBox="0 0 24 24" fill="none"><path d="M5 12l5 5L20 7" stroke="#16a34a" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"/></svg>
                  <span class="fp-dropzone__filename">{{ dokumenFiles.pasfoto.name }}</span>
                  <button class="fp-dropzone__remove" @click.stop="removeFile('pasfoto')" type="button">Hapus</button>
                </template>
              </div>
              <input ref="pick_pasfoto" type="file" accept=".jpg,.jpeg,.png" style="display:none" @change="onFilePick('pasfoto', $event)" />
            </div>

            <!-- 6. BPJS / KIS -->
            <div class="fp-doc">
              <div class="fp-doc__label">
                BPJS / KIS <span class="req">*</span>
              </div>
              <div
                class="fp-dropzone"
                :class="{ 'fp-dropzone--done': dokumenFiles.bpjs_kis }"
                @click="triggerPick('bpjs_kis')"
              >
                <template v-if="!dokumenFiles.bpjs_kis">
                  <svg width="18" height="18" viewBox="0 0 24 24" fill="none" class="fp-dropzone__icon"><path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4" stroke="currentColor" stroke-width="1.8" stroke-linecap="round"/><polyline points="17 8 12 3 7 8" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"/><line x1="12" y1="3" x2="12" y2="15" stroke="currentColor" stroke-width="1.8" stroke-linecap="round"/></svg>
                  <span class="fp-dropzone__text">Klik atau drag file ke sini</span>
                  <span class="fp-dropzone__hint">Foto kartu atau PDF</span>
                </template>
                <template v-else>
                  <svg width="16" height="16" viewBox="0 0 24 24" fill="none"><path d="M5 12l5 5L20 7" stroke="#16a34a" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"/></svg>
                  <span class="fp-dropzone__filename">{{ dokumenFiles.bpjs_kis.name }}</span>
                  <button class="fp-dropzone__remove" @click.stop="removeFile('bpjs_kis')" type="button">Hapus</button>
                </template>
              </div>
              <input ref="pick_bpjs_kis" type="file" accept=".jpg,.jpeg,.png,.pdf" style="display:none" @change="onFilePick('bpjs_kis', $event)" />
            </div>
          </div>
        </div>

        <!-- ── STEP 4: Review & Kirim ── -->
        <div v-else-if="step === 3" class="fp-section">
          <div class="fp-section__head">
            <div class="fp-section__title">Review dan Kirim</div>
            <div class="fp-section__sub">Periksa kembali data sebelum mengirim pengajuan</div>
          </div>

          <div class="fp-review-block">
            <div class="fp-review-block__head">Data Diri</div>
            <div class="fp-review-grid">
              <div class="fp-review-cell"><span>Nama Lengkap</span><strong>{{ diri.nama_lengkap }}</strong></div>
              <div class="fp-review-cell"><span>Jenis Kelamin</span><strong>{{ diri.jenis_kelamin === 'laki_laki' ? 'Laki-laki' : 'Perempuan' }}</strong></div>
              <div class="fp-review-cell"><span>Tempat Lahir</span><strong>{{ diri.tempat_lahir }}</strong></div>
              <div class="fp-review-cell"><span>Tanggal Lahir</span><strong>{{ formatDate(diri.tanggal_lahir) }}</strong></div>
              <div class="fp-review-cell"><span>No. HP</span><strong>{{ diri.no_hp }}</strong></div>
              <div class="fp-review-cell"><span>Email</span><strong>{{ diri.email }}</strong></div>
              <div class="fp-review-cell fp-review-cell--full"><span>Alamat</span><strong>{{ diri.alamat }}</strong></div>
            </div>
          </div>

          <div class="fp-review-block">
            <div class="fp-review-block__head">Data Akademis</div>
            <div class="fp-review-grid">
              <div class="fp-review-cell"><span>Kategori</span><strong>{{ labelKategori(akademis.kategori_magang) }}</strong></div>
              <div class="fp-review-cell"><span>Institusi</span><strong>{{ akademis.asal_institusi }}</strong></div>
              <div class="fp-review-cell"><span>{{ labelJurusan }}</span><strong>{{ akademis.jurusan }}</strong></div>
              <div class="fp-review-cell" v-if="!isPenelitian"><span>{{ isSmk ? 'Kelas' : 'Semester' }}</span><strong>{{ akademis.kelas_semester }}</strong></div>
              <div class="fp-review-cell" v-if="isPenelitian"><span>Judul Penelitian</span><strong>{{ akademis.kelas_semester }}</strong></div>
              <div class="fp-review-cell"><span>{{ labelNomorInduk }}</span><strong>{{ akademis.nomor_induk }}</strong></div>
            </div>
          </div>

          <div class="fp-review-block">
            <div class="fp-review-block__head">Dokumen</div>
            <div class="fp-review-docs">
              <div v-for="d in docsReviewList" :key="d.jenis" class="fp-review-doc">
                <svg v-if="d.done" width="15" height="15" viewBox="0 0 24 24" fill="none"><path d="M9 12l2 2 4-4" stroke="#16a34a" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"/><circle cx="12" cy="12" r="10" stroke="#16a34a" stroke-width="2"/></svg>
                <svg v-else-if="d.required" width="15" height="15" viewBox="0 0 24 24" fill="none"><circle cx="12" cy="12" r="10" stroke="#ef4444" stroke-width="2"/><path d="M12 8v4M12 16h.01" stroke="#ef4444" stroke-width="2" stroke-linecap="round"/></svg>
                <svg v-else width="15" height="15" viewBox="0 0 24 24" fill="none"><circle cx="12" cy="12" r="10" stroke="#9ca3af" stroke-width="2"/><path d="M12 8v4M12 16h.01" stroke="#9ca3af" stroke-width="2" stroke-linecap="round"/></svg>
                <span :class="d.done ? 'fp-review-doc--ok' : d.required ? 'fp-review-doc--miss' : 'fp-review-doc--skip'">
                  {{ d.label }}{{ !d.required && !d.done ? ' (opsional)' : '' }}
                </span>
              </div>
            </div>
          </div>

          <!-- Upload progress -->
          <div v-if="uploadStatus" class="fp-progress">
            <div class="fp-progress__bar"><div class="fp-progress__fill" :style="{ width: uploadPct + '%' }"></div></div>
            <div class="fp-progress__text">{{ uploadStatus }}</div>
          </div>
        </div>

        <!-- Error banner -->
        <div v-if="submitError" class="fp-error">{{ submitError }}</div>

        <!-- ── FOOTER / AKSI ── -->
        <div class="fp-footer">
          <button v-if="step > 0" class="fp-btn-back" @click="step--" :disabled="submitting">Kembali</button>
          <div v-else class="fp-footer__spacer"></div>

          <span class="fp-footer__counter">Langkah {{ step + 1 }} dari {{ STEPS.length }}</span>

          <button
            v-if="step < 3"
            class="fp-btn-next"
            @click="nextStep"
          >Lanjut</button>
          <button
            v-else
            class="fp-btn-next"
            :disabled="submitting"
            @click="submitForm"
          >
            <span v-if="submitting" class="fp-spinner"></span>
            {{ submitting ? 'Mengirim...' : 'Kirim Pengajuan' }}
          </button>
        </div>
      </div>
    </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from "vue";
import api from "@/lib/api";

const STEPS = ["Data Diri", "Akademis", "Dokumen", "Review & Kirim"];

const todayMax = new Date().toISOString().slice(0, 10);

const step = ref(0);
const submitted = ref(false);
const submitting = ref(false);
const submitError = ref<string | null>(null);
const uploadStatus = ref("");
const uploadPct = ref(0);
const docsUploaded = ref(0);

// Input bebas untuk kategori "Lainnya"
const kategoriLainnya = ref("");

const diri = reactive({
  nama_lengkap: "",
  tempat_lahir: "",
  tanggal_lahir: "",
  jenis_kelamin: "",
  no_hp: "",
  email: "",
  alamat: "",
});

const akademis = reactive({
  kategori_magang: "",
  asal_institusi: "",
  jurusan: "",
  kelas_semester: "",
  nomor_induk: "",
});

const dokumenFiles = reactive<Record<string, File | null>>({
  proposal_magang: null,
  surat_pengantar: null,
  ktp: null,
  ktm: null,
  pasfoto: null,
  bpjs_kis: null,
});

// ── Computed helpers kategori ──
const isSmk       = computed(() => akademis.kategori_magang === "smk");
const isPT        = computed(() => ["d3", "s1", "s2"].includes(akademis.kategori_magang));
const isPenelitian = computed(() => akademis.kategori_magang === "penelitian");
const isLainnya   = computed(() => akademis.kategori_magang === "lainnya");

// Max semester sesuai jenjang
const semesterMax = computed(() => {
  if (akademis.kategori_magang === "d3") return 6;
  if (akademis.kategori_magang === "s2") return 4;
  return 8; // s1 default
});

// Label dinamis
const labelInstitusi = computed(() => {
  if (isSmk.value) return "Nama Sekolah";
  if (isPenelitian.value) return "Asal Institusi / Lembaga";
  if (isLainnya.value) return "Nama Institusi / Organisasi";
  return "Nama Universitas / Politeknik";
});

const placeholderInstitusi = computed(() => {
  if (isSmk.value) return "Nama lengkap sekolah";
  if (isPenelitian.value) return "Nama universitas / lembaga riset";
  if (isLainnya.value) return "Nama institusi atau organisasi";
  return "Nama universitas atau politeknik";
});

const labelJurusan = computed(() => {
  if (isSmk.value) return "Kompetensi Keahlian";
  if (isPenelitian.value) return "Bidang Penelitian";
  if (isLainnya.value) return "Bidang / Jurusan";
  return "Program Studi";
});

const placeholderJurusan = computed(() => {
  if (isSmk.value) return "Contoh: Teknik Komputer dan Jaringan";
  if (isPenelitian.value) return "Contoh: Lingkungan Hidup, Pulp & Paper";
  if (isLainnya.value) return "Bidang atau jurusan yang relevan";
  return "Contoh: Teknik Informatika";
});

const labelNomorInduk = computed(() => {
  if (isSmk.value) return "NISN";
  if (isPenelitian.value) return "Nomor Identitas (NIK / NIM / NIP)";
  return "NIM";
});

const placeholderNomorInduk = computed(() => {
  if (isSmk.value) return "Nomor Induk Siswa Nasional";
  if (isPenelitian.value) return "Nomor identitas resmi";
  return "Nomor Induk Mahasiswa";
});

// Reset fields saat kategori berubah
function onKategoriChange() {
  akademis.asal_institusi = "";
  akademis.jurusan = "";
  akademis.kelas_semester = "";
  akademis.nomor_induk = "";
  kategoriLainnya.value = "";
}

// Refs untuk hidden file inputs
const pick_proposal_magang = ref<HTMLInputElement | null>(null);
const pick_surat_pengantar  = ref<HTMLInputElement | null>(null);
const pick_ktp              = ref<HTMLInputElement | null>(null);
const pick_ktm              = ref<HTMLInputElement | null>(null);
const pick_pasfoto          = ref<HTMLInputElement | null>(null);
const pick_bpjs_kis         = ref<HTMLInputElement | null>(null);

function triggerPick(jenis: string) {
  const refs: Record<string, any> = {
    proposal_magang: pick_proposal_magang,
    surat_pengantar: pick_surat_pengantar,
    ktp: pick_ktp,
    ktm: pick_ktm,
    pasfoto: pick_pasfoto,
    bpjs_kis: pick_bpjs_kis,
  };
  refs[jenis]?.value?.click();
}

function onFilePick(jenis: string, event: Event) {
  const input = event.target as HTMLInputElement;
  dokumenFiles[jenis] = input.files?.[0] ?? null;
  input.value = "";
}

function removeFile(jenis: string) {
  dokumenFiles[jenis] = null;
}

// KTM wajib hanya untuk SMK dan PT (D3/S1/S2)
const isKtmRequired = computed(() => isSmk.value || isPT.value);

// Label dinamis dokumen
const labelProposal = computed(() =>
  isPenelitian.value ? "Proposal Penelitian" : "Proposal Magang"
);
const labelSuratPengantar = computed(() => {
  if (isSmk.value) return "Surat Pengantar dari Sekolah";
  if (isPT.value)  return "Surat Pengantar dari Universitas";
  return "Surat Pengantar dari Institusi";
});
const labelKtm = computed(() =>
  isSmk.value ? "Kartu Pelajar" : "KTM (Kartu Tanda Mahasiswa)"
);

const docsReviewList = computed(() => [
  { jenis: "proposal_magang", label: labelProposal.value,       required: true,                done: !!dokumenFiles.proposal_magang },
  { jenis: "surat_pengantar", label: labelSuratPengantar.value, required: true,                done: !!dokumenFiles.surat_pengantar },
  { jenis: "ktp",             label: "KTP / Kartu Identitas",   required: true,                done: !!dokumenFiles.ktp },
  { jenis: "ktm",             label: labelKtm.value,            required: isKtmRequired.value, done: !!dokumenFiles.ktm },
  { jenis: "pasfoto",         label: "Pasfoto 3x4",             required: true,                done: !!dokumenFiles.pasfoto },
  { jenis: "bpjs_kis",        label: "BPJS / KIS",              required: true,                done: !!dokumenFiles.bpjs_kis },
]);

function formatDate(d: string) {
  if (!d) return "-";
  return new Date(d).toLocaleDateString("id-ID", { day: "2-digit", month: "long", year: "numeric" });
}

function labelKategori(v: string) {
  const MAP: Record<string, string> = {
    smk:        "SMK / Sekolah Menengah Kejuruan",
    d3:         "D3 (Diploma Tiga)",
    s1:         "S1 (Sarjana)",
    s2:         "S2 (Magister)",
    penelitian: "Penelitian / Riset",
    lainnya:    kategoriLainnya.value ? `Lainnya — ${kategoriLainnya.value}` : "Lainnya",
  };
  return MAP[v] ?? v;
}

// Map nilai frontend → enum backend
function backendKategori(v: string): string {
  return v;
}

function filterNoHp(e: Event) {
  const input = e.target as HTMLInputElement;
  const digits = input.value.replace(/\D/g, "");
  diri.no_hp = digits;
}

function validateStep0(): string | null {
  const nama = diri.nama_lengkap.trim();
  if (!nama)                          return "Nama lengkap wajib diisi";
  if (nama.length < 3)                return "Nama lengkap minimal 3 karakter";
  if (!/^[a-zA-Z\s'.,-]+$/.test(nama)) return "Nama lengkap hanya boleh berisi huruf";

  const tempat = diri.tempat_lahir.trim();
  if (!tempat)                        return "Tempat lahir wajib diisi";
  if (tempat.length < 3)              return "Tempat lahir minimal 3 karakter";
  if (!/^[a-zA-Z\s'.,-]+$/.test(tempat)) return "Tempat lahir hanya boleh berisi huruf";

  if (!diri.tanggal_lahir)            return "Tanggal lahir wajib diisi";
  const lahir = new Date(diri.tanggal_lahir);
  const now   = new Date();
  const usia  = now.getFullYear() - lahir.getFullYear();
  if (usia < 14)                      return "Usia minimal 14 tahun untuk mendaftar magang";
  if (usia > 70)                      return "Tanggal lahir tidak valid";

  if (!diri.jenis_kelamin)            return "Jenis kelamin wajib dipilih";

  const hp = diri.no_hp.replace(/\D/g, "");
  if (!hp)                            return "Nomor HP wajib diisi";
  if (!/^\d+$/.test(hp))              return "Nomor HP hanya boleh berisi angka";
  if (hp.length < 10 || hp.length > 15) return "Nomor HP harus 10–15 digit";
  if (!/^(08|628)/.test(hp))          return "Nomor HP harus diawali 08 atau 628";

  const email = diri.email.trim();
  if (!email)                         return "Email aktif wajib diisi";
  if (!/^[^\s@]+@[^\s@]+\.[^\s@]{2,}$/.test(email)) return "Format email tidak valid";

  const alamat = diri.alamat.trim();
  if (!alamat)                        return "Alamat lengkap wajib diisi";
  if (alamat.length < 20)             return "Alamat lengkap terlalu singkat (minimal 20 karakter)";

  return null;
}

function validateStep1(): string | null {
  if (!akademis.kategori_magang)       return "Kategori magang wajib dipilih";
  if (isLainnya.value && !kategoriLainnya.value.trim()) return "Sebutkan kategori Anda";
  if (isLainnya.value && kategoriLainnya.value.trim().length < 3) return "Kategori minimal 3 karakter";

  const institusi = akademis.asal_institusi.trim();
  if (!institusi)                      return `${labelInstitusi.value} wajib diisi`;
  if (institusi.length < 5)            return `${labelInstitusi.value} minimal 5 karakter`;

  const jurusan = akademis.jurusan.trim();
  if (!jurusan)                        return `${labelJurusan.value} wajib diisi`;
  if (jurusan.length < 3)              return `${labelJurusan.value} minimal 3 karakter`;

  if (isSmk.value && !akademis.kelas_semester)  return "Kelas wajib dipilih";
  if (isPT.value && !akademis.kelas_semester)   return "Semester wajib dipilih";
  if (isPenelitian.value && !akademis.kelas_semester.trim()) return "Judul / Topik penelitian wajib diisi";
  if (isPenelitian.value && akademis.kelas_semester.trim().length < 10) return "Judul penelitian minimal 10 karakter";

  const nomorInduk = akademis.nomor_induk.trim();
  if (!nomorInduk)                     return `${labelNomorInduk.value} wajib diisi`;
  if (nomorInduk.length < 5)           return `${labelNomorInduk.value} minimal 5 karakter`;
  if (!/^[a-zA-Z0-9\-./]+$/.test(nomorInduk)) return `${labelNomorInduk.value} hanya boleh berisi angka, huruf, atau tanda hubung`;

  return null;
}

function validateStep2(): string | null {
  if (!dokumenFiles.proposal_magang) return `${labelProposal.value} wajib diupload`;
  if (!dokumenFiles.surat_pengantar) return `${labelSuratPengantar.value} wajib diupload`;
  if (!dokumenFiles.ktp)             return "KTP / Kartu Identitas wajib diupload";
  if (isKtmRequired.value && !dokumenFiles.ktm) return `${labelKtm.value} wajib diupload`;
  if (!dokumenFiles.pasfoto)         return "Pasfoto wajib diupload";
  if (!dokumenFiles.bpjs_kis)        return "BPJS / KIS wajib diupload";
  return null;
}

function nextStep() {
  submitError.value = null;
  const validators = [validateStep0, validateStep1, validateStep2];
  const err = validators[step.value]?.();
  if (err) { submitError.value = err; return; }
  step.value++;
}

async function submitForm() {
  submitting.value = true;
  submitError.value = null;
  uploadStatus.value = "Mengirim data pengajuan...";
  uploadPct.value = 10;

  try {
    const res = await api.post("/api/pengajuan/publik", {
      step1: { ...diri },
      step2: {
        ...akademis,
        kategori_magang: backendKategori(akademis.kategori_magang),
        kelas_semester: isLainnya.value ? (kategoriLainnya.value || "Lainnya") : akademis.kelas_semester,
      },
    });

    const pengajuanId = res.data?.data?.id;
    if (!pengajuanId) throw new Error("ID pengajuan tidak diterima");

    const filesToUpload: Array<{ jenis: string; file: File }> = [];
    for (const [jenis, file] of Object.entries(dokumenFiles)) {
      if (file) filesToUpload.push({ jenis, file });
    }

    for (let i = 0; i < filesToUpload.length; i++) {
      const { jenis, file } = filesToUpload[i];
      uploadStatus.value = `Mengupload dokumen ${i + 1}/${filesToUpload.length}...`;
      uploadPct.value = 20 + Math.round(((i + 1) / filesToUpload.length) * 75);

      const fd = new FormData();
      fd.append("file", file);
      fd.append("jenis", jenis);
      fd.append("pengajuan_id", pengajuanId);

      await api.post("/api/dokumen/upload-publik", fd, {
        headers: { "Content-Type": "multipart/form-data" },
      });
    }

    docsUploaded.value = filesToUpload.length;
    submitted.value = true;
  } catch (e: any) {
    submitError.value = e.response?.data?.message ?? "Terjadi kesalahan. Silakan coba lagi.";
  } finally {
    submitting.value = false;
    uploadStatus.value = "";
  }
}
</script>

<style scoped>
/* ── page ── */
.fp-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #0b1c30 0%, #0f2d14 50%, #0b1c30 100%);
  font-family: "Inter", "Segoe UI", sans-serif;
  font-size: 14px;
  color: #111827;
}

/* ── navbar ── */
.fp-nav {
  background: #fff;
  border-bottom: 1px solid #e5e7eb;
  position: sticky; top: 0; z-index: 40;
}
.fp-nav__inner {
  max-width: 900px; margin: 0 auto;
  padding: 0 24px; height: 56px;
  display: flex; align-items: center; justify-content: space-between;
}
.fp-nav__brand {
  display: flex; align-items: center; gap: 8px;
  text-decoration: none; font-size: 14px; font-weight: 700; color: #111827;
}
.fp-nav__brand strong { color: #16a34a; }
.fp-nav__logo { height: 26px; width: auto; }
.fp-nav__right { display: flex; align-items: center; gap: 16px; }
.fp-nav__link {
  display: flex; align-items: center; gap: 5px;
  font-size: 13px; font-weight: 500; color: #6b7280;
  text-decoration: none; transition: color 0.15s;
}
.fp-nav__link:hover { color: #1a3f28; }
.fp-nav__btn {
  font-size: 13px; font-weight: 600; color: #1a3f28;
  text-decoration: none; border: 1.5px solid #1a3f28;
  border-radius: 7px; padding: 6px 14px;
  transition: all 0.15s;
}
.fp-nav__btn:hover { background: #1a3f28; color: #fff; }

/* ── success ── */
.fp-success {
  display: flex; align-items: center; justify-content: center;
  min-height: 100vh; padding: 40px 20px;
}
.fp-success__card {
  background: #fff; border-radius: 16px; padding: 48px 40px;
  max-width: 480px; width: 100%; text-align: center;
  box-shadow: 0 4px 24px rgba(0,0,0,0.08);
}
.fp-success__check {
  width: 64px; height: 64px; background: #16a34a; border-radius: 50%;
  display: flex; align-items: center; justify-content: center; margin: 0 auto 20px;
}
.fp-success__card h2 { font-size: 20px; font-weight: 700; margin: 0 0 12px; color: #111827; }
.fp-success__card p  { font-size: 14px; color: #6b7280; line-height: 1.7; margin: 0 0 24px; }
.fp-success__meta { display: flex; flex-direction: column; gap: 8px; margin-bottom: 28px; text-align: left; }
.fp-success__meta-item {
  display: flex; align-items: center; gap: 8px;
  font-size: 13px; color: #374151; background: #f0fdf4;
  border-radius: 8px; padding: 10px 12px;
}
.fp-success__btn {
  display: inline-block; background: #1a3f28; color: #fff;
  text-decoration: none; border-radius: 8px; padding: 11px 28px;
  font-size: 14px; font-weight: 600; transition: opacity 0.15s;
}
.fp-success__btn:hover { opacity: 0.88; }

/* ── main form wrap ── */
.fp-wrap {
  min-height: 100vh; display: flex; flex-direction: column;
  align-items: center; padding: 40px 20px;
}

/* ── header ── */
.fp-header {
  width: 100%; max-width: 900px;
  background: linear-gradient(135deg, #1a3f28 0%, #2d6a44 100%);
  border-radius: 14px 14px 0 0;
  padding: 24px 32px;
}
.fp-header__inner {}
.fp-header__title { font-size: 18px; font-weight: 700; color: #fff; }
.fp-header__sub   { font-size: 13px; color: #a7f3d0; margin-top: 3px; }

/* ── body ── */
.fp-body {
  width: 100%; max-width: 900px;
  background: #fff;
  border-radius: 0 0 14px 14px;
  padding: 28px 32px 0;
  box-shadow: 0 4px 24px rgba(0,0,0,0.07);
}

/* ── stepper ── */
.fp-stepper {
  display: flex; align-items: center;
  padding-bottom: 28px;
  border-bottom: 1px solid #e5e7eb;
  margin-bottom: 28px;
  overflow-x: auto;
}
.fp-stepper__col {
  display: flex; flex-direction: column; align-items: center; gap: 6px;
  flex-shrink: 0;
}
.fp-stepper__circle {
  width: 32px; height: 32px; border-radius: 50%;
  border: 2px solid #d1d5db;
  display: flex; align-items: center; justify-content: center;
  font-size: 13px; font-weight: 700; color: #9ca3af; background: #fff;
  transition: all 0.2s; flex-shrink: 0;
}
.fp-stepper__circle--active { border-color: #1a3f28; background: #1a3f28; color: #fff; }
.fp-stepper__circle--done   { border-color: #1a3f28; background: #1a3f28; color: #fff; }
.fp-stepper__label {
  font-size: 11px; font-weight: 600; color: #9ca3af;
  white-space: nowrap;
}
.fp-stepper__label--active { color: #1a3f28; }
.fp-stepper__label--done   { color: #374151; }
.fp-stepper__line {
  flex: 1; height: 2px; background: #e5e7eb; margin: 0 4px;
  margin-bottom: 20px;
  min-width: 40px;
  transition: background 0.2s;
}
.fp-stepper__line--done { background: #1a3f28; }

/* ── section ── */
.fp-section { }
.fp-section__head { margin-bottom: 20px; }
.fp-section__title { font-size: 15px; font-weight: 700; color: #111827; }
.fp-section__sub   { font-size: 13px; color: #6b7280; margin-top: 2px; }

/* ── form grid ── */
.fp-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 16px; }
.fp-field { display: flex; flex-direction: column; gap: 5px; }
.fp-field--full { grid-column: 1 / -1; }
.fp-field label { font-size: 12.5px; font-weight: 600; color: #374151; }
.req { color: #ef4444; }
.fp-field input, .fp-field select, .fp-field textarea {
  border: 1px solid #d1d5db; border-radius: 8px;
  padding: 9px 12px; font-size: 13.5px; font-family: inherit;
  color: #111827; outline: none; background: #fff;
  transition: border-color 0.15s;
}
.fp-field input:focus, .fp-field select:focus, .fp-field textarea:focus {
  border-color: #1a3f28; box-shadow: 0 0 0 3px rgba(26,63,40,0.08);
}
.fp-field textarea { resize: vertical; }

/* ── notice ── */
.fp-notice {
  display: flex; align-items: center; gap: 8px;
  background: #eff6ff; border: 1px solid #bfdbfe;
  border-radius: 8px; padding: 12px 14px;
  font-size: 13px; color: #1d4ed8; margin-top: 16px;
}

/* ── doc list (upload) ── */
.fp-doc-list { display: flex; flex-direction: column; gap: 16px; }
.fp-doc { display: flex; flex-direction: column; gap: 6px; }
.fp-doc__label {
  display: flex; align-items: center; gap: 8px;
  font-size: 13px; font-weight: 600; color: #374151;
}
.fp-badge {
  font-size: 11px; font-weight: 600; padding: 2px 8px;
  border-radius: 100px; letter-spacing: 0.03em;
}
.fp-badge--required { background: #fef2f2; color: #dc2626; }
.fp-badge--optional { background: #f3f4f6; color: #6b7280; }

/* ── dropzone ── */
.fp-dropzone {
  border: 1.5px dashed #d1d5db; border-radius: 8px;
  padding: 16px 20px; cursor: pointer;
  display: flex; align-items: center; gap: 10px;
  transition: border-color 0.15s, background 0.15s;
  background: #fff;
}
.fp-dropzone:hover { border-color: #4ade80; background: #f9fafb; }
.fp-dropzone--done { border-color: #4ade80; background: #f0fdf4; }
.fp-dropzone__icon { color: #9ca3af; flex-shrink: 0; }
.fp-dropzone__text { font-size: 13.5px; font-weight: 500; color: #374151; }
.fp-dropzone__hint { font-size: 12px; color: #9ca3af; margin-left: 2px; }
.fp-dropzone__filename {
  font-size: 13px; font-weight: 500; color: #16a34a;
  flex: 1; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;
}
.fp-dropzone__remove {
  background: none; border: 1px solid #e5e7eb; border-radius: 6px;
  font-size: 12px; color: #6b7280; cursor: pointer; padding: 3px 10px;
  transition: all 0.15s; font-family: inherit;
}
.fp-dropzone__remove:hover { background: #fef2f2; border-color: #fca5a5; color: #dc2626; }

/* ── review ── */
.fp-review-block { border: 1px solid #e5e7eb; border-radius: 10px; overflow: hidden; margin-bottom: 12px; }
.fp-review-block__head {
  background: #f9fafb; padding: 9px 16px;
  font-size: 11px; font-weight: 700; color: #374151;
  text-transform: uppercase; letter-spacing: 0.08em;
  border-bottom: 1px solid #e5e7eb;
}
.fp-review-grid {
  display: grid; grid-template-columns: 1fr 1fr;
  gap: 0; padding: 12px 16px;
}
.fp-review-cell { padding: 6px 0; }
.fp-review-cell--full { grid-column: 1 / -1; }
.fp-review-cell span    { display: block; font-size: 11px; color: #9ca3af; font-weight: 500; }
.fp-review-cell strong  { display: block; font-size: 13.5px; color: #111827; font-weight: 600; margin-top: 1px; }
.fp-review-docs {
  display: grid; grid-template-columns: 1fr 1fr;
  gap: 8px; padding: 14px 16px;
}
.fp-review-doc {
  display: flex; align-items: center; gap: 8px;
  font-size: 13px; font-weight: 500;
}
.fp-review-doc--ok   { color: #16a34a; }
.fp-review-doc--miss { color: #dc2626; }
.fp-review-doc--skip { color: #9ca3af; }

/* ── progress ── */
.fp-progress { margin: 16px 0 4px; }
.fp-progress__bar { height: 5px; background: #e5e7eb; border-radius: 100px; overflow: hidden; margin-bottom: 8px; }
.fp-progress__fill { height: 100%; background: #16a34a; border-radius: 100px; transition: width 0.3s; }
.fp-progress__text { font-size: 12.5px; color: #16a34a; font-weight: 600; text-align: center; }

/* ── error ── */
.fp-error {
  margin-top: 14px; background: #fef2f2; border: 1px solid #fca5a5;
  border-radius: 8px; padding: 10px 14px; font-size: 13px; color: #dc2626;
}

/* ── footer ── */
.fp-footer {
  display: flex; align-items: center; justify-content: space-between;
  padding: 20px 0 24px;
  margin-top: 24px;
  border-top: 1px solid #f3f4f6;
}
.fp-footer__spacer { width: 80px; }
.fp-footer__counter { font-size: 13px; color: #6b7280; font-weight: 500; }
.fp-btn-back {
  background: none; border: 1px solid #d1d5db; border-radius: 8px;
  padding: 9px 22px; font-size: 13.5px; font-weight: 600;
  color: #374151; cursor: pointer; font-family: inherit;
  transition: border-color 0.15s;
}
.fp-btn-back:hover:not(:disabled) { border-color: #9ca3af; }
.fp-btn-back:disabled { opacity: 0.5; cursor: not-allowed; }
.fp-btn-next {
  background: #16a34a; color: #fff; border: none;
  border-radius: 8px; padding: 10px 28px;
  font-size: 13.5px; font-weight: 700; cursor: pointer;
  font-family: inherit; display: inline-flex; align-items: center; gap: 6px;
  transition: background 0.15s;
}
.fp-btn-next:hover:not(:disabled) { background: #15803d; }
.fp-btn-next:disabled { opacity: 0.6; cursor: not-allowed; }

/* ── spinner ── */
.fp-spinner {
  width: 14px; height: 14px;
  border: 2px solid rgba(255,255,255,0.35);
  border-top-color: #fff;
  border-radius: 50%;
  animation: fp-spin 0.6s linear infinite;
}
@keyframes fp-spin { to { transform: rotate(360deg); } }

/* ── responsive ── */
@media (max-width: 640px) {
  .fp-header, .fp-body { border-radius: 0; }
  .fp-header { padding: 20px; }
  .fp-body   { padding: 20px 16px 0; }
  .fp-grid   { grid-template-columns: 1fr; }
  .fp-review-grid, .fp-review-docs { grid-template-columns: 1fr; }
  .fp-wrap { padding: 0; }
}
</style>
