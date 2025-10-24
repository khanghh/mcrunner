<script setup>
// No need to import icons - they are globally registered
</script>

<template>
  <header class="bg-dark px-6 py-4 flex justify-between items-center shadow-custom z-10 relative">
    <div class="logo flex items-center gap-3 font-bold text-xl text-text">
      <BoxIcon class="text-primary text-2xl" />
      <span>Minecraft Server Editor</span>
    </div>
    <div class="server-status flex items-center gap-4 bg-light px-4 py-2 rounded-lg text-sm">
      <div class="status-indicator w-2.5 h-2.5 rounded-full bg-success shadow-[0_0_8px_#2ecc71] animate-pulse-custom"></div>
      <span>Server Status: Online</span>
      <div class="player-count flex items-center gap-1.5">
        <UsersIcon class="w-4 h-4" />
        <span>3/20</span>
      </div>
    </div>
  </header>
</template>

<style scoped></style>