<template>
  <div class="compute-layout">
    <router-view :key="$route.fullPath" />
  </div>
</template>

<script setup lang="ts">
import { useRoute } from 'vue-router'
const route = useRoute()
// This is a layout component that renders child routes
// The :key ensures re-render when child route changes
</script>

<style scoped>
.compute-layout {
  padding: 0;
  margin: 0;
}
</style>