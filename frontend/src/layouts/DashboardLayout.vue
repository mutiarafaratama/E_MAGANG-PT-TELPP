<template>
  <div class="app-shell" :class="{ 'sidebar-collapsed': collapsed }">

    <!-- ── SIDEBAR ── -->
    <aside class="sidebar" :class="{ 'sidebar--open': mobileOpen }">

      <!-- Logo + Hamburger -->
      <div class="sidebar-header">
        <div class="sidebar-logo">
          <div class="sidebar-logo__mark">
            <img src="/logotel.png" alt="PT TELPP" class="sidebar-logo__img" />
          </div>
          <div v-if="!collapsed" class="sidebar-logo__text">
            <span class="sidebar-logo__name">e-Magang</span>
            <span class="sidebar-logo__company">PT TELPP</span>
          </div>
        </div>
        <button class="sidebar-toggle" @click="toggleCollapse" :title="collapsed ? 'Perluas sidebar' : 'Ciutkan sidebar'">
          <svg v-if="!collapsed" width="16" height="16" viewBox="0 0 24 24" fill="none">
            <line x1="3" y1="6" x2="21" y2="6" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
            <line x1="3" y1="12" x2="21" y2="12" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
            <line x1="3" y1="18" x2="21" y2="18" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
          </svg>
          <svg v-else width="16" height="16" viewBox="0 0 24 24" fill="none">
            <path d="M9 18l6-6-6-6" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
        </button>
        <button class="sidebar-mobile-close" @click="mobileOpen = false">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none"><path d="M18 6L6 18M6 6l12 12" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>
        </button>
      </div>

      <!-- User Card -->
      <div class="sidebar-user" @click="goToProfile" style="cursor:pointer" title="Profil Saya">
        <div class="sidebar-user__avatar">{{ initials }}</div>
        <div v-if="!collapsed" class="sidebar-user__info">
          <div class="sidebar-user__name">{{ user?.nama_lengkap }}</div>
          <div class="sidebar-user__role">{{ roleName }}</div>
        </div>
      </div>

      <!-- Nav Groups -->
      <nav class="sidebar-nav">
        <template v-for="(group, gi) in navGroups" :key="gi">
          <div v-if="!collapsed && group.label" class="sidebar-nav__label">{{ group.label }}</div>
          <div v-else-if="collapsed && gi > 0" class="sidebar-nav__divider"></div>
          <template v-for="item in group.items" :key="item.key">
            <!-- Parent item with children (expandable) -->
            <template v-if="item.children?.length">
              <button
                class="nav-btn nav-btn--parent"
                :class="{ 'nav-btn--active': isChildActive(item), 'nav-btn--expanded': expandedKeys.has(item.key) }"
                @click="handleNavClick(item)"
                :title="collapsed ? item.label : ''"
              >
                <span class="nav-btn__icon">
                  <span v-html="item.icon"></span>
                </span>
                <span v-if="item.dot" class="nav-btn__dot"></span>
                <span v-if="!collapsed" class="nav-btn__label">{{ item.label }}</span>
                <span v-if="!collapsed && item.badge" class="nav-btn__badge">{{ item.badge }}</span>
                <svg v-if="!collapsed" class="nav-btn__chevron" width="12" height="12" viewBox="0 0 24 24" fill="none">
                  <path d="M6 9l6 6 6-6" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"/>
                </svg>
              </button>
              <!-- Children -->
              <div v-if="!collapsed && expandedKeys.has(item.key)" class="nav-children">
                <button
                  v-for="child in item.children"
                  :key="child.key"
                  class="nav-btn nav-btn--child"
                  :class="{ 'nav-btn--active': activeTab === child.key }"
                  @click="handleNavClick(child)"
                >
                  <span class="nav-btn__child-dot"></span>
                  <span class="nav-btn__label">{{ child.label }}</span>
                  <span v-if="child.badge" class="nav-btn__badge">{{ child.badge }}</span>
                  <span v-if="child.dot" class="nav-btn__dot nav-btn__dot--inline"></span>
                </button>
              </div>
            </template>
            <!-- Regular item -->
            <button
              v-else
              class="nav-btn"
              :class="{ 'nav-btn--active': activeTab === item.key }"
              @click="handleNavClick(item)"
              :title="collapsed ? item.label : ''"
            >
              <span class="nav-btn__icon">
                <span v-html="item.icon"></span>
              </span>
              <span v-if="item.dot" class="nav-btn__dot"></span>
              <span v-if="!collapsed" class="nav-btn__label">{{ item.label }}</span>
              <span v-if="!collapsed && item.badge" class="nav-btn__badge">{{ item.badge }}</span>
            </button>
          </template>
        </template>
      </nav>

      <!-- Logout -->
      <button class="sidebar-logout" @click="logout" :title="collapsed ? 'Keluar' : ''">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none">
          <path d="M9 21H5a2 2 0 01-2-2V5a2 2 0 012-2h4" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
          <polyline points="16 17 21 12 16 7" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          <line x1="21" y1="12" x2="9" y2="12" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
        </svg>
        <span v-if="!collapsed">Keluar</span>
      </button>
    </aside>

    <!-- Mobile overlay -->
    <div v-if="mobileOpen" class="sidebar-overlay" @click="mobileOpen = false"></div>

    <!-- ── MAIN AREA ── -->
    <div class="main-area">

      <!-- Topbar -->
      <header class="topbar">
        <button class="topbar-hamburger" @click="mobileOpen = true">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
            <line x1="3" y1="6" x2="21" y2="6" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
            <line x1="3" y1="12" x2="21" y2="12" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
            <line x1="3" y1="18" x2="21" y2="18" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
          </svg>
        </button>

        <div class="topbar-breadcrumb">
          <span class="topbar-breadcrumb__app">e-Magang</span>
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none"><path d="M9 18l6-6-6-6" stroke="#9ca3af" stroke-width="2" stroke-linecap="round"/></svg>
          <span class="topbar-breadcrumb__page">{{ currentPageLabel }}</span>
        </div>

        <div class="topbar-right">
          <div class="topbar-role-chip">
            <span class="topbar-role-chip__dot"></span>
            {{ roleName }}
          </div>

          <!-- Bell Notifikasi -->
          <div class="notif-wrap" ref="notifWrap">
            <button class="notif-btn" @click.stop="toggleNotif" :title="notifBadge ? `${notifBadge} notifikasi baru` : 'Notifikasi'">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none"><path d="M18 8A6 6 0 006 8c0 7-3 9-3 9h18s-3-2-3-9" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"/><path d="M13.73 21a2 2 0 01-3.46 0" stroke="currentColor" stroke-width="1.8" stroke-linecap="round"/></svg>
              <span v-if="notifBadge > 0" class="notif-badge">{{ notifBadge > 99 ? '99+' : notifBadge }}</span>
            </button>

            <div v-if="notifOpen" class="notif-dropdown">
              <div class="notif-dropdown__head">
                <span class="notif-dropdown__title">Notifikasi</span>
                <button v-if="notifBadge > 0" class="notif-baca-semua" @click="markAllRead">Tandai semua dibaca</button>
              </div>
              <div v-if="notifLoading" class="notif-loading">
                <div class="notif-spinner"></div>
              </div>
              <div v-else-if="notifList.length === 0" class="notif-empty">
                <svg width="30" height="30" viewBox="0 0 24 24" fill="none"><path d="M18 8A6 6 0 006 8c0 7-3 9-3 9h18s-3-2-3-9" stroke="#d1d5db" stroke-width="1.5"/><path d="M13.73 21a2 2 0 01-3.46 0" stroke="#d1d5db" stroke-width="1.5"/></svg>
                <span>Belum ada notifikasi</span>
              </div>
              <div v-else class="notif-list">
                <div
                  v-for="n in notifList"
                  :key="n.id"
                  :class="['notif-item', !n.is_read && 'notif-item--unread']"
                  @click="markRead(n)"
                >
                  <div v-if="!n.is_read" class="notif-item__dot"></div>
                  <div v-else class="notif-item__dot notif-item__dot--read"></div>
                  <div class="notif-item__body">
                    <div class="notif-item__title">{{ n.judul }}</div>
                    <div class="notif-item__msg">{{ n.pesan }}</div>
                    <div class="notif-item__time">{{ formatNotifTime(n.created_at) }}</div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div class="topbar-avatar" :title="user?.nama_lengkap" @click="goToProfile" style="cursor:pointer">{{ initials }}</div>
        </div>
      </header>

      <!-- Page Content -->
      <main class="page-content">
        <slot />
      </main>
    </div>
  </div>

  <!-- Modal Peringatan Ganti Password -->
  <Teleport to="body">
    <div v-if="showPwWarning" class="pw-warn-backdrop">
      <div class="pw-warn-box">
        <div class="pw-warn-icon">
          <svg width="36" height="36" viewBox="0 0 24 24" fill="none">
            <rect x="3" y="11" width="18" height="11" rx="2" stroke="#f59e0b" stroke-width="2"/>
            <path d="M7 11V7a5 5 0 0110 0v4" stroke="#f59e0b" stroke-width="2" stroke-linecap="round"/>
            <circle cx="12" cy="16" r="1.5" fill="#f59e0b"/>
          </svg>
        </div>
        <div class="pw-warn-title">Segera Ganti Password</div>
        <div class="pw-warn-msg">
          Akun Anda masih menggunakan <strong>password default</strong> yang dikirimkan melalui email.
          Demi keamanan, harap ganti password Anda sekarang.
        </div>
        <div class="pw-warn-actions">
          <button class="pw-warn-btn pw-warn-btn--later" @click="showPwWarning = false">Nanti</button>
          <button class="pw-warn-btn pw-warn-btn--now" @click="goToProfileFromWarning">
            <svg width="13" height="13" viewBox="0 0 24 24" fill="none"><rect x="3" y="11" width="18" height="11" rx="2" stroke="currentColor" stroke-width="2"/><path d="M7 11V7a5 5 0 0110 0v4" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>
            Ubah Sekarang
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from "vue";
import { useRouter } from "vue-router";
import { useAuth } from "@/hooks/useAuth";
import api from "@/lib/api";
import { setupPushNotifications } from "@/services/fcm";

export interface NavItem {
  key: string;
  label: string;
  icon: string;
  badge?: string | number;
  dot?: boolean;
  route?: string;
  children?: NavItem[];
}

export interface NavGroup {
  label?: string;
  items: NavItem[];
}

const props = defineProps<{
  navGroups: NavGroup[];
  roleName: string;
  defaultTab?: string;
}>();

const emit = defineEmits<{
  (e: "tab-change", tab: string): void;
}>();

const { user, logout } = useAuth();
const router = useRouter();

const activeTab = ref(props.defaultTab ?? props.navGroups[0]?.items[0]?.key ?? "");
const collapsed = ref(false);
const mobileOpen = ref(false);
const expandedKeys = ref<Set<string>>(new Set());

const initials = computed(() => {
  const name = user.value?.nama_lengkap ?? "";
  return name.split(" ").map((n: string) => n[0]).slice(0, 2).join("").toUpperCase();
});

const allItems = computed(() => props.navGroups.flatMap(g => g.items.flatMap(i => i.children ? i.children : [i])));
const currentPageLabel = computed(() => {
  const flat = props.navGroups.flatMap(g => g.items);
  const parent = flat.find(i => i.children?.some(c => c.key === activeTab.value));
  if (parent) return parent.label;
  return flat.find(i => i.key === activeTab.value)?.label ?? allItems.value.find(i => i.key === activeTab.value)?.label ?? "Dashboard";
});

function isChildActive(item: NavItem): boolean {
  return !!item.children?.some(c => c.key === activeTab.value);
}

function toggleCollapse() {
  collapsed.value = !collapsed.value;
}

function toggleExpand(key: string) {
  if (expandedKeys.value.has(key)) {
    expandedKeys.value.delete(key);
  } else {
    expandedKeys.value.add(key);
  }
  expandedKeys.value = new Set(expandedKeys.value);
}

function handleNavClick(item: NavItem) {
  if (item.children?.length) {
    if (collapsed.value) collapsed.value = false;
    toggleExpand(item.key);
    return;
  }
  if (item.route) {
    router.push(item.route);
  } else {
    activeTab.value = item.key;
    mobileOpen.value = false;
    emit("tab-change", item.key);
  }
}

function goToProfile() {
  activeTab.value = "profil";
  mobileOpen.value = false;
  emit("tab-change", "profil");
}

// ── Notifikasi ─────────────────────────────────────────────
interface Notif {
  id: string;
  judul: string;
  pesan: string;
  tipe: string;
  referensi_id: string | null;
  is_read: boolean;
  created_at: string;
}

// ── Password warning ────────────────────────────────────────
const showPwWarning = ref(false);

async function checkPasswordChanged() {
  try {
    const r = await api.get("/api/auth/me");
    if (r.data?.data?.password_changed === false) {
      showPwWarning.value = true;
    }
  } catch { /* silent */ }
}

function goToProfileFromWarning() {
  showPwWarning.value = false;
  goToProfile();
}

const notifOpen    = ref(false);
const notifList    = ref<Notif[]>([]);
const notifBadge   = ref(0);
const notifLoading = ref(false);
const notifWrap    = ref<HTMLElement | null>(null);

async function fetchNotifBadge() {
  try {
    const r = await api.get("/api/notifikasi/badge");
    notifBadge.value = r.data?.data?.total_unread ?? 0;
  } catch { /* silent */ }
}

async function fetchNotifList() {
  notifLoading.value = true;
  try {
    const r = await api.get("/api/notifikasi");
    notifList.value = r.data?.data ?? [];
  } catch { /* silent */ } finally {
    notifLoading.value = false;
  }
}

async function toggleNotif() {
  notifOpen.value = !notifOpen.value;
  if (notifOpen.value) fetchNotifList();
}

function markRead(n: Notif) {
  if (!n.is_read) {
    n.is_read = true;
    notifBadge.value = Math.max(0, notifBadge.value - 1);
    api.patch(`/api/notifikasi/${n.id}/read`).catch(() => {});
  }
}

function markAllRead() {
  notifList.value.forEach(n => { n.is_read = true; });
  notifBadge.value = 0;
  api.patch("/api/notifikasi/read-all").catch(() => {});
}

function formatNotifTime(iso: string) {
  const d = new Date(iso);
  const diff = (Date.now() - d.getTime()) / 1000;
  if (diff < 60)    return "Baru saja";
  if (diff < 3600)  return `${Math.floor(diff / 60)} menit lalu`;
  if (diff < 86400) return `${Math.floor(diff / 3600)} jam lalu`;
  return d.toLocaleDateString("id-ID", { day: "numeric", month: "short" });
}

function onClickOutside(e: MouseEvent) {
  if (notifWrap.value && !notifWrap.value.contains(e.target as Node)) {
    notifOpen.value = false;
  }
}

let pollInterval: ReturnType<typeof setInterval> | null = null;

// Tangani navigasi dari klik notifikasi background (via SW postMessage)
function onSwMessage(event: MessageEvent) {
  if (event.data?.type === 'NAVIGATE' && event.data.route) {
    router.push(event.data.route).catch(() => {});
  }
}

onMounted(() => {
  document.addEventListener("click", onClickOutside);
  fetchNotifBadge();
  checkPasswordChanged();
  pollInterval = setInterval(fetchNotifBadge, 30_000);

  // Inisialisasi FCM push notification — minta izin & daftarkan token
  setupPushNotifications().catch(() => {});

  // Dengarkan pesan navigasi dari service worker (klik notif background)
  if ('serviceWorker' in navigator) {
    navigator.serviceWorker.addEventListener('message', onSwMessage);
  }
});

onUnmounted(() => {
  document.removeEventListener("click", onClickOutside);
  if (pollInterval) clearInterval(pollInterval);
  if ('serviceWorker' in navigator) {
    navigator.serviceWorker.removeEventListener('message', onSwMessage);
  }
});

defineExpose({ activeTab });
</script>

<style scoped>
* { box-sizing: border-box; }

/* ── SHELL ── */
.app-shell {
  display: flex;
  min-height: 100vh;
  background: #f1f5f0;
  font-family: "Poppins", "Inter", sans-serif;
  --sidebar-w: 256px;
  --sidebar-w-collapsed: 68px;
  --green-dark: #0d2818;
  --green-mid: #1a4228;
  --green-accent: #48AF4A;
  --green-light: #86efac;
  --green-muted: rgba(134,239,172,0.55);
}

/* ── SIDEBAR ── */
.sidebar {
  width: var(--sidebar-w);
  flex-shrink: 0;
  background: var(--green-dark);
  display: flex;
  flex-direction: column;
  position: sticky;
  top: 0;
  height: 100vh;
  overflow-y: auto;
  overflow-x: hidden;
  transition: width 0.22s cubic-bezier(0.4,0,0.2,1);
  scrollbar-width: none;
  z-index: 40;
}
.sidebar::-webkit-scrollbar { display: none; }
.app-shell.sidebar-collapsed .sidebar { width: var(--sidebar-w-collapsed); }

/* Header */
.sidebar-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 18px 14px 16px;
  border-bottom: 1px solid rgba(255,255,255,0.06);
  gap: 8px;
  min-height: 68px;
}

.sidebar-logo {
  display: flex;
  align-items: center;
  gap: 10px;
  overflow: hidden;
  min-width: 0;
}

.sidebar-logo__mark { flex-shrink: 0; display: flex; }
.sidebar-logo__img {
  height: 36px;
  width: 36px;
  object-fit: contain;
  border-radius: 6px;
  background: #fff;
  padding: 3px;
}

.sidebar-logo__text {
  display: flex;
  flex-direction: column;
  overflow: hidden;
}
.sidebar-logo__name {
  font-size: 13px;
  font-weight: 700;
  color: #fff;
  white-space: nowrap;
  line-height: 1.2;
}
.sidebar-logo__company {
  font-size: 10px;
  font-weight: 600;
  color: var(--green-light);
  white-space: nowrap;
  letter-spacing: 0.08em;
  text-transform: uppercase;
}

.sidebar-toggle {
  flex-shrink: 0;
  background: rgba(255,255,255,0.06);
  border: none;
  color: rgba(255,255,255,0.6);
  width: 30px;
  height: 30px;
  border-radius: 7px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: background 0.15s, color 0.15s;
}
.sidebar-toggle:hover { background: rgba(255,255,255,0.12); color: #fff; }

.sidebar-mobile-close {
  display: none;
  flex-shrink: 0;
  background: rgba(255,255,255,0.06);
  border: none;
  color: rgba(255,255,255,0.6);
  width: 30px;
  height: 30px;
  border-radius: 7px;
  align-items: center;
  justify-content: center;
  cursor: pointer;
}

/* User Card */
.sidebar-user {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 14px 14px;
  border-bottom: 1px solid rgba(255,255,255,0.06);
  overflow: hidden;
  min-height: 64px;
}

.sidebar-user__avatar {
  width: 38px;
  height: 38px;
  border-radius: 10px;
  background: linear-gradient(135deg, #48AF4A, #1a5c20);
  color: #fff;
  font-size: 13px;
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  border: 2px solid rgba(134,239,172,0.3);
}

.sidebar-user__info { overflow: hidden; min-width: 0; }
.sidebar-user__name {
  font-size: 12.5px;
  font-weight: 600;
  color: #fff;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.sidebar-user__role {
  font-size: 10.5px;
  color: var(--green-light);
  margin-top: 1px;
  white-space: nowrap;
}

/* Nav */
.sidebar-nav {
  flex: 1;
  padding: 10px 8px;
  display: flex;
  flex-direction: column;
  gap: 1px;
}

.sidebar-nav__label {
  font-size: 9.5px;
  font-weight: 700;
  color: var(--green-muted);
  letter-spacing: 0.1em;
  text-transform: uppercase;
  padding: 12px 8px 4px;
  white-space: nowrap;
  overflow: hidden;
}

.sidebar-nav__divider {
  height: 1px;
  background: rgba(255,255,255,0.06);
  margin: 6px 4px;
}

.nav-btn {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 9px 10px;
  border-radius: 9px;
  background: none;
  border: none;
  color: rgba(255,255,255,0.6);
  font-size: 12.5px;
  font-family: "Poppins", sans-serif;
  font-weight: 500;
  cursor: pointer;
  text-align: left;
  transition: background 0.14s, color 0.14s;
  width: 100%;
  white-space: nowrap;
  overflow: visible;
  position: relative;
}
.nav-btn:hover { background: rgba(255,255,255,0.07); color: #fff; }
.nav-btn--active {
  background: rgba(72,175,74,0.18);
  color: var(--green-light);
}
.nav-btn--active::before {
  content: "";
  position: absolute;
  left: 0;
  top: 20%;
  bottom: 20%;
  width: 3px;
  border-radius: 0 3px 3px 0;
  background: var(--green-accent);
}
.nav-btn--active:hover { background: rgba(72,175,74,0.24); }

.nav-btn__icon { display: flex; align-items: center; flex-shrink: 0; width: 16px; }
.nav-btn__label { flex: 1; overflow: hidden; text-overflow: ellipsis; }
.nav-btn__badge {
  background: #ef4444;
  color: #fff;
  font-size: 9.5px;
  font-weight: 700;
  border-radius: 100px;
  padding: 1px 6px;
  flex-shrink: 0;
}

/* Expandable parent nav */
.nav-btn__chevron {
  flex-shrink: 0;
  color: rgba(255,255,255,0.35);
  transition: transform 0.2s ease;
}
.nav-btn--expanded .nav-btn__chevron { transform: rotate(180deg); }
.nav-btn--parent.nav-btn--active .nav-btn__chevron { color: var(--green-light); }

/* Children container */
.nav-children {
  display: flex;
  flex-direction: column;
  gap: 1px;
  margin: 2px 0 4px 0;
  padding-left: 14px;
  border-left: 1.5px solid rgba(134,239,172,0.2);
  margin-left: 18px;
}

/* Child nav item */
.nav-btn--child {
  padding: 7px 10px 7px 8px;
  font-size: 12px;
  border-radius: 7px;
}
.nav-btn__child-dot {
  width: 5px;
  height: 5px;
  border-radius: 50%;
  background: rgba(255,255,255,0.25);
  flex-shrink: 0;
  margin-right: 2px;
  transition: background 0.15s;
}
.nav-btn--child.nav-btn--active .nav-btn__child-dot { background: var(--green-accent); }
.nav-btn--child:hover .nav-btn__child-dot { background: rgba(255,255,255,0.5); }

/* Red dot indicator */
.nav-btn__dot {
  position: absolute;
  top: 5px;
  left: 23px;
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #ef4444;
  border: 1.5px solid #1a5c20;
  pointer-events: none;
  z-index: 1;
}
.nav-btn__dot--inline {
  position: static;
  width: 7px;
  height: 7px;
  border: 1.5px solid rgba(26,92,32,0.6);
  margin-left: auto;
  flex-shrink: 0;
}

/* Logout */
.sidebar-logout {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 18px;
  margin: 6px 8px;
  border-radius: 9px;
  background: none;
  border: none;
  color: rgba(255,255,255,0.35);
  font-size: 12.5px;
  font-family: "Poppins", sans-serif;
  font-weight: 500;
  cursor: pointer;
  transition: color 0.14s, background 0.14s;
  white-space: nowrap;
  overflow: hidden;
}
.sidebar-logout:hover { color: #fca5a5; background: rgba(239,68,68,0.1); }

/* ── MAIN AREA ── */
.main-area {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-width: 0;
  overflow: hidden;
}

/* Topbar */
.topbar {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 0 20px;
  height: 58px;
  background: #fff;
  border-bottom: 1px solid #e9f0e9;
  position: sticky;
  top: 0;
  z-index: 30;
  flex-shrink: 0;
  box-shadow: 0 1px 3px rgba(13,40,24,0.06);
}

.topbar-hamburger {
  display: none;
  background: none;
  border: none;
  cursor: pointer;
  color: #374151;
  padding: 6px;
  border-radius: 8px;
  flex-shrink: 0;
}
.topbar-hamburger:hover { background: #f3f4f6; }

.topbar-breadcrumb {
  display: flex;
  align-items: center;
  gap: 6px;
  flex: 1;
  min-width: 0;
}
.topbar-breadcrumb__app {
  font-size: 12px;
  color: #9ca3af;
  font-weight: 500;
  white-space: nowrap;
}
.topbar-breadcrumb__page {
  font-size: 14px;
  font-weight: 700;
  color: #111827;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.topbar-right {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-shrink: 0;
}

.topbar-role-chip {
  display: flex;
  align-items: center;
  gap: 6px;
  background: #f0fdf4;
  border: 1px solid #bbf7d0;
  color: #16a34a;
  font-size: 11px;
  font-weight: 600;
  padding: 4px 12px;
  border-radius: 100px;
}
.topbar-role-chip__dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: #48AF4A;
}

.topbar-avatar {
  width: 34px;
  height: 34px;
  border-radius: 9px;
  background: linear-gradient(135deg, #48AF4A, #1a5c20);
  color: #fff;
  font-size: 12px;
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  border: 2px solid #d1fae5;
}

/* ── Notifikasi Dropdown ───────────────────────────────── */
.notif-wrap { position: relative; }
.notif-btn {
  position: relative; background: none; border: none; cursor: pointer;
  padding: 5px; color: #6b7280; display: flex; align-items: center;
  border-radius: 8px; transition: background .15s, color .15s;
}
.notif-btn:hover { background: #f0fdf4; color: #16a34a; }
.notif-badge {
  position: absolute; top: -3px; right: -5px;
  background: #ef4444; color: #fff; font-size: 9px; font-weight: 800;
  min-width: 16px; height: 16px; border-radius: 100px; padding: 0 4px;
  display: flex; align-items: center; justify-content: center;
  border: 1.5px solid #fff; line-height: 1;
}
.notif-dropdown {
  position: absolute; top: calc(100% + 10px); right: -4px; z-index: 9000;
  width: 340px; background: #fff; border-radius: 14px;
  box-shadow: 0 10px 36px rgba(0,0,0,0.14), 0 2px 8px rgba(0,0,0,0.06);
  border: 1px solid #e9f5e9; overflow: hidden;
  animation: notif-pop .14s ease;
}
@keyframes notif-pop {
  from { opacity: 0; transform: translateY(-8px) scale(0.98); }
  to   { opacity: 1; transform: none; }
}
.notif-dropdown__head {
  display: flex; align-items: center; justify-content: space-between;
  padding: 14px 16px 12px; border-bottom: 1px solid #f0faf0;
}
.notif-dropdown__title { font-size: 13px; font-weight: 700; color: #111827; }
.notif-baca-semua {
  background: none; border: none; cursor: pointer; padding: 0;
  font-size: 11.5px; font-weight: 600; color: #48AF4A;
  font-family: inherit; transition: opacity .15s;
}
.notif-baca-semua:hover { opacity: .75; }
.notif-empty {
  display: flex; flex-direction: column; align-items: center; gap: 8px;
  padding: 32px 20px; color: #9ca3af; font-size: 12.5px;
}
.notif-loading { display: flex; justify-content: center; padding: 28px; }
.notif-spinner {
  width: 22px; height: 22px; border: 2.5px solid #e5e7eb;
  border-top-color: #48AF4A; border-radius: 50%;
  animation: nspin .7s linear infinite;
}
@keyframes nspin { to { transform: rotate(360deg); } }
.notif-list { max-height: 340px; overflow-y: auto; }
.notif-item {
  display: flex; gap: 10px; align-items: flex-start;
  padding: 11px 16px; cursor: pointer;
  border-bottom: 1px solid #f9fafb; transition: background .12s;
}
.notif-item:last-child { border-bottom: none; }
.notif-item:hover { background: #f9fafb; }
.notif-item--unread { background: #f0fdf4; }
.notif-item--unread:hover { background: #e9faf0; }
.notif-item__dot {
  width: 7px; height: 7px; border-radius: 50%; margin-top: 5px; flex-shrink: 0;
  background: #22c55e;
}
.notif-item__dot--read { background: transparent; }
.notif-item__body { flex: 1; min-width: 0; }
.notif-item__title { font-size: 12.5px; font-weight: 700; color: #111827; margin-bottom: 2px; }
.notif-item__msg { font-size: 11.5px; color: #6b7280; line-height: 1.5; }
.notif-item__time { font-size: 10.5px; color: #9ca3af; margin-top: 4px; }

/* Page Content */
.page-content {
  flex: 1;
  padding: 24px;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

/* Mobile overlay */
.sidebar-overlay {
  display: none;
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,0.45);
  z-index: 39;
  backdrop-filter: blur(2px);
}

/* ── RESPONSIVE ── */
@media (max-width: 900px) {
  .sidebar {
    position: fixed;
    left: calc(-1 * var(--sidebar-w));
    top: 0;
    height: 100vh;
    transition: left 0.25s ease, width 0.22s cubic-bezier(0.4,0,0.2,1);
    z-index: 50;
  }
  .sidebar--open {
    left: 0 !important;
    width: var(--sidebar-w) !important;
  }
  .sidebar-mobile-close { display: flex; }
  .sidebar-toggle { display: none; }
  .sidebar-overlay { display: block; }
  .topbar-hamburger { display: flex; }
  .app-shell.sidebar-collapsed .sidebar { width: var(--sidebar-w); left: calc(-1 * var(--sidebar-w)); }
}

@media (max-width: 600px) {
  .page-content { padding: 16px; }
  .topbar { padding: 0 14px; }
  .topbar-role-chip { display: none; }

  /* Notif dropdown: full-width, fixed di bawah topbar agar tidak overflow */
  .notif-dropdown {
    position: fixed;
    top: 58px;
    left: 8px;
    right: 8px;
    width: auto;
    max-height: calc(100dvh - 80px);
    border-radius: 12px;
  }

  .notif-list {
    max-height: calc(100dvh - 180px);
  }
}

/* ── Password Warning Modal ── */
.pw-warn-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
  backdrop-filter: blur(3px);
}

.pw-warn-box {
  background: #fff;
  border-radius: 16px;
  padding: 32px 28px 24px;
  width: min(420px, calc(100vw - 32px));
  box-shadow: 0 20px 60px rgba(0,0,0,0.18);
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  gap: 12px;
}

.pw-warn-icon {
  width: 64px;
  height: 64px;
  border-radius: 50%;
  background: #f0fdf4;
  color: #48AF4A;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 4px;
}

.pw-warn-title {
  font-size: 1.125rem;
  font-weight: 700;
  color: #111827;
}

.pw-warn-msg {
  font-size: 0.875rem;
  color: #6b7280;
  line-height: 1.6;
}

.pw-warn-msg strong {
  color: #1a5c20;
}

.pw-warn-actions {
  display: flex;
  gap: 10px;
  margin-top: 8px;
  width: 100%;
}

.pw-warn-btn {
  flex: 1;
  padding: 10px 16px;
  border-radius: 8px;
  font-size: 0.875rem;
  font-weight: 600;
  cursor: pointer;
  border: none;
  transition: background 0.15s, transform 0.1s;
}

.pw-warn-btn:active { transform: scale(0.98); }

.pw-warn-btn--later {
  background: #f1f5f9;
  color: #6b7280;
}

.pw-warn-btn--later:hover { background: #e2e8f0; }

.pw-warn-btn--now {
  background: #48AF4A;
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
}

.pw-warn-btn--now:hover { background: #48AF4A; }
</style>
