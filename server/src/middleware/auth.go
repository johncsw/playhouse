package middleware

import (
	"context"
	"errors"
	"net/http"
	"playhouse-server/auth"
	"playhouse-server/response"
)

func AuthHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionToken := r.Header.Get("Authorization")
		if sessionToken == "" {
			sessionToken = r.URL.Query().Get("token")
		}

		isValid, sessionID := auth.IsSessionTokenValid(sessionToken)
		if !isValid {
			panic(
				response.Error{
					Code:  http.StatusForbidden,
					Cause: errors.New("not a valid token"),
				})
		}

		ctx := context.WithValue(r.Context(), "sessionID", sessionID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
