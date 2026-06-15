package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type OCRResponse struct {
	ParsedResults []struct {
		ParsedText string `json:"ParsedText"`
	} `json:"ParsedResults"`
}

// OCR uploaded file
func ExtractTextFromImage(
	file multipart.File,
	fileHeader *multipart.FileHeader,
) (string, error) {

	body := &bytes.Buffer{}

	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile(
		"file",
		fileHeader.Filename,
	)

	if err != nil {
		return "", err
	}

	_, err = io.Copy(part, file)

	if err != nil {
		return "", err
	}

	// OCR.Space API Key
	writer.WriteField("apikey", os.Getenv("OCR_SPACE_API_KEY"))

	writer.Close()

	req, err := http.NewRequest(
		"POST",
		"https://api.ocr.space/parse/image",
		body,
	)

	if err != nil {
		return "", err
	}

	req.Header.Set(
		"Content-Type",
		writer.FormDataContentType(),
	)

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	var result OCRResponse

	err = json.NewDecoder(resp.Body).Decode(&result)

	if err != nil {
		return "", err
	}

	if len(result.ParsedResults) == 0 {
		return "", nil
	}

	var fullText string

	for _, page := range result.ParsedResults {
		fullText += page.ParsedText + "\n"
	}

	return fullText, nil
}

// OCR all split PDFs
func ExtractTextFromPDFChunks(
	pdfParts []string,
) (string, error) {

	var fullText string

	for _, pdfPath := range pdfParts {

		text, err := ExtractTextFromPDFFile(pdfPath)

		if err != nil {
			continue
		}
		fmt.Println("Processing:", pdfPath)
		fmt.Println("Text Length:", len(text))

		fullText += text + "\n"
	}

	return fullText, nil
}
func ExtractTextFromPDFFile(
	pdfPath string,
) (string, error) {

	file, err := os.Open(pdfPath)

	if err != nil {
		return "", err
	}

	defer file.Close()

	body := &bytes.Buffer{}

	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile(
		"file",
		filepath.Base(pdfPath),
	)

	if err != nil {
		return "", err
	}

	_, err = io.Copy(part, file)

	if err != nil {
		return "", err
	}

	writer.WriteField("apikey", os.Getenv("OCR_SPACE_API_KEY"))

	writer.Close()

	req, err := http.NewRequest(
		"POST",
		"https://api.ocr.space/parse/image",
		body,
	)

	if err != nil {
		return "", err
	}

	req.Header.Set(
		"Content-Type",
		writer.FormDataContentType(),
	)

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	var result OCRResponse

	err = json.NewDecoder(resp.Body).Decode(&result)

	if err != nil {
		return "", err
	}

	var fullText string

	for _, page := range result.ParsedResults {
		fullText += page.ParsedText + "\n"
	}

	return fullText, nil
}
