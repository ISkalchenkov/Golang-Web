package session

import (
	"fmt"
	"net/http"
	"redditclone/internal/user"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

type JWTSessionManager struct {
	secret   []byte
	tokenTTL time.Duration
}

func NewJWTSessionManager(secret string, tokenTTL time.Duration) *JWTSessionManager {
	return &JWTSessionManager{
		secret:   []byte(secret),
		tokenTTL: tokenTTL,
	}
}

func (sm *JWTSessionManager) Check(r *http.Request) (*Session, error) {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return nil, ErrNoAuthToken
	}
	tokenString := strings.TrimPrefix(auth, "Bearer ")

	hashSecretGetter := func(token *jwt.Token) (interface{}, error) {
		method, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok || method.Alg() != "HS256" {
			return nil, fmt.Errorf("bad sign method")
		}
		return sm.secret, nil
	}

	token, err := jwt.Parse(tokenString, hashSecretGetter)
	if err != nil || !token.Valid {
		return nil, ErrBadToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrBadToken
	}
	userData, ok := claims["user"].(map[string]interface{})
	if !ok {
		return nil, ErrBadToken
	}
	userID, ok := userData["id"].(float64)
	if !ok {
		return nil, ErrBadToken
	}
	sess := &Session{
		Token:  tokenString,
		UserID: uint64(userID),
	}
	return sess, nil
}

func (sm *JWTSessionManager) Create(user *user.User) (*Session, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": map[string]interface{}{
			"username": user.Username,
			"id":       user.ID,
		},
		"iat": now.Unix(),
		"exp": now.Add(sm.tokenTTL).Unix(),
	})
	tokenString, err := token.SignedString(sm.secret)
	if err != nil {
		return nil, fmt.Errorf("token signing failed: %w", err)
	}
	sess := &Session{
		Token:  tokenString,
		UserID: user.ID,
	}
	return sess, nil
}
