package services

var TimePerQuestion = map[string]map[string]int{
	"Maths": {
		"Easy":   90,
		"Medium": 120,
		"Hard":   150,
	},
	"Reasoning": {
		"Easy":   60,
		"Medium": 75,
		"Hard":   90,
	},
	"English": {
		"Easy":   45,
		"Medium": 60,
		"Hard":   75,
	},
	"GK": {
		"Easy":   45,
		"Medium": 60,
		"Hard":   75,
	},
	"CA": {
		"Easy":   45,
		"Medium": 60,
		"Hard":   75,
	},
	"Science": {
		"Easy":   60,
		"Medium": 75,
		"Hard":   90,
	},
}

func CalculateDuration(category, difficulty string, totalQuestions int) int {

	secondsPerQuestion := 45 // default

	if categoryMap, ok := TimePerQuestion[category]; ok {
		if seconds, ok := categoryMap[difficulty]; ok {
			secondsPerQuestion = seconds
		}
	}

	totalSeconds := totalQuestions * secondsPerQuestion

	// Round up to the next minute
	return (totalSeconds + 59) / 60
}
