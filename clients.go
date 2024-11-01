package main

import (
	"time"

	"github.com/gorilla/websocket"
)

type Clients struct {
	clientsmap map[int]*Client
}

type Client struct {
	ID           int
	RoomID       int
	Conn         *websocket.Conn
	LastActivity time.Time
}

func (c *Clients) AddNewUser(clientID int, conn *websocket.Conn) {
	mu.Lock()
	user := &Client{ID: clientID, Conn: conn, LastActivity: time.Now().Add(10 * time.Second), RoomID: -1}
	c.clientsmap[clientID] = user
	subMainMenu[clientID] = user
	countuser++
	id++
	mu.Unlock()
}

func (c *Clients) DeleteUser(clientID int) {
	mu.Lock()
	delete(c.clientsmap, clientID)
	countuser--
	mu.Unlock()
}
