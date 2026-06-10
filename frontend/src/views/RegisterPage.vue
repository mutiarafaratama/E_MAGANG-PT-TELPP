<template>
  <div class="auth-page">
    <div class="auth-bg"></div>
    <div class="auth-overlay"></div>

    <div class="auth-container">
      <div class="auth-brand">
        <img src="/logotel.png" alt="PT TELPP" class="auth-brand__logo" />
        <span>e-Magang <strong>PT TELPP</strong></span>
      </div>

      <div class="auth-card">
        <div class="auth-card__header">
          <h1>Buat akun baru</h1>
          <p>Gabung untuk bisa mengajukan magang</p>
        </div>

        <form @submit.prevent="handleRegister" class="auth-form">
          <div class="form-field">
            <label>Nama Lengkap</label>
            <div class="field-input">
              <span class="field-icon">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none"><path d="M20 21v-2a4 4 0 00-4-4H8a4 4 0 00-4 4v2" stroke="#9ca3af" stroke-width="1.8" stroke-linecap="round"/><circle cx="12" cy="7" r="4" stroke="#9ca3af" stroke-width="1.8"/></svg>
              </span>
              <input
                v-model="form.nama_lengkap"
                type="text"
                placeholder="Nama sesuai KTP"
                autocomplete="name"
                required
              />
            </div>
          </div>

          <div class="form-field">
            <label>Email</label>
            <div class="field-input">
              <span class="field-icon">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none"><path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z" stroke="#9ca3af" stroke-width="1.8"/><polyline points="22,6 12,13 2,6" stroke="#9ca3af" stroke-width="1.8"/></svg>
              </span>
              <input
                v-model="form.email"
                type="email"
                placeholder="contoh@email.com"
                autocomplete="email"
                required
              />
            </div>
          </div>

          <div class="form-field">
            <label>No Telpon</label>
            <div class="field-input">
              <span class="field-icon">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none"><path d="M22 16.92v3a2 2 0 01-2.18 2 19.79 19.79 0 01-8.63-3.07A19.5 19.5 0 013.07 11 19.79 19.79 0 01.22 2.37a2 2 0 012-2.18H5a2 2 0 012 1.72 12.84 12.84 0 00.7 2.81 2 2 0 01-.45 2.11L6.09 7.91a16 16 0 006 6l1.27-1.27a2 2 0 012.11-.45 12.84 12.84 0 002.81.7A2 2 0 0122 16.92z" stroke="#9ca3af" stroke-width="1.8"/></svg>
              </span>
              <input
                v-model="form.no_hp"
                type="tel"
                placeholder="08xxxxxxxxxx"
                autocomplete="tel"
              />
            </div>
          </div>

          <div class="form-row">
            <div class="form-field">
              <label>Password</label>
              <div class="field-input">
                <span class="field-icon">
                  <svg width="16" height="16" viewBox="0 0 24 24" fill="none"><rect x="3" y="11" width="18" height="11" rx="2" ry="2" stroke="#9ca3af" stroke-width="1.8"/><path d="M7 11V7a5 5 0 0110 0v4" stroke="#9ca3af" stroke-width="1.8"/></svg>
                </span>
                <input
                  v-model="form.password"
                  :type="showPass ? 'text' : 'password'"
                  placeholder="Min. 8 karakter"
                  autocomplete="new-password"
                  required
                  minlength="8"
                />
                <button type="button" class="field-eye" @click="showPass = !showPass" tabindex="-1">
                  <svg v-if="!showPass" width="14" height="14" viewBox="0 0 24 24" fill="none"><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z" stroke="#9ca3af" stroke-width="1.8"/><circle cx="12" cy="12" r="3" stroke="#9ca3af" stroke-width="1.8"/></svg>
                  <svg v-else width="14" height="14" viewBox="0 0 24 24" fill="none"><path d="M17.94 17.94A10.07 10.07 0 0112 20c-7 0-11-8-11-8a18.45 18.45 0 015.06-5.94M9.9 4.24A9.12 9.12 0 0112 4c7 0 11 8 11 8a18.5 18.5 0 01-2.16 3.19m-6.72-1.07a3 3 0 11-4.24-4.24" stroke="#9ca3af" stroke-width="1.8" stroke-linecap="round"/><line x1="1" y1="1" x2="23" y2="23" stroke="#9ca3af" stroke-width="1.8" stroke-linecap="round"/></svg>
                </button>
              </div>
            </div>

            <div class="form-field">
              <label>Konfirmasi ulang</label>
              <div class="field-input">
                <span class="field-icon">
                  <svg width="16" height="16" viewBox="0 0 24 24" fill="none"><rect x="3" y="11" width="18" height="11" rx="2" ry="2" stroke="#9ca3af" stroke-width="1.8"/><path d="M7 11V7a5 5 0 0110 0v4" stroke="#9ca3af" stroke-width="1.8"/></svg>
                </span>
                <input
                  v-model="form.konfirmasi"
                  :type="showKonfirmasi ? 'text' : 'password'"
                  placeholder="Ulangi password"
                  autocomplete="new-password"
                  required
                />
                <button type="button" class="field-eye" @click="showKonfirmasi = !showKonfirmasi" tabindex="-1">
                  <svg v-if="!showKonfirmasi" width="14" height="14" viewBox="0 0 24 24" fill="none"><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z" stroke="#9ca3af" stroke-width="1.8"/><circle cx="12" cy="12" r="3" stroke="#9ca3af" stroke-width="1.8"/></svg>
                  <svg v-else width="14" height="14" viewBox="0 0 24 24" fill="none"><path d="M17.94 17.94A10.07 10.07 0 0112 20c-7 0-11-8-11-8a18.45 18.45 0 015.06-5.94M9.9 4.24A9.12 9.12 0 0112 4c7 0 11 8 11 8a18.5 18.5 0 01-2.16 3.19m-6.72-1.07a3 3 0 11-4.24-4.24" stroke="#9ca3af" stroke-width="1.8" stroke-linecap="round"/><line x1="1" y1="1" x2="23" y2="23" stroke="#9ca3af" stroke-width="1.8" stroke-linecap="round"/></svg>
                </button>
              </div>
              <span v-if="passwordMismatch" class="field-error">Password tidak cocok</span>
            </div>
          </div>

          <div v-if="error" class="alert alert--error">{{ error }}</div>

          <button type="submit" class="btn-primary" :disabled="loading || passwordMismatch">
            <svg v-if="loading" class="spinner-icon" width="16" height="16" viewBox="0 0 24 24" fill="none"><circle cx="12" cy="12" r="10" stroke="rgba(255,255,255,0.3)" stroke-width="2.5"/><path d="M12 2a10 10 0 0110 10" stroke="white" stroke-width="2.5" stroke-linecap="round"/></svg>
            {{ loading ? "Mendaftarkan..." : "Daftar" }}
          </button>
        </form>

        <p class="auth-footer-text">
          Sudah Punya Akun?
          <router-link to="/login">Login Di Sini</router-link>
        </p>

        <router-link to="/" class="auth-back-link">← Kembali ke Beranda</router-link>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from "vue";
import { useAuth } from "@/hooks/useAuth";

const { register, loading, error } = useAuth();

const form = ref({
  nama_lengkap: "",
  email: "",
  no_hp: "",
  password: "",
  konfirmasi: "",
});

const showPass = ref(false);
const showKonfirmasi = ref(false);

const passwordMismatch = computed(
  () => form.value.konfirmasi.length > 0 && form.value.password !== form.value.konfirmasi
);

async function handleRegister() {
  if (passwordMismatch.value) return;
  await register({
    nama_lengkap: form.value.nama_lengkap,
    email: form.value.email,
    password: form.value.password,
  });
}
</script>

<style scoped>
* { box-sizing: border-box; }

.auth-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  font-family: "Poppins", sans-serif;
  position: relative;
  padding: 24px 16px;
}

.auth-bg {
  position: fixed;
  inset: 0;
  background: linear-gradient(135deg, #0b1c30 0%, #0f2d14 50%, #0b1c30 100%);
  z-index: 0;
}

.auth-overlay {
  position: fixed;
  inset: 0;
  background:
    radial-gradient(ellipse 70% 60% at 20% 40%, rgba(72,175,74,0.12) 0%, transparent 60%),
    radial-gradient(ellipse 50% 50% at 80% 70%, rgba(11,28,48,0.8) 0%, transparent 70%);
  z-index: 1;
}

.auth-container {
  position: relative;
  z-index: 2;
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 100%;
  max-width: 420px;
  gap: 20px;
}

.auth-brand {
  display: flex;
  align-items: center;
  gap: 10px;
  color: #fff;
  font-size: 15px;
  font-weight: 700;
  text-decoration: none;
}
.auth-brand__logo {
  height: 36px;
  width: auto;
  object-fit: contain;
  border-radius: 6px;
  background: #fff;
  padding: 3px 5px;
}
.auth-brand strong { color: #86efac; }

.auth-card {
  width: 100%;
  background: #f0f0f0;
  border-radius: 18px;
  padding: 32px 28px 24px;
  box-shadow: 0 24px 60px rgba(0,0,0,0.45), 0 0 0 1px rgba(255,255,255,0.06);
}

.auth-card__header {
  text-align: center;
  margin-bottom: 24px;
}

.auth-card__header h1 {
  font-size: 20px;
  font-weight: 700;
  color: #111827;
  margin-bottom: 5px;
}

.auth-card__header p {
  font-size: 13px;
  color: #6b7280;
}

.auth-form {
  display: flex;
  flex-direction: column;
  gap: 14px;
  margin-bottom: 18px;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}

.form-field {
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.form-field label {
  font-size: 11.5px;
  font-weight: 600;
  color: #374151;
}

.field-input {
  position: relative;
  display: flex;
  align-items: center;
}

.field-icon {
  position: absolute;
  left: 11px;
  display: flex;
  align-items: center;
  pointer-events: none;
  z-index: 1;
}

.field-input input {
  width: 100%;
  padding: 10px 11px 10px 35px;
  background: #ffffff;
  border: 1.5px solid #e5e7eb;
  border-radius: 9px;
  font-size: 13px;
  font-family: "Poppins", sans-serif;
  color: #111827;
  outline: none;
  transition: border-color 0.15s, box-shadow 0.15s;
}

.field-input input:focus {
  border-color: #48AF4A;
  box-shadow: 0 0 0 3px rgba(72,175,74,0.12);
}

.field-input input::placeholder { color: #b4b9c1; }

.field-eye {
  position: absolute;
  right: 10px;
  background: none;
  border: none;
  cursor: pointer;
  padding: 0;
  display: flex;
  align-items: center;
}

.field-error {
  font-size: 11px;
  color: #dc2626;
  margin-top: 2px;
}

.alert {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 14px;
  border-radius: 9px;
  font-size: 13px;
  font-weight: 500;
}

.alert--error { background: #fef2f2; color: #dc2626; border: 1px solid #fecaca; }

.btn-primary {
  width: 100%;
  padding: 12px;
  background: #48AF4A;
  color: #fff;
  border: none;
  border-radius: 10px;
  font-size: 14px;
  font-weight: 700;
  font-family: "Poppins", sans-serif;
  cursor: pointer;
  transition: opacity 0.15s, transform 0.1s;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  margin-top: 4px;
}

.btn-primary:hover:not(:disabled) { opacity: 0.88; transform: translateY(-1px); }
.btn-primary:active:not(:disabled) { transform: translateY(0); }
.btn-primary:disabled { opacity: 0.55; cursor: not-allowed; }

.spinner-icon { animation: spin 0.8s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }

.auth-footer-text {
  text-align: center;
  font-size: 13px;
  color: #6b7280;
  margin: 0;
}

.auth-footer-text a {
  color: #48AF4A;
  font-weight: 600;
  text-decoration: none;
}

.auth-footer-text a:hover { text-decoration: underline; }

.auth-back-link {
  display: block;
  text-align: center;
  margin-top: 14px;
  font-size: 12px;
  color: rgba(255,255,255,0.45);
  text-decoration: none;
  transition: color 0.15s;
}

.auth-back-link:hover { color: rgba(255,255,255,0.75); }

@media (max-width: 480px) {
  .auth-card { padding: 24px 16px 20px; }
  .form-row { grid-template-columns: 1fr; }
}
</style>
