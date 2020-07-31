package handler

import (
	"encoding/json"
	"github.com/toky03/toky-finance-accounting-service/model"
	"github.com/toky03/toky-finance-accounting-service/service"
	"net/http"
)

// BookRealmService interface to define Contract
type BookRealmService interface {
	FindAllBookRealms() ([]model.BookRealmDTO, error)
	CreateBookRealm(model.BookRealmDTO, string) error
}

type UserService interface {
	CreateUser(model.ApplicationUserDTO) error
	ReadAllUsers() ([]model.ApplicationUserDTO, error)
}

// BookRealmHandler implementaion of Handler
type BookRealmHandler struct {
	BookRealmService BookRealmService
	UserService      UserService
}

func CreateBookRealmHandler() *BookRealmHandler {
	return &BookRealmHandler{
		BookRealmService: service.CreateBookService(),
		UserService:      service.CreateApplicationUserService(),
	}
}

func (h *BookRealmHandler) ReadBookRealms(w http.ResponseWriter, r *http.Request) {
	bookRealms, err := h.BookRealmService.FindAllBookRealms()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}
	js, err := json.Marshal(bookRealms)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}

func (h *BookRealmHandler) CreateBookRealm(w http.ResponseWriter, r *http.Request) {
	var bookRealm model.BookRealmDTO

	err := json.NewDecoder(r.Body).Decode(&bookRealm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	err = h.BookRealmService.CreateBookRealm(bookRealm, r.Header.Get("Username"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *BookRealmHandler) UpdateBookRealm(w http.ResponseWriter, r *http.Request) {

}

func (h *BookRealmHandler) ReadAccountingUsers(w http.ResponseWriter, r *http.Request) {
	// TODO limit einfügen
	applicationUsers, err := h.UserService.ReadAllUsers()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}
	js, err := json.Marshal(applicationUsers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (h *BookRealmHandler) CreateUser(w http.ResponseWriter, r *http.Request) {

	var applicationUser model.ApplicationUserDTO

	err := json.NewDecoder(r.Body).Decode(&applicationUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	err = h.UserService.CreateUser(applicationUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusCreated)

}
