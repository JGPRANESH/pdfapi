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

	ocrText, err := services.ExtractText(pdfPath)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
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

	questions, err := services.ParseMCQs(ocrText)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Generate metadata from extracted questions
	metadata, err := services.GenerateQuizMetadata(questions)

	if err != nil {
		fmt.Println("METADATA ERROR:", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	fmt.Printf("Generated Metadata: %+v\n", metadata)

	// Save quiz
	err = services.SaveQuiz(
		config.FirestoreClient,
		fileHeader.Filename,
		questions,
		metadata,
	)

	if err != nil {
		fmt.Println("SAVE ERROR:", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	fmt.Println("✅ Quiz saved successfully")
	fmt.Println("Questions count:", len(questions))

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
		"message":    "Questions generated successfully 🚀",
		"filename":   fileHeader.Filename,
		"ocr_length": len(ocrText),
		"questions":  questions,
	})
}
