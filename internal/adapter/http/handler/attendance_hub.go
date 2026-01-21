package handler

import (
	"sync"

	"github.com/gorilla/websocket"
)

type AttendanceHub struct {
	mu      sync.Mutex
	clients map[*websocket.Conn]struct{}
}

func NewAttendanceHub() *AttendanceHub {
	return &AttendanceHub{
		clients: make(map[*websocket.Conn]struct{}),
	}
}

func (h *AttendanceHub) Add(conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[conn] = struct{}{}
}

func (h *AttendanceHub) Remove(conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.clients, conn)
}

func (h *AttendanceHub) Broadcast(message []byte) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for conn := range h.clients {
		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			_ = conn.Close()
			delete(h.clients, conn)
		}
	}
}
