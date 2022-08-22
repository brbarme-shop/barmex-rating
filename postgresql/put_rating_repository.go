package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/brbarme-shop/brbarmex-rating/rating"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

const (
	sqlSelectRatingAverage = `SELECT rating_hash_id, rating_item_id, rating_avg, rating_start_i, rating_start_i_count, rating_start_ii, rating_start_ii_count, rating_start_iii, rating_start_iii_count, rating_start_iv, rating_start_iv_count, rating_start_x, rating_start_x_count FROM ratings_avarages WHERE rating_item_id = $1`
	sqlUpdateRating        = `UPDATE ratings_avarages SET rating_avg=$1, rating_start_i_count=$2 ,rating_start_ii_count=$3, rating_start_iii_count=$4, rating_start_iv_count=$5, rating_start_x_count=$6 WHERE rating_hash_id=$7 AND rating_item_id=$8`
)

type ratingAverageTable struct {
	rating_hash_id         string
	rating_item_id         string
	rating_avg             float64
	rating_start_i         int64
	rating_start_i_count   int64
	rating_start_ii        int64
	rating_start_ii_count  int64
	rating_start_iii       int64
	rating_start_iii_count int64
	rating_start_iv        int64
	rating_start_iv_count  int64
	rating_start_x         int64
	rating_start_x_count   int64
}

type repository struct {
	db *sql.DB
}

// UpdateRating implements rating.PutRatingRepository
func (*repository) UpdateRating(ctx context.Context, itemId string, average float64, star, count int64) error {
	panic("unimplemented")
}

func (r *repository) ReadByItemId(ctx context.Context, itemId string) (*rating.RatingAverage, error) {

	defer func() {
		if recover := recover(); recover != nil {
			return
		}
	}()

	rows, err := r.db.QueryContext(ctx, sqlSelectRatingAverage, itemId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, rating.ErrRatingNotFound
		}
		return nil, err
	}

	defer rows.Close()

	if rows.Next() {

		var ratingSelect ratingAverageTable
		err = rows.Scan(&ratingSelect.rating_hash_id,
			&ratingSelect.rating_item_id,
			&ratingSelect.rating_avg,
			&ratingSelect.rating_start_i,
			&ratingSelect.rating_start_i_count,
			&ratingSelect.rating_start_ii,
			&ratingSelect.rating_start_ii_count,
			&ratingSelect.rating_start_iii,
			&ratingSelect.rating_start_iii_count,
			&ratingSelect.rating_start_iv,
			&ratingSelect.rating_start_iv_count,
			&ratingSelect.rating_start_x,
			&ratingSelect.rating_start_x_count)

		if err != nil {
			return nil, err
		}

		var avg rating.RatingAverage

		avg.Ratings = make([]rating.Rating, 0, 5)
		avg.Ratings = append(avg.Ratings, rating.Rating{Star: ratingSelect.rating_start_i, Count: ratingSelect.rating_start_i_count})
		avg.Ratings = append(avg.Ratings, rating.Rating{Star: ratingSelect.rating_start_ii, Count: ratingSelect.rating_start_ii_count})
		avg.Ratings = append(avg.Ratings, rating.Rating{Star: ratingSelect.rating_start_iii, Count: ratingSelect.rating_start_iii_count})
		avg.Ratings = append(avg.Ratings, rating.Rating{Star: ratingSelect.rating_start_iv, Count: ratingSelect.rating_start_iv_count})
		avg.Ratings = append(avg.Ratings, rating.Rating{Star: ratingSelect.rating_start_x, Count: ratingSelect.rating_start_x_count})

		//	avg.Id = ratingSelect.rating_hash_id
		avg.ItemId = ratingSelect.rating_item_id
		avg.Average = ratingSelect.rating_avg

		return &avg, nil
	}

	return nil, nil
}

func (r *repository) Update(ctx context.Context, ratingAverage *rating.RatingAverage) error {

	defer func() {
		if recover := recover(); recover != nil {
			return
		}
	}()

	var ratingInsert ratingAverageTable
	for i := range ratingAverage.Ratings {

		switch ratingAverage.Ratings[i].Star {
		case 1:
			ratingInsert.rating_start_i_count = ratingAverage.Ratings[i].Count
		case 2:
			ratingInsert.rating_start_ii_count = ratingAverage.Ratings[i].Count
		case 3:
			ratingInsert.rating_start_iii_count = ratingAverage.Ratings[i].Count
		case 4:
			ratingInsert.rating_start_iv_count = ratingAverage.Ratings[i].Count
		case 5:
			ratingInsert.rating_start_x_count = ratingAverage.Ratings[i].Count
		}
	}

	//ratingInsert.rating_hash_id = ratingAverage.Id
	ratingInsert.rating_avg = ratingAverage.Average
	ratingInsert.rating_item_id = ratingAverage.ItemId

	rows, err := r.db.ExecContext(ctx, sqlUpdateRating, ratingInsert.rating_avg,
		ratingInsert.rating_start_i_count,
		ratingInsert.rating_start_ii_count,
		ratingInsert.rating_start_iii_count,
		ratingInsert.rating_start_iv_count,
		ratingInsert.rating_start_x_count,
		ratingInsert.rating_hash_id,
		ratingInsert.rating_item_id,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := rows.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected <= 0 {
		return fmt.Errorf("pau")
	}

	return nil
}

func (r *repository) PutNewRating(ctx context.Context, itemId string, average float64, star, count int64) error {

	defer func() {
		if recover := recover(); recover != nil {
			return
		}
	}()

	var sqlSmt *sql.Stmt
	var err error
	switch star {
	case 1:
		sqlSmt, err = r.db.PrepareContext(ctx, ``)
	case 2:
		sqlSmt, err = r.db.PrepareContext(ctx, ``)
	case 3:
		sqlSmt, err = r.db.PrepareContext(ctx, ``)
	case 4:
		sqlSmt, err = r.db.PrepareContext(ctx, ``)
	case 5:
		sqlSmt, err = r.db.PrepareContext(ctx, ``)
	default:
		err = errors.New("")
	}

	if err != nil {
		return nil
	}

	sqlResult, err := sqlSmt.ExecContext(ctx, uuid.NewString(), itemId, average, star, count)
	if err != nil {
		return nil
	}

	if err != nil {
		return nil
	}

	rowsAffected, err := sqlResult.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected <= 0 {
		err = fmt.Errorf("pau")
	}

	return err
}

func NewRatingRepository(db *sql.DB) rating.PutRatingRepository {
	return &repository{db: db}
}
