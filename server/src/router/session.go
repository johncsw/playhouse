package router

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"playhouse-server/repository"
	"playhouse-server/util"
)

func NewSessionRouter(f *repository.Factory) *chi.Mux {
	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		r.Post("/", newSessionHandler(f))
	})
	return r
}

func newSessionHandler(f *repository.Factory) http.HandlerFunc {
	repo := f.NewSessionRepo()

	return func(w http.ResponseWriter, r *http.Request) {
		s := repo.NewSession()
		token := util.GenJWT(s.ID, s.DueAt)
		headers := map[string]string{
			"Authorization": token,
		}
		util.ReturnSuccess(nil, headers, w)
	}
}
