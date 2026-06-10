import { createRouter, createWebHistory } from "vue-router";
import LandingPage from "@/views/LandingPage.vue";

const routes = [
  {
    path: "/",
    name: "landing",
    component: LandingPage,
  },
  {
    path: "/login",
    name: "login",
    component: () => import("@/views/LoginPage.vue"),
    meta: { guest: true },
  },
  {
    path: "/register",
    name: "register",
    component: () => import("@/views/RegisterPage.vue"),
    meta: { guest: true },
  },
  {
    path: "/dashboard/peserta",
    name: "dashboard-peserta",
    component: () => import("@/views/DashboardPeserta.vue"),
    meta: { requiresAuth: true, role: "peserta" },
  },
  {
    path: "/dashboard/hrd",
    name: "dashboard-hrd",
    component: () => import("@/views/DashboardHRD.vue"),
    meta: { requiresAuth: true, role: "hrd" },
  },
  {
    path: "/dashboard/admin",
    name: "dashboard-admin",
    component: () => import("@/views/DashboardAdmin.vue"),
    meta: { requiresAuth: true, role: "admin" },
  },
  {
    path: "/admin/landing-settings",
    name: "admin-landing-settings",
    component: () => import("@/views/admin/LandingSettings.vue"),
    meta: { requiresAuth: true, role: "admin" },
  },
  {
    path: "/dashboard",
    redirect: () => {
      const userStr = localStorage.getItem("user");
      const user = userStr ? JSON.parse(userStr) : null;
      if (user?.role === "admin") return "/dashboard/admin";
      if (user?.role === "hrd") return "/dashboard/hrd";
      return "/dashboard/peserta";
    },
  },
  {
    path: "/:pathMatch(.*)*",
    redirect: "/",
  },
];

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
  scrollBehavior(to, from, savedPosition) {
    if (to.hash) return { el: to.hash, behavior: "smooth" };
    if (savedPosition) return savedPosition;
    return { top: 0 };
  },
});

router.beforeEach((to) => {
  const token = localStorage.getItem("access_token");
  const userStr = localStorage.getItem("user");
  const user = userStr ? JSON.parse(userStr) : null;

  if (to.meta.requiresAuth && !token) {
    return { name: "login" };
  }

  if (to.meta.guest && token) {
    if (user?.role === "admin") return { name: "dashboard-admin" };
    if (user?.role === "hrd") return { name: "dashboard-hrd" };
    return { name: "dashboard-peserta" };
  }

  if (to.meta.role && user?.role !== to.meta.role) {
    if (user?.role === "admin") return { name: "dashboard-admin" };
    if (user?.role === "hrd") return { name: "dashboard-hrd" };
    return { name: "dashboard-peserta" };
  }
});

export default router;
