package services

import (
	"context"
	"time"

	"pdfapi/config"
	"pdfapi/models"

	"cloud.google.com/go/firestore"
)

func SaveMetadata(data models.FileMetadata) error {

	client, err := config.App.Firestore(context.Background())

	if err != nil {
		return err
	}

	_, _, err = client.Collection("files").Add(
		context.Background(),
		data,
	)

	return err
}
func SavePDFContent(
	client *firestore.Client,
	fileName string,
	content string,
) error {

	ctx := context.Background()

	doc := map[string]interface{}{
		"fileName":  fileName,
		"content":   content,
		"createdAt": time.Now(),
	}

	_, _, err := client.
		Collection("pdf_contents").
		Add(ctx, doc)

	return err
}
