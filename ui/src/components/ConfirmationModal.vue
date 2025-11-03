<script setup lang="ts">
import { Icon } from "@iconify/vue";
import Button from "./Button.vue";

export interface Props {
  title: string;
  message: string;
  confirmText?: string;
  cancelText?: string;
  isDangerous?: boolean;
  isOpen: boolean;
  isLoading?: boolean;
}

withDefaults(defineProps<Props>(), {
  confirmText: "Confirm",
  cancelText: "Cancel",
  isDangerous: false,
  isLoading: false,
});

const emit = defineEmits<{
  confirm: [];
  cancel: [];
}>();
</script>

<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="isOpen" class="modal-overlay" @click="emit('cancel')">
        <div class="modal-content" @click.stop>
          <div class="modal-header">
            <Icon
              icon="mdi:help-circle"
              class="modal-icon"
              :class="{ 'modal-icon--danger': isDangerous }"
            />
            <button class="modal-close" @click="emit('cancel')">
              <Icon icon="mdi:close" />
            </button>
          </div>

          <div class="modal-body">
            <h2 class="modal-title">{{ title }}</h2>
            <p class="modal-message">{{ message }}</p>
          </div>

          <div class="modal-footer">
            <Button
              variant="secondary"
              size="md"
              @click="emit('cancel')"
              :disabled="isLoading"
            >
              {{ cancelText }}
            </Button>
            <Button
              :variant="isDangerous ? 'danger' : 'primary'"
              size="md"
              :loading="isLoading"
              @click="emit('confirm')"
            >
              {{ confirmText }}
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
  font-size: 1.5rem;
  color: var(--color-info);
  transition: color var(--transition-base);
}

.modal-icon--danger {
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
