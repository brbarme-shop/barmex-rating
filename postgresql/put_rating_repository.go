package postgresql

import (
	"context"
	"database/sql"

	"github.com/brbarme-shop/brbarmex-rating/rating"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type repository struct {
	db *sql.DB
}

func (r *repository) UpdateRating(ctx context.Context, itemId string, average float64, star, count int64) error {

	sqlSmt, err := r.prepareToUpdateNewRating(ctx, star)
	if err != nil {
		return nil
	}

	defer sqlSmt.Close()

	sqlResult, err := sqlSmt.ExecContext(ctx, average, count, itemId)
	if err != nil {
		return nil
	}

	rowsAffected, err := sqlResult.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected <= 0 {
		return rating.ErrFailedToPutNewRating
	}

	return err
}

func (r *repository) ReadByItemId(ctx context.Context, itemId string) (*rating.RatingAverage, error) {

	sql := `
SELECT 
	rating_avg, rating_start_i, rating_start_i_count, rating_start_ii, rating_start_ii_count, rating_start_iii, rating_start_iii_count, rating_start_iv, rating_start_iv_count, rating_start_x, rating_start_x_count
FROM ratings_avarages
WHERE rating_item_id = $1
`
	row, err := r.db.QueryContext(ctx, sql, itemId)
	if err != nil {
		return nil, err
	}

	defer row.Close()

	if row.Next() {

		var avg float64
		var start_i, start_i_count, start_ii, start_ii_count, start_iii, start_iii_count, start_iv, start_iv_count, start_x, start_x_count int64

		err = row.Scan(&avg, &start_i, &start_i_count, &start_ii, &start_ii_count, &start_iii, &start_iii_count, &start_iv, &start_iv_count, &start_x, &start_x_count)
		if err != nil {
			return nil, err
		}

		_ratings := make([]rating.Rating, 0, 4)
		_ratings = append(_ratings, rating.Rating{Star: start_i, Count: start_i_count})
		_ratings = append(_ratings, rating.Rating{Star: start_ii, Count: start_ii_count})
		_ratings = append(_ratings, rating.Rating{Star: start_iii, Count: start_iii_count})
		_ratings = append(_ratings, rating.Rating{Star: start_iv, Count: start_iv_count})
		_ratings = append(_ratings, rating.Rating{Star: start_x, Count: start_x_count})

		ratingAverage := &rating.RatingAverage{ItemId: itemId, Average: avg, Ratings: _ratings}
		return ratingAverage, nil
	}

	return nil, rating.ErrRatingNotFound

}

func (r *repository) PutNewRating(ctx context.Context, itemId string, star int64) error {

	sqlSmt, err := r.prepareToInsertNewRating(ctx, star)
	if err != nil {
		return nil
	}

	defer sqlSmt.Close()

	sqlResult, err := sqlSmt.ExecContext(ctx, uuid.NewString(), itemId)
	if err != nil {
		return err
	}

	rowsAffected, err := sqlResult.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected <= 0 {
		return rating.ErrFailedToPutNewRating
	}

	return err
}

func (r *repository) prepareToUpdateNewRating(ctx context.Context, star int64) (*sql.Stmt, error) {

	var sql string
	switch star {
	case 1:
		sql = `UPDATE ratings_avarages
		SET rating_avg=$1, rating_start_i_count=$2
		WHERE rating_item_id=$3`
	case 2:
		sql = `UPDATE ratings_avarages
		SET rating_avg=$1, rating_start_ii_count=$2
		WHERE rating_item_id=$3`
	case 3:
		sql = `UPDATE ratings_avarages
		SET rating_avg=$1, rating_start_iii_count=$2
		WHERE rating_item_id=$3`
	case 4:
		sql = `UPDATE ratings_avarages
		SET rating_avg=$1, rating_start_iv_count=$2
		WHERE rating_item_id=$3`
	case 5:
		sql = `UPDATE ratings_avarages
		SET rating_avg=$1, rating_start_x_count=$2
		WHERE rating_item_id=$3`
	default:
		return nil, rating.ErrStartNotIdentifier
	}

	return r.db.PrepareContext(ctx, sql)
}

func (r *repository) prepareToInsertNewRating(ctx context.Context, star int64) (*sql.Stmt, error) {

	var sql string
	switch star {
	case 1:
		sql = `
		INSERT INTO ratings_avarages
				(rating_hash_id, rating_item_id, rating_avg, rating_start_i, rating_start_i_count, rating_start_ii, rating_start_ii_count, rating_start_iii, rating_start_iii_count, rating_start_iv, rating_start_iv_count, rating_start_x, rating_start_x_count)
		VALUES($1, $2, 0.01, 1, 1, 2, 0, 3, 0, 4, 0, 5, 0)`
	case 2:
		sql = `
		INSERT INTO ratings_avarages
				(rating_hash_id, rating_item_id, rating_avg, rating_start_i, rating_start_i_count, rating_start_ii, rating_start_ii_count, rating_start_iii, rating_start_iii_count, rating_start_iv, rating_start_iv_count, rating_start_x, rating_start_x_count)
		VALUES($1, $2, 0.01, 1, 0, 2, 1, 3, 0, 4, 0, 5, 0)`
	case 3:
		sql = `
		INSERT INTO ratings_avarages
				(rating_hash_id, rating_item_id, rating_avg, rating_start_i, rating_start_i_count, rating_start_ii, rating_start_ii_count, rating_start_iii, rating_start_iii_count, rating_start_iv, rating_start_iv_count, rating_start_x, rating_start_x_count)
		VALUES($1, $2, 0.01, 1, 0, 2, 0, 3, 1, 4, 0, 5, 0)`
	case 4:
		sql = `
		INSERT INTO ratings_avarages
				(rating_hash_id, rating_item_id, rating_avg, rating_start_i, rating_start_i_count, rating_start_ii, rating_start_ii_count, rating_start_iii, rating_start_iii_count, rating_start_iv, rating_start_iv_count, rating_start_x, rating_start_x_count)
		VALUES($1, $2, 0.01, 1, 0, 2, 0, 3, 0, 4, 1, 5, 0)`
	case 5:
		sql = `
		INSERT INTO ratings_avarages
				(rating_hash_id, rating_item_id, rating_avg, rating_start_i, rating_start_i_count, rating_start_ii, rating_start_ii_count, rating_start_iii, rating_start_iii_count, rating_start_iv, rating_start_iv_count, rating_start_x, rating_start_x_count)
		VALUES($1, $2, 0.01, 1, 0, 2, 0, 3, 0, 4, 0, 5, 1)`
	default:
		return nil, rating.ErrStartNotIdentifier
	}

	return r.db.PrepareContext(ctx, sql)
}

func NewRatingRepository(db *sql.DB) rating.PutRatingRepository {
	return &repository{db: db}
}
