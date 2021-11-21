package api

type (
	Response struct {
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
		Error   error       `json:"error"`
	}

	TokenPayload struct {
		UserID     int
		MerchantID int
		Email      string
	}
)
