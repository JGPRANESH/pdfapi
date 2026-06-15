package config

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

var App *firebase.App
var FirestoreClient *firestore.Client

func InitFirebase() {
	ctx := context.Background()

	opt := option.WithCredentialsFile("serviceAccountKey.json")

	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalf("Firebase init error: %v", err)
	}

	App = app

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("Firestore init error: %v", err)
	}

	FirestoreClient = client

	log.Println("✅ Firebase & Firestore initialized")
}
