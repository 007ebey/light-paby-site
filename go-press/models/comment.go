package models

import (
	"database/sql"
	"time"
	"word_press/database"
)

type Comment struct {
	ID int
	PostID int
	Author string
	Email string
	Body string
	Status string
	CreatedAt time.Time
}

func CreateComment(postID int, author, email, body string) error {
	query := `
	 INSERT INTO comments (post_id, author, email, body, status, created)
	 VALUES (?,?,?,?,?,?)
	`
	_, err := database.DB.Exec(
		query,
		postID,
		author,
		email,
		body,
		"approved",
		time.Now()
	)

	return err
}

func GetCommentsByPost(postID int) ([]Comment, error) {
	query := `
	SELECT id, post_id, author, email, body, status, created
	FROM comments
	WHERE post_id = ? AND status = 'approved'
	ORDER BY created DESC
	`

	rows, err := database.DB.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var c Comment
		err := rows.Scan(
			&c.ID,
			&c.PostID,
			&c.Author,
			&c.Email,
			&c.Body,
			&c.Status,
			&c.CreatedAt
		)

		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
}

func GetAllComments() ([]Comment, error) {
	rows, err := database.DB.Query(
		`
		SELECT id, post_id, author, email, body, status, created
		FROM comments
		ORDER BY created DESC
		`
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var c Comment
		rows.Scan(
			&c.ID,
			&c.PostID,
			&c.Author,
			&c.Email,
			&c.Body,
			&c.Status,
			&c.CreatedAt
		)
		comments = append(comments, c)
	}

	return comments, nil
}

func DeleteComment(id int) error {
	_, err := database.DB.Exec(
		"DELETE FROM comments WHERE id= ?",
		id
	)
	return err
}

func ApproveComment(id int) error {
	_, err := database.DB.Exec(
		"UPDATE comments SET status = 'approved' WHERE id = ?",
		id
	)
	return err
}