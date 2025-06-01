package api

import (
    "github.com/gorilla/mux"
)

func NewRouter(handler *Handler) *mux.Router {
    router := mux.NewRouter()

    router.Handle("/quotes", handler.HandleGetQuoteByAuthor()).Methods("GET").Queries("author", "{author}")
    router.Handle("/quotes", handler.HandleGetQuotes()).Methods("GET")
    router.Handle("/quotes", handler.HandlePostQuote()).Methods("POST")
    router.Handle("/quotes/random", handler.HandleGetRandQuote()).Methods("GET")
    router.Handle("/quotes/{id}", handler.HandleDeleteQuote()).Methods("DELETE")

    return router
}