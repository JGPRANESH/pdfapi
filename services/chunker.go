package services

import (
	"fmt"
	"pdfapi/config"
	"pdfapi/models"
	"strings"
)

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

//

func SaveQuizVariants(
	fileName string,
	examName string,
	questions []models.Question,
) {

	chunkSizes := []int{10, 5}

	for _, size := range chunkSizes {

		chunks := GenerateQuestionChunks(
			questions,
			size,
			size/2,
			true,
		)

		for i, chunk := range chunks {

			metadata, err := GenerateQuizMetadata(chunk)
			if err != nil {
				continue
			}

			variantName := fmt.Sprintf(
				"%s_%d_part_%d",
				fileName,
				size,
				i+1,
			)

			err = SaveQuiz(
				config.FirestoreClient,
				variantName,
				chunk,
				metadata,
				examName,
			)

			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
