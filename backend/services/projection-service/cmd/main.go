package main

import (
	"log"

	"wakeup-tracker/backend/services/projection-service/internal/infrastructure"
)

func main() {
	server := infrastructure.NewServer()
	log.Printf("projection-service listening on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
