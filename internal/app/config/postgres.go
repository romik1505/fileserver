package config

import (
	"context"
	"log"

	"database/sql"

	_ "github.com/lib/pq"
	"github.com/romik1505/fileserver/internal/app/store"
)

func NewPostgresConnection(ctx context.Context) store.Storage {
	conn, err := sql.Open("postgres", GetValue(PostgresConnection))
	if err != nil {
		log.Fatalln(err.Error())
	}

	return store.Storage{
		DB: conn,
	}
}
