package services

import (
	"regexp"
	"strings"
)

var (
	pageFooterRegex = regexp.MustCompile(`^\d+\s*\|\s*.*$`)
	pageHeaderRegex = regexp.MustCompile(`^\d+\s*[-–]\s*\d+$`)
)

func CleanExtractedText(text string) string {

	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")

	lines := strings.Split(text, "\n")

	var cleaned []string

	for _, line := range lines {

		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		if pageFooterRegex.MatchString(line) {
			continue
		}

		if pageHeaderRegex.MatchString(line) {
			continue
		}

		if strings.Contains(line, "Rapid Revision") {
			continue
		}

		if strings.Contains(line, "CA -") {
			continue
		}

		if strings.Contains(line, "Page ") {
			continue
		}

		cleaned = append(cleaned, line)
	}

	return strings.Join(cleaned, "\n")
}
