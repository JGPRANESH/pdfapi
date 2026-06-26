package services

import "pdfapi/models"

func GenerateQuestionChunks(
	questions []models.Question,
	size int,
	minLastChunkSize int,
	mergeSmallLastChunk bool,
) [][]models.Question {

	if size <= 0 {
		size = 20
	}

	if minLastChunkSize <= 0 {
		minLastChunkSize = size / 2
	}

	var result [][]models.Question

	for i := 0; i < len(questions); {

		remaining := len(questions) - i

		// Last chunk
		if remaining <= size {

			// Merge the last chunk into the previous one
			if mergeSmallLastChunk &&
				remaining < minLastChunkSize &&
				len(result) > 0 {

				result[len(result)-1] = append(
					result[len(result)-1],
					questions[i:]...,
				)

			} else {

				// Keep the last chunk as it is
				result = append(
					result,
					questions[i:],
				)
			}

			break
		}

		// Normal chunk
		result = append(
			result,
			questions[i:i+size],
		)

		i += size
	}

	return result
}
