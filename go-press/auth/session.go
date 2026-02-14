package auth

import (
	"net/http"
	"github.com/gorilla/sessions"
	"word_press/models"
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
	session, err := Store.Get(r, SessionName)
	if err != nil {
		return nil, err
	}
	userID, ok := session.Values["user_id"].(int)
	if !ok {
		return nil, nil
	}
	return models.GetUserByID(userID)
}

func IsLoggedIn(r *http.Request) bool {
	session, err := Store.Get(r, SessionName)
	if err != nil {
		return false
	}
	_, ok := session.Values[
		"user_id"]
	return ok
}

func RequireLogin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !IsLoggedIn(r) {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}