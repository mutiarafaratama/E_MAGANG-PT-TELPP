import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    name: 'landing',
    component: () => import('@/views/public/LandingPage.vue')
  },
  {
    path: '/login',
    name: 'login',
    component: () => import('@/views/public/LoginPage.vue')
  },
  {
    path: '/register',
    name: 'register',
    component: () => import('@/views/public/RegisterPage.vue')
  },
  {
    path: '/daftar',
    name: 'daftar',
    component: () => import('@/views/public/FormPengajuan.vue')
  },
  {
    path: '/faq',
    name: 'faq',
    component: () => import('@/views/public/FAQPage.vue')
  },
  {
    path: '/pengajuan',
    name: 'pengajuan',
    component: () => import('@/views/public/FormPengajuan.vue')
  },
  // Admin routes
  {
    path: '/admin',
    name: 'admin',
    component: () => import('@/views/admin/DashboardAdmin.vue'),
    meta: { requiresAuth: true, roles: ['admin'] }
  },
  {
    path: '/admin/users',
    name: 'admin-users',
    component: () => import('@/views/admin/AdminUsers.vue'),
    meta: { requiresAuth: true, roles: ['admin'] }
  },
  {
    path: '/admin/hero',
    name: 'admin-hero',
    component: () => import('@/views/admin/AdminHero.vue'),
    meta: { requiresAuth: true, roles: ['admin'] }
  },
  {
    path: '/admin/jadwal',
    name: 'admin-jadwal',
    component: () => import('@/views/admin/AdminJadwal.vue'),
    meta: { requiresAuth: true, roles: ['admin'] }
  },
  {
    path: '/admin/dokumen',
    name: 'admin-dokumen',
    component: () => import('@/views/admin/AdminDokumen.vue'),
    meta: { requiresAuth: true, roles: ['admin'] }
  },
  {
    path: '/admin/knowledge',
    name: 'admin-knowledge',
    component: () => import('@/views/admin/AdminKnowledge.vue'),
    meta: { requiresAuth: true, roles: ['admin'] }
  },
  {
    path: '/admin/landing',
    name: 'admin-landing',
    component: () => import('@/views/admin/DashboardAdmin.vue'),
    meta: { requiresAuth: true, roles: ['admin'] }
  },
  // HRD/Staff routes
  {
    path: '/staff',
    name: 'staff',
    component: () => import('@/views/staff/DashboardHRD.vue'),
    meta: { requiresAuth: true, roles: ['hrd'] }
  },
  {
    path: '/staff/verifikasi',
    name: 'staff-verifikasi',
    component: () => import('@/views/staff/HRDVerifikasi.vue'),
    meta: { requiresAuth: true, roles: ['hrd'] }
  },
  {
    path: '/staff/penempatan',
    name: 'staff-penempatan',
    component: () => import('@/views/staff/HRDPenempatanMenunggu.vue'),
    meta: { requiresAuth: true, roles: ['hrd'] }
  },
  {
    path: '/staff/berlangsung',
    name: 'staff-berlangsung',
    component: () => import('@/views/staff/HRDBerlangsung.vue'),
    meta: { requiresAuth: true, roles: ['hrd'] }
  },
  {
    path: '/staff/penilaian',
    name: 'staff-penilaian',
    component: () => import('@/views/staff/HRDPenilaian.vue'),
    meta: { requiresAuth: true, roles: ['hrd'] }
  },
  {
    path: '/staff/sertifikat',
    name: 'staff-sertifikat',
    component: () => import('@/views/staff/HRDSertifikat.vue'),
    meta: { requiresAuth: true, roles: ['hrd'] }
  },
  {
    path: '/staff/chat',
    name: 'staff-chat',
    component: () => import('@/views/staff/HRDChat.vue'),
    meta: { requiresAuth: true, roles: ['hrd'] }
  },
  {
    path: '/staff/laporan',
    name: 'staff-laporan',
    component: () => import('@/views/staff/HRDLaporanPeserta.vue'),
    meta: { requiresAuth: true, roles: ['hrd'] }
  },
  {
    path: '/staff/absen',
    name: 'staff-absen',
    component: () => import('@/views/staff/HRDRekapAbsen.vue'),
    meta: { requiresAuth: true, roles: ['hrd'] }
  },
  {
    path: '/staff/perpanjangan',
    name: 'staff-perpanjangan',
    component: () => import('@/views/staff/PerpanjanganHRDView.vue'),
    meta: { requiresAuth: true, roles: ['hrd'] }
  },
  // Peserta/User routes
  {
    path: '/dashboard',
    name: 'dashboard',
    component: () => import('@/views/user/DashboardPeserta.vue'),
    meta: { requiresAuth: true, roles: ['peserta'] }
  },
  {
    path: '/dashboard/profil',
    name: 'peserta-profil',
    component: () => import('@/views/user/PesertaProfil.vue'),
    meta: { requiresAuth: true, roles: ['peserta'] }
  },
  {
    path: '/dashboard/absensi',
    name: 'peserta-absensi',
    component: () => import('@/views/user/PesertaAbsensi.vue'),
    meta: { requiresAuth: true, roles: ['peserta'] }
  },
  {
    path: '/dashboard/chat',
    name: 'peserta-chat',
    component: () => import('@/views/user/PesertaChat.vue'),
    meta: { requiresAuth: true, roles: ['peserta'] }
  },
  {
    path: '/dashboard/dokumen',
    name: 'peserta-dokumen',
    component: () => import('@/views/user/PesertaDokumen.vue'),
    meta: { requiresAuth: true, roles: ['peserta'] }
  },
  {
    path: '/dashboard/nilai',
    name: 'peserta-nilai',
    component: () => import('@/views/user/NilaiView.vue'),
    meta: { requiresAuth: true, roles: ['peserta'] }
  },
  {
    path: '/dashboard/sertifikat',
    name: 'peserta-sertifikat',
    component: () => import('@/views/user/SertifikatView.vue'),
    meta: { requiresAuth: true, roles: ['peserta'] }
  },
  {
    path: '/dashboard/laporan',
    name: 'peserta-laporan',
    component: () => import('@/views/user/LaporanMagangView.vue'),
    meta: { requiresAuth: true, roles: ['peserta'] }
  },
  {
    path: '/dashboard/surat',
    name: 'peserta-surat',
    component: () => import('@/views/user/SuratBalasanView.vue'),
    meta: { requiresAuth: true, roles: ['peserta'] }
  },
  {
    path: '/dashboard/perpanjangan',
    name: 'peserta-perpanjangan',
    component: () => import('@/views/user/PerpanjanganView.vue'),
    meta: { requiresAuth: true, roles: ['peserta'] }
  },
  // Catch all - redirect to landing
  {
    path: '/:pathMatch(.*)*',
    redirect: '/'
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior(to, from, savedPosition) {
    if (to.hash) return { el: to.hash, behavior: 'smooth' }
    if (savedPosition) return savedPosition
    return { top: 0 }
  }
})

router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('access_token')
  const storedUser = localStorage.getItem('user')
  let user = null

  if (storedUser) {
    try {
      user = JSON.parse(storedUser)
    } catch {
      localStorage.removeItem('user')
    }
  }

  if (to.meta.requiresAuth) {
    if (!token) {
      next('/login')
      return
    }

    if (to.meta.roles && user) {
      const userRole = user.role?.toLowerCase()
      if (!to.meta.roles.includes(userRole)) {
        if (userRole === 'admin') next('/admin')
        else if (userRole === 'hrd') next('/staff')
        else next('/dashboard')
        return
      }
    }
  }

  // Redirect logged-in user away from guest pages
  if ((to.name === 'login' || to.name === 'register') && token && user) {
    const role = user.role?.toLowerCase()
    if (role === 'admin') next('/admin')
    else if (role === 'hrd') next('/staff')
    else next('/dashboard')
    return
  }

  next()
})

export default router
