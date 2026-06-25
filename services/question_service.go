package services

import (
	"context"
	"fmt"
	"pdfapi/models"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/google/uuid"
)

type QuizDocument struct {
	FileName        string    `firestore:"fileName"`
	CreatedAt       time.Time `firestore:"createdAt"`
	DurationMinutes int       `firestore:"durationMinutes"`
	MaxMarks        int       `firestore:"maxMarks"`
	TotalQuestions  int       `firestore:"totalQuestions"`
	Type            string    `firestore:"type"`

	Title       string `firestore:"title"`
	Explanation string `firestore:"explanation"`
	Category    string `firestore:"category"`
	Difficulty  string `firestore:"difficulty"`

	Questions []models.Question `firestore:"questions"`
}

func SaveQuiz(
	client *firestore.Client,
	fileName string,
	questions []models.Question,
	metadata *models.QuizMetadata,
) error {

	if metadata == nil {
		return fmt.Errorf("metadata is nil")
	}

	ctx := context.Background()

	quiz := QuizDocument{
		FileName:        fileName,
		CreatedAt:       time.Now(),
		DurationMinutes: 60, // Default duration of 60 minutes
		TotalQuestions:  len(questions),
		MaxMarks:        20, // Default max marks
		Type:            "mock",

		Title:       metadata.Title,
		Explanation: metadata.Description,
		Category:    metadata.Category,
		Difficulty:  metadata.Difficulty,

		Questions: questions,
	}

	quizID := uuid.New().String()

	println("Saving quiz ID:", quizID)

	_, err := client.
		Collection("quizzes").
		Doc(quizID).
		Set(ctx, quiz)

	if err != nil {
		return err
	}

	doc, err := client.
		Collection("quizzes").
		Doc(quizID).
		Get(ctx)

	if err != nil {
		println("READ FAILED:", err.Error())
	} else {
		println("DOCUMENT FOUND:", doc.Ref.ID)
	}

	println("Saved quiz ID:", quizID)

	return nil
}
