<template>
  <div class="server-config">
    <div class="config-header">
      <h3>Server Properties</h3>
      <div class="config-actions-container">
        <Button
          v-if="!isEditing"
          size="sm"
          variant="secondary"
          icon="mdi:pencil"
          @click="isEditing = true"
        >
          Edit
        </Button>
        <div v-else class="config-actions">
          <Button
            size="sm"
            variant="success"
            icon="mdi:check"
            :disabled="isSaving"
            @click="handleSave"
          >
            Save
          </Button>
          <Button
            size="sm"
            variant="secondary"
            icon="mdi:close"
            @click="handleCancel"
          >
            Cancel
          </Button>
        </div>
      </div>
    </div>

    <div v-if="isLoading" class="config-loading">
      <p>Loading configuration...</p>
    </div>

    <div v-else-if="error" class="config-error">
      <p>{{ error }}</p>
      <Button size="sm" variant="secondary" @click="loadProperties">
        Retry
      </Button>
    </div>

    <div v-else class="config-content">
      <!-- Dynamically rendered properties -->
      <div class="properties-container">
        <div
          v-for="key in Object.keys(editingProps)"
          :key="key"
          class="property-row"
        >
          <div class="property-left">
            <span class="property-key">{{ formatKeyName(key) }}</span>
            <span v-if="getPropertyHint(key)" class="property-hint">{{
              getPropertyHint(key)
            }}</span>
          </div>
          <div class="property-right">
            <input
              v-model="editingProps[key]"
              :readonly="!isEditing"
              :type="getPropertyType(key)"
              :placeholder="getPropertyPlaceholder(key)"
              class="property-value"
            />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import Button from "./Button.vue";

interface Props {
  serverId: string;
}

const props = defineProps<Props>();

const isLoading = ref(false);
const isEditing = ref(false);
const isSaving = ref(false);
const error = ref("");
const originalProps = ref<Record<string, string>>({});
const editingProps = ref<Record<string, string>>({});

// Metadata für bestimmte Properties
const propertyMetadata: Record<
  string,
  {
    hint?: string;
    type?: "text" | "number" | "boolean";
    placeholder?: string;
  }
> = {
  motd: { hint: "Message of the Day", placeholder: "A Minecraft Server" },
  difficulty: { hint: "Game difficulty level", type: "text" },
  gamemode: { hint: "Default game mode", type: "text" },
  "max-players": {
    hint: "Maximum number of players",
    type: "number",
    placeholder: "20",
  },
  "level-seed": { hint: "World seed (empty for random)", placeholder: "" },
  pvp: { hint: "Enable/disable PvP", type: "text" },
  "spawn-protection": {
    hint: "Spawn protection radius in blocks",
    type: "number",
  },
  "online-mode": { hint: "Require online account", type: "text" },
  "view-distance": { hint: "Render distance (3-32)", type: "number" },
  "simulation-distance": { hint: "Simulation distance", type: "number" },
  "max-tick-time": { hint: "Max time per tick in ms", type: "number" },
  "allow-nether": { hint: "Enable Nether dimension", type: "text" },
  "allow-flight": { hint: "Allow flight in survival mode", type: "text" },
  "enable-command-block": { hint: "Enable command blocks", type: "text" },
  "enable-query": { hint: "Enable query protocol", type: "text" },
  "enable-rcon": { hint: "Enable remote console", type: "text" },
  "white-list": { hint: "Enable whitelist", type: "text" },
  hardcore: { hint: "Hardcore mode", type: "text" },
  "level-name": { hint: "World name", placeholder: "world" },
  "level-type": { hint: "World generation type", type: "text" },
  "server-port": { hint: "Server port", type: "number" },
  "server-ip": { hint: "Server bind address", placeholder: "" },
  "resource-pack": { hint: "Resource pack URL", placeholder: "" },
  "player-idle-timeout": {
    hint: "Idle timeout in minutes (0 = disabled)",
    type: "number",
  },
};

// Formatiere den Key-Namen für die Anzeige
const formatKeyName = (key: string): string => {
  return key
    .split("-")
    .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
    .join(" ");
};

// Gib den Hint für eine Property zurück
const getPropertyHint = (key: string): string | undefined => {
  return propertyMetadata[key]?.hint;
};

// Gib den Input-Typ zurück
const getPropertyType = (key: string): string => {
  return propertyMetadata[key]?.type || "text";
};

// Gib den Placeholder zurück
const getPropertyPlaceholder = (key: string): string => {
  return propertyMetadata[key]?.placeholder || "";
};

const loadProperties = async () => {
  isLoading.value = true;
  error.value = "";
  try {
    const response = await fetch(
      `http://localhost:8484/servers/${props.serverId}/properties`
    );
    if (!response.ok) {
      throw new Error(`Failed to load properties: ${response.statusText}`);
    }
    const data = await response.json();
    originalProps.value = { ...data };
    editingProps.value = { ...data };
  } catch (err) {
    error.value =
      err instanceof Error ? err.message : "Failed to load properties";
    console.error("Error loading properties:", err);
  } finally {
    isLoading.value = false;
  }
};

const handleSave = async () => {
  isSaving.value = true;
  error.value = "";
  try {
    const response = await fetch(
      `http://localhost:8484/servers/${props.serverId}/properties`,
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(editingProps.value),
      }
    );
    if (!response.ok) {
      throw new Error(`Failed to save properties: ${response.statusText}`);
    }
    originalProps.value = { ...editingProps.value };
    isEditing.value = false;
  } catch (err) {
    error.value =
      err instanceof Error ? err.message : "Failed to save properties";
    console.error("Error saving properties:", err);
  } finally {
    isSaving.value = false;
  }
};

const handleCancel = () => {
  editingProps.value = { ...originalProps.value };
  isEditing.value = false;
};

onMounted(() => {
  loadProperties();
});
</script>

<style scoped>
.server-config {
  background: var(--color-surface);
  border: var(--border-width-thin) var(--border-style) var(--color-border);
  box-shadow: var(--shadow-md);
  overflow: hidden;
}

.config-header {
  display: flex;
  align-items: center;
  gap: var(--space-md);
  padding: var(--space-lg);
  border-bottom: var(--border-width-thin) var(--border-style)
    var(--color-border);
  background: linear-gradient(
    135deg,
    var(--color-background) 0%,
    var(--color-surface) 100%
  );
}

.config-icon {
  font-size: 1.25rem;
  color: var(--color-primary);
}

.config-header h3 {
  margin: 0;
  flex: 1;
  font-size: var(--font-size-lg);
  font-weight: var(--font-weight-semibold);
  color: var(--color-text);
}

.config-actions-container {
  display: flex;
  gap: var(--space-sm);
  align-items: center;
}

.config-actions {
  display: flex;
  gap: var(--space-sm);
  align-items: center;
}

.config-loading,
.config-error {
  padding: var(--space-xl);
  text-align: center;
  color: var(--color-text-secondary);
}

.config-error {
  background: rgba(239, 68, 68, 0.05);
  color: var(--color-error);
  border-top: var(--border-width-thin) var(--border-style) var(--color-error);
}

.config-content {
  padding: var(--space-lg);
}

.properties-container {
  display: flex;
  flex-direction: column;
  gap: var(--space-md);
}

.property-row {
  display: grid;
  grid-template-columns: 220px 1fr;
  gap: var(--space-md);
  align-items: center;
  padding: var(--space-md);
  background: var(--color-background);
  border-radius: var(--radius-md);
  border: var(--border-width-thin) var(--border-style) var(--color-border);
  transition: all var(--transition-base);
}

.property-row:hover {
  border-color: var(--color-primary);
  box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.1);
}

.property-left {
  display: flex;
  flex-direction: column;
  gap: var(--space-xs);
}

.property-key {
  font-size: var(--font-size-sm);
  color: var(--color-text);
  font-weight: var(--font-weight-semibold);
  word-break: break-word;
}

.property-hint {
  font-size: var(--font-size-xs);
  color: var(--color-text-secondary);
  font-weight: var(--font-weight-normal);
}

.property-right {
  display: flex;
  align-items: center;
}

.property-value {
  width: 100%;
  padding: var(--space-md);
  border: var(--border-width-thin) var(--border-style) var(--color-border);
  border-radius: var(--radius-md);
  background: var(--color-surface);
  color: var(--color-text);
  font-family: inherit;
  font-size: var(--font-size-sm);
  transition: all var(--transition-base);
}

.property-value:focus {
  outline: none;
  border-color: var(--color-primary);
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.property-value:disabled,
.property-value[readonly] {
  background: var(--color-hover);
  cursor: not-allowed;
  opacity: 0.7;
}
</style>
