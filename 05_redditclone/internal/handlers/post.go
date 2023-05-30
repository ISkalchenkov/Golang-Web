package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"redditclone/internal/comment"
	"redditclone/internal/post"
	"redditclone/internal/session"
	"redditclone/internal/user"
	"redditclone/internal/utils"
	"redditclone/internal/vote"
	"sort"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type PostHandler struct {
	Logger      *zap.SugaredLogger
	PostRepo    post.PostRepo
	UserRepo    user.UserRepo
	CommentRepo comment.CommentRepo
	VoteRepo    vote.VoteRepo
}

type AddPostReq struct {
	Category string `json:"category" valid:"required"`
	Title    string `json:"title" valid:"required"`
	Type     string `json:"type" valid:"in(link|text),required"`
	Text     string `json:"text" valid:"optional"`
	URL      string `json:"url" valid:"url, optional"`
}

type AddCommentReq struct {
	Comment string `json:"comment" valid:"required"`
}

type PostResp struct {
	ID               uint64         `json:"id"`
	Views            uint64         `json:"views"`
	Score            int64          `json:"score"`
	UpvotePercentage uint8          `json:"upvotePercentage"`
	Author           *Author        `json:"author"`
	Title            string         `json:"title"`
	Type             string         `json:"type"`
	Text             string         `json:"text,omitempty"`
	URL              string         `json:"url,omitempty"`
	Comments         []*CommentResp `json:"comments"`
	Category         string         `json:"category"`
	Created          time.Time      `json:"created"`
	Votes            []*VoteResp    `json:"votes"`
}

type Author struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
}

type CommentResp struct {
	ID      uint64    `json:"id"`
	Body    string    `json:"body"`
	Author  *Author   `json:"author"`
	Created time.Time `json:"created"`
}

type VoteResp struct {
	UserID uint64 `json:"user"`
	Value  int    `json:"vote"`
}

func (h *PostHandler) AddPost(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.Logger.Errorf("ReadAll failed: %v", err)
		utils.JSONMessage(w, "failed to read payload", http.StatusInternalServerError)
		return
	}

	req := &AddPostReq{}
	err = json.Unmarshal(body, req)
	if err != nil {
		utils.JSONMessage(w, "unmarshal payload failed", http.StatusBadRequest)
		return
	}

	_, err = govalidator.ValidateStruct(req)
	if err != nil {
		utils.JSONMessage(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Type == "link" && req.URL == "" {
		utils.JSONMessage(w, "post has link type, url field required", http.StatusBadRequest)
		return
	}

	if req.Type == "text" && req.Text == "" {
		utils.JSONMessage(w, "post has text type, text field required", http.StatusBadRequest)
		return
	}

	sess, err := session.SessionFromContext(r.Context())
	if err != nil {
		h.Logger.Errorf("SessionFromContext failed: %v", err)
		utils.JSONMessage(w, "getting session from context failed", http.StatusInternalServerError)
		return
	}

	post := &post.Post{
		Views:    0,
		AuthorID: sess.UserID,
		Title:    req.Title,
		Type:     req.Type,
		Category: req.Category,
		Created:  time.Now(),
	}

	if req.Type == "link" {
		post.URL = req.URL
	}

	if req.Type == "text" {
		post.Text = req.Text
	}

	postID, err := h.PostRepo.Add(post)
	if err != nil {
		h.Logger.Errorf("Add post failed: %v", err)
		utils.JSONMessage(w, "creating post failed", http.StatusInternalServerError)
		return
	}

	err = h.VoteRepo.Upvote(postID, sess.UserID)
	if err != nil {
		h.Logger.Errorf("Upvote failed: %v", err)
		utils.JSONMessage(w, "upvote failed", http.StatusInternalServerError)
		return
	}

	resp, err := h.getPostRespByPost(post)
	if err != nil {
		h.Logger.Errorf("getPostRespByPost failed: %v", err)
		utils.JSONMessage(w, "failed to get post", http.StatusInternalServerError)
		return
	}
	utils.JSONSend(w, resp, http.StatusCreated)
}

func (h *PostHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	postResps, err := h.getAllPostResps()
	if err != nil {
		h.Logger.Errorf("getAllPostResps failed: %v", err)
		utils.JSONMessage(w, "failed to get posts", http.StatusInternalServerError)
		return
	}
	sortByScore(postResps)
	utils.JSONSend(w, postResps, http.StatusOK)
}

func (h *PostHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID, err := strconv.ParseUint(vars["post_id"], 10, 64)
	if err != nil {
		utils.JSONMessage(w, "bad post id", http.StatusBadRequest)
		return
	}

	_, err = h.PostRepo.IncrementViews(postID)
	if err != nil {
		h.Logger.Errorf("IncrementViews failed: %v", err)
	}

	postResp, err := h.getPostRespByPostID(postID)
	if err != nil {
		switch {
		case errors.Is(err, post.ErrNoPost):
			utils.JSONMessage(w, "post does not exist", http.StatusUnprocessableEntity)
		default:
			h.Logger.Errorf("getPostRespByPostID failed: %v", err)
			utils.JSONMessage(w, "failed to get post", http.StatusInternalServerError)
		}
		return
	}
	utils.JSONSend(w, postResp, http.StatusOK)
}

func (h *PostHandler) GetByCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	category, ok := vars["category"]
	if !ok {
		utils.JSONMessage(w, "failed to determine the category", http.StatusInternalServerError)
		return
	}

	posts, err := h.PostRepo.GetByCategory(category)
	if err != nil {
		h.Logger.Errorf("GetByCategory failed: %v", err)
		utils.JSONMessage(w, "failed to get posts", http.StatusInternalServerError)
		return
	}

	postResps, err := h.getPostRespsByPosts(posts)
	if err != nil {
		h.Logger.Errorf("getPostRespsByPosts failed: %v", err)
		utils.JSONMessage(w, "failed to get posts", http.StatusInternalServerError)
		return
	}
	sortByScore(postResps)
	utils.JSONSend(w, postResps, http.StatusOK)
}

func (h *PostHandler) GetByUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username, ok := vars["username"]
	if !ok {
		utils.JSONMessage(w, "failed to determine the username", http.StatusInternalServerError)
		return
	}

	u, err := h.UserRepo.GetByUsername(username)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrNoUser):
			utils.JSONMessage(w, "user does not exist", http.StatusUnprocessableEntity)
		default:
			h.Logger.Errorf("GetByUsername failed: %v", err)
			utils.JSONMessage(w, "failed to get user", http.StatusInternalServerError)
		}
		return
	}

	posts, err := h.PostRepo.GetByAuthorID(u.ID)
	if err != nil {
		h.Logger.Errorf("GetByAuthorID failed: %v", err)
		utils.JSONMessage(w, "failed to get posts", http.StatusInternalServerError)
		return
	}

	postResps, err := h.getPostRespsByPosts(posts)
	if err != nil {
		h.Logger.Errorf("getPostRespsByPosts failed: %v", err)
		utils.JSONMessage(w, "failed to get posts", http.StatusInternalServerError)
		return
	}
	sortByScore(postResps)
	utils.JSONSend(w, postResps, http.StatusOK)
}

func (h *PostHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID, err := strconv.ParseUint(vars["post_id"], 10, 64)
	if err != nil {
		utils.JSONMessage(w, "bad post id", http.StatusBadRequest)
		return
	}

	sess, err := session.SessionFromContext(r.Context())
	if err != nil {
		h.Logger.Errorf("SessionFromContext failed: %v", err)
		utils.JSONMessage(w, "getting session from context failed", http.StatusInternalServerError)
		return
	}

	p, err := h.PostRepo.GetByID(postID)
	if err != nil {
		switch {
		case errors.Is(err, post.ErrNoPost):
			utils.JSONMessage(w, "post does not exist", http.StatusUnprocessableEntity)
		default:
			h.Logger.Errorf("GetById failed: %v", err)
			utils.JSONMessage(w, "post does not exist", http.StatusNotFound)
		}
		return
	}

	if sess.UserID != p.AuthorID {
		utils.JSONMessage(w, "you are not the author of this post", http.StatusForbidden)
		return
	}

	if _, err = h.PostRepo.Delete(postID); err != nil {
		h.Logger.Errorf("Delete post failed: %v", err)
		utils.JSONMessage(w, "failed to delete post", http.StatusInternalServerError)
		return
	}

	if err = h.CommentRepo.DeleteAllByPostID(postID); err != nil {
		h.Logger.Errorf("DeleteAllByPostID failed: %v", err)
	}

	if err = h.VoteRepo.DeleteAllByPostID(postID); err != nil {
		h.Logger.Errorf("DeleteAllByPostID failed: %v", err)
	}

	utils.JSONMessage(w, "success", http.StatusOK)
}

func (h *PostHandler) AddComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID, err := strconv.ParseUint(vars["post_id"], 10, 64)
	if err != nil {
		utils.JSONMessage(w, "bad post id", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.Logger.Errorf("ReadAll failed: %v", err)
		utils.JSONMessage(w, "failed to read payload", http.StatusInternalServerError)
		return
	}

	req := &AddCommentReq{}
	err = json.Unmarshal(body, req)
	if err != nil {
		utils.JSONMessage(w, "unmarshal payload failed", http.StatusBadRequest)
		return
	}

	_, err = govalidator.ValidateStruct(req)
	if err != nil {
		utils.JSONMessage(w, err.Error(), http.StatusBadRequest)
		return
	}

	sess, err := session.SessionFromContext(r.Context())
	if err != nil {
		h.Logger.Errorf("SessionFromContext failed: %v", err)
		utils.JSONMessage(w, "getting session from context failed", http.StatusInternalServerError)
		return
	}
	comment := &comment.Comment{
		AuthorID: sess.UserID,
		PostID:   postID,
		Body:     req.Comment,
		Created:  time.Now(),
	}
	_, err = h.CommentRepo.Add(comment)
	if err != nil {
		h.Logger.Errorf("Add comment failed: %v", err)
		utils.JSONMessage(w, "creating comment failed", http.StatusInternalServerError)
		return
	}

	postResp, err := h.getPostRespByPostID(postID)
	if err != nil {
		switch {
		case errors.Is(err, post.ErrNoPost):
			utils.JSONMessage(w, "post does not exist", http.StatusUnprocessableEntity)
		default:
			h.Logger.Errorf("getPostRespByPostID failed: %v", err)
			utils.JSONMessage(w, "failed to get post", http.StatusInternalServerError)
		}
		return
	}
	utils.JSONSend(w, postResp, http.StatusCreated)
}

func (h *PostHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID, err := strconv.ParseUint(vars["post_id"], 10, 64)
	if err != nil {
		utils.JSONMessage(w, "bad post id", http.StatusBadRequest)
		return
	}

	commentID, err := strconv.ParseUint(vars["comment_id"], 10, 64)
	if err != nil {
		utils.JSONMessage(w, "bad comment id", http.StatusBadRequest)
		return
	}

	sess, err := session.SessionFromContext(r.Context())
	if err != nil {
		h.Logger.Errorf("SessionFromContext failed: %v", err)
		utils.JSONMessage(w, "getting session from context failed", http.StatusInternalServerError)
		return
	}

	c, err := h.CommentRepo.GetByID(commentID)
	if err != nil {
		switch {
		case errors.Is(err, comment.ErrNoComment):
			utils.JSONMessage(w, "post does not exist", http.StatusUnprocessableEntity)
		default:
			h.Logger.Errorf("GetByID failed: %v", err)
			utils.JSONMessage(w, "failed to get comment", http.StatusInternalServerError)
		}
		return
	}

	if sess.UserID != c.AuthorID {
		utils.JSONMessage(w, "you are not the author of this comment", http.StatusForbidden)
		return
	}

	_, err = h.CommentRepo.Delete(commentID)
	if err != nil {
		h.Logger.Errorf("Delete comment failed: %v", err)
		utils.JSONMessage(w, "failed to delete comment", http.StatusInternalServerError)
		return
	}
	postResp, err := h.getPostRespByPostID(postID)
	if err != nil {
		switch {
		case errors.Is(err, post.ErrNoPost):
			utils.JSONMessage(w, "post does not exist", http.StatusUnprocessableEntity)
		default:
			h.Logger.Errorf("getPostRespByPostID failed: %v", err)
			utils.JSONMessage(w, "failed to get post", http.StatusInternalServerError)
		}
		return
	}
	utils.JSONSend(w, postResp, http.StatusOK)
}

func (h *PostHandler) Upvote(w http.ResponseWriter, r *http.Request) {
	h.vote(w, r, h.VoteRepo.Upvote)
}

func (h *PostHandler) Unvote(w http.ResponseWriter, r *http.Request) {
	h.vote(w, r, h.VoteRepo.Unvote)
}

func (h *PostHandler) Downvote(w http.ResponseWriter, r *http.Request) {
	h.vote(w, r, h.VoteRepo.Downvote)
}

func (h *PostHandler) vote(w http.ResponseWriter, r *http.Request, voteFunc func(postID, userID uint64) error) {
	vars := mux.Vars(r)
	postID, err := strconv.ParseUint(vars["post_id"], 10, 64)
	if err != nil {
		utils.JSONMessage(w, "bad post id", http.StatusBadRequest)
		return
	}
	sess, err := session.SessionFromContext(r.Context())
	if err != nil {
		h.Logger.Errorf("SessionFromContext failed: %v", err)
		utils.JSONMessage(w, "getting session from context failed", http.StatusInternalServerError)
		return
	}
	p, err := h.PostRepo.GetByID(postID)
	if err != nil {
		switch {
		case errors.Is(err, post.ErrNoPost):
			utils.JSONMessage(w, "post does not exist", http.StatusUnprocessableEntity)
		default:
			h.Logger.Errorf("GetById failed: %v", err)
			utils.JSONMessage(w, "post does not exist", http.StatusNotFound)
		}
		return
	}

	err = voteFunc(postID, sess.UserID)
	if err != nil {
		h.Logger.Errorf("voteFunc failed: %v", err)
		utils.JSONMessage(w, "failed to vote", http.StatusInternalServerError)
		return
	}

	postResp, err := h.getPostRespByPost(p)
	if err != nil {
		h.Logger.Errorf("getPostRespByPost failed: %v", err)
		utils.JSONMessage(w, "failed to get post", http.StatusInternalServerError)
		return
	}
	utils.JSONSend(w, postResp, http.StatusOK)
}

func (h *PostHandler) getAllPostResps() ([]*PostResp, error) {
	posts, err := h.PostRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("GetAll failed: %w", err)
	}
	postResps, err := h.getPostRespsByPosts(posts)
	if err != nil {
		return nil, fmt.Errorf("getPostRespsByPosts failed: %w", err)
	}
	return postResps, nil
}

func (h *PostHandler) getPostRespsByPosts(posts []*post.Post) ([]*PostResp, error) {
	postResps := []*PostResp{}
	for _, post := range posts {
		postResp, err := h.getPostRespByPost(post)
		if err != nil {
			return nil, fmt.Errorf("getPostRespByPost failed: %w", err)
		}
		postResps = append(postResps, postResp)
	}
	return postResps, nil
}

func (h *PostHandler) getPostRespByPostID(id uint64) (*PostResp, error) {
	post, err := h.PostRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("GetByID failed: %w", err)
	}
	postResp, err := h.getPostRespByPost(post)
	if err != nil {
		return nil, fmt.Errorf("getPostRespByPost failed: %w", err)
	}
	return postResp, nil
}

func (h *PostHandler) getPostRespByPost(post *post.Post) (*PostResp, error) {
	comments, err := h.getCommentRespsByPostID(post.ID)
	if err != nil {
		return nil, fmt.Errorf("getCommentRespsByPostID failed: %w", err)
	}
	votes, err := h.getVoteRespsByPostID(post.ID)
	if err != nil {
		return nil, fmt.Errorf("getVoteRespsByPostID failed: %w", err)
	}
	author, err := h.getAuthorByUserID(post.AuthorID)
	if err != nil {
		return nil, fmt.Errorf("getAuthorByUserID failed: %w", err)
	}
	score := calculateScore(votes)
	upvotePercentage := calculateUpvotePercentage(votes)

	postResp := &PostResp{
		ID:               post.ID,
		Score:            score,
		Views:            post.Views,
		UpvotePercentage: upvotePercentage,
		Author:           author,
		Title:            post.Title,
		Type:             post.Type,
		Text:             post.Text,
		URL:              post.URL,
		Comments:         comments,
		Category:         post.Category,
		Created:          post.Created,
		Votes:            votes,
	}
	return postResp, nil
}

func (h *PostHandler) getCommentRespsByPostID(id uint64) ([]*CommentResp, error) {
	comments, err := h.CommentRepo.GetByPostID(id)
	if err != nil {
		return nil, fmt.Errorf("GetByPostID failed: %w", err)
	}
	commentResps := []*CommentResp{}
	for _, comment := range comments {
		author, err := h.getAuthorByUserID(comment.AuthorID)
		if err != nil {
			return nil, fmt.Errorf("getAuthorByUserID failed: %w", err)
		}
		c := &CommentResp{
			ID:      comment.ID,
			Author:  author,
			Body:    comment.Body,
			Created: comment.Created,
		}
		commentResps = append(commentResps, c)
	}
	return commentResps, nil
}

func (h *PostHandler) getVoteRespsByPostID(id uint64) ([]*VoteResp, error) {
	votes, err := h.VoteRepo.GetByPostID(id)
	if err != nil {
		return nil, fmt.Errorf("GetByPostID failed: %w", err)
	}
	voteResps := []*VoteResp{}
	for _, vote := range votes {
		v := &VoteResp{
			UserID: vote.UserID,
			Value:  vote.Value,
		}
		voteResps = append(voteResps, v)
	}
	return voteResps, nil
}

func (h *PostHandler) getAuthorByUserID(id uint64) (*Author, error) {
	u, err := h.UserRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("GetByID failed: %w", err)
	}
	author := &Author{
		ID:       u.ID,
		Username: u.Username,
	}
	return author, nil
}

func calculateScore(votes []*VoteResp) int64 {
	var score int64 = 0
	for _, vote := range votes {
		score += int64(vote.Value)
	}
	return score
}

func calculateUpvotePercentage(votes []*VoteResp) uint8 {
	votesNumber := len(votes)
	if votesNumber == 0 {
		return 0
	}
	var upvoteNumber float64 = 0
	for _, v := range votes {
		if v.Value == vote.Upvote {
			upvoteNumber++
		}
	}
	upvotePercentage := (upvoteNumber / float64(votesNumber)) * 100
	upvotePercentage = math.Round(upvotePercentage)
	return uint8(upvotePercentage)
}

func sortByScore(posts []*PostResp) {
	less := func(i, j int) bool {
		return posts[i].Score > posts[j].Score
	}
	sort.SliceStable(posts, less)
}
