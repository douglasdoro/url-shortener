package main

import (
	"encoding/json"
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

type Redirect struct {
	stats chan string
}

func init() {
	port = 3000
	baseUrl = fmt.Sprintf("http://localhost:%d", port)

	url.ConfigRepository(url.NewRopositoryInMemory())
}

func main() {
	stats := make(chan string)
	defer close(stats)
	go recordStatistics(stats)

	http.HandleFunc("/api/tiny", UrlShortener)
	http.HandleFunc("/api/stats/", ShowStats)
	http.Handle("/r/", &Redirect{stats: stats})

	logger("Initialize the server on port %d", port)
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

	logger("%s was shorted to %s with success", url.Destination, shortUrl)

	respondWith(
		w,
		status,
		Headers{
			"Location": shortUrl,
			"Link":     fmt.Sprintf("<%s/api/stats/%s>; rel=\"stats\"", baseUrl, url.Id)})

}

func (red *Redirect) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	findUrlAndExecute(w, r, func(url *url.Url) {
		http.Redirect(w, r, url.Destination, http.StatusMovedPermanently)

		red.stats <- url.Id
	})
}

func ShowStats(w http.ResponseWriter, r *http.Request) {
	findUrlAndExecute(w, r, func(url *url.Url) {
		json, err := json.Marshal(url.Stats())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		respondWithJSON(w, string(json), Headers{"Content-Type": "application/json"})
	})
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

func respondWithJSON(w http.ResponseWriter, responseJSON string, headers Headers) {
	for k, v := range headers {
		w.Header().Set(k, v)
	}

	fmt.Fprintf(w, responseJSON)
}

func recordStatistics(ids <-chan string) {
	for id := range ids {
		url.RegisterClick(id)
		logger("Click recorded with success. ID: %s \n", id)
	}
}

func findUrlAndExecute(w http.ResponseWriter, r *http.Request, executor func(*url.Url)) {
	path := strings.Split(r.URL.Path, "/")
	id := path[len(path)-1]

	if url := url.Find(id); url != nil {
		executor(url)
	} else {
		http.NotFound(w, r)
	}
}

func logger(format string, values ...interface{}) {
	log.Printf(fmt.Sprintf("%s\n", format), values...)
}
