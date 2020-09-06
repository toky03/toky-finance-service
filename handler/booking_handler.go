package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/context"
	"github.com/toky03/toky-finance-accounting-service/model"
	"github.com/toky03/toky-finance-accounting-service/service"
)

// BookRealmService interface to define Contract
type bookRealmService interface {
	FindBookRealmsPermittedForUser(userId string) ([]model.BookRealmDTO, model.TokyError)
	CreateBookRealm(model.BookRealmDTO, string) model.TokyError
}

type userService interface {
	CreateUser(model.ApplicationUserDTO) model.TokyError
	SearchUsers(limit, searchTerm string) ([]model.ApplicationUserDTO, model.TokyError)
	FindUserByUsername(userName string) (model.ApplicationUserDTO, model.TokyError)
}

// BookRealmHandler implementaion of Handler
type BookRealmHandler struct {
	BookRealmService bookRealmService
	UserService      userService
}

func CreateBookRealmHandler() *BookRealmHandler {
	return &BookRealmHandler{
		BookRealmService: service.CreateBookService(),
		UserService:      service.CreateApplicationUserService(),
	}
}

func (h *BookRealmHandler) ReadBookRealms(w http.ResponseWriter, r *http.Request) {
	userName := context.Get(r, "user-id")
	bookRealms, err := h.BookRealmService.FindBookRealmsPermittedForUser(userName.(string))
	if model.IsExisting(err) {
		handleError(err, w)
		return
	}
	js, marshalErr := json.Marshal(bookRealms)
	if marshalErr != nil {
		http.Error(w, marshalErr.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}

func (h *BookRealmHandler) CreateBookRealm(w http.ResponseWriter, r *http.Request) {
	var bookRealm model.BookRealmDTO

	decoderError := json.NewDecoder(r.Body).Decode(&bookRealm)
	if decoderError != nil {
		http.Error(w, decoderError.Error(), http.StatusUnprocessableEntity)
		return
	}
	userName := context.Get(r, "user-id")

	createRealmErr := h.BookRealmService.CreateBookRealm(bookRealm, fmt.Sprint(userName))
	if model.IsExisting(createRealmErr) {
		handleError(createRealmErr, w)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *BookRealmHandler) UpdateBookRealm(w http.ResponseWriter, r *http.Request) {

}

func (h *BookRealmHandler) ReadAccountingUsers(w http.ResponseWriter, r *http.Request) {
	queries := r.URL.Query()
	limit := queries.Get("limit")
	searchTerm := queries.Get("searchTerm")
	applicationUsers, err := h.UserService.SearchUsers(limit, searchTerm)
	if model.IsExisting(err) {
		handleError(err, w)
	}
	js, marshalError := json.Marshal(applicationUsers)
	if marshalError != nil {
		http.Error(w, marshalError.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (h *BookRealmHandler) CreateUser(w http.ResponseWriter, r *http.Request) {

	var applicationUser model.ApplicationUserDTO

	decoderError := json.NewDecoder(r.Body).Decode(&applicationUser)
	if decoderError != nil {
		http.Error(w, decoderError.Error(), http.StatusUnprocessableEntity)
		return
	}

	createUserError := h.UserService.CreateUser(applicationUser)
	if model.IsExisting(createUserError) {
		handleError(createUserError, w)
		return
	}
	w.WriteHeader(http.StatusCreated)

}
