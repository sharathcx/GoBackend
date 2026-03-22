package websocket

const ChatPageHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Chat</title>
<style>
  * { margin: 0; padding: 0; box-sizing: border-box; }
  body {
    font-family: 'Inter', -apple-system, sans-serif;
    background: #060e20;
    color: #dee5ff;
    height: 100vh;
    display: flex;
  }

  /* Sidebar */
  .sidebar {
    width: 280px;
    background: #091328;
    border-right: 1px solid rgba(255,255,255,0.05);
    display: flex;
    flex-direction: column;
  }
  .sidebar-header {
    padding: 20px;
    border-bottom: 1px solid rgba(255,255,255,0.05);
  }
  .sidebar-header h2 {
    font-size: 18px;
    color: #a3a6ff;
    margin-bottom: 12px;
  }
  .login-form, .create-room-form {
    display: flex;
    gap: 8px;
  }
  .login-form input, .create-room-form input {
    flex: 1;
    background: #0f1930;
    border: 1px solid rgba(255,255,255,0.1);
    color: #dee5ff;
    padding: 8px 12px;
    border-radius: 6px;
    font-size: 13px;
    outline: none;
  }
  .login-form input:focus, .create-room-form input:focus {
    border-color: #a3a6ff;
  }
  .btn {
    background: #494bd7;
    color: white;
    border: none;
    padding: 8px 14px;
    border-radius: 6px;
    cursor: pointer;
    font-size: 13px;
    font-weight: 600;
    white-space: nowrap;
  }
  .btn:hover { background: #5c5ee0; }
  .btn:disabled { opacity: 0.5; cursor: not-allowed; }
  .btn-small {
    padding: 6px 10px;
    font-size: 12px;
  }
  .btn-danger { background: #a70138; }
  .btn-danger:hover { background: #c0144a; }

  .room-list {
    flex: 1;
    overflow-y: auto;
    padding: 12px;
  }
  .room-list-title {
    font-size: 10px;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 1.5px;
    color: #6d758c;
    padding: 8px;
    margin-top: 8px;
  }
  .room-item {
    padding: 10px 12px;
    border-radius: 6px;
    cursor: pointer;
    font-size: 14px;
    color: #a3aac4;
    display: flex;
    justify-content: space-between;
    align-items: center;
    transition: background 0.15s;
  }
  .room-item:hover { background: rgba(255,255,255,0.05); }
  .room-item.active {
    background: rgba(163,166,255,0.1);
    color: #a3a6ff;
    font-weight: 600;
    border-right: 2px solid #a3a6ff;
  }

  .connection-status {
    padding: 12px 20px;
    border-top: 1px solid rgba(255,255,255,0.05);
    font-size: 12px;
    display: flex;
    align-items: center;
    gap: 8px;
  }
  .status-dot {
    width: 8px; height: 8px;
    border-radius: 50%;
    background: #6d758c;
  }
  .status-dot.connected { background: #40ceed; box-shadow: 0 0 8px rgba(64,206,237,0.5); }
  .status-dot.error { background: #ff6e84; }

  /* Main chat area */
  .main {
    flex: 1;
    display: flex;
    flex-direction: column;
  }
  .chat-header {
    padding: 16px 24px;
    border-bottom: 1px solid rgba(255,255,255,0.05);
    background: rgba(6,14,32,0.8);
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
  .chat-header h3 { font-size: 16px; }
  .online-users {
    font-size: 12px;
    color: #6d758c;
  }

  .messages {
    flex: 1;
    overflow-y: auto;
    padding: 20px 24px;
    display: flex;
    flex-direction: column;
    gap: 4px;
  }
  .message {
    max-width: 70%;
    padding: 10px 14px;
    border-radius: 12px;
    font-size: 14px;
    line-height: 1.4;
    word-wrap: break-word;
  }
  .message.mine {
    align-self: flex-end;
    background: #494bd7;
    color: white;
    border-bottom-right-radius: 4px;
  }
  .message.theirs {
    align-self: flex-start;
    background: #141f38;
    border-bottom-left-radius: 4px;
  }
  .message .sender {
    font-size: 11px;
    font-weight: 700;
    color: #a3a6ff;
    margin-bottom: 2px;
  }
  .message.mine .sender { color: rgba(255,255,255,0.7); }
  .message .time {
    font-size: 10px;
    color: rgba(255,255,255,0.3);
    margin-top: 4px;
    text-align: right;
  }
  .message.system {
    align-self: center;
    background: none;
    color: #6d758c;
    font-size: 12px;
    font-style: italic;
    padding: 4px;
    max-width: 100%;
  }

  .typing-indicator {
    padding: 4px 24px;
    font-size: 12px;
    color: #6d758c;
    font-style: italic;
    min-height: 24px;
  }

  .message-input-area {
    padding: 16px 24px;
    border-top: 1px solid rgba(255,255,255,0.05);
    background: #091328;
  }
  .message-input-form {
    display: flex;
    gap: 8px;
  }
  .message-input-form input {
    flex: 1;
    background: #0f1930;
    border: 1px solid rgba(255,255,255,0.1);
    color: #dee5ff;
    padding: 12px 16px;
    border-radius: 8px;
    font-size: 14px;
    outline: none;
  }
  .message-input-form input:focus { border-color: #a3a6ff; }

  .no-room {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: center;
    color: #40485d;
    font-size: 16px;
  }
</style>
</head>
<body>

<div class="sidebar">
  <div class="sidebar-header">
    <h2>Chat</h2>
    <div id="auth-section">
      <div class="login-form" id="token-form">
        <input type="text" id="token-input" placeholder="Paste JWT token...">
        <button class="btn btn-small" onclick="connect()">Connect</button>
      </div>
    </div>
    <div class="create-room-form" style="margin-top:8px;">
      <input type="text" id="room-name-input" placeholder="New room name...">
      <button class="btn btn-small" onclick="createRoom()">Create</button>
    </div>
    <div class="create-room-form" style="margin-top:8px;">
      <input type="text" id="join-room-input" placeholder="Room ID (ROM_xxxx)">
      <button class="btn btn-small" onclick="joinByID()">Join</button>
    </div>
  </div>

  <div class="room-list" id="room-list">
    <div class="room-list-title">Your Rooms</div>
  </div>

  <div class="connection-status">
    <div class="status-dot" id="status-dot"></div>
    <span id="status-text">Disconnected</span>
  </div>
</div>

<div class="main">
  <div id="no-room-view" class="no-room">Select or create a room to start chatting</div>

  <div id="chat-view" style="display:none; flex:1; flex-direction:column;">
    <div class="chat-header">
      <h3 id="room-name">Room</h3>
      <span class="online-users" id="online-users"></span>
    </div>
    <div class="messages" id="messages"></div>
    <div class="typing-indicator" id="typing-indicator"></div>
    <div class="message-input-area">
      <form class="message-input-form" onsubmit="sendMessage(event)">
        <input type="text" id="message-input" placeholder="Type a message..." autocomplete="off">
        <button class="btn" type="submit">Send</button>
      </form>
    </div>
  </div>
</div>

<script>
let ws = null;
let token = '';
let currentRoomID = null;
let myUserID = '';
let typingTimeout = null;

// Parse JWT to get user info (no verification — just decode payload)
function parseJWT(t) {
  try {
    const payload = JSON.parse(atob(t.split('.')[1]));
    return payload;
  } catch { return null; }
}

function setStatus(state, text) {
  const dot = document.getElementById('status-dot');
  const txt = document.getElementById('status-text');
  dot.className = 'status-dot ' + state;
  txt.textContent = text;
}

function connect() {
  token = document.getElementById('token-input').value.trim();
  if (!token) return;

  const claims = parseJWT(token);
  if (claims) myUserID = claims.user_id;

  const protocol = location.protocol === 'https:' ? 'wss:' : 'ws:';
  ws = new WebSocket(protocol + '//' + location.host + '/ws?token=' + token);

  ws.onopen = () => {
    setStatus('connected', 'Connected as ' + (claims ? claims.first_name : 'user'));
    loadRooms();
  };

  ws.onmessage = (e) => {
    const msg = JSON.parse(e.data);
    handleEvent(msg);
  };

  ws.onclose = () => setStatus('', 'Disconnected');
  ws.onerror = () => setStatus('error', 'Connection error');
}

function handleEvent(msg) {
  switch (msg.event) {
    case 'message':
      if (msg.room_id === currentRoomID) {
        appendMessage(msg.data);
      }
      break;

    case 'user_joined':
      if (msg.room_id === currentRoomID) {
        appendSystem(msg.data.username + ' joined the room');
      }
      break;

    case 'user_left':
      if (msg.room_id === currentRoomID) {
        appendSystem(msg.data.username + ' left the room');
      }
      break;

    case 'online_users':
      if (msg.room_id === currentRoomID) {
        const names = msg.data.map(u => u.username);
        document.getElementById('online-users').textContent = names.length + ' online: ' + names.join(', ');
      }
      break;

    case 'typing':
      if (msg.room_id === currentRoomID && msg.data.user_id !== myUserID) {
        showTyping(msg.data.username);
      }
      break;

    case 'error':
      appendSystem('Error: ' + msg.data.message);
      break;
  }
}

function appendMessage(data) {
  const container = document.getElementById('messages');
  const div = document.createElement('div');
  const isMine = data.sender_id === myUserID;
  div.className = 'message ' + (isMine ? 'mine' : 'theirs');
  const time = new Date(data.created_at).toLocaleTimeString([], {hour: '2-digit', minute: '2-digit'});
  div.innerHTML = (isMine ? '' : '<div class="sender">' + data.sender_id + '</div>') +
    '<div>' + escapeHtml(data.content) + '</div>' +
    '<div class="time">' + time + '</div>';
  container.appendChild(div);
  container.scrollTop = container.scrollHeight;
}

function appendSystem(text) {
  const container = document.getElementById('messages');
  const div = document.createElement('div');
  div.className = 'message system';
  div.textContent = text;
  container.appendChild(div);
  container.scrollTop = container.scrollHeight;
}

function showTyping(username) {
  const el = document.getElementById('typing-indicator');
  el.textContent = username + ' is typing...';
  clearTimeout(typingTimeout);
  typingTimeout = setTimeout(() => { el.textContent = ''; }, 3000);
}

function sendMessage(e) {
  e.preventDefault();
  const input = document.getElementById('message-input');
  const content = input.value.trim();
  if (!content || !ws || !currentRoomID) return;

  ws.send(JSON.stringify({ action: 'send_message', room_id: currentRoomID, content: content }));
  input.value = '';
}

// Send typing indicator on keypress
document.addEventListener('DOMContentLoaded', () => {
  const input = document.getElementById('message-input');
  let lastTypingSent = 0;
  input.addEventListener('input', () => {
    if (!ws || !currentRoomID) return;
    const now = Date.now();
    if (now - lastTypingSent > 2000) {
      ws.send(JSON.stringify({ action: 'typing', room_id: currentRoomID }));
      lastTypingSent = now;
    }
  });
});

// REST API calls for room management
async function apiCall(method, path, body) {
  const opts = {
    method: method,
    headers: { 'Authorization': 'Bearer ' + token, 'Content-Type': 'application/json' },
  };
  if (body) opts.body = JSON.stringify(body);
  const res = await fetch(path, opts);
  return res.json();
}

async function loadRooms() {
  const res = await apiCall('GET', '/chat/rooms');
  const list = document.getElementById('room-list');
  list.innerHTML = '<div class="room-list-title">Your Rooms</div>';
  if (res.data) {
    res.data.forEach(room => {
      const div = document.createElement('div');
      div.className = 'room-item' + (room.room_id === currentRoomID ? ' active' : '');
      div.innerHTML = '<span>' + escapeHtml(room.name) + '</span>';
      div.onclick = () => joinRoom(room.room_id, room.name);
      list.appendChild(div);
    });
  }
}

async function createRoom() {
  const input = document.getElementById('room-name-input');
  const name = input.value.trim();
  if (!name || !token) return;

  const res = await apiCall('POST', '/chat/rooms', { name: name });
  if (res.data) {
    input.value = '';
    loadRooms();
    joinRoom(res.data.room_id, res.data.name);
  }
}

function joinByID() {
  const input = document.getElementById('join-room-input');
  const roomID = input.value.trim();
  if (!roomID || !ws) return;
  input.value = '';
  joinRoom(roomID, roomID);
}

function joinRoom(roomID, roomName) {
  // Leave current room first
  if (currentRoomID && ws) {
    ws.send(JSON.stringify({ action: 'leave_room', room_id: currentRoomID }));
  }

  currentRoomID = roomID;
  document.getElementById('room-name').textContent = roomName;
  document.getElementById('messages').innerHTML = '';
  document.getElementById('online-users').textContent = '';
  document.getElementById('no-room-view').style.display = 'none';
  document.getElementById('chat-view').style.display = 'flex';

  // Load message history via REST
  loadMessages(roomID);

  // Join via WebSocket
  if (ws) {
    ws.send(JSON.stringify({ action: 'join_room', room_id: roomID }));
  }

  // Update sidebar active state
  document.querySelectorAll('.room-item').forEach(el => el.classList.remove('active'));
  event.currentTarget.classList.add('active');
}

async function loadMessages(roomID) {
  const res = await apiCall('GET', '/chat/rooms/' + roomID + '/messages');
  if (res.data && res.data.length > 0) {
    // Messages come newest-first from the API, reverse for display
    res.data.reverse().forEach(msg => appendMessage(msg));
  }
}

function escapeHtml(text) {
  const div = document.createElement('div');
  div.textContent = text;
  return div.innerHTML;
}
</script>
</body>
</html>`
