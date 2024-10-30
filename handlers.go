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
func onRegister(clientID int) {
	user := clients[clientID]
	message := &Register{Type: "register", UserID: user.ID}
	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("onRegister: fail to Marshal json: %s", err)
	}
	err = user.Conn.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		log.Printf("onRegister fail to WriteMessage: %s", err)
	}
}

func OnHeartbeat(clientID int, done <-chan bool, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-time.After(10 * time.Second):
			mu.Lock()
			client, exists := clients[clientID]
			if !exists {
				mu.Unlock()
				return
			}

			if time.Since(client.LastActivity) > 10*time.Second {
				log.Printf("Client %s didn't respond to heartbeat, disconnecting", clients[clientID].Conn.RemoteAddr())
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
	user := clients[clientID]

	for {
		_, message, err := user.Conn.ReadMessage()
		if err != nil {
			log.Printf("user disconected: %s", user.Conn.RemoteAddr())
			mu.Lock()
			delete(clients, clientID)
			countuser--
			OncounterNotify(countuser)
			mu.Unlock()
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
			fmt.Printf("heartbeat from: %s\n", user.Conn.RemoteAddr())
			mu.Lock()
			user.LastActivity = time.Now()
			mu.Unlock()

		case "subMainMenu":
			subuser := &SubMain{}
			json.Unmarshal(message, subuser)
			if subuser.Subscription {
				subMainMenu[clientID] = user
			} else {
				delete(subMainMenu, clientID)
			}
		case "findInterlocutor":
			continue

		default:
			fmt.Printf("невідомий тип повідомлення: %s", typemessage.Type)
		}

	}
}

func OncounterNotify(count int) {
	for _, conn := range subMainMenu {

		sMM := &UpdateCountUser{Type: "subMainMenu", Count: count}
		data, _ := json.Marshal(sMM)
		conn.Conn.WriteMessage(websocket.TextMessage, data)
	}
}
