package handler

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	jwtauthhandler "github.com/toky03/jwt-auth-handler"
	"github.com/toky03/toky-finance-accounting-service/model"
	"github.com/toky03/toky-finance-accounting-service/service"
)

type accountingService interface {
	ReadBookIdFromAccount(accountId string) (string, model.TokyError)
	ReadBookIdFromBooking(bookingId string) (string, model.TokyError)
}
type authenticationHandlerImpl struct {
	userService       userService
	accountingService accountingService
	jwtAuthService    jwtauthhandler.JwtHandler
}

func CreateAuthenticationHandler() *authenticationHandlerImpl {
	openIDProvider := os.Getenv("OPENID_JWKS_URL")
	if openIDProvider == "" {
		panic("OPENID_JWKS URL must be provided")
	}
	jwtHandler, err := jwtauthhandler.CreateJwtHandler(openIDProvider)
	if err != nil {
		panic("jwt Handler could not have been initialized")
	}
	return &authenticationHandlerImpl{
		userService:       service.CreateApplicationUserService(),
		accountingService: service.CreateAccountingService(),
		jwtAuthService:    jwtHandler,
	}
}

func (h *authenticationHandlerImpl) HasWritePermissions(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := context.Get(r, "user-id")
		if userId == "" {
			return
		}
		var err model.TokyError
		vars := mux.Vars(r)
		bookID := vars["bookID"]
		accountID := vars["accountID"]
		bookingID := vars["bookingID"]
		if accountID != "" && bookingID == "" {
			bookID, err = h.accountingService.ReadBookIdFromAccount(accountID)
			if model.IsExisting(err) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Could not Read Book Id from account Id %s"))
				return
			}

		} else if bookingID != "" {
			bookID, err = h.accountingService.ReadBookIdFromBooking(bookingID)
			if model.IsExisting(err) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Could not Read Book Id from booking Id"))
				return
			}
		}
		isPermitted, err := h.userService.HasWriteAccessFromBook(userId.(string), bookID)
		if model.IsExisting(err) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Could not Read Write Permissions"))
			return
		}
		if isPermitted {
			next.ServeHTTP(w, r)
			return
		} else {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("User is not allowed to write to this Book"))
			next.ServeHTTP(w, r)
			return
		}
	})

}

func (h *authenticationHandlerImpl) IsOwner(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := context.Get(r, "user-id")
		if userId == "" {
			return
		}
		vars := mux.Vars(r)
		bookID := vars["bookID"]
		isPermitted, err := h.userService.IsOwnerOfBook(userId.(string), bookID)
		if model.IsExisting(err) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Could not Read Write Permissions"))
			return
		}
		if isPermitted {
			next.ServeHTTP(w, r)
			return
		} else {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("User is not allowed to modify this Book"))
			next.ServeHTTP(w, r)
			return
		}
	})
}

func (h *authenticationHandlerImpl) AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		if len(authHeader) != 2 {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Malformed Token"))
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
				return
			}
			// access := claims["resource_access"]
			userName := claims["preferred_username"]
			applicationUser, userServiceErr := h.userService.FindUserByUsername(userName.(string))
			if model.IsExisting(userServiceErr) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(userServiceErr.ErrorMessage()))
				return
			}
			context.Set(r, "user-id", applicationUser.UserID)
			next.ServeHTTP(w, r)
		}

	})
}
