package routes

import (
	"pdfapi/handlers"
	"pdfapi/services"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {

	r.POST("/upload", handlers.UploadFileHandler)
	r.POST("/generate", handlers.GenerateQuizHandler)

	r.GET("/test-notification", func(c *gin.Context) {

		err := services.SendTopicNotification(
			"live_test",
			"🧠 New Challenge Dropped",
			"New Mock Test Available 🎯",
		)

		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"message": "notification sent",
		})
	})
}
