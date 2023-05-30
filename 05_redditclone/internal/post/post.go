package post

import "time"

type Post struct {
	ID       uint64
	AuthorID uint64
	Views    uint64

	Title    string
	Type     string
	Text     string
	URL      string
	Category string

	Created time.Time
}

type PostRepo interface {
	GetAll() ([]*Post, error)
	GetByCategory(category string) ([]*Post, error)
	GetByAuthorID(authorID uint64) ([]*Post, error)
	GetByID(id uint64) (*Post, error)
	Add(post *Post) (uint64, error)
	IncrementViews(id uint64) (uint64, error)
	Delete(id uint64) (bool, error)
}
