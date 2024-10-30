package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

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
	clients     = make(map[int]*Client)
	queueUsers  = make([]*Client, 0, 2)
	mu          sync.Mutex
)

type Client struct {
	ID           int
	Conn         *websocket.Conn
	LastActivity time.Time
}

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
	clientID := id
	fmt.Printf("айді юзера %d\n", clientID)

	// Додаємо клієнта до мапи
	mu.Lock()
	user := &Client{ID: clientID, Conn: conn, LastActivity: time.Now().Add(10 * time.Second)}
	clients[clientID] = user
	subMainMenu[clientID] = user
	countuser++
	id++
	fmt.Println(clients)
	mu.Unlock()

	// надсилай користувачу його айді
	onRegister(clientID)
	OncounterNotify(countuser)

	done := make(chan bool) // змінна для контролю з'єднання (завершає горутину)
	var wg sync.WaitGroup

	//функція для надсилання hearbeat кожні 10 сек
	wg.Add(1)
	go OnHeartbeat(clientID, done, &wg)

	// обробка всіх повідомлень
	wg.Add(1)
	go ListenClient(clientID, done, &wg)
	wg.Wait()

}
