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
            <div
              v-if="server.state === 'running' && server.players"
              class="status-item"
            >
              <span class="status-label">Players Online:</span>
              <span class="status-value"
                >{{ server.players.current }}/{{ server.players.max }}</span
              >
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

      <!-- Server Console & Events Container -->
      <div class="console-events-container">
        <div class="console-events-tabs">
          <button
            class="tab-button"
            :class="{ 'tab-button--active': activeTab === 'console' }"
            @click="activeTab = 'console'"
          >
            <Icon icon="mdi:console" />
            Console
          </button>
          <button
            class="tab-button"
            :class="{ 'tab-button--active': activeTab === 'events' }"
            @click="activeTab = 'events'"
          >
            <Icon icon="mdi:history" />
            Recent Events
          </button>
          <button
            class="tab-button"
            :class="{ 'tab-button--active': activeTab === 'config' }"
            @click="activeTab = 'config'"
          >
            <Icon icon="mdi:cog" />
            Server Configuration
          </button>
        </div>

        <div class="console-events-content">
          <!-- Console Tab -->
          <div v-show="activeTab === 'console'" class="tab-pane">
            <div class="console-wrapper">
              <div class="console-output" ref="consoleOutputRef">
                <div v-if="serverLogs.length === 0" class="console-placeholder">
                  <p>Waiting for server output...</p>
                </div>
                <div v-else>
                  <div
                    v-for="(log, idx) in serverLogs"
                    :key="idx"
                    class="console-line"
                    :class="{ 'console-line--stderr': log.stream === 'stderr' }"
                  >
                    <span class="console-time">{{
                      formatLogTime(log.timestamp)
                    }}</span>
                    <span
                      class="console-stream"
                      :class="`console-stream--${log.stream}`"
                      >[{{ log.stream }}]</span
                    >
                    <span class="console-message">{{ log.message }}</span>
                  </div>
                </div>
              </div>
              <CommandConsole
                :server-id="serverId"
                :is-running="server.state === 'running'"
              />
            </div>
          </div>

          <!-- Events Tab -->
          <div v-show="activeTab === 'events'" class="tab-pane">
            <div class="events-list">
              <div v-if="recentLogs.length === 0" class="events-empty">
                <Icon icon="mdi:information-outline" />
                <p>No events yet</p>
              </div>
              <div v-else class="events-scroll">
                <div
                  v-for="(log, idx) in recentLogs"
                  :key="idx"
                  class="event-item"
                  :class="`event-item--${log.type}`"
                >
                  <span class="event-time">{{
                    formatLogTime(log.timestamp)
                  }}</span>
                  <span class="event-type">{{ log.type }}</span>
                  <span class="event-message">{{ log.message }}</span>
                </div>
              </div>
            </div>
          </div>

          <!-- Config Tab -->
          <div v-show="activeTab === 'config'" class="tab-pane">
            <ServerConfig :server-id="serverId" />
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
    <Loader size="48px" color="#ffffff" :thickness="4" :duration="0.8" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed, watch, nextTick } from "vue";
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
  subscribeToServerUpdates,
} from "../helper";
import Button from "../components/Button.vue";
import ConfirmationModal from "../components/ConfirmationModal.vue";
import Hero from "../components/Hero.vue";
import ServerConfig from "../components/ServerConfig.vue";
import CommandConsole from "../components/CommandConsole.vue";
import Loader from "../components/Loader.vue";
import type { ServerInfo } from "../types";

const route = useRoute();
const router = useRouter();
const server = ref<ServerInfo | null>(null);
const serverId = route.params.id as string;
const showDeleteModal = ref(false);
const isDeleting = ref(false);
const activeTab = ref<"console" | "events" | "config">("console");
const consoleOutputRef = ref<HTMLDivElement | null>(null);
const recentLogs = ref<
  Array<{ type: string; message: string; timestamp: number; stream?: string }>
>([]);
let eventSource: EventSource | null = null;
let unsubscribeUpdates: (() => void) | null = null;

// Filter server logs (logs with stream property)
const serverLogs = computed(() => {
  return recentLogs.value.filter((log) => log.stream);
});

// Auto-scroll console to bottom when new logs arrive
watch(
  () => serverLogs.value.length,
  async () => {
    await nextTick();
    if (consoleOutputRef.value) {
      consoleOutputRef.value.scrollTop = consoleOutputRef.value.scrollHeight;
    }
  }
);

const getStateIcon = (state: string) => {
  if (state === "running") return "mdi:play-circle";
  if (state === "stopped") return "mdi:stop-circle";
  return "mdi:pause-circle";
};

const formatLogTime = (timestamp: number): string => {
  const date = new Date(timestamp);
  return date.toLocaleTimeString();
};

const addLog = (type: string, message: string, stream?: string) => {
  recentLogs.value.push({
    type,
    message,
    timestamp: Date.now(),
    stream,
  });
  // Keep only last 50 logs
  if (recentLogs.value.length > 50) {
    recentLogs.value.shift();
  }
};

onMounted(async () => {
  console.log("[ServerDetail] Subscribing to events for server:", serverId);

  // Subscribe to server events for state changes and logs
  eventSource = subscribeToServerEvents(serverId, {
    onLog: (logEvent: any) => {
      console.log("[ServerDetail] Log event received:", logEvent);
      if (logEvent.data && logEvent.data.stream) {
        addLog(logEvent.data.stream, logEvent.data.line, logEvent.data.stream);
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
    onError: (error: Event) => {
      console.error("[ServerDetail] Event source error:", error);
    },
  });

  // Subscribe to real-time server info updates (for uptime, player counts, etc.)
  unsubscribeUpdates = subscribeToServerUpdates(serverId, (updatedInfo) => {
    console.log("[ServerDetail] Received server info update:", updatedInfo);
    if (server.value) {
      // Update all fields with fresh data
      Object.assign(server.value, updatedInfo);
    }
  });

  // Only load if not already set (e.g., when navigating from /servers)
  if (!server.value) {
    try {
      console.log("[ServerDetail] Loading server data on demand...");
      server.value = await getServer(serverId);
    } catch (error) {
      console.error("Failed to load server:", error);
    }
  }

  console.log("[ServerDetail] Event subscription created");
});

onUnmounted(() => {
  if (eventSource) {
    eventSource.close();
  }
  if (unsubscribeUpdates) {
    unsubscribeUpdates();
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

.config-container {
  margin-bottom: var(--section-gap);
  background: var(--color-surface);
  border: var(--border-width-thin) solid var(--color-border);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-md);
  overflow: hidden;
}

.config-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--space-lg);
  background: linear-gradient(
    135deg,
    var(--color-background) 0%,
    var(--color-surface) 100%
  );
  border-bottom: var(--border-width-thin) solid var(--color-border);
  cursor: pointer;
  transition: all var(--transition-base);
  user-select: none;
}

.config-header:hover {
  background: linear-gradient(
    135deg,
    var(--color-surface) 0%,
    var(--color-background) 100%
  );
}

.config-title {
  display: flex;
  align-items: center;
  gap: var(--space-md);
  margin: 0;
  font-size: var(--font-size-lg);
  font-weight: var(--font-weight-semibold);
  color: var(--color-text);
}

.config-toggle-icon {
  font-size: 1.25rem;
  color: var(--color-primary);
  transition: transform var(--transition-base);
}

.config-icon {
  font-size: 1.25rem;
  color: var(--color-primary);
}

.config-content {
  max-height: 2000px;
  overflow: hidden;
  transition: all var(--transition-base);
  padding: var(--space-lg);
}

[v-show="false"] {
  max-height: 0 !important;
  padding: 0 !important;
  overflow: hidden;
}

.console-container {
  margin-bottom: var(--section-gap);
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
.console-events-container {
  margin-bottom: var(--section-gap);
  background: var(--color-surface);
  border: var(--border-width-thin) solid var(--color-border);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-md);
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.console-events-tabs {
  display: flex;
  gap: 0;
  background: linear-gradient(
    135deg,
    var(--color-background) 0%,
    var(--color-surface) 100%
  );
  border-bottom: var(--border-width-thin) solid var(--color-border);
  padding: var(--space-sm);
}

.tab-button {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-sm);
  flex: 1;
  padding: var(--space-md) var(--space-lg);
  background: transparent;
  border: none;
  color: var(--color-text-secondary);
  cursor: pointer;
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-medium);
  transition: all var(--transition-base);
  border-radius: var(--radius-md);
}

.tab-button:hover {
  background: var(--color-border);
  color: var(--color-text);
}

.tab-button--active {
  background: var(--color-primary);
  color: white;
}

.console-events-content {
  display: flex;
  flex-direction: column;
  min-height: 300px;
}

.tab-pane {
  flex: 1;
  display: flex;
  flex-direction: column;
  animation: fadeIn var(--transition-base);
}

.console-wrapper {
  display: flex;
  flex-direction: column;
  flex: 1;
  gap: var(--space-md);
}

.console-output {
  flex: 1;
  background: var(--color-background);
  border-radius: var(--radius-md);
  padding: var(--space-md);
  font-family: var(--font-mono);
  font-size: var(--font-size-sm);
  overflow-y: auto;
  border: 1px solid var(--color-border);
  max-height: 300px;
  line-height: 1.5;
}

.console-placeholder {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: var(--color-text-secondary);
  text-align: center;
}

.console-line {
  display: flex;
  gap: var(--space-sm);
  padding: var(--space-xs) 0;
  color: var(--color-text);
  white-space: pre-wrap;
  word-break: break-word;
}

.console-line--stderr {
  color: rgb(239, 68, 68);
}

.console-time {
  color: var(--color-text-secondary);
  min-width: 70px;
  opacity: 0.6;
}

.console-stream {
  min-width: 65px;
  font-weight: 600;
  text-transform: uppercase;
  font-size: 0.75rem;
}

.console-stream--stdout {
  color: rgb(34, 197, 94);
}

.console-stream--stderr {
  color: rgb(239, 68, 68);
}

.console-message {
  flex: 1;
}

@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

.events-list {
  display: flex;
  flex-direction: column;
  flex: 1;
  padding: var(--space-md);
}

.events-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: var(--color-text-secondary);
  font-size: var(--font-size-sm);
  gap: var(--space-md);
  flex: 1;
}

.events-empty svg {
  font-size: 2rem;
  opacity: 0.5;
}

.events-scroll {
  display: flex;
  flex-direction: column;
  gap: var(--space-xs);
  max-height: 400px;
  overflow-y: auto;
}

.event-item {
  display: grid;
  grid-template-columns: 70px auto 1fr;
  gap: var(--space-md);
  padding: var(--space-md);
  background: var(--color-background);
  border-radius: var(--radius-md);
  border-left: 3px solid transparent;
  font-family: var(--font-mono);
  font-size: var(--font-size-xs);
  line-height: 1.4;
  transition: all var(--transition-base);
}

.event-item:hover {
  background: var(--color-hover);
}

.event-time {
  color: var(--color-text-secondary);
  font-size: var(--font-size-xs);
}

.event-type {
  font-weight: 600;
  padding: 0 var(--space-sm);
  border-radius: var(--radius-sm);
  background: var(--color-border);
  color: var(--color-text);
  text-transform: uppercase;
  font-size: 0.7rem;
  display: flex;
  align-items: center;
  justify-content: center;
  white-space: nowrap;
}

.event-message {
  color: var(--color-text);
  word-break: break-word;
}

.event-item--server\.started {
  border-left-color: var(--color-success);
}

.event-item--server\.started .event-type {
  background: rgba(34, 197, 94, 0.2);
  color: rgb(34, 197, 94);
}

.event-item--server\.stopped,
.event-item--server\.exited {
  border-left-color: var(--color-error);
}

.event-item--server\.stopped .event-type,
.event-item--server\.exited .event-type {
  background: rgba(239, 68, 68, 0.2);
  color: rgb(239, 68, 68);
}

.event-item--server\.crashed {
  border-left-color: var(--color-error);
}

.event-item--server\.crashed .event-type {
  background: rgba(239, 68, 68, 0.2);
  color: rgb(239, 68, 68);
}

.event-item--server\.log {
  border-left-color: var(--color-primary);
}

.event-item--server\.log .event-type {
  background: rgba(59, 130, 246, 0.2);
  color: rgb(59, 130, 246);
}

.loading {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  background: transparent;
}
</style>
