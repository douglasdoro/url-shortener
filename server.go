package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/douglasdoro/url-shortener/url"
)

var (
	port    int
	baseUrl string
)

type Headers map[string]string

func init() {
	port = 3000
	baseUrl = fmt.Sprintf("http://localhost:%d", port)

	url.ConfigRepository(url.NewRopositoryInMemory())
}

func main() {
	http.HandleFunc("/api/tiny", UrlShortener)
	http.HandleFunc("/r/", Redirect)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func UrlShortener(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		respondWith(w, http.StatusMethodNotAllowed, Headers{"Allow": "POST"})
	}

	url, new, err := url.FindOrCreateUrl(extractUrl(r))
	if err != nil {
		respondWith(w, http.StatusBadRequest, nil)
		return
	}

	var status int
	if new {
		status = http.StatusCreated
	} else {
		status = http.StatusOK
	}

	shortUrl := fmt.Sprintf("%s/r/%s", baseUrl, url.Id)

	respondWith(w, status, Headers{"Location": shortUrl})

	return
}

func Redirect(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	urlId := path[len(path)-1]

	if url := url.Find(urlId); url != nil {
		http.Redirect(w, r, url.Destination, http.StatusMovedPermanently)
	} else {
		http.NotFound(w, r)
	}
}

func extractUrl(r *http.Request) string {
	url := make([]byte, r.ContentLength, r.ContentLength)
	r.Body.Read(url)
	return string(url)
}

func respondWith(w http.ResponseWriter, status int, headers Headers) {
	for k, v := range headers {
		w.Header().Set(k, v)
	}
	w.WriteHeader(status)
}
