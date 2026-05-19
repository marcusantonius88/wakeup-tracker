package main

import (
	"log"

	"wakeup-tracker/backend/services/analytics-service/internal/infrastructure"
)

func main() {
	server := infrastructure.NewServer()
	log.Printf("analytics-service listening on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
