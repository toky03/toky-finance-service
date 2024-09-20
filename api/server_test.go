package api

import (
	"net/http"
	"reflect"
	"testing"
)

type MockBookHandler struct {
}

func TestCreateServer(t *testing.T) {
	type args struct {
		bookHandler           BookHandler
		monitoringHandler     MonitoringHandler
		accountingHandler     AccountingHandler
		authenticationHandler AuthenticationHandler
	}
	tests := []struct {
		name string
		args args
		want *Server
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateServer(tt.args.bookHandler, tt.args.monitoringHandler, tt.args.accountingHandler, tt.args.authenticationHandler); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateServer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_Start(t *testing.T) {
	type fields struct {
		bookHandler           BookHandler
		monitoringHandler     MonitoringHandler
		accountingHandler     AccountingHandler
		authenticationHandler AuthenticationHandler
	}
	tests := []struct {
		name   string
		fields fields
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
			s.Start()
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
			if got := s.authMonitoring(tt.args.handlerFunc, tt.args.additionalMiddlewares...); !reflect.DeepEqual(got, tt.want) {
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
			if got := combineMiddlewares(tt.args.handlerFunc, tt.args.middlewares...); !reflect.DeepEqual(got, tt.want) {
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
