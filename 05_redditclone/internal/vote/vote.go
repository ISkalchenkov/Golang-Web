package vote

type Vote struct {
	ID     uint64
	UserID uint64
	PostID uint64
	Value  int
}

type VoteRepo interface {
	GetByPostID(postID uint64) ([]*Vote, error)
	Upvote(postID, userID uint64) error
	Downvote(postID, userID uint64) error
	Unvote(postID, userID uint64) error
	DeleteAllByPostID(postID uint64) error
}

const (
	Downvote = -1
	Upvote   = 1
)
