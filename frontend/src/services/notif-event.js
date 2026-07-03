// Mini event emitter untuk notifikasi foreground
// Dipakai oleh fcm.js dan useAppWS.js untuk emit;
// DashboardLayout.vue listen dan tampilkan toast

const _listeners = new Set()

export function emitNotif(payload) {
  _listeners.forEach(fn => {
    try { fn(payload) } catch (_) {}
  })
}

// Daftarkan listener, kembalikan fungsi untuk unsubscribe
export function onNotif(fn) {
  _listeners.add(fn)
  return () => _listeners.delete(fn)
}
