package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

func onConnect(conn *websocket.Conn) {

	log.Printf("user connected: %s", conn.RemoteAddr())
}

func OnHeartbeat(clientID int, done <-chan bool, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-time.After(10 * time.Second):
			// log.Println("send heartbeat to:", clientID)
			mu.Lock()
			client, exists := clients.clientsmap[clientID]
			if !exists {
				mu.Unlock()
				return
			}
			// log.Printf("client[%d] last activity:%d", clientID, time.Now().Second()-client.LastActivity.Second())
			if time.Since(client.LastActivity) > 10*time.Second {

				log.Printf("client[%d] didn't respond to heartbeat, disconnecting", clientID)
				client.Conn.Close()
				mu.Unlock()
				return
			}
			mu.Unlock()

			// Відправляємо heartbeat повідомлення
			heartbeat := &BaseMessage{Type: "heartbeat"}
			data, _ := json.Marshal(heartbeat)
			err := client.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				log.Printf("Failed to send heartbeat to client %d: %s", clientID, err)
				return
			}
		case <-done:
			return

		}
	}
}

func ListenClient(clientID int, done chan<- bool, wg *sync.WaitGroup) {
	defer wg.Done()
	mu.Lock()
	user := clients.clientsmap[clientID]
	mu.Unlock()
	for {
		_, message, err := user.Conn.ReadMessage()
		if err != nil {
			log.Printf("user disconected: %s", user.Conn.RemoteAddr())
			// сповістити користувача по кімнаті що інший покинув
			if user.RoomID != -1 {
				rooms.Rooms[user.RoomID].RoomDeletionNotice(user)
				rooms.DeleteRoom(user.RoomID)
			}
			clients.DeleteUser(clientID)
			OncounterNotify(countuser)
			done <- true
			break
		}

		typemessage := &BaseMessage{}
		err = json.Unmarshal(message, typemessage)
		if err != nil {
			log.Printf("Failed to unmarshal JSON: %s", err)
			continue
		}

		switch typemessage.Type {
		case "heartbeat":
			// log.Printf("heartbeat from: %d\n", user.ID)
			user.LastActivity = time.Now().Add(10 * time.Second)

		case "subMainMenu":
			subuser := &SubMain{}
			json.Unmarshal(message, subuser)
			if subuser.Subscription {
				mu.Lock()
				subMainMenu[clientID] = user
				mu.Unlock()
			} else {
				mu.Lock()
				delete(subMainMenu, clientID)
				mu.Unlock()
			}
		case "findInterlocutor":
			interlocutor, inqueue := queueUsers.AddtoQueue(user)
			if !inqueue {
				rooms.CreateRoom(clientID)
			} else {

				queueUsers.DeleteFromQueue()
				rooms.AddToRoom(interlocutor.RoomID, clientID)
				data, _ := json.Marshal(&FindInterlocutor{Type: "findInterlocutor"})
				for _, conn := range rooms.Rooms[interlocutor.RoomID].Clients {
					conn.Conn.WriteMessage(websocket.TextMessage, data)
				}
			}
		case "stopFindingInterlocutor":
			queueUsers.DeleteFromQueue()
		case "message":
			log.Printf("message from[%d] to room[%d]", user.ID, user.RoomID)

			rooms.Rooms[user.RoomID].SendMessage(user, message)
		default:
			fmt.Printf("невідомий тип повідомлення: %s", typemessage.Type)
		}

	}
}

func OncounterNotify(count int) {
	mu.Lock()
	for _, conn := range subMainMenu {

		sMM := &UpdateCountUser{Type: "subMainMenu", Count: count}
		data, _ := json.Marshal(sMM)
		conn.Conn.WriteMessage(websocket.TextMessage, data)
	}
	mu.Unlock()
}
