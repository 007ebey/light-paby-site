package handlers

import (
	"word_press/models"
	"word_press/auth"
	"word_press/utils"
	"word_press/services"
	"net/http"
	"github.com/gorilla/mux"
    "strconv"
)

type AdminHandler struct {
	PostService services.PostService
	ImageService services.ImageService
	ContactService services.ContactService
}

var getCurrentUser = auth.GetCurrentUser
var render = RenderWithOpts

func NewAdminHandler(s services.PostService, is services.ImageService, cs services.ContactService) *AdminHandler {
	return &AdminHandler{
		PostService: s,
		ImageService: is,
		ContactService: cs,
	}
}

func (h *AdminHandler) AdminPosts(w http.ResponseWriter, r *http.Request) {
	user, _ := getCurrentUser(r)

	if user == nil {
	  http.Error(w, "Unauthorized", 401)
	  return
    }

	posts, err := h.PostService.GetAllPosts()
	if err != nil {
		http.Error(w, "Failed to load posts", http.StatusInternalServerError)
		return
	}

	render(w, RenderOptions{
		Page:   "admin/posts.html",
		Header: "posts-header.html",
		Footer: "posts-footer.html",
		Data: map[string]interface{}{
			"SiteTitle": "Manage Posts",
		    "Posts": posts,
		    "User": user,
		},
	})
}

func (h *AdminHandler) AdminCreatePost(w http.ResponseWriter, r *http.Request) {
	user, _ := getCurrentUser(r)

	if user == nil {
	  http.Error(w, "Unauthorized", 401)
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
			render(w, RenderOptions{
		          Page:   "admin/create-post.html",
		          Header: "posts-header.html",
		          Footer: "posts-footer.html",
		          Data: data,
	        })
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
			AuthorID:        user.ID,
		}

		err = h.PostService.CreatePost(post)
		if err != nil {
			data["Error"] = err.Error()
			render(w, RenderOptions{
		          Page:   "admin/create-post.html",
		          Header: "posts-header.html",
		          Footer: "posts-footer.html",
		          Data: data,
	        })
			return
		}

		http.Redirect(w, r, "/admin/posts", http.StatusSeeOther)
		return
	}

	render(w, RenderOptions{
		Page:   "admin/create-post.html",
		Header: "posts-header.html",
		Footer: "posts-footer.html",
		Data: data,
	})
}

func (h *AdminHandler) AdminEditPost(w http.ResponseWriter, r *http.Request) {
	user, _ := getCurrentUser(r)

	if user == nil {
	  http.Error(w, "Unauthorized", 401)
	  return
    }

	idStr := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(idStr)

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

		title := r.FormValue("title")
		slug := r.FormValue("slug")
		content := r.FormValue("content")
		status := r.FormValue("status")
		excerpt := r.FormValue("excerpt")

		imagePath := post.FeaturedImage

		file, handler, err := r.FormFile("featured_image")
		if err == nil {
			defer file.Close()

			imagePath, err = h.ImageService.SaveImage(file, handler.Filename, handler.Size)
			if err != nil {
				http.Error(w, "Image upload failed", http.StatusBadRequest)
				return
			}
		}

		post.Title = title
		post.Slug = slug
		post.Content = content
		post.Status = status
		post.Excerpt = excerpt
		post.FeaturedImage = imagePath
		post.AuthorID = user.ID

		err = h.PostService.UpdatePost(post)
		if err != nil {
			data["Error"] = err.Error()
			render(w, RenderOptions{
		      Page:   "admin/edit-post.html",
		      Header: "posts-header.html",
		      Footer: "posts-footer.html",
		      Data: data,
	        })
			return
		}

		http.Redirect(w, r, "/admin/posts", http.StatusSeeOther)
		return
	}

	render(w, RenderOptions{
		Page:   "admin/edit-post.html",
		Header: "posts-header.html",
		Footer: "posts-footer.html",
		Data: data,
	})
}

func (h *AdminHandler) AdminDeletePost(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// 🔥 get post first (for image cleanup)
	post, err := h.PostService.GetPostByID(id)
	if err != nil || post == nil {
		http.NotFound(w, r)
		return
	}

	// 🔥 delete post via service
	err = h.PostService.DeletePost(id)
	if err != nil {
		http.Error(w, "Failed to delete post", http.StatusInternalServerError)
		return
	}

	// 🔥 optional: delete image (IMPORTANT)
	if post.FeaturedImage != "" {
		h.ImageService.DeleteImage(post.FeaturedImage)
	}

	http.Redirect(w, r, "/admin/posts", http.StatusSeeOther)
}

func (h *AdminHandler) ManageQueries(w http.ResponseWriter, r *http.Request) {

	contacts, err := h.ContactService.GetAllContacts()
	if err != nil {
		http.Error(w, "Failed to load contacts", http.StatusInternalServerError)
		return
	}

	render(w, RenderOptions{
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

	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil || id == 0 {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

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

	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil || id == 0 {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := h.ContactService.DeleteContact(id); err != nil {
		http.Error(w, "Failed to delete contact", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/contact", http.StatusSeeOther)
}
