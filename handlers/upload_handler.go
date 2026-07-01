package handlers

import (
	"fmt"
	"net/http"
	"os"

	"pdfapi/config"
	"pdfapi/services"

	"github.com/gin-gonic/gin"
)

func UploadFileHandler(c *gin.Context) {

	// Get uploaded PDF
	_, fileHeader, err := c.Request.FormFile("file")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	examName := c.PostForm("examName")

	// Save PDF
	pdfPath := "uploads/" + fileHeader.Filename
	defer os.Remove(pdfPath)

	if err := c.SaveUploadedFile(fileHeader, pdfPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Split PDF into 3-page chunks
	// pdfParts, err := services.SplitPDF(pdfPath)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	fmt.Println("🎶🎶🎶🎶 Generated:")
	// OCR all chunks and merge text
	// ocrText, err := services.ExtractTextFromPDFChunks(pdfParts)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// ocrText, err := services.ExtractText(pdfPath)

	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"error": err.Error(),
	// 	})
	// 	return
	// }
	fmt.Println("🎶🎶 Generated:")
	// Create chunks
	// chunks := services.CreateChunks(ocrText)

	// fmt.Println("Total Chunks:", len(chunks))

	// provider := &embeddings.BGEProvider{
	// 	BaseURL: "http://127.0.0.1:8000",
	// }

	// embedService := embeddings.NewService(provider)

	// chunkEmbeddings, err := embeddings.GenerateChunkEmbeddings(
	// 	fileHeader.Filename,
	// 	chunks,
	// 	embedService,
	// )

	// if err != nil {
	// 	fmt.Println("warning: embeddings failed:", err)
	// 	chunkEmbeddings = []embeddings.ChunkEmbedding{}
	// }

	// fmt.Println("Embeddings Generated:", len(chunkEmbeddings))

	// savedCount, err := embeddings.StoreUniqueEmbeddings(
	// 	"embeddings.json",
	// 	chunkEmbeddings,
	// 	0.90,
	// )

	// if err != nil {
	// 	fmt.Println("embedding storage error:", err)
	// }

	// fmt.Println("New Embeddings Saved:", savedCount)

	// Generate questions using Groq
	// questions, err := services.GenerateQuestions(ocrText)

	// allQuestion, err := services.ParseMCQs(ocrText)

	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"error": err.Error(),
	// 	})
	// 	return
	// }
	// Step 1
	ocrText, err := services.ExtractText(pdfPath)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	allQuestions, err := services.ParseMCQs(ocrText)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Modes:
	// all
	// random
	//chunk

	mode := c.DefaultPostForm(
		"mode",
		services.ModeAll,
	)

	count := 20

	if countStr := c.PostForm("count"); countStr != "" {
		fmt.Sscanf(countStr, "%d", &count)
	}

	generatedQuestions := 0
	totalChunksCreated := 0

	switch mode {

	case services.ModeAll:
		generatedQuestions = len(allQuestions)

		metadata, err := services.GenerateQuizMetadata(
			allQuestions,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		err = services.SaveQuiz(
			config.FirestoreClient,
			fileHeader.Filename,
			allQuestions,
			metadata,
			examName,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		services.SaveQuizVariants(
			fileHeader.Filename,
			examName,
			allQuestions,
		)

	case services.ModeRandom:

		randomQuestions := services.SelectQuestions(
			allQuestions,
			services.ModeRandom,
			count,
		)
		generatedQuestions = len(randomQuestions)

		metadata, err := services.GenerateQuizMetadata(
			randomQuestions,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		err = services.SaveQuiz(
			config.FirestoreClient,
			fileHeader.Filename,
			randomQuestions,
			metadata,
			examName,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		services.SaveQuizVariants(
			fileHeader.Filename,
			examName,
			allQuestions,
		)

	case services.ModeChunk:

		chunks := services.GenerateQuestionChunks(
			allQuestions,
			count,
			0,
			false,
		)

		totalChunksCreated = len(chunks)
		generatedQuestions = len(allQuestions)

		for i, chunk := range chunks {

			metadata, err := services.GenerateQuizMetadata(
				chunk,
			)

			if err != nil {
				continue
			}

			fileName := fmt.Sprintf(
				"%s_part_%d",
				fileHeader.Filename,
				i+1,
			)

			err = services.SaveQuiz(
				config.FirestoreClient,
				fileName,
				chunk,
				metadata,
				examName,
			)

			if err != nil {
				fmt.Println(err)
			}

		}
		services.SaveQuizVariants(
			fileHeader.Filename,
			examName,
			allQuestions,
		)

	default:

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid mode",
		})
		return
	}

	title, body := services.GetRandomQuizNotification(examName)

	// services.SendTopicNotificationAsync(
	// 	"exam_news",
	// 	title,
	// 	body,
	// )
	services.SendTokenNotificationAsync(
		"ds0gUTOoQUqlyX7ggyosQ5:APA91bGnXumna765ZOc7WVr2Qbl3Ol0GBYgd_gSuyh8nj9hqJph3w8QCrw_CysvZ3gxWr_o-JVlsGvhFDe9ZCzTVHjS3O4Zcwg1ecThby07SDB4YOgzYzjM",
		title,
		body,
	)

	// err = services.SaveQuiz(
	// 	config.FirestoreClient,
	// 	fileHeader.Filename,
	// 	questions,
	// )

	// if err != nil {

	// 	fmt.Println("SAVE ERROR:", err)

	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"error": err.Error(),
	// 	})

	// 	return
	// }

	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"error": err.Error(),
	// 	})
	// 	return
	// }

	// Convert JSON string to JSON object

	// Final Response

	c.JSON(http.StatusOK, gin.H{
		"message":                 "Quiz generated successfully 🚀",
		"filename":                fileHeader.Filename,
		"mode":                    mode,
		"totalExtractedQuestions": len(allQuestions),
		"generatedQuestions":      generatedQuestions,
		"totalChunksCreated":      totalChunksCreated,
		"count":                   count,
	})
}
