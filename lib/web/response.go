package web

import (
	"encoding/json"
	"net/http"
)

func JSON(w http.ResponseWriter, status int, v interface{}) {
	res, err := json.Marshal(v)
	if err != nil {
		ServerError(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(res)
}

func Ok(w http.ResponseWriter, msg string) {
	var data struct {
		Msg string `json:"message"`
	}
	data.Msg = msg
	JSON(w, http.StatusOK, data)
}

func BadRequest(w http.ResponseWriter, msg string) {
	var data struct {
		Error string `json:"error"`
	}
	data.Error = msg
	JSON(w, http.StatusBadRequest, data)
}

func NotFound(w http.ResponseWriter, msg string) {
	var data struct {
		Error string `json:"error"`
	}
	data.Error = msg
	JSON(w, http.StatusNotFound, data)
}

func ServerError(w http.ResponseWriter, msg string) {
	var data struct {
		Error string `json:"error"`
	}
	data.Error = msg
	JSON(w, http.StatusInternalServerError, data)
}

func UnavaliableError(w http.ResponseWriter, msg string) {
	var data struct {
		Error string `json:"error"`
	}
	data.Error = msg
	JSON(w, http.StatusServiceUnavailable, data)
}
