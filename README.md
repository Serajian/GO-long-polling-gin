# Long-Polling with Gin (Go)

A lightweight and clean implementation of **Long-Polling** using the **Gin** framework in Go.  
This project demonstrates how to build a near real-time communication pattern between clients and the server.

<b>This is an educational project for learning long-polling, which can help understand socket-based communication.</b>

---

## 🧠 How It Works (Summary)

1) ##### A client calls /poll/:id and waits for a message.

2) ##### The server creates a channel (buffer 1) for that id.

3) ##### When /send/:id is called, the message is delivered via the channel and the client receives it. The client entry is then removed.

4) ##### If no message arrives within 30 seconds, or if the client disconnects, the channel is closed and the client is removed.

---

## 🚀 How to use

first clone:
```bash
git clone  https://github.com/Serajian/GO-long-polling-gin.git
cd GO-long-polling-gin
```
get requires:
```bash
go mod tidy
```

run app:
```bash
go run main.go
```

### The server will start on port 8090

### 🧪 Quick Test

Terminal 1:
```bash
curl -N http://localhost:8090/poll/u1
```

Terminal 2 (within 30s):
```bash
curl -X POST http://localhost:8090/send/u1 \
  -H "Content-Type: application/json" \
  -d '{"message":"ping!"}'

```
### You should see {"message":"ping!"} in Terminal 1.

---

## 🔌 API Endpoints
## 1) Receive messages (Long-Poll)
### GET /poll/:id

:id is the unique client identifier (e.g., user-123).

The server waits up to 30 seconds. If no message is sent, a 504 timeout is returned.

Successful response:
```json
HTTP/1.1 200 OK
{
  "message": "hello from server"
}
```

Timeout:
```json
HTTP/1.1 504 Gateway Timeout
{
  "error": "timeout"
}

```

## 2) Send a message to a client
### POST /send/:id
Request body (JSON):
```json
{
  "message": "hello from server"
}

```
Response:
```json
HTTP/1.1 200 OK
{
  "status": "message sent"
}

```

Note: If no client is currently waiting with that id, the message will be dropped (by design of long-polling).

---

🗺️ Comparison with Other Real-Time Patterns

<ul>Simple Polling: Client sends requests repeatedly (high overhead).</ul>

<ul>Long-Polling (this project): One request stays open until message or timeout.</ul>

<ul>WebSocket: Persistent two-way connection; best for high-frequency low-latency messaging.</ul>

<ul>SSE (Server-Sent Events): One-way streaming from server to client.</ul>