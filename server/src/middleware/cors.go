package middleware

import (
	"github.com/go-chi/cors"
	"net/http"
	"playhouse-server/util"
)

// This might not needed when we use server-side randering front-end framework
func CORSHandler(next http.Handler) http.Handler {
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{util.Env{}.CORS_ALLOWED_WEBSITE()}, // Use your client's url here
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})

	return cors.Handler(next)
}
