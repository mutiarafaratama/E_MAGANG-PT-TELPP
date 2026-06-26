<template>
  <DashboardLayout :nav-groups="navGroups" role-name="Super Admin" default-tab="beranda" ref="layout">
    <template #default>

      <!-- ── BERANDA ── -->
      <template v-if="activeTab === 'beranda'">
        <div class="welcome-banner">
          <div>
            <div class="welcome-banner__greeting">Selamat datang, {{ firstName }}! </div>
            <div class="welcome-banner__sub">Kelola seluruh sistem e-Magang PT TELPP dari sini.</div>
          </div>
          <div class="welcome-banner__badge">ADMIN</div>
        </div>

        <div class="stats-row">
          <div class="stat-card">
            <div class="stat-card__icon" style="background:#f0fdf4;color:#3b82f6">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none"><path d="M17 21v-2a4 4 0 00-4-4H5a4 4 0 00-4 4v2" stroke="currentColor" stroke-width="2"/><circle cx="9" cy="7" r="4" stroke="currentColor" stroke-width="2"/></svg>
            </div>
            <div>
              <div class="stat-card__label">Total Peserta</div>
              <div class="stat-card__value">{{ statsLoading ? '…' : stats.total_peserta }}</div>
              <div class="stat-card__sub">Terdaftar di sistem</div>
            </div>
          </div>
          <div class="stat-card">
            <div class="stat-card__icon" style="background:#fefce8;color:#ca8a04">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none"><circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2"/><polyline points="12 6 12 12 16 14" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>
            </div>
            <div>
              <div class="stat-card__label">Pengajuan Menunggu</div>
              <div class="stat-card__value">{{ statsLoading ? '…' : stats.pengajuan_menunggu }}</div>
              <div class="stat-card__sub">Perlu diverifikasi</div>
            </div>
          </div>
          <div class="stat-card">
            <div class="stat-card__icon" style="background:#f0fdf4;color:#16a34a">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none"><polyline points="20 6 9 17 4 12" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>
            </div>
            <div>
              <div class="stat-card__label">Magang Aktif</div>
              <div class="stat-card__value">{{ statsLoading ? '…' : stats.magang_aktif }}</div>
              <div class="stat-card__sub">Sedang berjalan</div>
            </div>
          </div>
          <div class="stat-card">
            <div class="stat-card__icon" style="background:#fff1f2;color:#e11d48">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none"><path d="M21 15a2 2 0 01-2 2H7l-4 4V5a2 2 0 012-2h14a2 2 0 012 2z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>
            </div>
            <div>
              <div class="stat-card__label">Tiket Chat</div>
              <div class="stat-card__value">{{ statsLoading ? '…' : stats.tiket_chat_menunggu }}</div>
              <div class="stat-card__sub">Menunggu balasan</div>
            </div>
          </div>
        </div>

        <!-- ── Akun Segera Dihapus ── -->
        <div class="card">
          <div class="card-header">
            <div>
              <h3 class="card-title">Jadwal Hapus Otomatis</h3>
              <p class="card-sub">Peserta yang selesai magang &gt; 23 hari lalu — akun dihapus dalam 7 hari ke depan</p>
            </div>
            <span class="jt-count" v-if="!jtLoading">{{ jtList.length }} akun</span>
          </div>

          <div v-if="jtLoading" class="empty-state"><div class="spinner"></div></div>
          <div v-else-if="jtList.length === 0" class="empty-state">
            <div class="empty-state__icon">
              <svg width="28" height="28" viewBox="0 0 24 24" fill="none"><path d="M17 21v-2a4 4 0 00-4-4H5a4 4 0 00-4 4v2" stroke="#d1d5db" stroke-width="1.5"/><circle cx="9" cy="7" r="4" stroke="#d1d5db" stroke-width="1.5"/></svg>
            </div>
            <p>Tidak ada akun yang dijadwalkan dihapus dalam 7 hari ke depan</p>
          </div>
          <div v-else class="table-wrap">
            <table class="data-table">
              <thead>
                <tr>
                  <th>Nama Peserta</th>
                  <th>Email</th>
                  <th>Asal Institusi</th>
                  <th>Divisi</th>
                  <th>Selesai Magang</th>
                  <th>Dihapus Pada</th>
                  <th>Sisa Hari</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="a in jtList" :key="a.user_id">
                  <td class="td-nama">{{ a.nama_lengkap }}</td>
                  <td style="color:#6b7280">{{ a.email }}</td>
                  <td>{{ a.asal_institusi || '—' }}</td>
                  <td>{{ a.divisi || '—' }}</td>
                  <td>{{ a.tanggal_selesai }}</td>
                  <td style="font-weight:600">{{ a.hapus_at }}</td>
                  <td>
                    <span :class="a.sisa_hari <= 2 ? 'status-pill status-pill--red' : 'status-pill status-pill--yellow'">
                      {{ a.sisa_hari }} hari
                    </span>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </template>

      <!-- ── MANAJEMEN USER ── -->
      <AdminUsers v-else-if="activeTab === 'users'" />

      <!-- ── HERO SECTION ── -->
      <AdminHero v-else-if="activeTab === 'hero'" />

      <!-- ── JADWAL & PERIODE ── -->
      <AdminJadwal v-else-if="activeTab === 'jadwal'" />

      <!-- ── LANDING ── -->
      <template v-else-if="activeTab === 'landing'">
        <div class="card">
          <div class="card-header">
            <div>
              <h3 class="card-title">Alur Pendaftaran</h3>
              <p class="card-sub">Tampil sebagai carousel di landing page — tiap item memiliki judul, paragraf, dan gambar</p>
            </div>
            <div style="display:flex;gap:8px;align-items:center">
              <a href="/" target="_blank" class="btn-preview-sm">Lihat Landing ↗</a>
              <button class="btn-green-sm" @click="openAlurModal(null)">+ Tambah Alur</button>
            </div>
          </div>

          <div v-if="alurLoading" class="empty-state"><div class="spinner"></div></div>
          <div v-else-if="alurList.length === 0" class="empty-state">
            <div class="empty-state__icon">
              <svg width="36" height="36" viewBox="0 0 24 24" fill="none">
                <path d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" stroke="#d1d5db" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
              </svg>
            </div>
            <p>Belum ada item alur. Klik <strong>+ Tambah Alur</strong> untuk mulai.</p>
          </div>
          <div v-else class="alur-admin-list">
            <div v-for="(item, idx) in alurList" :key="item.id" class="alur-admin-item">
              <div class="alur-admin-item__num">{{ idx + 1 }}</div>
              <div class="alur-admin-item__img" v-if="item.gambar_url">
                <img :src="item.gambar_url" alt="" />
              </div>
              <div class="alur-admin-item__img alur-admin-item__img--empty" v-else>
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
                  <path d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" stroke="#d1d5db" stroke-width="2" stroke-linecap="round"/>
                </svg>
              </div>
              <div class="alur-admin-item__content">
                <div class="alur-admin-item__judul">{{ item.judul }}</div>
                <div class="alur-admin-item__paragraf">{{ item.paragraf }}</div>
              </div>
              <div class="alur-admin-item__actions">
                <button class="btn-aksi btn-aksi--ghost" @click="openAlurModal(item)">Edit</button>
                <button class="btn-aksi btn-aksi--red" @click="deleteAlur(item)">Hapus</button>
              </div>
            </div>
          </div>
        </div>
      </template>

      <!-- ── FAQ ── -->
      <template v-else-if="activeTab === 'faq'">
        <div class="card">
          <div class="card-header">
            <div>
              <h3 class="card-title">Kelola FAQ</h3>
              <p class="card-sub">FAQ yang aktif tampil di landing page secara otomatis</p>
            </div>
            <button class="btn-green-sm" @click="openFaqForm(null)">+ Tambah FAQ</button>
          </div>

          <div v-if="faqLoading" class="empty-state"><div class="spinner"></div></div>
          <div v-else-if="faqList.length === 0" class="empty-state">
            <div class="empty-state__icon"><svg width="36" height="36" viewBox="0 0 24 24" fill="none"><circle cx="12" cy="12" r="10" stroke="#d1d5db" stroke-width="1.5"/><path d="M9.09 9a3 3 0 015.83 1c0 2-3 3-3 3" stroke="#d1d5db" stroke-width="1.5" stroke-linecap="round"/><line x1="12" y1="17" x2="12.01" y2="17" stroke="#d1d5db" stroke-width="2" stroke-linecap="round"/></svg></div>
            <p>Belum ada FAQ. Klik "+ Tambah FAQ" untuk menambah pertanyaan.</p>
          </div>
          <div v-else-if="faqList.length > 0" class="faq-admin-list">
            <div v-for="f in faqList" :key="f.id" class="faq-admin-item">
              <div class="faq-admin-item__num">{{ String(f.urutan).padStart(2,'0') }}</div>
              <div class="faq-admin-item__content">
                <div class="faq-admin-item__q">{{ f.pertanyaan }}</div>
                <div class="faq-admin-item__a">{{ f.jawaban }}</div>
              </div>
              <div class="faq-admin-item__actions">
                <span
                  :class="f.is_active ? 'status-pill status-pill--green' : 'status-pill status-pill--gray'"
                  style="cursor:pointer" title="Klik untuk toggle"
                  @click="toggleFaqItem(f)"
                >{{ f.is_active ? 'Aktif' : 'Nonaktif' }}</span>
                <button class="btn-aksi btn-aksi--ghost" @click="openFaqForm(f)">Edit</button>
                <button class="btn-aksi btn-aksi--red" @click="deleteFaqItem(f.id!)">Hapus</button>
              </div>
            </div>
          </div>
        </div>
      </template>

      <!-- ── JAM ABSEN ── -->
      <template v-else-if="activeTab === 'jamabsen'">
        <div class="card">
          <div class="card-header">
            <h3 class="card-title">Pengaturan Jam Absensi</h3>
            <span class="badge-info">Berlaku untuk semua peserta aktif</span>
          </div>
          <div v-if="cfgLoading" class="empty-state"><div class="spinner"></div></div>
          <div v-else class="jam-form">
            <div class="jam-section">
              <div class="jam-section__label">Sesi Absen Masuk</div>
              <div class="jam-row">
                <div class="jam-field">
                  <label>Buka</label>
                  <input type="text" v-model="jamForm.jam_masuk_buka" class="jam-input" placeholder="HH:MM" maxlength="5" @input="autoFormatTime('jam_masuk_buka', $event)" />
                </div>
                <span class="jam-sep">–</span>
                <div class="jam-field">
                  <label>Tutup</label>
                  <input type="text" v-model="jamForm.jam_masuk_tutup" class="jam-input" placeholder="HH:MM" maxlength="5" @input="autoFormatTime('jam_masuk_tutup', $event)" />
                </div>
              </div>
            </div>
            <div class="jam-section">
              <div class="jam-section__label">Sesi Absen Pulang</div>
              <div class="jam-row">
                <div class="jam-field">
                  <label>Buka</label>
                  <input type="text" v-model="jamForm.jam_pulang_buka" class="jam-input" placeholder="HH:MM" maxlength="5" @input="autoFormatTime('jam_pulang_buka', $event)" />
                </div>
                <span class="jam-sep">–</span>
                <div class="jam-field">
                  <label>Tutup</label>
                  <input type="text" v-model="jamForm.jam_pulang_tutup" class="jam-input" placeholder="HH:MM" maxlength="5" @input="autoFormatTime('jam_pulang_tutup', $event)" />
                </div>
              </div>
            </div>

            <div v-if="cfgError" class="cfg-error">{{ cfgError }}</div>
            <div v-if="cfgSuccess" class="cfg-success">{{ cfgSuccess }}</div>
            <div class="jam-actions">
              <button class="btn-green" @click="saveJam" :disabled="cfgSaving">
                {{ cfgSaving ? 'Menyimpan...' : 'Simpan Pengaturan' }}
              </button>
            </div>
          </div>
        </div>
      </template>

      <!-- ── KELOLA DIVISI ── -->
      <template v-else-if="activeTab === 'divisi'">
        <div class="card">
          <div class="card-header">
            <div>
              <h3 class="card-title">Kelola Divisi</h3>
              <p class="card-sub">Daftar divisi/unit kerja — termasuk geofencing absensi per-divisi</p>
            </div>
            <button class="btn-green-sm" @click="openDivisiModal(null)">+ Tambah Divisi</button>
          </div>
          <div v-if="divisiLoading" class="empty-state"><div class="spinner"></div></div>
          <div v-else-if="divisiList.length === 0" class="empty-state">
            <div class="empty-state__icon">
              <svg width="36" height="36" viewBox="0 0 24 24" fill="none"><rect x="2" y="7" width="20" height="14" rx="2" stroke="#d1d5db" stroke-width="1.5"/><path d="M16 7V5a2 2 0 00-2-2h-4a2 2 0 00-2 2v2" stroke="#d1d5db" stroke-width="1.5"/></svg>
            </div>
            <p>Belum ada divisi. Tambahkan divisi pertama.</p>
          </div>
          <div v-else class="table-wrap">
            <table class="data-table">
              <thead>
                <tr><th>Nama Divisi</th><th>Urutan</th><th>Geofencing</th><th>Status</th><th>Dibuat</th><th>Aksi</th></tr>
              </thead>
              <tbody>
                <tr v-for="d in divisiList" :key="d.id">
                  <td class="td-nama">{{ d.nama }}</td>
                  <td><span class="urutan-badge">{{ d.urutan }}</span></td>
                  <td>
                    <div v-if="d.geo_lat != null" class="geo-badge geo-badge--set">
                      <svg width="11" height="11" viewBox="0 0 24 24" fill="none"><path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7z" fill="currentColor"/><circle cx="12" cy="9" r="2.5" fill="#fff"/></svg>
                      {{ d.nama_lokasi || d.nama }}&nbsp;({{ d.geo_radius }}m)
                    </div>
                    <span v-else class="geo-badge geo-badge--none">Belum diset</span>
                  </td>
                  <td>
                    <span :class="d.is_active ? 'status-pill status-pill--green' : 'status-pill status-pill--gray'">
                      {{ d.is_active ? 'Aktif' : 'Nonaktif' }}
                    </span>
                  </td>
                  <td>{{ formatDate(d.created_at) }}</td>
                  <td>
                    <div class="aksi-cell">
                      <button class="btn-aksi btn-aksi--ghost" @click="openDivisiModal(d)">Edit</button>
                      <button
                        :class="d.is_active ? 'btn-aksi btn-aksi--warn' : 'btn-aksi btn-aksi--green'"
                        @click="toggleDivisi(d)"
                      >{{ d.is_active ? 'Nonaktifkan' : 'Aktifkan' }}</button>
                      <button class="btn-aksi btn-aksi--red" @click="deleteDivisi(d)">Hapus</button>
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
        <!-- info card -->
        <div class="geo-info-card">
          <svg width="15" height="15" viewBox="0 0 24 24" fill="none"><circle cx="12" cy="12" r="10" stroke="#3b82f6" stroke-width="2"/><line x1="12" y1="8" x2="12" y2="12" stroke="#3b82f6" stroke-width="2" stroke-linecap="round"/><line x1="12" y1="16" x2="12.01" y2="16" stroke="#3b82f6" stroke-width="2" stroke-linecap="round"/></svg>
          <span><strong>Geofencing per-divisi:</strong> Jika divisi memiliki geofencing, peserta di divisi tersebut hanya bisa absen dari area yang ditentukan. Divisi tanpa geofencing tidak memerlukan validasi lokasi saat absensi.</span>
        </div>
      </template>

      <!-- ── KELOLA DOKUMEN ── -->
      <AdminDokumen v-else-if="activeTab === 'dokumen'" />

      <!-- ── KNOWLEDGE BASE ── -->
      <AdminKnowledge v-else-if="activeTab === 'knowledge'" />

      <!-- ── PROFIL ── -->
      <template v-else-if="activeTab === 'profil'">
        <StaffProfil />
      </template>

    </template>

  </DashboardLayout>

  <!-- ── MODAL TAMBAH/EDIT FAQ ── -->
  <Teleport to="body">
    <div v-if="faqFormVisible" class="modal-overlay" @click.self="closeFaqForm">
      <div class="divisi-modal" style="width:min(520px,100%)">
        <div class="divisi-modal__header">
          <h3 class="divisi-modal__title">{{ faqEditTarget ? 'Edit FAQ' : 'Tambah FAQ Baru' }}</h3>
          <button class="modal-close" @click="closeFaqForm">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none"><path d="M18 6L6 18M6 6l12 12" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"/></svg>
          </button>
        </div>
        <div class="divisi-modal__body" style="gap:14px">
          <div class="dform-group">
            <label class="dform-label">Pertanyaan <span class="dform-req">*</span></label>
            <input v-model="faqFormData.pertanyaan" class="dform-input" placeholder="Tulis pertanyaan yang sering diajukan…" :disabled="faqSaving" />
          </div>
          <div class="dform-group">
            <label class="dform-label">Jawaban <span class="dform-req">*</span></label>
            <textarea v-model="faqFormData.jawaban" class="dform-input faq-textarea" rows="4" placeholder="Tulis jawaban yang jelas dan informatif…" :disabled="faqSaving"></textarea>
          </div>
          <div class="faq-form-row">
            <div class="dform-group">
              <label class="dform-label">Urutan</label>
              <input v-model.number="faqFormData.urutan" type="number" min="1" class="dform-input dform-input--sm" :disabled="faqSaving" />
            </div>
            <div class="dform-group">
              <label class="dform-label">Status</label>
              <select v-model="faqFormData.is_active" class="dform-input dform-input--sm" :disabled="faqSaving">
                <option :value="true">Aktif</option>
                <option :value="false">Nonaktif</option>
              </select>
            </div>
          </div>
          <div v-if="faqFormError" class="dform-error">{{ faqFormError }}</div>
        </div>
        <div class="divisi-modal__footer">
          <button class="btn-cancel" @click="closeFaqForm" :disabled="faqSaving">Batal</button>
          <button class="btn-green" @click="saveFaqItem" :disabled="faqSaving">
            <div v-if="faqSaving" class="spinner-sm"></div>
            {{ faqSaving ? 'Menyimpan…' : (faqEditTarget ? 'Simpan Perubahan' : 'Tambah FAQ') }}
          </button>
        </div>
      </div>
    </div>
  </Teleport>

  <!-- ── MODAL TAMBAH/EDIT ALUR ── -->
  <Teleport to="body">
    <div v-if="showAlurModal" class="modal-overlay" @click.self="showAlurModal = false">
      <div class="divisi-modal" style="width:min(520px,100%)">
        <div class="divisi-modal__header">
          <h3 class="divisi-modal__title">{{ alurForm.id ? 'Edit Item Alur' : 'Tambah Item Alur' }}</h3>
          <button class="modal-close" @click="showAlurModal = false">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none"><path d="M18 6L6 18M6 6l12 12" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"/></svg>
          </button>
        </div>
        <div class="divisi-modal__body" style="gap:14px">
          <div class="dform-group">
            <label class="dform-label">Judul <span class="dform-req">*</span></label>
            <input v-model="alurForm.judul" type="text" class="dform-input" placeholder="contoh: Isi Formulir Pendaftaran" :disabled="alurSaving"/>
          </div>
          <div class="dform-group">
            <label class="dform-label">Paragraf / Deskripsi</label>
            <textarea v-model="alurForm.paragraf" class="dform-input" rows="3" placeholder="Penjelasan singkat tentang langkah ini..." :disabled="alurSaving" style="resize:vertical"></textarea>
          </div>
          <div class="dform-group">
            <label class="dform-label">Gambar <span class="dform-opt">(opsional)</span></label>
            <div class="alur-upload-area" @click="triggerFileInput" :class="{ 'alur-upload-area--uploading': alurUploading }">
              <input
                ref="alurFileInput"
                type="file"
                accept=".jpg,.jpeg,.png,.webp,.gif"
                style="display:none"
                @change="onAlurFileChange"
              />
              <template v-if="alurUploadPreview">
                <img :src="alurUploadPreview" class="alur-upload-preview" alt="Preview" />
                <div class="alur-upload-change">Klik untuk ganti gambar</div>
              </template>
              <template v-else-if="alurForm.gambar_url && !alurUploadPreview">
                <img :src="alurForm.gambar_url" class="alur-upload-preview" alt="Gambar saat ini" />
                <div class="alur-upload-change">Klik untuk ganti gambar</div>
              </template>
              <template v-else>
                <div class="alur-upload-icon">
                  <svg width="28" height="28" viewBox="0 0 24 24" fill="none"><path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4M17 8l-5-5-5 5M12 3v12" stroke="#94a3b8" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>
                </div>
                <div class="alur-upload-text">Klik untuk upload gambar</div>
                <div class="alur-upload-hint">JPG, PNG, WebP · Maks 5 MB</div>
              </template>
              <div v-if="alurUploading" class="alur-upload-overlay">
                <div class="spinner"></div>
              </div>
            </div>
            <div v-if="alurUploadError" class="dform-error" style="margin-top:6px">{{ alurUploadError }}</div>
          </div>
          <div class="dform-group">
            <label class="dform-label">Urutan Tampil <span class="dform-opt">(opsional)</span></label>
            <input v-model.number="alurForm.urutan" type="number" class="dform-input" min="1" :disabled="alurSaving" style="width:100px"/>
          </div>
          <div v-if="alurFormError" class="dform-error">{{ alurFormError }}</div>
        </div>
        <div class="divisi-modal__footer">
          <button class="btn-cancel" @click="showAlurModal = false" :disabled="alurSaving">Batal</button>
          <button class="btn-green" @click="saveAlur" :disabled="alurSaving || !alurForm.judul.trim()">
            <div v-if="alurSaving" class="spinner-sm"></div>
            {{ alurSaving ? 'Menyimpan…' : (alurForm.id ? 'Simpan Perubahan' : 'Tambah Alur') }}
          </button>
        </div>
      </div>
    </div>
  </Teleport>

  <!-- ── MODAL TAMBAH/EDIT DIVISI ── -->
  <Teleport to="body">
    <div v-if="showDivisiModal" class="modal-overlay" @click.self="showDivisiModal = false">
      <div class="divisi-modal" style="width:min(560px,100%);max-height:90vh;overflow-y:auto">
        <div class="divisi-modal__header">
          <h3 class="divisi-modal__title">{{ divisiForm.id ? 'Edit Divisi' : 'Tambah Divisi' }}</h3>
          <button class="modal-close" @click="showDivisiModal = false">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none"><path d="M18 6L6 18M6 6l12 12" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"/></svg>
          </button>
        </div>
        <div class="divisi-modal__body" style="gap:14px">
          <!-- Nama & Urutan -->
          <div class="dform-group">
            <label class="dform-label">Nama Divisi <span class="dform-req">*</span></label>
            <input v-model="divisiForm.nama" type="text" class="dform-input" placeholder="contoh: IT / Sistem Informasi" :disabled="divisiSaving"/>
          </div>
          <div class="dform-group">
            <label class="dform-label">Urutan Tampil <span class="dform-opt">(opsional)</span></label>
            <input v-model.number="divisiForm.urutan" type="number" class="dform-input" placeholder="0" min="0" :disabled="divisiSaving"/>
            <span class="dform-hint">Angka kecil tampil lebih atas di dropdown</span>
          </div>

          <!-- Geofencing section -->
          <div class="geo-section">
            <label class="geo-toggle-row" @click="divisiForm.geoEnabled = !divisiForm.geoEnabled">
              <div class="geo-toggle-info">
                <div class="geo-toggle-title">
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none"><path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7z" stroke="currentColor" stroke-width="2" stroke-linejoin="round"/><circle cx="12" cy="9" r="2.5" stroke="currentColor" stroke-width="2"/></svg>
                  Geofencing Absensi Divisi
                </div>
                <div class="geo-toggle-sub">Jika diaktifkan, peserta divisi ini hanya bisa absen dari lokasi ini</div>
              </div>
              <div :class="divisiForm.geoEnabled ? 'geo-sw geo-sw--on' : 'geo-sw'">
                <div class="geo-sw__knob"></div>
              </div>
            </label>

            <div v-if="divisiForm.geoEnabled" class="geo-fields">
              <div class="dform-group">
                <label class="dform-label">Nama Lokasi <span class="dform-opt">(opsional)</span></label>
                <input v-model="divisiForm.namaLokasi" type="text" class="dform-input" placeholder="contoh: Gedung Produksi Lantai 2" :disabled="divisiSaving"/>
                <span class="dform-hint">Ditampilkan ke peserta saat absensi ditolak</span>
              </div>

              <div class="geo-coords-row">
                <div class="dform-group" style="flex:1">
                  <label class="dform-label">Latitude <span class="dform-req">*</span></label>
                  <input v-model.number="divisiForm.geoLat" type="number" step="0.000001" class="dform-input" placeholder="-3.432194" :disabled="divisiSaving"/>
                </div>
                <div class="dform-group" style="flex:1">
                  <label class="dform-label">Longitude <span class="dform-req">*</span></label>
                  <input v-model.number="divisiForm.geoLng" type="number" step="0.000001" class="dform-input" placeholder="104.035361" :disabled="divisiSaving"/>
                </div>
              </div>

              <div class="dform-group">
                <label class="dform-label">Radius Absensi (meter) <span class="dform-req">*</span></label>
                <input v-model.number="divisiForm.geoRadius" type="number" min="50" max="5000" class="dform-input" placeholder="500" :disabled="divisiSaving"/>
                <span class="dform-hint">Peserta harus berada dalam radius ini dari titik koordinat di atas</span>
              </div>

              <button type="button" class="btn-use-location" @click="useMyLocation" :disabled="gettingLocation">
                <svg width="13" height="13" viewBox="0 0 24 24" fill="none"><circle cx="12" cy="12" r="3" stroke="currentColor" stroke-width="2"/><path d="M12 2v3M12 19v3M2 12h3M19 12h3" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>
                {{ gettingLocation ? 'Mengambil lokasi…' : 'Gunakan Lokasi Saya Sekarang' }}
              </button>
              <div v-if="locationError" class="dform-error">{{ locationError }}</div>
            </div>
          </div>

          <div v-if="divisiFormError" class="dform-error">{{ divisiFormError }}</div>
        </div>
        <div class="divisi-modal__footer">
          <button class="btn-cancel" @click="showDivisiModal = false" :disabled="divisiSaving">Batal</button>
          <button class="btn-green" @click="saveDivisi" :disabled="divisiSaving || !divisiForm.nama.trim()">
            <div v-if="divisiSaving" class="spinner-sm"></div>
            {{ divisiSaving ? 'Menyimpan…' : (divisiForm.id ? 'Simpan Perubahan' : 'Tambah Divisi') }}
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from "vue";
import DashboardLayout from "@/layouts/DashboardLayout.vue";
import AdminUsers      from "./AdminUsers.vue";
import AdminDokumen    from "./AdminDokumen.vue";
import AdminHero       from "./AdminHero.vue";
import AdminJadwal     from "./AdminJadwal.vue";
import AdminKnowledge  from "./AdminKnowledge.vue";
import StaffProfil     from "@/views/staff/StaffProfil.vue";
import { useAuth } from "@/hooks/useAuth";
import api from "@/lib/api";

const { user } = useAuth();
const layout   = ref<InstanceType<typeof DashboardLayout> | null>(null);
const activeTab = computed(() => layout.value?.activeTab ?? "beranda");
const firstName = computed(() => user.value?.nama_lengkap?.split(" ")[0] ?? "");


const ICON = {
  home:   `<svg width="16" height="16" viewBox="0 0 24 24" fill="none"><path d="M3 9l9-7 9 7v11a2 2 0 01-2 2H5a2 2 0 01-2-2z" stroke="currentColor" stroke-width="2"/><polyline points="9 22 9 12 15 12 15 22" stroke="currentColor" stroke-width="2"/></svg>`,
  users:  `<svg width="16" height="16" viewBox="0 0 24 24" fill="none"><path d="M17 21v-2a4 4 0 00-4-4H5a4 4 0 00-4 4v2" stroke="currentColor" stroke-width="2"/><circle cx="9" cy="7" r="4" stroke="currentColor" stroke-width="2"/></svg>`,
  period: `<svg width="16" height="16" viewBox="0 0 24 24" fill="none"><rect x="3" y="4" width="18" height="18" rx="2" ry="2" stroke="currentColor" stroke-width="2"/><line x1="3" y1="10" x2="21" y2="10" stroke="currentColor" stroke-width="2"/></svg>`,
  globe:  `<svg width="16" height="16" viewBox="0 0 24 24" fill="none"><circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2"/><line x1="2" y1="12" x2="22" y2="12" stroke="currentColor" stroke-width="2"/><path d="M12 2a15.3 15.3 0 014 10 15.3 15.3 0 01-4 10 15.3 15.3 0 01-4-10 15.3 15.3 0 014-10z" stroke="currentColor" stroke-width="2"/></svg>`,
  faq:    `<svg width="16" height="16" viewBox="0 0 24 24" fill="none"><circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2"/><path d="M9.09 9a3 3 0 015.83 1c0 2-3 3-3 3" stroke="currentColor" stroke-width="2" stroke-linecap="round"/><line x1="12" y1="17" x2="12.01" y2="17" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>`,
  stats:  `<svg width="16" height="16" viewBox="0 0 24 24" fill="none"><line x1="18" y1="20" x2="18" y2="10" stroke="currentColor" stroke-width="2" stroke-linecap="round"/><line x1="12" y1="20" x2="12" y2="4" stroke="currentColor" stroke-width="2" stroke-linecap="round"/><line x1="6" y1="20" x2="6" y2="14" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>`,
  clock:  `<svg width="16" height="16" viewBox="0 0 24 24" fill="none"><circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2"/><polyline points="12 6 12 12 16 14" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>`,
};

const ICON_DIVISI   = `<svg width="16" height="16" viewBox="0 0 24 24" fill="none"><rect x="2" y="7" width="20" height="14" rx="2" stroke="currentColor" stroke-width="2"/><path d="M16 7V5a2 2 0 00-2-2h-4a2 2 0 00-2 2v2" stroke="currentColor" stroke-width="2"/></svg>`;
const ICON_DOKUMEN  = `<svg width="16" height="16" viewBox="0 0 24 24" fill="none"><path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z" stroke="currentColor" stroke-width="2"/><polyline points="14 2 14 8 20 8" stroke="currentColor" stroke-width="2"/><line x1="8" y1="13" x2="16" y2="13" stroke="currentColor" stroke-width="2"/><line x1="8" y1="17" x2="16" y2="17" stroke="currentColor" stroke-width="2"/></svg>`;

const navGroups = [
  { items: [{ key: "beranda", label: "Beranda", icon: ICON.home }] },
  {
    label: "Manajemen",
    items: [
      { key: "users",    label: "Manajemen User", icon: ICON.users },
      { key: "divisi",   label: "Kelola Divisi",  icon: ICON_DIVISI },
      { key: "jamabsen", label: "Jam Absensi",     icon: ICON.clock },
    ],
  },
  {
    label: "Konten Web",
    items: [
      { key: "hero",    label: "Hero Section",     icon: `<svg width="16" height="16" viewBox="0 0 24 24" fill="none"><rect x="3" y="3" width="18" height="18" rx="3" stroke="currentColor" stroke-width="2"/><path d="M3 15l5-5 4 4 3-3 4 4" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/><circle cx="8.5" cy="8.5" r="1.5" fill="currentColor"/></svg>` },
      { key: "jadwal",  label: "Jadwal & Periode", icon: `<svg width="16" height="16" viewBox="0 0 24 24" fill="none"><path d="M8 6h13M8 12h13M8 18h13M3 6h.01M3 12h.01M3 18h.01" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>` },
      { key: "landing", label: "Alur Pendaftaran", icon: ICON.globe },
      { key: "faq",     label: "Kelola FAQ",       icon: ICON.faq },
    ],
  },
  {
    label: "Helpdesk",
    items: [
      { key: "knowledge", label: "Knowledge Base", icon: `<svg width="16" height="16" viewBox="0 0 24 24" fill="none"><path d="M12 2a10 10 0 100 20A10 10 0 0012 2z" stroke="currentColor" stroke-width="2"/><path d="M9.09 9a3 3 0 015.83 1c0 2-3 3-3 3" stroke="currentColor" stroke-width="2" stroke-linecap="round"/><line x1="12" y1="17" x2="12.01" y2="17" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>` },
    ],
  },
  {
    label: "Laporan",
    items: [
      { key: "dokumen", label: "Kelola Dokumen", icon: ICON_DOKUMEN },
    ],
  },
  {
    label: "Akun",
    items: [
      { key: "profil", label: "Profil Saya", icon: `<svg width="16" height="16" viewBox="0 0 24 24" fill="none"><path d="M20 21v-2a4 4 0 00-4-4H8a4 4 0 00-4 4v2" stroke="currentColor" stroke-width="2"/><circle cx="12" cy="7" r="4" stroke="currentColor" stroke-width="2"/></svg>` },
    ],
  },
];

// ── Jam Absensi Config ────────────────────────────────────────
const jamForm  = ref({
  jam_masuk_buka: '07:30', jam_masuk_tutup: '08:00',
  jam_pulang_buka: '15:00', jam_pulang_tutup: '16:00',
});
const cfgLoading = ref(false);
const cfgSaving  = ref(false);
const cfgError   = ref('');
const cfgSuccess = ref('');

async function fetchJam() {
  cfgLoading.value = true;
  try {
    const r = await api.get('/api/absensi/config');
    const d = r.data?.data;
    if (d) jamForm.value = {
      jam_masuk_buka: d.jam_masuk_buka, jam_masuk_tutup: d.jam_masuk_tutup,
      jam_pulang_buka: d.jam_pulang_buka, jam_pulang_tutup: d.jam_pulang_tutup,
    };
  } catch { /* gunakan default */ } finally { cfgLoading.value = false; }
}

function autoFormatTime(field: 'jam_masuk_buka' | 'jam_masuk_tutup' | 'jam_pulang_buka' | 'jam_pulang_tutup', e: Event) {
  let v = (e.target as HTMLInputElement).value.replace(/[^0-9]/g, '');
  if (v.length > 4) v = v.slice(0, 4);
  if (v.length >= 3) v = v.slice(0, 2) + ':' + v.slice(2);
  jamForm.value[field] = v;
  (e.target as HTMLInputElement).value = v;
}

async function saveJam() {
  cfgSaving.value = true; cfgError.value = ''; cfgSuccess.value = '';
  const timeRe = /^([01]\d|2[0-3]):[0-5]\d$/;
  const timeFields = ['jam_masuk_buka', 'jam_masuk_tutup', 'jam_pulang_buka', 'jam_pulang_tutup'] as const;
  if (timeFields.some(f => !timeRe.test(jamForm.value[f]))) {
    cfgError.value = 'Format jam tidak valid. Gunakan HH:MM (contoh: 07:30, 16:00)';
    cfgSaving.value = false;
    return;
  }
  try {
    await api.put('/api/admin/absensi/config', jamForm.value);
    cfgSuccess.value = 'Pengaturan berhasil disimpan.';
    setTimeout(() => { cfgSuccess.value = ''; }, 3000);
  } catch (e: any) {
    cfgError.value = e.response?.data?.message ?? 'Gagal menyimpan';
  } finally { cfgSaving.value = false; }
}

// ── Stats Beranda ───────────────────────────────────────────────
const statsLoading = ref(false);
const stats = ref({
  total_peserta: 0,
  pengajuan_menunggu: 0,
  magang_aktif: 0,
  tiket_chat_menunggu: 0,
});

async function fetchStats() {
  statsLoading.value = true;
  try {
    const r = await api.get('/api/admin/stats');
    const d = r.data;
    stats.value = {
      total_peserta:       d.total_peserta       ?? 0,
      pengajuan_menunggu:  d.pengajuan_menunggu  ?? 0,
      magang_aktif:        d.magang_aktif        ?? 0,
      tiket_chat_menunggu: d.tiket_chat_menunggu ?? 0,
    };
  } catch { /* ignore */ }
  finally { statsLoading.value = false; }
}

// ── Akun Jatuh Tempo ────────────────────────────────────────────
interface AkunJatuhTempo {
  user_id: string;
  nama_lengkap: string;
  email: string;
  asal_institusi: string;
  divisi: string | null;
  tanggal_selesai: string;
  sisa_hari: number;
  hapus_at: string;
}

const jtList    = ref<AkunJatuhTempo[]>([]);
const jtLoading = ref(false);

async function fetchJatuhTempo() {
  jtLoading.value = true;
  try {
    const r = await api.get('/api/admin/cleanup/preview?days=7');
    jtList.value = Array.isArray(r.data?.data) ? r.data.data : [];
  } catch { jtList.value = []; }
  finally { jtLoading.value = false; }
}

onMounted(() => { fetchStats(); fetchJatuhTempo(); });

watch(activeTab, (tab) => {
  if (tab === 'jamabsen') fetchJam();
  if (tab === 'divisi') fetchDivisi();
  if (tab === 'landing') fetchAlur();
  if (tab === 'faq') fetchFaqAdmin();
});

// ── Kelola FAQ ─────────────────────────────────────────────────
interface FaqItem { id?: string; pertanyaan: string; jawaban: string; urutan: number; is_active?: boolean; }

const faqList        = ref<FaqItem[]>([]);
const faqLoading     = ref(false);
const faqFormVisible = ref(false);
const faqSaving      = ref(false);
const faqFormError   = ref('');
const faqEditTarget  = ref<FaqItem | null>(null);
const faqFormData    = ref<FaqItem>({ pertanyaan: '', jawaban: '', urutan: 1, is_active: true });

async function fetchFaqAdmin() {
  faqLoading.value = true;
  try {
    const r = await api.get('/api/landing/faq/all');
    faqList.value = Array.isArray(r.data?.data) ? r.data.data : [];
  } catch { faqList.value = []; }
  finally { faqLoading.value = false; }
}

function openFaqForm(item: FaqItem | null) {
  faqFormError.value = '';
  faqEditTarget.value = item;
  faqFormData.value = item
    ? { ...item }
    : { pertanyaan: '', jawaban: '', urutan: faqList.value.length + 1, is_active: true };
  faqFormVisible.value = true;
}

function closeFaqForm() {
  faqFormVisible.value = false;
  faqEditTarget.value = null;
  faqFormError.value = '';
}

async function saveFaqItem() {
  if (!faqFormData.value.pertanyaan.trim() || !faqFormData.value.jawaban.trim()) {
    faqFormError.value = 'Pertanyaan dan jawaban wajib diisi'; return;
  }
  faqSaving.value = true; faqFormError.value = '';
  try {
    await api.post('/api/landing/faq', faqFormData.value);
    closeFaqForm();
    await fetchFaqAdmin();
  } catch (e: any) {
    faqFormError.value = e.response?.data?.message ?? 'Gagal menyimpan FAQ';
  } finally { faqSaving.value = false; }
}

async function deleteFaqItem(id: string) {
  if (!confirm('Hapus FAQ ini dari landing page?')) return;
  try { await api.delete(`/api/landing/faq/${id}`); await fetchFaqAdmin(); } catch { /* ignore */ }
}

async function toggleFaqItem(item: FaqItem) {
  try {
    await api.post('/api/landing/faq', { ...item, is_active: !item.is_active });
    await fetchFaqAdmin();
  } catch { /* ignore */ }
}

// ── Alur landing page ──────────────────────────────────────────
interface AlurItem { id: string; judul: string; paragraf: string; gambar_url: string; urutan: number; }
const alurList      = ref<AlurItem[]>([]);
const alurLoading   = ref(false);
const showAlurModal = ref(false);
const alurSaving    = ref(false);
const alurFormError = ref('');
const alurForm      = ref<{ id: string | null; judul: string; paragraf: string; gambar_url: string; urutan: number }>({
  id: null, judul: '', paragraf: '', gambar_url: '', urutan: 1,
});
const alurFileInput   = ref<HTMLInputElement | null>(null);
const alurUploadPreview = ref('');
const alurUploading   = ref(false);
const alurUploadError = ref('');

async function fetchAlur() {
  alurLoading.value = true;
  try {
    const r = await api.get('/api/landing/alur');
    alurList.value = Array.isArray(r.data?.data) ? r.data.data : [];
  } catch { alurList.value = []; }
  finally { alurLoading.value = false; }
}

function openAlurModal(item: AlurItem | null) {
  alurFormError.value = '';
  alurUploadPreview.value = '';
  alurUploadError.value = '';
  if (item) {
    alurForm.value = { id: item.id, judul: item.judul, paragraf: item.paragraf, gambar_url: item.gambar_url, urutan: item.urutan };
  } else {
    alurForm.value = { id: null, judul: '', paragraf: '', gambar_url: '', urutan: alurList.value.length + 1 };
  }
  showAlurModal.value = true;
}

function triggerFileInput() {
  if (!alurUploading.value) alurFileInput.value?.click();
}

async function onAlurFileChange(e: Event) {
  const file = (e.target as HTMLInputElement).files?.[0];
  if (!file) return;
  alurUploadError.value = '';
  alurUploadPreview.value = URL.createObjectURL(file);
  alurUploading.value = true;
  try {
    const fd = new FormData();
    fd.append('gambar', file);
    const r = await api.post('/api/admin/alur/upload', fd, { headers: { 'Content-Type': 'multipart/form-data' } });
    alurForm.value.gambar_url = r.data?.data?.url ?? '';
  } catch (err: any) {
    alurUploadError.value = err.response?.data?.message ?? 'Gagal upload gambar';
    alurUploadPreview.value = '';
  } finally {
    alurUploading.value = false;
    (e.target as HTMLInputElement).value = '';
  }
}

async function saveAlur() {
  if (!alurForm.value.judul.trim()) { alurFormError.value = 'Judul wajib diisi'; return; }
  alurSaving.value = true; alurFormError.value = '';
  try {
    const payload = {
      judul: alurForm.value.judul.trim(),
      paragraf: alurForm.value.paragraf.trim(),
      gambar_url: alurForm.value.gambar_url.trim(),
      urutan: alurForm.value.urutan,
    };
    if (alurForm.value.id) {
      await api.put(`/api/admin/alur/${alurForm.value.id}`, payload);
    } else {
      await api.post('/api/admin/alur', payload);
    }
    showAlurModal.value = false;
    await fetchAlur();
  } catch (e: any) {
    alurFormError.value = e.response?.data?.message ?? 'Gagal menyimpan item alur';
  } finally { alurSaving.value = false; }
}

async function deleteAlur(item: AlurItem) {
  if (!confirm(`Hapus alur "${item.judul}"?`)) return;
  try {
    await api.delete(`/api/admin/alur/${item.id}`);
    await fetchAlur();
  } catch { /* ignore */ }
}

// ── Kelola Divisi ─────────────────────────────────────────────
interface Divisi {
  id: string; nama: string; is_active: boolean; urutan: number; created_at: string;
  geo_lat?: number | null; geo_lng?: number | null; geo_radius?: number | null; nama_lokasi?: string | null;
}

interface DivisiForm {
  id: string | null; nama: string; urutan: number;
  geoEnabled: boolean; namaLokasi: string;
  geoLat: number | null; geoLng: number | null; geoRadius: number | null;
}

const divisiList      = ref<Divisi[]>([]);
const divisiLoading   = ref(false);
const showDivisiModal = ref(false);
const divisiSaving    = ref(false);
const divisiFormError = ref('');
const gettingLocation = ref(false);
const locationError   = ref('');
const divisiForm      = ref<DivisiForm>({ id: null, nama: '', urutan: 0, geoEnabled: false, namaLokasi: '', geoLat: null, geoLng: null, geoRadius: 500 });

async function fetchDivisi() {
  divisiLoading.value = true;
  try {
    const r = await api.get('/api/admin/divisi');
    divisiList.value = Array.isArray(r.data?.data) ? r.data.data : [];
  } catch { divisiList.value = []; }
  finally { divisiLoading.value = false; }
}

function openDivisiModal(d: Divisi | null) {
  divisiFormError.value = '';
  locationError.value = '';
  if (d) {
    divisiForm.value = {
      id: d.id, nama: d.nama, urutan: d.urutan,
      geoEnabled: d.geo_lat != null,
      namaLokasi: d.nama_lokasi ?? '',
      geoLat: d.geo_lat ?? null,
      geoLng: d.geo_lng ?? null,
      geoRadius: d.geo_radius ?? 500,
    };
  } else {
    divisiForm.value = { id: null, nama: '', urutan: divisiList.value.length + 1, geoEnabled: false, namaLokasi: '', geoLat: null, geoLng: null, geoRadius: 500 };
  }
  showDivisiModal.value = true;
}

async function saveDivisi() {
  if (!divisiForm.value.nama.trim()) return;
  divisiFormError.value = '';

  const f = divisiForm.value;

  // Validasi geofencing jika diaktifkan
  if (f.geoEnabled) {
    if (f.geoLat === null || f.geoLng === null) {
      divisiFormError.value = 'Latitude dan Longitude wajib diisi jika geofencing diaktifkan';
      return;
    }
    if (!f.geoRadius || f.geoRadius < 50) {
      divisiFormError.value = 'Radius minimal 50 meter';
      return;
    }
  }

  const payload: Record<string, unknown> = {
    nama: f.nama.trim(),
    urutan: f.urutan,
    geo_lat: f.geoEnabled ? f.geoLat : null,
    geo_lng: f.geoEnabled ? f.geoLng : null,
    geo_radius: f.geoEnabled ? f.geoRadius : null,
    nama_lokasi: f.geoEnabled && f.namaLokasi.trim() ? f.namaLokasi.trim() : null,
  };

  divisiSaving.value = true;
  try {
    if (f.id) {
      await api.patch(`/api/admin/divisi/${f.id}`, payload);
    } else {
      await api.post('/api/admin/divisi', payload);
    }
    showDivisiModal.value = false;
    await fetchDivisi();
  } catch (e: any) {
    divisiFormError.value = e.response?.data?.message ?? 'Gagal menyimpan divisi';
  } finally { divisiSaving.value = false; }
}

function useMyLocation() {
  locationError.value = '';
  if (!navigator.geolocation) {
    locationError.value = 'Browser tidak mendukung GPS. Masukkan koordinat secara manual.';
    return;
  }
  gettingLocation.value = true;
  navigator.geolocation.getCurrentPosition(
    (pos) => {
      divisiForm.value.geoLat = parseFloat(pos.coords.latitude.toFixed(8));
      divisiForm.value.geoLng = parseFloat(pos.coords.longitude.toFixed(8));
      gettingLocation.value = false;
    },
    (err) => {
      locationError.value = 'Gagal mengambil lokasi: ' + (err.message || 'izin ditolak');
      gettingLocation.value = false;
    },
    { enableHighAccuracy: true, timeout: 10000 }
  );
}

async function toggleDivisi(d: Divisi) {
  try {
    await api.patch(`/api/admin/divisi/${d.id}/toggle`, { is_active: !d.is_active });
    await fetchDivisi();
  } catch { /* ignore */ }
}

async function deleteDivisi(d: Divisi) {
  if (!confirm(`Hapus divisi "${d.nama}"? Peserta yang sudah ditempatkan tidak akan terpengaruh.`)) return;
  try {
    await api.delete(`/api/admin/divisi/${d.id}`);
    await fetchDivisi();
  } catch { /* ignore */ }
}

function formatDate(s: string) {
  return new Date(s).toLocaleDateString('id-ID', { day: '2-digit', month: 'short', year: 'numeric' });
}
</script>

<style scoped>
.welcome-banner { background: linear-gradient(135deg, #0d2818 0%, #1a5c20 100%); border-radius: 14px; padding: 24px 28px; color: #fff; display: flex; align-items: center; justify-content: space-between; gap: 16px; }
.welcome-banner__greeting { font-size: 17px; font-weight: 700; }
.welcome-banner__sub      { font-size: 12.5px; color: rgba(255,255,255,0.65); margin-top: 4px; }
.welcome-banner__badge    { background: rgba(72,175,74,0.25); border: 1px solid rgba(134,239,172,0.4); color: #86efac; font-size: 11px; font-weight: 800; letter-spacing: 0.12em; padding: 6px 14px; border-radius: 100px; }

.stats-row { display: grid; grid-template-columns: repeat(4,1fr); gap: 14px; }
.stat-card { background: #fff; border-radius: 12px; padding: 16px; display: flex; align-items: center; gap: 12px; border: 1px solid #e9f5e9; box-shadow: 0 1px 3px rgba(13,40,24,0.05); }
.stat-card__icon  { width: 38px; height: 38px; border-radius: 10px; display: flex; align-items: center; justify-content: center; flex-shrink: 0; }
.stat-card__label { font-size: 11px; color: #6b7280; font-weight: 500; }
.stat-card__value { font-size: 20px; font-weight: 700; color: #111827; margin-top: 1px; }
.stat-card__sub   { font-size: 10.5px; color: #9ca3af; }

.jt-count { background: #f0fdf4; color: #16a34a; font-size: 11px; font-weight: 700; padding: 3px 12px; border-radius: 100px; border: 1px solid #d1fae5; white-space: nowrap; }

.card { background: #fff; border-radius: 14px; border: 1px solid #e9f5e9; box-shadow: 0 1px 3px rgba(13,40,24,0.05); overflow: hidden; }
.card-header { display: flex; align-items: center; justify-content: space-between; padding: 16px 20px; border-bottom: 1px solid #f0faf0; }
.card-title  { font-size: 13.5px; font-weight: 700; color: #111827; margin: 0; }

.empty-state { display: flex; flex-direction: column; align-items: center; padding: 44px 24px; gap: 12px; text-align: center; }
.empty-state__icon { width: 72px; height: 72px; background: #f9fafb; border-radius: 50%; display: flex; align-items: center; justify-content: center; }
.empty-state p { font-size: 13px; color: #9ca3af; line-height: 1.7; margin: 0; }

.btn-green    { background: #48AF4A; color: #fff; border: none; border-radius: 9px; padding: 10px 22px; font-size: 13px; font-weight: 600; font-family: "Poppins",sans-serif; cursor: pointer; }
.btn-green-sm { background: #48AF4A; color: #fff; border: none; border-radius: 8px; padding: 6px 14px; font-size: 12px; font-weight: 600; font-family: "Poppins",sans-serif; cursor: pointer; white-space: nowrap; }

.jam-divider { height: 1px; background: #f0faf0; margin: 4px 0; }
.geo-desc { font-size: 12px; color: #6b7280; margin-bottom: 14px; line-height: 1.6; }
.geo-hint { font-size: 11.5px; color: #9ca3af; margin-top: 12px; line-height: 1.6; }
.jam-input--wide { width: 140px; }
.jam-input--sm   { width: 100px; }

@media (max-width: 700px) { .stats-row, .quick-grid { grid-template-columns: 1fr 1fr; } }
@media (max-width: 420px) { .stats-row, .quick-grid { grid-template-columns: 1fr; } }

/* ── Alur admin list ── */
.btn-preview-sm {
  display: inline-flex; align-items: center; gap: 5px;
  padding: 5px 12px; font-size: 12px; font-weight: 600;
  border: 1.5px solid #48AF4A; color: #48AF4A; border-radius: 8px;
  text-decoration: none; font-family: "Poppins", sans-serif;
  transition: background 0.15s;
}
.btn-preview-sm:hover { background: #f0fdf4; }

.alur-admin-list { display: flex; flex-direction: column; }
.alur-admin-item {
  display: flex; align-items: center; gap: 14px;
  padding: 14px 20px; border-bottom: 1px solid #f9fafb;
  transition: background 0.1s;
}
.alur-admin-item:last-child { border-bottom: none; }
.alur-admin-item:hover { background: #fafffe; }

.alur-admin-item__num {
  width: 28px; height: 28px; flex-shrink: 0; border-radius: 50%;
  background: #e8f5e9; color: #16a34a; font-size: 12px; font-weight: 800;
  display: flex; align-items: center; justify-content: center;
}
.alur-admin-item__img {
  width: 52px; height: 52px; flex-shrink: 0; border-radius: 10px; overflow: hidden;
  background: #f3f4f6; display: flex; align-items: center; justify-content: center;
}
.alur-admin-item__img img { width: 100%; height: 100%; object-fit: cover; }
.alur-admin-item__img--empty { background: #f9fafb; }

.alur-admin-item__content { flex: 1; min-width: 0; }
.alur-admin-item__judul {
  font-size: 13.5px; font-weight: 700; color: #111827;
  white-space: nowrap; overflow: hidden; text-overflow: ellipsis;
}
.alur-admin-item__paragraf {
  font-size: 12px; color: #6b7280; margin-top: 3px; line-height: 1.5;
  display: -webkit-box; -webkit-line-clamp: 2; -webkit-box-orient: vertical; overflow: hidden;
}
.alur-admin-item__actions { display: flex; gap: 6px; flex-shrink: 0; }

/* ── Jam Absen Form ───────────────────────── */
.badge-info { font-size: 11px; font-weight: 600; color: #0d2818; background: #f0fdf4; border: 1px solid #bbf7d0; border-radius: 100px; padding: 4px 12px; }
.jam-form { padding: 24px; display: flex; flex-direction: column; gap: 24px; }
.jam-section__label { font-size: 12px; font-weight: 700; color: #374151; text-transform: uppercase; letter-spacing: .06em; margin-bottom: 12px; }
.jam-row { display: flex; align-items: flex-end; gap: 12px; }
.jam-field { display: flex; flex-direction: column; gap: 4px; }
.jam-field label { font-size: 11.5px; font-weight: 600; color: #6b7280; }
.jam-input { border: 1.5px solid #e5e7eb; border-radius: 9px; padding: 9px 12px; font-size: 14px; font-family: "Poppins",sans-serif; outline: none; transition: border-color .15s; color: #111827; background: #fff; }
.jam-input:focus { border-color: #48AF4A; }
.jam-sep { font-size: 16px; font-weight: 600; color: #9ca3af; padding-bottom: 8px; }
.jam-actions { display: flex; justify-content: flex-end; padding-top: 4px; }
.cfg-error   { background: #fff1f2; border: 1px solid #fecdd3; color: #be123c; font-size: 12.5px; padding: 8px 14px; border-radius: 8px; }
.cfg-success { background: #f0fdf4; border: 1px solid #bbf7d0; color: #16a34a; font-size: 12.5px; padding: 8px 14px; border-radius: 8px; }

/* ── Alur image upload area ── */
.alur-upload-area {
  position: relative; border: 2px dashed #d1d5db; border-radius: 12px;
  min-height: 140px; display: flex; flex-direction: column;
  align-items: center; justify-content: center; gap: 6px;
  cursor: pointer; transition: border-color 0.15s, background 0.15s;
  overflow: hidden; background: #fafafa;
}
.alur-upload-area:hover { border-color: #48AF4A; background: #f0fdf4; }
.alur-upload-area--uploading { cursor: wait; pointer-events: none; }
.alur-upload-icon { color: #94a3b8; }
.alur-upload-text { font-size: 13px; font-weight: 600; color: #374151; }
.alur-upload-hint { font-size: 11px; color: #9ca3af; }
.alur-upload-preview {
  width: 100%; max-height: 200px; object-fit: cover;
  border-radius: 10px; display: block;
}
.alur-upload-change {
  position: absolute; bottom: 0; left: 0; right: 0;
  background: rgba(0,0,0,0.5); color: #fff; font-size: 11px;
  font-weight: 600; text-align: center; padding: 6px; letter-spacing: 0.03em;
}
.alur-upload-overlay {
  position: absolute; inset: 0; background: rgba(255,255,255,0.75);
  display: flex; align-items: center; justify-content: center;
}

/* Spinner (shared) */
.spinner { width: 36px; height: 36px; border: 3px solid #e5e7eb; border-top-color: #48AF4A; border-radius: 50%; animation: spin .8s linear infinite; margin: 40px auto; }
.spinner-sm { width: 14px; height: 14px; border: 2px solid rgba(255,255,255,0.35); border-top-color: #fff; border-radius: 50%; animation: spin .7s linear infinite; display: inline-block; }
@keyframes spin { to { transform: rotate(360deg); } }

/* ── Divisi table ── */
.card-sub { font-size: 11.5px; color: #9ca3af; margin: 2px 0 0; }
.table-wrap { overflow-x: auto; }
.data-table { width: 100%; border-collapse: collapse; font-size: 13px; }
.data-table th { padding: 11px 16px; text-align: left; font-size: 10.5px; font-weight: 600; color: #6b7280; background: #f9fafb; border-bottom: 1px solid #f1f5f9; text-transform: uppercase; letter-spacing: .04em; white-space: nowrap; }
.data-table td { padding: 13px 16px; border-bottom: 1px solid #f9fafb; color: #374151; vertical-align: middle; }
.td-nama { font-weight: 600; color: #111827; }
.urutan-badge { background: #f3f4f6; color: #374151; font-size: 11px; font-weight: 700; padding: 2px 8px; border-radius: 6px; }
.status-pill { display: inline-flex; align-items: center; font-size: 11px; font-weight: 700; padding: 3px 10px; border-radius: 100px; }
.status-pill--green  { background: #dcfce7; color: #16a34a; }
.status-pill--gray   { background: #f3f4f6; color: #6b7280; }
.status-pill--red    { background: #fee2e2; color: #dc2626; }
.status-pill--yellow { background: #fef9c3; color: #a16207; }
.aksi-cell { display: flex; align-items: center; gap: 6px; flex-wrap: wrap; }
.btn-aksi { border: none; border-radius: 7px; padding: 5px 11px; font-size: 11.5px; font-weight: 700; font-family: inherit; cursor: pointer; white-space: nowrap; transition: opacity .15s; }
.btn-aksi--ghost { background: #f3f4f6; color: #374151; }
.btn-aksi--ghost:hover { background: #e5e7eb; }
.btn-aksi--green { background: #dcfce7; color: #16a34a; }
.btn-aksi--green:hover { background: #bbf7d0; }
.btn-aksi--warn  { background: #fef9c3; color: #92400e; }
.btn-aksi--warn:hover { background: #fde68a; }
.btn-aksi--red   { background: #fee2e2; color: #dc2626; }
.btn-aksi--red:hover { background: #fecaca; }

/* ── Divisi modal ── */
.modal-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.5); backdrop-filter: blur(3px); z-index: 300; display: flex; align-items: center; justify-content: center; padding: 20px; }
.divisi-modal { background: #fff; border-radius: 18px; width: min(440px, 100%); box-shadow: 0 24px 80px rgba(0,0,0,0.22); display: flex; flex-direction: column; overflow: hidden; }
.divisi-modal__header { display: flex; align-items: center; justify-content: space-between; padding: 20px 24px 14px; border-bottom: 1px solid #f0faf0; }
.divisi-modal__title  { font-size: 15px; font-weight: 700; color: #111827; margin: 0; }
.divisi-modal__body   { padding: 20px 24px; display: flex; flex-direction: column; gap: 16px; }
.divisi-modal__footer { display: flex; gap: 10px; padding: 14px 24px; border-top: 1px solid #f0faf0; }
.modal-close { background: #f3f4f6; border: none; border-radius: 8px; width: 30px; height: 30px; display: flex; align-items: center; justify-content: center; cursor: pointer; color: #6b7280; }
.modal-close:hover { background: #e5e7eb; }
.dform-group { display: flex; flex-direction: column; gap: 5px; }
.dform-label { font-size: 12px; font-weight: 600; color: #374151; }
.dform-req   { color: #dc2626; }
.dform-opt   { color: #9ca3af; font-weight: 400; }
.dform-input { border: 1px solid #e5e7eb; border-radius: 9px; padding: 9px 12px; font-size: 13px; font-family: inherit; color: #111827; outline: none; transition: border-color .15s; }
.dform-input:focus { border-color: #48AF4A; box-shadow: 0 0 0 3px rgba(72,175,74,.12); }
.dform-input:disabled { background: #f9fafb; color: #9ca3af; cursor: not-allowed; }
.dform-hint  { font-size: 11px; color: #9ca3af; }
.dform-error { font-size: 12.5px; color: #dc2626; background: #fef2f2; border: 1px solid #fecaca; border-radius: 8px; padding: 9px 13px; }
.btn-cancel { flex: 1; background: #f3f4f6; color: #374151; border: none; border-radius: 9px; padding: 10px 16px; font-size: 13px; font-weight: 600; font-family: inherit; cursor: pointer; }
.btn-cancel:hover:not(:disabled) { background: #e5e7eb; }
.btn-cancel:disabled { opacity: .5; cursor: default; }
.btn-green  { flex: 1; background: #48AF4A; color: #fff; border: none; border-radius: 9px; padding: 10px 22px; font-size: 13px; font-weight: 600; font-family: "Poppins",sans-serif; cursor: pointer; display: flex; align-items: center; justify-content: center; gap: 7px; }
.btn-green:hover:not(:disabled) { background: #48AF4A; }
.btn-green:disabled { opacity: .5; cursor: not-allowed; }

/* ── FAQ admin ── */
.faq-form-box {
  border-bottom: 1px solid #f0faf0; padding: 20px 20px 24px;
  background: #f9fffe; display: flex; flex-direction: column; gap: 14px;
}
.faq-form-title { font-size: 13px; font-weight: 700; color: #374151; }
.faq-textarea { resize: vertical; min-height: 90px; font-family: inherit; line-height: 1.6; }
.dform-input--sm { max-width: 120px; }
.faq-form-row { display: flex; gap: 16px; }
.faq-form-actions { display: flex; align-items: center; gap: 10px; justify-content: flex-end; }

.faq-admin-list { display: flex; flex-direction: column; }
.faq-admin-item {
  display: flex; align-items: flex-start; gap: 14px;
  padding: 16px 20px; border-bottom: 1px solid #f9fafb; transition: background 0.1s;
}
.faq-admin-item:last-child { border-bottom: none; }
.faq-admin-item:hover { background: #fafffe; }
.faq-admin-item__num {
  flex-shrink: 0; width: 28px; height: 28px; margin-top: 2px; border-radius: 8px;
  background: #e8f5e9; color: #16a34a; font-size: 11px; font-weight: 800;
  display: flex; align-items: center; justify-content: center;
}
.faq-admin-item__content { flex: 1; min-width: 0; display: flex; flex-direction: column; gap: 5px; }
.faq-admin-item__q { font-size: 13.5px; font-weight: 700; color: #111827; line-height: 1.4; }
.faq-admin-item__a {
  font-size: 12px; color: #6b7280; line-height: 1.5;
  display: -webkit-box; -webkit-line-clamp: 2; -webkit-box-orient: vertical; overflow: hidden;
}
.faq-admin-item__actions { display: flex; align-items: center; gap: 8px; flex-shrink: 0; margin-top: 2px; }

/* ── Geofencing divisi ── */
.geo-badge {
  display: inline-flex; align-items: center; gap: 4px;
  font-size: 11px; font-weight: 600; padding: 3px 9px; border-radius: 100px;
  white-space: nowrap;
}
.geo-badge--set  { background: #ecfdf5; color: #065f46; }
.geo-badge--none { background: #f3f4f6; color: #9ca3af; font-style: italic; }

.geo-info-card {
  display: flex; align-items: flex-start; gap: 9px;
  background: #f0fdf4; border: 1px solid #bbf7d0; border-radius: 10px;
  padding: 12px 16px; font-size: 12px; color: #1e40af; line-height: 1.6;
}
.geo-info-card svg { flex-shrink: 0; margin-top: 2px; }
.geo-info-card strong { font-weight: 700; }
.geo-info-card em { font-style: normal; font-weight: 600; }

.geo-section { border: 1.5px solid #e5e7eb; border-radius: 12px; overflow: hidden; }
.geo-toggle-row {
  display: flex; align-items: center; justify-content: space-between; gap: 12px;
  padding: 14px 16px; cursor: pointer; user-select: none;
  background: #f9fafb; transition: background 0.15s;
}
.geo-toggle-row:hover { background: #f0fdf4; }
.geo-toggle-title { display: flex; align-items: center; gap: 6px; font-size: 13px; font-weight: 700; color: #111827; }
.geo-toggle-sub   { font-size: 11.5px; color: #6b7280; margin-top: 2px; }
.geo-toggle-info  { flex: 1; }

/* Toggle switch */
.geo-sw { width: 36px; height: 20px; background: #d1d5db; border-radius: 100px; flex-shrink: 0; position: relative; transition: background 0.2s; }
.geo-sw--on { background: #48AF4A; }
.geo-sw__knob { position: absolute; top: 2px; left: 2px; width: 16px; height: 16px; background: #fff; border-radius: 50%; transition: transform 0.2s; box-shadow: 0 1px 3px rgba(0,0,0,0.2); }
.geo-sw--on .geo-sw__knob { transform: translateX(16px); }

.geo-fields { padding: 16px; display: flex; flex-direction: column; gap: 13px; border-top: 1.5px solid #e5e7eb; background: #fff; }
.geo-coords-row { display: flex; gap: 12px; }

.btn-use-location {
  display: inline-flex; align-items: center; gap: 6px;
  background: #f0fdf4; border: 1.5px solid #86efac; color: #16a34a;
  border-radius: 8px; padding: 8px 14px; font-size: 12px; font-weight: 600;
  font-family: "Poppins", sans-serif; cursor: pointer; transition: background 0.15s;
  align-self: flex-start;
}
.btn-use-location:hover:not(:disabled) { background: #dcfce7; }
.btn-use-location:disabled { opacity: 0.5; cursor: not-allowed; }
</style>
