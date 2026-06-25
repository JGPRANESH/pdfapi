package handlers

import (
	"io"
	"net/http"
	"strings"

	"pdfapi/models"
	"pdfapi/services/generation"

	"github.com/gin-gonic/gin"
)

func GenerateQuizHandler(c *gin.Context) {

	var req models.GenerateQuizRequest

	// Try JSON first
	if err := c.ShouldBindJSON(&req); err != nil {

		// If JSON fails, try plain text
		body, readErr := io.ReadAll(c.Request.Body)
		if readErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": readErr.Error(),
			})
			return
		}

		req.Topic = strings.TrimSpace(string(body))
	}

	// Validation
	if req.Topic == "" && req.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "topic or content is required",
		})
		return
	}

	// Defaults
	if req.Count <= 0 {
		req.Count = 10
	}

	if req.Difficulty == "" {
		req.Difficulty = "medium"
	}

	questions, err := generation.GenerateQuiz(req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"questions": questions,
	})
}
