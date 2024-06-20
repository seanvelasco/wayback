package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

//var c = make(chan []byte)

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
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
	}
	//c <- body
	return body
}

func handleLanding(w http.ResponseWriter, r *http.Request) {

	//w.Header().Set("Content-Type", "text/html; charset=utf-8")
}

func handleAppendOrRetrieveResource(w http.ResponseWriter, r *http.Request) {
	resourceUrl := r.PathValue("resource")
	validUrl, err := validateUrl(resourceUrl)
	if err != nil {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
	}
	go fetchResource(validUrl)

	w.Write([]byte("Resource " + validUrl + " was fetched successfully."))

}

func handleRetrieveResource(w http.ResponseWriter, r *http.Request) {
}

func main() {

	//go func() {
	//	for data := range c {
	//		fmt.Print(string(data))
	//	}
	//}()

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
