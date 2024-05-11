package server

import (
	"net/http"

	clienthub "github.com/code-buddy/internal/client_hub"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
	r.Header.Del("Sec-WebSocket-Extensions")
	conn, err := upgrader.Upgrade(w, r, r.Header)
	if err != nil {
		return
	}

	var hub *clienthub.Hub

	h := r.PathValue("hubId")
	if h != "" {
		hub = clienthub.GetHub(h)
	} else {
		hub, err = clienthub.CreateHub()
	}

	if err != nil || hub == nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	u := r.PathValue("userName")
	if u == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	clienthub.InitClient(u, hub, conn)

}
