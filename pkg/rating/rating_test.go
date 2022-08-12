package rating

import "testing"

func TestAverage(t *testing.T) {

	for _, tb := range tableTest {
		got := Average(tb.ratingAvg...)
		if got != tb.avgExpect {
			t.Errorf("Average() returned %f but was expected %f", got, tb.avgExpect)
		}
	}

}

func BenchmarkAverage(b *testing.B) {

	for _, tb := range tableTest {
		for i := 0; i < b.N; i++ {
			got := Average(tb.ratingAvg...)
			if got != tb.avgExpect {
				b.Errorf("Average() returned %f but was expected %f", got, tb.avgExpect)
			}
		}
	}

}

var tableTest = []struct {
	ratingAvg []RatingAvg
	avgExpect float64
}{
	{
		ratingAvg: []RatingAvg{
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
