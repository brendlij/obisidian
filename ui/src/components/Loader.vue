<!-- components/Loader.vue -->
<template>
  <span
    class="loader"
    :style="{
      '--size': size,
      '--color': color,
      '--thickness': thickness + 'px',
      '--duration': duration + 's',
      'animation-play-state': paused ? 'paused' : 'running',
    }"
    role="status"
    :aria-label="ariaLabel"
  />
</template>

<script setup lang="ts">
const props = withDefaults(
  defineProps<{
    size?: string;
    color?: string;
    thickness?: number;
    duration?: number;
    paused?: boolean;
    ariaLabel?: string;
  }>(),
  {
    size: "48px",
    color: "currentColor",
    thickness: 5,
    duration: 1,
    paused: false,
    ariaLabel: "Loading…",
  }
);
</script>

<style scoped>
.loader {
  width: var(--size);
  height: var(--size);
  border: var(--thickness) solid var(--color);
  border-bottom-color: transparent;
  border-radius: 50%;
  display: inline-block;
  box-sizing: border-box;
  animation: rotation var(--duration) linear infinite;
}
@keyframes rotation {
  from {
    transform: rotate(0);
  }
  to {
    transform: rotate(360deg);
  }
}
</style>
