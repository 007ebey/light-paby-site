package handlers

import (
	"word_press/models"
	"word_press/auth"
	"net/http"
	"github.com/gorilla/mux"
	"strings"
    "strconv"
)


func AdminPosts(w http.ResponseWriter, r *http.Request) {
	posts, _ := models.GetAllPosts()

	user, _ := auth.GetCurrentUser(r) 

	Render(w, "admin/posts.html", map[string]interface{}{
	    "SiteTitle": "Manage Posts",
	    "Posts": posts,
	    "User": user,
    })
}

func AdminCreatePost(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.GetCurrentUser(r)

	data := map[string]interface{}{
		"SiteTitle": "Create Post",
		"User": user,
	}

	if r.Method == http.MethodPost {
		title := r.FormValue("title")
		slug := r.FormValue("slug")
		content := r.FormValue("content")
		status := r.FormValue("status")

		if title == "" || content == "" {
			data["Error"] = "Title and Content required"
			data["Form"] = map[string]string{
				"title": title,
				"slug": slug,
				"content": content,
			}
			Render(w, "admin/create-post.html", data)
			return
		}

		if slug == "" {
			slug = strings.ToLower(strings.ReplaceAll(title, " ", "-"))
		}

		err := models.CreatePost(title, slug, content, status, user.ID)

		if err != nil {
			data["Error"] = err.Error()
			Render(w, "admin/create-post.html", data)
			return
		}

		http.Redirect(w, r, "/admin/posts", http.StatusSeeOther)
		return
	}

	Render(w, "admin/create-post.html", data)
}

func AdminEditPost(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.GetCurrentUser(r)

	idStr := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(idStr)

	post, err := models.GetPostByID(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	data := map[string]interface{}{
		"User": user,
		"Post": post,
	}

	if r.Method == http.MethodPost {

		title := r.FormValue("title")
		slug  := r.FormValue("slug")
		content := r.FormValue("content")
		status := r.FormValue("status")
		
		err := models.UpdatePost(id, title, slug, content, status)
		if err != nil {
			data["Error"] = err.Error()
			Render(w, "admin/edit-post.html", data)
			return
		}
		http.Redirect(w, r, "/admin/posts", http.StatusSeeOther)
		return
	}

	Render(w, "admin/edit-post.html", data)
}

func BlogPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]

	post, err := models.GetPostBySlug(slug)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	Render(w, "blog-post.html", map[string]interface{}{
		"Post": post,
	})
}


