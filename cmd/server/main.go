package main

import (
	"fmt"
	"log"
	"os"

	"github.com/danzBraham/halo-suster/internal/helpers"
	"github.com/danzBraham/halo-suster/internal/infrastructures/db"
	"github.com/danzBraham/halo-suster/internal/infrastructures/server"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	dbpool, err := db.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer dbpool.Close()

	helpers.NewValidate()

	address := fmt.Sprintf("%s:%s", os.Getenv("APP_HOST"), os.Getenv("APP_PORT"))
	server := server.NewAPIServer(address, dbpool)
	if err := server.Launch(); err != nil {
		log.Fatal(err)
	}
}
