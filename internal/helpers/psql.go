package helpers

import "github.com/Masterminds/squirrel"

func QueryBuilder() squirrel.StatementBuilderType {
	return squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
}
