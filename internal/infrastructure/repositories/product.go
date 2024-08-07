package repositories

import (
	"Fridger/internal/domain/interfaces/repositories"
	"Fridger/internal/domain/models"
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
			Columns("name", "gtin", "cis", "category", "expiration_date", "is_active", "created_at").
			Values(product.Name, product.Gtin, product.Cis, product.Category, product.ExpirationDate, product.IsActive, product.CreatedAt).
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

func (r *productRepo) FindByCis(ctx context.Context, cis string) (*models.Product, error) {
	sql, params, err := helpers.QueryBuilder().
		Select("name", "gtin", "cis", "category", "expiration_date", "is_active", "created_at").
		From("product").
		Where(squirrel.Eq{"cis": cis}).
		ToSql()

	if err != nil {
		return nil, err
	}

	product := &models.Product{}
	err = pgxscan.Get(ctx, r.pool, product, sql, params...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return product, nil
}

func (r *productRepo) DeleteByCis(ctx context.Context, cis string) error {
	sql, params, err := helpers.QueryBuilder().
		Delete("product").
		Where(squirrel.Eq{"cis": cis}).
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
