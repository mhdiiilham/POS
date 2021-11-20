package user

import "context"

type Repository interface {
	FindUserByEmail(ctx context.Context, email string) (*User, error)
}

type Service interface {
	Login(ctx context.Context, email, password string) (accessToken string, err error)
}
