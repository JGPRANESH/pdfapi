package models

type QuizMetadata struct {
	Title       string `json:"title" firestore:"title"`
	Explanation string `json:"explanation" firestore:"explanation"`
	Category    string `json:"category" firestore:"category"`
	Difficulty  string `json:"difficulty" firestore:"difficulty"`
	Description string `json:"description" firestore:"description"`
}
