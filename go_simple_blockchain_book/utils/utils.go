package utils

import (
	"encoding/json"
	"net/http"
)

func RespondBadRequest(w http.ResponseWriter, err error) {
  w.WriteHeader(http.StatusBadRequest)
  w.Write([]byte(err.Error()))
}

func RespondJson(w http.ResponseWriter, data interface{}, status int) {
  w.Header().Add("Content-Type", "application/json")
  w.WriteHeader(status)
  if err := json.NewEncoder(w).Encode(data); err != nil {
    RespondInternalError(w, err)
  }
}

func RespondInternalError(w http.ResponseWriter, err error) {
  w.WriteHeader(http.StatusInternalServerError)
  w.Write([]byte(err.Error()))
}
