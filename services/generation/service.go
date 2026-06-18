package generation

import (
	"encoding/json"
	"fmt"

	"pdfapi/models"
	"pdfapi/services"
)

func GenerateQuiz(
	req models.GenerateQuizRequest,
) ([]models.Question, error) {

	prompt := BuildPrompt(req)

	response, err := services.GenerateContent(prompt)
	if err != nil {
		return nil, err
	}

	var aiQuestions []models.AIQuestion

	err = json.Unmarshal(
		[]byte(response),
		&aiQuestions,
	)

	if err != nil {
		return nil, err
	}

	var questions []models.Question

	for i, q := range aiQuestions {

		correctIndex := 0

		switch q.Answer {
		case "A":
			correctIndex = 0
		case "B":
			correctIndex = 1
		case "C":
			correctIndex = 2
		case "D":
			correctIndex = 3
		}

		questions = append(
			questions,
			models.Question{
				ID:                 fmt.Sprintf("%d", i+1),
				QuestionText:       q.Question,
				Options:            q.Options,
				CorrectOptionIndex: correctIndex,
				Difficulty:         req.Difficulty,
				Description:        q.Description,
			},
		)
	}

	return questions, nil
}
