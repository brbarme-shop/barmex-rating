package rating

import "math"

type RatingAvg struct {
	OverallRating float64
	TotalRating   float64
}

func Average(avgs ...RatingAvg) float64 {

	var r, v float64
	for _, avg := range avgs {
		r += avg.OverallRating * avg.TotalRating
		v += avg.TotalRating
	}

	return math.Floor((r/v)*100) / 100
}
