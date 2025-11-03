<template>
  <div class="server-detail" v-if="server">
    <Hero
      title="Server"
      :subtitle="server.config.name"
      showBackButton
      backTo="/servers"
    />

    <div class="server-detail__container">
      <div class="detail-grid">
        <div class="info-card">
          <div class="info-card__header">
            <Icon icon="mdi:information" class="info-card__icon" />
            <h3>Server Info</h3>
          </div>
          <div class="info-card__content">
            <div class="info-item">
              <span class="info-label">ID:</span>
              <span class="info-value">{{ server.config.id }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">Type:</span>
              <span class="info-value">{{ server.config.type }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">Version:</span>
              <span class="info-value">{{ server.config.version }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">Port:</span>
              <span class="info-value">{{ server.config.port }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">Memory:</span>
              <span class="info-value">{{ server.config.memoryMb }}MB</span>
            </div>
            <div class="info-item">
              <span class="info-label">Path:</span>
              <span class="info-value">{{ server.config.path }}</span>
            </div>
          </div>
        </div>

        <div class="status-card">
          <div class="status-card__header">
            <Icon icon="mdi:pulse" class="status-card__icon" />
            <h3>Status</h3>
          </div>
          <div class="status-card__content">
            <div class="status-item">
              <span class="status-label">State:</span>
              <span
                class="status-badge"
                :class="`status-badge--${server.state}`"
              >
                <Icon :icon="getStateIcon(server.state)" />
                {{ server.state }}
              </span>
            </div>
            <div class="status-item">
              <span class="status-label">PID:</span>
              <span class="status-value">{{ server.pid || "N/A" }}</span>
            </div>
            <div class="status-item">
              <span class="status-label">Uptime:</span>
              <span class="status-value">{{
                formatUptime(server.uptimeSec)
              }}</span>
            </div>
            <div v-if="server.lastExitErr" class="status-item">
              <span class="status-label">Last Error:</span>
              <span class="status-value error">{{ server.lastExitErr }}</span>
            </div>
          </div>
        </div>
      </div>

      <div class="controls">
        <Button
          v-if="server.state === 'stopped'"
          size="lg"
          icon="mdi:play-circle"
          variant="success"
          @click="handleStart"
        >
          Start Server
        </Button>
        <Button
          v-else
          size="lg"
          icon="mdi:stop-circle"
          variant="warning"
          @click="handleStop"
        >
          Stop Server
        </Button>
        <Button
          v-if="server.state === 'running'"
          size="lg"
          icon="mdi:restart"
          variant="secondary"
          @click="handleRestart"
        >
          Restart Server
        </Button>
        <Button
          size="lg"
          icon="mdi:delete"
          variant="danger"
          :disabled="server.state !== 'stopped'"
          :title="
            server.state !== 'stopped' ? 'Server must be stopped to delete' : ''
          "
          @click="showDeleteModal = true"
        >
          Delete
        </Button>
      </div>

      <!-- Recent Logs -->
      <div class="logs-container">
        <div class="logs-header">
          <Icon icon="mdi:history" class="logs-icon" />
          <h3>Recent Events</h3>
          <Button
            v-if="!showLogs"
            size="sm"
            variant="secondary"
            icon="mdi:eye"
            @click="showLogs = true"
          >
            Show
          </Button>
          <Button
            v-else
            size="sm"
            variant="secondary"
            icon="mdi:eye-off"
            @click="showLogs = false"
          >
            Hide
          </Button>
        </div>
        <div v-if="showLogs" class="logs-content">
          <div v-if="recentLogs.length === 0" class="logs-empty">
            <p>No events yet</p>
          </div>
          <div v-else class="logs-list">
            <div
              v-for="(log, idx) in recentLogs"
              :key="idx"
              class="log-entry"
              :class="`log-${log.type}`"
            >
              <div class="log-timestamp">
                {{ formatLogTime(log.timestamp) }}
              </div>
              <div class="log-type">{{ log.type }}</div>
              <div class="log-message">{{ log.message }}</div>
            </div>
          </div>
        </div>
      </div>

      <ConfirmationModal
        :is-open="showDeleteModal"
        title="Delete Server"
        :message="`Are you sure you want to delete '${server.config.name}'? This action cannot be undone.`"
        confirm-text="Delete"
        cancel-text="Cancel"
        :is-dangerous="true"
        :is-loading="isDeleting"
        @confirm="handleDelete"
        @cancel="showDeleteModal = false"
      />
    </div>
  </div>
  <div v-else class="loading">
    <p>Loading...</p>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from "vue";
import { Icon } from "@iconify/vue";
import { useRoute, useRouter } from "vue-router";
import {
  getServer,
  startServer,
  stopServer,
  restartServer,
  deleteServer,
  formatUptime,
  subscribeToServerEvents,
} from "../helper";
import Button from "../components/Button.vue";
import ConfirmationModal from "../components/ConfirmationModal.vue";
import Hero from "../components/Hero.vue";
import type { ServerInfo } from "../types";

const route = useRoute();
const router = useRouter();
const server = ref<ServerInfo | null>(null);
const serverId = route.params.id as string;
const showDeleteModal = ref(false);
const isDeleting = ref(false);
const showLogs = ref(true);
const recentLogs = ref<
  Array<{ type: string; message: string; timestamp: number }>
>([]);
let eventSource: EventSource | null = null;

const getStateIcon = (state: string) => {
  if (state === "running") return "mdi:play-circle";
  if (state === "stopped") return "mdi:stop-circle";
  return "mdi:pause-circle";
};

const formatLogTime = (timestamp: number): string => {
  const date = new Date(timestamp);
  return date.toLocaleTimeString();
};

const addLog = (type: string, message: string) => {
  recentLogs.value.unshift({
    type,
    message,
    timestamp: Date.now(),
  });
  // Keep only last 50 logs
  if (recentLogs.value.length > 50) {
    recentLogs.value.pop();
  }
};

onMounted(async () => {
  try {
    server.value = await getServer(serverId);
  } catch (error) {
    console.error("Failed to load server:", error);
  }

  console.log("[ServerDetail] Subscribing to events for server:", serverId);

  // Subscribe to server events
  eventSource = subscribeToServerEvents(serverId, {
    onLog: (logEvent: any) => {
      console.log("[ServerDetail] Log event received:", logEvent);
      if (logEvent.data && logEvent.data.stream) {
        addLog(logEvent.data.stream, logEvent.data.line);
      }
    },
    onStateChange: async (type: string, eventData: any) => {
      console.log(
        "[ServerDetail] State change event received:",
        type,
        eventData
      );
      addLog(type, `Server ${type.split(".")[1]}`);
      // Refresh server data
      try {
        server.value = await getServer(serverId);
        console.log("[ServerDetail] Server data refreshed:", server.value);
      } catch (error) {
        console.error("Failed to refresh server:", error);
      }
    },
    onError: (error) => {
      console.error("[ServerDetail] Event source error:", error);
    },
  });

  console.log("[ServerDetail] Event subscription created");
});

onUnmounted(() => {
  if (eventSource) {
    eventSource.close();
  }
});

const handleStart = async () => {
  try {
    await startServer(serverId);
  } catch (error) {
    console.error("Failed to start server:", error);
  }
};

const handleStop = async () => {
  try {
    await stopServer(serverId);
  } catch (error) {
    console.error("Failed to stop server:", error);
  }
};

const handleRestart = async () => {
  try {
    await restartServer(serverId);
  } catch (error) {
    console.error("Failed to restart server:", error);
  }
};

const handleDelete = async () => {
  isDeleting.value = true;
  try {
    await deleteServer(serverId);
    showDeleteModal.value = false;
    router.push("/servers");
  } catch (error) {
    console.error("Failed to delete server:", error);
  } finally {
    isDeleting.value = false;
  }
};
</script>

<style scoped>
.server-detail {
  flex: 1;
}

.server-detail__container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 var(--section-padding-x) var(--section-padding-y);
}

.detail-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
  gap: var(--space-xl);
  margin-bottom: var(--section-gap);
}

.info-card,
.status-card {
  background: var(--color-surface);
  border: var(--border-width-thin) var(--border-style) var(--color-border);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-md);
  transition: all var(--transition-base);
}

.info-card:hover,
.status-card:hover {
  box-shadow: var(--shadow-lg);
  border-color: var(--color-primary);
}

.info-card__header,
.status-card__header {
  display: flex;
  align-items: center;
  gap: var(--space-md);
  padding: var(--space-lg);
  border-bottom: var(--border-width-thin) var(--border-style)
    var(--color-border);
}

.info-card__icon,
.status-card__icon {
  font-size: 1.25rem;
  color: var(--color-primary);
}

.info-card h3,
.status-card h3 {
  margin: 0;
  font-size: var(--font-size-lg);
  font-weight: var(--font-weight-semibold);
  color: var(--color-text);
}

.info-card__content,
.status-card__content {
  padding: var(--space-lg);
  display: flex;
  flex-direction: column;
  gap: var(--space-md);
}

.info-item,
.status-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: var(--space-md);
}

.info-label,
.status-label {
  font-weight: var(--font-weight-medium);
  color: var(--color-text-secondary);
  font-size: var(--font-size-sm);
}

.info-value,
.status-value {
  color: var(--color-text);
  font-family: var(--font-mono);
  font-size: var(--font-size-sm);
}

.status-value.error {
  color: var(--color-error);
}

.status-badge {
  display: inline-flex;
  align-items: center;
  gap: var(--space-xs);
  padding: var(--space-xs) var(--space-md);
  border-radius: var(--radius-md);
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-semibold);
  text-transform: capitalize;
}

.status-badge :deep(svg) {
  font-size: 0.875rem;
}

.status-badge--running {
  background: rgba(16, 185, 129, 0.15);
  color: var(--color-success);
}

.status-badge--stopped {
  background: rgba(239, 68, 68, 0.15);
  color: var(--color-error);
}

.status-badge--starting {
  background: rgba(245, 158, 11, 0.15);
  color: var(--color-warning);
}

.status-badge--crashed {
  background: rgba(239, 68, 68, 0.15);
  color: var(--color-error);
}

.controls {
  display: flex;
  gap: var(--space-md);
  margin-bottom: var(--section-gap);
}

.logs-container {
  background: var(--color-surface);
  border: var(--border-width-thin) var(--border-style) var(--color-border);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-md);
  overflow: hidden;
}

.logs-header {
  display: flex;
  align-items: center;
  gap: var(--space-md);
  padding: var(--space-lg);
  border-bottom: var(--border-width-thin) var(--border-style)
    var(--color-border);
  background: linear-gradient(
    135deg,
    var(--color-background) 0%,
    var(--color-surface) 100%
  );
}

.logs-icon {
  font-size: 1.25rem;
  color: var(--color-primary);
}

.logs-header h3 {
  margin: 0;
  flex: 1;
  font-size: var(--font-size-lg);
  font-weight: var(--font-weight-semibold);
  color: var(--color-text);
}

.logs-content {
  max-height: 400px;
  overflow-y: auto;
  padding: var(--space-md);
}

.logs-empty {
  text-align: center;
  padding: var(--space-xl) var(--space-lg);
  color: var(--color-text-secondary);
  font-size: var(--font-size-sm);
}

.logs-list {
  display: flex;
  flex-direction: column-reverse;
  gap: var(--space-xs);
}

.log-entry {
  display: flex;
  gap: var(--space-sm);
  padding: var(--space-md);
  background: var(--color-background);
  border-radius: var(--radius-md);
  border-left: 3px solid transparent;
  font-family: var(--font-mono);
  font-size: var(--font-size-xs);
  line-height: 1.4;
  transition: all var(--transition-base);
}

.log-entry:hover {
  background: var(--color-hover);
}

.log-entry.log-stdout {
  border-left-color: var(--color-primary);
  color: var(--color-text);
}

.log-entry.log-stderr {
  border-left-color: var(--color-error);
  color: var(--color-error);
}

.log-entry.log-server\.log {
  border-left-color: var(--color-warning);
  color: var(--color-warning);
}

.log-entry.log-server\.started {
  border-left-color: var(--color-success);
  color: var(--color-success);
}

.log-entry.log-server\.stopped,
.log-entry.log-server\.exited {
  border-left-color: var(--color-error);
  color: var(--color-error);
}

.log-entry.log-server\.crashed {
  border-left-color: var(--color-error);
  color: var(--color-error);
  font-weight: var(--font-weight-semibold);
}

.log-timestamp {
  color: var(--color-text-secondary);
  font-size: var(--font-size-xs);
  opacity: 0.7;
  flex-shrink: 0;
}

.log-type {
  color: var(--color-primary);
  font-weight: var(--font-weight-semibold);
  flex-shrink: 0;
  text-transform: uppercase;
  font-size: 0.75rem;
}

.log-message {
  flex: 1;
  word-break: break-word;
  white-space: pre-wrap;
}

.loading {
  text-align: center;
  padding: var(--section-padding-y) var(--section-padding-x);
  color: var(--color-text-secondary);
}
</style>
