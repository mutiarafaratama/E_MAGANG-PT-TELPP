<template>
  <div>
    <div class="profil-hero">
      <div class="profil-hero__av">{{ profilInitials }}</div>
      <div class="profil-hero__info">
        <div class="profil-hero__name">{{ user?.nama_lengkap }}</div>
        <div class="profil-hero__chips">
          <span class="profil-chip profil-chip--role">{{ roleLabel }}</span>
          <span class="profil-chip profil-chip--ok">Aktif</span>
        </div>
      </div>
      <button class="btn-ubah-pw" @click="showPasswordModal = true">
        <svg width="13" height="13" viewBox="0 0 24 24" fill="none"><rect x="3" y="11" width="18" height="11" rx="2" stroke="currentColor" stroke-width="2"/><path d="M7 11V7a5 5 0 0110 0v4" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>
        Ubah Password
      </button>
    </div>

    <div class="profil-cols">
      <div class="profil-card">
        <div class="profil-card__head">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none"><path d="M20 21v-2a4 4 0 00-4-4H8a4 4 0 00-4 4v2" stroke="#48AF4A" stroke-width="2"/><circle cx="12" cy="7" r="4" stroke="#48AF4A" stroke-width="2"/></svg>
          Informasi Akun
        </div>
        <div v-if="profileLoading" class="profil-skeleton">
          <div class="profil-skeleton__bar"></div>
          <div class="profil-skeleton__bar profil-skeleton__bar--short"></div>
          <div class="profil-skeleton__bar"></div>
        </div>
        <div v-else class="profil-field-list">
          <div class="profil-field">
            <span class="profil-field__label">Nama Lengkap</span>
            <span class="profil-field__value">{{ profileData?.nama_lengkap ?? user?.nama_lengkap ?? '—' }}</span>
          </div>
          <div class="profil-field">
            <span class="profil-field__label">Email</span>
            <span class="profil-field__value">{{ profileData?.email ?? user?.email ?? '—' }}</span>
          </div>
          <div class="profil-field">
            <span class="profil-field__label">Role</span>
            <span class="profil-field__value">{{ roleLabel }}</span>
          </div>
          <div class="profil-field">
            <span class="profil-field__label">Status Akun</span>
            <span class="profil-status-badge profil-status-badge--ok">Aktif</span>
          </div>
          <div class="profil-field">
            <span class="profil-field__label">Bergabung Sejak</span>
            <span class="profil-field__value">{{ profileData?.created_at ? formatTanggal(profileData.created_at) : '—' }}</span>
          </div>
        </div>
      </div>

      <div class="profil-card">
        <div class="profil-card__head">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none"><rect x="3" y="11" width="18" height="11" rx="2" stroke="#48AF4A" stroke-width="2"/><path d="M7 11V7a5 5 0 0110 0v4" stroke="#48AF4A" stroke-width="2" stroke-linecap="round"/></svg>
          Keamanan Akun
        </div>
        <div class="profil-field-list">
          <div class="profil-field">
            <span class="profil-field__label">Password</span>
            <span class="profil-field__value pw-dots">••••••••</span>
          </div>
          <div class="profil-field">
            <span class="profil-field__label">Autentikasi</span>
            <span class="profil-status-badge profil-status-badge--ok">Password</span>
          </div>
        </div>
        <div class="profil-security-action">
          <button class="btn-ubah-pw-secondary" @click="showPasswordModal = true">
            <svg width="13" height="13" viewBox="0 0 24 24" fill="none"><rect x="3" y="11" width="18" height="11" rx="2" stroke="currentColor" stroke-width="2"/><path d="M7 11V7a5 5 0 0110 0v4" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>
            Ganti Password
          </button>
          <p class="profil-security-hint">Gunakan password yang kuat — min 8 karakter dengan huruf, angka, dan karakter spesial.</p>
        </div>
      </div>
    </div>
  </div>

  <Teleport to="body">
    <div v-if="showPasswordModal" class="modal-backdrop" @click.self="closePasswordModal">
      <div class="modal-box modal-pw">
        <div class="modal-pw__head">
          <div class="modal-pw__title">Ubah Password</div>
          <button class="modal-pw__close" @click="closePasswordModal">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none"><path d="M18 6L6 18M6 6l12 12" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>
          </button>
        </div>
        <div class="pw-form">
          <div class="pw-field">
            <label class="pw-label">Password Lama <span class="ap-required">*</span></label>
            <input v-model="pwLama" type="password" class="pw-input" placeholder="Masukkan password saat ini" autocomplete="current-password" />
          </div>
          <div class="pw-field">
            <label class="pw-label">Password Baru <span class="ap-required">*</span></label>
            <input v-model="pwBaru" type="password" class="pw-input" placeholder="Min 8 karakter" autocomplete="new-password" />
            <span class="pw-hint">Min 8 karakter · harus ada huruf, angka, dan karakter spesial (!@#$%^&*...)</span>
          </div>
          <div class="pw-field">
            <label class="pw-label">Konfirmasi Password Baru <span class="ap-required">*</span></label>
            <input v-model="pwKonfirmasi" type="password" class="pw-input" placeholder="Ulangi password baru" autocomplete="new-password" @keyup.enter="changePassword" />
          </div>
          <div v-if="pwError"   class="pw-alert pw-alert--err">{{ pwError }}</div>
          <div v-if="pwSuccess" class="pw-alert pw-alert--ok">{{ pwSuccess }}</div>
        </div>
        <div class="modal-box__actions">
          <button class="btn-cancel" @click="closePasswordModal">Batal</button>
          <button class="btn-confirm" @click="changePassword" :disabled="pwLoading">
            {{ pwLoading ? 'Menyimpan...' : 'Simpan Password' }}
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import { useAuth } from "@/hooks/useAuth";
import api from "@/lib/api";

const { user } = useAuth();

const profilInitials = computed(() => {
  const name = user.value?.nama_lengkap ?? "";
  return name.split(" ").map((n: string) => n[0]).slice(0, 2).join("").toUpperCase() || "?";
});

const roleLabel = computed(() => {
  const r = user.value?.role ?? "";
  if (r === "admin") return "Super Admin";
  if (r === "hrd")   return "Staff HRD";
  return r;
});

const profileData    = ref<any>(null);
const profileLoading = ref(false);

const showPasswordModal = ref(false);
const pwLama       = ref("");
const pwBaru       = ref("");
const pwKonfirmasi = ref("");
const pwLoading    = ref(false);
const pwError      = ref("");
const pwSuccess    = ref("");

function formatTanggal(iso: string) {
  if (!iso) return "—";
  return new Date(iso).toLocaleDateString("id-ID", { day: "numeric", month: "long", year: "numeric" });
}

async function fetchData() {
  profileLoading.value = true;
  try {
    const r = await api.get("/api/auth/me");
    profileData.value = r.data?.data ?? null;
  } catch { /* silent */ } finally {
    profileLoading.value = false;
  }
}

function validateNewPassword(s: string): boolean {
  if (s.length < 8) return false;
  const special = "!@#$%^&*()_+=.,><?/";
  let hasLetter = false, hasDigit = false, hasSpecial = false;
  for (const c of s) {
    if (/[a-zA-Z]/.test(c))  hasLetter  = true;
    else if (/[0-9]/.test(c)) hasDigit  = true;
    else if (special.includes(c)) hasSpecial = true;
  }
  return hasLetter && hasDigit && hasSpecial;
}

async function changePassword() {
  pwError.value = ""; pwSuccess.value = "";
  if (!pwLama.value || !pwBaru.value || !pwKonfirmasi.value) {
    pwError.value = "Semua field wajib diisi"; return;
  }
  if (pwBaru.value !== pwKonfirmasi.value) {
    pwError.value = "Konfirmasi password baru tidak cocok"; return;
  }
  if (!validateNewPassword(pwBaru.value)) {
    pwError.value = "Password baru harus min 8 karakter, mengandung huruf, angka, dan karakter spesial"; return;
  }
  pwLoading.value = true;
  try {
    await api.post("/api/auth/change-password", { old_password: pwLama.value, new_password: pwBaru.value });
    pwSuccess.value = "Password berhasil diubah!";
    pwLama.value = ""; pwBaru.value = ""; pwKonfirmasi.value = "";
    setTimeout(() => { showPasswordModal.value = false; pwSuccess.value = ""; }, 2000);
  } catch (e: any) {
    pwError.value = e.response?.data?.message ?? "Gagal mengubah password";
  } finally {
    pwLoading.value = false;
  }
}

function closePasswordModal() {
  showPasswordModal.value = false;
  pwLama.value = ""; pwBaru.value = ""; pwKonfirmasi.value = "";
  pwError.value = ""; pwSuccess.value = "";
}

onMounted(fetchData);
</script>

<style scoped>
.profil-hero {
  background: linear-gradient(135deg, #0d2818 0%, #1a5c20 100%);
  border-radius: 14px; padding: 24px 28px; margin-bottom: 0;
  display: flex; align-items: center; gap: 18px;
}
.profil-hero__av {
  width: 62px; height: 62px; border-radius: 16px; flex-shrink: 0;
  background: rgba(134,239,172,0.2); border: 2px solid rgba(134,239,172,0.4);
  color: #86efac; font-size: 22px; font-weight: 800;
  display: flex; align-items: center; justify-content: center;
}
.profil-hero__info { flex: 1; min-width: 0; }
.profil-hero__name { font-size: 17px; font-weight: 700; color: #fff; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.profil-hero__chips { display: flex; gap: 8px; flex-wrap: wrap; margin-top: 8px; }
.profil-chip { font-size: 11px; font-weight: 600; padding: 3px 10px; border-radius: 100px; }
.profil-chip--role { background: rgba(255,255,255,0.12); color: #d1fae5; }
.profil-chip--ok   { background: rgba(134,239,172,0.2); color: #86efac; border: 1px solid rgba(134,239,172,0.4); }
.btn-ubah-pw {
  display: flex; align-items: center; gap: 6px; flex-shrink: 0;
  background: rgba(255,255,255,0.1); border: 1px solid rgba(255,255,255,0.2);
  color: #fff; border-radius: 9px; padding: 8px 14px;
  font-size: 12px; font-weight: 600; font-family: "Poppins", sans-serif; cursor: pointer;
}
.btn-ubah-pw:hover { background: rgba(255,255,255,0.18); }

.profil-cols { display: grid; grid-template-columns: 1fr 1fr; gap: 14px; margin-top: 14px; }
@media (max-width: 640px) { .profil-cols { grid-template-columns: 1fr; } }

.profil-card { background: #fff; border-radius: 14px; border: 1px solid #e9f5e9; box-shadow: 0 1px 3px rgba(13,40,24,0.05); overflow: hidden; }
.profil-card__head { display: flex; align-items: center; gap: 8px; padding: 14px 18px; font-size: 12.5px; font-weight: 700; color: #374151; border-bottom: 1px solid #f0faf0; }

.profil-field-list { padding: 4px 0 8px; }
.profil-field { display: flex; align-items: baseline; justify-content: space-between; gap: 12px; padding: 9px 18px; border-bottom: 1px solid #f9fafb; }
.profil-field:last-child { border-bottom: none; }
.profil-field__label { font-size: 11.5px; color: #9ca3af; font-weight: 500; flex-shrink: 0; }
.profil-field__value { font-size: 13px; font-weight: 600; color: #111827; text-align: right; }
.pw-dots { letter-spacing: 3px; color: #6b7280; }

.profil-status-badge { font-size: 11px; font-weight: 700; padding: 3px 10px; border-radius: 100px; }
.profil-status-badge--ok  { background: #f0fdf4; color: #16a34a; }

.profil-security-action { padding: 14px 18px; border-top: 1px solid #f0faf0; display: flex; flex-direction: column; gap: 8px; }
.btn-ubah-pw-secondary {
  display: inline-flex; align-items: center; gap: 6px;
  background: #f0fdf4; border: 1px solid #bbf7d0; color: #16a34a;
  border-radius: 9px; padding: 8px 14px; font-size: 12px; font-weight: 600;
  font-family: "Poppins", sans-serif; cursor: pointer; align-self: flex-start;
}
.btn-ubah-pw-secondary:hover { background: #dcfce7; }
.profil-security-hint { font-size: 11px; color: #9ca3af; margin: 0; line-height: 1.6; }

.profil-skeleton { padding: 16px 18px; display: flex; flex-direction: column; gap: 10px; }
.profil-skeleton__bar { height: 14px; background: #f3f4f6; border-radius: 6px; animation: pulse 1.5s ease-in-out infinite; }
.profil-skeleton__bar--short { width: 60%; }
@keyframes pulse { 0%,100% { opacity: 1; } 50% { opacity: .5; } }

.modal-backdrop { position: fixed; inset: 0; background: rgba(0,0,0,0.4); display: flex; align-items: center; justify-content: center; z-index: 1000; }
.modal-box { background: #fff; border-radius: 16px; padding: 0 24px 24px; width: 100%; max-width: 420px; box-shadow: 0 20px 60px rgba(0,0,0,0.15); }
.modal-pw__head { display: flex; align-items: center; justify-content: space-between; padding: 20px 0 16px; }
.modal-pw__title { font-size: 16px; font-weight: 700; color: #111827; }
.modal-pw__close { background: #f3f4f6; border: none; border-radius: 8px; width: 30px; height: 30px; display: flex; align-items: center; justify-content: center; cursor: pointer; color: #6b7280; }
.modal-pw__close:hover { background: #e5e7eb; }
.pw-form { display: flex; flex-direction: column; gap: 14px; }
.pw-field { display: flex; flex-direction: column; gap: 5px; }
.pw-label { font-size: 12px; font-weight: 600; color: #374151; }
.pw-input { border: 1.5px solid #e5e7eb; border-radius: 9px; padding: 9px 12px; font-size: 13px; font-family: "Poppins", sans-serif; outline: none; transition: border-color .15s; }
.pw-input:focus { border-color: #48AF4A; }
.pw-hint { font-size: 11px; color: #9ca3af; }
.ap-required { color: #ef4444; }
.pw-alert { padding: 8px 12px; border-radius: 8px; font-size: 12.5px; }
.pw-alert--err { background: #fff1f2; color: #be123c; border: 1px solid #fecdd3; }
.pw-alert--ok  { background: #f0fdf4; color: #16a34a; border: 1px solid #bbf7d0; }
.modal-box__actions { display: flex; gap: 10px; justify-content: flex-end; margin-top: 18px; }
.btn-cancel  { background: #f3f4f6; color: #374151; border: none; border-radius: 9px; padding: 9px 20px; font-size: 13px; font-weight: 600; font-family: "Poppins",sans-serif; cursor: pointer; }
.btn-confirm { background: #48AF4A; color: #fff; border: none; border-radius: 9px; padding: 9px 20px; font-size: 13px; font-weight: 600; font-family: "Poppins",sans-serif; cursor: pointer; }
.btn-confirm:disabled { background: #d1d5db; cursor: default; }

@media (max-width: 500px) {
  .profil-hero { flex-wrap: wrap; }
  .btn-ubah-pw { width: 100%; justify-content: center; }
}
</style>
