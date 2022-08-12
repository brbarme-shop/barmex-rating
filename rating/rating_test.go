package rating

import (
	"context"
	"testing"
)

var tableTest = []struct {
	ratingAvg []Average
	avgExpect float64
}{
	{
		ratingAvg: []Average{
			{
				OverallRating: 5,
				Ratings:       252,
			},
			{
				OverallRating: 4,
				Ratings:       124,
			},
			{
				OverallRating: 3,
				Ratings:       40,
			},
			{
				OverallRating: 2,
				Ratings:       29,
			},
			{
				OverallRating: 1,
				Ratings:       33,
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

func TestPutRating(t *testing.T) {

	db := &repositoryMock{
		GetItemByProductIdMock: func() (*RatingItem, error) {
			return &RatingItem{
				RatingId: "fake-hash-id",
				Item:     "fake-hash-product-item-id",
				Avg:      4.11,
				Averages: []Average{
					{
						OverallRating: 5,
						Ratings:       252,
					},
					{
						OverallRating: 4,
						Ratings:       124,
					},
					{
						OverallRating: 3,
						Ratings:       40,
					},
					{
						OverallRating: 2,
						Ratings:       29,
					},
					{
						OverallRating: 1,
						Ratings:       33,
					},
				},
			}, nil
		},
		SaveRatingItemMock: func() error {
			return nil
		},
	}

	err := PutRating(context.Background(), &RatingInput{
		ProductId:     "fake-hash-product-item-id",
		OverallRating: 5,
	}, db)

	if err != nil {
		t.Fail()
	}
}

type repositoryMock struct {
	GetItemByProductIdMock func() (*RatingItem, error)
	SaveRatingItemMock     func() error
}

func (r *repositoryMock) GetItemByProductId(ctx context.Context, productId string) (*RatingItem, error) {
	return r.GetItemByProductIdMock()
}

func (r *repositoryMock) SaveRatingItem(ctx context.Context, ratingProduct *RatingItem) error {
	return r.SaveRatingItemMock()
}
