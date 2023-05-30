package note

import (
	"sync"

	"github.com/pkg/errors"
)

var (
	ErrNoNote = errors.New("note not found")
)

type NoteMemoryRepository struct {
	data   []*Note
	lastID uint64
	sync.RWMutex
}

func NewMemoryRepo() *NoteMemoryRepository {
	return &NoteMemoryRepository{
		data: make([]*Note, 0, 10),
	}
}

func (repo *NoteMemoryRepository) GetAll() ([]*Note, error) {
	repo.RLock()
	defer repo.RUnlock()
	return repo.data, nil
}

func (repo *NoteMemoryRepository) GetByID(id uint64) (*Note, error) {
	repo.RLock()
	defer repo.RUnlock()
	for _, note := range repo.data {
		if id == note.ID {
			return note, nil
		}
	}
	return nil, ErrNoNote
}

func (repo *NoteMemoryRepository) Add(note *Note) (uint64, error) {
	repo.Lock()
	defer repo.Unlock()
	repo.lastID++
	note.ID = repo.lastID
	repo.data = append(repo.data, note)
	return repo.lastID, nil
}

func (repo *NoteMemoryRepository) Update(newNote *Note) (*Note, error) {
	repo.Lock()
	defer repo.Unlock()
	for _, note := range repo.data {
		if newNote.ID != note.ID {
			continue
		}
		note.Text = newNote.Text
		note.UpdatedAt = newNote.UpdatedAt
		return note, nil
	}
	return nil, ErrNoNote
}

func (repo *NoteMemoryRepository) Delete(id uint64) error {
	repo.RLock()
	defer repo.RUnlock()
	i := -1
	for idx, note := range repo.data {
		if id != note.ID {
			continue
		}
		i = idx
		break
	}
	if i < 0 {
		return ErrNoNote
	}

	if i < len(repo.data)-1 {
		copy(repo.data[i:], repo.data[i+1:])
	}
	repo.data[len(repo.data)-1] = nil
	repo.data = repo.data[:len(repo.data)-1]
	return nil
}
