package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	jwtauthhandler "github.com/toky03/jwt-auth-handler"
	"github.com/toky03/toky-finance-accounting-service/model"
)

type customUserIdKey string

const (
	USER_ID customUserIdKey = "x-user-id"
)

type accountingService interface {
	ReadBookIdFromAccount(accountId string) (string, model.TokyError)
	ReadBookIdFromBooking(bookingId string) (string, model.TokyError)
}

type authenticationHandlerImpl struct {
	userService        userService
	accountingService  accountingService
	jwtAuthService     jwtauthhandler.JwtHandler
	openIDBaseUrl      string
	openIDClientSecret string
	openIDClientID     string
}

func CreateAuthenticationHandler(accountingService accountingService, userService userService) *authenticationHandlerImpl {
	openIDProvider := os.Getenv("OPENID_JWKS_URL")
	if openIDProvider == "" {
		panic("OPENID_JWKS URL must be provided")
	}
	var externalOpenIDProvider string
	if os.Getenv("OPENID_JWKS_EXTERNAL_URL") != "" {
		externalOpenIDProvider = os.Getenv("OPENID_JWKS_EXTERNAL_URL")
	} else {
		externalOpenIDProvider = openIDProvider
	}
	openIDClientSecret := os.Getenv("ID_PROVIDER_CLIENT_SECRET")
	if openIDClientSecret == "" {
		panic("ID_PROVIDER_CLIENT_SECRET URL must be provided")
	}
	openIDClientId := os.Getenv("ID_PROVIDER_CLIENT_ID")
	if openIDClientId == "" {
		panic("ID_PROVIDER_CLIENT_ID URL must be provided")
	}

	jwtHandler, err := jwtauthhandler.CreateJwtHandler(openIDProvider + "/certs")
	if err != nil {
		panic("jwt Handler could not have been initialized")
	}
	return &authenticationHandlerImpl{
		userService:        userService,
		accountingService:  accountingService,
		jwtAuthService:     jwtHandler,
		openIDBaseUrl:      externalOpenIDProvider,
		openIDClientSecret: openIDClientSecret,
		openIDClientID:     openIDClientId,
	}
}

func (h *authenticationHandlerImpl) JwksUrl(w http.ResponseWriter, r *http.Request) {

	loginInformation := model.LoginInformationDto{
		AuthUrl:      h.openIDBaseUrl + "/token",
		ClientSecret: h.openIDClientSecret,
		ClientId:     h.openIDClientID,
	}

	js, marshalError := json.Marshal(loginInformation)
	if marshalError != nil {
		http.Error(w, marshalError.Error(), http.StatusUnprocessableEntity)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (h *authenticationHandlerImpl) HasWritePermissions(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId, ok := r.Context().Value(USER_ID).(string)
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Missing " + string(USER_ID)))
		}
		var err model.TokyError
		bookID := r.PathValue("bookID")
		accountID := r.PathValue("accountID")
		bookingID := r.PathValue("bookingID")
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
		isPermitted, err := h.userService.HasWriteAccessFromBook(userId, bookID)
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
		userId, ok := r.Context().Value(USER_ID).(string)
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Missing " + string(USER_ID)))
		}
		bookID := r.PathValue("bookID")
		isPermitted, err := h.userService.IsOwnerOfBook(userId, bookID)
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
						fmt.Printf("signing method not expected\n")
						return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
					}
					return &rsaKey, nil
				})
				if token.Valid {
					err = nil
					break
				} else {
					err = errParse
				}
			}
			if err != nil {
				log.Printf("Error %v \n", err)
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
			ctx := context.WithValue(r.Context(), USER_ID, applicationUser.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		}

	})
}
