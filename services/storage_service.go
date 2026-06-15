package services

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"

	"pdfapi/config"
)

func UploadFile(file multipart.File, fileHeader *multipart.FileHeader) (string, error) {

	storageClient, err := config.App.Storage(context.Background())

	if err != nil {
		return "", err
	}

	bucket, err := storageClient.DefaultBucket()

	if err != nil {
		return "", err
	}

	object := bucket.Object(fileHeader.Filename)

	writer := object.NewWriter(context.Background())

	_, err = io.Copy(writer, file)

	if err != nil {
		return "", err
	}

	err = writer.Close()

	if err != nil {
		return "", err
	}

	url := fmt.Sprintf(
		"https://storage.googleapis.com/%s/%s",
		bucket.BucketName(),
		fileHeader.Filename,
	)

	return url, nil
}
