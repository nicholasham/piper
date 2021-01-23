package piper

func removeDuplicates(stages []Stage) []Stage {
	check := make(map[Stage]int)
	var result []Stage
	for _, stage := range stages {
		check[stage] = 1
	}
	for stage, _ := range check {
		result = append(result, stage)
	}
	return result
}
