package database

import (
	"context"
	"database/sql"
	"time"

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

	dbPingErr := db.Ping()
	if dbPingErr != nil {
		logger.Error(ctx, ops, dbPingErr.Error())
		return db, dbPingErr
	}

	db.SetConnMaxLifetime(5 * time.Minute)
	return db, nil
}
