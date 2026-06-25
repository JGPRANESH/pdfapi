package services

import (
	"math/rand"
	"pdfapi/models"
	"strconv"
	"time"
)

const (
	ModeAll    = "all"
	ModeRandom = "random"
	ModeChunk  = "chunk"
)

func SelectQuestions(
	questions []models.Question,
	mode string,
	count int,
) []models.Question {

	if len(questions) == 0 {
		return []models.Question{}
	}

	switch mode {

	case ModeRandom:

		if count <= 0 {
			count = 20
		}

		if count > len(questions) {
			count = len(questions)
		}

		copied := make([]models.Question, len(questions))
		copy(copied, questions)

		r := rand.New(rand.NewSource(time.Now().UnixNano()))

		r.Shuffle(len(copied), func(i, j int) {
			copied[i], copied[j] = copied[j], copied[i]
		})

		selected := copied[:count]

		for i := range selected {
			selected[i].ID = strconv.Itoa(i + 1)
		}

		return selected

	case ModeAll:
		fallthrough

	default:

		copied := make([]models.Question, len(questions))
		copy(copied, questions)

		for i := range copied {
			copied[i].ID = strconv.Itoa(i + 1)
		}

		return copied
	}
}
