package repositories

import (
	"Fridger/internal/domain/interfaces/repositories"
	"Fridger/internal/domain/models"
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

type productRepo struct {
	pool *pgxpool.Pool
}

func NewProductRepo(pool *pgxpool.Pool) repositories.ProductRepo {
	return &productRepo{pool: pool}
}

func (r *productRepo) Save(ctx context.Context, product *models.Product) error {
	sql, params, err :=
		squirrel.
			StatementBuilder.
			PlaceholderFormat(squirrel.Dollar).
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
