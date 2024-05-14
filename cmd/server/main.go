package main

import (
	"fmt"
	"log"
	"os"

	"github.com/danzBraham/halo-suster/internal/infrastructures/server"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Panic("Error loading .env file")
	}

	address := fmt.Sprintf("%s:%s", os.Getenv("APP_HOST"), os.Getenv("APP_PORT"))
	server := server.NewAPIServer(address, nil)
	if err := server.Launch(); err != nil {
		log.Panic(err)
	}
}
