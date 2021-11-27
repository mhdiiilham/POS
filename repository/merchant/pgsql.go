package merchant

import (
	"context"
	"database/sql"
	"time"

	"github.com/mhdiiilham/POS/entity/merchant"
	"github.com/mhdiiilham/POS/pkg/logger"
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, entity merchant.Merchant) (int, error) {
	const ops = "repository.merchant.Create"
	var id int64

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Error(ctx, ops, "error begin db tx: %v", err)
		return 0, err
	}

	err = tx.QueryRowContext(ctx, createNewMerchant, entity.Name, time.Now()).Scan(&id)
	if err != nil {
		tx.Rollback()
		logger.Error(ctx, ops, "error inserting merchant to db %v", err)
		return 0, err
	}

	tx.Commit()
	return int(id), nil
}
