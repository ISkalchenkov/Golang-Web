package handlers

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"gitlab.com/homework-ci-cd/internal/note"
)

func TestGetAll(t *testing.T) {
	type TestCase struct {
		StatusCode int
		URL        string
		Result     []note.Note
	}
	cases := []TestCase{
		{
			StatusCode: 200,
			URL:        "/note",
			Result: []note.Note{
				{
					ID:   1,
					Text: "send a message",
				},
				{
					ID:   2,
					Text: "buy food",
				},
			},
		},
		{
			StatusCode: 200,
			URL:        "/note?order_field=id",
			Result: []note.Note{
				{
					ID:   1,
					Text: "send a message",
				},
				{
					ID:   2,
					Text: "buy food",
				},
			},
		},
		{
			StatusCode: 200,
			URL:        "/note?order_field=text",
			Result: []note.Note{
				{
					ID:   2,
					Text: "buy food",
				},
				{
					ID:   1,
					Text: "send a message",
				},
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	st := note.NewMockNoteRepo(ctrl)
	service := &NoteHandler{
		NoteRepo: st,
	}

	resultNotes := []*note.Note{
		{
			ID:   1,
			Text: "send a message",
		},
		{
			ID:   2,
			Text: "buy food",
		},
	}

	for caseNum, item := range cases {
		st.EXPECT().GetAll().Return(resultNotes, nil)
		req := httptest.NewRequest("GET", item.URL, nil)
		w := httptest.NewRecorder()

		service.GetAll(w, req)
		resp := w.Result()
		defer resp.Body.Close()

		notes := []note.Note{}
		body, _ := io.ReadAll(resp.Body)
		_ = json.Unmarshal(body, &notes)
		if !reflect.DeepEqual(notes, item.Result) {
			t.Errorf("[%d] wrong result, expected %#v, got %#v", caseNum, item.Result, notes)
		}
		if resp.StatusCode != item.StatusCode {
			t.Errorf("[%d] wrong status code, expected %d, got %d", caseNum, item.StatusCode, resp.StatusCode)
		}
	}
}
