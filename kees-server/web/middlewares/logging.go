package middlewares

import (
	"github.com/gorilla/handlers"
	"net/http"
	"os"
)

func AddLogging(next http.Handler) http.Handler {
	return handlers.CombinedLoggingHandler(os.Stdout, next)
}
