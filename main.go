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
	UpdateBooking(http.ResponseWriter, *http.Request)
	DeleteBooking(http.ResponseWriter, *http.Request)
	CreateAccount(w http.ResponseWriter, r *http.Request)
	UpdateAccount(w http.ResponseWriter, r *http.Request)
	DeleteAccount(http.ResponseWriter, *http.Request)
	SaveAccountOption(w http.ResponseWriter, r *http.Request)
	ReadClosingStatements(w http.ResponseWriter, r *http.Request)
}

type BookHandler interface {
	ReadBookRealms(w http.ResponseWriter, r *http.Request)
	CreateBookRealm(w http.ResponseWriter, r *http.Request)
	UpdateBookRealm(w http.ResponseWriter, r *http.Request)
	DeleteBookRealm(w http.ResponseWriter, r *http.Request)
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
	IsOwner(next http.Handler) http.Handler
	JwksUrl(w http.ResponseWriter, r *http.Request)
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

	// TODO bei den Updates und deltes muss noch geprüft werden, ob die Entsprechende Entity auch im richtigen BookRealm ist
	r.Handle("/metrics", monitoringHandler.MetricsHandler())
	r.HandleFunc("/login-info", authenticationHandler.JwksUrl).Methods("GET")
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/book", bookHandler.ReadBookRealms).Methods("GET")
	api.HandleFunc("/book", bookHandler.CreateBookRealm).Methods("POST")
	api.Handle("/book/{bookID}", authenticationHandler.IsOwner(http.HandlerFunc(bookHandler.UpdateBookRealm))).Methods("PUT")
	api.Handle("/book/{bookID}", authenticationHandler.IsOwner(http.HandlerFunc(bookHandler.DeleteBookRealm))).Methods("DELETE")
	api.HandleFunc("/book/{bookID}", bookHandler.ReadBookRealmById).Methods("GET")
	api.HandleFunc("/book/{bookID}/account", accountingHandler.ReadAccounts).Methods("GET")
	api.Handle("/book/{bookID}/account", authenticationHandler.HasWritePermissions(http.HandlerFunc(accountingHandler.CreateAccount))).Methods("POST")
	api.Handle("/book/{bookID}/account/{accountID}", authenticationHandler.HasWritePermissions(http.HandlerFunc(accountingHandler.UpdateAccount))).Methods("PUT")
	api.Handle("/book/{bookID}/account/{accountID}", authenticationHandler.HasWritePermissions(http.HandlerFunc(accountingHandler.DeleteAccount))).Methods("DELETE")
	api.HandleFunc("/book/{bookID}/accountOption", accountingHandler.ReadAccountOptions).Methods("GET")
	api.HandleFunc("/book/{bookID}/closingStatements", accountingHandler.ReadClosingStatements).Methods("GET")
	api.HandleFunc("/book/{bookID}/booking", accountingHandler.ReadBookings).Methods("GET")
	api.Handle("/book/{bookID}/booking", authenticationHandler.HasWritePermissions(http.HandlerFunc(accountingHandler.CreateBooking))).Methods("POST")
	api.Handle("/book/{bookID}/booking/{bookingID}", authenticationHandler.HasWritePermissions(http.HandlerFunc(accountingHandler.UpdateBooking))).Methods("PUT")
	api.Handle("/book/{bookID}/booking/{bookingID}", authenticationHandler.HasWritePermissions(http.HandlerFunc(accountingHandler.DeleteBooking))).Methods("DELETE")

	api.HandleFunc("/user", bookHandler.CreateUser).Methods("POST")
	api.HandleFunc("/user", bookHandler.ReadAccountingUsers).Methods("GET")

	api.Use(authenticationHandler.AuthenticationMiddleware)

	//r.Use(monitoringHandler.MeasureRequest)

	log.Fatal(http.ListenAndServe(":3001", r))

}
