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

// RatingInput is to input in PutRating service function
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

// PutRating exec business logic to add rating count of the product
func PutRating(ctx context.Context, ratingInput *RatingInput, repository RatingItemRepository) error {

	if ratingInput == nil {
		return ErrRatingInputInvalid
	}

	rating, err := repository.GetRatingItemByItemId(ctx, ratingInput.ProductId)
	if err != nil {
		return err
	}

	if rating == nil {

		average := Average{OverallRating: ratingInput.OverallRating, Ratings: float64(1)}
		rating = &RatingItem{
			Item:     ratingInput.ProductId,
			Avg:      calcAVG(average),
			Averages: []Average{average},
		}

		err = repository.SaveRatingItem(ctx, rating)
		return err
	}

	avgExist := false
	for i := 0; i < len(rating.Averages); i++ {

		if rating.Averages[i].OverallRating == ratingInput.OverallRating {

			rating.Averages[i].OverallRating += 1
			avgExist = true
			break
		}

	}

	if !avgExist {
		rating.Averages = append(rating.Averages, Average{OverallRating: ratingInput.OverallRating, Ratings: float64(1)})
	}

	rating.Avg = calcAVG(rating.Averages...)
	err = repository.SaveRatingItem(ctx, rating)
	return err
}
