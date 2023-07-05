package router

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"playhouse-server/auth"
	"playhouse-server/util"
)

func newSessionRouter() *chi.Mux {
	authenticator := auth.NewSessionAuthenticator()
	wrapper := util.NewWrapper()

	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		r.Post("/", newSessionHandler(authenticator, wrapper))
	})
	return r
}

func newSessionHandler(authenticator *auth.SessionAuthenticator, wrapper *util.Wrapper) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		token := authenticator.InitializeSession()
		headers := map[string]string{
			"Authorization": token,
		}
		wrapper.SuccessfulResponse(nil, headers, w)
	}
}
