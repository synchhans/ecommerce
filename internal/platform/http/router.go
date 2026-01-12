package httpx

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Router struct {
	mux *chi.Mux
}

func NewRouter() *Router {
	r := chi.NewRouter()

	// Security headers
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.Header().Set("X-Frame-Options", "DENY")
			next.ServeHTTP(w, req)
		})
	})

	// Standard middlewares
	BaseMiddlewares(r)

	// CORS
	r.Use(SimpleCORS)

	return &Router{mux: r}
}

func (r *Router) Mux() *chi.Mux { return r.mux }
