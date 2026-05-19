package main

import (
	"log"

	"wakeup-tracker/backend/services/device-validation-service/internal/infrastructure"
)

func main() {
	server := infrastructure.NewServer()
	log.Printf("device-validation-service listening on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
