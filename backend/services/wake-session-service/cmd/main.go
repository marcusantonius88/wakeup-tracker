package main

import (
	"log"

	"wakeup-tracker/backend/services/wake-session-service/internal/infrastructure"
)

func main() {
	server := infrastructure.NewServer()
	log.Printf("wake-session-service listening on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
