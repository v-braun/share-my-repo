// package share-my-repo ....
package main

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/v-braun/hero-scrape"
	"github.com/v-braun/share-my-repo/strategy"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})

	bndAddr := ":3001"
	r := mux.NewRouter()

	r.HandleFunc("/{user}/{repo}", scrapeHandler).Methods("GET")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./bin/")))

	log.WithFields(log.Fields{
		"bind-addr": bndAddr,
	}).Info("start webserver")

	err := http.ListenAndServe(bndAddr, r)
	if err != nil {
		panic(err)
	}
}

func scrapeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if vars == nil {
		http.NotFound(w, r)
		return
	}

	usr, ok := vars["user"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	repo, ok := vars["repo"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	u := "https://github.com/" + usr + "/" + repo
	parsedURL, err := url.Parse(u)
	if parsedURL == nil || err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res, err := http.Get(parsedURL.String())
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer res.Body.Close()

	result, _ := heroscrape.ScrapeWithStrategy(parsedURL, res.Body, strategy.NewGitHubStrategy())
	if result == nil {
		http.Redirect(w, r, u, http.StatusTemporaryRedirect)
		return
	}

	js, err := json.Marshal(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}
