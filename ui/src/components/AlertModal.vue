<script setup lang="ts">
import { Icon } from "@iconify/vue";
import Button from "./Button.vue";

export interface Props {
  title: string;
  message: string;
  type?: "info" | "success" | "warning" | "error";
  icon?: string;
  isOpen: boolean;
}

withDefaults(defineProps<Props>(), {
  type: "info",
});

const emit = defineEmits<{
  close: [];
}>();

const getTypeIcon = (type: string) => {
  switch (type) {
    case "success":
      return "mdi:check-circle";
    case "warning":
      return "mdi:alert-circle";
    case "error":
      return "mdi:close-circle";
    default:
      return "mdi:information";
  }
};

const getTypeColor = (type: string) => {
  switch (type) {
    case "success":
      return "success";
    case "warning":
      return "warning";
    case "error":
      return "error";
    default:
      return "info";
  }
};
</script>

<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="isOpen" class="modal-overlay" @click="emit('close')">
        <div class="modal-content" @click.stop>
          <div class="modal-header">
            <div
              class="modal-icon"
              :class="`modal-icon--${getTypeColor(type)}`"
            >
              <Icon :icon="icon || getTypeIcon(type)" />
            </div>
            <button class="modal-close" @click="emit('close')">
              <Icon icon="mdi:close" />
            </button>
          </div>

          <div class="modal-body">
            <h2 class="modal-title">{{ title }}</h2>
            <p class="modal-message">{{ message }}</p>
          </div>

          <div class="modal-footer">
            <Button variant="primary" size="md" @click="emit('close')">
              Close
            </Button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: var(--z-modal);
}

.modal-content {
  background: var(--color-surface);
  border: var(--border-width-thin) var(--border-style) var(--color-border);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-xl);
  max-width: 400px;
  width: 90%;
  animation: slideUp var(--transition-base) ease-out;
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--space-xl);
  border-bottom: var(--border-width-thin) var(--border-style)
    var(--color-border);
}

.modal-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 2.5rem;
  height: 2.5rem;
  border-radius: var(--radius-lg);
  font-size: 1.25rem;
}

.modal-icon--info {
  background: rgba(59, 130, 246, 0.1);
  color: var(--color-info);
}

.modal-icon--success {
  background: rgba(16, 185, 129, 0.1);
  color: var(--color-success);
}

.modal-icon--warning {
  background: rgba(245, 158, 11, 0.1);
  color: var(--color-warning);
}

.modal-icon--error {
  background: rgba(239, 68, 68, 0.1);
  color: var(--color-error);
}

.modal-close {
  background: none;
  border: none;
  color: var(--color-text-secondary);
  cursor: pointer;
  font-size: 1.25rem;
  padding: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: color var(--transition-base);
}

.modal-close:hover {
  color: var(--color-text);
}

.modal-body {
  padding: var(--space-xl);
}

.modal-title {
  font-size: var(--font-size-lg);
  font-weight: var(--font-weight-semibold);
  margin: 0 0 var(--space-md) 0;
}

.modal-message {
  font-size: var(--font-size-base);
  color: var(--color-text-secondary);
  margin: 0;
  line-height: 1.6;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: var(--space-md);
  padding: var(--space-xl);
  border-top: var(--border-width-thin) var(--border-style) var(--color-border);
}

.modal-enter-active,
.modal-leave-active {
  transition: opacity var(--transition-base);
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

@keyframes slideUp {
  from {
    transform: translateY(20px);
    opacity: 0;
  }
  to {
    transform: translateY(0);
    opacity: 1;
  }
}
</style>
