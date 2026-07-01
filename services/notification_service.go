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

func SendTokenNotification(token, title, body string) error {
	ctx := context.Background()
	client, err := config.App.Messaging(ctx)
	if err != nil {
		return fmt.Errorf("failed to create messaging client: %w", err)
	}
	message := &messaging.Message{
		Token: token,
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Data: map[string]string{
			"click_action": "FLUTTER_NOTIFICATION_CLICK",
		},
		Android: &messaging.AndroidConfig{
			Priority: "high",
		},
		APNS: &messaging.APNSConfig{
			Headers: map[string]string{
				"apns-priority": "10",
			},
		},
	}

	_, err = client.Send(ctx, message)
	if err != nil {
		return err
	}

	return nil
}
func SendTopicNotificationAsync(topic, title, body string) {

	go func() {
		err := SendTopicNotification(topic, title, body)
		if err != nil {
			fmt.Printf("Failed to send notification: %v\n", err)
		}
	}()
}

func SendTokenNotificationAsync(token, title, body string) {
	go func() {
		err := SendTokenNotification(token, title, body)
		if err != nil {
			fmt.Printf("Failed to send token notification: %v\n", err)
		}
	}()
}
