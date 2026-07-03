<template>
  <RouterView />
  <InstallBanner />
</template>

<script setup lang="ts">
import { onMounted } from "vue";
import { RouterView } from "vue-router";
import { setupPushNotifications } from "@/services/fcm";
import { startExpiryTimer, isTokenExpired } from "@/hooks/useAuth";
import InstallBanner from "@/components/InstallBanner.vue";

onMounted(async () => {
  const token = localStorage.getItem("access_token");
  if (token) {
    if (isTokenExpired()) {
      localStorage.removeItem("access_token");
      localStorage.removeItem("user");
      localStorage.removeItem("token_expires_at");
      window.location.href = "/login?expired=1";
      return;
    }
    startExpiryTimer();
    await setupPushNotifications();
  }
});
</script>
