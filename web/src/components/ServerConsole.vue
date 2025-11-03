<script setup>
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { Terminal } from '@xterm/xterm';
import { FitAddon } from '@xterm/addon-fit';
import '@xterm/xterm/css/xterm.css';

const terminalContainer = ref(null);
let term = null;
let fitAddon = null;
let resizeObserver = null;
let ws = null;

// Base64 helpers
const enc = new TextEncoder();
const dec = new TextDecoder();
function b64FromBytes(u8) { 
  let bin = ''; 
  for (let i = 0; i < u8.length; i++) bin += String.fromCharCode(u8[i]); 
  return btoa(bin); 
}
function bytesFromB64(b64) { 
  const bin = atob(b64); 
  const u8 = new Uint8Array(bin.length); 
  for (let i = 0; i < bin.length; i++) u8[i] = bin.charCodeAt(i); 
  return u8; 
}

function connectWebSocket() {
  const wsProtocol = location.protocol === 'https:' ? 'wss' : 'ws';
  const wsPath = '/ws';
  const wsUrl = `${wsProtocol}://${location.host}${wsPath}`;
  
  console.log('Connecting to:', wsUrl);
  ws = new WebSocket(wsUrl);
  
  ws.onopen = () => {
    console.log('WebSocket connected');
    // Subscribe to minecraft server session with replay
    ws.send(JSON.stringify({ type: 'subscribe', sessions: ['minecraft'], replay: true }));
    // Send initial resize
    sendResize();
  };
  
  ws.onmessage = (ev) => {
    let msg;
    try {
      msg = JSON.parse(typeof ev.data === 'string' ? ev.data : dec.decode(new Uint8Array(ev.data)));
    } catch {
      return;
    }
    
    if (msg.type === 'output' || msg.type === 'buffer') {
      const u8 = bytesFromB64(msg.data || '');
      const text = dec.decode(u8);
      if (term) term.write(text);
    } else if (msg.type === 'error') {
      console.warn('Server error:', msg);
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
  ws.send(JSON.stringify({ type: 'resize', sessionId: 'minecraft', cols, rows }));
  console.log(`Terminal resized to ${cols}x${rows}`);
}

function sendInput(data) {
  if (!ws || ws.readyState !== WebSocket.OPEN) return;
  const u8 = enc.encode(data);
  ws.send(JSON.stringify({ type: 'input', sessionId: 'minecraft', data: b64FromBytes(u8) }));
}

const clearLogs = () => {
  if (term) {
    term.clear();
    fitAddon.fit();
  }
};

const stopServer = () => {
  // TODO: Implement server stop logic
  console.log('Stop server');
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