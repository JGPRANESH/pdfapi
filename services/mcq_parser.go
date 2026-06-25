package services

import (
	"errors"
	"fmt"
	"pdfapi/models"
	"regexp"
	"strings"
)

var (
	questionRegex = regexp.MustCompile(`^(\d+)[.)]\s*(.+)$`)

	optionRegex = regexp.MustCompile(`^([A-D])[.):]\s*(.+)$`)

	answerRegex = regexp.MustCompile(
		`(?i)^(?:correct\s+)?answer\s*[:\-]?\s*(?:option\s*)?([A-D])\.?$`,
	)
)

func ParseMCQs(text string) ([]models.Question, error) {

	text = CleanExtractedText(text)

	lines := strings.Split(text, "\n")

	var questions []models.Question
	var current *models.Question

	currentField := ""

	saveCurrent := func() {
		if current == nil {
			return
		}

		current.QuestionText = cleanText(current.QuestionText)

		for i := range current.Options {
			current.Options[i] = cleanText(current.Options[i])
		}

		if err := ValidateQuestion(*current); err != nil {
			fmt.Printf("Skipping Question %s : %v\n", current.ID, err)
			return
		}

		questions = append(questions, *current)
	}

	for _, raw := range lines {

		line := strings.TrimSpace(raw)

		if line == "" {
			continue
		}

		//-----------------------------------------
		// Question Start
		//-----------------------------------------

		if m := questionRegex.FindStringSubmatch(line); m != nil {

			saveCurrent()

			current = &models.Question{
				ID:                 m[1],
				QuestionText:       m[2],
				Options:            make([]string, 4),
				CorrectOptionIndex: -1,
				Difficulty:         "Medium",
				Description:        "",
			}

			currentField = "question"
			continue
		}

		if current == nil {
			continue
		}

		//-----------------------------------------
		// Option
		//-----------------------------------------

		if m := optionRegex.FindStringSubmatch(line); m != nil {

			index := int(m[1][0] - 'A')

			current.Options[index] = m[2]

			currentField = fmt.Sprintf("option%d", index)

			continue
		}

		//-----------------------------------------
		// Answer
		//-----------------------------------------

		if m := answerRegex.FindStringSubmatch(line); m != nil {

			current.CorrectOptionIndex = int(m[1][0] - 'A')
			currentField = ""

			continue
		}

		//-----------------------------------------
		// Multi-line Support
		//-----------------------------------------
		if shouldIgnoreLine(line) {
			continue
		}
		switch currentField {

		case "question":
			current.QuestionText += " " + line

		case "option0":
			current.Options[0] += " " + line

		case "option1":
			current.Options[1] += " " + line

		case "option2":
			current.Options[2] += " " + line

		case "option3":
			current.Options[3] += " " + line
		}
	}

	saveCurrent()

	if len(questions) == 0 {
		return nil, errors.New("no valid questions found")
	}

	return questions, nil
}

// func validateQuestion(q models.Question) error {

// 	if strings.TrimSpace(q.QuestionText) == "" {
// 		return errors.New("empty question")
// 	}

// 	if len(q.Options) != 4 {
// 		return errors.New("must have exactly 4 options")
// 	}

// 	for i, option := range q.Options {

// 		if strings.TrimSpace(option) == "" {
// 			return fmt.Errorf("option %c is empty", 'A'+i)
// 		}
// 	}

// 	if q.CorrectOptionIndex < 0 || q.CorrectOptionIndex > 3 {
// 		return errors.New("invalid answer")
// 	}

// 	return nil
// }

func cleanText(text string) string {

	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.ReplaceAll(text, "\r", " ")
	text = strings.Join(strings.Fields(text), " ")

	return strings.TrimSpace(text)
}
func shouldIgnoreLine(line string) bool {

	line = strings.TrimSpace(line)

	if line == "" {
		return true
	}

	if pageFooterRegex.MatchString(line) {
		return true
	}

	if pageHeaderRegex.MatchString(line) {
		return true
	}

	if strings.Contains(line, "Rapid Revision") {
		return true
	}

	if strings.Contains(line, "CA -") {
		return true
	}

	return false
}
