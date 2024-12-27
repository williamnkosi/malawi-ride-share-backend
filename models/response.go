package models

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Error   interface{} `json:"error"`
}

func (r *Response) SuccessfulResponse(w http.ResponseWriter, message string, data interface{}) *Response {
	res := Response{
		Status:  http.StatusOK,
		Message: message,
		Data:    data,
	}

	sendJSON(w, res, http.StatusOK)

	return r
}

func (r *Response) FailedResponse(w http.ResponseWriter, statusCode int, message string, errDetails interface{}) *Response {
	res := Response{
		Status:  statusCode,
		Message: message,
		Error:   errDetails,
	}
	sendJSON(w, res, statusCode)

	return r
}

func sendJSON(w http.ResponseWriter, payload Response, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
