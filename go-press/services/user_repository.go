package services

import "word_press/models"

type userRepo struct{}

type UserRepository interface {
	GetByUsername(username string) (*models.User, error)
	Create(username, email, password, role string) error
}

func NewUserRepository() UserRepository {
	return &userRepo{}
}

func (r *userRepo) GetByUsername(username string) (*models.User, error) {
	return models.GetUserByUserName(username)
}

func (r *userRepo) Create(username, email, password, role string) error {
	return models.CreateUser(username, email, password, role)
}