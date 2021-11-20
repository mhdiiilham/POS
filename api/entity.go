package api

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Error   error       `json:"error"`
}

type HealthCheckResponse struct {
	DBConnected bool  `json:"DBConnected"`
	Error       error `json:"error"`
}
