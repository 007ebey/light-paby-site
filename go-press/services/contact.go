package services

import (
	"word_press/models"
	"errors"
)

type ContactService interface {
	GetAllContacts() ([]models.Contact, error)
	MarkAsRead(id int) error
	DeleteContact(id int) error
}

type contactService struct {
	repo ContactRepository
}

func NewContactService(repo ContactRepository) ContactService {
	return &contactService{repo: repo}
}

func (s *contactService) GetAllContacts() ([]models.Contact, error) {
	return s.repo.GetAll()
}

func (s *contactService) MarkAsRead(id int) error {
	if id == 0 {
		return errors.New("invalid contact ID")
	}
	return s.repo.UpdateStatus(id, "read")
}

func (s *contactService) DeleteContact(id int) error {
	if id == 0 {
		return errors.New("invalid contact ID")
	}
	return s.repo.Delete(id)
}