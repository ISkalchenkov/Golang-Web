package vote

import (
	"sync"
)

type VoteMemoryRepository struct {
	lastID uint64
	data   []*Vote
	sync.RWMutex
}

func NewMemoryRepo() *VoteMemoryRepository {
	return &VoteMemoryRepository{
		data: make([]*Vote, 0, 10),
	}
}

func (repo *VoteMemoryRepository) GetByPostID(postID uint64) ([]*Vote, error) {
	repo.RLock()
	defer repo.RUnlock()
	votes := []*Vote{}
	for _, vote := range repo.data {
		if postID == vote.PostID {
			votes = append(votes, vote)
		}
	}
	return votes, nil
}

func (repo *VoteMemoryRepository) Upvote(postID, userID uint64) error {
	repo.Lock()
	defer repo.Unlock()
	if ok := repo.updateValue(Upvote, postID, userID); ok {
		return nil
	}
	repo.add(Upvote, postID, userID)
	return nil
}

func (repo *VoteMemoryRepository) Downvote(postID, userID uint64) error {
	repo.Lock()
	defer repo.Unlock()
	if ok := repo.updateValue(Downvote, postID, userID); ok {
		return nil
	}
	repo.add(Downvote, postID, userID)
	return nil
}

func (repo *VoteMemoryRepository) Unvote(postID, userID uint64) error {
	repo.Lock()
	defer repo.Unlock()
	repo.delete(postID, userID)
	return nil
}

func (repo *VoteMemoryRepository) DeleteAllByPostID(postID uint64) error {
	repo.Lock()
	defer repo.Unlock()
	undeletedVotes := []*Vote{}
	for _, vote := range repo.data {
		if postID != vote.PostID {
			undeletedVotes = append(undeletedVotes, vote)
		}
	}
	repo.data = undeletedVotes
	return nil
}

func (repo *VoteMemoryRepository) updateValue(value int, postID, userID uint64) bool {
	for _, vote := range repo.data {
		if postID != vote.PostID || userID != vote.UserID {
			continue
		}
		vote.Value = value
		return true
	}
	return false
}

func (repo *VoteMemoryRepository) add(value int, postID, userID uint64) {
	repo.lastID++
	vote := &Vote{
		ID:     repo.lastID,
		UserID: userID,
		PostID: postID,
		Value:  value,
	}
	repo.data = append(repo.data, vote)
}

func (repo *VoteMemoryRepository) delete(postID, userID uint64) {
	i := -1
	for idx, vote := range repo.data {
		if postID != vote.PostID || userID != vote.UserID {
			continue
		}
		i = idx
	}
	if i < 0 {
		return
	}

	if i < len(repo.data)-1 {
		copy(repo.data[i:], repo.data[i+1:])
	}
	repo.data[len(repo.data)-1] = nil
	repo.data = repo.data[:len(repo.data)-1]
}
