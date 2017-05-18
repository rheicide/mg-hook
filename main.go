package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	router := NewRouter()
	addr := os.Getenv("ADDR")

	if addr == "" {
		addr = ":8080"
	}

	log.Printf("Listening at %s", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}
