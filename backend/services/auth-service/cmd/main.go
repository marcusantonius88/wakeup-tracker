package main

import (
	"log"

	"wakeup-tracker/backend/services/auth-service/internal/infrastructure"
)

func main() {
	server := infrastructure.NewServer()
	log.Printf("auth-service listening on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
