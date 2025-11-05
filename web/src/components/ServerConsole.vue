<script setup>
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { Terminal } from '@xterm/xterm';
import { FitAddon } from '@xterm/addon-fit';
import '@xterm/xterm/css/xterm.css';

import { Message, MessageType, PtyBuffer, PtyInput, PtyResize } from '@/generated/message.js';

const wsUrl = 'ws://172.17.0.3:3000/ws'; // will be set in connectWebSocket
const terminalContainer = ref(null);
let term = null;
let fitAddon = null;
let resizeObserver = null;
let ws = null;

function connectWebSocket() {
  console.log('Connecting to:', wsUrl);
  ws = new WebSocket(wsUrl);

  ws.onopen = () => {
    console.log('WebSocket connected');
    term.clear();
    sendResize();
  };

  ws.onmessage = async (ev) => {
    const data = await ev.data.arrayBuffer();
    let msg = Message.decode(new Uint8Array(data));
    if (msg.type === MessageType.PTY_BUFFER) {
      const ptyBuffer = msg.ptyBuffer;
      if (ptyBuffer) {
        term.write(ptyBuffer.data);
      }
    } else if (msg.type === MessageType.ERROR) {
      alert('Server error:', msg.error);
    }
  };

  ws.onclose = () => {
    console.log('WebSocket disconnected, reconnecting...');
    setTimeout(connectWebSocket, 1000);
  };

  ws.onerror = () => {
    console.log('WebSocket error');
  };
}

function sendResize() {
  if (!ws || ws.readyState !== WebSocket.OPEN || !term || !fitAddon) return;

  try {
    fitAddon.fit();
  } catch (e) {
    console.error('Fit error:', e);
  }

  const cols = term.cols || 80;
  const rows = term.rows || 24;
  const msg = Message.create({
    type: MessageType.PTY_RESIZE,
    ptyResize: PtyResize.create({ cols, rows })
  });
  ws.send(Message.encode(msg).finish());
  console.log(`Terminal resized to ${cols}x${rows}`);
}

function sendInput(data) {
  if (!ws || ws.readyState !== WebSocket.OPEN) return;
  const msg = Message.create({
    type: MessageType.PTY_INPUT,
    ptyInput: PtyInput.create({ data })
  });
  ws.send(Message.encode(msg).finish());
}

const clearLogs = () => {
  if (term) {
    term.clear();
    fitAddon.fit();
  }
};

const restartServer = () => {
  // TODO: Implement server restart logic
  console.log('Restart server');
};

onMounted(() => {
  term = new Terminal({
    convertEol: true,
    cursorBlink: true,
    fontFamily: 'ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", monospace',
    fontSize: 13,
    theme: { background: '#111111' }
  });

  fitAddon = new FitAddon();
  term.loadAddon(fitAddon);
  term.open(terminalContainer.value);

  // Fit terminal and setup resize observer
  resizeObserver = new ResizeObserver(() => {
    sendResize();
  });
  resizeObserver.observe(terminalContainer.value);
  fitAddon.fit();

  // Connect terminal input to WebSocket
  term.onData((data) => {
    sendInput(data);
  });

  // Connect WebSocket
  connectWebSocket();
});

onBeforeUnmount(() => {
  if (resizeObserver) {
    resizeObserver.disconnect();
  }
  if (ws) {
    ws.close();
  }
  if (term) {
    term.dispose();
  }
});

</script>

<template>
  <div class="flex flex-col h-full w-full">
    <div
      class="terminal-header bg-dark px-5 py-3.5 border-b border-border border-gray-500 flex justify-between items-center">
      <div class="terminal-title flex items-center gap-2.5 font-medium">
        <TerminalIcon class="w-4 h-4" />
        <span>Server Console</span>
      </div>
      <div class="terminal-controls flex gap-2.5">
        <button
          class="btn btn-clear px-4 py-2 border-none rounded-lg cursor-pointer flex items-center gap-1.5 transition-all font-medium text-sm bg-secondary text-white hover:bg-[#2c4a6b] hover:-translate-y-0.5"
          @click="clearLogs">
          <Trash2Icon class="w-4 h-4" />
          <span>Clear</span>
        </button>
        <button
          class="btn btn-stop px-4 py-2 border-none rounded-lg cursor-pointer flex items-center gap-1.5 transition-all font-medium text-sm bg-secondary text-white hover:bg-[#2c4a6b] hover:-translate-y-0.5"
          @click="stopServer">
          <SquareIcon class="w-4 h-4" />
          <span>Stop</span>
        </button>
        <button
          class="btn btn-restart px-4 py-2 border-none rounded-lg cursor-pointer flex items-center gap-1.5 transition-all font-medium text-sm bg-secondary text-white hover:bg-[#2c4a6b] hover:-translate-y-0.5"
          @click="restartServer">
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
  min-height: 0;
  /* allow children to shrink properly in flexbox */
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
