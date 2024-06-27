package main

import (
	"Goblin/db"
	"log"
	"net/http"
)

func main() {

	server := &http.Server{
		Addr: ":8000",
	}

	http.HandleFunc("/", handleLanding)
	http.HandleFunc("/{timestamp}/{resource}", handleRetrieveResource)
	http.HandleFunc("/resource", handleAppendOrRetrieveResource)

	if err := db.Init(); err != nil {
		log.Fatal(err)
	}

	if _, err := db.DB.Exec("CREATE TABLE IF NOT EXISTS web_resources (id TEXT PRIMARY KEY, timestamp TIMESTAMP, url TEXT, content_type TEXT )"); err != nil {
		log.Fatal(err)
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
