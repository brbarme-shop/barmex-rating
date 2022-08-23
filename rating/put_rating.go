package rating

import (
	"context"
	"errors"
	"math"
)

var (
	ErrPutRatingInputInvalid    = errors.New("the RatingAverageInput is invalid. The struct and your fields cann't be NIL. Check all fields or see more about in documentation")
	ErrRatingNotFound           = errors.New("the RatingAverageInput not exists")
	ErrAverageRepositoryInvalid = errors.New("the AverageRepository can't be NIL")
	ErrStarIdNotFound           = errors.New("the StarId not containt star")
	ErrFailedToPutNewRating     = errors.New("a")
	ErrStartNotIdentifier       = errors.New("")
)

type PutRatingRepository interface {
	PutNewRating(ctx context.Context, itemId string, star int64) error
	ReadByItemId(ctx context.Context, itemId string) (*RatingAverage, error)
	UpdateRating(ctx context.Context, itemId string, average float64, star, count int64) error
}

type PutRatingInput struct {
	ItemId string `json:"item_id"`
	Star   int64  `json:"star"`
}

type Rating struct {
	Star  int64
	Count int64
}

type RatingAverage struct {
	ItemId  string
	Average float64
	Ratings []Rating
}

func (r *RatingAverage) calcAVG() {

	var avg, avgCounts float64
	for i := range r.Ratings {
		avg += float64(r.Ratings[i].Star) * float64(r.Ratings[i].Count)
		avgCounts += float64(r.Ratings[i].Count)
	}

	r.Average = math.Floor((avg/avgCounts)*100) / 100
}

func PutRating(ctx context.Context, ratingInput *PutRatingInput, db PutRatingRepository) error {

	isInvalid := ratingInput == nil || len(ratingInput.ItemId) == 0 || ratingInput.Star <= 0 || ratingInput.Star > 5
	if isInvalid {
		return ErrPutRatingInputInvalid
	}

	ratingData, err := db.ReadByItemId(ctx, ratingInput.ItemId)
	if err != nil {
		if !errors.Is(err, ErrRatingNotFound) {
			return err
		}
	}

	if ratingData == nil {
		err = db.PutNewRating(ctx, ratingInput.ItemId, ratingInput.Star)
		return err
	}

	for i := range ratingData.Ratings {
		if ratingData.Ratings[i].Star == ratingInput.Star {
			ratingData.Ratings[i].Count += 1
			ratingData.calcAVG()
			return db.UpdateRating(ctx, ratingData.ItemId, ratingData.Average, ratingData.Ratings[i].Star, ratingData.Ratings[i].Count)
		}
	}

	ratingData.Ratings = append(ratingData.Ratings, Rating{Star: ratingInput.Star, Count: 1})
	ratingData.calcAVG()

	err = db.UpdateRating(ctx, ratingData.ItemId, ratingData.Average, ratingInput.Star, 1)
	return err
}
