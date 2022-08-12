package rating

import "math"

// average is the representation of the evaluation.
type average struct {

	// overallRating is the basis considered for the rating.
	// for example if your base was 5 stars, the overall rating should be 5
	overallRating float64

	// ratings represents the total value of ratings of an overallRating.
	// for example, for example out of 100 reviews are 5 stars so you have 100 reviews from overallRating.
	ratings float64
}

// calcAVG calcular a media de avaliações
func calcAVG(avgs ...average) float64 {

	var avg, votes float64

	for _, _avg := range avgs {
		avg += _avg.overallRating * _avg.ratings
		votes += _avg.ratings
	}

	result := math.Floor((avg/votes)*100) / 100
	return result
}
