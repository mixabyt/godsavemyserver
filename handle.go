package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var (
	countuser   = 0
	id          = 0
	subMainMenu = make(map[int]*Client) // підписка на лічильник користувачів
	clients     = &Clients{clientsmap: make(map[int]*Client)}
	queueUsers  = &QueueUsers{Queue: make([]*Client, 0, 1)}
	rooms       = &Rooms{Rooms: make(map[int]*Room)}
	mu          sync.Mutex
)

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("fail upgrade to web socket: %s", err)
	}
	defer conn.Close()
	onConnect(conn)
	handlemessage(conn)
}

func handlemessage(conn *websocket.Conn) {
	mu.Lock()
	clientID := id
	mu.Unlock()
	fmt.Printf("айді юзера %d\n", clientID)

	// Додаємо клієнта до мапи
	clients.AddNewUser(clientID, conn)

	OncounterNotify(countuser) // оновлюй лічильник юзерів

	done := make(chan bool) // змінна для контролю з'єднання (завершає горутину)
	var wg sync.WaitGroup

	//функція для надсилання hearbeat кожні 10 сек
	wg.Add(1)
	go OnHeartbeat(clientID, done, &wg)

	// обробка всіх повідомлень
	wg.Add(1)
	go ListenClient(clientID, done, &wg)
	wg.Wait()
	fmt.Println("end for some user")

}
