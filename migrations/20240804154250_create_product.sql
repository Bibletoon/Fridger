-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS product (
    id bigserial primary key,
    name text not null,
    gtin bigint not null,
    serial text null,
    category text,
    expiration_date timestamptz,
    is_active boolean default true,
    created_at timestamptz,
    deleted_at timestamptz null 
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS product;
-- +goose StatementEnd
