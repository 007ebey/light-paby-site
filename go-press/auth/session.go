package auth

import (
	"net/http"
	"github.com/gorilla/sessions"
	"word_press/models"
	"log"
	"fmt"
)

var Store = sessions.NewCookieStore([]byte(`
2L_ap(xE[?*k&DA*k.?h<{1.Azg-OgR@lo9|p&M$cD+~gIdy+}.jNkLK-V$BCd7*
`))

func init() {
	Store.Options = &sessions.Options{
		Path: "/",
		MaxAge: 86400 * 7,
		HttpOnly: true,
		Secure: false}
}

const SessionName = "word_press_session"

func LoginUser(w http.ResponseWriter, r *http.Request, user *models.User) error {

	session, err := Store.Get(r, SessionName)

	if err != nil {
		return err
	}

	session.Values["user_id"] = user.ID
	session.Values["username"] = user.Username
	session.Values["role"] = user.Role

	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,
	}

	return session.Save(r, w)
}

func LogoutUser(w http.ResponseWriter, r *http.Request) error {
	session, err := Store.Get(r, SessionName)
	if err != nil {
		return err
	}
	session.Options.MaxAge = -1
	return session.Save(r, w)
}

func GetCurrentUser(r *http.Request) (*models.User, error) {
	if Store == nil {
		return nil, fmt.Errorf("store is nil")
	}

	session, err := Store.Get(r, SessionName)
	if err != nil {
		return nil, err
	}

	if session == nil {
		return nil, fmt.Errorf("session is nil")
	}

	id, ok := session.Values["user_id"]
	if !ok {
		return nil, fmt.Errorf("user_id missing")
	}

	var userID int
	switch v := id.(type) {
	case int:
		userID = v
	case int64:
		userID = int(v)
	case float64:
		userID = int(v)
	default:
		return nil, fmt.Errorf("invalid user_id type: %T", v)
	}

	if userID <= 0 {
		return nil, fmt.Errorf("invalid user_id")
	}

	user, err := models.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}



func IsLoggedIn(r *http.Request) bool {
	session, err := Store.Get(r, SessionName)
	if err != nil {
		log.Println("failed to get session:", err)
		return false
	}

	userVal, ok := session.Values["user_id"]
	if !ok {
		return false
	}

	var userID int

	switch v := userVal.(type) {
	case int:
		userID = v
	case int64:
		userID = int(v)
	case float64: // happens with JSON/session decoding
		userID = int(v)
	default:
		log.Println("unexpected user_id type:", v)
		return false
	}

	return userID > 0
}

func RequireLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if !IsLoggedIn(r) {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func RequireAdmin(next http.Handler) http.Handler {
	log.Println("debug")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		user, err := GetCurrentUser(r)
		if err != nil || user == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		if user.Role != "admin" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}