package handlers

import (
	"fmt"
	"net/http"

	"pdfapi/models"
	"pdfapi/services"

	"github.com/gin-gonic/gin"
)

func NotificationHandler(c *gin.Context) {

	var req models.NotificationRequest

	// Parse JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	fmt.Println("========== Notification Request ==========")
	fmt.Println("Token:", req.Token)
	fmt.Println("Topic:", req.Topic)

	// Validate required fields
	// if req.Topic == "" || req.Title == "" || req.Body == "" {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"error": "topic, title and body are required",
	// 	})
	// 	return
	// }

	if req.Token == "" || req.Title == "" || req.Body == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "token, title and body are required",
		})
		return
	}

	// Send notification
	if err := services.SendTokenNotification(
		req.Token,
		req.Title,
		req.Body,
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Notification sent successfully 🚀",
	})

}
