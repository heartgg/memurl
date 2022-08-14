package main

import (
	"context"
	"log"
	"time"

	"github.com/heartgg/memurl/service/db"
)

type URL struct {
	CreatedOn    time.Time `firestore:"created_on"`
	MemorableURL string    `firestore:"memorable"`
	OriginalURL  string    `firestore:"original"`
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
	urlDoc := URL{}
	if err := doc.DataTo(&urlDoc); err != nil {
		log.Fatalln(err)
	}
	log.Println(urlDoc)
}
