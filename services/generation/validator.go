package generation

import "pdfapi/models"

func ValidateQuestions(
	questions []models.Question,
) []models.Question {

	var valid []models.Question

	for _, q := range questions {

		if len(q.Options) != 4 {
			continue
		}

		valid = append(valid, q)
	}

	return valid
}
