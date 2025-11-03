<script setup lang="ts">
import { ref } from "vue";
import { Icon } from "@iconify/vue";
import { executeCommand } from "../helper";

interface Props {
  serverId: string;
  isRunning: boolean;
}

const props = defineProps<Props>();

const commandInput = ref("");
const commandHistory = ref<string[]>([]);
const historyIndex = ref(-1);
const isExecuting = ref(false);
const feedback = ref<{ type: "success" | "error"; message: string } | null>(
  null
);

const executeCmd = async () => {
  const cmd = commandInput.value.trim();
  if (!cmd) return;

  // Add to history
  commandHistory.value.push(cmd);
  historyIndex.value = -1;

  isExecuting.value = true;
  feedback.value = null;

  try {
    await executeCommand(props.serverId, cmd);
    feedback.value = {
      type: "success",
      message: `Command executed: ${cmd}`,
    };
    commandInput.value = "";
  } catch (error) {
    feedback.value = {
      type: "error",
      message: `Error: ${
        error instanceof Error ? error.message : String(error)
      }`,
    };
  } finally {
    isExecuting.value = false;
    // Clear feedback after 3 seconds
    setTimeout(() => {
      feedback.value = null;
    }, 3000);
  }
};

const handleKeyDown = (e: KeyboardEvent) => {
  if (e.key === "Enter") {
    executeCmd();
  } else if (e.key === "ArrowUp") {
    e.preventDefault();
    if (historyIndex.value < commandHistory.value.length - 1) {
      historyIndex.value++;
      commandInput.value =
        commandHistory.value[
          commandHistory.value.length - 1 - historyIndex.value
        ];
    }
  } else if (e.key === "ArrowDown") {
    e.preventDefault();
    if (historyIndex.value > 0) {
      historyIndex.value--;
      commandInput.value =
        commandHistory.value[
          commandHistory.value.length - 1 - historyIndex.value
        ];
    } else if (historyIndex.value === 0) {
      historyIndex.value = -1;
      commandInput.value = "";
    }
  }
};

const clearHistory = () => {
  commandHistory.value = [];
  historyIndex.value = -1;
};
</script>

<template>
  <div class="command-console">
    <div class="console-header">
      <div class="console-title">
        <Icon icon="mdi:console" class="console-icon" />
        <h3>Server Console</h3>
      </div>
      <button
        class="console-clear-btn"
        @click="clearHistory"
        title="Clear history"
      >
        <Icon icon="mdi:trash-can-outline" />
      </button>
    </div>

    <div class="console-content">
      <div class="command-history">
        <div
          v-for="(cmd, idx) in commandHistory"
          :key="idx"
          class="history-item"
        >
          <span class="history-prompt">></span>
          <span class="history-cmd">{{ cmd }}</span>
        </div>
      </div>

      <div v-if="feedback" :class="`feedback feedback--${feedback.type}`">
        {{ feedback.message }}
      </div>
    </div>

    <div class="console-input">
      <span class="input-prompt">></span>
      <input
        v-model="commandInput"
        type="text"
        class="command-input"
        placeholder="Enter command (e.g., say Hello World)"
        :disabled="!props.isRunning || isExecuting"
        @keydown="handleKeyDown"
      />
      <button
        class="execute-btn"
        @click="executeCmd"
        :disabled="!props.isRunning || isExecuting || !commandInput.trim()"
        :title="!props.isRunning ? 'Server must be running' : 'Execute command'"
      >
        <Icon icon="mdi:send" />
      </button>
    </div>

    <div v-if="!props.isRunning" class="console-disabled">
      <Icon icon="mdi:information-outline" />
      <span>Server must be running to execute commands</span>
    </div>
  </div>
</template>

<style scoped>
.command-console {
  display: flex;
  flex-direction: column;
  gap: var(--space-md);
  padding: var(--space-lg);
  background: var(--color-surface);
  border: var(--border-width-thin) solid var(--color-border);
  border-radius: var(--radius-md);
  font-family: "Courier New", monospace;
  font-size: 0.875rem;
}

.console-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-bottom: var(--space-md);
  border-bottom: var(--border-width-thin) solid var(--color-border);
}

.console-title {
  display: flex;
  align-items: center;
  gap: var(--space-sm);
}

.console-icon {
  color: var(--color-primary);
  font-size: 1.25rem;
}

.console-title h3 {
  margin: 0;
  font-size: 1rem;
  font-weight: 600;
}

.console-clear-btn {
  background: transparent;
  border: none;
  color: var(--color-text-secondary);
  cursor: pointer;
  padding: var(--space-xs) var(--space-sm);
  border-radius: var(--radius-sm);
  display: flex;
  align-items: center;
  transition: all var(--transition-base);
}

.console-clear-btn:hover {
  background: var(--color-border);
  color: var(--color-text);
}

.console-content {
  display: flex;
  flex-direction: column;
  gap: var(--space-sm);
  background: var(--color-background);
  border: var(--border-width-thin) solid var(--color-border);
  border-radius: var(--radius-sm);
  padding: var(--space-md);
  max-height: 300px;
  overflow-y: auto;
  min-height: 100px;
}

.command-history {
  display: flex;
  flex-direction: column;
  gap: var(--space-xs);
}

.history-item {
  display: flex;
  gap: var(--space-sm);
  color: var(--color-text-secondary);
}

.history-prompt {
  color: var(--color-primary);
  font-weight: 600;
  min-width: 1.5rem;
}

.history-cmd {
  color: var(--color-text);
  word-break: break-all;
}

.feedback {
  padding: var(--space-sm) var(--space-md);
  border-radius: var(--radius-sm);
  font-size: 0.875rem;
}

.feedback--success {
  background: rgba(34, 197, 94, 0.1);
  color: rgb(34, 197, 94);
}

.feedback--error {
  background: rgba(239, 68, 68, 0.1);
  color: rgb(239, 68, 68);
}

.console-input {
  display: flex;
  align-items: center;
  gap: var(--space-sm);
  background: var(--color-background);
  border: var(--border-width-thin) solid var(--color-border);
  border-radius: var(--radius-sm);
  padding: var(--space-md);
}

.input-prompt {
  color: var(--color-primary);
  font-weight: 600;
  min-width: 1.5rem;
}

.command-input {
  flex: 1;
  background: transparent;
  border: none;
  color: var(--color-text);
  font-family: "Courier New", monospace;
  font-size: 0.875rem;
  outline: none;
}

.command-input::placeholder {
  color: var(--color-text-secondary);
}

.command-input:disabled {
  color: var(--color-text-secondary);
  cursor: not-allowed;
}

.execute-btn {
  background: var(--color-primary);
  color: white;
  border: none;
  padding: var(--space-sm) var(--space-md);
  border-radius: var(--radius-sm);
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: var(--space-xs);
  transition: all var(--transition-base);
}

.execute-btn:hover:not(:disabled) {
  background: var(--color-primary-dark, rgb(59, 130, 246));
  transform: translateY(-2px);
}

.execute-btn:disabled {
  background: var(--color-border);
  cursor: not-allowed;
  opacity: 0.5;
}

.console-disabled {
  display: flex;
  align-items: center;
  gap: var(--space-sm);
  padding: var(--space-md);
  background: rgba(239, 68, 68, 0.1);
  border: var(--border-width-thin) solid rgba(239, 68, 68, 0.3);
  border-radius: var(--radius-sm);
  color: rgb(239, 68, 68);
  font-size: 0.875rem;
}
</style>
