package main

import (
	"context"
	"log"
	"time"

	firebase "firebase.google.com/go"

	"google.golang.org/api/option"
)

type URL struct {
	CreatedOn    time.Time `firestore:"created_on"`
	MemorableURL string    `firestore:"memorable"`
	OriginalURL  string    `firestore:"original"`
}

func main() {
	// TODO change this when deploying to firebase
	opt := option.WithCredentialsFile("../memurl/serviceAccountKey.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)

	client, err := app.Firestore(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

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
