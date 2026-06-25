package services

import "pdfapi/models"

func GenerateQuestionChunks(
	questions []models.Question,
	size int,
) [][]models.Question {

	if size <= 0 {
		size = 20
	}

	var result [][]models.Question

	for i := 0; i < len(questions); i += size {

		end := i + size

		if end > len(questions) {
			end = len(questions)
		}

		result = append(result, questions[i:end])
	}

	return result
}
