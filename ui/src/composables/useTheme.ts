import { ref, computed } from "vue";
import type { Theme } from "../types/theme";

// Cache for dynamically loaded themes
const themesCache = ref<Record<string, Theme>>({});
const availableThemes = ref<string[]>([]);
const currentTheme = ref<Theme | null>(null);
const isDarkMode = ref<boolean>(false);
const themesLoaded = ref<boolean>(false);

/**
 * Simple TOML parser for theme files
 */
function parseToml(content: string): Record<string, any> {
  const result: Record<string, any> = {};
  const lines = content.split("\n");
  let currentSection: string | null = null;

  for (const line of lines) {
    const trimmed = line.trim();

    // Skip empty lines and comments
    if (!trimmed || trimmed.startsWith("#")) continue;

    // Handle sections [section.subsection]
    const sectionMatch = trimmed.match(/^\[([^\]]+)\]$/);
    if (sectionMatch) {
      currentSection = sectionMatch[1];
      // Initialize section path
      const parts = currentSection.split(".");
      let obj = result;
      for (let i = 0; i < parts.length - 1; i++) {
        if (!obj[parts[i]]) obj[parts[i]] = {};
        obj = obj[parts[i]];
      }
      if (!obj[parts[parts.length - 1]]) {
        obj[parts[parts.length - 1]] = {};
      }
      continue;
    }

    // Handle key-value pairs
    const kvMatch = trimmed.match(/^([^=]+)=\s*"([^"]*)"\s*$/);
    if (kvMatch) {
      const [, key, value] = kvMatch;
      if (currentSection) {
        const parts = currentSection.split(".");
        let obj = result;
        for (const part of parts) {
          if (!obj[part]) obj[part] = {};
          obj = obj[part];
        }
        obj[key.trim()] = value;
      } else {
        result[key.trim()] = value;
      }
    }
  }

  return result;
}

/**
 * Convert parsed TOML to Theme object
 */
function parseTomlToTheme(content: string, filename: string): Theme | null {
  try {
    const parsed = parseToml(content);

    const theme: Theme = {
      name: parsed.name || filename.replace(".toml", ""),
      author: parsed.author || "Unknown",
      description: parsed.description || "",
      light: {
        colors: parsed.light?.colors || {},
      },
      dark: {
        colors: parsed.dark?.colors || {},
      },
    };

    return theme;
  } catch (error) {
    console.error(`[useTheme] Failed to parse TOML ${filename}:`, error);
    return null;
  }
}

/**
 * Load all themes from public/themes folder
 */
async function loadAvailableThemes(): Promise<void> {
  if (themesLoaded.value) return;

  try {
    console.log("[useTheme] Loading themes from public/themes/...");

    // List of theme files to try loading
    const themeFiles = ["default.toml", "nord.toml", "polar-mist.toml"];
    const loaded: string[] = [];

    for (const file of themeFiles) {
      try {
        const response = await fetch(`/themes/${file}`);
        if (!response.ok) {
          console.warn(`[useTheme] Theme file not found: ${file}`);
          continue;
        }

        const content = await response.text();
        const theme = parseTomlToTheme(content, file);

        if (theme) {
          themesCache.value[theme.name] = theme;
          loaded.push(theme.name);
          console.log(`[useTheme] Loaded theme: ${theme.name}`);
        }
      } catch (error) {
        console.error(`[useTheme] Error loading ${file}:`, error);
      }
    }

    availableThemes.value = loaded;
    themesLoaded.value = true;

    console.log(`[useTheme] Loaded ${loaded.length} themes:`, loaded);
  } catch (error) {
    console.error("[useTheme] Failed to load themes:", error);
  }
}

export function useTheme() {
  /**
   * Apply theme colors to CSS variables
   */
  function applyTheme(theme: Theme, darkMode: boolean) {
    const colors = darkMode ? theme.dark.colors : theme.light.colors;
    const root = document.documentElement;

    Object.entries(colors).forEach(([key, value]) => {
      root.style.setProperty(`--color-${key}`, value);
    });

    // Store current theme
    currentTheme.value = theme;
    isDarkMode.value = darkMode;

    // Persist to localStorage
    localStorage.setItem("theme-name", theme.name);
    localStorage.setItem("theme-dark-mode", darkMode.toString());

    console.log(
      `[useTheme] Theme applied: ${theme.name} (${darkMode ? "dark" : "light"})`
    );
  }

  /**
   * Switch between light and dark mode
   */
  function toggleDarkMode(darkMode?: boolean) {
    if (!currentTheme.value) {
      console.warn("[useTheme] No theme loaded yet");
      return;
    }

    const newMode = darkMode !== undefined ? darkMode : !isDarkMode.value;
    applyTheme(currentTheme.value, newMode);
  }

  /**
   * Set theme by name
   */
  function setTheme(themeName: string, darkMode: boolean = false) {
    const theme = themesCache.value[themeName];
    if (!theme) {
      console.error(`[useTheme] Theme not found: ${themeName}`);
      console.error(
        "[useTheme] Available themes:",
        Object.keys(themesCache.value)
      );
      return;
    }

    applyTheme(theme, darkMode);
  }

  /**
   * Initialize theme on app startup
   */
  async function initializeTheme() {
    console.log("[useTheme] Initializing themes...");

    // Load all themes first
    await loadAvailableThemes();

    // Get saved theme from localStorage or use default
    const savedThemeName = localStorage.getItem("theme-name") || "default";
    const savedDarkMode =
      localStorage.getItem("theme-dark-mode") === "true" ||
      window.matchMedia("(prefers-color-scheme: dark)").matches;

    // Apply the theme
    const themeToApply = themesCache.value[savedThemeName];
    if (themeToApply) {
      applyTheme(themeToApply, savedDarkMode);
    } else {
      console.warn(
        `[useTheme] Saved theme not found: ${savedThemeName}, using first available`
      );
      const firstTheme = Object.values(themesCache.value)[0];
      if (firstTheme) {
        applyTheme(firstTheme, savedDarkMode);
      } else {
        console.error("[useTheme] No themes loaded!");
      }
    }
  }

  /**
   * Get list of available themes
   */
  function getAvailableThemes(): string[] {
    return availableThemes.value;
  }

  /**
   * Get theme by name
   */
  function getTheme(themeName: string) {
    return themesCache.value[themeName];
  }

  /**
   * Get current theme info
   */
  const theme = computed(() => currentTheme.value);
  const isDark = computed(() => isDarkMode.value);
  const themes = computed(() => availableThemes.value);

  return {
    theme,
    isDark,
    themes,
    setTheme,
    toggleDarkMode,
    initializeTheme,
    getAvailableThemes,
    getTheme,
  };
}
