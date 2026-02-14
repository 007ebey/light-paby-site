package models

import (
  "word_press/database"
  "time"
  "database/sql"
)

type Post struct {
	ID int 
	Title string
	Slug string 
	Content string
	Status string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func GetPublishedPosts() ([]Post, error) {
	rows, _ := database.DB.Query("SELECT id,title,slug,content FROM posts WHERE status='published'")
	var posts []Post
	for rows.Next() {
		var p Post 
		rows.Scan(&p.ID, &p.Title, &p.Slug, &p.Content)
		posts = append(posts, p)
	}
	return posts, nil
}

func GetPostByID(id int) (*Post, error) {
	query :=
	`
	SELECT id, title, slug, content, status, created, updated
	FROM posts
	WHERE id = ?
	LIMIT 1
	`

	var p Post
	err := database.DB.QueryRow(query, id).Scan(
		&p.ID,
		&p.Title,
		&p.Slug,
		&p.Content,
		&p.Status,
		&p.CreatedAt,
		&p.UpdatedAt)
	if err != nil {

		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &p, nil
}

func GetPostsByStatus(status string) ([]Post, error) {
	query := `
	SELECT id, title, slug, content, status, created, updated
	FROM posts
	WHERE status = ?
	ORDER BY created DESC
	`

	rows, err := database.DB.Query(query, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post

	for rows.Next() {
		var p Post
		err := rows.Scan(
			&p.ID,
			&p.Title,
			&p.Slug,
			&p.Content,
			&p.Status,
			&p.CreatedAt,
			&p.UpdatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}

func GetPostsPaginated(limit, offset int) ([]Post, error) {
	query := `
	 SELECT id, title, slug, content, status, created, updated
	 FROM posts
	 WHERE status = 'published'
	 ORDER BY created DESC
	 LIMIT ? OFFSET ?
	`

	rows, err := database.DB.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var posts []Post
	for rows.Next() {
		var p Post
		rows.Scan(
			&p.ID,
			&p.Title,
			&p.Slug,
			&p.Content,
			&p.Status,
			&p.CreatedAt,
			&p.UpdatedAt)
		posts = append(posts, p)
	}
	return posts, nil
}

func SearchPosts(keywords string) ([]Post, error) {
	query := 
	`
	SELECT id, title, slug, content, status, created, updated
	FROM posts
	WHERE status = 'published'
	AND (
	   title LIKE ?
	   OR content LIKE ?
	)
	ORDER BY created DESC
	`

	search := "%" + keywords + "%"

	rows, err := database.DB.Query(query, search, search)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var posts []Post
	for rows.Next() {
		var p Post
		rows.Scan(
			&p.ID,
			&p.Title,
			&p.Slug,
			&p.Content,
			&p.Status,
			&p.CreatedAt,
			&p.UpdatedAt)

		posts = append(posts, p)
	}
	
	return posts, nil
}

func GetPostCount(status string) (int, error) {
	query :=
	`
	SELECT COUNT(*)
	FROM posts
	WHERE status = ?
	`
	var count int 
	err := database.DB.QueryRow(query, status).Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil
}