package middlewares

import (
	"net/http"

	"kees/server/constants"
)

func AddBaseHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-Kees-Version", constants.Version)
		w.Header().Add("X-Kees-Commit", constants.GitCommit)
		next.ServeHTTP(w, r)
	})
}

func AddJSONHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		next.ServeHTTP(w, r)
	})
}
