<template>
  <div class="server-create">
    <Hero title="Create Server" subtitle="Set up a new Minecraft server" />

    <div class="server-create__container">
      <form @submit.prevent="handleSubmit" class="server-create__form">
        <!-- Server Name -->
        <div class="form-group">
          <label for="name" class="form-label">Server Name</label>
          <input
            id="name"
            v-model="form.name"
            type="text"
            class="form-input"
            placeholder="My Minecraft Server"
            required
          />
          <p class="form-help">The display name for your server</p>
        </div>

        <!-- Port -->
        <div class="form-group">
          <label for="port" class="form-label">Port</label>
          <input
            id="port"
            v-model.number="form.port"
            type="number"
            class="form-input"
            placeholder="0 (auto)"
            min="0"
            max="65535"
            required
          />
          <p class="form-help">
            Server port (0 = auto-select free port, default: 25565)
          </p>
        </div>

        <!-- Server Type -->
        <div class="form-group">
          <label for="type" class="form-label">Server Type</label>
          <select id="type" v-model="form.type" class="form-select" required>
            <option value="vanilla">Vanilla</option>
            <option value="paper">Paper</option>
            <option value="fabric">Fabric</option>
          </select>
          <p class="form-help">Choose the server software</p>
        </div>

        <!-- Version -->
        <div class="form-group">
          <label for="version" class="form-label">Version</label>
          <select
            id="version"
            v-model="form.version"
            class="form-select"
            :disabled="versionsLoading"
            required
          >
            <option value="" disabled>
              {{
                versionsLoading ? "Loading versions..." : "Select a version..."
              }}
            </option>
            <option
              v-for="version in availableVersions"
              :key="version"
              :value="version"
            >
              {{ version }}
            </option>
          </select>
          <p class="form-help">
            {{
              versionsLoading
                ? "Loading available versions..."
                : `${form.type} versions`
            }}
          </p>
        </div>

        <!-- Max Memory (MB) -->
        <div class="form-group">
          <label for="memoryMb" class="form-label">Max Memory (MB)</label>
          <input
            id="memoryMb"
            v-model.number="form.memoryMb"
            type="number"
            class="form-input"
            placeholder="2048"
            min="256"
            required
          />
          <p class="form-help">Maximum heap memory allocation in megabytes</p>
        </div>

        <!-- EULA Agreement -->
        <div class="form-group">
          <label class="form-checkbox">
            <input v-model="form.eula" type="checkbox" required />
            <span>I agree to the Minecraft EULA</span>
          </label>
          <p class="form-help">
            You must agree to the Minecraft End User License Agreement
          </p>
        </div>

        <!-- Form Actions -->
        <div class="form-actions">
          <button
            type="button"
            class="btn btn--secondary"
            @click="$router.back()"
            :disabled="isLoading"
          >
            Cancel
          </button>
          <button type="submit" class="btn btn--primary" :disabled="isLoading">
            <span v-if="isLoading">Creating...</span>
            <span v-else>Create Server</span>
          </button>
        </div>

        <!-- Error Message -->
        <div v-if="errorMessage" class="error-message">
          <p>{{ errorMessage }}</p>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from "vue";
import { useRouter } from "vue-router";
import { createServer } from "../helper";
import Hero from "../components/Hero.vue";
import type { CreateServerRequest } from "../types";

const router = useRouter();

const isLoading = ref(false);
const errorMessage = ref("");
const versionsLoading = ref(false);
const availableVersions = ref<string[]>([]);

const form = ref<CreateServerRequest>({
  name: "",
  port: 0,
  type: "vanilla",
  version: "",
  memoryMb: 2048,
  eula: false,
});

// Load versions when component mounts or server type changes
const loadVersions = async () => {
  if (!form.value.type) return;

  versionsLoading.value = true;
  errorMessage.value = "";

  try {
    console.log("[ServerCreate] Fetching versions for:", form.value.type);
    const response = await fetch(
      `http://localhost:8484/versions?type=${form.value.type}`
    );

    if (!response.ok) {
      throw new Error(`HTTP ${response.status}`);
    }

    const data = await response.json();
    availableVersions.value = data.versions || [];

    // Set first version as default
    if (availableVersions.value.length > 0) {
      form.value.version = availableVersions.value[0];
    }

    console.log(
      `[ServerCreate] Loaded ${availableVersions.value.length} versions`
    );
  } catch (error) {
    console.error("[ServerCreate] Failed to fetch versions:", error);
    errorMessage.value = "Failed to load server versions. Please try again.";
  } finally {
    versionsLoading.value = false;
  }
};

// Watch for server type changes and reload versions
watch(
  () => form.value.type,
  () => {
    loadVersions();
  }
);

// Load versions on component mount
onMounted(() => {
  loadVersions();
});

const handleSubmit = async () => {
  if (
    !form.value.name ||
    form.value.port === null ||
    form.value.port === undefined ||
    !form.value.version
  ) {
    errorMessage.value = "Please fill in all required fields";
    return;
  }

  isLoading.value = true;
  errorMessage.value = "";

  try {
    console.log("[ServerCreate] Creating server:", form.value);
    await createServer(form.value);
    console.log("[ServerCreate] Server created successfully");
    // Redirect to servers page
    router.push("/servers");
  } catch (error) {
    console.error("[ServerCreate] Failed to create server:", error);
    errorMessage.value =
      error instanceof Error
        ? error.message
        : "Failed to create server. Please try again.";
  } finally {
    isLoading.value = false;
  }
};
</script>

<style scoped>
.server-create {
  min-height: 100vh;
}

.server-create__container {
  max-width: 600px;
  margin: 0 auto;
  padding: var(--space-xl) var(--section-padding-x);
}

.server-create__form {
  background: var(--color-surface);
  border: var(--border-width-thin) var(--border-style) var(--color-border);
  border-radius: var(--radius-lg);
  padding: var(--space-xl);
}

/* Form Groups */
.form-group {
  margin-bottom: var(--space-lg);
}

.form-label {
  display: block;
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-semi-bold);
  color: var(--color-text);
  margin-bottom: var(--space-sm);
}

.form-input,
.form-select,
.form-textarea {
  width: 100%;
  padding: var(--space-md);
  border: var(--border-width-thin) var(--border-style) var(--color-border);
  border-radius: var(--radius-md);
  background: var(--color-background);
  color: var(--color-text);
  font-size: var(--font-size-sm);
  font-family: inherit;
  transition: all var(--transition-base);
}

.form-input:hover,
.form-select:hover,
.form-textarea:hover {
  border-color: var(--color-primary);
}

.form-input:focus,
.form-select:focus,
.form-textarea:focus {
  outline: none;
  border-color: var(--color-primary);
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.form-textarea {
  resize: vertical;
}

.form-help {
  font-size: var(--font-size-xs);
  color: var(--color-text-secondary);
  margin-top: var(--space-xs);
}

.form-checkbox {
  display: flex;
  align-items: center;
  gap: var(--space-sm);
  cursor: pointer;
  font-size: var(--font-size-sm);
}

.form-checkbox input[type="checkbox"] {
  width: 18px;
  height: 18px;
  cursor: pointer;
  accent-color: var(--color-primary);
}

/* Form Actions */
.form-actions {
  display: flex;
  gap: var(--space-md);
  margin-top: var(--space-xl);
  padding-top: var(--space-lg);
  border-top: var(--border-width-thin) var(--border-style) var(--color-border);
}

.btn {
  flex: 1;
  padding: var(--space-md) var(--space-lg);
  border: none;
  border-radius: var(--radius-md);
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-semi-bold);
  cursor: pointer;
  transition: all var(--transition-base);
  text-transform: none;
}

.btn--primary {
  background: var(--color-primary);
  color: white;
}

.btn--primary:hover:not(:disabled) {
  background: var(--color-secondary);
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
}

.btn--primary:active:not(:disabled) {
  transform: translateY(0);
}

.btn--secondary {
  background: var(--color-surface-alt);
  color: var(--color-text);
  border: var(--border-width-thin) var(--border-style) var(--color-border);
}

.btn--secondary:hover:not(:disabled) {
  background: var(--color-border);
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

/* Error Message */
.error-message {
  margin-top: var(--space-lg);
  padding: var(--space-md);
  background: rgba(239, 68, 68, 0.1);
  border: var(--border-width-thin) var(--border-style) var(--color-error);
  border-radius: var(--radius-md);
  color: var(--color-error);
  font-size: var(--font-size-sm);
}

.error-message p {
  margin: 0;
}
</style>
