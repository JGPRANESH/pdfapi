package config

import (
	"context"
	"encoding/base64"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

var App *firebase.App
var FirestoreClient *firestore.Client

func InitFirebase() {
	ctx := context.Background()

	projectID := os.Getenv("FIREBASE_PROJECT_ID")
	if projectID == "" {
		log.Fatal("FIREBASE_PROJECT_ID is not set")
	}

	encodedCreds := os.Getenv("FIREBASE_CREDENTIALS_BASE64")
	if encodedCreds == "" {
		log.Fatal("FIREBASE_CREDENTIALS_BASE64 is not set")
	}

	credJSON, err := base64.StdEncoding.DecodeString(encodedCreds)
	if err != nil {
		log.Fatalf("Failed to decode Firebase credentials: %v", err)
	}

	conf := &firebase.Config{
		ProjectID: projectID,
	}

	opt := option.WithCredentialsJSON(credJSON)

	app, err := firebase.NewApp(ctx, conf, opt)
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
