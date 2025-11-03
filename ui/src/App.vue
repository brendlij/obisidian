<script setup lang="ts">
import { onMounted } from "vue";
import { useTheme } from "./composables/useTheme";
import Header from "./components/Header.vue";

const { initializeTheme } = useTheme();

onMounted(async () => {
  await initializeTheme();
});
</script>

<template>
  <Header />
  <main class="main-container">
    <router-view />
  </main>
</template>

<style scoped>
.main-container {
  min-height: calc(100vh - 60px);
  background: var(--color-background);
}
</style>
