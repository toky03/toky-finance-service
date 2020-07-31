package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/toky03/toky-finance-accounting-service/model"
	"github.com/toky03/toky-finance-accounting-service/service"
	"net/http"
)

type AccountingService interface {
	ReadAccountsFromBook(string) ([]model.AccountTableDTO, error)
	CreateAccount(account model.AccountTableDTO) error
	CreateBooking(booking model.BookingDTO) error
}

// BookRealmHandler implementaion of Handler
type AccountingHandlerImpl struct {
	AccountingService AccountingService
	UserService       UserService
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
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}
	js, err := json.Marshal(accounts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}

func (h *AccountingHandlerImpl) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var account model.AccountTableDTO

	err := json.NewDecoder(r.Body).Decode(&account)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	err = h.AccountingService.CreateAccount(account)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusCreated)
}
func (h *AccountingHandlerImpl) ReadBookings(w http.ResponseWriter, r *http.Request) {

}
func (h *AccountingHandlerImpl) CreateBooking(w http.ResponseWriter, r *http.Request) {
	var booking model.BookingDTO

	err := json.NewDecoder(r.Body).Decode(&booking)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	err = h.AccountingService.CreateBooking(booking)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusCreated)

}
func (h *AccountingHandlerImpl) SaveAccountOption(w http.ResponseWriter, r *http.Request) {

}
