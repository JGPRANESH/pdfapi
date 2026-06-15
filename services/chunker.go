package services

import "strings"

func CreateChunks(text string) []string {

	const chunkSize = 200 // words per chunk

	words := strings.Fields(text)

	var chunks []string

	for i := 0; i < len(words); i += chunkSize {

		end := i + chunkSize

		if end > len(words) {
			end = len(words)
		}

		chunk := strings.Join(words[i:end], " ")

		chunks = append(chunks, chunk)
	}

	return chunks
}
