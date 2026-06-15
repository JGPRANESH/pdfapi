package services

import (
	"os"
	"path/filepath"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

func SplitPDF(pdfPath string) ([]string, error) {

	outputDir := "tmp"

	err := os.MkdirAll(outputDir, os.ModePerm)

	if err != nil {
		return nil, err
	}

	err = api.SplitFile(
		pdfPath,
		outputDir,
		3, // split every 3 pages
		nil,
	)

	if err != nil {
		return nil, err
	}

	files, err := filepath.Glob(outputDir + "/*.pdf")

	if err != nil {
		return nil, err
	}

	return files, nil
}
