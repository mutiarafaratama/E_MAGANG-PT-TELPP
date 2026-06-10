<template>
  <DashboardLayout :nav-groups="navGroups" role-name="Staff HRD" default-tab="beranda" ref="layout">
    <template #default>

      <!-- ══════════════════════════════════════════════════════
           BERANDA
      ══════════════════════════════════════════════════════ -->
      <template v-if="activeTab === 'beranda'">
        <div class="welcome-banner">
          <div>
            <div class="welcome-banner__greeting">Halo, {{ firstName }}! 👋</div>
            <div class="welcome-banner__sub">Kelola seleksi dan pelaksanaan magang dari sini.</div>
          </div>
          <div class="welcome-banner__icon">
            <svg width="36" height="36" viewBox="0 0 24 24" fill="none"><path d="M17 21v-2a4 4 0 00-4-4H5a4 4 0 00-4 4v2" stroke="rgba(255,255,255,0.4)" stroke-width="1.5"/><circle cx="9" cy="7" r="4" stroke="rgba(255,255,255,0.4)" stroke-width="1.5"/><path d="M23 21v-2a4 4 0 00-3-3.87" stroke="#86efac" stroke-width="1.5" stroke-linecap="round"/><path d="M16 3.13a4 4 0 010 7.75" stroke="#86efac" stroke-width="1.5" stroke-linecap="round"/></svg>
          </div>
        </div>

        <div class="stats-row">
          <div class="stat-card">
            <div class="stat-card__icon" style="background:#fff7ed;color:#ea580c">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none"><path d="M17 21v-2a4 4 0 00-4-4H5a4 4 0 00-4 4v2" stroke="currentColor" stroke-width="2"/><circle cx="9" cy="7" r="4" stroke="currentColor" stroke-width="2"/></svg>
            </div>
            <div>
              <div class="stat-card__label">Total Pengajuan</div>
              <div class="stat-card__value">{{ statsLoading ? '…' : totalAll }}</div>
            </div>
          </div>
          <div class="stat-card">
            <div class="stat-card__icon" style="background:#fefce8;color:#ca8a04">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none"><circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2"/><polyline points="12 6 12 12 16 14" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>
            </div>
            <div>
              <div class="stat-card__label">Menunggu Verifikasi</div>
              <div class="stat-card__value">{{ statsLoading ? '…' : totalPending }}</div>
            </div>
          </div>
          <div class="stat-card">
            <div class="stat-card__icon" style="background:#f0fdf4;color:#16a34a">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none"><polyline points="20 6 9 17 4 12" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>
            </div>
            <div>
              <div class="stat-card__label">Diterima</div>
              <div class="stat-card__value">{{ statsLoading ? '…' : totalDiterima }}</div>
            </div>
          </div>
          <div class="stat-card">
            <div class="stat-card__icon" style="background:#fef2f2;color:#dc2626">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none"><circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2"/><line x1="15" y1="9" x2="9" y2="15" stroke="currentColor" stroke-width="2" stroke-linecap="round"/><line x1="9" y1="9" x2="15" y2="15" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>
            </div>
            <div>
              <div class="stat-card__label">Ditolak</div>
              <div class="stat-card__value">{{ statsLoading ? '…' : totalDitolak }}</div>
            </div>
          </div>
        </div>

        <!-- Perlu ditindaklanjuti -->
        <div class="card">
          <div class="card-header">
            <h3 class="card-title">Pengajuan Perlu Ditindaklanjuti</h3>
            <span class="badge badge--yellow">{{ totalPending }} Pending</span>
          </div>
          <div v-if="statsLoading" class="empty-state">
            <div class="spinner"></div>
          </div>
          <div v-else-if="pendingItems.length === 0" class="empty-state">
            <div class="empty-state__icon"><svg width="36" height="36" viewBox="0 0 24 24" fill="none"><path d="M17 21v-2a4 4 0 00-4-4H5a4 4 0 00-4 4v2" stroke="#d1d5db" stroke-width="1.5"/><circle cx="9" cy="7" r="4" stroke="#d1d5db" stroke-width="1.5"/></svg></div>
            <p>Tidak ada pengajuan yang perlu diproses saat ini.</p>
          </div>
          <div v-else class="table-wrap">
            <table class="data-table">
              <thead><tr><th>Nama</th><th>Institusi</th><th>Kategori</th><th>Tanggal</th><th>Status</th><th>Aksi</th></tr></thead>
              <tbody>
                <tr v-for="p in pendingItems.slice(0,5)" :key="p.id">
                  <td class="td-name">{{ p.nama_lengkap }}</td>
                  <td>{{ p.asal_institusi }}</td>
                  <td>{{ formatKategori(p.kategori_magang) }}</td>
                  <td>{{ formatDate(p.created_at) }}</td>
                  <td><span :class="statusBadgeClass(p.status)">{{ formatStatus(p.status) }}</span></td>
                  <td>
                    <button class="btn-detail" @click="openDetail(p.id)">Tinjau</button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
          <div v-if="pendingItems.length > 5" class="card-footer">
            <button class="btn-link" @click="goToVerifikasi">Lihat semua {{ pendingItems.length }} pengajuan →</button>
          </div>
        </div>
      </template>

      <!-- ══════════════════════════════════════════════════════
           VERIFIKASI BERKAS
      ══════════════════════════════════════════════════════ -->
      <template v-else-if="activeTab === 'verifikasi'">
        <div class="card">
          <div class="card-header">
            <h3 class="card-title">Daftar Pengajuan Masuk</h3>
            <div class="card-header-actions">
              <div class="search-box">
                <svg class="search-icon" width="14" height="14" viewBox="0 0 24 24" fill="none"><circle cx="11" cy="11" r="8" stroke="currentColor" stroke-width="2"/><path d="M21 21l-4.35-4.35" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>
                <input v-model="searchQuery" class="search-input" placeholder="Cari nama / institusi…" />
              </div>
              <button class="btn-green-sm" @click="fetchPengajuan">
                <svg width="13" height="13" viewBox="0 0 24 24" fill="none"><path d="M23 4v6h-6" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/><path d="M1 20v-6h6" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/><path d="M3.51 9a9 9 0 0114.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0020.49 15" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>
                Refresh
              </button>
            </div>
          </div>

          <!-- Status filter pills -->
          <div class="filter-pills">
            <button
              v-for="f in statusFilters" :key="f.value"
              class="pill"
              :class="{ 'pill--active': filterStatus === f.value }"
              @click="filterStatus = f.value; currentPage = 1"
            >{{ f.label }}
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
                  <th>Nama Peserta</th>
                  <th>Institusi</th>
                  <th>Jurusan</th>
                  <th>Kategori</th>
                  <th>Tanggal Daftar</th>
                  <th>Status</th>
                  <th>Aksi</th>
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
                      <template v-if="['diajukan','menunggu_verifikasi','diproses'].includes(p.status)">
                        <button class="btn-approve" @click="openAction(p, 'diterima')">Terima</button>
                        <button class="btn-reject" @click="openAction(p, 'ditolak')">Tolak</button>
                      </template>
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
      </template>

      <!-- ── Other tabs (peserta, absensi, penilaian, sertifikat, chat) stay as before ── -->
      <template v-else-if="activeTab === 'peserta'">
        <div class="card">
          <div class="card-header"><h3 class="card-title">Peserta Magang Aktif</h3></div>
          <div class="empty-state">
            <div class="empty-state__icon"><svg width="36" height="36" viewBox="0 0 24 24" fill="none"><path d="M17 21v-2a4 4 0 00-4-4H5a4 4 0 00-4 4v2" stroke="#d1d5db" stroke-width="1.5"/><circle cx="9" cy="7" r="4" stroke="#d1d5db" stroke-width="1.5"/></svg></div>
            <p>Belum ada peserta magang aktif.</p>
          </div>
        </div>
      </template>

      <template v-else-if="activeTab === 'absensi'">
        <div class="card">
          <div class="card-header">
            <h3 class="card-title">Rekap Absensi</h3>
            <button class="btn-green-sm">Export PDF</button>
          </div>
          <div class="empty-state">
            <div class="empty-state__icon"><svg width="36" height="36" viewBox="0 0 24 24" fill="none"><rect x="3" y="4" width="18" height="18" rx="2" ry="2" stroke="#d1d5db" stroke-width="1.5"/><line x1="3" y1="10" x2="21" y2="10" stroke="#d1d5db" stroke-width="1.5"/></svg></div>
            <p>Belum ada data absensi.</p>
          </div>
        </div>
      </template>

      <template v-else-if="activeTab === 'penilaian'">
        <div class="card">
          <div class="card-header"><h3 class="card-title">Penilaian Peserta</h3></div>
          <div class="empty-state">
            <div class="empty-state__icon"><svg width="36" height="36" viewBox="0 0 24 24" fill="none"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2" stroke="#d1d5db" stroke-width="1.5"/></svg></div>
            <p>Belum ada peserta yang perlu dinilai.</p>
          </div>
        </div>
      </template>

      <template v-else-if="activeTab === 'sertifikat'">
        <div class="card">
          <div class="card-header"><h3 class="card-title">Generate Sertifikat</h3></div>
          <div class="empty-state">
            <div class="empty-state__icon"><svg width="36" height="36" viewBox="0 0 24 24" fill="none"><circle cx="12" cy="8" r="6" stroke="#d1d5db" stroke-width="1.5"/><path d="M15.477 12.89L17 22l-5-3-5 3 1.523-9.11" stroke="#d1d5db" stroke-width="1.5"/></svg></div>
            <p>Belum ada sertifikat yang dapat di-generate.</p>
          </div>
        </div>
      </template>

      <template v-else-if="activeTab === 'chat'">
        <div class="card">
          <div class="card-header"><h3 class="card-title">Chat Helpdesk</h3><span class="badge badge--gray">0 Tiket</span></div>
          <div class="empty-state">
            <div class="empty-state__icon"><svg width="36" height="36" viewBox="0 0 24 24" fill="none"><path d="M21 15a2 2 0 01-2 2H7l-4 4V5a2 2 0 012-2h14a2 2 0 012 2z" stroke="#d1d5db" stroke-width="1.5"/></svg></div>
            <p>Belum ada tiket masuk dari peserta.</p>
          </div>
        </div>
      </template>

    </template>
  </DashboardLayout>

  <!-- ══════════════════════════════════════════════════════
       DETAIL DRAWER
  ══════════════════════════════════════════════════════ -->
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

        <div v-if="detailLoading" class="drawer-body drawer-loading">
          <div class="spinner"></div>
        </div>

        <div v-else-if="detailData" class="drawer-body">
          <!-- Status banner -->
          <div class="detail-status-banner" :class="`status-banner--${detailData.status}`">
            <span :class="statusBadgeClass(detailData.status)">{{ formatStatus(detailData.status) }}</span>
            <span v-if="detailData.catatan_hrd" class="catatan-hrd">{{ detailData.catatan_hrd }}</span>
          </div>

          <!-- Info grid -->
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

          <!-- Dokumen peserta -->
          <div class="detail-section">
            <div class="detail-section__title">Dokumen Peserta</div>
            <div v-if="dokumenLoading" class="doc-loading"><div class="spinner spinner--sm"></div></div>
            <div v-else-if="dokumenList.length === 0" class="doc-empty">Belum ada dokumen diunggah.</div>
            <div v-else class="doc-list">
              <div v-for="doc in dokumenList" :key="doc.id" class="doc-item">
                <div class="doc-icon">
                  <svg width="16" height="16" viewBox="0 0 24 24" fill="none"><path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z" stroke="currentColor" stroke-width="2"/><polyline points="14 2 14 8 20 8" stroke="currentColor" stroke-width="2"/></svg>
                </div>
                <div class="doc-info">
                  <div class="doc-name">{{ doc.nama_file }}</div>
                  <div class="doc-meta">{{ formatJenisDok(doc.jenis) }} · {{ formatSize(doc.ukuran_bytes) }}</div>
                </div>
                <a class="btn-download" :href="`/api/dokumen/${doc.id}/download`" target="_blank">
                  <svg width="13" height="13" viewBox="0 0 24 24" fill="none"><path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4" stroke="currentColor" stroke-width="2" stroke-linecap="round"/><polyline points="7 10 12 15 17 10" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/><line x1="12" y1="15" x2="12" y2="3" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>
                </a>
              </div>
            </div>
          </div>

          <!-- Upload dokumen balasan -->
          <div class="detail-section">
            <div class="detail-section__title">Upload Dokumen Balasan</div>
            <div class="upload-area" @dragover.prevent @drop.prevent="onDrop">
              <div v-if="!uploadFile" class="upload-placeholder" @click="triggerFileInput">
                <svg width="28" height="28" viewBox="0 0 24 24" fill="none"><path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4" stroke="#9ca3af" stroke-width="1.5" stroke-linecap="round"/><polyline points="17 8 12 3 7 8" stroke="#9ca3af" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/><line x1="12" y1="3" x2="12" y2="15" stroke="#9ca3af" stroke-width="1.5" stroke-linecap="round"/></svg>
                <p>Klik atau seret file ke sini</p>
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
            <input ref="fileInputRef" type="file" accept=".pdf,.doc,.docx" style="display:none" @change="onFileChange" />
            <button
              class="btn-upload"
              :disabled="!uploadFile || uploadLoading"
              @click="submitUpload"
            >
              <svg v-if="!uploadLoading" width="14" height="14" viewBox="0 0 24 24" fill="none"><path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4" stroke="currentColor" stroke-width="2" stroke-linecap="round"/><polyline points="17 8 12 3 7 8" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/><line x1="12" y1="3" x2="12" y2="15" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>
              <div v-else class="spinner spinner--sm spinner--white"></div>
              {{ uploadLoading ? 'Mengupload…' : 'Upload Surat Balasan' }}
            </button>
            <div v-if="uploadSuccess" class="upload-msg upload-msg--ok">✓ Dokumen berhasil diupload.</div>
            <div v-if="uploadError" class="upload-msg upload-msg--err">{{ uploadError }}</div>
          </div>
        </div>

        <!-- Approve / Reject footer -->
        <div v-if="detailData && ['diajukan','menunggu_verifikasi','diproses'].includes(detailData.status)" class="drawer-footer">
          <div class="catatan-row">
            <label class="catatan-label">Catatan HRD (opsional)</label>
            <textarea v-model="actionCatatan" class="catatan-input" rows="2" placeholder="Tambahkan catatan untuk peserta…"></textarea>
          </div>
          <div class="footer-btns">
            <button class="btn-reject-lg" :disabled="actionLoading" @click="submitAction('ditolak')">
              <div v-if="actionLoading && pendingAction === 'ditolak'" class="spinner spinner--sm spinner--white"></div>
              <svg v-else width="15" height="15" viewBox="0 0 24 24" fill="none"><circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2"/><line x1="15" y1="9" x2="9" y2="15" stroke="currentColor" stroke-width="2" stroke-linecap="round"/><line x1="9" y1="9" x2="15" y2="15" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>
              Tolak Pengajuan
            </button>
            <button class="btn-approve-lg" :disabled="actionLoading" @click="submitAction('diterima')">
              <div v-if="actionLoading && pendingAction === 'diterima'" class="spinner spinner--sm spinner--white"></div>
              <svg v-else width="15" height="15" viewBox="0 0 24 24" fill="none"><polyline points="20 6 9 17 4 12" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>
              Terima Pengajuan
            </button>
          </div>
          <div v-if="actionError" class="action-error">{{ actionError }}</div>
        </div>
      </div>
    </div>
  </Teleport>

</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from "vue";
import DashboardLayout from "@/layouts/DashboardLayout.vue";
import { useAuth } from "@/hooks/useAuth";
import api from "@/lib/api";

// ── auth ──────────────────────────────────────────────────────────
const { user } = useAuth();
const layout = ref<InstanceType<typeof DashboardLayout> | null>(null);
const activeTab = computed(() => layout.value?.activeTab ?? "beranda");
const firstName = computed(() => user.value?.nama_lengkap?.split(" ")[0] ?? "");

// ── types ─────────────────────────────────────────────────────────
interface Pengajuan {
  id: string;
  nama_lengkap: string;
  email: string;
  no_hp: string;
  tempat_lahir: string;
  tanggal_lahir: string;
  jenis_kelamin: string;
  alamat: string;
  asal_institusi: string;
  jurusan: string;
  kelas_semester: string;
  nomor_induk: string;
  kategori_magang: string;
  status: string;
  catatan_hrd: string | null;
  created_at: string;
}

interface Dokumen {
  id: string;
  jenis: string;
  nama_file: string;
  ukuran_bytes: number | null;
  uploaded_at: string;
}

// ── pengajuan list state ──────────────────────────────────────────
const allPengajuan = ref<Pengajuan[]>([]);
const listLoading = ref(false);
const listError = ref<string | null>(null);
const statsLoading = ref(false);
const filterStatus = ref("semua");
const searchQuery = ref("");
const currentPage = ref(1);
const PAGE_SIZE = 10;

// ── detail state ──────────────────────────────────────────────────
const showDetail = ref(false);
const detailData = ref<Pengajuan | null>(null);
const detailLoading = ref(false);
const dokumenList = ref<Dokumen[]>([]);
const dokumenLoading = ref(false);

// ── action state ──────────────────────────────────────────────────
const actionCatatan = ref("");
const actionLoading = ref(false);
const pendingAction = ref<string>("");
const actionError = ref<string | null>(null);

// ── upload state ──────────────────────────────────────────────────
const uploadFile = ref<File | null>(null);
const uploadLoading = ref(false);
const uploadSuccess = ref(false);
const uploadError = ref<string | null>(null);
const fileInputRef = ref<HTMLInputElement | null>(null);

// ── computed stats ────────────────────────────────────────────────
const totalAll      = computed(() => allPengajuan.value.length);
const totalPending  = computed(() => allPengajuan.value.filter(p => ["diajukan","menunggu_verifikasi","diproses"].includes(p.status)).length);
const totalDiterima = computed(() => allPengajuan.value.filter(p => p.status === "diterima").length);
const totalDitolak  = computed(() => allPengajuan.value.filter(p => p.status === "ditolak").length);
const pendingItems  = computed(() => allPengajuan.value.filter(p => ["diajukan","menunggu_verifikasi","diproses"].includes(p.status)));

const statusFilters = computed(() => [
  { value: "semua",               label: "Semua",          count: allPengajuan.value.length },
  { value: "diajukan",            label: "Baru Diajukan",  count: allPengajuan.value.filter(p => p.status === "diajukan").length },
  { value: "menunggu_verifikasi", label: "Menunggu Verifikasi", count: allPengajuan.value.filter(p => p.status === "menunggu_verifikasi").length },
  { value: "diproses",            label: "Diproses",       count: allPengajuan.value.filter(p => p.status === "diproses").length },
  { value: "diterima",            label: "Diterima",       count: totalDiterima.value },
  { value: "ditolak",             label: "Ditolak",        count: totalDitolak.value },
]);

const filteredList = computed(() => {
  let list = allPengajuan.value;
  if (filterStatus.value !== "semua") list = list.filter(p => p.status === filterStatus.value);
  const q = searchQuery.value.trim().toLowerCase();
  if (q) list = list.filter(p => p.nama_lengkap.toLowerCase().includes(q) || p.asal_institusi.toLowerCase().includes(q));
  return list;
});

const totalPages = computed(() => Math.max(1, Math.ceil(filteredList.value.length / PAGE_SIZE)));
const paginatedList = computed(() => {
  const start = (currentPage.value - 1) * PAGE_SIZE;
  return filteredList.value.slice(start, start + PAGE_SIZE);
});

// reset page when filter/search changes
watch([filterStatus, searchQuery], () => { currentPage.value = 1; });

// ── data fetching ──────────────────────────────────────────────────
async function fetchPengajuan() {
  listLoading.value = true;
  statsLoading.value = true;
  listError.value = null;
  try {
    const res = await api.get("/api/hrd/pengajuan?page=1&limit=200");
    const body = res.data;
    allPengajuan.value = body.data ?? body ?? [];
  } catch (e: any) {
    listError.value = e.response?.data?.message ?? "Gagal memuat data pengajuan.";
  } finally {
    listLoading.value = false;
    statsLoading.value = false;
  }
}

async function openDetail(id: string) {
  showDetail.value = true;
  detailLoading.value = true;
  detailData.value = null;
  dokumenList.value = [];
  actionCatatan.value = "";
  actionError.value = null;
  uploadFile.value = null;
  uploadSuccess.value = false;
  uploadError.value = null;
  try {
    const res = await api.get(`/api/hrd/pengajuan/${id}`);
    detailData.value = res.data.data ?? res.data;
    fetchDokumen(id);
  } catch {
    closeDetail();
  } finally {
    detailLoading.value = false;
  }
}

async function fetchDokumen(pengajuanId: string) {
  dokumenLoading.value = true;
  try {
    const res = await api.get(`/api/hrd/dokumen/pengajuan/${pengajuanId}`);
    dokumenList.value = res.data.data ?? [];
  } catch {
    dokumenList.value = [];
  } finally {
    dokumenLoading.value = false;
  }
}

function closeDetail() {
  showDetail.value = false;
  detailData.value = null;
  uploadFile.value = null;
}

async function submitAction(type: string) {
  if (!detailData.value) return;
  actionLoading.value = true;
  pendingAction.value = type;
  actionError.value = null;
  try {
    await api.patch(`/api/hrd/pengajuan/${detailData.value.id}/status`, {
      status: type,
      catatan: actionCatatan.value || undefined,
    });
    // Refresh data
    await fetchPengajuan();
    // Refresh detail
    const res = await api.get(`/api/hrd/pengajuan/${detailData.value.id}`);
    detailData.value = res.data.data ?? res.data;
    actionCatatan.value = "";
  } catch (e: any) {
    actionError.value = e.response?.data?.message ?? "Gagal memperbarui status.";
  } finally {
    actionLoading.value = false;
    pendingAction.value = "";
  }
}

function openAction(p: Pengajuan, type: string) {
  openDetail(p.id);
}

function goToVerifikasi() {
  if (layout.value) (layout.value as any).activeTab = "verifikasi";
}

// ── upload ─────────────────────────────────────────────────────────
function triggerFileInput() { fileInputRef.value?.click(); }

function onFileChange(e: Event) {
  const input = e.target as HTMLInputElement;
  if (input.files?.[0]) { uploadFile.value = input.files[0]; uploadSuccess.value = false; uploadError.value = null; }
}

function onDrop(e: DragEvent) {
  const file = e.dataTransfer?.files[0];
  if (file) { uploadFile.value = file; uploadSuccess.value = false; uploadError.value = null; }
}

async function submitUpload() {
  if (!uploadFile.value || !detailData.value) return;
  uploadLoading.value = true;
  uploadSuccess.value = false;
  uploadError.value = null;
  try {
    const form = new FormData();
    form.append("file", uploadFile.value);
    form.append("jenis", "surat_balasan");
    form.append("pengajuan_id", detailData.value.id);
    await api.post("/api/hrd/dokumen/upload", form, { headers: { "Content-Type": "multipart/form-data" } });
    uploadFile.value = null;
    uploadSuccess.value = true;
    fetchDokumen(detailData.value.id);
  } catch (e: any) {
    uploadError.value = e.response?.data?.message ?? "Gagal mengupload dokumen.";
  } finally {
    uploadLoading.value = false;
  }
}

// ── formatting ─────────────────────────────────────────────────────
function formatDate(d: string | null) {
  if (!d) return "-";
  return new Date(d).toLocaleDateString("id-ID", { day: "2-digit", month: "short", year: "numeric" });
}

function formatStatus(s: string) {
  const m: Record<string,string> = {
    diajukan: "Diajukan",
    menunggu_verifikasi: "Menunggu Verifikasi",
    diproses: "Diproses",
    diterima: "Diterima",
    ditolak: "Ditolak",
  };
  return m[s] ?? s;
}

function statusBadgeClass(s: string) {
  const base = "status-badge ";
  if (s === "diterima") return base + "status-badge--green";
  if (s === "ditolak")  return base + "status-badge--red";
  if (s === "diproses") return base + "status-badge--blue";
  return base + "status-badge--yellow";
}

function formatKategori(k: string) {
  const m: Record<string,string> = {
    siswa_smk: "Siswa SMK",
    mahasiswa_d3: "Mahasiswa D3",
    mahasiswa_s1: "Mahasiswa S1",
    penelitian: "Penelitian",
  };
  return m[k] ?? k;
}

function formatJK(jk: string) {
  return jk === "l" ? "Laki-laki" : jk === "p" ? "Perempuan" : jk;
}

function formatJenisDok(j: string) {
  const m: Record<string,string> = {
    proposal_magang: "Proposal",
    ktp: "KTP", ktm: "KTM", pasfoto: "Pas Foto",
    bpjs_kis: "BPJS/KIS", surat_balasan: "Surat Balasan",
    laporan_magang: "Laporan", sertifikat: "Sertifikat",
  };
  return m[j] ?? j;
}

function formatSize(b: number | null) {
  if (!b) return "";
  if (b < 1024) return b + " B";
  if (b < 1048576) return (b/1024).toFixed(1) + " KB";
  return (b/1048576).toFixed(1) + " MB";
}

// ── nav ────────────────────────────────────────────────────────────
const ICON = {
  home:     `<svg width="16" height="16" viewBox="0 0 24 24" fill="none"><path d="M3 9l9-7 9 7v11a2 2 0 01-2 2H5a2 2 0 01-2-2z" stroke="currentColor" stroke-width="2"/><polyline points="9 22 9 12 15 12 15 22" stroke="currentColor" stroke-width="2"/></svg>`,
  verify:   `<svg width="16" height="16" viewBox="0 0 24 24" fill="none"><path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z" stroke="currentColor" stroke-width="2"/><polyline points="14 2 14 8 20 8" stroke="currentColor" stroke-width="2"/><polyline points="9 15 11 17 15 13" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>`,
  users:    `<svg width="16" height="16" viewBox="0 0 24 24" fill="none"><path d="M17 21v-2a4 4 0 00-4-4H5a4 4 0 00-4 4v2" stroke="currentColor" stroke-width="2"/><circle cx="9" cy="7" r="4" stroke="currentColor" stroke-width="2"/></svg>`,
  calendar: `<svg width="16" height="16" viewBox="0 0 24 24" fill="none"><rect x="3" y="4" width="18" height="18" rx="2" ry="2" stroke="currentColor" stroke-width="2"/><line x1="3" y1="10" x2="21" y2="10" stroke="currentColor" stroke-width="2"/></svg>`,
  star:     `<svg width="16" height="16" viewBox="0 0 24 24" fill="none"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2" stroke="currentColor" stroke-width="2"/></svg>`,
  award:    `<svg width="16" height="16" viewBox="0 0 24 24" fill="none"><circle cx="12" cy="8" r="6" stroke="currentColor" stroke-width="2"/><path d="M15.477 12.89L17 22l-5-3-5 3 1.523-9.11" stroke="currentColor" stroke-width="2"/></svg>`,
  chat:     `<svg width="16" height="16" viewBox="0 0 24 24" fill="none"><path d="M21 15a2 2 0 01-2 2H7l-4 4V5a2 2 0 012-2h14a2 2 0 012 2z" stroke="currentColor" stroke-width="2"/></svg>`,
};

const navGroups = computed(() => [
  { items: [{ key: "beranda", label: "Beranda", icon: ICON.home }] },
  {
    label: "Seleksi & Verifikasi",
    items: [
      { key: "verifikasi", label: "Verifikasi Berkas", icon: ICON.verify, badge: totalPending.value || undefined },
      { key: "peserta",    label: "Peserta Aktif",     icon: ICON.users },
    ],
  },
  {
    label: "Pelaksanaan",
    items: [
      { key: "absensi",    label: "Rekap Absensi",     icon: ICON.calendar },
      { key: "penilaian",  label: "Penilaian Peserta", icon: ICON.star },
      { key: "sertifikat", label: "Sertifikat",        icon: ICON.award },
    ],
  },
  { label: "Komunikasi", items: [{ key: "chat", label: "Chat Helpdesk", icon: ICON.chat }] },
]);

// ── lifecycle ──────────────────────────────────────────────────────
onMounted(() => { fetchPengajuan(); });
watch(activeTab, (tab) => { if (tab === "verifikasi") fetchPengajuan(); });
</script>

<style scoped>
/* ── welcome ─────────────────────────────────────────────────────── */
.welcome-banner { background: linear-gradient(135deg, #0d2818 0%, #1a5c20 100%); border-radius: 14px; padding: 24px 28px; color: #fff; display: flex; align-items: center; justify-content: space-between; gap: 16px; }
.welcome-banner__greeting { font-size: 17px; font-weight: 700; }
.welcome-banner__sub { font-size: 12.5px; color: rgba(255,255,255,0.65); margin-top: 4px; }
.welcome-banner__icon { opacity: 0.8; }

/* ── stats ───────────────────────────────────────────────────────── */
.stats-row { display: grid; grid-template-columns: repeat(4,1fr); gap: 14px; }
.stat-card { background: #fff; border-radius: 12px; padding: 16px; display: flex; align-items: center; gap: 12px; border: 1px solid #e9f5e9; box-shadow: 0 1px 3px rgba(13,40,24,0.05); }
.stat-card__icon { width: 38px; height: 38px; border-radius: 10px; display: flex; align-items: center; justify-content: center; flex-shrink: 0; }
.stat-card__label { font-size: 11px; color: #6b7280; font-weight: 500; }
.stat-card__value { font-size: 20px; font-weight: 700; color: #111827; margin-top: 1px; }

/* ── card ────────────────────────────────────────────────────────── */
.card { background: #fff; border-radius: 14px; border: 1px solid #e9f5e9; box-shadow: 0 1px 3px rgba(13,40,24,0.05); overflow: hidden; }
.card-header { display: flex; align-items: center; justify-content: space-between; padding: 16px 20px; border-bottom: 1px solid #f0faf0; gap: 12px; flex-wrap: wrap; }
.card-title { font-size: 13.5px; font-weight: 700; color: #111827; margin: 0; }
.card-header-actions { display: flex; align-items: center; gap: 8px; }
.card-footer { padding: 12px 20px; border-top: 1px solid #f3f4f6; text-align: center; }
.btn-link { background: none; border: none; color: #48AF4A; font-size: 12.5px; font-weight: 600; cursor: pointer; font-family: inherit; }
.btn-link:hover { text-decoration: underline; }

/* ── search ──────────────────────────────────────────────────────── */
.search-box { position: relative; display: flex; align-items: center; }
.search-icon { position: absolute; left: 10px; color: #9ca3af; pointer-events: none; }
.search-input { border: 1px solid #e5e7eb; border-radius: 8px; padding: 7px 12px 7px 30px; font-size: 12.5px; font-family: inherit; color: #374151; width: 220px; outline: none; }
.search-input:focus { border-color: #48AF4A; }

/* ── filter pills ────────────────────────────────────────────────── */
.filter-pills { display: flex; gap: 6px; padding: 12px 20px; border-bottom: 1px solid #f3f4f6; flex-wrap: wrap; }
.pill { background: #f9fafb; border: 1px solid #e5e7eb; border-radius: 100px; padding: 5px 13px; font-size: 12px; font-weight: 500; cursor: pointer; color: #6b7280; font-family: inherit; display: flex; align-items: center; gap: 5px; transition: all 0.15s; }
.pill:hover { border-color: #48AF4A; color: #16a34a; }
.pill--active { background: #f0fdf4; border-color: #48AF4A; color: #16a34a; font-weight: 600; }
.pill-count { background: #e5e7eb; color: #6b7280; border-radius: 100px; padding: 0 6px; font-size: 10.5px; font-weight: 700; }
.pill--active .pill-count { background: #bbf7d0; color: #15803d; }

/* ── table ───────────────────────────────────────────────────────── */
.table-wrap { overflow-x: auto; }
.data-table { width: 100%; border-collapse: collapse; font-size: 13px; }
.data-table th { padding: 11px 16px; text-align: left; font-size: 10.5px; font-weight: 600; color: #6b7280; background: #f9fafb; border-bottom: 1px solid #f1f5f9; text-transform: uppercase; letter-spacing: 0.04em; white-space: nowrap; }
.data-table td { padding: 13px 16px; border-bottom: 1px solid #f9fafb; color: #374151; vertical-align: middle; }
.table-empty { text-align: center; color: #9ca3af; padding: 32px 16px; }
.td-name { min-width: 160px; }

.name-cell { display: flex; align-items: center; gap: 10px; }
.name-avatar { width: 32px; height: 32px; border-radius: 8px; background: linear-gradient(135deg, #48AF4A, #2d8f30); color: #fff; font-size: 13px; font-weight: 700; display: flex; align-items: center; justify-content: center; flex-shrink: 0; }
.name-text { font-weight: 600; color: #111827; font-size: 12.5px; }
.name-sub { font-size: 11px; color: #9ca3af; }

.tag { background: #eff6ff; color: #1d4ed8; border-radius: 6px; padding: 2px 8px; font-size: 11px; font-weight: 600; white-space: nowrap; }

/* ── status badges ───────────────────────────────────────────────── */
.status-badge { padding: 3px 10px; border-radius: 100px; font-size: 11px; font-weight: 600; white-space: nowrap; }
.status-badge--yellow { background: #fefce8; color: #ca8a04; }
.status-badge--green  { background: #f0fdf4; color: #16a34a; }
.status-badge--red    { background: #fef2f2; color: #dc2626; }
.status-badge--blue   { background: #eff6ff; color: #2563eb; }

.badge { padding: 3px 10px; border-radius: 100px; font-size: 11px; font-weight: 600; }
.badge--yellow { background: #fefce8; color: #ca8a04; }
.badge--gray   { background: #f3f4f6; color: #6b7280; }

/* ── action buttons in table ─────────────────────────────────────── */
.action-btns { display: flex; gap: 5px; align-items: center; flex-wrap: nowrap; }
.btn-detail  { background: #f3f4f6; color: #374151; border: none; border-radius: 7px; padding: 5px 11px; font-size: 11.5px; font-weight: 600; font-family: inherit; cursor: pointer; white-space: nowrap; }
.btn-detail:hover { background: #e5e7eb; }
.btn-approve { background: #f0fdf4; color: #16a34a; border: 1px solid #bbf7d0; border-radius: 7px; padding: 5px 11px; font-size: 11.5px; font-weight: 600; font-family: inherit; cursor: pointer; white-space: nowrap; }
.btn-approve:hover { background: #dcfce7; }
.btn-reject  { background: #fef2f2; color: #dc2626; border: 1px solid #fecaca; border-radius: 7px; padding: 5px 11px; font-size: 11.5px; font-weight: 600; font-family: inherit; cursor: pointer; white-space: nowrap; }
.btn-reject:hover { background: #fee2e2; }

/* ── pagination ──────────────────────────────────────────────────── */
.pagination { display: flex; align-items: center; justify-content: center; gap: 12px; padding: 14px 20px; border-top: 1px solid #f3f4f6; }
.page-btn { background: #f3f4f6; border: none; border-radius: 7px; width: 30px; height: 30px; font-size: 16px; cursor: pointer; display: flex; align-items: center; justify-content: center; color: #374151; }
.page-btn:disabled { opacity: 0.4; cursor: default; }
.page-btn:not(:disabled):hover { background: #e5e7eb; }
.page-info { font-size: 12.5px; color: #6b7280; }

/* ── empty / loading ─────────────────────────────────────────────── */
.empty-state { display: flex; flex-direction: column; align-items: center; padding: 44px 24px; gap: 12px; text-align: center; }
.empty-state__icon { width: 72px; height: 72px; background: #f9fafb; border-radius: 50%; display: flex; align-items: center; justify-content: center; }
.empty-state p { font-size: 13px; color: #9ca3af; line-height: 1.7; margin: 0; }
.text-red { color: #dc2626; font-size: 13px; }

.spinner { width: 24px; height: 24px; border: 2.5px solid #e5e7eb; border-top-color: #48AF4A; border-radius: 50%; animation: spin 0.7s linear infinite; }
.spinner--sm { width: 14px; height: 14px; border-width: 2px; }
.spinner--white { border-color: rgba(255,255,255,0.3); border-top-color: #fff; }
@keyframes spin { to { transform: rotate(360deg); } }

/* ── drawer ──────────────────────────────────────────────────────── */
.drawer-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.45); backdrop-filter: blur(2px); z-index: 100; display: flex; justify-content: flex-end; }
.drawer { width: min(560px, 100vw); height: 100vh; background: #fff; display: flex; flex-direction: column; box-shadow: -4px 0 24px rgba(0,0,0,0.12); overflow: hidden; }

.drawer-header { display: flex; align-items: flex-start; justify-content: space-between; padding: 20px 24px 16px; border-bottom: 1px solid #f0faf0; flex-shrink: 0; }
.drawer-title { font-size: 16px; font-weight: 700; color: #111827; margin: 0 0 2px; }
.drawer-sub { font-size: 12.5px; color: #6b7280; margin: 0; }
.drawer-close { background: #f3f4f6; border: none; border-radius: 8px; width: 32px; height: 32px; display: flex; align-items: center; justify-content: center; cursor: pointer; color: #6b7280; flex-shrink: 0; margin-top: 2px; }
.drawer-close:hover { background: #e5e7eb; color: #111827; }

.drawer-body { flex: 1; overflow-y: auto; padding: 20px 24px; display: flex; flex-direction: column; gap: 20px; }
.drawer-loading { align-items: center; justify-content: center; }

/* ── detail sections ─────────────────────────────────────────────── */
.detail-status-banner { display: flex; align-items: center; gap: 10px; padding: 10px 14px; background: #f9fafb; border-radius: 10px; }
.catatan-hrd { font-size: 12px; color: #6b7280; font-style: italic; }
.status-banner--diterima { background: #f0fdf4; }
.status-banner--ditolak  { background: #fef2f2; }

.detail-section { display: flex; flex-direction: column; gap: 10px; }
.detail-section__title { font-size: 11px; font-weight: 700; color: #9ca3af; text-transform: uppercase; letter-spacing: 0.08em; }
.info-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 10px; }
.info-item { display: flex; flex-direction: column; gap: 2px; }
.info-item--full { grid-column: 1 / -1; }
.info-label { font-size: 11px; color: #9ca3af; font-weight: 500; }
.info-val { font-size: 13px; color: #111827; font-weight: 500; }

/* ── dokumen list ────────────────────────────────────────────────── */
.doc-loading { display: flex; justify-content: center; padding: 16px; }
.doc-empty { font-size: 12.5px; color: #9ca3af; text-align: center; padding: 12px; background: #f9fafb; border-radius: 8px; }
.doc-list { display: flex; flex-direction: column; gap: 6px; }
.doc-item { display: flex; align-items: center; gap: 10px; padding: 10px 12px; background: #f9fafb; border-radius: 9px; border: 1px solid #f3f4f6; }
.doc-icon { width: 32px; height: 32px; background: #eff6ff; border-radius: 7px; display: flex; align-items: center; justify-content: center; color: #2563eb; flex-shrink: 0; }
.doc-info { flex: 1; min-width: 0; }
.doc-name { font-size: 12.5px; font-weight: 600; color: #111827; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.doc-meta { font-size: 11px; color: #9ca3af; margin-top: 1px; }
.btn-download { width: 28px; height: 28px; background: #f3f4f6; border-radius: 7px; display: flex; align-items: center; justify-content: center; color: #374151; text-decoration: none; flex-shrink: 0; }
.btn-download:hover { background: #e5e7eb; }

/* ── upload area ─────────────────────────────────────────────────── */
.upload-area { border: 2px dashed #d1d5db; border-radius: 10px; cursor: pointer; transition: border-color 0.15s; }
.upload-area:hover { border-color: #48AF4A; }
.upload-placeholder { display: flex; flex-direction: column; align-items: center; gap: 6px; padding: 24px 16px; }
.upload-placeholder p { font-size: 13px; color: #374151; font-weight: 500; margin: 0; }
.upload-placeholder span { font-size: 11px; color: #9ca3af; }
.upload-selected { display: flex; align-items: center; gap: 10px; padding: 12px 16px; }
.upload-filename { flex: 1; font-size: 12.5px; font-weight: 500; color: #111827; min-width: 0; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.upload-remove { background: none; border: none; color: #9ca3af; cursor: pointer; display: flex; align-items: center; padding: 2px; }
.upload-remove:hover { color: #dc2626; }
.btn-upload { width: 100%; margin-top: 10px; background: #48AF4A; color: #fff; border: none; border-radius: 9px; padding: 10px 16px; font-size: 13px; font-weight: 600; font-family: inherit; cursor: pointer; display: flex; align-items: center; justify-content: center; gap: 7px; transition: background 0.15s; }
.btn-upload:hover:not(:disabled) { background: #3d9e3f; }
.btn-upload:disabled { opacity: 0.5; cursor: default; }
.upload-msg { font-size: 12px; margin-top: 6px; padding: 7px 12px; border-radius: 7px; }
.upload-msg--ok  { background: #f0fdf4; color: #16a34a; }
.upload-msg--err { background: #fef2f2; color: #dc2626; }

/* ── drawer footer ───────────────────────────────────────────────── */
.drawer-footer { flex-shrink: 0; padding: 16px 24px; border-top: 1px solid #f0faf0; background: #fff; display: flex; flex-direction: column; gap: 12px; }
.catatan-label { font-size: 11.5px; font-weight: 600; color: #374151; display: block; margin-bottom: 5px; }
.catatan-input { width: 100%; border: 1px solid #e5e7eb; border-radius: 8px; padding: 8px 12px; font-size: 12.5px; font-family: inherit; color: #374151; resize: none; outline: none; box-sizing: border-box; }
.catatan-input:focus { border-color: #48AF4A; }
.footer-btns { display: flex; gap: 10px; }
.btn-reject-lg  { flex: 1; background: #fef2f2; color: #dc2626; border: 1px solid #fecaca; border-radius: 9px; padding: 11px 16px; font-size: 13px; font-weight: 600; font-family: inherit; cursor: pointer; display: flex; align-items: center; justify-content: center; gap: 7px; }
.btn-reject-lg:hover:not(:disabled)  { background: #fee2e2; }
.btn-approve-lg { flex: 1; background: #48AF4A; color: #fff; border: none; border-radius: 9px; padding: 11px 16px; font-size: 13px; font-weight: 600; font-family: inherit; cursor: pointer; display: flex; align-items: center; justify-content: center; gap: 7px; }
.btn-approve-lg:hover:not(:disabled) { background: #3d9e3f; }
.btn-approve-lg:disabled, .btn-reject-lg:disabled { opacity: 0.5; cursor: default; }
.action-error { font-size: 12px; color: #dc2626; background: #fef2f2; padding: 7px 12px; border-radius: 7px; }

.btn-green-sm { background: #48AF4A; color: #fff; border: none; border-radius: 8px; padding: 6px 14px; font-size: 12px; font-weight: 600; font-family: "Poppins", sans-serif; cursor: pointer; white-space: nowrap; display: flex; align-items: center; gap: 5px; }
.btn-green-sm:hover { background: #3d9e3f; }

@media (max-width: 700px) { .stats-row { grid-template-columns: 1fr 1fr; } }
@media (max-width: 420px) { .stats-row { grid-template-columns: 1fr; } }
@media (max-width: 600px) { .info-grid { grid-template-columns: 1fr; } .drawer { width: 100vw; } }
</style>
