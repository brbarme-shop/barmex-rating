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

	var err error
	var sqlSmt *sql.Stmt

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

	rowsAffected, err := sqlResult.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected <= 0 {
		return fmt.Errorf("pau")
	}

	return err
}

func NewRatingRepository(db *sql.DB) rating.PutRatingRepository {
	return &repository{db: db}
}
