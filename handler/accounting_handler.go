package handler

import (
	"encoding/json"
	"net/http"

	"github.com/toky03/toky-finance-accounting-service/model"
)

type AccountingService interface {
	ReadAccountsFromBook(string) ([]model.AccountTableDTO, model.TokyError)
	ReadAccountOptionsFromBook(string) ([]model.AccountOptionDTO, model.TokyError)
	ReadBookings(string) ([]model.BookingDTO, model.TokyError)
	CreateAccount(bookID string, account model.AccountOptionDTO) model.TokyError
	UpdateAccount(accountID string, account model.AccountOptionDTO) model.TokyError
	DeleteAccount(accountID string) model.TokyError
	CreateBooking(booking model.BookingDTO) model.TokyError
	UpdateBooking(bookingID string, booking model.BookingDTO) model.TokyError
	DeleteBooking(bookingID string) model.TokyError
	ReadClosingStatements(bookID string) (model.ClosingSheetStatements, model.TokyError)
}

// BookRealmHandler implementaion of Handler
type accountingHandlerImpl struct {
	AccountingService AccountingService
	UserService       userService
}

func CreateAccountingHandler(accountingService AccountingService, userService userService) *accountingHandlerImpl {
	return &accountingHandlerImpl{
		AccountingService: accountingService,
		UserService:       userService,
	}
}

func (h *accountingHandlerImpl) ReadAccounts(w http.ResponseWriter, r *http.Request) {
	bookID := r.PathValue("bookID")
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

func (h *accountingHandlerImpl) ReadAccountOptions(w http.ResponseWriter, r *http.Request) {
	bookID := r.PathValue("bookID")
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

func (h *accountingHandlerImpl) ReadClosingStatements(w http.ResponseWriter, r *http.Request) {
	bookID := r.PathValue("bookID")
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

func (h *accountingHandlerImpl) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var account model.AccountOptionDTO
	bookID := r.PathValue("bookID")

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

func (h *accountingHandlerImpl) UpdateAccount(w http.ResponseWriter, r *http.Request) {
	var account model.AccountOptionDTO
	accountID := r.PathValue("accountID")

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
func (h *accountingHandlerImpl) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	accountID := r.PathValue("accountID")
	accountDeletionError := h.AccountingService.DeleteAccount(accountID)
	if model.IsExisting(accountDeletionError) {
		handleError(accountDeletionError, w)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (h *accountingHandlerImpl) ReadBookings(w http.ResponseWriter, r *http.Request) {

	bookID := r.PathValue("bookID")
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
func (h *accountingHandlerImpl) CreateBooking(w http.ResponseWriter, r *http.Request) {
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

func (h *accountingHandlerImpl) UpdateBooking(w http.ResponseWriter, r *http.Request) {
	var booking model.BookingDTO
	bookingId := r.PathValue("bookID")

	err := json.NewDecoder(r.Body).Decode(&booking)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	bookingUpdateError := h.AccountingService.UpdateBooking(bookingId, booking)
	if model.IsExisting(bookingUpdateError) {
		handleError(bookingUpdateError, w)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *accountingHandlerImpl) DeleteBooking(w http.ResponseWriter, r *http.Request) {
	bookingId := r.PathValue("bookID")
	bookingDeletionError := h.AccountingService.DeleteBooking(bookingId)
	if model.IsExisting(bookingDeletionError) {
		handleError(bookingDeletionError, w)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (h *accountingHandlerImpl) SaveAccountOption(w http.ResponseWriter, r *http.Request) {

}
