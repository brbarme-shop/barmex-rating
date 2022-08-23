package rating

import (
	"context"
	"errors"
	"fmt"
	"testing"
)

var ctx = context.TODO()

func TestPutRating(t *testing.T) {

	t.Run("Must calculate the average when the item exists but the star is not", func(t *testing.T) {

		data := &RatingAverage{
			ItemId:  "658029a7-33da-4997-aeec-5e37947a1d1f",
			Average: 0,
			Ratings: []Rating{
				{
					Star:  5,
					Count: 0,
				},
			},
		}

		db := &repositoryMock{
			readByItemId: func() (*RatingAverage, error) {
				return data, nil
			},
			updateRating: func() error {
				return nil
			},
		}

		inputValid := &PutRatingInput{
			ItemId: "658029a7-33da-4997-aeec-5e37947a1d1f",
			Star:   1,
		}

		err := PutRating(ctx, inputValid, db)
		if err != nil {
			t.Fatalf("expected: NIL, got: %v", err)
		}

		if data.Average == 0 {
			t.Fatal("the mean was not calculated")
		}

	})

	t.Run("Must calculate average when item and star exist", func(t *testing.T) {

		data := &RatingAverage{
			ItemId:  "658029a7-33da-4997-aeec-5e37947a1d1f",
			Average: 0,
			Ratings: []Rating{
				{
					Star:  5,
					Count: 0,
				},
			},
		}

		db := &repositoryMock{
			readByItemId: func() (*RatingAverage, error) {
				return data, nil
			},
			updateRating: func() error {
				return nil
			},
		}

		inputValid := &PutRatingInput{
			ItemId: "658029a7-33da-4997-aeec-5e37947a1d1f",
			Star:   5,
		}

		err := PutRating(ctx, inputValid, db)
		if err != nil {
			t.Fatalf("expected: NIL, got: %v", err)
		}

		if data.Average == 0 {
			t.Fatal("the mean was not calculated")
		}

	})

	t.Run("Should abort process when `PutNewRating` failed and returns an error", func(t *testing.T) {

		db := &repositoryMock{
			readByItemId: func() (*RatingAverage, error) {
				return nil, ErrRatingNotFound
			},
			putNewRating: func() error {
				return fmt.Errorf("failed to put-new-rating in database")
			},
		}

		inputValid := &PutRatingInput{
			ItemId: "658029a7-33da-4997-aeec-5e37947a1d1f",
			Star:   5,
		}

		err := PutRating(ctx, inputValid, db)
		if err == nil {
			t.Fatalf("expected: failed to put-new-rating in database, got: %v", err)
		}

	})

	t.Run("Should abort process when `ReadByItemId` returns an error other than ErrRatingNotFound", func(t *testing.T) {

		db := &repositoryMock{
			readByItemId: func() (*RatingAverage, error) {
				return nil, fmt.Errorf("any error other than ErrRatingNotFound")
			},
		}

		inputValid := &PutRatingInput{
			ItemId: "658029a7-33da-4997-aeec-5e37947a1d1f",
			Star:   5,
		}

		err := PutRating(ctx, inputValid, db)

		if err == nil {
			t.Fatalf("expected: %v, got: NIL", ErrRatingNotFound)
		}

		if err != nil {
			if errors.Is(err, ErrPutRatingInputInvalid) {
				t.Fatalf("expected: any error other than ErrRatingNotFound, got: %v", err)
			}
		}
	})

	t.Run("Should abort process when input is invalid and return error ErrPutRatingInputInvalid", func(t *testing.T) {

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
