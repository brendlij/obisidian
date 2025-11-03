<template>
  <div class="theme-switcher">
    <!-- Theme selector dropdown -->
    <select
      v-model="selectedTheme"
      class="theme-select"
      :title="`Current theme: ${selectedTheme}`"
      @change="changeTheme"
    >
      <option v-for="theme in availableThemes" :key="theme" :value="theme">
        {{ theme }}
      </option>
    </select>

    <!-- Dark mode toggle -->
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
import { ref, onMounted, watch } from "vue";
import { useTheme } from "../composables/useTheme";
import { Icon } from "@iconify/vue";

const { isDark, toggleDarkMode, setTheme, getAvailableThemes, theme } =
  useTheme();

const selectedTheme = ref<string>("default");
const availableThemes = ref<string[]>([]);

onMounted(async () => {
  // Get available themes after they're loaded
  availableThemes.value = getAvailableThemes();

  // Get current theme from composable
  if (theme.value) {
    selectedTheme.value = theme.value.name;
  }

  console.log("[ThemeSwitcher] Available themes:", availableThemes.value);
});

// Watch for theme changes from outside (e.g., from settings)
watch(theme, (newTheme) => {
  if (newTheme) {
    selectedTheme.value = newTheme.name;
  }
});

const toggleTheme = () => {
  toggleDarkMode();
};

const changeTheme = () => {
  setTheme(selectedTheme.value, isDark.value);
};
</script>

<style scoped>
.theme-switcher {
  display: flex;
  align-items: center;
  gap: var(--space-sm);
}

.theme-select {
  border: var(--border-width-thin) var(--border-style) var(--color-border);
  border-radius: var(--radius-md);
  padding: var(--space-sm) var(--space-md);
  cursor: pointer;
  color: var(--color-text);
  background: var(--color-surface);
  font-size: var(--font-size-sm);
  transition: all var(--transition-base);
  text-transform: capitalize;
}

.theme-select:hover {
  border-color: var(--color-primary);
  background: var(--color-surface-alt);
}

.theme-select:focus {
  outline: none;
  border-color: var(--color-primary);
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.theme-toggle {
  border: var(--border-width-thin) var(--border-style) var(--color-border);
  border-radius: var(--radius-md);
  padding: var(--space-sm);
  cursor: pointer;
  color: var(--color-text);
  background: var(--color-surface);
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
