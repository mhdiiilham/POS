package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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

func (r *repository) Get(ctx context.Context, merchantID int, opts *user.RepositoryGetUserPaginationOptions) (users []user.User, totalData int, err error) {
	const ops = "user.repository.Get"
	total := struct {
		totalUser int64 `db:"totalUsers"`
	}{}

	logger.Info(ctx, ops, "get users of merchant %d", merchantID)
	query := getUserByMerchantID

	if opts != nil {
		query = fmt.Sprintf(`%s AND "User".id > %d LIMIT %d`, getUserByMerchantID, opts.Cursor, opts.Limit)
	}

	row := r.db.QueryRowContext(ctx, countAllUsersInMerchantID, merchantID)
	err = row.Scan(&total.totalUser)
	if err != nil {
		logger.Error(ctx, ops, "unexpected error %v", err)
		return
	}

	rows, err := r.db.QueryContext(ctx, query, merchantID)
	if err != nil {
		logger.Error(ctx, ops, "unexpected error: %v", err)
		return
	}
	defer rows.Close()

	logger.Info(ctx, ops, "scanning get users result rows")
	for rows.Next() {
		var u user.User
		errScan := rows.Scan(
			&u.ID,
			&u.MerchantID,
			&u.Email,
			&u.FirstName,
			&u.LastName,
		)
		if errScan != nil {
			logger.Error(ctx, ops, "unexpected error while scanning rows %v", err)
			err = errScan
			return
		}
		users = append(users, u)
	}
	logger.Info(ctx, ops, "scanning rows completed")

	totalData = int(total.totalUser)
	return
}

func (r *repository) Remove(ctx context.Context, userID int) (err error) {
	const ops = "repository.user.Remove"
	var tx *sql.Tx
	var res sql.Result
	var rowsAffected int64

	tx, err = r.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Error(ctx, ops, "error trying to begin db tx: %v", err)
		return
	}

	res, err = tx.ExecContext(ctx, deleteUserFromID, time.Now(), userID)
	if err != nil {
		tx.Rollback()
		return
	}

	rowsAffected, err = res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return
	}

	if rowsAffected == 0 {
		tx.Rollback()
		err = sql.ErrNoRows
		return
	}

	return tx.Commit()
}

func (r *repository) GetUser(ctx context.Context, userID int) (entity user.User, err error) {
	const ops = "repository.user.GetUser"

	err = r.db.QueryRowContext(ctx, getUser, userID).Scan(
		&entity.ID,
		&entity.MerchantID,
		&entity.Email,
		&entity.FirstName,
		&entity.LastName,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.DeletedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = user.ErrUserNotFound
			return
		}

		logger.Error(ctx, ops, "error r.db.QueryRowContext %v", err)
		return
	}

	return
}
