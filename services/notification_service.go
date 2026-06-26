package services

import (
	"context"
	"fmt"

	"pdfapi/config"

	"firebase.google.com/go/messaging"
)

func SendTopicNotification(
	topic string,
	title string,
	body string,
) error {

	ctx := context.Background()

	client, err := config.App.Messaging(ctx)
	if err != nil {
		return fmt.Errorf("failed to create messaging client: %w", err)
	}

	message := &messaging.Message{
		Topic: topic,
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
	}

	_, err = client.Send(ctx, message)
	if err != nil {
		return fmt.Errorf("failed to send notification: %w", err)
	}

	return nil
}
