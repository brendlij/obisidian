export interface ThemeColors {
  primary: string;
  secondary: string;
  accent: string;
  background: string;
  surface: string;
  surfaceAlt: string;
  text: string;
  textSecondary: string;
  border: string;
  error: string;
  success: string;
  warning: string;
  info: string;
}

export interface ThemeMode {
  colors: ThemeColors;
}

export interface Theme {
  name: string;
  author: string;
  description: string;
  light: ThemeMode;
  dark: ThemeMode;
}
