package main

import (
	"log"
)

func main() {
	r := setupRouter()
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start the URLShortener server: %v", err)
	}
}
