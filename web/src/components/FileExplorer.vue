<script setup>
import { ref, onMounted } from 'vue'

const currentFile = ref('server.properties')
const fileContent = ref(`# Minecraft server properties
# Last modified: 2023-10-15

server-port=25565
max-players=20
level-name=world
online-mode=true
difficulty=normal
enable-command-block=true
motd=A Minecraft Server
pvp=true
allow-flight=false
enforce-whitelist=false
spawn-protection=16
view-distance=10
server-ip=

# Game settings
gamemode=survival
hardcore=false
allow-nether=true
spawn-animals=true
spawn-npcs=true
spawn-monsters=true
generate-structures=true

# Server resource pack
resource-pack=`)

const fileTree = ref([
  {
    name: 'config',
    type: 'folder',
    path: 'config',
    open: false,
    children: [
      { name: 'server.properties', type: 'file', path: 'config/server.properties' },
      { name: 'bukkit.yml', type: 'file', path: 'config/bukkit.yml' },
      { name: 'spigot.yml', type: 'file', path: 'config/spigot.yml' }
    ]
  },
  {
    name: 'plugins',
    type: 'folder',
    path: 'plugins',
    open: false,
    children: [
      {
        name: 'WorldEdit',
        type: 'folder',
        path: 'plugins/WorldEdit',
        open: false,
        children: [
          { name: 'config.yml', type: 'file', path: 'plugins/WorldEdit/config.yml' }
        ]
      },
      {
        name: 'Essentials',
        type: 'folder',
        path: 'plugins/Essentials',
        open: false,
        children: [
          { name: 'config.yml', type: 'file', path: 'plugins/Essentials/config.yml' }
        ]
      }
    ]
  },
  { name: 'whitelist.json', type: 'file', path: 'whitelist.json' },
  { name: 'ops.json', type: 'file', path: 'ops.json' },
  { name: 'server.jar', type: 'file', path: 'server.jar' }
])

const sampleFiles = {
  'config/server.properties': `# Minecraft server properties
max-players=20
motd=My Server`,
  'config/bukkit.yml': '# bukkit config...\n',
  'config/spigot.yml': '# spigot config...\n',
  'plugins/WorldEdit/config.yml': '# worldedit config...\n',
  'plugins/Essentials/config.yml': '# essentials config...\n',
  'whitelist.json': '[\n  "Player123"\n]',
  'ops.json': '[\n  { "uuid": "0000-0000", "name": "Admin" }\n]'
}

const toggleFolder = (folder) => {
  folder.open = !folder.open
}

const openFile = (file) => {
  currentFile.value = file.name
  fileContent.value = sampleFiles[file.path] || `// ${file.name} is empty or binary`
}

const reloadFile = () => {
  // Simulate reload
  console.log('Reloading file')
}

const saveFile = () => {
  // Simulate save
  console.log('Saving file')
}
</script>

<template>
  <div class="file-explorer flex h-full w-full">
    <!-- File Tree -->
    <div class="file-tree w-72 bg-dark border-r border-border overflow-y-auto p-4 custom-scrollbar">
      <ul>
        <li v-for="item in fileTree" :key="item.path" :class="{ folder: item.type === 'folder', file: item.type === 'file' }">
          <div
            v-if="item.type === 'folder'"
            class="flex items-center py-2 px-3 cursor-pointer rounded-lg transition-all hover:bg-light ml-0"
            @click="toggleFolder(item)"
          >
            <FolderIcon class="w-4 inline-block mr-2 text-center text-primary" />
            <span class="name font-semibold text-primary">{{ item.name }}</span>
          </div>
          <div
            v-else
            class="flex items-center py-2 px-3 cursor-pointer rounded-lg transition-all hover:bg-light ml-0"
            @click="openFile(item)"
          >
            <FileCodeIcon v-if="item.name.endsWith('.yml') || item.name.endsWith('.properties') || item.name.endsWith('.jar')" class="w-4 inline-block mr-2 text-center text-text" />
            <FileIcon v-else class="w-4 inline-block mr-2 text-center text-text" />
            <span class="name text-text">{{ item.name }}</span>
          </div>
          <ul v-if="item.type === 'folder' && item.open">
            <li v-for="child in item.children" :key="child.path" :class="{ folder: child.type === 'folder', file: child.type === 'file' }">
              <div
                v-if="child.type === 'folder'"
                class="flex items-center py-2 px-3 cursor-pointer rounded-lg transition-all hover:bg-light ml-4"
                @click="toggleFolder(child)"
              >
                <FolderIcon class="w-4 inline-block mr-2 text-center text-primary" />
                <span class="name font-semibold text-primary">{{ child.name }}</span>
              </div>
              <div
                v-else
                class="flex items-center py-2 px-3 cursor-pointer rounded-lg transition-all hover:bg-light ml-4"
                @click="openFile(child)"
              >
                <FileCodeIcon v-if="child.name.endsWith('.yml') || child.name.endsWith('.properties')" class="w-4 inline-block mr-2 text-center text-text" />
                <FileIcon v-else class="w-4 inline-block mr-2 text-center text-text" />
                <span class="name text-text">{{ child.name }}</span>
              </div>
              <ul v-if="child.type === 'folder' && child.open">
                <li v-for="grandchild in child.children" :key="grandchild.path" class="file">
                  <div
                    class="flex items-center py-2 px-3 cursor-pointer rounded-lg transition-all hover:bg-light ml-8"
                    @click="openFile(grandchild)"
                  >
                    <FileCodeIcon v-if="grandchild.name.endsWith('.yml')" class="w-4 inline-block mr-2 text-center text-text" />
                    <FileIcon v-else class="w-4 inline-block mr-2 text-center text-text" />
                    <span class="name text-text">{{ grandchild.name }}</span>
                  </div>
                </li>
              </ul>
            </li>
          </ul>
        </li>
      </ul>
    </div>

    <!-- Editor -->
    <div class="editor-container flex-1 flex flex-col overflow-hidden">
      <div class="editor-header bg-dark px-5 py-3.5 border-b border-border flex justify-between items-center">
        <div class="file-info flex items-center gap-2.5 font-medium">
          <FileCodeIcon class="w-4 h-4" />
          <span>{{ currentFile }}</span>
        </div>
        <div class="editor-actions flex gap-2.5">
          <button class="btn btn-secondary px-4 py-2 border-none rounded-lg cursor-pointer flex items-center gap-1.5 transition-all font-medium text-sm bg-secondary text-white hover:bg-[#2c4a6b] hover:-translate-y-0.5" @click="reloadFile">
            <RotateCcwIcon class="w-4 h-4" />
            <span>Reload</span>
          </button>
          <button class="btn btn-primary px-4 py-2 border-none rounded-lg cursor-pointer flex items-center gap-1.5 transition-all font-medium text-sm bg-primary text-white hover:bg-primary-dark hover:-translate-y-0.5" @click="saveFile">
            <SaveIcon class="w-4 h-4" />
            <span>Save</span>
          </button>
        </div>
      </div>
      <div class="code-editor flex-1 bg-darker p-5 overflow-auto">
        <textarea
          v-model="fileContent"
          class="w-full h-full bg-dark text-text border border-border rounded-lg p-5 font-code text-sm resize-none outline-none focus:border-primary focus:shadow-[0_0_0_2px_rgba(0,188,140,0.2)] transition-all"
        ></textarea>
      </div>
    </div>
  </div>
</template>

<style scoped></style>