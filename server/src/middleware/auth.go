package middleware

import (
	"errors"
	"net/http"
	"playhouse-server/auth"
	"playhouse-server/responsebody"
)

func AuthHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionToken := r.Header.Get("Authorization")
		if sessionToken == "" {
			sessionToken = r.URL.Query().Get("token")
		}

		authenticator := auth.NewSessionAuthenticator()

		isNotValid := !authenticator.IsJWTValid(sessionToken)
		if isNotValid {
			panic(
				responsebody.ResponseErr{
					Code:    http.StatusForbidden,
					ErrBody: errors.New("not a valid token"),
				})
		}

		next.ServeHTTP(w, r)
	})
}
