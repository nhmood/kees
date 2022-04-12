package middlewares

import (
	"kees-server/constants"
	"net/http"
)

func AddBaseHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("X-Kees-Version", constants.Version)
		w.Header().Add("X-Kees-Commit", constants.GitCommit)
		next.ServeHTTP(w, r)
	})
}
