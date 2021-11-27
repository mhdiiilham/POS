package merchant

type Merchant struct {
	ID   int    `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
	Logo string `db:"name" json:"logo"`
}
