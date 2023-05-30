package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"gitlab.com/homework-ci-cd/internal/note"
	"gitlab.com/homework-ci-cd/internal/utils"
)

type NoteHandler struct {
	NoteRepo note.NoteRepo
}

type NoteReq struct {
	Text string `json:"text" valid:"required"`
}

type Options struct {
	OrderField string `schema:"order_field" valid:"in(id|text|created_at|updated_at), optional"`
}

func (h *NoteHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)

	opt := &Options{}
	if err := decoder.Decode(opt, r.URL.Query()); err != nil {
		utils.JSONError(w, "decode query params failed", http.StatusInternalServerError)
		return
	}

	if _, err := govalidator.ValidateStruct(opt); err != nil {
		utils.JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	notes, err := h.NoteRepo.GetAll()
	if err != nil {
		utils.JSONError(w, "failed to get notes", http.StatusInternalServerError)
		return
	}

	if opt.OrderField == "" {
		utils.JSONSend(w, notes, http.StatusOK)
		return
	}

	data := make([]*note.Note, len(notes))
	copy(data, notes)
	sortNotes(data, opt.OrderField)
	utils.JSONSend(w, data, http.StatusOK)
}

func (h *NoteHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	noteID, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		utils.JSONError(w, "bad note id", http.StatusBadRequest)
		return
	}

	n, err := h.NoteRepo.GetByID(noteID)
	if err != nil {
		switch {
		case errors.Is(err, note.ErrNoNote):
			utils.JSONError(w, "note does not exist", http.StatusUnprocessableEntity)
		default:
			utils.JSONError(w, "failed to get note", http.StatusInternalServerError)
		}
		return
	}

	utils.JSONSend(w, n, http.StatusOK)
}
func (h *NoteHandler) Add(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.JSONError(w, "failed to read payload", http.StatusInternalServerError)
		return
	}

	req := &NoteReq{}
	if err = json.Unmarshal(body, req); err != nil {
		utils.JSONError(w, "unmarshal payload failed", http.StatusBadRequest)
		return
	}

	if _, err = govalidator.ValidateStruct(req); err != nil {
		utils.JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	n := &note.Note{
		Text:      req.Text,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if _, err = h.NoteRepo.Add(n); err != nil {
		utils.JSONError(w, "failed to add new note", http.StatusInternalServerError)
		return
	}

	utils.JSONSend(w, n, http.StatusCreated)
}

func (h *NoteHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	noteID, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		utils.JSONError(w, "bad note id", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.JSONError(w, "failed to read payload", http.StatusInternalServerError)
		return
	}

	req := &NoteReq{}
	if err = json.Unmarshal(body, req); err != nil {
		utils.JSONError(w, "unmarshal payload failed", http.StatusBadRequest)
		return
	}

	if _, err = govalidator.ValidateStruct(req); err != nil {
		utils.JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	n := &note.Note{
		ID:        noteID,
		Text:      req.Text,
		UpdatedAt: time.Now(),
	}

	if n, err = h.NoteRepo.Update(n); err != nil {
		switch {
		case errors.Is(err, note.ErrNoNote):
			utils.JSONError(w, "note does not exist", http.StatusUnprocessableEntity)
		default:
			utils.JSONError(w, "failed to update note", http.StatusInternalServerError)
		}
		return
	}

	utils.JSONSend(w, n, http.StatusOK)

}
func (h *NoteHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	noteID, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		utils.JSONError(w, "bad note id", http.StatusBadRequest)
		return
	}

	if err := h.NoteRepo.Delete(noteID); err != nil {
		switch {
		case errors.Is(err, note.ErrNoNote):
			utils.JSONError(w, "note does not exist", http.StatusUnprocessableEntity)
		default:
			utils.JSONError(w, "failed to delete note", http.StatusInternalServerError)
		}
		return
	}

	success := map[string]string{"message": "success"}
	utils.JSONSend(w, success, http.StatusOK)
}

func sortNotes(data []*note.Note, orderField string) {
	less := func(i, j int) bool {
		switch orderField {
		case "id":
			return data[i].ID < data[j].ID
		case "text":
			return strings.ToLower(data[i].Text) < strings.ToLower(data[j].Text)
		case "created_at":
			return data[i].CreatedAt.Unix() < data[j].CreatedAt.Unix()
		case "updated_at":
			return data[i].UpdatedAt.Unix() < data[j].UpdatedAt.Unix()
		default:
			return data[i].ID < data[j].ID
		}
	}
	sort.SliceStable(data, less)
}
