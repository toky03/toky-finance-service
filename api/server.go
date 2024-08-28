package api

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
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

type Server struct {
	bookHandler           BookHandler
	monitoringHandler     MonitoringHandler
	accountingHandler     AccountingHandler
	authenticationHandler AuthenticationHandler
}

func CreateServer(bookHandler BookHandler, monitoringHandler MonitoringHandler, accountingHandler AccountingHandler, authenticationHandler AuthenticationHandler) *Server {

	return &Server{
		bookHandler:           bookHandler,
		monitoringHandler:     monitoringHandler,
		accountingHandler:     accountingHandler,
		authenticationHandler: authenticationHandler,
	}
}

func (s *Server) Start() {

	r := mux.NewRouter()

	// TODO bei den Updates und deltes muss noch geprüft werden, ob die Entsprechende Entity auch im richtigen BookRealm ist
	r.Handle("/metrics", s.monitoringHandler.MetricsHandler())
	r.HandleFunc("/login-info", s.authenticationHandler.JwksUrl).Methods("GET")
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/book", s.bookHandler.ReadBookRealms).Methods("GET")
	api.HandleFunc("/book", s.bookHandler.CreateBookRealm).Methods("POST")
	api.Handle("/book/{bookID}", s.authenticationHandler.IsOwner(http.HandlerFunc(s.bookHandler.UpdateBookRealm))).Methods("PUT")
	api.Handle("/book/{bookID}", s.authenticationHandler.IsOwner(http.HandlerFunc(s.bookHandler.DeleteBookRealm))).Methods("DELETE")
	api.HandleFunc("/book/{bookID}", s.bookHandler.ReadBookRealmById).Methods("GET")
	api.HandleFunc("/book/{bookID}/account", s.accountingHandler.ReadAccounts).Methods("GET")
	api.Handle("/book/{bookID}/account", s.authenticationHandler.HasWritePermissions(http.HandlerFunc(s.accountingHandler.CreateAccount))).Methods("POST")
	api.Handle("/book/{bookID}/account/{accountID}", s.authenticationHandler.HasWritePermissions(http.HandlerFunc(s.accountingHandler.UpdateAccount))).Methods("PUT")
	api.Handle("/book/{bookID}/account/{accountID}", s.authenticationHandler.HasWritePermissions(http.HandlerFunc(s.accountingHandler.DeleteAccount))).Methods("DELETE")
	api.HandleFunc("/book/{bookID}/accountOption", s.accountingHandler.ReadAccountOptions).Methods("GET")
	api.HandleFunc("/book/{bookID}/closingStatements", s.accountingHandler.ReadClosingStatements).Methods("GET")
	api.HandleFunc("/book/{bookID}/booking", s.accountingHandler.ReadBookings).Methods("GET")
	api.Handle("/book/{bookID}/booking", s.authenticationHandler.HasWritePermissions(http.HandlerFunc(s.accountingHandler.CreateBooking))).Methods("POST")
	api.Handle("/book/{bookID}/booking/{bookingID}", s.authenticationHandler.HasWritePermissions(http.HandlerFunc(s.accountingHandler.UpdateBooking))).Methods("PUT")
	api.Handle("/book/{bookID}/booking/{bookingID}", s.authenticationHandler.HasWritePermissions(http.HandlerFunc(s.accountingHandler.DeleteBooking))).Methods("DELETE")

	api.HandleFunc("/user", s.bookHandler.CreateUser).Methods("POST")
	api.HandleFunc("/user", s.bookHandler.ReadAccountingUsers).Methods("GET")

	api.Use(s.authenticationHandler.AuthenticationMiddleware)

	r.Use(s.monitoringHandler.MeasureRequest)

	log.Fatal(http.ListenAndServe(":3001", r))
}
