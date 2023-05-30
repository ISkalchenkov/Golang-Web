package user

import (
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"sync"
)

var (
	ErrNoUser        = errors.New("no user found")
	ErrBadPass       = errors.New("invalid password")
	ErrUsernameTaken = errors.New("username is already taken")
)

type UserMemoryRepository struct {
	data   map[string]*User
	lastID uint64
	sync.RWMutex
}

func NewMemoryRepo() *UserMemoryRepository {
	return &UserMemoryRepository{
		data: make(map[string]*User),
	}
}

func (repo *UserMemoryRepository) GetByID(id uint64) (*User, error) {
	repo.RLock()
	defer repo.RUnlock()
	for _, u := range repo.data {
		if id == u.ID {
			return u, nil
		}
	}
	return nil, ErrNoUser
}

func (repo *UserMemoryRepository) GetByUsername(username string) (*User, error) {
	repo.RLock()
	defer repo.RUnlock()
	for _, u := range repo.data {
		if username == u.Username {
			return u, nil
		}
	}
	return nil, ErrNoUser
}

func (repo *UserMemoryRepository) Authorize(username, password string) (*User, error) {
	repo.RLock()
	defer repo.RUnlock()
	u, ok := repo.data[username]
	if !ok {
		return nil, ErrNoUser
	}

	if hashPassword(password, []byte(u.salt)) != u.password {
		return nil, ErrBadPass
	}

	return u, nil
}

func (repo *UserMemoryRepository) Registrate(username, password string) (*User, error) {
	repo.Lock()
	defer repo.Unlock()
	_, ok := repo.data[username]
	if ok {
		return nil, ErrUsernameTaken
	}
	repo.lastID++
	salt, err := generateSalt(32)
	if err != nil {
		return nil, fmt.Errorf("registrate failed: %w", err)
	}
	u := &User{
		ID:       repo.lastID,
		Username: username,
		password: hashPassword(password, salt),
		salt:     string(salt),
	}
	repo.data[username] = u
	return u, nil
}

func generateSalt(size int) ([]byte, error) {
	salt := make([]byte, size)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, fmt.Errorf("generate salt failed: %w", err)
	}
	return salt, nil
}

func hashPassword(password string, salt []byte) string {
	passwordBytes := []byte(password)
	passwordBytes = append(passwordBytes, salt...)
	hashedPasswordBytes := sha256.Sum256(passwordBytes)
	return fmt.Sprintf("%x", hashedPasswordBytes)
}
