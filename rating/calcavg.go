package rating

import "math"

func calcAVG(avgs ...AverageScore) float64 {

	var avgScore, avgScorePoints float64
	for i := range avgs {
		avgScore += float64(avgs[i].Star) * float64(avgs[i].ScorePoint)
		avgScorePoints += float64(avgs[i].ScorePoint)
	}

	result := math.Floor((avgScore/avgScorePoints)*100) / 100
	return result
}
