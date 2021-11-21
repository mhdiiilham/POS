package user

import (
	"context"
	"database/sql"
	"time"

	"github.com/mhdiiilham/POS/entity/user"
	"github.com/mhdiiilham/POS/pkg/logger"
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *repository {
	return &repository{
		db: db,
	}
}

func (r *repository) FindUserByEmail(ctx context.Context, email string) (*user.User, error) {
	const ops = "repository.user.FindUserByEmail"
	var entity user.User

	row := r.db.QueryRowContext(ctx, findUserByEmail, email)
	err := row.Scan(
		&entity.ID,
		&entity.MerchantID,
		&entity.Email,
		&entity.Password,
		&entity.FirstName,
		&entity.LastName,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.DeletedAt,
	)
	if err != nil {
		logger.Error(ctx, ops, "trying to find user by email err: %v", err)
		return nil, err
	}
	return &entity, nil
}

func (r *repository) Create(ctx context.Context, entity user.User) (id int64, err error) {
	const ops = "repository.user.Create"
	now := time.Now()

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Error(ctx, ops, "error trying to begin transaction %v", err)
		return
	}

	logger.Info(ctx, ops, "creating new user")
	err = tx.QueryRowContext(
		ctx,
		insertUser,
		entity.Email,
		entity.FirstName,
		entity.LastName,
		entity.Password,
		entity.MerchantID,
		now,
	).Scan(&id)
	if err != nil {
		tx.Rollback()
		logger.Error(ctx, ops, "error trying to insert to db: %v", err)
		return
	}

	tx.Commit()
	return
}
