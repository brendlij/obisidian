<template>
  <div class="settings">
    <Hero title="Settings" subtitle="Manage your application preferences" />

    <div class="settings__container">
      <div class="settings__grid">
        <div class="settings__card">
          <div class="settings__card-header">
            <Icon icon="mdi:palette" class="settings__card-icon" />
            <h2 class="settings__card-title">Theme</h2>
          </div>
          <div class="settings__section">
            <label class="settings__label">Select Theme</label>
            <select
              v-model="selectedTheme"
              @change="switchTheme"
              class="settings__select"
            >
              <option
                v-for="theme in availableThemes"
                :key="theme"
                :value="theme"
              >
                {{ theme.charAt(0).toUpperCase() + theme.slice(1) }}
              </option>
            </select>
          </div>
        </div>

        <div class="settings__card">
          <div class="settings__card-header">
            <Icon icon="mdi:api" class="settings__card-icon" />
            <h2 class="settings__card-title">API Configuration</h2>
          </div>
          <div class="settings__section">
            <label class="settings__label">Backend URL</label>
            <input
              v-model="apiUrl"
              type="text"
              placeholder="http://localhost:8484"
              class="settings__input"
            />
          </div>
          <div class="settings__actions">
            <Button
              size="md"
              icon="mdi:content-save"
              variant="primary"
              @click="saveSettings"
            >
              Save Configuration
            </Button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { Icon } from "@iconify/vue";
import { useTheme } from "../composables/useTheme";
import Button from "../components/Button.vue";
import Hero from "../components/Hero.vue";

const apiUrl = ref(localStorage.getItem("apiUrl") || "http://localhost:8484");
const selectedTheme = ref(localStorage.getItem("theme-name") || "default");
const { setTheme, isDark, getAvailableThemes } = useTheme();
const availableThemes = ref<string[]>([]);

onMounted(() => {
  availableThemes.value = getAvailableThemes();
});

const saveSettings = () => {
  localStorage.setItem("apiUrl", apiUrl.value);
  alert("Settings saved!");
};

const switchTheme = () => {
  setTheme(selectedTheme.value, isDark.value);
};
</script>

<style scoped>
.settings {
  flex: 1;
}

.settings__container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 var(--section-padding-x) var(--section-padding-y);
}

.settings__grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
  gap: var(--space-xl);
}

.settings__card {
  background: var(--color-surface);
  border: var(--border-width-thin) var(--border-style) var(--color-border);
  border-radius: var(--radius-lg);
  padding: var(--space-xl);
  box-shadow: var(--shadow-md);
  transition: all var(--transition-base);
}

.settings__card:hover {
  box-shadow: var(--shadow-lg);
  border-color: var(--color-primary);
}

.settings__card-header {
  display: flex;
  align-items: center;
  gap: var(--space-md);
  margin-bottom: var(--space-lg);
}

.settings__card-icon {
  font-size: 1.5rem;
  color: var(--color-primary);
}

.settings__card-title {
  font-size: var(--font-size-lg);
  font-weight: var(--font-weight-semibold);
  margin: 0;
  color: var(--color-text);
}

.settings__section {
  margin-bottom: var(--space-lg);
}

.settings__label {
  display: block;
  margin-bottom: var(--space-sm);
  font-weight: var(--font-weight-medium);
  color: var(--color-text);
  font-size: var(--font-size-sm);
}

.settings__input,
.settings__select {
  width: 100%;
  padding: var(--space-md);
  background: var(--color-surface-alt);
  border: var(--border-width-thin) var(--border-style) var(--color-border);
  border-radius: var(--radius-md);
  color: var(--color-text);
  font-size: var(--font-size-base);
  font-family: var(--font-family);
  transition: all var(--transition-base);
}

.settings__input:focus,
.settings__select:focus {
  outline: none;
  border-color: var(--color-primary);
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.settings__actions {
  display: flex;
  gap: var(--space-md);
  padding-top: var(--space-lg);
  border-top: var(--border-width-thin) var(--border-style) var(--color-border);
}
</style>
