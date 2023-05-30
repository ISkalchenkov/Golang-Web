package comment

import (
	"errors"
	"sync"
)

var (
	ErrNoComment = errors.New("no comment found")
)

type CommentMemoryRepository struct {
	lastID uint64
	data   []*Comment
	sync.RWMutex
}

func NewMemoryRepo() *CommentMemoryRepository {
	return &CommentMemoryRepository{
		data: make([]*Comment, 0, 10),
	}
}

func (repo *CommentMemoryRepository) GetByID(id uint64) (*Comment, error) {
	repo.RLock()
	defer repo.RUnlock()
	for _, comment := range repo.data {
		if id == comment.ID {
			return comment, nil
		}
	}
	return nil, ErrNoComment
}

func (repo *CommentMemoryRepository) GetByPostID(id uint64) ([]*Comment, error) {
	repo.RLock()
	defer repo.RUnlock()
	comments := []*Comment{}
	for _, comment := range repo.data {
		if id == comment.PostID {
			comments = append(comments, comment)
		}
	}
	return comments, nil
}

func (repo *CommentMemoryRepository) Add(comment *Comment) (uint64, error) {
	repo.Lock()
	defer repo.Unlock()
	repo.lastID++
	comment.ID = repo.lastID
	repo.data = append(repo.data, comment)
	return repo.lastID, nil
}

func (repo *CommentMemoryRepository) Delete(id uint64) (bool, error) {
	repo.Lock()
	defer repo.Unlock()
	i := -1
	for idx, comment := range repo.data {
		if id != comment.ID {
			continue
		}
		i = idx
	}
	if i < 0 {
		return false, nil
	}

	if i < len(repo.data)-1 {
		copy(repo.data[i:], repo.data[i+1:])
	}
	repo.data[len(repo.data)-1] = nil
	repo.data = repo.data[:len(repo.data)-1]
	return true, nil
}

func (repo *CommentMemoryRepository) DeleteAllByPostID(postID uint64) error {
	repo.Lock()
	defer repo.Unlock()
	undeletedComments := []*Comment{}
	for _, comment := range repo.data {
		if postID != comment.PostID {
			undeletedComments = append(undeletedComments, comment)
		}
	}
	repo.data = undeletedComments
	return nil
}
