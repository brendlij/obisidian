<script setup lang="ts">
import { Icon } from "@iconify/vue";

export interface Props {
  icon?: string;
  title: string;
  value: string | number;
  color?: "primary" | "success" | "error" | "warning" | "info";
}

withDefaults(defineProps<Props>(), {
  color: "primary",
});
</script>

<template>
  <div class="stat-card" :class="`stat-card--${color}`">
    <div v-if="icon" class="stat-card__icon">
      <Icon :icon="icon" />
    </div>
    <div class="stat-card__content">
      <h3 class="stat-card__title">{{ title }}</h3>
      <p class="stat-card__value">{{ value }}</p>
    </div>
  </div>
</template>

<style scoped>
.stat-card {
  display: flex;
  align-items: center;
  gap: var(--space-lg);
  padding: var(--space-xl);
  background: var(--color-surface);
  border: var(--border-width-thin) var(--border-style) var(--color-border);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-md);
  transition: all var(--transition-base);
}

.stat-card:hover {
  box-shadow: var(--shadow-lg);
  border-color: var(--color-primary);
  transform: translateY(-2px);
}

.stat-card__icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 3rem;
  height: 3rem;
  border-radius: var(--radius-lg);
  font-size: 1.5rem;
  background: var(--color-surface-alt);
}

.stat-card--primary .stat-card__icon {
  background: rgba(102, 126, 234, 0.2);
  color: var(--color-primary);
}

.stat-card--success .stat-card__icon {
  background: rgba(16, 185, 129, 0.2);
  color: var(--color-success);
}

.stat-card--error .stat-card__icon {
  background: rgba(239, 68, 68, 0.2);
  color: var(--color-error);
}

.stat-card--warning .stat-card__icon {
  background: rgba(245, 158, 11, 0.2);
  color: var(--color-warning);
}

.stat-card--info .stat-card__icon {
  background: rgba(59, 130, 246, 0.2);
  color: var(--color-info);
}

.stat-card__content {
  flex: 1;
}

.stat-card__title {
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-medium);
  color: var(--color-text-secondary);
  text-transform: uppercase;
  letter-spacing: 0.05em;
  margin-bottom: var(--space-xs);
}

.stat-card__value {
  font-size: var(--font-size-2xl);
  font-weight: var(--font-weight-bold);
  margin: 0;
}

.stat-card--primary .stat-card__value {
  color: var(--color-primary);
}

.stat-card--success .stat-card__value {
  color: var(--color-success);
}

.stat-card--error .stat-card__value {
  color: var(--color-error);
}

.stat-card--warning .stat-card__value {
  color: var(--color-warning);
}

.stat-card--info .stat-card__value {
  color: var(--color-info);
}
</style>
