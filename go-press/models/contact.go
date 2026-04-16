package models

import (
	"time"
	"word_press/database"
	"errors"
)

type Contact struct {
	ID        int
	Name      string
	Email     string
	Message   string
	Status    string
	IPAddress string
	CreatedAt time.Time
}

func CreateContact(name, email, message, ip string) error {
	query := `
	INSERT INTO contacts (name, email, message, status, ip_address, created)
	VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := database.DB.Exec(
		query,
		name,
		email,
		message,
		"new", // important: lifecycle starts here
		ip,
		time.Now(),
	)

	return err
}

func GetContactsByStatus(status string) ([]Contact, error) {
	if database.DB == nil {
		return nil, errors.New("DB not initialized")
	}

	rows, err := database.DB.Query(`
		SELECT id, name, email, message, status, ip_address, created
		FROM contacts
		WHERE status = ?
		ORDER BY created DESC
	`, status)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contacts []Contact

	for rows.Next() {
		var c Contact
		err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.Email,
			&c.Message,
			&c.Status,
			&c.IPAddress,
			&c.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		contacts = append(contacts, c)
	}

	return contacts, nil
}

func GetAllContacts() ([]Contact, error) {
	rows, err := database.DB.Query(`
		SELECT id, name, email, message, status, ip_address, created
		FROM contacts
		ORDER BY created DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contacts []Contact

	for rows.Next() {
		var c Contact
		err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.Email,
			&c.Message,
			&c.Status,
			&c.IPAddress,
			&c.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		contacts = append(contacts, c)
	}

	// 🔴 You forgot this (important)
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return contacts, nil
}

func UpdateContactStatus(id int, status string) error {
	_, err := database.DB.Exec(
		"UPDATE contacts SET status = ? WHERE id = ?",
		status,
		id,
	)
	return err
}

func DeleteContact(id int) error {
	_, err := database.DB.Exec(
		"DELETE FROM contacts WHERE id = ?",
		id,
	)
	return err
}