import { ref, computed } from "vue";
import type { Theme } from "../types/theme";

// Embedded themes - no need to fetch from files
const EMBEDDED_THEMES: Record<string, Theme> = {
  default: {
    name: "default",
    author: "MCS Manager Team",
    description:
      "Default light/dark theme with vibrant blue and purple accents",
    light: {
      colors: {
        primary: "#667eea",
        secondary: "#764ba2",
        accent: "#f093fb",
        background: "#ffffff",
        surface: "#f8f9ff",
        surfaceAlt: "#f0f1ff",
        text: "#0f0f1e",
        textSecondary: "#6b7280",
        border: "#e5e7eb",
        error: "#ef4444",
        success: "#10b981",
        warning: "#f59e0b",
        info: "#3b82f6",
      },
    },
    dark: {
      colors: {
        primary: "#667eea",
        secondary: "#764ba2",
        accent: "#f093fb",
        background: "#0f0f1e",
        surface: "#1a1a2e",
        surfaceAlt: "#16213e",
        text: "#ffffff",
        textSecondary: "#b0b0b0",
        border: "#2a2a3e",
        error: "#ef4444",
        success: "#10b981",
        warning: "#f59e0b",
        info: "#3b82f6",
      },
    },
  },
  nord: {
    name: "nord",
    author: "MCS Manager Team",
    description: "Arctic, north-bluish color palette",
    light: {
      colors: {
        primary: "#88c0d0",
        secondary: "#81a1c1",
        accent: "#b48ead",
        background: "#eceff4",
        surface: "#e5e9f0",
        surfaceAlt: "#d8dee9",
        text: "#2e3440",
        textSecondary: "#434c5e",
        border: "#d8dee9",
        error: "#bf616a",
        success: "#a3be8c",
        warning: "#ebcb8b",
        info: "#81a1c1",
      },
    },
    dark: {
      colors: {
        primary: "#88c0d0",
        secondary: "#81a1c1",
        accent: "#b48ead",
        background: "#2e3440",
        surface: "#3b4252",
        surfaceAlt: "#434c5e",
        text: "#eceff4",
        textSecondary: "#d8dee9",
        border: "#434c5e",
        error: "#bf616a",
        success: "#a3be8c",
        warning: "#ebcb8b",
        info: "#81a1c1",
      },
    },
  },
};

const currentTheme = ref<Theme | null>(null);
const isDarkMode = ref<boolean>(false);

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
      `Theme applied: ${theme.name} (${darkMode ? "dark" : "light"})`
    );
  }

  /**
   * Switch between light and dark mode
   */
  function toggleDarkMode(darkMode?: boolean) {
    if (!currentTheme.value) {
      console.warn("No theme loaded yet");
      return;
    }

    const newMode = darkMode !== undefined ? darkMode : !isDarkMode.value;
    applyTheme(currentTheme.value, newMode);
  }

  /**
   * Set theme by name
   */
  function setTheme(themeName: string, darkMode: boolean = false) {
    const theme = EMBEDDED_THEMES[themeName];
    if (!theme) {
      console.error(`Theme "${themeName}" not found`);
      return;
    }

    applyTheme(theme, darkMode);
  }

  /**
   * Initialize theme from localStorage or use default
   */
  function initializeTheme() {
    const savedThemeName = localStorage.getItem("theme-name") || "default";
    const savedDarkMode = localStorage.getItem("theme-dark-mode") === "true";

    setTheme(savedThemeName, savedDarkMode);
  }

  /**
   * Get available themes
   */
  function getAvailableThemes(): string[] {
    return Object.keys(EMBEDDED_THEMES);
  }

  /**
   * Get current theme info
   */
  const theme = computed(() => currentTheme.value);
  const isDark = computed(() => isDarkMode.value);

  return {
    theme,
    isDark,
    setTheme,
    toggleDarkMode,
    initializeTheme,
    getAvailableThemes,
  };
}
