<template>
  <div>
    <!-- ── ABSENSI HARI INI ─────────────────────────────────── -->
    <div class="card">
      <div class="card-header">
        <h3 class="card-title">Absensi Hari Ini</h3>
        <span class="today-label">{{ todayLabel }}</span>
      </div>

      <div v-if="absensiLoading" class="empty-state"><div class="spinner"></div></div>

      <div v-else-if="absensiState === 'no_magang'" class="empty-state">
        <div class="empty-state__icon"><svg width="36" height="36" viewBox="0 0 24 24" fill="none"><rect x="3" y="4" width="18" height="18" rx="2" stroke="#d1d5db" stroke-width="1.5"/><line x1="3" y1="10" x2="21" y2="10" stroke="#d1d5db" stroke-width="1.5"/></svg></div>
        <p>Absensi tersedia setelah jadwal magang dikonfirmasi oleh HRD.</p>
      </div>

      <div v-else-if="absensiState === 'magang_selesai'" class="absensi-panel">
        <div class="ap-icon" style="background:#f0fdf4;color:#16a34a">
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none"><path d="M22 11.08V12a10 10 0 11-5.93-9.14" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/><polyline points="22 4 12 14.01 9 11.01" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>
        </div>
        <div class="ap-title">Periode Magang Telah Selesai</div>
        <div class="ap-desc">
          <span v-if="pelaksanaanSaya?.status === 'upload_laporan'">Silakan upload laporan akhir magang kamu.</span>
          <span v-else-if="pelaksanaanSaya?.status === 'penilaian'">Laporan sedang dalam proses penilaian.</span>
          <span v-else>Selamat! Proses magang kamu sudah selesai.</span>
        </div>
      </div>

      <div v-else-if="absensiState === 'locked'" class="absensi-panel">
        <div class="ap-icon ap-icon--locked"><svg width="24" height="24" viewBox="0 0 24 24" fill="none"><rect x="3" y="11" width="18" height="11" rx="2" stroke="currentColor" stroke-width="2"/><path d="M7 11V7a5 5 0 0110 0v4" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg></div>
        <div class="ap-title">Sesi Absen Masuk Belum Dibuka</div>
        <div class="ap-desc">Dibuka pukul <strong>{{ cfg?.jam_masuk_buka }}</strong> WIB</div>
        <div class="ap-countdown">{{ countdownMasukText }}</div>
        <button v-if="pelaksanaanSaya?.status === 'aktif'" class="btn-ajukan-sm" @click="openIzinModal">
          <svg width="11" height="11" viewBox="0 0 24 24" fill="none"><path d="M12 5v14M5 12h14" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"/></svg>
          Ajukan Izin/Sakit
        </button>
      </div>

      <div v-else-if="absensiState === 'checkin_open'" class="absensi-panel">
        <!-- Geo block state: divisi tidak aktif atau geofence belum dikonfigurasi -->
        <template v-if="geoBlockCode">
          <div class="ap-icon ap-icon--blocked">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none">
              <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2"/>
              <line x1="4.93" y1="4.93" x2="19.07" y2="19.07" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
            </svg>
          </div>
          <div class="ap-title" v-if="geoBlockCode === 'divisi_inactive'">Divisi Tidak Aktif</div>
          <div class="ap-title" v-else-if="geoBlockCode === 'divisi_not_set'">Divisi Belum Diatur</div>
          <div class="ap-title" v-else-if="geoBlockCode === 'divisi_not_found'">Divisi Tidak Ditemukan</div>
          <div class="ap-title" v-else-if="geoBlockCode === 'outside_geofence'">Di Luar Lokasi Kantor</div>
          <div class="ap-title" v-else-if="geoBlockCode === 'pelaksanaan_not_found'">Data Magang Tidak Ditemukan</div>
          <div class="ap-title" v-else>Lokasi Absensi Belum Dikonfigurasi</div>
          <div class="ap-desc ap-desc--block">{{ geoBlockMsg }}</div>
          <button class="btn-absen btn-absen--refresh" @click="fetchAbsensi">
            <svg width="13" height="13" viewBox="0 0 24 24" fill="none"><path d="M23 4v6h-6" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/><path d="M1 20v-6h6" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/><path d="M3.51 9a9 9 0 0114.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0020.49 15" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>
            Refresh Status
          </button>
        </template>
        <!-- Normal checkin state -->
        <template v-else>
          <div class="ap-icon ap-icon--open"><svg width="24" height="24" viewBox="0 0 24 24" fill="none"><polyline points="9 11 12 14 22 4" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/><path d="M21 12v7a2 2 0 01-2 2H5a2 2 0 01-2-2V5a2 2 0 012-2h11" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg></div>
          <div class="ap-title">Sesi Absen Masuk Dibuka</div>
          <div class="ap-desc">{{ cfg?.jam_masuk_buka }} – {{ cfg?.jam_masuk_tutup }} WIB</div>
          <div v-if="absensiError" class="ap-error">{{ absensiError }}</div>
          <div v-if="gpsError" class="ap-error ap-error--with-retry">
            <span>{{ gpsError }}</span>
            <button class="ap-retry-btn" @click="openCheckinModal">Coba Lagi</button>
          </div>
          <div class="ap-btn-row">
            <button class="btn-absen" :disabled="gpsLoading" @click="openCheckinModal">
              <span v-if="gpsLoading" class="btn-spinner"></span>
              {{ gpsLoading ? ('Mencari GPS...' + (gpsLiveAccuracy > 0 ? ' ±' + gpsLiveAccuracy + 'm' : '')) : 'Absen Masuk Sekarang' }}
            </button>
            <button v-if="pelaksanaanSaya?.status === 'aktif'" class="btn-ajukan" @click="openIzinModal">
              <svg width="12" height="12" viewBox="0 0 24 24" fill="none"><path d="M12 5v14M5 12h14" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"/></svg>
              Ajukan
            </button>
          </div>
        </template>
      </div>

      <div v-else-if="absensiState === 'waiting_pulang'" class="absensi-panel">
        <div class="ap-done-row"><span class="ap-badge ap-badge--masuk">✓ Masuk</span><span class="ap-done-time">{{ formatJam(todayAbsensi?.jam_masuk) }} WIB</span></div>
        <div class="ap-icon ap-icon--wait" style="margin-top:16px"><svg width="24" height="24" viewBox="0 0 24 24" fill="none"><circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2"/><polyline points="12 6 12 12 16 14" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg></div>
        <div class="ap-title">Menunggu Sesi Absen Pulang</div>
        <div class="ap-desc">Dibuka pukul <strong>{{ cfg?.jam_pulang_buka }}</strong> WIB</div>
        <div class="ap-countdown">{{ countdownPulangText }}</div>
        <button v-if="pelaksanaanSaya?.status === 'aktif'" class="btn-ajukan-sm" @click="openIzinModal">
          <svg width="11" height="11" viewBox="0 0 24 24" fill="none"><path d="M12 5v14M5 12h14" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"/></svg>
          Ajukan Izin/Sakit
        </button>
      </div>

      <div v-else-if="absensiState === 'checkout_open'" class="absensi-panel">
        <div class="ap-done-row"><span class="ap-badge ap-badge--masuk">✓ Masuk</span><span class="ap-done-time">{{ formatJam(todayAbsensi?.jam_masuk) }} WIB</span></div>
        <div class="ap-form">
          <div class="ap-form-label">Catatan Kegiatan Hari Ini <span class="ap-required">*</span></div>
          <textarea v-model="kegiatanInput" class="ap-textarea" placeholder="Tulis satu kegiatan per baris, contoh:&#10;Membuat laporan mingguan&#10;Rapat koordinasi tim&#10;Input data ke sistem" rows="4"></textarea>
          <div class="ap-form-hint">Tulis satu kegiatan per baris</div>
        </div>
        <div v-if="absensiError" class="ap-error">{{ absensiError }}</div>
        <div v-if="gpsCheckoutError" class="ap-error ap-error--with-retry">
          <span>{{ gpsCheckoutError }}</span>
          <button class="ap-retry-btn" @click="openCheckoutModal">Coba Lagi</button>
        </div>
        <div class="ap-btn-row">
          <button class="btn-absen btn-absen--pulang" :disabled="checkingOut || gpsCheckoutLoading || !kegiatanInput.trim()" @click="openCheckoutModal">
            <span v-if="gpsCheckoutLoading" class="btn-spinner"></span>
            {{ gpsCheckoutLoading ? ('Mencari GPS...' + (gpsCheckoutLiveAccuracy > 0 ? ' ±' + gpsCheckoutLiveAccuracy + 'm' : '')) : (checkingOut ? 'Memproses...' : 'Absen Pulang') }}
          </button>
          <button v-if="pelaksanaanSaya?.status === 'aktif'" class="btn-ajukan" @click="openIzinModal">
            <svg width="12" height="12" viewBox="0 0 24 24" fill="none"><path d="M12 5v14M5 12h14" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"/></svg>
            Ajukan
          </button>
        </div>
      </div>

      <div v-else-if="absensiState === 'done'" class="absensi-panel absensi-panel--done">
        <div class="ap-done-rows">
          <div class="ap-done-row"><span class="ap-badge ap-badge--masuk">✓ Masuk</span><span class="ap-done-time">{{ formatJam(todayAbsensi?.jam_masuk) }} WIB</span></div>
          <div class="ap-done-row"><span class="ap-badge ap-badge--pulang">✓ Pulang</span><span class="ap-done-time">{{ formatJam(todayAbsensi?.jam_keluar) }} WIB</span></div>
        </div>
        <div class="ap-kegiatan">
          <div class="ap-kegiatan__label">Kegiatan Hari Ini</div>
          <ol v-if="kegiatanPoin(todayAbsensi?.kegiatan).length" class="kegiatan-ol">
            <li v-for="(poin, i) in kegiatanPoin(todayAbsensi?.kegiatan)" :key="i">{{ poin }}</li>
          </ol>
          <div v-else class="ap-kegiatan__text">–</div>
        </div>
        <button v-if="pelaksanaanSaya?.status === 'aktif'" class="btn-ajukan-sm" style="margin-top:8px" @click="openIzinModal">
          <svg width="11" height="11" viewBox="0 0 24 24" fill="none"><path d="M12 5v14M5 12h14" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"/></svg>
          Ajukan Izin/Sakit
        </button>
      </div>

      <div v-else-if="absensiState === 'missed_checkin' || absensiState === 'missed_checkout'" class="absensi-panel">
        <div class="ap-icon ap-icon--miss"><svg width="24" height="24" viewBox="0 0 24 24" fill="none"><circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2"/><line x1="12" y1="8" x2="12" y2="12" stroke="currentColor" stroke-width="2" stroke-linecap="round"/><line x1="12" y1="16" x2="12.01" y2="16" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg></div>
        <div class="ap-title">{{ absensiState === 'missed_checkin' ? 'Sesi Absen Masuk Terlewat' : 'Sesi Absen Pulang Terlewat' }}</div>
        <div class="ap-desc">Jika ada keterangan izin atau sakit, ajukan di bawah ini.</div>
        <button v-if="pelaksanaanSaya?.status === 'aktif'" class="btn-ajukan" style="margin-top:4px" @click="openIzinModal">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none"><path d="M12 5v14M5 12h14" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"/></svg>
          Ajukan Izin/Sakit
        </button>
      </div>
    </div>

    <Transition name="toast">
      <div v-if="absensiSuccess" class="ap-toast">{{ absensiSuccess }}</div>
    </Transition>

    <!-- ── REKAP ABSENSI ─────────────────────────────────────── -->
    <div class="card" style="margin-top:14px">
      <div class="card-header">
        <h3 class="card-title">Rekap Absensi</h3>
        <div class="header-actions">
          <button v-if="pelaksanaanSaya?.status === 'aktif'" class="btn-outline-sm" @click="showIzinPanel = true">
            <svg width="12" height="12" viewBox="0 0 24 24" fill="none"><path d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>
            Riwayat Izin
          </button>
          <button class="btn-green-sm" @click="openPDFModal" :disabled="pdfLoading">
            <span v-if="pdfLoading" class="btn-spinner-sm" style="border-color:rgba(22,163,74,.3);border-top-color:#16a34a"></span>
            <svg v-else width="12" height="12" viewBox="0 0 24 24" fill="none"><path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/><polyline points="7 10 12 15 17 10" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/><line x1="12" y1="15" x2="12" y2="3" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>
            {{ pdfLoading ? 'Memuat...' : 'Unduh PDF' }}
          </button>
        </div>
      </div>

      <!-- Stat chips -->
      <div class="rekap-stats">
        <div class="stat-chip stat-chip--hadir"><span class="stat-num">{{ rekap.hadir }}</span><span class="stat-lbl">Hadir</span></div>
        <div class="stat-chip stat-chip--izin"><span class="stat-num">{{ rekap.izin }}</span><span class="stat-lbl">Izin</span></div>
        <div class="stat-chip stat-chip--sakit"><span class="stat-num">{{ rekap.sakit }}</span><span class="stat-lbl">Sakit</span></div>
        <div class="stat-chip stat-chip--alpha"><span class="stat-num">{{ rekap.alpha }}</span><span class="stat-lbl">Alpha</span></div>
      </div>

      <!-- Full date range table -->
      <div v-if="tabelHarian.length > 0" class="abs-table-wrap">
        <table class="abs-table">
          <thead>
            <tr>
              <th style="width:36px">No</th>
              <th>Tanggal</th>
              <th style="width:48px">Hari</th>
              <th style="width:68px">Masuk</th>
              <th style="width:68px">Pulang</th>
              <th style="width:72px">Status</th>
              <th>Kegiatan</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="row in tabelHarian" :key="row.tanggal" :class="{ 'row-today': row.isToday }">
              <td class="td-no">{{ row.no }}</td>
              <td class="td-tanggal">{{ formatTanggalShort(row.tanggal) }}</td>
              <td class="td-hari">{{ row.hari }}</td>
              <td class="td-jam">{{ row.jamMasuk }}</td>
              <td class="td-jam">{{ row.jamKeluar }}</td>
              <td>
                <span v-if="row.status !== 'belum'" :class="['ket-badge', `ket-badge--${row.status}`]">
                  {{ ({ hadir:'Hadir', izin:'Izin', sakit:'Sakit', alpha:'Alpha' } as Record<string,string>)[row.status] ?? row.status }}
                </span>
                <span v-else class="ket-badge ket-badge--belum">–</span>
              </td>
              <td class="td-kegiatan">
                <ul v-if="kegiatanPoin(row.kegiatan).length" class="kegiatan-ul">
                  <li v-for="(poin, i) in kegiatanPoin(row.kegiatan)" :key="i">{{ poin }}</li>
                </ul>
                <span v-else class="td-empty">–</span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div v-else-if="absensiLoading" class="empty-state" style="padding:28px"><div class="spinner"></div></div>

      <div v-else class="empty-state" style="padding:28px">
        <p>Jadwal magang belum tersedia. Rekap akan muncul setelah HRD mengkonfirmasi jadwal.</p>
      </div>
    </div>
  </div>

  <!-- ── MODAL CHECKIN ──────────────────────────────────────── -->
  <Teleport to="body">
    <div v-if="showCheckinModal" class="modal-backdrop" @click.self="showCheckinModal = false">
      <div class="modal-box">
        <div class="modal-box__title">Konfirmasi Absen Masuk</div>
        <div class="modal-box__desc">Kamu akan dicatat hadir pada pukul <strong>{{ currentTimeWIB }} WIB</strong>.</div>
        <div class="modal-box__geo">
          <svg width="13" height="13" viewBox="0 0 24 24" fill="none"><path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5S10.62 6.5 12 6.5s2.5 1.12 2.5 2.5S13.38 11.5 12 11.5z" fill="#16a34a"/></svg>
          Koordinat: {{ userLat.toFixed(5) }}, {{ userLng.toFixed(5) }}
          <span class="modal-box__acc" :class="userAccuracy <= 50 ? 'modal-box__acc--good' : 'modal-box__acc--ok'">
            ±{{ Math.round(userAccuracy) }} m
          </span>
        </div>
        <div class="modal-box__actions">
          <button class="btn-cancel" @click="showCheckinModal = false">Batal</button>
          <button class="btn-confirm" @click="doCheckin" :disabled="checkingIn">{{ checkingIn ? 'Memproses...' : 'Ya, Absen Sekarang' }}</button>
        </div>
      </div>
    </div>
  </Teleport>

  <!-- ── MODAL CHECKOUT (GPS) ────────────────────────────────── -->
  <Teleport to="body">
    <div v-if="showCheckoutModal" class="modal-backdrop" @click.self="showCheckoutModal = false">
      <div class="modal-box">
        <div class="modal-box__title">Konfirmasi Absen Pulang</div>
        <div class="modal-box__desc">Kamu akan dicatat pulang pada pukul <strong>{{ currentTimeWIB }} WIB</strong>.</div>
        <div class="modal-box__geo">
          <svg width="13" height="13" viewBox="0 0 24 24" fill="none"><path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5S10.62 6.5 12 6.5s2.5 1.12 2.5 2.5S13.38 11.5 12 11.5z" fill="#ea580c"/></svg>
          Koordinat: {{ checkoutLat.toFixed(5) }}, {{ checkoutLng.toFixed(5) }}
          <span class="modal-box__acc" :class="checkoutAccuracy <= 50 ? 'modal-box__acc--good' : 'modal-box__acc--ok'">
            ±{{ Math.round(checkoutAccuracy) }} m
          </span>
        </div>
        <div class="modal-box__actions">
          <button class="btn-cancel" @click="showCheckoutModal = false">Batal</button>
          <button class="btn-confirm btn-confirm--pulang" @click="doCheckout" :disabled="checkingOut">{{ checkingOut ? 'Memproses...' : 'Ya, Absen Pulang' }}</button>
        </div>
      </div>
    </div>
  </Teleport>

  <!-- ── MODAL AJUKAN IZIN / SAKIT ──────────────────────────── -->
  <Teleport to="body">
    <Transition name="modal-fade">
      <div v-if="showIzinModal" class="modal-backdrop" @click.self="closeIzinModal">
        <div class="modal-box modal-box--izin">
          <div class="modal-box__header">
            <div class="modal-box__title">Ajukan Izin / Sakit</div>
            <button class="modal-close-btn" @click="closeIzinModal">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none"><path d="M18 6L6 18M6 6l12 12" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"/></svg>
            </button>
          </div>

          <div class="modal-field">
            <label class="modal-label">Jenis Ketidakhadiran <span class="ap-required">*</span></label>
            <div class="jenis-tabs">
              <button :class="['jenis-tab', izinForm.jenis === 'izin' && 'jenis-tab--active']"
                @click="izinForm.jenis = 'izin'">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none"><rect x="3" y="4" width="18" height="18" rx="2" stroke="currentColor" stroke-width="2"/><path d="M8 2v4M16 2v4M3 10h18" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>
                Izin
              </button>
              <button :class="['jenis-tab', izinForm.jenis === 'sakit' && 'jenis-tab--active jenis-tab--sakit']"
                @click="izinForm.jenis = 'sakit'">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none"><path d="M12 2a10 10 0 100 20A10 10 0 0012 2z" stroke="currentColor" stroke-width="2"/><path d="M12 8v4M12 16h.01" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>
                Sakit
              </button>
            </div>
          </div>

          <div class="modal-field">
            <label class="modal-label">Tanggal <span class="ap-required">*</span></label>
            <input v-model="izinForm.tanggal" type="date" class="modal-input" :min="minTanggalIzin" />
            <div class="modal-hint">Boleh memilih tanggal hari ini atau ke depan</div>
          </div>

          <div class="modal-field">
            <label class="modal-label">Alasan <span class="ap-required">*</span></label>
            <textarea v-model="izinForm.alasan" class="modal-textarea" rows="3"
              :placeholder="izinForm.jenis === 'sakit' ? 'Contoh: Demam tinggi dan dianjurkan istirahat oleh dokter' : 'Contoh: Ada keperluan keluarga mendadak'">
            </textarea>
          </div>

          <Transition name="slide-down">
            <div class="modal-field">
              <label class="modal-label">
                <template v-if="izinForm.jenis === 'izin'">Surat Izin <span class="ap-required">*</span></template>
                <template v-else>Surat Sakit <span class="ap-required">*</span></template>
              </label>
              <div class="upload-zone" :class="{ 'upload-zone--has-file': izinForm.buktiFile }"
                @click="triggerFileInput" @dragover.prevent @drop.prevent="onFileDrop">
                <input ref="fileInputRef" type="file" accept=".jpg,.jpeg,.png,.pdf" style="display:none" @change="onFileChange" />
                <template v-if="!izinForm.buktiFile">
                  <div class="upload-zone__icon"><svg width="22" height="22" viewBox="0 0 24 24" fill="none"><path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/><polyline points="17 8 12 3 7 8" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/><line x1="12" y1="3" x2="12" y2="15" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg></div>
                  <div class="upload-zone__text">Klik atau seret file ke sini</div>
                  <div class="upload-zone__hint">JPG, PNG, atau PDF · Maks. 10MB</div>
                </template>
                <template v-else>
                  <div class="upload-zone__preview">
                    <svg width="18" height="18" viewBox="0 0 24 24" fill="none"><path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z" stroke="#16a34a" stroke-width="2"/><polyline points="14 2 14 8 20 8" stroke="#16a34a" stroke-width="2"/></svg>
                    <span class="upload-zone__filename">{{ izinForm.buktiFile.name }}</span>
                    <button class="upload-zone__remove" @click.stop="removeFile">
                      <svg width="12" height="12" viewBox="0 0 24 24" fill="none"><path d="M18 6L6 18M6 6l12 12" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"/></svg>
                    </button>
                  </div>
                </template>
              </div>
            </div>
          </Transition>

          <div v-if="izinError" class="ap-error" style="margin-top:4px">{{ izinError }}</div>

          <div class="modal-box__actions" style="margin-top:18px">
            <button class="btn-cancel" @click="closeIzinModal">Batal</button>
            <button class="btn-confirm" @click="submitIzin"
              :disabled="submittingIzin || !izinForm.tanggal || !izinForm.jenis || izinForm.alasan.length < 5 || !izinForm.buktiFile">
              <span v-if="submittingIzin" class="btn-spinner-sm"></span>
              {{ submittingIzin ? 'Mengirim...' : 'Kirim Pengajuan' }}
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>

  <!-- ── MODAL SURAT SAKIT VIEWER ────────────────────────────── -->
  <Teleport to="body">
    <div v-if="buktiModal.show" class="bukti-overlay" @click.self="buktiModal.show = false">
      <div class="bukti-box">
        <div class="bukti-header">
          <span class="bukti-title">Surat Sakit</span>
          <div class="bukti-header-actions">
            <a :href="buktiModal.url" download class="btn-outline-sm" title="Download">
              <svg width="12" height="12" viewBox="0 0 24 24" fill="none"><path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/><polyline points="7 10 12 15 17 10" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/><line x1="12" y1="15" x2="12" y2="3" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>
              Download
            </a>
            <button class="modal-close-btn" @click="buktiModal.show = false">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none"><path d="M18 6L6 18M6 6l12 12" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"/></svg>
            </button>
          </div>
        </div>
        <div class="bukti-body">
          <img v-if="buktiModal.type === 'image'" :src="buktiModal.url" class="bukti-img" alt="Surat Sakit" />
          <iframe v-else-if="buktiModal.type === 'pdf'" :src="buktiModal.url" class="bukti-iframe" />
          <div v-else class="bukti-unsupported">
            <svg width="40" height="40" viewBox="0 0 24 24" fill="none"><path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z" stroke="#9ca3af" stroke-width="1.5"/><polyline points="14 2 14 8 20 8" stroke="#9ca3af" stroke-width="1.5"/></svg>
            <span>Format file tidak dapat ditampilkan</span>
            <a :href="buktiModal.url" download class="btn-green-sm">Download File</a>
          </div>
        </div>
      </div>
    </div>
  </Teleport>

  <!-- ── PDF MODAL ─────────────────────────────────────────────── -->
  <Teleport to="body">
    <Transition name="modal-fade">
      <div v-if="showPDFModal" class="modal-backdrop" @click.self="closePDFModal">
        <div class="modal-box modal-box--pdf">
          <div class="modal-box__header">
            <div class="modal-box__title">Rekap Absensi PDF</div>
            <div style="display:flex;gap:8px;align-items:center">
              <a v-if="pdfBlobUrl" :href="pdfBlobUrl" download="Rekap_Absensi.pdf" class="btn-confirm" style="text-decoration:none;font-size:12px;padding:7px 16px">
                <svg width="12" height="12" viewBox="0 0 24 24" fill="none"><path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/><polyline points="7 10 12 15 17 10" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/><line x1="12" y1="15" x2="12" y2="3" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>
                Download
              </a>
              <button class="modal-close-btn" @click="closePDFModal">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none"><line x1="18" y1="6" x2="6" y2="18" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"/><line x1="6" y1="6" x2="18" y2="18" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"/></svg>
              </button>
            </div>
          </div>
          <div class="pdf-modal-body">
            <iframe v-if="pdfBlobUrl" :src="pdfBlobUrl" class="pdf-modal-iframe"></iframe>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>

  <!-- ── RIWAYAT IZIN / SAKIT — SIDE PANEL + BOTTOM SHEET ─────── -->
  <Teleport to="body">
    <Transition name="panel-fade">
      <div v-if="showIzinPanel" class="panel-wrapper" @click.self="showIzinPanel = false">
        <Transition name="panel-slide">
          <div v-if="showIzinPanel" class="panel-body">
            <div class="panel-header">
              <div>
                <div class="panel-title">Riwayat Izin &amp; Sakit</div>
                <div class="panel-sub">{{ izinList.length }} pengajuan</div>
              </div>
              <div class="panel-header-right">
                <button class="btn-green-sm" @click="openIzinModal" style="font-size:11px">+ Ajukan Baru</button>
                <button class="modal-close-btn" @click="showIzinPanel = false">
                  <svg width="16" height="16" viewBox="0 0 24 24" fill="none"><path d="M18 6L6 18M6 6l12 12" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"/></svg>
                </button>
              </div>
            </div>
            <div class="panel-content">
              <div v-if="izinLoading" class="empty-state" style="padding:28px"><div class="spinner"></div></div>
              <div v-else-if="izinList.length === 0" class="empty-state" style="padding:40px">
                <div class="empty-state__icon"><svg width="36" height="36" viewBox="0 0 24 24" fill="none"><path d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" stroke="#d1d5db" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/></svg></div>
                <p>Belum ada riwayat pengajuan izin/sakit.</p>
              </div>
              <div v-else class="izin-list">
                <div v-for="item in izinList" :key="item.id" class="izin-row">
                  <div class="izin-row__top">
                    <span :class="['izin-badge', `izin-badge--${item.jenis}`]">{{ item.jenis === 'izin' ? 'Izin' : 'Sakit' }}</span>
                    <span class="izin-tanggal">{{ formatTanggalShort(item.tanggal) }}</span>
                    <span :class="['izin-status', `izin-status--${item.status}`]">
                      {{ { pending: 'Menunggu', disetujui: 'Disetujui', ditolak: 'Ditolak' }[item.status] ?? item.status }}
                    </span>
                  </div>
                  <div class="izin-alasan">{{ item.alasan }}</div>
                  <div v-if="item.catatan_hrd" class="izin-catatan">
                    <svg width="11" height="11" viewBox="0 0 24 24" fill="none"><path d="M21 15a2 2 0 01-2 2H7l-4 4V5a2 2 0 012-2h14a2 2 0 012 2z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>
                    {{ item.catatan_hrd }}
                  </div>
                  <div v-if="item.bukti_path" class="izin-bukti-row">
                    <button class="bukti-link" @click="lihatBukti(item.bukti_path!)">
                      <svg width="11" height="11" viewBox="0 0 24 24" fill="none"><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z" stroke="currentColor" stroke-width="2"/><circle cx="12" cy="12" r="3" stroke="currentColor" stroke-width="2"/></svg>
                      Lihat Surat Sakit
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </Transition>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from "vue";
import api from "@/lib/api";

type AbsensiCfg = {
  jam_masuk_buka: string; jam_masuk_tutup: string;
  jam_pulang_buka: string; jam_pulang_tutup: string;
};

interface IzinSakitItem {
  id: string;
  tanggal: string;
  jenis: string;
  alasan: string;
  status: string;
  bukti_path?: string | null;
  catatan_hrd?: string | null;
}

const cfg            = ref<AbsensiCfg | null>(null);
const absensiList    = ref<any[]>([]);
const rekap          = ref({ hadir: 0, izin: 0, sakit: 0, alpha: 0 });
const absensiLoading = ref(false);
const showCheckinModal = ref(false);
const kegiatanInput  = ref('');
const checkingIn     = ref(false);
const checkingOut    = ref(false);
const absensiError   = ref('');
const absensiSuccess = ref('');

const gpsLoading        = ref(false);
const gpsError          = ref('');
const gpsLiveAccuracy   = ref(0);
const userLat           = ref(0);
const userLng           = ref(0);
const userDistanceM     = ref(0);

// GPS checkout
const gpsCheckoutLoading       = ref(false);
const gpsCheckoutError         = ref('');
const gpsCheckoutLiveAccuracy  = ref(0);
const checkoutLat              = ref(0);
const checkoutLng              = ref(0);
const checkoutAccuracy         = ref(0);
const showCheckoutModal        = ref(false);

// GPS checkin accuracy
const userAccuracy = ref(0);

const geoBlockCode   = ref('');
const geoBlockMsg    = ref('');

// Koordinat divisi — didapat dari /geo-status saat halaman dimuat
const divisiGeoLat    = ref<number | null>(null);
const divisiGeoLng    = ref<number | null>(null);
const divisiGeoRadius = ref<number | null>(null);

const pelaksanaanSaya = ref<any>(null);

const izinList       = ref<IzinSakitItem[]>([]);
const izinLoading    = ref(false);
const showIzinModal  = ref(false);
const showIzinPanel  = ref(false);
const submittingIzin = ref(false);
const izinError      = ref('');
const izinForm       = ref<{ tanggal: string; jenis: string; alasan: string; buktiFile: File | null }>({
  tanggal: '', jenis: 'izin', alasan: '', buktiFile: null
});
const fileInputRef   = ref<HTMLInputElement | null>(null);

const buktiModal = ref<{ show: boolean; url: string; type: 'image' | 'pdf' | 'other' }>({
  show: false, url: '', type: 'other'
});

const showPDFModal = ref(false);
const pdfBlobUrl   = ref('');
const pdfLoading   = ref(false);

const nowWIB = ref(getNowWIB());
let clockTimer: ReturnType<typeof setInterval> | null = null;

function getNowWIB(): Date {
  const utcMs = Date.now() + new Date().getTimezoneOffset() * 60000;
  return new Date(utcMs + 7 * 3600000);
}

function parseWinTime(hhMM: string, base: Date): Date {
  if (!hhMM) return base;
  const [h, m] = hhMM.split(':').map(Number);
  const d = new Date(base);
  d.setHours(h, m, 0, 0);
  return d;
}

const todayStr = computed(() => {
  const d = nowWIB.value;
  return `${d.getFullYear()}-${String(d.getMonth()+1).padStart(2,'0')}-${String(d.getDate()).padStart(2,'0')}`;
});

const todayLabel = computed(() =>
  nowWIB.value.toLocaleDateString('id-ID', { weekday:'long', day:'numeric', month:'long', year:'numeric' })
);

const currentTimeWIB = computed(() => {
  const d = nowWIB.value;
  return `${String(d.getHours()).padStart(2,'0')}:${String(d.getMinutes()).padStart(2,'0')}`;
});

const minTanggalIzin = computed(() => todayStr.value);

const todayAbsensi = computed(() =>
  absensiList.value.find(a => (a.tanggal || '').slice(0, 10) === todayStr.value) ?? null
);

const absensiState = computed(() => {
  if (absensiLoading.value) return 'loading';
  if (!pelaksanaanSaya.value) return 'no_magang';
  const st = pelaksanaanSaya.value.status;
  // Magang belum dikonfirmasi sama sekali
  if (st === 'menunggu_mulai') return 'no_magang';
  // Magang sudah selesai / dalam proses penilaian / upload laporan
  if (st === 'upload_laporan' || st === 'penilaian' || st === 'selesai') return 'magang_selesai';
  // Hanya status 'aktif' yang tampilkan kontrol absensi
  if (st !== 'aktif') return 'no_magang';
  if (!cfg.value) return 'loading';
  const now = nowWIB.value;
  const masukBuka   = parseWinTime(cfg.value.jam_masuk_buka,   now);
  const masukTutup  = parseWinTime(cfg.value.jam_masuk_tutup,  now);
  const pulangBuka  = parseWinTime(cfg.value.jam_pulang_buka,  now);
  const pulangTutup = parseWinTime(cfg.value.jam_pulang_tutup, now);
  const ta = todayAbsensi.value;
  if (ta?.jam_keluar) return 'done';
  if (ta?.jam_masuk) {
    if (now >= pulangBuka && now <= pulangTutup) return 'checkout_open';
    if (now < pulangBuka) return 'waiting_pulang';
    return 'missed_checkout';
  }
  if (now >= masukBuka && now <= masukTutup) return 'checkin_open';
  if (now < masukBuka) return 'locked';
  return 'missed_checkin';
});

const HARI_SHORT = ['Min','Sen','Sel','Rab','Kam','Jum','Sab'];

const tabelHarian = computed(() => {
  if (!pelaksanaanSaya.value?.tanggal_mulai || !pelaksanaanSaya.value?.tanggal_selesai) return [];
  const mulaiRaw   = pelaksanaanSaya.value.tanggal_mulai as string;
  const selesaiRaw = pelaksanaanSaya.value.tanggal_selesai as string;
  const mulaiD   = new Date(mulaiRaw);
  const selesaiD = new Date(selesaiRaw);
  const today = todayStr.value;
  const rows: any[] = [];
  let no = 1;
  const cur = new Date(Date.UTC(mulaiD.getUTCFullYear(), mulaiD.getUTCMonth(), mulaiD.getUTCDate()));
  const end = new Date(Date.UTC(selesaiD.getUTCFullYear(), selesaiD.getUTCMonth(), selesaiD.getUTCDate()));
  while (cur <= end) {
    const dow = cur.getUTCDay();
    if (dow !== 0 && dow !== 6) {
      const dateStr = cur.toISOString().slice(0, 10);
      const absensi = absensiList.value.find(a => (a.tanggal || '').slice(0, 10) === dateStr);
      let status = dateStr < today ? 'alpha' : 'belum';
      if (absensi) status = absensi.keterangan;
      rows.push({
        no: no++,
        tanggal: dateStr,
        hari: HARI_SHORT[dow],
        jamMasuk:  formatJam(absensi?.jam_masuk),
        jamKeluar: formatJam(absensi?.jam_keluar),
        kegiatan:  absensi?.kegiatan   || '',
        status,
        isToday: dateStr === today,
      });
    }
    cur.setUTCDate(cur.getUTCDate() + 1);
  }
  return rows;
});

function formatJam(t: string | null | undefined): string {
  if (!t) return '–';
  return t.slice(0, 5);
}

function kegiatanPoin(text: string | null | undefined): string[] {
  if (!text || text === '–') return [];
  return text.split('\n').map(l => l.trim()).filter(l => l.length > 0);
}

function formatCountdown(ms: number): string {
  if (ms <= 0) return '–';
  const s = Math.floor(ms / 1000);
  const h = Math.floor(s / 3600);
  const m = Math.floor((s % 3600) / 60);
  const sec = s % 60;
  if (h > 0) return `${h} jam ${m} menit lagi`;
  if (m > 0) return `${m} menit ${sec} detik lagi`;
  return `${sec} detik lagi`;
}

const countdownMasukText = computed(() => {
  if (!cfg.value) return '';
  return formatCountdown(parseWinTime(cfg.value.jam_masuk_buka, nowWIB.value).getTime() - nowWIB.value.getTime());
});

const countdownPulangText = computed(() => {
  if (!cfg.value) return '';
  return formatCountdown(parseWinTime(cfg.value.jam_pulang_buka, nowWIB.value).getTime() - nowWIB.value.getTime());
});

async function fetchAbsensi() {
  absensiLoading.value = true;
  absensiError.value = '';
  try {
    const [r0, r1, r2, r3] = await Promise.allSettled([
      api.get('/api/pelaksanaan/saya'),
      api.get('/api/absensi/config'),
      api.get('/api/absensi/saya'),
      api.get('/api/absensi/geo-status'),
    ]);
    if (r0.status === 'fulfilled') pelaksanaanSaya.value = r0.value.data?.data ?? null;
    if (r1.status === 'fulfilled') cfg.value = r1.value.data?.data;
    if (r2.status === 'fulfilled') {
      absensiList.value = r2.value.data?.data?.list ?? [];
      rekap.value = r2.value.data?.data?.rekap ?? { hadir:0, izin:0, sakit:0, alpha:0 };
    }
    if (r3.status === 'fulfilled') {
      const geo = r3.value.data;
      if (geo && !geo.can_checkin) {
        geoBlockCode.value  = geo.code ?? '';
        geoBlockMsg.value   = geo.message ?? '';
        divisiGeoLat.value  = null;
        divisiGeoLng.value  = null;
        divisiGeoRadius.value = null;
      } else if (geo) {
        geoBlockCode.value    = '';
        geoBlockMsg.value     = '';
        divisiGeoLat.value    = geo.geo_lat ?? null;
        divisiGeoLng.value    = geo.geo_lng ?? null;
        divisiGeoRadius.value = geo.geo_radius ?? null;
      }
    }
  } finally {
    absensiLoading.value = false;
  }
}

async function fetchIzin() {
  izinLoading.value = true;
  try {
    const r = await api.get('/api/izin-sakit/saya');
    izinList.value = Array.isArray(r.data?.data) ? r.data.data : [];
  } catch {
    izinList.value = [];
  } finally {
    izinLoading.value = false;
  }
}

function haversineM(lat1: number, lng1: number, lat2: number, lng2: number): number {
  const R = 6371000;
  const toRad = (d: number) => d * Math.PI / 180;
  const dLat = toRad(lat2 - lat1);
  const dLng = toRad(lng2 - lng1);
  const a = Math.sin(dLat / 2) ** 2 + Math.cos(toRad(lat1)) * Math.cos(toRad(lat2)) * Math.sin(dLng / 2) ** 2;
  return R * 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));
}

// ── Shared GPS helper ────────────────────────────────────────────────────────
// Batas akurasi yang diterima: ≤300 m — mencakup GPS satelit & WiFi positioning.
// Akurasi > 300 m (cell tower murni tanpa WiFi) ditolak.
//
// Strategi:
//  1. watchPosition high-accuracy (timeout:Infinity)
//  2. Panggil onProgress(accuracy) tiap update — untuk live feedback di UI.
//  3. Terima SEGERA jika accuracy ≤ 50 m (GPS bagus).
//  4. Jika accuracy 50–300 m, tunggu 5 detik refinement lalu terima.
//  5. Hard timer 60 s: pakai posisi TERBAIK yang ada — meski accuracy > 300 m.
//     Tidak reject karena akurasi — geofence jarak di backend yang menentukan.
//     Hanya reject jika GPS benar-benar tidak mengembalikan posisi sama sekali.

function acquireGPS(onProgress?: (accuracy: number) => void): Promise<GeolocationPosition> {
  return new Promise((resolve, reject) => {
    if (!navigator.geolocation) { reject(new Error('no_geolocation')); return; }

    let bestPos: GeolocationPosition | null = null;
    let watchId    = -1;
    let settled    = false;
    let refineTimer: ReturnType<typeof setTimeout> | null = null;
    let hardTimer:   ReturnType<typeof setTimeout>;

    const cleanup = () => {
      if (watchId !== -1) { navigator.geolocation.clearWatch(watchId); watchId = -1; }
      if (refineTimer) { clearTimeout(refineTimer); refineTimer = null; }
      clearTimeout(hardTimer);
    };

    const done = (pos: GeolocationPosition | null, errCode?: string) => {
      if (settled) return;
      settled = true;
      cleanup();
      if (pos) resolve(pos);
      else     reject(new Error(errCode ?? 'unavailable'));
    };

    watchId = navigator.geolocation.watchPosition(
      (pos) => {
        // Selalu simpan posisi paling akurat yang pernah didapat
        if (!bestPos || pos.coords.accuracy < bestPos.coords.accuracy) bestPos = pos;

        // Kirim update akurasi real-time ke UI
        onProgress?.(pos.coords.accuracy);

        // Akurasi sangat bagus → terima langsung
        if (pos.coords.accuracy <= 50) { done(bestPos); return; }

        // Akurasi sedang (50–300 m) → tunggu 5 detik refinement
        if (pos.coords.accuracy <= 300 && !refineTimer) {
          refineTimer = setTimeout(() => { done(bestPos); }, 5000);
        }
        // Akurasi > 300 m → terus tunggu, semoga membaik
      },
      (err) => {
        if (settled) return;
        if (err.code === 1 /* PERMISSION_DENIED */) { done(null, 'permission_denied'); return; }
        // Error lain: gunakan posisi terbaik jika ada, jika tidak → unavailable
        if (bestPos) { done(bestPos); return; }
        done(null, 'unavailable');
      },
      { enableHighAccuracy: true, timeout: Infinity, maximumAge: 0 },
    );

    // Hard timer 60 s: gunakan posisi TERBAIK yang ada, apapun akurasinya.
    // Cek akurasi TIDAK dilakukan di sini — backend yang menentukan via geofence.
    hardTimer = setTimeout(() => {
      if (settled) return;
      if (bestPos) { done(bestPos); return; }
      done(null, 'unavailable');
    }, 60000);
  });
}

function openCheckinModal() {
  gpsError.value        = '';
  absensiError.value    = '';
  gpsLiveAccuracy.value = 0;
  if (!navigator.geolocation) {
    gpsError.value = 'Browser tidak mendukung GPS. Gunakan Chrome/Firefox versi terbaru.';
    return;
  }
  gpsLoading.value = true;

  acquireGPS((acc) => { gpsLiveAccuracy.value = Math.round(acc); })
    .then((pos) => {
      userLat.value      = pos.coords.latitude;
      userLng.value      = pos.coords.longitude;
      userAccuracy.value = pos.coords.accuracy;

      // Cek jarak client-side sebelum modal muncul
      if (divisiGeoLat.value !== null && divisiGeoLng.value !== null && divisiGeoRadius.value !== null) {
        const jarakM = haversineM(userLat.value, userLng.value, divisiGeoLat.value, divisiGeoLng.value);
        if (jarakM > divisiGeoRadius.value) {
          gpsError.value = `Lokasi Anda terlalu jauh dari kantor (${Math.round(jarakM)} m, radius ${Math.round(divisiGeoRadius.value)} m). Harap absen dari lokasi kantor.`;
          return; // jangan tampilkan modal
        }
      }

      showCheckinModal.value = true;
    })
    .catch((err: Error) => {
      if (err.message === 'permission_denied') {
        gpsError.value = 'Izin lokasi ditolak. Buka pengaturan browser → izinkan akses lokasi, lalu tekan Coba Lagi.';
      } else if (err.message === 'no_geolocation') {
        gpsError.value = 'Browser tidak mendukung GPS. Gunakan Chrome/Firefox versi terbaru.';
      } else {
        gpsError.value = 'Lokasi tidak terdeteksi. Pastikan GPS & izin lokasi aktif di browser, lalu tekan Coba Lagi.';
      }
    })
    .finally(() => { gpsLoading.value = false; });
}

async function doCheckin() {
  checkingIn.value = true;
  absensiError.value = '';
  try {
    await api.post('/api/absensi/checkin', {
      latitude:  userLat.value,
      longitude: userLng.value,
      accuracy:  userAccuracy.value,
    });
    showCheckinModal.value = false;
    showToast('Absen masuk berhasil!');
    await fetchAbsensi();
  } catch (e: any) {
    showCheckinModal.value = false;
    const code = e.response?.data?.error ?? '';
    const msg  = e.response?.data?.message ?? 'Gagal absen masuk';
    const geoErrorCodes = ['divisi_inactive', 'geofence_not_configured', 'outside_geofence',
                           'divisi_not_set', 'divisi_not_found', 'pelaksanaan_not_found'];
    if (geoErrorCodes.includes(code)) {
      geoBlockCode.value = code;
      geoBlockMsg.value  = msg;
    } else {
      absensiError.value = msg;
    }
  } finally {
    checkingIn.value = false;
  }
}

function openCheckoutModal() {
  if (!kegiatanInput.value.trim()) return;
  gpsCheckoutError.value         = '';
  absensiError.value             = '';
  gpsCheckoutLiveAccuracy.value  = 0;
  if (!navigator.geolocation) {
    gpsCheckoutError.value = 'Browser tidak mendukung GPS. Gunakan Chrome/Firefox versi terbaru.';
    return;
  }
  gpsCheckoutLoading.value = true;

  acquireGPS((acc) => { gpsCheckoutLiveAccuracy.value = Math.round(acc); })
    .then((pos) => {
      checkoutLat.value      = pos.coords.latitude;
      checkoutLng.value      = pos.coords.longitude;
      checkoutAccuracy.value = pos.coords.accuracy;

      // Cek jarak client-side sebelum modal pulang muncul
      if (divisiGeoLat.value !== null && divisiGeoLng.value !== null && divisiGeoRadius.value !== null) {
        const jarakM = haversineM(checkoutLat.value, checkoutLng.value, divisiGeoLat.value, divisiGeoLng.value);
        if (jarakM > divisiGeoRadius.value) {
          gpsCheckoutError.value = `Lokasi Anda terlalu jauh dari kantor (${Math.round(jarakM)} m, radius ${Math.round(divisiGeoRadius.value)} m). Harap absen dari lokasi kantor.`;
          return; // jangan tampilkan modal
        }
      }

      showCheckoutModal.value = true;
    })
    .catch((err: Error) => {
      if (err.message === 'permission_denied') {
        gpsCheckoutError.value = 'Izin lokasi ditolak. Buka pengaturan browser → izinkan akses lokasi, lalu tekan Coba Lagi.';
      } else if (err.message === 'no_geolocation') {
        gpsCheckoutError.value = 'Browser tidak mendukung GPS. Gunakan Chrome/Firefox versi terbaru.';
      } else {
        gpsCheckoutError.value = 'Lokasi tidak terdeteksi. Pastikan GPS & izin lokasi aktif di browser, lalu tekan Coba Lagi.';
      }
    })
    .finally(() => { gpsCheckoutLoading.value = false; });
}

async function doCheckout() {
  if (!kegiatanInput.value.trim()) return;
  showCheckoutModal.value = false;
  checkingOut.value = true;
  absensiError.value = '';
  try {
    await api.patch('/api/absensi/checkout', {
      kegiatan:  kegiatanInput.value.trim(),
      latitude:  checkoutLat.value,
      longitude: checkoutLng.value,
      accuracy:  checkoutAccuracy.value,
    });
    kegiatanInput.value = '';
    showToast('Absen pulang berhasil!');
    await fetchAbsensi();
  } catch (e: any) {
    const code = (e.response?.data?.error ?? '') as string;
    const msg  = e.response?.data?.message ?? 'Gagal absen pulang';
    const geoErrorCodes = ['divisi_inactive', 'geofence_not_configured', 'outside_geofence',
                           'divisi_not_set', 'divisi_not_found', 'pelaksanaan_not_found'];
    if (geoErrorCodes.includes(code)) {
      geoBlockCode.value = code;
      geoBlockMsg.value  = msg;
    } else {
      absensiError.value = msg;
    }
  } finally {
    checkingOut.value = false;
  }
}

function openIzinModal() {
  izinForm.value = { tanggal: todayStr.value, jenis: 'izin', alasan: '', buktiFile: null };
  izinError.value = '';
  showIzinModal.value = true;
}

function closeIzinModal() {
  showIzinModal.value = false;
  izinError.value = '';
}

function triggerFileInput() { fileInputRef.value?.click(); }

function onFileChange(e: Event) {
  const input = e.target as HTMLInputElement;
  if (input.files?.[0]) izinForm.value.buktiFile = input.files[0];
}

function onFileDrop(e: DragEvent) {
  const file = e.dataTransfer?.files?.[0];
  if (file) izinForm.value.buktiFile = file;
}

function removeFile() { izinForm.value.buktiFile = null; resetFileInput(); }
function resetFileInput() { if (fileInputRef.value) fileInputRef.value.value = ''; }

async function submitIzin() {
  if (!izinForm.value.tanggal || !izinForm.value.jenis || izinForm.value.alasan.length < 5) return;
  submittingIzin.value = true;
  izinError.value = '';
  try {
    const fd = new FormData();
    fd.append('tanggal', izinForm.value.tanggal);
    fd.append('jenis', izinForm.value.jenis);
    fd.append('alasan', izinForm.value.alasan.trim());
    if (izinForm.value.buktiFile) {
      fd.append('bukti', izinForm.value.buktiFile);
    }
    await api.post('/api/izin-sakit', fd, { headers: { 'Content-Type': 'multipart/form-data' } });
    closeIzinModal();
    showToast(`Pengajuan ${izinForm.value.jenis} berhasil dikirim!`);
    await fetchIzin();
  } catch (e: any) {
    izinError.value = e.response?.data?.message ?? 'Gagal mengirim pengajuan';
  } finally {
    submittingIzin.value = false;
  }
}

function lihatBukti(path: string) {
  const ext = path.split('.').pop()?.toLowerCase() ?? '';
  const type = ['jpg','jpeg','png','webp','gif'].includes(ext) ? 'image' : ext === 'pdf' ? 'pdf' : 'other';
  buktiModal.value = { show: true, url: `/uploads/${path}`, type };
}

function showToast(msg: string) {
  absensiSuccess.value = msg;
  setTimeout(() => { absensiSuccess.value = ''; }, 3000);
}

function formatTanggalShort(iso: string) {
  if (!iso) return '–';
  const d = new Date(iso.slice(0, 10) + 'T00:00:00');
  return d.toLocaleDateString('id-ID', { day:'numeric', month:'short', year:'numeric' });
}

async function openPDFModal() {
  if (pdfLoading.value) return;
  pdfLoading.value = true;
  try {
    const r = await api.get('/api/absensi/saya/pdf', { responseType: 'blob' });
    if (pdfBlobUrl.value) URL.revokeObjectURL(pdfBlobUrl.value);
    pdfBlobUrl.value = URL.createObjectURL(new Blob([r.data], { type: 'application/pdf' }));
    showPDFModal.value = true;
  } catch {
    showToast('Gagal memuat PDF. Coba lagi.');
  } finally {
    pdfLoading.value = false;
  }
}

function closePDFModal() {
  showPDFModal.value = false;
}

onMounted(() => {
  fetchAbsensi();
  fetchIzin();
  clockTimer = setInterval(() => { nowWIB.value = getNowWIB(); }, 1000);
});

onUnmounted(() => {
  if (clockTimer) clearInterval(clockTimer);
  if (pdfBlobUrl.value) URL.revokeObjectURL(pdfBlobUrl.value);
});
</script>

<style scoped>
.card { background: #fff; border-radius: 14px; border: 1px solid #e9f5e9; box-shadow: 0 1px 3px rgba(13,40,24,0.05); overflow: hidden; }
.card-header { display: flex; align-items: center; justify-content: space-between; padding: 16px 20px; border-bottom: 1px solid #f0faf0; gap: 12px; }
.card-title  { font-size: 13.5px; font-weight: 700; color: #111827; margin: 0 0 2px; }
.today-label { font-size: 11.5px; color: #9ca3af; }
.header-actions { display: flex; align-items: center; gap: 8px; flex-shrink: 0; }

.btn-green-sm { background: #f0fdf4; color: #16a34a; border: 1.5px solid #bbf7d0; border-radius: 8px; padding: 6px 13px; font-size: 12px; font-weight: 600; cursor: pointer; font-family: inherit; text-decoration: none; display: inline-flex; align-items: center; gap: 5px; white-space: nowrap; }
.btn-green-sm:hover { background: #dcfce7; }
.btn-outline-sm { background: #f9fafb; color: #374151; border: 1.5px solid #e5e7eb; border-radius: 8px; padding: 6px 13px; font-size: 12px; font-weight: 600; cursor: pointer; font-family: inherit; text-decoration: none; display: inline-flex; align-items: center; gap: 5px; white-space: nowrap; }
.btn-outline-sm:hover { background: #f3f4f6; border-color: #d1d5db; }

.empty-state { display: flex; flex-direction: column; align-items: center; padding: 44px 24px; gap: 12px; text-align: center; }
.empty-state__icon { width: 64px; height: 64px; background: #f9fafb; border-radius: 50%; display: flex; align-items: center; justify-content: center; }
.empty-state p { font-size: 13px; color: #9ca3af; line-height: 1.7; margin: 0; }

.absensi-panel { display: flex; flex-direction: column; align-items: center; padding: 28px 24px 24px; gap: 10px; text-align: center; }
.absensi-panel--done { padding: 20px 24px; }

.ap-icon { width: 56px; height: 56px; border-radius: 50%; display: flex; align-items: center; justify-content: center; margin-bottom: 2px; }
.ap-icon--locked  { background: #f3f4f6; color: #9ca3af; }
.ap-icon--open    { background: #f0fdf4; color: #16a34a; }
.ap-icon--wait    { background: #fffbeb; color: #d97706; }
.ap-icon--blocked { background: #fff1f2; color: #be123c; }
.ap-desc--block { max-width: 300px; color: #dc2626; }
.btn-absen--refresh {
  background: #f3f4f6; color: #374151; border: 1.5px solid #e5e7eb;
  display: inline-flex; align-items: center; gap: 7px; margin-top: 4px;
}
.btn-absen--refresh:hover { background: #e5e7eb; }
.ap-icon--miss   { background: #fff7ed; color: #ea580c; }

.ap-title { font-size: 14px; font-weight: 700; color: #111827; }
.ap-desc  { font-size: 12.5px; color: #6b7280; }
.ap-countdown { background: #f0fdf4; border: 1px solid #bbf7d0; color: #16a34a; font-size: 12.5px; font-weight: 600; padding: 5px 16px; border-radius: 100px; margin-top: 4px; }
.ap-error { background: #fff1f2; border: 1px solid #fecdd3; color: #be123c; font-size: 12px; padding: 8px 14px; border-radius: 8px; width: 100%; max-width: 460px; text-align: left; line-height: 1.6; }
.ap-required { color: #ef4444; }

.ap-done-rows { display: flex; flex-direction: column; gap: 8px; align-self: stretch; }
.ap-done-row  { display: flex; align-items: center; gap: 10px; }
.ap-badge { font-size: 11px; font-weight: 700; padding: 4px 10px; border-radius: 100px; }
.ap-badge--masuk  { background: #f0fdf4; color: #16a34a; border: 1px solid #bbf7d0; }
.ap-badge--pulang { background: #f0fdf4; color: #1a5c20; border: 1px solid #bbf7d0; }
.ap-done-time { font-size: 13px; font-weight: 600; color: #374151; }

.ap-kegiatan { background: #f9fafb; border-radius: 10px; padding: 12px 14px; align-self: stretch; text-align: left; margin-top: 8px; }
.ap-kegiatan__label { font-size: 10.5px; font-weight: 600; color: #9ca3af; text-transform: uppercase; letter-spacing: .05em; margin-bottom: 6px; }
.ap-kegiatan__text  { font-size: 13px; color: #374151; line-height: 1.6; }

.kegiatan-ol { margin: 0; padding-left: 18px; display: flex; flex-direction: column; gap: 4px; }
.kegiatan-ol li { font-size: 13px; color: #374151; line-height: 1.6; padding-left: 2px; }

.ap-form { align-self: stretch; text-align: left; }
.ap-form-label { font-size: 12px; font-weight: 600; color: #374151; margin-bottom: 6px; }
.ap-form-hint  { font-size: 11px; color: #9ca3af; margin-top: 5px; }
.ap-textarea { width: 100%; border: 1.5px solid #e5e7eb; border-radius: 10px; padding: 10px 13px; font-size: 13px; font-family: inherit; resize: vertical; outline: none; color: #111827; box-sizing: border-box; line-height: 1.7; }
.ap-textarea:focus { border-color: #48AF4A; }

.kegiatan-ul { margin: 0; padding-left: 14px; list-style: disc; display: flex; flex-direction: column; gap: 2px; }
.kegiatan-ul li { font-size: 11.5px; color: #374151; line-height: 1.5; padding-left: 2px; }
.td-empty { color: #d1d5db; }

.ap-btn-row { display: flex; align-items: center; gap: 10px; margin-top: 4px; flex-wrap: wrap; justify-content: center; }
.btn-absen { background: #48AF4A; color: #fff; border: none; border-radius: 10px; padding: 11px 28px; font-size: 13.5px; font-weight: 700; cursor: pointer; display: flex; align-items: center; gap: 8px; font-family: inherit; }
.btn-absen:disabled { opacity: .6; cursor: not-allowed; }
.btn-absen--pulang { background: #1a5c20; }
.btn-ajukan { background: #f9fafb; color: #374151; border: 1.5px solid #e5e7eb; border-radius: 10px; padding: 11px 18px; font-size: 13px; font-weight: 600; cursor: pointer; display: flex; align-items: center; gap: 6px; font-family: inherit; white-space: nowrap; }
.btn-ajukan:hover { background: #f0fdf4; border-color: #bbf7d0; color: #16a34a; }
.btn-ajukan-sm { background: transparent; color: #6b7280; border: 1px dashed #d1d5db; border-radius: 8px; padding: 5px 14px; font-size: 12px; font-weight: 500; cursor: pointer; display: inline-flex; align-items: center; gap: 5px; font-family: inherit; margin-top: 8px; }
.btn-ajukan-sm:hover { background: #f9fafb; color: #374151; border-color: #9ca3af; }
.ap-error--with-retry { display: flex; flex-direction: column; gap: 8px; }
.ap-retry-btn { align-self: flex-start; background: #fff; color: #be123c; border: 1.5px solid #fecdd3; border-radius: 7px; padding: 5px 14px; font-size: 12px; font-weight: 600; cursor: pointer; font-family: inherit; }
.ap-retry-btn:hover { background: #fff1f2; }
.btn-spinner { width: 14px; height: 14px; border: 2px solid rgba(255,255,255,.3); border-top-color: #fff; border-radius: 50%; animation: spin .7s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }

.ap-toast { position: fixed; bottom: 24px; left: 50%; transform: translateX(-50%); background: #1a2e1a; color: #fff; font-size: 13px; font-weight: 600; padding: 10px 22px; border-radius: 100px; box-shadow: 0 4px 16px rgba(0,0,0,.18); z-index: 9999; white-space: nowrap; }
.toast-enter-active, .toast-leave-active { transition: all .3s; }
.toast-enter-from, .toast-leave-to { opacity: 0; transform: translate(-50%, 10px); }

/* ── Rekap stats ── */
.rekap-stats { display: flex; gap: 0; border-bottom: 1px solid #f0faf0; }
.stat-chip { flex: 1; display: flex; flex-direction: column; align-items: center; padding: 16px 8px; border-right: 1px solid #f0faf0; }
.stat-chip:last-child { border-right: none; }
.stat-num { font-size: 24px; font-weight: 800; line-height: 1; }
.stat-lbl { font-size: 10.5px; color: #9ca3af; margin-top: 3px; font-weight: 500; }
.stat-chip--hadir .stat-num { color: #16a34a; }
.stat-chip--izin  .stat-num { color: #ca8a04; }
.stat-chip--sakit .stat-num { color: #1a5c20; }
.stat-chip--alpha .stat-num { color: #dc2626; }

/* ── Absensi table ── */
.abs-table-wrap { overflow-x: auto; }
.abs-table { width: 100%; border-collapse: collapse; font-size: 12.5px; }
.abs-table th { padding: 9px 12px; font-size: 10px; font-weight: 600; color: #6b7280; background: #f9fafb; border-bottom: 1px solid #f1f5f9; text-transform: uppercase; letter-spacing: .04em; text-align: left; white-space: nowrap; }
.abs-table td { padding: 9px 12px; border-bottom: 1px solid #f9fafb; color: #374151; vertical-align: middle; }
.abs-table tbody tr:last-child td { border-bottom: none; }
.abs-table tbody tr:hover { background: #fafff9; }
.row-today { background: #f0fdf4 !important; }
.td-no { color: #9ca3af; font-size: 11px; text-align: center; width: 36px; }
.td-tanggal { font-weight: 600; white-space: nowrap; }
.td-hari { color: #6b7280; font-size: 11.5px; white-space: nowrap; }
.td-jam { font-size: 12px; color: #374151; font-feature-settings: 'tnum'; white-space: nowrap; }
.td-kegiatan { font-size: 11.5px; color: #6b7280; max-width: 200px; }

.ket-badge { font-size: 10.5px; font-weight: 600; padding: 2px 8px; border-radius: 100px; white-space: nowrap; }
.ket-badge--hadir { background: #f0fdf4; color: #16a34a; border: 1px solid #bbf7d0; }
.ket-badge--izin  { background: #fffbeb; color: #b45309; border: 1px solid #fde68a; }
.ket-badge--sakit { background: #f0fdf4; color: #0d2818; border: 1px solid #bbf7d0; }
.ket-badge--alpha { background: #fff1f2; color: #be123c; border: 1px solid #fecdd3; }
.ket-badge--belum { color: #d1d5db; }

/* ── Modal base ── */
.modal-backdrop { position: fixed; inset: 0; background: rgba(0,0,0,.45); display: flex; align-items: center; justify-content: center; z-index: 1000; padding: 16px; backdrop-filter: blur(2px); }
.modal-box { background: #fff; border-radius: 18px; width: 100%; max-width: 440px; padding: 24px; box-shadow: 0 20px 60px rgba(0,0,0,.18); }
.modal-box--izin { max-width: 460px; }
.modal-box--pdf { max-width: 900px; width: 95vw; padding: 0; overflow: hidden; border-radius: 16px; }
.modal-box--pdf .modal-box__header { padding: 14px 18px; border-bottom: 1px solid #f0f0f0; border-radius: 16px 16px 0 0; }
.pdf-modal-body { height: 78vh; background: #f3f4f6; }
.pdf-modal-iframe { width: 100%; height: 100%; border: none; display: block; }
.modal-box__header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 20px; }
.modal-box__title { font-size: 16px; font-weight: 700; color: #111827; }
.modal-close-btn { background: #f3f4f6; border: none; border-radius: 8px; width: 30px; height: 30px; display: flex; align-items: center; justify-content: center; cursor: pointer; color: #6b7280; flex-shrink: 0; }
.modal-close-btn:hover { background: #e5e7eb; }
.modal-box__desc { font-size: 13.5px; color: #374151; margin-bottom: 16px; line-height: 1.6; }
.modal-box__geo  { display: flex; align-items: center; gap: 6px; background: #f0fdf4; border-radius: 8px; padding: 8px 12px; font-size: 12px; color: #16a34a; margin-bottom: 16px; flex-wrap: wrap; }
.modal-box__acc { margin-left: auto; font-size: 11px; font-weight: 700; padding: 2px 8px; border-radius: 20px; white-space: nowrap; }
.modal-box__acc--good { background: #dcfce7; color: #16a34a; }
.modal-box__acc--ok   { background: #fef9c3; color: #854d0e; }
.modal-box__actions { display: flex; gap: 10px; justify-content: flex-end; margin-top: 6px; }

.modal-field { margin-bottom: 16px; }
.modal-label { display: block; font-size: 12px; font-weight: 600; color: #374151; margin-bottom: 6px; }
.modal-label--optional { font-size: 11px; font-weight: 400; color: #9ca3af; margin-left: 4px; }
.modal-hint { font-size: 11px; color: #9ca3af; margin-top: 4px; }
.modal-input { width: 100%; border: 1.5px solid #e5e7eb; border-radius: 10px; padding: 9px 13px; font-size: 13px; font-family: inherit; outline: none; color: #111827; box-sizing: border-box; }
.modal-input:focus { border-color: #48AF4A; }
.modal-textarea { width: 100%; border: 1.5px solid #e5e7eb; border-radius: 10px; padding: 10px 13px; font-size: 13px; font-family: inherit; resize: vertical; outline: none; color: #111827; box-sizing: border-box; }
.modal-textarea:focus { border-color: #48AF4A; }

.jenis-tabs { display: flex; gap: 8px; }
.jenis-tab { flex: 1; display: flex; align-items: center; justify-content: center; gap: 7px; border: 1.5px solid #e5e7eb; background: #f9fafb; color: #6b7280; font-size: 13px; font-weight: 600; padding: 10px 16px; border-radius: 10px; cursor: pointer; font-family: inherit; transition: all .15s; }
.jenis-tab:hover { border-color: #d1d5db; background: #f3f4f6; }
.jenis-tab--active { border-color: #48AF4A; background: #f0fdf4; color: #16a34a; }
.jenis-tab--sakit.jenis-tab--active { border-color: #3b82f6; background: #f0fdf4; color: #0d2818; }

.upload-zone { border: 2px dashed #e5e7eb; border-radius: 12px; padding: 20px; display: flex; flex-direction: column; align-items: center; gap: 8px; cursor: pointer; transition: all .2s; text-align: center; min-height: 90px; justify-content: center; }
.upload-zone:hover { border-color: #93c5fd; background: #f0f9ff; }
.upload-zone--has-file { border-color: #86efac; background: #f0fdf4; border-style: solid; }
.upload-zone__icon { color: #9ca3af; }
.upload-zone__text { font-size: 13px; font-weight: 600; color: #374151; }
.upload-zone__hint { font-size: 11px; color: #9ca3af; }
.upload-zone__preview { display: flex; align-items: center; gap: 8px; width: 100%; }
.upload-zone__filename { font-size: 12.5px; font-weight: 600; color: #16a34a; flex: 1; text-align: left; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.upload-zone__remove { background: #fee2e2; border: none; border-radius: 6px; width: 22px; height: 22px; display: flex; align-items: center; justify-content: center; cursor: pointer; color: #dc2626; flex-shrink: 0; }
.upload-zone__remove:hover { background: #fecaca; }

.btn-cancel  { background: #f3f4f6; color: #374151; border: none; border-radius: 10px; padding: 10px 20px; font-size: 13px; font-weight: 600; cursor: pointer; font-family: inherit; }
.btn-cancel:hover { background: #e5e7eb; }
.btn-confirm { background: #48AF4A; color: #fff; border: none; border-radius: 10px; padding: 10px 22px; font-size: 13px; font-weight: 600; cursor: pointer; font-family: inherit; display: flex; align-items: center; gap: 7px; }
.btn-confirm:disabled { opacity: .55; cursor: not-allowed; }
.btn-confirm:not(:disabled):hover { background: #48AF4A; }
.btn-confirm--pulang { background: #48AF4A; }
.btn-confirm--pulang:not(:disabled):hover { background: #3d9e3f; }
.btn-spinner-sm { width: 13px; height: 13px; border: 2px solid rgba(255,255,255,.35); border-top-color: #fff; border-radius: 50%; animation: spin .7s linear infinite; }

.modal-fade-enter-active, .modal-fade-leave-active { transition: all .25s; }
.modal-fade-enter-from, .modal-fade-leave-to { opacity: 0; }
.modal-fade-enter-from .modal-box, .modal-fade-leave-to .modal-box { transform: scale(.96) translateY(8px); }
.slide-down-enter-active, .slide-down-leave-active { transition: all .2s; overflow: hidden; }
.slide-down-enter-from, .slide-down-leave-to { opacity: 0; max-height: 0; margin-bottom: 0; }
.slide-down-enter-to, .slide-down-leave-from { opacity: 1; max-height: 200px; }
.spinner { width: 24px; height: 24px; border: 3px solid #e5e7eb; border-top-color: #48AF4A; border-radius: 50%; animation: spin .7s linear infinite; }

/* ── Surat Sakit Modal ── */
.bukti-overlay { position: fixed; inset: 0; z-index: 1200; background: rgba(0,0,0,.65); backdrop-filter: blur(3px); display: flex; align-items: center; justify-content: center; padding: 20px; }
.bukti-box { background: #fff; border-radius: 16px; box-shadow: 0 24px 60px rgba(0,0,0,.25); display: flex; flex-direction: column; max-width: 860px; width: 100%; max-height: 90vh; overflow: hidden; }
.bukti-header { display: flex; align-items: center; justify-content: space-between; padding: 14px 18px; border-bottom: 1px solid #f0f0f0; flex-shrink: 0; }
.bukti-title { font-size: 14px; font-weight: 700; color: #111827; }
.bukti-header-actions { display: flex; align-items: center; gap: 8px; }
.bukti-body { flex: 1; overflow: auto; display: flex; align-items: center; justify-content: center; background: #f9fafb; min-height: 300px; }
.bukti-img { max-width: 100%; max-height: 75vh; object-fit: contain; display: block; }
.bukti-iframe { width: 100%; height: 75vh; border: none; display: block; }
.bukti-unsupported { display: flex; flex-direction: column; align-items: center; gap: 12px; padding: 40px; color: #6b7280; font-size: 13px; }

/* ── Side Panel / Bottom Sheet ── */
.panel-wrapper { position: fixed; inset: 0; z-index: 1100; }
.panel-body {
  position: fixed;
  top: 0; right: 0; bottom: 0;
  width: 380px;
  max-width: 100vw;
  background: #fff;
  box-shadow: -4px 0 24px rgba(0,0,0,.12);
  display: flex;
  flex-direction: column;
  z-index: 1101;
}

@media (max-width: 640px) {
  .panel-wrapper::before {
    content: '';
    position: fixed;
    inset: 0;
    background: rgba(0,0,0,.4);
  }
  .panel-body {
    top: auto;
    left: 0; right: 0; bottom: 0;
    width: 100%;
    height: 82vh;
    border-radius: 20px 20px 0 0;
    box-shadow: 0 -4px 24px rgba(0,0,0,.15);
  }
}

.panel-header { display: flex; align-items: center; justify-content: space-between; padding: 16px 18px; border-bottom: 1px solid #f0f0f0; flex-shrink: 0; }
.panel-title { font-size: 14px; font-weight: 700; color: #111827; }
.panel-sub { font-size: 11.5px; color: #9ca3af; margin-top: 1px; }
.panel-header-right { display: flex; align-items: center; gap: 8px; }
.panel-content { flex: 1; overflow-y: auto; }

.izin-list { padding: 4px 0; }
.izin-row { padding: 14px 18px; border-bottom: 1px solid #f9fafb; }
.izin-row:last-child { border-bottom: none; }
.izin-row__top { display: flex; align-items: center; gap: 8px; flex-wrap: wrap; margin-bottom: 5px; }
.izin-badge { font-size: 10.5px; font-weight: 700; padding: 2px 8px; border-radius: 100px; white-space: nowrap; flex-shrink: 0; }
.izin-badge--izin  { background: #fffbeb; color: #b45309; border: 1px solid #fde68a; }
.izin-badge--sakit { background: #f0fdf4; color: #0d2818; border: 1px solid #bbf7d0; }
.izin-tanggal { font-size: 12.5px; font-weight: 600; color: #111827; flex: 1; }
.izin-status  { font-size: 10.5px; font-weight: 600; padding: 2px 8px; border-radius: 100px; white-space: nowrap; }
.izin-status--pending   { background: #fef9c3; color: #a16207; border: 1px solid #fde047; }
.izin-status--disetujui { background: #f0fdf4; color: #16a34a; border: 1px solid #86efac; }
.izin-status--ditolak   { background: #fff1f2; color: #be123c; border: 1px solid #fecdd3; }
.izin-alasan { font-size: 12px; color: #6b7280; margin-bottom: 4px; line-height: 1.5; }
.izin-catatan { display: flex; align-items: center; gap: 5px; font-size: 11px; color: #9ca3af; margin-bottom: 6px; }
.izin-bukti-row { margin-top: 6px; }
.bukti-link { display: inline-flex; align-items: center; gap: 5px; font-size: 11.5px; font-weight: 600; color: #1a5c20; background: #f0fdf4; border: 1px solid #bbf7d0; border-radius: 100px; padding: 3px 10px; cursor: pointer; font-family: inherit; }
.bukti-link:hover { background: #dcfce7; }

/* Panel transitions */
.panel-fade-enter-active, .panel-fade-leave-active { transition: opacity .25s; }
.panel-fade-enter-from, .panel-fade-leave-to { opacity: 0; }
.panel-slide-enter-active, .panel-slide-leave-active { transition: transform .3s cubic-bezier(.4,0,.2,1); }
.panel-slide-enter-from, .panel-slide-leave-to { transform: translateX(100%); }

@media (max-width: 640px) {
  .panel-slide-enter-from, .panel-slide-leave-to { transform: translateY(100%); }
}
</style>
