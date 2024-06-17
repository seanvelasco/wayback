package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var c = make(chan []byte)

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

func fetchResource(url string) []byte {
	res, err := http.Get(url)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("client: got response!\n")
	fmt.Printf("client: status code: %d\n", res.StatusCode)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}
	c <- body
	return body
}

func handleLanding(w http.ResponseWriter, r *http.Request) {
}

func handleAppendOrRetrieveResource(w http.ResponseWriter, r *http.Request) {
	resourceUrl := r.PathValue("resource")
	validUrl, err := validateUrl(resourceUrl)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	go fetchResource(validUrl)

	w.Write([]byte("Resource " + validUrl + " was fetched successfully."))

}

func handleRetrieveResource(w http.ResponseWriter, r *http.Request) {
}

func main() {

	go func() {
		for data := range c {
			fmt.Print(string(data))
		}
	}()

	server := &http.Server{
		Addr: ":8000",
	}
	http.HandleFunc("/", handleLanding)
	http.HandleFunc("/{timestamp}/{resource}", handleRetrieveResource)
	http.HandleFunc("/{resource}", handleAppendOrRetrieveResource)

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}

}
