package models

type NotificationRequest struct {
	Topic string `json:"topic"`
	Title string `json:"title"`
	Body  string `json:"body"`
}
