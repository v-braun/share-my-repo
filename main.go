// package share-my-repo ....
package main

import (
	"encoding/json"
	"html/template"
	"net/http"
	"net/url"

	"github.com/v-braun/hero-scrape"
	"github.com/v-braun/share-my-repo/strategy"
	"github.com/v-braun/share-my-repo/tpl"

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

	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./bin/assets/"))))
	r.HandleFunc("/api/{user}/{repo}", scrapeApiHandler).Methods("GET")
	r.HandleFunc("/{user}/{repo}", scrapeHtmlHandler).Methods("GET")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./bin/")))

	log.WithFields(log.Fields{
		"bind-addr": bndAddr,
	}).Info("start webserver")

	err := http.ListenAndServe(bndAddr, r)
	if err != nil {
		panic(err)
	}
}

func scrape(w http.ResponseWriter, r *http.Request) (user string, repo string, result *strategy.GHStrategy) {
	vars := mux.Vars(r)
	if vars == nil {
		http.NotFound(w, r)
		return "", "", nil
	}

	user, ok := vars["user"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return "", "", nil
	}

	repo, ok = vars["repo"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return "", "", nil
	}

	u := "https://github.com/" + user + "/" + repo
	parsedURL, err := url.Parse(u)
	if parsedURL == nil || err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return "", "", nil
	}

	res, err := http.Get(parsedURL.String())
	if err != nil {
		http.NotFound(w, r)
		return "", "", nil
	}
	defer res.Body.Close()

	gh := strategy.NewGitHubStrategy()
	strat, _ := heroscrape.ScrapeWithStrategy(parsedURL, res.Body, gh)
	if strat == nil {
		http.Redirect(w, r, u, http.StatusTemporaryRedirect)
		return "", "", nil
	}

	return user, repo, gh
}

func scrapeApiHandler(w http.ResponseWriter, r *http.Request) {
	_, _, gh := scrape(w, r)
	js, err := json.Marshal(gh)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func scrapeHtmlHandler(w http.ResponseWriter, r *http.Request) {
	user, repo, gh := scrape(w, r)
	model := tpl.NewModel(user, repo, gh.EndResult)

	tmpl, err := template.New("gh-redirect").Parse(tpl.GetTemplate())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(w, model)
}
