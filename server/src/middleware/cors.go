package middleware

import (
	"net/http"
	"os"
)

// This might not needed when we use server-side randering front-end framework
func CORSHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		webSiteAllowed := os.Getenv("CORS_WEBSITE_ALLOWED")
		w.Header().Set("Access-Control-Allow-Origin", webSiteAllowed)
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE, HEAD")
		w.Header().Set("Access-Control-Expose-Headers", "Authorization")
		next.ServeHTTP(w, r)
	})
}
