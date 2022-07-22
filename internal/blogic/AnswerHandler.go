package blogic

func orderly(userAnswer string, answers []string, maxPoint int) int {
	for _, val := range answers {
		if val == userAnswer {
			return maxPoint
		}
	}
	return 0
}
