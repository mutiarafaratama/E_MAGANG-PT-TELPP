<template>
  <div class="card">
    <div class="card-header">
      <h3 class="card-title">Daftar Pengajuan Masuk</h3>
      <div class="card-header-actions">
        <div class="search-box">
          <svg class="search-icon" width="14" height="14" viewBox="0 0 24 24" fill="none">
            <circle cx="11" cy="11" r="8" stroke="currentColor" stroke-width="2"/>
            <path d="M21 21l-4.35-4.35" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
          </svg>
          <input v-model="searchQuery" class="search-input" placeholder="Cari nama / institusi…"/>
        </div>
        <button class="btn-green-sm" @click="fetchPengajuan">
          <svg width="13" height="13" viewBox="0 0 24 24" fill="none"><path d="M23 4v6h-6" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/><path d="M1 20v-6h6" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/><path d="M3.51 9a9 9 0 0114.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0020.49 15" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>
          Refresh
        </button>
      </div>
    </div>

    <!-- Date filters -->
    <div class="date-filters">
      <div class="date-filter-group">
        <svg width="13" height="13" viewBox="0 0 24 24" fill="none"><rect x="3" y="4" width="18" height="18" rx="2" stroke="currentColor" stroke-width="2"/><line x1="3" y1="10" x2="21" y2="10" stroke="currentColor" stroke-width="2"/><line x1="8" y1="2" x2="8" y2="6" stroke="currentColor" stroke-width="2" stroke-linecap="round"/><line x1="16" y1="2" x2="16" y2="6" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>
        <select v-model="filterBulan" class="date-select" @change="currentPage = 1">
          <option value="">Semua Bulan</option>
          <option v-for="b in bulanList" :key="b.value" :value="b.value">{{ b.label }}</option>
        </select>
      </div>
      <div class="date-filter-group">
        <select v-model="filterTahun" class="date-select" @change="currentPage = 1">
          <option value="">Semua Tahun</option>
          <option v-for="y in tahunList" :key="y" :value="y">{{ y }}</option>
        </select>
      </div>
      <button v-if="filterBulan || filterTahun" class="btn-clear-date" @click="filterBulan=''; filterTahun=''; currentPage=1">
        <svg width="11" height="11" viewBox="0 0 24 24" fill="none"><path d="M18 6L6 18M6 6l12 12" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"/></svg>
        Reset
      </button>
    </div>

    <!-- Status filter pills -->
    <div class="filter-pills">
      <button v-for="f in statusFilters" :key="f.value" class="pill"
        :class="{ 'pill--active': filterStatus === f.value }" @click="filterStatus = f.value; currentPage = 1">
        {{ f.label }}
        <span v-if="f.count !== undefined" class="pill-count">{{ f.count }}</span>
      </button>
    </div>

    <!-- Table -->
    <div v-if="listLoading" class="empty-state"><div class="spinner"></div></div>
    <div v-else-if="listError" class="empty-state">
      <p class="text-red">{{ listError }}</p>
      <button class="btn-green-sm" @click="fetchPengajuan">Coba lagi</button>
    </div>
    <div v-else-if="filteredList.length === 0" class="empty-state">
      <div class="empty-state__icon"><svg width="36" height="36" viewBox="0 0 24 24" fill="none"><path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z" stroke="#d1d5db" stroke-width="1.5"/><polyline points="14 2 14 8 20 8" stroke="#d1d5db" stroke-width="1.5"/></svg></div>
      <p>Tidak ada pengajuan {{ filterStatus !== 'semua' ? 'dengan status ini' : '' }}.</p>
    </div>
    <div v-else class="table-wrap">
      <table class="data-table">
        <thead>
          <tr>
            <th>Nama Peserta</th><th>Institusi</th><th>Jurusan</th>
            <th>Kategori</th><th>Tanggal Daftar</th><th>Status</th><th>Aksi</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="p in paginatedList" :key="p.id">
            <td class="td-name">
              <div class="name-cell">
                <div class="name-avatar">{{ p.nama_lengkap[0] }}</div>
                <div>
                  <div class="name-text">{{ p.nama_lengkap }}</div>
                  <div class="name-sub">{{ p.email }}</div>
                </div>
              </div>
            </td>
            <td>{{ p.asal_institusi }}</td>
            <td>{{ p.jurusan }}</td>
            <td><span class="tag">{{ formatKategori(p.kategori_magang) }}</span></td>
            <td>{{ formatDate(p.created_at) }}</td>
            <td><span :class="statusBadgeClass(p.status)">{{ formatStatus(p.status) }}</span></td>
            <td>
              <div class="action-btns">
                <button class="btn-detail" @click="openDetail(p.id)">Detail</button>
                <button class="btn-hapus" :title="`Hapus ${p.nama_lengkap}`" @click="konfirmasiHapus(p)">
                  <svg width="13" height="13" viewBox="0 0 24 24" fill="none"><polyline points="3 6 5 6 21 6" stroke="currentColor" stroke-width="2" stroke-linecap="round"/><path d="M19 6l-1 14a2 2 0 01-2 2H8a2 2 0 01-2-2L5 6" stroke="currentColor" stroke-width="2" stroke-linecap="round"/><path d="M10 11v6M14 11v6" stroke="currentColor" stroke-width="2" stroke-linecap="round"/><path d="M9 6V4a1 1 0 011-1h4a1 1 0 011 1v2" stroke="currentColor" stroke-width="2"/></svg>
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Pagination -->
    <div v-if="totalPages > 1" class="pagination">
      <button class="page-btn" :disabled="currentPage === 1" @click="currentPage--">‹</button>
      <span class="page-info">Halaman {{ currentPage }} dari {{ totalPages }}</span>
      <button class="page-btn" :disabled="currentPage === totalPages" @click="currentPage++">›</button>
    </div>
  </div>

  <!-- ── DETAIL DRAWER ──────────────────────────────────────── -->
  <Teleport to="body">
    <div v-if="showDetail" class="drawer-overlay" @click.self="closeDetail">
      <div class="drawer">
        <div class="drawer-header">
          <div>
            <h2 class="drawer-title">Detail Pengajuan</h2>
            <p class="drawer-sub" v-if="detailData">{{ detailData.nama_lengkap }}</p>
          </div>
          <button class="drawer-close" @click="closeDetail">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none"><path d="M18 6L6 18M6 6l12 12" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>
          </button>
        </div>

        <div v-if="detailLoading" class="drawer-body drawer-loading"><div class="spinner"></div></div>

        <div v-else-if="detailData" class="drawer-body">
          <!-- Status banner -->
          <div class="detail-status-banner" :class="`status-banner--${detailData.status}`">
            <span :class="statusBadgeClass(detailData.status)">{{ formatStatus(detailData.status) }}</span>
            <span v-if="detailData.catatan_hrd" class="catatan-hrd">{{ detailData.catatan_hrd }}</span>
          </div>

          <!-- Data Pribadi -->
          <div class="detail-section">
            <div class="detail-section__title">Data Pribadi</div>
            <div class="info-grid">
              <div class="info-item"><span class="info-label">Nama Lengkap</span><span class="info-val">{{ detailData.nama_lengkap }}</span></div>
              <div class="info-item"><span class="info-label">Email</span><span class="info-val">{{ detailData.email }}</span></div>
              <div class="info-item"><span class="info-label">No. HP</span><span class="info-val">{{ detailData.no_hp }}</span></div>
              <div class="info-item"><span class="info-label">Jenis Kelamin</span><span class="info-val">{{ formatJK(detailData.jenis_kelamin) }}</span></div>
              <div class="info-item"><span class="info-label">Tempat Lahir</span><span class="info-val">{{ detailData.tempat_lahir }}</span></div>
              <div class="info-item"><span class="info-label">Tanggal Lahir</span><span class="info-val">{{ formatDate(detailData.tanggal_lahir) }}</span></div>
              <div class="info-item info-item--full"><span class="info-label">Alamat</span><span class="info-val">{{ detailData.alamat }}</span></div>
            </div>
          </div>

          <!-- Data Akademik -->
          <div class="detail-section">
            <div class="detail-section__title">Data Akademik</div>
            <div class="info-grid">
              <div class="info-item"><span class="info-label">Asal Institusi</span><span class="info-val">{{ detailData.asal_institusi }}</span></div>
              <div class="info-item"><span class="info-label">Jurusan</span><span class="info-val">{{ detailData.jurusan }}</span></div>
              <div class="info-item"><span class="info-label">Kelas / Semester</span><span class="info-val">{{ detailData.kelas_semester }}</span></div>
              <div class="info-item"><span class="info-label">Nomor Induk</span><span class="info-val">{{ detailData.nomor_induk }}</span></div>
              <div class="info-item"><span class="info-label">Kategori Magang</span><span class="info-val">{{ formatKategori(detailData.kategori_magang) }}</span></div>
              <div class="info-item"><span class="info-label">Tanggal Daftar</span><span class="info-val">{{ formatDate(detailData.created_at) }}</span></div>
            </div>
          </div>

          <!-- Dokumen Peserta -->
          <div class="detail-section">
            <div class="detail-section__title-row">
              <span class="detail-section__title">Dokumen Peserta</span>
              <span v-if="dokumenPeserta.length > 0" class="doc-count-badge">{{ dokumenPeserta.length }} file</span>
            </div>
            <div v-if="dokumenLoading" class="doc-loading"><div class="spinner spinner--sm"></div></div>
            <div v-else-if="dokumenPeserta.length === 0" class="doc-empty">
              <svg width="28" height="28" viewBox="0 0 24 24" fill="none"><path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z" stroke="#d1d5db" stroke-width="1.5"/><polyline points="14 2 14 8 20 8" stroke="#d1d5db" stroke-width="1.5"/></svg>
              Belum ada dokumen diunggah oleh peserta.
            </div>
            <div v-else class="doc-list">
              <div v-for="doc in dokumenPeserta" :key="doc.id" class="doc-item">
                <div class="doc-icon" :class="`doc-icon--${docTypeColor(doc.jenis)}`">
                  <svg width="15" height="15" viewBox="0 0 24 24" fill="none"><path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z" stroke="currentColor" stroke-width="2"/><polyline points="14 2 14 8 20 8" stroke="currentColor" stroke-width="2"/><line x1="16" y1="13" x2="8" y2="13" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/><line x1="16" y1="17" x2="8" y2="17" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/></svg>
                </div>
                <div class="doc-info">
                  <div class="doc-name">{{ doc.nama_file }}</div>
                  <div class="doc-meta">
                    <span class="doc-jenis-tag" :class="`doc-jenis-tag--${docTypeColor(doc.jenis)}`">{{ formatJenisDok(doc.jenis) }}</span>
                    <span v-if="doc.ukuran_bytes">{{ formatSize(doc.ukuran_bytes) }}</span>
                    <span>{{ formatDate(doc.uploaded_at) }}</span>
                  </div>
                </div>
                <div class="doc-actions">
                  <button class="btn-preview" @click="previewDoc(doc)">
                    <svg width="14" height="14" viewBox="0 0 24 24" fill="none"><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z" stroke="currentColor" stroke-width="2"/><circle cx="12" cy="12" r="3" stroke="currentColor" stroke-width="2"/></svg>
                  </button>
                  <button class="btn-download" :class="{ 'btn-download--loading': downloadingId === doc.id }" @click="downloadDoc(doc)">
                    <div v-if="downloadingId === doc.id" class="spinner spinner--sm"></div>
                    <svg v-else width="14" height="14" viewBox="0 0 24 24" fill="none"><path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4" stroke="currentColor" stroke-width="2" stroke-linecap="round"/><polyline points="7 10 12 15 17 10" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/><line x1="12" y1="15" x2="12" y2="3" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>
                  </button>
                </div>
              </div>
            </div>
          </div>

          <!-- Surat Balasan HRD -->
          <div class="detail-section surat-balasan-section">
            <div class="detail-section__title-row">
              <svg width="13" height="13" viewBox="0 0 24 24" fill="none" style="color:#1a5c20"><path d="M20 4H4a2 2 0 00-2 2v12a2 2 0 002 2h16a2 2 0 002-2V6a2 2 0 00-2-2z" stroke="currentColor" stroke-width="2"/><polyline points="22 6 12 13 2 6" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>
              <span class="detail-section__title" style="color:#0d2818;">Surat Balasan HRD</span>
              <span v-if="suratBalasanList.length > 0" class="doc-count-badge">{{ suratBalasanList.length }} file</span>
            </div>
            <div v-if="!dokumenLoading && suratBalasanList.length > 0" class="doc-list">
              <div v-for="doc in suratBalasanList" :key="doc.id" class="doc-item doc-item--surat">
                <div class="doc-icon doc-icon--green">
                  <svg width="15" height="15" viewBox="0 0 24 24" fill="none"><path d="M20 4H4a2 2 0 00-2 2v12a2 2 0 002 2h16a2 2 0 002-2V6a2 2 0 00-2-2z" stroke="currentColor" stroke-width="2"/><polyline points="22 6 12 13 2 6" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>
                </div>
                <div class="doc-info">
                  <div class="doc-name">{{ doc.nama_file }}</div>
                  <div class="doc-meta">
                    <span class="doc-jenis-tag doc-jenis-tag--green">Surat Balasan</span>
                    <span v-if="doc.ukuran_bytes">{{ formatSize(doc.ukuran_bytes) }}</span>
                    <span>{{ formatDate(doc.uploaded_at) }}</span>
                  </div>
                </div>
                <div class="doc-actions">
                  <button class="btn-preview" @click="previewDoc(doc)">
                    <svg width="14" height="14" viewBox="0 0 24 24" fill="none"><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z" stroke="currentColor" stroke-width="2"/><circle cx="12" cy="12" r="3" stroke="currentColor" stroke-width="2"/></svg>
                  </button>
                  <button class="btn-download" :class="{ 'btn-download--loading': downloadingId === doc.id }" @click="downloadDoc(doc)">
                    <div v-if="downloadingId === doc.id" class="spinner spinner--sm"></div>
                    <svg v-else width="14" height="14" viewBox="0 0 24 24" fill="none"><path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4" stroke="currentColor" stroke-width="2" stroke-linecap="round"/><polyline points="7 10 12 15 17 10" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/><line x1="12" y1="15" x2="12" y2="3" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>
                  </button>
                </div>
              </div>
            </div>
            <div v-else-if="dokumenLoading" class="doc-loading"><div class="spinner spinner--sm"></div></div>

            <div class="upload-area" @dragover.prevent @drop.prevent="onDrop">
              <div v-if="!uploadFile" class="upload-placeholder" @click="triggerFileInput">
                <svg width="28" height="28" viewBox="0 0 24 24" fill="none"><path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4" stroke="#9ca3af" stroke-width="1.5" stroke-linecap="round"/><polyline points="17 8 12 3 7 8" stroke="#9ca3af" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/><line x1="12" y1="3" x2="12" y2="15" stroke="#9ca3af" stroke-width="1.5" stroke-linecap="round"/></svg>
                <p>{{ suratBalasanList.length > 0 ? 'Upload ulang / revisi surat' : 'Upload Surat Balasan Magang' }}</p>
                <span>PDF, DOC, DOCX — maks. 10MB</span>
              </div>
              <div v-else class="upload-selected">
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none"><path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z" stroke="#48AF4A" stroke-width="2"/><polyline points="14 2 14 8 20 8" stroke="#48AF4A" stroke-width="2"/></svg>
                <span class="upload-filename">{{ uploadFile.name }}</span>
                <button class="upload-remove" @click="uploadFile = null">
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none"><path d="M18 6L6 18M6 6l12 12" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>
                </button>
              </div>
            </div>
            <input ref="fileInputRef" type="file" accept=".pdf,.doc,.docx" style="display:none" @change="onFileChange"/>
            <button class="btn-upload" :disabled="!uploadFile || uploadLoading" @click="submitUpload">
              <svg v-if="!uploadLoading" width="14" height="14" viewBox="0 0 24 24" fill="none"><path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4" stroke="currentColor" stroke-width="2" stroke-linecap="round"/><polyline points="17 8 12 3 7 8" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/><line x1="12" y1="3" x2="12" y2="15" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>
              <div v-else class="spinner spinner--sm spinner--white"></div>
              {{ uploadLoading ? 'Mengupload…' : 'Upload Surat Balasan' }}
            </button>
            <div v-if="uploadSuccess" class="upload-msg upload-msg--ok">✓ Surat balasan berhasil diupload.</div>
            <div v-if="uploadError" class="upload-msg upload-msg--err">{{ uploadError }}</div>
          </div>
        </div>

        <!-- Approve/Reject footer -->
        <div v-if="detailData && ['diajukan','menunggu_verifikasi','diproses'].includes(detailData.status)" class="drawer-footer">
          <div v-if="!hasSuratBalasan" class="surat-warning">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none"><circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2"/><line x1="12" y1="8" x2="12" y2="12" stroke="currentColor" stroke-width="2" stroke-linecap="round"/><circle cx="12" cy="16" r="1" fill="currentColor"/></svg>
            Upload <strong>Surat Balasan Magang</strong> terlebih dahulu sebelum menerima atau menolak.
          </div>
          <div v-else class="auto-email-info">
            <svg width="13" height="13" viewBox="0 0 24 24" fill="none"><path d="M20 4H4a2 2 0 00-2 2v12a2 2 0 002 2h16a2 2 0 002-2V6a2 2 0 00-2-2z" stroke="currentColor" stroke-width="2"/><polyline points="22 6 12 13 2 6" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>
            Email otomatis dikirim ke <strong>{{ detailData.email }}</strong> setelah keputusan ini.
          </div>
          <div class="catatan-row">
            <label class="catatan-label">Catatan HRD (opsional)</label>
            <textarea v-model="actionCatatan" class="catatan-input" rows="2" placeholder="Tambahkan catatan untuk peserta…"></textarea>
          </div>
          <div class="footer-btns">
            <button class="btn-reject-lg" :disabled="actionLoading || !hasSuratBalasan" @click="submitAction('ditolak')">
              <div v-if="actionLoading && pendingAction === 'ditolak'" class="spinner spinner--sm spinner--white"></div>
              <svg v-else width="15" height="15" viewBox="0 0 24 24" fill="none"><circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2"/><line x1="15" y1="9" x2="9" y2="15" stroke="currentColor" stroke-width="2" stroke-linecap="round"/><line x1="9" y1="9" x2="15" y2="15" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>
              Tolak Pengajuan
            </button>
            <button class="btn-approve-lg" :disabled="actionLoading || !hasSuratBalasan" @click="submitAction('diterima')">
              <div v-if="actionLoading && pendingAction === 'diterima'" class="spinner spinner--sm spinner--white"></div>
              <svg v-else width="15" height="15" viewBox="0 0 24 24" fill="none"><polyline points="20 6 9 17 4 12" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>
              Terima Pengajuan
            </button>
          </div>
          <div v-if="actionError" class="action-error">{{ actionError }}</div>
        </div>

        <!-- Kirim akun footer -->
        <div v-if="detailData && detailData.status === 'diterima'" class="drawer-footer drawer-footer--kirim">
          <template v-if="detailData.akun_terkirim_at">
            <div class="kirim-akun-sent">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none"><circle cx="12" cy="12" r="10" fill="#16a34a"/><path d="M8 12l3 3 5-6" stroke="#fff" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>
              Email &amp; akun dikirim pada {{ new Date(detailData.akun_terkirim_at).toLocaleDateString('id-ID', { day:'numeric', month:'long', year:'numeric' }) }}
            </div>
          </template>
          <template v-else>
            <div class="kirim-akun-info">
              <div class="kirim-akun-info__title">Kirim Ulang Akun Login</div>
              <div class="kirim-akun-info__sub">Kirim ulang email kredensial ke <strong>{{ detailData.email }}</strong></div>
            </div>
            <button class="btn-kirim-akun" :disabled="kirimAkunLoading" @click="kirimAkun">
              <div v-if="kirimAkunLoading" class="spinner spinner--sm spinner--white"></div>
              <svg v-else width="15" height="15" viewBox="0 0 24 24" fill="none"><path d="M22 2L11 13" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/><path d="M22 2L15 22 11 13 2 9l20-7z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>
              {{ kirimAkunLoading ? 'Mengirim…' : 'Kirim Ulang Email Akun' }}
            </button>
            <div v-if="kirimAkunError" class="action-error">{{ kirimAkunError }}</div>
            <div v-if="kirimAkunDone" class="action-success">✓ Email berhasil dikirim ulang!</div>
          </template>
        </div>
      </div>
    </div>
  </Teleport>

  <!-- ── KONFIRMASI HAPUS ───────────────────────────────────── -->
  <Teleport to="body">
    <div v-if="hapusTarget" class="modal-overlay" @click.self="hapusTarget = null">
      <div class="confirm-modal">
        <div class="confirm-modal__icon">
          <svg width="28" height="28" viewBox="0 0 24 24" fill="none"><circle cx="12" cy="12" r="10" stroke="#dc2626" stroke-width="2"/><line x1="12" y1="8" x2="12" y2="12" stroke="#dc2626" stroke-width="2.5" stroke-linecap="round"/><circle cx="12" cy="16" r="1" fill="#dc2626"/></svg>
        </div>
        <div class="confirm-modal__title">Hapus Pengajuan?</div>
        <div class="confirm-modal__sub">
          Pengajuan dari <strong>{{ hapusTarget.nama_lengkap }}</strong> akan dihapus permanen beserta seluruh dokumen terkait. Tindakan ini tidak dapat dibatalkan.
        </div>
        <div class="confirm-modal__actions">
          <button class="btn-cancel-modal" :disabled="hapusLoading" @click="hapusTarget = null">Batal</button>
          <button class="btn-confirm-hapus" :disabled="hapusLoading" @click="eksekusiHapus">
            <div v-if="hapusLoading" class="spinner spinner--sm spinner--white"></div>
            <svg v-else width="14" height="14" viewBox="0 0 24 24" fill="none"><polyline points="3 6 5 6 21 6" stroke="currentColor" stroke-width="2" stroke-linecap="round"/><path d="M19 6l-1 14a2 2 0 01-2 2H8a2 2 0 01-2-2L5 6" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>
            {{ hapusLoading ? 'Menghapus…' : 'Ya, Hapus' }}
          </button>
        </div>
        <div v-if="hapusError" class="confirm-modal__error">{{ hapusError }}</div>
      </div>
    </div>
  </Teleport>

  <!-- ── PREVIEW LIGHTBOX ──────────────────────────────────── -->
  <Teleport to="body">
    <div v-if="previewState.show" class="preview-overlay" @click.self="closePreview">
      <div class="preview-box">
        <div class="preview-header">
          <span class="preview-filename">{{ previewState.nama }}</span>
          <div class="preview-header-actions">
            <button class="preview-btn-dl" @click="downloadDoc(previewState.docRef!)" title="Download">
              <svg width="15" height="15" viewBox="0 0 24 24" fill="none"><path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4" stroke="currentColor" stroke-width="2" stroke-linecap="round"/><polyline points="7 10 12 15 17 10" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/><line x1="12" y1="15" x2="12" y2="3" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>
            </button>
            <button class="preview-close" @click="closePreview">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none"><path d="M18 6L6 18M6 6l12 12" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>
            </button>
          </div>
        </div>
        <div class="preview-body">
          <div v-if="previewState.loading" class="preview-loading"><div class="spinner"></div></div>
          <div v-else-if="previewState.error" class="preview-error">
            <svg width="32" height="32" viewBox="0 0 24 24" fill="none"><circle cx="12" cy="12" r="10" stroke="#9ca3af" stroke-width="1.5"/><line x1="15" y1="9" x2="9" y2="15" stroke="#9ca3af" stroke-width="1.5" stroke-linecap="round"/><line x1="9" y1="9" x2="15" y2="15" stroke="#9ca3af" stroke-width="1.5" stroke-linecap="round"/></svg>
            <p>Tidak dapat menampilkan file ini. Gunakan tombol download.</p>
          </div>
          <img v-else-if="previewState.type === 'image'" :src="previewState.url" class="preview-image" alt="preview"/>
          <iframe v-else-if="previewState.type === 'pdf'" :src="previewState.url" class="preview-iframe"/>
          <div v-else class="preview-error">
            <svg width="40" height="40" viewBox="0 0 24 24" fill="none"><path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z" stroke="#9ca3af" stroke-width="1.5"/><polyline points="14 2 14 8 20 8" stroke="#9ca3af" stroke-width="1.5"/></svg>
            <p>Preview tidak tersedia untuk tipe file ini.</p>
            <button class="btn-green-sm" @click="downloadDoc(previewState.docRef!)">Download File</button>
          </div>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from "vue";
import api from "@/lib/api";

const emit = defineEmits<{ updated: [] }>();

// ── interfaces ──────────────────────────────────────────────────────
interface Pengajuan {
  id: string; nama_lengkap: string; email: string; no_hp: string;
  tempat_lahir: string; tanggal_lahir: string; jenis_kelamin: string;
  alamat: string; asal_institusi: string; jurusan: string;
  kelas_semester: string; nomor_induk: string; kategori_magang: string;
  status: string; catatan_hrd: string | null; akun_terkirim_at?: string | null;
  created_at: string;
}
interface Dokumen {
  id: string; jenis: string; nama_file: string;
  ukuran_bytes: number | null; mime_type: string | null; uploaded_at: string;
}
interface PreviewState {
  show: boolean; loading: boolean; error: boolean;
  url: string; type: "image" | "pdf" | "other"; nama: string; docRef: Dokumen | null;
}

// ── list state ──────────────────────────────────────────────────────
const allPengajuan  = ref<Pengajuan[]>([]);
const listLoading   = ref(false);
const listError     = ref<string | null>(null);
const filterStatus  = ref("menunggu_verifikasi");
const filterBulan   = ref("");
const filterTahun   = ref("");
const searchQuery   = ref("");
const currentPage   = ref(1);
const PAGE_SIZE     = 10;

const bulanList = [
  { value:"1",label:"Januari"},{value:"2",label:"Februari"},{value:"3",label:"Maret"},
  {value:"4",label:"April"},{value:"5",label:"Mei"},{value:"6",label:"Juni"},
  {value:"7",label:"Juli"},{value:"8",label:"Agustus"},{value:"9",label:"September"},
  {value:"10",label:"Oktober"},{value:"11",label:"November"},{value:"12",label:"Desember"},
];
const tahunList = computed(() => {
  const years = new Set(allPengajuan.value.map(p => new Date(p.created_at).getFullYear()));
  const cur = new Date().getFullYear(); years.add(cur); years.add(cur-1);
  return Array.from(years).sort((a,b) => b-a);
});
const statusFilters = computed(() => [
  { value:"semua",               label:"Semua",                count: allPengajuan.value.length },
  { value:"menunggu_verifikasi", label:"Menunggu Verifikasi",  count: allPengajuan.value.filter(p=>p.status==="diajukan"||p.status==="menunggu_verifikasi").length },
  { value:"diproses",            label:"Diproses",             count: allPengajuan.value.filter(p=>p.status==="diproses").length },
  { value:"diterima",            label:"Diterima",             count: allPengajuan.value.filter(p=>p.status==="diterima").length },
  { value:"ditolak",             label:"Ditolak",              count: allPengajuan.value.filter(p=>p.status==="ditolak").length },
]);
const filteredList = computed(() => {
  let list = allPengajuan.value;
  if (filterStatus.value === "menunggu_verifikasi") {
    list = list.filter(p => p.status === "diajukan" || p.status === "menunggu_verifikasi");
  } else if (filterStatus.value !== "semua") {
    list = list.filter(p => p.status === filterStatus.value);
  }
  if (filterBulan.value) list = list.filter(p => (new Date(p.created_at).getMonth()+1) === Number(filterBulan.value));
  if (filterTahun.value) list = list.filter(p => new Date(p.created_at).getFullYear() === Number(filterTahun.value));
  const q = searchQuery.value.trim().toLowerCase();
  if (q) list = list.filter(p => p.nama_lengkap.toLowerCase().includes(q) || p.asal_institusi.toLowerCase().includes(q));
  return list;
});
const totalPages   = computed(() => Math.max(1, Math.ceil(filteredList.value.length / PAGE_SIZE)));
const paginatedList = computed(() => filteredList.value.slice((currentPage.value-1)*PAGE_SIZE, currentPage.value*PAGE_SIZE));
watch([filterStatus, filterBulan, filterTahun, searchQuery], () => { currentPage.value = 1; });

// ── detail state ────────────────────────────────────────────────────
const showDetail      = ref(false);
const detailData      = ref<Pengajuan | null>(null);
const detailLoading   = ref(false);
const dokumenList     = ref<Dokumen[]>([]);
const dokumenLoading  = ref(false);
const actionCatatan   = ref("");
const actionLoading   = ref(false);
const pendingAction   = ref("");
const actionError     = ref<string | null>(null);
const dokumenPeserta  = computed(() => dokumenList.value.filter(d => d.jenis !== "surat_balasan"));
const suratBalasanList = computed(() => dokumenList.value.filter(d => d.jenis === "surat_balasan"));
const hasSuratBalasan = computed(() => suratBalasanList.value.length > 0);
const kirimAkunLoading = ref(false);
const kirimAkunDone   = ref(false);
const kirimAkunError  = ref<string | null>(null);
const downloadingId   = ref<string | null>(null);
const uploadFile      = ref<File | null>(null);
const uploadLoading   = ref(false);
const uploadSuccess   = ref(false);
const uploadError     = ref<string | null>(null);
const fileInputRef    = ref<HTMLInputElement | null>(null);
const hapusTarget     = ref<Pengajuan | null>(null);
const hapusLoading    = ref(false);
const hapusError      = ref<string | null>(null);
const previewState    = ref<PreviewState>({ show:false, loading:false, error:false, url:"", type:"other", nama:"", docRef:null });

// ── fetch ───────────────────────────────────────────────────────────
async function fetchPengajuan() {
  listLoading.value = true; listError.value = null;
  try {
    const res = await api.get("/api/pengajuan?page=1&limit=200");
    allPengajuan.value = Array.isArray(res.data.data) ? res.data.data : [];
  } catch (e: any) {
    listError.value = e.response?.data?.message ?? "Gagal memuat data pengajuan.";
  } finally { listLoading.value = false; }
}

async function openDetail(id: string) {
  showDetail.value = true; detailLoading.value = true;
  detailData.value = null; dokumenList.value = [];
  actionCatatan.value = ""; actionError.value = null;
  uploadFile.value = null; uploadSuccess.value = false; uploadError.value = null;
  kirimAkunLoading.value = false; kirimAkunDone.value = false; kirimAkunError.value = null;
  try {
    const res = await api.get(`/api/pengajuan/${id}`);
    detailData.value = res.data.data ?? res.data;
    fetchDokumen(id);
  } catch { closeDetail(); } finally { detailLoading.value = false; }
}

async function fetchDokumen(pengajuanId: string) {
  dokumenLoading.value = true;
  try {
    const res = await api.get(`/api/dokumen/pengajuan/${pengajuanId}`);
    dokumenList.value = res.data.data ?? [];
  } catch { dokumenList.value = []; } finally { dokumenLoading.value = false; }
}

function closeDetail() { showDetail.value = false; detailData.value = null; uploadFile.value = null; }

async function submitAction(type: string) {
  if (!detailData.value) return;
  actionLoading.value = true; pendingAction.value = type; actionError.value = null;
  try {
    await api.patch(`/api/pengajuan/${detailData.value.id}/status`, { status: type, catatan: actionCatatan.value || undefined });
    await fetchPengajuan();
    const res = await api.get(`/api/pengajuan/${detailData.value.id}`);
    detailData.value = res.data.data ?? res.data;
    actionCatatan.value = "";
    emit("updated");
  } catch (e: any) { actionError.value = e.response?.data?.message ?? "Gagal memperbarui status."; }
  finally { actionLoading.value = false; pendingAction.value = ""; }
}

async function kirimAkun() {
  if (!detailData.value) return;
  kirimAkunLoading.value = true; kirimAkunError.value = null; kirimAkunDone.value = false;
  try {
    await api.post(`/api/pengajuan/${detailData.value.id}/kirim-akun`);
    kirimAkunDone.value = true;
    const res = await api.get(`/api/pengajuan/${detailData.value.id}`);
    detailData.value = res.data.data ?? res.data;
    await fetchPengajuan();
  } catch (e: any) { kirimAkunError.value = e.response?.data?.message ?? "Gagal mengirim akun. Coba lagi."; }
  finally { kirimAkunLoading.value = false; }
}

function konfirmasiHapus(p: Pengajuan) { hapusTarget.value = p; hapusError.value = null; }

async function eksekusiHapus() {
  if (!hapusTarget.value) return;
  hapusLoading.value = true; hapusError.value = null;
  try {
    await api.delete(`/api/pengajuan/${hapusTarget.value.id}`);
    hapusTarget.value = null; await fetchPengajuan(); emit("updated");
  } catch (e: any) { hapusError.value = e.response?.data?.message ?? "Gagal menghapus pengajuan."; }
  finally { hapusLoading.value = false; }
}

async function previewDoc(doc: Dokumen) {
  previewState.value = { show:true, loading:true, error:false, url:"", type:"other", nama:doc.nama_file, docRef:doc };
  try {
    const res = await api.get(`/api/dokumen/${doc.id}/download`, { responseType:"blob" });
    const mime = doc.mime_type ?? "application/octet-stream";
    const blob = new Blob([res.data], { type: mime });
    const url = URL.createObjectURL(blob);
    let type: "image"|"pdf"|"other" = "other";
    if (mime.startsWith("image/")) type = "image";
    else if (mime === "application/pdf") type = "pdf";
    previewState.value = { ...previewState.value, loading:false, url, type };
  } catch { previewState.value = { ...previewState.value, loading:false, error:true }; }
}

function closePreview() {
  if (previewState.value.url) URL.revokeObjectURL(previewState.value.url);
  previewState.value = { show:false, loading:false, error:false, url:"", type:"other", nama:"", docRef:null };
}

function triggerFileInput() { fileInputRef.value?.click(); }
function onFileChange(e: Event) {
  const f = (e.target as HTMLInputElement).files?.[0];
  if (f) { uploadFile.value = f; uploadSuccess.value = false; uploadError.value = null; }
}
function onDrop(e: DragEvent) {
  const f = e.dataTransfer?.files[0];
  if (f) { uploadFile.value = f; uploadSuccess.value = false; uploadError.value = null; }
}
async function submitUpload() {
  if (!uploadFile.value || !detailData.value) return;
  uploadLoading.value = true; uploadSuccess.value = false; uploadError.value = null;
  try {
    const form = new FormData();
    form.append("file", uploadFile.value); form.append("jenis", "surat_balasan");
    form.append("pengajuan_id", detailData.value.id);
    await api.post("/api/dokumen/upload", form, { headers: { "Content-Type": "multipart/form-data" } });
    uploadFile.value = null; uploadSuccess.value = true;
    fetchDokumen(detailData.value.id);
  } catch (e: any) { uploadError.value = e.response?.data?.message ?? "Gagal mengupload dokumen."; }
  finally { uploadLoading.value = false; }
}

async function downloadDoc(doc: Dokumen) {
  if (downloadingId.value) return; downloadingId.value = doc.id;
  try {
    const res = await api.get(`/api/dokumen/${doc.id}/download`, { responseType:"blob" });
    const url = URL.createObjectURL(new Blob([res.data], { type: doc.mime_type ?? "application/octet-stream" }));
    const a = Object.assign(document.createElement("a"), { href: url, download: doc.nama_file });
    document.body.appendChild(a); a.click(); document.body.removeChild(a);
    URL.revokeObjectURL(url);
  } catch {} finally { downloadingId.value = null; }
}

// ── format helpers ──────────────────────────────────────────────────
function formatDate(d: string | null) {
  if (!d) return "-";
  return new Date(d).toLocaleDateString("id-ID", { day:"2-digit", month:"short", year:"numeric" });
}
function formatStatus(s: string) {
  return ({ diajukan:"Diajukan", menunggu_verifikasi:"Menunggu Verifikasi", diproses:"Diproses", diterima:"Diterima", ditolak:"Ditolak" } as Record<string,string>)[s] ?? s;
}
function statusBadgeClass(s: string) {
  const base = "status-badge ";
  if (s === "diterima") return base + "status-badge--green";
  if (s === "ditolak")  return base + "status-badge--red";
  if (s === "diproses") return base + "status-badge--blue";
  return base + "status-badge--yellow";
}
function formatKategori(k: string) {
  return ({
    siswa_smk:"Siswa SMK", mahasiswa_d3:"Mahasiswa D3", mahasiswa_s1:"Mahasiswa S1",
    smk:"SMK / Sederajat", d3_s1_s2:"D3 / S1 / S2",
    d3:"D3", s1:"S1", s2:"S2",
    penelitian:"Penelitian", lainnya:"Lainnya",
  } as Record<string,string>)[k] ?? k;
}
function formatJK(jk: string) { return jk === "l" ? "Laki-laki" : jk === "p" ? "Perempuan" : jk; }
function formatJenisDok(j: string) {
  return ({ proposal_magang:"Proposal", surat_pengantar:"Surat Pengantar", ktp:"KTP", ktm:"KTM", pasfoto:"Pas Foto", bpjs_kis:"BPJS/KIS", surat_balasan:"Surat Balasan", laporan_magang:"Laporan", sertifikat:"Sertifikat" } as Record<string,string>)[j] ?? j;
}
function docTypeColor(jenis: string) {
  return ({ proposal_magang:"red", surat_pengantar:"teal", ktp:"blue", ktm:"blue", pasfoto:"purple", bpjs_kis:"green", surat_balasan:"orange", laporan_magang:"indigo", sertifikat:"gold" } as Record<string,string>)[jenis] ?? "gray";
}
function formatSize(b: number | null) {
  if (!b) return "";
  if (b < 1024) return b + " B";
  if (b < 1048576) return (b/1024).toFixed(1) + " KB";
  return (b/1048576).toFixed(1) + " MB";
}

onMounted(fetchPengajuan);

// expose for parent if needed
defineExpose({ fetchPengajuan });
</script>

<style scoped>
/* ── shared card ──── */
.card { background:#fff; border-radius:14px; border:1px solid #e9f5e9; box-shadow:0 1px 3px rgba(13,40,24,0.05); overflow:hidden; }
.card-header { display:flex; align-items:center; justify-content:space-between; padding:16px 20px; border-bottom:1px solid #f0faf0; gap:12px; flex-wrap:wrap; }
.card-title { font-size:13.5px; font-weight:700; color:#111827; margin:0; }
.card-header-actions { display:flex; align-items:center; gap:8px; }
/* ── search ── */
.search-box { position:relative; display:flex; align-items:center; }
.search-icon { position:absolute; left:10px; color:#9ca3af; pointer-events:none; }
.search-input { border:1px solid #e5e7eb; border-radius:8px; padding:7px 12px 7px 30px; font-size:12.5px; font-family:inherit; color:#374151; width:220px; outline:none; }
.search-input:focus { border-color:#48AF4A; }
/* ── date filters ── */
.date-filters { display:flex; align-items:center; gap:8px; padding:10px 20px; border-bottom:1px solid #f3f4f6; flex-wrap:wrap; }
.date-filter-group { display:flex; align-items:center; gap:6px; background:#f9fafb; border:1px solid #e5e7eb; border-radius:8px; padding:5px 10px; color:#6b7280; }
.date-select { background:none; border:none; font-size:12.5px; font-family:inherit; color:#374151; cursor:pointer; outline:none; }
.btn-clear-date { background:none; border:1px solid #e5e7eb; border-radius:7px; padding:5px 10px; font-size:11.5px; color:#9ca3af; cursor:pointer; display:flex; align-items:center; gap:4px; font-family:inherit; }
.btn-clear-date:hover { border-color:#dc2626; color:#dc2626; background:#fef2f2; }
/* ── filter pills ── */
.filter-pills { display:flex; gap:6px; padding:12px 20px; border-bottom:1px solid #f3f4f6; flex-wrap:wrap; }
.pill { background:#f9fafb; border:1px solid #e5e7eb; border-radius:100px; padding:5px 13px; font-size:12px; font-weight:500; cursor:pointer; color:#6b7280; font-family:inherit; display:flex; align-items:center; gap:5px; transition:all 0.15s; }
.pill:hover { border-color:#48AF4A; color:#16a34a; }
.pill--active { background:#f0fdf4; border-color:#48AF4A; color:#16a34a; font-weight:600; }
.pill-count { background:#e5e7eb; color:#6b7280; border-radius:100px; padding:0 6px; font-size:10.5px; font-weight:700; }
.pill--active .pill-count { background:#bbf7d0; color:#16a34a; }
/* ── table ── */
.table-wrap { overflow-x:auto; }
.data-table { width:100%; border-collapse:collapse; font-size:13px; }
.data-table th { padding:11px 16px; text-align:left; font-size:10.5px; font-weight:600; color:#6b7280; background:#f9fafb; border-bottom:1px solid #f1f5f9; text-transform:uppercase; letter-spacing:0.04em; white-space:nowrap; }
.data-table td { padding:13px 16px; border-bottom:1px solid #f9fafb; color:#374151; vertical-align:middle; }
.td-name { min-width:160px; }
.name-cell { display:flex; align-items:center; gap:10px; }
.name-avatar { width:32px; height:32px; border-radius:8px; background:linear-gradient(135deg,#48AF4A,#1a5c20); color:#fff; font-size:13px; font-weight:700; display:flex; align-items:center; justify-content:center; flex-shrink:0; }
.name-text { font-weight:600; color:#111827; font-size:12.5px; }
.name-sub { font-size:11px; color:#9ca3af; }
.tag { background:#f0fdf4; color:#0d2818; border-radius:6px; padding:2px 8px; font-size:11px; font-weight:600; white-space:nowrap; }
/* ── status badges ── */
.status-badge { padding:3px 10px; border-radius:100px; font-size:11px; font-weight:600; white-space:nowrap; }
.status-badge--yellow { background:#fefce8; color:#ca8a04; }
.status-badge--green  { background:#f0fdf4; color:#16a34a; }
.status-badge--red    { background:#fef2f2; color:#dc2626; }
.status-badge--blue   { background:#f0fdf4; color:#1a5c20; }
/* ── action btns ── */
.action-btns { display:flex; gap:5px; align-items:center; }
.btn-detail { background:#f3f4f6; color:#374151; border:none; border-radius:7px; padding:5px 11px; font-size:11.5px; font-weight:600; font-family:inherit; cursor:pointer; white-space:nowrap; }
.btn-detail:hover { background:#e5e7eb; }
.btn-hapus { width:28px; height:28px; background:#fef2f2; border:1px solid #fecaca; border-radius:7px; display:flex; align-items:center; justify-content:center; color:#dc2626; cursor:pointer; flex-shrink:0; transition:all 0.15s; padding:0; }
.btn-hapus:hover { background:#fee2e2; border-color:#f87171; }
/* ── pagination ── */
.pagination { display:flex; align-items:center; justify-content:center; gap:12px; padding:14px 20px; border-top:1px solid #f3f4f6; }
.page-btn { background:#f3f4f6; border:none; border-radius:7px; width:30px; height:30px; font-size:16px; cursor:pointer; display:flex; align-items:center; justify-content:center; color:#374151; }
.page-btn:disabled { opacity:0.4; cursor:default; }
.page-btn:not(:disabled):hover { background:#e5e7eb; }
.page-info { font-size:12.5px; color:#6b7280; }
/* ── empty/loading ── */
.empty-state { display:flex; flex-direction:column; align-items:center; padding:44px 24px; gap:12px; text-align:center; }
.empty-state__icon { width:72px; height:72px; background:#f9fafb; border-radius:50%; display:flex; align-items:center; justify-content:center; }
.empty-state p { font-size:13px; color:#9ca3af; line-height:1.7; margin:0; }
.text-red { color:#dc2626; font-size:13px; }
.spinner { width:24px; height:24px; border:2.5px solid #e5e7eb; border-top-color:#48AF4A; border-radius:50%; animation:spin 0.7s linear infinite; }
.spinner--sm { width:14px; height:14px; border-width:2px; }
.spinner--white { border-color:rgba(255,255,255,0.3); border-top-color:#fff; }
@keyframes spin { to { transform:rotate(360deg); } }
/* ── btn-green-sm ── */
.btn-green-sm { background:#48AF4A; color:#fff; border:none; border-radius:8px; padding:6px 14px; font-size:12px; font-weight:600; font-family:inherit; cursor:pointer; white-space:nowrap; display:flex; align-items:center; gap:5px; }
.btn-green-sm:hover { background:#48AF4A; }
/* ── drawer ── */
.drawer-overlay { position:fixed; inset:0; background:rgba(0,0,0,0.45); backdrop-filter:blur(2px); z-index:100; display:flex; justify-content:flex-end; }
.drawer { width:min(560px,100vw); height:100vh; background:#fff; display:flex; flex-direction:column; box-shadow:-4px 0 24px rgba(0,0,0,0.12); overflow:hidden; }
.drawer-header { display:flex; align-items:flex-start; justify-content:space-between; padding:20px 24px 16px; border-bottom:1px solid #f0faf0; flex-shrink:0; }
.drawer-title { font-size:16px; font-weight:700; color:#111827; margin:0 0 2px; }
.drawer-sub { font-size:12.5px; color:#6b7280; margin:0; }
.drawer-close { background:#f3f4f6; border:none; border-radius:8px; width:32px; height:32px; display:flex; align-items:center; justify-content:center; cursor:pointer; color:#6b7280; flex-shrink:0; margin-top:2px; }
.drawer-close:hover { background:#e5e7eb; color:#111827; }
.drawer-body { flex:1; overflow-y:auto; padding:20px 24px; display:flex; flex-direction:column; gap:20px; }
.drawer-loading { align-items:center; justify-content:center; }
/* ── detail ── */
.detail-status-banner { display:flex; align-items:center; gap:10px; padding:10px 14px; background:#f9fafb; border-radius:10px; }
.catatan-hrd { font-size:12px; color:#6b7280; font-style:italic; }
.status-banner--diterima { background:#f0fdf4; }
.status-banner--ditolak  { background:#fef2f2; }
.detail-section { display:flex; flex-direction:column; gap:10px; }
.detail-section__title { font-size:11px; font-weight:700; color:#9ca3af; text-transform:uppercase; letter-spacing:0.08em; }
.detail-section__title-row { display:flex; align-items:center; gap:8px; }
.info-grid { display:grid; grid-template-columns:1fr 1fr; gap:10px; }
.info-item { display:flex; flex-direction:column; gap:2px; }
.info-item--full { grid-column:1/-1; }
.info-label { font-size:11px; color:#9ca3af; font-weight:500; }
.info-val { font-size:13px; color:#111827; font-weight:500; }
/* ── dokumen ── */
.doc-count-badge { background:#1a5c20; color:#fff; font-size:10px; font-weight:700; padding:2px 8px; border-radius:100px; }
.doc-loading { display:flex; justify-content:center; padding:16px; }
.doc-empty { display:flex; flex-direction:column; align-items:center; gap:8px; font-size:12.5px; color:#9ca3af; text-align:center; padding:20px 12px; background:#f9fafb; border-radius:10px; border:1.5px dashed #e5e7eb; }
.doc-list { display:flex; flex-direction:column; gap:7px; }
.doc-item { display:flex; align-items:center; gap:10px; padding:10px 12px; background:#fff; border-radius:10px; border:1px solid #e9f5e9; transition:border-color 0.15s,box-shadow 0.15s; }
.doc-item:hover { border-color:#86efac; box-shadow:0 2px 8px rgba(26,92,32,0.07); }
.doc-item--surat { border-color:#bbf7d0; background:#f8fffe; }
.doc-icon { width:34px; height:34px; border-radius:8px; display:flex; align-items:center; justify-content:center; flex-shrink:0; }
.doc-icon--red{background:#fef2f2;color:#dc2626} .doc-icon--blue{background:#f0fdf4;color:#1a5c20}
.doc-icon--purple{background:#f5f3ff;color:#7c3aed} .doc-icon--green{background:#f0fdf4;color:#16a34a}
.doc-icon--orange{background:#fff7ed;color:#ea580c} .doc-icon--indigo{background:#eef2ff;color:#4338ca}
.doc-icon--gold{background:#fefce8;color:#ca8a04} .doc-icon--gray{background:#f3f4f6;color:#6b7280}
.doc-info { flex:1; min-width:0; }
.doc-name { font-size:12.5px; font-weight:600; color:#111827; white-space:nowrap; overflow:hidden; text-overflow:ellipsis; }
.doc-meta { display:flex; align-items:center; flex-wrap:wrap; gap:5px; margin-top:3px; font-size:11px; color:#9ca3af; }
.doc-jenis-tag { padding:1px 7px; border-radius:100px; font-size:10.5px; font-weight:600; }
.doc-jenis-tag--red{background:#fef2f2;color:#dc2626} .doc-jenis-tag--blue{background:#f0fdf4;color:#1a5c20}
.doc-jenis-tag--purple{background:#f5f3ff;color:#7c3aed} .doc-jenis-tag--green{background:#f0fdf4;color:#16a34a}
.doc-jenis-tag--orange{background:#fff7ed;color:#ea580c} .doc-jenis-tag--indigo{background:#eef2ff;color:#4338ca}
.doc-jenis-tag--gold{background:#fefce8;color:#ca8a04} .doc-jenis-tag--gray{background:#f3f4f6;color:#6b7280}
.doc-actions { display:flex; gap:5px; flex-shrink:0; }
.btn-preview { width:32px; height:32px; background:#f0fdf4; border:1px solid #bbf7d0; border-radius:8px; display:flex; align-items:center; justify-content:center; color:#1a5c20; cursor:pointer; flex-shrink:0; transition:all 0.15s; padding:0; }
.btn-preview:hover { background:#dcfce7; border-color:#93c5fd; }
.btn-download { width:32px; height:32px; background:#f0fdf4; border:1px solid #bbf7d0; border-radius:8px; display:flex; align-items:center; justify-content:center; color:#16a34a; cursor:pointer; flex-shrink:0; transition:all 0.15s; }
.btn-download:hover { background:#dcfce7; border-color:#86efac; }
.btn-download--loading { opacity:0.6; cursor:wait; }
/* ── upload ── */
.surat-balasan-section { background:#f0fdf4; border:1.5px solid #bbf7d0; border-radius:12px; padding:14px; }
.surat-balasan-section .detail-section__title { color:#0d2818; }
.surat-balasan-section .detail-section__title-row { margin-bottom:10px; }
.surat-balasan-section .upload-area { border-color:#86efac; }
.surat-balasan-section .upload-area:hover { border-color:#48AF4A; }
.upload-area { border:2px dashed #d1d5db; border-radius:10px; cursor:pointer; transition:border-color 0.15s; }
.upload-area:hover { border-color:#48AF4A; }
.upload-placeholder { display:flex; flex-direction:column; align-items:center; gap:6px; padding:24px 16px; }
.upload-placeholder p { font-size:13px; color:#374151; font-weight:500; margin:0; }
.upload-placeholder span { font-size:11px; color:#9ca3af; }
.upload-selected { display:flex; align-items:center; gap:10px; padding:12px 16px; }
.upload-filename { flex:1; font-size:12.5px; font-weight:500; color:#111827; min-width:0; overflow:hidden; text-overflow:ellipsis; white-space:nowrap; }
.upload-remove { background:none; border:none; color:#9ca3af; cursor:pointer; display:flex; align-items:center; padding:2px; }
.upload-remove:hover { color:#dc2626; }
.btn-upload { width:100%; margin-top:10px; background:#48AF4A; color:#fff; border:none; border-radius:9px; padding:10px 16px; font-size:13px; font-weight:600; font-family:inherit; cursor:pointer; display:flex; align-items:center; justify-content:center; gap:7px; transition:background 0.15s; }
.btn-upload:hover:not(:disabled) { background:#48AF4A; }
.btn-upload:disabled { opacity:0.5; cursor:default; }
.upload-msg { font-size:12px; margin-top:6px; padding:7px 12px; border-radius:7px; }
.upload-msg--ok  { background:#f0fdf4; color:#16a34a; }
.upload-msg--err { background:#fef2f2; color:#dc2626; }
/* ── drawer footer ── */
.drawer-footer { flex-shrink:0; padding:16px 24px; border-top:1px solid #f0faf0; background:#fff; display:flex; flex-direction:column; gap:12px; }
.drawer-footer--kirim { background:#f0fdf4; border-top:1px solid #bbf7d0; }
.surat-warning { display:flex; align-items:flex-start; gap:8px; background:#fffbeb; border:1px solid #fde68a; border-radius:8px; padding:10px 14px; font-size:12.5px; color:#92400e; line-height:1.5; }
.surat-warning svg { flex-shrink:0; margin-top:1px; color:#d97706; }
.auto-email-info { display:flex; align-items:center; gap:8px; background:#f0fdf4; border:1px solid #bbf7d0; border-radius:8px; padding:9px 14px; font-size:12px; color:#166534; line-height:1.5; }
.auto-email-info svg { flex-shrink:0; color:#16a34a; }
.catatan-label { font-size:11.5px; font-weight:600; color:#374151; display:block; margin-bottom:5px; }
.catatan-input { width:100%; border:1px solid #e5e7eb; border-radius:8px; padding:8px 12px; font-size:12.5px; font-family:inherit; color:#374151; resize:none; outline:none; box-sizing:border-box; }
.catatan-input:focus { border-color:#48AF4A; }
.footer-btns { display:flex; gap:10px; }
.btn-reject-lg { flex:1; background:#fef2f2; color:#dc2626; border:1px solid #fecaca; border-radius:9px; padding:11px 16px; font-size:13px; font-weight:600; font-family:inherit; cursor:pointer; display:flex; align-items:center; justify-content:center; gap:7px; }
.btn-reject-lg:hover:not(:disabled) { background:#fee2e2; }
.btn-approve-lg { flex:1; background:#48AF4A; color:#fff; border:none; border-radius:9px; padding:11px 16px; font-size:13px; font-weight:600; font-family:inherit; cursor:pointer; display:flex; align-items:center; justify-content:center; gap:7px; }
.btn-approve-lg:hover:not(:disabled) { background:#48AF4A; }
.btn-approve-lg:disabled,.btn-reject-lg:disabled { opacity:0.45; cursor:not-allowed; }
.action-error { font-size:12px; color:#dc2626; background:#fef2f2; padding:7px 12px; border-radius:7px; }
.action-success { background:#dcfce7; border:1px solid #86efac; border-radius:8px; padding:8px 12px; font-size:12px; color:#16a34a; font-weight:600; }
.kirim-akun-sent { display:flex; align-items:center; gap:8px; font-size:12.5px; color:#16a34a; font-weight:600; }
.kirim-akun-info__title { font-size:13px; font-weight:700; color:#0d2818; margin-bottom:3px; }
.kirim-akun-info__sub { font-size:11.5px; color:#6b7280; line-height:1.5; }
.btn-kirim-akun { background:#0d2818; color:#fff; border:none; border-radius:9px; padding:10px 18px; font-size:13px; font-weight:600; font-family:inherit; cursor:pointer; display:flex; align-items:center; gap:8px; transition:background 0.15s; width:100%; justify-content:center; }
.btn-kirim-akun:hover:not(:disabled) { background:#1a5c20; }
.btn-kirim-akun:disabled { opacity:0.6; cursor:not-allowed; }
/* ── modal hapus ── */
.modal-overlay { position:fixed; inset:0; background:rgba(0,0,0,0.5); backdrop-filter:blur(3px); z-index:200; display:flex; align-items:center; justify-content:center; padding:20px; }
.confirm-modal { background:#fff; border-radius:16px; padding:28px 28px 24px; max-width:400px; width:100%; box-shadow:0 20px 60px rgba(0,0,0,0.18); display:flex; flex-direction:column; align-items:center; gap:12px; text-align:center; }
.confirm-modal__icon { width:56px; height:56px; border-radius:50%; background:#fef2f2; display:flex; align-items:center; justify-content:center; }
.confirm-modal__title { font-size:16px; font-weight:700; color:#111827; }
.confirm-modal__sub { font-size:13px; color:#6b7280; line-height:1.6; }
.confirm-modal__sub strong { color:#111827; }
.confirm-modal__actions { display:flex; gap:10px; width:100%; margin-top:4px; }
.btn-cancel-modal { flex:1; background:#f3f4f6; color:#374151; border:none; border-radius:9px; padding:10px 16px; font-size:13px; font-weight:600; font-family:inherit; cursor:pointer; }
.btn-cancel-modal:hover:not(:disabled) { background:#e5e7eb; }
.btn-confirm-hapus { flex:1; background:#dc2626; color:#fff; border:none; border-radius:9px; padding:10px 16px; font-size:13px; font-weight:600; font-family:inherit; cursor:pointer; display:flex; align-items:center; justify-content:center; gap:7px; }
.btn-confirm-hapus:hover:not(:disabled) { background:#b91c1c; }
.btn-confirm-hapus:disabled,.btn-cancel-modal:disabled { opacity:0.5; cursor:default; }
.confirm-modal__error { font-size:12px; color:#dc2626; background:#fef2f2; padding:7px 12px; border-radius:7px; width:100%; text-align:left; }
/* ── preview lightbox ── */
.preview-overlay { position:fixed; inset:0; background:rgba(0,0,0,0.75); backdrop-filter:blur(4px); z-index:300; display:flex; align-items:center; justify-content:center; padding:20px; }
.preview-box { background:#fff; border-radius:14px; width:min(860px,100%); max-height:90vh; display:flex; flex-direction:column; overflow:hidden; box-shadow:0 24px 80px rgba(0,0,0,0.3); }
.preview-header { display:flex; align-items:center; justify-content:space-between; padding:14px 18px; border-bottom:1px solid #f3f4f6; gap:12px; flex-shrink:0; }
.preview-filename { font-size:13px; font-weight:600; color:#111827; white-space:nowrap; overflow:hidden; text-overflow:ellipsis; min-width:0; }
.preview-header-actions { display:flex; gap:6px; flex-shrink:0; }
.preview-btn-dl { background:#f0fdf4; border:1px solid #bbf7d0; color:#16a34a; border-radius:8px; width:32px; height:32px; display:flex; align-items:center; justify-content:center; cursor:pointer; transition:all 0.15s; }
.preview-btn-dl:hover { background:#dcfce7; }
.preview-close { background:#f3f4f6; border:none; border-radius:8px; width:32px; height:32px; display:flex; align-items:center; justify-content:center; cursor:pointer; color:#6b7280; }
.preview-close:hover { background:#e5e7eb; color:#111827; }
.preview-body { flex:1; overflow:auto; min-height:300px; display:flex; align-items:center; justify-content:center; background:#f9fafb; }
.preview-image { max-width:100%; max-height:calc(90vh - 70px); object-fit:contain; display:block; }
.preview-iframe { width:100%; height:calc(90vh - 70px); border:none; }
.preview-loading { display:flex; align-items:center; justify-content:center; padding:60px; }
.preview-error { display:flex; flex-direction:column; align-items:center; gap:12px; padding:48px 24px; text-align:center; color:#9ca3af; font-size:13px; }
@media (max-width:600px) { .info-grid { grid-template-columns:1fr; } .drawer { width:100vw; } }
</style>
