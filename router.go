package main

import "github.com/gorilla/mux"

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	for _, route := range routes {
		handler := Logger(route.HandlerFunc)
		router.Methods(route.Method).Path(route.Pattern).Handler(handler)
	}

	return router
}
