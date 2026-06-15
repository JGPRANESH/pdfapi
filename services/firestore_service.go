package services

import (
	"context"

	"pdfapi/config"
	"pdfapi/models"
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
