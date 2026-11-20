package main

import (
	"log"
	"net/http"

	"go-lang-restapi/api"
	"go-lang-restapi/expense"
)

func main() {
	store := expense.NewMemoryStore()
	router := api.NewRouter(store)

	log.Println("Server listening on http://127.0.0.1:8000")
	log.Fatal(http.ListenAndServe("127.0.0.1:8000", router))
}
