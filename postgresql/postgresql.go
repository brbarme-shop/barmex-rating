package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"

	"github.com/brbarme-shop/brbarmex-rating/rating"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
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
	m  sync.Mutex
	db *sql.DB
}

func (r *repository) ReadRatingAverages(ctx context.Context, itemId string) (*rating.RatingAverage, error) {

	defer func() {
		if recover := recover(); recover != nil {
			return
		}
	}()

	r.m.Lock()
	defer r.m.Unlock()

	rows, err := r.db.QueryContext(ctx, `SELECT rating_hash_id, rating_item_id, rating_avg, rating_start_i, rating_start_i_count, rating_start_ii, rating_start_ii_count, rating_start_iii, rating_start_iii_count, rating_start_iv, rating_start_iv_count, rating_start_x, rating_start_x_count
	FROM ratings_avarages
	WHERE rating_item_id = $1`, itemId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, rating.ErrAverageNotExist
		}
		return nil, err
	}

	defer rows.Close()
	var avg rating.RatingAverage

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

		avg.Id = ratingSelect.rating_hash_id
		avg.ItemId = ratingSelect.rating_item_id
		avg.Average = ratingSelect.rating_avg
		avg.Ratings = append(avg.Ratings, rating.Rating{Star: ratingSelect.rating_start_i, Count: ratingSelect.rating_start_i_count})
		avg.Ratings = append(avg.Ratings, rating.Rating{Star: ratingSelect.rating_start_ii, Count: ratingSelect.rating_start_ii_count})
		avg.Ratings = append(avg.Ratings, rating.Rating{Star: ratingSelect.rating_start_iii, Count: ratingSelect.rating_start_iii_count})
		avg.Ratings = append(avg.Ratings, rating.Rating{Star: ratingSelect.rating_start_iv, Count: ratingSelect.rating_start_iv_count})
		avg.Ratings = append(avg.Ratings, rating.Rating{Star: ratingSelect.rating_start_x, Count: ratingSelect.rating_start_x_count})

		return &avg, nil
	}

	return nil, nil
}

func (r *repository) UpdateRatingAverage(ctx context.Context, ratingAverage *rating.RatingAverage) error {

	defer func() {
		if recover := recover(); recover != nil {
			return
		}
	}()

	r.m.Lock()
	defer r.m.Unlock()

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

	ratingInsert.rating_hash_id = ratingAverage.Id
	ratingInsert.rating_avg = ratingAverage.Average
	ratingInsert.rating_item_id = ratingAverage.ItemId

	rows, err := r.db.ExecContext(ctx, `UPDATE ratings_avarages
	SET rating_avg=$1, rating_start_i_count=$2 ,rating_start_ii_count=$3, rating_start_iii_count=$4, rating_start_iv_count=$5, rating_start_x_count=$6
	WHERE rating_hash_id=$7 AND rating_item_id=$8
	`, ratingInsert.rating_avg,
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

func (r *repository) CreateRatingAverage(ctx context.Context, ratingAverage *rating.RatingAverage) error {

	defer func() {
		if recover := recover(); recover != nil {
			return
		}
	}()

	r.m.Lock()
	defer r.m.Unlock()

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	var ratingInsert ratingAverageTable
	for i := range ratingAverage.Ratings {

		switch ratingAverage.Ratings[i].Star {
		case 1:
			ratingInsert.rating_start_i = 1
			ratingInsert.rating_start_ii = 2
			ratingInsert.rating_start_iii = 3
			ratingInsert.rating_start_iv = 4
			ratingInsert.rating_start_x = 5
			ratingInsert.rating_start_i_count = 1
		case 2:
			ratingInsert.rating_start_i = 1
			ratingInsert.rating_start_ii = 2
			ratingInsert.rating_start_iii = 3
			ratingInsert.rating_start_iv = 4
			ratingInsert.rating_start_x = 5
			ratingInsert.rating_start_ii_count = 1
		case 3:
			ratingInsert.rating_start_i = 1
			ratingInsert.rating_start_ii = 2
			ratingInsert.rating_start_iii = 3
			ratingInsert.rating_start_iv = 4
			ratingInsert.rating_start_x = 5
			ratingInsert.rating_start_iii_count = 1
		case 4:
			ratingInsert.rating_start_i = 1
			ratingInsert.rating_start_ii = 2
			ratingInsert.rating_start_iii = 3
			ratingInsert.rating_start_iv = 4
			ratingInsert.rating_start_x = 5
			ratingInsert.rating_start_iv_count = 1
		case 5:
			ratingInsert.rating_start_i = 1
			ratingInsert.rating_start_ii = 2
			ratingInsert.rating_start_iii = 3
			ratingInsert.rating_start_iv = 4
			ratingInsert.rating_start_x = 5
			ratingInsert.rating_start_x_count = 1
		}
	}

	ratingInsert.rating_hash_id = uuid.NewString()
	ratingInsert.rating_avg = ratingAverage.Average
	ratingInsert.rating_item_id = ratingAverage.ItemId

	if row, err := tx.ExecContext(ctx,
		`INSERT INTO ratings_avarages
	(rating_hash_id, rating_item_id, rating_avg, rating_start_i, rating_start_i_count, rating_start_ii, rating_start_ii_count, rating_start_iii, rating_start_iii_count, rating_start_iv, rating_start_iv_count, rating_start_x, rating_start_x_count)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`,
		ratingInsert.rating_hash_id,
		ratingInsert.rating_item_id,
		ratingInsert.rating_avg,
		ratingInsert.rating_start_i,
		ratingInsert.rating_start_i_count,
		ratingInsert.rating_start_ii,
		ratingInsert.rating_start_ii_count,
		ratingInsert.rating_start_iii,
		ratingInsert.rating_start_iii_count,
		ratingInsert.rating_start_iv,
		ratingInsert.rating_start_iv_count,
		ratingInsert.rating_start_x,
		ratingInsert.rating_start_x_count); err != nil {
		tx.Rollback()
		return err
	} else {

		rowsAffected, err := row.RowsAffected()
		if err != nil {
			tx.Rollback()
			return err
		}

		if rowsAffected <= 0 {
			return fmt.Errorf("pau")
		}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return err
}

func NewRatingRepository(db *sql.DB) rating.RatingAverageRepository {
	return &repository{db: db}
}
