package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Hanya load .env jika bukan di production
	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	NewServer().Run()
}
