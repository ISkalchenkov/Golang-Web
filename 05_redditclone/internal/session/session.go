package session

import (
	"context"
	"errors"
	"net/http"
	"redditclone/internal/user"
)

type Session struct {
	Token  string
	UserID uint64
}

type SessionManager interface {
	Check(r *http.Request) (*Session, error)
	Create(user *user.User) (*Session, error)
}

var (
	ErrNoAuth      = errors.New("no session found")
	ErrNoAuthToken = errors.New("no token in authorization header")
	ErrBadToken    = errors.New("bad token")
)

type sessKey string

var SessionKey sessKey = "sessionKey"

func ContextWithSession(ctx context.Context, sess *Session) context.Context {
	return context.WithValue(ctx, SessionKey, sess)
}

func SessionFromContext(ctx context.Context) (*Session, error) {
	sess, ok := ctx.Value(SessionKey).(*Session)
	if !ok || sess == nil {
		return nil, ErrNoAuth
	}
	return sess, nil
}
