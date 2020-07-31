package main

import (
	"net/http"

	"github.com/toky03/toky-finance-accounting-service/handler"
)

type AccountingHandler interface {
	ReadAccounts(w http.ResponseWriter, r *http.Request)
	ReadBookings(w http.ResponseWriter, r *http.Request)
	SaveBooking(w http.ResponseWriter, r *http.Request)
	SaveAccountOption(w http.ResponseWriter, r *http.Request)
}

type BookHandler interface {
	ReadBookRealms(w http.ResponseWriter, r *http.Request)
	CreateBookRealm(w http.ResponseWriter, r *http.Request)
	UpdateBookRealm(w http.ResponseWriter, r *http.Request)
	ReadAccountingUsers(w http.ResponseWriter, r *http.Request)
	CreateUser(w http.ResponseWriter, r *http.Request)
}

func main() {
	var bookHandler BookHandler

	bookHandler = handler.CreateBookRealmHandler()

	http.HandleFunc("/api/books", bookHandler.ReadBookRealms)
	http.HandleFunc("/api/users", bookHandler.CreateUser)
	http.ListenAndServe(":3000", nil)

}
