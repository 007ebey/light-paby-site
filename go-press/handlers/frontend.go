package handlers

import (
	"log"
	"net/http"
	"strconv"
    "strings"
	"word_press/auth"
	"word_press/models"
	"word_press/services"
)

const pageLimit = 10

type PageHandler struct {
	PostService services.PostService
    ContactService services.ContactService
	GetUser func(*http.Request) (*models.User, error)
	Render  func(http.ResponseWriter, RenderOptions)
}

func NewPageHandler(ps services.PostService, cs services.ContactService) *PageHandler {
	return &PageHandler{
		PostService: ps,
		ContactService: cs,
		GetUser:     auth.GetCurrentUser,
		Render:      RenderWithOpts,
	}
}

func (h *PageHandler) Home(w http.ResponseWriter, r *http.Request) {
	success := r.URL.Query().Get("success")

	data := map[string]interface{}{
		"SiteTitle": "Pastor Aby & Pastor Smitha George",
	}

	if success == "1" {
		data["Success"] = "Your message has been sent successfully!"
	}

	h.Render(w, RenderOptions{
		Page:   "index.html",
		Header: "home-header.html",
		Footer: "index-footer.html",
		Data:   data,
	})
}

func (h *PageHandler) Blog(w http.ResponseWriter, r *http.Request) {
	page := 1

	if p, err := strconv.Atoi(r.URL.Query().Get("page")); err == nil && p > 0 {
		page = p
	}

	posts, totalPages, err := h.PostService.GetPublishedPosts(page, pageLimit)
	if err != nil {
		log.Println("Blog fetch error:", err)
		http.Error(w, "Unable to load the blog at this time", http.StatusInternalServerError)
		return
	}

	user, err := h.GetUser(r)
	if err != nil {
		log.Println("user fetch error:", err)
	}

	showAdminLinks := user != nil &&
		(user.Role == "admin" || user.Role == "editor")

	h.Render(w, RenderOptions{
		Page:   "blog.html",
		Header: "home-header.html",
		Footer: "default",
		Data: map[string]interface{}{
			"SiteTitle":      "Blog",
			"Posts":          posts,
			"Page":           page,
			"TotalPages":     totalPages,
			"ShowAdminLinks": showAdminLinks,
		},
	})
}

func (h *PageHandler) Contact(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		success := r.URL.Query().Get("success")

		data := map[string]interface{}{
			"SiteTitle": "Pastor Aby & Pastor Smitha",
		}

		if success == "1" {
			data["Success"] = "Your message has been sent successfully!"
		}

		h.Render(w, RenderOptions{
			Page:   "index.html",
			Header: "home-header.html",
			Footer: "index-footer.html",
			Data:   data,
		})
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if r.FormValue("form_type") != "contact" {
		http.Error(w, "Invalid form submission", http.StatusBadRequest)
		return
	}

	name := strings.TrimSpace(r.FormValue("name"))
	email := strings.TrimSpace(r.FormValue("email"))
	message := strings.TrimSpace(r.FormValue("message"))

	if name == "" || email == "" || message == "" {
		h.renderContactError(w, "All fields are required", name, email, message)
		return
	}

	if !strings.Contains(email, "@") {
		h.renderContactError(w, "Invalid email address", name, email, message)
		return
	}

	ip := r.RemoteAddr

	err := h.ContactService.CreateContact(name, email, message, ip)
	if err != nil {
		log.Println("Failed to save contact:", err)
		h.renderContactError(w, "Something went wrong. Please try again", name, email, message)
		return
	}

	http.Redirect(w, r, "/?success=1#contact-form", http.StatusSeeOther)
}

func (h *PageHandler) renderContactError(
	w http.ResponseWriter,
	msg, name, email, message string,
) {
	h.Render(w, RenderOptions{
		Page:   "index.html",
		Header: "home-header.html",
		Footer: "index-footer.html",
		Data: map[string]interface{}{
			"SiteTitle": "Pastor Aby & Pastor Smitha",
			"Error":     msg,
			"FormData": map[string]string{
				"name":    name,
				"email":   email,
				"message": message,
			},
		},
	})
}