package rating

import (
	"context"
	"errors"
	"math"
)

// errors customized
var (
	ErrRatingInputInvalid = errors.New("the RatingInput is invalid. The struct cannot be nil")
)

// RatingItemRepository interface must be implemented by a database service.
type RatingItemRepository interface {
	GetRatingItemByItemId(ctx context.Context, itemId string) (*RatingItem, error)
	SaveRatingItem(ctx context.Context, ratingItem *RatingItem) error
}

// RatingItem is the representation of the model data persistence.
type RatingItem struct {
	RatingId int64
	Item     string
	Avg      float64
	Averages []Average
}

type RatingInput struct {
	ProductId     string
	OverallRating float64
}

// Average is the representation of the evaluation.
type Average struct {

	// OverallRating is the basis considered for the rating.
	// for example if your base was 5 stars, the overall rating should be 5
	OverallRating float64

	// Ratings represents the total value of Ratings of an overallRating.
	// for example, for example out of 100 reviews are 5 stars so you have 100 reviews from overallRating.
	Ratings float64
}

// calcAVG calculate the average of ratings.
func calcAVG(avgs ...Average) float64 {

	var avg, votes float64
	for _, _avg := range avgs {
		avg += _avg.OverallRating * _avg.Ratings
		votes += _avg.Ratings
	}

	result := math.Floor((avg/votes)*100) / 100
	return result
}

func PutRating(ctx context.Context, input *RatingInput, repository RatingItemRepository) error {

	if input == nil {
		return ErrRatingInputInvalid
	}

	ratingItem, err := repository.GetRatingItemByItemId(ctx, input.ProductId)
	if err != nil {
		return err
	}

	if ratingItem == nil {

		average := Average{OverallRating: input.OverallRating, Ratings: float64(1)}
		ratingItem = &RatingItem{
			Item:     input.ProductId,
			Avg:      calcAVG(average),
			Averages: []Average{average},
		}

		err = repository.SaveRatingItem(ctx, ratingItem)
		return err
	}

	avgExist := false

	for i := 0; i < len(ratingItem.Averages); i++ {

		if ratingItem.Averages[i].OverallRating == input.OverallRating {

			ratingItem.Averages[i].OverallRating += 1
			avgExist = true
			break
		}

	}

	if !avgExist {
		ratingItem.Averages = append(ratingItem.Averages, Average{OverallRating: input.OverallRating, Ratings: float64(1)})
	}

	ratingItem.Avg = calcAVG(ratingItem.Averages...)
	err = repository.SaveRatingItem(ctx, ratingItem)
	return err
}
