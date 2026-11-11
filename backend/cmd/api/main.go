package main

import (
	"fmt"
	"log"
	"net/http"
)

func httpServer(w http.ResponseWriter, r *http.Request){
	fmt.Fprint(w, "Testing HTTP Server");
}

func main() {
	http.HandleFunc("/", httpServer)
	port := ":6000"

	fmt.Printf("Server is running on %s", port)

	err := http.ListenAndServe(port, nil);

	if err != nil {
		log.Fatalf("Server Wont Start %v", err);
	}
}