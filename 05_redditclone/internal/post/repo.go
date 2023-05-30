package post

import (
	"errors"
	"sync"
)

var (
	ErrNoPost = errors.New("no post found")
)

type PostMemoryRepository struct {
	lastID uint64
	data   []*Post
	sync.RWMutex
}

func NewMemoryRepo() *PostMemoryRepository {
	return &PostMemoryRepository{
		data: make([]*Post, 0, 10),
	}
}

func (repo *PostMemoryRepository) GetAll() ([]*Post, error) {
	repo.RLock()
	defer repo.RUnlock()
	return repo.data, nil
}

func (repo *PostMemoryRepository) GetByCategory(category string) ([]*Post, error) {
	repo.RLock()
	defer repo.RUnlock()
	posts := []*Post{}
	for _, post := range repo.data {
		if category == post.Category {
			posts = append(posts, post)
		}
	}
	return posts, nil
}

func (repo *PostMemoryRepository) GetByAuthorID(authorID uint64) ([]*Post, error) {
	repo.RLock()
	defer repo.RUnlock()
	posts := []*Post{}
	for _, post := range repo.data {
		if authorID == post.AuthorID {
			posts = append(posts, post)
		}
	}
	return posts, nil
}

func (repo *PostMemoryRepository) GetByID(id uint64) (*Post, error) {
	repo.RLock()
	defer repo.RUnlock()
	for _, post := range repo.data {
		if id == post.ID {
			return post, nil
		}
	}
	return nil, ErrNoPost
}

func (repo *PostMemoryRepository) Add(post *Post) (uint64, error) {
	repo.Lock()
	defer repo.Unlock()
	repo.lastID++
	post.ID = repo.lastID
	repo.data = append(repo.data, post)
	return repo.lastID, nil
}

func (repo *PostMemoryRepository) IncrementViews(id uint64) (uint64, error) {
	repo.Lock()
	defer repo.Unlock()
	for _, post := range repo.data {
		if id == post.ID {
			post.Views++
			return post.Views, nil
		}
	}
	return 0, ErrNoPost
}

func (repo *PostMemoryRepository) Delete(id uint64) (bool, error) {
	repo.Lock()
	defer repo.Unlock()
	i := -1
	for idx, post := range repo.data {
		if id != post.ID {
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
