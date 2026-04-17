package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"word_press/models"
)

//
// ===== TESTS =====
//

func TestLogin_GET(t *testing.T) {
	mockAuth := &MockAuthService{}
	mockPost := &MockPostService{}

	h := NewAuthHandler(mockAuth, mockPost)

	h.GetUser = func(r *http.Request) (*models.User, error) {
		return nil, nil
	}

	h.Render = func(w http.ResponseWriter, opts RenderOptions) {
		w.WriteHeader(http.StatusOK)
	}

	req := httptest.NewRequest("GET", "/login", nil)
	rr := httptest.NewRecorder()

	h.Login(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

//
// LOGIN SUCCESS
//

func TestLogin_Success(t *testing.T) {
	mockAuth := &MockAuthService{
		LoginFn: func(u, p string) (*models.User, error) {
			return &models.User{ID: 1}, nil
		},
	}

	h := NewAuthHandler(mockAuth, nil)

	h.GetUser = func(r *http.Request) (*models.User, error) {
		return nil, nil
	}

	h.Render = func(w http.ResponseWriter, opts RenderOptions) {
		w.WriteHeader(http.StatusOK)
	}

	form := strings.NewReader("username=test&password=pass")
	req := httptest.NewRequest("POST", "/login", form)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	h.Login(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

//
// LOGIN INVALID
//

func TestLogin_Invalid(t *testing.T) {
	mockAuth := &MockAuthService{
		LoginFn: func(u, p string) (*models.User, error) {
			return nil, errors.New("invalid")
		},
	}

	h := NewAuthHandler(mockAuth, nil)

	h.GetUser = func(r *http.Request) (*models.User, error) {
		return nil, nil
	}

	h.Render = func(w http.ResponseWriter, opts RenderOptions) {
		w.WriteHeader(http.StatusOK)
	}

	form := strings.NewReader("username=test&password=wrong")
	req := httptest.NewRequest("POST", "/login", form)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	h.Login(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

//
// REGISTER SUCCESS
//

func TestRegister_Success(t *testing.T) {
	mockAuth := &MockAuthService{
		RegisterFn: func(u, e, p, r string) error {
			return nil
		},
		LoginFn: func(u, p string) (*models.User, error) {
			return &models.User{ID: 1}, nil
		},
	}

	h := NewAuthHandler(mockAuth, nil)

	h.Render = func(w http.ResponseWriter, opts RenderOptions) {
		w.WriteHeader(http.StatusOK)
	}

	form := strings.NewReader("username=a&email=a@test.com&password=1&confirm_password=1")
	req := httptest.NewRequest("POST", "/register", form)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	h.Register(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Fatalf("expected redirect, got %d", rr.Code)
	}
}

//
// REGISTER PASSWORD MISMATCH
//

func TestRegister_PasswordMismatch(t *testing.T) {
	h := NewAuthHandler(&MockAuthService{}, nil)

	h.Render = func(w http.ResponseWriter, opts RenderOptions) {
		w.WriteHeader(http.StatusOK)
	}

	form := strings.NewReader("username=a&email=a@test.com&password=1&confirm_password=2")
	req := httptest.NewRequest("POST", "/register", form)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	h.Register(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

//
// LOGOUT
//

func TestLogout(t *testing.T) {
	h := NewAuthHandler(&MockAuthService{}, nil)

	req := httptest.NewRequest("GET", "/logout", nil)
	rr := httptest.NewRecorder()

	h.Logout(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Fatalf("expected redirect, got %d", rr.Code)
	}
}