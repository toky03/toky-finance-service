package api

import (
	"log"
	"net/http"
	"strings"
)

type middleware func(http.Handler) http.Handler

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

	r := http.NewServeMux()

	// TODO bei den Updates und deltes muss noch geprüft werden, ob die Entsprechende Entity auch im richtigen BookRealm ist
	r.Handle("GET /metrics", s.monitoringHandler.MetricsHandler())
	r.HandleFunc("GET /login-info", s.authenticationHandler.JwksUrl)
	api := Subrouter(r, "/api")
	api.Handle("GET /book", s.authMonitoring(s.bookHandler.ReadBookRealms))
	api.Handle("POST /book", s.authMonitoring(s.bookHandler.CreateBookRealm))
	api.Handle("PUT /book/{bookID}", s.authMonitoring(http.HandlerFunc(s.bookHandler.UpdateBookRealm), s.authenticationHandler.IsOwner))
	api.Handle("DELTE /book/{bookID}", s.authMonitoring(http.HandlerFunc(s.bookHandler.DeleteBookRealm), s.authenticationHandler.IsOwner))
	api.Handle("GET /book/{bookID}", s.authMonitoring(s.bookHandler.ReadBookRealmById))
	api.Handle("GET /book/{bookID}/account", s.authMonitoring(s.accountingHandler.ReadAccounts))
	api.Handle("POST /book/{bookID}/account", s.authMonitoring(http.HandlerFunc(s.accountingHandler.CreateAccount), s.authenticationHandler.HasWritePermissions))
	api.Handle("PUT /book/{bookID}/account/{accountID}", s.authMonitoring(http.HandlerFunc(s.accountingHandler.UpdateAccount), s.authenticationHandler.HasWritePermissions))
	api.Handle("DELETE /book/{bookID}/account/{accountID}", s.authMonitoring(http.HandlerFunc(s.accountingHandler.DeleteAccount), s.authenticationHandler.HasWritePermissions))
	api.Handle("GET /book/{bookID}/accountOption", s.authMonitoring(s.accountingHandler.ReadAccountOptions))
	api.Handle("GET /book/{bookID}/closingStatements", s.authMonitoring(s.accountingHandler.ReadClosingStatements))
	api.Handle("GET /book/{bookID}/booking", s.authMonitoring(s.accountingHandler.ReadBookings))
	api.Handle("POST /book/{bookID}/booking", s.authMonitoring(http.HandlerFunc(s.accountingHandler.CreateBooking), s.authenticationHandler.HasWritePermissions))
	api.Handle("PUT /book/{bookID}/booking/{bookingID}", s.authMonitoring(http.HandlerFunc(s.accountingHandler.UpdateBooking), s.authenticationHandler.HasWritePermissions))
	api.Handle("DELETE /book/{bookID}/booking/{bookingID}", s.authMonitoring(http.HandlerFunc(s.accountingHandler.DeleteBooking), s.authenticationHandler.HasWritePermissions))
	api.Handle("POST /user", s.authMonitoring(s.bookHandler.CreateUser))
	api.Handle("GET /user", s.authMonitoring(s.bookHandler.ReadAccountingUsers))

	log.Fatal(http.ListenAndServe(":3001", r))
}

func (s *Server) authMonitoring(handlerFunc http.HandlerFunc, additionalMiddlewares ...middleware) http.Handler {
	middlewares := append(additionalMiddlewares, s.authenticationHandler.AuthenticationMiddleware, s.monitoringHandler.MeasureRequest)
	return combineMiddlewares(handlerFunc, middlewares...)
}

func combineMiddlewares(handlerFunc http.HandlerFunc, middlewares ...middleware) http.Handler {

	handler := http.Handler(handlerFunc)

	if len(middlewares) < 1 {
		return handler
	}

	wrappedHandler := handler
	for _, mw := range middlewares {
		wrappedHandler = mw(wrappedHandler)
	}
	return wrappedHandler
}

func Subrouter(router *http.ServeMux, route string) *http.ServeMux {
	sr := http.NewServeMux()
	route = strings.TrimSuffix(route, "/")
	router.Handle(route, removePrefix(sr, route))
	router.Handle(route+"/", removePrefix(sr, route))
	return sr
}

func removePrefix(h http.Handler, prefix string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		r.URL.Path = "/" + strings.TrimPrefix(strings.TrimPrefix(path, prefix), "/")
		h.ServeHTTP(w, r)
		r.URL.Path = path
	})
}
