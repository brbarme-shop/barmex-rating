package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"sync"

	"github.com/brbarme-shop/brbarmex-rating/rating"
	_ "github.com/lib/pq"
)

type repository struct {
	sourceName string
	m          sync.Mutex
}

func (r *repository) CreateAverage(ctx context.Context, average *rating.Average) error {
	return nil
}

func (r *repository) ReadAverages(ctx context.Context, itemId string) (*rating.Average, error) {

	r.m.Lock()
	defer r.m.Unlock()

	db, err := sql.Open("postgres", r.sourceName)
	if err != nil {
		return nil, err
	}

	defer db.Close()

	rows, err := db.QueryContext(ctx, `SELECT rating_hash_id, rating_item_id, rating_avg FROM ratings WHERE rating_item_id = $1`, itemId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, rating.ErrAverageNotExist
		}
		return nil, err
	}

	defer rows.Close()
	var avg *rating.Average

	for rows.Next() {
		err = rows.Scan(&avg.RatingId, &avg.ItemId, &avg.Avg)
		if err != nil {
			return nil, err
		}
	}

	rows, err = db.QueryContext(ctx, `SELECT rs.rating_star_id, rs.rating_star, ra.rating_count  FROM ratings r 
		LEFT JOIN ratings_averages ra on ra.rating_id  = r.rating_id 
		LEFT JOIN ratings_stars rs on rs.rating_star_id = ra.rating_star_id  
		WHERE r.rating_item_id = $1
		GROUP BY  rs.rating_star_id, rs.rating_star , r.rating_avg , ra.rating_count`, itemId)

	if err != nil {
		return nil, err
	}

	var avgScore *rating.AverageScore
	for rows.Next() {
		err = rows.Scan(&avgScore.StarId, &avgScore.Star, &avgScore.ScorePoint)
		if err != nil {
			return nil, err
		}

		avg.AverageScore = append(avg.AverageScore, *avgScore)
	}

	return avg, nil
}

func (r *repository) ReadStar(ctx context.Context, startid string) (int, error) {

	r.m.Lock()
	defer r.m.Unlock()

	db, err := sql.Open("postgres", r.sourceName)
	if err != nil {
		return -1, err
	}

	defer db.Close()

	rows, err := db.QueryContext(ctx, `SELECT rating_star FROM ratings_stars WHERE rating_star_id = $1`, startid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return -1, rating.ErrStarIdNotFound
		}
		return -1, err
	}

	defer rows.Close()
	var star int

	for rows.Next() {
		err = rows.Scan(&star)
		if err != nil {
			return -1, err
		}
	}

	return star, nil
}

func (r *repository) UpdateAverage(ctx context.Context, average *rating.Average) error {
	return nil
}

func NewRatingRepository(sourceName string) rating.AverageRepository {
	return &repository{sourceName: sourceName}
}

// func (r *repository) SaveRating(ctx context.Context, ratingItem *rating.RatingItem) error {

// 	db, err := sql.Open("postgres", r.dataSourceName)
// 	if err != nil {
// 		return err
// 	}

// 	defer db.Close()

// 	tx, err := db.BeginTx(ctx, nil)
// 	if err != nil {
// 		return err
// 	}

// 	sql := `INSERT INTO ratings
// 	(rating_item, rating_avg)
// 	VALUES($1, $2) RETURNING rating_id`

// 	row := tx.QueryRowContext(ctx, sql, ratingItem.Item, ratingItem.Avg)

// 	var rating_id int64
// 	err = row.Scan(&rating_id)
// 	if err != nil {
// 		tx.Rollback()
// 		return err
// 	}

// 	sql = `INSERT INTO ratings_average
// 	(rating_id, rating_average_overall_rating, rating_average_ratings)
// 	VALUES($1, $2, $3);
// 	`

// 	for _, avg := range ratingItem.Averages {
// 		_, err = tx.ExecContext(ctx, sql, rating_id, avg.OverallRating, avg.Ratings)
// 		if err != nil {
// 			tx.Rollback()
// 			return err
// 		}
// 	}

// 	err = tx.Commit()
// 	if err != nil {
// 		tx.Rollback()
// 		return err
// 	}

// 	return err
// }

// func (r *repository) UpdateRating(ctx context.Context, ratingItem *rating.RatingItem) error {

// 	db, err := sql.Open("postgres", r.dataSourceName)
// 	if err != nil {
// 		return err
// 	}

// 	defer db.Close()

// 	tx, err := db.BeginTx(ctx, nil)
// 	if err != nil {
// 		return err
// 	}

// 	sql := `UPDATE ratings
// 	SET rating_avg=$1
// 	WHERE rating_id=$2;
// 	`
// 	_, err = tx.ExecContext(ctx, sql, ratingItem.RatingId, ratingItem.Avg)
// 	if err != nil {
// 		tx.Rollback()
// 		return err
// 	}

// 	sql = `UPDATE ratings_average
// 	SET rating_average_ratings=$1
// 	WHERE rating_id=$2 AND rating_average_overall_rating=$3;
// 	`

// 	for i := range ratingItem.Averages {
// 		_, err = tx.ExecContext(ctx, sql, ratingItem.RatingId, ratingItem.Averages[i].Ratings, ratingItem.Averages[i].Ratings)
// 		if err != nil {
// 			tx.Rollback()
// 			return err
// 		}
// 	}

// 	err = tx.Commit()
// 	if err != nil {
// 		tx.Rollback()
// 		return err
// 	}

// 	return err
// }

// func NewPostgreSqlRepository(strConnection string) rating.RatingItemRepository {
// 	return &repository{dataSourceName: strConnection}
// }
