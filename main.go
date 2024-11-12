package main

import (
	"log"
	"net/http"

	"github.com/IsraelTeo/api-paw/db"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {

	r := mux.NewRouter()

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loanding .env main")
	}

	if err := db.Connection(); err != nil {
		log.Fatalf("Error trying to connect with database: %v", err)
	}

	if err := db.MigrateDataBase(); err != nil {
		log.Fatalf("Error migrating database: %v", err)
	}
	log.Println("Database migration successful")

	log.Println("Starting server on port 8080...")

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
