package handlers

import (
	"net/http"
	"word_press/models"
	"word_press/auth"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"SiteTitle": "Login",
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
			Render(w, "login.html", data)
			return
		}

		err = auth.LoginUser(w, r, user)
		if err != nil {
			data["Error"] = "Failed to create sesson"
			Render(w, "login.html", data)
			return
		}

		data["User"] = user
        
		Render(w, "login.html", data)
		return
	}

	Render(w, "login.html", data)
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