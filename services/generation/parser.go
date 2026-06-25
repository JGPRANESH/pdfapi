package generation

import "strings"

func CleanResponse(response string) string {

	response = strings.TrimSpace(response)

	response = strings.TrimPrefix(
		response,
		"```json",
	)

	response = strings.TrimPrefix(
		response,
		"```",
	)

	response = strings.TrimSuffix(
		response,
		"```",
	)

	return strings.TrimSpace(response)
}
