<script setup>
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { Terminal } from '@xterm/xterm';
import { FitAddon } from '@xterm/addon-fit';
import '@xterm/xterm/css/xterm.css';

const terminalContainer = ref(null);
let term = null;
let fitAddon = null;
let resizeObserver = null;

onMounted(() => {
  term = new Terminal();
  fitAddon = new FitAddon();
  resizeObserver = new ResizeObserver(() => {
    fitAddon.fit();
  });
  resizeObserver.observe(terminalContainer.value);
  term.loadAddon(fitAddon);
  term.open(terminalContainer.value);
  term.write('Welcome to the Xterm.js terminal in Vue 3!\r\n');

  term.onData(e => {
    term.write(e); 
  });
});

onBeforeUnmount(() => {
  if (resizeObserver) {
    resizeObserver.disconnect();
  }
  if (term) {
    term.dispose();
  }
});

</script>

<template>
  <div class="flex flex-col h-full w-full">
    <div class="terminal-header bg-dark px-5 py-3.5 border-b border-border border-gray-500 flex justify-between items-center">
      <div class="terminal-title flex items-center gap-2.5 font-medium">
        <TerminalIcon class="w-4 h-4" />
        <span>Server Console</span>
      </div>
      <div class="terminal-controls flex gap-2.5">
        <button class="btn btn-clear px-4 py-2 border-none rounded-lg cursor-pointer flex items-center gap-1.5 transition-all font-medium text-sm bg-secondary text-white hover:bg-[#2c4a6b] hover:-translate-y-0.5" @click="clearLogs">
          <Trash2Icon class="w-4 h-4" />
          <span>Clear</span>
        </button>
        <button class="btn btn-stop px-4 py-2 border-none rounded-lg cursor-pointer flex items-center gap-1.5 transition-all font-medium text-sm bg-secondary text-white hover:bg-[#2c4a6b] hover:-translate-y-0.5" @click="stopServer">
          <SquareIcon class="w-4 h-4" />
          <span>Stop</span>
        </button>
        <button class="btn btn-restart px-4 py-2 border-none rounded-lg cursor-pointer flex items-center gap-1.5 transition-all font-medium text-sm bg-secondary text-white hover:bg-[#2c4a6b] hover:-translate-y-0.5" @click="restartServer">
          <RotateCcwIcon class="w-4 h-4" />
          <span>Restart</span>
        </button>
      </div>
    </div>
    <div ref="terminalContainer" class="terminal-container"></div>
  </div>
</template>

<style scoped>
.terminal-container {
  /* allow this container to take remaining space in the column flex layout */
  flex: 1 1 0%;
  min-height: 0; /* allow children to shrink properly in flexbox */
  position: relative;
}

/* Ensure xterm elements fill the container when possible */
.terminal-container .xterm,
.terminal-container .xterm-viewport,
.terminal-container .xterm-screen,
.terminal-container .xterm-rows {
  height: 100% !important;
}

/* make the underlying canvas/textarea fill the area too */
.terminal-container .xterm-canvas,
.terminal-container .xterm-text-layer,
.terminal-container .xterm-viewport {
  height: 100% !important;
}
</style>