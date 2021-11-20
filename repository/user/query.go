package user

var (
	findUserByEmail = `
		SELECT
			id,
			merchant_id,
			email,
			password,
			firstname,
			lastname,
			created_at,
			updated_at,
			deleted_at
		FROM "User"
		Where "email"=$1 AND "deleted_at" IS NULL LIMIT 1
	`
)
