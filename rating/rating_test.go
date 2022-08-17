package rating

import (
	"context"
	"errors"
	"fmt"
	"testing"
)

func TestPutRatingAverageTest(t *testing.T) {

	ctx := context.Background()

	t.Run("Must stop execution when repository fails to get averages", func(t *testing.T) {

		err := PutRatingAverage(ctx, &AverageInput{
			ItemId:  "dummy",
			StartId: "dummy",
		}, &repositoryMock{
			readAveragesMock: func() (*Average, error) { return nil, fmt.Errorf("failed") },
		})

		if err == nil {
			t.Fail()
		}

	})

	t.Run("Must stop execution when repository fails to get star by starId", func(t *testing.T) {

		dbMock := &repositoryMock{
			readAveragesMock: func() (*Average, error) { return nil, ErrAverageNotExist },
			readStarMock:     func() (int, error) { return 0, fmt.Errorf("failed") },
		}

		err := PutRatingAverage(ctx, &AverageInput{
			ItemId:  "dummy",
			StartId: "dummy",
		}, dbMock)

		if err == nil {
			t.Fail()
		}

	})

	t.Run("Must stop execution when repository fails to save new average when don`t exist rating", func(t *testing.T) {

		dbMock := &repositoryMock{
			readAveragesMock:  func() (*Average, error) { return nil, ErrAverageNotExist },
			readStarMock:      func() (int, error) { return 0, nil },
			createAverageMock: func() error { return fmt.Errorf("failed") },
		}

		err := PutRatingAverage(ctx, &AverageInput{
			ItemId:  "dummy",
			StartId: "dummy",
		}, dbMock)

		if err == nil {
			t.Fail()
		}

	})

	t.Run("It should run successfully when there is some data in the database", func(t *testing.T) {

		dbMock := &repositoryMock{
			readAveragesMock:  func() (*Average, error) { return nil, ErrAverageNotExist },
			readStarMock:      func() (int, error) { return 0, nil },
			createAverageMock: func() error { return nil },
		}

		input := &AverageInput{
			ItemId:  "0798112345321",
			StartId: "8925314f-3dc0-48b3-8a2e-2778350f28cf",
		}

		err := PutRatingAverage(ctx, input, dbMock)
		if err == nil {
			t.Fail()
		}

	})

	t.Run("Must successfully save a new average when it doesn't exist", func(t *testing.T) {

		dbMock := &repositoryMock{
			readAveragesMock:  func() (*Average, error) { return nil, ErrAverageNotExist },
			readStarMock:      func() (int, error) { return 0, nil },
			createAverageMock: func() error { return nil },
		}

		input := &AverageInput{
			ItemId:  "0798112345321",
			StartId: "8925314f-3dc0-48b3-8a2e-2778150f28cf",
		}

		err := PutRatingAverage(ctx, input, dbMock)
		if err == nil {
			t.Fail()
		}

	})

	t.Run("Must stop process when input is invalid", func(t *testing.T) {

		err := PutRatingAverage(ctx, nil, nil)

		if err == nil {
			t.Fail()
		}

		if !errors.Is(err, ErrRatingAverageInputInvalid) {
			t.Fail()
		}

	})

	t.Run("Must successfully calculate average rating", func(t *testing.T) {
		avg := &Average{RatingId: "8925314f-3dc0-48b3-8a2e-2778350f28xf", ItemId: "0798112345321", Avg: 0, AverageScore: []AverageScore{{StarId: "8925314f-3dc0-48b3-8a2e-2778350f28cf", Star: 5, ScorePoint: 251}, {StarId: "9025314f-3dc0-48b3-8a2e-2778350f28dg", Star: 4, ScorePoint: 124}, {StarId: "0025314f-3dc0-48b3-8a2e-2778350f28eh", Star: 3, ScorePoint: 40}, {StarId: "1125314f-3dc0-48b3-8a2e-2778350f28gi", Star: 2, ScorePoint: 29}, {StarId: "2125314f-3dc0-48b3-8a2e-2778350f28hj", Star: 1, ScorePoint: 33}}}

		dbMock := &repositoryMock{
			createAverageMock: func() error { return nil },
			readStarMock:      func() (int, error) { return (5), nil },
			updateAverageMock: func() error { return nil },
			readAveragesMock:  func() (*Average, error) { return avg, nil },
		}

		input := &AverageInput{
			ItemId:  "0798112345321",
			StartId: "892531ef-3dc0-48b3-8a2e-2778350f28cf",
		}

		avgExpect := 4.11

		err := PutRatingAverage(ctx, input, dbMock)
		if err != nil {
			t.Fail()
		}

		if avg.Avg != avgExpect {
			t.Fail()
		}
	})

	t.Run("Must successfully calculate average rating and stop proccess when repository interface fails to update average", func(t *testing.T) {
		avg := &Average{RatingId: "8925314f-3dc0-48b3-8a2e-2778350f28xf", ItemId: "0798112345321", Avg: 0, AverageScore: []AverageScore{{StarId: "8925314f-3dc0-48b3-8a2e-2778350f28cf", Star: 5, ScorePoint: 251}, {StarId: "9025314f-3dc0-48b3-8a2e-2778350f28dg", Star: 4, ScorePoint: 124}, {StarId: "0025314f-3dc0-48b3-8a2e-2778350f28eh", Star: 3, ScorePoint: 40}, {StarId: "1125314f-3dc0-48b3-8a2e-2778350f28gi", Star: 2, ScorePoint: 29}, {StarId: "2125314f-3dc0-48b3-8a2e-2778350f28hj", Star: 1, ScorePoint: 33}}}

		dbMock := &repositoryMock{
			createAverageMock: func() error { return nil },
			readStarMock:      func() (int, error) { return (5), nil },
			updateAverageMock: func() error { return fmt.Errorf("failed") },
			readAveragesMock:  func() (*Average, error) { return avg, nil },
		}

		input := &AverageInput{
			ItemId:  "0798112345321",
			StartId: "892531ef-3dc0-48b3-8a2e-2778350f28cf",
		}

		avgExpect := 4.11

		err := PutRatingAverage(ctx, input, dbMock)
		if err != nil {
			t.Fail()
		}

		if avg.Avg != avgExpect {
			t.Fail()
		}
	})

	t.Run("Must successfully calculate average rating and stop proccess when repository interface fails to get star", func(t *testing.T) {
		avg := &Average{RatingId: "8925314f-3dc0-48b3-8a2e-2778350f28xf", ItemId: "0798112345321", Avg: 0, AverageScore: []AverageScore{{StarId: "8925314f-3dc0-48b3-8a2e-2778350f28cf", Star: 5, ScorePoint: 251}, {StarId: "9025314f-3dc0-48b3-8a2e-2778350f28dg", Star: 4, ScorePoint: 124}, {StarId: "0025314f-3dc0-48b3-8a2e-2778350f28eh", Star: 3, ScorePoint: 40}, {StarId: "1125314f-3dc0-48b3-8a2e-2778350f28gi", Star: 2, ScorePoint: 29}, {StarId: "2125314f-3dc0-48b3-8a2e-2778350f28hj", Star: 1, ScorePoint: 33}}}

		dbMock := &repositoryMock{
			createAverageMock: func() error { return nil },
			readStarMock:      func() (int, error) { return 0, fmt.Errorf("failed") },
			updateAverageMock: func() error { return nil },
			readAveragesMock:  func() (*Average, error) { return avg, nil },
		}

		input := &AverageInput{
			ItemId:  "0798112345321",
			StartId: "892531ef-3dc0-48b3-8a2e-2778350f28cf",
		}

		avgExpect := 4.11

		err := PutRatingAverage(ctx, input, dbMock)
		if err != nil {
			t.Fail()
		}

		if avg.Avg != avgExpect {
			t.Fail()
		}
	})

	t.Run("Must successfully calculate average rating and stop proccess when repository interface fails to update average", func(t *testing.T) {
		avg := &Average{RatingId: "8925314f-3dc0-48b3-8a2e-2778350f28xf", ItemId: "0798112345321", Avg: 0, AverageScore: []AverageScore{{StarId: "8925314f-3dc0-48b3-8a2e-2778350f28cf", Star: 5, ScorePoint: 251}, {StarId: "9025314f-3dc0-48b3-8a2e-2778350f28dg", Star: 4, ScorePoint: 124}, {StarId: "0025314f-3dc0-48b3-8a2e-2778350f28eh", Star: 3, ScorePoint: 40}, {StarId: "1125314f-3dc0-48b3-8a2e-2778350f28gi", Star: 2, ScorePoint: 29}, {StarId: "2125314f-3dc0-48b3-8a2e-2778350f28hj", Star: 1, ScorePoint: 33}}}

		dbMock := &repositoryMock{
			createAverageMock: func() error { return nil },
			readStarMock:      func() (int, error) { return (5), nil },
			updateAverageMock: func() error { return nil },
			readAveragesMock:  func() (*Average, error) { return avg, nil },
		}

		avgExpect := 4.11

		err := PutRatingAverage(ctx, &AverageInput{
			ItemId:  "0798112345321",
			StartId: "8925314f-3dc0-48b3-8a2e-2778350f28cf",
		}, dbMock)

		if err != nil {
			t.Fail()
		}

		if avg.Avg != avgExpect {
			t.Fail()
		}
	})
}

type repositoryMock struct {
	createAverageMock func() error
	readAveragesMock  func() (*Average, error)
	readStarMock      func() (int, error)
	updateAverageMock func() error
}

func (r *repositoryMock) CreateAverage(ctx context.Context, average *Average) error {
	return r.createAverageMock()
}

func (r *repositoryMock) ReadAverages(ctx context.Context, itemId string) (*Average, error) {
	return r.readAveragesMock()
}

func (r *repositoryMock) ReadStar(ctx context.Context, startid string) (int, error) {
	return r.readStarMock()
}

func (r *repositoryMock) UpdateAverage(ctx context.Context, average *Average) error {
	return r.updateAverageMock()
}
