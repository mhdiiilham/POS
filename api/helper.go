package api

import (
	"encoding/json"
	"net/http"
)

func SuccessResponse(w http.ResponseWriter, msg string, data interface{}, httpCode int) {
	JSON, err := json.Marshal(Response{
		Code:    httpCode,
		Message: msg,
		Data:    data,
		Error:   nil,
	})
	if err != nil {
		UnknownErrorResponse(w, err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)
	w.Write(JSON)
}

func FailedResponse(w http.ResponseWriter, err error, httpCode int) {
	JSON, jsonErr := json.Marshal(Response{
		Code:    httpCode,
		Message: err.Error(),
		Data:    nil,
		Error:   err,
	})
	if jsonErr != nil {
		UnknownErrorResponse(w, err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)
	w.Write(JSON)
}

func UnknownErrorResponse(w http.ResponseWriter, err error) {
	httpCode := http.StatusInternalServerError
	resp := Response{
		Code:    httpCode,
		Message: "internal server error",
		Error:   err,
	}

	j, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)
	w.Write(j)
}
