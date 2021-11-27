package merchant

import "context"

type Repository interface {
	Create(ctx context.Context, entity Merchant) (id int, err error)
}
