package api

import (
	"encoding/json"
	"net/http"
)

func SuccessResponse(w http.ResponseWriter, httpStatus int) {
	w.WriteHeader(httpStatus)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"status": "ok"}`))

}

func FailureResponse(w http.ResponseWriter, httpStatus int, err error) {
	body, _ := json.Marshal(struct {
		Status string `json:"status"`
		Error  error  `json:"error"`
	}{
		Status: "failure",
		Error:  err,
	})

	w.WriteHeader(httpStatus)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(body)
}
