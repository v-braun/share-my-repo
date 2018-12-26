// package share-my-repo ....
package main


import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/gorilla/mux"
)


func main() {
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})

	bndAddr := ":3001"
	r := mux.NewRouter()

	r.HandleFunc("/ping", pingHandler).Methods("GET")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./bin/")))

	log.WithFields(log.Fields{
		"bind-addr": bndAddr,
	}).Info("start webserver")
  
	err := http.ListenAndServe(bndAddr, r)
	if err != nil {
		panic(err)
	}
}

type ResultObj struct{
	Message string
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	log.Debug("call pingHandler")

	result := new(ResultObj)
	result.Message = "pong"

	js, err := json.Marshal(result)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}
