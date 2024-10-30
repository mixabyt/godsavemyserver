package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	file, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalln(err)
	}
	log.SetOutput(file)

	http.HandleFunc("/ws", HandleWebSocket)

	log.Println("WebSocket server starting on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
