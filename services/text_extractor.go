package services

import (
	"strings"

	"github.com/ledongthuc/pdf"
)

func ExtractNativeText(pdfPath string) (string, error) {

	f, r, err := pdf.Open(pdfPath)

	if err != nil {
		return "", err
	}

	defer f.Close()

	var builder strings.Builder

	totalPage := r.NumPage()

	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {

		page := r.Page(pageIndex)

		if page.V.IsNull() {
			continue
		}

		text, err := page.GetPlainText(nil)

		if err != nil {
			continue
		}

		builder.WriteString(text)
		builder.WriteString("\n")
	}

	return builder.String(), nil
}
