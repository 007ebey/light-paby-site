package handlers

import (
	"word_press/models"
	"word_press/auth"
	"word_press/services"
	"net/http"
	"github.com/gorilla/mux"
    "strconv"
	"log"
	"strings"
)

type BlogHandler struct {
	PostService    services.PostService
	CommentService services.CommentService
}

func NewBlogHandler(ps services.PostService, cs services.CommentService) *BlogHandler {
	return &BlogHandler{
		PostService:    ps,
		CommentService: cs,
	}
}

func (h *BlogHandler) BlogPost(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.GetCurrentUser(r)

	slug := mux.Vars(r)["slug"]

	post, err := h.PostService.GetPostBySlug(slug)
	if err != nil || post == nil {
		http.NotFound(w, r)
		return
	}

	comments, err := h.CommentService.GetCommentsByPost(post.ID)
	if err != nil {
		// don't break page, just log
		log.Println("comments error:", err)
		comments = []models.Comment{}
	}

	RenderWithOpts(w, RenderOptions{
		Page:   "blog-post.html",
		Header: "blog-header.html",
		Footer: "default",
		Data: map[string]interface{}{
			"Post":     post,
			"Comments": comments,
			"User":     user,
		},
	})
}

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

func (h *BlogHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user, err := auth.GetCurrentUser(r)
	if err != nil {
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

	// 🔐 Authorization
	isOwner := comment.Email == user.Email
	isAdmin := user.Role == "admin"

	if !isOwner && !isAdmin {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	err = h.CommentService.DeleteComment(id)
	if err != nil {
		http.Error(w, "Failed to delete comment", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/blog/"+slug+"#comments", http.StatusSeeOther)
}
