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

	insertUser = `
		INSERT INTO public."User" (email, firstname, lastname, "password", merchant_id, created_at, updated_at, deleted_at)
		VALUES($1, $2, $3, $4, $5, $6, $6, null) RETURNING id;
	`
)
