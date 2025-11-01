<script setup>
import { ref } from 'vue'
import Header from './Header.vue'
import Tabs from './Tabs.vue'

const activeTab = ref('file-explorer')
</script>

<template>
  <div class="app-container flex flex-col h-screen">
    <Header />
    <Tabs v-model:activeTab="activeTab" />
    <div class="content flex-1 overflow-hidden relative bg-darker">
      <slot :activeTab="activeTab" />
    </div>
  </div>
</template>

<style scoped></style>