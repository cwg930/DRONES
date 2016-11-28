package controllers

import (
	"net/http"
)

func HandleCORS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin","*")
	w.Header().Set("Access-Control-Allow-Methods","GET, POST")
	w.Header().Set("Access-Control-Allow-Headers","Authorization")
	w.WriteHeader(http.StatusOK)
}
