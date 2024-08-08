package repositories

import (
	"Fridger/internal/domain/interfaces/repositories"
	"Fridger/internal/domain/models"
	errors2 "Fridger/internal/errors"
	"Fridger/internal/helpers"
	"context"
	"errors"
	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type productRepo struct {
	pool *pgxpool.Pool
}

func NewProductRepo(pool *pgxpool.Pool) repositories.ProductRepo {
	return &productRepo{pool: pool}
}

func (r *productRepo) Add(ctx context.Context, product *models.Product) error {
	sql, params, err :=
		helpers.QueryBuilder().
			Insert("product").
			Columns("name", "gtin", "serial", "category", "expiration_date", "is_active", "created_at").
			Values(
				product.Name,
				product.Gtin,
				product.Serial,
				product.Category,
				product.ExpirationDate,
				product.IsActive,
				product.CreatedAt).
			ToSql()

	if err != nil {
		return err
	}

	_, err = r.pool.Exec(ctx, sql, params...)
	if err != nil {
		return err
	}

	return nil
}

func (r *productRepo) GetBySerial(ctx context.Context, serial string) (*models.Product, error) {
	sql, params, err := helpers.QueryBuilder().
		Select("name", "gtin", "serial", "category", "expiration_date", "is_active", "created_at").
		From("product").
		Where(squirrel.Eq{"serial": serial}).
		ToSql()

	if err != nil {
		return nil, err
	}

	product := &models.Product{}
	err = pgxscan.Get(ctx, r.pool, product, sql, params...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors2.ErrNotFound
		}

		return nil, err
	}

	return product, nil
}

func (r *productRepo) DeleteBySerial(ctx context.Context, serial string) error {
	sql, params, err := helpers.QueryBuilder().
		Update("product").
		Set("is_active", false).
		Where(squirrel.Eq{"serial": serial}).
		ToSql()

	if err != nil {
		return err
	}

	_, err = r.pool.Exec(ctx, sql, params...)
	if err != nil {
		return err
	}

	return nil
}
