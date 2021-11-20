package user

import (
	"errors"
	"time"
)

type User struct {
	ID         int        `db:"id"`
	MerchantID int        `db:"merchant_id"`
	Email      string     `db:"email"`
	Password   string     `db:"password"`
	FirstName  string     `db:"firstname"`
	LastName   *string    `db:"lastname"`
	CreatedAt  time.Time  `db:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at"`
	DeletedAt  *time.Time `db:"deleted_at"`
}

var (
	ErrEmptyEmailAndPassword  error = errors.New("email or/and password can't be empty")
	ErrInvalidEmailAndPasword error = errors.New("invalid email or/and password")
)
