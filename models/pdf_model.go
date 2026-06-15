package models

import "time"

type PDFDocument struct {
	FileName  string    `json:"fileName" firestore:"fileName"`
	Content   string    `json:"content" firestore:"content"`
	CreatedAt time.Time `json:"createdAt" firestore:"createdAt"`
}
