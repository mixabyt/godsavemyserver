package main

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

var (
	RoomsID = 0
)

type Rooms struct {
	Rooms map[int]*Room
}

type Room struct {
	Clients map[int]*Client
}

func (r *Rooms) CreateRoom(clientID int) {
	mu.Lock()
	defer mu.Unlock()

	user := clients.clientsmap[clientID]
	newRoom := &Room{Clients: make(map[int]*Client)}

	newRoom.Clients[clientID] = user
	r.Rooms[RoomsID] = newRoom
	user.RoomID = RoomsID

	log.Printf("id[%d] create and add to room[%d]\n", user.ID, RoomsID)
	RoomsID++
}

func (r *Rooms) AddToRoom(roomID, clientID int) {
	mu.Lock()
	defer mu.Unlock()
	user := clients.clientsmap[clientID]
	user.RoomID = roomID
	room := r.Rooms[roomID]

	room.Clients[user.ID] = user
	log.Printf("id[%d] add to room[%d]\n", clientID, roomID)
}

func (r *Rooms) DeleteRoom(roomID int) {
	delete(r.Rooms, roomID)
}

func (r *Room) SendMessage(user *Client, text []byte) {

	for _, v := range r.Clients {
		log.Println("id:", v.ID)
		if v != user {
			v.Conn.WriteMessage(websocket.TextMessage, text)

		}
	}
}

func (r *Room) RoomDeletionNotice(user *Client) {
	end := &DeleteNotice{Type: "roomDeletionNotice"}
	data, err := json.Marshal(end)
	if err != nil {
		log.Println("fail marshal RoomDeletionNotice()")
	}
	for _, v := range r.Clients {
		if v != user {
			v.RoomID = -1
			v.Conn.WriteMessage(websocket.TextMessage, data)
		}
	}
}
