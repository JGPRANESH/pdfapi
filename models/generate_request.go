package models

type GenerateQuizRequest struct {
	Topic      string `json:"topic"`
	Content    string `json:"content"`
	Difficulty string `json:"difficulty"`
	Count      int    `json:"count"`
}
