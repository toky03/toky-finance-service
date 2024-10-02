package api

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func createCallNamesHandlerMap(
	bookHandler *MockBookHandler,
	monitoringHandler *MockMonitoringHandler,
	accountingHandler *MockAccountingHandler,
	authenticationHandler *MockAuthenticationHandler,
) map[string]Mock {
	return map[string]Mock{
		"readAccounts":             accountingHandler,
		"readAccountOptions":       accountingHandler,
		"readBookings":             accountingHandler,
		"createBooking":            accountingHandler,
		"updateBooking":            accountingHandler,
		"deleteBooking":            accountingHandler,
		"createAccount":            accountingHandler,
		"updateAccount":            accountingHandler,
		"deleteAccount":            accountingHandler,
		"saveAccountOption":        accountingHandler,
		"readClosingStatements":    accountingHandler,
		"readBookRealms":           bookHandler,
		"createBookRealm":          bookHandler,
		"updateBookRealm":          bookHandler,
		"deleteBookRealm":          bookHandler,
		"readAccountingUsers":      bookHandler,
		"createUser":               bookHandler,
		"readBookRealmById":        bookHandler,
		"monitoringHandler":        monitoringHandler,
		"measureRequest":           monitoringHandler,
		"authenticationMiddleware": authenticationHandler,
		"hasWritePermissions":      authenticationHandler,
		"isOwner":                  authenticationHandler,
		"jwksUrl":                  authenticationHandler,
	}
}

func TestServer_Start(t *testing.T) {
	type fields struct {
		requestType           string
		requestUrl            string
		excpectedCallsInOrder []string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "Test metrics",
			fields: fields{
				requestType:           "GET",
				requestUrl:            "/metrics",
				excpectedCallsInOrder: []string{"monitoringHandler"},
			},
		},
		{
			name: "Test jwks url",
			fields: fields{
				requestType:           "GET",
				requestUrl:            "/login-info",
				excpectedCallsInOrder: []string{"jwksUrl"},
			},
		},
		{
			name: "Test readBookRealms",
			fields: fields{
				requestType:           "GET",
				requestUrl:            "/api/book",
				excpectedCallsInOrder: []string{"measureRequest", "authenticationMiddleware", "readBookRealms"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bookHandler := MockBookHandler{}
			accountingHandler := MockAccountingHandler{}
			monitoringHandler := MockMonitoringHandler{}
			authenticationHandler := MockAuthenticationHandler{}

			handlerCallMap := createCallNamesHandlerMap(&bookHandler, &monitoringHandler, &accountingHandler, &authenticationHandler)

			accountingHandler.resetCalls()
			bookHandler.resetCalls()
			monitoringHandler.resetCalls()
			authenticationHandler.resetCalls()

			s := &Server{
				bookHandler:           &bookHandler,
				monitoringHandler:     &monitoringHandler,
				accountingHandler:     &accountingHandler,
				authenticationHandler: &authenticationHandler,
			}
			s.RegisterHandlers()

			w := httptest.NewRecorder()

			lastTime := time.Now()

			request := httptest.NewRequest(tt.fields.requestType, tt.fields.requestUrl, nil)

			s.router.ServeHTTP(w, request)

			for _, callName := range tt.fields.excpectedCallsInOrder {
				call, ok := handlerCallMap[callName].popFirstCall()
				if !ok {
					t.Errorf("expected call %s, but it was not called", callName)
					return
				}
				if call.name != callName {
					t.Errorf("expected call %s, but got %s", callName, call.name)
				}
				if call.time.Before(lastTime) {
					t.Errorf("expected call %s to be after %v, but it was before", call.name, lastTime)
				}
				lastTime = call.time
			}

			// calls should now be empty as every call was checked
			callsBookHandler := bookHandler.readCalls()
			if len(callsBookHandler) > 0 {
				t.Errorf("expected no more calls for bookHandler, but got %v", callsBookHandler)
			}
			callsMonitoringHandler := monitoringHandler.readCalls()
			if len(callsMonitoringHandler) > 0 {
				t.Errorf("expected no more calls for monitoringHandler, but got %v", callsMonitoringHandler)
			}
			callsAccountingHandler := accountingHandler.readCalls()
			if len(callsAccountingHandler) > 0 {
				t.Errorf("expected no more calls for accountingHandler, but got %v", callsAccountingHandler)
			}
			callsAuthenticationHandler := authenticationHandler.readCalls()
			if len(callsAuthenticationHandler) > 0 {
				t.Errorf("expected no more calls for authenticationHandler, but got %v", callsAuthenticationHandler)
			}
		})
	}
}

func TestServer_authMonitoring(t *testing.T) {
	type fields struct {
		bookHandler           BookHandler
		monitoringHandler     MonitoringHandler
		accountingHandler     AccountingHandler
		authenticationHandler AuthenticationHandler
	}
	type args struct {
		handlerFunc           http.HandlerFunc
		additionalMiddlewares []middleware
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   http.Handler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				bookHandler:           tt.fields.bookHandler,
				monitoringHandler:     tt.fields.monitoringHandler,
				accountingHandler:     tt.fields.accountingHandler,
				authenticationHandler: tt.fields.authenticationHandler,
			}
			if got := s.authMonitoring(tt.args.handlerFunc, tt.args.additionalMiddlewares...); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("Server.authMonitoring() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_combineMiddlewares(t *testing.T) {
	type args struct {
		handlerFunc http.HandlerFunc
		middlewares []middleware
	}
	tests := []struct {
		name string
		args args
		want http.Handler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := combineMiddlewares(tt.args.handlerFunc, tt.args.middlewares...); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("combineMiddlewares() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSubrouter(t *testing.T) {
	type args struct {
		router *http.ServeMux
		route  string
	}
	tests := []struct {
		name string
		args args
		want *http.ServeMux
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Subrouter(tt.args.router, tt.args.route); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Subrouter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_removePrefix(t *testing.T) {
	type args struct {
		h      http.Handler
		prefix string
	}
	tests := []struct {
		name string
		args args
		want http.Handler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := removePrefix(tt.args.h, tt.args.prefix); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("removePrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}
