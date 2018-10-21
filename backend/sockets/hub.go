package sockets

import (
	"encoding/json"
	"fmt"
	"log"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clientsByRoom map[string]map[*Client]User

	// Inbound messages from the clients.
	broadcast chan Message

	// Register requests from the clients.
	register chan Subscription

	// Unregister requests from clients.
	unregister chan Subscription
}

type User struct {
	name    string `json:"name,omitempty"`
	isAdmin bool   `json:"isAdmin,omitempty"`
}

func NewUser(name string) User {
	return User{
		name:    name,
		isAdmin: false,
	}
}

type Subscription struct {
	room   string
	user   User
	client *Client
}

type Message struct {
	room string
	user string
	data []byte
}

type RoomStatus struct {
	Room   string `json:"room"`
	Action string `json:"action"`
	Users  []map[string]string `json:"users"`
}

func NewHub() *Hub {
	return &Hub{
		broadcast:     make(chan Message),
		register:      make(chan Subscription),
		unregister:    make(chan Subscription),
		clientsByRoom: make(map[string]map[*Client]User),
	}
}

func (hub *Hub) Run() {
	for {
		select {
		case subscription := <-hub.register:

			var usr = NewUser(subscription.user.name)
			log.Printf("Hub subscription register: %v %v", subscription.room, usr.name)
			clients := hub.clientsByRoom[subscription.room]
			if clients == nil {
				usr.isAdmin = true /*first user is admin */
				fmt.Print(subscription)
				clients = make(map[*Client]User)
				fmt.Printf("clients ######### \n")
				fmt.Print(clients)
				hub.clientsByRoom[subscription.room] = clients
			}
			hub.clientsByRoom[subscription.room][subscription.client] = usr
			hub.broadcastRoomStatus(subscription.room, "Register: "+usr.name)

		case subscription := <-hub.unregister:
			var usr = NewUser(subscription.user.name)
			log.Printf("Hub subscription unregister: %v %v", subscription.room, usr.name)
			clients := hub.clientsByRoom[subscription.room]
			if clients != nil {
				if _, ok := clients[subscription.client]; ok {
					hub.removeClient(subscription.client, subscription.room)
				}
			}
			hub.broadcastRoomStatus(subscription.room, "Unregister: "+usr.name)

		case message := <-hub.broadcast:
			log.Printf("Hub message broadcast: %v %v", message.room, message.user)
			hub.broadcastData(message.data, message.room)
		}

	}
}

func (hub *Hub) getRoomUsers(room string) []map[string]string {
	var users []map[string]string
	clients := hub.clientsByRoom[room]
	for _, user := range clients {
		x := map[string]string {"name":user.name, "isAdmin":"false"}
		users = append(users, x)

		log.Printf("Room users: %v", user.name)

		//users = append(users, `{name:%d, isAdmin:user.isAdmin}`,user.name )
	}
	log.Printf("Room users: %v", users)


	return users
}

func (hub *Hub) broadcastRoomStatus(room string, action string) {
	fmt.Printf("############ hub.getRoomUsers(room) \n")
	fmt.Println(hub.getRoomUsers(room))

	log.Printf("Broadcast room status: %v", room)
	roomStatus := &RoomStatus{
		Room:   room,
		Action: action,
		Users:  hub.getRoomUsers(room),
	}
	roomStatusJson, _ := json.Marshal(roomStatus)
	fmt.Println(roomStatusJson)
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
