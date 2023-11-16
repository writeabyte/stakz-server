package main

import (
	"fmt"
	"log"
	"net/http"
)

interface{}

func main() {
	// Define a function to handle HTTP requests
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
    if r. {
      
    }
	})

	// Start the server and listen on port 8080
	log.Fatalln(http.ListenAndServe(":8080", nil))
	fmt.Println("Server listening on :8080")
}
