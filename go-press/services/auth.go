package services

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"word_press/models"
)

type AuthService interface {
	Login(username, password string) (*models.User, error)
	Register(username, email, password, role string) error
}

type authService struct{
	repo UserRepository
}

func NewAuthService(repo UserRepository) AuthService {
	return &authService{
		repo: repo,
	}
}

func (s *authService) Register(username, email, password, role string) error {

	existing, err := s.repo.GetByUsername(username)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New("user already exists")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = s.repo.Create(username, email, string(hashed), role)
	if err != nil {
		return errors.New("failed to create user")
	}

	return nil
}

func (s *authService) Login(username, password string) (*models.User, error) {

	user, err := s.repo.GetByUsername(username)
	if err != nil || user == nil {
		return nil, errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}