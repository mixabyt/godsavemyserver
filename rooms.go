package main

import "log"

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

	RoomsID++
	log.Printf("create and add to room [%s]\n", user.Conn.RemoteAddr())
}

func (r *Rooms) AddToRoom(roomID, clientID int) {
	mu.Lock()
	defer mu.Unlock()
	user := clients.clientsmap[clientID]
	room := r.Rooms[roomID]

	room.Clients[user.ID] = user
	log.Printf("added to room %s", user.Conn.RemoteAddr())
}
