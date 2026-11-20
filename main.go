package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"go-lang-restapi/api"
	"go-lang-restapi/database"
	"go-lang-restapi/expense"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	db, err := database.OpenPostgres(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		return err
	}
	defer db.Close()

	log.Println("PostgreSQL connection verified")

	store := expense.NewMemoryStore()
	router := api.NewRouter(store)

	log.Println("Server listening on http://127.0.0.1:8000")
	return http.ListenAndServe("127.0.0.1:8000", router)
}
