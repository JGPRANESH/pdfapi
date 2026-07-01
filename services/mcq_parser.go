package services

import (
	"pdfapi/models"
	"regexp"
	"strings"
)

func ParseMCQs(text string) ([]models.Question, error) {

	var questions []models.Question

	pattern := regexp.MustCompile(
		`(?s)(\d+)\.\s*(.*?)\s*A\.\s*(.*?)\s*B\.\s*(.*?)\s*C\.\s*(.*?)\s*D\.\s*(.*?)\s*(?:Answer|Correct Answer)\s*:\s*(?:Option\s*)?([A-D])`,
	)

	matches := pattern.FindAllStringSubmatch(text, -1)

	for _, m := range matches {

		q := models.Question{
			ID:           strings.TrimSpace(m[1]),
			QuestionText: cleanText(m[2]),
			Options: []string{
				cleanText(m[3]),
				cleanText(m[4]),
				cleanText(m[5]),
				cleanText(m[6]),
			},
			Difficulty: "Medium",
		}

		switch m[7] {
		case "A":
			q.CorrectOptionIndex = 0
		case "B":
			q.CorrectOptionIndex = 1
		case "C":
			q.CorrectOptionIndex = 2
		case "D":
			q.CorrectOptionIndex = 3
		}

		questions = append(questions, q)
	}

	return questions, nil
}
func cleanText(text string) string {

	text = strings.ReplaceAll(text, "\r", " ")

	text = strings.ReplaceAll(text, "\n", " ")

	text = strings.Join(strings.Fields(text), " ")

	return strings.TrimSpace(text)
}
