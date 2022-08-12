package rating_test

import (
	"testing"

	"github.com/brbarme-shop/brbarmex-rating/rating"
)

var tableTest = []struct {
	ratingAvg []rating.RatingAvg
	avgExpect float64
}{
	{
		ratingAvg: []rating.RatingAvg{
			{
				OverallRating: 5,
				TotalRating:   252,
			},
			{
				OverallRating: 4,
				TotalRating:   124,
			},
			{
				OverallRating: 3,
				TotalRating:   40,
			},
			{
				OverallRating: 2,
				TotalRating:   29,
			},
			{
				OverallRating: 1,
				TotalRating:   33,
			},
		},
		avgExpect: 4.11,
	},
}

func TestAverage(t *testing.T) {

	for _, tb := range tableTest {
		got := rating.Average(tb.ratingAvg...)
		if got != tb.avgExpect {
			t.Errorf("Average() returned %f but was expected %f", got, tb.avgExpect)
		}
	}

}

func BenchmarkAverage(b *testing.B) {

	for _, tb := range tableTest {
		for i := 0; i < b.N; i++ {
			got := rating.Average(tb.ratingAvg...)
			if got != tb.avgExpect {
				b.Errorf("Average() returned %f but was expected %f", got, tb.avgExpect)
			}
		}
	}

}
