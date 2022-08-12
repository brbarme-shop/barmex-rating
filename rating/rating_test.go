package rating

import (
	"testing"
)

var tableTest = []struct {
	ratingAvg []average
	avgExpect float64
}{
	{
		ratingAvg: []average{
			{
				overallRating: 5,
				ratings:       252,
			},
			{
				overallRating: 4,
				ratings:       124,
			},
			{
				overallRating: 3,
				ratings:       40,
			},
			{
				overallRating: 2,
				ratings:       29,
			},
			{
				overallRating: 1,
				ratings:       33,
			},
		},
		avgExpect: 4.11,
	},
}

func TestAverage(t *testing.T) {

	for _, tb := range tableTest {
		got := calcAVG(tb.ratingAvg...)
		if got != tb.avgExpect {
			t.Errorf("Average() returned %f but was expected %f", got, tb.avgExpect)
		}
	}

}

func BenchmarkAverage(b *testing.B) {

	for _, tb := range tableTest {
		for i := 0; i < b.N; i++ {
			got := calcAVG(tb.ratingAvg...)
			if got != tb.avgExpect {
				b.Errorf("Average() returned %f but was expected %f", got, tb.avgExpect)
			}
		}
	}

}
