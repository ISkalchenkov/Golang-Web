package middleware

import (
	"errors"
	"net/http"
	"redditclone/internal/session"
	"redditclone/internal/utils"

	"go.uber.org/zap"
)

type MiddlewareFunc func(http.HandlerFunc) http.HandlerFunc

func NewAuthMiddleware(logger *zap.SugaredLogger, sm session.SessionManager) MiddlewareFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return auth(logger, sm, next)
	}
}

func auth(logger *zap.SugaredLogger, sm session.SessionManager, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sess, err := sm.Check(r)
		if err != nil {
			switch {
			case errors.Is(err, session.ErrBadToken):
				fallthrough
			case errors.Is(err, session.ErrNoAuthToken):
				utils.JSONMessage(w, "unauthorized", http.StatusUnauthorized)
			default:
				logger.Errorf("session check failed: %v", err)
				utils.JSONMessage(w, "session check failed", http.StatusInternalServerError)
			}
			return
		}
		ctx := session.ContextWithSession(r.Context(), sess)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
