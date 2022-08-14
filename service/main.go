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
	r.HandleFunc("/get_url", getUrl)
	r.HandleFunc("/{link}", redirect)
	http.ListenAndServe(":3000", r)
}

func getUrl(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	w.Header().Set("Content-Type", "application/json")
	memorableUrl, expiration, err := db.MapURL(client, r.Form["user-url"][0])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Panicln(err)
	}
	w.WriteHeader(http.StatusCreated)
	data := ResponseURL{
		Expiration:   expiration,
		MemorableURL: memorableUrl,
	}
	resp, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Panicln(err)
	}
	w.Write(resp)
}

func redirect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	link, ok := vars["link"]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
	}
	originalUrl, err := db.RetrieveURL(client, link)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	}
	http.Redirect(w, r, originalUrl, http.StatusSeeOther)
}
