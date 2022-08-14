package db

import (
	"context"
	"errors"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type DatabaseURL struct {
	Expiration   time.Time `firestore:"expiration"`
	MemorableURL string    `firestore:"memorable"`
	OriginalURL  string    `firestore:"original"`
}

func Init() (*firestore.Client, error) {
	opt := option.WithCredentialsFile("./gcp_key.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, err
	}

	client, err := app.Firestore(context.Background())
	if err != nil {
		return nil, err
	}
	return client, nil
}

func MapURLs(client *firestore.Client, originalUrl string, generatedUrl string) (string, error) {
	qIter := client.Collection("urls").Where("memorable", "==", generatedUrl).Documents(context.Background())
	_, err := qIter.Next()
	if err != iterator.Done {
		// MapURLs(client, originalUrl, generatedUrl)
		return "", errors.New("Duplicate generated URL found")
	}
	docSnap, err := client.Collection("urls").NewDoc().Set(context.Background(), DatabaseURL{
		Expiration:   time.Now().Add(24 * time.Hour),
		OriginalURL:  originalUrl,
		MemorableURL: generatedUrl})
	if err != nil {
		return "", err
	}
	log.Println(docSnap)
	return generatedUrl, err
}
