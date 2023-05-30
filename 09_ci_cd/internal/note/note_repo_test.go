package note

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

// go test ./... -coverprofile=cover.out && go tool cover -html=cover.out -o cover.html && rm cover.out

func getInitialData() []*Note {
	initialData := []*Note{
		{
			ID:        1,
			Text:      "send a message",
			CreatedAt: time.Date(2023, time.January, 1, 0, 0, 0, 0, time.Local),
			UpdatedAt: time.Date(2023, time.January, 1, 0, 0, 0, 0, time.Local),
		},
		{
			ID:        2,
			Text:      "buy food",
			CreatedAt: time.Date(2023, time.February, 1, 0, 0, 0, 0, time.Local),
			UpdatedAt: time.Date(2023, time.February, 1, 0, 0, 0, 0, time.Local),
		},
		{
			ID:        3,
			Text:      "listen to music",
			CreatedAt: time.Date(2023, time.March, 1, 0, 0, 0, 0, time.Local),
			UpdatedAt: time.Date(2023, time.March, 1, 0, 0, 0, 0, time.Local),
		},
	}
	return initialData
}

func initNoteRepo(notes []*Note) *NoteMemoryRepository {
	noteRepo := NewMemoryRepo()

	for _, note := range notes {
		n := &Note{}
		*n = *note
		noteRepo.data = append(noteRepo.data, n)
	}
	noteRepo.lastID = uint64(len(notes))

	return noteRepo
}

func TestGetAll(t *testing.T) {
	initialData := getInitialData()
	noteRepo := initNoteRepo(initialData)
	notes, err := noteRepo.GetAll()

	if err != nil {
		t.Errorf("unexpected error: %#v", err)
	}
	if !reflect.DeepEqual(notes, initialData) {
		t.Errorf("wrong result, expected %#v, got %#v", initialData, notes)
	}
}

func TestGetByID(t *testing.T) {
	type TestCase struct {
		ID     uint64
		Result *Note
		Err    error
	}
	cases := []TestCase{
		{
			ID: 2,
			Result: &Note{
				ID:        2,
				Text:      "buy food",
				CreatedAt: time.Date(2023, time.February, 1, 0, 0, 0, 0, time.Local),
				UpdatedAt: time.Date(2023, time.February, 1, 0, 0, 0, 0, time.Local),
			},
			Err: nil,
		},
		{
			ID:     5,
			Result: nil,
			Err:    ErrNoNote,
		},
	}

	initialData := getInitialData()
	noteRepo := initNoteRepo(initialData)

	for caseNum, item := range cases {
		note, err := noteRepo.GetByID(item.ID)
		if !errors.Is(err, item.Err) {
			t.Errorf("[%d] wrong error, expected %#v, got %#v", caseNum, item.Err, err)
		}
		if !reflect.DeepEqual(note, item.Result) {
			t.Errorf("[%d] wrong result, expected %#v, got %#v", caseNum, item.Result, note)
		}
	}
}

func TestAdd(t *testing.T) {
	noteRepo := NewMemoryRepo()
	newNote := &Note{
		Text:      "send a message",
		CreatedAt: time.Date(2023, time.January, 1, 0, 0, 0, 0, time.Local),
		UpdatedAt: time.Date(2023, time.January, 1, 0, 0, 0, 0, time.Local),
	}
	id, err := noteRepo.Add(newNote)
	if err != nil {
		t.Errorf("unexpected error: %#v", err)
	}
	if id != 1 {
		t.Errorf("bad id, expected %v, got %v", id, 1)
	}
}

func TestUpdate(t *testing.T) {
	type TestCase struct {
		NewNote *Note
		Result  *Note
		Err     error
	}
	cases := []TestCase{
		{
			NewNote: &Note{
				ID:        2,
				Text:      "New text",
				UpdatedAt: time.Date(2023, time.April, 1, 0, 0, 0, 0, time.Local),
			},
			Result: &Note{
				ID:        2,
				Text:      "New text",
				CreatedAt: time.Date(2023, time.February, 1, 0, 0, 0, 0, time.Local),
				UpdatedAt: time.Date(2023, time.April, 1, 0, 0, 0, 0, time.Local),
			},
			Err: nil,
		},
		{
			NewNote: &Note{
				ID:        5,
				Text:      "New text",
				UpdatedAt: time.Date(2023, time.April, 1, 0, 0, 0, 0, time.Local),
			},
			Result: nil,
			Err:    ErrNoNote,
		},
	}

	initialData := getInitialData()
	noteRepo := initNoteRepo(initialData)

	for caseNum, item := range cases {
		note, err := noteRepo.Update(item.NewNote)
		if !errors.Is(err, item.Err) {
			t.Errorf("[%d] wrong error, expected %#v, got %#v", caseNum, item.Err, err)
		}
		if !reflect.DeepEqual(note, item.Result) {
			t.Errorf("[%d] wrong result, expected %#v, got %#v", caseNum, item.Result, note)
		}
	}
}

func TestDelete(t *testing.T) {
	type TestCase struct {
		ID  uint64
		Err error
	}
	cases := []TestCase{
		{
			ID:  2,
			Err: nil,
		},
		{
			ID:  2,
			Err: ErrNoNote,
		},
	}

	initialData := getInitialData()
	noteRepo := initNoteRepo(initialData)

	for caseNum, item := range cases {
		err := noteRepo.Delete(item.ID)
		if !errors.Is(err, item.Err) {
			t.Errorf("[%d] wrong error, expected %#v, got %#v", caseNum, item.Err, err)
		}
	}
}
