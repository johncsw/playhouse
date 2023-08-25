package router

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"playhouse-server/auth"
	"playhouse-server/response"
)

func newSessionRouter() *chi.Mux {

	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		r.Post("/", newSessionHandler())
	})
	return r
}

func newSessionHandler() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		token := auth.CreateSessionToken()
		builder := response.Builder{Writer: w}
		builder.Header(map[string]string{
			"Authorization": token,
		}).Status(http.StatusCreated).BuildWithBytes(nil)
	}
}
