package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/tajimyradov/quotes-api/docs"
	"github.com/tajimyradov/quotes-api/handlers"
	"github.com/tajimyradov/quotes-api/storage"
)

func main() {
	db := storage.NewMemoryStorage()
	handler := handlers.NewQuoteHandler(db)

	r := mux.NewRouter()

	r.HandleFunc("/quotes", handler.CreateQuote).Methods("POST")
	r.HandleFunc("/quotes", handler.GetAllQuotes).Methods("GET")
	r.HandleFunc("/quotes/random", handler.GetRandomQuote).Methods("GET")
	r.HandleFunc("/quotes/{id:[0-9]+}", handler.DeleteQuote).Methods("DELETE")

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
