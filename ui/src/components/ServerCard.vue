<script setup lang="ts">
import { ref, reactive, onMounted, onUnmounted } from "vue";
import { useRouter } from "vue-router";
import { Icon } from "@iconify/vue";
import {
  startServer,
  stopServer,
  restartServer,
  subscribeToServerUpdates,
} from "../helper";
import type { ServerInfo } from "../types";

export interface Props {
  server: ServerInfo;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  delete: [];
}>();

const router = useRouter();
const showMenu = ref(false);
const serverData = reactive<ServerInfo>(
  JSON.parse(JSON.stringify(props.server))
);
let unsubscribe: (() => void) | null = null;

const goToDetail = () => {
  router.push(`/servers/${serverData.config.id}`);
};

const handleDelete = () => {
  emit("delete");
  showMenu.value = false;
};

const handleStart = async () => {
  try {
    console.log("[ServerCard] Starting server:", serverData.config.id);
    await startServer(serverData.config.id);
    console.log("[ServerCard] Server started successfully");
  } catch (error) {
    console.error("Failed to start server:", error);
  }
};

const handleStop = async () => {
  try {
    console.log("[ServerCard] Stopping server:", serverData.config.id);
    await stopServer(serverData.config.id);
    console.log("[ServerCard] Server stopped successfully");
  } catch (error) {
    console.error("Failed to stop server:", error);
  }
};

const handleRestart = async () => {
  try {
    console.log("[ServerCard] Restarting server:", serverData.config.id);
    await restartServer(serverData.config.id);
    console.log("[ServerCard] Server restarted successfully");
  } catch (error) {
    console.error("Failed to restart server:", error);
  }
};

const getStatusColor = (status: string) => {
  if (status === "running") return "success";
  if (status === "stopped") return "error";
  return "warning";
};

const getStatusIcon = (status: string) => {
  if (status === "running") return "mdi:play-circle";
  if (status === "stopped") return "mdi:stop-circle";
  return "mdi:pause-circle";
};

const formatUptime = (seconds: number) => {
  if (seconds < 60) return `${seconds}s`;
  if (seconds < 3600) return `${Math.floor(seconds / 60)}m`;
  if (seconds < 86400) return `${Math.floor(seconds / 3600)}h`;
  return `${Math.floor(seconds / 86400)}d`;
};

onMounted(() => {
  // Subscribe to real-time updates
  unsubscribe = subscribeToServerUpdates(
    serverData.config.id,
    (updatedInfo) => {
      console.log("[ServerCard] Received server update:", updatedInfo);
      Object.assign(serverData, updatedInfo);
    }
  );
});

onUnmounted(() => {
  if (unsubscribe) {
    unsubscribe();
  }
});
</script>

<template>
  <div class="server-card">
    <div class="server-card__header">
      <div class="server-card__title-section" @click="goToDetail">
        <Icon icon="mdi:minecraft" class="server-card__server-icon" />
        <div>
          <h3 class="server-card__title">{{ serverData.config.name }}</h3>
          <p class="server-card__meta">{{ serverData.config.version }}</p>
        </div>
      </div>
      <div class="server-card__header-actions">
        <div
          class="server-card__status"
          :class="`server-card__status--${getStatusColor(serverData.state)}`"
        >
          <Icon :icon="getStatusIcon(serverData.state)" />
          <span>{{ serverData.state }}</span>
        </div>
        <div class="server-card__menu" @click.stop>
          <button
            class="server-card__menu-button"
            @click="showMenu = !showMenu"
            :aria-label="showMenu ? 'Close menu' : 'Open menu'"
          >
            <Icon icon="hugeicons:more-horizontal-square-02" />
          </button>
          <div v-if="showMenu" class="server-card__menu-dropdown">
            <button
              class="server-card__menu-item server-card__menu-item--danger"
              :disabled="serverData.state !== 'stopped'"
              @click="handleDelete"
              :title="
                serverData.state !== 'stopped'
                  ? 'Server must be stopped to delete'
                  : ''
              "
            >
              <Icon icon="mdi:delete" />
              Delete
            </button>
          </div>
        </div>
      </div>
    </div>

    <div class="server-card__details">
      <div class="server-card__detail">
        <Icon icon="mdi:information" />
        <span>{{ serverData.config.type }}</span>
      </div>
      <div class="server-card__detail">
        <Icon icon="mdi:network" />
        <span>Port: {{ serverData.config.port }}</span>
      </div>
      <div class="server-card__detail">
        <Icon icon="mdi:clock-outline" />
        <span>Uptime: {{ formatUptime(serverData.uptimeSec) }}</span>
      </div>
      <div
        v-if="serverData.state === 'running' && serverData.players"
        class="server-card__detail"
      >
        <Icon icon="mdi:account-multiple" />
        <span
          >Players: {{ serverData.players.current }}/{{
            serverData.players.max
          }}</span
        >
      </div>
    </div>

    <div class="server-card__actions">
      <button
        v-if="serverData.state === 'stopped'"
        class="action-button action-button--success"
        @click="handleStart"
      >
        <Icon icon="mdi:play" />
        Start
      </button>
      <button
        v-else
        class="action-button action-button--warning"
        @click="handleStop"
      >
        <Icon icon="mdi:stop" />
        Stop
      </button>
      <button
        v-if="serverData.state === 'running'"
        class="action-button action-button--secondary"
        @click="handleRestart"
      >
        <Icon icon="mdi:restart" />
        Restart
      </button>
    </div>
  </div>
</template>

<style scoped>
.server-card {
  display: flex;
  flex-direction: column;
  gap: var(--space-lg);
  padding: var(--space-xl);
  background: var(--color-surface);
  border: var(--border-width-thin) var(--border-style) var(--color-border);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-md);
  transition: all var(--transition-base);
}

.server-card:hover {
  box-shadow: var(--shadow-lg);
  border-color: var(--color-primary);
}

.server-card__header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: var(--space-md);
}

.server-card__header-actions {
  display: flex;
  align-items: flex-start;
  gap: var(--space-md);
}

.server-card__title-section {
  display: flex;
  align-items: center;
  gap: var(--space-md);
  flex: 1;
  min-width: 0;
  cursor: pointer;
  transition: opacity var(--transition-base);
}

.server-card__title-section:hover {
  opacity: 0.8;
}

.server-card__server-icon {
  font-size: 1.5rem;
  color: var(--color-primary);
  flex-shrink: 0;
}

.server-card__title {
  font-size: var(--font-size-lg);
  font-weight: var(--font-weight-semibold);
  margin: 0 0 var(--space-xs) 0;
}

.server-card__meta {
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
  margin: 0;
}

.server-card__status {
  display: inline-flex;
  align-items: center;
  gap: var(--space-xs);
  padding: var(--space-xs) var(--space-md);
  border-radius: var(--radius-md);
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-semibold);
  text-transform: capitalize;
  background: var(--color-surface-alt);
  flex-shrink: 0;
}

.server-card__status--success {
  color: var(--color-success);
  background: rgba(16, 185, 129, 0.1);
}

.server-card__status--error {
  color: var(--color-error);
  background: rgba(239, 68, 68, 0.1);
}

.server-card__status--warning {
  color: var(--color-warning);
  background: rgba(245, 158, 11, 0.1);
}

.server-card__menu {
  position: relative;
  flex-shrink: 0;
}

.server-card__menu-button {
  background: none;
  border: none;
  padding: var(--space-xs);
  font-size: 1.25rem;
  color: var(--color-text-secondary);
  cursor: pointer;
  border-radius: var(--radius-md);
  transition: all var(--transition-base);
  display: flex;
  align-items: center;
  justify-content: center;
}

.server-card__menu-button:hover {
  background: var(--color-surface-alt);
  color: var(--color-primary);
}

.server-card__menu-dropdown {
  position: absolute;
  top: 100%;
  right: 0;
  background: var(--color-surface);
  border: var(--border-width-thin) var(--border-style) var(--color-border);
  border-radius: var(--radius-md);
  box-shadow: var(--shadow-lg);
  z-index: var(--z-dropdown);
  min-width: 150px;
  margin-top: var(--space-sm);
}

.server-card__menu-item {
  display: flex;
  align-items: center;
  gap: var(--space-md);
  width: 100%;
  padding: var(--space-md);
  background: none;
  border: none;
  color: var(--color-text);
  cursor: pointer;
  font-size: var(--font-size-base);
  transition: all var(--transition-base);
  text-align: left;
}

.server-card__menu-item:first-child {
  border-radius: var(--radius-md) var(--radius-md) 0 0;
}

.server-card__menu-item:last-child {
  border-radius: 0 0 var(--radius-md) var(--radius-md);
}

.server-card__menu-item:hover:not(:disabled) {
  background: var(--color-surface-alt);
}

.server-card__menu-item--danger {
  color: var(--color-error);
}

.server-card__menu-item--danger:hover:not(:disabled) {
  background: rgba(239, 68, 68, 0.1);
}

.server-card__menu-item:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  color: var(--color-text-secondary);
}

.server-card__details {
  display: flex;
  flex-direction: column;
  gap: var(--space-sm);
  padding-top: var(--space-md);
  border-top: var(--border-width-thin) var(--border-style) var(--color-border);
}

.server-card__detail {
  display: flex;
  align-items: center;
  gap: var(--space-sm);
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
}

.server-card__detail :deep(svg) {
  font-size: 1rem;
  color: var(--color-primary);
}

.server-card__actions {
  display: flex;
  gap: var(--space-sm);
  padding-top: var(--space-md);
  border-top: var(--border-width-thin) var(--border-style) var(--color-border);
}

.action-button {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-sm);
  padding: var(--space-sm) var(--space-md);
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-semibold);
  border-radius: var(--radius-lg);
  border: var(--border-width-thin) var(--border-style) transparent;
  transition: all var(--transition-base);
  cursor: pointer;
  white-space: nowrap;
  user-select: none;
  background: var(--color-surface-alt);
  color: var(--color-text);
}

.action-button:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
  border-color: currentColor;
}

.action-button:active {
  transform: translateY(0);
}

.action-button--success {
  background: var(--color-success);
  color: white;
  border-color: var(--color-success);
}

.action-button--success:hover {
  background: color-mix(in srgb, var(--color-success) 90%, black);
}

.action-button--warning {
  background: var(--color-warning);
  color: white;
  border-color: var(--color-warning);
}

.action-button--warning:hover {
  background: color-mix(in srgb, var(--color-warning) 90%, black);
}

.action-button--secondary {
  background: var(--color-primary);
  color: white;
  border-color: var(--color-primary);
}

.action-button--secondary:hover {
  background: color-mix(in srgb, var(--color-primary) 90%, black);
}

.action-button :deep(svg) {
  font-size: 1.1em;
}
</style>
