package main

import (
	"log"

	//"github.com/joho/godotenv"
)

func main() {
//	err := godotenv.Load()if err != nil {
//		log.Fatalf("Failed to load env file")}

	r := SetupRouter()
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start the URLShortener server: %v", err)
	}
}
