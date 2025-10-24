import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import './style.css'

// Import and register Lucide icons globally
import {
  Box,
  Users,
  Folder,
  Terminal,
  FileCode,
  File,
  RotateCcw,
  Save,
  Trash2,
  Square,
  Send
} from 'lucide-vue-next'

const app = createApp(App)

// Register icons globally
app.component('BoxIcon', Box)
app.component('UsersIcon', Users)
app.component('FolderIcon', Folder)
app.component('TerminalIcon', Terminal)
app.component('FileCodeIcon', FileCode)
app.component('FileIcon', File)
app.component('RotateCcwIcon', RotateCcw)
app.component('SaveIcon', Save)
app.component('Trash2Icon', Trash2)
app.component('SquareIcon', Square)
app.component('SendIcon', Send)

app.use(router)

app.mount('#app')
