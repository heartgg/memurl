package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/heartgg/memurl/service/db"
)

type ResponseURL struct {
	Expiration   time.Time `json:"expiration"`
	MemorableURL string    `json:"url"`
}

func main() {
	client, err := db.Init()
	if err != nil {
		log.Fatalln(err)
	}
	doc, err := client.Collection("urls").Doc("hEgtQQqS7yjymkqF8ZHk").Get(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
	urlDoc := db.DatabaseURL{}
	if err := doc.DataTo(&urlDoc); err != nil {
		log.Fatalln(err)
	}
	log.Println(urlDoc)

	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/generate_url", generateUrl)
	http.ListenAndServe(":3000", nil)
}

func generateUrl(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	log.Println(r.Form)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	data, err := json.Marshal(ResponseURL{Expiration: time.Now().Add(24 * time.Hour), MemorableURL: "test"})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Panicln(err)
	}
	w.Write(data)
}
