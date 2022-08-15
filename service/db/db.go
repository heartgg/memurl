package db

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/heartgg/memurl/service/generator"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type DatabaseURL struct {
	Expiration   time.Time `firestore:"expiration"`
	MemorableURL string    `firestore:"memorable"`
	OriginalURL  string    `firestore:"original"`
}

// Initiates the firestore db and returns the firestore client
func Init() (*firestore.Client, error) {
	opt := option.WithCredentialsFile("./gcp_key.json")
	ctx := context.Background()
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, err
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// Maps the original URL to a generated URL,
func MapURL(client *firestore.Client, originalUrl string) (DatabaseURL, error) {
	data := DatabaseURL{}
	ctx := context.Background()
	qIter := client.Collection("urls").Where("original", "==", originalUrl).Documents(ctx)
	docSnap, err := qIter.Next()
	if docSnap != nil {
		docSnap.DataTo(&data)
		return data, nil
	}
	generatedUrl := generator.GenerateURL()
	data = DatabaseURL{
		Expiration:   time.Now().Add(24 * time.Hour),
		OriginalURL:  originalUrl,
		MemorableURL: generatedUrl,
	}
	_, err = client.Collection("urls").NewDoc().Set(ctx, data)
	if err != nil {
		generator.BreakURL(generatedUrl)
		return data, err
	}
	return data, nil
}

// Gets the original URL using the memorable URL
func RetrieveURL(client *firestore.Client, memorableUrl string) (string, error) {
	ctx := context.Background()
	qIter := client.Collection("urls").Where("memorable", "==", memorableUrl).Documents(ctx)
	data := DatabaseURL{}
	docSnap, err := qIter.Next()
	if docSnap != nil {
		docSnap.DataTo(&data)
		return data.OriginalURL, nil
	}
	return "", err
}

// Wipes the URL collection.
// Adjusted from Firestore docs
func WipeDB(client *firestore.Client) error {
	ctx := context.Background()
	for {
		// Get a batch of documents
		iter := client.Collection("urls").Limit(100).Documents(ctx)
		numDeleted := 0
		// Iterate through the documents, adding
		// a delete operation for each one to a
		// WriteBatch.
		batch := client.Batch()
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return err
			}

			batch.Delete(doc.Ref)
			numDeleted++
		}
		// If there are no documents to delete,
		// the process is over.
		if numDeleted == 0 {
			return nil
		}
		_, err := batch.Commit(ctx)
		if err != nil {
			return err
		}
	}
}
