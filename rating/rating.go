package rating

import (
	"context"
	"errors"
)

var (
	ErrRatingAverageInputInvalid = errors.New("the RatingAverageInput is invalid. The struct and your fields cann't be NIL. Check all fields or see more about in documentation")
	ErrAverageNotExist           = errors.New("the RatingAverageInput not exists")
	ErrAverageRepositoryInvalid  = errors.New("the AverageRepository can't be NIL")
	ErrStarIdNotFound            = errors.New("the StarId not containt star")
)

type AverageRepository interface {
	CreateAverage(ctx context.Context, average *Average) error
	ReadAverages(ctx context.Context, itemId string) (*Average, error)
	ReadStar(ctx context.Context, startid string) (int, error)
	UpdateAverage(ctx context.Context, average *Average) error
}

func PutRatingAverage(ctx context.Context, averageInput *AverageInput, db AverageRepository) error {

	defer func() {
		if recover := recover(); recover != nil {
			return
		}
	}()

	if averageInput == nil || len(averageInput.ItemId) == 0 || len(averageInput.StartId) == 0 {
		return ErrRatingAverageInputInvalid
	}

	averages, err := db.ReadAverages(ctx, averageInput.ItemId)
	if err != nil {
		if !errors.Is(err, ErrAverageNotExist) {
			return err
		}
	}

	if averages == nil || len(averages.AverageScore) == 0 {

		var star int
		star, err = db.ReadStar(ctx, averageInput.StartId)
		if err != nil {
			return err
		}

		averages = &Average{
			ItemId: averageInput.ItemId,
			AverageScore: []AverageScore{{
				ScorePoint: 1,
				StarId:     averageInput.ItemId,
				Star:       star,
			}},
		}

		averages.Avg = calcAVG(averages.AverageScore...)

		err = db.CreateAverage(ctx, averages)
		if err != nil {
			return err
		}

		return err
	}

	var updated bool
	for i := range averages.AverageScore {
		if averages.AverageScore[i].StarId == averageInput.StartId {
			averages.AverageScore[i].ScorePoint += 1
			updated = true
			break
		}
	}

	if !updated {

		var star int
		star, err = db.ReadStar(ctx, averageInput.StartId)
		if err != nil {
			return err
		}

		averages.AverageScore = append(averages.AverageScore, AverageScore{
			ScorePoint: 1,
			StarId:     averageInput.StartId,
			Star:       star,
		})
	}

	averages.Avg = calcAVG(averages.AverageScore...)

	err = db.UpdateAverage(ctx, averages)
	if err != nil {
		return err
	}

	return err
}
