package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"sync"

	"github.com/brbarme-shop/brbarmex-rating/rating"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type repository struct {
	sourceName string
	m          sync.Mutex
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
		INNER JOIN ratings_averages ra ON ra.rating_id      = r.rating_id 
		INNER JOIN ratings_stars    rs ON rs.rating_star_id = ra.rating_star_id  
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

	r.m.Lock()
	defer r.m.Unlock()

	db, err := sql.Open("postgres", r.sourceName)
	if err != nil {
		return err
	}

	defer db.Close()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	sql := `UPDATE ratings
	SET rating_avg=$1
	WHERE rating_item_id=$2 AND rating_hash_id=$3
	`

	rows, err := tx.ExecContext(ctx, sql, average.Avg, average.RatingId, average.ItemId)
	if err != nil {
		tx.Rollback()
		return err
	}

	rowsAffected, err := rows.RowsAffected()
	if err != nil || rowsAffected <= 0 {
		tx.Rollback()
		return err
	}

	sql = `UPDATE ratings_averages
	SET rating_count=$1;
	WHERE rating_id=$2
	`

	for i := range average.AverageScore {
		rows, err := tx.ExecContext(ctx, sql, average.RatingId, average.AverageScore[i].ScorePoint)
		if err != nil {
			tx.Rollback()
			return err
		}

		rowsAffected, err := rows.RowsAffected()
		if err != nil || rowsAffected <= 0 {
			tx.Rollback()
			return err
		}
	}

	return err
}

func (r *repository) CreateAverage(ctx context.Context, average *rating.Average) error {

	r.m.Lock()
	defer r.m.Unlock()

	db, err := sql.Open("postgres", r.sourceName)
	if err != nil {
		return err
	}

	defer db.Close()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	sql := `INSERT INTO ratings
	(rating_hash_id, rating_item_id, rating_avg)
	VALUES($1, $2, $3)
	RETURNING rating_id
	`

	ratingHashId := uuid.NewString()
	var ratingId string

	row := tx.QueryRowContext(ctx, sql, ratingHashId, average.ItemId, average.Avg)
	err = row.Scan(&ratingId)
	if err != nil {
		tx.Rollback()
		return err
	}

	sql = `INSERT INTO ratings_averages
	(rating_id, rating_star_id, rating_count)
	VALUES($1, $2, $3);
	`

	for i := range average.AverageScore {
		_, err = tx.ExecContext(ctx, sql, ratingHashId, average.AverageScore[i].StarId, average.AverageScore[i].ScorePoint)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return err
}

// func NewPostgreSqlRepository(strConnection string) rating.RatingItemRepository {
// 	return &repository{dataSourceName: strConnection}
// }

func NewRatingRepository(sourceName string) rating.AverageRepository {
	return &repository{sourceName: sourceName}
}
