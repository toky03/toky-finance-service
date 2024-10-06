package handler

import (
	"encoding/json"
	"net/http"

	"github.com/toky03/toky-finance-accounting-service/model"
)

// BookRealmService interface to define Contract
type bookRealmService interface {
	FindBookRealmsPermittedForUser(userId string) ([]model.BookRealmDTO, model.TokyError)
	CreateBookRealm(model.BookRealmDTO, string) model.TokyError
	FindBookRealmById(bookId string) (bookRealmDto model.BookRealmDTO, err model.TokyError)
	DeleteBookRealm(bookId string) model.TokyError
	UpdateBookRealm(bookRealm model.BookRealmDTO, bookID string) model.TokyError
}

type userService interface {
	CreateUser(model.ApplicationUserDTO) model.TokyError
	SearchUsers(limit, searchTerm string) ([]model.ApplicationUserDTO, model.TokyError)
	FindUserByUsername(userName string) (model.ApplicationUserDTO, model.TokyError)
	HasWriteAccessFromBook(userId, bookId string) (bool, model.TokyError)
	IsOwnerOfBook(userId, bookId string) (bool, model.TokyError)
}

// bookRealmHandler implementaion of Handler
type bookRealmHandler struct {
	bookRealmService bookRealmService
	userService      userService
}

func CreateBookRealmHandler(bookRealmService bookRealmService, userService userService) *bookRealmHandler {
	return &bookRealmHandler{
		bookRealmService: bookRealmService,
		userService:      userService,
	}
}

func (h *bookRealmHandler) ReadBookRealmById(w http.ResponseWriter, r *http.Request) {
	bookID := r.PathValue("bookID")
	bookRealm, err := h.bookRealmService.FindBookRealmById(bookID)
	if model.IsExisting(err) {
		handleError(err, w)
		return
	}
	js, marshalErr := json.Marshal(bookRealm)
	if marshalErr != nil {
		http.Error(w, marshalErr.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}

func (h *bookRealmHandler) ReadBookRealms(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value(USER_ID).(string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing " + USER_ID))
	}
	bookRealms, err := h.bookRealmService.FindBookRealmsPermittedForUser(userId)
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

func (h *bookRealmHandler) CreateBookRealm(w http.ResponseWriter, r *http.Request) {
	var bookRealm model.BookRealmDTO

	decoderError := json.NewDecoder(r.Body).Decode(&bookRealm)
	if decoderError != nil {
		http.Error(w, decoderError.Error(), http.StatusUnprocessableEntity)
		return
	}
	userId, ok := r.Context().Value(USER_ID).(string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing " + USER_ID))
	}
	createRealmErr := h.bookRealmService.CreateBookRealm(bookRealm, userId)
	if model.IsExisting(createRealmErr) {
		handleError(createRealmErr, w)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *bookRealmHandler) UpdateBookRealm(w http.ResponseWriter, r *http.Request) {
	var bookRealm model.BookRealmDTO
	bookID := r.PathValue("bookID")
	decoderError := json.NewDecoder(r.Body).Decode(&bookRealm)
	if decoderError != nil {
		http.Error(w, decoderError.Error(), http.StatusUnprocessableEntity)
		return
	}

	if bookID != bookRealm.BookID {
		http.Error(w, "bookID Path param and Payload BookID do not match",
			http.StatusUnprocessableEntity)
		return
	}

	updateRealmErr := h.bookRealmService.UpdateBookRealm(bookRealm, bookID)
	if model.IsExisting(updateRealmErr) {
		handleError(updateRealmErr, w)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *bookRealmHandler) DeleteBookRealm(w http.ResponseWriter, r *http.Request) {
	bookID := r.PathValue("bookID")
	err := h.bookRealmService.DeleteBookRealm(bookID)
	if model.IsExisting(err) {
		handleError(err, w)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *bookRealmHandler) ReadAccountingUsers(w http.ResponseWriter, r *http.Request) {
	queries := r.URL.Query()
	limit := queries.Get("limit")
	searchTerm := queries.Get("searchTerm")
	applicationUsers, err := h.userService.SearchUsers(limit, searchTerm)
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

func (h *bookRealmHandler) CreateUser(w http.ResponseWriter, r *http.Request) {

	var applicationUser model.ApplicationUserDTO

	decoderError := json.NewDecoder(r.Body).Decode(&applicationUser)
	if decoderError != nil {
		http.Error(w, decoderError.Error(), http.StatusUnprocessableEntity)
		return
	}

	createUserError := h.userService.CreateUser(applicationUser)
	if model.IsExisting(createUserError) {
		handleError(createUserError, w)
		return
	}
	w.WriteHeader(http.StatusCreated)

}
