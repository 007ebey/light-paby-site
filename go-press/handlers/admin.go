package handlers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"word_press/auth"
	"word_press/models"
	"word_press/services"
	"word_press/utils"
)

type AdminHandler struct {
	PostService    services.PostService
	ImageService   services.ImageService
	ContactService services.ContactService

	GetUser func(*http.Request) (*models.User, error)
	Render  func(http.ResponseWriter, RenderOptions)
}

func NewAdminHandler(ps services.PostService, is services.ImageService, cs services.ContactService) *AdminHandler {
	return &AdminHandler{
		PostService:    ps,
		ImageService:   is,
		ContactService: cs,
		GetUser:        auth.GetCurrentUser,
		Render:         RenderWithOpts,
	}
}

//
// ===== COMMON RENDER =====
//

func (h *AdminHandler) render(w http.ResponseWriter, page string, data map[string]interface{}) {
	h.Render(w, RenderOptions{
		Page:   page,
		Header: "posts-header.html",
		Footer: "posts-footer.html",
		Data:   data,
	})
}

//
// ===== POSTS =====
//

func (h *AdminHandler) AdminPosts(w http.ResponseWriter, r *http.Request) {
	user, _ := h.GetUser(r)

	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	posts, err := h.PostService.GetAllPosts()
	if err != nil {
		http.Error(w, "Failed to load posts", http.StatusInternalServerError)
		return
	}

	h.render(w, "admin/posts.html", map[string]interface{}{
		"SiteTitle": "Manage Posts",
		"Posts":     posts,
		"User":      user,
	})
}

func (h *AdminHandler) AdminCreatePost(w http.ResponseWriter, r *http.Request) {
	user, _ := h.GetUser(r)

	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	data := map[string]interface{}{
		"SiteTitle": "Create Post",
		"User":      user,
	}

	if r.Method == http.MethodPost {
		title := r.FormValue("title")
		slug := r.FormValue("slug")
		content := r.FormValue("content")
		status := r.FormValue("status")
		excerpt := r.FormValue("excerpt")

		if title == "" || content == "" {
			data["Error"] = "Title and Content required"
			data["Form"] = map[string]string{
				"title":   title,
				"slug":    slug,
				"content": content,
			}
			h.render(w, "admin/create-post.html", data)
			return
		}

		if slug == "" {
			slug = utils.GenerateSlug(title)
		}

		var imagePath string

		file, handler, err := r.FormFile("featured_image")
		if err == nil {
			defer file.Close()

			imagePath, err = h.ImageService.SaveImage(file, handler.Filename, handler.Size)
			if err != nil {
				http.Error(w, "Image upload failed", http.StatusBadRequest)
				return
			}
		}

		post := &models.Post{
			Title:         title,
			Slug:          slug,
			Content:       content,
			Status:        status,
			Excerpt:       excerpt,
			FeaturedImage: imagePath,
			AuthorID:      user.ID,
		}

		if err := h.PostService.CreatePost(post); err != nil {
			data["Error"] = err.Error()
			h.render(w, "admin/create-post.html", data)
			return
		}

		http.Redirect(w, r, "/admin/posts", http.StatusSeeOther)
		return
	}

	h.render(w, "admin/create-post.html", data)
}

func (h *AdminHandler) AdminEditPost(w http.ResponseWriter, r *http.Request) {
	user, _ := h.GetUser(r)

	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	post, err := h.PostService.GetPostByID(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	data := map[string]interface{}{
		"User": user,
		"Post": post,
	}

	if r.Method == http.MethodPost {

		post.Title = r.FormValue("title")
		post.Slug = r.FormValue("slug")
		post.Content = r.FormValue("content")
		post.Status = r.FormValue("status")
		post.Excerpt = r.FormValue("excerpt")

		file, handler, err := r.FormFile("featured_image")
		if err == nil {
			defer file.Close()

			imagePath, err := h.ImageService.SaveImage(file, handler.Filename, handler.Size)
			if err != nil {
				http.Error(w, "Image upload failed", http.StatusBadRequest)
				return
			}
			post.FeaturedImage = imagePath
		}

		post.AuthorID = user.ID

		if err := h.PostService.UpdatePost(post); err != nil {
			data["Error"] = err.Error()
			h.render(w, "admin/edit-post.html", data)
			return
		}

		http.Redirect(w, r, "/admin/posts", http.StatusSeeOther)
		return
	}

	h.render(w, "admin/edit-post.html", data)
}

func (h *AdminHandler) AdminDeletePost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.NotFound(w, r)
		return
	}

	post, err := h.PostService.GetPostByID(id)
	if err != nil || post == nil {
		http.NotFound(w, r)
		return
	}

	if err := h.PostService.DeletePost(id); err != nil {
		http.Error(w, "Failed to delete post", http.StatusInternalServerError)
		return
	}

	if post.FeaturedImage != "" {
		h.ImageService.DeleteImage(post.FeaturedImage)
	}

	http.Redirect(w, r, "/admin/posts", http.StatusSeeOther)
}

//
// ===== CONTACTS =====
//

func (h *AdminHandler) ManageQueries(w http.ResponseWriter, r *http.Request) {

	contacts, err := h.ContactService.GetAllContacts()
	if err != nil {
		http.Error(w, "Failed to load contacts", http.StatusInternalServerError)
		return
	}

	h.Render(w, RenderOptions{
		Page:   "admin/contact.html",
		Header: "contact-header.html",
		Footer: "contact-footer.html",
		Data: map[string]interface{}{
			"SiteTitle": "Manage Queries",
			"Contacts":  contacts,
		},
	})
}

func (h *AdminHandler) MarkContactRead(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, _ := strconv.Atoi(r.FormValue("id"))

	if err := h.ContactService.MarkAsRead(id); err != nil {
		http.Error(w, "Failed to update contact", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/contact", http.StatusSeeOther)
}

func (h *AdminHandler) DeleteContact(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, _ := strconv.Atoi(r.FormValue("id"))

	if err := h.ContactService.DeleteContact(id); err != nil {
		http.Error(w, "Failed to delete contact", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/contact", http.StatusSeeOther)
}