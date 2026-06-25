package services

import (
	"errors"
	"fmt"
	"pdfapi/models"
	"strings"
)

func ValidateQuestion(q models.Question) error {

	if strings.TrimSpace(q.QuestionText) == "" {
		return errors.New("empty question")
	}

	if len(q.Options) != 4 {
		return errors.New("must have four options")
	}

	for i, option := range q.Options {

		if strings.TrimSpace(option) == "" {
			return fmt.Errorf("option %c is empty", 'A'+i)
		}
	}

	if q.CorrectOptionIndex == -1 {
		return errors.New("missing answer")
	}

	return nil
}
