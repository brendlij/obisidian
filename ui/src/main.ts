import { createApp } from "vue";
import router from "./router";
import "./style.css";
import App from "./App.vue";
import { useTheme } from "./composables/useTheme";

const app = createApp(App);

// Initialize theme before mounting
const { initializeTheme } = useTheme();
initializeTheme().then(() => {
  console.log("[main] Theme initialized");
  app.use(router);
  app.mount("#app");
});
