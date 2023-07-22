package router

import (
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	ourMiddleware "playhouse-server/middleware"
)

func NewRootRouter() *chi.Mux {
	r := chi.NewRouter()

	// ErrorHandler must be put at first, as a response would go through middlewares from bottom to top
	r.Use(ourMiddleware.ErrorHandler)
	r.Use(ourMiddleware.CORSHandler)
	r.Use(chiMiddleware.Logger)
	r.Mount("/session", newSessionRouter())
	r.Mount("/upload", newUploadRouter())
	r.Mount("/video", newVideoRouter())
	return r
}
