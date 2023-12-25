package triprepo

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"trip/models"

	sq "github.com/Masterminds/squirrel"
)

type Params struct {
	Id      string
	UserId  string
	OfferId string
	Status  string
}

type TripRepo struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *TripRepo {
	return &TripRepo{db: db}
}

func (r *TripRepo) WithNewTx(ctx context.Context, opts *sql.TxOptions, f func(ctx context.Context) error) error {
	tx, err := r.db.BeginTx(ctx, opts)
	if err != nil {
		return err
	}

	defer func() {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil && !errors.Is(rollbackErr, sql.ErrTxDone) {
			err = rollbackErr
		}
	}()

	fErr := f(ctx)
	if fErr != nil {
		_ = tx.Rollback()
		return fErr
	}

	return tx.Commit()
}

func (r *TripRepo) Get(ctx context.Context, params *Params) ([]models.Trip, error) {
	query := sq.Select("*").
		From("trips").
		PlaceholderFormat(sq.Dollar)

	if params.Id != "" {
		query = query.Where(sq.Eq{"author": params.Id})
	}

	if params.UserId != "" {
		query = query.Where(sq.Eq{"user_id": params.UserId})
	}
	if params.OfferId != "" {
		query = query.Where(sq.Eq{"offer_id": params.OfferId})
	}
	if params.Status != "" {
		query = query.Where(sq.Eq{"status": params.Status})
	}

	sql, args, err := query.ToSql()
	rows, err := r.db.QueryxContext(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	news := make([]models.Trip, 0)

	for rows.Next() {
		n := models.Trip{}

		if err = rows.StructScan(&n); err != nil {
			return nil, err
		}

		news = append(news, n)
	}

	return news, nil
}

func (r *TripRepo) Create(ctx context.Context, n *models.Trip) error {
	sql, args, err := sq.
		Insert("trips").Columns("id", "user_id", "offer_id", "status").
		Values(n.Id, n.UserId, n.OfferId, n.Status).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	var id int
	row := r.db.QueryRowContext(ctx, sql, args...)
	if err = row.Scan(&id); err != nil {
		return err
	}

	return nil
}
