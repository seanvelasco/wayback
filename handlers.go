package main

import (
	"Goblin/db"
	"fmt"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type WebResource struct {
	id          uuid.UUID
	timestamp   time.Time
	url         string
	contentType string
}

func validateUrl(unvalidated string) (string, error) {
	validated, err := url.ParseRequestURI(unvalidated)
	if err != nil {
		if !strings.HasPrefix(unvalidated, "http") {
			unvalidated = "https://" + unvalidated
		}
		validated, err = url.ParseRequestURI(unvalidated)
		if err != nil {
			return "", err
		}
	}
	return validated.String(), nil
}

func archiveResource(url string, resource io.Reader, contentType string) error {

	webResourceObject := WebResource{
		id:          uuid.New(),
		timestamp:   time.Now().UTC(),
		url:         url,
		contentType: contentType,
	}

	res, err := db.DB.Exec("INSERT INTO web_resources (id, timestamp, url, content_type) VALUES ($1, $2, $3, $4)",
		webResourceObject.id,
		webResourceObject.timestamp,
		webResourceObject.url,
		webResourceObject.contentType)

	if err != nil {
		log.Fatal(err)
		return err
	}

	fmt.Println(res.RowsAffected())

	fmt.Println(webResourceObject)
	// save to db
	// save to blob storage
	return nil
}

func fetchResource(url string) error {

	res, err := http.Get(url)

	if err != nil {
		return fmt.Errorf("failed to fetch resource: %w", err)
	}

	defer res.Body.Close()

	contentType := res.Header.Get("Content-Type")

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid status code: %d", res.StatusCode)
	}

	return archiveResource(url, res.Body, contentType)
}

func handleLanding(w http.ResponseWriter, r *http.Request) {

	//w.Header().Set("Content-Type", "text/html; charset=utf-8")
}

func handleAppendOrRetrieveResource(w http.ResponseWriter, r *http.Request) {
	//resourceUrl := r.PathValue("resource")

	resourceUrl := r.URL.Query().Get("url")

	validUrl, err := validateUrl(resourceUrl)

	if err != nil {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	go func() {
		err := fetchResource(validUrl)
		if err != nil {
			log.Printf("Error fetching resource: %v", err)
		}
	}()

	w.WriteHeader(http.StatusAccepted)
	//w.Write([]byte("Resource " + validUrl + " was fetched successfully."))
}

func handleRetrieveResource(w http.ResponseWriter, r *http.Request) {
}
