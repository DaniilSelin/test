package api

import (
    "github.com/gorilla/mux"

    "net/http"
)

func NewRouter(handler *Handler) *mux.Router {
    router := mux.NewRouter()
    
    router.HandleFunc("/quotes", handler.HandleGetQuoteByAuthor()).Methods("GET").Queries("author", "{author}")
    router.HandleFunc("/quotes", handler.HandleGetQuote()).Methods("GET")
    router.HandleFunc("/quotes", handler.HandlePostQuote()).Methods("POST")
    router.HandleFunc("/quotes/random", handler.HandleGetRandQuote()).Methods("GET")
    router.HandleFunc("/quotes", handler.HandleDeleteQuote()).Methods("DELETE").Queries("id", "{id}")

    return router
}
