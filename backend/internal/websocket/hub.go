package websocket

import (
	"encoding/json"
	"log"
	"sync"
	"time"
)

type Hub struct {
	// Registered clients per room
	Rooms map[uint]map[*Client]bool

	// Register requests from clients
	Register chan *Client

	// Unregister requests from clients
	Unregister chan *Client

	// Broadcast messages to clients in a room
	Broadcast chan *Message

	// Mutex for thread-safe operations
	mu sync.RWMutex

	// Online users tracking
	OnlineUsers map[uint]bool
}

func NewHub() *Hub {
	return &Hub{
		Rooms:       make(map[uint]map[*Client]bool),
		Register:    make(chan *Client),
		Unregister:  make(chan *Client),
		Broadcast:   make(chan *Message, 256),
		OnlineUsers: make(map[uint]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			if h.Rooms[client.RoomID] == nil {
				h.Rooms[client.RoomID] = make(map[*Client]bool)
			}
			h.Rooms[client.RoomID][client] = true
			h.OnlineUsers[client.UserID] = true
			h.mu.Unlock()

			// Notify room that user joined
			joinMsg := &Message{
				Type:      "join",
				RoomID:    client.RoomID,
				UserID:    client.UserID,
				Username:  client.Username,
				Timestamp: time.Now(),
			}
			h.broadcastToRoom(client.RoomID, joinMsg)

			log.Printf("Client registered: User %d in Room %d", client.UserID, client.RoomID)

		case client := <-h.Unregister:
			h.mu.Lock()
			if clients, ok := h.Rooms[client.RoomID]; ok {
				if _, exists := clients[client]; exists {
					delete(clients, client)
					close(client.Send)

					// Remove room if empty
					if len(clients) == 0 {
						delete(h.Rooms, client.RoomID)
					}

					// Check if user is still online in other rooms
					stillOnline := false
					for _, roomClients := range h.Rooms {
						for c := range roomClients {
							if c.UserID == client.UserID {
								stillOnline = true
								break
							}
						}
						if stillOnline {
							break
						}
					}
					if !stillOnline {
						delete(h.OnlineUsers, client.UserID)
					}
				}
			}
			h.mu.Unlock()

			// Notify room that user left
			leaveMsg := &Message{
				Type:      "leave",
				RoomID:    client.RoomID,
				UserID:    client.UserID,
				Username:  client.Username,
			}
			h.broadcastToRoom(client.RoomID, leaveMsg)

			log.Printf("Client unregistered: User %d from Room %d", client.UserID, client.RoomID)

		case message := <-h.Broadcast:
			h.broadcastToRoom(message.RoomID, message)
		}
	}
}

func (h *Hub) broadcastToRoom(roomID uint, message *Message) {
	h.mu.RLock()
	clients := h.Rooms[roomID]
	h.mu.RUnlock()

	if clients == nil {
		return
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}

	for client := range clients {
		select {
		case client.Send <- messageBytes:
		default:
			// Client's send channel is full, close and remove
			close(client.Send)
			h.mu.Lock()
			delete(h.Rooms[roomID], client)
			h.mu.Unlock()
		}
	}
}

func (h *Hub) GetRoomClients(roomID uint) []*Client {
	h.mu.RLock()
	defer h.mu.RUnlock()

	clients := h.Rooms[roomID]
	if clients == nil {
		return []*Client{}
	}

	result := make([]*Client, 0, len(clients))
	for client := range clients {
		result = append(result, client)
	}
	return result
}

func (h *Hub) GetOnlineUsers() []uint {
	h.mu.RLock()
	defer h.mu.RUnlock()

	users := make([]uint, 0, len(h.OnlineUsers))
	for userID := range h.OnlineUsers {
		users = append(users, userID)
	}
	return users
}

func (h *Hub) IsUserOnline(userID uint) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.OnlineUsers[userID]
}

func (h *Hub) GetRoomUserCount(roomID uint) int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	clients := h.Rooms[roomID]
	if clients == nil {
		return 0
	}
	return len(clients)
}
