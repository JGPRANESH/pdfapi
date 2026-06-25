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
Analyze the following quiz questions.

Generate:

1. Title
2. small detail Description 
3. Category
4. Difficulty (Easy, Medium, Hard)

Return ONLY valid JSON:

{
  "title":"",
  "description":"",
  "category":"",
  "difficulty":""
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
	fmt.Println("========== RAW GROQ RESPONSE ==========")
	fmt.Println(string(body))
	var groqResp GroqResponse

	if err := json.Unmarshal(body, &groqResp); err != nil {
		return "", err
	}

	if len(groqResp.Choices) == 0 {
		return "", fmt.Errorf("empty response")
	}

	content := groqResp.Choices[0].Message.Content
	fmt.Println("========== CONTENT ==========")
	fmt.Println(content)

	content = strings.ReplaceAll(content, "```json", "")
	content = strings.ReplaceAll(content, "```", "")
	content = strings.TrimSpace(content)

	return content, nil
}
func BuildMCQPrompt(text string) string {

	return fmt.Sprintf(`
You are an expert MCQ extraction engine.

Extract ALL MCQs from the text.

Return ONLY valid JSON.

Schema:

[
  {
    "id":"1",
    "questionText":"",
    "options":["","","",""],
    "correctOptionIndex":0,
    "difficulty":"Medium",
    "description":""
  }
]

Rules:

- Return ALL questions.
- Never skip questions.
- Never merge questions.
- Exactly four options.
- Ignore page headers.
- Ignore page footers.
- Ignore advertisements.
- Ignore page numbers.
- Ignore explanations.
- Do not generate markdown.
- Return ONLY JSON.

TEXT:

%s
`, text)

}

func GenerateQuestionsFromText(text string) ([]models.Question, error) {

	cleanText := CleanExtractedText(text)

	chunks := CreateChunks(cleanText)

	var allQuestions []models.Question

	for i, chunk := range chunks {

		fmt.Printf("\n========== Processing Chunk %d/%d ==========\n", i+1, len(chunks))

		prompt := BuildMCQPrompt(chunk)

		var response string
		var err error

		// Retry this chunk up to 3 times
		for retry := 1; retry <= 3; retry++ {

			response, err = GenerateContent(prompt)

			if err == nil {
				break
			}

			fmt.Printf(
				"Chunk %d failed (Attempt %d/3): %v\n",
				i+1,
				retry,
				err,
			)

			time.Sleep(3 * time.Second)
		}

		if err != nil {
			fmt.Printf("Skipping chunk %d after retries\n", i+1)
			continue
		}

		response = ExtractJSONArray(response)

		var questions []models.Question

		if err := json.Unmarshal([]byte(response), &questions); err != nil {

			fmt.Printf(
				"Chunk %d returned invalid JSON: %v\n",
				i+1,
				err,
			)

			continue
		}

		for _, q := range questions {

			if q.Difficulty == "" {
				q.Difficulty = "Medium"
			}

			if err := ValidateQuestion(q); err != nil {

				fmt.Println("Skipping invalid question:", err)

				continue
			}

			allQuestions = append(allQuestions, q)
		}

		fmt.Printf(
			"Chunk %d completed (%d questions)\n",
			i+1,
			len(questions),
		)

		// Wait before the next request to avoid rate limits
		time.Sleep(2500 * time.Millisecond)
	}

	if len(allQuestions) == 0 {
		return nil, fmt.Errorf("no questions extracted")
	}

	fmt.Printf(
		"\nTotal Questions Extracted: %d\n",
		len(allQuestions),
	)

	return allQuestions, nil
}
func ExtractJSONArray(text string) string {

	text = strings.TrimSpace(text)

	// Remove markdown fences
	text = strings.ReplaceAll(text, "```json", "")
	text = strings.ReplaceAll(text, "```", "")

	start := strings.Index(text, "[")
	end := strings.LastIndex(text, "]")

	if start == -1 || end == -1 || end < start {
		return text
	}

	return strings.TrimSpace(text[start : end+1])
}
