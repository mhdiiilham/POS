package user

import "context"

type Repository interface {
	FindUserByEmail(ctx context.Context, email string) (*User, error)
	Create(ctx context.Context, entity User) (id int64, err error)
	Get(ctx context.Context, merchantID int, opts *RepositoryGetUserPaginationOptions) (users []User, totalData int, err error)
}
