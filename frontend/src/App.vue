<template>
  <RouterView />
</template>

<script setup lang="ts">
import { onMounted } from "vue";
import { RouterView } from "vue-router";
import { setupPushNotifications } from "@/services/fcm";

onMounted(async () => {
  const token = localStorage.getItem("access_token");
  if (token) {
    await setupPushNotifications();
  }
});
</script>
