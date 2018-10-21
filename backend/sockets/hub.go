package sockets

import (
	"encoding/json"
	"log"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clientsByRoom map[string]map[*Client]string

	// Inbound messages from the clients.
	broadcast chan Message

	// Register requests from the clients.
	register chan Subscription

	// Unregister requests from clients.
	unregister chan Subscription
}

type Subscription struct {
	room   string
	user   string
	client *Client
}

type Message struct {
	room string
	user string
	data []byte
}

type RoomStatus struct {
	Room   string   `json:"room"`
	Action string   `json:"action"`
	Users  []string `json:"users"`
}

func NewHub() *Hub {
	return &Hub{
		broadcast:     make(chan Message),
		register:      make(chan Subscription),
		unregister:    make(chan Subscription),
		clientsByRoom: make(map[string]map[*Client]string),
	}
}

func (hub *Hub) Run() {
	for {
		select {
		case subscription := <-hub.register:
			log.Printf("Hub subscription register: %v %v", subscription.room, subscription.user)
			clients := hub.clientsByRoom[subscription.room]
			if clients == nil {
				clients = make(map[*Client]string)
				hub.clientsByRoom[subscription.room] = clients
			}
			hub.clientsByRoom[subscription.room][subscription.client] = subscription.user
			hub.broadcastRoomStatus(subscription.room, "Register: " + subscription.user)

		case subscription := <-hub.unregister:
			log.Printf("Hub subscription unregister: %v %v", subscription.room, subscription.user)
			clients := hub.clientsByRoom[subscription.room]
			if clients != nil {
				if _, ok := clients[subscription.client]; ok {
					hub.removeClient(subscription.client, subscription.room)
				}
			}
			hub.broadcastRoomStatus(subscription.room, "Unregister: " + subscription.user)

		case message := <-hub.broadcast:
			log.Printf("Hub message broadcast: %v %v", message.room, message.user)
			hub.broadcastData(message.data, message.room)
		}

	}
}

func (hub *Hub) getRoomUsers(room string) []string {
	var users []string
	clients := hub.clientsByRoom[room]
	for _, user := range clients {
		users = append(users, user)
	}
	log.Printf("Room users: %v", users)
	return users
}

func (hub *Hub) broadcastRoomStatus(room string, action string) {
	log.Printf("Broadcast room status: %v", room)
	roomStatus := &RoomStatus{
		Room:   room,
		Action: action,
		Users:  hub.getRoomUsers(room)}
	roomStatusJson, _ := json.Marshal(roomStatus)
	roomStatusData := []byte(roomStatusJson)
	hub.broadcastData(roomStatusData, room)
}

func (hub *Hub) broadcastData(data []byte, room string) {
	clients := hub.clientsByRoom[room]
	for client := range clients {
		select {
		case client.send <- data:
		default:
			hub.removeClient(client, room)
		}
	}
}

func (hub *Hub) removeClient(client *Client, room string) {
	close(client.send)
	clients := hub.clientsByRoom[room]
	delete(clients, client)
	if len(clients) == 0 {
		delete(hub.clientsByRoom, room)
	}
}
