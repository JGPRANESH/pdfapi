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

	println("Saving quiz ID:", quizID)

	_, err := client.
		Collection("quizzes").
		Doc(quizID).
		Set(ctx, quiz)

	doc, err := client.
		Collection("quizzes").
		Doc(quizID).
		Get(ctx)

	if err != nil {
		println("READ FAILED:", err.Error())
	} else {
		println("DOCUMENT FOUND:", doc.Ref.ID)
	}

	if err != nil {
		return err
	}

	println("Saved quiz ID:", quizID)
	return nil
}
