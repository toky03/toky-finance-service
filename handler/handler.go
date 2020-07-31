package handler

import (
	"net/http"

	"github.com/toky03/toky-finance-accounting-service/model"
	"github.com/toky03/toky-finance-accounting-service/service"
)

// BookRealmService interface to define Contract
type BookRealmService interface {
	FindAllBookRealms() []model.BookRealmDTO
	CreateBookRealm(model.BookRealmDTO, string) error
}

type UserService interface {
	CreateUser(model.ApplicationUserDTO) error
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

}

func (h *BookRealmHandler) CreateBookRealm(w http.ResponseWriter, r *http.Request) {

}

func (h *BookRealmHandler) UpdateBookRealm(w http.ResponseWriter, r *http.Request) {

}

func (h *BookRealmHandler) ReadAccountingUsers(w http.ResponseWriter, r *http.Request) {

}

func (h *BookRealmHandler) CreateUser(w http.ResponseWriter, r *http.Request) {

}
