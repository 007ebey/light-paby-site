package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"word_press/models"
)

func newPageHandler(
	ps *MockPostService,
	cs *MockContactService,
) *PageHandler {
	return &PageHandler{
		PostService:    ps,
		ContactService: cs,
		GetUser: func(r *http.Request) (*models.User, error) {
			return &models.User{Email: "test@test.com", Role: "user"}, nil
		},
		Render: func(w http.ResponseWriter, opts RenderOptions) {},
	}
}

func TestHome_SuccessMessage(t *testing.T) {
	h := newPageHandler(nil, nil)

	var success string

	h.Render = func(w http.ResponseWriter, opts RenderOptions) {
		data := opts.Data.(map[string]interface{}) // ✅ FIX
		if v, ok := data["Success"]; ok {
			success = v.(string)
		}
	}

	req := httptest.NewRequest("GET", "/?success=1", nil)
	rr := httptest.NewRecorder()

	h.Home(rr, req)

	if success == "" {
		t.Fatal("expected success message")
	}
}

func TestBlog_Success(t *testing.T) {
	called := false

	mockPost := &MockPostService{
		GetPublishedPostsFn: func(page, limit int) ([]models.Post, int, error) {
			called = true
			if page != 1 {
				t.Fatalf("wrong page")
			}
			return []models.Post{{ID: 1}}, 2, nil
		},
	}

	h := newPageHandler(mockPost, nil)

	var data map[string]interface{}

	h.Render = func(w http.ResponseWriter, opts RenderOptions) {
		data = opts.Data.(map[string]interface{}) // ✅ FIX
	}

	req := httptest.NewRequest("GET", "/blog", nil)
	rr := httptest.NewRecorder()

	h.Blog(rr, req)

	if !called {
		t.Fatal("expected service call")
	}

	if data["TotalPages"].(int) != 2 {
		t.Fatal("wrong total pages")
	}
}

func TestBlog_Error(t *testing.T) {
	mockPost := &MockPostService{
		GetPublishedPostsFn: func(page, limit int) ([]models.Post, int, error) {
			return nil, 0, errors.New("fail")
		},
	}

	h := newPageHandler(mockPost, nil)

	req := httptest.NewRequest("GET", "/blog", nil)
	rr := httptest.NewRecorder()

	h.Blog(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 got %d", rr.Code)
	}
}

func TestContact_GET(t *testing.T) {
	h := newPageHandler(nil, nil)

	called := false

	h.Render = func(w http.ResponseWriter, opts RenderOptions) {
		called = true
	}

	req := httptest.NewRequest("GET", "/contact", nil)
	rr := httptest.NewRecorder()

	h.Contact(rr, req)

	if !called {
		t.Fatal("expected render")
	}
}

func TestContact_MethodNotAllowed(t *testing.T) {
	h := newPageHandler(nil, nil)

	req := httptest.NewRequest("PUT", "/contact", nil)
	rr := httptest.NewRecorder()

	h.Contact(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405 got %d", rr.Code)
	}
}

func TestContact_InvalidFormType(t *testing.T) {
	h := newPageHandler(nil, nil)

	req := httptest.NewRequest("POST", "/contact",
		strings.NewReader("form_type=wrong"),
	)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	h.Contact(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d", rr.Code)
	}
}

func TestContact_ValidationError(t *testing.T) {
	h := newPageHandler(nil, nil)

	rendered := false

	h.Render = func(w http.ResponseWriter, opts RenderOptions) {
		rendered = true
		data := opts.Data.(map[string]interface{}) // ✅ FIX
		if data["Error"] == "" {
			t.Fatal("expected error message")
		}
	}

	req := httptest.NewRequest("POST", "/contact",
		strings.NewReader("form_type=contact&name=&email=&message="),
	)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	h.Contact(rr, req)

	if !rendered {
		t.Fatal("expected render")
	}
}

func TestContact_InvalidEmail(t *testing.T) {
	h := newPageHandler(nil, nil)

	rendered := false

	h.Render = func(w http.ResponseWriter, opts RenderOptions) {
		rendered = true
	}

	req := httptest.NewRequest("POST", "/contact",
		strings.NewReader("form_type=contact&name=a&email=bad&message=msg"),
	)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	h.Contact(rr, req)

	if !rendered {
		t.Fatal("expected validation render")
	}
}

func TestContact_Success(t *testing.T) {
	called := false

	mockContact := &MockContactService{
		CreateFn: func(name, email, msg, ip string) error {
			called = true
			return nil
		},
	}

	h := newPageHandler(nil, mockContact)

	req := httptest.NewRequest("POST", "/contact",
		strings.NewReader("form_type=contact&name=a&email=a@test.com&message=msg"),
	)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	h.Contact(rr, req)

	if !called {
		t.Fatal("expected CreateContact call")
	}

	if rr.Code != http.StatusSeeOther {
		t.Fatalf("expected redirect got %d", rr.Code)
	}
}

func TestContact_ServiceFailure(t *testing.T) {
	mockContact := &MockContactService{
		CreateFn: func(name, email, msg, ip string) error {
			return errors.New("fail")
		},
	}

	h := newPageHandler(nil, mockContact)

	rendered := false

	h.Render = func(w http.ResponseWriter, opts RenderOptions) {
		rendered = true
	}

	req := httptest.NewRequest("POST", "/contact",
		strings.NewReader("form_type=contact&name=a&email=a@test.com&message=msg"),
	)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	h.Contact(rr, req)

	if !rendered {
		t.Fatal("expected error render")
	}
}