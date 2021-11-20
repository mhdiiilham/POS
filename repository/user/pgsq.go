package user

import (
	"context"
	"database/sql"

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
