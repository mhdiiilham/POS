package user

import (
	"errors"
	"time"
)

type User struct {
	ID         int        `db:"id" json:"id"`
	MerchantID int        `db:"merchant_id" json:"merchantID"`
	Email      string     `db:"email" json:"email"`
	Password   string     `db:"password" json:"-"`
	FirstName  string     `db:"firstname" json:"firstname"`
	LastName   *string    `db:"lastname" json:"lastname"`
	CreatedAt  time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt  *time.Time `db:"deleted_at" json:"-"`
}

var (
	ErrEmptyEmailAndPassword   error = errors.New("email or/and password can't be empty")
	ErrInvalidEmailAndPasword  error = errors.New("invalid email or/and password")
	ErrInvalidCreateParameters error = errors.New("failed creating user due to invalid parameters")
	ErrEmailNotUnique          error = errors.New("email is already registered")
)
