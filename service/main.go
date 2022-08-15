package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gorilla/mux"
	"github.com/heartgg/memurl/service/db"
	"github.com/heartgg/memurl/service/generator"
)

type ResponseURL struct {
	Expiration   time.Time `json:"expiration"`
	MemorableURL string    `json:"url"`
}

var client *firestore.Client

func main() {
	var err error
	client, err = db.Init()
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()
	if err := generator.LoadWords(); err != nil {
		log.Fatalln(err)
	}

	r := mux.NewRouter()

	r.Handle("/", http.FileServer(http.Dir("./static")))
	r.HandleFunc("/get_url", getUrlHandler)
	r.HandleFunc("/u/{link}", redirectHandler)
	http.ListenAndServe(":3000", r)
}

// Handler that maps user's url to a generated url and returns the generated url and expiration date
func getUrlHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	w.Header().Set("Content-Type", "application/json")
	url := r.Form["user-url"][0]
	if len(url) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	databaseUrl, err := db.MapURL(client, url)
	if err != nil {
		log.Printf("[getUrl handler] MapURL threw error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data := ResponseURL{
		Expiration:   databaseUrl.Expiration,
		MemorableURL: databaseUrl.MemorableURL,
	}
	resp, err := json.Marshal(data)
	if err != nil {
		log.Printf("[getUrl handler] Data marshal threw error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
	return
}

// Handler that redirects using the memorable link to original link
func redirectHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	link, ok := vars["link"]
	if !ok {
		log.Printf("[redirect handler] Issue getting memorable link: %v", vars)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	originalUrl, err := db.RetrieveURL(client, link)
	if err != nil {
		log.Printf("[redirect handler] RetrieveURL threw error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, originalUrl, http.StatusSeeOther)
}
