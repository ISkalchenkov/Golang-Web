package comment

import "time"

type Comment struct {
	ID       uint64
	AuthorID uint64
	PostID   uint64
	Body     string
	Created  time.Time
}

type CommentRepo interface {
	GetByID(id uint64) (*Comment, error)
	GetByPostID(postID uint64) ([]*Comment, error)
	Add(comment *Comment) (uint64, error)
	Delete(id uint64) (bool, error)
	DeleteAllByPostID(postID uint64) error
}
