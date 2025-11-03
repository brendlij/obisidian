<script setup lang="ts">
import { Icon } from "@iconify/vue";
import ThemeSwitcher from "./ThemeSwitcher.vue";
</script>

<template>
  <header class="header">
    <div class="header__content">
      <router-link to="/" class="header__logo">
        <Icon icon="hugeicons:server-stack-03" class="header__logo-icon" />
        <span>Obsidian</span>
      </router-link>

      <nav class="header__nav">
        <router-link
          to="/"
          :class="{ active: $route.path === '/' }"
          class="header__nav-link"
        >
          <Icon icon="mdi:home" />
          <span>Dashboard</span>
        </router-link>
        <router-link
          to="/servers"
          :class="{ active: $route.path.startsWith('/servers') }"
          class="header__nav-link"
        >
          <Icon icon="hugeicons:server-stack-03" />
          <span>Servers</span>
        </router-link>
        <router-link
          to="/settings"
          :class="{ active: $route.path === '/settings' }"
          class="header__nav-link"
        >
          <Icon icon="mdi:cog" />
          <span>Settings</span>
        </router-link>
      </nav>

      <ThemeSwitcher />
    </div>
  </header>
</template>

<style scoped>
.header {
  background: var(--color-surface);
  border-bottom: var(--border-width-thin) var(--border-style)
    var(--color-border);
  position: sticky;
  top: 0;
  z-index: var(--z-sticky);
}

.header__content {
  max-width: 1400px;
  margin: 0 auto;
  padding: var(--space-lg) var(--section-padding-x);
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: var(--space-xl);
}

.header__logo {
  display: flex;
  align-items: center;
  gap: var(--space-md);
  font-size: var(--font-size-lg);
  font-weight: var(--font-weight-bold);
  color: var(--color-primary);
  text-decoration: none;
  transition: color var(--transition-base);
  white-space: nowrap;
}

.header__logo:hover {
  color: var(--color-secondary);
}

.header__logo-icon {
  font-size: 1.5rem;
}

.header__nav {
  display: flex;
  gap: var(--space-lg);
  flex: 1;
  margin-left: var(--space-xl);
}

.header__nav-link {
  display: flex;
  align-items: center;
  gap: var(--space-sm);
  color: var(--color-text-secondary);
  text-decoration: none;
  font-size: var(--font-size-sm);
  transition: color var(--transition-base);
  padding: var(--space-sm) var(--space-md);
  border-radius: var(--radius-md);
  white-space: nowrap;
}

.header__nav-link:hover {
  color: var(--color-text);
}

.header__nav-link.active {
  color: var(--color-primary);
  background: rgba(102, 126, 234, 0.1);
}

.header__nav-link :deep(svg) {
  font-size: 1rem;
}
</style>
