package db

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/heartgg/memurl/service/generator"
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

func MapURL(client *firestore.Client, originalUrl string) (string, time.Time, error) {
	data := DatabaseURL{}
	qIter := client.Collection("urls").Where("original", "==", originalUrl).Documents(context.Background())
	docSnap, err := qIter.Next()
	if docSnap != nil {
		docSnap.DataTo(&data)
		return data.MemorableURL, data.Expiration, nil
	}
	generatedUrl := generator.GenerateURL()
	data = DatabaseURL{
		Expiration:   time.Now().Add(24 * time.Hour),
		OriginalURL:  originalUrl,
		MemorableURL: generatedUrl,
	}
	_, err = client.Collection("urls").NewDoc().Set(context.Background(), data)
	if err != nil {
		generator.BreakURL(generatedUrl)
		return "", time.Now(), err
	}
	return data.MemorableURL, data.Expiration, nil
}
