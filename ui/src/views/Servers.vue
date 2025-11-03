<template>
  <div class="servers">
    <Hero
      title="Manage Servers"
      subtitle="View and manage all your Minecraft servers"
    />

    <div class="servers__container">
      <div class="servers__actions">
        <router-link to="/servers/create" class="servers__create-link">
          <Button size="md" icon="mdi:plus" variant="primary">
            Create Server
          </Button>
        </router-link>
      </div>

      <div v-if="servers.length === 0" class="servers__empty">
        <p>No servers created yet. Create your first server to get started!</p>
      </div>

      <div v-else class="servers-grid">
        <ServerCard
          v-for="server in servers"
          :key="server.config.id"
          :server="server"
          @updated="loadServers"
          @delete="confirmDelete(server.config.id, server.config.name)"
        />
      </div>

      <ConfirmationModal
        :is-open="deleteModal.isOpen"
        title="Delete Server"
        :message="`Are you sure you want to delete '${deleteModal.serverName}'? This action cannot be undone.`"
        confirm-text="Delete"
        cancel-text="Cancel"
        :is-dangerous="true"
        :is-loading="deleteModal.isLoading"
        @confirm="performDelete"
        @cancel="deleteModal.isOpen = false"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onUnmounted } from "vue";
import {
  listServers,
  deleteServer as apiDeleteServer,
  subscribeToServerEvents,
} from "../helper";
import Button from "../components/Button.vue";
import ServerCard from "../components/ServerCard.vue";
import ConfirmationModal from "../components/ConfirmationModal.vue";
import Hero from "../components/Hero.vue";
import type { ServerInfo } from "../types";

const servers = ref<ServerInfo[]>([]);
let eventSources: Map<string, EventSource> = new Map();

const deleteModal = reactive({
  isOpen: false,
  serverId: "",
  serverName: "",
  isLoading: false,
});

const loadServers = async () => {
  try {
    servers.value = await listServers();
    // Subscribe to events for each server
    subscribeToServerStates();
  } catch (error) {
    console.error("Failed to load servers:", error);
  }
};

const subscribeToServerStates = () => {
  // Close old subscriptions
  eventSources.forEach((es) => es.close());
  eventSources.clear();

  // Subscribe to each server
  servers.value.forEach((server) => {
    const es = subscribeToServerEvents(server.config.id, {
      onStateChange: async () => {
        // Reload all servers when any state changes
        try {
          servers.value = await listServers();
        } catch (error) {
          console.error("Failed to refresh servers:", error);
        }
      },
      onError: (error) => {
        console.error("Event source error:", error);
      },
    });
    eventSources.set(server.config.id, es);
  });
};

const confirmDelete = (id: string, name: string) => {
  deleteModal.serverId = id;
  deleteModal.serverName = name;
  deleteModal.isOpen = true;
};

const performDelete = async () => {
  deleteModal.isLoading = true;
  try {
    await apiDeleteServer(deleteModal.serverId);
    deleteModal.isOpen = false;
    await loadServers();
  } catch (error) {
    console.error("Failed to delete server:", error);
  } finally {
    deleteModal.isLoading = false;
  }
};

onMounted(loadServers);

onUnmounted(() => {
  eventSources.forEach((es) => es.close());
  eventSources.clear();
});
</script>

<style scoped>
.servers {
  flex: 1;
}

.servers__container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 var(--section-padding-x) var(--section-padding-y);
}

.servers__actions {
  display: flex;
  gap: var(--space-md);
  margin-bottom: var(--section-gap);
}

.servers__create-link {
  text-decoration: none;
}

.servers__empty {
  text-align: center;
  padding: var(--section-padding-y) var(--section-padding-x);
  background: var(--color-surface);
  border: var(--border-width-thin) dashed var(--color-border);
  border-radius: var(--radius-lg);
  color: var(--color-text-secondary);
}

.servers-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: var(--space-xl);
}
</style>
