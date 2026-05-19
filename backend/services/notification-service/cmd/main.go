package main

import (
	"log"

	"wakeup-tracker/backend/services/notification-service/internal/infrastructure"
)

func main() {
	server := infrastructure.NewServer()
	log.Printf("notification-service listening on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
