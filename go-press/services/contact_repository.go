package services

import (
	"word_press/models"
)

type ContactRepository interface {
	GetAll() ([]models.Contact, error)
	UpdateStatus(id int, status string) error
	Delete(id int) error
	CreateContact(name, email, message, ip string) error
}

type contactRepo struct{}

func NewContactRepository() ContactRepository {
	return &contactRepo{}
}

func (r *contactRepo) GetAll() ([]models.Contact, error) {
	return models.GetAllContacts()
}

func (r *contactRepo) UpdateStatus(id int, status string) error {
	return models.UpdateContactStatus(id, status)
}

func (r *contactRepo) Delete(id int) error {
	return models.DeleteContact(id)
}

func (r *contactRepo) CreateContact(name, email, message, ip string) error {
	return models.CreateContact(name, email, message, ip)
}