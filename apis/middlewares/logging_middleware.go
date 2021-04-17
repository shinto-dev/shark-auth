package middlewares

import (
	"net/http"

	"shark-auth/foundation/logging"
)

func LoggingMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := logging.NewContext(r.Context())
		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}
