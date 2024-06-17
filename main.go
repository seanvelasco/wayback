package main

import (
	"fmt"
	"net/http"
)

func handleLanding(w http.ResponseWriter, r *http.Request) {
}

func handleAppendOrRetrieveWebResource(w http.ResponseWriter, r *http.Request) {
}

func handleRetrieveResource(w http.ResponseWriter, r *http.Request) {
}

func main() {
	server := &http.Server{
		Addr: ":8000",
	}
	http.HandleFunc("/", handleLanding)
	http.HandleFunc("/{timestamp}/{resource}", handleRetrieveResource)
	http.HandleFunc("/{resource}", handleAppendOrRetrieveWebResource)

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}
