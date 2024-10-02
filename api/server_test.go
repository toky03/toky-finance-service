package api

import (
	"net/http/httptest"
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
		{
			name: "Test createBookRealm",
			fields: fields{
				requestType:           "POST",
				requestUrl:            "/api/book",
				excpectedCallsInOrder: []string{"measureRequest", "authenticationMiddleware", "createBookRealm"},
			},
		},
		{
			name: "Test updateBookRealm",
			fields: fields{
				requestType:           "PUT",
				requestUrl:            "/api/book/123",
				excpectedCallsInOrder: []string{"measureRequest", "authenticationMiddleware", "isOwner", "updateBookRealm"},
			},
		},
		{
			name: "Test deleteBookRealm",
			fields: fields{
				requestType:           "DELETE",
				requestUrl:            "/api/book/123",
				excpectedCallsInOrder: []string{"measureRequest", "authenticationMiddleware", "isOwner", "deleteBookRealm"},
			},
		},
		{
			name: "Test readBookRealmById",
			fields: fields{
				requestType:           "GET",
				requestUrl:            "/api/book/123",
				excpectedCallsInOrder: []string{"measureRequest", "authenticationMiddleware", "readBookRealmById"},
			},
		},
		{
			name: "Test readAccounts",
			fields: fields{
				requestType:           "GET",
				requestUrl:            "/api/book/123/account",
				excpectedCallsInOrder: []string{"measureRequest", "authenticationMiddleware", "readAccounts"},
			},
		},
		{
			name: "Test createAccount",
			fields: fields{
				requestType:           "POST",
				requestUrl:            "/api/book/123/account",
				excpectedCallsInOrder: []string{"measureRequest", "authenticationMiddleware", "hasWritePermissions", "createAccount"},
			},
		},
		{
			name: "Test updateAccount",
			fields: fields{
				requestType:           "PUT",
				requestUrl:            "/api/book/123/account/456",
				excpectedCallsInOrder: []string{"measureRequest", "authenticationMiddleware", "hasWritePermissions", "updateAccount"},
			},
		},
		{
			name: "Test deleteAccount",
			fields: fields{
				requestType:           "DELETE",
				requestUrl:            "/api/book/123/account/456",
				excpectedCallsInOrder: []string{"measureRequest", "authenticationMiddleware", "hasWritePermissions", "deleteAccount"},
			},
		},
		{
			name: "Test readAccountOptions",
			fields: fields{
				requestType:           "GET",
				requestUrl:            "/api/book/123/accountOption",
				excpectedCallsInOrder: []string{"measureRequest", "authenticationMiddleware", "readAccountOptions"},
			},
		},
		{
			name: "Test readClosingStatements",
			fields: fields{
				requestType:           "GET",
				requestUrl:            "/api/book/123/closingStatements",
				excpectedCallsInOrder: []string{"measureRequest", "authenticationMiddleware", "readClosingStatements"},
			},
		},
		{
			name: "Test readBookings",
			fields: fields{
				requestType:           "GET",
				requestUrl:            "/api/book/123/booking",
				excpectedCallsInOrder: []string{"measureRequest", "authenticationMiddleware", "readBookings"},
			},
		},
		{
			name: "Test createBooking",
			fields: fields{
				requestType:           "POST",
				requestUrl:            "/api/book/123/booking",
				excpectedCallsInOrder: []string{"measureRequest", "authenticationMiddleware", "hasWritePermissions", "createBooking"},
			},
		},
		{
			name: "Test updateBooking",
			fields: fields{
				requestType:           "PUT",
				requestUrl:            "/api/book/123/booking/789",
				excpectedCallsInOrder: []string{"measureRequest", "authenticationMiddleware", "hasWritePermissions", "updateBooking"},
			},
		},
		{
			name: "Test deleteBooking",
			fields: fields{
				requestType:           "DELETE",
				requestUrl:            "/api/book/123/booking/789",
				excpectedCallsInOrder: []string{"measureRequest", "authenticationMiddleware", "hasWritePermissions", "deleteBooking"},
			},
		},
		{
			name: "Test createUser",
			fields: fields{
				requestType:           "POST",
				requestUrl:            "/api/user",
				excpectedCallsInOrder: []string{"measureRequest", "authenticationMiddleware", "createUser"},
			},
		},
		{
			name: "Test readUser",
			fields: fields{
				requestType:           "GET",
				requestUrl:            "/api/user",
				excpectedCallsInOrder: []string{"measureRequest", "authenticationMiddleware", "readAccountingUsers"},
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
					t.Errorf("%s: expected call %s, but it was not called", tt.name, callName)
					return
				}
				if call.name != callName {
					t.Errorf("%s: expected call %s, but got %s", tt.name, callName, call.name)
				}
				if call.time.Before(lastTime) {
					t.Errorf("%s: expected call %s to be after %v, but it was before", tt.name, call.name, lastTime)
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
