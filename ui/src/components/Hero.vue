<script setup lang="ts">
import Button from "./Button.vue";
import { Icon } from "@iconify/vue";

export interface Props {
  title: string;
  subtitle: string;
  showActionButton?: boolean;
  actionButtonText?: string;
  actionButtonIcon?: string;
  actionButtonTo?: string;
  showBackButton?: boolean;
  backTo?: string;
}

withDefaults(defineProps<Props>(), {
  showActionButton: false,
  actionButtonText: "Manage Servers",
  actionButtonIcon: "hugeicons:server-stack-03",
  actionButtonTo: "/servers",
  showBackButton: false,
  backTo: "/",
});
</script>

<template>
  <div class="hero">
    <div class="hero__content">
      <router-link v-if="showBackButton" :to="backTo" class="hero__back">
        <Icon icon="mdi:arrow-left" />
        Back
      </router-link>
      <h1 class="hero__title">{{ title }}</h1>
      <p class="hero__subtitle">{{ subtitle }}</p>
    </div>
    <router-link
      v-if="showActionButton"
      :to="actionButtonTo"
      class="hero__action"
    >
      <Button :icon="actionButtonIcon" size="lg" variant="primary">
        {{ actionButtonText }}
      </Button>
    </router-link>
  </div>
</template>

<style scoped>
.hero {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: var(--space-2xl);
  padding: var(--section-padding-y) var(--section-padding-x);
  margin-bottom: var(--section-gap);
  border-bottom: var(--border-width-thin) var(--border-style)
    var(--color-border);
}

.hero__content {
  flex: 1;
}

.hero__back {
  display: inline-flex;
  align-items: center;
  gap: var(--space-sm);
  color: var(--color-primary);
  text-decoration: none;
  font-size: var(--font-size-sm);
  transition: color var(--transition-base);
  margin-bottom: var(--space-md);
}

.hero__back:hover {
  color: var(--color-secondary);
}

.hero__back :deep(svg) {
  font-size: 1rem;
}

.hero__action {
  text-decoration: none;
  flex-shrink: 0;
}

.hero__title {
  font-size: var(--font-size-3xl);
  font-weight: var(--font-weight-bold);
  margin: 0 0 var(--space-sm) 0;
  color: var(--color-text);
}

.hero__subtitle {
  font-size: var(--font-size-base);
  color: var(--color-text-secondary);
  margin: 0;
}
</style>
