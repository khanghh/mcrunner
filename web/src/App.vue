<script setup>
import { ref } from 'vue'
import Header from './components/Header.vue'
import Tabs from './components/Tabs.vue'
import CodeEditor from './components/CodeEditor.vue'
import Terminal from './components/Terminal.vue'

const activeTab = ref('code-editor')
</script>

<template>
  <div class="app-container flex flex-col h-screen">
    <Header />
    <Tabs v-model:activeTab="activeTab" />
    <div class="content flex-1 overflow-hidden relative bg-darker">
      <iframe class="w-full h-full border-0" src="/vscode" v-show="activeTab === 'code-editor'" ></iframe>
      <Terminal v-show="activeTab === 'terminal'" />
    </div>
  </div>
</template>

<style scoped></style>
