package store

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
)

type Storage struct {
	*sql.DB
}

func (s Storage) Builder() sq.StatementBuilderType {
	return sq.StatementBuilder.PlaceholderFormat(sq.Dollar).RunWith(s.DB)
}
