package handler

import (
	"net/http"

	"github.com/gorilla/websocket"
)

const attendanceAPIKey = "Dptaspen@25!"

var attendanceUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) AttendanceWebsocket(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("key") != attendanceAPIKey {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := attendanceUpgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	h.hub.Add(conn)
	defer h.hub.Remove(conn)

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}
}
