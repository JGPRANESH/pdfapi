package services

import (
	"log"
	"strings"
)

func ExtractText(pdfPath string) (string, error) {

	log.Println("Trying native extraction...")

	text, err := ExtractNativeText(pdfPath)

	if err == nil {

		cleanText := strings.TrimSpace(text)

		log.Printf("Native text length: %d\n", len(cleanText))

		if len(cleanText) > 500 {

			log.Println("✅ Using native extraction")

			return cleanText, nil
		}
	}

	log.Println("⚠️ Falling back to OCR...")

	pdfParts, err := SplitPDF(pdfPath)

	if err != nil {
		return "", err
	}

	return ExtractTextFromPDFChunks(pdfParts)
}
