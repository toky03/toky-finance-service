package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/toky03/toky-finance-accounting-service/model"
	"github.com/toky03/toky-finance-accounting-service/service"
)

type AccountingService interface {
	ReadAccountsFromBook(string) ([]model.AccountTableDTO, model.TokyError)
	ReadAccountOptionsFromBook(string) ([]model.AccountOptionDTO, model.TokyError)
	ReadBookings(string) ([]model.BookingDTO, model.TokyError)
	CreateAccount(bookID string, account model.AccountOptionDTO) model.TokyError
	UpdateAccount(accountID string, account model.AccountOptionDTO) model.TokyError
	CreateBooking(booking model.BookingDTO) model.TokyError
	ReadClosingStatements(bookID string) (model.ClosingSheetStatements, model.TokyError)
}

// BookRealmHandler implementaion of Handler
type AccountingHandlerImpl struct {
	AccountingService AccountingService
	UserService       userService
}

func CreateAccountingHandler() *AccountingHandlerImpl {
	return &AccountingHandlerImpl{
		AccountingService: service.CreateAccountingService(),
		UserService:       service.CreateApplicationUserService(),
	}
}

func (h *AccountingHandlerImpl) ReadAccounts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookID := vars["bookID"]
	accounts, err := h.AccountingService.ReadAccountsFromBook(bookID)
	if model.IsExistingNotFoundError(err) {
		handleError(err, w)
		return
	}
	js, marshalError := json.Marshal(accounts)
	if err != nil {
		http.Error(w, marshalError.Error(), http.StatusUnprocessableEntity)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (h *AccountingHandlerImpl) ReadAccountOptions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookID := vars["bookID"]
	accounts, err := h.AccountingService.ReadAccountOptionsFromBook(bookID)
	if model.IsExisting(err) {
		handleError(err, w)
		return
	}
	js, marshalError := json.Marshal(accounts)
	if marshalError != nil {
		http.Error(w, marshalError.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}

func (h *AccountingHandlerImpl) ReadClosingStatements(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookID := vars["bookID"]
	ClosingSheetStatements, err := h.AccountingService.ReadClosingStatements(bookID)
	if model.IsExisting(err) {
		handleError(err, w)
		return
	}
	js, marshalError := json.Marshal(ClosingSheetStatements)
	if marshalError != nil {
		http.Error(w, marshalError.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (h *AccountingHandlerImpl) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var account model.AccountOptionDTO
	vars := mux.Vars(r)
	bookID := vars["bookID"]

	decoderError := json.NewDecoder(r.Body).Decode(&account)
	if decoderError != nil {
		http.Error(w, decoderError.Error(), http.StatusBadRequest)
		return
	}

	accountCreationError := h.AccountingService.CreateAccount(bookID, account)
	if model.IsExisting(accountCreationError) {
		handleError(accountCreationError, w)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *AccountingHandlerImpl) UpdateAccount(w http.ResponseWriter, r *http.Request) {
	var account model.AccountOptionDTO
	vars := mux.Vars(r)
	accountID := vars["accountID"]

	decoderError := json.NewDecoder(r.Body).Decode(&account)
	if decoderError != nil {
		http.Error(w, decoderError.Error(), http.StatusBadRequest)
		return
	}

	accountCreationError := h.AccountingService.UpdateAccount(accountID, account)
	if model.IsExisting(accountCreationError) {
		handleError(accountCreationError, w)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *AccountingHandlerImpl) ReadBookings(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	bookID := vars["bookID"]
	bookings, err := h.AccountingService.ReadBookings(bookID)
	if model.IsExisting(err) {
		handleError(err, w)
		return
	}
	js, marshalError := json.Marshal(bookings)
	if marshalError != nil {
		http.Error(w, marshalError.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
func (h *AccountingHandlerImpl) CreateBooking(w http.ResponseWriter, r *http.Request) {
	var booking model.BookingDTO

	err := json.NewDecoder(r.Body).Decode(&booking)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	bookingCreationError := h.AccountingService.CreateBooking(booking)
	if model.IsExisting(bookingCreationError) {
		handleError(bookingCreationError, w)
		return
	}
	w.WriteHeader(http.StatusCreated)

}
func (h *AccountingHandlerImpl) SaveAccountOption(w http.ResponseWriter, r *http.Request) {

}
