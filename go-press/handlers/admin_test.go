package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"word_press/models"
)

func TestAdminPosts_Success(t *testing.T) {

	mockPost := &MockPostService{
		GetAllFn: func() ([]models.Post, error) {
			return []models.Post{{ID: 1}}, nil
		},
	}

	h := NewAdminHandler(mockPost, nil, nil)

	// ✅ inject instead of global override
	h.GetUser = func(r *http.Request) (*models.User, error) {
		return &models.User{ID: 1}, nil
	}
	h.Render = func(w http.ResponseWriter, opts RenderOptions) {
		w.WriteHeader(http.StatusOK)
	}

	req := httptest.NewRequest("GET", "/admin/posts", nil)
	rr := httptest.NewRecorder()

	h.AdminPosts(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

func TestAdminPosts_Error(t *testing.T) {

	mockPost := &MockPostService{
		GetAllFn: func() ([]models.Post, error) {
			return nil, errors.New("db error")
		},
	}

	h := NewAdminHandler(mockPost, nil, nil)

	// ✅ inject instead of global override
	h.GetUser = func(r *http.Request) (*models.User, error) {
		return &models.User{ID: 1}, nil
	}
	h.Render = func(w http.ResponseWriter, opts RenderOptions) {
		w.WriteHeader(http.StatusOK)
	}


	req := httptest.NewRequest("GET", "/admin/posts", nil)
	rr := httptest.NewRecorder()

	h.AdminPosts(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rr.Code)
	}
}

func TestAdminCreatePost_Success(t *testing.T) {

	mockPost := &MockPostService{
		CreateFn: func(p *models.Post) error {
			if p.Title != "Test" {
				t.Fatalf("unexpected title")
			}
			return nil
		},
	}

	h := NewAdminHandler(mockPost, &MockImageService{}, nil)

	// ✅ inject instead of global override
	h.GetUser = func(r *http.Request) (*models.User, error) {
		return &models.User{ID: 1}, nil
	}
	h.Render = func(w http.ResponseWriter, opts RenderOptions) {
		w.WriteHeader(http.StatusOK)
	}


	form := strings.NewReader("title=Test&content=Body")
	req := httptest.NewRequest("POST", "/admin/posts/create", form)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	h.AdminCreatePost(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Fatalf("expected redirect, got %d", rr.Code)
	}
}

func TestAdminCreatePost_ValidationFail(t *testing.T) {

	h := NewAdminHandler(nil, nil, nil)

	// ✅ inject instead of global override
	h.GetUser = func(r *http.Request) (*models.User, error) {
		return &models.User{ID: 1}, nil
	}
	h.Render = func(w http.ResponseWriter, opts RenderOptions) {
		w.WriteHeader(http.StatusOK)
	}


	form := strings.NewReader("title=&content=")
	req := httptest.NewRequest("POST", "/admin/posts/create", form)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	h.AdminCreatePost(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}