# WebSocket Streaming Design — stdout/stderr, Commands, and File-Change Events

## Goal
- Stream Minecraft server `stdout` and `stderr` to connected web clients via a single WebSocket endpoint (`/ws`).
- Accept command input from clients and forward to the running Java process (`ServerRunner`).
- Watch predefined server files (e.g., `server.properties` or the server jar) and emit a `"file_changed"` event on modification.
- Support multiple clients, backpressure handling, and optional authentication.

---

## High-Level Architecture

- **WebSocket Endpoint:** `/ws`
- **Server Components:**
  - **Runner (`ServerRunner`)**
    - Channels:
      '''go
      OutputChan() <-chan string // stdout lines
      ErrorChan() <-chan string  // stderr lines
      '''
    - Commands:
      '''go
      SendCommand(cmd string) error
      '''
  - **WebSocketHub**
    - Accepts connections and maintains a list of clients.
    - Broadcasts messages from `OutputChan`/`ErrorChan`.
  - **Command Handler**
    - Receives `"command"` messages and calls `ServerRunner.SendCommand(command)`.
  - **File Watcher**
    - Watches predefined files and emits `"file_changed"` through the hub on changes.
    - Debounces rapid events (e.g., 500ms).

---

## Message Contract (JSON over WebSocket)

### Outgoing (Server → Client)
'''json
{ "type": "stdout", "line": "..." }
{ "type": "stderr", "line": "..." }
{ "type": "event", "name": "file_changed", "path": "/path/to/file", "timestamp": 1690000000 }
{ "type": "info", "message": "Server restarting" }
{ "type": "error", "message": "description" }
'''

### Incoming (Client → Server)
'''json
{ "type": "command", "command": "say hello" }
{ "type": "ping" }
{ "type": "auth", "token": "..." } // optional, for handshake authentication
'''

---

## Protocol Rules

- **Delivery:** One JSON message per WebSocket frame.
- **Backpressure:** Use non-blocking sends with per-client buffer; on overflow, drop oldest messages or close the client (configurable).
- **Ordering:** Preserve ordering per stream (`stdout`/`stderr`), interleave streams as they arrive.
- **Heartbeat:** Support ping/pong or require clients to send periodic pings.

---

## File Watching Behavior

- Use a file watch library (e.g., `fsnotify`) for efficient notifications.
- Emit `"file_changed"` on first detection:
  '''json
  { "type": "event", "name": "file_changed", "path": "/path/to/file", "timestamp": 1690000000 }
  '''
- Coalesce multiple events in a short window (debounce ~500ms).
- Optionally include a snapshot or hash of the file.

---

## Error Handling & Edge Cases

- **Resource limits:** Limit max clients; refuse connections with a clear close reason.
- **Slow clients:** Use bounded per-client buffer. Drop logs on overflow; close clients for control/command overflow.
- **Server restart:** Send `"info"` before shutdown; send resume/synced after restart.
- **Unauthorized clients:** Reject on upgrade if using token/JWT or send `"error"` for commands.

---

## Security

- Use TLS (`wss://`) in production.
- Require authentication for command messages; allow read-only streaming without auth.
- Rate-limit commands per client.

---

## APIs & Helpers

### ServerRunner (Go)
'''go
OutputChan() <-chan string
ErrorChan() <-chan string
SendCommand(cmd string) error
'''

### WebSocketHub (Go)
'''go
Broadcast(msg ServerMessage)
RegisterClient(conn *websocket.Conn, options ClientOptions)
UnregisterClient(conn)
'''

---

## Client-Side Responsibilities (JS)

- Connect to `ws(s)://host/ws`.
- Optionally send auth token.
- Listen for messages, parse JSON, append log lines to a scrolling UI.
- Send commands using JSON:
  '''json
  { "type": "command", "command": "..." }
  '''
- Reconnect with backoff on disconnect.

---

## Suggested Libraries

### Go (Server)
- WebSocket: `github.com/gorilla/websocket` or `github.com/gofiber/websocket/v2`
- File Watch: `github.com/fsnotify/fsnotify`
- HTTP: `net/http + gorilla/mux` or `github.com/gofiber/fiber/v2`
- JSON: `encoding/json`

### JavaScript (Client)
- WebSocket: native API
- Reconnection helper: `reconnecting-websocket` or manual backoff
- UI: any framework (React, Svelte, Vue) or plain HTML with virtual list

---

## Example Handling

- **stdout line:**
  '''js
  hub.Broadcast({ type: 'stdout', line: line })
  '''
- **Client command:**
  '''js
  if (client.authenticated) {
      runner.SendCommand(command)
      client.Write({ type: 'info', message: 'command accepted' })
  } else {
      client.Write({ type: 'error', message: 'unauthorized' })
  }
  '''

---

## Testing

- Unit tests for hub broadcast/backpressure logic.
- Integration tests using headless WebSocket clients.
- File-watch integration tests: touch watched file, assert `"file_changed"` emitted.

---

## Open Decisions

- Authentication method (JWT, session cookie, ephemeral token).
- Max connected clients and per-client buffer sizes.
- Support per-client subscriptions (e.g., only stdout or stderr).

---

> **Note:** This is a design document (no code) and should be implemented under `internal/web` or `internal/handlers`.
