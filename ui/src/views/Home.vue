<template>
  <div class="home">
    <Hero
      title="Obsidian Dashboard"
      subtitle="Manage your Minecraft servers"
      showActionButton
      actionButtonText="Manage Servers"
      actionButtonIcon="hugeicons:server-stack-03"
      actionButtonTo="/servers"
    />

    <div class="stats-grid">
      <StatCard
        title="Total Servers"
        :value="servers.length"
        icon="mdi:server"
        color="primary"
      />
      <StatCard
        title="Running"
        :value="runningServers"
        icon="mdi:play-circle"
        color="success"
      />
      <StatCard
        title="Stopped"
        :value="stoppedServers"
        icon="mdi:stop-circle"
        color="error"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import { listServers } from "../helper";
import StatCard from "../components/StatCard.vue";
import Hero from "../components/Hero.vue";
import type { ServerInfo } from "../types";

const servers = ref<ServerInfo[]>([]);

const runningServers = computed(
  () => servers.value.filter((s) => s.state === "running").length
);

const stoppedServers = computed(
  () => servers.value.filter((s) => s.state === "stopped").length
);

onMounted(async () => {
  try {
    servers.value = await listServers();
  } catch (error) {
    console.error("Failed to load servers:", error);
  }
});
</script>

<style scoped>
.home {
  max-width: 1400px;
  margin: 0 auto;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: var(--space-xl);
  padding: 0 var(--section-padding-x) var(--section-padding-y);
}
</style>
