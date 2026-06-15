package models

type QuizMetadata struct {
	Title       string `json:"title" firestore:"title"`
	Description string `json:"description" firestore:"description"`
	Category    string `json:"category" firestore:"category"`
	Difficulty  string `json:"difficulty" firestore:"difficulty"`
}
