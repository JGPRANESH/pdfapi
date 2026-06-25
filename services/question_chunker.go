package services

import "pdfapi/models"

func GenerateChunks(
	questions []models.Question,
	size int,
) [][]models.Question {

	if size <= 0 {
		size = 20
	}

	var chunks [][]models.Question

	for i := 0; i < len(questions); i += size {

		end := i + size

		if end > len(questions) {
			end = len(questions)
		}

		chunks = append(
			chunks,
			questions[i:end],
		)
	}

	return chunks
}
