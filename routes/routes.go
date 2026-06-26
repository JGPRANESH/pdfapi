package routes

import (
	"net/http"
	"pdfapi/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Server is running 🚀",
		})
	})

	// PDF MCQ Extraction
	router.POST("/upload", handlers.UploadFileHandler)

	// AI MCQ Generation
	router.POST("/generate", handlers.GenerateQuizHandler)
	// notification
	router.POST("/notification", handlers.NotificationHandler)
}
