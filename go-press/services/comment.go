package services

import  ( 
	"word_press/models"
	"errors"
	"strings"
)

type CommentService interface {
	CreateComment(postID int, author, email, body string) error
	DeleteComment(id int) error
	GetCommentsByPost(postID int) ([]models.Comment, error)
	GetCommentByID(postID, commentID int) (*models.Comment, error)
}

type commentService struct{}

func NewCommentService() CommentService {
	return &commentService{}
}

func (s *commentService) GetCommentsByPost(postID int) ([]models.Comment, error) {
	return models.GetCommentsByPost(postID)
}

func (s *commentService) CreateComment(postID int, author, email, body string) error {
	if author == "" || email == "" || body == "" {
		return errors.New("missing fields")
	}

	if !strings.Contains(email, "@") {
		return errors.New("invalid email")
	}

	return models.CreateComment(postID, author, email, body)
}

func (s *commentService) GetCommentByID(postID, commentID int) (*models.Comment, error) {
	comments, err := models.GetCommentsByPost(postID)
	if err != nil {
		return nil, err
	}

	for _, c := range comments {
		if c.ID == commentID {
			return &c, nil
		}
	}

	return nil, errors.New("not found")
}

func (s *commentService) DeleteComment(id int) error {

	if id == 0 {
		return errors.New("invalid comment ID")
	}

	return models.DeleteComment(id)
}