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

type repository struct {
	db *sql.DB
}

// UpdateRating implements rating.PutRatingRepository
func (*repository) UpdateRating(ctx context.Context, itemId string, average float64, star, count int64) error {
	panic("unimplemented")
}

func (r *repository) ReadByItemId(ctx context.Context, itemId string) (*rating.RatingAverage, error) {
	panic("unimplemented")
}

func (r *repository) PutNewRating(ctx context.Context, itemId string, average float64, star, count int64) error {

	sqlSmt, err := r.prepareToInserNewRating(ctx, itemId, average, star, count)
	if err != nil {
		return nil
	}

	sqlResult, err := sqlSmt.ExecContext(ctx, uuid.NewString(), itemId, average, star, count)
	if err != nil {
		return nil
	}

	rowsAffected, err := sqlResult.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected <= 0 {
		return fmt.Errorf("pau")
	}

	return err
}

func (r *repository) prepareToInserNewRating(ctx context.Context, itemId string, average float64, star, count int64) (*sql.Stmt, error) {

	var err error
	var sqlStmt *sql.Stmt

	switch star {
	case 1:
		sqlStmt, err = r.db.PrepareContext(ctx, ``)
	case 2:
		sqlStmt, err = r.db.PrepareContext(ctx, ``)
	case 3:
		sqlStmt, err = r.db.PrepareContext(ctx, ``)
	case 4:
		sqlStmt, err = r.db.PrepareContext(ctx, ``)
	case 5:
		sqlStmt, err = r.db.PrepareContext(ctx, ``)
	default:
		err = errors.New("")
	}

	return sqlStmt, err
}

func NewRatingRepository(db *sql.DB) rating.PutRatingRepository {
	return &repository{db: db}
}
