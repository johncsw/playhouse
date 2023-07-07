package middleware

import (
	"errors"
	"net/http"
	"playhouse-server/auth"
	"playhouse-server/util"
)

func AuthHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionToken := r.Header.Get("Authorization")

		authenticator := auth.NewSessionAuthenticator()

		isNotValid := !authenticator.IsJWTValid(sessionToken)
		if isNotValid {
			panic(
				util.ResponseErr{
					Code:    http.StatusForbidden,
					ErrBody: errors.New("not a valid token"),
				})
		}

		next.ServeHTTP(w, r)
	})
}
