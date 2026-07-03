// Suara notifikasi menggunakan Web Audio API
// Otomatis mengikuti volume perangkat — tidak berbunyi jika device silent/mute

let _ctx = null

function getCtx() {
  if (!_ctx || _ctx.state === 'closed') {
    _ctx = new (window.AudioContext || window.webkitAudioContext)()
  }
  // Resume jika suspended (browser suspend AudioContext sebelum ada interaksi user)
  if (_ctx.state === 'suspended') {
    _ctx.resume().catch(() => {})
  }
  return _ctx
}

// Nada "ding" dua-tone — lembut dan tidak mengganggu
export function playNotifSound() {
  try {
    const ctx  = getCtx()
    const now  = ctx.currentTime

    function tone(freq, startTime, duration, volume = 0.25) {
      const osc  = ctx.createOscillator()
      const gain = ctx.createGain()

      osc.type = 'sine'
      osc.frequency.setValueAtTime(freq, startTime)

      gain.gain.setValueAtTime(0, startTime)
      gain.gain.linearRampToValueAtTime(volume, startTime + 0.01)
      gain.gain.exponentialRampToValueAtTime(0.001, startTime + duration)

      osc.connect(gain)
      gain.connect(ctx.destination)

      osc.start(startTime)
      osc.stop(startTime + duration)
    }

    // Dua nada pendek: 880 Hz → 1100 Hz (seperti "ding-ding" notif)
    tone(880,  now,        0.18)
    tone(1100, now + 0.15, 0.22)
  } catch (_) {
    // Browser tidak support atau audio diblokir — silent fail
  }
}
