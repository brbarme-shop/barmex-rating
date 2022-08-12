package rating

import "math"

type RatingAvg struct {
	OverallRating float64
	TotalRating   float64
}

func Average(ratings ...RatingAvg) float64 {

	var totalAVG, votesAVG float64
	for _, rating := range ratings {
		totalAVG += rating.OverallRating * rating.TotalRating
		votesAVG += rating.TotalRating
	}

	return math.Floor((totalAVG/votesAVG)*100) / 100
}
