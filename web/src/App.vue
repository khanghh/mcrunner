<script setup>
import { ref } from 'vue'
import Header from './components/Header.vue'
import Tabs from './components/Tabs.vue'
import ServerConsole from './components/ServerConsole.vue'

const activeTab = ref('code-editor')
</script>

<template>
  <div class="app-container flex flex-col h-screen">
    <Header />
    <Tabs v-model:activeTab="activeTab" />
    <div class="content flex-1 overflow-hidden relative bg-darker">
      <iframe class="w-full h-full border-0" src="/vscode" v-show="activeTab === 'code-editor'"></iframe>
      <ServerConsole v-show="activeTab === 'console'" />
    </div>
  </div>
</template>

<style scoped></style>
