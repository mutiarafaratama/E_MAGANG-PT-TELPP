<template>
  <div class="auth-page">
    <div class="auth-bg"></div>
    <div class="auth-overlay"></div>

    <div class="auth-container">
      <div class="auth-brand">
        <img src="/logo_emagang.png" alt="PT TELPP" class="auth-brand__logo" />
        <span>e-Magang <strong>PT TELPP</strong></span>
      </div>

      <div class="auth-card">
        <div v-if="!token" class="state-error">
          <svg width="40" height="40" viewBox="0 0 24 24" fill="none">
            <circle cx="12" cy="12" r="10" fill="#fee2e2" stroke="#dc2626" stroke-width="1.5"/>
            <path d="M12 8v4m0 4h.01" stroke="#dc2626" stroke-width="2" stroke-linecap="round"/>
          </svg>
          <h2>Tautan Tidak Valid</h2>
          <p>Tautan reset kata sandi tidak ditemukan atau sudah kedaluwarsa.</p>
          <router-link to="/lupa-kata-sandi" class="btn-primary">Minta Tautan Baru</router-link>
        </div>

        <template v-else>
          <div class="auth-card__header">
            <div class="icon-wrap">
              <svg width="28" height="28" viewBox="0 0 24 24" fill="none">
                <rect x="3" y="11" width="18" height="11" rx="2" stroke="#48AF4A" stroke-width="2"/>
                <path d="M7 11V7a5 5 0 0110 0v4" stroke="#48AF4A" stroke-width="2"/>
                <circle cx="12" cy="16" r="1.5" fill="#48AF4A"/>
              </svg>
            </div>
            <h1>Buat Kata Sandi Baru</h1>
            <p>Kata sandi harus minimal 8 karakter dan mengandung huruf, angka, serta karakter spesial.</p>
          </div>

          <Transition name="fade-slide" mode="out-in">
            <div v-if="sukses" key="sukses" class="sukses-state">
              <div class="sukses-icon">
                <svg width="36" height="36" viewBox="0 0 24 24" fill="none">
                  <circle cx="12" cy="12" r="10" fill="#dcfce7" stroke="#16a34a" stroke-width="1.5"/>
                  <path d="M7 13l3.5 3.5L17 8.5" stroke="#16a34a" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"/>
                </svg>
              </div>
              <h2>Kata Sandi Diperbarui!</h2>
              <p>Kata sandi Anda berhasil diubah. Silakan masuk dengan kata sandi baru.</p>
              <router-link to="/login?reset=1" class="btn-primary">Masuk Sekarang</router-link>
            </div>

            <div v-else key="form" class="auth-form">
              <div class="form-field">
                <label>Kata Sandi Baru</label>
                <div class="field-input">
                  <span class="field-icon">
                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none">
                      <rect x="3" y="11" width="18" height="11" rx="2" ry="2" stroke="#9ca3af" stroke-width="1.8"/>
                      <path d="M7 11V7a5 5 0 0110 0v4" stroke="#9ca3af" stroke-width="1.8"/>
                    </svg>
                  </span>
                  <input
                    v-model="form.password"
                    :type="showPass ? 'text' : 'password'"
                    placeholder="Minimal 8 karakter"
                    autocomplete="new-password"
                  />
                  <button type="button" class="field-eye" @click="showPass = !showPass" tabindex="-1">
                    <svg v-if="!showPass" width="16" height="16" viewBox="0 0 24 24" fill="none"><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z" stroke="#9ca3af" stroke-width="1.8"/><circle cx="12" cy="12" r="3" stroke="#9ca3af" stroke-width="1.8"/></svg>
                    <svg v-else width="16" height="16" viewBox="0 0 24 24" fill="none"><path d="M17.94 17.94A10.07 10.07 0 0112 20c-7 0-11-8-11-8a18.45 18.45 0 015.06-5.94M9.9 4.24A9.12 9.12 0 0112 4c7 0 11 8 11 8a18.5 18.5 0 01-2.16 3.19m-6.72-1.07a3 3 0 11-4.24-4.24" stroke="#9ca3af" stroke-width="1.8" stroke-linecap="round"/><line x1="1" y1="1" x2="23" y2="23" stroke="#9ca3af" stroke-width="1.8" stroke-linecap="round"/></svg>
                  </button>
                </div>
                <div class="strength-bar">
                  <div class="strength-bar__fill" :class="strengthClass" :style="{ width: strengthPct + '%' }"></div>
                </div>
                <p class="strength-label" :class="strengthClass">{{ strengthLabel }}</p>
              </div>

              <div class="form-field">
                <label>Konfirmasi Kata Sandi Baru</label>
                <div class="field-input">
                  <span class="field-icon">
                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none">
                      <rect x="3" y="11" width="18" height="11" rx="2" ry="2" stroke="#9ca3af" stroke-width="1.8"/>
                      <path d="M7 11V7a5 5 0 0110 0v4" stroke="#9ca3af" stroke-width="1.8"/>
                    </svg>
                  </span>
                  <input
                    v-model="form.konfirmasi"
                    :type="showKonfirmasi ? 'text' : 'password'"
                    placeholder="Ulangi kata sandi baru"
                    autocomplete="new-password"
                    @keyup.enter="simpan"
                  />
                  <button type="button" class="field-eye" @click="showKonfirmasi = !showKonfirmasi" tabindex="-1">
                    <svg v-if="!showKonfirmasi" width="16" height="16" viewBox="0 0 24 24" fill="none"><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z" stroke="#9ca3af" stroke-width="1.8"/><circle cx="12" cy="12" r="3" stroke="#9ca3af" stroke-width="1.8"/></svg>
                    <svg v-else width="16" height="16" viewBox="0 0 24 24" fill="none"><path d="M17.94 17.94A10.07 10.07 0 0112 20c-7 0-11-8-11-8a18.45 18.45 0 015.06-5.94M9.9 4.24A9.12 9.12 0 0112 4c7 0 11 8 11 8a18.5 18.5 0 01-2.16 3.19m-6.72-1.07a3 3 0 11-4.24-4.24" stroke="#9ca3af" stroke-width="1.8" stroke-linecap="round"/><line x1="1" y1="1" x2="23" y2="23" stroke="#9ca3af" stroke-width="1.8" stroke-linecap="round"/></svg>
                  </button>
                </div>
              </div>

              <Transition name="err-fade">
                <div v-if="error" class="alert alert--error">{{ error }}</div>
              </Transition>

              <button type="button" class="btn-primary" :disabled="loading" @click="simpan">
                <svg v-if="loading" class="spinner-icon" width="16" height="16" viewBox="0 0 24 24" fill="none">
                  <circle cx="12" cy="12" r="10" stroke="rgba(255,255,255,0.3)" stroke-width="2.5"/>
                  <path d="M12 2a10 10 0 0110 10" stroke="white" stroke-width="2.5" stroke-linecap="round"/>
                </svg>
                {{ loading ? 'Menyimpan...' : 'Simpan Kata Sandi Baru' }}
              </button>
            </div>
          </Transition>
        </template>

        <router-link to="/login" class="auth-back-link">← Kembali ke Login</router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import axios from 'axios'

const route = useRoute()
const token = ref('')

onMounted(() => {
  token.value = route.query.token || ''
})

const form = ref({ password: '', konfirmasi: '' })
const showPass = ref(false)
const showKonfirmasi = ref(false)
const loading = ref(false)
const error = ref(null)
const sukses = ref(false)

const strengthScore = computed(() => {
  const p = form.value.password
  if (!p) return 0
  let s = 0
  if (p.length >= 8)  s++
  if (/[A-Z]/.test(p)) s++
  if (/[0-9]/.test(p)) s++
  if (/[^A-Za-z0-9]/.test(p)) s++
  return s
})

const strengthPct   = computed(() => [0, 25, 50, 75, 100][strengthScore.value])
const strengthClass = computed(() => ['', 'weak', 'fair', 'good', 'strong'][strengthScore.value])
const strengthLabel = computed(() => ['', 'Lemah', 'Sedang', 'Baik', 'Kuat'][strengthScore.value])

async function simpan() {
  if (loading.value) return
  error.value = null

  if (!form.value.password) {
    error.value = 'Kata sandi wajib diisi.'
    return
  }
  if (form.value.password !== form.value.konfirmasi) {
    error.value = 'Konfirmasi kata sandi tidak cocok.'
    return
  }
  if (form.value.password.length < 8) {
    error.value = 'Kata sandi minimal 8 karakter.'
    return
  }

  loading.value = true
  try {
    await axios.post('/api/auth/reset-password', {
      token: token.value,
      new_password: form.value.password
    })
    sukses.value = true
  } catch (e) {
    error.value = e.response?.data?.message || 'Gagal menyimpan kata sandi. Tautan mungkin sudah kedaluwarsa.'
  } finally {
    loading.value = false
  }
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
  position: fixed; inset: 0;
  background: linear-gradient(135deg, #0b1c30 0%, #0f2d14 50%, #0b1c30 100%);
  z-index: 0;
}

.auth-overlay {
  position: fixed; inset: 0;
  background:
    radial-gradient(ellipse 70% 60% at 20% 40%, rgba(72,175,74,0.12) 0%, transparent 60%),
    radial-gradient(ellipse 50% 50% at 80% 70%, rgba(11,28,48,0.8) 0%, transparent 70%);
  z-index: 1;
}

.auth-container {
  position: relative; z-index: 2;
  display: flex; flex-direction: column; align-items: center;
  width: 100%; max-width: 400px; gap: 20px;
}

.auth-brand {
  display: flex; align-items: center; gap: 10px;
  color: #fff; font-size: 15px; font-weight: 700; text-decoration: none;
}
.auth-brand__logo {
  height: 36px; width: auto; object-fit: contain;
  border-radius: 6px; background: #fff; padding: 3px 5px;
}
.auth-brand strong { color: #86efac; }

.auth-card {
  width: 100%; background: #f0f0f0; border-radius: 18px;
  padding: 36px 32px 28px;
  box-shadow: 0 24px 60px rgba(0,0,0,0.45), 0 0 0 1px rgba(255,255,255,0.06);
}

.auth-card__header { text-align: center; margin-bottom: 28px; }

.icon-wrap {
  display: flex; align-items: center; justify-content: center;
  width: 56px; height: 56px; background: #dcfce7;
  border-radius: 50%; margin: 0 auto 16px;
}

.auth-card__header h1 { font-size: 20px; font-weight: 700; color: #111827; margin: 0 0 8px; }
.auth-card__header p  { font-size: 13px; color: #6b7280; line-height: 1.6; margin: 0; }

.auth-form {
  display: flex; flex-direction: column; gap: 16px; margin-bottom: 20px;
}

.form-field { display: flex; flex-direction: column; gap: 6px; }
.form-field label { font-size: 12px; font-weight: 600; color: #374151; }

.field-input { position: relative; display: flex; align-items: center; }

.field-icon {
  position: absolute; left: 12px; display: flex;
  align-items: center; pointer-events: none; z-index: 1;
}

.field-input input {
  width: 100%; padding: 11px 36px 11px 38px;
  background: #ffffff; border: 1.5px solid #e5e7eb;
  border-radius: 10px; font-size: 13.5px;
  font-family: "Poppins", sans-serif; color: #111827;
  outline: none; transition: border-color 0.15s, box-shadow 0.15s;
}
.field-input input:focus { border-color: #48AF4A; box-shadow: 0 0 0 3px rgba(72,175,74,0.12); }
.field-input input::placeholder { color: #b4b9c1; }

.field-eye {
  position: absolute; right: 12px; background: none; border: none;
  cursor: pointer; padding: 0; display: flex; align-items: center; color: #9ca3af;
}

.strength-bar {
  height: 4px; background: #e5e7eb; border-radius: 2px; overflow: hidden; margin-top: 2px;
}
.strength-bar__fill { height: 100%; border-radius: 2px; transition: width 0.3s, background 0.3s; }
.strength-bar__fill.weak  { background: #ef4444; }
.strength-bar__fill.fair  { background: #f59e0b; }
.strength-bar__fill.good  { background: #3b82f6; }
.strength-bar__fill.strong { background: #16a34a; }

.strength-label { font-size: 11px; margin: 0; }
.strength-label.weak   { color: #ef4444; }
.strength-label.fair   { color: #f59e0b; }
.strength-label.good   { color: #3b82f6; }
.strength-label.strong { color: #16a34a; }

.alert {
  display: flex; align-items: center; gap: 8px;
  padding: 10px 14px; border-radius: 9px; font-size: 13px; font-weight: 500;
}
.alert--error { background: #fef2f2; color: #dc2626; border: 1px solid #fecaca; }

.btn-primary {
  width: 100%; padding: 12px; background: #48AF4A; color: #fff;
  border: none; border-radius: 10px; font-size: 14px; font-weight: 700;
  font-family: "Poppins", sans-serif; cursor: pointer;
  transition: opacity 0.15s, transform 0.1s;
  display: flex; align-items: center; justify-content: center;
  gap: 8px; text-decoration: none; margin-top: 4px;
}
.btn-primary:hover:not(:disabled) { opacity: 0.88; transform: translateY(-1px); }
.btn-primary:active:not(:disabled) { transform: translateY(0); }
.btn-primary:disabled { opacity: 0.55; cursor: not-allowed; }

.spinner-icon { animation: spin 0.8s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }

.sukses-state {
  display: flex; flex-direction: column; align-items: center;
  text-align: center; gap: 12px; margin-bottom: 20px;
}
.sukses-icon { margin-bottom: 4px; }
.sukses-state h2 { font-size: 18px; font-weight: 700; color: #111827; margin: 0; }
.sukses-state p  { font-size: 13.5px; color: #374151; line-height: 1.7; margin: 0; }

.state-error {
  display: flex; flex-direction: column; align-items: center;
  text-align: center; gap: 12px; padding: 8px 0 20px;
}
.state-error h2 { font-size: 18px; font-weight: 700; color: #111827; margin: 0; }
.state-error p  { font-size: 13.5px; color: #6b7280; margin: 0; }

.err-fade-enter-active, .err-fade-leave-active { transition: opacity 0.2s ease, transform 0.2s ease; }
.err-fade-enter-from { opacity: 0; transform: translateY(-4px); }
.err-fade-leave-to   { opacity: 0; transform: translateY(-4px); }

.fade-slide-enter-active, .fade-slide-leave-active { transition: opacity 0.25s ease, transform 0.25s ease; }
.fade-slide-enter-from { opacity: 0; transform: translateY(8px); }
.fade-slide-leave-to   { opacity: 0; transform: translateY(-8px); }

.auth-back-link {
  display: block; text-align: center; margin-top: 14px;
  font-size: 12px; color: rgba(255,255,255,0.45);
  text-decoration: none; transition: color 0.15s;
}
.auth-back-link:hover { color: rgba(255,255,255,0.75); }

@media (max-width: 480px) {
  .auth-card { padding: 28px 20px 22px; }
}
</style>
