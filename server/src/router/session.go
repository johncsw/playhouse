package router

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"playhouse-server/auth"
	"playhouse-server/responsebody"
)

func newSessionRouter() *chi.Mux {
	authenticator := auth.NewSessionAuthenticator()

	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		r.Post("/", newSessionHandler(authenticator))
	})
	return r
}

func newSessionHandler(authenticator *auth.SessionAuthenticator) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		token := authenticator.InitializeSession()
		wrapper := responsebody.Wrapper{Writer: w}
		wrapper.Header(map[string]string{
			"Authorization": token,
		}).Status(http.StatusCreated).RawBody(nil)
	}
}
