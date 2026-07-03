<template>
  <div class="rekap-root">

    <!-- ── PANEL IZIN & SAKIT (persetujuan pending) ────────────── -->
    <div class="card">
      <div class="card-header">
        <div>
          <h3 class="card-title">Persetujuan Izin &amp; Sakit</h3>
          <p class="card-sub">Pengajuan dari peserta yang menunggu konfirmasi</p>
        </div>
        <button class="btn-green-sm" @click="fetchIzinSakit">
          <svg width="13" height="13" viewBox="0 0 24 24" fill="none"><path d="M23 4v6h-6" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/><path d="M1 20v-6h6" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/><path d="M3.51 9a9 9 0 0114.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0020.49 15" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>
          Refresh
        </button>
      </div>

      <div class="filter-bar">
        <div class="filter-pills">
          <button v-for="f in izinFilters" :key="f.key"
            :class="['filter-pill', izinActiveFilter === f.key && 'filter-pill--active']"
            @click="izinActiveFilter = f.key">{{ f.label }}</button>
        </div>
      </div>

      <div v-if="izinLoading" class="empty-state"><div class="spinner"></div></div>
      <div v-else-if="izinError" class="empty-state"><p style="color:#dc2626">{{ izinError }}</p></div>
      <div v-else-if="filteredIzin.length === 0" class="empty-state">
        <div class="empty-state__icon">
          <svg width="32" height="32" viewBox="0 0 24 24" fill="none"><path d="M9 12l2 2 4-4" stroke="#d1d5db" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/><circle cx="12" cy="12" r="10" stroke="#d1d5db" stroke-width="1.5"/></svg>
        </div>
        <p>{{ izinActiveFilter === 'pending' ? 'Tidak ada pengajuan yang menunggu.' : 'Tidak ada data.' }}</p>
      </div>
      <div v-else class="table-wrap">
        <table class="data-table">
          <thead>
            <tr>
              <th>Peserta</th>
              <th>Tanggal</th>
              <th>Jenis</th>
              <th>Alasan</th>
              <th>Bukti</th>
              <th>Status</th>
              <th>Aksi</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="item in filteredIzin" :key="item.id">
              <td>
                <div class="name-cell">
                  <div class="name-avatar">{{ item.nama_peserta[0] }}</div>
                  <div>
                    <div class="name-text">{{ item.nama_peserta }}</div>
                    <div class="name-sub">{{ item.divisi ?? '–' }}</div>
                  </div>
                </div>
              </td>
              <td style="white-space:nowrap;font-weight:600">{{ fmtDate(item.tanggal) }}</td>
              <td>
                <span :class="['jenis-badge', `jenis-badge--${item.jenis}`]">
                  {{ item.jenis === 'izin' ? 'Izin' : 'Sakit' }}
                </span>
              </td>
              <td class="alasan-cell">{{ item.alasan }}</td>
              <td>
                <button v-if="item.bukti_path" class="bukti-link" @click="lihatBukti(item.bukti_path!)">
                  <svg width="11" height="11" viewBox="0 0 24 24" fill="none"><path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z" stroke="currentColor" stroke-width="2"/><polyline points="14 2 14 8 20 8" stroke="currentColor" stroke-width="2"/></svg>
                  Lihat
                </button>
                <span v-else class="name-sub">–</span>
              </td>
              <td>
                <span v-if="item.status === 'pending'"    class="status-badge status-badge--pending">Menunggu</span>
                <span v-else-if="item.status === 'disetujui'" class="status-badge status-badge--ok">Disetujui</span>
                <span v-else class="status-badge status-badge--tolak">Ditolak</span>
              </td>
              <td>
                <div v-if="item.status === 'pending'" class="aksi-cell">
                  <button class="btn-aksi btn-aksi--green" :disabled="processingId === item.id"
                    @click="approve(item)">
                    {{ processingId === item.id ? '…' : '✓ Setujui' }}
                  </button>
                  <button class="btn-aksi btn-aksi--red" :disabled="processingId === item.id"
                    @click="openTolakModal(item)">
                    ✕ Tolak
                  </button>
                </div>
                <span v-else class="name-sub">–</span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- ── REKAP ABSENSI PESERTA ────────────────────────────────── -->
    <div class="card">
      <div class="card-header">
        <div>
          <h3 class="card-title">Rekap Absensi Peserta</h3>
          <p class="card-sub">Ringkasan kehadiran seluruh peserta magang</p>
        </div>
        <button class="btn-green-sm" @click="fetchRekap">
          <svg width="13" height="13" viewBox="0 0 24 24" fill="none"><path d="M23 4v6h-6" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/><path d="M1 20v-6h6" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/><path d="M3.51 9a9 9 0 0114.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0020.49 15" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>
          Refresh
        </button>
      </div>

      <div class="filter-bar">
        <div class="filter-pills">
          <button v-for="f in filters" :key="f.key"
            :class="['filter-pill', activeFilter === f.key && 'filter-pill--active']"
            @click="activeFilter = f.key">{{ f.label }}</button>
        </div>
        <input v-model="search" type="text" class="search-input" placeholder="Cari nama peserta…"/>
      </div>

      <div v-if="!loading && rows.length" class="stat-chips">
        <div class="stat-chip stat-chip--green">
          <span class="stat-chip__val">{{ totalHadir }}</span>
          <span class="stat-chip__lbl">Total Hadir</span>
        </div>
        <div class="stat-chip stat-chip--yellow">
          <span class="stat-chip__val">{{ totalIzin }}</span>
          <span class="stat-chip__lbl">Izin</span>
        </div>
        <div class="stat-chip stat-chip--blue">
          <span class="stat-chip__val">{{ totalSakit }}</span>
          <span class="stat-chip__lbl">Sakit</span>
        </div>
        <div class="stat-chip stat-chip--red">
          <span class="stat-chip__val">{{ totalAlpha }}</span>
          <span class="stat-chip__lbl">Alpha</span>
        </div>
      </div>

      <div v-if="loading" class="empty-state"><div class="spinner"></div></div>
      <div v-else-if="error" class="empty-state">
        <p style="color:#dc2626">{{ error }}</p>
        <button class="btn-green-sm" @click="fetchRekap">Coba lagi</button>
      </div>
      <div v-else-if="filteredRows.length === 0" class="empty-state">
        <div class="empty-state__icon">
          <svg width="36" height="36" viewBox="0 0 24 24" fill="none"><rect x="3" y="4" width="18" height="18" rx="2" stroke="#d1d5db" stroke-width="1.5"/><line x1="3" y1="10" x2="21" y2="10" stroke="#d1d5db" stroke-width="1.5"/></svg>
        </div>
        <p>{{ search ? 'Tidak ada peserta yang cocok.' : 'Belum ada data absensi.' }}</p>
      </div>
      <div v-else class="table-wrap">
        <table class="data-table">
          <thead>
            <tr>
              <th>Peserta</th>
              <th>Divisi</th>
              <th>Periode</th>
              <th>Status</th>
              <th style="text-align:center">H</th>
              <th style="text-align:center">I</th>
              <th style="text-align:center">S</th>
              <th style="text-align:center">A</th>
              <th style="text-align:center">% Hadir</th>
              <th>Aksi</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="r in filteredRows" :key="r.pelaksanaan_id"
              :class="selectedRow?.pelaksanaan_id === r.pelaksanaan_id ? 'tr-selected' : ''">
              <td>
                <div class="name-cell">
                  <div class="name-avatar">{{ r.nama_lengkap[0] }}</div>
                  <div>
                    <div class="name-text">{{ r.nama_lengkap }}</div>
                    <div class="name-sub">{{ r.asal_institusi }}</div>
                  </div>
                </div>
              </td>
              <td>
                <span v-if="r.divisi" class="tag">{{ r.divisi }}</span>
                <span v-else class="name-sub">–</span>
              </td>
              <td style="white-space:nowrap;font-size:12px">
                <div>{{ fmtDate(r.tanggal_mulai) }}</div>
                <div class="name-sub">s/d {{ fmtDate(r.tanggal_selesai) }}</div>
              </td>
              <td><span :class="statusClass(r.status)">{{ fmtStatus(r.status) }}</span></td>
              <td style="text-align:center"><span class="abs-num abs-num--green">{{ r.hadir }}</span></td>
              <td style="text-align:center"><span class="abs-num abs-num--yellow">{{ r.izin }}</span></td>
              <td style="text-align:center"><span class="abs-num abs-num--blue">{{ r.sakit }}</span></td>
              <td style="text-align:center"><span class="abs-num abs-num--red">{{ r.alpha }}</span></td>
              <td style="text-align:center">
                <div class="pct-bar-wrap">
                  <div class="pct-bar"><div class="pct-bar-fill" :style="{ width: persen(r) + '%' }"></div></div>
                </div>
                <span class="pct-label">{{ persen(r) }}%</span>
              </td>
              <td>
                <div class="aksi-cell">
                  <button class="btn-aksi btn-aksi--ghost" @click="openDetail(r)">
                    <svg width="12" height="12" viewBox="0 0 24 24" fill="none"><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z" stroke="currentColor" stroke-width="2"/><circle cx="12" cy="12" r="3" stroke="currentColor" stroke-width="2"/></svg>
                    Detail
                  </button>
                  <button class="btn-aksi btn-aksi--blue"
                    :disabled="pdfLoadingId === r.pelaksanaan_id"
                    @click="openPDFModal(r.pelaksanaan_id)">
                    <span v-if="pdfLoadingId === r.pelaksanaan_id" class="btn-spinner"></span>
                    <svg v-else width="12" height="12" viewBox="0 0 24 24" fill="none"><path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z" stroke="currentColor" stroke-width="2"/><polyline points="14 2 14 8 20 8" stroke="currentColor" stroke-width="2"/></svg>
                    PDF
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

  </div>

  <!-- ── SIDE PANEL DETAIL ────────────────────────────────────── -->
  <Teleport to="body">
    <Transition name="side-panel">
      <div v-if="selectedRow" class="side-overlay" @click.self="selectedRow = null">
        <div class="side-panel">

          <!-- Header -->
          <div class="sp-header">
            <div class="sp-header__avatar">{{ selectedRow.nama_lengkap[0] }}</div>
            <div class="sp-header__info">
              <div class="sp-header__name">{{ selectedRow.nama_lengkap }}</div>
              <div class="sp-header__meta">
                <span v-if="selectedRow.divisi" class="sp-tag">{{ selectedRow.divisi }}</span>
                <span class="sp-header__periode">
                  {{ fmtDate(selectedRow.tanggal_mulai) }} – {{ fmtDate(selectedRow.tanggal_selesai) }}
                </span>
              </div>
              <!-- WA Pembimbing shortcut -->
              <a v-if="selectedRow.wa_pembimbing"
                :href="`https://wa.me/${selectedRow.wa_pembimbing.replace(/\D/g,'').replace(/^0/,'62')}`"
                target="_blank" rel="noopener" class="sp-wa-link" title="Konfirmasi ke pembimbing via WhatsApp">
                <svg width="13" height="13" viewBox="0 0 24 24" fill="none">
                  <path d="M17.472 14.382c-.297-.149-1.758-.867-2.03-.967-.273-.099-.471-.148-.67.15-.197.297-.767.966-.94 1.164-.173.199-.347.223-.644.075-.297-.15-1.255-.463-2.39-1.475-.883-.788-1.48-1.761-1.653-2.059-.173-.297-.018-.458.13-.606.134-.133.298-.347.446-.52.149-.174.198-.298.298-.497.099-.198.05-.371-.025-.52-.075-.149-.669-1.612-.916-2.207-.242-.579-.487-.5-.669-.51-.173-.008-.371-.01-.57-.01-.198 0-.52.074-.792.372-.272.297-1.04 1.016-1.04 2.479 0 1.462 1.065 2.875 1.213 3.074.149.198 2.096 3.2 5.077 4.487.709.306 1.262.489 1.694.625.712.227 1.36.195 1.871.118.571-.085 1.758-.719 2.006-1.413.248-.694.248-1.289.173-1.413-.074-.124-.272-.198-.57-.347z" fill="#25D366"/>
                  <path d="M12 2C6.477 2 2 6.477 2 12c0 1.89.525 3.66 1.438 5.168L2 22l4.954-1.418A9.955 9.955 0 0012 22c5.523 0 10-4.477 10-10S17.523 2 12 2z" stroke="#25D366" stroke-width="1.5"/>
                </svg>
                Konfirmasi Pembimbing
              </a>
            </div>
            <button class="sp-close" @click="selectedRow = null">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none"><path d="M18 6L6 18M6 6l12 12" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"/></svg>
            </button>
          </div>

          <!-- Rekap chips -->
          <div class="sp-rekap">
            <div class="sp-rekap__item sp-rekap__item--green">
              <div class="sp-rekap__num">{{ selectedRow.hadir }}</div>
              <div class="sp-rekap__lbl">Hadir</div>
            </div>
            <div class="sp-rekap__item sp-rekap__item--yellow">
              <div class="sp-rekap__num">{{ selectedRow.izin }}</div>
              <div class="sp-rekap__lbl">Izin</div>
            </div>
            <div class="sp-rekap__item sp-rekap__item--blue">
              <div class="sp-rekap__num">{{ selectedRow.sakit }}</div>
              <div class="sp-rekap__lbl">Sakit</div>
            </div>
            <div class="sp-rekap__item sp-rekap__item--red">
              <div class="sp-rekap__num">{{ selectedRow.alpha }}</div>
              <div class="sp-rekap__lbl">Alpha</div>
            </div>
          </div>

          <!-- Progress bar kehadiran -->
          <div class="sp-progress-section">
            <div class="sp-progress-label">
              <span>Persentase Kehadiran</span>
              <strong :class="persen(selectedRow) >= 80 ? 'pct-good' : 'pct-warn'">{{ persen(selectedRow) }}%</strong>
            </div>
            <div class="sp-progress-track">
              <div class="sp-progress-fill"
                :class="persen(selectedRow) >= 80 ? 'sp-progress-fill--green' : 'sp-progress-fill--yellow'"
                :style="{ width: persen(selectedRow) + '%' }"></div>
            </div>
          </div>

          <!-- Aksi PDF + Input Manual -->
          <div class="sp-actions">
            <button class="sp-btn-pdf"
              :disabled="pdfLoadingId === selectedRow.pelaksanaan_id"
              @click="openPDFModal(selectedRow.pelaksanaan_id)">
              <span v-if="pdfLoadingId === selectedRow.pelaksanaan_id" class="btn-spinner btn-spinner--green"></span>
              <template v-else>
                <svg width="13" height="13" viewBox="0 0 24 24" fill="none"><path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z" stroke="currentColor" stroke-width="2"/><polyline points="14 2 14 8 20 8" stroke="currentColor" stroke-width="2"/><line x1="16" y1="13" x2="8" y2="13" stroke="currentColor" stroke-width="2" stroke-linecap="round"/><line x1="16" y1="17" x2="8" y2="17" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>
                Unduh PDF Rekap
              </template>
            </button>
            <button class="sp-btn-manual" @click="openManualModal">
              <svg width="13" height="13" viewBox="0 0 24 24" fill="none"><path d="M12 5v14M5 12h14" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"/></svg>
              Input Absensi Manual
            </button>
          </div>

          <!-- Divider -->
          <div class="sp-divider">
            <span>Riwayat Absensi</span>
          </div>

          <!-- Tabel absensi harian (generate dari date range) -->
          <div v-if="detailLoading" class="sp-loading">
            <div class="spinner"></div>
            <span>Memuat data…</span>
          </div>
          <div v-else-if="detailError" class="sp-error-msg">{{ detailError }}</div>
          <div v-else-if="tabelHarian.length === 0" class="sp-empty">
            <svg width="36" height="36" viewBox="0 0 24 24" fill="none"><rect x="3" y="4" width="18" height="18" rx="2" stroke="#d1d5db" stroke-width="1.5"/><line x1="3" y1="10" x2="21" y2="10" stroke="#d1d5db" stroke-width="1.5"/></svg>
            <p>Belum ada catatan absensi</p>
          </div>
          <div v-else class="sp-table-wrap">
            <table class="sp-table">
              <thead>
                <tr>
                  <th style="width:32px">No</th>
                  <th>Tanggal</th>
                  <th style="width:34px">Hari</th>
                  <th style="width:52px">Masuk</th>
                  <th style="width:52px">Pulang</th>
                  <th style="width:62px">Status</th>
                  <th>Kegiatan</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="row in tabelHarian" :key="row.tanggal" :class="{ 'row-today': row.isToday }">
                  <td class="td-no">{{ row.no }}</td>
                  <td class="td-date">{{ fmtDateShort(row.tanggal) }}</td>
                  <td class="td-hari">{{ row.hari }}</td>
                  <td class="td-time">{{ row.jamMasuk }}</td>
                  <td class="td-time">{{ row.jamKeluar }}</td>
                  <td>
                    <div style="display:flex;flex-direction:column;gap:3px;align-items:flex-start">
                      <span v-if="row.status !== 'belum'" :class="['ket-badge', `ket-badge--${row.status}`]">
                        {{ ({ hadir:'Hadir', izin:'Izin', sakit:'Sakit', alpha:'Alpha' } as Record<string,string>)[row.status] ?? row.status }}
                      </span>
                      <span v-else class="ket-badge ket-badge--belum">–</span>
                      <span v-if="row.isManual" class="ket-badge-manual" :title="row.catatanManual ?? 'Input manual oleh HRD'">Manual</span>
                    </div>
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

        </div>
      </div>
    </Transition>
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

  <!-- ── TOAST ────────────────────────────────────────────────── -->
  <Teleport to="body">
    <Transition name="toast">
      <div v-if="toastMsg" class="toast-msg">{{ toastMsg }}</div>
    </Transition>
  </Teleport>

  <!-- ── MODAL TOLAK ─────────────────────────────────────────── -->
  <Teleport to="body">
    <div v-if="showTolakModal" class="modal-backdrop" @click.self="closeTolakModal">
      <div class="modal-box">
        <div class="modal-title">Tolak Pengajuan</div>
        <div class="modal-desc">
          Pengajuan <strong>{{ tolakTarget?.jenis }}</strong> dari <strong>{{ tolakTarget?.nama_peserta }}</strong>
          untuk tanggal <strong>{{ tolakTarget ? fmtDate(tolakTarget.tanggal) : '' }}</strong>.
        </div>
        <div class="modal-field">
          <label class="modal-label">Catatan (opsional)</label>
          <textarea v-model="tolakCatatan" class="modal-textarea" rows="3"
            placeholder="Tulis alasan penolakan jika diperlukan..."></textarea>
        </div>
        <div class="modal-actions">
          <button class="btn-cancel" @click="closeTolakModal">Batal</button>
          <button class="btn-red" @click="submitTolak" :disabled="processingId === tolakTarget?.id">
            {{ processingId === tolakTarget?.id ? '…' : 'Ya, Tolak' }}
          </button>
        </div>
      </div>
    </div>
  </Teleport>

  <!-- ── Modal Input Absensi Manual ───────────────────────── -->
  <Teleport to="body">
    <Transition name="modal-fade">
      <div v-if="showManualModal" class="modal-backdrop" @click.self="showManualModal = false">
        <div class="modal-box">
          <div class="modal-box__header">
            <div class="modal-box__title">Input Absensi Manual</div>
            <button class="modal-close-btn" @click="showManualModal = false">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none"><path d="M18 6L6 18M6 6l12 12" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"/></svg>
            </button>
          </div>
          <div class="modal-sub">Peserta: <strong>{{ selectedRow?.nama_lengkap }}</strong></div>

          <!-- Jenis toggle -->
          <div class="jenis-toggle">
            <button :class="['jenis-btn', manualForm.jenis === 'masuk' ? 'jenis-btn--active' : '']"
              @click="manualForm.jenis = 'masuk'; manualError = ''">
              <svg width="13" height="13" viewBox="0 0 24 24" fill="none"><path d="M15 3h4a2 2 0 012 2v14a2 2 0 01-2 2h-4" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/><polyline points="10 17 15 12 10 7" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/><line x1="15" y1="12" x2="3" y2="12" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>
              Jam Masuk
            </button>
            <button :class="['jenis-btn', manualForm.jenis === 'pulang' ? 'jenis-btn--active' : '']"
              @click="manualForm.jenis = 'pulang'; manualError = ''">
              <svg width="13" height="13" viewBox="0 0 24 24" fill="none"><path d="M9 21H5a2 2 0 01-2-2V5a2 2 0 012-2h4" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/><polyline points="16 17 21 12 16 7" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/><line x1="21" y1="12" x2="9" y2="12" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>
              Jam Pulang
            </button>
          </div>

          <div class="modal-field">
            <label class="modal-label">Tanggal <span style="color:#dc2626">*</span></label>
            <input v-model="manualForm.tanggal" type="date" class="modal-input" />
          </div>

          <!-- Masuk: hanya jam masuk -->
          <div v-if="manualForm.jenis === 'masuk'" class="modal-field">
            <label class="modal-label">Jam Masuk <span style="color:#dc2626">*</span></label>
            <input v-model="manualForm.jam_masuk" type="time" class="modal-input" />
          </div>

          <!-- Pulang: jam keluar + kegiatan -->
          <template v-if="manualForm.jenis === 'pulang'">
            <div class="modal-field">
              <label class="modal-label">Jam Keluar <span style="color:#dc2626">*</span></label>
              <input v-model="manualForm.jam_keluar" type="time" class="modal-input" />
            </div>
            <div class="modal-field">
              <label class="modal-label">Kegiatan</label>
              <textarea v-model="manualForm.kegiatan" class="modal-textarea" rows="3"
                placeholder="Tulis kegiatan yang dilakukan hari ini"></textarea>
            </div>
          </template>

          <div class="modal-field">
            <label class="modal-label">Catatan / Alasan</label>
            <input v-model="manualForm.catatan_manual" type="text" class="modal-input"
              placeholder="Contoh: Peserta lupa absen, konfirmasi via WA" />
          </div>
          <div v-if="manualError" class="modal-error">{{ manualError }}</div>
          <div class="modal-actions" style="margin-top:18px">
            <button class="btn-cancel" @click="showManualModal = false" :disabled="submittingManual">Batal</button>
            <button class="btn-confirm" @click="submitManual" :disabled="submittingManual || !manualForm.tanggal">
              <span v-if="submittingManual" class="btn-spinner btn-spinner--white"></span>
              {{ submittingManual ? 'Menyimpan...' : 'Simpan Absensi' }}
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>

  <!-- ── Modal Bukti Surat Sakit ───────────────────────────── -->
  <Teleport to="body">
    <div v-if="buktiModal.show" class="bukti-overlay" @click.self="closeBukti">
      <div class="bukti-box">
        <div class="bukti-header">
          <span class="bukti-title">Bukti Pengajuan</span>
          <div class="bukti-header-actions">
            <a :href="buktiModal.url" download class="btn-aksi btn-aksi--blue" title="Download">
              <svg width="13" height="13" viewBox="0 0 24 24" fill="none"><path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4" stroke="currentColor" stroke-width="2" stroke-linecap="round"/><polyline points="7 10 12 15 17 10" stroke="currentColor" stroke-width="2" stroke-linecap="round"/><line x1="12" y1="15" x2="12" y2="3" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>
              Download
            </a>
            <button class="sp-close" @click="closeBukti">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none"><path d="M18 6L6 18M6 6l12 12" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>
            </button>
          </div>
        </div>
        <div class="bukti-body">
          <img v-if="buktiModal.type === 'image'" :src="buktiModal.url" class="bukti-img" alt="Surat Sakit" />
          <iframe v-else-if="buktiModal.type === 'pdf'" :src="buktiModal.url" class="bukti-iframe" />
          <div v-else class="bukti-unsupported">
            <svg width="40" height="40" viewBox="0 0 24 24" fill="none"><path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z" stroke="#9ca3af" stroke-width="1.5"/><polyline points="14 2 14 8 20 8" stroke="#9ca3af" stroke-width="1.5"/></svg>
            <p>File tidak dapat ditampilkan di sini.</p>
            <a :href="buktiModal.url" download class="btn-aksi btn-aksi--blue">Download File</a>
          </div>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from "vue";
import api from "@/lib/api";
import { useAppWS } from "@/composables/useAppWS";

interface RekapRow {
  pelaksanaan_id:  string;
  nama_lengkap:    string;
  asal_institusi:  string;
  kategori_magang: string;
  divisi:          string | null;
  wa_pembimbing:   string | null;
  pembimbing:      string | null;
  tanggal_mulai:   string;
  tanggal_selesai: string;
  status:          string;
  hadir:           number;
  izin:            number;
  sakit:           number;
  alpha:           number;
  pending_approval: number;
}
interface AbsensiItem {
  id:             string;
  tanggal:        string;
  jam_masuk:      string | null;
  jam_keluar:     string | null;
  keterangan:     string;
  kegiatan:       string | null;
  is_manual:      boolean;
  catatan_manual: string | null;
}
interface IzinSakitItem {
  id:           string;
  tanggal:      string;
  jenis:        string;
  alasan:       string;
  status:       string;
  nama_peserta: string;
  divisi:       string | null;
  bukti_path:   string | null;
  catatan_hrd?: string | null;
}

const HARI = ['Min', 'Sen', 'Sel', 'Rab', 'Kam', 'Jum', 'Sab'];

// ── Rekap state ─────────────────────────────────────────────
const rows    = ref<RekapRow[]>([]);
const loading = ref(false);
const error   = ref<string | null>(null);
const search       = ref("");
const activeFilter = ref("semua");
const selectedRow  = ref<RekapRow | null>(null);
const detailList   = ref<AbsensiItem[]>([]);
const detailLoading = ref(false);
const detailError  = ref<string>("");

const filters = [
  { key: "semua",          label: "Semua" },
  { key: "aktif",          label: "Aktif" },
  { key: "upload_laporan", label: "Upload Laporan" },
  { key: "penilaian",      label: "Penilaian" },
  { key: "selesai",        label: "Selesai" },
];

// ── Izin/Sakit state ─────────────────────────────────────────
const izinList        = ref<IzinSakitItem[]>([]);
const izinLoading     = ref(false);
const izinError       = ref<string | null>(null);
const izinActiveFilter = ref("pending");
const processingId    = ref<string | null>(null);
const showTolakModal  = ref(false);
const tolakTarget     = ref<IzinSakitItem | null>(null);
const tolakCatatan    = ref("");

// ── Bukti modal state ─────────────────────────────────────────
const buktiModal = ref<{ show: boolean; url: string; type: 'image' | 'pdf' | 'other' }>({
  show: false, url: "", type: "other",
});

// ── PDF modal state ───────────────────────────────────────────
const showPDFModal  = ref(false);
const pdfBlobUrl    = ref<string>("");
const pdfLoadingId  = ref<string | null>(null);

// ── Input Absensi Manual state ────────────────────────────────
const showManualModal  = ref(false);
const submittingManual = ref(false);
const manualError      = ref("");
const manualForm = ref({
  jenis:          "masuk" as "masuk" | "pulang",
  tanggal:        "",
  jam_masuk:      "08:00",
  jam_keluar:     "16:00",
  kegiatan:       "",
  catatan_manual: "",
});

const izinFilters = [
  { key: "pending",   label: "Menunggu" },
  { key: "disetujui", label: "Disetujui" },
  { key: "ditolak",   label: "Ditolak" },
  { key: "semua",     label: "Semua" },
];

// ── Computed ─────────────────────────────────────────────────
const filteredRows = computed(() => {
  let list = rows.value;
  if (activeFilter.value !== "semua") list = list.filter(r => r.status === activeFilter.value);
  if (search.value.trim()) {
    const q = search.value.toLowerCase();
    list = list.filter(r => r.nama_lengkap.toLowerCase().includes(q) || (r.divisi ?? "").toLowerCase().includes(q));
  }
  return list;
});

const totalHadir = computed(() => filteredRows.value.reduce((s, r) => s + r.hadir, 0));
const totalIzin  = computed(() => filteredRows.value.reduce((s, r) => s + r.izin, 0));
const totalSakit = computed(() => filteredRows.value.reduce((s, r) => s + r.sakit, 0));
const totalAlpha = computed(() => filteredRows.value.reduce((s, r) => s + r.alpha, 0));

const filteredIzin = computed(() => {
  if (izinActiveFilter.value === "semua") return izinList.value;
  return izinList.value.filter(i => i.status === izinActiveFilter.value);
});

// ── Generate tabel harian dari date range ────────────────────
const tabelHarian = computed(() => {
  if (!selectedRow.value) return [];
  const mulaiRaw   = selectedRow.value.tanggal_mulai;
  const selesaiRaw = selectedRow.value.tanggal_selesai;
  if (!mulaiRaw || !selesaiRaw) return [];

  const mulaiD   = new Date(mulaiRaw);
  const selesaiD = new Date(selesaiRaw);
  const todayStr = new Date().toISOString().slice(0, 10);
  const result: any[] = [];
  let no = 1;

  const cur = new Date(Date.UTC(mulaiD.getUTCFullYear(), mulaiD.getUTCMonth(), mulaiD.getUTCDate()));
  const end = new Date(Date.UTC(selesaiD.getUTCFullYear(), selesaiD.getUTCMonth(), selesaiD.getUTCDate()));

  while (cur <= end) {
    const dow = cur.getUTCDay();
    if (dow !== 0 && dow !== 6) {
      const dateStr = cur.toISOString().slice(0, 10);
      const absensi = detailList.value.find(a => (a.tanggal || '').slice(0, 10) === dateStr);
      let status = dateStr < todayStr ? 'alpha' : 'belum';
      if (absensi) status = absensi.keterangan;
      result.push({
        no:        no++,
        tanggal:   dateStr,
        hari:      HARI[dow],
        jamMasuk:  fmtJam(absensi?.jam_masuk),
        jamKeluar: fmtJam(absensi?.jam_keluar),
        kegiatan:  absensi?.kegiatan || '',
        status,
        isToday:   dateStr === todayStr,
        isManual:  absensi?.is_manual ?? false,
        catatanManual: absensi?.catatan_manual ?? null,
      });
    }
    cur.setUTCDate(cur.getUTCDate() + 1);
  }
  return result;
});

// ── Data fetch ──────────────────────────────────────────────
async function fetchRekap() {
  loading.value = true; error.value = null;
  try {
    const r = await api.get("/api/absensi/rekap");
    rows.value = Array.isArray(r.data?.data) ? r.data.data : [];
    // Sinkronisasi selectedRow agar counter hadir/izin/sakit/alpha ikut update
    if (selectedRow.value) {
      const updated = rows.value.find(r => r.pelaksanaan_id === selectedRow.value!.pelaksanaan_id);
      if (updated) selectedRow.value = updated;
    }
  } catch (e: any) {
    error.value = e.response?.data?.message ?? "Gagal memuat data rekap";
  } finally { loading.value = false; }
}

async function fetchIzinSakit() {
  izinLoading.value = true; izinError.value = null;
  try {
    const r = await api.get("/api/izin-sakit");
    izinList.value = Array.isArray(r.data?.data) ? r.data.data : [];
  } catch (e: any) {
    izinError.value = e.response?.data?.message ?? "Gagal memuat data izin/sakit";
  } finally { izinLoading.value = false; }
}

async function openDetail(row: RekapRow) {
  selectedRow.value = row;
  detailList.value  = [];
  detailLoading.value = true;
  detailError.value   = "";
  try {
    const r = await api.get(`/api/absensi/pelaksanaan/${row.pelaksanaan_id}`);
    const data = r.data?.data;
    detailList.value = Array.isArray(data?.list) ? data.list : (Array.isArray(data) ? data : []);
  } catch (e: any) {
    detailError.value = e.response?.data?.message ?? "Gagal memuat detail absensi";
  } finally { detailLoading.value = false; }
}

// ── Input Absensi Manual ─────────────────────────────────────
function openManualModal() {
  manualForm.value = {
    jenis:          "masuk",
    tanggal:        new Date().toISOString().slice(0, 10),
    jam_masuk:      "08:00",
    jam_keluar:     "16:00",
    kegiatan:       "",
    catatan_manual: "",
  };
  manualError.value    = "";
  showManualModal.value = true;
}

async function submitManual() {
  if (!selectedRow.value) return;
  const jenis = manualForm.value.jenis;
  if (!manualForm.value.tanggal) {
    manualError.value = "Tanggal wajib diisi";
    return;
  }
  if (jenis === "masuk" && !manualForm.value.jam_masuk) {
    manualError.value = "Jam masuk wajib diisi";
    return;
  }
  if (jenis === "pulang" && !manualForm.value.jam_keluar) {
    manualError.value = "Jam keluar wajib diisi";
    return;
  }
  submittingManual.value = true;
  manualError.value = "";
  try {
    if (jenis === "masuk") {
      await api.post("/api/absensi/manual", {
        pelaksanaan_id: selectedRow.value.pelaksanaan_id,
        tanggal:        manualForm.value.tanggal,
        jam_masuk:      manualForm.value.jam_masuk,
        catatan_manual: manualForm.value.catatan_manual.trim(),
      });
    } else {
      await api.patch("/api/absensi/manual-pulang", {
        pelaksanaan_id: selectedRow.value.pelaksanaan_id,
        tanggal:        manualForm.value.tanggal,
        jam_keluar:     manualForm.value.jam_keluar,
        kegiatan:       manualForm.value.kegiatan.trim(),
        catatan_manual: manualForm.value.catatan_manual.trim(),
      });
    }
    showManualModal.value = false;
    showToast("Absensi manual berhasil diinput");
    await openDetail(selectedRow.value);
    fetchRekap();
  } catch (e: any) {
    manualError.value = e.response?.data?.message ?? "Gagal menyimpan absensi manual";
  } finally {
    submittingManual.value = false;
  }
}

// ── PDF modal ────────────────────────────────────────────────
async function openPDFModal(pelaksanaanId: string) {
  if (pdfLoadingId.value === pelaksanaanId) return;
  pdfLoadingId.value = pelaksanaanId;
  try {
    const r = await api.get(`/api/absensi/pelaksanaan/${pelaksanaanId}/pdf`, { responseType: 'blob' });
    if (pdfBlobUrl.value) URL.revokeObjectURL(pdfBlobUrl.value);
    pdfBlobUrl.value = URL.createObjectURL(new Blob([r.data], { type: 'application/pdf' }));
    showPDFModal.value = true;
  } catch {
    showToast('Gagal memuat PDF. Coba lagi.');
  } finally {
    pdfLoadingId.value = null;
  }
}

function closePDFModal() {
  showPDFModal.value = false;
}

// ── Izin/Sakit actions ───────────────────────────────────────
async function approve(item: IzinSakitItem) {
  processingId.value = item.id;
  try {
    await api.patch(`/api/izin-sakit/${item.id}/approve`, {});
    item.status = "disetujui";
    showToast(`Pengajuan ${item.jenis} ${item.nama_peserta} disetujui`);
    fetchRekap();
  } catch (e: any) {
    alert(e.response?.data?.message ?? "Gagal menyetujui");
  } finally { processingId.value = null; }
}

function openTolakModal(item: IzinSakitItem) {
  tolakTarget.value  = item;
  tolakCatatan.value = "";
  showTolakModal.value = true;
}

function closeTolakModal() {
  showTolakModal.value = false;
  tolakTarget.value    = null;
  tolakCatatan.value   = "";
}

async function submitTolak() {
  if (!tolakTarget.value) return;
  processingId.value = tolakTarget.value.id;
  try {
    await api.patch(`/api/izin-sakit/${tolakTarget.value.id}/tolak`, { catatan_hrd: tolakCatatan.value });
    tolakTarget.value.status = "ditolak";
    closeTolakModal();
    showToast("Pengajuan ditolak");
  } catch (e: any) {
    alert(e.response?.data?.message ?? "Gagal menolak");
  } finally { processingId.value = null; }
}

// ── Bukti modal ──────────────────────────────────────────────
function lihatBukti(path: string) {
  const ext  = path.split(".").pop()?.toLowerCase() ?? "";
  const type = ["jpg","jpeg","png","gif","webp"].includes(ext) ? "image"
             : ext === "pdf" ? "pdf" : "other";
  buktiModal.value = { show: true, url: `/uploads/${path}`, type };
}
function closeBukti() {
  buktiModal.value = { show: false, url: "", type: "other" };
}

// ── Toast ────────────────────────────────────────────────────
const toastMsg = ref("");
function showToast(msg: string) {
  toastMsg.value = msg;
  setTimeout(() => { toastMsg.value = ""; }, 3000);
}

// ── Helpers ─────────────────────────────────────────────────
function persen(r: RekapRow) {
  const total = r.hadir + r.izin + r.sakit + r.alpha;
  if (!total) return 0;
  return Math.round((r.hadir / total) * 100);
}
function fmtDate(d: string) {
  if (!d) return '–';
  return new Date(d).toLocaleDateString("id-ID", { day:"2-digit", month:"short", year:"numeric" });
}
function fmtDateShort(d: string) {
  if (!d) return '–';
  return new Date(d + 'T00:00:00').toLocaleDateString("id-ID", { day:"2-digit", month:"short" });
}
function fmtJam(t: string | null | undefined): string {
  if (!t) return '–';
  return t.slice(0, 5);
}
function kegiatanPoin(text: string | null | undefined): string[] {
  if (!text || text === '–') return [];
  return text.split('\n').map(l => l.trim()).filter(l => l.length > 0);
}
function fmtStatus(s: string) {
  return ({ aktif:"Aktif", upload_laporan:"Upload Lap.", penilaian:"Penilaian", selesai:"Selesai", menunggu_mulai:"Belum Mulai" } as Record<string,string>)[s] ?? s;
}
function statusClass(s: string) {
  if (s === "aktif")          return "sp-badge sp-badge--green";
  if (s === "upload_laporan") return "sp-badge sp-badge--blue";
  if (s === "penilaian")      return "sp-badge sp-badge--orange";
  if (s === "selesai")        return "sp-badge sp-badge--gray";
  return "sp-badge sp-badge--gray";
}

// ── WebSocket realtime ────────────────────────────────────────
const { connect: wsConnect, disconnect: wsDisconnect, subscribe: wsSubscribe } = useAppWS();
let wsUnsub: (() => void) | null = null;

onMounted(() => {
  fetchRekap();
  fetchIzinSakit();

  wsConnect();
  wsUnsub = wsSubscribe((msg: any) => {
    if (msg.type === 'notifikasi' && msg.data?.tipe === 'absensi_manual') {
      fetchRekap();
      if (selectedRow.value) {
        openDetail(selectedRow.value);
      }
    }
  });
});

onUnmounted(() => {
  if (wsUnsub) wsUnsub();
  wsDisconnect();
});
</script>

<style scoped>
.rekap-root { display: flex; flex-direction: column; gap: 16px; }

.card { background: #fff; border-radius: 14px; border: 1px solid #e9f5e9; box-shadow: 0 1px 3px rgba(13,40,24,0.05); overflow: hidden; }
.card-header { display: flex; align-items: center; justify-content: space-between; padding: 16px 20px; border-bottom: 1px solid #f0faf0; gap: 12px; flex-wrap: wrap; }
.card-title { font-size: 13.5px; font-weight: 700; color: #111827; margin: 0; }
.card-sub { font-size: 11.5px; color: #9ca3af; margin: 2px 0 0; }

.filter-bar { display: flex; align-items: center; justify-content: space-between; gap: 12px; padding: 12px 20px; border-bottom: 1px solid #f9fafb; flex-wrap: wrap; }
.filter-pills { display: flex; gap: 6px; flex-wrap: wrap; }
.filter-pill { border: 1.5px solid #e5e7eb; background: #fff; color: #6b7280; font-size: 12px; font-weight: 600; padding: 5px 13px; border-radius: 100px; cursor: pointer; font-family: inherit; transition: all .15s; }
.filter-pill--active { border-color: #48AF4A; background: #f0fdf4; color: #16a34a; }
.filter-pill:hover:not(.filter-pill--active) { border-color: #d1d5db; background: #f9fafb; }
.search-input { border: 1.5px solid #e5e7eb; border-radius: 9px; padding: 7px 13px; font-size: 12.5px; font-family: inherit; outline: none; color: #111827; min-width: 200px; }
.search-input:focus { border-color: #48AF4A; }

.stat-chips { display: flex; gap: 10px; padding: 12px 20px; flex-wrap: wrap; }
.stat-chip { display: flex; flex-direction: column; align-items: center; background: #f9fafb; border-radius: 10px; padding: 8px 18px; border: 1px solid #f1f5f9; min-width: 70px; }
.stat-chip__val { font-size: 20px; font-weight: 700; }
.stat-chip__lbl { font-size: 10.5px; color: #9ca3af; font-weight: 500; }
.stat-chip--green .stat-chip__val { color: #16a34a; }
.stat-chip--yellow .stat-chip__val { color: #ca8a04; }
.stat-chip--blue .stat-chip__val { color: #1a5c20; }
.stat-chip--red .stat-chip__val { color: #dc2626; }

.table-wrap { overflow-x: auto; }
.data-table { width: 100%; border-collapse: collapse; font-size: 13px; }
.data-table th { padding: 10px 14px; text-align: left; font-size: 10.5px; font-weight: 600; color: #6b7280; background: #f9fafb; border-bottom: 1px solid #f1f5f9; text-transform: uppercase; letter-spacing: .04em; white-space: nowrap; }
.data-table td { padding: 12px 14px; border-bottom: 1px solid #f9fafb; color: #374151; vertical-align: middle; }
.tr-selected td { background: #f0fdf4; }

.name-cell { display: flex; align-items: center; gap: 10px; }
.name-avatar { width: 32px; height: 32px; border-radius: 50%; background: linear-gradient(135deg,#48AF4A,#2d7a2e); color: #fff; font-size: 13px; font-weight: 700; display: flex; align-items: center; justify-content: center; flex-shrink: 0; }
.name-text { font-size: 13px; font-weight: 600; color: #111827; }
.name-sub  { font-size: 11.5px; color: #9ca3af; }
.tag { background: #f0fdf4; color: #16a34a; border: 1px solid #bbf7d0; border-radius: 100px; font-size: 11.5px; font-weight: 600; padding: 3px 9px; }

.abs-num { font-size: 14px; font-weight: 700; padding: 2px 8px; border-radius: 6px; }
.abs-num--green  { color: #16a34a; background: #f0fdf4; }
.abs-num--yellow { color: #ca8a04; background: #fefce8; }
.abs-num--blue   { color: #1a5c20; background: #f0fdf4; }
.abs-num--red    { color: #dc2626; background: #fff1f2; }

.pct-bar-wrap { width: 60px; margin: 0 auto 3px; }
.pct-bar  { height: 5px; background: #f1f5f9; border-radius: 100px; overflow: hidden; }
.pct-bar-fill { height: 100%; background: linear-gradient(90deg, #48AF4A, #16a34a); border-radius: 100px; transition: width .4s; }
.pct-label { font-size: 11px; font-weight: 600; color: #374151; }

.sp-badge { font-size: 11px; font-weight: 600; padding: 3px 9px; border-radius: 100px; }
.sp-badge--green  { background: #f0fdf4; color: #16a34a; border: 1px solid #86efac; }
.sp-badge--blue   { background: #f0fdf4; color: #0d2818; border: 1px solid #bbf7d0; }
.sp-badge--orange { background: #fff7ed; color: #c2410c; border: 1px solid #fed7aa; }
.sp-badge--gray   { background: #f9fafb; color: #6b7280; border: 1px solid #e5e7eb; }

.aksi-cell { display: flex; gap: 6px; align-items: center; }
.btn-aksi { font-size: 11.5px; font-weight: 600; padding: 5px 11px; border-radius: 7px; cursor: pointer; font-family: inherit; border: 1.5px solid transparent; display: inline-flex; align-items: center; gap: 4px; white-space: nowrap; }
.btn-aksi:disabled { opacity: .5; cursor: not-allowed; }
.btn-aksi--ghost { background: #f9fafb; color: #374151; border-color: #e5e7eb; }
.btn-aksi--ghost:hover:not(:disabled) { background: #f0fdf4; color: #16a34a; border-color: #bbf7d0; }
.btn-aksi--green { background: #f0fdf4; color: #16a34a; border-color: #86efac; }
.btn-aksi--green:hover:not(:disabled) { background: #dcfce7; }
.btn-aksi--red   { background: #fff1f2; color: #be123c; border-color: #fecdd3; }
.btn-aksi--red:hover:not(:disabled) { background: #ffe4e6; }
.btn-aksi--blue  { background: #f0fdf4; color: #0d2818; border-color: #bbf7d0; text-decoration: none; }
.btn-aksi--blue:hover:not(:disabled) { background: #dcfce7; }

.btn-spinner { width: 11px; height: 11px; border: 2px solid rgba(29,78,216,.2); border-top-color: #0d2818; border-radius: 50%; animation: spin .7s linear infinite; display: inline-block; }
.btn-spinner--green { border-color: rgba(21,128,61,.2); border-top-color: #16a34a; }

.jenis-badge { font-size: 11px; font-weight: 700; padding: 3px 9px; border-radius: 100px; }
.jenis-badge--izin  { background: #fffbeb; color: #b45309; border: 1px solid #fde68a; }
.jenis-badge--sakit { background: #f0fdf4; color: #0d2818; border: 1px solid #bbf7d0; }

.alasan-cell { font-size: 12.5px; color: #6b7280; max-width: 200px; }

.status-badge { font-size: 11px; font-weight: 600; padding: 3px 9px; border-radius: 100px; }
.status-badge--pending { background: #fef9c3; color: #a16207; border: 1px solid #fde047; }
.status-badge--ok      { background: #f0fdf4; color: #16a34a; border: 1px solid #86efac; }
.status-badge--tolak   { background: #fff1f2; color: #be123c; border: 1px solid #fecdd3; }

.bukti-link { display: inline-flex; align-items: center; gap: 4px; font-size: 11px; font-weight: 600; color: #1a5c20; background: #f0fdf4; border: 1px solid #bbf7d0; border-radius: 100px; padding: 2px 8px; white-space: nowrap; cursor: pointer; font-family: inherit; }
.bukti-link:hover { background: #dcfce7; }

/* ── Bukti Modal ── */
.bukti-overlay { position: fixed; inset: 0; z-index: 600; background: rgba(0,0,0,.65); backdrop-filter: blur(3px); display: flex; align-items: center; justify-content: center; padding: 20px; }
.bukti-box { background: #fff; border-radius: 16px; box-shadow: 0 24px 60px rgba(0,0,0,.25); display: flex; flex-direction: column; max-width: 860px; width: 100%; max-height: 90vh; overflow: hidden; }
.bukti-header { display: flex; align-items: center; justify-content: space-between; padding: 16px 20px; border-bottom: 1px solid #f0faf0; flex-shrink: 0; }
.bukti-title { font-size: 14px; font-weight: 700; color: #111827; }
.bukti-header-actions { display: flex; align-items: center; gap: 8px; }
.bukti-body { flex: 1; overflow: auto; display: flex; align-items: center; justify-content: center; background: #f9fafb; min-height: 300px; }
.bukti-img { max-width: 100%; max-height: 75vh; object-fit: contain; display: block; }
.bukti-iframe { width: 100%; height: 75vh; border: none; display: block; }
.bukti-unsupported { display: flex; flex-direction: column; align-items: center; gap: 12px; padding: 40px; color: #6b7280; font-size: 13px; }

/* ── PDF Modal ── */
.modal-box--pdf { max-width: 900px; width: 95vw; height: 90vh; display: flex; flex-direction: column; padding: 0; }
.modal-box__header { display: flex; align-items: center; justify-content: space-between; padding: 16px 20px; border-bottom: 1px solid #f0faf0; flex-shrink: 0; }
.modal-box__title { font-size: 14px; font-weight: 700; color: #111827; }
.modal-close-btn { background: #f3f4f6; border: none; border-radius: 8px; width: 32px; height: 32px; display: flex; align-items: center; justify-content: center; cursor: pointer; color: #6b7280; }
.modal-close-btn:hover { background: #e5e7eb; color: #111827; }
.pdf-modal-body { flex: 1; overflow: hidden; background: #f3f4f6; }
.pdf-modal-iframe { width: 100%; height: 100%; border: none; display: block; }
.btn-confirm { background: #48AF4A; color: #fff; border: none; border-radius: 8px; padding: 8px 16px; font-size: 13px; font-weight: 600; cursor: pointer; font-family: inherit; display: inline-flex; align-items: center; gap: 6px; }
.btn-confirm:hover { background: #48AF4A; }

.empty-state { display: flex; flex-direction: column; align-items: center; padding: 40px 24px; gap: 10px; text-align: center; }
.empty-state__icon { width: 64px; height: 64px; background: #f9fafb; border-radius: 50%; display: flex; align-items: center; justify-content: center; }
.empty-state p { font-size: 13px; color: #9ca3af; margin: 0; }

/* ── Side Panel ─────────────────────────────────────────────── */
.side-overlay {
  position: fixed; inset: 0; z-index: 500;
  background: rgba(0, 0, 0, 0.35);
  backdrop-filter: blur(2px);
  display: flex; justify-content: flex-end;
}
.side-panel {
  width: 520px; max-width: 95vw;
  height: 100%;
  background: #fff;
  display: flex; flex-direction: column;
  box-shadow: -8px 0 40px rgba(0,0,0,0.14);
  overflow: hidden;
}

.sp-header {
  display: flex; align-items: center; gap: 14px;
  padding: 20px 22px 18px;
  border-bottom: 1px solid #f0faf0;
  background: linear-gradient(135deg, #f0fdf4 0%, #ffffff 60%);
  flex-shrink: 0;
}
.sp-header__avatar {
  width: 44px; height: 44px; border-radius: 50%;
  background: linear-gradient(135deg, #48AF4A, #2d7a2e);
  color: #fff; font-size: 18px; font-weight: 700;
  display: flex; align-items: center; justify-content: center;
  flex-shrink: 0; box-shadow: 0 2px 8px rgba(72,175,74,.3);
}
.sp-header__info { flex: 1; min-width: 0; }
.sp-header__name { font-size: 15px; font-weight: 700; color: #111827; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.sp-header__meta { display: flex; align-items: center; gap: 8px; margin-top: 4px; flex-wrap: wrap; }
.sp-tag { background: #f0fdf4; color: #16a34a; border: 1px solid #bbf7d0; border-radius: 100px; font-size: 11px; font-weight: 600; padding: 2px 8px; }
.sp-wa-link { display: inline-flex; align-items: center; gap: 5px; margin-top: 6px; font-size: 11.5px; font-weight: 600; color: #16a34a; background: #f0fdf4; border: 1px solid #bbf7d0; border-radius: 100px; padding: 3px 10px 3px 7px; text-decoration: none; transition: background 0.15s, color 0.15s; }
.sp-wa-link:hover { background: #dcfce7; color: #15803d; }
.sp-header__periode { font-size: 11.5px; color: #9ca3af; }
.sp-close {
  background: #f3f4f6; border: none; border-radius: 9px;
  width: 34px; height: 34px; display: flex; align-items: center; justify-content: center;
  cursor: pointer; color: #6b7280; flex-shrink: 0; transition: all .15s;
}
.sp-close:hover { background: #e5e7eb; color: #111827; }

.sp-rekap {
  display: grid; grid-template-columns: repeat(4, 1fr);
  border-bottom: 1px solid #f0faf0; flex-shrink: 0;
}
.sp-rekap__item {
  display: flex; flex-direction: column; align-items: center;
  padding: 14px 8px; border-right: 1px solid #f0faf0;
}
.sp-rekap__item:last-child { border-right: none; }
.sp-rekap__num { font-size: 22px; font-weight: 800; }
.sp-rekap__lbl { font-size: 10.5px; color: #9ca3af; font-weight: 500; margin-top: 2px; }
.sp-rekap__item--green .sp-rekap__num { color: #16a34a; }
.sp-rekap__item--yellow .sp-rekap__num { color: #ca8a04; }
.sp-rekap__item--blue .sp-rekap__num { color: #1a5c20; }
.sp-rekap__item--red .sp-rekap__num { color: #dc2626; }

.sp-progress-section { padding: 14px 22px 10px; flex-shrink: 0; }
.sp-progress-label { display: flex; justify-content: space-between; align-items: center; font-size: 12px; color: #6b7280; margin-bottom: 7px; }
.pct-good { color: #16a34a; }
.pct-warn { color: #d97706; }
.sp-progress-track { height: 7px; background: #f1f5f9; border-radius: 100px; overflow: hidden; }
.sp-progress-fill { height: 100%; border-radius: 100px; transition: width .5s cubic-bezier(.4,0,.2,1); }
.sp-progress-fill--green  { background: linear-gradient(90deg, #48AF4A, #16a34a); }
.sp-progress-fill--yellow { background: linear-gradient(90deg, #f59e0b, #d97706); }

.sp-actions { padding: 10px 22px 14px; flex-shrink: 0; }
.sp-btn-pdf {
  display: inline-flex; align-items: center; gap: 7px;
  background: #f0fdf4; color: #16a34a;
  border: 1.5px solid #bbf7d0; border-radius: 9px;
  font-size: 12.5px; font-weight: 600; padding: 8px 16px;
  cursor: pointer; font-family: inherit; transition: all .15s;
}
.sp-btn-pdf:hover:not(:disabled) { background: #dcfce7; border-color: #86efac; }
.sp-btn-pdf:disabled { opacity: .55; cursor: not-allowed; }
.sp-btn-manual {
  display: inline-flex; align-items: center; gap: 7px;
  background: #fff; color: #374151;
  border: 1.5px solid #e5e7eb; border-radius: 9px;
  font-size: 12.5px; font-weight: 600; padding: 8px 16px;
  cursor: pointer; font-family: inherit; transition: all .15s;
}
.sp-btn-manual:hover { background: #f0fdf4; color: #16a34a; border-color: #bbf7d0; }

.sp-divider {
  display: flex; align-items: center; padding: 0 22px;
  flex-shrink: 0;
}
.sp-divider::before, .sp-divider::after {
  content: ''; flex: 1; height: 1px; background: #f0faf0;
}
.sp-divider span {
  font-size: 10.5px; font-weight: 700; color: #9ca3af;
  text-transform: uppercase; letter-spacing: .06em;
  padding: 0 12px; white-space: nowrap;
}

/* ── Tabel harian ── */
.sp-table-wrap { flex: 1; overflow-y: auto; padding-bottom: 16px; }
.sp-table { width: 100%; border-collapse: collapse; font-size: 12px; }
.sp-table th {
  position: sticky; top: 0; z-index: 1;
  padding: 8px 10px;
  font-size: 9.5px; font-weight: 700; color: #9ca3af;
  background: #f9fafb; border-bottom: 1px solid #f1f5f9;
  text-transform: uppercase; letter-spacing: .05em; text-align: left;
  white-space: nowrap;
}
.sp-table td { padding: 8px 10px; border-bottom: 1px solid #f9fafb; color: #374151; vertical-align: middle; }
.sp-table tr:last-child td { border-bottom: none; }
.sp-table tr.row-today td { background: #fffbeb; }
.td-no { font-size: 11px; color: #9ca3af; text-align: center; }
.td-date { font-weight: 600; white-space: nowrap; }
.td-hari { font-size: 11px; color: #6b7280; }
.td-time { color: #6b7280; white-space: nowrap; }
.td-kegiatan { max-width: 120px; }
.kegiatan-ul { margin: 0; padding-left: 14px; list-style: disc; }
.kegiatan-ul li { font-size: 11.5px; color: #6b7280; line-height: 1.5; }
.td-empty { color: #d1d5db; }

.ket-badge { font-size: 10.5px; font-weight: 600; padding: 2px 7px; border-radius: 100px; white-space: nowrap; }
.ket-badge--hadir  { background: #f0fdf4; color: #16a34a; border: 1px solid #bbf7d0; }
.ket-badge--izin   { background: #fffbeb; color: #b45309; border: 1px solid #fde68a; }
.ket-badge--sakit  { background: #f0fdf4; color: #0d2818; border: 1px solid #bbf7d0; }
.ket-badge--alpha  { background: #fff1f2; color: #be123c; border: 1px solid #fecdd3; }
.ket-badge--belum  { color: #d1d5db; }
.ket-badge-manual  { font-size: 9px; font-weight: 700; padding: 1px 5px; border-radius: 4px; background: #eff6ff; color: #1e40af; border: 1px solid #bfdbfe; white-space: nowrap; cursor: default; }

.sp-loading { display: flex; flex-direction: column; align-items: center; padding: 40px; gap: 12px; color: #9ca3af; font-size: 13px; }
.sp-error-msg { padding: 20px 22px; color: #dc2626; font-size: 13px; }
.sp-empty { display: flex; flex-direction: column; align-items: center; padding: 40px; gap: 12px; color: #9ca3af; }
.sp-empty p { font-size: 13px; margin: 0; }

/* ── Transition side panel ── */
.side-panel-enter-active { transition: all .3s cubic-bezier(.4,0,.2,1); }
.side-panel-leave-active { transition: all .25s cubic-bezier(.4,0,.2,1); }
.side-panel-enter-from .side-panel { transform: translateX(100%); }
.side-panel-leave-to .side-panel { transform: translateX(100%); }
.side-panel-enter-from { opacity: 0; }
.side-panel-leave-to { opacity: 0; }

/* ── Toast ── */
.toast-msg { position: fixed; bottom: 24px; left: 50%; transform: translateX(-50%); background: #1a2e1a; color: #fff; font-size: 13px; font-weight: 600; padding: 10px 22px; border-radius: 100px; box-shadow: 0 4px 16px rgba(0,0,0,.18); z-index: 9999; white-space: nowrap; }
.toast-enter-active, .toast-leave-active { transition: all .3s; }
.toast-enter-from, .toast-leave-to { opacity: 0; transform: translate(-50%, 10px); }

/* ── Modal (Tolak) ── */
.modal-backdrop { position: fixed; inset: 0; background: rgba(0,0,0,.45); display: flex; align-items: center; justify-content: center; z-index: 1000; padding: 16px; backdrop-filter: blur(2px); }
.modal-box { background: #fff; border-radius: 18px; width: 100%; max-width: 420px; padding: 24px; box-shadow: 0 20px 60px rgba(0,0,0,.18); }
.modal-title { font-size: 16px; font-weight: 700; color: #111827; margin-bottom: 8px; }
.modal-desc  { font-size: 13px; color: #6b7280; margin-bottom: 18px; line-height: 1.6; }
.modal-field { margin-bottom: 14px; }
.modal-label { display: block; font-size: 12px; font-weight: 600; color: #374151; margin-bottom: 5px; }
.modal-textarea { width: 100%; border: 1.5px solid #e5e7eb; border-radius: 10px; padding: 10px 13px; font-size: 13px; font-family: inherit; resize: vertical; outline: none; color: #111827; box-sizing: border-box; }
.modal-textarea:focus { border-color: #48AF4A; }
.modal-input { width: 100%; border: 1.5px solid #e5e7eb; border-radius: 10px; padding: 9px 13px; font-size: 13px; font-family: inherit; outline: none; color: #111827; box-sizing: border-box; }
.modal-input:focus { border-color: #48AF4A; }
.modal-sub { font-size: 12.5px; color: #6b7280; margin-bottom: 16px; }
.jenis-toggle { display: flex; gap: 8px; margin-bottom: 14px; }
.jenis-btn { flex: 1; display: inline-flex; align-items: center; justify-content: center; gap: 6px; padding: 7px 12px; border-radius: 8px; border: 1.5px solid #d1d5db; background: #f9fafb; color: #6b7280; font-size: 13px; font-weight: 500; cursor: pointer; transition: all 0.15s; }
.jenis-btn:hover { border-color: #6ee7b7; background: #f0fdf4; color: #059669; }
.jenis-btn--active { border-color: #16a34a; background: #f0fdf4; color: #16a34a; font-weight: 600; }
.modal-error { font-size: 12px; color: #dc2626; background: #fff1f2; border: 1px solid #fecdd3; border-radius: 8px; padding: 8px 12px; margin-top: 4px; }
.modal-actions { display: flex; gap: 10px; justify-content: flex-end; }
.btn-spinner--white { width: 11px; height: 11px; border: 2px solid rgba(255,255,255,.3); border-top-color: #fff; border-radius: 50%; animation: spin .7s linear infinite; display: inline-block; }
.btn-cancel { background: #f3f4f6; color: #374151; border: none; border-radius: 10px; padding: 10px 20px; font-size: 13px; font-weight: 600; cursor: pointer; font-family: inherit; }
.btn-cancel:hover { background: #e5e7eb; }
.btn-red { background: #dc2626; color: #fff; border: none; border-radius: 10px; padding: 10px 20px; font-size: 13px; font-weight: 600; cursor: pointer; font-family: inherit; }
.btn-red:hover { background: #b91c1c; }
.btn-red:disabled { opacity: .55; cursor: not-allowed; }

/* ── Transition modal ── */
.modal-fade-enter-active, .modal-fade-leave-active { transition: opacity .2s; }
.modal-fade-enter-from, .modal-fade-leave-to { opacity: 0; }

.btn-green-sm { background: #f0fdf4; color: #16a34a; border: 1.5px solid #bbf7d0; border-radius: 8px; padding: 6px 14px; font-size: 12px; font-weight: 600; cursor: pointer; font-family: inherit; text-decoration: none; display: inline-flex; align-items: center; gap: 5px; }
.btn-green-sm:hover { background: #dcfce7; }

.spinner { width: 26px; height: 26px; border: 3px solid #e5e7eb; border-top-color: #48AF4A; border-radius: 50%; animation: spin .7s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }
</style>
