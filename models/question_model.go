package models

type Question struct {
	ID                 string   `json:"id" firestore:"id"`
	QuestionText       string   `json:"questionText" firestore:"questionText"`
	Options            []string `json:"options" firestore:"options"`
	CorrectOptionIndex int      `json:"correctOptionIndex" firestore:"correctOptionIndex"`
	Difficulty         string   `json:"difficulty" firestore:"difficulty"`
}
