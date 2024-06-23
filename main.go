package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

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

func archiveResource(id string, resource io.Reader) error {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	fmt.Println(timestamp)
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

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid status code: %d", res.StatusCode)
	}

	return archiveResource(url, res.Body)
}

func handleLanding(w http.ResponseWriter, r *http.Request) {

	//w.Header().Set("Content-Type", "text/html; charset=utf-8")
}

func handleAppendOrRetrieveResource(w http.ResponseWriter, r *http.Request) {
	resourceUrl := r.PathValue("resource")
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

func main() {

	server := &http.Server{
		Addr: ":8000",
	}
	http.HandleFunc("/", handleLanding)
	http.HandleFunc("/{timestamp}/{resource}", handleRetrieveResource)
	http.HandleFunc("/{resource}", handleAppendOrRetrieveResource)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
