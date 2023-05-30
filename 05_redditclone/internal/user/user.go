package user

type User struct {
	ID       uint64
	Username string
	password string
	salt     string
}

type UserRepo interface {
	Authorize(username, password string) (*User, error)
	Registrate(username, password string) (*User, error)
	GetByID(id uint64) (*User, error)
	GetByUsername(username string) (*User, error)
}
