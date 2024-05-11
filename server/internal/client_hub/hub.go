package clienthub

import (
	"sync"

	"github.com/google/uuid"
)

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

// TODO: use some cache like redis to get the hub data
type HubMapper struct {
	hubMap map[string]*Hub
	mu     sync.Mutex
}

func (hm *HubMapper) Set(key uuid.UUID, value *Hub) {
	hm.mu.Lock()
	hm.hubMap[key.String()] = value
	hm.mu.Unlock()
}

func (hm *HubMapper) Get(key string) *Hub {
	hm.mu.Lock()
	defer hm.mu.Unlock()
	return hm.hubMap[key]
}

var HubMap = &HubMapper{hubMap: map[string]*Hub{}, mu: sync.Mutex{}}

func CreateHub() (*Hub, error) {
	hubId, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	hub := &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}

	HubMap.Set(hubId, hub)
	go hub.start()
	return hub, nil
}

func GetHub(hubId string) *Hub {
	return HubMap.Get(hubId)
}

func (h *Hub) start() {
	if h == nil {
		return
	}
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send) //TODO: why this channel is not closeable before deleting entry
			}
		case msg := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- msg:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
