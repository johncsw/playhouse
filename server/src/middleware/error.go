package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"playhouse-server/response"
)

func ErrorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				var responseErr response.Error
				switch e := err.(type) {
				case response.Error:
					responseErr = e
				default:
					errMessage := fmt.Sprintf("Unknown Error: %v", err)
					responseErr = response.Error{
						Code:  http.StatusInternalServerError,
						Cause: errors.New(errMessage),
					}
				}

				log.Print(responseErr.Cause.Error())

				w.WriteHeader(responseErr.Code)
				err := json.NewEncoder(w).Encode(map[string]string{"error": responseErr.Error()})
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			}
		}()
		next.ServeHTTP(w, r)
	})
}
