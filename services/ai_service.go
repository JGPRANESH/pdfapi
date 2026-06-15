package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"pdfapi/models"
	"strings"
	"time"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type GroqRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type GroqResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`

	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// type Question struct {
// 	ID                 string   `json:"id"`
// 	QuestionText       string   `json:"questionText"`
// 	Options            []string `json:"options"`
// 	CorrectOptionIndex int      `json:"correctOptionIndex"`
// 	Difficulty         string   `json:"difficulty"`

// }

type QuestionsResponse struct {
	Questions []models.Question `json:"questions"`
}

func ParseQuestions(jsonString string) ([]models.Question, error) {

	var response QuestionsResponse

	err := json.Unmarshal([]byte(jsonString), &response)
	if err != nil {
		return nil, err
	}

	return response.Questions, nil

}

func GenerateQuestions(documentText string) ([]models.Question, error) {

	if strings.TrimSpace(documentText) == "" {
		return nil, fmt.Errorf("document text is empty")
	}

	log.Printf("Document Length=%d", len(documentText))

	chunks := CreateChunks(documentText)

	log.Printf(
		"Question Generation Chunks=%d",
		len(chunks),
	)

	// For now process only first 3 chunks.
	// We will improve this later.
	if len(chunks) > 10 {
		chunks = chunks[:10]
	}

	var allQuestions []models.Question

	for i, chunk := range chunks {

		log.Printf(
			"Generating questions from chunk %d/%d",
			i+1,
			len(chunks),
		)

		response, err := GenerateQuestionsFromText(chunk)

		if err != nil {

			log.Printf(
				"Chunk %d failed: %v",
				i+1,
				err,
			)

			continue
		}

		questions, err := ParseQuestions(response)

		if err != nil {

			log.Printf(
				"Chunk %d parse failed: %v",
				i+1,
				err,
			)

			continue
		}

		allQuestions = append(
			allQuestions,
			questions...,
		)
	}

	if len(allQuestions) == 0 {
		return nil, fmt.Errorf(
			"no questions generated from any chunk",
		)
	}

	// Keep only 10 questions
	if len(allQuestions) > 15 {
		allQuestions = allQuestions[:15]
	}
	return allQuestions, nil
}

func GenerateQuestionsFromText(text string) (string, error) {

	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("GROQ_API_KEY not found")
	}

	url := "https://api.groq.com/openai/v1/chat/completions"

	prompt := `
Generate exactly 2 unique MCQs from the provided content
Return only JSON:

{
  "questions":[
    {
      "id":"1",
      "questionText":"Question?",
      "options":["A","B","C","D"],
      "correctOptionIndex":0,
      "difficulty":"medium"
    }
  ]
}

CONTENT:
` + text

	requestBody := GroqRequest{
		Model: "openai/gpt-oss-20b",
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		url,
		bytes.NewBuffer(jsonData),
	)

	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 90 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	log.Printf(
		"Groq Status=%d",
		resp.StatusCode,
	)

	if resp.StatusCode != http.StatusOK {

		return "", fmt.Errorf(
			"groq request failed status=%d body=%s",
			resp.StatusCode,
			string(body),
		)
	}

	var groqResp GroqResponse

	err = json.Unmarshal(body, &groqResp)
	if err != nil {
		return "", err
	}

	if groqResp.Error != nil {
		return "", fmt.Errorf(
			"groq error: %s",
			groqResp.Error.Message,
		)
	}

	if len(groqResp.Choices) == 0 {
		return "", fmt.Errorf(
			"groq returned zero choices: %s",
			string(body),
		)
	}

	content := groqResp.Choices[0].Message.Content

	content = strings.ReplaceAll(content, "```json", "")
	content = strings.ReplaceAll(content, "```", "")
	content = strings.TrimSpace(content)

	return content, nil

}
