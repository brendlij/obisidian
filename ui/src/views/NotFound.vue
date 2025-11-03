<template>
  <div class="not-found">
    <div class="content">
      <h1>404</h1>
      <h2>Page Not Found</h2>
      <p>The page you're looking for doesn't exist.</p>

      <router-link to="/">
        <MyButton
          size="lg"
          icon="mdi:home"
          label="Back to Home"
          variant="primary"
        />
      </router-link>
    </div>
  </div>
</template>

<script setup lang="ts">
import MyButton from "../components/MyButton.vue";
</script>

<style scoped>
.not-found {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  padding: 2rem;
}

.content {
  text-align: center;
}

h1 {
  font-size: 6rem;
  margin: 0;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

h2 {
  font-size: 2rem;
  margin: 0.5rem 0;
}

p {
  color: rgba(255, 255, 255, 0.7);
  margin-bottom: 2rem;
}
</style>
