package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"playhouse-server/repository"
)

func main() {
	// Test if server works
	//r := chi.NewRouter()
	//r.Use(middleware.Logger)
	//r.Get("/", func(w http.ResponseWriter, r *http.Request) {
	//	w.Write([]byte("Hello World!"))
	//})
	//http.ListenAndServe(":2345", r)

	if err := godotenv.Load("../conf/.env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	db := repository.NewDB()
	// Execute a raw SQL command.
	result := db.Connection.Raw("select 1;")
	if result.Error != nil {
		log.Fatalf("failed to execute raw SQL: %v", result.Error)
	}

	fmt.Println("Database connection successful.")
}
