<script setup>
import { ref, onMounted, nextTick } from 'vue'

const commandInput = ref('')
const logs = ref([
  { time: '12:34:56', type: 'info', message: 'Starting Minecraft server on *:25565' },
  { time: '12:34:56', type: 'info', message: 'Preparing level "world"' },
  { time: '12:34:57', type: 'info', message: 'Loaded 7 recipes' },
  { time: '12:34:58', type: 'info', message: 'Loaded 1184 advancements' },
  { time: '12:35:01', type: 'success', message: 'Done (4.123s)! For help, type "help"' },
  { time: '12:35:05', type: 'info', message: 'Player123 joined the game' },
  { time: '12:35:10', type: 'info', message: 'CreeperHunter joined the game' },
  { time: '12:35:15', type: 'info', message: 'MineCrafter99 joined the game' },
  { time: '12:36:22', type: 'warning', message: 'Can\'t keep up! Is the server overloaded?' },
  { time: '12:37:45', type: 'info', message: 'Player123 issued server command: /time set day' },
  { time: '12:38:10', type: 'info', message: 'CreeperHunter was slain by Zombie' },
  { time: '12:40:33', type: 'info', message: 'MineCrafter99 issued server command: /give Player123 minecraft:diamond 5' }
])

const logsContainer = ref(null)
const stopped = ref(false)

const getCurrentTime = () => {
  const now = new Date()
  return `${now.getHours().toString().padStart(2, '0')}:${now.getMinutes().toString().padStart(2, '0')}:${now.getSeconds().toString().padStart(2, '0')}`
}

const appendLog = (message, type = 'info') => {
  logs.value.push({
    time: getCurrentTime(),
    type,
    message
  })
  nextTick(() => {
    if (logsContainer.value) {
      logsContainer.value.scrollTop = logsContainer.value.scrollHeight
    }
  })
}

const sendCommand = async () => {
  const command = commandInput.value.trim()
  if (!command) return

  // Add command to logs
  appendLog(command, 'command')

  // Clear input
  commandInput.value = ''

  // Send command to server
  try {
    const response = await fetch('/api/v1/server/command', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ command }),
    })
    if (!response.ok) {
      throw new Error('Failed to send command')
    }
    const data = await response.json()
    console.log('Command sent:', data)
  } catch (error) {
    console.error('Error sending command:', error)
    appendLog('Failed to send command', 'error')
  }
}

const clearLogs = () => {
  logs.value = []
}

const stopServer = () => {
  stopped.value = true
  commandInput.value.disabled = true
  appendLog('Server stopped', 'error')
}

const restartServer = () => {
  if (!stopped.value) {
    appendLog('Server is restarting...', 'info')
  }
  stopped.value = false
  commandInput.value.disabled = false
  // Simulate server restart sequence
  appendLog('Stopping server...', 'info')
  setTimeout(() => appendLog('Starting server...', 'info'), 800)
  setTimeout(() => appendLog('Done (0.512s)! For help, type "help"', 'success'), 1600)
}

const handleKeyPress = (event) => {
  if (event.key === 'Enter') {
    sendCommand()
  }
}

// Simulate server logs
onMounted(() => {
  const interval = setInterval(() => {
    const logTypes = ['info', 'info', 'info', 'info', 'warning']
    const logMessages = [
      'A zombie was slain by Player123',
      'CreeperHunter mined diamond ore',
      'MineCrafter99 placed a block',
      'Server performance is normal',
      'Can\'t keep up! Is the server overloaded?'
    ]

    const randomType = logTypes[Math.floor(Math.random() * logTypes.length)]
    const randomMessage = logMessages[Math.floor(Math.random() * logMessages.length)]

    appendLog(randomMessage, randomType)
  }, 10000) // Add a log every 10 seconds

  return () => clearInterval(interval)
})
</script>

<template>
  <div class="terminal-container flex flex-col h-full w-full">
    <div class="terminal-header bg-dark px-5 py-3.5 border-b border-border flex justify-between items-center">
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
    <div ref="logsContainer" class="terminal-logs flex-1 p-5 overflow-y-auto font-code text-sm bg-[#0d1117] border-b border-border custom-scrollbar">
      <div
        v-for="log in logs"
        :key="`${log.time}-${log.message}`"
        class="log-entry mb-2"
        :class="{
          'text-info': log.type === 'info',
          'text-accent': log.type === 'error',
          'text-success': log.type === 'success',
          'text-warning': log.type === 'warning',
          'text-text': log.type === 'command'
        }"
      >
        [{{ log.time }} {{ log.type.toUpperCase() }}]: {{ log.message }}
      </div>
    </div>
    <div class="terminal-input-container p-5 flex gap-3">
      <input
        v-model="commandInput"
        type="text"
        class="terminal-input flex-1 bg-dark text-text border border-border rounded-lg px-4 py-3 font-code outline-none transition-all focus:border-primary focus:shadow-[0_0_0_2px_rgba(0,188,140,0.2)]"
        placeholder="Enter server command..."
        @keypress="handleKeyPress"
      >
      <button class="btn-terminal bg-primary text-white border-none rounded-lg px-5 py-3 cursor-pointer transition-all font-medium flex items-center gap-1.5 hover:bg-primary-dark hover:-translate-y-0.5" @click="sendCommand">
        <SendIcon class="w-4 h-4" />
        <span>Send</span>
      </button>
    </div>
  </div>
</template>

<style scoped></style>