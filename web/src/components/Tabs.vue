<script setup>
import { ref } from 'vue'

const props = defineProps({
  activeTab: String
})

const emit = defineEmits(['update:activeTab'])

const setActiveTab = (tab) => {
  emit('update:activeTab', tab)
}
</script>

<template>
  <div class="tabs flex bg-dark border-gray-700 border-b">
    <div
      class="tab px-6 py-3.5 cursor-pointer border-b-3 border-transparent transition-all flex items-center gap-2 font-medium text-text-muted hover:text-text hover:bg-white/5"
      :class="{ 'border-gray-500 text-white bg-gray-800': activeTab === 'code-editor' }"
      @click="setActiveTab('code-editor')"
    >
      <FolderIcon class="w-4 h-4" />
      <span>Code Editor</span>
    </div>
    <div
      class="tab px-6 py-3.5 cursor-pointer border-b-3 border-transparent transition-all flex items-center gap-2 font-medium text-text-muted hover:text-text hover:bg-white/5"
      :class="{ 'border-gray-500 text-white bg-gray-800': activeTab === 'terminal' }"
      @click="setActiveTab('terminal')"
    >
      <TerminalIcon class="w-4 h-4" />
      <span>Server Terminal</span>
    </div>
  </div>
</template>

<style scoped></style>
