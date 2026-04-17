package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"log"

	"github.com/gorilla/mux"

	"word_press/models"
	"word_press/auth"
	"word_press/services"
)

type BlogHandler struct {
	PostService    services.PostService
	CommentService services.CommentService

	GetUser func(*http.Request) (*models.User, error)
	Render  func(http.ResponseWriter, RenderOptions)
}

func NewBlogHandler(ps services.PostService, cs services.CommentService) *BlogHandler {
	return &BlogHandler{
		PostService:    ps,
		CommentService: cs,
		GetUser:        auth.GetCurrentUser,
		Render:         RenderWithOpts,
	}
}

//
// ===== COMMON RENDER =====
//

func (h *BlogHandler) render(w http.ResponseWriter, page string, data map[string]interface{}) {
	h.Render(w, RenderOptions{
		Page:   page,
		Header: "blog-header.html",
		Footer: "default",
		Data:   data,
	})
}

//
// ===== BLOG POST =====
//

func (h *BlogHandler) BlogPost(w http.ResponseWriter, r *http.Request) {
	user, _ := h.GetUser(r)

	slug := mux.Vars(r)["slug"]

	post, err := h.PostService.GetPostBySlug(slug)
	if err != nil || post == nil {
		http.NotFound(w, r)
		return
	}

	comments, err := h.CommentService.GetCommentsByPost(post.ID)
	if err != nil {
		log.Println("comments error:", err)
		comments = []models.Comment{}
	}

	h.render(w, "blog-post.html", map[string]interface{}{
		"Post":     post,
		"Comments": comments,
		"User":     user,
	})
}

//
// ===== CREATE COMMENT =====
//

func (h *BlogHandler) CreateComment(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	slug := mux.Vars(r)["slug"]

	post, err := h.PostService.GetPostBySlug(slug)
	if err != nil || post == nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	author := strings.TrimSpace(r.FormValue("author"))
	email := strings.TrimSpace(r.FormValue("email"))
	body := strings.TrimSpace(r.FormValue("body"))

	err = h.CommentService.CreateComment(post.ID, author, email, body)
	if err != nil {
		http.Redirect(w, r, "/blog/"+slug+"?error=1", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/blog/"+slug+"?comment=success#comments", http.StatusSeeOther)
}

//
// ===== DELETE COMMENT =====
//

func (h *BlogHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user, err := h.GetUser(r)
	if err != nil || user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	id, err := strconv.Atoi(r.FormValue("id"))
	slug := r.FormValue("slug")

	if err != nil || id == 0 || slug == "" {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	post, err := h.PostService.GetPostBySlug(slug)
	if err != nil || post == nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	comment, err := h.CommentService.GetCommentByID(post.ID, id)
	if err != nil {
		http.Error(w, "Comment not found", http.StatusNotFound)
		return
	}

	isOwner := comment.Email == user.Email
	isAdmin := user.Role == "admin"

	if !isOwner && !isAdmin {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	if err := h.CommentService.DeleteComment(id); err != nil {
		http.Error(w, "Failed to delete comment", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/blog/"+slug+"#comments", http.StatusSeeOther)
}