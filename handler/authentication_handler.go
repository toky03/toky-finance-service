package handler

import (
	"net/http"

	"github.com/gorilla/context"
)

func AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context.Set(r, "user-id", "1")
		next.ServeHTTP(w, r)
	})
}
