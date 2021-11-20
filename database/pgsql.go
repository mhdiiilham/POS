package database

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/mhdiiilham/POS/pkg/logger"
)

func NewPostgreSQLConnection(ctx context.Context, dns string) (*sql.DB, error) {
	const ops = "database.NewPostgreSQLConnection"
	db, openErr := sql.Open("postgres", dns)
	if openErr != nil {
		logger.Error(ctx, ops, "error trying to open database: %v", openErr)
		return nil, openErr
	}

	logger.Info(ctx, ops, "connected to database")
	return db, nil
}
