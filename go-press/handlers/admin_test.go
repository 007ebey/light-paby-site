package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
    "io"
	"word_press/models"
	"word_press/auth"
	"word_press/services"
)

//
// ===== MOCKS (COMPLETE) =====
//

type MockPostService struct {
	GetAllFn     func() ([]models.Post, error)
	GetByIDFn    func(int) (*models.Post, error)
	GetBySlugFn  func(string) (*models.Post, error)
	CreateFn     func(*models.Post) error
	UpdateFn     func(*models.Post) error
	DeleteFn     func(int) error
}

func (m *MockPostService) GetAllPosts() ([]models.Post, error) {
	return m.GetAllFn()
}

func (m *MockPostService) GetPostByID(id int) (*models.Post, error) {
	return m.GetByIDFn(id)
}

func (m *MockPostService) GetPostBySlug(slug string) (*models.Post, error) {
	if m.GetBySlugFn != nil {
		return m.GetBySlugFn(slug)
	}
	return nil, nil
}

func (m *MockPostService) CreatePost(p *models.Post) error {
	return m.CreateFn(p)
}

func (m *MockPostService) UpdatePost(p *models.Post) error {
	if m.UpdateFn != nil {
		return m.UpdateFn(p)
	}
	return nil
}

func (m *MockPostService) DeletePost(id int) error {
	return m.DeleteFn(id)
}

//
// IMAGE SERVICE
//

type MockImageService struct {
	SaveFn   func(io.Reader, string, int64) (string, error)
	DeleteFn func(string) error
}

func (m *MockImageService) SaveImage(file io.Reader, filename string, size int64) (string, error) {
	if m.SaveFn != nil {
		return m.SaveFn(file, filename, size)
	}
	return "", nil
}

func (m *MockImageService) DeleteImage(path string) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(path)
	}
	return nil
}

//
// ===== COMPILE-TIME SAFETY (CRITICAL) =====
//

var _ services.PostService = (*MockPostService)(nil)
var _ services.ImageService = (*MockImageService)(nil)

//
// ===== TESTS =====
//

func TestAdminPosts_Success(t *testing.T) {

	// override global (NOT :=)
	getCurrentUser = func(r *http.Request) (*models.User, error) {
		return &models.User{ID: 1}, nil
	}
	defer func() { getCurrentUser = auth.GetCurrentUser }()

	render = func(w http.ResponseWriter, opts RenderOptions) {
	  w.WriteHeader(http.StatusOK)
    }
    defer func() { render = RenderWithOpts }() 

	mockPost := &MockPostService{
		GetAllFn: func() ([]models.Post, error) {
			return []models.Post{{ID: 1}}, nil
		},
	}

	h := NewAdminHandler(mockPost, nil, nil)

	req := httptest.NewRequest("GET", "/admin/posts", nil)
	rr := httptest.NewRecorder()

	h.AdminPosts(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

func TestAdminPosts_Error(t *testing.T) {

	getCurrentUser = func(r *http.Request) (*models.User, error) {
		return &models.User{ID: 1}, nil
	}

	render = func(w http.ResponseWriter, opts RenderOptions) {
	  w.WriteHeader(http.StatusOK)
    }
    defer func() { render = RenderWithOpts }() 

	mockPost := &MockPostService{
		GetAllFn: func() ([]models.Post, error) {
			return nil, errors.New("db error")
		},
	}

	h := NewAdminHandler(mockPost, nil, nil)

	req := httptest.NewRequest("GET", "/admin/posts", nil)
	rr := httptest.NewRecorder()

	h.AdminPosts(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rr.Code)
	}
}

func TestAdminCreatePost_Success(t *testing.T) {

	getCurrentUser = func(r *http.Request) (*models.User, error) {
		return &models.User{ID: 1}, nil
	}

	render = func(w http.ResponseWriter, opts RenderOptions) {
	  w.WriteHeader(http.StatusOK)
    }
    defer func() { render = RenderWithOpts }() 

	mockPost := &MockPostService{
		CreateFn: func(p *models.Post) error {
			if p.Title != "Test" {
				t.Fatalf("unexpected title")
			}
			return nil
		},
	}

	h := NewAdminHandler(mockPost, &MockImageService{}, nil)

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

	getCurrentUser = func(r *http.Request) (*models.User, error) {
		return &models.User{ID: 1}, nil
	}

	render = func(w http.ResponseWriter, opts RenderOptions) {
	  w.WriteHeader(http.StatusOK)
    }
    defer func() { render = RenderWithOpts }() 

	h := NewAdminHandler(nil, nil, nil)

	form := strings.NewReader("title=&content=")
	req := httptest.NewRequest("POST", "/admin/posts/create", form)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	h.AdminCreatePost(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}