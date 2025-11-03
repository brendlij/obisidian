<script setup lang="ts">
import { Icon } from "@iconify/vue";
import { computed } from "vue";

export interface Props {
  size?: "sm" | "md" | "lg" | "xl" | "xxl";
  icon?: string;
  iconPosition?: "left" | "top" | "right" | "bottom";
  label?: string;
  secondaryLabel?: string;
  variant?: "primary" | "secondary" | "danger" | "ghost" | "accent";
  disabled?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  size: "md",
  variant: "primary",
  disabled: false,
  iconPosition: "left",
});

defineEmits<{
  click: [];
}>();

const sizeConfig = {
  sm: {
    classes: "px-3 py-1.5 text-xs font-semibold",
    iconSize: "16px",
    gap: "gap-1.5",
  },
  md: {
    classes: "px-4 py-2 text-sm font-semibold",
    iconSize: "20px",
    gap: "gap-2",
  },
  lg: {
    classes: "px-6 py-2.5 text-base font-semibold",
    iconSize: "24px",
    gap: "gap-2.5",
  },
  xl: {
    classes: "px-8 py-3 text-lg font-semibold",
    iconSize: "28px",
    gap: "gap-3",
  },
  xxl: {
    classes: "px-10 py-4 text-xl font-semibold",
    iconSize: "32px",
    gap: "gap-3.5",
  },
};

const variantClasses = {
  primary:
    "bg-gradient-to-br from-blue-600 to-blue-700 text-white shadow-lg hover:shadow-xl hover:from-blue-700 hover:to-blue-800 active:shadow-md",
  secondary:
    "bg-gradient-to-br from-gray-600 to-gray-700 text-white shadow-lg hover:shadow-xl hover:from-gray-700 hover:to-gray-800 active:shadow-md",
  danger:
    "bg-gradient-to-br from-red-600 to-red-700 text-white shadow-lg hover:shadow-xl hover:from-red-700 hover:to-red-800 active:shadow-md",
  accent:
    "bg-gradient-to-br from-purple-600 to-pink-600 text-white shadow-lg hover:shadow-xl hover:from-purple-700 hover:to-pink-700 active:shadow-md",
  ghost:
    "bg-transparent text-gray-700 border border-gray-300 hover:bg-gray-50 active:bg-gray-100",
};

const currentSize = computed(() => sizeConfig[props.size]);

const buttonClasses = computed(() => [
  "inline-flex",
  props.iconPosition === "left" || props.iconPosition === "right"
    ? "flex-row"
    : "flex-col",
  props.iconPosition === "right" ? "flex-row-reverse" : "",
  "items-center justify-center rounded-xl font-semibold",
  "transition-all duration-200 ease-out",
  "focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500",
  "active:scale-95",
  currentSize.value.classes,
  currentSize.value.gap,
  variantClasses[props.variant],
  props.disabled &&
    "opacity-50 cursor-not-allowed hover:shadow-lg active:scale-100",
]);
</script>

<template>
  <button :class="buttonClasses" :disabled="disabled" @click="$emit('click')">
    <Icon v-if="icon" :icon="icon" :width="currentSize.iconSize" />
    <div
      v-if="label || secondaryLabel"
      :class="[
        'flex',
        props.iconPosition === 'left' || props.iconPosition === 'right'
          ? 'flex-col'
          : 'flex-col',
        'items-center justify-center',
      ]"
    >
      <span v-if="label" class="font-semibold leading-tight block">{{
        label
      }}</span>
      <span
        v-if="secondaryLabel"
        class="text-xs opacity-60 font-normal leading-tight block"
        >{{ secondaryLabel }}</span
      >
    </div>
  </button>
</template>

<style scoped>
button {
  user-select: none;
  -webkit-user-select: none;
  -webkit-appearance: none;
  appearance: none;
  border: none;
  cursor: pointer;
  outline: none;
}

button:disabled {
  cursor: not-allowed;
}

button:active {
  transform: scale(0.95);
}

/* Smooth gradient animation on hover */
button {
  position: relative;
  overflow: hidden;
}

button::before {
  content: "";
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: rgba(255, 255, 255, 0.1);
  transition: left 0.3s ease;
  pointer-events: none;
}

button:hover::before {
  left: 100%;
}
</style>
