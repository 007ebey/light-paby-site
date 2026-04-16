package handlers

import (
	"net/http"
	"word_press/models"
	"word_press/auth"
	"golang.org/x/crypto/bcrypt"
	"log"
	"word_press/services"
)

type AuthHandler struct {
	Service services.AuthService
}

func NewAuthHandler(s services.AuthService) *AuthHandler {
	return &AuthHandler{
		Service: s,
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	currentUser, _ := auth.GetCurrentUser(r)

	data := map[string]interface{}{
		"SiteTitle": "Login",
	}

	// keep this for now (we’ll refactor later)
	recentPosts, err := models.GetRecentPosts(3)
	if err == nil {
		data["RecentPosts"] = recentPosts
	}

	// 🔴 Already logged in
	if currentUser != nil {
		data["User"] = currentUser
		renderLogin(w, data)
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
			renderLogin(w, data)
			return
		}

		// ✅ Session still belongs here
		if err := auth.LoginUser(w, r, user); err != nil {
			data["Error"] = "Failed to create session"
			renderLogin(w, data)
			return
		}

		data["User"] = user
		renderLogin(w, data)
		return
	}

	// 🔴 Default GET
	renderLogin(w, data)
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
			Render(w, "register.html", data)
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
			Render(w, "register.html", data)
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

	Render(w, "register.html", data)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if err := auth.LogoutUser(w, r); err != nil {
		http.Error(w, "Unable to logout", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	currentUser, _ := auth.GetCurrentUser(r)

	data := map[string]interface{}{
		"SiteTitle": "Login",
	}

	recentPosts, err2 := models.GetRecentPosts(3)

	log.Println(err2)
    
	// 🔴 Critical: handle already logged-in users FIRST
	if currentUser != nil {
		data["User"] = currentUser

		if err2 == nil {
			log.Println(err2)
			log.Println(recentPosts)
			data["RecentPosts"] = recentPosts
		}

		renderLogin(w, data)
		return
	}

	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		user, err := models.AuthenticateUser(username, password)
		if err != nil {
			data["Error"] = "Invalid username or password"
			data["Form"] = map[string]string{
				"username": username,
			}
			renderLogin(w, data)
			return
		}

		if err := auth.LoginUser(w, r, user); err != nil {
			data["Error"] = "Failed to create session"
			renderLogin(w, data)
			return
		}

		if err2 == nil {
			log.Println(err2)
			log.Println(recentPosts)
			data["RecentPosts"] = recentPosts
		}

		data["User"] = user
		renderLogin(w, data)
		return
	}

	renderLogin(w, data)
}

// 🔴 Extract rendering (you were repeating yourself everywhere)
func renderLogin(w http.ResponseWriter, data map[string]interface{}) {
	RenderWithOpts(w, RenderOptions{
		Page:   "login.html",
		Header: "login-header.html",
		Footer: "login-footer.html",
		Data:   data,
	})
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {

	data := map[string]interface{}{
		"SiteTitle": "Register",
	}

	if r.Method == http.MethodPost {

		username := r.FormValue("username")
		email    := r.FormValue("email")
		password := r.FormValue("password")
		confirm  := r.FormValue("confirm_password")

		//Validation
		if password != confirm {
			data["Error"] = "Passwords do not match"

			data["Form"] = map[string]string{
				"username": username,
				"email":    email,
			}
			Render(w, "register.html", data)
			return
		}

		hash, err := bcrypt.GenerateFromPassword(
			[]byte(password),
			bcrypt.DefaultCost,
		)

		if err != nil {
			data["Error"] = "Failed to create user"
			Render(w, "register.html", data)
			return
		}

		err = models.CreateUser(username, email, string(hash), "user")

		if err != nil {
			data["Error"] = "Username or email already exists"
			data["Form"] = map[string]string{
				"username": username,
				"email": email,
			}

			Render(w, "register.html", data)
			return
		}

		user, _ := models.GetUserByUserName(username)
		auth.LoginUser(w, r, user)

        http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	Render(w, "register.html", data)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	err := auth.LogoutUser(w, r)
	if err != nil {
		http.Error(w, "Unable to Logout", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}