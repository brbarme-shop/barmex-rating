package postgresql

import (
	"context"
	"database/sql"

	"github.com/brbarme-shop/brbarmex-rating/rating"
	_ "github.com/lib/pq"
)

// repository is the representation of the persistence service
type repository struct {
	dataSourceName string
	*sql.DB
	*sql.Tx
}

func (r *repository) GetRatingItemByItemId(ctx context.Context, itemId string) (*rating.RatingItem, error) {

	var err error
	r.DB, err = sql.Open("postgres", r.dataSourceName)
	if err != nil {
		return nil, err
	}

	defer r.DB.Close()

	var ratingItem *rating.RatingItem
	err = r.DB.QueryRowContext(ctx, `SELECT * FROM ratings WHERE rating_item = $1`, itemId).Scan(&ratingItem)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	err = r.DB.QueryRowContext(ctx, `SELECT * FROM ratings_average WHERE rating_id = $1`, ratingItem.RatingId).Scan(&ratingItem.Averages)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}

	return ratingItem, nil
}

func (r *repository) SaveRatingItem(ctx context.Context, ratingItem *rating.RatingItem) error {

	var err error
	r.DB, err = sql.Open("postgres", r.dataSourceName)
	if err != nil {
		return err
	}

	defer r.DB.Close()

	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	sql := `INSERT INTO ratings 
	(rating_item, rating_avg) 
	VALUES($1, $2) RETURNING rating_id`

	row := tx.QueryRowContext(ctx, sql, ratingItem.Item, ratingItem.Avg)

	var rating_id int64
	err = row.Scan(&rating_id)
	if err != nil {
		tx.Rollback()
		return err
	}

	sql = `INSERT INTO ratings_average
	(rating_id, rating_average_overall_rating, rating_average_ratings)
	VALUES($1, $2, $3);
	`

	for _, avg := range ratingItem.Averages {
		_, err = tx.ExecContext(ctx, sql, rating_id, avg.OverallRating, avg.Ratings)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	return err
}

func NewPostgreSqlRepository(strConnection string) rating.RatingItemRepository {
	return &repository{dataSourceName: strConnection}
}
