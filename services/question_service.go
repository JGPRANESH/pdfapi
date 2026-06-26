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
	FileName                  string    `firestore:"fileName"`
	CreatedAt                 time.Time `firestore:"createdAt"`
	DurationMinutes           int       `firestore:"durationMinutes"`
	MaxMarks                  int       `firestore:"maxMarks"`
	TotalQuestions            int       `firestore:"totalQuestions"`
	Type                      string    `firestore:"type"`
	Description               string    `firestore:"description"`
	EachQuestionMarks         int       `firestore:"eachQuestionMarks"`
	EachQuestionNegativeMarks int       `firestore:"eachQuestionNegativeMarks"`

	Title       string `firestore:"title"`
	Explanation string `firestore:"explanation"`
	Category    string `firestore:"category"`
	Difficulty  string `firestore:"difficulty"`
	ExamName    string `firestore:"examName"`

	Questions []models.Question `firestore:"questions"`
}

func SaveQuiz(
	client *firestore.Client,
	fileName string,
	questions []models.Question,
	metadata *models.QuizMetadata,
	examName string,
) error {

	if metadata == nil {
		return fmt.Errorf("metadata is nil")
	}

	ctx := context.Background()

	quiz := QuizDocument{
		FileName:  fileName,
		CreatedAt: time.Now(),
		DurationMinutes: CalculateDuration(
			metadata.Category,
			metadata.Difficulty,
			len(questions),
		),
		TotalQuestions:            len(questions),
		MaxMarks:                  len(questions),
		Type:                      "mock",
		Description:               metadata.Description,
		ExamName:                  examName,
		EachQuestionMarks:         1,
		EachQuestionNegativeMarks: -1,

		Title:       metadata.Title,
		Explanation: metadata.Explanation,
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
