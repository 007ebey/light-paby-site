package services

import (
	"word_press/models"
	"errors"
)

type postService struct {
	repo PostRepository
}

type PostService interface {
	GetAllPosts() ([]models.Post, error)
	GetPostByID(id int) (*models.Post, error)
	GetPostBySlug(slug string) (*models.Post, error)
	CreatePost(post *models.Post) error
	UpdatePost(post *models.Post) error
	DeletePost(id int) error
}

func NewPostService(repo PostRepository) PostService {
	return &postService{repo: repo}
}

func (s *postService) GetAllPosts() ([]models.Post, error) {
	return s.repo.GetAll()
}

func (s *postService) GetPostByID(id int) (*models.Post, error) {
	return s.repo.GetByID(id)
}

func (s *postService) GetPostBySlug(slug string) (*models.Post, error) {
	return s.repo.GetBySlug(slug)
}

func (s *postService) CreatePost(post *models.Post) error {

	// basic guard (don’t trust handlers blindly)
	if post.Title == "" || post.Content == "" {
		return errors.New("title and content are required")
	}

	return s.repo.Create(post)
}

func (s *postService) UpdatePost(post *models.Post) error {
	if post.ID == 0 {
		return errors.New("invalid post ID")
	}

	return s.repo.Update(post)
}

func (s *postService) DeletePost(id int) error {
	if id == 0 {
		return errors.New("invalid post ID")
	}

	return s.repo.Delete(id)
}