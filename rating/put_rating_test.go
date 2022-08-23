package rating

import (
	"context"
	"errors"
	"fmt"
	"testing"
)

func TestPutRating(t *testing.T) {

	t.Run("Should abort process when input is invalid and return error ErrPutRatingInputInvalid", func(t *testing.T) {

		ctx := context.TODO()

		var testTable = []struct {
			inputinvalid *PutRatingInput
		}{
			{
				inputinvalid: &PutRatingInput{},
			},
			{
				inputinvalid: &PutRatingInput{ItemId: "", Star: 0},
			},
			{
				inputinvalid: &PutRatingInput{ItemId: "", Star: -1},
			},
			{
				inputinvalid: &PutRatingInput{ItemId: "", Star: 6},
			},
			{
				inputinvalid: &PutRatingInput{ItemId: "", Star: 1},
			},
			{
				inputinvalid: &PutRatingInput{ItemId: "658029a7-33da-4997-aeec-5e37947a1d1f", Star: 0},
			},
			{
				inputinvalid: &PutRatingInput{ItemId: "658029a7-33da-4997-aeec-5e37947a1d1f", Star: -1},
			},
			{
				inputinvalid: &PutRatingInput{ItemId: "658029a7-33da-4997-aeec-5e37947a1d1f", Star: 6},
			},
			{
				inputinvalid: nil,
			},
		}

		for _, tb := range testTable {

			err := PutRating(ctx, tb.inputinvalid, &repositoryMock{readByItemId: func() (*RatingAverage, error) { return nil, fmt.Errorf("failed") }})
			if err == nil {
				t.Fatalf("expected: %v, got: NIL", ErrPutRatingInputInvalid)
			}

			if err != nil {
				if !errors.Is(err, ErrPutRatingInputInvalid) {
					t.Fatalf("expected: %v, got: %v", ErrPutRatingInputInvalid, err)
				}
			}

		}

	})

}

type repositoryMock struct {
	putNewRating func() error
	readByItemId func() (*RatingAverage, error)
	updateRating func() error
}

func (r *repositoryMock) PutNewRating(ctx context.Context, itemId string, star int64) error {
	return r.putNewRating()
}

func (r *repositoryMock) ReadByItemId(ctx context.Context, itemId string) (*RatingAverage, error) {
	return r.readByItemId()
}

func (r *repositoryMock) UpdateRating(ctx context.Context, itemId string, average float64, star, count int64) error {
	return r.updateRating()
}
