<script setup lang="ts">
import { computed } from "vue";
import { Icon } from "@iconify/vue";

export interface Props {
  variant?:
    | "primary"
    | "secondary"
    | "danger"
    | "success"
    | "warning"
    | "ghost";
  size?: "sm" | "md" | "lg";
  icon?: string;
  disabled?: boolean;
  loading?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  variant: "primary",
  size: "md",
  disabled: false,
  loading: false,
});

defineEmits<{
  click: [];
}>();

const buttonClass = computed(() => [
  "btn",
  `btn--${props.variant}`,
  `btn--${props.size}`,
  props.disabled && "btn--disabled",
  props.loading && "btn--loading",
]);
</script>

<template>
  <button
    :class="buttonClass"
    :disabled="disabled || loading"
    @click="$emit('click')"
  >
    <Icon v-if="icon && !loading" :icon="icon" />
    <Icon v-if="loading" icon="mdi:loading" class="btn-loading-icon" />
    <span><slot /></span>
  </button>
</template>

<style scoped>
.btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-sm);
  font-weight: var(--font-weight-semibold);
  border-radius: var(--radius-lg);
  border: var(--border-width-thin) var(--border-style) transparent;
  transition: all var(--transition-base);
  cursor: pointer;
  white-space: nowrap;
  user-select: none;
}

/* Sizes */
.btn--sm {
  padding: var(--space-sm) var(--space-md);
  font-size: var(--font-size-sm);
}

.btn--md {
  padding: var(--space-md) var(--space-lg);
  font-size: var(--font-size-base);
}

.btn--lg {
  padding: var(--space-lg) var(--space-xl);
  font-size: var(--font-size-lg);
}

/* Variants */
.btn--primary {
  background: var(--color-primary);
  color: white;
  box-shadow: var(--shadow-md);
}

.btn--primary:hover:not(.btn--disabled) {
  background: var(--color-secondary);
  box-shadow: var(--shadow-hover);
  transform: translateY(-2px);
}

.btn--primary:active:not(.btn--disabled) {
  transform: translateY(0);
}

.btn--secondary {
  background: var(--color-surface);
  color: var(--color-text);
  border-color: var(--color-border);
}

.btn--secondary:hover:not(.btn--disabled) {
  background: var(--color-surface-alt);
  border-color: var(--color-primary);
}

.btn--danger {
  background: var(--color-error);
  color: white;
}

.btn--danger:hover:not(.btn--disabled) {
  background: #dc2626;
}

.btn--success {
  background: var(--color-success);
  color: white;
}

.btn--success:hover:not(.btn--disabled) {
  background: #059669;
}

.btn--warning {
  background: var(--color-warning);
  color: white;
}

.btn--warning:hover:not(.btn--disabled) {
  background: #d97706;
}

.btn--ghost {
  background: transparent;
  color: var(--color-text);
  border-color: var(--color-border);
}

.btn--ghost:hover:not(.btn--disabled) {
  background: var(--color-surface);
  border-color: var(--color-primary);
}

/* States */
.btn--disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn--loading {
  pointer-events: none;
}

.btn-loading-icon {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}
</style>
