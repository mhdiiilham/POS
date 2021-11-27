package merchant

import "time"

type Merchant struct {
	ID        int        `db:"id" json:"id"`
	Name      string     `db:"name" json:"name"`
	Logo      string     `db:"name" json:"logo"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"-"`
}
