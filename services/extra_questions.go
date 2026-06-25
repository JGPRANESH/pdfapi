package services

import (
	"math/rand"
	"pdfapi/models"
	"strconv"
)

func GetRandomExtraQuestions(
	allQuestions []models.Question,
	count int,
) []models.Question {

	if count <= 0 {
		return []models.Question{}
	}

	if count > len(allQuestions) {
		count = len(allQuestions)
	}

	questions := make([]models.Question, len(allQuestions))
	copy(questions, allQuestions)

	rand.Shuffle(len(questions), func(i, j int) {
		questions[i], questions[j] = questions[j], questions[i]
	})

	extra := questions[:count]

	for i := range extra {
		extra[i].ID = strconv.Itoa(i + 1)
	}

	return extra
}
