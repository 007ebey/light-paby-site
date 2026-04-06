package handlers

import (
	"net/http"
	"strings"
	"word_press/models"
    "word_press/auth"
	"strconv"
)

func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	slug := pathParts[2]
	// Get post by slug
	post, err := models.GetPostBySlug(slug)
    if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	author := strings.TrimSpace(r.FormValue("author"))
	email  := strings.TrimSpace(r.FormValue("email"))
	body   := strings.TrimSpace(r.FormValue("body"))

	if author == "" || email == "" || body == "" {
		http.Redirect(w, r, "/blog/" + slug + "?error=1", http.StatusSeeOther)
		return
	}

	if !strings.Contains(email, "@") {
		http.Redirect(w, r, "/blog/" + slug + "?error=invalid_email", http.StatusSeeOther)
		return
	}

	err = models.CreateComment(post.ID, author, email, body) 

	if err != nil {
		http.Error(w, "Failed to save comment", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/blog/"+ slug + "?comment=success#comments", http.StatusSeeOther)
}

func DeleteCommentHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	user, err := auth.GetCurrentUser(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	//Read inputs
	idStr := r.FormValue("id")
	slug := r.FormValue("slug")

	if idStr == "" || slug == "" {
		http.Error(w, "Missing data", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	post, err := models.GetPostBySlug(slug)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	comments, err := models.GetCommentsByPost(post.ID)
	if err != nil {
        http.Error(w, "Failed to load comments", http.StatusInternalServerError)
	}

	// Find the comment
	var target *models.Comment
	for _, c := range comments {
		if c.ID == id {
			target = &c
			break
		}
	}

	if target == nil {
		http.Error(w, "Comment not found", http.StatusNotFound)
		return
	}

	// Authorization
	isOwner := target.Email == user.Email
	isAdmin := user.Role == "admin"

	if !isOwner && !isAdmin {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Delete
	err = models.DeleteComment(id)
	if err != nil {
		http.Error(w, "Failed to delete comment", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/blog/"+slug+"#comments", http.StatusSeeOther)
}