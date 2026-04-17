package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"

	"word_press/handlers"
	"word_press/models"
)

func newTestHandler(ps *handlers.MockPostService, cs *handlers.MockCommentService) *handlers.BlogHandler {
	return &handlers.BlogHandler{
		PostService:    ps,
		CommentService: cs,
		GetUser: func(r *http.Request) (*models.User, error) {
			return &models.User{Email: "test@test.com", Role: "user"}, nil
		},
		Render: func(w http.ResponseWriter, opts handlers.RenderOptions) {},
	}
}

func TestBlogPost_Success(t *testing.T) {
	mockPost := &handlers.MockPostService{
		GetBySlugFn: func(slug string) (*models.Post, error) {
			return &models.Post{ID: 1, Title: "Test"}, nil
		},
	}

	mockComment := &handlers.MockCommentService{
		GetByPostFn: func(postID int) ([]models.Comment, error) {
			return []models.Comment{{ID: 1}}, nil
		},
	}

	var rendered bool

	h := newTestHandler(mockPost, mockComment)
	h.Render = func(w http.ResponseWriter, opts handlers.RenderOptions) {
		rendered = true
		if opts.Page != "blog-post.html" {
			t.Errorf("wrong template")
		}
	}

	req := httptest.NewRequest("GET", "/blog/test", nil)
	req = mux.SetURLVars(req, map[string]string{"slug": "test"})
	rr := httptest.NewRecorder()

	h.BlogPost(rr, req)

	if !rendered {
		t.Fatal("expected render to be called")
	}
}

func TestBlogPost_NotFound(t *testing.T) {
	mockPost := &handlers.MockPostService{
		GetBySlugFn: func(slug string) (*models.Post, error) {
			return nil, nil
		},
	}

	h := newTestHandler(mockPost, nil)

	req := httptest.NewRequest("GET", "/blog/test", nil)
	req = mux.SetURLVars(req, map[string]string{"slug": "test"})
	rr := httptest.NewRecorder()

	h.BlogPost(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404 got %d", rr.Code)
	}
}

func TestCreateComment_Success(t *testing.T) {
	mockPost := &handlers.MockPostService{
		GetBySlugFn: func(slug string) (*models.Post, error) {
			return &models.Post{ID: 1}, nil
		},
	}

	called := false

	mockComment := &handlers.MockCommentService{
		CreateFn: func(postID int, a, e, b string) error {
			called = true
			if postID != 1 {
				t.Fatalf("wrong postID")
			}
			return nil
		},
	}

	h := newTestHandler(mockPost, mockComment)

	req := httptest.NewRequest("POST", "/blog/test",
		strings.NewReader("author=a&email=e&body=b"),
	)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req = mux.SetURLVars(req, map[string]string{"slug": "test"})

	rr := httptest.NewRecorder()

	h.CreateComment(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Fatalf("expected redirect got %d", rr.Code)
	}

	if !called {
		t.Fatal("expected CreateComment to be called")
	}

	location := rr.Header().Get("Location")
	if !strings.Contains(location, "comment=success") {
		t.Fatalf("unexpected redirect location: %s", location)
	}
}

func TestCreateComment_MethodNotAllowed(t *testing.T) {
	h := newTestHandler(nil, nil)

	req := httptest.NewRequest("GET", "/blog/test", nil)
	rr := httptest.NewRecorder()

	h.CreateComment(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405 got %d", rr.Code)
	}
}

func TestDeleteComment_Unauthorized(t *testing.T) {
	h := newTestHandler(nil, nil)
	h.GetUser = func(r *http.Request) (*models.User, error) {
		return nil, nil
	}

	req := httptest.NewRequest("POST", "/delete", nil)
	rr := httptest.NewRecorder()

	h.DeleteComment(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 got %d", rr.Code)
	}
}

func TestDeleteComment_Success(t *testing.T) {
	mockPost := &handlers.MockPostService{
		GetBySlugFn: func(slug string) (*models.Post, error) {
			return &models.Post{ID: 1}, nil
		},
	}

	called := false

	mockComment := &handlers.MockCommentService{
		GetByIDFn: func(postID, id int) (*models.Comment, error) {
			return &models.Comment{Email: "test@test.com"}, nil
		},
		DeleteFn: func(id int) error {
			called = true
			return nil
		},
	}

	h := newTestHandler(mockPost, mockComment)

	req := httptest.NewRequest("POST", "/delete",
		strings.NewReader("id=1&slug=test"),
	)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	h.DeleteComment(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Fatalf("expected redirect got %d", rr.Code)
	}

	if !called {
		t.Fatal("expected DeleteComment to be called")
	}
}