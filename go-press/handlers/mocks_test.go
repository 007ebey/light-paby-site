package handlers

import (
	"io"
	"word_press/models"
	"word_press/services"
)

type MockPostService struct {
	GetAllFn    func() ([]models.Post, error)
	GetByIDFn   func(int) (*models.Post, error)
	GetBySlugFn func(string) (*models.Post, error)
	CreateFn    func(*models.Post) error
	UpdateFn    func(*models.Post) error
	DeleteFn    func(int) error
	GetRecentFn func(int) ([]models.Post, error)
	GetPublishedPostsFn func(page, limit int) ([]models.Post, int, error)
}

func (m *MockPostService) GetPublishedPosts(page, limit int) ([]models.Post, int, error) {
	if m.GetPublishedPostsFn != nil {
		return m.GetPublishedPostsFn(page, limit)
	}
	return nil, 0, nil
}

func (m *MockPostService) GetAllPosts() ([]models.Post, error) {
	return m.GetAllFn()
}

func (m *MockPostService) GetPostByID(id int) (*models.Post, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(id)
	}
	return nil, nil
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

func (m *MockPostService) GetRecentPosts(limit int) ([]models.Post, error) {
	if m.GetRecentFn != nil {
		return m.GetRecentFn(limit)
	}
	return nil, nil
}

//
// IMAGE MOCK
//

type MockImageService struct {
	SaveFn   func(io.Reader, string, int64) (string, error)
	DeleteFn func(string) error
}

func (m *MockImageService) SaveImage(file io.Reader, name string, size int64) (string, error) {
	if m.SaveFn != nil {
		return m.SaveFn(file, name, size)
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
// ===== MOCKS =====
//

type MockAuthService struct {
	LoginFn    func(string, string) (*models.User, error)
	RegisterFn func(string, string, string, string) error
}

func (m *MockAuthService) Login(u, p string) (*models.User, error) {
	return m.LoginFn(u, p)
}

func (m *MockAuthService) Register(u, e, p, r string) error {
	return m.RegisterFn(u, e, p, r)
}

type MockCommentService struct {
	GetByPostFn   func(int) ([]models.Comment, error)
	CreateFn      func(int, string, string, string) error
	GetByIDFn     func(int, int) (*models.Comment, error)
	DeleteFn      func(int) error
}

func (m *MockCommentService) GetCommentsByPost(postID int) ([]models.Comment, error) {
	if m.GetByPostFn != nil {
		return m.GetByPostFn(postID)
	}
	return nil, nil
}

func (m *MockCommentService) CreateComment(postID int, author, email, body string) error {
	return m.CreateFn(postID, author, email, body)
}

func (m *MockCommentService) GetCommentByID(postID, id int) (*models.Comment, error) {
	return m.GetByIDFn(postID, id)
}

func (m *MockCommentService) DeleteComment(id int) error {
	return m.DeleteFn(id)
}

type MockContactService struct {
	CreateFn       func(string, string, string, string) error
	GetAllFn       func() ([]models.Contact, error)
	MarkAsReadFn   func(int) error
	DeleteFn       func(int) error
}

func (m *MockContactService) CreateContact(name, email, message, ip string) error {
	if m.CreateFn != nil {
		return m.CreateFn(name, email, message, ip)
	}
	return nil
}

func (m *MockContactService) GetAllContacts() ([]models.Contact, error) {
	if m.GetAllFn != nil {
		return m.GetAllFn()
	}
	return nil, nil
}

func (m *MockContactService) MarkAsRead(id int) error {
	if m.MarkAsReadFn != nil {
		return m.MarkAsReadFn(id)
	}
	return nil
}

func (m *MockContactService) DeleteContact(id int) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(id)
	}
	return nil
}

var _ services.PostService = (*MockPostService)(nil)
var _ services.ImageService = (*MockImageService)(nil)
var _ services.CommentService = (*MockCommentService)(nil)
var _ services.ContactService = (*MockContactService)(nil)