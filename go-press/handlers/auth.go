package handlers

import (
	"net/http"
	"word_press/auth"
	"word_press/services"
	"word_press/models"
)

type AuthHandler struct {
	Service services.AuthService
	PostService services.PostService
	GetUser     func(*http.Request) (*models.User, error)
	Render      func(http.ResponseWriter, RenderOptions)
}

func NewAuthHandler(s services.AuthService,ps services.PostService) *AuthHandler {
	return &AuthHandler{
		Service: s,
		PostService: ps,
		GetUser: auth.GetCurrentUser,
		Render: RenderWithOpts,
	}
}

func (h *AuthHandler) render(w http.ResponseWriter, page string, data map[string]interface{}) {
	h.Render(w, RenderOptions{
		Page:   page,
		Header: "login-header.html",
		Footer: "login-footer.html",
		Data:   data,
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	currentUser, _ := h.GetUser(r)

	data := map[string]interface{}{
		"SiteTitle": "Login",
	}

	if h.PostService != nil {
	  recentPosts, err := h.PostService.GetRecentPosts(3)
	  if err == nil {
		data["RecentPosts"] = recentPosts
	  }
    }

	// 🔴 Already logged in
	if currentUser != nil {
		data["User"] = currentUser
		h.render(w, "login.html", data)
		return
	}

	// 🔴 Handle POST
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		// ✅ Use AuthService (no models here)
		user, err := h.Service.Login(username, password)
		if err != nil {
			data["Error"] = "Invalid username or password"
			data["Form"] = map[string]string{
				"username": username,
			}
			h.render(w, "login.html", data)
			return
		}

		// ✅ Session still belongs here
		if err := auth.LoginUser(w, r, user); err != nil {
			data["Error"] = "Failed to create session"
			h.render(w, "login.html", data)
			return
		}

		data["User"] = user
		h.render(w, "login.html", data)
		return
	}

	// 🔴 Default GET
	h.render(w, "login.html", data)
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {

	data := map[string]interface{}{
		"SiteTitle": "Register",
	}

	if r.Method == http.MethodPost {

		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")
		confirm := r.FormValue("confirm_password")

		// ✅ Validation stays in handler
		if password != confirm {
			data["Error"] = "Passwords do not match"
			data["Form"] = map[string]string{
				"username": username,
				"email":    email,
			}
			h.render(w, "register.html", data)
			return
		}

		// ✅ Use AuthService (no bcrypt, no models here)
		err := h.Service.Register(username, email, password, "user")
		if err != nil {
			data["Error"] = err.Error()
			data["Form"] = map[string]string{
				"username": username,
				"email":    email,
			}
			h.render(w, "register.html", data)
			return
		}

		// ⚠️ You STILL need user for login → this is a leak
		user, _ := h.Service.Login(username, password)

		// ✅ Session stays in handler
		if user != nil {
			auth.LoginUser(w, r, user)
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	h.render(w, "register.html", data)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if err := auth.LogoutUser(w, r); err != nil {
		http.Error(w, "Unable to logout", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

