package handler

import (
	"net/http"

	"github.com/toky03/toky-finance-accounting-service/model"
)

func handleError(error model.TokyError, w http.ResponseWriter) {

	if error.IsTechnicalError() {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(error.ErrorMessage()))
		return
	}
	if model.IsExistingNotFoundError(error) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(error.ErrorMessage()))
		return
	}
	if model.IsExistingValidationError(error) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(error.ErrorMessage()))
		return
	}
	if model.IsExistingBuisnessError(error) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(error.ErrorMessage()))
		return
	}
}
