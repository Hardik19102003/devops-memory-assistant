package ai

func CalculateConfidence(distance float64) string {

	if distance < 10 {
		return "HIGH"
	}

	if distance < 25 {
		return "MEDIUM"
	}

	return "LOW"
}