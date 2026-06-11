package main

import (
	"log"
	"net/http"
)

func main() {
	InitRedis()

	http.HandleFunc("/ws", HandleWebSocket)

	log.Println("🚀 Conferencing MVP Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
