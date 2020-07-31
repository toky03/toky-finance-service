package main

import (
	"github.com/gorilla/mux"
	"github.com/toky03/toky-finance-accounting-service/handler"
	"log"
	"net/http"
)

type AccountingHandler interface {
	ReadAccounts(w http.ResponseWriter, r *http.Request)
	ReadBookings(w http.ResponseWriter, r *http.Request)
	CreateBooking(w http.ResponseWriter, r *http.Request)
	SaveAccountOption(w http.ResponseWriter, r *http.Request)
}

type BookHandler interface {
	ReadBookRealms(w http.ResponseWriter, r *http.Request)
	CreateBookRealm(w http.ResponseWriter, r *http.Request)
	UpdateBookRealm(w http.ResponseWriter, r *http.Request)
	ReadAccountingUsers(w http.ResponseWriter, r *http.Request)
	CreateUser(w http.ResponseWriter, r *http.Request)
}

type MonitoringHandler interface {
	MetricsHandler() http.Handler
	MeasureRequest(http.Handler) http.Handler
}

func main() {
	var bookHandler BookHandler
	var monitoringHandler MonitoringHandler
	var accountingHandler AccountingHandler

	bookHandler = handler.CreateBookRealmHandler()
	monitoringHandler = handler.CreateMonitoringHandler()
	accountingHandler = handler.CreateAccountingHandler()

	r := mux.NewRouter()
	r.Use(monitoringHandler.MeasureRequest)
	r.Handle("/metrics", monitoringHandler.MetricsHandler())
	r.HandleFunc("/api/book", bookHandler.ReadBookRealms).Methods("GET")
	r.HandleFunc("/api/book", bookHandler.CreateBookRealm).Methods("POST")
	r.HandleFunc("/api/book/{bookID}", bookHandler.UpdateBookRealm).Methods("PUT")
	r.HandleFunc("/api/book/{bookID}/account", accountingHandler.ReadAccounts).Methods("GET")
	r.HandleFunc("/api/book/{bookID}/account/{accountID}", accountingHandler.CreateBooking).Methods("POST")
	r.HandleFunc("/api/user", bookHandler.CreateUser).Methods("POST")
	r.HandleFunc("/api/user", bookHandler.ReadAccountingUsers).Methods("GET")
	log.Fatal(http.ListenAndServe(":3000", r))

}
