package merchant

var (
	createNewMerchant = `
		INSERT INTO public."Merchant" ("name", created_at, updated_at)
		VALUES($1, $2, $2) RETURNING id;
	`
)
