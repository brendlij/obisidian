<script setup lang="ts">
import * as obisidan from "./helper";

function listServers() {
  obisidan.listServers().then((servers) => {
    console.log("Servers:", servers);
  });
}
</script>

<template>
  <div id="app">
    <h1>Welcome to Obisidan UI</h1>
    <button @click="listServers">List Servers</button>
  </div>
</template>

<style scoped></style>
