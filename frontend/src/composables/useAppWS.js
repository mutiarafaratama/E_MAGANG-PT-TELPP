const _handlers = new Set()
let _ws = null
let _reconnectTimer = null
let _refCount = 0

function _forceLogout() {
  localStorage.removeItem('access_token')
  localStorage.removeItem('user')
  if (_ws) { _ws.onclose = null; _ws.close(); _ws = null }
  if (_reconnectTimer) { clearTimeout(_reconnectTimer); _reconnectTimer = null }
  _refCount = 0
  _handlers.clear()
  window.location.href = '/login?reason=akun_dihapus'
}

function _connect() {
  if (_ws && (_ws.readyState === WebSocket.OPEN || _ws.readyState === WebSocket.CONNECTING)) return
  const token = localStorage.getItem('access_token')
  if (!token) return
  const proto = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  _ws = new WebSocket(`${proto}//${window.location.host}/api/ws?token=${token}`)
  _ws.onmessage = (e) => {
    try {
      const msg = JSON.parse(e.data)
      if (msg.type === 'force_logout') { _forceLogout(); return }
      _handlers.forEach(h => h(msg))
    } catch {}
  }
  _ws.onclose = () => {
    _ws = null
    if (_refCount > 0) {
      _reconnectTimer = setTimeout(_connect, 3000)
    }
  }
}

export function useAppWS() {
  function connect() {
    _refCount++
    _connect()
  }

  function disconnect() {
    _refCount = Math.max(0, _refCount - 1)
    if (_refCount === 0) {
      if (_reconnectTimer) clearTimeout(_reconnectTimer)
      if (_ws) { _ws.onclose = null; _ws.close(); _ws = null }
    }
  }

  function subscribe(handler) {
    _handlers.add(handler)
    return () => _handlers.delete(handler)
  }

  return { connect, disconnect, subscribe }
}
