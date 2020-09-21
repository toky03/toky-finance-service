package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/toky03/toky-finance-accounting-service/handler"
)

type AccountingHandler interface {
	ReadAccounts(w http.ResponseWriter, r *http.Request)
	ReadAccountOptions(w http.ResponseWriter, r *http.Request)
	ReadBookings(w http.ResponseWriter, r *http.Request)
	CreateBooking(w http.ResponseWriter, r *http.Request)
	CreateAccount(w http.ResponseWriter, r *http.Request)
	SaveAccountOption(w http.ResponseWriter, r *http.Request)
	ReadClosingStatements(w http.ResponseWriter, r *http.Request)
}

type BookHandler interface {
	ReadBookRealms(w http.ResponseWriter, r *http.Request)
	CreateBookRealm(w http.ResponseWriter, r *http.Request)
	UpdateBookRealm(w http.ResponseWriter, r *http.Request)
	ReadAccountingUsers(w http.ResponseWriter, r *http.Request)
	CreateUser(w http.ResponseWriter, r *http.Request)
	ReadBookRealmById(w http.ResponseWriter, r *http.Request)
}

type MonitoringHandler interface {
	MetricsHandler() http.Handler
	MeasureRequest(http.Handler) http.Handler
}

type AuthenticationHandler interface {
	AuthenticationMiddleware(http.Handler) http.Handler
	HasWritePermissions(next http.Handler) http.Handler
}

func main() {
	var bookHandler BookHandler
	var monitoringHandler MonitoringHandler
	var accountingHandler AccountingHandler
	var authenticationHandler AuthenticationHandler

	bookHandler = handler.CreateBookRealmHandler()
	monitoringHandler = handler.CreateMonitoringHandler()
	accountingHandler = handler.CreateAccountingHandler()
	authenticationHandler = handler.CreateAuthenticationHandler()

	err := handler.CreateAndRegisterUserBatchService()
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()

	r.Handle("/metrics", monitoringHandler.MetricsHandler())
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/book", bookHandler.ReadBookRealms).Methods("GET")
	api.HandleFunc("/book", bookHandler.CreateBookRealm).Methods("POST")
	api.HandleFunc("/book/{bookID}", bookHandler.UpdateBookRealm).Methods("PUT")
	api.HandleFunc("/book/{bookID}", bookHandler.ReadBookRealmById).Methods("GET")
	api.HandleFunc("/book/{bookID}/account", accountingHandler.ReadAccounts).Methods("GET")
	api.HandleFunc("/book/{bookID}/accountOption", accountingHandler.ReadAccountOptions).Methods("GET")
	api.HandleFunc("/book/{bookID}/closingStatements", accountingHandler.ReadClosingStatements).Methods("GET")
	api.Handle("/book/{bookID}/account", authenticationHandler.HasWritePermissions(http.HandlerFunc(accountingHandler.CreateAccount))).Methods("POST")
	api.Handle("/book/{bookID}/booking", authenticationHandler.HasWritePermissions(http.HandlerFunc(accountingHandler.CreateBooking))).Methods("POST")
	api.HandleFunc("/user", bookHandler.CreateUser).Methods("POST")
	api.HandleFunc("/user", bookHandler.ReadAccountingUsers).Methods("GET")
	api.Use(authenticationHandler.AuthenticationMiddleware)
	r.Use(monitoringHandler.MeasureRequest)

	log.Fatal(http.ListenAndServe(":3001", r))

}
