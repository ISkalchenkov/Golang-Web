package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"redditclone/internal/session"
	"redditclone/internal/user"
	"redditclone/internal/utils"

	"github.com/asaskevich/govalidator"
	"go.uber.org/zap"
)

type UserHandler struct {
	Logger   *zap.SugaredLogger
	UserRepo user.UserRepo
	Sessions session.SessionManager
}

type AuthReq struct {
	Username string `json:"username" valid:"alphanum,required,length(1|20)"`
	Password string `json:"password" valid:"printableascii,required,length(8|20)"`
}

type AuthResp struct {
	Token string `json:"token"`
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	req := h.getAuthReq(w, r)
	if req == nil {
		return
	}

	u, err := h.UserRepo.Authorize(req.Username, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrNoUser):
			fallthrough
		case errors.Is(err, user.ErrBadPass):
			utils.JSONMessage(w, "username or password invalid", http.StatusUnauthorized)
		default:
			h.Logger.Errorf("Authorize failed: %v", err)
			utils.JSONMessage(w, "authorization failed", http.StatusInternalServerError)
		}
		return
	}

	sess, err := h.Sessions.Create(u)
	if err != nil {
		h.Logger.Errorf("Create session failed: %v", err)
		utils.JSONMessage(w, "session creation failed", http.StatusInternalServerError)
		return
	}

	resp := AuthResp{Token: sess.Token}
	utils.JSONSend(w, resp, http.StatusOK)
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	req := h.getAuthReq(w, r)
	if req == nil {
		return
	}

	u, err := h.UserRepo.Registrate(req.Username, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrUsernameTaken):
			utils.JSONMessage(w, "username already exists", http.StatusUnprocessableEntity)
		default:
			h.Logger.Errorf("Registrate failed: %v", err)
			utils.JSONMessage(w, "registration failed", http.StatusInternalServerError)
		}
		return
	}

	sess, err := h.Sessions.Create(u)
	if err != nil {
		h.Logger.Errorf("Create session failed: %v", err)
		utils.JSONMessage(w, "session creation failed", http.StatusInternalServerError)
		return
	}

	resp := AuthResp{Token: sess.Token}
	utils.JSONSend(w, resp, http.StatusCreated)
}

func (h *UserHandler) getAuthReq(w http.ResponseWriter, r *http.Request) *AuthReq {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.Logger.Errorf("ReadAll failed: %v", err)
		utils.JSONMessage(w, "failed to read payload", http.StatusInternalServerError)
		return nil
	}
	req := &AuthReq{}
	err = json.Unmarshal(body, req)
	if err != nil {
		utils.JSONMessage(w, "unmarshal payload failed", http.StatusBadRequest)
		return nil
	}
	_, err = govalidator.ValidateStruct(req)
	if err != nil {
		utils.JSONSend(w, err.Error(), http.StatusBadRequest)
		return nil
	}
	return req
}
