package services

import (
	"context"
	"pdfapi/models"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/google/uuid"
)

type QuizDocument struct {
	FileName       string            `firestore:"fileName"`
	CreatedAt      time.Time         `firestore:"createdAt"`
	TotalQuestions int               `firestore:"totalQuestions"`
	Questions      []models.Question `firestore:"questions"`
}

func SaveQuiz(
	client *firestore.Client,
	fileName string,
	questions []models.Question,
) error {

	ctx := context.Background()

	quiz := QuizDocument{
		FileName:       fileName,
		CreatedAt:      time.Now(),
		TotalQuestions: len(questions),
		Questions:      questions,
	}
	quizID := uuid.New().String()

	_, err := client.
		Collection("quizzes").
		Doc(quizID).
		Set(ctx, quiz)

	if err != nil {
		return err
	}
	return nil
}
