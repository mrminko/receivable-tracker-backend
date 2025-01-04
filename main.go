package main

import (
	"database/sql"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/mrminko/receivable-tracker/internal/database"
	"log"
	"net/http"
	"os"
)

type DBQuery struct {
	db *database.Queries
}

func main() {
	fmt.Println("Initializing...")
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error when loading env file.")
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatalf("PORT is not defined")
	}

	dbURL := os.Getenv("GOOSE_DBSTRING")
	if dbURL == "" {
		log.Fatal("GOOSE_DBSTRING not found in env.")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalln("Error when connecting with database.", err)
	}
	fmt.Println("Database connection established.")
	defer conn.Close()

	db := database.New(conn)
	Query := DBQuery{
		db: db,
	}

	router := chi.NewRouter()
	router.Get("/users/all", Query.getAllUsers)
	router.Post("/users/create", Query.createUser)
	router.Post("/users/delete/{userId}", Query.deleteUser)
	router.Post("/users/update/{userId}", Query.updateUser)

	router.Get("/receivables/all", Query.getAllReceivables)
	router.Post("/receivables/create", Query.createReceivable)
	router.Post("/receivables/delete/{receivableId}", Query.deleteReceivable)
	router.Post("/receivables/update/{receivableId}", Query.updateReceivable)

	srv := http.Server{
		Handler: router,
		Addr:    ":" + port,
	}
	fmt.Printf("Server listening on port %v", port)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatalf("Error when initiating server: %v", err)
	}
}
