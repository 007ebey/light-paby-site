package services

import "word_press/models"

type postRepo struct{}

type PostRepository interface {
	GetAll() ([]models.Post, error)
	GetByID(id int) (*models.Post, error)
	GetBySlug(slug string) (*models.Post, error)
	Create(post *models.Post) error
	Update(post *models.Post) error
	Delete(id int) error
	GetRecentPosts(count int) ([]models.Post)
}

func NewPostRepository() PostRepository {
	return &postRepo{}
}

func (r *postRepo) GetAll() ([]models.Post, error) {
	return models.GetAllPosts()
}

func (r *postRepo) GetByID(id int) (*models.Post, error) {
	return models.GetPostByID(id)
}

func (r *postRepo) GetBySlug(slug string) (*models.Post, error) {
	return models.GetPostBySlug(slug)
}

func (r *postRepo) Create(post *models.Post) error {
	return models.CreatePost(
		post.Title,
		post.Slug,
		post.Content,
		post.Status,
		post.AuthorID,
		post.Excerpt,
		post.FeaturedImage,
	)
}

func (r *postRepo) Update(post *models.Post) error {
	return models.UpdatePost(
		post.ID,
		post.Title,
		post.Slug,
		post.Content,
		post.Status,
		post.Excerpt,
		post.AuthorID,
		post.FeaturedImage,
	)
}

func (r *postRepo) Delete(id int) error {
	return models.DeletePost(id)
}

func (r *postRepo) GetRecentPosts(count int) []models.Post {
	posts, err := models.GetRecentPosts(count)
	if err != nil {
		// ⚠️ swallowing error — not ideal
		return []models.Post{}
	}
	return posts
}