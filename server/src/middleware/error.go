package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

type ResponseErr struct {
	Code    int
	ErrBody error
}

func (h ResponseErr) Error() string {
	return h.ErrBody.Error()
}

func ErrorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				var responseErr ResponseErr
				switch e := err.(type) {
				case ResponseErr:
					responseErr = e
				default:
					errMessage := fmt.Sprintf("Unknown Error: %v", err)
					responseErr = ResponseErr{
						Code:    http.StatusInternalServerError,
						ErrBody: errors.New(errMessage),
					}
					log.Fatalf(errMessage)
				}

				w.WriteHeader(responseErr.Code)
				err := json.NewEncoder(w).Encode(map[string]string{"error": responseErr.Error()})
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					log.Fatalf(err.Error())
				}
			}
		}()
		next.ServeHTTP(w, r)
	})
}
