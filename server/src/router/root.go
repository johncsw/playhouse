package router

import (
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	ourMiddleware "playhouse-server/middleware"
	"playhouse-server/repository"
)

func NewRootRouter(f *repository.Factory) *chi.Mux {
	r := chi.NewRouter()
	r.Use(chiMiddleware.Logger)
	r.Use(ourMiddleware.ErrorHandler)
	r.Use(ourMiddleware.CORSHandler)
	r.Mount("/session", NewSessionRouter(f))
	return r
}
