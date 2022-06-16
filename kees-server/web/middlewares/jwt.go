package middlewares

import (
	"context"
	"net/http"

	"kees/server/helpers"
)

func ValidateJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-Kees-JWT-Token")
		helpers.Debug(token)

		jwt, err := helpers.ValidateJWT(token)
		if err != nil {
			helpers.Halt(w, http.StatusBadRequest, "Invalid JWT", nil)
			return
		}
		helpers.Debug(jwt)

		// Ugly hack - can't seem to figure out how to get jwt.MapClaims
		// to turn into a proper map[string]interface{} after passing through
		// middleware (either panic (1st) or compile error (2nd)
		// -> interface conversion: interface {} is jwt.MapClaims, not map[string]interface {}
		// -> invalid operation: jwt["username"] (type interface {} does not support indexing)
		// Resort to converting from jwt.MapClaims -> JSON string -> map[string]interface
		data := helpers.ToInterface(jwt)

		// Store the JWT payload in the request context and continue
		ctxWithUser := context.WithValue(r.Context(), "jwt", data)
		r = r.WithContext(ctxWithUser)

		next.ServeHTTP(w, r)

	})
}
