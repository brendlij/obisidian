<template>
  <div class="theme-switcher">
    <button
      class="theme-toggle"
      :title="`Switch to ${isDark ? 'light' : 'dark'} mode`"
      @click="toggleTheme"
    >
      <Icon
        :icon="isDark ? 'hugeicons:sun-03' : 'hugeicons:moon-02'"
        width="24"
      />
    </button>
  </div>
</template>

<script setup lang="ts">
import { useTheme } from "../composables/useTheme";
import { Icon } from "@iconify/vue";

const { isDark, toggleDarkMode } = useTheme();

const toggleTheme = () => {
  toggleDarkMode();
};
</script>

<style scoped>
.theme-switcher {
  display: flex;
  align-items: center;
}

.theme-toggle {
  border: var(--border-width-thin) var(--border-style) var(--color-border);
  border-radius: var(--radius-md);
  padding: var(--space-sm);
  cursor: pointer;
  color: var(--color-text);
  transition: all var(--transition-base);
  display: flex;
  align-items: center;
  justify-content: center;
}

.theme-toggle:hover {
  background: var(--color-primary);
  color: white;
  border-color: var(--color-primary);
  transform: scale(1.05);
}

.theme-toggle:active {
  transform: scale(0.95);
}
</style>
