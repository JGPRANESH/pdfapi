package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

func GenerateQuizMetadata(
	questions []models.Question,
) (*models.QuizMetadata, error) {

	var questionText strings.Builder

	for i, q := range questions {
		questionText.WriteString(
			fmt.Sprintf(
				"%d. %s\n",
				i+1,
				q.QuestionText,
			),
		)
	}

	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GROQ_API_KEY not found")
	}

	url := "https://api.groq.com/openai/v1/chat/completions"

	prompt := `
Analyze the following quiz questions and generate quiz metadata.

Generate:
1. Title
2. Short Description
3. Category
4. Difficulty (Easy, Medium, Hard)

Category MUST be exactly ONE of these values:
- GK
- CA
- Reasoning
- English
- Science
- Maths

Rules:
- Choose the category based only on the provided PDF content.
- Do NOT invent a new category.
- If multiple categories seem applicable, select the single best match.
- Return ONLY one of the six allowed category values.

Return ONLY valid JSON:

{
  "title": "",
  "description": "",
  "category": "",
  "difficulty": ""
}

QUESTIONS:

` + questionText.String()

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
		return nil, err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		url,
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set(
		"Authorization",
		"Bearer "+apiKey,
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	client := &http.Client{
		Timeout: 90 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"groq request failed status=%d body=%s",
			resp.StatusCode,
			string(body),
		)
	}

	var groqResp GroqResponse

	err = json.Unmarshal(body, &groqResp)
	if err != nil {
		return nil, err
	}

	if groqResp.Error != nil {
		return nil, fmt.Errorf(
			"groq error: %s",
			groqResp.Error.Message,
		)
	}

	if len(groqResp.Choices) == 0 {
		return nil, fmt.Errorf(
			"no response from Groq",
		)
	}

	content := groqResp.Choices[0].Message.Content

	content = strings.ReplaceAll(
		content,
		"```json",
		"",
	)

	content = strings.ReplaceAll(
		content,
		"```",
		"",
	)

	content = strings.TrimSpace(content)

	fmt.Println("========== GROQ RESPONSE ==========")
	fmt.Println(content)

	var metadata models.QuizMetadata

	err = json.Unmarshal(
		[]byte(content),
		&metadata,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to parse metadata response: %v\nresponse=%s",
			err,
			content,
		)
	}

	return &metadata, nil
}

// question generation from topic or text

func GenerateContent(prompt string) (string, error) {
	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("GROQ_API_KEY not found")
	}

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
		"https://api.groq.com/openai/v1/chat/completions",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return "", err
	}

	req.Header.Set(
		"Authorization",
		"Bearer "+apiKey,
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

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

	var groqResp GroqResponse

	if err := json.Unmarshal(body, &groqResp); err != nil {
		return "", err
	}

	if len(groqResp.Choices) == 0 {
		return "", fmt.Errorf("empty response")
	}

	content := groqResp.Choices[0].Message.Content

	content = strings.ReplaceAll(content, "```json", "")
	content = strings.ReplaceAll(content, "```", "")
	content = strings.TrimSpace(content)

	return content, nil
}
