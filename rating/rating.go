package rating

import (
	"context"
	"errors"
	"math"
)

var (
	ErrRatingAverageInputInvalid = errors.New("the RatingAverageInput is invalid. The struct and your fields cann't be NIL. Check all fields or see more about in documentation")
	ErrAverageNotExist           = errors.New("the RatingAverageInput not exists")
	ErrAverageRepositoryInvalid  = errors.New("the AverageRepository can't be NIL")
	ErrStarIdNotFound            = errors.New("the StarId not containt star")
)

type RatingAverageRepository interface {
	CreateRatingAverage(ctx context.Context, average *RatingAverage) error
	ReadRatingAverages(ctx context.Context, itemId string) (*RatingAverage, error)
	UpdateRatingAverage(ctx context.Context, average *RatingAverage) error
}

type RatingInput struct {
	ItemId string `json:"item_id"`
	Star   int64  `json:"star"`
}

type Rating struct {
	Star  int64
	Count int64
}

type RatingAverage struct {
	Id      string
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

func PutRating(ctx context.Context, ratingInput *RatingInput, db RatingAverageRepository) error {

	if ratingInput == nil || len(ratingInput.ItemId) == 0 || ratingInput.Star == 0 {
		return ErrRatingAverageInputInvalid
	}

	_rating, err := db.ReadRatingAverages(ctx, ratingInput.ItemId)
	if err != nil {
		if !errors.Is(err, ErrAverageNotExist) {
			return err
		}
	}

	if _rating == nil {
		_rating = &RatingAverage{}
		_rating.ItemId = ratingInput.ItemId
		_rating.Ratings = append(_rating.Ratings, Rating{Star: ratingInput.Star, Count: 1})
		_rating.calcAVG()
		return db.CreateRatingAverage(ctx, _rating)
	}

	for i := range _rating.Ratings {
		if _rating.Ratings[i].Star == ratingInput.Star {
			_rating.Ratings[i].Count += 1
			_rating.calcAVG()
			return db.UpdateRatingAverage(ctx, _rating)
		}
	}

	_rating.Ratings = append(_rating.Ratings, Rating{Star: ratingInput.Star, Count: 1})
	_rating.calcAVG()
	return db.CreateRatingAverage(ctx, _rating)
}
