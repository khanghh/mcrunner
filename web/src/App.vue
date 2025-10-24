<script setup>
import { ref } from 'vue'
import Header from './components/Header.vue'
import Tabs from './components/Tabs.vue'
import FileExplorer from './components/FileExplorer.vue'
import Terminal from './components/Terminal.vue'

const activeTab = ref('file-explorer')
</script>

<template>
  <div class="app-container flex flex-col h-screen">
    <Header />
    <Tabs v-model:activeTab="activeTab" />
    <div class="content flex-1 overflow-hidden relative bg-darker">
      <FileExplorer v-show="activeTab === 'file-explorer'" />
      <Terminal v-show="activeTab === 'terminal'" />
    </div>
  </div>
</template>

<style scoped></style>
