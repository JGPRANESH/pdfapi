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

	client, err := config.App.Messaging(context.Background())
	if err != nil {
		return err
	}

	message := &messaging.Message{
		Topic: "exam_news",
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
	}

	response, err := client.Send(context.Background(), message)

	fmt.Println("FCM Response:", response)
	fmt.Println("FCM Error:", err)

	return err
}
