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
      <div class="sidebar-user">
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
          <button
            v-for="item in group.items"
            :key="item.key"
            class="nav-btn"
            :class="{ 'nav-btn--active': activeTab === item.key }"
            @click="handleNavClick(item)"
            :title="collapsed ? item.label : ''"
          >
            <span class="nav-btn__icon" v-html="item.icon"></span>
            <span v-if="!collapsed" class="nav-btn__label">{{ item.label }}</span>
            <span v-if="!collapsed && item.badge" class="nav-btn__badge">{{ item.badge }}</span>
          </button>
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
          <div class="topbar-avatar" :title="user?.nama_lengkap">{{ initials }}</div>
        </div>
      </header>

      <!-- Page Content -->
      <main class="page-content">
        <slot />
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from "vue";
import { useRouter } from "vue-router";
import { useAuth } from "@/hooks/useAuth";

export interface NavItem {
  key: string;
  label: string;
  icon: string;
  badge?: string | number;
  route?: string;
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

const initials = computed(() => {
  const name = user.value?.nama_lengkap ?? "";
  return name.split(" ").map((n: string) => n[0]).slice(0, 2).join("").toUpperCase();
});

const allItems = computed(() => props.navGroups.flatMap(g => g.items));
const currentPageLabel = computed(() => allItems.value.find(i => i.key === activeTab.value)?.label ?? "Dashboard");

function toggleCollapse() {
  collapsed.value = !collapsed.value;
}

function handleNavClick(item: NavItem) {
  if (item.route) {
    router.push(item.route);
  } else {
    activeTab.value = item.key;
    mobileOpen.value = false;
    emit("tab-change", item.key);
  }
}

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
  background: linear-gradient(135deg, #48AF4A, #2d8f30);
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
  overflow: hidden;
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
  background: linear-gradient(135deg, #48AF4A, #2d8f30);
  color: #fff;
  font-size: 12px;
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  border: 2px solid #d1fae5;
}

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
}
</style>
