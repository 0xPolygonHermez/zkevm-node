package database

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

func NewSQLDB(dbName, dbUser, dbPassword, dbHost, dbPort string) (*pgxpool.Pool, error) {
	return pgxpool.Connect(context.Background(), "postgres://"+dbUser+":"+dbPassword+"@"+dbHost+":"+dbPort+"/"+dbName)
}
