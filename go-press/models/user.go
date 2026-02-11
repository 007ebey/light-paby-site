package models

import (
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"word_press/database"
)

type User struct {
	ID int
	Username string
	Email    string
	Password string
	Role     string
	CreatedAt time.Time
}

func CreateUser(username, email, password, role string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword(
	[]byte(password),
	 bcrypt.DefaultCost
	)

	if err != nil {
		return err
	}
	query := `
	INSERT INTO users (username, email, password, role, created)
	VALUES (?, ?, ?, ?, ?)
	`

	_, err = database.DB.Exec(
		query,
		username,
		email,
		string(hashedPassword),
		role,
		time.Now()
	)

	return err
}

func AuthenticateUser(username, password string) (*User, error) {
	query := 
	`
	SELECT id, username, email, password, role, created
	FROM users
	WHERE username = ?
	LIMIT 1
	`
	var u User

	err := database.DB.QueryRow(query, username)
	            .Scan(
					&u.ID,
					&u.Username,
					&u.Email,
					&u.Password,
					&u.Role,
					&u.CreatedAt
				)
   
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("User not found")
		}

		return nil, err
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(u.Password),
		[]byte(password)
	)

	if err != nil {
		return nil, errors.New("invalid password")
	}

	return &u, nil
}

func GetUserByID(id int) (*User, error) {
	query := `
	  SELECT id, username, email, password, role, created
	  FROM users
	  WHERE id = ?
	`

	var u User

	err := database.DB.QueryRow(query, id).Scan(
		&u.ID,
		&u.Username,
		&u.Email,
		&u.Password,
		&u.Role,
		&u.CreatedAt
	)

	if err != nil {

		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &u, nil
}

func GetUserByUserName(username string) (*User, error) {
	query := `
	  SELECT id, username, email, password, role, created
	  FROM users
	  WHERE username = ?
	`
	var u User

	err := database.DB.QueryRow(query, username).Scan(
		&u.ID,
		&u.Username,
		&u.Email,
		&u.Password,
		&u.Role,
		&u.CreatedAt
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func GetAllUsers() ([]User, error) {
	rows, err := database.DB.Query(`
	  SELECT id, username, email, password, role, created
	  FROM users
	  ORDER BY created DESC
	`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []User

	for rows.Next() {
		var u User
		rows.Scan(
			&u.ID,
			&u.Username,
			&u.Email,
			&u.Password,
			&u.Role,
			&u.CreatedAt
		)
		users = append(users, u)
	}
	return users, nil
}

func UpdateUserRole(userID int, role string) error {
  _, err := database.DB.Exec(
	"UPDATE users SET role = ? WHERE id = ?",
	role,
	userID
  )
  return err
}

func DeleteUser(userID int) error {
	_, err := database.DB.Exec(
		"DELETE FROM users WHERE id = ?",
		userID
	)
	return err
}

func UpdatePassword(userID int, newPassword string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(newPassword),
		bcrypt.DefaultCost
	)

	if err != nil {
		return err
	}

	_, err := database.DB.Exec(
		"UPDATE users SET password = ? WHERE id = ?",
		string(hashedPassword),
		userID
	)

	return err
}

