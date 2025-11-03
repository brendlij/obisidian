<script setup lang="ts">
import { Icon } from "@iconify/vue";

export interface Props {
  title: string;
  subtitle: string;
  showBackButton?: boolean;
  backTo?: string;
}

withDefaults(defineProps<Props>(), {
  showBackButton: false,
  backTo: "/",
});
</script>

<template>
  <div class="page-hero">
    <div class="page-hero__content">
      <router-link v-if="showBackButton" :to="backTo" class="page-hero__back">
        <Icon icon="mdi:arrow-left" />
        Back
      </router-link>
      <h1 class="page-hero__title">{{ title }}</h1>
      <p class="page-hero__subtitle">{{ subtitle }}</p>
    </div>
  </div>
</template>

<style scoped>
.page-hero {
  padding: var(--section-padding-y) var(--section-padding-x);
  margin-bottom: var(--section-gap);
  border-bottom: var(--border-width-thin) var(--border-style)
    var(--color-border);
}

.page-hero__content {
  max-width: 1400px;
  margin: 0 auto;
}

.page-hero__back {
  display: inline-flex;
  align-items: center;
  gap: var(--space-sm);
  color: var(--color-primary);
  text-decoration: none;
  font-size: var(--font-size-sm);
  transition: color var(--transition-base);
  margin-bottom: var(--space-md);
}

.page-hero__back:hover {
  color: var(--color-secondary);
}

.page-hero__back :deep(svg) {
  font-size: 1rem;
}

.page-hero__title {
  font-size: var(--font-size-3xl);
  font-weight: var(--font-weight-bold);
  margin: 0 0 var(--space-sm) 0;
  color: var(--color-text);
}

.page-hero__subtitle {
  font-size: var(--font-size-base);
  color: var(--color-text-secondary);
  margin: 0;
}
</style>
