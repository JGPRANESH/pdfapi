package generation

import (
	"fmt"
	"pdfapi/models"
)

func BuildPrompt(req models.GenerateQuizRequest) string {

	source := req.Topic

	if req.Content != "" {
		source = req.Content
	}

	return fmt.Sprintf(`
Generate %d multiple-choice questions.

Difficulty: %s

Source:
%s

Rules:
1. Return valid JSON only.
2. 4 options per question.
3. Exactly one correct answer.
4. No markdown.
5. No explanations.

Format:
Rules:
- No markdown
- No code fences
- No java blocks
- Return code snippets as plain text
- Keep questionText readable
- Return valid JSON only

[
 {
   "question":"...",
   "options":["A","B","C","D"],
   "answer":"..."
   "description":"Short explanation of why the answer is correct."
 }
   For every question:
- Provide a short explanation.
- Description should explain why the correct answer is correct.
- Maximum 2-3 sentences.
]
`,
		req.Count,
		req.Difficulty,
		source,
	)
}
