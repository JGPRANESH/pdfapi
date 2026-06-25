package services

import (
	"regexp"
	"strings"
)

var questionBoundaryRegex = regexp.MustCompile(`(?m)^\s*\d+[.)]\s+`)

func CreateChunks(text string) []string {

	text = CleanExtractedText(text)

	matches := questionBoundaryRegex.FindAllStringIndex(text, -1)

	if len(matches) == 0 {
		return []string{text}
	}

	var questions []string

	for i := range matches {

		start := matches[i][0]

		var end int

		if i == len(matches)-1 {
			end = len(text)
		} else {
			end = matches[i+1][0]
		}

		q := strings.TrimSpace(text[start:end])

		if q != "" {
			questions = append(questions, q)
		}
	}

	const questionsPerChunk = 4

	var chunks []string

	for i := 0; i < len(questions); i += questionsPerChunk {

		end := i + questionsPerChunk

		if end > len(questions) {
			end = len(questions)
		}

		chunks = append(
			chunks,
			strings.Join(questions[i:end], "\n\n"),
		)
	}

	return chunks
}
