package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// You likely already have router setup in main.go
func setupTestRouter() http.Handler {
	return SetupRouter() // adjust if your function name differs
}

func TestHomePage(t *testing.T) {
	r := setupTestRouter()

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestBlogPage(t *testing.T) {
	r := setupTestRouter()

	req := httptest.NewRequest("GET", "/blog", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestBlogPagination(t *testing.T) {
	r := setupTestRouter()

	req := httptest.NewRequest("GET", "/blog?page=2", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 for page=2, got %d", w.Code)
	}
}

func TestStaticFiles(t *testing.T) {
	r := setupTestRouter()

	req := httptest.NewRequest("GET", "/style/images/logo.png", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("static file not served, got %d", w.Code)
	}
}

func TestNotFound(t *testing.T) {
	r := setupTestRouter()

	req := httptest.NewRequest("GET", "/does-not-exist", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestAdminRouteUnauthorized(t *testing.T) {
	r := setupTestRouter()

	req := httptest.NewRequest("GET", "/admin/posts", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized && w.Code != http.StatusForbidden {
		t.Errorf("expected auth failure, got %d", w.Code)
	}
}

func TestTemplateRender(t *testing.T) {
	r := setupTestRouter()

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	body := w.Body.String()

	if body == "" {
		t.Fatal("expected HTML response, got empty body")
	}

	if !contains(body, "<html") {
		t.Error("response does not look like HTML")
	}
}

// helper
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && (string([]byte(s)[0:len(substr)]) == substr || contains(s[1:], substr)))
}