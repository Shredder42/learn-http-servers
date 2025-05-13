package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	const port = "8080"

	serverMux := http.NewServeMux()

	server := &http.Server{
		Addr:    ":" + port,
		Handler: serverMux,
	}
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}

	log.Printf("Serving on port: %s\n", port)

}
