package router

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"playhouse-server/auth"
	"playhouse-server/responsebody"
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
		wrapper := responsebody.Wrapper{Writer: w}
		wrapper.Header(map[string]string{
			"Authorization": token,
		}).Status(http.StatusCreated).RawBody(nil)
	}
}
