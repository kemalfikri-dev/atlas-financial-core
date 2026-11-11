package main

import (
	"fmt"
	"log"
	"net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, World!")
}

func main() {
	http.HandleFunc("/", homeHandler)

	port := ":6000"
	fmt.Printf("Server is running on %s", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server Wont Start: %v", err)
	}
}
