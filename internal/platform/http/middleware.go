package httpx

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func BaseMiddlewares(r chiLike) {
	// Request ID
	r.Use(middleware.RequestID)

	// Recover from panics
	r.Use(middleware.Recoverer)

	// Timeout (global)
	r.Use(middleware.Timeout(15 * time.Second))

	// Minimal logging
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, req.ProtoMajor)

			next.ServeHTTP(ww, req)

			lat := time.Since(start)
			rid := middleware.GetReqID(req.Context())
			log.Printf("rid=%s method=%s path=%s status=%d bytes=%d latency=%s",
				rid, req.Method, req.URL.Path, ww.Status(), ww.BytesWritten(), lat)
		})
	})
}

// small adapter so we can pass chi.Router or *chi.Mux
type chiLike interface {
	Use(middlewares ...func(http.Handler) http.Handler)
}
