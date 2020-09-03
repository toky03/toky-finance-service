package handler

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	jwtauthhandler "github.com/toky03/jwt-auth-handler"
	"github.com/toky03/toky-finance-accounting-service/model"
	"github.com/toky03/toky-finance-accounting-service/service"
)

type AuthenticationHandlerImpl struct {
	userService    userService
	jwtAuthService jwtauthhandler.JwtHandler
}

func CreateAuthenticationHandler() *AuthenticationHandlerImpl {
	openIDProvider := os.Getenv("OPENID_JWKS_URL")
	if openIDProvider == "" {
		panic("OPENID_JWKS URL must be provided")
	}
	jwtHandler, err := jwtauthhandler.CreateJwtHandler(openIDProvider)
	if err != nil {
		panic("jwt Handler could not have been initialized")
	}
	return &AuthenticationHandlerImpl{
		userService:    service.CreateApplicationUserService(),
		jwtAuthService: jwtHandler,
	}
}

func (h *AuthenticationHandlerImpl) AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		if len(authHeader) != 2 {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Malformed Token"))
			next.ServeHTTP(w, r)
			return
		} else {
			jwtToken := authHeader[1]
			claims := jwt.MapClaims{}
			rsaKeys := h.jwtAuthService.ReadPublicKeys()
			var err error
			for _, rsaKey := range rsaKeys {
				token, errParse := jwt.ParseWithClaims(jwtToken, claims, func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
						return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
					}
					return &rsaKey, nil
				})
				if token.Valid {
					break
				} else {
					err = errParse
				}
			}
			if err != nil {
				log.Printf("Error %v \n", err)
				log.Printf("token %v", jwtToken)
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Invalid Token"))
				next.ServeHTTP(w, r)
				return
			}
			// access := claims["resource_access"]
			userName := claims["preferred_username"]
			applicationUser, userServiceErr := h.userService.FindUserByUsername(userName.(string))
			if model.IsExisting(userServiceErr) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(userServiceErr.ErrorMessage()))
				next.ServeHTTP(w, r)
				return
			}
			context.Set(r, "user-id", applicationUser.UserID)
			next.ServeHTTP(w, r)
		}

	})
}
