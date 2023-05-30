package middleware

import (
	"fmt"
	"net/http"
	"redditclone/internal/utils"

	"go.uber.org/zap"
)

func Panic(logger *zap.SugaredLogger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				errStr := fmt.Sprint(err)
				logger.Errorw(errStr,
					"method", r.Method,
					"remote_addr", r.RemoteAddr,
					"url", r.URL,
				)
				utils.JSONMessage(w, "internal server error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
