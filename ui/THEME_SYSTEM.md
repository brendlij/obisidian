# Theme System Documentation

## Overview

The theme system allows users to define custom themes via TOML files with light and dark mode support. Themes are loaded dynamically and applied via CSS variables.

## Theme Structure

### TOML Format

```toml
name = "Theme Name"
author = "Your Name"
description = "Theme description"

[light]
[light.colors]
primary = "#667eea"
secondary = "#764ba2"
accent = "#f093fb"
background = "#ffffff"
surface = "#f8f9ff"
surfaceAlt = "#f0f1ff"
text = "#0f0f1e"
textSecondary = "#6b7280"
border = "#e5e7eb"
error = "#ef4444"
success = "#10b981"
warning = "#f59e0b"
info = "#3b82f6"

[dark]
[dark.colors]
primary = "#667eea"
secondary = "#764ba2"
accent = "#f093fb"
background = "#0f0f1e"
surface = "#1a1a2e"
surfaceAlt = "#16213e"
text = "#ffffff"
textSecondary = "#b0b0b0"
border = "#2a2a3e"
error = "#ef4444"
success = "#10b981"
warning = "#f59e0b"
info = "#3b82f6"
```

## Available Colors

| Color           | Purpose               |
| --------------- | --------------------- |
| `primary`       | Main brand color      |
| `secondary`     | Secondary accent      |
| `accent`        | Tertiary highlight    |
| `background`    | Main background       |
| `surface`       | Card/panel background |
| `surfaceAlt`    | Alternative surface   |
| `text`          | Primary text          |
| `textSecondary` | Secondary text        |
| `border`        | Border color          |
| `error`         | Error states          |
| `success`       | Success states        |
| `warning`       | Warning states        |
| `info`          | Info states           |

## Usage in Components

### Using CSS Variables

```css
.my-element {
  background-color: var(--color-primary);
  color: var(--color-text);
  border: 1px solid var(--color-border);
}
```

### Using useTheme Composable

```vue
<script setup lang="ts">
import { useTheme } from "@/composables/useTheme";

const { theme, isDark, toggleDarkMode, setTheme } = useTheme();

// Toggle between light and dark
const handleToggle = () => {
  toggleDarkMode();
};

// Switch to a specific theme
const switchTheme = async () => {
  await setTheme("/themes/my-theme.toml", false);
};
</script>

<template>
  <div>
    <p>Current theme: {{ theme?.name }}</p>
    <p>Dark mode: {{ isDark }}</p>
    <button @click="handleToggle">Toggle Dark Mode</button>
  </div>
</template>
```

## Creating Custom Themes

1. Create a new `.toml` file in `ui/src/themes/`
2. Follow the structure from `schema.toml`
3. Define your color palette for both light and dark modes
4. Place the file in the public/themes directory or import it

## Theme Switcher Component

The `ThemeSwitcher.vue` component provides a UI toggle:

```vue
<ThemeSwitcher />
```

It appears in the navbar by default and allows users to switch between light and dark modes.

## localStorage

Themes are persisted to localStorage:

- `theme-name`: The name of the current theme
- `theme-dark-mode`: Boolean indicating dark mode status

## API Reference

### useTheme()

#### `setTheme(path: string | Theme, darkMode?: boolean): Promise<void>`

Load and apply a theme from file path or Theme object.

#### `toggleDarkMode(darkMode?: boolean): void`

Toggle between light and dark mode, or set to specific mode.

#### `loadThemeFromFile(path: string): Promise<Theme>`

Load a theme from a TOML file.

#### `parseTheme(tomlContent: string): Promise<Theme>`

Parse TOML content string to Theme object.

#### `initializeTheme(): Promise<void>`

Initialize theme from localStorage or use default.

#### `getAvailableThemes(): Promise<string[]>`

Get list of available themes.

#### `theme: computed<Theme | null>`

Current theme object.

#### `isDark: computed<boolean>`

Whether dark mode is enabled.

## Examples

### Example: Create a "Cyberpunk" Theme

```toml
name = "Cyberpunk"
author = "Theme Designer"
description = "Neon cyberpunk theme with bright accents"

[light]
[light.colors]
primary = "#00ff00"
secondary = "#ff00ff"
accent = "#00ffff"
# ... rest of colors

[dark]
[dark.colors]
primary = "#00ff00"
secondary = "#ff00ff"
accent = "#00ffff"
background = "#0a0e27"
# ... rest of colors
```

## Tips

- Keep contrast ratios in mind for accessibility
- Test both light and dark modes
- Use the CSS variable names consistently
- Add your themes to `public/themes/` for production
