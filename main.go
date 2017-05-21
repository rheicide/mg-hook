package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8080"
	}

	log.Printf("Starting server at %s", addr)
	log.Fatal(http.ListenAndServe(addr, NewRouter()))
}
